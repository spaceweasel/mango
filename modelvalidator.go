package mango

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidateFunc is the signature for implementing custom model validators.
type ValidateFunc func(m interface{}) (map[string][]ValidationFailure, bool)

func newModelValidator(handler ValidationHandler) contextModelValidator {
	return contextModelValidator{
		validationHandler: handler,
		customValidators:  make(map[reflect.Type]ValidateFunc),
	}
}

// ModelValidator is the interface describing a model validator.
type ModelValidator interface {
	AddCustomValidator(m interface{}, fn ValidateFunc)
	Validate(m interface{}) (map[string][]ValidationFailure, bool)
}

type contextModelValidator struct {
	validationHandler ValidationHandler
	customValidators  map[reflect.Type]ValidateFunc
}

func (mv contextModelValidator) AddCustomValidator(m interface{}, fn ValidateFunc) {
	mv.customValidators[reflect.TypeOf(m)] = fn
}

func (mv contextModelValidator) Validate(m interface{}) (map[string][]ValidationFailure, bool) {
	results := make(map[string][]ValidationFailure)
	var rv reflect.Value
	if reflect.TypeOf(m).Kind() == reflect.Ptr {
		rv = reflect.ValueOf(m).Elem()
	} else {
		rv = reflect.ValueOf(m)
	}

	t := rv.Type()
	cv, exists := mv.customValidators[t]
	if exists {
		return cv(rv.Interface())
	}

	for i := 0; i < rv.NumField(); i++ {
		fieldType := rv.Type().Field(i)
		if len(fieldType.PkgPath) > 0 {
			continue
		}
		value := rv.Field(i)
		constraints := fieldType.Tag.Get("validate")

		details, ok := mv.validateProperty(fieldType.Name, value, constraints)
		if ok {
			continue
		}
		for k, v := range details {
			results[k] = v
		}
	}
	return results, len(results) == 0
}

// TODO: update docs - Seeing this error?
//
//  panic: alpha validator can only validate strings not, *string
//
//  try using the pointer validator *alpha
//
//
func (mv contextModelValidator) validateProperty(name string, rv reflect.Value, constraints string) (map[string][]ValidationFailure, bool) {
	results := make(map[string][]ValidationFailure)
	tests := mv.validationHandler.ParseConstraints(constraints)
	ignore := false
	for constraint := range tests {
		if constraint == "ignorecontents" {
			ignore = true
			break
		}
	}
	switch rv.Kind() {
	case reflect.Struct:
		if ignore {
			break
		}
		details, ok := mv.Validate(rv.Interface())
		if ok {
			break
		}
		for k, v := range details {
			results[name+"."+k] = v
		}
	case reflect.Ptr:
		// Split constraints into ones targetted at the pointer itself
		// and ones for the dereferenced object. Constraints for the
		// dereferenced value will be prefixed with '*', so that needs
		// to be stripped off before being sent to the validator.
		derefTests := ""
		ptrTests := ""
		for constraint, args := range tests {
			params := ""
			if len(args) > 0 {
				params += "(" + strings.Join(args, ",") + ")"
			}
			if len(constraint) > 0 && constraint[0] == '*' {
				derefTests = constraint[1:] + params + ","
			} else {
				ptrTests = constraint + params + ","
			}
		}
		ptrTests = strings.Trim(ptrTests, ",")
		derefTests = strings.Trim(derefTests, ",")
		// Validate the pointer itself...
		if len(ptrTests) > 0 {
			fails, ok := mv.validationHandler.IsValid(rv.Interface(), ptrTests)
			if !ok {
				results[name] = fails
			}
		}

		// Now validate the dereferenced value...
		if rv.IsNil() || ignore {
			break
		}
		val := rv.Elem()
		details, ok := mv.validateProperty(name, val, derefTests)
		if ok {
			break
		}
		for k, v := range details {
			results[k] = v
		}

	case reflect.Slice, reflect.Array:
		if constraints != "" {
			// validate the array/slice "as a container" first
			fails, ok := mv.validationHandler.IsValid(rv.Interface(), constraints)
			if !ok {
				results[name] = fails
			}
		}
		// now validate each element, but only if they're structs and
		// ignorecontents isn't set
		if ignore || rv.Type().Elem().Kind() != reflect.Struct {
			break
		}
		for i := 0; i < rv.Len(); i++ {
			details, ok := mv.Validate(rv.Index(i).Interface())
			if ok {
				continue
			}
			for k, v := range details {
				results[fmt.Sprintf("%s[%d].%s", name, i, k)] = v
			}
		}
	case reflect.Map:
		if constraints != "" {
			// validate the map "as a container" first
			fails, ok := mv.validationHandler.IsValid(rv.Interface(), constraints)
			if !ok {
				results[name] = fails
			}
		}
		// now validate each element (value, not the key), but only
		// if they're structs and ignorecontents isn't set
		if ignore || rv.Type().Elem().Kind() != reflect.Struct {
			break
		}
		for _, key := range rv.MapKeys() {
			details, ok := mv.Validate(rv.MapIndex(key).Interface())
			if ok {
				continue
			}
			for k, v := range details {
				results[fmt.Sprintf("%s[%v].%s", name, key, k)] = v
			}
		}

	case reflect.String,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		if constraints == "" {
			break
		}
		fails, ok := mv.validationHandler.IsValid(rv.Interface(), constraints)
		if !ok {
			results[name] = fails
		}
		//case reflect.Bool:
	}
	return results, len(results) == 0
}

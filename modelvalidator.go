package mango

import "reflect"

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
		value := rv.Field(i)
		fieldType := rv.Type().Field(i)
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

func (mv contextModelValidator) validateProperty(name string, rv reflect.Value, constraints string) (map[string][]ValidationFailure, bool) {
	results := make(map[string][]ValidationFailure)
	switch rv.Kind() {
	case reflect.Struct:
		// TODO: check if 'nested' constraint
		details, ok := mv.Validate(rv.Interface())
		if ok {
			break
		}
		for k, v := range details {
			results[name+"."+k] = v
		}
	case reflect.Ptr:
		val := rv.Elem()
		return mv.validateProperty(name, val, constraints)
	case reflect.String,
		reflect.Map,
		reflect.Array,
		reflect.Slice,
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

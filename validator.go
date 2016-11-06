package mango

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationFailure holds details about a validation failure. Code will
// contain the validator type and Message a user friendly description of
// the reason for the failure.
type ValidationFailure struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Validator is the interface that wraps the basic Validate method.
// Validators are used to validate models in addition to sections of a URL path
// which match the pattern of an entry in the routing tree.
type Validator interface {
	// Validate tests if val matches the validation rules. The validation test
	// may involve constraint specific args.
	Validate(val interface{}, args []string) bool
	// Type returns the constraint name used in routing patterns.
	// ValidationHandler will use this to locate the correct validator.
	Type() string
	// FailureMsg returns a string with a readable message about the validation failure.
	FailureMsg() string
}

// ValidationHandler is the interface for a collection of Validators.
// The IsValid method can be used to validate a model property or checking a
// route parameter value against its constraint.
// New validators can be added individually using AddValidator or as a
// collection using AddValidators.
type ValidationHandler interface {
	AddValidator(v Validator)
	AddValidators(validators []Validator)
	IsValid(val interface{}, constraints string) ([]ValidationFailure, bool)
	ParseConstraints(constraints string) map[string][]string
}

type elementValidationHandler struct {
	validators map[string]Validator
}

func (r *elementValidationHandler) AddValidator(v Validator) {
	if _, ok := r.validators[v.Type()]; ok {
		panic("conflicting constraint type: " + v.Type())
	}
	r.validators[v.Type()] = v
}

func (r *elementValidationHandler) AddValidators(validators []Validator) {
	for _, v := range validators {
		r.AddValidator(v)
	}
}

func (r *elementValidationHandler) IsValid(val interface{}, constraints string) (fails []ValidationFailure, ok bool) {
	// Split constraints at commas, but need to consider some may
	// have parameters which also have commas, e.g. range(3,8).
	tests := r.ParseConstraints(constraints)
	for name, args := range tests {
		// ignorecontents is a special case instruction rather than constraint
		if name == "ignorecontents" {
			continue
		}
		var v Validator
		v, ok = r.validators[name]
		if !ok {
			panic(fmt.Sprintf("unknown constraint: %s", name))
		}
		if !v.Validate(val, args) {
			if len(args) > 0 {
				name += "(" + strings.Join(args, ",") + ")"
			}
			fails = append(fails, ValidationFailure{name, v.FailureMsg()})
		}
	}
	ok = len(fails) == 0
	return
}

func (r *elementValidationHandler) ParseConstraints(constraints string) map[string][]string {
	results := make(map[string][]string)
	brace := 0
	args := []string{}
	buf := make([]byte, len(constraints))
	b := 0
	name := ""
	for i := 0; i < len(constraints); i++ {
		if constraints[i] == '(' {
			brace++
		}
		if constraints[i] == ')' {
			brace--
			if brace < 0 {
				panic(fmt.Sprintf("illegal constraint format: %s", constraints))
			}
			continue
		}

		if constraints[i] == ',' || constraints[i] == '(' {
			arg := strings.TrimSpace(string(buf[:b]))

			if name == "" {
				if arg == "" {
					panic(fmt.Sprintf("illegal constraint format: %s", constraints))
				}
				name = arg
			} else {
				args = append(args, arg)
			}
			b = 0

			if brace == 0 {
				// Must be between constraints.
				// Add what we have and reset.
				results[name] = args
				name = ""
				args = []string{}
			}
			continue
		}
		buf[b] = constraints[i]
		b++
	}
	arg := strings.TrimSpace(string(buf[:b]))
	if name == "" {
		name = arg
	} else {
		args = append(args, arg)
	}
	if name != "" {
		results[name] = args
	}
	if brace != 0 {
		panic(fmt.Sprintf("illegal constraint format: %s", constraints))
	}
	return results
}

func newValidationHandler() ValidationHandler {
	v := elementValidationHandler{}
	v.validators = make(map[string]Validator)
	v.AddValidators(getDefaultValidators())
	return &v
}

// EmptyValidator is the default validator used to validate parameters where
// no constraint has been stipulated. It returns true in all cases
type EmptyValidator struct{}

// Validate returns true in all cases. This is the default validator.
func (v *EmptyValidator) Validate(val interface{}, args []string) bool {
	return true
}

// Type returns the constraint name. This is an empty string to
// ensure this validator is selected when no constraint has been
// specified in the route pattern parameter.
func (v *EmptyValidator) Type() string {
	return ""
}

// FailureMsg returns a string with a readable message about the validation failure.
// As this validator never fails, this method just returns an empty string.
func (v *EmptyValidator) FailureMsg() string {
	return ""
}

// Int32Validator tests for 32 bit integer values.
type Int32Validator struct{}

// Validate tests for 32 bit integer values.
// Returns true if val is a string containing an integer in the range -2147483648 to 2147483647
// Validate panics if val is not a string.
func (v *Int32Validator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("int32 validator can only validate strings not, %T", val))
	}
	_, err := strconv.ParseInt(s, 10, 32)
	return err == nil
}

// Type returns the constraint name (int32).
func (v *Int32Validator) Type() string {
	return "int32"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *Int32Validator) FailureMsg() string {
	return "must be a 32 bit integer."
}

// Int64Validator tests for 64 bit integer values.
type Int64Validator struct{}

// Validate tests for 64 bit integer values.
// Returns true if val is a string containing an integer in the range -9223372036854775808 to 9223372036854775807
// Validate panics if val is not a string.
func (v *Int64Validator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("int64 validator can only validate strings not, %T", val))
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// Type returns the constraint name (int64).
func (v *Int64Validator) Type() string {
	return "int64"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *Int64Validator) FailureMsg() string {
	return "must be a 64 bit integer."
}

// AlphaValidator tests for a sequence containing only alpha characters.
type AlphaValidator struct{}

// Validate tests for alpha values.
// Returns true if val is a string containing only characters in the ranges a-z or A-Z.
// Validate panics if val is not a string.
func (v *AlphaValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("alpha validator can only validate strings not, %T", val))
	}
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (alpha).
func (v *AlphaValidator) Type() string {
	return "alpha"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *AlphaValidator) FailureMsg() string {
	return "must contain only alpha characters."
}

// DigitsValidator tests for a sequence of digits.
type DigitsValidator struct{}

// Validate tests for digit values.
// Returns true if val is a string containing only digits 0-9.
// Validate panics if val is not a string.
func (v *DigitsValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("digits validator can only validate strings not, %T", val))
	}
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (digits).
func (v *DigitsValidator) Type() string {
	return "digits"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *DigitsValidator) FailureMsg() string {
	return "must contain only digit characters."
}

// Hex32Validator tests for 32 bit hex values.
type Hex32Validator struct{}

// Validate tests for 32 bit hex values.
// Returns true if val is a hexadecimal string in the range -80000000 to 7FFFFFFF.
// The test is not case sensitive, i.e. 3ef42bc7 and 3EF42BC7 will both return true.
// Validate panics if val is not a string.
func (v *Hex32Validator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("hex32 validator can only validate strings not, %T", val))
	}
	_, err := strconv.ParseInt(s, 16, 32)
	return err == nil
}

// Type returns the constraint name (hex32).
func (v *Hex32Validator) Type() string {
	return "hex32"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *Hex32Validator) FailureMsg() string {
	return "must be a 32 bit hexadecimal value."
}

// Hex64Validator tests for 64 bit hex values.
type Hex64Validator struct{}

// Validate tests for 64 bit hex values.
// Returns true if val is a hexadecimal string in the range -8000000000000000 to 7FFFFFFFFFFFFFFF.
// The test is not case sensitive, i.e. 3ef42bc7 and 3EF42BC7 will both return true.
// Validate panics if val is not a string.
func (v *Hex64Validator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("hex64 validator can only validate strings not, %T", val))
	}
	_, err := strconv.ParseInt(s, 16, 64)
	return err == nil
}

// Type returns the constraint name (hex64).
func (v *Hex64Validator) Type() string {
	return "hex64"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *Hex64Validator) FailureMsg() string {
	return "must be a 64 bit hexadecimal value."
}

// HexValidator tests for a sequence of hexadecimal characters.
type HexValidator struct{}

// Validate tests for hex values.
// Returns true if if val is a string containing only hex characters, (i.e. 0-9, a-e, A-F).
// Validate panics if val is not a string.
func (v *HexValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("hex validator can only validate strings not, %T", val))
	}
	re := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (hex).
func (v *HexValidator) Type() string {
	return "hex"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *HexValidator) FailureMsg() string {
	return "must contain only hexadecimal characters."
}

// UUIDValidator tests for UUIDs.
type UUIDValidator struct{}

// Validate tests for UUID values.
// Returns true if val is a string in one of the following formats:
//   xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//   {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
//   (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
//   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//   {xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx}
//   (xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
//   XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
//   {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
//   (XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX)
//   XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
//   {XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX}
//   (XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX)
//
// where X and x represent upper and lowercase hexadecimal values respectively.
//
// Valid UUID examples:
//  {58D5E212-165B-4CA0-909B-C86B9CEE0111}
//  {58d5e212-165b-4ca0-909b-c86b9cee0111}
//  58D5E212165B4CA0909BC86B9CEE0111
//
// Validate panics if val is not a string.
func (v *UUIDValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("uuid validator can only validate strings not, %T", val))
	}
	str := `^[{|\(]?[0-9a-fA-F]{8}[-]?([0-9a-fA-F]{4}[-]?){3}[0-9a-fA-F]{12}[\)|}]?$`
	re := regexp.MustCompile(str)
	if !re.MatchString(s) {
		return false
	}
	// ensure if we start or finish with a bookend, there is a matching one
	switch s[0] {
	case '{':
		return s[len(s)-1] == '}'
	case '(':
		return s[len(s)-1] == ')'
	}
	return s[len(s)-1] != ')' && s[len(s)-1] != '}'
}

// Type returns the constraint name (uuid).
func (v *UUIDValidator) Type() string {
	return "uuid"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *UUIDValidator) FailureMsg() string {
	return "must be a valid UUID."
}

// AlphaNumValidator tests for a sequence containing only alphanumeric characters.
type AlphaNumValidator struct{}

// Validate tests for alphanumeric values.
// Returns true if val is a string containing only characters in the ranges a-z, A-Z or 0-9.
// Validate panics if val is not a string.
func (v *AlphaNumValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("alphanum validator can only validate strings not, %T", val))
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (alphanum).
func (v *AlphaNumValidator) Type() string {
	return "alphanum"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *AlphaNumValidator) FailureMsg() string {
	return "must contain only alphanumeric characters."
}

// SentenceValidator tests for a sequence containing only characters normally found
// in a sentence.
type SentenceValidator struct{}

// Validate tests for characters normally found in a sentence.
// Returns true if val is a string containing only characters in the ranges a-z,
// A-Z or 0-9, plus spaces and punctuation characters ,.:;!? single and double quotes,
// braces and hyphens. Note that underscores are not included in the permissble set.
// Validate panics if val is not a string.
func (v *SentenceValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("sentence validator can only validate strings not, %T", val))
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9 ,.()'":;!?-]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (sentence).
func (v *SentenceValidator) Type() string {
	return "sentence"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *SentenceValidator) FailureMsg() string {
	return "must contain only characters found in a sentence."
}

// PhoneValidator tests for a telephone number.
type PhoneValidator struct{}

// Validate tests for phone number.
// Returns true if val is a string containing only characters found in a telephone
// number. e.g. +44 1752 123456, (1752) 123456 or 01752 123456
// Validate panics if val is not a string.
func (v *PhoneValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("phone validator can only validate strings not, %T", val))
	}
	re := regexp.MustCompile(`^(\+[0-9]+)?\s?((\([0-9]+\))|([0-9]+))?\s?[0-9]+\s?[0-9]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (phone).
func (v *PhoneValidator) Type() string {
	return "phone"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *PhoneValidator) FailureMsg() string {
	return "must contain a telephone number."
}

// PrefixValidator tests for a specified prefix.
type PrefixValidator struct {
	prefix string
}

// Validate tests for a prefix.
// Returns true if val is a string starting with the prefix specified in params.
// Validate panics if val is not a string.
func (v *PrefixValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("prefix validator can only validate strings not, %T", val))
	}

	v.prefix = ""
	if len(params) > 0 {
		v.prefix = params[0]
	}

	return strings.HasPrefix(s, v.prefix)
}

// Type returns the constraint name (prefix).
func (v *PrefixValidator) Type() string {
	return "prefix"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *PrefixValidator) FailureMsg() string {
	return fmt.Sprintf("must have the prefix %q.", v.prefix)
}

// SuffixValidator tests for a specified suffix.
type SuffixValidator struct {
	suffix string
}

// Validate tests for a suffix.
// Returns true if val is a string ending with the suffix specified in params.
// Validate panics if val is not a string.
func (v *SuffixValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("suffix validator can only validate strings not, %T", val))
	}
	v.suffix = ""
	if len(params) > 0 {
		v.suffix = params[0]
	}
	return strings.HasSuffix(s, v.suffix)
}

// Type returns the constraint name (suffix).
func (v *SuffixValidator) Type() string {
	return "suffix"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *SuffixValidator) FailureMsg() string {
	return fmt.Sprintf("must have the suffix %q.", v.suffix)
}

// MinValidator tests for a minumum numeric value.
type MinValidator struct {
	min string
}

// Validate tests for a minimum numerical value.
// Returns true if val is a number greater or equal to the value specified in params.
// Validate panics if val is not a number or supplied params argument is not a number.
func (v *MinValidator) Validate(val interface{}, params []string) bool {
	number, ok := normalizeNumber(val)
	if !ok {
		panic(fmt.Sprintf("min validator can only validate numbers not, %T", val))
	}
	v.min = ""
	if len(params) == 0 {
		panic("missing parameter for MinValidator")
	}
	v.min = params[0]
	switch reflect.TypeOf(number).Kind() {
	case reflect.Int64:
		p, err := strconv.ParseInt(params[0], 10, 64)
		if err == nil {
			return number.(int64) >= p
		}
	case reflect.Uint64:
		p, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil {
			return number.(uint64) >= p
		}
	case reflect.Float64:
		p, err := strconv.ParseFloat(params[0], 64)
		if err == nil {
			return number.(float64) >= p
		}
	}
	panic("non-numerical parameter used in MinValidator")
}

// Type returns the constraint name (min).
func (v *MinValidator) Type() string {
	return "min"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *MinValidator) FailureMsg() string {
	return fmt.Sprintf("must not be less than %v.", v.min)
}

// MaxValidator tests for a maxumum numeric value.
type MaxValidator struct {
	max string
}

// Validate tests for a maximum numerical value.
// Returns true if val is a number lower or equal to the value specified in params.
// Validate panics if val is not a number or supplied params argument is not a number.
func (v *MaxValidator) Validate(val interface{}, params []string) bool {
	number, ok := normalizeNumber(val)
	if !ok {
		panic(fmt.Sprintf("max validator can only validate numbers not, %T", val))
	}
	v.max = ""
	if len(params) == 0 {
		panic("missing parameter for MaxValidator")
	}
	v.max = params[0]
	switch reflect.TypeOf(number).Kind() {
	case reflect.Int64:
		p, err := strconv.ParseInt(params[0], 10, 64)
		if err == nil {
			return number.(int64) <= p
		}
	case reflect.Uint64:
		p, err := strconv.ParseUint(params[0], 10, 64)
		if err == nil {
			return number.(uint64) <= p
		}
	case reflect.Float64:
		p, err := strconv.ParseFloat(params[0], 64)
		if err == nil {
			return number.(float64) <= p
		}
	}
	panic("non-numerical parameter used in MaxValidator")
}

// Type returns the constraint name (max).
func (v *MaxValidator) Type() string {
	return "max"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *MaxValidator) FailureMsg() string {
	return fmt.Sprintf("must not be greater than %v.", v.max)
}

// RangeValidator tests for a numerical value in a given range.
type RangeValidator struct {
	min string
	max string
}

// Validate tests for a numerical value in a given range.
// Returns true if val is a number between the lower and upper limits specified in params.
// Example: range(2,4) - returns true if input is between 2 and 4 (inclusive).
// RangeValidator accepts all numeric types, e.g. range(3.123, 23456.89).
// Validate panics if val is not a number or supplied params arguments are not a numbers.
func (v *RangeValidator) Validate(val interface{}, params []string) bool {
	number, ok := normalizeNumber(val)
	if !ok {
		panic(fmt.Sprintf("range validator can only validate numbers not, %T", val))
	}

	v.min = ""
	v.max = ""
	if len(params) != 2 {
		panic("missing parameters for RangeValidator")
	}
	v.min = params[0]
	v.max = params[1]
	switch reflect.TypeOf(number).Kind() {
	case reflect.Int64:
		l, errl := strconv.ParseInt(params[0], 10, 64)
		u, erru := strconv.ParseInt(params[1], 10, 64)
		if errl != nil || erru != nil {
			break
		}
		return number.(int64) >= l && number.(int64) <= u
	case reflect.Uint64:
		l, errl := strconv.ParseUint(params[0], 10, 64)
		u, erru := strconv.ParseUint(params[1], 10, 64)
		if errl != nil || erru != nil {
			break
		}
		return number.(uint64) >= l && number.(uint64) <= u
	case reflect.Float64:
		l, errl := strconv.ParseFloat(params[0], 64)
		u, erru := strconv.ParseFloat(params[1], 64)
		if errl != nil || erru != nil {
			break
		}
		return number.(float64) >= l && number.(float64) <= u
	}
	panic("non-numerical parameters used in RangeValidator")
}

// Type returns the constraint name (range).
func (v *RangeValidator) Type() string {
	return "range"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *RangeValidator) FailureMsg() string {
	return fmt.Sprintf("must be between %v and %v.", v.min, v.max)
}

// LenMinValidator tests for a minimim length of String, Array, Slice or Map.
type LenMinValidator struct {
	min string
}

// Validate tests for a minimim length.
// Returns true if length of val is greater or equal to the value specified in params.
// Validate panics if val is not a String, Array, Slice or Map, or if supplied params argument is not an integer.
func (v *LenMinValidator) Validate(val interface{}, params []string) bool {
	if len(params) == 0 {
		panic("missing parameter for LenMinValidator")
	}
	v.min = params[0]
	l, err := strconv.Atoi(params[0])
	if err != nil {
		panic("non-integer parameter used in LenMinValidator")
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return reflect.ValueOf(val).Len() >= l
	default:
		panic(fmt.Sprintf("lenmin validator can only validate strings, arrays, slices and maps, not %T", val))
	}
}

// Type returns the constraint name (lenmin).
func (v *LenMinValidator) Type() string {
	return "lenmin"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *LenMinValidator) FailureMsg() string {
	return fmt.Sprintf("must contain at least %v elements.", v.min)
}

// LenMaxValidator tests for a maximum length of String, Array, Slice or Map.
type LenMaxValidator struct {
	max string
}

// Validate tests for a maximum length.
// Returns true if length of val is lower or equal to the value specified in params.
// Validate panics if val is not a String, Array, Slice or Map, or if supplied params argument is not an integer.
func (v *LenMaxValidator) Validate(val interface{}, params []string) bool {
	if len(params) == 0 {
		panic("missing parameter for LenMaxValidator")
	}
	v.max = params[0]
	u, err := strconv.Atoi(params[0])
	if err != nil {
		panic("non-integer parameter used in LenMaxValidator")
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return reflect.ValueOf(val).Len() <= u
	default:
		panic(fmt.Sprintf("lenmax validator can only validate strings, arrays, slices and maps, not %T", val))
	}
}

// Type returns the constraint name (lenmax).
func (v *LenMaxValidator) Type() string {
	return "lenmax"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *LenMaxValidator) FailureMsg() string {
	return fmt.Sprintf("must not contain more than %v elements.", v.max)
}

// LenRangeValidator tests for length of String, Array, Slice or Map, in a given range.
type LenRangeValidator struct {
	min string
	max string
}

// Validate tests for a length in a given range.
// Returns true if length of val is between the lower and upper limits specified in params.
// Validate panics if val is not a String, Array, Slice or Map, or if supplied params arguments are not an integer.
func (v *LenRangeValidator) Validate(val interface{}, params []string) bool {
	if len(params) != 2 {
		panic("missing parameters for LenRangeValidator")
	}
	v.min = params[0]
	v.max = params[1]
	l, errl := strconv.Atoi(params[0])
	u, erru := strconv.Atoi(params[1])
	if errl != nil || erru != nil {
		panic("non-integer parameters used in LenRangeValidator")
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		lenval := reflect.ValueOf(val).Len()
		return lenval >= l && lenval <= u
	default:
		panic(fmt.Sprintf("lenrange validator can only validate strings, arrays, slices and maps, not %T", val))
	}
}

// Type returns the constraint name (lenrange).
func (v *LenRangeValidator) Type() string {
	return "lenrange"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *LenRangeValidator) FailureMsg() string {
	return fmt.Sprintf("must contain between %v and %v elements.", v.min, v.max)
}

// ContainsValidator tests whether a container holds a specific string.
type ContainsValidator struct {
	required string
}

// Validate tests for a existence of a string within another string, Array, Slice
// or Map (keys).
// Returns true if val is a String, Array, Slice or Map containing the string
// specified in params. Contains is case-sensitive.
// Validate panics if val is not a String, Array, Slice or Map.
func (v *ContainsValidator) Validate(val interface{}, params []string) bool {
	v.required = ""
	if len(params) > 0 {
		v.required = params[0]
	}

	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			el := reflect.ValueOf(key.Interface())
			if el.Kind() == reflect.Ptr {
				el = el.Elem()
			}
			if el.Kind() != reflect.String {
				panic(fmt.Sprintf("contains validator can only validate maps with keys of string, not %T", key.Interface()))
			}
			if el.String() == v.required {
				return true
			}
		}
		return false
	case reflect.Array, reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			el := rv.Index(i)
			if el.Kind() == reflect.Ptr {
				el = el.Elem()
			}
			if el.Kind() != reflect.String {
				panic(fmt.Sprintf("contains validator can only validate arrays and slices of string, not %T", el.Interface()))
			}
			if el.String() == v.required {
				return true
			}
		}
		return false
	case reflect.String:
		return strings.Contains(val.(string), v.required)
	default:
		panic(fmt.Sprintf("contains validator can only validate strings, arrays, slices and maps, not %T", val))
	}
}

// Type returns the constraint name (contains).
func (v *ContainsValidator) Type() string {
	return "contains"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *ContainsValidator) FailureMsg() string {
	return fmt.Sprintf("must contain %q.", v.required)
}

// InSetValidator tests a value is in a set.
type InSetValidator struct {
	set string
}

// Validate tests for a value within a set of values.
// Returns true if val is a string or int within the set specified in params.
// Validate panics if val is not a string or int.
func (v *InSetValidator) Validate(val interface{}, params []string) bool {

	v.set = strings.Join(params, ", ")
	switch reflect.TypeOf(val).Kind() {
	case reflect.String:
		for _, p := range params {
			s := strings.TrimSpace(p)
			if val.(string) == s {
				return true
			}
		}
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		v := reflect.ValueOf(val).Int()

		for _, p := range params {
			s := strings.TrimSpace(p)
			i, err := strconv.ParseInt(s, 10, 64) //strconv.Atoi(s)
			if err != nil {
				panic("non-integer parameter used in InSetValidator")
			}
			if v == i {
				return true
			}
		}
	default:
		panic(fmt.Sprintf("inset validator can only validate strings and ints, not %T", val))
	}
	return false
}

// Type returns the constraint name (inset).
func (v *InSetValidator) Type() string {
	return "inset"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *InSetValidator) FailureMsg() string {
	return fmt.Sprintf("must be in the set [%s].", v.set)
}

// NotEmptyValidator tests for an empty String, Array, Slice or Map.
type NotEmptyValidator struct{}

// Validate tests for an empty String, Array, Slice or Map.
// Returns true if val is String, Array, Slice or Map with elements.
// Equivlent to (and shorthand for) minlen(1).
// Validate panics if val is not a String, Array, Slice or Map.
func (v *NotEmptyValidator) Validate(val interface{}, params []string) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return reflect.ValueOf(val).Len() > 0
	default:
		panic(fmt.Sprintf("notempty validator can only validate strings, arrays, slices and maps, not %T", val))
	}
}

// Type returns the constraint name (notempty).
func (v *NotEmptyValidator) Type() string {
	return "notempty"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *NotEmptyValidator) FailureMsg() string {
	return "must not be empty."
}

// NotZeroValidator tests for a value of zero.
type NotZeroValidator struct{}

// Validate tests for a numerical value of zero.
// Returns true if val is a number not equal to zero.
// Validate panics if val is not a number.
func (v *NotZeroValidator) Validate(val interface{}, params []string) bool {
	number, ok := normalizeNumber(val)
	if !ok {
		panic(fmt.Sprintf("notzero validator can only validate numbers not, %T", val))
	}

	switch reflect.TypeOf(number).Kind() {
	case reflect.Int64:
		return number.(int64) != 0
	case reflect.Uint64:
		return number.(uint64) != 0
	case reflect.Float64:
		return number.(float64) != 0
	}
	return false
}

// Type returns the constraint name (notzero).
func (v *NotZeroValidator) Type() string {
	return "notzero"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *NotZeroValidator) FailureMsg() string {
	return "must not be zero."
}

// NotNilValidator tests for an uninitialized map or slice, or nil pointer.
type NotNilValidator struct{}

// Validate tests for an uninitialized map or slice, or nil pointer.
// Returns true if val is an initialized map or slice, or non-nil pointer.
// Validate panics if val is not a map, slice or pointer.
func (v *NotNilValidator) Validate(val interface{}, params []string) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		return reflect.ValueOf(val).Pointer() != 0

	default:
		panic(fmt.Sprintf("notnil validator can only validate maps, slices and pointers, not %T", val))
	}
}

// Type returns the constraint name (notnil).
func (v *NotNilValidator) Type() string {
	return "notnil"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *NotNilValidator) FailureMsg() string {
	return "must not be nil."
}

// EmailValidator tests for a valid email address
type EmailValidator struct{}

// Validate tests for an email address.
// Returns true if val is an email address.
// Validate panics if val is not a string.
func (v *EmailValidator) Validate(val interface{}, params []string) bool {
	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("email validator can only validate strings not, %T", val))
	}
	//re := regexp.MustCompile(`^(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})$`)
	// use HTML5 email regex:
	const exp = `^[a-zA-Z0-9.!#$%&’*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`
	re := regexp.MustCompile(exp)
	return re.MatchString(s)
}

// Type returns the constraint name (email).
func (v *EmailValidator) Type() string {
	return "email"
}

// FailureMsg returns a string with a readable message about the validation failure.
func (v *EmailValidator) FailureMsg() string {
	return "must be a valid email address."
}

func normalizeNumber(i interface{}) (interface{}, bool) {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return reflect.ValueOf(i).Int(), true
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return reflect.ValueOf(i).Uint(), true
	case reflect.Float32,
		reflect.Float64:
		return reflect.ValueOf(i).Float(), true
	}
	return nil, false
}

func getDefaultValidators() []Validator {
	return []Validator{
		&EmptyValidator{},
		&Int32Validator{},
		&Int64Validator{},
		&AlphaValidator{},
		&DigitsValidator{},
		&Hex32Validator{},
		&Hex64Validator{},
		&HexValidator{},
		&UUIDValidator{},
		&AlphaNumValidator{},
		&PrefixValidator{},
		&SuffixValidator{},
		&MinValidator{},
		&MaxValidator{},
		&RangeValidator{},
		&LenMinValidator{},
		&LenMaxValidator{},
		&LenRangeValidator{},
		&ContainsValidator{},
		&InSetValidator{},
		&NotEmptyValidator{},
		&NotZeroValidator{},
		&NotNilValidator{},
		&SentenceValidator{},
		&PhoneValidator{},
	}
}

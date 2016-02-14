package mango

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParamValidator is the interface describing a route parameter validator.
// ParamValidators are used to validate sections of a URL path which match
// the pattern of an entry in the routing tree.
type ParamValidator interface {
	// Validate tests if s matches the validation rules. The validation test
	// may involve constraint specific args.

	Validate(s string, args []string) bool
	// Type returns the constraint name used in routing patterns.
	// RouteParamValidators will use this to locate the correct validator.
	Type() string
}

// RouteParamValidators is the interface for a collection of ParamValidators.
// A route parameter value can be checked against its constraint using the
// IsValid method. New validators can be added individually using AddValidator
// or as a collection using AddValidators.
type RouteParamValidators interface {
	AddValidator(v ParamValidator)
	AddValidators(validators []ParamValidator)
	IsValid(s, constraint string) bool
}

type parameterValidators struct {
	validators map[string]ParamValidator
}

func (r *parameterValidators) AddValidator(v ParamValidator) {
	r.validators[v.Type()] = v
}

func (r *parameterValidators) AddValidators(validators []ParamValidator) {
	for _, v := range validators {
		r.AddValidator(v)
	}
}

func (r *parameterValidators) IsValid(s, constraint string) bool {
	var args []string
	if strings.HasSuffix(constraint, ")") {
		i := strings.IndexByte(constraint, byte('('))
		if i < 1 {
			panic(fmt.Sprintf("illegal constraint format: %s", constraint))
		}
		args = strings.Split(constraint[i+1:len(constraint)-1], ",")
		for i, p := range args {
			args[i] = strings.TrimSpace(p)
		}
		constraint = constraint[:i]
	} else if strings.IndexByte(constraint, byte('(')) > -1 {
		panic(fmt.Sprintf("illegal constraint format: %s", constraint))
	}

	constraint = strings.TrimSpace(constraint)
	v, ok := r.validators[constraint]
	if !ok {
		panic(fmt.Sprintf("unknown constraint: %s", constraint))
	}
	return v.Validate(s, args)
}

func newParameterValidators() RouteParamValidators {
	v := parameterValidators{}
	v.validators = make(map[string]ParamValidator)
	v.AddValidators(getDefaultValidators())
	return &v
}

// EmptyValidator is the default validator used to validate parameters where
// no constraint has been stipulated. It returns true in all cases
type EmptyValidator struct{}

// Validate returns true in all cases. This is the default validator.
func (EmptyValidator) Validate(s string, args []string) bool {
	return true
}

// Type returns the constraint name. This is an empty string to
// ensure this valiadtor is selected when no constraint has been
// specified in the route pattern parameter.
func (EmptyValidator) Type() string {
	return ""
}

// Int32Validator tests for 32 bit integer values.
type Int32Validator struct{}

// Validate tests for 32 bit integer values.
// Returns true if s is an integer in the range -2147483648 to 2147483647
func (Int32Validator) Validate(s string, params []string) bool {
	_, err := strconv.ParseInt(s, 10, 32)
	return err == nil
}

// Type returns the constraint name (int32).
func (Int32Validator) Type() string {
	return "int32"
}

// Int64Validator tests for 64 bit integer values.
type Int64Validator struct{}

// Validate tests for 64 bit integer values.
// Returns true if s is an integer in the range -9223372036854775808 to 9223372036854775807
func (Int64Validator) Validate(s string, params []string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// Type returns the constraint name (int64).
func (Int64Validator) Type() string {
	return "int64"
}

// AlphaValidator tests for a sequence containing only alpha characters.
type AlphaValidator struct{}

// Validate tests for alpha values.
// Returns true if s contains only characters in the ranges a-z and A-Z.
func (AlphaValidator) Validate(s string, params []string) bool {
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (alpha).
func (AlphaValidator) Type() string {
	return "alpha"
}

// DigitsValidator tests for a sequence of digits.
type DigitsValidator struct{}

// Validate tests for digit values.
// Returns true if s contains only digits 0-9.
func (DigitsValidator) Validate(s string, params []string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (digits).
func (DigitsValidator) Type() string {
	return "digits"
}

// Hex32Validator tests for 32 bit hex values.
type Hex32Validator struct{}

// Validate tests for 32 bit hex values.
// Returns true if s is hexadecimal in the range -80000000 to 7FFFFFFF.
// The test is not case sensitive, i.e. 3ef42bc7 and 3EF42BC7 will both return true.
func (Hex32Validator) Validate(s string, params []string) bool {
	_, err := strconv.ParseInt(s, 16, 32)
	return err == nil
}

// Type returns the constraint name (hex32).
func (Hex32Validator) Type() string {
	return "hex32"
}

// Hex64Validator tests for 64 bit hex values.
type Hex64Validator struct{}

// Validate tests for 64 bit hex values.
// Returns true if s is hexadecimal in the range -8000000000000000 to 7FFFFFFFFFFFFFFF.
// The test is not case sensitive, i.e. 3ef42bc7 and 3EF42BC7 will both return true.
func (Hex64Validator) Validate(s string, params []string) bool {
	_, err := strconv.ParseInt(s, 16, 64)
	return err == nil
}

// Type returns the constraint name (hex64).
func (Hex64Validator) Type() string {
	return "hex64"
}

// HexValidator tests for a sequence of hexadecimal characters.
type HexValidator struct{}

// Validate tests for hex values.
// Returns true if s contains only hex characters, (i.e. 0-9, a-e, A-F).
func (HexValidator) Validate(s string, params []string) bool {
	re := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	return re.MatchString(s)
}

// Type returns the constraint name (hex).
func (HexValidator) Type() string {
	return "hex"
}

// UUIDValidator tests for UUIDs.
type UUIDValidator struct{}

// Validate tests for UUID values.
// Returns true if s is in one of the following formats:
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
func (UUIDValidator) Validate(s string, params []string) bool {
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
func (UUIDValidator) Type() string {
	return "uuid"
}

func getDefaultValidators() []ParamValidator {
	return []ParamValidator{
		EmptyValidator{},
		Int32Validator{},
		Int64Validator{},
		AlphaValidator{},
		DigitsValidator{},
		Hex32Validator{},
		Hex64Validator{},
		HexValidator{},
		UUIDValidator{},
	}
}

/*
TODO: Add validators for these

alphanumeric, bool, float32, float64, uint32, uint64:

b, err := strconv.ParseBool("true")
f, err := strconv.ParseFloat("3.1415", 64)
uint64, err := strconv.ParseUint("42", 10, 64)

parameterised constraints:

min(minimum) - Allows only integer values with the specified minimum value.
max(maximum) - Allows only integer values with the specified maximum value.
range(minimum, maximum) - Allows only integer values within the specified range. (Between minimum and maximum)
minlength(length) - Allows only values longer than the specified minimum length.
maxlength(length) - Allows only values shorter that the maximum length.
length(minimum, maximum) - Allows only values with length within the specified range. (Between minimum and maximum)
*/

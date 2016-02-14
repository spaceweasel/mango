package mango

import "testing"

func TestEmptyValidatorHasEmptyStringType(t *testing.T) {
	want := ""

	v := EmptyValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Validator type = %q, want %q", got, want)
	}
}

func TestEmptyValidatorReturnsTrueForAllInput(t *testing.T) {
	want := true

	v := EmptyValidator{}
	args := []string{}
	got := v.Validate("mango", args)

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

//int32
func TestInt32ValidatorType(t *testing.T) {
	want := "int32"

	v := Int32Validator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestInt32Validator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"34566", []string{}, true, "Int"},
		{"mango", []string{}, false, "Non-Int"},
		{"2147483647", []string{}, true, "== MaxInt32"},
		{"2147483648", []string{}, false, "> MaxInt32"},
		{"-2147483648", []string{}, true, "== MaxInt32"},
		{"-2147483649", []string{}, false, "< MaxInt32"},
	}

	v := Int32Validator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//int64
func TestInt64ValidatorType(t *testing.T) {
	want := "int64"

	v := Int64Validator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestInt64Validator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"34566", []string{}, true, "Int"},
		{"mango", []string{}, false, "Non-Int"},
		{"9223372036854775807", []string{}, true, "== MaxInt64"},
		{"9223372036854775808", []string{}, false, "> MaxInt64"},
		{"-9223372036854775808", []string{}, true, "== MaxInt64"},
		{"-9223372036854775809", []string{}, false, "< MaxInt64"},
	}

	v := Int64Validator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//alpha
func TestAlphaValidatorType(t *testing.T) {
	want := "alpha"

	v := AlphaValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestAlphaValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"KHKSDFHIASHJHGJHKJHKHKJHKKDEWQW", []string{}, true, "Uppercase"},
		{"bdekjaskjdjhgbjhgjhmnbmmgjhfaksj", []string{}, true, "Lowercase"},
		{"KHasdasdKSDFHasdadadIASDEWQWfdsf", []string{}, true, "Mixedcase"},
		{"7997665764359698", []string{}, false, "Digits"},
		{"3A456DE63A456DE63A456DE6", []string{}, false, "Hex"},
		{"bdek jaskjdjhg bjhgj  hmnbmm gjhfaksj", []string{}, false, "Spaces"},
		{"bde!kjask,jdjhgbjhgjhm?nbmmgjh;faksj", []string{}, false, "Punctuation"},
		{"bdekj_askjdj_hgbjhgjhmnbm_mgjhfaksj", []string{}, false, "Underscores"},
		{"bdek-jaskjd-jhgbjhgj-hmnbmmg-jhfaksj", []string{}, false, "Hyphens"},
		{"bdekja.skjdjhg.bjhgjhmnbmm.gjhfaksj", []string{}, false, "Periods"},
	}

	v := AlphaValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//digits
func TestDigitsValidatorType(t *testing.T) {
	want := "digits"

	v := DigitsValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestDigitsValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"923372036233720368547753372036807", []string{}, true, "Digits"},
		{"92233ASD72036fhghgf854775808", []string{}, false, "Alpha"},
		{"92233-72036-854775808", []string{}, false, "Hyphens"},
		{"92233.72036854775808", []string{}, false, "Periods"},
		{"92E23F358A08B8976D", []string{}, false, "Hex"},
		{"92337203 623372036854775337 2036807", []string{}, false, "Spaces"},
	}

	v := DigitsValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//hex32
func TestHex32ValidatorType(t *testing.T) {
	want := "hex32"

	v := Hex32Validator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestHex32Validator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"3A456DE6", []string{}, true, "Uppercase"},
		{"3a456de6", []string{}, true, "Lowercase"},
		{"79959698", []string{}, true, "Digits only"},
		{"abdafec", []string{}, true, "a-f only"},
		{"mango", []string{}, false, "Non-Hex"},
		{"7FFFFFFF", []string{}, true, "== MaxHex32"},
		{"80000000", []string{}, false, "> MaxHex32"},
		{"-80000000", []string{}, true, "== MaxHex32"},
		{"-80000001", []string{}, false, "< MaxHex32"},
	}

	v := Hex32Validator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//Hex64
func TestHex64ValidatorType(t *testing.T) {
	want := "hex64"

	v := Hex64Validator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestHex64Validator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"3A456DE6", []string{}, true, "Uppercase"},
		{"3a456de6", []string{}, true, "Lowercase"},
		{"2342341319", []string{}, true, "Digits only"},
		{"abeedafecb", []string{}, true, "a-f only"},
		{"mango", []string{}, false, "Non-Hex"},
		{"7FFFFFFFFFFFFFFF", []string{}, true, "== MaxHex64"},
		{"8000000000000000", []string{}, false, "> MaxHex64"},
		{"-8000000000000000", []string{}, true, "== MaxHex64"},
		{"-8000000000000001", []string{}, false, "< MaxHex64"},
	}

	v := Hex64Validator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//Hex
func TestHexValidatorType(t *testing.T) {
	want := "hex"

	v := HexValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestHexValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"3A4561BA23EF8DCDE6", []string{}, true, "Uppercase"},
		{"3a456de61ab23ef8dc", []string{}, true, "Lowercase"},
		{"987aFE0A8dcd98a2b3eF", []string{}, true, "Mixedcase"},
		{"7997665764359698", []string{}, true, "Digits only"},
		{"abaffcedaefabdafec", []string{}, true, "a-f only"},
		{"ABAFFCEDAEFABDAFEC", []string{}, true, "A-F only"},
		{"mango", []string{}, false, "Non-Hex"},
		{"aba ffceda efab dafec", []string{}, false, "Spaces"},
		{"aba?ffc!edae;fabda,fec", []string{}, false, "Punctuation"},
		{"abaf_fce_daefabd_afec", []string{}, false, "Underscores"},
		{"ab-affcedae-fabd-afec", []string{}, false, "Hyphens"},
		{"aba.ffced.aefabda.fec", []string{}, false, "Periods"},
	}

	v := HexValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

//UUID
func TestUUIDValidatorType(t *testing.T) {
	want := "uuid"

	v := UUIDValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestUUIDValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"7A8CA1EA-B53B-4231-9260-D33F652F1ED9", []string{}, true, "Uppercase"},
		{"7a8ca1ea-b53b-4231-9260-d33f652f1ed9", []string{}, true, "Lowercase"},
		{"7a8ca1EA-B53B-4231-9260-d33F652f1ed9", []string{}, true, "Mixedcase"},
		{"7a8ca1eab53b42319260d33f652f1ed9", []string{}, true, "Hyphenless"},
		{"7a8ca1-eab53b42-3192-60d33f-652f1ed9", []string{}, false, "Incorrect hyphen position"},
		{"7a8ca1EA.B53B-4231-9260-Z33F652T1ed9", []string{}, false, "Non Hex or Hyphen"},
		{"7a8ca1eab53b42319260d33f652", []string{}, false, "Too short"},
		{"7a8ca1eab53b42319260d33f652f1ed9ab45ce", []string{}, false, "Too long"},
		{"{7a8ca1ea-b53b-4231-9260-d33f652f1ed9}", []string{}, true, "Hyphens and curly braces"},
		{"{7a8ca1eab53b42319260d33f652f1ed9}", []string{}, true, "Curly braces"},
		{"(7a8ca1ea-b53b-4231-9260-d33f652f1ed9)", []string{}, true, "Hyphens and plain braces"},
		{"(7a8ca1eab53b42319260d33f652f1ed9)", []string{}, true, "Plain braces"},
		{"{7a8ca1ea-b53b-4231-9260-d33f652f1ed9)", []string{}, false, "Curly start and plain end"},
		{"(7a8ca1ea-b53b-4231-9260-d33f652f1ed9}", []string{}, false, "Plain start and curly end"},
		{"{7a8ca1ea-b53b-4231-9260-d33f652f1ed9", []string{}, false, "Curly start"},
		{"7a8ca1ea-b53b-4231-9260-d33f652f1ed9}", []string{}, false, "Curly end"},
		{"7a8ca1ea-b53b-4231-9260-d33f652f1ed9)", []string{}, false, "Plain end"},
		{"(7a8ca1ea-b53b-4231-9260-d33f652f1ed9", []string{}, false, "Plain start "},
	}

	v := UUIDValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate(%q) = %v (%s)", test.input, got, test.comment)
		}
	}
}

// End of validators

func TestNewParameterValidatorsReturnsWithInitialisedValidatorsMap(t *testing.T) {
	want := len(getDefaultValidators())

	pv := newParameterValidators()

	cpv := pv.(*parameterValidators)
	got := len(cpv.validators)

	if got != want {
		t.Errorf("Validator count = %d, want %d", got, want)
	}
}

func TestNewParameterValidatorsHasDefaultValidators(t *testing.T) {
	pv := newParameterValidators()
	cpv := pv.(*parameterValidators)
	for _, dv := range getDefaultValidators() {
		v, ok := cpv.validators[dv.Type()]
		want := dv.Type()
		if !ok {
			t.Errorf("Validator Type = nil, want %q", want)
			continue
		}
		got := v.Type()
		if got != want {
			t.Errorf("Validator Type = %q, want %q", got, want)
		}
	}
}

func TestAddValidatorsAddsToCollection(t *testing.T) {
	want := 3

	pv := parameterValidators{}
	pv.validators = make(map[string]ParamValidator)
	v := []ParamValidator{EmptyValidator{}, testValidator1{}, testValidator2{}}
	pv.AddValidators(v)

	got := len(pv.validators)

	if got != want {
		t.Errorf("Validator count = %d, want %d", got, want)
	}
}

func TestIsValidPanicsWhenUnknownConstraint(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pv := newParameterValidators()
	pv.IsValid("validator1", "test1")
}

func TestIsValidPanicsWithCorrectErrorMessageWhenUnknownConstraint(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "unknown constraint: test1"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	pv := newParameterValidators()
	pv.IsValid("validator1", "test1")
}

func TestIsValidDoesNotPanicIWhenKnownConstraint(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()

	pv := newParameterValidators()
	pv.AddValidator(testValidator1{})
	pv.IsValid("validator1", "test1")
}

func TestIsValidPassesParamValueToValidator(t *testing.T) {
	want := true

	pv := newParameterValidators()
	pv.AddValidator(testValidator1{})

	got := pv.IsValid("validator1", "test1")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	got = pv.IsValid("validator2", "test1")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidParsesConstraintArgsWhenSingleArg(t *testing.T) {
	want := true

	pv := newParameterValidators()
	pv.AddValidator(testValidator2{})

	got := pv.IsValid("paramValue", "test2(arg)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	got = pv.IsValid("paramValue", "test2(6)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidParsesConstraintArgsWhenMultipleArgs(t *testing.T) {
	want := true

	pv := newParameterValidators()
	pv.AddValidator(testValidator3{})

	got := pv.IsValid("paramValue", "test3(6,arg2)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	got = pv.IsValid("paramValue", "test3(6,arg)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidTrimsSpaceAroundConstraintArgs(t *testing.T) {
	want := true

	pv := newParameterValidators()
	pv.AddValidator(testValidator3{})

	got := pv.IsValid("paramValue", "test3(  6 ,  arg2 )")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	got = pv.IsValid("paramValue", "test3( 6 , arg )")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidPanicsWhenConstraintMalformed(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pv := newParameterValidators()
	pv.AddValidator(testValidator3{})

	pv.IsValid("paramValue", "test3 arg2 )")
}

func TestIsValidPanicsWithCorrectErrorMessageWhenConstraintWithOnlyClosingBracket(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "illegal constraint format: test3 arg2)"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	pv := newParameterValidators()
	pv.AddValidator(testValidator3{})

	pv.IsValid("paramValue", "test3 arg2)")
}

func TestIsValidPanicsWithCorrectErrorMessageWhenConstraintWithOnlyOpeningBracket(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "illegal constraint format: test3(arg2"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	pv := newParameterValidators()
	pv.AddValidator(testValidator3{})

	pv.IsValid("paramValue", "test3(arg2")
}

func TestIsValidTrimsSpaceAroundConstraintName(t *testing.T) {
	want := true

	pv := newParameterValidators()
	pv.AddValidator(testValidator3{})

	got := pv.IsValid("paramValue", "  test3  (6,arg2)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	got = pv.IsValid("paramValue", "  test3  (6,arg)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

type testValidator1 struct{}

func (testValidator1) Validate(s string, args []string) bool {
	return s == "validator1"
}

func (testValidator1) Type() string {
	return "test1"
}

type testValidator2 struct{}

func (testValidator2) Validate(s string, args []string) bool {
	return len(args) == 1 && args[0] == "arg"
}

func (testValidator2) Type() string {
	return "test2"
}

type testValidator3 struct{}

func (testValidator3) Validate(s string, args []string) bool {
	return len(args) == 2 && args[1] == "arg2"
}

func (testValidator3) Type() string {
	return "test3"
}

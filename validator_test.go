package mango

import (
	"fmt"
	"testing"
)

func TestEmptyValidatorHasEmptyStringType(t *testing.T) {
	want := ""

	v := EmptyValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Validator type = %q, want %q", got, want)
	}
}

func TestEmptyValidatorFailureMessage(t *testing.T) {
	want := ""

	v := EmptyValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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

func TestInt32ValidatorFailureMessage(t *testing.T) {
	want := "Must be a 32 bit integer."

	v := Int32Validator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q int32 = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestInt32ValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "int32 validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := Int32Validator{}
	v.Validate(32, []string{})
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

func TestInt64ValidatorFailureMessage(t *testing.T) {
	want := "Must be a 64 bit integer."

	v := Int64Validator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q int64 = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestInt64ValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "int64 validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := Int64Validator{}
	v.Validate(32, []string{})
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

func TestAlphaValidatorFailureMessage(t *testing.T) {
	want := "Must contain only alpha characters."

	v := AlphaValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q alpha = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestAlphaValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "alpha validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := AlphaValidator{}
	v.Validate(32, []string{})
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

func TestDigitsValidatorFailureMessage(t *testing.T) {
	want := "Must contain only digit characters."

	v := DigitsValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q digits = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestDigitsValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "digits validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := DigitsValidator{}
	v.Validate(32, []string{})
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

func TestHex32ValidatorFailureMessage(t *testing.T) {
	want := "Must be a 32 bit hexadecimal value."

	v := Hex32Validator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q hex32 = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestHex32ValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "hex32 validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := Hex32Validator{}
	v.Validate(32, []string{})
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
func TestHex64ValidatorFailureMessage(t *testing.T) {
	want := "Must be a 64 bit hexadecimal value."

	v := Hex64Validator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q hex64 = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestHex64ValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "hex64 validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := Hex64Validator{}
	v.Validate(32, []string{})
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
func TestHexValidatorFailureMessage(t *testing.T) {
	want := "Must contain only hexadecimal characters."

	v := HexValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q hex = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestHexValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "hex validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := HexValidator{}
	v.Validate(32, []string{})
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
func TestUUIDValidatorFailureMessage(t *testing.T) {
	want := "Must be a valid UUID."

	v := UUIDValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
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
			t.Errorf("Validate (%s): %q uuid = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestUUIDValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "uuid validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := UUIDValidator{}
	v.Validate(32, []string{})
}

//alphanum
func TestAlphaNumValidatorType(t *testing.T) {
	want := "alphanum"

	v := AlphaNumValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestAlphaNumValidatorFailureMessage(t *testing.T) {
	want := "Must contain only alphanumeric characters."

	v := AlphaNumValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}

func TestAlphaNumValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"K3HAS34DASDKSHA75SDA1DADIASE48WF7DSF", []string{}, true, "Uppercase"},
		{"k3has34dasdksha75sda1dadiase48wf7dsf", []string{}, true, "Lowercase"},
		{"K3Has34dasdKSHa75sda1dadIASE48Wf7dsf", []string{}, true, "Mixedcase"},
		{"KHasdasdKSDFHasdadadIASDEWQWfdsf", []string{}, true, "Alpha only"},
		{"799766576435969875764448", []string{}, true, "Digits only"},
		{"3A456DE63A456DE63A456DE6", []string{}, true, "Hex"},
		{"K3Has 34dasdKSHa 75sda1da  dIASE 48Wf7dsf", []string{}, false, "Spaces"},
		{"K3!Has34dasd?KSHa75sd,a1da;dIASE48Wf7dsf", []string{}, false, "Punctuation"},
		{"K3Has34_dasdK_SHa75sda1dad_IASE48Wf7_dsf", []string{}, false, "Underscores"},
		{"K3Has3-4dasdKSHa7-5sda1dadIA-SE48W-f7dsf", []string{}, false, "Hyphens"},
		{"K3Has34.dasdKS.Ha75sda1d.adIASE4.8Wf7dsf", []string{}, false, "Periods"},
	}

	v := AlphaNumValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %q alphanum = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestAlphaNumValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "alphanum validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := AlphaNumValidator{}
	v.Validate(32, []string{})
}

//prefix
func TestPrefixValidatorType(t *testing.T) {
	want := "prefix"

	v := PrefixValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestPrefixValidatorFailureMessage(t *testing.T) {
	want := "Must have the correct prefix."

	v := PrefixValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}

func TestPrefixValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"CHEESEBICYCLE", []string{"Cheese"}, false, "Uppercase"},
		{"cheesebicycle", []string{"Cheese"}, false, "Lowercase"},
		{"CheeseBicycle", []string{"Cheese"}, true, "Mixedcase"},
		{"CheeseBicycle", []string{"heese"}, false, "Offset"},
		{"CheeseBicycle", []string{"Bicycle"}, false, "Suffix"},
		{"Ch33seBicycle", []string{"Ch33se"}, true, "Digits"},
		{"Cheese Bicycle", []string{"Cheese"}, true, "End Spaces"},
		{"Che eseBicycle", []string{"Che ese"}, true, "Mid Spaces"},
		{"Cheese!Bicycle", []string{"Cheese!B"}, true, "Punctuation"},
		{"Cheese_Bicycle", []string{"Cheese_B"}, true, "Underscores"},
		{"Cheese-Bicycle", []string{"Cheese-B"}, true, "Hyphens"},
		{"Cheese.Bicycle", []string{"Cheese.B"}, true, "Periods"},
		{"CheeseBicycle", []string{}, true, "Empty"},
	}

	v := PrefixValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %q prefix(%s) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestPrefixValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "prefix validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := PrefixValidator{}
	v.Validate(32, []string{})
}

//suffix
func TestSuffixValidatorType(t *testing.T) {
	want := "suffix"

	v := SuffixValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestSuffixValidatorFailureMessage(t *testing.T) {
	want := "Must have the correct suffix."

	v := SuffixValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}

func TestSuffixValidator(t *testing.T) {
	var tests = []struct {
		input   string
		args    []string
		want    bool
		comment string
	}{
		{"CHEESEBICYCLE", []string{"Bicycle"}, false, "Uppercase"},
		{"cheesebicycle", []string{"Bicycle"}, false, "Lowercase"},
		{"CheeseBicycle", []string{"Bicycle"}, true, "Mixedcase"},
		{"CheeseBicycle", []string{"Bicycl"}, false, "Offset"},
		{"CheeseBicycle", []string{"Cheese"}, false, "Prefix"},
		{"Cheese81cycl3", []string{"81cycl3"}, true, "Digits"},
		{"Cheese Bicycle", []string{"Bicycle"}, true, "Start Spaces"},
		{"CheeseBi cycle", []string{"Bi cycle"}, true, "Mid Spaces"},
		{"Cheese!Bicycle", []string{"e!Bicycle"}, true, "Punctuation"},
		{"Cheese_Bicycle", []string{"e_Bicycle"}, true, "Underscores"},
		{"Cheese-Bicycle", []string{"e-Bicycle"}, true, "Hyphens"},
		{"Cheese.Bicycle", []string{"e.Bicycle"}, true, "Periods"},
		{"CheeseBicycle", []string{}, true, "Empty"},
	}

	v := SuffixValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %q suffix(%s) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestSuffixValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "suffix validator can only validate strings not, int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := SuffixValidator{}
	v.Validate(32, []string{})
}

//min
func TestMinValidatorType(t *testing.T) {
	want := "min"

	v := MinValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestMinValidatorFailureMessage(t *testing.T) {
	want := "Must not be less than the minimum permitted."

	v := MinValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}

func TestMinValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{234, []string{"56"}, true, "int"},
		{45, []string{"56"}, false, "int"},
		{int8(127), []string{"56"}, true, "int8"},
		{int8(45), []string{"56"}, false, "int8"},
		{int16(234), []string{"56"}, true, "int16"},
		{int16(45), []string{"56"}, false, "int16"},
		{int32(234), []string{"56"}, true, "int32"},
		{int32(45), []string{"56"}, false, "int32"},
		{int64(234), []string{"56"}, true, "int64"},
		{int64(45), []string{"56"}, false, "int64"},
		{uint(234), []string{"56"}, true, "uint"},
		{uint(45), []string{"56"}, false, "uint"},
		{uint8(234), []string{"56"}, true, "uint8"},
		{uint8(45), []string{"56"}, false, "uint8"},
		{uint16(234), []string{"56"}, true, "uint16"},
		{uint16(45), []string{"56"}, false, "uint16"},
		{uint32(234), []string{"56"}, true, "uint32"},
		{uint32(45), []string{"56"}, false, "uint32"},
		{uint64(234), []string{"56"}, true, "uint64"},
		{uint64(45), []string{"56"}, false, "uint64"},
		{float32(234.76), []string{"56.765"}, true, "float32"},
		{float32(45.657), []string{"56.654"}, false, "float32"},
		{float64(234.564), []string{"56.654"}, true, "float64"},
		{float64(45.54), []string{"56.45"}, false, "float64"},
		{-234, []string{"-56"}, false, "-int"},
		{-45, []string{"-56"}, true, "-int"},
		{56, []string{"56"}, true, "int border"},
		{55, []string{"56"}, false, "int border"},
		{float64(56.654), []string{"56.654"}, true, "float64 border"},
		{float64(56.653), []string{"56.654"}, false, "float64 border"},
	}

	v := MinValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v min(%v) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestMinValidatorPanicsWhenInputNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "min validator can only validate numbers not, string"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := MinValidator{}
	v.Validate("mango", []string{"3"})
}

func TestMinValidatorPanicsWhenParameterNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-numerical parameter used in MinValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := MinValidator{}
	v.Validate(345, []string{"mango"})
}

func TestMinValidatorPanicsWhenParameterMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameter for MinValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := MinValidator{}
	v.Validate(345, []string{})
}

//max
func TestMaxValidatorType(t *testing.T) {
	want := "max"

	v := MaxValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestMaxValidatorFailureMessage(t *testing.T) {
	want := "Must not be greater than the maximum permitted."

	v := MaxValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestMaxValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{234, []string{"56"}, false, "int"},
		{45, []string{"56"}, true, "int"},
		{int8(127), []string{"56"}, false, "int8"},
		{int8(45), []string{"56"}, true, "int8"},
		{int16(234), []string{"56"}, false, "int16"},
		{int16(45), []string{"56"}, true, "int16"},
		{int32(234), []string{"56"}, false, "int32"},
		{int32(45), []string{"56"}, true, "int32"},
		{int64(234), []string{"56"}, false, "int64"},
		{int64(45), []string{"56"}, true, "int64"},
		{uint(234), []string{"56"}, false, "uint"},
		{uint(45), []string{"56"}, true, "uint"},
		{uint8(234), []string{"56"}, false, "uint8"},
		{uint8(45), []string{"56"}, true, "uint8"},
		{uint16(234), []string{"56"}, false, "uint16"},
		{uint16(45), []string{"56"}, true, "uint16"},
		{uint32(234), []string{"56"}, false, "uint32"},
		{uint32(45), []string{"56"}, true, "uint32"},
		{uint64(234), []string{"56"}, false, "uint64"},
		{uint64(45), []string{"56"}, true, "uint64"},
		{float32(234.76), []string{"56.765"}, false, "float32"},
		{float32(45.657), []string{"56.654"}, true, "float32"},
		{float64(234.564), []string{"56.654"}, false, "float64"},
		{float64(45.54), []string{"56.45"}, true, "float64"},
		{-234, []string{"-56"}, true, "-int"},
		{-45, []string{"-56"}, false, "-int"},
		{56, []string{"56"}, true, "int border"},
		{57, []string{"56"}, false, "int border"},
		{float64(56.654), []string{"56.654"}, true, "float64 border"},
		{float64(56.655), []string{"56.654"}, false, "float64 border"},
	}

	v := MaxValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v max(%v) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestMaxValidatorPanicsWhenInputNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "max validator can only validate numbers not, string"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := MaxValidator{}
	v.Validate("mango", []string{"3"})
}

func TestMaxValidatorPanicsWhenParameterNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-numerical parameter used in MaxValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := MaxValidator{}
	v.Validate(345, []string{"mango"})
}

func TestMaxValidatorPanicsWhenParameterMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameter for MaxValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := MaxValidator{}
	v.Validate(345, []string{})
}

// range
func TestRangeValidatorType(t *testing.T) {
	want := "range"

	v := RangeValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestRangeValidatorFailureMessage(t *testing.T) {
	want := "Must be within the permitted range."

	v := RangeValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestRangeValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{234, []string{"56", "350"}, true, "int"},
		{45, []string{"56", "350"}, false, "int"},
		{int8(127), []string{"56", "350"}, true, "int8"},
		{int8(45), []string{"56", "350"}, false, "int8"},
		{int16(234), []string{"56", "350"}, true, "int16"},
		{int16(45), []string{"56", "350"}, false, "int16"},
		{int32(234), []string{"56", "350"}, true, "int32"},
		{int32(45), []string{"56", "350"}, false, "int32"},
		{int64(234), []string{"56", "350"}, true, "int64"},
		{int64(45), []string{"56", "350"}, false, "int64"},
		{uint(234), []string{"56", "350"}, true, "uint"},
		{uint(45), []string{"56", "350"}, false, "uint"},
		{uint8(234), []string{"56", "350"}, true, "uint8"},
		{uint8(45), []string{"56", "350"}, false, "uint8"},
		{uint16(234), []string{"56", "350"}, true, "uint16"},
		{uint16(45), []string{"56", "350"}, false, "uint16"},
		{uint32(234), []string{"56", "350"}, true, "uint32"},
		{uint32(45), []string{"56", "350"}, false, "uint32"},
		{uint64(234), []string{"56", "350"}, true, "uint64"},
		{uint64(45), []string{"56", "350"}, false, "uint64"},
		{float32(234.76), []string{"56.765", "341.456"}, true, "float32"},
		{float32(45.657), []string{"56.654", "341.456"}, false, "float32"},
		{float64(234.564), []string{"56.654", "341.456"}, true, "float64"},
		{float64(45.54), []string{"56.45", "341.456"}, false, "float64"},
		{-234, []string{"-56", "-35"}, false, "-int"},
		{-45, []string{"-56", "-35"}, true, "-int"},
		{56, []string{"56", "350"}, true, "int lower limit"},
		{350, []string{"56", "350"}, true, "int upper limit"},
		{float64(56.654), []string{"56.654", "341.456"}, true, "float64 lower limit"},
		{float64(341.456), []string{"56.654", "341.456"}, true, "float64 upper limit"},
	}

	v := RangeValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v range(%v,%v) = %v, want %v", test.comment, test.input, test.args[0], test.args[1], got, test.want)
		}
	}
}

func TestRangeValidatorPanicsWhenInputNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "range validator can only validate numbers not, string"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate("mango", []string{"3"})
}

func TestRangeValidatorPanicsWhenSingleParameterNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-numerical parameters used in RangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate(345, []string{"mango", "999"})
}

func TestRangeValidatorPanicsWhenBothParametersNotNumbers(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-numerical parameters used in RangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate(345, []string{"mango", "cheese"})
}

func TestRangeValidatorPanicsWhenValIsUintAndSingleParameterNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-numerical parameters used in RangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate(uint(345), []string{"mango", "999"})
}

func TestRangeValidatorPanicsWhenValIsFloatAndSingleParameterNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-numerical parameters used in RangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate(345.756, []string{"mango", "999"})
}

func TestRangeValidatorPanicsWhenOneParameterMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameters for RangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate(345, []string{"2"})
}

func TestRangeValidatorPanicsWhenBothParametersMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameters for RangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := RangeValidator{}
	v.Validate(345, []string{})
}

//lenmin
func TestLenMinValidatorType(t *testing.T) {
	want := "lenmin"

	v := LenMinValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestLenMinValidatorFailureMessage(t *testing.T) {
	want := "Must not contain fewer elements than minimum permitted."

	v := LenMinValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestLenMinValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{"abcdefg", []string{"3"}, true, "string"},
		{"abcdefg", []string{"9"}, false, "string"},
		{"abcdefg", []string{"7"}, true, "string limit"},
		{[5]int{1, 2, 3, 4, 5}, []string{"3"}, true, "array"},
		{[5]int{1, 2, 3, 4, 5}, []string{"9"}, false, "array"},
		{[5]int{1, 2, 3, 4, 5}, []string{"5"}, true, "array limit"},
		{[]int{1, 2, 3, 4, 5}, []string{"3"}, true, "slice"},
		{[]int{1, 2, 3, 4, 5}, []string{"9"}, false, "slice"},
		{[]int{1, 2, 3, 4, 5}, []string{"5"}, true, "slice limit"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"2"}, true, "map"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"9"}, false, "map"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"3"}, true, "map limit"},
	}

	v := LenMinValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v lenmin(%v) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestLenMinValidatorPanicsWhenInputNotStringArraySliceOrMap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "lenmin validator can only validate strings, arrays, slices and maps, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenMinValidator{}
	v.Validate(76554, []string{"3"})
}

func TestLenMinValidatorPanicsWhenParameterNotInteger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-integer parameter used in LenMinValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenMinValidator{}
	v.Validate(345, []string{"22.4"})
}

func TestLenMinValidatorPanicsWhenParameterMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameter for LenMinValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenMinValidator{}
	v.Validate(345, []string{})
}

//lenmax
func TestLenMaxValidatorType(t *testing.T) {
	want := "lenmax"

	v := LenMaxValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestLenMaxValidatorFailureMessage(t *testing.T) {
	want := "Must not contain more elements than the maximum permitted."

	v := LenMaxValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestLenMaxValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{"abcdefg", []string{"3"}, false, "string"},
		{"abcdefg", []string{"9"}, true, "string"},
		{"abcdefg", []string{"7"}, true, "string limit"},
		{[5]int{1, 2, 3, 4, 5}, []string{"3"}, false, "array"},
		{[5]int{1, 2, 3, 4, 5}, []string{"9"}, true, "array"},
		{[5]int{1, 2, 3, 4, 5}, []string{"5"}, true, "array limit"},
		{[]int{1, 2, 3, 4, 5}, []string{"3"}, false, "slice"},
		{[]int{1, 2, 3, 4, 5}, []string{"9"}, true, "slice"},
		{[]int{1, 2, 3, 4, 5}, []string{"5"}, true, "slice limit"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"2"}, false, "map"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"9"}, true, "map"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"3"}, true, "map limit"},
	}

	v := LenMaxValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v lenmax(%v) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestLenMaxValidatorPanicsWhenInputNotStringArraySliceOrMap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "lenmax validator can only validate strings, arrays, slices and maps, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenMaxValidator{}
	v.Validate(76554, []string{"3"})
}

func TestLenMaxValidatorPanicsWhenParameterNotInteger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-integer parameter used in LenMaxValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenMaxValidator{}
	v.Validate(345, []string{"34.6"})
}

func TestLenMaxValidatorPanicsWhenParameterMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameter for LenMaxValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenMaxValidator{}
	v.Validate(345, []string{})
}

//lenrange
func TestLenRangeValidatorType(t *testing.T) {
	want := "lenrange"

	v := LenRangeValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestLenRangeValidatorFailureMessage(t *testing.T) {
	want := "Must have a quantity of elements within the permitted range."

	v := LenRangeValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestLenRangeValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{"abcdefg", []string{"3", "9"}, true, "string"},
		{"ab", []string{"3", "9"}, false, "string too short"},
		{"abcdefgoiu", []string{"3", "9"}, false, "string too long"},
		{"abcdefg", []string{"3", "7"}, true, "string upper limit"},
		{"abc", []string{"3", "9"}, true, "string lower limit"},
		{[5]int{1, 2, 3, 4, 5}, []string{"3", "9"}, true, "array"},
		{[5]int{1, 2, 3, 4, 5}, []string{"6", "9"}, false, "array too short"},
		{[5]int{1, 2, 3, 4, 5}, []string{"3", "4"}, false, "array too long"},
		{[5]int{1, 2, 3, 4, 5}, []string{"3", "5"}, true, "array upper limit"},
		{[5]int{1, 2, 3, 4, 5}, []string{"5", "9"}, true, "array lower limit"},
		{[]int{1, 2, 3, 4, 5}, []string{"3", "9"}, true, "slice"},
		{[]int{1, 2, 3, 4, 5}, []string{"6", "9"}, false, "slice too short"},
		{[]int{1, 2, 3, 4, 5}, []string{"3", "4"}, false, "slice too long"},
		{[]int{1, 2, 3, 4, 5}, []string{"3", "5"}, true, "slice upper limit"},
		{[]int{1, 2, 3, 4, 5}, []string{"5", "9"}, true, "slice lower limit"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"2", "9"}, true, "map"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"5", "9"}, false, "map too short"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"1", "2"}, false, "map too long"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"1", "3"}, true, "map upper limit"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"3", "9"}, true, "map lower limit"},
	}

	v := LenRangeValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v lenrange(%v) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestLenRangeValidatorPanicsWhenInputNotStringArraySliceOrMap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "lenrange validator can only validate strings, arrays, slices and maps, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenRangeValidator{}
	v.Validate(76554, []string{"3", "3000"})
}

func TestLenRangeValidatorPanicsWhenOneParameterNotInteger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-integer parameters used in LenRangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenRangeValidator{}
	v.Validate(345, []string{"34.6", "22"})
}

func TestLenRangeValidatorPanicsWhenBothParametersNotInteger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-integer parameters used in LenRangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenRangeValidator{}
	v.Validate(345, []string{"34.6", "11.45"})
}

func TestLenRangeValidatorPanicsWhenOneParameterMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameters for LenRangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenRangeValidator{}
	v.Validate(345, []string{"34"})
}

func TestLenRangeValidatorPanicsWhenBothParametersMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "missing parameters for LenRangeValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := LenRangeValidator{}
	v.Validate(345, []string{})
}

//contins
func TestContainValidatorType(t *testing.T) {
	want := "contains"

	v := ContainsValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestContainsValidatorFailureMessage(t *testing.T) {
	want := "Must contain a specific string."

	v := ContainsValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}

func TestContainsValidator(t *testing.T) {
	a := "1"
	b := "2"
	c := "3"
	d := "4"
	e := "5"

	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{"CHEESEBICYCLE", []string{"Bicycle"}, false, "Uppercase"},
		{"cheesebicycle", []string{"Bicycle"}, false, "Lowercase"},
		{"CheeseBicycle", []string{"Bicycle"}, true, "End"},
		{"CheeseBicycle", []string{"Bicycles"}, false, "ExtraSuffix"},
		{"CheeseBicycle", []string{"Cheese"}, true, "Start"},
		{"Cheese81cycl3", []string{"81cycl3"}, true, "Digits"},
		{"Cheese Bicycle", []string{"se Bi"}, true, "Mid Spaces"},
		{"CheeseBicycle", []string{"MyCheese"}, false, "ExtraPrefix"},
		{"seBi", []string{"CheeseBicycle"}, false, "Reverse Subset"},
		{"CheeseBicycle", []string{""}, true, "EmptyTest"},
		{"CheeseBicycle", []string{}, true, "MissingTest"},
		{"", []string{"MyCheese"}, false, "Empty"},
		{"", []string{}, true, "BothEmpty"},
		{[0]string{}, []string{}, false, "emptyarray"},
		{[0]string{}, []string{"3"}, false, "emptyarray"},
		{[5]string{"1", "2", "3", "4", "5"}, []string{"3"}, true, "array"},
		{[5]string{"1", "2", "3", "4", "5"}, []string{"6"}, false, "array"},
		{[5]string{"123", "234", "345", "456", "567"}, []string{"4"}, false, "whole string array"},
		{[5]string{"1", "2", "3", "4", "5"}, []string{"234"}, false, "whole string array"},
		{[]string{}, []string{}, false, "emptyslice"},
		{[]string{}, []string{"3"}, false, "emptyslice"},
		{[]string{"1", "2", "3", "4", "5"}, []string{"3"}, true, "slice"},
		{[]string{"1", "2", "3", "4", "5"}, []string{"6"}, false, "slice"},
		{[]string{"123", "234", "345", "456", "567"}, []string{"4"}, false, "whole string slice"},
		{[]string{"1", "2", "3", "4", "5"}, []string{"234"}, false, "whole string slice"},
		{[5]*string{&a, &b, &c, &d, &e}, []string{"3"}, true, "pointer array"},
		{[]*string{&a, &b, &c, &d, &e}, []string{"3"}, true, "pointer slice"},
		{map[string]int{}, []string{}, false, "emptymap"},
		{map[string]int{"a": 1, "b": 2, "c": 3}, []string{"b"}, true, "map"},
		{map[string]int{"a": 1, "b": 2, "c": 3}, []string{"e"}, false, "map"},
		{map[string]int{"abc": 1, "bcd": 2, "cde": 3}, []string{"b"}, false, "whole string map"},
		{map[string]int{"a": 1, "b": 2, "c": 3}, []string{"abc"}, false, "whole string map"},
		{map[*string]int{&a: 1, &b: 2, &c: 3}, []string{"3"}, true, "pointer map"},
	}

	v := ContainsValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %q contains(%s) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestContainsValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "contains validator can only validate strings, arrays, slices and maps, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := ContainsValidator{}
	v.Validate(32, []string{})
}

func TestContainsValidatorPanicsWhenInputNotSliceString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "contains validator can only validate arrays and slices of string, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := ContainsValidator{}
	v.Validate([]int{1, 2, 3, 4, 5}, []string{"3"})
}

func TestContainsValidatorPanicsWhenInputNotMapStringKey(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "contains validator can only validate maps with keys of string, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := ContainsValidator{}
	v.Validate(map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}, []string{"3"})
}

//inset
func TestInSetValidatorType(t *testing.T) {
	want := "inset"

	v := InSetValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}

func TestInSetValidatorFailureMessage(t *testing.T) {
	want := "Must be in the permitted set."

	v := InSetValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}

func TestInSetValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{"CHEESE", []string{"Cheese", "Bicycle", "Mango"}, false, "Uppercase"},
		{"cheese", []string{"Cheese", "Bicycle", "Mango"}, false, "Lowercase"},
		{"Cheese", []string{"Cheese", "Bicycle", "Mango"}, true, "Match"},
		{"Bicycle", []string{"Cheese", "Bicycle", "Mango"}, true, "Middle"},
		{"Bicycle", []string{"Cheese", " Bicycle ", "Mango"}, true, "Whitespace"},
		{"", []string{"Cheese", "", "Mango"}, true, "Empty"},
		{"", []string{"Cheese", "Mango"}, false, "Missing"},
		{2, []string{"1", "2", "3"}, true, "Digits"},
		{5, []string{"1", "2", "3"}, false, "Digits"},
		{2, []string{"1", " 2 ", "3"}, true, "Whitespace"},
	}

	v := InSetValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %q inset(%s) = %v, want %v", test.comment, test.input, test.args[0], got, test.want)
		}
	}
}

func TestInSetValidatorPanicsWhenInputNotString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "inset validator can only validate strings and ints, not float64"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := InSetValidator{}
	v.Validate(32.6, []string{"45", "67"})
}

func TestInSetValidatorPanicsWhenInputIsIntAndParamsAreNotInts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "non-integer parameter used in InSetValidator"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := InSetValidator{}
	v.Validate(32, []string{"45", "67", "45.7"})
}

//notempty
func TestNotEmptyValidatorType(t *testing.T) {
	want := "notempty"

	v := NotEmptyValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestNotEmptyValidatorFailureMessage(t *testing.T) {
	want := "Must not be empty."

	v := NotEmptyValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestNotEmptyValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{"", []string{}, false, "string"},
		{"abcdefg", []string{}, true, "string"},
		{[0]int{}, []string{}, false, "array"},
		{[5]int{1, 2, 3, 4, 5}, []string{}, true, "array"},
		{[]int{}, []string{}, false, "slice"},
		{[]int{1, 2, 3, 4, 5}, []string{}, true, "slice"},
		{map[int]string{}, []string{}, false, "map"},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{}, true, "map"},
	}

	v := NotEmptyValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v notempty = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestNotEmptyValidatorPanicsWhenInputNotStringArraySliceOrMap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "notempty validator can only validate strings, arrays, slices and maps, not int"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := NotEmptyValidator{}
	v.Validate(76554, []string{})
}

//notzero
func TestNotZeroValidatorType(t *testing.T) {
	want := "notzero"

	v := NotZeroValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestNotZeroValidatorFailureMessage(t *testing.T) {
	want := "Must not be zero."

	v := NotZeroValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestNotZeroValidator(t *testing.T) {
	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{234, []string{}, true, "int"},
		{0, []string{}, false, "int"},
		{int8(127), []string{}, true, "int8"},
		{int8(0), []string{}, false, "int8"},
		{int16(234), []string{}, true, "int16"},
		{int16(0), []string{}, false, "int16"},
		{int32(234), []string{}, true, "int32"},
		{int32(0), []string{}, false, "int32"},
		{int64(234), []string{}, true, "int64"},
		{int64(0), []string{}, false, "int64"},
		{uint(234), []string{}, true, "uint"},
		{uint(0), []string{}, false, "uint"},
		{uint8(234), []string{}, true, "uint8"},
		{uint8(0), []string{}, false, "uint8"},
		{uint16(234), []string{}, true, "uint16"},
		{uint16(0), []string{}, false, "uint16"},
		{uint32(234), []string{}, true, "uint32"},
		{uint32(0), []string{}, false, "uint32"},
		{uint64(234), []string{}, true, "uint64"},
		{uint64(0), []string{"56"}, false, "uint64"},
		{float32(234.76), []string{}, true, "float32"},
		{float32(0), []string{}, false, "float32"},
		{float64(234.564), []string{}, true, "float64"},
		{float64(0), []string{}, false, "float64"},
		{-45, []string{}, true, "-int"},
		{float64(-56.654), []string{}, true, "-float64"},
		{1, []string{}, true, "int border"},
		{float64(0.0000000001), []string{}, true, "float64 border"},
	}

	v := NotZeroValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v notzero = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestNotZeroValidatorPanicsWhenInputNotNumber(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "notzero validator can only validate numbers not, string"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := NotZeroValidator{}
	v.Validate("mango", []string{"3"})
}

//notnil
func TestNotNilValidatorType(t *testing.T) {
	want := "notnil"

	v := NotNilValidator{}
	got := v.Type()

	if got != want {
		t.Errorf("Valid = %q, want %q", got, want)
	}
}
func TestNotNilValidatorFailureMessage(t *testing.T) {
	want := "Must not be nil."

	v := NotNilValidator{}
	got := v.FailureMsg()

	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}
}
func TestNotNilValidator(t *testing.T) {
	str := "asdd"
	strP := &str
	num := 345
	numP := &num
	strtString := struct {
		Name string
		P    *string
	}{
		"test", nil,
	}

	strStrut := struct {
		Name       string
		Dimensions *struct {
			Width  float32
			Length float32
		}
	}{
		"Mango", nil,
	}

	var unMap map[string]int
	inMap := make(map[string]int)
	var unSl []int
	inSl := []int{}

	var tests = []struct {
		input   interface{}
		args    []string
		want    bool
		comment string
	}{
		{strP, []string{}, true, "string ptr"},
		{numP, []string{}, true, "number ptr"},
		{&strtString, []string{}, true, "struct ptr"},
		{strtString.P, []string{}, false, "struct nil string ptr prop"},
		{strStrut.Dimensions, []string{}, false, "struct nil struct ptr prop nil ptr"},
		{unMap, []string{}, false, "uninitialized map"},
		{inMap, []string{}, true, "initialized map"},
		{unSl, []string{}, false, "uninitialized slice"},
		{inSl, []string{}, true, "initialized slice"},
	}

	v := NotNilValidator{}

	for _, test := range tests {
		if got := v.Validate(test.input, test.args); got != test.want {
			t.Errorf("Validate (%s): %v notnil = %v, want %v", test.comment, test.input, got, test.want)
		}
	}
}

func TestNotNilValidatorPanicsWhenInputNotPointer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "notnil validator can only validate maps, slices and pointers, not string"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	v := NotNilValidator{}
	v.Validate("sdfsdf", []string{})
}

// ********************
//
//  End of validators
//
// ********************

func TestNewParameterValidatorsReturnsWithInitialisedValidatorsMap(t *testing.T) {
	want := len(getDefaultValidators())

	pv := newValidationHandler()

	cpv := pv.(*elementValidationHandler)
	got := len(cpv.validators)

	if got != want {
		t.Errorf("Validator count = %d, want %d", got, want)
	}
}

func TestNewParameterValidatorsHasDefaultValidators(t *testing.T) {
	pv := newValidationHandler()
	cpv := pv.(*elementValidationHandler)
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

	pv := elementValidationHandler{}
	pv.validators = make(map[string]Validator)
	v := []Validator{EmptyValidator{}, testValidator1{}, testValidator2{}}
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

	pv := newValidationHandler()
	pv.IsValid("validator1", "test1")
}

func TestIsValidDoesNotPanicWhenUnknownConstraintIsIgnoreContents(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()

	pv := newValidationHandler()
	pv.IsValid("validator1", "ignorecontents")
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

	pv := newValidationHandler()
	pv.IsValid("validator1", "test1")
}

func TestIsValidDoesNotPanicIWhenKnownConstraint(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()

	pv := newValidationHandler()
	pv.AddValidator(testValidator1{})
	pv.IsValid("validator1", "test1")
}

func TestIsValidPassesParamValueToValidator(t *testing.T) {
	want := true

	pv := newValidationHandler()
	pv.AddValidator(testValidator1{})

	_, got := pv.IsValid("validator1", "test1")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	_, got = pv.IsValid("validator2", "test1")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidParsesConstraintArgsWhenSingleArg(t *testing.T) {
	want := true

	pv := newValidationHandler()
	pv.AddValidator(testValidator2{})

	_, got := pv.IsValid("paramValue", "test2(arg)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	_, got = pv.IsValid("paramValue", "test2(6)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidParsesConstraintArgsWhenMultipleArgs(t *testing.T) {
	want := true

	pv := newValidationHandler()
	pv.AddValidator(testValidator3{})

	_, got := pv.IsValid("paramValue", "test3(6,arg2)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	_, got = pv.IsValid("paramValue", "test3(6,arg)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestIsValidTrimsSpaceAroundConstraintArgs(t *testing.T) {
	want := true

	pv := newValidationHandler()
	pv.AddValidator(testValidator3{})

	_, got := pv.IsValid("paramValue", "test3(  6 ,  arg2 )")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	_, got = pv.IsValid("paramValue", "test3( 6 , arg )")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestAddValidatorPanicsWhenConstraintTypeExists(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pv := newValidationHandler()
	pv.AddValidator(testValidator1{})
	pv.AddValidator(testValidator1DuplicateTypeCode{})
}

func TestAddValidatorsPanicsWhenConstraintTypeExists(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pv := newValidationHandler()
	pv.AddValidators([]Validator{
		testValidator1{},
		testValidator1DuplicateTypeCode{},
	})
}

func TestAddValidatorPanicsWithCorrectErrorMessageWhenConstraintTypeExists(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: test1"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	pv := newValidationHandler()
	pv.AddValidator(testValidator1{})
	pv.AddValidator(testValidator1DuplicateTypeCode{})
}

func TestAddValidatorsPanicsWithCorrectErrorMessageWhenConstraintTypeExists(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: test1"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	pv := newValidationHandler()
	pv.AddValidators([]Validator{
		testValidator1{},
		testValidator1DuplicateTypeCode{},
	})
}

func TestIsValidPanicsWhenConstraintMalformed(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pv := newValidationHandler()
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

	pv := newValidationHandler()
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

	pv := newValidationHandler()
	pv.AddValidator(testValidator3{})

	pv.IsValid("paramValue", "test3(arg2")
}

func TestIsValidTrimsSpaceAroundConstraintName(t *testing.T) {
	want := true

	pv := newValidationHandler()
	pv.AddValidator(testValidator3{})

	_, got := pv.IsValid("paramValue", "  test3  (6,arg2)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}

	want = false
	_, got = pv.IsValid("paramValue", "  test3  (6,arg)")

	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestParseConstraints(t *testing.T) {

	var tests = []struct {
		input     string
		wantCount int
		results   map[string]string
	}{
		{"alpha", 1, map[string]string{"alpha": "[]"}},
		{"min(1)", 1, map[string]string{"min": "[1]"}},
		{"range(12.3,45.6)", 1, map[string]string{"range": "[12.3 45.6]"}},
		{" alpha ", 1, map[string]string{"alpha": "[]"}},
		{"min ( 1)", 1, map[string]string{"min": "[1]"}},
		{"range(12.3 ,45.6 )", 1, map[string]string{"range": "[12.3 45.6]"}},
		{"alpha,", 1, map[string]string{"alpha": "[]"}},
		{"min(1),", 1, map[string]string{"min": "[1]"}},
		{"range(12.3,45.6),", 1, map[string]string{"range": "[12.3 45.6]"}},
		{"alpha,min(1)", 2, map[string]string{"alpha": "[]", "min": "[1]"}},
		{"alpha,min(1),range(23.4, 6754)", 3, map[string]string{"alpha": "[]", "min": "[1]", "range": "[23.4 6754]"}},
		{"min(1),range(23.4, 6754),alpha", 3, map[string]string{"alpha": "[]", "min": "[1]", "range": "[23.4 6754]"}},
		{"alpha,min(1),", 2, map[string]string{"alpha": "[]", "min": "[1]"}},
		{"alpha,min(1),range(23.4, 6754),", 3, map[string]string{"alpha": "[]", "min": "[1]", "range": "[23.4 6754]"}},
	}

	pv := newValidationHandler()
	for _, test := range tests {

		parsed := pv.ParseConstraints(test.input)
		gotCount := len(parsed)
		if gotCount != test.wantCount {
			t.Errorf("Constraint count = %d, want %d", gotCount, test.wantCount)
			return
		}
		for name, wantArgs := range test.results {
			args, ok := parsed[name]
			if !ok {
				t.Errorf("Parsed result missing (%s), want %q", test.input, name)
			}
			got := fmt.Sprintf("%s", args)
			if got != wantArgs {
				t.Errorf("Parsed (%s): %q, want %q", test.input, got, wantArgs)
			}
		}
	}

}

func TestParseConstraintsPanicsWhenConstraintHasDoubleCommas(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "illegal constraint format: alpha,,min(6)"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	pv := newValidationHandler()
	pv.ParseConstraints("alpha,,min(6)")
}

func getParsedConstraintsKeys(m map[string][]string) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

type testValidator1 struct{}

func (testValidator1) Validate(i interface{}, args []string) bool {
	return i.(string) == "validator1"
}

func (testValidator1) Type() string {
	return "test1"
}
func (testValidator1) FailureMsg() string {
	return "must test1"
}

type testValidator2 struct{}

func (testValidator2) Validate(i interface{}, args []string) bool {
	return len(args) == 1 && args[0] == "arg"
}

func (testValidator2) Type() string {
	return "test2"
}
func (testValidator2) FailureMsg() string {
	return "must test2"
}

type testValidator3 struct{}

func (testValidator3) Validate(i interface{}, args []string) bool {
	return len(args) == 2 && args[1] == "arg2"
}

func (testValidator3) Type() string {
	return "test3"
}
func (testValidator3) FailureMsg() string {
	return "must test3"
}

type testValidator1DuplicateTypeCode struct{}

func (testValidator1DuplicateTypeCode) Validate(i interface{}, args []string) bool {
	return i.(string) == "testValidator1DuplicateTypeCode"
}

func (testValidator1DuplicateTypeCode) Type() string {
	return "test1"
}
func (testValidator1DuplicateTypeCode) FailureMsg() string {
	return ""
}

type dupValidator struct{}

func (dupValidator) Validate(i interface{}, args []string) bool {
	return true
}

func (dupValidator) Type() string {
	return "int32"
}
func (dupValidator) FailureMsg() string {
	return ""
}

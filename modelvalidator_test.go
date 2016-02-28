package mango

import (
	"reflect"
	"testing"
)

func TestValidateModelByValueDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha"`
		LastName  string `validate:"alpha"`
	}{
		"Jeff", "Mango",
	}
	validator.Validate(test)
}

func TestValidateModelByRefDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha"`
		LastName  string `validate:"alpha"`
	}{
		"Jeff", "Mango",
	}
	validator.Validate(&test)
}

func TestValidateStringsModelWithNoErrorsReturnsOK(t *testing.T) {
	want := true
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha"`
		LastName  string `validate:"alpha"`
	}{
		"Jeff", "Mango",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateStringsModelWithErrorReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha"`
		LastName  string `validate:"alpha"`
	}{
		"564", "9876",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateStringsModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha"`
		LastName  string `validate:"alpha"`
	}{
		"564", "9876",
	}

	details, _ := validator.Validate(&test)

	want := "alpha"
	fails, ok := details["FirstName"]
	if !ok {
		t.Errorf("FirstName Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}

	fails, ok = details["LastName"]
	if !ok {
		t.Errorf("LastName Validate fail count = 0, want 1")
	}
	got = fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateStringsModelWithMultipleErrorsReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha,prefix(sheep)"`
		LastName  string `validate:"alpha,suffix(cheese)"`
	}{
		"5kjh64", "jh9876",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateStringsModelWithMultipleErrorsReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string `validate:"alpha,prefix(sheep)"`
		LastName  string `validate:"alpha,suffix(cheese)"`
	}{
		"56kj4", "987kjh6",
	}

	details, _ := validator.Validate(&test)

	want := "alpha,prefix(sheep)"
	fails, ok := details["FirstName"]
	if !ok {
		t.Errorf("FirstName Validate fail count = 0, want 2")
	}
	got := fails[0].Code + "," + fails[1].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
	want = "alpha,suffix(cheese)"
	fails, ok = details["LastName"]
	if !ok {
		t.Errorf("LastName Validate fail count = 0, want 1")
	}
	got = fails[0].Code + "," + fails[1].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateIntegerModelWithErrorReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name string `validate:"alpha"`
		Age  int    `validate:"max(67)"`
	}{
		"Mango", 9876,
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateIntegerModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name string `validate:"alpha"`
		Age  int    `validate:"max(67)"`
	}{
		"Mango", 9876,
	}

	details, _ := validator.Validate(&test)

	want := "max(67)"

	fails, ok := details["Age"]
	if !ok {
		t.Errorf("Age Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateUint8ModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name string `validate:"alpha"`
		Age  uint8  `validate:"min(67)"`
	}{
		"Mango", 56,
	}

	details, ok := validator.Validate(&test)
	if ok {
		t.Errorf("Validate result = true, want false")
		return
	}

	want := "min(67)"
	fails, ok := details["Age"]
	if !ok {
		t.Errorf("Age Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateFloatModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string  `validate:"alpha"`
		Amount float32 `validate:"min(67.786)"`
	}{
		"Mango", 56.3379,
	}

	details, _ := validator.Validate(test)

	want := "min(67.786)"
	fails, ok := details["Amount"]
	if !ok {
		t.Errorf("Amount Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateModelInnerStructWithNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name       string `validate:"alpha"`
		Dimensions struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			124.6, 61.8,
		},
	}

	_, ok := validator.Validate(&test)

	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelInnerStructWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name       string `validate:"alpha"`
		Dimensions struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			24.6, 53.8,
		},
	}

	details, _ := validator.Validate(&test)

	want := "min(50)"
	fails, ok := details["Dimensions.Width"]
	if !ok {
		t.Errorf("Validate error Dimensions.Width details missing")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateModelInnerPointerStructWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name       string `validate:"alpha"`
		Dimensions *struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", &struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			24.6, 53.8,
		},
	}

	details, _ := validator.Validate(&test)

	want := "min(50)"
	fails, ok := details["Dimensions.Width"]
	if !ok {
		t.Errorf("Validate error Dimensions.Width details missing")
		return
	}

	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateStringModelWithoutConstraints(t *testing.T) {
	lastname := "9876"
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string
		LastName  *string
	}{
		"Mango", &lastname,
	}

	_, ok := validator.Validate(&test)

	want := true
	got := ok

	if got != want {
		t.Errorf("Validate errors = %t, want %t", got, want)
	}
}

func TestValidateStringPointerModelWithErrorReturnsDetailsInMap(t *testing.T) {
	lastname := "9876"
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string  `validate:"alpha"`
		LastName  *string `validate:"alpha"`
	}{
		"564", &lastname,
	}

	details, _ := validator.Validate(&test)

	want := "alpha"
	fails, ok := details["FirstName"]
	if !ok {
		t.Errorf("FirstName Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}

	fails, ok = details["LastName"]
	if !ok {
		t.Errorf("LastName Validate fail count = 0, want 1")
	}
	got = fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateStringNilPointerModelDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string  `validate:"alpha"`
		LastName  *string `validate:"alpha"`
	}{
		FirstName: "564", LastName: nil,
	}

	validator.Validate(&test)
}

func TestValidateModelInnerPointerNilStructDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name       string `validate:"alpha"`
		Dimensions *struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		Name: "Mango", Dimensions: nil,
	}

	validator.Validate(&test)
}

func TestValidateMapModelWithErrorReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string         `validate:"alpha"`
		Items map[string]int `validate:"notnil"`
	}{
		Name: "Mango",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateMapModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string         `validate:"alpha"`
		Items map[string]int `validate:"notnil"`
	}{
		Name: "Mango",
	}

	details, ok := validator.Validate(&test)
	if ok {
		t.Errorf("Validate result = true, want false")
		return
	}

	want := "notnil"
	fails, ok := details["Items"]
	if !ok {
		t.Errorf("Items Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateSliceModelWithErrorReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string   `validate:"alpha"`
		Items []string `validate:"notnil"`
	}{
		Name: "Mango",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateSliceModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string   `validate:"alpha"`
		Items []string `validate:"notempty"`
	}{
		Name: "Mango",
	}

	details, ok := validator.Validate(&test)
	if ok {
		t.Errorf("Validate result = true, want false")
		return
	}

	want := "notempty"
	fails, ok := details["Items"]
	if !ok {
		t.Errorf("Items Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateArrayModelWithErrorReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string    `validate:"alpha"`
		Items [5]string `validate:"lenmin(6)"`
	}{
		Name: "Mango",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidateArrayModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string    `validate:"alpha"`
		Items [5]string `validate:"lenmin(6)"`
	}{
		Name: "Mango",
	}

	details, ok := validator.Validate(&test)
	if ok {
		t.Errorf("Validate result = true, want false")
		return
	}

	want := "lenmin(6)"
	fails, ok := details["Items"]
	if !ok {
		t.Errorf("Items Validate fail count = 0, want 1")
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestModelValidatorAddCustomValidator(t *testing.T) {
	type model struct {
		Name string
		Age  int
	}
	custValidator := func(m interface{}) (map[string][]ValidationFailure, bool) {
		results := make(map[string][]ValidationFailure)
		cm := m.(model)
		if cm.Name != "Mango" {
			results["Name"] = []ValidationFailure{
				{Code: "wrongname",
					Message: "Name must be Mango"},
			}
		}
		return results, len(results) == 0
	}

	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
		customValidators:  make(map[reflect.Type]ValidateFunc),
	}
	validator.AddCustomValidator(model{}, custValidator)

	testModel := model{"Mingo", 45}
	_, ok := validator.Validate(testModel)
	want := false
	got := ok
	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
		return
	}
}

func TestCustomValidatorWithPointerModel(t *testing.T) {
	type model struct {
		Name string
		Age  int
	}
	custValidator := func(m interface{}) (map[string][]ValidationFailure, bool) {
		results := make(map[string][]ValidationFailure)
		cm := m.(model)
		if cm.Name != "Mango" {
			results["Name"] = []ValidationFailure{
				{Code: "wrongname",
					Message: "Name must be Mango"},
			}
		}
		return results, len(results) == 0
	}

	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
		customValidators:  make(map[reflect.Type]ValidateFunc),
	}
	validator.AddCustomValidator(model{}, custValidator)

	testModel := model{"Mingo", 45}
	_, ok := validator.Validate(&testModel)
	want := false
	got := ok
	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
		return
	}
}

func TestCustomValidatorreturnsFailureDetails(t *testing.T) {
	type model struct {
		Name string
		Age  int
	}
	custValidator := func(m interface{}) (map[string][]ValidationFailure, bool) {
		results := make(map[string][]ValidationFailure)
		cm := m.(model)
		if cm.Name != "Mango" {
			results["Name"] = []ValidationFailure{
				{Code: "wrongname",
					Message: "Name must be Mango"},
			}
		}
		return results, len(results) == 0
	}

	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
		customValidators:  make(map[reflect.Type]ValidateFunc),
	}
	validator.AddCustomValidator(model{}, custValidator)

	testModel := model{"Mingo", 45}
	details, ok := validator.Validate(&testModel)
	if ok {
		return
	}
	want := "Name must be Mango"
	got := details["Name"][0].Message
	if got != want {
		t.Errorf("Validate message = %q, want %q", got, want)
		return
	}
}

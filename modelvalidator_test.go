package mango

import (
	"reflect"
	"testing"
	"time"
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

func TestValidateModelWithUnexportedFieldDoesNotPanic(t *testing.T) {
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
		nickname  string
	}{
		"Jeff", "Mango", "jango",
	}
	validator.Validate(&test)
}

func TestValidateModelWithLocalTypeWithUnexportedFieldDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()

	type WithUnexportedField struct {
		explode int
	}

	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName   string `validate:"alpha"`
		LastName    string `validate:"alpha"`
		DangerField WithUnexportedField
	}{
		"Jeff", "Mango", WithUnexportedField{},
	}
	validator.Validate(&test)
}

func TestValidateModelWithPackageTypeWithUnexportedFieldDoesNotPanic(t *testing.T) {
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
		DOB       time.Time
	}{
		"Jeff", "Mango", time.Now(),
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

	match := func(fails []ValidationFailure, s string) string {
		for _, v := range fails {
			if v.Code == s {
				return s
			}
		}
		return ""
	}

	details, _ := validator.Validate(&test)

	fails, ok := details["FirstName"]
	if !ok {
		t.Errorf("FirstName Validate fail count = 0, want 2")
	}

	want := "alpha"
	got := match(fails, want)
	if got != want {
		t.Errorf("Validate error (FirstName) = %q, want %q", got, want)
	}
	want = "prefix(sheep)"
	got = match(fails, want)
	if got != want {
		t.Errorf("Validate error (FirstName) = %q, want %q", got, want)
	}

	fails, ok = details["LastName"]
	if !ok {
		t.Errorf("LastName Validate fail count = 0, want 1")
	}
	want = "alpha"
	got = match(fails, want)
	if got != want {
		t.Errorf("Validate error (LastName) = %q, want %q", got, want)
	}
	want = "suffix(cheese)"
	got = match(fails, want)
	if got != want {
		t.Errorf("Validate error (LastName) = %q, want %q", got, want)
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

func TestValidateModelWithSimpleJsonTagsReturnsTagNameInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name string `json:"name" validate:"alpha"`
		Age  int    `json:"age" validate:"max(67)"`
	}{
		"Mango", 9876,
	}

	details, _ := validator.Validate(&test)

	want := "max(67)"

	fails, ok := details["age"]
	if !ok {
		t.Errorf("'age' not in map")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateModelWithComplexJsonTagsReturnsTagNameInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name string `json:"name,omitempty" validate:"alpha"`
		Age  int    `json:"age,omitempty" validate:"max(67)"`
	}{
		"Mango", 9876,
	}

	details, _ := validator.Validate(&test)

	want := "max(67)"

	fails, ok := details["age"]
	if !ok {
		t.Errorf("'age' not in map")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateNestedModelWithComplexJsonTagsReturnsTagNameInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}

	type nested struct {
		Name string `json:"name,omitempty" validate:"alpha"`
		Age  int    `json:"age,omitempty" validate:"max(67)"`
	}

	type outer struct {
		Person nested `json:"person, omitempty"`
	}

	test := outer{
		Person: nested{"Mango", 9876},
	}

	details, _ := validator.Validate(&test)

	want := "max(67)"

	fails, ok := details["person.age"]
	if !ok {
		t.Errorf("'person.age' not in map")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateCollectionModelWithComplexJsonTagsReturnsTagNameInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}

	type nested struct {
		Name string `json:"name,omitempty" validate:"alpha"`
		Age  int    `json:"age,omitempty" validate:"max(67)"`
	}

	type outer struct {
		Persons []nested `json:"persons, omitempty"`
	}

	test := outer{
		Persons: []nested{{"Mango", 9876}},
	}

	details, _ := validator.Validate(&test)

	want := "max(67)"

	fails, ok := details["persons[0].age"]
	if !ok {
		t.Errorf("'age' not in map")
		return
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
		LastName  *string `validate:"*alpha"`
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

func TestValidatePointerModelWithParameterConstraintAndErrorReturnsDetailsInMap(t *testing.T) {
	lastname := "Mango"
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		FirstName string  `validate:"alpha"`
		LastName  *string `validate:"*lenmin(6)"`
	}{
		"Jeff", &lastname,
	}

	details, _ := validator.Validate(&test)

	want := "lenmin(6)"
	fails, ok := details["LastName"]
	if !ok {
		t.Errorf("LastName Validate fail count = 0, want 1")
	}
	got := fails[0].Code

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
		LastName  *string `validate:"*alpha"`
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

func TestValidatePointerModelWithErrorReturnsNotOK(t *testing.T) {
	want := false
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string  `validate:"alpha"`
		Items *string `validate:"notnil"`
	}{
		Name: "Mango",
	}

	_, got := validator.Validate(&test)

	if got != want {
		t.Errorf("Validate result = %t, want %t", got, want)
	}
}

func TestValidatePointerModelWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name  string  `validate:"alpha"`
		Items *string `validate:"notnil"`
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

func TestCustomValidatorReturnsFailureDetails(t *testing.T) {
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

func TestValidateModelWithArrayValidatesElementsWithNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes [2]struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", [2]struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			{84.6, 531.8}, {124.6, 111.8},
		},
	}

	_, ok := validator.Validate(&test)
	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelWithArrayValidatesElementsWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes [2]struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", [2]struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			{84.6, 531.8}, {24.6, 111.8},
		},
	}
	details, _ := validator.Validate(&test)

	want := "min(50)"
	fails, ok := details["Shapes[1].Width"]
	if !ok {
		t.Errorf("Validate error Dimensions.Width details missing")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateModelWithSliceValidatesElementsWithNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes []struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", []struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			{84.6, 531.8}, {124.6, 111.8},
		},
	}

	_, ok := validator.Validate(&test)
	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelWithSliceValidatesElementsWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes []struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}
	}{
		"Mango", []struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			{84.6, 531.8}, {24.6, 111.8},
		},
	}
	details, _ := validator.Validate(&test)

	want := "min(50)"
	fails, ok := details["Shapes[1].Width"]
	if !ok {
		t.Errorf("Validate error Dimensions.Width details missing")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateModelWithStringKeyMapValidatesElementsWithNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	type shape struct {
		Width  float32 `validate:"min(50)"`
		Length float32 `validate:"min(60)"`
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes map[string]shape
	}{
		"Mango", map[string]shape{
			"square":    shape{85.8, 85.8},
			"rectangle": shape{124.6, 111.8},
		},
	}

	_, ok := validator.Validate(&test)
	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelWithStringKeyMapValidatesElementsWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	type shape struct {
		Width  float32 `validate:"min(50)"`
		Length float32 `validate:"min(60)"`
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes map[string]shape
	}{
		"Mango", map[string]shape{
			"square":    shape{85.8, 85.8},
			"rectangle": shape{124.6, 11.8},
		},
	}
	details, _ := validator.Validate(&test)
	want := "min(60)"
	fails, ok := details["Shapes[rectangle].Length"]
	if !ok {
		t.Errorf("Validate error Dimensions.Width details missing")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

func TestValidateModelWithIntKeyMapValidatesElementsWithErrorReturnsDetailsInMap(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	type shape struct {
		Width  float32 `validate:"min(50)"`
		Length float32 `validate:"min(60)"`
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes map[int]shape
	}{
		"Mango", map[int]shape{
			34:      shape{85.8, 85.8},
			9454985: shape{124.6, 11.8},
		},
	}
	details, _ := validator.Validate(&test)
	want := "min(60)"
	fails, ok := details["Shapes[9454985].Length"]
	if !ok {
		t.Errorf("Validate error Dimensions.Width details missing")
		return
	}
	got := fails[0].Code

	if got != want {
		t.Errorf("Validate error = %q, want %q", got, want)
	}
}

// ignorecontents

func TestValidateModelInnerStructWhenIgnoreContentsReturnsNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name       string `validate:"alpha"`
		Dimensions struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		} `validate:"ignorecontents"`
	}{
		"Mango", struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			1.6, 1.8,
		},
	}

	_, ok := validator.Validate(&test)

	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelInnerPointerStructWhenIgnoreContentsReturnsNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name       string `validate:"alpha"`
		Dimensions *struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		} `validate:"ignorecontents"`
	}{
		"Mango", &struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			24.6, 53.8,
		},
	}

	_, ok := validator.Validate(&test)

	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelWithArrayWhenIgnoreContentsReturnsNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes [2]struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		} `validate:"ignorecontents"`
	}{
		"Mango", [2]struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			{4.6, 31.8}, {14.6, 11.8},
		},
	}

	_, ok := validator.Validate(&test)
	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelWithSliceWhenIgnoreContentsReturnsNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	test := struct {
		Name   string `validate:"alpha"`
		Shapes []struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		} `validate:"ignorecontents"`
	}{
		"Mango", []struct {
			Width  float32 `validate:"min(50)"`
			Length float32 `validate:"min(60)"`
		}{
			{4.6, 51.8}, {14.6, 11.8},
		},
	}

	_, ok := validator.Validate(&test)
	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

func TestValidateModelWithMapWhenIgnoreContentsReturnsNoErrors(t *testing.T) {
	validator := contextModelValidator{
		validationHandler: newValidationHandler(),
	}
	type shape struct {
		Width  float32 `validate:"min(50)"`
		Length float32 `validate:"min(60)"`
	}
	test := struct {
		Name   string           `validate:"alpha"`
		Shapes map[string]shape `validate:"ignorecontents"`
	}{
		"Mango", map[string]shape{
			"square":    shape{5.8, 5.8},
			"rectangle": shape{14.6, 11.8},
		},
	}

	_, ok := validator.Validate(&test)
	want := true
	got := ok

	if got != want {
		t.Errorf("Validate ok = %t, want %t", got, want)
	}
}

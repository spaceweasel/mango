package mango

import (
	"fmt"
	"reflect"
	"testing"
)

func testFunc(c *Context)  {}
func testFunc2(c *Context) {}
func testFunc3(c *Context) {}

func TestAddHandlerFuncAddsToTreeRoot(t *testing.T) {
	want := 1
	tree := tree{}
	tree.AddHandlerFunc("/Sleep", "GET", testFunc)
	got := len(tree.Root().children)
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddHandlerFuncSetsNodeLabel(t *testing.T) {
	want := "/Sleep"
	tree := tree{}
	tree.AddHandlerFunc("/Sleep", "GET", testFunc)
	got := tree.Root().children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingTwoHandlersWithCompletelyDifferentRoutesAddsTwoNodesToTreeRoot(t *testing.T) {
	want := 2
	tree := tree{}
	tree.AddHandlerFunc("Sleep", "GET", testFunc)
	tree.AddHandlerFunc("Fish", "GET", testFunc)
	got := len(tree.Root().children)
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerWithRouteStartingWithExistingRouteAddsChildToExistingRoute(t *testing.T) {
	want := 1
	tree := tree{}
	tree.AddHandlerFunc("Sleep", "GET", testFunc)
	tree.AddHandlerFunc("Sleepers", "GET", testFunc)
	got := len(tree.Root().children[0].children)
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerWithRouteStartingWithExistingRouteSplitsNodeWithCorrectLabel(t *testing.T) {
	want := "Sleep"
	tree := tree{}
	tree.AddHandlerFunc("Sleep", "GET", testFunc)
	tree.AddHandlerFunc("Sleepers", "GET", testFunc)
	got := tree.Root().children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingHandlerWithRouteStartingWithExistingRouteAddsChildWithCorrectLabel(t *testing.T) {
	want := "ers"
	tree := tree{}
	tree.AddHandlerFunc("Sleep", "GET", testFunc)
	tree.AddHandlerFunc("Sleepers", "GET", testFunc)
	got := tree.Root().children[0].children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingHandlerWithRouteSubstringOfExistingRouteAddsChildWithCorrectLabel(t *testing.T) {
	want := "ers"
	tree := tree{}
	tree.AddHandlerFunc("Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("Sleep", "GET", testFunc)
	got := tree.Root().children[0].children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingMultipleHandlersWithDifferentRoutesAddsMultipleNodes(t *testing.T) {
	//
	// - /S - ugar
	//      - leep - ers
	//             - ing
	//             - y
	//
	want := 6
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepy", "GET", testFunc)
	tree.AddHandlerFunc("/Sleeping", "GET", testFunc)
	tree.AddHandlerFunc("/Sugar", "GET", testFunc)
	got := tree.GetStats().totalNodes
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerStoresHandlerInNode(t *testing.T) {
	want := "testFunc"
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	kvp := tree.Root().children[0].handlers
	f := kvp["GET"]
	name := extractFnName(f)
	got := name
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestAddingHandlerToSameRouteDoesNotAddNode(t *testing.T) {
	want := 1
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepers", "POST", testFunc2)
	got := tree.GetStats().totalNodes
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerToSameRouteUsesExisitingNode(t *testing.T) {
	want := "testFunc2"
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepers", "POST", testFunc2)
	kvp := tree.Root().children[0].handlers
	f := kvp["POST"]
	name := extractFnName(f)
	got := name
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestAddingHandlerForDuplicateRouteMethodPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc2)
}

func TestAddingHandlerForDuplicateRouteMethodPanicsWithCorrectMessage(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "duplicate route handler method: \"GET /Sleepers\""
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc2)
}

func TestRetrievingHandlersForNonexistentPathReturnsFalse(t *testing.T) {
	want := false
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	_, _, got := tree.HandlerFuncs("/Sugar")
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestRetrievingHandlersForExistentPathReturnsTrue(t *testing.T) {
	want := true
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	_, _, got := tree.HandlerFuncs("/Sleep")
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestRetrievingHandlersForMatchingPathReturnsNonEmptyMap(t *testing.T) {
	want := 1
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	handlers, _, _ := tree.HandlerFuncs("/Sleep")
	got := len(handlers)
	if got != want {
		t.Errorf("Handlers count = %d, want %d", got, want)
	}
}

func TestRetrievingHandlersForMatchingPathReturnsMapWithMulipleHandlers(t *testing.T) {
	want := 2
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "POST", testFunc2)
	handlers, _, _ := tree.HandlerFuncs("/Sleep")
	got := len(handlers)
	if got != want {
		t.Errorf("Handlers count = %d, want %d", got, want)
	}
}

func TestRetrievingCorrectHandlerForMethod(t *testing.T) {
	want := "testFunc2"
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "POST", testFunc2)
	handlers, _, _ := tree.HandlerFuncs("/Sleep")
	h := handlers["POST"]
	name := extractFnName(h)
	got := name
	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
}

func TestSettingRouteParamValidators(t *testing.T) {
	validator := mockRouteParamsValidator{valid: true}

	tree := tree{}
	tree.SetRouteParamValidators(validator)

	want := reflect.TypeOf(validator).Name()
	got := reflect.TypeOf(tree.paramValidator).Name()
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestTreeAddRouteParamValidator(t *testing.T) {
	want := true

	tree := tree{}
	v := &parameterValidators{}
	v.validators = make(map[string]ParamValidator)
	tree.paramValidator = v

	tree.AddRouteParamValidator(Int32Validator{})
	valid := tree.paramValidator.IsValid("123", "int32")
	got := valid
	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

func TestTreeAddRouteParamValidators(t *testing.T) {
	want := true

	tree := tree{}
	v := &parameterValidators{}
	v.validators = make(map[string]ParamValidator)
	tree.paramValidator = v

	tree.AddRouteParamValidators([]ParamValidator{Int32Validator{}})
	valid := tree.paramValidator.IsValid("123", "int32")
	got := valid
	if got != want {
		t.Errorf("Valid = %t, want %t", got, want)
	}
}

type dupValidator struct{}

func (dupValidator) Validate(s string, args []string) bool {
	return true
}

func (dupValidator) Type() string {
	return "int32"
}

func TestTreeAddRouteParamValidatorPanicsIfConstraintConflicts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: int32"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	tree := tree{}
	v := &parameterValidators{}
	v.validators = make(map[string]ParamValidator)
	tree.paramValidator = v

	tree.AddRouteParamValidator(Int32Validator{})
	tree.AddRouteParamValidator(dupValidator{})
}

func TestTreeAddRouteParamValidatorsPanicsIfConstraintConflicts(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "conflicting constraint type: int32"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()

	tree := tree{}
	v := &parameterValidators{}
	v.validators = make(map[string]ParamValidator)
	tree.paramValidator = v

	tree.AddRouteParamValidators([]ParamValidator{
		Int32Validator{},
		dupValidator{},
	})
}

func TestNewTreeSetsRouteParamValidators(t *testing.T) {
	want := reflect.TypeOf(&parameterValidators{}).String()
	tree := newTree()
	if tree.paramValidator == nil {
		t.Errorf("RouteParamValidators type = \"<nil>\", want %q", want)
		return
	}
	got := reflect.TypeOf(tree.paramValidator).String()
	if got != want {
		t.Errorf("RouteParamValidators type = %q, want %q", got, want)
	}
}

func TestRetrievingSingleRouteParameter(t *testing.T) {
	want := "45"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/Sleepers/{sleeperID}", "GET", testFunc)
	//tree.Print()
	_, params, _ := tree.HandlerFuncs("/Sleepers/45")
	p := params["sleeperID"]
	got := p
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestRouteIsIgnoredWhenTerminatingParameterIsInvalid(t *testing.T) {
	want := false
	validator := mockRouteParamsValidator{valid: false}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/Sleepers/{sleeperID}", "GET", testFunc)
	_, _, ok := tree.HandlerFuncs("/Sleepers/45")
	got := ok
	if got != want {
		t.Errorf("Value = %t, want %t", got, want)
	}
}

func TestRouteIsIgnoredWhenMidplacedParameterIsInvalid(t *testing.T) {
	want := false
	validator := mockRouteParamsValidator{valid: false}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/Sleepers/{sleeperID}/books", "GET", testFunc)
	_, _, ok := tree.HandlerFuncs("/Sleepers/45/books")
	got := ok
	if got != want {
		t.Errorf("Value = %t, want %t", got, want)
	}
}

func TestRetrievingMultipleRouteParameters(t *testing.T) {
	want := "45greenant"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/Sleepers/{sleeperID}/eyecolor/{color}/{insect}", "GET", testFunc)
	_, params, _ := tree.HandlerFuncs("/Sleepers/45/eyecolor/green/ant")
	p := params["sleeperID"]
	p += params["color"]
	p += params["insect"]

	got := p
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestRetrievingMultipleRouteParametersWhenManyRoutes(t *testing.T) {
	want := "45greenant"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}
	tree.AddHandlerFunc("/Cheese/{sleeperID}", "GET", testFunc)
	tree.AddHandlerFunc("/{sleeperID}/eggs", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepers/{sleeperID}", "GET", testFunc)
	tree.AddHandlerFunc("/Sleepers/{sleeperID}/eyecolor/{color}", "GET", testFunc2)
	tree.AddHandlerFunc("/Sleepers/{sleeperID}/eyecolor/{color}/{insect}", "GET", testFunc3)
	//tree.Print()
	_, params, _ := tree.HandlerFuncs("/Sleepers/45/eyecolor/green/ant")
	p := params["sleeperID"]
	p += params["color"]
	p += params["insect"]

	got := p
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestParametizedNodeIsLastToMatchIfParametizedRouteAddedLast(t *testing.T) {
	want := "testFunc"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}
	tree.AddHandlerFunc("/eyecolor/green", "GET", testFunc)
	tree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	handlers, params, _ := tree.HandlerFuncs("/eyecolor/green")

	h := handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = ""
	got = params["color"]
	if got != want {
		t.Errorf("color = %q, want %q", got, want)
	}
}

func TestParametizedNodeIsLastToMatchIfParametizedRouteAddedFirst(t *testing.T) {
	want := "testFunc"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	tree.AddHandlerFunc("/eyecolor/green", "GET", testFunc)
	handlers, params, _ := tree.HandlerFuncs("/eyecolor/green")

	h := handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = ""
	got = params["color"]
	if got != want {
		t.Errorf("color = %q, want %q", got, want)
	}
}

func TestConstrainedParametizedNodeMatchesNonconstrainedIfAddedFirst(t *testing.T) {
	want := "testFunc"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/{color:alpha}", "GET", testFunc)
	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc3)
	tree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	handlers, params, _ := tree.HandlerFuncs("/eyecolor/green")

	h := handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = "green"
	got = params["color"]
	if got != want {
		t.Errorf("color = %q, want %q", got, want)
	}
}

func TestConstrainedParametizedNodeMatchesNonconstrainedIfAddedLast(t *testing.T) {
	want := "testFunc"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc3)
	tree.AddHandlerFunc("/eyecolor/{color:alpha}", "GET", testFunc)
	handlers, params, _ := tree.HandlerFuncs("/eyecolor/green")

	h := handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = "green"
	got = params["color"]
	if got != want {
		t.Errorf("color = %q, want %q", got, want)
	}
}

func TestMatchingPathElementShorterThanLabelReturnsFalseIfNoOtherMatches(t *testing.T) {
	want := false
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	tree.AddHandlerFunc("/eyecolor/greenish", "GET", testFunc3)
	handlers, _, got := tree.HandlerFuncs("/eyecolor/green")
	h := handlers["GET"]
	name := extractFnName(h)
	fmt.Println(name)
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestAddHandlerFuncPanicsWhenMismatchingParameterBraces(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	tree := tree{}
	tree.AddHandlerFunc("/eyecolor/{color", "GET", testFunc)
}

func TestAddHandlerFuncPanicsWithCorrectMessageWhenMismatchingParameterBraces(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "invalid route syntax: {color"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		}
	}()

	tree := tree{}
	tree.AddHandlerFunc("/eyecolor/{color", "GET", testFunc)
}

// mocks

type mockRouteParamsValidator struct {
	valid bool
}

func (r mockRouteParamsValidator) AddValidator(v ParamValidator) {

}
func (r mockRouteParamsValidator) AddValidators(validators []ParamValidator) {

}
func (r mockRouteParamsValidator) IsValid(s, constraint string) bool {
	return r.valid
}

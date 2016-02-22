package mango

import (
	"reflect"
	"strconv"
	"strings"
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

func TestAddGetResourceetsNodeLabel(t *testing.T) {
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
	_, got := tree.GetResource("/Sugar")
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestRetrievingHandlersForExistentPathReturnsTrue(t *testing.T) {
	want := true
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	_, got := tree.GetResource("/Sleep")
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestRetrievingHandlersForMatchingPathReturnsNonEmptyMap(t *testing.T) {
	want := 1
	tree := tree{}
	tree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	tree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	resource, _ := tree.GetResource("/Sleep")
	got := len(resource.Handlers)
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
	resource, _ := tree.GetResource("/Sleep")
	got := len(resource.Handlers)
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
	resource, _ := tree.GetResource("/Sleep")
	h := resource.Handlers["POST"]
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
	resource, _ := tree.GetResource("/Sleepers/45")
	p := resource.RouteParams["sleeperID"]
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
	_, ok := tree.GetResource("/Sleepers/45")
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
	_, ok := tree.GetResource("/Sleepers/45/books")
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
	resource, _ := tree.GetResource("/Sleepers/45/eyecolor/green/ant")
	p := resource.RouteParams["sleeperID"]
	p += resource.RouteParams["color"]
	p += resource.RouteParams["insect"]

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
	resource, _ := tree.GetResource("/Sleepers/45/eyecolor/green/ant")
	p := resource.RouteParams["sleeperID"]
	p += resource.RouteParams["color"]
	p += resource.RouteParams["insect"]

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
	resource, _ := tree.GetResource("/eyecolor/green")

	h := resource.Handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = ""
	got = resource.RouteParams["color"]
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
	resource, _ := tree.GetResource("/eyecolor/green")

	h := resource.Handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = ""
	got = resource.RouteParams["color"]
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
	resource, _ := tree.GetResource("/eyecolor/green")

	h := resource.Handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = "green"
	got = resource.RouteParams["color"]
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
	resource, _ := tree.GetResource("/eyecolor/green")

	h := resource.Handlers["GET"]
	name := extractFnName(h)
	got := name

	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
	want = "green"
	got = resource.RouteParams["color"]
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
	_, got := tree.GetResource("/eyecolor/green")
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

func TestCORSConfigAppliedToResource(t *testing.T) {
	want := "http://greencheese.com"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	tree.AddHandlerFunc("/eyecolor/blue", "POST", testFunc3)
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
	}
	tree.SetCORS("/eyecolor/blue", config)

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v", config)
		return
	}
	got := strings.Join(res.CORSConfig.Origins, ", ")
	if got != want {
		t.Errorf("Origins = %q, want %q", got, want)
	}
}

func TestTreeGetResourceUsesGlobalCORSConfigWhenResourceConfigIsNil(t *testing.T) {
	want := "http://greencheese.com"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	tree.AddHandlerFunc("/eyecolor/blue", "POST", testFunc3)
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
	}
	tree.SetGlobalCORS(config)

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v", config)
		return
	}
	got := strings.Join(res.CORSConfig.Origins, ", ")
	if got != want {
		t.Errorf("Origins = %q, want %q", got, want)
	}
}

func TestTreeGetResourceDoesNotUseGlobalCORSConfigWhenResourceConfigIsNotNil(t *testing.T) {
	want := "http://bluecheese.com"
	validator := mockRouteParamsValidator{valid: true}
	tree := tree{paramValidator: validator}

	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	tree.AddHandlerFunc("/eyecolor/blue", "POST", testFunc3)
	rConfig := CORSConfig{
		Origins: []string{"http://bluecheese.com"},
	}
	tree.SetCORS("/eyecolor/blue", rConfig)
	gConfig := CORSConfig{
		Origins: []string{"http://greencheese.com"},
	}
	tree.SetGlobalCORS(gConfig)

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v", rConfig)
		return
	}
	got := strings.Join(res.CORSConfig.Origins, ", ")
	if got != want {
		t.Errorf("Origins = %q, want %q", got, want)
	}
}

////

func TestTreeAddCORSWhenNoGlobalConfig(t *testing.T) {
	tree := tree{}
	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	tree.AddCORS("/eyecolor/blue", rConfig)

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"http://bluecheese.com",
			"Origins",
			func(c *CORSConfig) string {
				return strings.Join(c.Origins, ", ")
			},
		},
		{
			"POST, PATCH",
			"Methods",
			func(c *CORSConfig) string {
				return strings.Join(c.Methods, ", ")
			},
		},
		{
			"X-Cheese, X-Mangoes",
			"Headers",
			func(c *CORSConfig) string {
				return strings.Join(c.Headers, ", ")
			},
		},
		{
			"X-Biscuits, X-Mangoes",
			"ExposedHeaders",
			func(c *CORSConfig) string {
				return strings.Join(c.ExposedHeaders, ", ")
			},
		},
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v\n", rConfig)
		return
	}

	for _, test := range tests {
		if got := test.fn(res.CORSConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestTreeAddCORSIncludesDataFromGlobalConfig(t *testing.T) {
	tree := tree{}
	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)

	gConfig := CORSConfig{
		Origins:        []string{"http://greencheese.com"},
		Methods:        []string{"PUT"},
		Headers:        []string{"X-Custard", "X-Fish"},
		ExposedHeaders: []string{"X-Onions"},
	}
	tree.SetGlobalCORS(gConfig)

	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	tree.AddCORS("/eyecolor/blue", rConfig)

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"http://greencheese.com, http://bluecheese.com",
			"Origins",
			func(c *CORSConfig) string {
				return strings.Join(c.Origins, ", ")
			},
		},
		{
			"PUT, POST, PATCH",
			"Methods",
			func(c *CORSConfig) string {
				return strings.Join(c.Methods, ", ")
			},
		},
		{
			"X-Custard, X-Fish, X-Cheese, X-Mangoes",
			"Headers",
			func(c *CORSConfig) string {
				return strings.Join(c.Headers, ", ")
			},
		},
		{
			"X-Onions, X-Biscuits, X-Mangoes",
			"ExposedHeaders",
			func(c *CORSConfig) string {
				return strings.Join(c.ExposedHeaders, ", ")
			},
		},
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v\n", rConfig)
		return
	}

	for _, test := range tests {
		if got := test.fn(res.CORSConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestTreeAddCORSRemovesDuplicates(t *testing.T) {
	tree := tree{}
	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)

	gConfig := CORSConfig{
		Origins:        []string{"http://greencheese.com", "http://bluecheese.com"},
		Methods:        []string{"PUT", "PATCH"},
		Headers:        []string{"X-Custard", "X-Fish", "X-Mangoes"},
		ExposedHeaders: []string{"X-Onions", "X-Biscuits"},
	}
	tree.SetGlobalCORS(gConfig)

	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	tree.AddCORS("/eyecolor/blue", rConfig)

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"http://greencheese.com, http://bluecheese.com",
			"Origins",
			func(c *CORSConfig) string {
				return strings.Join(c.Origins, ", ")
			},
		},
		{
			"PUT, PATCH, POST",
			"Methods",
			func(c *CORSConfig) string {
				return strings.Join(c.Methods, ", ")
			},
		},
		{
			"X-Custard, X-Fish, X-Mangoes, X-Cheese",
			"Headers",
			func(c *CORSConfig) string {
				return strings.Join(c.Headers, ", ")
			},
		},
		{
			"X-Onions, X-Biscuits, X-Mangoes",
			"ExposedHeaders",
			func(c *CORSConfig) string {
				return strings.Join(c.ExposedHeaders, ", ")
			},
		},
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v\n", rConfig)
		return
	}

	for _, test := range tests {
		if got := test.fn(res.CORSConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestTreeAddCORSCredentialsAndMaxAgeOverrideGlobalSetting(t *testing.T) {
	tree := tree{}
	tree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)

	gConfig := CORSConfig{
		Credentials: false,
		MaxAge:      32,
	}
	tree.SetGlobalCORS(gConfig)

	rConfig := CORSConfig{
		Credentials: true,
		MaxAge:      45,
	}
	tree.AddCORS("/eyecolor/blue", rConfig)

	tests := []struct {
		want string
		name string
		fn   func(c *CORSConfig) string
	}{
		{
			"true",
			"Credentials",
			func(c *CORSConfig) string {
				return strconv.FormatBool(c.Credentials)
			},
		},
		{
			"45",
			"MaxAge",
			func(c *CORSConfig) string {
				return strconv.Itoa(c.MaxAge)
			},
		},
	}

	res, _ := tree.GetResource("/eyecolor/blue")
	if res.CORSConfig == nil {
		t.Errorf("CORSConfig = <nil>, want %v\n", rConfig)
		return
	}

	for _, test := range tests {
		if got := test.fn(res.CORSConfig); got != test.want {
			t.Errorf("CORSConfig.%s = %q, want %q", test.name, got, test.want)
		}
	}
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

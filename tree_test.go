package mango

import (
	"strconv"
	"strings"
	"testing"
)

func testFunc(c *Context)  {}
func testFunc2(c *Context) {}
func testFunc3(c *Context) {}

func TestAddHandlerFuncAddsToTreeRoot(t *testing.T) {
	want := 1
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc)
	got := len(testTree.Root().children)
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddGetResourceetsNodeLabel(t *testing.T) {
	want := "/Sleep"
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc)
	got := testTree.Root().children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingTwoHandlersWithCompletelyDifferentRoutesAddsTwoNodesToTreeRoot(t *testing.T) {
	want := 2
	testTree := tree{}
	testTree.AddHandlerFunc("Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("Fish", "GET", testFunc)
	got := len(testTree.Root().children)
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerWithRouteStartingWithExistingRouteAddsChildToExistingRoute(t *testing.T) {
	want := 1
	testTree := tree{}
	testTree.AddHandlerFunc("Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("Sleepers", "GET", testFunc)
	got := len(testTree.Root().children[0].children)
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerWithRouteStartingWithExistingRouteSplitsNodeWithCorrectLabel(t *testing.T) {
	want := "Sleep"
	testTree := tree{}
	testTree.AddHandlerFunc("Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("Sleepers", "GET", testFunc)
	got := testTree.Root().children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingHandlerWithRouteStartingWithExistingRouteAddsChildWithCorrectLabel(t *testing.T) {
	want := "ers"
	testTree := tree{}
	testTree.AddHandlerFunc("Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("Sleepers", "GET", testFunc)
	got := testTree.Root().children[0].children[0].label
	if got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}

func TestAddingHandlerWithRouteSubstringOfExistingRouteAddsChildWithCorrectLabel(t *testing.T) {
	want := "ers"
	testTree := tree{}
	testTree.AddHandlerFunc("Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("Sleep", "GET", testFunc)
	got := testTree.Root().children[0].children[0].label
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
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepy", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleeping", "GET", testFunc)
	testTree.AddHandlerFunc("/Sugar", "GET", testFunc)
	got := testTree.GetStats().totalNodes
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerStoresHandlerInNode(t *testing.T) {
	want := "testFunc"
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	kvp := testTree.Root().children[0].handlers
	f := kvp["GET"]
	name := extractFnName(f)
	got := name
	if got != want {
		t.Errorf("Handler function = %q, want %q", got, want)
	}
}

func TestAddingHandlerToSameRouteDoesNotAddNode(t *testing.T) {
	want := 1
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepers", "POST", testFunc2)
	got := testTree.GetStats().totalNodes
	if got != want {
		t.Errorf("Node count = %d, want %d", got, want)
	}
}

func TestAddingHandlerToSameRouteUsesExisitingNode(t *testing.T) {
	want := "testFunc2"
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepers", "POST", testFunc2)
	kvp := testTree.Root().children[0].handlers
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

	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc2)
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

	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc2)
}

func TestRetrievingHandlersForNonexistentPathReturnsFalse(t *testing.T) {
	want := false
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	_, got := testTree.GetResource("/Sugar")
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestRetrievingHandlersForExistentPathReturnsTrue(t *testing.T) {
	want := true
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	_, got := testTree.GetResource("/Sleep")
	if got != want {
		t.Errorf("Result = %t, want %t", got, want)
	}
}

func TestRetrievingHandlersForMatchingPathReturnsNonEmptyMap(t *testing.T) {
	want := 1
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc2)
	resource, _ := testTree.GetResource("/Sleep")
	got := len(resource.Handlers)
	if got != want {
		t.Errorf("Handlers count = %d, want %d", got, want)
	}
}

func TestRetrievingHandlersForMatchingPathReturnsMapWithMulipleHandlers(t *testing.T) {
	want := 2
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "POST", testFunc2)
	resource, _ := testTree.GetResource("/Sleep")
	got := len(resource.Handlers)
	if got != want {
		t.Errorf("Handlers count = %d, want %d", got, want)
	}
}

func TestRetrievingCorrectHandlerForMethod(t *testing.T) {
	want := "testFunc2"
	testTree := tree{}
	testTree.AddHandlerFunc("/Sleepers", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleep", "POST", testFunc2)
	resource, _ := testTree.GetResource("/Sleep")
	h := resource.Handlers["POST"]
	name := extractFnName(h)
	got := name
	if got != want {
		t.Errorf("Handler = %q, want %q", got, want)
	}
}

func TestRetrievingSingleRouteParameter(t *testing.T) {
	want := "45"
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/Sleepers/{sleeperID}", "GET", testFunc)
	resource, _ := testTree.GetResource("/Sleepers/45")
	p := resource.RouteParams["sleeperID"]
	got := p
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestRouteIsIgnoredWhenTerminatingParameterIsInvalid(t *testing.T) {
	want := false
	validator := mockValidationHandler{valid: false}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/Sleepers/{sleeperID}", "GET", testFunc)
	_, ok := testTree.GetResource("/Sleepers/45")
	got := ok
	if got != want {
		t.Errorf("Value = %t, want %t", got, want)
	}
}

func TestRouteIsIgnoredWhenMidplacedParameterIsInvalid(t *testing.T) {
	want := false
	validator := mockValidationHandler{valid: false}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/Sleepers/{sleeperID}/books", "GET", testFunc)
	_, ok := testTree.GetResource("/Sleepers/45/books")
	got := ok
	if got != want {
		t.Errorf("Value = %t, want %t", got, want)
	}
}

func TestRetrievingMultipleRouteParameters(t *testing.T) {
	want := "45greenant"
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/Sleepers/{sleeperID}/eyecolor/{color}/{insect}", "GET", testFunc)
	resource, _ := testTree.GetResource("/Sleepers/45/eyecolor/green/ant")
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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}
	testTree.AddHandlerFunc("/Cheese/{sleeperID}", "GET", testFunc)
	testTree.AddHandlerFunc("/{sleeperID}/eggs", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepers/{sleeperID}", "GET", testFunc)
	testTree.AddHandlerFunc("/Sleepers/{sleeperID}/eyecolor/{color}", "GET", testFunc2)
	testTree.AddHandlerFunc("/Sleepers/{sleeperID}/eyecolor/{color}/{insect}", "GET", testFunc3)
	resource, _ := testTree.GetResource("/Sleepers/45/eyecolor/green/ant")
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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}
	testTree.AddHandlerFunc("/eyecolor/green", "GET", testFunc)
	testTree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	resource, _ := testTree.GetResource("/eyecolor/green")

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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	testTree.AddHandlerFunc("/eyecolor/green", "GET", testFunc)
	resource, _ := testTree.GetResource("/eyecolor/green")

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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/{color:alpha}", "GET", testFunc)
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc3)
	testTree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	resource, _ := testTree.GetResource("/eyecolor/green")

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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/{color}", "GET", testFunc2)
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc3)
	testTree.AddHandlerFunc("/eyecolor/{color:alpha}", "GET", testFunc)
	resource, _ := testTree.GetResource("/eyecolor/green")

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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	testTree.AddHandlerFunc("/eyecolor/greenish", "GET", testFunc3)
	_, got := testTree.GetResource("/eyecolor/green")
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

	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/{color", "GET", testFunc)
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

	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/{color", "GET", testFunc)
}

func TestCORSConfigAppliedToResource(t *testing.T) {
	want := "http://greencheese.com"
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	testTree.AddHandlerFunc("/eyecolor/blue", "POST", testFunc3)
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
	}
	testTree.SetCORS("/eyecolor/blue", config)

	res, _ := testTree.GetResource("/eyecolor/blue")
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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	testTree.AddHandlerFunc("/eyecolor/blue", "POST", testFunc3)
	config := CORSConfig{
		Origins: []string{"http://greencheese.com"},
	}
	testTree.SetGlobalCORS(config)

	res, _ := testTree.GetResource("/eyecolor/blue")
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
	validator := mockValidationHandler{valid: true}
	testTree := tree{validators: validator}

	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	testTree.AddHandlerFunc("/eyecolor/blue", "POST", testFunc3)
	rConfig := CORSConfig{
		Origins: []string{"http://bluecheese.com"},
	}
	testTree.SetCORS("/eyecolor/blue", rConfig)
	gConfig := CORSConfig{
		Origins: []string{"http://greencheese.com"},
	}
	testTree.SetGlobalCORS(gConfig)

	res, _ := testTree.GetResource("/eyecolor/blue")
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
func TestTreeAddCORSPanicsWhenNoGlobalConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			want := "GlobalCORS has not been set"
			got := r
			if got != want {
				t.Errorf("Error message = %q, want %q", got, want)
			}
		} else {
			t.Errorf("The code did not panic")
		}
	}()
	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	testTree.AddCORS("/eyecolor/blue", rConfig)
}

func TestTreeSetCORSWhenNoGlobalConfig(t *testing.T) {
	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)
	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	testTree.SetCORS("/eyecolor/blue", rConfig)

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

	res, _ := testTree.GetResource("/eyecolor/blue")
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
	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)

	gConfig := CORSConfig{
		Origins:        []string{"http://greencheese.com"},
		Methods:        []string{"PUT"},
		Headers:        []string{"X-Custard", "X-Fish"},
		ExposedHeaders: []string{"X-Onions"},
	}
	testTree.SetGlobalCORS(gConfig)

	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	testTree.AddCORS("/eyecolor/blue", rConfig)

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

	res, _ := testTree.GetResource("/eyecolor/blue")
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
	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)

	gConfig := CORSConfig{
		Origins:        []string{"http://greencheese.com", "http://bluecheese.com"},
		Methods:        []string{"PUT", "PATCH"},
		Headers:        []string{"X-Custard", "X-Fish", "X-Mangoes"},
		ExposedHeaders: []string{"X-Onions", "X-Biscuits"},
	}
	testTree.SetGlobalCORS(gConfig)

	rConfig := CORSConfig{
		Origins:        []string{"http://bluecheese.com"},
		Methods:        []string{"POST", "PATCH"},
		Headers:        []string{"X-Cheese", "X-Mangoes"},
		ExposedHeaders: []string{"X-Biscuits", "X-Mangoes"},
		Credentials:    true,
		MaxAge:         45,
	}
	testTree.AddCORS("/eyecolor/blue", rConfig)

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

	res, _ := testTree.GetResource("/eyecolor/blue")
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
	testTree := tree{}
	testTree.AddHandlerFunc("/eyecolor/blue", "GET", testFunc2)

	gConfig := CORSConfig{
		Credentials: false,
		MaxAge:      32,
	}
	testTree.SetGlobalCORS(gConfig)

	rConfig := CORSConfig{
		Credentials: true,
		MaxAge:      45,
	}
	testTree.AddCORS("/eyecolor/blue", rConfig)

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

	res, _ := testTree.GetResource("/eyecolor/blue")
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

func TestTreeStructureWithSingleStaticRoute(t *testing.T) {
	want := `>Label: "/Cheese/sleeper"	Handlers [GET: testFunc]	 ParamNames []
`
	testTree := tree{}
	testTree.AddHandlerFunc("/Cheese/sleeper", "GET", testFunc)
	got := testTree.Structure()
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestTreeStructureWithTwoStaticRoutesWithDifferentInitialSegment(t *testing.T) {
	want := `>Label: "/"		` + `
	>Label: "Onions/spring"	Handlers [GET: testFunc2]	 ParamNames []` + `
	>Label: "Cheese/sleeper"	Handlers [GET: testFunc]	 ParamNames []
`
	testTree := tree{}
	testTree.AddHandlerFunc("/Cheese/sleeper", "GET", testFunc)
	testTree.AddHandlerFunc("/Onions/spring", "GET", testFunc2)
	got := testTree.Structure()
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestTreeStructureWithTwoStaticRoutesWithParametizedSegments(t *testing.T) {
	want := `>Label: "/"		` + `
	>Label: "Onions/"		` + `
		>Param: ""		` + `
			>Label: "/spring"	Handlers [GET: testFunc2]	 ParamNames [season,]
	>Label: "Cheese/"		` + `
		>Param: ""		` + `
			>Label: "/sleeper/"		` + `
				>Param: ""	Handlers [GET: testFunc]	 ParamNames [eyeball,desire,]
`
	testTree := tree{}
	testTree.AddHandlerFunc("/Cheese/{desire}/sleeper/{eyeball}", "GET", testFunc)
	testTree.AddHandlerFunc("/Onions/{season}/spring", "GET", testFunc2)
	got := testTree.Structure()
	if got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

// mocks

type mockValidationHandler struct {
	valid bool
}

func (r mockValidationHandler) AddValidator(v Validator) {

}
func (r mockValidationHandler) AddValidators(validators []Validator) {

}
func (r mockValidationHandler) IsValid(val interface{}, constraints string) ([]ValidationFailure, bool) {
	return []ValidationFailure{}, r.valid
}
func (r mockValidationHandler) ParseConstraints(constraints string) map[string][]string {
	return make(map[string][]string)
}

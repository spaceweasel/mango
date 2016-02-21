package mango

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestContextRespondReturnsResponseStruct(t *testing.T) {
	want := reflect.TypeOf(&Response{})
	c := Context{}
	r := c.Respond()
	got := reflect.TypeOf(r)
	if got != want {
		t.Errorf("Response type = %s, want %s", got, want)
	}
}

func TestContextRespondReturnsResponseHoldingCorrectContext(t *testing.T) {
	want := 404
	c := Context{status: 404}
	r := c.Respond()
	got := r.context.status
	if got != want {
		t.Errorf("Response context status = %d, want %d", got, want)
	}
}

func TestResponseWithModelSetsContextModel(t *testing.T) {
	type testModel struct {
		a string
		b int
	}
	want := "cheese24"
	model := testModel{"cheese", 24}
	c := Context{}
	r := c.Respond().WithModel(model)
	m := r.context.model.(testModel)
	got := fmt.Sprintf("%s%d", m.a, m.b)
	if got != want {
		t.Errorf("Model = %q, want %q", got, want)
	}
}

func TestResponseWithStatusSetsContextStatus(t *testing.T) {
	want := 404
	c := Context{}
	r := c.Respond().WithStatus(404)
	got := r.context.status
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseWithModelReturnsItself(t *testing.T) {
	c := Context{}
	r := c.Respond()
	wmr := r.WithModel("Cheese")
	if wmr == nil {
		t.Errorf("WithModel return got nil, want Response")
	}
	if wmr != r {
		t.Errorf("WithModel did not return Response")
	}
}

func TestResponseWithStatusReturnsItself(t *testing.T) {
	c := Context{}
	r := c.Respond()
	wmr := r.WithStatus(500)
	if wmr == nil {
		t.Errorf("WithModel return got nil, want Response")
	}
	if wmr != r {
		t.Errorf("WithModel did not return Response")
	}
}

func TestResponseWithHeaderSetsResponseWriterHeader(t *testing.T) {
	want := "/somewhere"
	w := httptest.NewRecorder()
	c := Context{
		Writer: w,
	}
	c.Respond().WithHeader("Location", "/somewhere")
	got := w.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location = %q, want %q", got, want)
	}
}

func TestResponseWithContentTypeSetsResponseWriterHeader(t *testing.T) {
	want := "text/html"
	w := httptest.NewRecorder()
	c := Context{
		Writer: w,
	}
	c.Respond().WithContentType("text/html")
	got := w.HeaderMap.Get("Content-Type")
	if got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}

func TestResponseChaining(t *testing.T) {
	w := httptest.NewRecorder()
	c := Context{
		Writer: w,
	}
	r := c.Respond().WithModel("Mango biscuits").WithStatus(204)
	r.WithHeader("Location", "/somewhere").WithContentType("application/json")
	want := 204
	got := c.status
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
	wants := "Mango biscuits"
	gots := c.model
	if gots != wants {
		t.Errorf("Model = %q, want %q", gots, wants)
	}
	wants = "/somewhere"
	gots = w.HeaderMap.Get("Location")
	if gots != wants {
		t.Errorf("Location = %q, want %q", gots, wants)
	}
	wants = "application/json"
	gots = w.HeaderMap.Get("Content-Type")
	if gots != wants {
		t.Errorf("Content-Type = %q, want %q", gots, wants)
	}
}

func TestRespondWithSetsStatusWhenCalledWithInteger(t *testing.T) {
	want := 204
	c := Context{}
	c.RespondWith(204)
	got := c.status
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestRespondWithSetsPayloadBufferWhenCalledWithString(t *testing.T) {
	want := "Mango biscuits"
	c := Context{}
	c.RespondWith("Mango biscuits")

	got := string(c.payload)
	if got != want {
		t.Errorf("Status = %q, want %q", got, want)
	}
}

func TestRespondWithSetsModelWhenCalledWithArgThatNotIntegerOrString(t *testing.T) {
	type testModel struct {
		a string
		b int
	}
	want := "cheese24"
	model := testModel{"cheese", 24}

	c := Context{}
	c.RespondWith(model)
	m := c.model.(testModel)

	got := fmt.Sprintf("%s%d", m.a, m.b)

	if got != want {
		t.Errorf("Model = %q, want %q", got, want)
	}
}

func TestUrlSchemeHostPrefixesHostWithHttp(t *testing.T) {
	want := "http://www.mango.biscuits"
	r := http.Request{
		Host: "www.mango.biscuits",
	}
	c := Context{Request: &r}
	got := c.urlSchemeHost()
	if got != want {
		t.Errorf("Scheme-host = %q, want %q", got, want)
	}
}

func TestUrlSchemeHostPrefixesHostWithHttpsWhenTLS(t *testing.T) {
	want := "https://www.mango.biscuits"
	r := http.Request{
		Host: "www.mango.biscuits",
		TLS:  &tls.ConnectionState{},
	}
	c := Context{Request: &r}
	got := c.urlSchemeHost()
	if got != want {
		t.Errorf("Scheme-host = %q, want %q", got, want)
	}
}

func TestAuthenticatedReturnsTrueWhenIdentity(t *testing.T) {
	want := true
	c := Context{Identity: BasicIdentity{}}
	got := c.Authenticated()
	if got != want {
		t.Errorf("Authenticated = %t, want %t", got, want)
	}
}

func TestAuthenticatedReturnsFalseWhenNoIdentity(t *testing.T) {
	want := false
	c := Context{}
	got := c.Authenticated()
	if got != want {
		t.Errorf("Authenticated = %t, want %t", got, want)
	}
}

func TestContentDecoderReturnsErrorIfNoDecoderBasedOnContentType(t *testing.T) {
	want := "no decoder for content-type: test/mango"
	ee := &mockEncoderEngine{}
	req, _ := http.NewRequest("POST", "someurl", nil)
	req.Header.Set("Content-Type", "test/mango")
	c := Context{
		Request:       req,
		encoderEngine: ee,
	}
	_, err := c.contentDecoder()

	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestContentDecoderReturnsNoErrorIfDecoderExistsForContentType(t *testing.T) {
	want := error(nil)
	ee := &mockEncoderEngine{}
	req, _ := http.NewRequest("POST", "someurl", nil)
	req.Header.Set("Content-Type", "test/test")
	c := Context{
		Request:       req,
		encoderEngine: ee,
	}
	_, err := c.contentDecoder()

	got := err
	if got != want {
		t.Errorf("Decoder = %v, want %v", got, want)
	}
}

func TestContentDecoderReturnsDecoderBasedOnContentType(t *testing.T) {
	want := "mockDecoder"
	ee := &mockEncoderEngine{}
	req, _ := http.NewRequest("POST", "someurl", nil)
	req.Header.Set("Content-Type", "test/test")
	c := Context{
		Request:       req,
		encoderEngine: ee,
	}
	decoder, _ := c.contentDecoder()

	got := reflect.TypeOf(decoder).Name()
	if got != want {
		t.Errorf("Decoder = %q, want %q", got, want)
	}
}

func TestBindReturnsErrorIfNoSuitableDecoder(t *testing.T) {
	want := "unable to bind: no decoder for content-type: test/mango"
	ee := &mockEncoderEngine{}
	req, _ := http.NewRequest("POST", "someurl", nil)
	req.Header.Set("Content-Type", "test/mango")
	c := Context{
		Request:       req,
		encoderEngine: ee,
	}
	m := ""
	err := c.Bind(&m)

	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestBindReturnsErrorIfDecodingError(t *testing.T) {
	want := "unable to bind: incorrect interface type, expected string"
	ee := &mockEncoderEngine{}
	json := `{"id":34,"name":"Mango"}`
	req, _ := http.NewRequest("POST", "someurl", bytes.NewBufferString(json))
	req.Header.Set("Content-Type", "test/test")
	c := Context{
		Request:       req,
		encoderEngine: ee,
	}
	type model struct{}
	m := &model{}
	err := c.Bind(m)

	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestBindReturnsNoErrorIfBindingSucceeds(t *testing.T) {
	want := error(nil)
	ee := &mockEncoderEngine{}
	json := `{"id":34,"name":"Mango"}`
	req, _ := http.NewRequest("POST", "someurl", bytes.NewBufferString(json))
	req.Header.Set("Content-Type", "test/test")
	c := Context{
		Request:       req,
		encoderEngine: ee,
	}
	m := ""
	err := c.Bind(&m)

	got := err
	if got != want {
		t.Errorf("Error = %v, want %v", got, want)
	}
}

func TestBindingWithJsonBody(t *testing.T) {
	want := "Mango-34"
	json := `{"id":34,"name":"Mango"}`
	r, _ := http.NewRequest("POST", "someurl", bytes.NewBufferString(json))
	r.Header.Set("Content-Type", "application/json")
	c := Context{
		Request:       r,
		encoderEngine: newEncoderEngine(),
	}
	type data struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	decoded := data{}
	c.Bind(&decoded)

	got := fmt.Sprintf("%s-%d", decoded.Name, decoded.Id)
	if got != want {
		t.Errorf("Bind() = %q, want %q", got, want)
	}
}

func TestAcceptableMediaTypesSplitsHeaderToSlice(t *testing.T) {
	want := "text/plain|text/html"
	req, _ := http.NewRequest("GET", "someurl", nil)
	req.Header.Set("Accept", "text/plain,text/html")
	c := Context{
		Request: req,
	}
	types := c.acceptableMediaTypes()
	got := strings.Join(types, "|")
	if got != want {
		t.Errorf("Media types = %q, want %q", got, want)
	}
}

func TestAcceptableMediaTypesReturnsTrimmedTypes(t *testing.T) {
	want := "text/plain|text/html"
	req, _ := http.NewRequest("GET", "someurl", nil)
	req.Header.Set("Accept", "text/plain , text/html ")
	c := Context{
		Request: req,
	}
	types := c.acceptableMediaTypes()
	got := strings.Join(types, "|")
	if got != want {
		t.Errorf("Media types = %q, want %q", got, want)
	}
}

func TestAcceptableMediaTypesAreSortedByWeightIfQualityFactorPresent(t *testing.T) {
	want := "text/xml|text/html|text/plain|*/*"
	req, _ := http.NewRequest("GET", "someurl", nil)
	req.Header.Set("Accept", "text/html;q=0.8, text/plain;q=0.6, text/xml;q=1.0, */*;q=0.5 ")
	c := Context{
		Request: req,
	}
	types := c.acceptableMediaTypes()
	got := strings.Join(types, "|")
	if got != want {
		t.Errorf("Media types = %q, want %q", got, want)
	}
}

func TestAcceptableMediaTypesAssumeValueOneIfQualityFactorMissing(t *testing.T) {
	want := "text/xml|text/html|text/plain|*/*"
	req, _ := http.NewRequest("GET", "someurl", nil)
	req.Header.Set("Accept", "text/html;q=0.8, text/plain;q=0.6, text/xml, */*;q=0.5 ")
	c := Context{
		Request: req,
	}
	types := c.acceptableMediaTypes()
	got := strings.Join(types, "|")
	if got != want {
		t.Errorf("Media types = %q, want %q", got, want)
	}
}

func TestAcceptableMediaTypesPreferSpecificRangeOverAnyRange(t *testing.T) {
	want := "text/xml|text/plain|text/html|text/*"
	req, _ := http.NewRequest("GET", "someurl", nil)
	req.Header.Set("Accept", "text/*,text/html,text/plain,text/xml")
	c := Context{
		Request: req,
	}
	types := c.acceptableMediaTypes()
	got := strings.Join(types, "|")
	if got != want {
		t.Errorf("Media types = %q, want %q", got, want)
	}
}

func TestWhenAcceptHeaderEmptyAndDefaultMediaTypeNotSetGetEncoderReturnsError(t *testing.T) {
	want := "no encoder for content-type: "
	req, _ := http.NewRequest("GET", "someurl", nil)
	c := Context{
		Request:       req,
		encoderEngine: &encoderEngine{},
	}

	_, _, err := c.GetEncoder()

	got := err.Error()
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestWhenAcceptHeaderEmptyAndDefaultMediaTypeSetGetEncoderReturnsNoError(t *testing.T) {
	want := error(nil)
	req, _ := http.NewRequest("GET", "someurl", nil)
	c := Context{
		Request:       req,
		encoderEngine: newEncoderEngine(),
	}
	c.encoderEngine.SetDefaultMediaType("application/json")
	_, _, err := c.GetEncoder()

	got := err
	if got != want {
		t.Errorf("Error = %v, want %v", got, want)
	}
}

func TestWhenAcceptHeaderEmptyGetEncoderUsesDefaultMediaType(t *testing.T) {
	want := "application/json"
	req, _ := http.NewRequest("GET", "someurl", nil)
	c := Context{
		Request:       req,
		encoderEngine: &encoderEngine{},
	}
	c.encoderEngine.SetDefaultMediaType("application/json")

	_, contentType, _ := c.GetEncoder()

	got := contentType
	if got != want {
		t.Errorf("ContentType = %q, want %q", got, want)
	}
}

func TestWhenAcceptHeaderStarSlashStarGetEncoderUsesDefaultMediaType(t *testing.T) {
	want := "application/json"
	req, _ := http.NewRequest("GET", "someurl", nil)
	req.Header.Set("Accept", "*/*")
	c := Context{
		Request:       req,
		encoderEngine: &encoderEngine{},
	}
	c.encoderEngine.SetDefaultMediaType("application/json")

	_, contentType, _ := c.GetEncoder()

	got := contentType
	if got != want {
		t.Errorf("ContentType = %q, want %q", got, want)
	}
}

func TestErrorUsesSuppliedStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	c := Context{Writer: w}
	c.Error("an error string", 404)
	want := 404
	got := w.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestErrorUsesSuppliedErrorMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c := Context{Writer: w}
	c.Error("an error string", 404)
	want := "an error string\n"
	got := w.Body.String()
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestErrorSetsContentTypeToTextPlain(t *testing.T) {
	w := httptest.NewRecorder()
	c := Context{Writer: w}
	c.Error("an error string", 404)
	want := "text/plain; charset=utf-8"
	got := w.HeaderMap.Get("Content-Type")
	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestRedirectUsesSuppliedStatusCode(t *testing.T) {
	req, _ := http.NewRequest("GET", "someurl", nil)
	w := httptest.NewRecorder()
	c := Context{
		Request: req,
		Writer:  w,
	}
	c.Redirect("/here/be/dragons", 301)
	want := 301
	got := w.Code
	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestRedirectUsesSuppliedUrlToSetLocationHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "someurl", nil)
	w := httptest.NewRecorder()
	c := Context{
		Request: req,
		Writer:  w,
	}
	c.Redirect("/here/be/dragons", 301)
	want := "/here/be/dragons"
	got := w.HeaderMap.Get("Location")
	if got != want {
		t.Errorf("Location = %q, want %q", got, want)
	}
}

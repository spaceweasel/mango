package mango

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestEncoderEngineEncodersReturnsCorrectType(t *testing.T) {
	want := reflect.TypeOf(&map[string]EncoderFunc{}).Name()
	ee := encoderEngine{}
	e := ee.Encoders
	got := reflect.TypeOf(&e).Name()
	if got != want {
		t.Errorf("Encoders type = %q, want %q", got, want)
	}
}

func TestEncoderEngineDecodersReturnsCorrectType(t *testing.T) {
	want := reflect.TypeOf(&map[string]DecoderFunc{}).Name()
	ee := encoderEngine{}
	d := ee.Decoders
	got := reflect.TypeOf(&d).Name()
	if got != want {
		t.Errorf("Encoders type = %q, want %q", got, want)
	}
}

func TestGetEncoderReturnsErrorWhenNoEncoderFound(t *testing.T) {
	want := "no encoder for content-type: test"
	ee := encoderEngine{}
	_, err := ee.GetEncoder(nil, "test")
	got := err.Error()
	if got != want {
		t.Errorf("GetEncoder() = %q, want %q", got, want)
	}
}

func TestGetDecoderReturnsErrorWhenNoDecoderFound(t *testing.T) {
	want := "no decoder for content-type: test"
	ee := encoderEngine{}
	_, err := ee.GetDecoder(nil, "test")
	got := err.Error()
	if got != want {
		t.Errorf("GetDecoder() = %q, want %q", got, want)
	}
}

func TestNewEncoderCreatesEncodersMap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewEncoderEngine() did not initialise Encoders map")
		}
	}()
	ee := newEncoderEngine()
	ee.Encoders["test"] = func(w io.Writer) Encoder {
		return nil
	}
}

func TestNewEncoderCreatesDecodersMap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewEncoderEngine() did not initialise Decoders map")
		}
	}()
	ee := newEncoderEngine()
	ee.Decoders["test"] = func(r io.Reader) Decoder {
		return nil
	}
}

func TestGetEncoderReturnsNoErrorWhenEncoderFound(t *testing.T) {
	want := error(nil)
	ee := newEncoderEngine()
	ee.Encoders["test"] = func(w io.Writer) Encoder {
		return nil
	}
	_, err := ee.GetEncoder(nil, "test")
	got := err
	if got != want {
		t.Errorf("GetEncoder() = %v, want %v", got, want)
	}
}

func TestGetEncoderReturnsEncoderFuncWithWriter(t *testing.T) {
	want := "to be encoded"
	ee := newEncoderEngine()

	ee.Encoders["test"] = func(w io.Writer) Encoder { return NewMockEncoder(w) }

	w := new(bytes.Buffer)
	enc, _ := ee.GetEncoder(w, "test")
	enc.Encode("to be encoded")
	got := w.String()
	if got != want {
		t.Errorf("Encoded string = %q, want %q", got, want)
	}
}

type mockEncoder struct {
	w io.Writer
}

func (m mockEncoder) Encode(i interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	m.w.Write([]byte(i.(string)))
	return nil
}

func NewMockEncoder(w io.Writer) Encoder {
	return mockEncoder{w: w}
}

func TestGetDecoderReturnsNoErrorWhenDecoderFound(t *testing.T) {
	want := error(nil)
	ee := newEncoderEngine()
	ee.Decoders["test"] = func(r io.Reader) Decoder {
		return nil
	}
	_, err := ee.GetDecoder(nil, "test")
	got := err
	if got != want {
		t.Errorf("GetEncoder() = %v, want %v", got, want)
	}
}

func TestGetDecoderReturnsDecoderFuncWithReader(t *testing.T) {
	want := "to be decoded"
	ee := newEncoderEngine()

	ee.Decoders["test"] = func(r io.Reader) Decoder { return NewMockDecoder(r) }

	r := bytes.NewBufferString("to be decoded")

	dec, _ := ee.GetDecoder(r, "test")
	decoded := new(string)
	dec.Decode(decoded)
	got := *decoded
	if got != want {
		t.Errorf("Decode() = %q, want %q", got, want)
	}
}

type mockDecoder struct {
	r io.Reader
}

func (m mockDecoder) Decode(i interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(m.r)
	s, ok := i.(*string)
	if !ok {
		return fmt.Errorf("incorrect interface type, expected string")
	}
	*s = buf.String()
	return nil
}

func NewMockDecoder(r io.Reader) Decoder {
	return mockDecoder{r: r}
}

type mockEncoderEngine struct {
	DecoderRequests  []string
	EncoderRequests  []string
	defaultMediaType string
}

func (e *mockEncoderEngine) GetDecoder(r io.Reader, ct string) (Decoder, error) {
	if ct == "test/test" {
		d := NewMockDecoder(r)
		return d, nil
	}
	return nil, fmt.Errorf("no decoder for content-type: %v", ct)
}

func (e *mockEncoderEngine) GetEncoder(w io.Writer, ct string) (Encoder, error) {
	if ct == "test/test" {
		e := NewMockEncoder(w)
		return e, nil
	}
	return nil, fmt.Errorf("no encoder for content-type: %v", ct)
}

func (e *mockEncoderEngine) DefaultMediaType() string {
	return e.defaultMediaType
}

func (e *mockEncoderEngine) SetDefaultMediaType(mt string) {
	e.defaultMediaType = mt
}

func (e *mockEncoderEngine) AddEncoderFunc(ct string, fn EncoderFunc) error {
	s := ct + "-" + extractFnName(fn)
	e.EncoderRequests = append(e.EncoderRequests, s)
	if ct == "error/error" {
		return fmt.Errorf(ct)
	}
	return nil
}

func (e *mockEncoderEngine) AddDecoderFunc(ct string, fn DecoderFunc) error {
	s := ct + "-" + extractFnName(fn)
	e.DecoderRequests = append(e.DecoderRequests, s)
	if ct == "error/error" {
		return fmt.Errorf(ct)
	}
	return nil
}

// test standard encoders

func TestGetEncoderReturnsJsonEncoderFuncWhenContentTypeApplicationJson(t *testing.T) {
	want := `{"id":34,"name":"Mango"}` + "\n"
	ee := newEncoderEngine()

	type data struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	d := data{34, "Mango"}

	w := new(bytes.Buffer)
	enc, _ := ee.GetEncoder(w, "application/json")
	enc.Encode(d)
	got := w.String()
	if got != want {
		t.Errorf("GetEncoder(\"application/json\") = %q, want %q", got, want)
	}
}

func TestGetDecoderReturnsJsonDecoderFuncWhenContentTypeApplicationJson(t *testing.T) {
	want := "Mango-34"
	ee := newEncoderEngine()

	type data struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	r := bytes.NewBufferString(`{"id":34,"name":"Mango"}`)

	dec, _ := ee.GetDecoder(r, "application/json")
	decoded := data{}
	dec.Decode(&decoded)
	got := fmt.Sprintf("%s-%d", decoded.Name, decoded.Id)
	if got != want {
		t.Errorf("GetDecoder(\"application/json\") = %q, want %q", got, want)
	}
}

func TestGetEncoderReturnsXmlEncoderFuncWhenContentTypeApplicationXml(t *testing.T) {
	want := `<data><id>34</id><name>Mango</name></data>`
	ee := newEncoderEngine()

	type data struct {
		Id   int    `xml:"id"`
		Name string `xml:"name"`
	}
	d := data{34, "Mango"}

	w := new(bytes.Buffer)
	enc, _ := ee.GetEncoder(w, "application/xml")
	enc.Encode(d)
	got := w.String()
	if got != want {
		t.Errorf("GetEncoder(\"application/xml\") = %q, want %q", got, want)
	}
}

func TestGetDecoderReturnsXmlDecoderFuncWhenContentTypeApplicationXml(t *testing.T) {
	want := "Mango-34"
	ee := newEncoderEngine()

	type data struct {
		Id   int    `xml:"id"`
		Name string `xml:"name"`
	}

	r := bytes.NewBufferString(`<data><id>34</id><name>Mango</name></data>`)

	dec, _ := ee.GetDecoder(r, "application/xml")
	decoded := data{}
	dec.Decode(&decoded)
	got := fmt.Sprintf("%s-%d", decoded.Name, decoded.Id)
	if got != want {
		t.Errorf("GetDecoder(\"application/xml\") = %q, want %q", got, want)
	}
}

func TestNewEncoderEngineSetsDefaultMediaType(t *testing.T) {
	want := DefaultMediaType
	ee := newEncoderEngine()

	got := ee.defaultMediaType
	if got != want {
		t.Errorf("defaultMediaType = %q, want %q", got, want)
	}
	got = ee.DefaultMediaType()
	if got != want {
		t.Errorf("DefaultMediaType() = %q, want %q", got, want)
	}
}

func TestSetDefaultMediaType(t *testing.T) {
	want := "text/html"
	ee := newEncoderEngine()
	ee.SetDefaultMediaType("text/html")
	got := ee.defaultMediaType
	if got != want {
		t.Errorf("Default media type = %q, want %q", got, want)
	}
}

func TestAddEncoderFuncAddsFunctionWithoutError(t *testing.T) {
	want := error(nil)
	fn := func(io.Writer) Encoder {
		return nil
	}
	ee := newEncoderEngine()
	err := ee.AddEncoderFunc("mango/test", fn)
	got := err
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestAddEncoderFuncReturnsErrorWhenConflictingContentType(t *testing.T) {
	want := "conflicts with existing encoder for content-type: application/json"
	fn := func(io.Writer) Encoder {
		return nil
	}
	ee := newEncoderEngine()
	err := ee.AddEncoderFunc("application/json", fn)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %v, want %v", got, want)
	}
}

func TestAddEncoderFuncAddsToInternalMap(t *testing.T) {
	want := "to be encoded"
	fn := func(w io.Writer) Encoder {
		return NewMockEncoder(w)
	}
	ee := newEncoderEngine()
	ee.AddEncoderFunc("mango/json", fn)
	w := new(bytes.Buffer)
	enc, err := ee.GetEncoder(w, "mango/json")
	if err != nil {
		t.Errorf("Encoder = <nil>, want mockEncoder")
		return
	}
	enc.Encode("to be encoded")
	got := w.String()
	if got != want {
		t.Errorf("Encoded string = %q, want %q", got, want)
	}
}

func TestAddDecoderFuncAddsFunctionWithoutError(t *testing.T) {
	want := error(nil)
	fn := func(io.Reader) Decoder {
		return nil
	}
	ee := newEncoderEngine()
	err := ee.AddDecoderFunc("mango/test", fn)
	got := err
	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestAddDecoderFuncReturnsErrorWhenConflictingContentType(t *testing.T) {
	want := "conflicts with existing decoder for content-type: application/json"
	fn := func(io.Reader) Decoder {
		return nil
	}
	ee := newEncoderEngine()
	err := ee.AddDecoderFunc("application/json", fn)
	got := err.Error()
	if got != want {
		t.Errorf("Error = %v, want %v", got, want)
	}
}

func TestAddDecoderFuncAddsToInternalMap(t *testing.T) {
	want := "to be decoded"
	fn := func(r io.Reader) Decoder {
		return NewMockDecoder(r)
	}
	ee := newEncoderEngine()
	ee.AddDecoderFunc("mango/json", fn)

	r := bytes.NewBufferString("to be decoded")

	dec, err := ee.GetDecoder(r, "mango/json")
	if err != nil {
		t.Errorf("Decoder = <nil>, want mockDecoder")
		return
	}
	decoded := new(string)
	dec.Decode(decoded)
	got := *decoded
	if got != want {
		t.Errorf("Encoded string = %q, want %q", got, want)
	}
}

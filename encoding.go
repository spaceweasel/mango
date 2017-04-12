package mango

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

const (
	// DefaultMediaType is the default media type that will be used
	// to encode response data if no other type is specified in the
	// request Accept header.
	// This default type can be overridden by calling the
	// SetDefaultMediaType method
	DefaultMediaType = "application/json"
)

// Encoder is the interface that wraps the basic Encode method.
// Encode is used for serializing response data.
type Encoder interface {
	Encode(v interface{}) error
}

// Decoder is the interface that wraps the basic Decode method.
// Decode is used for de-serializing request data.
type Decoder interface {
	Decode(v interface{}) error
}

// EncoderEngine is the interface for the encoder/decoder management.
type EncoderEngine interface {
	GetDecoder(r io.Reader, ct string) (Decoder, error)
	GetEncoder(w io.Writer, ct string) (Encoder, error)
	DefaultMediaType() string
	SetDefaultMediaType(mt string)
	AddEncoderFunc(ct string, fn EncoderFunc) error
	AddDecoderFunc(ct string, fn DecoderFunc) error
}

// EncoderFunc is a function that returns an Encoder pre-injected
// with an io.Writer
type EncoderFunc func(io.Writer) Encoder

// DecoderFunc is a function that returns a Decoder pre-injected
// with an io.Reader
type DecoderFunc func(io.Reader) Decoder

type encoderEngine struct {
	Encoders         map[string]EncoderFunc
	Decoders         map[string]DecoderFunc
	defaultMediaType string
}

// AddEncoderFunc adds an EncoderFunc fn for the specified content-type ct.
// If an EncoderFunc pre-exists for content-type ct, then fn will not be added
// and AddEncoderFunc will return an error. Successful addition return nil.
func (e *encoderEngine) AddEncoderFunc(ct string, fn EncoderFunc) error {
	if _, ok := e.Encoders[ct]; ok {
		return fmt.Errorf("conflicts with existing encoder for content-type: %v", ct)
	}
	e.Encoders[ct] = fn
	return nil
}

// AddDecoderFunc adds a DecoderFunc fn for the specified content-type ct.
// If a DecoderFunc pre-exists for content-type ct, then fn will not be added
// and AddDecoderFunc will return an error. Successful addition return nil.
func (e *encoderEngine) AddDecoderFunc(ct string, fn DecoderFunc) error {
	if _, ok := e.Decoders[ct]; ok {
		return fmt.Errorf("conflicts with existing decoder for content-type: %v", ct)
	}
	e.Decoders[ct] = fn
	return nil
}

// GetDecoder returns a Decoder for the specified content-type (ct). The
// decoder will have the supplied io.Reader pre-injected, so decoding
// simply requires calling the Decode method, supplying the target model
// as the only parameter.
// If no suitable decoder can be found, then an error will be returned.
func (e *encoderEngine) GetDecoder(r io.Reader, ct string) (Decoder, error) {
	f, ok := e.Decoders[ct]
	if !ok {
		return nil, fmt.Errorf("no decoder for content-type: %v", ct)
	}
	return f(r), nil
}

// GetEncoder returns an Encoder for the specified content-type (ct). The
// encoder will have the supplied io.Writer pre-injected, so encoding
// simply requires calling the Encode method, supplying the model to be
// encoded as the only parameter.
// If no suitable encoder can be found, then an error will be returned.
func (e *encoderEngine) GetEncoder(w io.Writer, ct string) (Encoder, error) {
	f, ok := e.Encoders[ct]
	if !ok {
		return nil, fmt.Errorf("no encoder for content-type: %v", ct)
	}
	return f(w), nil
}

// DefaultMediaType is the default media type that will be used for
// encoding responses, if the request has not stipulated a preference in
// the Accept header.
func (e *encoderEngine) DefaultMediaType() string {
	return e.defaultMediaType
}

// SetDefaultMediaType sets the default media type.
func (e *encoderEngine) SetDefaultMediaType(mt string) {
	e.defaultMediaType = mt
}

func newEncoderEngine() *encoderEngine {
	e := encoderEngine{}
	e.defaultMediaType = DefaultMediaType
	e.Encoders = make(map[string]EncoderFunc)
	e.Decoders = make(map[string]DecoderFunc)

	e.Decoders["application/json"] = func(r io.Reader) Decoder {
		return json.NewDecoder(r)
	}
	e.Decoders["application/xml"] = func(r io.Reader) Decoder {
		return xml.NewDecoder(r)
	}
	e.Encoders["application/json"] = func(w io.Writer) Encoder {
		return json.NewEncoder(w)
	}
	e.Encoders["application/xml"] = func(w io.Writer) Encoder {
		return xml.NewEncoder(w)
	}

	return &e
}

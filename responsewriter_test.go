package mango

import (
	"compress/flate"
	"compress/gzip"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestNewResponseWriterSetsStatusToOK(t *testing.T) {
	want := 200
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	got := resp.status

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteHeaderCallsInternalWriteHeader(t *testing.T) {
	want := 404
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.WriteHeader(404)

	got := w.Code

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteHeaderSetsStatus(t *testing.T) {
	want := 404
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.WriteHeader(404)
	got := resp.status

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteCallsInternalWrite(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterWriteUpdatesByteCount(t *testing.T) {
	want := 22
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)
	got := resp.byteCount

	if got != want {
		t.Errorf("ByteCount = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteReturnsBytesWritten(t *testing.T) {
	want := 22
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	var bytes = []byte("mangoes in the morning")
	written, _ := resp.Write(bytes)
	got := written

	if got != want {
		t.Errorf("Written = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteAddsToBytesWritten(t *testing.T) {
	want := 54
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)
	bytes = []byte(" are like peaches in the evening")
	resp.Write(bytes)
	got := resp.byteCount

	if got != want {
		t.Errorf("Total bytes written = %d, want %d", got, want)
	}
}

func TestResponseWriterSetsHeader(t *testing.T) {
	want := "application/mango"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.Header().Set("Content-Type", "application/mango")

	got := w.HeaderMap.Get("Content-Type")

	if got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}

func TestResponseWriterWriteHeaderSetsHeadersSentFlag(t *testing.T) {
	want := true
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.WriteHeader(404)
	got := resp.headersSent

	if got != want {
		t.Errorf("HeadersSent = %t, want %t", got, want)
	}
}

func TestResponseWriterWriteHeaderSetsRespondedFlag(t *testing.T) {
	want := true
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.WriteHeader(404)
	got := resp.responded

	if got != want {
		t.Errorf("Responded = %t, want %t", got, want)
	}
}

func TestResponseWriterWriteSetsRespondedFlag(t *testing.T) {
	want := true
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := resp.responded

	if got != want {
		t.Errorf("Responded = %t, want %t", got, want)
	}
}

func TestResponseWriterWriteHeaderIsIgnoredWhenReadonly(t *testing.T) {
	want := 200
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.readonly = true
	resp.WriteHeader(404)
	got := w.Code

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteHeaderIsIgnoredWhenHeadersSent(t *testing.T) {
	want := 200
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.headersSent = true
	resp.WriteHeader(404)
	got := w.Code

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteSendsNoBytesWhenReadonly(t *testing.T) {
	want := 0
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.readonly = true
	var bytes = []byte("mangoes in the morning")
	bc, _ := resp.Write(bytes)

	got := bc

	if got != want {
		t.Errorf("Body = %d, want %d", got, want)
	}
}

func TestResponseWriterWriteReturnsErrorWhenReadonly(t *testing.T) {
	want := "write method has been called already"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.readonly = true
	var bytes = []byte("mangoes in the morning")
	_, err := resp.Write(bytes)

	if err == nil {
		t.Errorf("Error = <nil>, want %q", want)
	}
	got := err.Error()

	if got != want {
		t.Errorf("Error = %q, want %q", got, want)
	}
}

func TestResponseWriterWriteIsIgnoredWhenReadonly(t *testing.T) {
	want := ""
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.readonly = true
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterDoesNotSetUnderlyingHeaderWhenReadonly(t *testing.T) {
	want := "application/mango"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "", 0)
	resp.Header().Set("Content-Type", "application/mango")
	resp.readonly = true
	resp.Header().Set("Content-Type", "application/biscuits")
	got := w.HeaderMap.Get("Content-Type")

	if got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}

func TestResponseWriterDoesNotCompressWhenMinLengthEqualToZero(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "gzip", 0)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterDoesNotCompressWhenMinLengthAboveZeroAndHigherThanBodySize(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "gzip", 100)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterDoesNotCompressWhenMinLengthAboveZeroAndLessThanBodySizeButNoMatchingEncoder(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "zipme", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterSetsContentEncodingHeaderWhenGzipAccepted(t *testing.T) {
	want := "gzip"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "gzip", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.HeaderMap.Get("Content-Encoding")

	if got != want {
		t.Errorf("Content encoding = %q, want %q", got, want)
	}
}

func TestResponseWriterSetsContentEncodingHeaderWhenDeflateAccepted(t *testing.T) {
	want := "deflate"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "deflate", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.HeaderMap.Get("Content-Encoding")

	if got != want {
		t.Errorf("Content encoding = %q, want %q", got, want)
	}
}

func TestResponseWriterSelectsFirstMatchingEncoderWhenMultipleAccepted(t *testing.T) {
	want := "deflate"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "deflate, gzip", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.HeaderMap.Get("Content-Encoding")

	if got != want {
		t.Errorf("Content encoding = %q, want %q", got, want)
	}
}

func TestResponseWriterHandlesSpacesInAcceptEncoding(t *testing.T) {
	want := "gzip"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, " gzip , deflate ", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.HeaderMap.Get("Content-Encoding")

	if got != want {
		t.Errorf("Content encoding = %q, want %q", got, want)
	}
}

func TestResponseWriterUsesGzipCompressionWhenGzipAcceptedAndDataAboveMinimum(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "gzip", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)
	r, _ := gzip.NewReader(w.Body)
	defer r.Close()
	s, _ := ioutil.ReadAll(r)

	got := string(s)

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterUsesDeflateCompressionWhenDeflateAcceptedAndDataAboveMinimum(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "deflate", 5)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)
	r := flate.NewReader(w.Body)
	defer r.Close()
	s, _ := ioutil.ReadAll(r)

	got := string(s)

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterDoesNotCompressWhenHeadersAlreadySent(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "gzip", 5)
	resp.WriteHeader(200)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestResponseWriterWriteReturnsCompressedBytesWrittenWhenGzipAccepted(t *testing.T) {
	want := 36
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "gzip", 5)
	var bytes = []byte("mango mango mango mango mango mango ")
	got, _ := resp.Write(bytes)

	if got >= want {
		t.Errorf("Bytes written = %d, want less than %d", got, want)
	}
}

func TestResponseWriterWriteReturnsCompressedBytesWrittenWhenDeflateAccepted(t *testing.T) {
	want := 36
	w := httptest.NewRecorder()
	resp := NewResponseWriter(w, "deflate", 5)
	var bytes = []byte("mango mango mango mango mango mango ")
	got, _ := resp.Write(bytes)

	if got >= want {
		t.Errorf("Bytes written = %d, want less than %d", got, want)
	}
}

package mango

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestResponseLogCommonFormat(t *testing.T) {
	want := "127.0.0.1 - frank [11/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	log := RequestLog{
		Start:         start,
		RemoteAddr:    "127.0.0.1",
		AccessRequest: "GET /apache_pb.gif HTTP/1.0",
		Status:        200,
		BytesOut:      2326,
		UserID:        "frank",
	}

	got := log.CommonFormat()

	if got != want {
		t.Errorf("Log = %q, want %q", got, want)
	}
}

func TestResponseLogCombinedFormat(t *testing.T) {
	want := "127.0.0.1 - frank [11/Oct/2000:13:55:36 -0700] " +
		"\"GET /apache_pb.gif HTTP/1.0\" 200 2326 " +
		"\"https://github.com/spaceweasel/mango\" " +
		"\"Mozilla/5.0 (Android; rv:12.0) Gecko/12.0 Firefox/12.0\""

	location := time.FixedZone("test", -25200)
	start := time.Date(2000, 10, 11, 13, 55, 36, 0, location)
	log := RequestLog{
		Start:         start,
		RemoteAddr:    "127.0.0.1",
		AccessRequest: "GET /apache_pb.gif HTTP/1.0",
		Status:        200,
		BytesOut:      2326,
		Referer:       "https://github.com/spaceweasel/mango",
		UserAgent:     "Mozilla/5.0 (Android; rv:12.0) Gecko/12.0 Firefox/12.0",
		UserID:        "frank",
	}

	got := log.CombinedFormat()

	if got != want {
		t.Errorf("Log = %q, want %q", got, want)
	}
}

func TestStopSetsFinishTimeToCurrentTime(t *testing.T) {
	now := time.Now().UTC()
	nowUTC = func() time.Time {
		return now
	}
	// defer func() {
	// 	nowUTC = func() time.Time {
	// 		return time.Now().UTC()
	// 	}
	// }()
	want := now

	log := RequestLog{}
	log.stop()

	got := log.Finish

	if got != want {
		t.Errorf("Finish = %v, want %v", got, want)
	}
}

func TestStopSetsDuration(t *testing.T) {
	now := time.Now().UTC()
	nowUTC = func() time.Time {
		return now
	}
	want := int64(5000000)

	log := RequestLog{}
	log.Start = now.Add(-time.Millisecond * 5)
	log.stop()

	got := log.Duration.Nanoseconds()

	if got != want {
		t.Errorf("Duration = %v, want %v", got, want)
	}
}

func TestNewRequestLogSetsHost(t *testing.T) {
	want := "github.com"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.Host

	if got != want {
		t.Errorf("Host = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsAccessRequest(t *testing.T) {
	want := "GET /spaceweasel/mango/stone.png HTTP/1.1"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.AccessRequest

	if got != want {
		t.Errorf("AccessRequest = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsUserAgent(t *testing.T) {
	want := "Mozilla/5.0 (Android; rv:12.0) Gecko/12.0 Firefox/12.0"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Android; rv:12.0) Gecko/12.0 Firefox/12.0")
	log := NewRequestLog(req)

	got := log.UserAgent

	if got != want {
		t.Errorf("UserAgent = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsRemoteAddr(t *testing.T) {
	want := "123.123.123.123:17654"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	req.RemoteAddr = "123.123.123.123:17654"
	log := NewRequestLog(req)

	got := log.RemoteAddr

	if got != want {
		t.Errorf("RemoteAddr = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsReferer(t *testing.T) {
	want := "http://www.google.com"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	req.Header.Set("Referer", "http://www.google.com")
	log := NewRequestLog(req)

	got := log.Referer

	if got != want {
		t.Errorf("Referer = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsUserID(t *testing.T) {
	want := "-"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.UserID

	if got != want {
		t.Errorf("UserID = %q, want %q", got, want)
	}
}

func TestNewRequestLogStartUsingNowFunc(t *testing.T) {
	now := time.Now().UTC()
	nowUTC = func() time.Time {
		return now
	}
	want := now

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.Start

	if got != want {
		t.Errorf("Start = %v, want %v", got, want)
	}
}

func TestNewWatchedResponseSetsStatusToOK(t *testing.T) {
	want := 200
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	got := resp.status

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteHeaderCallsInternalWriteHeader(t *testing.T) {
	want := 404
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.WriteHeader(404)

	got := w.Code

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteHeaderSetsStatus(t *testing.T) {
	want := 404
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.WriteHeader(404)
	got := resp.status

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteCallsInternalWrite(t *testing.T) {
	want := "mangoes in the morning"
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %q, want %q", got, want)
	}
}

func TestWatchedResponseWriteUpdatesByteCount(t *testing.T) {
	want := 22
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)
	got := resp.byteCount

	if got != want {
		t.Errorf("ByteCount = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteReturnsBytesWritten(t *testing.T) {
	want := 22
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	var bytes = []byte("mangoes in the morning")
	written, _ := resp.Write(bytes)
	got := written

	if got != want {
		t.Errorf("Written = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteAddsToBytesWritten(t *testing.T) {
	want := 54
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)
	bytes = []byte(" are like peaches in the evening")
	resp.Write(bytes)
	got := resp.byteCount

	if got != want {
		t.Errorf("Total bytes written = %d, want %d", got, want)
	}
}

func TestWatchedResponseSetsHeader(t *testing.T) {
	want := "application/mango"
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.Header().Set("Content-Type", "application/mango")

	got := w.HeaderMap.Get("Content-Type")

	if got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}

func TestWatchedResponseWriteHeaderSetsRespondedFlag(t *testing.T) {
	want := true
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.WriteHeader(404)
	got := resp.responded

	if got != want {
		t.Errorf("Responded = %t, want %t", got, want)
	}
}

func TestWatchedResponseWriteSetsRespondedFlag(t *testing.T) {
	want := true
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := resp.responded

	if got != want {
		t.Errorf("Responded = %t, want %t", got, want)
	}
}

func TestWatchedResponseWriteHeaderIsIgnoredWhenReadonly(t *testing.T) {
	want := 200
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.readonly = true
	resp.WriteHeader(404)
	got := w.Code

	if got != want {
		t.Errorf("Status = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteSendsNoBytesWhenReadonly(t *testing.T) {
	want := 0
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.readonly = true
	var bytes = []byte("mangoes in the morning")
	bc, _ := resp.Write(bytes)

	got := bc

	if got != want {
		t.Errorf("Body = %d, want %d", got, want)
	}
}

func TestWatchedResponseWriteReturnsErrorWhenReadonly(t *testing.T) {
	want := "write method has been called already"
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
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

func TestWatchedResponseWriteIsIgnoredWhenReadonly(t *testing.T) {
	want := ""
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.readonly = true
	var bytes = []byte("mangoes in the morning")
	resp.Write(bytes)

	got := w.Body.String()

	if got != want {
		t.Errorf("Body = %t, want %t", got, want)
	}
}

func TestWatchedResponseDoesNotSetUnderlyingHeaderWhenReadonly(t *testing.T) {
	want := "application/mango"
	w := httptest.NewRecorder()
	resp := NewWatchedResponse(w)
	resp.Header().Set("Content-Type", "application/mango")
	resp.readonly = true
	resp.Header().Set("Content-Type", "application/biscuits")
	got := w.HeaderMap.Get("Content-Type")

	if got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
}

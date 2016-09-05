package mango

import (
	"net/http"
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

func TestNowUtcReturnsTimeNow(t *testing.T) {
	want := 0
	now := time.Now().UTC()
	nowFromFn := nowUTC()
	secs := nowFromFn.Sub(now).Seconds()
	got := int(secs)

	if got != want {
		t.Errorf("Difference in seconds = %d, want %d", got, want)
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

func TestNewRequestLogSetsMethod(t *testing.T) {
	want := "GET"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.Method

	if got != want {
		t.Errorf("Method = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsURI(t *testing.T) {
	want := "/spaceweasel/mango/stone.png"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.URI

	if got != want {
		t.Errorf("URI = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsProtocol(t *testing.T) {
	want := "HTTP/1.1"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.Protocol

	if got != want {
		t.Errorf("Protocol = %q, want %q", got, want)
	}
}

func TestNewRequestLogSetsHeader(t *testing.T) {
	want := "Onions"

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	req.Header.Set("X-SomeHeader", "Onions")
	log := NewRequestLog(req)

	got := log.Header("X-SomeHeader")

	if got != want {
		t.Errorf("Header = %q, want %q", got, want)
	}
}

func TestRequestLogHeaderReturnsEmptyStringWhenMissingHeader(t *testing.T) {
	want := ""

	req, _ := http.NewRequest("GET", "https://github.com/spaceweasel/mango/stone.png", nil)
	log := NewRequestLog(req)

	got := log.Header("X-SomeMissingHeader")

	if got != want {
		t.Errorf("Header = %q, want %q", got, want)
	}
}

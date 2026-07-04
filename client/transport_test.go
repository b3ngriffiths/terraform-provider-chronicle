package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func newTestClient(t *testing.T, attempts uint) *Client {
	t.Helper()

	cli, err := NewClient(RegionEurope, "test", context.Background(), WithRequestAttempts(attempts))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return cli
}

func TestSendRequestDoesNotRetryPermanentErrors(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		http.Error(w, `{"error": {"message": "feed not found"}}`, http.StatusNotFound)
	}))
	defer server.Close()

	cli := newTestClient(t, 5)
	_, err := sendRequest(cli, server.Client(), "GET", "test", server.URL, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if got := atomic.LoadInt32(&requests); got != 1 {
		t.Errorf("expected 1 request for 404 response, got %d", got)
	}

	var apiErr *ChronicleAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected ChronicleAPIError, got %T: %v", err, err)
	}
	if apiErr.HTTPStatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.HTTPStatusCode)
	}
	if apiErr.Message != "feed not found" {
		t.Errorf("expected API error message to be preserved, got %q", apiErr.Message)
	}
}

func TestSendRequestRetriesServerErrors(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&requests, 1) < 3 {
			http.Error(w, `{"error": {"message": "boom"}}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	cli := newTestClient(t, 5)
	res, err := sendRequest(cli, server.Client(), "GET", "test", server.URL, nil)
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if string(res) != `{"ok": true}` {
		t.Errorf("unexpected response body: %s", res)
	}
	if got := atomic.LoadInt32(&requests); got != 3 {
		t.Errorf("expected 3 requests, got %d", got)
	}
}

func TestSendRequestExhaustedRetriesReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error": {"message": "unavailable"}}`, http.StatusServiceUnavailable)
	}))
	defer server.Close()

	cli := newTestClient(t, 2)
	_, err := sendRequest(cli, server.Client(), "GET", "test", server.URL, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *ChronicleAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected ChronicleAPIError, got %T: %v", err, err)
	}
	if apiErr.HTTPStatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", apiErr.HTTPStatusCode)
	}
}

func TestSendRequestNilHTTPClient(t *testing.T) {
	cli := newTestClient(t, 1)
	_, err := sendRequest(cli, nil, "GET", "test", "https://example.invalid", nil)
	if err == nil {
		t.Fatal("expected error for missing credentials, got nil")
	}
}

func TestIsRetryableError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"transport error", errors.New("connection reset"), true},
		{"api 429", &ChronicleAPIError{HTTPStatusCode: 429}, true},
		{"api 500", &ChronicleAPIError{HTTPStatusCode: 500}, true},
		{"api 400", &ChronicleAPIError{HTTPStatusCode: 400}, false},
		{"api 404", &ChronicleAPIError{HTTPStatusCode: 404}, false},
	}

	for _, tc := range cases {
		if got := isRetryableError(tc.err); got != tc.want {
			t.Errorf("%s: isRetryableError = %v, want %v", tc.name, got, tc.want)
		}
	}
}

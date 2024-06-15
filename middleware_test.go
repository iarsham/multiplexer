package multiplexer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test-middleware", "true")
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func TestHttpMiddleware(t *testing.T) {
	mux := New(http.NewServeMux(), "/api")
	dynamic := NewChain(testMiddleware)
	mux.Handle("GET /hello", dynamic.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("test-middleware") != "true" {
		t.Fatalf("expected: true, got %s", resp.Header.Get("test-middleware"))
	}
}

func TestHttpAppendMiddleware(t *testing.T) {
	mux := New(http.NewServeMux(), "/api")
	dynamic := NewChain(testMiddleware)
	protected := dynamic.Append(authMiddleware)
	mux.Handle("GET /hello", protected.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", resp.StatusCode)
	}
}

package multiplexer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpNotFound(t *testing.T) {
	mux := New(http.NewServeMux(), "/")
	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("custom not found handler"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/blah-blah")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "custom not found handler" {
		t.Fatalf("expected: custom not found handler, got %s", string(body))
	}
}

func TestHttpMethodNotAllowed(t *testing.T) {
	mux := New(http.NewServeMux(), "/")
	mux.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("custom method not allowed handler"))
	})
	mux.HandleFunc("GET api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/hello", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) == "hello" {
		t.Fatalf("expected: custom method not allowed handler, got %s", string(body))
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected: %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("test-middleware", "true")
		next.ServeHTTP(w, r)
	})
}

func TestMiddleware(t *testing.T) {
	mux := New(http.NewServeMux(), "/")
	dynamic := NewChain(testMiddleware)
	mux.Handle("GET api/hello", dynamic.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
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

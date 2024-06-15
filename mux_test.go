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
	mux := New(http.NewServeMux(), "/api")
	mux.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("custom method not allowed handler"))
	})
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
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

func TestHttpPathWithSlash(t *testing.T) {
	mux := New(http.NewServeMux(), "/api")
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/hello/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("expected with / 200 but got 404")
	}
}

func TestHttpGroupAPI(t *testing.T) {
	mux := New(http.NewServeMux(), "/api")
	usersGroup := mux.Group("/users")
	usersGroup.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user is here"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/users/fetch")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "user is here" {
		t.Fatalf("expected: user is here, got %s", string(body))
	}
}

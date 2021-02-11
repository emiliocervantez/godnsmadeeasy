package godnsmadeeasy

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiRequestInvalidEndpoint(t *testing.T) {
	request := Request{}
	client := NewClient("___https://invalid.url", "", "")
	_, _, err := client.apiRequest(request)
	var e *ErrApiRequest
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestApiRequestInvalidMethod(t *testing.T) {
	request := Request{
		Method: "___INVALID_METHOD",
	}
	client := NewClient("https://aaa.aa", "", "")
	_, _, err := client.apiRequest(request)
	var e *ErrApiRequest
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestApiRequestInvalidRequest(t *testing.T) {
	request := Request{
		Method:  http.MethodGet,
		Path:    "/",
		Queries: nil,
		Body:    []byte(""),
	}
	client := NewClient("https://aaa.aa", "", "")
	_, _, err := client.apiRequest(request)
	var e *ErrApiRequest
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestApiRequestStatus500(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer testServer.Close()
	request := Request{
		Method:  http.MethodGet,
		Path:    "/",
		Queries: map[string]string{},
		Body:    []byte("aaa"),
	}
	client := NewClient(testServer.URL, "", "")
	_, _, err := client.apiRequest(request)
	var e *ErrApiRequest
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestApiRequestStatus403(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer testServer.Close()
	request := Request{
		Method:  http.MethodGet,
		Path:    "/domains/12345",
		Queries: map[string]string{},
		Body:    []byte("aaa"),
	}
	client := NewClient(testServer.URL, "", "")
	_, _, err := client.apiRequest(request)
	var e *ErrApiRequest
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

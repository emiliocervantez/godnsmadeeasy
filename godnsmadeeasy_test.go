package godnsmadeeasy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApiRequestInvalidEndpoint(t *testing.T) {
	request := Request{}
	client := NewClient("___https://invalid.url", "", "")
	_, err := client.apiRequest(request)
	if err == nil && strings.Contains(err.Error(), "unable to parse endpoint") == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestApiRequestInvalidMethod(t *testing.T) {
	request := Request{
		Method: "___INVALID_METHOD",
	}
	client := NewClient("https://aaa.aa", "", "")
	_, err := client.apiRequest(request)
	if err == nil && strings.Contains(err.Error(), "invalid request method") == false {
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
	_, err := client.apiRequest(request)
	if err == nil && strings.Contains(err.Error(), "unable to do request") == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestApiRequestInvalidResponseStatus(t *testing.T) {
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
	_, err := client.apiRequest(request)
	if err == nil && strings.Contains(err.Error(), "http response status 500 is not 200 or 201") == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestGetSingleDomainStruct(t *testing.T) {
	respBody, _ := json.Marshal("{ \"Name\": \"bla-bla-bla\" }")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(respBody)
	}))
	defer testServer.Close()
	client := NewClient(testServer.URL, "", "")
	_, err := client.GetSingleDomainById(666)
	if err == nil && strings.Contains(err.Error(), "unable to json-unmarshal response body") == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestGetSingleDomainApiRequest(t *testing.T) {
	client := NewClient("https://aaa.aa", "", "")
	_, err := client.GetSingleDomainById(666)
	if err == nil && strings.Contains(err.Error(), "api request error") == false {
		t.Errorf("Test error: %v", err)
	}
}

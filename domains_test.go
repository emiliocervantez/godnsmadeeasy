package godnsmadeeasy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

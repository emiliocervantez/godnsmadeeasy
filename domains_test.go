package godnsmadeeasy

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSingleDomainNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte{})
	}))
	defer testServer.Close()
	client := NewClient(testServer.URL, "", "")
	_, err := client.GetSingleDomainById(666)
	var e *ErrDomainIdNotFound
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

func TestGetSingleDomainApiRequest(t *testing.T) {
	client := NewClient("https://aaa.aa", "", "")
	_, err := client.GetSingleDomainById(666)
	var e *ErrApiRequest
	if errors.As(err, &e) == false {
		t.Errorf("Test error: %v", err)
	}
}

package godnsmadeeasy

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	Endpoint   string
	Key        string
	Secret     string
	HttpClient *http.Client
}

type Request struct {
	Method  string
	Path    string
	Queries map[string]string
	Body    []byte
}

func NewClient(endpoint string, key string, secret string) *Client {
	client := Client{
		Endpoint:   endpoint,
		Key:        key,
		Secret:     secret,
		HttpClient: &http.Client{Timeout: time.Second * 10},
	}
	return &client
}

func getHttpTime() string {
	loc, _ := time.LoadLocation("GMT")
	return time.Now().In(loc).Format("Mon, 02 Jan 2006 15:04:05 GMT") // RFC2616
}

func getHmac(secret string, message string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

type ErrApiRequest struct {
	Err error
}

func (e *ErrApiRequest) Error() string {
	return fmt.Sprintf("api error: %v", e.Err)
}

func (client *Client) apiRequest(request Request) (int, []byte, error) {
	uri, err := url.Parse(client.Endpoint)
	if err != nil {
		return 0, []byte{}, &ErrApiRequest{Err: err}
	}
	uri.Path = path.Join(uri.Path, request.Path)
	q := uri.Query()
	for k, v := range request.Queries {
		q.Set(k, v)
	}
	uri.RawQuery = q.Encode()
	reqDate := getHttpTime()
	reqHmac := getHmac(client.Secret, reqDate)
	if request.Method != "GET" && request.Method != "POST" && request.Method != "PUT" && request.Method != "DELETE" {
		return 0, []byte{}, &ErrApiRequest{Err: fmt.Errorf("invalid request method %v", request.Method)}
	}
	req, _ := http.NewRequest(request.Method, uri.String(), bytes.NewBuffer(request.Body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-dnsme-apiKey", client.Key)
	req.Header.Add("x-dnsme-requestDate", reqDate)
	req.Header.Add("x-dnsme-hmac", reqHmac)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return 0, []byte{}, &ErrApiRequest{Err: err}
	}
	if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
		return 0, []byte{}, &ErrApiRequest{Err: fmt.Errorf("server error %v", resp.Status)}
	}
	if resp.StatusCode == 403 {
		return 0, []byte{}, &ErrApiRequest{Err: fmt.Errorf("auth error %v", resp.Status)}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, error(nil)
}

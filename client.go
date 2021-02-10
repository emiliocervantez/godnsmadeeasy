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

func (client *Client) apiRequest(request Request) ([]byte, error) {
	uri, err := url.Parse(client.Endpoint)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to parse endpoint: %v", err)
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
		return []byte{}, fmt.Errorf("invalid request method %v", request.Method)
	}
	req, _ := http.NewRequest(request.Method, uri.String(), bytes.NewBuffer(request.Body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-dnsme-apiKey", client.Key)
	req.Header.Add("x-dnsme-requestDate", reqDate)
	req.Header.Add("x-dnsme-hmac", reqHmac)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to do request: %v", err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return []byte{}, fmt.Errorf("http response status %v is not 200 or 201", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body, error(nil)
}

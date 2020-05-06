package godnsmadeeasy

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type client struct {
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

type DomainsList struct {
	TotalRecords  int           `json:"totalRecords"`
	TotalPackages int           `json:"totalPackages"`
	Page          int           `json:"page"`
	Data          []DomainShort `json:"data"`
}

type DomainShort struct {
	Created            int      `json:"created"`
	FolderId           int      `json:"folderId"`
	GtdEnabled         bool     `json:"gtdEnabled"`
	Updated            int      `json:"updated"`
	ProcessMulti       bool     `json:"processMulti"`
	ActiveThirdParties []string `json:"activeThirdParties"`
	PendingActionId    int      `json:"pendingActionId"`
	Name               string   `json:"name"`
	Id                 int      `json:"id"`
}

type Domain struct {
	Name                string              `json:"name"`
	Id                  int                 `json:"id"`
	Created             int                 `json:"created"`
	DelegateNameServers []string            `json:"delegateNameServers"`
	FolderId            int                 `json:"folderId"`
	GtdEnabled          bool                `json:"gtdEnabled"`
	NameServers         []NameServers       `json:"nameServers"`
	Updated             int                 `json:"updated"`
	ProcessMulti        bool                `json:"processMulti"`
	ActiveThirdParties  []string            `json:"activeThirdParties"`
	PendingActionId     int                 `json:"pendingActionId"`
	VanityId            int                 `json:"vanityId"`
	VanityNameServers   []VanityNameServers `json:"vanityNameServers"`
}

type NameServers struct {
	Ipv6 string `json:"ipv6"`
	Ipv4 string `json:"ipv4"`
	Fqdn string `json:"fqdn"`
}

type VanityNameServers struct {
	Fqdn string `json:"fqdn"`
}

type NewDomain struct {
	Name string `json:"name"`
}

type RecordsList struct {
	TotalRecords int      `json:"totalRecords"`
	TotalPages   int      `json:"totalPages"`
	Data         []Record `json:"data"`
	Page         int      `json:"page"`
}

type Record struct {
	Source      int    `json:"source"`
	Ttl         int    `json:"ttl"`
	GtdLocation string `json:"gtdLocation"`
	SourceId    int    `json:"sourceId"`
	Failover    bool   `json:"failover"`
	Monitor     bool   `json:"monitor"`
	HardLink    bool   `json:"hardLink"`
	DynamicDns  bool   `json:"dynamicDns"`
	Failed      bool   `json:"failed"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Id          int    `json:"id"`
	Type        string `json:"type"`
}

type NewRecord struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	GtdLocation string `json:"gtdLocation"`
	Ttl         int    `json:"ttl"`
}

type UpdateRecord struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Id          int    `json:"id"`
	GtdLocation string `json:"gtdLocation"`
	Ttl         int    `json:"ttl"`
}

func NewClient(endpoint string, key string, secret string) *client {
	client := client{
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

func (client *client) apiRequest(request Request) ([]byte, error) {
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

// Get single domain by id
func (client *client) GetSingleDomainById(domainId int) (Domain, error) {
	request := Request{
		http.MethodGet,
		fmt.Sprintf("/dns/managed/%v", domainId),
		map[string]string{},
		[]byte(""),
	}
	var domain Domain
	body, err := client.apiRequest(request)
	if err != nil {
		return domain, fmt.Errorf("api request error: %v", err)
	}
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return domain, fmt.Errorf("unable to json-unmarshal response body: %v", err)
	}
	return domain, err
}

// Get single domain by name
func (client *client) GetSingleDomainByName(domainName string) (Domain, error) {
	request := Request{
		http.MethodGet,
		fmt.Sprintf("/dns/managed/name"),
		map[string]string{"domainname": domainName},
		[]byte(""),
	}
	var domain Domain
	body, err := client.apiRequest(request)
	if err != nil {
		return domain, fmt.Errorf("api request error: %v", err)
	}
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return domain, fmt.Errorf("unable to json-unmarshal response body: %v", err)
	}
	return domain, err
}

// Add single domain
func (client *client) AddSingleDomain(domainName string) (Domain, error) {
	reqBody := NewDomain{Name: domainName}
	reqBodyBytes, _ := json.Marshal(reqBody)
	request := Request{
		http.MethodPost,
		fmt.Sprintf("/dns/managed/"),
		map[string]string{},
		reqBodyBytes,
	}
	var domain Domain
	body, err := client.apiRequest(request)
	if err != nil {
		return domain, fmt.Errorf("api request error: %v", err)
	}
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return domain, fmt.Errorf("unable to json-unmarshal response body: %v", err)
	}
	return domain, err
}

// Delete single domain
func (client *client) DeleteSingleDomain(domainId int) error {
	request := Request{
		http.MethodDelete,
		fmt.Sprintf("/dns/managed/%v", domainId),
		map[string]string{},
		[]byte(""),
	}
	_, err := client.apiRequest(request)
	if err != nil {
		return fmt.Errorf("api request error: %v", err)
	}
	return err
}

// Get all domains
func (client *client) GetAllDomains(index string, order string) (DomainsList, error) {
	var domains DomainsList
	if index != "name" && index != "updated" && index != "id" && index != "folder" {
		return domains, fmt.Errorf("index %v is not in name, updated, id, folder", index)
	}
	if order != "DESC" && order != "ASC" {
		return domains, fmt.Errorf("order %v is not DESC or ASC", order)
	}
	request := Request{
		http.MethodGet,
		fmt.Sprintf("/dns/managed/"),
		map[string]string{"sidx": index, "sord": order},
		[]byte(""),
	}
	body, err := client.apiRequest(request)
	if err != nil {
		return domains, fmt.Errorf("api request error: %v", err)
	}
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return domains, fmt.Errorf("unable to json-unmarshal response body: %v", err)
	}
	return domains, err
}

// get all records
func (client *client) GetAllRecords(domainId int) (RecordsList, error) {
	request := Request{
		http.MethodGet,
		fmt.Sprintf("/dns/managed/%v/records", domainId),
		map[string]string{},
		[]byte(""),
	}
	var records RecordsList
	body, err := client.apiRequest(request)
	if err != nil {
		return records, fmt.Errorf("api request error: %v", err)
	}
	err = json.Unmarshal(body, &records)
	if err != nil {
		return records, fmt.Errorf("unable to json-unmarshal response body: %v", err)
	}
	return records, err
}

// add record
func (client *client) AddRecord(domainId int,
	recordName string,
	recordType string,
	recordValue string,
	recordGtdLocation string,
	recordTtl int) (Record, error) {
	reqBody := NewRecord{
		recordName,
		recordType,
		recordValue,
		recordGtdLocation,
		recordTtl,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	request := Request{
		http.MethodPost,
		fmt.Sprintf("/dns/managed/%v/records", domainId),
		map[string]string{},
		reqBodyBytes,
	}
	var record Record
	body, err := client.apiRequest(request)
	if err != nil {
		return record, fmt.Errorf("api request error: %v", err)
	}
	err = json.Unmarshal(body, &record)
	if err != nil {
		return record, fmt.Errorf("unable to json-unmarshal response body: %v", err)
	}
	return record, err
}

// update record
func (client *client) UpdateRecord(domainId int,
	recordName string,
	recordType string,
	recordValue string,
	recordId int,
	recordGtdLocation string,
	recordTtl int) error {
	reqBody := UpdateRecord{
		recordName,
		recordType,
		recordValue,
		recordId,
		recordGtdLocation,
		recordTtl,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	request := Request{
		http.MethodPut,
		fmt.Sprintf("/dns/managed/%v/records/%v", domainId, recordId),
		map[string]string{},
		reqBodyBytes,
	}
	_, err := client.apiRequest(request)
	if err != nil {
		return fmt.Errorf("api request error: %v", err)
	}
	return err
}

// delete record
func (client *client) DeleteRecord(domainId int, recordId int) error {
	request := Request{
		http.MethodDelete,
		fmt.Sprintf("/dns/managed/%v/records/%v", domainId, recordId),
		map[string]string{},
		[]byte(""),
	}
	_, err := client.apiRequest(request)
	if err != nil {
		return fmt.Errorf("api request error: %v", err)
	}
	return err
}

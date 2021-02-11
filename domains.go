package godnsmadeeasy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

type ErrDomainNotFound struct {
	Domain string
	Err    error
}

func (e *ErrDomainNotFound) Error() string {
	return fmt.Sprintf("domain %v not found", e.Domain)
}

type ErrDomainIdNotFound struct {
	Id  int
	Err error
}

func (e *ErrDomainIdNotFound) Error() string {
	return fmt.Sprintf("domain id %v not found", e.Id)
}

type ErrDomainExists struct {
	Domain string
	Err    error
}

func (e *ErrDomainExists) Error() string {
	return fmt.Sprintf("domain %v already exists", e.Domain)
}

type ErrDomainIdPending struct {
	Id  int
	Err error
}

func (e *ErrDomainIdPending) Error() string {
	return fmt.Sprintf("domain id %v is pending", e.Id)
}

type ErrFormat struct {
	Err error
}

func (e *ErrFormat) Error() string {
	return fmt.Sprintf("request format error: %v", e.Err)
}

// Get single domain by id
func (client *Client) GetSingleDomainById(domainId int) (Domain, error) {
	request := Request{
		http.MethodGet,
		fmt.Sprintf("/dns/managed/%v", domainId),
		map[string]string{},
		[]byte(""),
	}
	var domain Domain
	status, body, err := client.apiRequest(request)
	if err != nil {
		return domain, err
	}
	if status == 404 {
		return domain, &ErrDomainIdNotFound{Id: domainId}
	}
	_ = json.Unmarshal(body, &domain)
	return domain, err
}

// Get single domain by name
func (client *Client) GetSingleDomainByName(domainName string) (Domain, error) {
	request := Request{
		http.MethodGet,
		fmt.Sprintf("/dns/managed/name"),
		map[string]string{"domainname": domainName},
		[]byte(""),
	}
	var domain Domain
	status, body, err := client.apiRequest(request)
	if err != nil {
		return domain, &ErrApiRequest{Err: err}
	}
	if status == 404 {
		return domain, &ErrDomainNotFound{Domain: domainName}
	}
	_ = json.Unmarshal(body, &domain)
	return domain, err
}

// Add single domain
func (client *Client) AddSingleDomain(domainName string) (Domain, error) {
	reqBody := NewDomain{Name: domainName}
	reqBodyBytes, _ := json.Marshal(reqBody)
	request := Request{
		http.MethodPost,
		fmt.Sprintf("/dns/managed/"),
		map[string]string{},
		reqBodyBytes,
	}
	var domain Domain
	status, body, err := client.apiRequest(request)
	if err != nil {
		return domain, &ErrApiRequest{Err: err}
	}
	if status == 400 {
		return domain, &ErrDomainExists{Domain: domainName}
	}
	_ = json.Unmarshal(body, &domain)
	return domain, err
}

// Delete single domain
func (client *Client) DeleteSingleDomain(domainId int) error {
	request := Request{
		http.MethodDelete,
		fmt.Sprintf("/dns/managed/%v", domainId),
		map[string]string{},
		[]byte(""),
	}
	status, _, err := client.apiRequest(request)
	if err != nil {
		return &ErrApiRequest{Err: err}
	}
	if status == 404 {
		return &ErrDomainIdNotFound{Id: domainId}
	}
	if status == 400 {
		return &ErrDomainIdPending{Id: domainId}
	}
	return err
}

// Get all domains
func (client *Client) GetAllDomains(index string, order string) (DomainsList, error) {
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
	_, body, err := client.apiRequest(request)
	if err != nil {
		return domains, &ErrApiRequest{Err: err}
	}
	_ = json.Unmarshal(body, &domains)
	return domains, err
}

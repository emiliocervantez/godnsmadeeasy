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

// Get single domain by id
func (client *Client) GetSingleDomainById(domainId int) (Domain, error) {
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
func (client *Client) GetSingleDomainByName(domainName string) (Domain, error) {
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
func (client *Client) DeleteSingleDomain(domainId int) error {
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

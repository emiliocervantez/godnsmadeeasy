package godnsmadeeasy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

// get all records
func (client *Client) GetAllRecords(domainId int) (RecordsList, error) {
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
func (client *Client) AddRecord(domainId int,
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
func (client *Client) UpdateRecord(domainId int,
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
func (client *Client) DeleteRecord(domainId int, recordId int) error {
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

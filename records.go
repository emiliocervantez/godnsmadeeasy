package godnsmadeeasy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type ErrRecordExists struct {
	Type  string
	Name  string
	Value string
	Err   error
}

func (e *ErrRecordExists) Error() string {
	return fmt.Sprintf("record [type: %v, name: %v, value: %v] already exists in dme", e.Type, e.Name, e.Value)
}

type ErrDomainIdOrRecordIdNotFound struct {
	DomainId int
	RecordId int
	Err      error
}

func (e *ErrDomainIdOrRecordIdNotFound) Error() string {
	return fmt.Sprintf("domain id %v or record id %v not found in dme", e.DomainId, e.RecordId)
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
	status, body, err := client.apiRequest(request)
	if err != nil {
		return records, &ErrApiRequest{Err: err}
	}
	if status == 404 {
		return records, &ErrDomainIdNotFound{Id: domainId}
	}
	_ = json.Unmarshal(body, &records)
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
	status, body, err := client.apiRequest(request)
	if err != nil {
		return record, &ErrApiRequest{Err: err}
	}
	if status == 400 {
		var eBody errBody
		_ = json.Unmarshal(body, &eBody)
		if strings.Contains(eBody.Error[0], "already exists") {
			return record, &ErrRecordExists{
				Type:  recordType,
				Name:  recordName,
				Value: recordValue,
			}
		} else {
			return record, &ErrFormat{Err: fmt.Errorf(eBody.Error[0])}
		}
	}
	_ = json.Unmarshal(body, &record)
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
	status, _, err := client.apiRequest(request)
	if err != nil {
		return &ErrApiRequest{Err: err}
	}
	if status == 404 {
		return &ErrDomainIdOrRecordIdNotFound{
			DomainId: domainId,
			RecordId: recordId,
		}
	}
	if status == 400 {
		return &ErrFormat{Err: err}
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
	status, _, err := client.apiRequest(request)
	if err != nil {
		return &ErrApiRequest{Err: err}
	}
	if status == 404 {
		return &ErrDomainIdOrRecordIdNotFound{
			DomainId: domainId,
			RecordId: recordId,
		}
	}
	return err
}

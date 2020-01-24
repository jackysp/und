package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	prefixAPI       = "https://www.namesilo.com/api/"
	fixedFormat     = "?version=1&type=xml&key=%s&domain=%s"
	getDNSFormat    = prefixAPI + "dnsListRecords" + fixedFormat
	updateDNSFormat = prefixAPI + "dnsUpdateRecord" + fixedFormat + "&rrid=%s&rrhost=%s&rrvalue=%s&rrttl=7207"
)

type Request struct {
	Operation string `xml:"operation"`
	IP        string `xml:"ip"`
}

type Status struct {
	Code   int    `xml:"code"`
	Detail string `xml:"detail"`
}

type DNSRecord struct {
	RecordID string `xml:"record_id"`
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Value    string `xml:"value"`
	TTL      int    `xml:"ttl"`
	Distance int    `xml:"distance"`
}

type ListReply struct {
	Status
	DNSRecords []DNSRecord `xml:"resource_record"`
}

type ListResp struct {
	XMLName   xml.Name  `xml:"namesilo"`
	Request   Request   `xml:"request"`
	ListReply ListReply `xml:"reply"`
}

type UpdateReply struct {
	Status
	RecordID string `xml:"record_id"`
}

type UpdateResp struct {
	XMLName     xml.Name    `xml:"namesilo"`
	Request     Request     `xml:"request"`
	UpdateReply UpdateReply `xml:"reply"`
}

func sendRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			log.Printf("[%v] resp.Body.Close() failed, error: %v", time.Now().Format(time.RFC3339), errClose)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("resp.StatusCode is not OK, code: %d, body: %s", resp.StatusCode, body))
	}
	return body, nil
}

func dnsList(domain, key string) (*ListResp, error) {
	url := fmt.Sprintf(getDNSFormat, key, domain)
	data, err := sendRequest(url)
	if err != nil {
		return nil, err
	}
	resp := &ListResp{}
	err = xml.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}
	log.Printf("[%v] list: %v", time.Now().Format(time.RFC3339), resp)
	if resp.ListReply.Code != 300 {
		return nil, errors.New(fmt.Sprintf("wrong response, code: %d, detail:%s", resp.ListReply.Code, resp.ListReply.Detail))
	}
	return resp, nil
}

func dnsUpdate(key, domain, recordID, host, value string) error {
	data, err := sendRequest(fmt.Sprintf(updateDNSFormat, key, domain, recordID, host, value))
	if err != nil {
		return err
	}
	resp := &UpdateResp{}
	if err = xml.Unmarshal(data, resp); err != nil {
		return err
	}
	log.Printf("[%v] list: %v", time.Now().Format(time.RFC3339), resp)
	if resp.UpdateReply.Code != 300 {
		return fmt.Errorf("code not ok, code:%d, msg:%s", resp.UpdateReply.Code, resp.UpdateReply.Detail)
	}
	return nil
}

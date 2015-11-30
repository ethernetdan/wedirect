package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ethernetdan/cloudflare"
)

const defaultDomain = "google.com"

func createOrGetRecord(domain, recordName string) (*cloudflare.Record, error) {
	records, err := cf.RetrieveRecordsByName(domain, recordName, false)
	if err != nil {
		return createRecord(domain, recordName)
	} else {
		for _, record := range records {
			if record.Type == "CNAME" {
				return &record, nil
			}
		}
		return createRecord(domain, recordName)
	}
}

func createRecord(domain, recordName string) (*cloudflare.Record, error) {
	record, err := cf.CreateRecord(domain, &cloudflare.CreateRecord{
		Type:    "CNAME",
		Name:    recordName,
		Content: defaultDomain,
		Ttl:     "1",
	})
	if err != nil {
		err = fmt.Errorf("Was unable to get or create new record for domain `%s` and record `%s`: %v", domain, recordName, err)
		log.Error(err)
		return nil, err
	}
	return record, nil
}

func updateRecord(domain, id, recordName, value string) error {
	return cf.UpdateRecord(domain, id, &cloudflare.UpdateRecord{
		Type:    "CNAME",
		Name:    recordName,
		Content: value,
		Ttl:     "1",
		Proxied: true,
	})
}

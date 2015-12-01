package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/melvinmt/firebase"
)

type Domain struct {
	Time time.Time
	Name string
}

type DomainStore struct {
	history *firebase.Reference
	domain  *firebase.Reference
}

func NewDomainStore(fbUrl, auth string) DomainStore {
	historyUrl := fmt.Sprintf("%s/history", fbUrl)
	domainUrl := fmt.Sprintf("%s/domain", fbUrl)
	return DomainStore{
		history: firebase.NewReference(historyUrl).Auth(auth),
		domain:  firebase.NewReference(domainUrl).Auth(auth),
	}
}

func (d DomainStore) Domain() (domain string, err error) {
	err = d.domain.Value(&domain)
	return domain, err
}

func (d DomainStore) History() (domains []Domain, err error) {
	var m map[string]Domain
	err = d.history.Value(&m)

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, k := range keys {
		domains = append(domains, m[k])
	}
	return domains, err
}

func (d DomainStore) Set(domainName string) error {
	domain := Domain{
		Time: time.Now(),
		Name: domainName,
	}

	err := d.domain.Write(domain.Name)
	if err != nil {
		return fmt.Errorf("Failed to set current domain: %v", err)
	}

	err = d.history.Push(domain)
	if err != nil {
		return fmt.Errorf("Failed to add domain to history: %v", err)
	}
	return nil
}

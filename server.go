package main // import "github.com/ethernetdan/wedirect"

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ethernetdan/cloudflare"
)

const defaultConfigFile = "config.json"

var (
	config   Config
	store    DomainStore
	cf       *cloudflare.Client
	recordId string
)

func init() {
	configFile := defaultConfigFile
	if configFileEnv := os.Getenv("CONFIG_FILE"); len(configFileEnv) != 0 {
		configFile = configFileEnv
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("Could not read configuration file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Panicf("Could not unmarshal configuration struct from file `%s`: %v", configFile, err)
	}
}

func main() {
	log.Infof("Starting up Popular DNS for domain `%s`...", config.Domain)

	cfClient, err := cloudflare.NewClient(config.CloudFlareEmail, config.CloudFlareToken)
	if err != nil {
		log.Panicf("Could not create CloudFlare client for email `%s`: %v", config.CloudFlareEmail, err)
	} else {
		cf = cfClient
	}

	record, err := createOrGetRecord(config.Domain, config.Domain)
	if err != nil {
		log.Panic(err)
	} else {
		recordId = record.Id
	}

	store = NewDomainStore(config.FirebaseURL, config.FirebaseAuth)

	http.Handle("/ui/dist/", http.StripPrefix("/ui/dist/", http.FileServer(http.Dir("ui/dist"))))
	http.HandleFunc("/set", set)
	http.HandleFunc("/", view)
	http.ListenAndServe(":8080", nil)

}

type PageData struct {
	Domain        string
	CurrentDomain string
	History       []Domain
}

func set(w http.ResponseWriter, r *http.Request) {
	domain := r.FormValue("domain")

	// check if valid domain
	if addrs, err := net.LookupIP(domain); err != nil {
		err = fmt.Errorf("Failed to resolve entered host(%s): %v", domain, err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(addrs) < 1 {
		err = fmt.Errorf("No IPs associated with host `%s`", domain)
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if current, err := store.Domain(); strings.ToLower(current) == strings.ToLower(domain) {
		err = fmt.Errorf("Already set to `%s`", domain)
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err := updateRecord(config.Domain, recordId, config.Domain, domain)
	if err != nil {
		err = fmt.Errorf("Failed to update record: %v", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = store.Set(domain)
	if err != nil {
		err = fmt.Errorf("Failed to persist update info: %v", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully switched to %s", domain)
}

func view(w http.ResponseWriter, r *http.Request) {
	currentDomain, err := store.Domain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	history, err := store.History()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Domain:        config.Domain,
		CurrentDomain: currentDomain,
		History:       history,
	}
	renderTemplate(w, "home", &pageData)
}

func renderTemplate(w http.ResponseWriter, tmpl string, d *PageData) {
	t, _ := template.ParseFiles("ui/dist/" + tmpl + ".html")
	t.Execute(w, d)
}

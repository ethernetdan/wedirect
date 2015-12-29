package main // import "github.com/ethernetdan/wedirect"

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
)

const defaultConfigFile = "config.json"

var (
	config Config
	store  DomainStore
	proxy  *WedirectProxy
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
	log.Infof("Starting up Wedirect for domain `%s`...", config.Domain)

	store = NewDomainStore(config.FirebaseURL, config.FirebaseAuth)

	// Setup proxy, use current if available
	destination := fmt.Sprintf("http://%s:8080", config.Domain)
	if current, err := store.Domain(); err == nil {
		destination = current
	}

	url, err := url.Parse(destination)
	proxy = NewWedirectProxy(url)
	if err != nil {
		log.Fatalf("Failed to parse URL from domain: %s", config.Domain)
	}

	go func() {
		http.ListenAndServe(":80", proxy)
	}()

	// Configure UI
	ui := http.NewServeMux()
	ui.Handle("/ui/dist/", http.StripPrefix("/ui/dist/", http.FileServer(http.Dir("ui/dist"))))
	ui.HandleFunc("/set", set)
	ui.HandleFunc("/", view)

	http.ListenAndServe(":8080", ui)
}

type PageData struct {
	Domain        string
	CurrentDomain string
	History       []Domain
}

func set(w http.ResponseWriter, r *http.Request) {
	urlInput := r.FormValue("domain")

	u, err := url.Parse(urlInput)
	if err != nil {
		err = fmt.Errorf("Could not parse `%s`: %v", urlInput, err)
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlStr := u.String()
	if current, err := store.Domain(); current == urlStr {
		err = fmt.Errorf("Already set to `%s`", current)
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// Change reverse proxy destination
	proxy.Change(u)

	err = store.Set(urlStr)
	if err != nil {
		err = fmt.Errorf("Failed to persist update info: %v", err)
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully switched to %s", urlStr)
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

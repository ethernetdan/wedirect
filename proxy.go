package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

type WedirectProxy struct {
	url   *url.URL
	proxy *httputil.ReverseProxy
}

func NewWedirectProxy(target *url.URL) *WedirectProxy {
	proxy := new(WedirectProxy)
	proxy.Change(target)
	return proxy
}

func (p *WedirectProxy) Change(target *url.URL) {
	log.Infof("Changing currently set URL to %v", target)
	p.url = target
	p.proxy = httputil.NewSingleHostReverseProxy(target)
}

func (p *WedirectProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.proxy.ServeHTTP(rw, req)
}

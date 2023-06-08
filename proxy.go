package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// var webServerUrl string = "https://httpbin.org/"
var rules RuleConfig

type SimpleProxy struct {
	Proxy *httputil.ReverseProxy
}

func NewProxy(urlRaw string) (*SimpleProxy, error) {

	origin, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}
	s := &SimpleProxy{httputil.NewSingleHostReverseProxy(origin)}
	// Modify requests
	// originalDirector := s.Proxy.Director
	// s.Proxy.Director = func(r *http.Request) {
	// 	originalDirector(r)
	// 	r.Header.Set("Some-Header", "Some Value")
	// }

	// // Modify response
	// s.Proxy.ModifyResponse = func(r *http.Response) error {
	// 	// Add a response header
	// 	r.Header.Set("Server", "CodeDodle")
	// 	return nil
	// }

	return s, nil
}

func (s *SimpleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print(r.RequestURI)
	if rule, ok := IsInURI(r.RequestURI); ok {
		AnalyzeRequest(r, &rule)
	} else {
		// block by default
		BlockRequest(&w)
	}
	s.Proxy.ServeHTTP(w, r)
}

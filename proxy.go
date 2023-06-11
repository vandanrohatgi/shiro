package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/charmbracelet/log"
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
	var isBlocked bool
	log.Info(r.RequestURI)
	if rule, ok := IsInURI(r.RequestURI); ok {
		isBlocked = IsRequestBlocked(r, &rule)
		log.Debug(rule)

	} else {
		// block by default
		isBlocked = true
	}
	if isBlocked {
		io.Copy(io.Discard, r.Body)
		defer r.Body.Close()
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Error("Request blocked")
	} else {
		s.Proxy.ServeHTTP(w, r)
	}
}

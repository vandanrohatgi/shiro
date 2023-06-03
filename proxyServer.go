package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var webServerUrl string = "https://httpbin.org/"
var proxyServerPort string = ":8080"
var rules RuleConfig = RuleConfig{Path: "./rules.yaml"}

type SimpleProxy struct {
	Proxy *httputil.ReverseProxy
}

func init() {
	log.Println("Initialising...")
	rules.IngestRules()
}

func isInURL(toCheck string) bool {
	for _, i := range rules.RulesArray {
		if toCheck == i.Path {
			return true
		}
	}
	return false
}

func NewProxy(urlRaw string) (*SimpleProxy, error) {

	origin, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}
	s := &SimpleProxy{httputil.NewSingleHostReverseProxy(origin)}
	//   // Modify requests
	// originalDirector := s.Proxy.Director
	// s.Proxy.Director = func(r *http.Request) {
	// 	originalDirector(r)
	// 	r.Header.Set("Some-Header", "Some Value")
	// }

	// Modify response
	// s.Proxy.ModifyResponse = func(r *http.Response) error {
	// 	// Add a response header
	// 	r.Header.Set("Server", "CodeDodle")
	// 	return nil
	// }

	return s, nil
}

func (s *SimpleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Do anything you want here
	// e.g. blacklisting IP, log time, modify headers, etc
	if isInURL(r.RequestURI) {
		io.Copy(io.Discard, r.Body)
		defer r.Body.Close()
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
	log.Printf("Proxy receives request.")
	log.Printf("Proxy forwards request to origin.")
	// s.Proxy.ServeHTTP(w, r)
	log.Printf("Origin server completes request.")
}

func main() {
	rules.PrintRules()
	log.Println("Starting Proxy...")
	proxy, err := NewProxy(webServerUrl)
	if err != nil {
		log.Fatal("Error creating proxy", err)
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(proxyServerPort, nil))
}

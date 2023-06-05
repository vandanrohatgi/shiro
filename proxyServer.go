package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// var webServerUrl string = "https://httpbin.org/"
var webServerUrl string = "http://localhost:8888"
var proxyServerPort string = ":8080"
var rules RuleConfig

type SimpleProxy struct {
	Proxy *httputil.ReverseProxy
}

func init() {
	log.Println("Initialising...")
	if Path, ok := os.LookupEnv("SHIROPATH"); !ok {
		rules.Path = "rules.yaml"
	} else {
		rules.Path = Path
	}
	rules.IngestRules()
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
	// if IsInURL(r.RequestURI) {
	// 	io.Copy(io.Discard, r.Body)
	// 	defer r.Body.Close()
	// 	http.Error(w, "Forbidden", http.StatusForbidden)
	// }
	s.Proxy.ServeHTTP(w, r)
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

package main

import (
	"flag"
	"log"
	"net/http"
)

var targetURL, proxyPort, path string

func init() {
	log.Println("Initialising...")

	flag.StringVar(&targetURL, "targetURL", "http://localhost:8888", "URL to proxy")
	flag.StringVar(&proxyPort, "proxyPort", "8080", "port to host the proxy")
	flag.StringVar(&path, "path", "rules.yaml", "path to the rules file")
	flag.Parse()

	rules.Path = path
	rules.IngestRules()

}

func main() {
	rules.PrintRules()
	log.Println("Starting Proxy...")
	proxy, err := NewProxy(targetURL)
	if err != nil {
		log.Fatal("Error creating proxy", err)
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

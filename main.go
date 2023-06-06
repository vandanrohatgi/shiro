package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var targetURL, proxyPort string

func init() {
	log.Println("Initialising...")
	if Path, ok := os.LookupEnv("SHIROPATH"); !ok {
		rules.Path = "rules.yaml"
	} else {
		rules.Path = Path
	}
	rules.IngestRules()
	flag.StringVar(&targetURL, "targetURL", "http://localhost:8888", "URL to proxy")
	flag.StringVar(&proxyPort, "proxyPort", "8080", "port to host the proxy")
	flag.Parse()
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

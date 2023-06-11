package main

import (
	"flag"
	"net/http"

	"github.com/charmbracelet/log"
)

var targetURL, proxyPort, path string
var verbose bool

func init() {
	log.Info("Initialising...")

	flag.StringVar(&targetURL, "targetURL", "https://httpbin.org/", "URL to proxy")
	flag.StringVar(&proxyPort, "proxyPort", "8080", "port to host the proxy")
	flag.StringVar(&path, "path", "rules.yaml", "path to the rules file")
	flag.BoolVar(&verbose, "verbose", false, "Output all type of logs")
	flag.Parse()

	rules.Path = path
	rules.IngestRules()
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

}

func main() {
	rules.PrintRules()
	log.Info("Starting Proxy...")
	proxy, err := NewProxy(targetURL)
	if err != nil {
		log.Fatal("Error creating proxy", err)
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

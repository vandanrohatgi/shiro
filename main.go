package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

var targetURL, proxyPort, path string
var verbose, monitor bool
var timeout int

func init() {
	log.Info("Initialising...")

	flag.StringVar(&targetURL, "targetURL", "https://httpbin.org/", "URL to proxy")
	flag.StringVar(&proxyPort, "proxyPort", "8080", "port to host the proxy")
	flag.StringVar(&path, "path", "rules.yaml", "path to the rules file")
	flag.BoolVar(&verbose, "verbose", false, "Output all type of logs")
	flag.IntVar(&timeout, "timeout", 5, "Timout for the proxy requests")
	flag.BoolVar(&monitor, "monitor", false, "Monitor proxy traffic and generate regex automatically")
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
	proxy, err := NewProxy(targetURL, time.Duration(timeout)*time.Second, monitor)
	if err != nil {
		log.Fatal("Error creating proxy", err)
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

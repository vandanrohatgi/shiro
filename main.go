package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmbracelet/log"
)

var targetURL, proxyPort, path string
var verbose, monitor bool
var timeout int
var ruleconfig RuleConfig

// init performs CLI flags, reads the rule file for use and sets log level
func init() {
	flag.StringVar(&targetURL, "targetURL", "https://httpbin.org/", "URL to proxy")
	flag.StringVar(&proxyPort, "proxyPort", "8080", "port to host the proxy")
	flag.StringVar(&path, "path", "rules.yaml", "path to the rules file")
	flag.BoolVar(&verbose, "verbose", false, "Output all type of logs")
	flag.IntVar(&timeout, "timeout", 5, "Timout for the proxy requests")
	flag.BoolVar(&monitor, "monitor", false, "Monitor proxy traffic and generate regex automatically")
	flag.Parse()

	log.Info("Initialising...")
	ruleconfig = RuleConfig{
		Path:  path,
		Rules: make(map[string]Rules),
	}
	ruleconfig.IngestRules()
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Ingested Rules: ", ruleconfig.Rules)

}

func main() {
	// Goroutine when application is run in monitoring mode.
	// To monitor for ctrl+c (SIGINT) and writes the monitored rules to a file.
	if monitor {
		go func() {
			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, os.Interrupt)
			<-sigchan
			log.Error("Initiating exit process...")

			ruleconfig.WriteRules()
			os.Exit(0)
		}()
	}

	log.Info("Starting Proxy...")
	proxy, err := NewProxy(targetURL, time.Duration(timeout)*time.Second, monitor)
	if err != nil {
		log.Fatal("Error creating proxy", err)
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

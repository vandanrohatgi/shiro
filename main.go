package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/charmbracelet/log"
)

var targetURL = "https://httpbin.org/"
var proxyPort = "8080"
var path = "rules.yaml"
var verbose = false
var monitor = false
var timeout = 10
var ruleconfig RuleConfig

func parseFlags() {
	flag.StringVar(&targetURL, "targetURL", targetURL, "URL to proxy")
	flag.StringVar(&proxyPort, "proxyPort", proxyPort, "port to host the proxy")
	flag.StringVar(&path, "path", path, "path to the rules file")
	flag.BoolVar(&verbose, "verbose", verbose, "Output all types of logs")
	flag.IntVar(&timeout, "timeout", timeout, "Timeout for the proxy requests")
	flag.BoolVar(&monitor, "monitor", monitor, "Monitor proxy traffic and generate rules automatically")
	flag.Parse()
}

func main() {
	log.Info("Initialising...")

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	ruleconfig = RuleConfig{
		Path:  path,
		Rules: make(map[string]Rules),
	}

	log.Infof("Reading rules from %s", ruleconfig.Path)
	ruleconfig.IngestRules()
	log.Debug("Ingested Rules: %v", ruleconfig.Rules)

	// if not in monitoring mode, compile regex for blocking
	if !monitor {
		log.Info("Compiling regex...")
		for URI, rule := range ruleconfig.Rules {
			rule.BodyRegex = *regexp.MustCompile(rule.Body)
			rule.MethodRegex = *regexp.MustCompile(rule.Method)
			rule.Headers.KeyRegex = *regexp.MustCompile(rule.Headers.Key)
			rule.Headers.ValueRegex = *regexp.MustCompile(rule.Headers.Value)
			ruleconfig.Rules[URI] = rule
		}
	}

	// CLI flags
	parseFlags()

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

	// Initilaizing finished

	log.Info("Starting Proxy...")
	proxy, err := NewProxy(targetURL, time.Duration(timeout)*time.Second, monitor)
	if err != nil {
		log.Fatal("Error creating proxy", err)
	}

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(":"+proxyPort, nil))
}

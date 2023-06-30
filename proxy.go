package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// var webServerUrl string = "https://httpbin.org/"
var isBlocked bool = true // Block by default
var err error
var mutex sync.Mutex

type SimpleProxy struct {
	Proxy                 *httputil.ReverseProxy
	Timeout               time.Duration
	MonitoringModeEnabled bool
}

// NewProxy returns an instance of SimpleProxy struct with defined configurations.
func NewProxy(urlRaw string, timeout time.Duration, monitor bool) (*SimpleProxy, error) {

	origin, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}
	s := &SimpleProxy{
		Proxy:                 httputil.NewSingleHostReverseProxy(origin),
		Timeout:               timeout,
		MonitoringModeEnabled: monitor,
	}

	// Set the client for the reverse proxy
	// Main job of this code is to set the timeout for connecting with target URL
	s.Proxy.Transport = &http.Transport{
		DialContext:           (&net.Dialer{Timeout: s.Timeout}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: s.Timeout,
	}
	return s, nil
}

func (s *SimpleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Show incoming request info
	log.Infof("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)

	// Update the request's context with the client's context
	// This code is for setting the time duration for the whole process of taking the request, connecting to target URL,
	// receiving response, processing and returning proxy response to client.
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()
	r = r.WithContext(ctx)

	// Check if URI has a rule for it
	rule, ruleExists := ruleconfig.Rules[r.RequestURI]

	if s.MonitoringModeEnabled {
		monitorRequest(r, rule)
		goto serve
	}

	// Blocking mode
	if ruleExists {
		isBlocked, err = IsRequestBlocked(r, rule)
		if err != nil { // Block if anything goes out of ordinary
			log.Error("Request blocked", err)
		}
	}

	if isBlocked {
		blockRequest(&w, r)
		return
	} else {
		goto serve
	}

serve:
	s.Proxy.ServeHTTP(w, r)
}

// monitorRequest inspects the request method, headers and body.
// Generates a new regex and updates the corresponding rule in ruleconfig.Rules map
func monitorRequest(r *http.Request, rule Rules) {
	rule = InspectMethod(r, rule)
	rule = InspectHeaders(r, rule)
	rule = InspectBody(r, rule)
	// Reverse proxy deals with requests in separate goroutines. Map is not thread safe.
	mutex.Lock()
	ruleconfig.Rules[r.RequestURI] = rule
	mutex.Unlock()
}

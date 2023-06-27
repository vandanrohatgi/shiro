package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// var webServerUrl string = "https://httpbin.org/"
var rules RuleConfig

type SimpleProxy struct {
	Proxy   *httputil.ReverseProxy
	Timeout time.Duration
	Monitor bool
}

// NewProxy returns an instance of SimpleProxy struct with defined configurations.
func NewProxy(urlRaw string, timeout time.Duration, monitor bool) (*SimpleProxy, error) {

	origin, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}
	s := &SimpleProxy{
		Proxy:   httputil.NewSingleHostReverseProxy(origin),
		Timeout: timeout,
		Monitor: monitor,
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
	log.Info(r.Method, r.RequestURI)

	// Update the request's context with the client's context
	// This code is for setting the time duration for the whole process of taking the request, connecting to target URL,
	// receiving response, processing and returning proxy response to client.
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()
	r = r.WithContext(ctx)

	// TODO: improve and optimize the code below. Looks like we cn write some functions to stop being redundant

	// Check if URI has a rule for it
	var rule Rules
	ruleIndex, ruleExists := IsInURI(r.RequestURI)
	if ruleExists {
		rule = rules.RulesArray[ruleIndex]
	}

	if s.Monitor {
		if ruleExists {
			// Generate Regex for body
			body, _ := io.ReadAll(r.Body)
			rules.RulesArray[ruleIndex].Body, _ = GenerateRegex([]string{
				rule.Body,
				string(body),
			})

			// Generate Regex for headers
			for header, value := range r.Header {
				rules.RulesArray[ruleIndex].Headers.Key, _ = GenerateRegex([]string{
					header,
					rule.Headers.Key,
				})
				rules.RulesArray[ruleIndex].Headers.Value, _ = GenerateRegex([]string{
					strings.Join(value, ","),
					rule.Headers.Value,
				})
			}

			//Generate Regex for Method
			rules.RulesArray[ruleIndex].Method, _ = GenerateRegex([]string{
				r.Method,
				rule.Method,
			})
		} else {
			// Generate URI Regex
			rule.URI = fmt.Sprintf("^%s$", r.RequestURI)

			// Generate Regex for body
			body, _ := io.ReadAll(r.Body)
			rule.Body, _ = GenerateRegex([]string{
				string(body),
			})

			// Generate Regex for headers
			for header, value := range r.Header {
				rule.Headers.Key, _ = GenerateRegex([]string{
					header,
				})
				rule.Headers.Value, _ = GenerateRegex([]string{
					strings.Join(value, ","),
				})
			}

			//Generate Regex for Method
			rule.Method, _ = GenerateRegex([]string{
				r.Method,
			})

			rules.RulesArray = append(rules.RulesArray, rule)
		}

		s.Proxy.ServeHTTP(w, r)
	} else {
		var isBlocked bool = true // Block by default
		var err error

		if ruleExists {
			isBlocked, err = IsRequestBlocked(r, rule)
			if err != nil {
				log.Error("Request blocked", err)
			}
		}

		if isBlocked {
			io.Copy(io.Discard, r.Body)
			defer r.Body.Close()
			http.Error(w, "Forbidden", http.StatusForbidden)

		} else {
			s.Proxy.ServeHTTP(w, r)
		}
	}

}

package main

import (
	"context"
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

func NewProxy(urlRaw string, timeout time.Duration, monitor bool) (*SimpleProxy, error) {

	origin, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}
	return &SimpleProxy{
		Proxy:   httputil.NewSingleHostReverseProxy(origin),
		Timeout: timeout,
		Monitor: monitor,
	}, nil

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
}

func (s *SimpleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Method, r.RequestURI)

	// Set the client for the reverse proxy
	s.Proxy.Transport = &http.Transport{
		DialContext:           (&net.Dialer{Timeout: s.Timeout}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: s.Timeout,
	}

	// Update the request's context with the client's context
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()
	r = r.WithContext(ctx)

	// Check if URI has a rule for it
	ruleIndex, ok := IsInURI(r.RequestURI)
	rule := rules.RulesArray[ruleIndex]

	if s.Monitor {
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
		// TODO: Write the updated rules to rule.path file
		s.Proxy.ServeHTTP(w, r)
	} else {
		var isBlocked bool = true // Block by default

		if ok {
			isBlocked, _ = IsRequestBlocked(r, rule)
			log.Debug(rule)
		}

		if isBlocked {
			io.Copy(io.Discard, r.Body)
			defer r.Body.Close()
			http.Error(w, "Forbidden", http.StatusForbidden)
			log.Error("Request blocked")
		} else {
			s.Proxy.ServeHTTP(w, r)
		}
	}

}

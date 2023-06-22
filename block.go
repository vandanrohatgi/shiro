package main

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
)

// IsInURI checks the incoming request URI with rules from rules.yaml if a rules exists for that URI
// returns the index of the rule from rulesArray and bool for if a rule was found or not
func IsInURI(toCheck string) (int, bool) {
	for i, rule := range rules.RulesArray {
		if ok, _ := regexp.MatchString(rule.URI, toCheck); ok {
			return i, true
		}
	}
	return 0, false
}

func IsRequestBlocked(r *http.Request, rule Rules) (bool, error) {
	// TODO regex match over all the fields
	// Check Body
	bodyDecision, err := checkBody(r, rule)
	if err != nil || bodyDecision {
		return true, err
	}

	//check headers
	headerDecision, err := checkHeaders(r, rule)
	if err != nil || headerDecision {
		return true, err
	}

	return false, nil
}

func checkBody(r *http.Request, rule Rules) (bool, error) {
	body, err := io.ReadAll(r.Body)
	log.Debug(string(body[:]))
	if err != nil {
		return true, err
	}
	if ok, _ := regexp.Match(rule.Body, body); !ok {
		return true, nil
	}
	return false, nil
}

func checkHeaders(r *http.Request, rule Rules) (bool, error) {
	log.Debug(r.Header)
	for key, value := range r.Header {
		valueString := strings.Join(value, ",")
		keyOk, _ := regexp.MatchString(rule.Headers.Key, key)
		valueOk, _ := regexp.MatchString(rule.Headers.Value, valueString)
		if !(keyOk && valueOk) {
			return true, nil
		}
	}
	return false, nil

}

package main

import (
	"fmt"
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
	// Check method
	methodDecision, err := checkMethod(r, rule)
	if err != nil || methodDecision {
		return true, err
	}

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

	// request not blocked
	return false, nil
}

// checkMethod takes the incoming request and the associated rule for making blocking decision
func checkMethod(r *http.Request, rule Rules) (bool, error) {
	if ok, err := regexp.MatchString(rule.Method, r.Method); !ok {
		return true, fmt.Errorf("request method %s violates defined method %s", r.Method, rule.Method)
	} else if err != nil {
		return true, err
	}
	return false, nil
}

// checkBody takes the incoming request and the associated rule for making blocking decision
func checkBody(r *http.Request, rule Rules) (bool, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return true, err
	}
	if ok, err := regexp.Match(rule.Body, body); !ok {
		return true, fmt.Errorf("request body %s violates defined body %s", string(body[:]), rule.Body)
	} else if err != nil {
		return true, err
	}
	return false, nil
}

// checkHeaders takes the incoming request and the associated rule for making blocking decision
func checkHeaders(r *http.Request, rule Rules) (bool, error) {
	log.Debug(r.Header)
	for key, value := range r.Header {
		valueString := strings.Join(value, ",")
		keyOk, keyErr := regexp.MatchString(rule.Headers.Key, key)
		valueOk, valErr := regexp.MatchString(rule.Headers.Value, valueString)
		if !(keyOk || valueOk) {
			return true, fmt.Errorf("request header %s:%s violates defined header %s:%s",
				key,
				value,
				rule.Headers.Key,
				rule.Headers.Value)
		} else if keyErr != nil {
			return true, keyErr
		} else if valErr != nil {
			return true, valErr
		}
	}
	return false, nil

}

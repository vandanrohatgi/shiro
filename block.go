package main

import (
	"io"
	"log"
	"net/http"
	"regexp"
)

func IsInURI(toCheck string) (Rules, bool) {
	for _, i := range rules.RulesArray {
		if ok, _ := regexp.Match(i.URI, []byte(toCheck)); ok {
			return i, true
		}
	}
	return Rules{}, false
}

func IsRequestBlocked(r *http.Request, rule *Rules) bool {
	// TODO regex match over all the fields
	// Check Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Not able to read body of the request", err)
	}
	if ok, _ := regexp.Match(rule.Body, body); !ok {
		return true
	}
	// check headers
	// if ok,_:=regexp.Match(rule.Headers,[]byte(r.Header[][]));ok{
	// 	return true
	// }
	return false
}

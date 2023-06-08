package main

import "net/http"

func IsInURI(toCheck string) (Rules, bool) {
	for _, i := range rules.RulesArray {
		if toCheck == i.Path {
			return i, true
		}
	}
	return rules.RulesArray[0], false
}

func AnalyzeRequest(r *http.Request, rule *Rules) bool {
	// TODO regex match over all the fields
	//regexp.Match(rule.)
	return false
}

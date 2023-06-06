package main

func IsInURI(toCheck string) (Rules, bool) {
	for _, i := range rules.RulesArray {
		if toCheck == i.Path {
			return i, true
		}
	}
	return Rules{AllowPatternURI: ".*", AllowPatternBody: ".*", AllowPatternHeaders: ".*"}, false
}

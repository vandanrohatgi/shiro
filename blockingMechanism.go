package main

func IsInURI(toCheck string) (string, bool) {
	for _, i := range rules.RulesArray {
		if toCheck == i.Path {
			return i.AllowPatternURI, true
		}
	}
	return "", false
}

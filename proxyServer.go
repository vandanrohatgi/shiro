package main

import (
	// "log"
	// "net/http"
	// "net/http/httputil"
	// "net/url"

	rules "github.com/vandanrohatgi/shiro_waf/ruleSetup"
)

var webServerUrl string = "https://httpbin.org/"
var proxyServerPort string = ":8080"

func main() {
	rules.PrintRules()
	// log.Println("Starting Proxy...")

	// origin, err := url.Parse(webServerUrl)
	// if err != nil {
	// 	panic(err)
	// }
	// proxy := httputil.NewSingleHostReverseProxy(origin)

	// http.Handle("/", proxy)
	// log.Fatal(http.ListenAndServe(proxyServerPort, nil))

}

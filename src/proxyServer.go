package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var webServerUrl string = "https://gorest.co.in/"
var proxyServerPort string = ":8080"

func main() {
	//go testServer.Server(webServerPort)
	fmt.Println("Starting Proxy...")
	//origin, err := url.Parse("http://localhost" + webServerPort)
	origin, err := url.Parse(webServerUrl)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(origin)

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(proxyServerPort, nil))
}

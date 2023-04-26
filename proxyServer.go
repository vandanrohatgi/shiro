package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/vandanrohatgi/shiro_waf/testServer"
)

var webServerPort string = ":8081"
var proxyServerPort string = ":8080"

func main() {
	go testServer.Server(webServerPort)
	fmt.Println("Starting Proxy...")
	origin, err := url.Parse("http://localhost" + webServerPort)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(origin)

	http.Handle("/", proxy)
	log.Fatal(http.ListenAndServe(proxyServerPort, nil))
}

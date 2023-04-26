package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[origin server] received request at: %s\n", time.Now())
		_, _ = fmt.Fprint(rw, "origin server response")
	})

	log.Fatal(http.ListenAndServe(":8081", originServerHandler))
}

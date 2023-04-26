package testServer

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Server(port string) {
	fmt.Println("Starting Web server...")
	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[origin server] received request at: %s\n", time.Now())
		_, _ = fmt.Fprint(rw, "origin server response")
	})

	log.Fatal(http.ListenAndServe(":8081", originServerHandler))
	fmt.Println("Web server started")
}

package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "[Origin server]\n\n")
		fmt.Fprintf(rw, "Header\n\n")
		for key, value := range r.Header {
			fmt.Fprintf(rw, "%q: %q\n", key, value)
		}

		fmt.Fprintf(rw, "\n\nBody\n\n")
		fmt.Fprintf(rw, "%q", r.Body)
	})
	http.ListenAndServe(":8888", nil)
}

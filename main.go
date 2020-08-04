package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	h := http.HandlerFunc(Echo)

	if err := http.ListenAndServe("localhost:8080", h); err != nil {
		log.Fatalf("error: listening and serving: %s", err)
	}
}

// Echo is a basic HTTP Handler.
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)
}

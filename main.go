package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// ===

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// ==
	// Start API Service
	api := http.Server{
		Addr:         "localhost:8080",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
		log.Fatalf("error: listening and serving: %s", err)
	}
}

// Echo is a basic HTTP Handler.
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)
}

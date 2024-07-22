package main

import (
	"fmt"
	"net/http"
	"os"

	"ascii/handlers"
)

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

	// Determine the port in use
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	fmt.Printf("Running server on port %s...\n", port)

	// start the server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// Handle different url paths.
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/":
		handlers.IndexHandler(w, r)

	case "/ascii-art":
		handlers.AsciiArtHandler(w, r)

	default:
		http.NotFound(w, r)
	}
}

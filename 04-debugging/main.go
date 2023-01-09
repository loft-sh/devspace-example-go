package main

import (
	"fmt"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received.")
		fmt.Fprint(w, "Hello, world!")
	})

	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

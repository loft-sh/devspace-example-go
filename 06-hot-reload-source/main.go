package main

import (
	"fmt"
	"net/http"

	"github.com/loft-sh/devspace-example-go/hello-world/pkg/auth"
	"github.com/loft-sh/devspace-example-go/hello-world/pkg/server"
)

const (
	port = "8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("App %s received request.\n", server.Name)
		fmt.Printf("Expected token %s ...\n", auth.Token)
		fmt.Fprint(w, "Hello, world!")
	})

	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

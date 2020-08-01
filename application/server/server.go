package main

import (
	"fmt"
	"log"
	"net/http"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}

func main() {
	handler := http.HandlerFunc(PlayerServer)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}

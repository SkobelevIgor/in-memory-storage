package main

import (
	"in-memory-storage/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", api.RequestHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

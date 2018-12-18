package api

import (
	"fmt"
	"in-memory-storage/store"
	"net/http"
)

func handleGet(w http.ResponseWriter, r *http.Request) {
	id, err := resolveID(r)
	if err != nil || id == "" {
		http.NotFound(w, r)
		return
	}

	rec := store.GetRecord(id)
	if rec == nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, jsonResponse(rec))
}

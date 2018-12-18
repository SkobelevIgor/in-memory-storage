package api

import (
	"in-memory-storage/store"
	"net/http"
)

func handlePut(w http.ResponseWriter, r *http.Request) {
	id, err := resolveID(r)
	if err != nil || id == "" {
		http.NotFound(w, r)
		return
	}
	js, err := getBodyJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = store.ReplaceRecord(id, js)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}

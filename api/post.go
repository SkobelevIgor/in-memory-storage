package api

import (
	"fmt"
	"in-memory-storage/store"
	"net/http"
)

type newRecordResponse struct {
	ID string `json:"id"`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	js, err := getBodyJSON(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := store.SaveRecord(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := newRecordResponse{ID: id}
	fmt.Fprint(w, jsonResponse(resp))
}

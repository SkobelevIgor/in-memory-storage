package api

import (
	"encoding/json"
	"fmt"
	"in-memory-storage/store"
	"io/ioutil"
	"net/http"
	"strings"
)

const storagePath = "storage"

type newRecordResponse struct {
	ID string `json:"id"`
}

func resolveID(r *http.Request) (id string) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) == 0 || string(path[1]) != "store" {
		return
	}
	if len(path) == 3 {
		id = path[2]
	}
	return
}

func getBodyJSON(r *http.Request) (inp json.RawMessage) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var checkJSON interface{}
	if len(body) > 0 && json.Unmarshal(body, &checkJSON) == nil {
		inp = json.RawMessage(body)
	}

	return
}

func jsonResponse(r interface{}) string {
	json, err := json.Marshal(r)
	if err != nil {
		// @TODO process error
	}
	return string(json)
}

// RequestHandler Handle all requests
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodPut:
		handlePut(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	default:
		http.NotFound(w, r)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id := resolveID(r)
	if id == "" {
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

func handlePost(w http.ResponseWriter, r *http.Request) {
	js := getBodyJSON(r)
	if js == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	id := store.SaveRecord(js)
	if id != "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	resp := newRecordResponse{ID: id}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, jsonResponse(resp))
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	id := resolveID(r)
	if id == "" {
		http.NotFound(w, r)
		return
	}
	js := getBodyJSON(r)
	if js == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if store.ReplaceRecord(id, js) {
		w.WriteHeader(http.StatusOK)
	} else {
		http.NotFound(w, r)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id := resolveID(r)
	if id == "" {
		http.NotFound(w, r)
		return
	}

	if store.DeleteRecord(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.NotFound(w, r)
	}
}

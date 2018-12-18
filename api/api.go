package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const storagePath = "storage"

func resolveID(r *http.Request) (id string, err error) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) == 0 || string(path[1]) != "store" {
		err = errors.New("Wrong route")
		return
	}
	if len(path) == 3 {
		id = path[2]
	}
	return
}

func getBodyJSON(r *http.Request) (json.RawMessage, error) {
	var inp json.RawMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		inp = json.RawMessage(body)
	}

	return inp, err
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

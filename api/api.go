package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"in-memory-storage/store"
	"io/ioutil"
	"net/http"
	"strings"
)

const storagePath = "storage"

func resolveID(path []string) (id string, err error) {
	if len(path) == 0 || string(path[1]) != "store" {
		err = errors.New("Wrong route")
		return
	}
	if len(path) == 3 {
		id = path[2]
	}
	return
}

func jsonResponse(r interface{}) []byte {
	json, err := json.Marshal(r)
	if err != nil {
		// @TODO process error
	}

	return json
}

// RequestHandler Handle all requests
func RequestHandler(w http.ResponseWriter, r *http.Request) {

	id, err := resolveID(
		strings.Split(r.URL.Path, "/"))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println("Error")
	}

	var inp interface{}
	var resp interface{}
	if len(body) > 0 {
		err := json.Unmarshal(body, &inp)
		if err != nil {
			fmt.Println("Could not parse json data")
		}
	}

	switch r.Method {
	case http.MethodGet:
		resp, err = store.GetRecord(id)
	case http.MethodPost:
		resp, err = store.SaveRecord(inp)
		// case http.MethodPut:
		// 	resp, err := store.ReplaceRecord(id, inp)
		// case http.MethodDelete:
		// 	resp, err := store.DeleteRecord(id)
		// default:
		http.NotFound(w, r)

		return
	}

	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte(err))
		// @TODO process error
		fmt.Println(err)

		return
	}

	w.Write(jsonResponse(resp))
}

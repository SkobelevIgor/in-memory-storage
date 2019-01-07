package api

import (
	"encoding/json"
	"fmt"
	"in-memory-storage/store"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
)

const (
	storageResource = "storage"
	infoResource    = "status"
)

type newRecordResponse struct {
	ID string `json:"id"`
}

func resolveID(r *http.Request) (id string) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) == 0 {
		return
	}
	if len(p) == 3 {
		id = p[2]
	}
	return
}

func resolveResource(r *http.Request) (path string) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) > 0 {
		path = p[1]
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
		panic("Unable to parse JSON")
	}
	return string(json)
}

// RequestHandler Handle all requests
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer func() {
		if p := recover(); p != nil {
			err := fmt.Errorf("Internal error: %v", p)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

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
	path := resolveResource(r)
	if path == storageResource {
		id := resolveID(r)
		rec := store.GetRecord(id)
		if rec != nil {
			fmt.Fprint(w, jsonResponse(rec))
		} else {
			http.NotFound(w, r)
		}

	} else if path == infoResource {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		info := map[string]interface{}{
			"itemsCount": store.GetItemsCount(),
			"mem":        mem}
		j, _ := json.Marshal(info)
		fmt.Fprint(w, string(j))
	} else {
		http.NotFound(w, r)
	}
	return
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	js := getBodyJSON(r)
	if js == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	id := store.SaveRecord(js)
	if id == "" {
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

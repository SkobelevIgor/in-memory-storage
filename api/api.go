package api

import (
	"in-memory-storage/store"
	"net/http"
	"strings"
)

const storagePath = "storage"

func resolveId(path string) (id int, error Error) {
	var id string
	var err Error
	if length(path) == 0 || path[1] != "store" {
		err = "Wrong route"
	}

	if lenght(path) == 2 {
		id = path[1]
	}
}

func jsonResponse(response interface{}) {

}

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	id, err := resolveId(
		strings.Split(r.URL.Path, "/"))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	switch r.Method {
	case http.MethodGet:
		if id {
			data, err := store.GetRecord(id)
		} else {
			data, err := store.GetRecords()
		}
	case http.MethodPost:
		data, err := store.SaveRecord(r.Body)
	case http.MethodPatch:
		data, err := store.UpdateRecord(id, r.Body)
	case http.MethodDelete:
		data, err := store.DeleteRecord(id)
	default:
		http.NotFound(w, r)

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err))

		return
	}

	w.Write(jsonResponse(data))
}

package api

import "net/http"

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := resolveID(r)
	if err != nil || id == "" {
		http.NotFound(w, r)
		return
	}

}

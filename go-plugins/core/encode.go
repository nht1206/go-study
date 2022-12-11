package core

import (
	"encoding/json"
	"net/http"
)

func DecodeRequest(w http.ResponseWriter, r *http.Request, v interface{}) (err error) {
	if err = json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	return
}

func EncodeResponse(w http.ResponseWriter, v interface{}, hasErr bool) {
	w.Header().Set("Content-Type", DefaultContentType)
	if hasErr {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(v)
}

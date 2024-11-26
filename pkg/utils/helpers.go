package utils

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Read UID param
func ReadUIDParam(r *http.Request) string {
	params := httprouter.ParamsFromContext(r.Context())
	uid := params.ByName("uid")
	return uid
}

// JSON envelope
type Envelope map[string]interface{}

// Write JSON body
func WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

package handlers

import (
	"encoding/json"
	"net/http"
)

// convertObject marshal a generic object on a []byte return.
func convertObject(object any) (response []byte, err error) {
	if response, err = json.Marshal(&object); err != nil {
		return make([]byte, 0), err
	}
	return response, nil
}

// writeResponse write the response byte input on writer.
func writeResponse(w http.ResponseWriter, object any) error {
	response, err := convertObject(object)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		return err
	}
	return nil
}

// handleError write down an error with code to the writer response.
func handleError(w http.ResponseWriter, code int, err error) (hasError bool) {
	hasError = err != nil
	if hasError {
		http.Error(w, err.Error(), code)
	}
	return hasError
}

// handleError write down an error with code to the writer response.
func writeError(w http.ResponseWriter, code int, err error) {
	http.Error(w, err.Error(), code)
}

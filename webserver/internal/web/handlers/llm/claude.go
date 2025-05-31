package llm

import (
	"encoding/json"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/llm"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

const (
	websocketBufferSize = 1024
)

// RequestBody represents the structure of a request payload.
type RequestBody struct {
	Request string `json:"request"`
}

// HandleClaude processes an HTTP POST request, decodes the request body,
// and communicates with the Claude LLM service.
func HandleClaude(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		system.HandleError(w, http.StatusBadRequest, err)
		return
	}

	client, err := llm.NewClient()
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	response, err := client.SendMessage(r.Context(), reqBody.Request)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, response)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

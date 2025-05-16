package llm

import (
	"encoding/json"
	"net/http"

	"github.com/knabben/observatio/webserver/internal/infra/llm"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

type RequestBody struct {
	Request string `json:"request"`
}

func HandleClaude(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		system.HandleError(w, http.StatusBadRequest, err)
		return
	}

	client := llm.NewClient(reqBody.Request)
	response, err := client.SendMessage(r.Context())
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}

	err = system.WriteResponse(w, response)
	if system.HandleError(w, http.StatusInternalServerError, err) {
		return
	}
}

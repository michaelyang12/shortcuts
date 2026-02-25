package handler

import (
	"encoding/json"
	"net/http"

	"github.com/michaelyang12/shortcuts/claude"
)

type textRequest struct {
	Text   string `json:"text"`
	Prompt string `json:"prompt"`
}

func Text(client *claude.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req textRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.Text == "" || req.Prompt == "" {
			writeError(w, http.StatusBadRequest, "text and prompt are required")
			return
		}

		result, err := client.Text(r.Context(), req.Text, req.Prompt)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeResult(w, result)
	}
}

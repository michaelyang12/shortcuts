package handler

import (
	"encoding/json"
	"net/http"

	"github.com/michaelyang12/shortcuts/claude"
	"github.com/michaelyang12/shortcuts/media"
)

type videoRequest struct {
	URL    string `json:"url"`
	Prompt string `json:"prompt"`
}

func Video(client *claude.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req videoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.URL == "" || req.Prompt == "" {
			writeError(w, http.StatusBadRequest, "url and prompt are required")
			return
		}

		frames, cleanup, err := media.ExtractFrames(req.URL)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer cleanup()

		result, err := client.Vision(r.Context(), frames, req.Prompt)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeResult(w, result)
	}
}

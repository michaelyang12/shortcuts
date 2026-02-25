package handler

import (
	"encoding/json"
	"net/http"

	"github.com/michaelyang12/shortcuts/claude"
)

type imageRequest struct {
	Image     string `json:"image"`      // base64-encoded
	MediaType string `json:"media_type"` // e.g. "image/jpeg", "image/png"
	Prompt    string `json:"prompt"`
}

func Image(client *claude.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req imageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.Image == "" || req.Prompt == "" {
			writeError(w, http.StatusBadRequest, "image and prompt are required")
			return
		}
		if req.MediaType == "" {
			req.MediaType = "image/jpeg"
		}

		result, err := client.VisionBase64(r.Context(), req.Image, req.MediaType, req.Prompt)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeResult(w, result)
	}
}

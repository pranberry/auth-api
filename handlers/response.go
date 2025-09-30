package handlers

import (
	"encoding/json"
	"auth-api/models"
	"net/http"
)

func writeResponse(w http.ResponseWriter, resp models.Response) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		resp.Error = err
		http.Error(w, resp.Message, resp.Status)
	}
	// Add log here
}

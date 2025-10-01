package handlers

import (
	"auth-api/models"
	"encoding/json"
	"net/http"
)

// Shared function to write json out to the
func WriteResponse(w http.ResponseWriter, resp *models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		resp.Error = err
		http.Error(w, resp.Message, resp.Status)
	}
}

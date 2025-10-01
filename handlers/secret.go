package handlers

import (
	"net/http"
)

// Serves the secret file
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/gif")
	w.WriteHeader(http.StatusOK)
	filePath := "assets/hamster_dance.gif"
	http.ServeFile(w, r, filePath)
}

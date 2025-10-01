package handlers

import "net/http"

// Serves the secret file
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	ServeStaticFile(w, r, "assets/hamster_dance.gif", "image/gif")
}

// Simple healths check handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, &Response{
		Message: "alive and well",
		Status:  http.StatusOK,
	})
}

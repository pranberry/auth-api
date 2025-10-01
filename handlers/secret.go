package handlers

import (
	"net/http"
)

// Serves the secret file
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/gif")
	w.WriteHeader(http.StatusOK)
	filePath := "/app/assets/hamster_dance.gif"
	http.ServeFile(w, r, filePath)
}

// Serve secret data
func SecretHandlerTest(w http.ResponseWriter, r *http.Request) {

	anonStruct := []struct{
		Secret string	`json:"secret,omitempty"`
		Num	int			`json:"num,omitempty"`
	}{{
		Secret :"Knock Knock",
		Num : 5,
	},{
		Secret: "Gabbagoul",
		Num: 6,
	}}

	WriteResponse(w, &Response{
		Message: "THIS IS THE super secret STRUCT. guard it with your LIFE!",
		Data: anonStruct,
		Error: nil,
		Status: http.StatusOK,
	})
}

// Simple healths check handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, &Response{
		Message: "alive and well",
		Status:  http.StatusOK,
	})
}

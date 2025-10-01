package handlers

import (
	"encoding/json"
	"net/http"
)

// Helper function to handle early exists (in middleware, etc)
func WriteResponse(w http.ResponseWriter, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		resp.Error = err
		http.Error(w, resp.Message, resp.Status)
	}
}

// Response structs carries some often needed fields for middleware
type Response struct {
	Message string `json:"message"`
	Error   error  `json:"-"`
	// http status of the response
	Status int `json:"-"`
	// JWTResponse only goes of if login was successful
	//Auth *auth.JWTResponse `json:"auth2,omitempty"`
	Data any               `json:"auth,omitempty"`
}

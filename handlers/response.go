package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
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
	Status  int    `json:"-"` // http status of the response
	Data    any    `json:"auth,omitempty"`
}

// ServeStaticFile ensures static files are served with the desired headers and
// graceful error handling via WriteResponse.
func ServeStaticFile(w http.ResponseWriter, r *http.Request, filePath, contentType string) {
	resolvedPath, info, err := resolveStaticFilePath(filePath)
	if err != nil {
		resp := &Response{
			Status:  http.StatusNotFound,
			Message: "requested asset not found",
			Error:   err,
		}
		if !errors.Is(err, os.ErrNotExist) {
			resp.Status = http.StatusInternalServerError
			resp.Message = "failed to access asset"
		}
		WriteResponse(w, resp)
		return
	}
	if info.IsDir() {
		WriteResponse(w, &Response{
			Status:  http.StatusNotFound,
			Message: "requested asset not found",
		})
		return
	}

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	http.ServeFile(w, r, resolvedPath)
}

func resolveStaticFilePath(filePath string) (string, os.FileInfo, error) {
	if filepath.IsAbs(filePath) {
		info, err := os.Stat(filePath)
		return filePath, info, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}

	dir := wd
	for {
		candidate := filepath.Join(dir, filePath)
		info, statErr := os.Stat(candidate)
		if statErr == nil {
			return candidate, info, nil
		}
		if !errors.Is(statErr, os.ErrNotExist) {
			return "", nil, statErr
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", nil, os.ErrNotExist
}

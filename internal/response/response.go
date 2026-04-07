package response

import (
	"encoding/json"
	"net/http"
)


type APIResponse struct {
	Status string `json:"status"`
  Message string `json:"message"`
  Data  any `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, StatusCode int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	json.NewEncoder(w).Encode(resp)
}

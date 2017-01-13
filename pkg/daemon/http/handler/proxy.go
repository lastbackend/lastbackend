package handler

import (
	"encoding/json"
	"net/http"
)

func ProxyToken(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("token")
}

package core

import (
	"encoding/json"
	"net/http"
)

func parseData(r *http.Request, w http.ResponseWriter, v any) bool {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		writeBadResponse(http.StatusBadRequest, w, "Invalid Json")
		return false
	}

	return true
}

func writeBadResponse(code int, w http.ResponseWriter, detail string) {

	(w).Header().Set("Content-Type", "application/json")
	(w).WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "error",
		"detail": detail,
	})

}

func writeGoodResponse(code int, w http.ResponseWriter, detail string, data map[string]any) {

	(w).Header().Set("Content-Type", "application/json")
	(w).WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "status",
		"detail": detail,
		"data":   data,
	})

}

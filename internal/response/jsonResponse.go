package response

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("\n", "    ")
	enc.Encode(payload)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

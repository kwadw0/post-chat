package json

import (
	"encoding/json"
	"net/http"
)


func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ReadJson is our "Safety Inspector" for incoming mail.
func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	// 1. Limit the size of the mail to 1 Megabyte
	// (We don't want to try and read a library!)
	maxBytes := 1_048_576 
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// 2. Setup the Decoder
	dec := json.NewDecoder(r.Body)
	
	// 3. Try to read the mail into our 'data' container
	return dec.Decode(data)
}

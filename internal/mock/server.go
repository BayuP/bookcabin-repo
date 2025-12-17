package mock

import (
	"log"
	"net/http"
	"os"
)

func ServeJSONFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("JSON file not found: %s", path)
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, path)
	}
}

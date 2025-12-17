package mock

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

func MockAirAsiaServer() *http.Server {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	jsonPath := filepath.Join(baseDir, "airasia.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/airasia/search", ServeJSONFile(jsonPath))

	server := &http.Server{
		Addr:    ":8081", // fixed port for curl
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Mock airasia server running at http://127.0.0.1:8081")
	return server
}

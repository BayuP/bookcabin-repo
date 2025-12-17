package mock

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

func MockGarudaServer() *http.Server {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	jsonPath := filepath.Join(baseDir, "garuda.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/garuda/search", ServeJSONFile(jsonPath))

	server := &http.Server{
		Addr:    ":8083", // fixed port for curl
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Mock Garuda server running at http://127.0.0.1:8083")
	return server
}

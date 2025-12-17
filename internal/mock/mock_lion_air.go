package mock

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

func MockLionServer() *http.Server {

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	jsonPath := filepath.Join(baseDir, "lion.json")

	mux := http.NewServeMux()
	mux.HandleFunc("/lion/search", ServeJSONFile(jsonPath))

	server := &http.Server{
		Addr:    ":8084", // fixed port for curl
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Mock Lion server running at http://127.0.0.1:8084")
	return server

}

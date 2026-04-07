package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/ykube-app/demo-app/internal/api"
	"github.com/ykube-app/demo-app/internal/db"
)

//go:embed all:static
var staticFS embed.FS

func main() {
	pool, err := db.Open()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	// API routes — handler internally routes /api/tasks and /api/tasks/{id}
	apiHandler := api.NewHandler(pool)
	mux.Handle("/api/tasks", apiHandler)
	mux.Handle("/api/tasks/", apiHandler)

	// Serve the embedded React app for everything else.
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("static fs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

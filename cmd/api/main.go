package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/findhouse/internal/api"
	"github.com/findhouse/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbPath := flag.String("db", "findhouse.db", "Path to SQLite database")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	// Inicializar DB
	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer database.Close()

	// Crear router y handlers
	r := chi.NewRouter()
	h := api.NewHandler(database)

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rutas
	r.Route("/api", func(r chi.Router) {
		r.Route("/properties", func(r chi.Router) {
			r.Get("/unrated", h.GetUnratedProperties)
			r.Get("/liked", h.GetLikedProperties)
			r.Put("/{id}/rate", h.RateProperty)
		})
	})

	// Iniciar servidor
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

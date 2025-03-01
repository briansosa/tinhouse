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
	"github.com/go-chi/cors"
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

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://*:5173"}, // URL de desarrollo de Vite
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	h := api.NewHandler(database)

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rutas
	r.Route("/api", func(r chi.Router) {
		r.Get("/properties/unrated", h.GetUnratedProperties)
		r.Get("/properties/liked", h.GetLikedProperties)
		r.Get("/properties/favorites", h.GetFavoriteProperties)
		r.Put("/properties/{id}/rate", h.RateProperty)
		r.Put("/properties/{id}/favorite", h.TogglePropertyFavorite)

		// Rutas para notas de propiedades
		r.Get("/properties/{id}/notes", h.GetPropertyNotes)
		r.Post("/properties/{id}/notes", h.AddPropertyNote)
		r.Delete("/properties/notes/{noteId}", h.DeletePropertyNote)
	})

	// Iniciar servidor
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

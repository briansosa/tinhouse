package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/findhouse/internal/db"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{db: db}
}

// PropertyResponse es la estructura de respuesta para las propiedades
type PropertyResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Code         string    `json:"code"`
	Price        string    `json:"price"`
	Location     string    `json:"location"`
	PropertyType string    `json:"property_type,omitempty"`
	ImageURL     string    `json:"image_url"`
	Images       []string  `json:"images,omitempty"`
	URL          string    `json:"url"`
	Description  string    `json:"description,omitempty"`
	LastUpdated  time.Time `json:"last_updated"`
	CreatedAt    time.Time `json:"created_at"`
	Details      Details   `json:"details"`
	Agency       Agency    `json:"agency"`
}

type Details struct {
	Bedrooms  *int     `json:"bedrooms,omitempty"`
	Bathrooms *int     `json:"bathrooms,omitempty"`
	Area      *float64 `json:"area,omitempty"`
	TotalArea *float64 `json:"total_area,omitempty"`
	LandArea  *float64 `json:"land_area,omitempty"`
	Rooms     *int     `json:"rooms,omitempty"`
	Floors    *int     `json:"floors,omitempty"`
	Garages   *int     `json:"garages,omitempty"`
	Status    *string  `json:"status,omitempty"`
	Expenses  *float64 `json:"expenses,omitempty"`
	Age       *string  `json:"age,omitempty"`
	FrontSize *float64 `json:"front_size,omitempty"`
	BackSize  *float64 `json:"back_size,omitempty"`
}

type Agency struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetUnratedProperties retorna las propiedades sin calificar
func (h *Handler) GetUnratedProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.db.GetUnratedProperties()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting unrated properties: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]PropertyResponse, 0, len(properties))
	for _, p := range properties {
		response = append(response, h.toPropertyResponse(&p))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": response,
		"total":      len(response),
	})
}

// GetLikedProperties retorna las propiedades con like
func (h *Handler) GetLikedProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.db.GetLikedProperties()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting liked properties: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]PropertyResponse, 0, len(properties))
	for _, p := range properties {
		response = append(response, h.toPropertyResponse(&p))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": response,
		"total":      len(response),
	})
}

// RateProperty califica una propiedad
func (h *Handler) RateProperty(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	var request struct {
		Rating string `json:"rating"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.db.RateProperty(propertyID, request.Rating); err != nil {
		http.Error(w, fmt.Sprintf("error rating property: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"property_id": propertyID,
		"rating":      request.Rating,
	})
}

// Helper para limpiar el precio
func cleanPrice(price string) string {
	// Eliminar saltos de línea y texto extra
	if idx := strings.Index(price, "\n"); idx != -1 {
		price = price[:idx]
	}
	return strings.TrimSpace(price)
}

// Helper para convertir Propiedad a PropertyResponse
func (h *Handler) toPropertyResponse(p *db.Propiedad) PropertyResponse {
	// Obtener inmobiliaria si es necesario
	var agency Agency
	if p.InmobiliariaID > 0 {
		if i, err := h.db.GetInmobiliariaByID(p.InmobiliariaID); err == nil {
			agency = Agency{
				ID:   i.ID,
				Name: i.Nombre,
			}
		}
	}

	// Convertir las imágenes
	var images []string
	if p.Imagenes != nil {
		images = *p.Imagenes
	}

	return PropertyResponse{
		ID:           p.ID,
		Title:        p.Titulo,
		Code:         p.Codigo,
		Price:        cleanPrice(p.Precio),
		Location:     getLocation(p),
		PropertyType: getString(p.TipoPropiedad),
		ImageURL:     p.ImagenURL,
		Images:       images,
		URL:          p.URL,
		Description:  getString(p.Descripcion),
		LastUpdated:  p.FechaScraping,
		CreatedAt:    p.CreatedAt,
		Details: Details{
			Bedrooms:  p.Dormitorios,
			Bathrooms: p.Banios,
			Area:      p.SuperficieCubierta,
			TotalArea: p.SuperficieTotal,
			LandArea:  p.SuperficieTerreno,
			Rooms:     p.Ambientes,
			Floors:    p.Plantas,
			Garages:   p.Cocheras,
			Status:    p.Situacion,
			Expenses:  p.Expensas,
			Age:       p.Antiguedad,
			FrontSize: p.Frente,
			BackSize:  p.Fondo,
		},
		Agency: agency,
	}
}

// Helper para obtener la ubicación completa
func getLocation(p *db.Propiedad) string {
	if p.Direccion != "" {
		return p.Direccion
	}
	if p.Ubicacion != nil {
		return *p.Ubicacion
	}
	return ""
}

// Helper para manejar punteros string
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

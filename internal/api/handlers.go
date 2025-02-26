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
	HasNotes     bool      `json:"has_notes"`
	IsFavorite   bool      `json:"is_favorite"`
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

// Note representa una nota de propiedad
type Note struct {
	ID         int64     `json:"id"`
	PropertyID int64     `json:"property_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

// GetUnratedProperties retorna las propiedades sin calificar
func (h *Handler) GetUnratedProperties(w http.ResponseWriter, r *http.Request) {
	// Parsear filtros de la solicitud
	filter, err := parsePropertyFilter(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing filters: %v", err), http.StatusBadRequest)
		return
	}

	properties, err := h.db.GetUnratedProperties(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting unrated properties: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]PropertyResponse, 0, len(properties))
	for _, p := range properties {
		resp := h.toPropertyResponse(&p)

		// Verificar si la propiedad tiene notas
		hasNotes, err := h.db.PropertyHasNotes(p.ID)
		if err == nil {
			resp.HasNotes = hasNotes
		}

		// Verificar si la propiedad es favorita
		isFavorite, err := h.db.IsPropertyFavorite(p.ID)
		if err == nil {
			resp.IsFavorite = isFavorite
		}

		response = append(response, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": response,
		"total":      len(response),
	})
}

// GetLikedProperties retorna las propiedades con like
func (h *Handler) GetLikedProperties(w http.ResponseWriter, r *http.Request) {
	// Parsear filtros de la solicitud
	filter, err := parsePropertyFilter(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing filters: %v", err), http.StatusBadRequest)
		return
	}

	properties, err := h.db.GetLikedProperties(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting liked properties: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]PropertyResponse, 0, len(properties))
	for _, p := range properties {
		resp := h.toPropertyResponse(&p)

		// Verificar si la propiedad tiene notas
		hasNotes, err := h.db.PropertyHasNotes(p.ID)
		if err == nil {
			resp.HasNotes = hasNotes
		}

		// Verificar si la propiedad es favorita
		isFavorite, err := h.db.IsPropertyFavorite(p.ID)
		if err == nil {
			resp.IsFavorite = isFavorite
		}

		response = append(response, resp)
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

// GetPropertyNotes obtiene las notas de una propiedad
func (h *Handler) GetPropertyNotes(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	notes, err := h.db.GetPropertyNotes(propertyID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting property notes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notes": notes,
	})
}

// AddPropertyNote agrega una nota a una propiedad
func (h *Handler) AddPropertyNote(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	var request struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Text == "" {
		http.Error(w, "note text cannot be empty", http.StatusBadRequest)
		return
	}

	note := &db.PropertyNote{
		PropertyID: propertyID,
		Text:       request.Text,
	}

	if err := h.db.AddPropertyNote(note); err != nil {
		http.Error(w, fmt.Sprintf("error adding property note: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"note":    note,
	})
}

// DeletePropertyNote elimina una nota de una propiedad
func (h *Handler) DeletePropertyNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := strconv.ParseInt(chi.URLParam(r, "noteId"), 10, 64)
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeletePropertyNote(noteID); err != nil {
		http.Error(w, fmt.Sprintf("error deleting property note: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"note_id": noteID,
	})
}

// TogglePropertyFavorite marca o desmarca una propiedad como favorita
func (h *Handler) TogglePropertyFavorite(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	var request struct {
		IsFavorite bool `json:"is_favorite"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.db.TogglePropertyFavorite(propertyID, request.IsFavorite); err != nil {
		http.Error(w, fmt.Sprintf("error toggling favorite: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"property_id": propertyID,
		"is_favorite": request.IsFavorite,
	})
}

// GetFavoriteProperties retorna las propiedades marcadas como favoritas
func (h *Handler) GetFavoriteProperties(w http.ResponseWriter, r *http.Request) {
	// Parsear filtros de la solicitud
	filter, err := parsePropertyFilter(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing filters: %v", err), http.StatusBadRequest)
		return
	}

	properties, err := h.db.GetFavoriteProperties(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting favorite properties: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]PropertyResponse, 0, len(properties))
	for _, p := range properties {
		resp := h.toPropertyResponse(&p)

		// Verificar si la propiedad tiene notas
		hasNotes, err := h.db.PropertyHasNotes(p.ID)
		if err == nil {
			resp.HasNotes = hasNotes
		}

		// Las propiedades favoritas siempre tienen IsFavorite = true
		resp.IsFavorite = true

		response = append(response, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"properties": response,
		"total":      len(response),
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
		Agency:     agency,
		IsFavorite: p.IsFavorite,
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

// parsePropertyFilter parsea los filtros de la solicitud
func parsePropertyFilter(r *http.Request) (*db.PropertyFilter, error) {
	// Si es GET, parsear de query params
	if r.Method == http.MethodGet {
		return parseFilterFromQueryParams(r)
	}

	// Si es POST, parsear del body
	if r.Method == http.MethodPost {
		var filter db.PropertyFilter
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			return nil, fmt.Errorf("error decoding filter: %v", err)
		}
		return &filter, nil
	}

	// Por defecto, retornar un filtro vacío
	return &db.PropertyFilter{}, nil
}

// parseFilterFromQueryParams parsea los filtros de los query params
func parseFilterFromQueryParams(r *http.Request) (*db.PropertyFilter, error) {
	filter := &db.PropertyFilter{
		Currency: "ARS", // Valor por defecto
	}

	// Tipo de propiedad
	if propertyType := r.URL.Query().Get("property_type"); propertyType != "" {
		filter.PropertyType = propertyType
	}

	// Ubicaciones
	if locations := r.URL.Query().Get("locations"); locations != "" {
		filter.Locations = strings.Split(locations, ",")
	}

	// Características
	if features := r.URL.Query().Get("features"); features != "" {
		filter.Features = strings.Split(features, ",")
	}

	// Precio mínimo
	if priceMin := r.URL.Query().Get("price_min"); priceMin != "" {
		min, err := strconv.ParseFloat(priceMin, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price_min: %v", err)
		}
		filter.PriceMin = &min
	}

	// Precio máximo
	if priceMax := r.URL.Query().Get("price_max"); priceMax != "" {
		max, err := strconv.ParseFloat(priceMax, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price_max: %v", err)
		}
		filter.PriceMax = &max
	}

	// Moneda
	if currency := r.URL.Query().Get("currency"); currency != "" {
		filter.Currency = currency
	}

	// Tamaño mínimo
	if sizeMin := r.URL.Query().Get("size_min"); sizeMin != "" {
		min, err := strconv.ParseFloat(sizeMin, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid size_min: %v", err)
		}
		filter.SizeMin = &min
	}

	// Tamaño máximo
	if sizeMax := r.URL.Query().Get("size_max"); sizeMax != "" {
		max, err := strconv.ParseFloat(sizeMax, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid size_max: %v", err)
		}
		filter.SizeMax = &max
	}

	// Ambientes
	if rooms := r.URL.Query().Get("rooms"); rooms != "" {
		roomsInt, err := strconv.Atoi(rooms)
		if err != nil {
			return nil, fmt.Errorf("invalid rooms: %v", err)
		}
		filter.Rooms = &roomsInt
	}

	// Baños
	if bathrooms := r.URL.Query().Get("bathrooms"); bathrooms != "" {
		bathroomsInt, err := strconv.Atoi(bathrooms)
		if err != nil {
			return nil, fmt.Errorf("invalid bathrooms: %v", err)
		}
		filter.Bathrooms = &bathroomsInt
	}

	// Antigüedad
	if antiquity := r.URL.Query().Get("antiquity"); antiquity != "" {
		antiquityInt, err := strconv.Atoi(antiquity)
		if err != nil {
			return nil, fmt.Errorf("invalid antiquity: %v", err)
		}
		filter.Antiquity = &antiquityInt
	}

	// Solo con notas
	if showOnlyWithNotes := r.URL.Query().Get("show_only_with_notes"); showOnlyWithNotes == "true" {
		filter.ShowOnlyWithNotes = true
	}

	return filter, nil
}

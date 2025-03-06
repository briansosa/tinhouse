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
	ID           int64               `json:"id"`
	Title        string              `json:"title"`
	Code         string              `json:"code"`
	Price        string              `json:"price"`
	Location     string              `json:"location"`
	PropertyType string              `json:"property_type,omitempty"`
	ImageURL     string              `json:"image_url"`
	Images       []string            `json:"images,omitempty"`
	URL          string              `json:"url"`
	Description  string              `json:"description,omitempty"`
	LastUpdated  time.Time           `json:"last_updated"`
	CreatedAt    time.Time           `json:"created_at"`
	Details      Details             `json:"details"`
	Agency       Agency              `json:"agency"`
	HasNotes     bool                `json:"has_notes"`
	IsFavorite   bool                `json:"is_favorite"`
	Features     map[string][]string `json:"features,omitempty"`
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
	Age       *int     `json:"age,omitempty"`
	FrontSize *float64 `json:"front_size,omitempty"`
	BackSize  *float64 `json:"back_size,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
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

	// Verificar si la propiedad tiene notas
	hasNotes, _ := h.db.PropertyHasNotes(p.ID)

	// Obtener características de la propiedad
	features, _ := h.db.GetPropertyFeaturesAsMap(p.ID)

	// Obtener el tipo de propiedad
	var propertyType string
	if p.TipoPropiedad != nil {
		// Obtener el código del tipo de propiedad
		code := h.getPropertyTypeCodeById(*p.TipoPropiedad)

		// Obtener el nombre del tipo de propiedad a partir del código
		if name, err := h.db.GetPropertyTypeNameByCode(code); err == nil {
			propertyType = name
		} else {
			// Si hay error, usar el código como fallback
			propertyType = code
		}
	}

	return PropertyResponse{
		ID:           p.ID,
		Title:        p.Titulo,
		Code:         p.Codigo,
		Price:        cleanPrice(p.Precio),
		Location:     getLocation(p),
		PropertyType: propertyType,
		ImageURL:     p.ImagenURL,
		Images:       images,
		URL:          p.URL,
		Description:  getString(p.Descripcion),
		LastUpdated:  p.UpdatedAt,
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
			Latitude:  p.Latitud,
			Longitude: p.Longitud,
		},
		Agency:     agency,
		HasNotes:   hasNotes,
		IsFavorite: p.IsFavorite,
		Features:   features,
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
	filter := &db.PropertyFilter{}

	// Tipo de propiedad (por código)
	if propertyType := r.URL.Query().Get("property_type"); propertyType != "" {
		filter.PropertyType = propertyType
		fmt.Printf("Filtro property_type (código): %s\n", propertyType)
	}

	// Tipo de propiedad (por ID)
	if propertyTypeID := r.URL.Query().Get("property_type_id"); propertyTypeID != "" {
		id, err := strconv.ParseInt(propertyTypeID, 10, 64)
		if err != nil {
			fmt.Printf("Error al convertir property_type_id a int64: %v\n", err)
		} else {
			filter.PropertyTypeID = &id
			fmt.Printf("Filtro property_type_id: %d\n", id)
		}
	}

	// Múltiples tipos de propiedad (por IDs)
	if propertyTypeIDs := r.URL.Query().Get("property_type_ids"); propertyTypeIDs != "" {
		idStrings := strings.Split(propertyTypeIDs, ",")
		var ids []int64

		for _, idStr := range idStrings {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				fmt.Printf("Error al convertir un ID de property_type_ids a int64: %v\n", err)
				continue
			}
			ids = append(ids, id)
		}

		if len(ids) > 0 {
			filter.PropertyTypeIDs = ids
			fmt.Printf("Filtro property_type_ids: %v\n", ids)
		}
	}

	// Ubicaciones
	if locations := r.URL.Query().Get("locations"); locations != "" {
		filter.Locations = strings.Split(locations, ",")
		fmt.Printf("Filtro locations: %v\n", filter.Locations)
	}

	// Características - manejar múltiples formatos posibles
	// 1. Formato de array: features[]=1&features[]=2
	// 2. Formato de cadena separada por comas: features=1,2
	features := r.URL.Query()["features[]"] // Intentar obtener como array
	if len(features) == 0 {
		// Si no hay resultados, intentar como parámetro simple
		if featuresStr := r.URL.Query().Get("features"); featuresStr != "" {
			features = strings.Split(featuresStr, ",")
		}
	}

	// Limpiar valores vacíos
	var cleanFeatures []string
	for _, f := range features {
		if f != "" {
			cleanFeatures = append(cleanFeatures, f)
		}
	}

	if len(cleanFeatures) > 0 {
		filter.Features = cleanFeatures
		fmt.Printf("Filtro features: %v\n", filter.Features)
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

	// Moneda - solo establecer si se proporciona explícitamente o si hay filtros de precio
	if currency := r.URL.Query().Get("currency"); currency != "" {
		filter.Currency = currency
	} else if filter.PriceMin != nil || filter.PriceMax != nil {
		// Si hay filtros de precio pero no se especificó moneda, usar ARS por defecto
		filter.Currency = "ARS"
	}

	// Tamaño mínimo (compatibilidad con versión anterior)
	if sizeMin := r.URL.Query().Get("size_min"); sizeMin != "" {
		min, err := strconv.ParseFloat(sizeMin, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid size_min: %v", err)
		}
		filter.SizeMin = &min
	}

	// Tamaño máximo (compatibilidad con versión anterior)
	if sizeMax := r.URL.Query().Get("size_max"); sizeMax != "" {
		max, err := strconv.ParseFloat(sizeMax, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid size_max: %v", err)
		}
		filter.SizeMax = &max
	}

	// Superficie Total mínima
	if totalAreaMin := r.URL.Query().Get("total_area_min"); totalAreaMin != "" {
		min, err := strconv.ParseFloat(totalAreaMin, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total_area_min: %v", err)
		}
		filter.TotalAreaMin = &min
	}

	// Superficie Total máxima
	if totalAreaMax := r.URL.Query().Get("total_area_max"); totalAreaMax != "" {
		max, err := strconv.ParseFloat(totalAreaMax, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total_area_max: %v", err)
		}
		filter.TotalAreaMax = &max
	}

	// Superficie Cubierta mínima
	if coveredAreaMin := r.URL.Query().Get("covered_area_min"); coveredAreaMin != "" {
		min, err := strconv.ParseFloat(coveredAreaMin, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid covered_area_min: %v", err)
		}
		filter.CoveredAreaMin = &min
	}

	// Superficie Cubierta máxima
	if coveredAreaMax := r.URL.Query().Get("covered_area_max"); coveredAreaMax != "" {
		max, err := strconv.ParseFloat(coveredAreaMax, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid covered_area_max: %v", err)
		}
		filter.CoveredAreaMax = &max
	}

	// Superficie Terreno mínima
	if landAreaMin := r.URL.Query().Get("land_area_min"); landAreaMin != "" {
		min, err := strconv.ParseFloat(landAreaMin, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid land_area_min: %v", err)
		}
		filter.LandAreaMin = &min
	}

	// Superficie Terreno máxima
	if landAreaMax := r.URL.Query().Get("land_area_max"); landAreaMax != "" {
		max, err := strconv.ParseFloat(landAreaMax, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid land_area_max: %v", err)
		}
		filter.LandAreaMax = &max
	}

	// Frente
	if front := r.URL.Query().Get("front"); front != "" {
		frontValue, err := strconv.ParseFloat(front, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid front: %v", err)
		}
		filter.Front = &frontValue
	}

	// Fondo
	if back := r.URL.Query().Get("back"); back != "" {
		backValue, err := strconv.ParseFloat(back, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid back: %v", err)
		}
		filter.Back = &backValue
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

	// Disposición
	if disposition := r.URL.Query().Get("disposition"); disposition != "" {
		filter.Disposition = strings.Split(disposition, ",")
	}

	// Orientación
	if orientation := r.URL.Query().Get("orientation"); orientation != "" {
		filter.Orientation = strings.Split(orientation, ",")
	}

	// Solo con notas
	if showOnlyWithNotes := r.URL.Query().Get("show_only_with_notes"); showOnlyWithNotes == "true" {
		filter.ShowOnlyWithNotes = true
	}

	return filter, nil
}

// GetAvailableFeatures retorna todas las características disponibles agrupadas por categoría
func (h *Handler) GetAvailableFeatures(w http.ResponseWriter, r *http.Request) {
	features, err := h.db.GetAllFeatures()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo características: %v", err), http.StatusInternalServerError)
		return
	}

	// Agrupar características por categoría
	categorizedFeatures := make(map[string][]map[string]interface{})
	for _, feature := range features {
		if _, exists := categorizedFeatures[feature.Category]; !exists {
			categorizedFeatures[feature.Category] = []map[string]interface{}{}
		}

		featureMap := map[string]interface{}{
			"id":   feature.ID, // Usar el ID numérico
			"name": feature.Name,
		}
		categorizedFeatures[feature.Category] = append(categorizedFeatures[feature.Category], featureMap)
	}

	// Convertir a formato esperado por el frontend
	var response []map[string]interface{}
	for category, features := range categorizedFeatures {
		categoryMap := map[string]interface{}{
			"id":       strings.ToLower(category),
			"name":     strings.ToUpper(category),
			"features": features,
		}
		response = append(response, categoryMap)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": response,
	})
}

// GetPropertyTypes devuelve los tipos de propiedades disponibles
func (h *Handler) GetPropertyTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.db.GetAllPropertyTypes()
	if err != nil {
		fmt.Printf("Error al obtener tipos de propiedad: %v\n", err)
		http.Error(w, fmt.Sprintf("Error al obtener tipos de propiedad: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(types)
}

// GetListValues devuelve los valores de una lista específica
func (h *Handler) GetListValues(w http.ResponseWriter, r *http.Request) {
	listName := chi.URLParam(r, "listName")
	if listName == "" {
		http.Error(w, "Nombre de lista no especificado", http.StatusBadRequest)
		return
	}

	values, err := h.db.GetListValuesByName(listName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener valores de la lista %s: %v", listName, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)
}

// getPropertyTypeCode obtiene el código del tipo de propiedad
func (h *Handler) getPropertyTypeCode(propertyType string) string {
	if propertyType == "" {
		return ""
	}

	// Consultar en la base de datos el código correspondiente al nombre
	query := `SELECT code FROM property_types WHERE name = ?`
	var code string

	err := h.db.QueryRow(query, propertyType).Scan(&code)
	if err != nil {
		// Si hay error, intentamos hacer una coincidencia aproximada
		fmt.Printf("Error al obtener código para tipo de propiedad '%s': %v\n", propertyType, err)

		// Mapeo básico para compatibilidad
		switch strings.ToLower(propertyType) {
		case "casa", "chalet", "casa chalet", "casa quinta", "quinta":
			return "house"
		case "departamento", "depto", "dpto", "dpto.", "depto.", "departamento con dependencia":
			return "apartment"
		case "ph", "p.h.", "p.h", "propiedad horizontal":
			return "ph"
		case "local", "local comercial", "fondo de comercio":
			return "local"
		case "oficina", "consultorio", "estudio":
			return "office"
		case "terreno", "lote", "lote de terreno", "fracción", "fraccion":
			return "land"
		case "galpón", "galpon", "depósito", "deposito", "nave industrial":
			return "warehouse"
		default:
			// Si no hay coincidencia, devolvemos el tipo en minúsculas como fallback
			return strings.ToLower(strings.ReplaceAll(propertyType, " ", "_"))
		}
	}

	return code
}

// getPropertyTypeCodeById obtiene el código del tipo de propiedad a partir del ID
func (h *Handler) getPropertyTypeCodeById(propertyType int64) string {
	// Consultar en la base de datos el código correspondiente al ID
	query := `SELECT code FROM property_types WHERE id = ?`
	var code string

	err := h.db.QueryRow(query, propertyType).Scan(&code)
	if err != nil {
		// Si hay error, intentamos hacer una coincidencia aproximada
		fmt.Printf("Error al obtener código para tipo de propiedad ID %d: %v\n", propertyType, err)

		// Mapeo básico para compatibilidad
		switch propertyType {
		case 1:
			return "house"
		case 2:
			return "apartment"
		case 3:
			return "ph"
		case 4:
			return "local"
		case 5:
			return "office"
		case 6:
			return "land"
		case 7:
			return "warehouse"
		default:
			// Si no hay coincidencia, devolvemos el ID como fallback
			return fmt.Sprintf("type_%d", propertyType)
		}
	}

	return code
}

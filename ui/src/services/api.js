import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getUnratedProperties = (filters = null) => {
  if (filters) {
    console.log("filters", filters);
    return api.get('/properties/unrated', { params: convertFilters(filters) });
  }
  return api.get('/properties/unrated');
};

export const getLikedProperties = (filters = null) => {
  if (filters) {
    console.log("filters", filters);
    return api.get('/properties/liked', { params: convertFilters(filters) });
  }
  return api.get('/properties/liked');
};

export const rateProperty = (id, rating) => api.put(`/properties/${id}/rate`, { rating });

// Función para marcar una propiedad como dislike
export const dislikeProperty = (id) => api.put(`/properties/${id}/rate`, { rating: 'dislike' });

// Función para marcar/desmarcar una propiedad como favorita
export const togglePropertyFavorite = (id, isFavorite) => api.put(`/properties/${id}/favorite`, { is_favorite: isFavorite });

// Función para obtener propiedades favoritas
export const getFavoriteProperties = (filters = null) => {
  if (filters) {
    console.log("filters", filters);
    return api.get('/properties/favorites', { params: convertFilters(filters) });
  }
  return api.get('/properties/favorites');
};

// Métodos para manejar notas de propiedades
export const getPropertyNotes = (propertyId) => api.get(`/properties/${propertyId}/notes`);

export const addPropertyNote = (propertyId, text) => api.post(`/properties/${propertyId}/notes`, { text });

export const deletePropertyNote = (noteId) => api.delete(`/properties/notes/${noteId}`);

// Función para convertir los filtros del frontend al formato del backend
const convertFilters = (filters) => {
  const backendFilters = {};

  if (filters.propertyType) {
    backendFilters.property_type = filters.propertyType;
  }

  if (filters.locations) {
    backendFilters.locations = filters.locations;
  }

  if (filters.features) {
    backendFilters.features = filters.features;
  }

  if (filters.priceRange?.currency) {
    backendFilters.currency = filters.priceRange.currency;
  }

  if (filters.showOnlyWithNotes) {
    backendFilters.show_only_with_notes = filters.showOnlyWithNotes;
  }

  if (filters.showOnlyFavorites) {
    backendFilters.show_only_favorites = filters.showOnlyFavorites;
  }

  // Convertir rango de precios
  if (filters.priceRange?.min) {
    backendFilters.price_min = filters.priceRange.min;
  }
  if (filters.priceRange?.max) {
    backendFilters.price_max = filters.priceRange.max;
  }

  // Convertir rango de tamaño
  if (filters.sizeRange?.min) {
    backendFilters.size_min = filters.sizeRange.min;
  }
  if (filters.sizeRange?.max) {
    backendFilters.size_max = filters.sizeRange.max;
  }

  // Convertir ambientes, baños y antigüedad
  if (filters.rooms) {
    backendFilters.rooms = filters.rooms;
  }
  if (filters.bathrooms) {
    backendFilters.bathrooms = filters.bathrooms;
  }
  if (filters.antiquity !== null && filters.antiquity !== undefined) {
    backendFilters.antiquity = filters.antiquity;
  }

  return backendFilters;
};

export default api;

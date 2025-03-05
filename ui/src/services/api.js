import axios from 'axios';

const api = axios.create({
  // baseURL: 'http://192.168.0.190:8080/api',
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getUnratedProperties = (filters = null, cancelToken = null) => {
  const config = {};
  
  if (cancelToken) {
    config.cancelToken = cancelToken;
  }
  
  if (filters) {
    config.params = convertFilters(filters);
    return api.get('/properties/unrated', config);
  }
  return api.get('/properties/unrated', config);
};

export const getLikedProperties = (filters = null, cancelToken = null) => {
  const config = {};
  
  if (cancelToken) {
    config.cancelToken = cancelToken;
  }
  
  if (filters) {
    config.params = convertFilters(filters);
    return api.get('/properties/liked', config);
  }
  return api.get('/properties/liked', config);
};

export const rateProperty = (id, rating) => api.put(`/properties/${id}/rate`, { rating });

// Función para marcar una propiedad como dislike
export const dislikeProperty = (id) => api.put(`/properties/${id}/rate`, { rating: 'dislike' });

// Función para marcar/desmarcar una propiedad como favorita
export const togglePropertyFavorite = (id, isFavorite) => api.put(`/properties/${id}/favorite`, { is_favorite: isFavorite });

// Función para obtener propiedades favoritas
export const getFavoriteProperties = (filters = null, cancelToken = null) => {
  const config = {};
  
  if (cancelToken) {
    config.cancelToken = cancelToken;
  }
  
  if (filters) {
    config.params = convertFilters(filters);
    return api.get('/properties/favorites', config);
  }
  return api.get('/properties/favorites', config);
};

// Métodos para manejar notas de propiedades
export const getPropertyNotes = (propertyId) => api.get(`/properties/${propertyId}/notes`);

export const addPropertyNote = (propertyId, text) => api.post(`/properties/${propertyId}/notes`, { text });

export const deletePropertyNote = (noteId) => api.delete(`/properties/notes/${noteId}`);

// Función para convertir los filtros del frontend al formato del backend
const convertFilters = (filters) => {
  if (!filters) return {};
  
  const backendFilters = {};

  if (filters.propertyType && filters.propertyType !== 'all') {
    backendFilters.property_type_id = filters.propertyType;
  }

  if (filters.locations && filters.locations.length > 0) {
    backendFilters.locations = filters.locations.join(',');
  }

  if (filters.features && filters.features.length > 0) {
    backendFilters.features = filters.features.join(',');
  }

  // Solo enviar la moneda si hay filtros de precio
  if ((filters.priceRange?.min || filters.priceRange?.max) && filters.priceRange?.currency) {
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

// Función para obtener todas las características disponibles
export const getAvailableFeatures = () => api.get('/features');

// Obtener todos los tipos de propiedad disponibles
export const getPropertyTypes = () => {
  return api.get('/property-types')
    .then(response => {
      return response;
    })
    .catch(error => {
      console.error("Error en getPropertyTypes:", error);
      throw error;
    });
};

export default api;

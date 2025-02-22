import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getUnratedProperties = () => api.get('/properties/unrated');
export const getLikedProperties = () => api.get('/properties/liked');
export const rateProperty = (id, rating) => api.put(`/properties/${id}/rate`, { rating });

export default api;

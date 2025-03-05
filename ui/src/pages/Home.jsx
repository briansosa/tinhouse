import { useEffect, useState, useRef } from 'react';
import { getUnratedProperties } from '../services/api';
import PropertyCard from '../components/PropertyCard/PropertyCard';
import axios from 'axios';
import api from '../services/api';

export default function Home({ setShowNavBar }) {
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);
  const [history, setHistory] = useState([]); // Historial de propiedades descartadas
  const [globalFilters, setGlobalFilters] = useState(null);
  const isInitialMount = useRef(true);
  const activeRequest = useRef(null);

  useEffect(() => {
    // Cargar los filtros globales del localStorage
    const savedFilters = localStorage.getItem('globalFilters');
    if (savedFilters) {
      try {
        const parsedFilters = JSON.parse(savedFilters);
        console.log("Filtros cargados del localStorage:", parsedFilters);
        setGlobalFilters(parsedFilters);
        // No cargar propiedades aquí, se hará en el siguiente useEffect
      } catch (error) {
        console.error("Error al parsear filtros:", error);
        // Si hay error, cargar sin filtros
        loadProperties(null);
      }
    } else {
      // Solo cargar propiedades si no hay filtros guardados
      loadProperties(null);
    }

    // Limpiar al desmontar
    return () => {
      if (activeRequest.current) {
        activeRequest.current.cancel();
      }
    };
  }, []);

  // Efecto para recargar propiedades cuando cambien los filtros globales
  useEffect(() => {
    // Evitar la carga en el montaje inicial si ya hay filtros
    if (isInitialMount.current) {
      isInitialMount.current = false;
      if (globalFilters) {
        console.log("Cargando propiedades con filtros iniciales:", globalFilters);
        loadProperties(globalFilters);
      }
      return;
    }

    if (globalFilters) {
      console.log("Filtros actualizados, recargando propiedades:", globalFilters);
      loadProperties(globalFilters);
    }
  }, [globalFilters]);

  // Efecto para detectar cambios en los filtros globales en localStorage
  useEffect(() => {
    const handleStorageChange = () => {
      const savedFilters = localStorage.getItem('globalFilters');
      if (savedFilters) {
        setGlobalFilters(JSON.parse(savedFilters));
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => {
      window.removeEventListener('storage', handleStorageChange);
    };
  }, []);

  const loadProperties = async (filters) => {
    setLoading(true);
    
    // Cancelar cualquier solicitud activa
    if (activeRequest.current) {
      activeRequest.current.cancel();
    }
    
    try {
      console.log("Cargando propiedades con filtros:", filters);
      const source = axios.CancelToken.source(); // Usar axios directamente
      activeRequest.current = source;
      
      const response = await getUnratedProperties(filters, source.token);
      console.log("Respuesta recibida:", response.data);
      
      // Solo actualizar si esta es la solicitud más reciente
      if (activeRequest.current === source) {
        setProperties(response.data.properties || []);
        activeRequest.current = null;
      }
    } catch (error) {
      if (!axios.isCancel(error)) {
        console.error('Error fetching properties:', error);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleRate = () => {
    // Guardamos la propiedad actual en el historial antes de descartarla
    setHistory(prev => [...prev, properties[0]]);
    // Removemos la propiedad actual
    setProperties(prev => prev.slice(1));
  };

  const handleUndo = () => {
    if (history.length === 0) return;
    
    // Recuperamos la última propiedad del historial
    const lastProperty = history[history.length - 1];
    
    // La agregamos al inicio de las propiedades activas
    setProperties(prev => [lastProperty, ...prev]);
    
    // La removemos del historial
    setHistory(prev => prev.slice(0, -1));
  };

  return (
    <div className="h-full flex flex-col w-full">
      {loading && properties.length === 0 ? (
        <div className="flex justify-center items-center h-full w-full">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      ) : properties.length > 0 ? (
        <PropertyCard 
          property={properties[0]}
          onRate={handleRate}
          onUndo={handleUndo}
          canUndo={history.length > 0}
          setShowNavBar={setShowNavBar}
        />
      ) : (
        <div className="flex-1 flex flex-col items-center justify-center text-center text-gray-600 w-full p-4">
          <p className="text-xl mb-4">No more properties to show</p>
          <p className="text-sm mb-8">Try adjusting your filters to see more properties</p>
          <button 
            onClick={() => {
              // Limpiar filtros y recargar
              localStorage.removeItem('globalFilters');
              setGlobalFilters(null);
              loadProperties(null);
            }}
            className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            Clear Filters
          </button>
        </div>
      )}
    </div>
  );
}

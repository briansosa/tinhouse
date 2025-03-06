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
  const initialLoadDone = useRef(false);
  const isMounted = useRef(true);

  useEffect(() => {
    const initialize = () => {
      // Cargar los filtros globales del localStorage
      const savedFilters = localStorage.getItem('globalFilters');
      
      if (savedFilters) {
        try {
          const parsedFilters = JSON.parse(savedFilters);
          setGlobalFilters(parsedFilters);
          
          if (isMounted.current) {
            loadProperties(parsedFilters);
          }
    
          return () => clearTimeout(timeoutId);
        } catch (error) {
          console.error("Error al parsear filtros:", error);
          
          if (isMounted.current) {
            loadProperties(null);
          }
        
          return () => clearTimeout(timeoutId);
        }
      } else {
        if (isMounted.current) {
          console.log("Cargando propiedades sin filtros después de delay");
          loadProperties(null);
        }
        
        return () => clearTimeout(timeoutId);
      }
    };
    
    // Marcar componente como montado
    isMounted.current = true;
    
    // Inicializar
    initialize();

    // Limpiar al desmontar
    return () => {
      isMounted.current = false;
      
      if (activeRequest.current) {
        activeRequest.current.cancel("Componente desmontado");
        activeRequest.current = null;
      }
    };
  }, []);

  // Efecto para recargar propiedades cuando cambien los filtros globales
  useEffect(() => {
    // Evitar la carga en el montaje inicial
    if (isInitialMount.current) {
      isInitialMount.current = false;
      return;
    }
    
    // Solo recargar si ya se hizo la carga inicial y hay cambios en los filtros
    if (initialLoadDone.current && globalFilters) {
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
    if (!isMounted.current) {
      return;
    }
    
    setLoading(true);
    
    // Cancelar cualquier solicitud activa
    if (activeRequest.current) {
      activeRequest.current.cancel("Cancelada por nueva petición");
      activeRequest.current = null;
    }
    
    try {
      // Crear un nuevo token de cancelación
      const source = axios.CancelToken.source();
      activeRequest.current = source;
      
      const response = await getUnratedProperties(filters, source.token);
      
      // Si el componente ya no está montado, no actualizar el estado
      if (!isMounted.current) {
        return;
      }
      
      // Verificar si esta petición sigue siendo la activa
      if (activeRequest.current === source) {
        if (response.data && Array.isArray(response.data.properties)) {
          setProperties(response.data.properties);
        } else {
          console.warn("Formato de respuesta inesperado:", response.data);
          setProperties([]);
        }
        // Limpiar la petición activa
        activeRequest.current = null;
      } 
    } catch (error) {
      // Si el componente ya no está montado, no actualizar el estado
      if (!isMounted.current) {
        return;
      }
      
      if (!axios.isCancel(error)) {
        console.error("Error al cargar propiedades:", error);
        // Solo actualizar el estado si no hay otra petición en curso
        if (!activeRequest.current || activeRequest.current.token === error.config?.cancelToken) {
          setProperties([]);
        }
      }
    } finally {
      // Solo actualizar el estado de carga si el componente sigue montado
      // y si esta era la última petición
      if (isMounted.current && (!activeRequest.current || activeRequest.current === null)) {
        setLoading(false);
      }
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

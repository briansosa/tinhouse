import { useEffect, useState } from 'react';
import { getUnratedProperties } from '../services/api';
import PropertyCard from '../components/PropertyCard/PropertyCard';

export default function Home({ setShowNavBar }) {
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);
  const [history, setHistory] = useState([]); // Historial de propiedades descartadas
  const [globalFilters, setGlobalFilters] = useState(null);

  useEffect(() => {
    // Cargar los filtros globales del localStorage
    const savedFilters = localStorage.getItem('globalFilters');
    if (savedFilters) {
      setGlobalFilters(JSON.parse(savedFilters));
    }
    
    loadProperties();
  }, []);

  // Efecto para recargar propiedades cuando cambien los filtros globales
  useEffect(() => {
    if (globalFilters) {
      loadProperties();
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

  const loadProperties = async () => {
    setLoading(true);
    try {
      const response = await getUnratedProperties(globalFilters);
      setProperties(response.data.properties || []);
    } catch (error) {
      console.error('Error fetching properties:', error);
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

  if (loading) {
    return <div className="flex justify-center items-center h-screen">Loading...</div>;
  }

  return (
    <div className="h-full">
      {properties.length > 0 ? (
        <PropertyCard 
          property={properties[0]}
          onRate={handleRate}
          onUndo={handleUndo}
          canUndo={history.length > 0}
          setShowNavBar={setShowNavBar}
        />
      ) : (
        <div className="text-center text-gray-600">No more properties to show</div>
      )}
    </div>
  );
}

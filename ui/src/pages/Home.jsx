import { useEffect, useState } from 'react';
import { getUnratedProperties } from '../services/api';
import PropertyCard from '../components/PropertyCard/PropertyCard';

export default function Home() {
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadProperties();
  }, []);

  const loadProperties = async () => {
    setLoading(true);
    try {
      const response = await getUnratedProperties();
      setProperties(response.data.properties || []);
    } catch (error) {
      console.error('Error fetching properties:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRate = () => {
    // Remover la propiedad actual y mostrar la siguiente
    setProperties(prev => prev.slice(1));
  };

  if (loading) {
    return <div className="flex justify-center items-center h-screen">Loading...</div>;
  }

  return (
    <div className="container mx-auto p-4">
      {properties.length > 0 ? (
        <PropertyCard 
          property={properties[0]}
          onRate={handleRate}
        />
      ) : (
        <div className="text-center text-gray-600">No more properties to show</div>
      )}
    </div>
  );
}

import { useState, useEffect } from 'react';
import { getAvailableFeatures } from '../../services/api';

const FeaturesFilter = ({ onChange, initialFeatures = [] }) => {
    // Asegurar que initialFeatures sea un array de strings
    const normalizedInitialFeatures = initialFeatures.map(id => String(id));
    
    const [selectedFeatures, setSelectedFeatures] = useState(normalizedInitialFeatures);
    const [featureCategories, setFeatureCategories] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    
    // Cargar características desde el backend
    useEffect(() => {
        const fetchFeatures = async () => {
            try {
                setIsLoading(true);
                const response = await getAvailableFeatures();
                
                // Normalizar los datos para asegurar que los IDs sean strings
                const normalizedCategories = response.data.categories.map(category => ({
                    ...category,
                    features: category.features.map(feature => ({
                        ...feature,
                        id: String(feature.id) // Asegurar que el ID sea string
                    }))
                }));
                
                setFeatureCategories(normalizedCategories || []);
                setError(null);
            } catch (err) {
                console.error('Error fetching features:', err);
                setError('No se pudieron cargar las características');
                // Usar categorías por defecto en caso de error
                setFeatureCategories(defaultFeatureCategories);
            } finally {
                setIsLoading(false);
            }
        };
        
        fetchFeatures();
    }, []);
    
    // Categorías por defecto en caso de error
    const defaultFeatureCategories = [
        {
            id: 'general',
            name: 'GENERAL',
            features: [
                { id: 'furnished', name: 'Amoblado' },
                { id: 'pets', name: 'Acepta mascotas' },
                { id: 'storage', name: 'Baulera' },
                { id: 'security', name: 'Seguridad 24hs' }
            ]
        },
        {
            id: 'amenities',
            name: 'AMENITIES',
            features: [
                { id: 'pool', name: 'Piscina' },
                { id: 'gym', name: 'Gimnasio' },
                { id: 'bbq', name: 'Parrilla' },
                { id: 'garden', name: 'Jardín' },
                { id: 'sum', name: 'SUM' }
            ]
        },
        {
            id: 'parking',
            name: 'ESTACIONAMIENTO',
            features: [
                { id: 'garage', name: 'Cochera fija' },
                { id: 'visitor_parking', name: 'Cochera visitantes' }
            ]
        },
        {
            id: 'services',
            name: 'SERVICIOS',
            features: [
                { id: 'gas', name: 'Gas natural' },
                { id: 'electricity', name: 'Electricidad' },
                { id: 'water', name: 'Agua corriente' },
                { id: 'internet', name: 'Internet' }
            ]
        }
    ];

    const handleFeatureToggle = (featureId) => {
        // featureId ya debe ser string
        const newFeatures = selectedFeatures.includes(featureId)
            ? selectedFeatures.filter(id => id !== featureId)
            : [...selectedFeatures, featureId];
        
        setSelectedFeatures(newFeatures);
        onChange(newFeatures);
    };

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4 text-center text-red-500">
                {error}
            </div>
        );
    }

    return (
        <div className="space-y-12 px-4">
            {featureCategories.map(category => (
                <div key={category.id} className="space-y-3">
                    <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase">
                        {category.name}
                    </h3>
                    <div className="grid grid-cols-2 gap-2">
                        {category.features.map(feature => (
                            <button
                                key={feature.id}
                                onClick={() => handleFeatureToggle(feature.id)}
                                className={`p-3 rounded-xl text-sm font-medium transition-all flex items-center gap-2 ${
                                    selectedFeatures.includes(feature.id)
                                        ? 'bg-blue-500 text-white'
                                        : 'bg-gray-800 text-gray-300'
                                }`}
                            >
                                {selectedFeatures.includes(feature.id) && (
                                    <svg className="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                    </svg>
                                )}
                                <span className="truncate">{feature.name}</span>
                            </button>
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );
};

export default FeaturesFilter; 
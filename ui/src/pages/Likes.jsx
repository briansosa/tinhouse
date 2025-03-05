import { useState, useEffect, useRef } from 'react';
import { getLikedProperties, togglePropertyFavorite, dislikeProperty } from '../services/api';
import LikedPropertyCard from '../components/LikedPropertyCard/LikedPropertyCard';
import PropertyNotes from '../components/PropertyNotes/PropertyNotes';
import PropertyDetails from '../components/PropertyDetails/PropertyDetails';
import Filters from '../components/Filters/Filters';
import FilterChips from '../components/Filters/FilterChips';
import axios from 'axios';

export default function Likes({ setShowNavBar }) {
    const [likedProperties, setLikedProperties] = useState([]);
    const [selectedProperty, setSelectedProperty] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showDetails, setShowDetails] = useState(false);
    const [showFilters, setShowFilters] = useState(false);
    const [showSortMenu, setShowSortMenu] = useState(false);
    const [sortBy, setSortBy] = useState('recent');
    const activeRequest = useRef(null);
    const [activeFilters, setActiveFilters] = useState({
        propertyType: null,
        showOnlyWithNotes: false,
        showOnlyFavorites: false,
        priceRange: {
            min: null,
            max: null,
            currency: null
        },
        locations: [],
        features: [],
        sizeRange: {
            min: null,
            max: null
        },
        rooms: null,
        bathrooms: null,
        antiquity: null
    });

    useEffect(() => {
        fetchLikedProperties();
        
        return () => {
            // Cancelar cualquier solicitud pendiente al desmontar
            if (activeRequest.current) {
                activeRequest.current.cancel();
            }
        };
    }, []);

    const fetchLikedProperties = async (filters = null) => {
        try {
            setIsLoading(true);
            
            // Cancelar cualquier solicitud activa
            if (activeRequest.current) {
                activeRequest.current.cancel();
            }
            
            // Crear un nuevo token de cancelación
            const source = axios.CancelToken.source();
            activeRequest.current = source;
            
            const response = await getLikedProperties(filters, source.token);
            
            // Solo actualizar si esta es la solicitud más reciente
            if (activeRequest.current === source) {
                setLikedProperties(response.data.properties || []);
                activeRequest.current = null;
            }
        } catch (err) {
            if (!axios.isCancel(err)) {
                setError('Error al cargar las propiedades');
                console.error('Error fetching liked properties:', err);
            }
        } finally {
            setIsLoading(false);
        }
    };

    const handleApplyFilters = (filters) => {
        setActiveFilters(filters);
        fetchLikedProperties(filters);
        setShowFilters(false);
        setShowNavBar(true);
    };

    const handleRemoveFilter = (filterName, defaultValue) => {
        const updatedFilters = {
            ...activeFilters,
            [filterName]: defaultValue
        };
        setActiveFilters(updatedFilters);
        fetchLikedProperties(updatedFilters);
    };

    const sortProperties = (properties) => {
        if (!properties || properties.length === 0) return [];
        
        const sortedProperties = [...properties];
        
        switch (sortBy) {
            case 'price_asc':
                return sortedProperties.sort((a, b) => {
                    const priceA = parseFloat(a.price.replace(/[^\d.-]/g, '')) || 0;
                    const priceB = parseFloat(b.price.replace(/[^\d.-]/g, '')) || 0;
                    return priceA - priceB;
                });
            case 'price_desc':
                return sortedProperties.sort((a, b) => {
                    const priceA = parseFloat(a.price.replace(/[^\d.-]/g, '')) || 0;
                    const priceB = parseFloat(b.price.replace(/[^\d.-]/g, '')) || 0;
                    return priceB - priceA;
                });
            case 'recent':
            default:
                return sortedProperties; // Ya vienen ordenados por fecha
        }
    };

    // Filtrar propiedades favoritas si está activado el filtro
    const filteredProperties = activeFilters.showOnlyFavorites 
        ? sortProperties(likedProperties).filter(property => property.is_favorite)
        : sortProperties(likedProperties);

    const hasActiveFilters = () => {
        return (
            activeFilters.propertyType !== 'all' ||
            activeFilters.showOnlyWithNotes ||
            activeFilters.showOnlyFavorites ||
            activeFilters.priceRange.min !== null ||
            activeFilters.priceRange.max !== null ||
            activeFilters.locations.length > 0 ||
            activeFilters.features.length > 0 ||
            activeFilters.sizeRange.min !== null ||
            activeFilters.sizeRange.max !== null ||
            activeFilters.rooms !== null ||
            activeFilters.bathrooms !== null ||
            activeFilters.antiquity !== null
        );
    };

    // Manejador para marcar/desmarcar favoritos
    const handleToggleFavorite = async (propertyId, isFavorite) => {
        try {
            await togglePropertyFavorite(propertyId, isFavorite);
            // Actualizar el estado local para reflejar el cambio
            setLikedProperties(prevProperties => 
                prevProperties.map(property => 
                    property.id === propertyId 
                        ? { ...property, is_favorite: isFavorite } 
                        : property
                )
            );
        } catch (error) {
            console.error('Error toggling favorite:', error);
        }
    };

    // Manejador para dar dislike a una propiedad
    const handleDislike = async (propertyId) => {
        try {
            await dislikeProperty(propertyId);
            // Actualizar la lista de propiedades eliminando la que recibió dislike
            setLikedProperties(prevProperties => 
                prevProperties.filter(property => property.id !== propertyId)
            );
            // Cerrar la vista de detalles o notas
            setSelectedProperty(null);
            setShowDetails(false);
            setShowNavBar(true);
        } catch (error) {
            console.error('Error marking property as dislike:', error);
        }
    };

    return (
        <div className="h-full bg-gray-950 dark:bg-gray-950">
            {!selectedProperty ? (
                <div className="h-full relative">
                    <div className="p-4">
                        <div className="flex items-center justify-between mb-4">
                            <h1 className="text-xl font-semibold dark:text-white">
                                Tus Favoritos
                            </h1>
                            <div className="flex items-center gap-2">
                                <button 
                                    onClick={() => {
                                        setShowFilters(true);
                                        setShowNavBar(false);
                                    }}
                                    className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
                                >
                                    <svg className="w-6 h-6 text-gray-600 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
                                    </svg>
                                </button>
                                <button 
                                    onClick={() => setShowSortMenu(prev => !prev)}
                                    className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 relative"
                                >
                                    <svg className="w-6 h-6 text-gray-600 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
                                    </svg>
                                </button>
                            </div>
                        </div>
                        
                        {showSortMenu && (
                            <div className="absolute right-16 top-16 bg-white dark:bg-gray-900 rounded-lg shadow-lg py-2 z-10">
                                {[
                                    { id: 'recent', label: 'Más reciente' },
                                    { id: 'price_asc', label: 'Menor precio' },
                                    { id: 'price_desc', label: 'Mayor precio' }
                                ].map(option => (
                                    <button
                                        key={option.id}
                                        className={`w-full px-4 py-2 text-left text-sm ${
                                            sortBy === option.id
                                                ? 'bg-blue-50 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400'
                                                : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                                        }`}
                                        onClick={() => {
                                            setSortBy(option.id);
                                            setShowSortMenu(false);
                                        }}
                                    >
                                        {option.label}
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                    
                    {/* Filtros rápidos de botones */}
                    <div className="px-4 mb-2 flex gap-2">
                        <button
                            className={`px-3 py-1.5 rounded-full text-sm ${
                                activeFilters.showOnlyWithNotes
                                    ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400'
                                    : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300'
                            }`}
                            onClick={() => {
                                const updatedFilters = {
                                    ...activeFilters,
                                    showOnlyWithNotes: !activeFilters.showOnlyWithNotes
                                };
                                setActiveFilters(updatedFilters);
                                fetchLikedProperties(updatedFilters);
                            }}
                        >
                            <span className="flex items-center gap-1">
                                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                                </svg>
                                Con notas
                            </span>
                        </button>
                        
                        <button
                            className={`px-3 py-1.5 rounded-full text-sm ${
                                activeFilters.showOnlyFavorites
                                    ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400'
                                    : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300'
                            }`}
                            onClick={() => {
                                const updatedFilters = {
                                    ...activeFilters,
                                    showOnlyFavorites: !activeFilters.showOnlyFavorites
                                };
                                setActiveFilters(updatedFilters);
                                fetchLikedProperties(updatedFilters);
                            }}
                        >
                            <span className="flex items-center gap-1">
                                <svg className="w-4 h-4" fill={activeFilters.showOnlyFavorites ? "currentColor" : "none"} viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                                </svg>
                                Favoritos
                            </span>
                        </button>
                    </div>
                    
                    {/* Filtros rápidos */}
                    <FilterChips 
                        filters={activeFilters} 
                        onRemove={(key, value) => {
                            const updatedFilters = { ...activeFilters, [key]: value };
                            setActiveFilters(updatedFilters);
                            fetchLikedProperties(updatedFilters);
                        }} 
                    />
                    
                    {/* Contenido principal - Reemplazamos el carousel por un div con scroll normal */}
                    <div 
                        className={`h-[calc(100%-9rem)] overflow-y-auto transition-all duration-300 ${
                            showFilters ? 'hidden' : ''
                        }`}
                    >
                        {isLoading ? (
                            <div className="flex justify-center items-center h-full">
                                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
                            </div>
                        ) : error ? (
                            <div className="text-center py-20 text-red-500">
                                {error}
                            </div>
                        ) : filteredProperties.length > 0 ? (
                            <div 
                                className="flex flex-col gap-4 px-4"
                            >
                                {filteredProperties.map(property => (
                                    <div key={property.id}>
                                        <LikedPropertyCard 
                                            property={property}
                                            onClick={() => {
                                                setSelectedProperty(property);
                                                setShowDetails(true);
                                                setShowNavBar(false);
                                            }}
                                            onToggleFavorite={handleToggleFavorite}
                                            onNotesClick={(property) => {
                                                setSelectedProperty(property);
                                                setShowDetails(false);
                                                setShowNavBar(false);
                                            }}
                                        />
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="flex flex-col items-center justify-center h-full text-gray-500 dark:text-gray-400 px-4 text-center">
                                <svg className="w-16 h-16 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                                </svg>
                                <p className="text-lg font-medium">No se encontraron propiedades</p>
                                <p className="mt-2">Prueba con otros filtros o agrega propiedades a tus favoritos</p>
                            </div>
                        )}
                    </div>

                    {showFilters && (
                        <div className="absolute inset-0 bg-gray-950">
                            {console.log("Montando componente Filters...")}
                            <Filters 
                                initialFilters={activeFilters}
                                onClose={() => {
                                    setShowFilters(false);
                                    setShowNavBar(true);
                                }}
                                onApplyFilters={handleApplyFilters}
                            />
                        </div>
                    )}
                </div>
            ) : showDetails ? (
                <PropertyDetails 
                    property={selectedProperty}
                    onClose={() => {
                        setShowDetails(false);
                        setSelectedProperty(null);
                        setShowNavBar(true);
                    }}
                    onDislike={() => handleDislike(selectedProperty.id)}
                    onNotesClick={() => {
                        setShowDetails(false);
                    }}
                />
            ) : (
                <PropertyNotes 
                    property={selectedProperty}
                    onClose={() => {
                        setSelectedProperty(null);
                        setShowNavBar(true);
                    }}
                    onImageClick={() => {
                        setShowDetails(true);
                    }}
                    onDislike={() => handleDislike(selectedProperty.id)}
                />
            )}
        </div>
    );
}

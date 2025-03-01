import ImageCarousel from '../ImageCarousel/ImageCarousel';
import { useState } from 'react';
import { dislikeProperty, togglePropertyFavorite } from '../../services/api';

export default function PropertyDetails({ property, onClose, onDislike, onNotesClick }) {
    const [showMenu, setShowMenu] = useState(false);
    const [isFavorite, setIsFavorite] = useState(property.is_favorite);
    const [isHeartHovered, setIsHeartHovered] = useState(false);

    const handleDislike = async () => {
        try {
            if (onDislike) {
                onDislike();
            } else {
                await dislikeProperty(property.id);
                setShowMenu(false);
                if (onClose) onClose();
            }
        } catch (err) {
            console.error('Error al marcar como dislike:', err);
        }
    };

    const handleToggleFavorite = async () => {
        try {
            await togglePropertyFavorite(property.id, !isFavorite);
            setIsFavorite(!isFavorite);
        } catch (err) {
            console.error('Error al cambiar estado de favorito:', err);
        }
    };

    return (
        <div className="h-full flex flex-col bg-white dark:bg-gray-950">
            {/* Contenido scrolleable */}
            <div className="flex-1 overflow-y-auto">
                {/* Carrusel de imágenes con botón de retroceso */}
                <div className="relative h-[40vh]">
                    <ImageCarousel property={property} />
                    
                    {onClose && (
                        <button 
                            onClick={onClose}
                            className={`absolute top-4 right-4 bg-white/30 hover:bg-white/50 dark:bg-black/30 dark:hover:bg-black/50 rounded-full p-3 transition-all transform z-20`}
                        >
                            <svg className="w-6 h-6" fill="none" stroke="white" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                            </svg>
                        </button>
                    )}
                </div>

                {/* Información detallada */}
                <div className="p-6 space-y-6">
                    {/* Información principal */}
                    <div>
                        <div className="flex justify-between items-start">
                            <h2 className="text-2xl font-bold dark:text-white">{property.title}</h2>
                            {onClose && (
                                <div className="flex items-center">
                                    {/* Botón de notas */}
                                    <div 
                                        className="p-2 rounded-full cursor-pointer transition-all duration-200"
                                        onClick={onNotesClick}
                                    >
                                        <svg 
                                            className={`w-6 h-6 transition-all duration-200 ${
                                                property.has_notes 
                                                    ? 'text-blue-500 fill-current' 
                                                    : 'fill-none stroke-gray-300 stroke-1'
                                            }`} 
                                            viewBox="0 0 24 24"
                                        >
                                            <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/>
                                        </svg>
                                    </div>
                                    
                                    {/* Botón de favorito */}
                                    <div 
                                        className="p-2 rounded-full cursor-pointer transition-all duration-200"
                                        onClick={handleToggleFavorite}
                                        onMouseEnter={() => setIsHeartHovered(true)}
                                        onMouseLeave={() => setIsHeartHovered(false)}
                                    >
                                        <svg 
                                            className={`w-6 h-6 transition-all duration-200 ${
                                                isFavorite 
                                                    ? 'text-red-500 fill-current' 
                                                : isHeartHovered 
                                                    ? 'text-red-400 fill-red-300' 
                                                    : 'fill-none stroke-gray-300 stroke-1'
                                            }`} 
                                            viewBox="0 0 24 24"
                                        >
                                            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
                                        </svg>
                                    </div>

                                    {/* Botón de menú (tres puntos) */}
                                    <div className="relative">
                                        <button 
                                            className="p-2 rounded-full cursor-pointer transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800 ml-1"
                                            onClick={() => setShowMenu(!showMenu)}
                                        >
                                            <svg className="w-6 h-6 dark:text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                                            </svg>
                                        </button>

                                        {/* Menú desplegable */}
                                        {showMenu && (
                                            <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg z-10">
                                                <div className="py-1">
                                                    <button 
                                                        className="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center"
                                                        onClick={handleDislike}
                                                    >
                                                        <svg className="w-5 h-5 mr-2 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.096c.5 0 .905-.405.905-.904 0-.715.211-1.413.608-2.008L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5" />
                                                        </svg>
                                                        Ya no me gusta
                                                    </button>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            )}
                        </div>
                        <span className="text-4xl font-bold dark:text-white block mt-2">{property.price}</span>
                        <div className="flex items-center gap-2 mt-2">
                            <p className="text-gray-600 dark:text-gray-300">{property.location}</p>
                            {property.property_type && (
                                <>
                                    <span className="text-gray-400 dark:text-gray-500">•</span>
                                    <p className="text-gray-600 dark:text-gray-300">{property.property_type}</p>
                                </>
                            )}
                        </div>
                        {property.agency?.name && (
                            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{property.agency.name}</p>
                        )}
                    </div>

                    {/* Características destacadas */}
                    <div className="grid grid-cols-3 gap-4 py-4 border-y border-gray-200 dark:border-gray-700">
                        <div className="text-center">
                            <p className="text-3xl font-semibold dark:text-white">{property.details.bedrooms}</p>
                            <p className="text-sm text-gray-600 dark:text-gray-400">Dormitorios</p>
                        </div>
                        <div className="text-center">
                            <p className="text-3xl font-semibold dark:text-white">{property.details.bathrooms}</p>
                            <p className="text-sm text-gray-600 dark:text-gray-400">Baños</p>
                        </div>
                        <div className="text-center">
                            <p className="text-3xl font-semibold dark:text-white">{property.details.total_area}</p>
                            <p className="text-sm text-gray-600 dark:text-gray-400">m² Totales</p>
                        </div>
                    </div>

                    {/* Detalles del inmueble */}
                    <div>
                        <h3 className="text-2xl font-bold mb-4 dark:text-white">Detalles del inmueble</h3>
                        
                        {/* A. Información General */}
                        <div className="mb-6">
                            <h4 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-3">Información General</h4>
                            <div className="space-y-3">
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Situación</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.status || 'Sin especificar'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Antigüedad</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.age || 'Sin especificar'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Plantas</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.floors || '-'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Ambientes</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.rooms || '-'}</span>
                                </div>
                            </div>
                        </div>

                        {/* B. Superficies */}
                        <div className="mb-6">
                            <h4 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-3">Superficies</h4>
                            <div className="space-y-3">
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Superficie Total</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.total_area ? `${property.details.total_area}m²` : '-'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Superficie Cubierta</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.area ? `${property.details.area}m²` : '-'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Superficie Terreno</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.land_area ? `${property.details.land_area}m²` : '-'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Frente x Fondo</span>
                                    <span className="font-medium dark:text-gray-200">
                                        {property.details.front_size && property.details.back_size 
                                            ? `${property.details.front_size}m x ${property.details.back_size}m`
                                            : '-'}
                                    </span>
                                </div>
                            </div>
                        </div>

                        {/* C. Características Adicionales */}
                        <div>
                            <h4 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-3">Características Adicionales</h4>
                            <div className="space-y-3">
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Cocheras</span>
                                    <span className="font-medium dark:text-gray-200">{property.details.garages || 'Sin cochera'}</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Expensas</span>
                                    <span className="font-medium dark:text-gray-200">
                                        {property.details.expenses ? `$${property.details.expenses}` : 'Sin expensas'}
                                    </span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-gray-600 dark:text-gray-400">Código</span>
                                    <span className="font-medium dark:text-gray-200">{property.code || '-'}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Descripción */}
                    {property.description && (
                        <div>
                            <h3 className="text-2xl font-bold mb-3 dark:text-white">Descripción</h3>
                            <p className="text-gray-700 dark:text-gray-300 leading-relaxed">{property.description}</p>
                        </div>
                    )}

                    {/* Botón de contacto */}
                    <a 
                        href={property.url} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        className="block w-full text-center bg-gradient-to-r from-blue-500 to-blue-600 text-white py-4 rounded-lg hover:from-blue-600 hover:to-blue-700 transition-all font-medium"
                    >
                        Ver en {property.agency?.name || 'sitio web'}
                    </a>
                </div>
            </div>
        </div>
    );
} 
import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';

export default function LikedPropertyCard({ property, onClick, onToggleFavorite }) {
    // Ya no necesitamos verificar el localStorage, la API nos indica si la propiedad tiene notas
    const hasNotes = property.has_notes;
    const isFavorite = property.is_favorite;
    const [isHovered, setIsHovered] = useState(false);

    // Función para manejar el clic en el botón de favorito
    const handleFavoriteClick = (e) => {
        e.stopPropagation(); // Evitar que el clic se propague al contenedor
        onToggleFavorite(property.id, !isFavorite);
    };

    return (
        <motion.div 
            className="flex items-start gap-3 bg-white dark:bg-gray-900 rounded-lg p-1.5 shadow-sm cursor-pointer select-none relative"
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            onClick={onClick}
        >
            {/* Contenedor de imagen */}
            <div className="relative w-16 h-16 flex-shrink-0 select-none">
                <div className="w-full h-full rounded-full overflow-hidden">
                    <img 
                        src={property.image_url} 
                        alt={property.title}
                        className="w-full h-full object-cover pointer-events-none"
                        draggable="false"
                    />
                </div>
            </div>

            {/* Información */}
            <div className="flex-1 min-w-0 select-none pointer-events-none py-0.5">
                <p className="text-base font-semibold dark:text-white truncate">{property.title}</p>
                <p className="text-sm text-gray-700 dark:text-gray-300 truncate">{property.location}</p>
                <p className="text-sm text-gray-700 dark:text-gray-300 truncate">{property.price}</p>
            </div>

            {/* Iconos de acciones (notas y favoritos) */}
            <div className="absolute bottom-1.5 right-1.5 flex items-center gap-2">
                {/* Indicador de notas */}
                {hasNotes && (
                    <div className="bg-blue-500 text-white p-1 rounded-full pointer-events-none">
                        <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/>
                        </svg>
                    </div>
                )}
                
                {/* Botón de favorito */}
                <div 
                    className="p-1 rounded-full cursor-pointer z-10 transition-all duration-200"
                    onClick={handleFavoriteClick}
                    onMouseEnter={() => setIsHovered(true)}
                    onMouseLeave={() => setIsHovered(false)}
                >
                    <svg 
                        className={`w-5 h-5 transition-all duration-200 ${
                            isFavorite 
                                ? 'text-red-500 fill-current' 
                                : isHovered 
                                    ? 'text-red-400 fill-red-300' 
                                    : 'fill-none stroke-red-300 stroke-1'
                        }`} 
                        viewBox="0 0 24 24"
                    >
                        <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
                    </svg>
                </div>
            </div>
        </motion.div>
    );
} 
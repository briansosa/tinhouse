import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';

export default function LikedPropertyCard({ property, onClick }) {
    // Ya no necesitamos verificar el localStorage, la API nos indica si la propiedad tiene notas
    const hasNotes = property.has_notes;

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
            </div>

            {/* Indicador de notas */}
            {hasNotes && (
                <div className="absolute bottom-1.5 right-1.5 bg-blue-500 text-white p-1 rounded-full pointer-events-none">
                    <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/>
                    </svg>
                </div>
            )}
        </motion.div>
    );
} 
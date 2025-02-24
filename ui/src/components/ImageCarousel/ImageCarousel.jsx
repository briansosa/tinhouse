import { useState, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

export default function ImageCarousel({ property }) {
    const [currentIndex, setCurrentIndex] = useState(0);

    // Mejorar la lógica para determinar qué imágenes mostrar
    const images = useMemo(() => {
        // Si hay imágenes en el array, usarlas
        if (property.images && Array.isArray(property.images) && property.images.length > 0) {
            return property.images;
        }
        
        // Si no hay imágenes pero hay image_url, usar esa
        if (property.image_url) {
            return [property.image_url];
        }
        
        // Si no hay ninguna imagen, usar placeholder
        return ['https://via.placeholder.com/400x300'];
    }, [property.images, property.image_url]);

    const nextImage = () => {
        if (currentIndex < images.length - 1) {
            setCurrentIndex(prev => prev + 1);
        }
    };

    const prevImage = () => {
        if (currentIndex > 0) {
            setCurrentIndex(prev => prev - 1);
        }
    };

    // No mostrar controles si solo hay una imagen
    const showControls = images.length > 1;

    return (
        <div className="relative w-full h-full overflow-hidden">
            <AnimatePresence mode="wait">
                <motion.img
                    key={currentIndex}
                    src={images[currentIndex]}
                    alt={property.title || 'Property image'}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.2 }}
                    className="absolute inset-0 w-full h-full object-cover"
                    draggable="false"
                    onError={(e) => {
                        e.target.src = 'https://via.placeholder.com/400x300';
                    }}
                />
            </AnimatePresence>
            
            {showControls && (
                <>
                    {/* Botones de navegación */}
                    <div className="absolute inset-0 flex items-center justify-between px-4">
                        <button
                            onClick={prevImage}
                            className="w-8 h-8 rounded-full bg-black/30 hover:bg-black/50 flex items-center justify-center transition-colors backdrop-blur-sm"
                            disabled={currentIndex === 0}
                        >
                            <svg className="w-5 h-5" fill="none" stroke="white" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M15 19l-7-7 7-7" />
                            </svg>
                        </button>

                        <button
                            onClick={nextImage}
                            className="w-8 h-8 rounded-full bg-black/30 hover:bg-black/50 flex items-center justify-center transition-colors backdrop-blur-sm"
                            disabled={currentIndex === images.length - 1}
                        >
                            <svg className="w-5 h-5" fill="none" stroke="white" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M9 5l7 7-7 7" />
                            </svg>
                        </button>
                    </div>

                    {/* Indicadores de posición */}
                    <div className="absolute top-4 left-0 right-0 flex justify-center gap-1">
                        {images.map((_, idx) => (
                            <div
                                key={idx}
                                className={`h-1 transition-all duration-300 rounded-full ${
                                    idx === currentIndex 
                                        ? 'w-6 bg-white' 
                                        : 'w-2 bg-white/50'
                                }`}
                            />
                        ))}
                    </div>
                </>
            )}
        </div>
    );
}

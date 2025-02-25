import { useState, useEffect } from 'react';
import { getLikedProperties } from '../services/api';
import LikedPropertyCard from '../components/LikedPropertyCard/LikedPropertyCard';
import PropertyNotes from '../components/PropertyNotes/PropertyNotes';
import { useDrag } from '@use-gesture/react';
import { motion, useMotionValue, animate } from 'framer-motion';
import PropertyDetails from '../components/PropertyDetails/PropertyDetails';

export default function Likes({ setShowNavBar }) {
    const [likedProperties, setLikedProperties] = useState([]);
    const [selectedProperty, setSelectedProperty] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [containerHeight, setContainerHeight] = useState(0);
    const [contentHeight, setContentHeight] = useState(0);
    const y = useMotionValue(0);
    const [showDetails, setShowDetails] = useState(false);

    useEffect(() => {
        const fetchLikedProperties = async () => {
            try {
                setIsLoading(true);
                const response = await getLikedProperties();
                setLikedProperties(response.data.properties || []);
            } catch (err) {
                setError('Error al cargar las propiedades');
                console.error('Error fetching liked properties:', err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchLikedProperties();
    }, []);

    useEffect(() => {
        const container = document.getElementById('carousel-container');
        const content = document.getElementById('carousel-content');
        if (container && content) {
            setContainerHeight(container.offsetHeight);
            setContentHeight(content.scrollHeight);
            content.style.touchAction = 'none';
        }
    }, [likedProperties]);

    const bind = useDrag(({ down, movement: [mx, my] }) => {
        const currentY = y.get();
        const targetY = currentY + my;
        
        // Límites del scroll
        const minY = -(contentHeight - containerHeight);
        const boundedY = Math.max(minY, Math.min(0, targetY));

        // Aplicar la posición
        y.set(boundedY);

        // Si soltamos, aplicar una animación simple para "asentar" el scroll
        if (!down) {
            animate(y, boundedY, {
                type: "spring",
                damping: 20,
                stiffness: 200,
                mass: 0.5
            });
        }
    }, {
        axis: 'y',
        filterTaps: true,
        bounds: {
            top: 0,
            bottom: -(contentHeight - containerHeight)
        },
        rubberband: 0.5
    });

    return (
        <div className="h-full bg-white dark:bg-gray-950">
            {!selectedProperty ? (
                <div className="h-full">
                    <div className="p-4">
                        <h1 className="text-xl font-semibold dark:text-white mb-4">
                            Tus Favoritos
                        </h1>
                    </div>
                    
                    <div id="carousel-container" className="h-[calc(100%-5rem)] overflow-hidden">
                        {isLoading ? (
                            <div>Cargando...</div>
                        ) : (
                            <motion.div 
                                id="carousel-content"
                                {...bind()}
                                style={{ y }}
                                className="flex flex-col gap-4 px-4 touch-none cursor-grab active:cursor-grabbing"
                            >
                                {likedProperties.map(property => (
                                    <div key={property.id}>
                                        <LikedPropertyCard
                                            property={property}
                                            onClick={() => {
                                                setSelectedProperty(property);
                                                setShowNavBar(false);
                                            }}
                                        />
                                    </div>
                                ))}
                            </motion.div>
                        )}
                    </div>
                </div>
            ) : showDetails ? (
                <PropertyDetails 
                    property={selectedProperty}
                    onClose={() => {
                        setShowDetails(false);
                        setShowNavBar(false);
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
                        setShowNavBar(false);
                    }}
                />
            )}
        </div>
    );
}

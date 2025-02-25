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
    const [containerWidth, setContainerWidth] = useState(0);
    const [contentWidth, setContentWidth] = useState(0);
    const x = useMotionValue(0);
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
            setContainerWidth(container.offsetWidth);
            setContentWidth(content.scrollWidth);
        }
    }, [likedProperties]);

    const bind = useDrag(({ down, movement: [mx], velocity, direction: [dx] }) => {
        const currentX = x.get();
        const targetX = currentX + mx;
        
        // Límites del scroll
        const minX = -contentWidth + containerWidth;
        const boundedX = Math.max(minX, Math.min(0, targetX));

        if (down) {
            // Durante el drag, movimiento más suave
            x.set(boundedX);
        } else {
            // Al soltar, calculamos la inercia basada en la velocidad
            const momentum = velocity * dx * 150; // Cambiamos direction por dx
            const finalX = Math.max(minX, Math.min(0, currentX + momentum));
            
            animate(x, finalX, {
                type: "spring",
                damping: 40,
                stiffness: 90,
                mass: 0.8,
                restDelta: 0.01,
                restSpeed: 0.01,
                velocity: velocity * dx * 2 // Cambiamos direction por dx
            });
        }
    }, {
        axis: 'x',
        filterTaps: true,
        from: () => [x.get(), 0],
        bounds: {
            left: -contentWidth + containerWidth,
            right: 0
        },
        rubberband: 0.5
    });

    return (
        <div className="h-full bg-white dark:bg-gray-950">
            {!selectedProperty ? (
                // Vista de cards
                <div className="h-full">
                    <div className="p-4">
                        <h1 className="text-xl font-semibold dark:text-white mb-4">
                            Tus Favoritos
                        </h1>
                    </div>
                    
                    <div id="carousel-container" className="px-4 overflow-hidden">
                        {isLoading ? (
                            <div>Cargando...</div>
                        ) : (
                            <motion.div 
                                id="carousel-content"
                                {...bind()}
                                style={{ x }}
                                className="flex gap-2 pb-4"
                            >
                                {likedProperties.map(property => (
                                    <LikedPropertyCard
                                        key={property.id}
                                        property={property}
                                        onClick={() => {
                                            setSelectedProperty(property)
                                            setShowNavBar(false)
                                        }}
                                    />
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

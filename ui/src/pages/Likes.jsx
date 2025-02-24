import { useState, useEffect } from 'react';
import { getLikedProperties } from '../services/api';
import LikedPropertyCard from '../components/LikedPropertyCard/LikedPropertyCard';
import { useDrag } from '@use-gesture/react';
import { motion, useMotionValue, animate } from 'framer-motion';

export default function Likes() {
    const [likedProperties, setLikedProperties] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [containerWidth, setContainerWidth] = useState(0);
    const [contentWidth, setContentWidth] = useState(0);
    const x = useMotionValue(0);

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
            const momentum = velocity * direction * 150; // Aumentamos el factor de inercia
            const finalX = Math.max(minX, Math.min(0, currentX + momentum));
            
            animate(x, finalX, {
                type: "spring",
                damping: 40,      // Más resistencia para control
                stiffness: 90,    // Menos rigidez para movimiento más suave
                mass: 0.8,        // Masa equilibrada
                restDelta: 0.01,  // Precisión del punto de reposo
                restSpeed: 0.01,  // Velocidad mínima para considerar reposo
                velocity: velocity * direction * 2 // Duplicamos la velocidad inicial
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
        rubberband: 0.5,  // Más elasticidad en los límites
        delay: 0,
    });

    return (
        <div className="h-full bg-white dark:bg-gray-950">
            <div className="p-4">
                <h1 className="text-xl font-semibold dark:text-white mb-4">
                    Tus Favoritos
                </h1>
            </div>
            
            <div id="carousel-container" className="px-4 overflow-hidden">
                {isLoading ? (
                    <div className="flex items-center justify-center w-full p-4">
                        <div className="text-gray-500 dark:text-gray-400">
                            Cargando propiedades...
                        </div>
                    </div>
                ) : error ? (
                    <div className="flex items-center justify-center w-full p-4">
                        <div className="text-red-500">{error}</div>
                    </div>
                ) : likedProperties.length === 0 ? (
                    <div className="flex items-center justify-center w-full p-4">
                        <div className="text-gray-500 dark:text-gray-400">
                            No tienes propiedades favoritas aún
                        </div>
                    </div>
                ) : (
                    <motion.div 
                        id="carousel-content"
                        {...bind()}
                        style={{ x }}
                        className="flex gap-2 pb-4 touch-none cursor-grab active:cursor-grabbing"
                    >
                        {likedProperties.map(property => (
                            <LikedPropertyCard
                                key={property.id}
                                property={property}
                            />
                        ))}
                    </motion.div>
                )}
            </div>
        </div>
    );
}

import { useState } from 'react';
import { motion, useMotionValue, useTransform, animate } from 'framer-motion';
import { useDrag } from '@use-gesture/react';
import { rateProperty } from '../../services/api';
import PropertyDetails from '../PropertyDetails/PropertyDetails';

export default function PropertyCard({ property, onRate, onUndo, canUndo }) {
  const [showDetails, setShowDetails] = useState(false);
  const [isRating, setIsRating] = useState(false);

  // Valores para la animación
  const x = useMotionValue(0);
  const rotate = useTransform(x, [200, -200], [10, -10]);
  
  // Indicadores con su propia opacidad
  const nopeOpacity = useTransform(x, [0, 100], [0, 1]); 
  const likeOpacity = useTransform(x, [-100, 0], [1, 0]);
  const dislikeScale = useTransform(x, [0, 125], [1, 1.5]);
  const likeScale = useTransform(x, [-125, 0], [1.5, 1]);

  const handleRate = async (rating) => {
    if (isRating) return;
    
    setIsRating(true);
    try {
      const direction = rating === 'like' ? -200 : 200;
      
      // Primero hacemos el PUT a la API
      await rateProperty(property.id, rating);
      
      // Luego la animación
      await animate(x, direction, { 
        type: "spring",
        stiffness: 1900,
        damping: 48,
        duration: 0.5,
        velocity: 5,
        onComplete: () => {
          onRate(rating);
          x.set(0);
        }
      });

    } catch (error) {
      console.error('Error rating property:', error);
      animate(x, 0, { duration: 0.1 });
    } finally {
      setIsRating(false);
    }
  };

  // Configuración del gesto de arrastre
  const bind = useDrag(({ down, movement: [mx], direction: [xDir], velocity, active }) => {
    if (isRating) return;

    if (down) {
      // Mientras arrastramos
      x.set(mx);
    } else {
      // Cuando soltamos
      const swipe = Math.abs(mx) > 60 || velocity > 0.4;
      if (swipe) {
        const rating = mx > 0 ? 'dislike' : 'like';
        handleRate(rating);
      } else {
        // Si no fue suficiente el swipe, volvemos al centro
        animate(x, 0, { 
          type: "spring",
          duration: 0.15,
          bounce: 0
        });
      }
    }
  }, {
    axis: 'x',
    filterTaps: true,
    bounds: { left: -300, right: 300 },
    rubberband: true,
    from: () => [x.get(), 0]
  });

  return (
    <div className="relative max-w-md mx-auto">
      <motion.div
        {...bind()}
        style={{ x, rotate }}
        className="relative bg-white dark:bg-gray-950 rounded-xl shadow-lg overflow-hidden h-[85vh] touch-none cursor-grab active:cursor-grabbing"
      >
        {/* Indicadores de Like/Dislike con su propia opacidad */}
        <motion.div 
          style={{ scale: likeScale, opacity: likeOpacity }}
          className="absolute top-8 right-8 z-10 border-4 border-green-500 text-green-500 px-8 py-2 rounded-full font-bold transform rotate-12"
        >
          LIKE
        </motion.div>
        <motion.div 
          style={{ scale: dislikeScale, opacity: nopeOpacity }}
          className="absolute top-8 left-8 z-10 border-4 border-red-500 text-red-500 px-8 py-2 rounded-full font-bold transform -rotate-12"
        >
          NOPE
        </motion.div>

        {/* Contenido de la card */}
        <div className="h-full overflow-y-auto select-none">
          <div className={`relative ${showDetails ? 'h-[40vh]' : 'h-[73vh]'} transition-all duration-300`}>
            <img 
              src={property.image_url || 'https://via.placeholder.com/400x300'} 
              alt={property.title}
              className="w-full h-full object-cover pointer-events-none"
              draggable="false"
            />
            
            {!showDetails && (
              <>
                <div className="absolute bottom-0 left-0 right-0 h-48 bg-gradient-to-t from-black/90 to-transparent pointer-events-none" />
                <div className="absolute bottom-2 left-0 right-0 px-4 text-white pointer-events-none">
                  <h2 className="text-2xl font-bold leading-tight">{property.title}</h2>
                  <p className="text-sm mt-1 opacity-90">{property.location}</p>
                  <p className="text-xl font-semibold mt-1">{property.price}</p>
                  <div className="flex items-center gap-2 mt-1 text-sm opacity-80">
                    <span>{property.details.bedrooms} dorm.</span>
                    <span>•</span>
                    <span>{property.details.bathrooms} baños</span>
                    <span>•</span>
                    <span>{property.details.area}m²</span>
                  </div>
                </div>
              </>
            )}

            <button 
              onClick={() => setShowDetails(!showDetails)}
              className={`absolute top-4 right-4 bg-white/30 hover:bg-white/50 dark:bg-black/30 dark:hover:bg-black/50 rounded-full p-3 transition-all transform ${showDetails ? 'rotate-180' : ''} z-20`}
            >
              <svg className="w-6 h-6" fill="none" stroke="white" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
              </svg>
            </button>
          </div>

          {showDetails && (
            <div className="pb-16">
              <PropertyDetails property={property} />
            </div>
          )}
        </div>
      </motion.div>

      {/* Botones de acción */}
      <div className={`
        absolute bottom-0 left-0 right-0 
        ${showDetails ? 'h-16 bg-gradient-to-t from-white dark:from-black via-white dark:via-black to-transparent' : ''}
      `}>
        <div className={`
          flex justify-around items-center px-4
          ${showDetails 
            ? 'h-full pb-6'
            : 'h-24 py-6'
          }
        `}>
          <button 
            className="rounded-full p-3 bg-gray-800/50 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none"
            disabled={!canUndo || isRating}
            onClick={onUndo}
          >
            <svg className="w-7 h-7" fill="none" stroke="white" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
          </button>
          <button 
            className="rounded-full p-4 bg-gray-800/50 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none"
            disabled={isRating}
            onClick={() => handleRate('dislike')}
          >
            <svg className="w-8 h-8" fill="#FF4458" viewBox="0 0 24 24">
              <path d="M14.8 12l3.6-3.6c.8-.8.8-2 0-2.8-.8-.8-2-.8-2.8 0L12 9.2 8.4 5.6c-.8-.8-2-.8-2.8 0-.8.8-.8 2 0 2.8L9.2 12l-3.6 3.6c-.8.8-.8 2 0 2.8.4.4.9.6 1.4.6.5 0 1-.2 1.4-.6l3.6-3.6 3.6 3.6c.4.4.9.6 1.4.6.5 0 1-.2 1.4-.6.8-.8.8-2 0-2.8L14.8 12z"/>
            </svg>
          </button>
          <button 
            className="rounded-full p-4 bg-gray-800/50 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none"
            disabled={isRating}
            onClick={() => handleRate('like')}
          >
            <svg className="w-8 h-8" fill="#00DC7D" viewBox="0 0 24 24">
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
} 
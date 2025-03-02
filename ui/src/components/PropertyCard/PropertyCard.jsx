import { useState, useEffect, useRef } from 'react';
import { motion, useMotionValue, useTransform, animate, AnimatePresence } from 'framer-motion';
import { useDrag } from '@use-gesture/react';
import { rateProperty } from '../../services/api';
import PropertyDetails from '../PropertyDetails/PropertyDetails';
import ImageCarousel from '../ImageCarousel/ImageCarousel';
import { Link, useLocation } from 'react-router-dom';

export default function PropertyCard({ property, onRate, onUndo, canUndo, setShowNavBar }) {
  const [showDetails, setShowDetails] = useState(false);
  const [isRating, setIsRating] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const isTogglingDetails = useRef(false);
  const dragStartTime = useRef(null);
  const lastDragPosition = useRef(0);
  const dragStartPosition = useRef(null);
  const hasMoved = useRef(false);

  // Constantes para la configuración del arrastre
  const CARD_WIDTH = 300; // Ancho aproximado de la tarjeta
  const ACTIVATION_THRESHOLD = CARD_WIDTH * 0.4; // 40% del ancho para activar like/dislike
  const MOVEMENT_THRESHOLD = 15; // Umbral para considerar que ha habido movimiento

  // Valores para la animación
  const x = useMotionValue(0);
  const rotate = useTransform(x, [CARD_WIDTH, -CARD_WIDTH], [15, -15]);
  
  // Indicadores con su propia opacidad y escala
  const nopeOpacity = useTransform(x, [0, ACTIVATION_THRESHOLD * 0.5, ACTIVATION_THRESHOLD], [0, 0.5, 1]); 
  const likeOpacity = useTransform(x, [-ACTIVATION_THRESHOLD, -ACTIVATION_THRESHOLD * 0.5, 0], [1, 0.5, 0]);
  
  // Escala para los indicadores - aumenta gradualmente hasta el umbral de activación
  const dislikeScale = useTransform(x, [0, ACTIVATION_THRESHOLD * 0.5, ACTIVATION_THRESHOLD], [1, 1.2, 1.5]);
  const likeScale = useTransform(x, [-ACTIVATION_THRESHOLD, -ACTIVATION_THRESHOLD * 0.5, 0], [1.5, 1.2, 1]);
  
  // Color de borde para indicar proximidad a la activación
  const likeBorderColor = useTransform(
    x, 
    [-ACTIVATION_THRESHOLD, -ACTIVATION_THRESHOLD * 0.7, -ACTIVATION_THRESHOLD * 0.3, 0], 
    ['#00DC7D', '#00DC7D', 'rgba(0, 220, 125, 0.7)', 'rgba(0, 220, 125, 0.5)']
  );
  
  const dislikeBorderColor = useTransform(
    x, 
    [0, ACTIVATION_THRESHOLD * 0.3, ACTIVATION_THRESHOLD * 0.7, ACTIVATION_THRESHOLD], 
    ['rgba(255, 68, 88, 0.5)', 'rgba(255, 68, 88, 0.7)', '#FF4458', '#FF4458']
  );

  // Efecto para manejar el estado de la barra de navegación cuando cambia showDetails
  useEffect(() => {
    setShowNavBar(!showDetails);
  }, [showDetails, setShowNavBar]);

  const handleRate = async (rating) => {
    if (isRating) {
      console.log('Rating cancelado: ya está en proceso');
      return;
    }
    
    console.log('Iniciando rating:', rating);
    setIsRating(true);
    try {
      const direction = rating === 'like' ? -CARD_WIDTH : CARD_WIDTH;
      
      // Primero hacemos el PUT a la API
      console.log('Llamando a API para rating:', rating);
      await rateProperty(property.id, rating);
      
      // Luego la animación
      console.log('Iniciando animación de rating');
      await animate(x, direction, { 
        type: "spring",
        stiffness: 1900,
        damping: 48,
        duration: 0.5,
        velocity: 5,
        onComplete: () => {
          console.log('Animación completada, llamando a onRate');
          onRate(rating);
          x.set(0);
        }
      });

    } catch (error) {
      console.error('Error rating property:', error);
      animate(x, 0, { duration: 0.1 });
    } finally {
      console.log('Finalizando proceso de rating');
      setIsRating(false);
      setShowNavBar(true);
      setShowDetails(false);
      setIsDragging(false);
      dragStartTime.current = null;
      lastDragPosition.current = 0;
      dragStartPosition.current = null;
      hasMoved.current = false;
    }
  };

  // Configuración del gesto de arrastre
  const bind = useDrag(({ down, movement: [mx], direction: [xDir], velocity, active, event, first, last, tap, memo }) => {
    // Si es un tap (clic sin arrastre), ignoramos el evento para el arrastre
    if (tap) {
      console.log('Evento de tap detectado, ignorando para arrastre');
      return;
    }
    
    if (isRating || isTogglingDetails.current) {
      console.log('Drag cancelado: isRating o isTogglingDetails activo');
      return;
    }
    
    // Registrar el tiempo y posición de inicio del arrastre
    if (first) {
      dragStartTime.current = Date.now();
      dragStartPosition.current = mx;
      hasMoved.current = false;
      console.log('Inicio de arrastre registrado:', { time: dragStartTime.current, position: dragStartPosition.current });
    }
    
    // Determinar si ha habido movimiento significativo
    if (Math.abs(mx) > MOVEMENT_THRESHOLD && !hasMoved.current) {
      hasMoved.current = true;
      console.log('Movimiento significativo detectado');
    }
    
    // Guardar la última posición de arrastre
    lastDragPosition.current = mx;
    
    console.log('Estado del drag:', { 
      down, 
      mx, 
      velocity, 
      active, 
      first,
      last,
      tap,
      hasMoved: hasMoved.current,
      isDragging: isDragging,
      isRating: isRating,
      isTogglingDetails: isTogglingDetails.current,
      dragTime: dragStartTime.current ? Date.now() - dragStartTime.current : null,
      activationThreshold: ACTIVATION_THRESHOLD,
      percentComplete: Math.abs(mx) / ACTIVATION_THRESHOLD * 100
    });
    
    // Actualizar el estado de arrastre solo si hay movimiento real
    if (active && hasMoved.current) {
      setIsDragging(true);
    } else if (!active) {
      // Cuando soltamos, esperamos un poco antes de resetear isDragging
      setTimeout(() => {
        setIsDragging(false);
      }, 100);
    }
    
    // Ocultar la barra de navegación solo cuando estamos arrastrando con movimiento real
    if (active && hasMoved.current && !showDetails) {
      setShowNavBar(false);
    } else if (!active && !showDetails) {
      // Restaurar la barra de navegación cuando dejamos de arrastrar
      setTimeout(() => {
        if (!showDetails && !isRating && !isTogglingDetails.current) {
          setShowNavBar(true);
        }
      }, 200);
    }

    if (down) {
      // Mientras arrastramos
      console.log('Arrastrando - mx:', mx, 'x actual:', x.get(), 'porcentaje:', Math.abs(mx) / ACTIVATION_THRESHOLD * 100 + '%');
      x.set(mx);
      
      // Si alcanzamos el umbral de activación mientras arrastramos, podemos activar el like/dislike
      if (Math.abs(mx) > ACTIVATION_THRESHOLD && hasMoved.current) {
        console.log('Umbral de activación alcanzado mientras arrastramos');
        // Opcional: podríamos activar aquí, pero mejor esperar a que suelte
      }
    } else if (last && dragStartTime.current) {
      // Solo procesamos el final del arrastre cuando realmente es el último evento y tenemos un tiempo de inicio válido
      const dragDuration = Date.now() - dragStartTime.current;
      console.log('Soltado - mx:', mx, 'velocity:', velocity[0], 'tiempo total:', dragDuration);
      
      // Verificar si fue un arrastre significativo
      const isQuickSwipe = Math.abs(velocity[0]) > 0.8;
      const isLongDrag = Math.abs(mx) > ACTIVATION_THRESHOLD;
      // Solo consideramos significativo si hay movimiento real y no ha pasado demasiado tiempo
      const isSignificantDrag = (isQuickSwipe || isLongDrag) && Math.abs(mx) > MOVEMENT_THRESHOLD * 2 && dragDuration < 1000 && hasMoved.current;
      
      console.log('Análisis de swipe:', { 
        isQuickSwipe, 
        isLongDrag, 
        dragDuration, 
        isSignificantDrag, 
        absMovement: Math.abs(mx),
        hasMoved: hasMoved.current,
        percentComplete: Math.abs(mx) / ACTIVATION_THRESHOLD * 100 + '%'
      });
      
      if (isSignificantDrag) {
        const rating = mx > 0 ? 'dislike' : 'like';
        console.log('Swipe detectado:', rating);
        handleRate(rating);
      } else {
        // Si no fue suficiente el swipe, volvemos al centro
        console.log('Swipe insuficiente, volviendo al centro');
        animate(x, 0, { 
          type: "spring",
          duration: 0.15,
          bounce: 0
        });
      }
      
      // Resetear el tiempo de inicio y otros estados
      dragStartTime.current = null;
      dragStartPosition.current = null;
      hasMoved.current = false;
    }
    
    return memo;
  }, {
    axis: 'x',
    filterTaps: true,
    bounds: { left: -CARD_WIDTH, right: CARD_WIDTH },
    rubberband: true,
    from: () => [x.get(), 0],
    threshold: MOVEMENT_THRESHOLD, // Umbral para detectar el inicio del arrastre
    delay: 0, // Sin retraso para iniciar el arrastre
    swipe: { distance: ACTIVATION_THRESHOLD, velocity: 0.4 }, // Configuración explícita de swipe
    pointer: { touch: true }, // Asegurar que funciona bien con eventos táctiles
  });

  // Función para manejar el toggle de detalles
  const handleToggleDetails = (e) => {
    // Si estamos arrastrando, no permitimos el toggle
    if (isDragging || hasMoved.current || Math.abs(lastDragPosition.current) > MOVEMENT_THRESHOLD) {
      console.log('Toggle detalles cancelado: estamos arrastrando', { 
        isDragging, 
        hasMoved: hasMoved.current,
        lastPosition: lastDragPosition.current 
      });
      return;
    }
    
    // Marcamos que estamos en proceso de cambiar el estado de detalles
    isTogglingDetails.current = true;
    console.log('Toggle detalles iniciado');
    
    // Cambiamos el estado de detalles
    setShowDetails(!showDetails);
    
    // Después de un tiempo, indicamos que ya no estamos cambiando el estado
    setTimeout(() => {
      isTogglingDetails.current = false;
      console.log('Toggle detalles completado');
    }, 300);
  };

  return (
    <div className="flex flex-col h-full w-full property-card-container">
      <div className="flex-1 relative w-full flex flex-col">
        <motion.div
          {...bind()}
          style={{ x, rotate }}
          className="relative bg-white dark:bg-gray-900/20 rounded-xl shadow-lg overflow-hidden flex-1 touch-none cursor-grab active:cursor-grabbing w-full"
          data-testid="property-card-draggable"
        >
          {/* Indicadores de Like/Dislike */}
          <motion.div 
            style={{ 
              scale: likeScale, 
              opacity: likeOpacity,
              borderColor: likeBorderColor
            }} 
            className="absolute top-8 right-8 z-10 border-4 text-green-500 px-8 py-2 rounded-full font-bold transform rotate-12"
          >
            LIKE
          </motion.div>
          <motion.div 
            style={{ 
              scale: dislikeScale, 
              opacity: nopeOpacity,
              borderColor: dislikeBorderColor
            }} 
            className="absolute top-8 left-8 z-10 border-4 text-red-500 px-8 py-2 rounded-full font-bold transform -rotate-12"
          >
            NOPE
          </motion.div>

          {/* Contenido de la card */}
          <div className="h-full overflow-y-auto select-none w-full">
              {!showDetails && (
            <div className={`relative ${showDetails ? 'h-[40vh]' : 'h-full'} transition-all duration-300 w-full`}>
              
                <>
                  <ImageCarousel property={property} />
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

            </div>
             )}
              <button 
                onClick={handleToggleDetails}
                className={`absolute top-4 right-4 bg-white/40 hover:bg-white/60 dark:bg-black/40 dark:hover:bg-black/60 rounded-full p-3 transition-all transform ${showDetails ? 'rotate-180' : ''} z-20`}
                data-testid="toggle-details-button"
                >
                <svg className="w-6 h-6" fill="none" stroke="white" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                </svg>
              </button>
             

            {showDetails && (
              <div className="w-full">
                <PropertyDetails property={property} />
              </div>
            )}
          </div>
        </motion.div>
      </div>

      {/* Botones de acción - Siempre visibles */}
      <motion.div 
        key="action-buttons"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.2 }}
        className="action-buttons py-2 flex justify-center gap-4 w-full mt-auto"
      >
        <button 
          className="rounded-full p-3 bg-gray-700/60 dark:bg-gray-800/80 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none" 
          disabled={!canUndo || isRating} 
          onClick={onUndo}
          data-testid="undo-button"
        >
          <svg className="w-7 h-7" fill="none" stroke="white" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
          </svg>
        </button>
        <button 
          className="rounded-full p-4 bg-gray-700/60 dark:bg-gray-800/80 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none" 
          disabled={isRating} 
          onClick={() => handleRate('dislike')}
          data-testid="dislike-button"
        >
          <svg className="w-8 h-8" fill="#FF4458" viewBox="0 0 24 24">
            <path d="M14.8 12l3.6-3.6c.8-.8.8-2 0-2.8-.8-.8-2-.8-2.8 0L12 9.2 8.4 5.6c-.8-.8-2-.8-2.8 0-.8.8-.8 2 0 2.8L9.2 12l-3.6 3.6c-.8.8-.8 2 0 2.8.4.4.9.6 1.4.6.5 0 1-.2 1.4-.6l3.6-3.6 3.6 3.6c.4.4.9.6 1.4.6.5 0 1-.2 1.4-.6.8-.8.8-2 0-2.8L14.8 12z"/>
          </svg>
        </button>
        <button 
          className="rounded-full p-4 bg-gray-700/60 dark:bg-gray-800/80 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none" 
          disabled={isRating} 
          onClick={() => handleRate('like')}
          data-testid="like-button"
        >
          <svg className="w-8 h-8" fill="#00DC7D" viewBox="0 0 24 24">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </button>
      </motion.div>
    </div>
  );
} 
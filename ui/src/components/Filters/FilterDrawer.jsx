import { motion, AnimatePresence, useMotionValue } from 'framer-motion';
import { useEffect, useRef, useState } from 'react';

const FilterDrawer = ({ 
    isOpen, 
    onClose, 
    title,
    options,
    selectedValue,
    onSelect,
    customContent
}) => {
    const contentRef = useRef(null);
    const [contentHeight, setContentHeight] = useState(0);
    const [windowHeight, setWindowHeight] = useState(0);

    // Calcular la altura del contenido y de la ventana cuando el drawer se abre
    useEffect(() => {
        if (isOpen && contentRef.current) {
            const updateHeights = () => {
                const contentHeight = contentRef.current.scrollHeight;
                const windowHeight = window.innerHeight;
                setContentHeight(contentHeight);
                setWindowHeight(windowHeight);
            };
            
            updateHeights();
            // Actualizar en caso de cambio de tamaño
            window.addEventListener('resize', updateHeights);
            return () => window.removeEventListener('resize', updateHeights);
        }
    }, [isOpen, options, customContent]);

    // Calcular la altura máxima del drawer (70% de la pantalla o el contenido si es menor)
    const drawerHeight = Math.min(contentHeight + 80, windowHeight * 0.7);
    
    return (
        <AnimatePresence>
            {isOpen && (
                <motion.div 
                    className="absolute inset-x-0 bottom-0 bg-gray-900 dark:bg-gray-900 rounded-t-2xl"
                    initial={{ y: "100%" }}
                    animate={{ y: 0 }}
                    exit={{ y: "100%" }}
                    transition={{ type: "spring", damping: 25, stiffness: 200 }}
                    style={{ 
                        zIndex: 10,
                        height: drawerHeight > 0 ? `${drawerHeight}px` : 'auto',
                        maxHeight: '70vh'
                    }}
                >
                    <div className="h-full flex flex-col">
                        {/* Indicador de arrastre */}
                        <div className="w-12 h-1 bg-gray-600 dark:bg-gray-600 rounded-full mx-auto mt-3 mb-2" />

                        {/* Header */}
                        <div className="p-4 border-b border-gray-700 dark:border-gray-700">
                            <div className="flex items-center gap-3">
                                <button 
                                    onClick={onClose}
                                    className="p-2 hover:bg-gray-800 dark:hover:bg-gray-800 rounded-full"
                                >
                                    <svg className="w-6 h-6 text-gray-300 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                    </svg>
                                </button>
                                <h3 className="text-xl font-bold text-white dark:text-white">{title}</h3>
                            </div>
                        </div>

                        {/* Contenido */}
                        <div className="flex-1 overflow-y-auto" ref={contentRef}>
                            {customContent ? customContent : (
                                <div className="p-4 space-y-2">
                                    {options?.map(option => (
                                        <button
                                            key={option.id}
                                            className={`w-full p-4 rounded-xl text-base font-medium transition-all ${
                                                selectedValue === option.id
                                                    ? 'bg-blue-500 text-white'
                                                    : 'bg-gray-800 dark:bg-gray-800 text-gray-300 dark:text-gray-300'
                                            }`}
                                            onClick={() => onSelect(option.id)}
                                        >
                                            {option.label}
                                        </button>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                </motion.div>
            )}
        </AnimatePresence>
    );
}

export default FilterDrawer; 
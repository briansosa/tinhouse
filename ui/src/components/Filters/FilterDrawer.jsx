import { motion, AnimatePresence } from 'framer-motion';

const FilterDrawer = ({ 
    isOpen, 
    onClose, 
    title,
    options,
    selectedValue,
    onSelect,
    customContent
}) => {
    return (
        <AnimatePresence>
            {isOpen && (
                <motion.div 
                    className="absolute inset-0 bg-gray-950 rounded-t-2xl"
                    initial={{ y: "100%" }}
                    animate={{ y: 0 }}
                    exit={{ y: "100%" }}
                    transition={{ type: "spring", damping: 25, stiffness: 200 }}
                    style={{ zIndex: 10 }}
                >
                    <div className="h-full flex flex-col">
                        {/* Indicador de arrastre */}
                        <div className="w-12 h-1 bg-gray-800 rounded-full mx-auto mt-3 mb-2" />

                        {/* Header */}
                        <div className="p-4 border-b border-gray-800">
                            <div className="flex items-center gap-3">
                                <button 
                                    onClick={onClose}
                                    className="p-2 hover:bg-gray-800 rounded-full"
                                >
                                    <svg className="w-6 h-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                    </svg>
                                </button>
                                <h3 className="text-xl font-bold text-white">{title}</h3>
                            </div>
                        </div>

                        {/* Contenido */}
                        <div className="flex-1 overflow-y-auto">
                            {customContent ? customContent : (
                                <div className="p-4 space-y-2">
                                    {options?.map(option => (
                                        <button
                                            key={option.id}
                                            className={`w-full p-4 rounded-xl text-base font-medium transition-all ${
                                                selectedValue === option.id
                                                    ? 'bg-blue-500 text-white'
                                                    : 'bg-gray-800 text-gray-300'
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
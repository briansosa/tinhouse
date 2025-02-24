import { motion } from 'framer-motion';

export default function LikedPropertyCard({ property }) {
    return (
        <motion.div 
            className="snap-start flex-shrink-0 w-28 cursor-pointer select-none"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
        >
            <div className="relative w-28 h-28 rounded-lg overflow-hidden">
                <img 
                    src={property.images[0]} 
                    alt={property.title}
                    className="w-full h-full object-cover pointer-events-none"
                    draggable="false"
                />
                <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-2">
                    <p className="text-white text-xs truncate pointer-events-none">
                        {property.location}
                    </p>
                </div>
            </div>
        </motion.div>
    );
} 
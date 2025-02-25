import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';

export default function LikedPropertyCard({ property, onClick }) {
    const [hasNotes, setHasNotes] = useState(false);

    useEffect(() => {
        const savedNotes = localStorage.getItem(`property-notes-${property.id}`);
        if (savedNotes) {
            const notes = JSON.parse(savedNotes);
            setHasNotes(notes.length > 0);
        }
    }, [property.id]);

    return (
        <motion.div 
            className="snap-start flex-shrink-0 w-28 cursor-pointer select-none"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={onClick}
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
                {hasNotes && (
                    <div className="absolute top-2 right-2 bg-blue-500 text-white p-1 rounded-full">
                        <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/>
                        </svg>
                    </div>
                )}
            </div>
        </motion.div>
    );
} 
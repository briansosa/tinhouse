import { Link, useLocation } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useRef, useEffect } from 'react';

export default function BottomNavBar({ show }) {
    const location = useLocation();
    const hasAnimated = useRef(false);
    
    // Definimos las variantes de animación
    const variants = {
        initial: { y: 100 },
        visible: { y: 0 },
        hidden: { y: 100, transition: { duration: 0 } }
    };

    return (
        <AnimatePresence mode="wait">
            <motion.div 
                key="bottomNav"
                variants={variants}
                initial={hasAnimated.current ? false : "initial"}
                animate={show ? "visible" : "hidden"}
                transition={hasAnimated.current ? { duration: 0 } : { type: "spring", stiffness: 300, damping: 30 }}
                onAnimationComplete={() => {
                    hasAnimated.current = true;
                }}
                className="absolute bottom-0 left-0 right-0 h-16 bg-gray-950 dark:bg-gray-950 border-t border-gray-800 dark:border-gray-800 z-10"
                style={{ maxWidth: "448px", margin: "0 auto" }}
            >
                <div className="flex justify-around items-center h-full px-2">
                    <Link 
                        to="/" 
                        className={`flex flex-col items-center transition-colors ${
                            location.pathname === '/' ? 'text-green-500' : 'text-gray-500 dark:text-gray-300 hover:text-green-500'
                        }`}
                    >
                        <svg className="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                        </svg>
                    </Link>

                    <Link 
                        to="/likes" 
                        className={`flex flex-col items-center transition-colors ${
                            location.pathname === '/likes' ? 'text-green-500' : 'text-gray-500 dark:text-gray-300 hover:text-green-500'
                        }`}
                    >
                        <svg className="w-7 h-7" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                        </svg>
                    </Link>

                    <Link 
                        to="/settings" 
                        className={`flex flex-col items-center text-gray-500 dark:text-gray-300 hover:text-green-500 transition-colors ${
                            location.pathname === '/settings' ? 'text-green-500' : 'text-gray-500 dark:text-gray-300 hover:text-green-500'
                        }`}
                    >
                        <svg className="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                        </svg>
                    </Link>
                </div>
            </motion.div>
        </AnimatePresence>
    );
} 
import { motion } from 'framer-motion';

export default function Layout({ children }) {
    return (
        <motion.div 
            className="relative bg-white dark:bg-gray-950 overflow-hidden h-full"
        >
            {children}
        </motion.div>
    );
} 
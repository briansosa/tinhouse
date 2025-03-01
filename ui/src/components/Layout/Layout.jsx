import { motion } from 'framer-motion';

export default function Layout({ children }) {
    return (
        <motion.div 
            className="flex flex-col h-full w-full bg-white dark:bg-gray-900/20 overflow-hidden"
        >
            {children}
        </motion.div>
    );
} 
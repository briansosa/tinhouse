import { Link, useLocation } from 'react-router-dom';

const NavItem = ({ to, icon: Icon, label }) => {
    const location = useLocation();
    const isActive = location.pathname === to;

    return (
        <Link 
            to={to} 
            className={`flex flex-col items-center p-2 transition-colors ${
                isActive 
                    ? 'text-green-500 dark:text-green-400' 
                    : 'text-gray-600 dark:text-gray-400 hover:text-green-500 dark:hover:text-green-400'
            }`}
        >
            <Icon className="w-6 h-6" />
            <span className="text-xs mt-1">{label}</span>
        </Link>
    );
};

export default function BottomNavBar() {
    return (
        <nav className="fixed bottom-0 left-0 right-0 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800 shadow-lg">
            <div className="max-w-md mx-auto flex justify-around items-center h-16 px-4">
                <NavItem 
                    to="/" 
                    icon={(props) => (
                        <svg {...props} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                        </svg>
                    )}
                    label="Inicio"
                />
                <NavItem 
                    to="/likes" 
                    icon={(props) => (
                        <svg {...props} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                        </svg>
                    )}
                    label="Favoritos"
                />
                <NavItem 
                    to="/filters" 
                    icon={(props) => (
                        <svg {...props} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                        </svg>
                    )}
                    label="Filtros"
                />
            </div>
        </nav>
    );
} 
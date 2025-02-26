import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Likes from './pages/Likes';
import Settings from './pages/Settings';
import Layout from './components/Layout/Layout';
import { useState, useEffect } from 'react';
import BottomNavBar from './components/BottomNavBar/BottomNavBar';

function App() {
  const [darkMode, setDarkMode] = useState(() => {
    // Chequeamos la preferencia guardada o del sistema
    return localStorage.getItem('theme') === 'dark' || 
      (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
  });

  const [showNavBar, setShowNavBar] = useState(true);

  useEffect(() => {
    // Actualizamos solo la clase en el html sin modificar localStorage
    if (darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [darkMode]);

  // Efecto para manejar el tema automático basado en la hora
  useEffect(() => {
    const themeMode = localStorage.getItem('themeMode');
    
    if (themeMode === 'auto') {
      const updateThemeByTime = () => {
        const hour = new Date().getHours();
        const isDarkHours = hour < 7 || hour > 19; // Oscuro entre 7pm y 7am
        setDarkMode(isDarkHours);
      };
      
      // Actualizar el tema inmediatamente
      updateThemeByTime();
      
      // Configurar un intervalo para verificar cada hora
      const interval = setInterval(updateThemeByTime, 60 * 60 * 1000);
      
      return () => clearInterval(interval);
    }
  }, []);

  return (
    <Router>
      <div className="min-h-screen bg-gray-950 dark:bg-gray-950">
        <div className="relative max-w-md mx-auto h-screen overflow-hidden">
          <Routes>
            <Route path="/" element={<Home setShowNavBar={setShowNavBar} />} />
            <Route path="/likes" element={<Likes setShowNavBar={setShowNavBar} />} />
            <Route path="/settings" element={<Settings setShowNavBar={setShowNavBar} />} />
          </Routes>
          <BottomNavBar show={showNavBar} />
        </div>
      </div>
    </Router>
  );
}

export default App;

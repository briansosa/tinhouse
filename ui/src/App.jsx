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

  // Efecto para calcular la altura real del viewport en dispositivos móviles
  useEffect(() => {
    const setVhVariable = () => {
      const vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty('--vh', `${vh}px`);
    };

    // Establecer el valor inicial
    setVhVariable();

    // Actualizar en cambios de tamaño o orientación
    window.addEventListener('resize', setVhVariable);
    window.addEventListener('orientationchange', setVhVariable);

    return () => {
      window.removeEventListener('resize', setVhVariable);
      window.removeEventListener('orientationchange', setVhVariable);
    };
  }, []);

  return (
    <Router>
      <div className={`app-container bg-gray-950 dark:bg-gray-950 ${!showNavBar ? 'no-navbar' : ''}`}>
        <div className="w-full h-full flex flex-col">
          <div className={`main-content ${!showNavBar ? 'pb-0' : ''}`}>
            <Routes>
              <Route path="/" element={<Home setShowNavBar={setShowNavBar} />} />
              <Route path="/likes" element={<Likes setShowNavBar={setShowNavBar} />} />
              <Route path="/settings" element={<Settings setShowNavBar={setShowNavBar} />} />
            </Routes>
          </div>
          {showNavBar && <BottomNavBar show={showNavBar} />}
        </div>
      </div>
    </Router>
  );
}

export default App;

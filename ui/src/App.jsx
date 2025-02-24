import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Likes from './pages/Likes';
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
    // Actualizamos la clase en el html y guardamos la preferencia
    if (darkMode) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }, [darkMode]);

  return (
    <Router>
      <div className="min-h-screen bg-gray-100 dark:bg-gray-950">
        <div className="relative max-w-md mx-auto h-screen">
          <Routes>
            <Route path="/" element={<Home setShowNavBar={setShowNavBar} />} />
            <Route path="/likes" element={<Likes setShowNavBar={setShowNavBar} />} />
          </Routes>
          <BottomNavBar show={showNavBar} />
        </div>

        {/* Botón para cambiar el tema */}
        <button 
          onClick={() => setDarkMode(!darkMode)}
          className="fixed top-4 right-4 p-2 rounded-full bg-gray-200 dark:bg-gray-700"
        >
          {darkMode ? '🌞' : '🌙'}
        </button>
      </div>
    </Router>
  );
}

export default App;

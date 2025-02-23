import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Likes from './pages/Likes';
import { useState, useEffect } from 'react';

function App() {
  const [darkMode, setDarkMode] = useState(() => {
    // Chequeamos la preferencia guardada o del sistema
    return localStorage.getItem('theme') === 'dark' || 
      (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
  });

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
      <div className="min-h-screen bg-gray-100 dark:bg-black transition-colors">
        <nav className="bg-white shadow-lg">
          <div className="max-w-6xl mx-auto px-4">
            <div className="flex justify-between">
              <div className="flex space-x-7">
                <div className="flex items-center py-4">
                  <Link to="/" className="text-lg font-semibold">FindHouse</Link>
                </div>
                <div className="flex items-center space-x-4">
                  <Link to="/" className="py-4 px-2 hover:text-gray-700">Home</Link>
                  <Link to="/likes" className="py-4 px-2 hover:text-gray-700">Likes</Link>
                </div>
              </div>
            </div>
          </div>
        </nav>

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/likes" element={<Likes />} />
        </Routes>

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

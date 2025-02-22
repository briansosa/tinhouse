import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './pages/Home';
import Likes from './pages/Likes';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
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
      </div>
    </Router>
  );
}

export default App;

import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './index.css'

// En desarrollo, StrictMode causa doble montaje/desmontaje
// Esto puede causar problemas con peticiones a la API
// Comentar la siguiente línea para habilitar StrictMode si es necesario
ReactDOM.createRoot(document.getElementById('root')).render(
  // <React.StrictMode>
    <App />
  // </React.StrictMode>,
)

import { useState } from 'react';
import { rateProperty } from '../../services/api';
import PropertyDetails from '../PropertyDetails/PropertyDetails';

export default function PropertyCard({ property, onRate }) {
  const [showDetails, setShowDetails] = useState(false);
  const [isRating, setIsRating] = useState(false);

  const handleRate = async (rating) => {
    if (isRating) return;
    
    setIsRating(true);
    try {
      await rateProperty(property.id, rating);
      onRate(rating);
    } catch (error) {
      console.error('Error rating property:', error);
    } finally {
      setIsRating(false);
    }
  };

  return (
    <div className="relative max-w-md mx-auto bg-white rounded-xl shadow-lg overflow-hidden">
      <div className={`relative transition-all duration-300 ${showDetails ? 'h-[40vh]' : 'h-[70vh]'}`}>
        <img 
          src={property.image_url || 'https://via.placeholder.com/400x300'} 
          alt={property.title}
          className="w-full h-full object-cover"
        />
        
        {!showDetails && (
          <>
            <div className="absolute bottom-0 left-0 right-0 h-48 bg-gradient-to-t from-black/90 to-transparent" />
            <div className="absolute bottom-2 left-0 right-0 px-4 text-white">
              <h2 className="text-2xl font-bold leading-tight">{property.title}</h2>
              <p className="text-sm mt-1 opacity-90">{property.location}</p>
              <p className="text-xl font-semibold mt-1">USD{property.price}</p>
              <div className="flex items-center gap-2 mt-1 text-sm opacity-80">
                <span>{property.details.bedrooms} dorm.</span>
                <span>•</span>
                <span>{property.details.bathrooms} baños</span>
                <span>•</span>
                <span>{property.details.area}m²</span>
              </div>
            </div>
          </>
        )}

        {/* Botón de más detalles */}
        <button 
          onClick={() => setShowDetails(!showDetails)}
          className={`absolute top-4 right-4 bg-white/30 hover:bg-white/50 rounded-full p-3 transition-all transform ${showDetails ? 'rotate-180' : ''}`}
        >
          <svg className="w-6 h-6" fill="none" stroke="white" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
          </svg>
        </button>
      </div>

      {/* Contenedor de detalles */}
      <div className={`transition-all duration-300 ${showDetails ? 'max-h-[60vh]' : 'max-h-0'} overflow-y-auto`}>
        {showDetails && <PropertyDetails property={property} />}
      </div>

      {/* Botones de acción en franja blanca */}
      <div className="flex justify-around items-center p-6 bg-white">
        <button 
          className="rounded-full p-4 hover:bg-gray-100 transition-colors disabled:opacity-50"
          disabled={isRating}
          onClick={() => console.log('undo')}
        >
          <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
          </svg>
        </button>
        <button 
          className="rounded-full p-4 hover:bg-gray-100 transition-colors disabled:opacity-50"
          disabled={isRating}
          onClick={() => handleRate('dislike')}
        >
          <svg className="w-12 h-12" fill="none" stroke="red" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <button 
          className="rounded-full p-4 hover:bg-gray-100 transition-colors disabled:opacity-50"
          disabled={isRating}
          onClick={() => handleRate('like')}
        >
          <svg className="w-12 h-12" fill="none" stroke="green" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
        </button>
      </div>
    </div>
  );
} 
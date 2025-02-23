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
    <div className="relative max-w-md mx-auto bg-white dark:bg-gray-950 rounded-xl shadow-lg overflow-hidden h-[85vh]">
      <div className="h-full overflow-y-auto">
        <div className={`relative ${showDetails ? 'h-[40vh]' : 'h-[73vh]'} transition-all duration-300`}>
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
                <p className="text-xl font-semibold mt-1">{property.price}</p>
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

          <button 
            onClick={() => setShowDetails(!showDetails)}
            className={`absolute top-4 right-4 bg-white/30 hover:bg-white/50 dark:bg-black/30 dark:hover:bg-black/50 rounded-full p-3 transition-all transform ${showDetails ? 'rotate-180' : ''}`}
          >
            <svg className="w-6 h-6" fill="none" stroke="white" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
            </svg>
          </button>
        </div>

        {showDetails && (
          <div className="pb-16">
            <PropertyDetails property={property} />
          </div>
        )}
      </div>

      <div className={`
        absolute bottom-0 left-0 right-0 
        ${showDetails ? 'h-16 bg-gradient-to-t from-white dark:from-black via-white dark:via-black to-transparent' : ''}
      `}>
        <div className={`
          flex justify-around items-center px-4
          ${showDetails 
            ? 'h-full pb-6'
            : 'h-24 py-6'
          }
        `}>
          <button 
            className="rounded-full p-3 bg-gray-800/50 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none"
            disabled={isRating}
            onClick={() => console.log('undo')}
          >
            <svg className="w-7 h-7" fill="none" stroke="white" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
          </button>
          <button 
            className="rounded-full p-4 bg-gray-800/50 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none"
            disabled={isRating}
            onClick={() => handleRate('dislike')}
          >
            <svg className="w-8 h-8" fill="#FF4458" viewBox="0 0 24 24">
              <path d="M14.8 12l3.6-3.6c.8-.8.8-2 0-2.8-.8-.8-2-.8-2.8 0L12 9.2 8.4 5.6c-.8-.8-2-.8-2.8 0-.8.8-.8 2 0 2.8L9.2 12l-3.6 3.6c-.8.8-.8 2 0 2.8.4.4.9.6 1.4.6.5 0 1-.2 1.4-.6l3.6-3.6 3.6 3.6c.4.4.9.6 1.4.6.5 0 1-.2 1.4-.6.8-.8.8-2 0-2.8L14.8 12z"/>
            </svg>
          </button>
          <button 
            className="rounded-full p-4 bg-gray-800/50 shadow-lg hover:shadow-xl transform hover:scale-110 transition-all duration-200 disabled:opacity-50 disabled:scale-100 disabled:shadow-none"
            disabled={isRating}
            onClick={() => handleRate('like')}
          >
            <svg className="w-8 h-8" fill="#00DC7D" viewBox="0 0 24 24">
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
} 
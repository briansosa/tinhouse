export default function PropertyDetails({ property }) {
  return (
    <div className="p-6 bg-white">
      {/* Detalles principales */}
      <div className="flex justify-around py-4 border-y border-gray-200">
        {property.details.bedrooms && (
          <div className="text-center">
            <p className="text-xl font-semibold">{property.details.bedrooms}</p>
            <p className="text-sm text-gray-600">Dormitorios</p>
          </div>
        )}
        {property.details.bathrooms && (
          <div className="text-center">
            <p className="text-xl font-semibold">{property.details.bathrooms}</p>
            <p className="text-sm text-gray-600">Baños</p>
          </div>
        )}
        {property.details.area && (
          <div className="text-center">
            <p className="text-xl font-semibold">{property.details.area}m²</p>
            <p className="text-sm text-gray-600">Superficie</p>
          </div>
        )}
      </div>

      {/* Ubicación */}
      <div className="mt-4">
        <h3 className="text-lg font-semibold">Ubicación</h3>
        <p className="text-gray-700">{property.location}</p>
      </div>

      {/* Descripción */}
      {property.description && (
        <div className="mt-4">
          <h3 className="text-lg font-semibold">Sobre la propiedad</h3>
          <p className="text-gray-700 whitespace-pre-line">{property.description}</p>
        </div>
      )}

      {/* Botón de ver más */}
      <a 
        href={property.url} 
        target="_blank" 
        rel="noopener noreferrer"
        className="block w-full text-center bg-gradient-to-r from-blue-500 to-blue-600 text-white py-3 rounded-lg hover:from-blue-600 hover:to-blue-700 transition-all mt-6"
      >
        Ver en {property.agency?.name || 'sitio web'}
      </a>
    </div>
  );
} 
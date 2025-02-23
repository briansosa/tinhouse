export default function PropertyDetails({ property }) {
  return (
    <div className="p-6 bg-white dark:bg-gray-950 space-y-6">
      {/* Información principal */}
      <div>
        <span className="text-4xl font-bold dark:text-white">{property.price}</span>
        <div className="flex items-center gap-2 mt-2">
          <p className="text-gray-600 dark:text-gray-300">{property.location}</p>
          {property.property_type && (
            <>
              <span className="text-gray-400 dark:text-gray-500">•</span>
              <p className="text-gray-600 dark:text-gray-300">{property.property_type}</p>
            </>
          )}
        </div>
        {property.agency?.name && (
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{property.agency.name}</p>
        )}
      </div>

      {/* Características destacadas */}
      <div className="grid grid-cols-3 gap-4 py-4 border-y border-gray-200 dark:border-gray-700">
        <div className="text-center">
          <p className="text-3xl font-semibold dark:text-white">{property.details.bedrooms}</p>
          <p className="text-sm text-gray-600 dark:text-gray-400">Dormitorios</p>
        </div>
        <div className="text-center">
          <p className="text-3xl font-semibold dark:text-white">{property.details.bathrooms}</p>
          <p className="text-sm text-gray-600 dark:text-gray-400">Baños</p>
        </div>
        <div className="text-center">
          <p className="text-3xl font-semibold dark:text-white">{property.details.total_area}</p>
          <p className="text-sm text-gray-600 dark:text-gray-400">m² Totales</p>
        </div>
      </div>

      {/* Detalles del inmueble */}
      <div>
        <h3 className="text-2xl font-bold mb-4 dark:text-white">Detalles del inmueble</h3>
        
        {/* A. Información General */}
        <div className="mb-6">
          <h4 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-3">Información General</h4>
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Situación</span>
              <span className="font-medium dark:text-gray-200">{property.details.status || 'Sin especificar'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Antigüedad</span>
              <span className="font-medium dark:text-gray-200">{property.details.age || 'Sin especificar'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Plantas</span>
              <span className="font-medium dark:text-gray-200">{property.details.floors || '-'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Ambientes</span>
              <span className="font-medium dark:text-gray-200">{property.details.rooms || '-'}</span>
            </div>
          </div>
        </div>

        {/* B. Superficies */}
        <div className="mb-6">
          <h4 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-3">Superficies</h4>
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Superficie Total</span>
              <span className="font-medium dark:text-gray-200">{property.details.total_area ? `${property.details.total_area}m²` : '-'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Superficie Cubierta</span>
              <span className="font-medium dark:text-gray-200">{property.details.area ? `${property.details.area}m²` : '-'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Superficie Terreno</span>
              <span className="font-medium dark:text-gray-200">{property.details.land_area ? `${property.details.land_area}m²` : '-'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Frente x Fondo</span>
              <span className="font-medium dark:text-gray-200">
                {property.details.front_size && property.details.back_size 
                  ? `${property.details.front_size}m x ${property.details.back_size}m`
                  : '-'}
              </span>
            </div>
          </div>
        </div>

        {/* C. Características Adicionales */}
        <div>
          <h4 className="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-3">Características Adicionales</h4>
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Cocheras</span>
              <span className="font-medium dark:text-gray-200">{property.details.garages || 'Sin cochera'}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Expensas</span>
              <span className="font-medium dark:text-gray-200">
                {property.details.expenses ? `$${property.details.expenses}` : 'Sin expensas'}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600 dark:text-gray-400">Código</span>
              <span className="font-medium dark:text-gray-200">{property.code || '-'}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Descripción */}
      {property.description && (
        <div>
          <h3 className="text-2xl font-bold mb-3 dark:text-white">Descripción</h3>
          <p className="text-gray-700 dark:text-gray-300 leading-relaxed">{property.description}</p>
        </div>
      )}

      {/* Botón de contacto */}
      <a 
        href={property.url} 
        target="_blank" 
        rel="noopener noreferrer"
        className="block w-full text-center bg-gradient-to-r from-blue-500 to-blue-600 text-white py-4 rounded-lg hover:from-blue-600 hover:to-blue-700 transition-all font-medium"
      >
        Ver en {property.agency?.name || 'sitio web'}
      </a>
    </div>
  );
} 
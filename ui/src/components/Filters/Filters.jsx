import { useState } from 'react';
import FilterDrawer from './FilterDrawer';
import PriceRangeFilter from './PriceRangeFilter';
import LocationFilter from './LocationFilter';
import FeaturesFilter from './FeaturesFilter';
import SizeRangeFilter from './SizeRangeFilter';
import RoomsFilter from './RoomsFilter';
import AntiquityFilter from './AntiquityFilter';

const Filters = ({ onClose, onApplyFilters, initialFilters }) => {
    const [filters, setFilters] = useState(initialFilters || {
        propertyType: 'all',
        showOnlyWithNotes: false,
        showOnlyFavorites: false,
        priceRange: {
            min: null,
            max: null,
            currency: 'ARS'
        },
        locations: [],
        features: [],
        sizeRange: {
            min: null,
            max: null
        },
        rooms: null,
        bathrooms: null,
        antiquity: null
    });

    const [showPropertyTypeDrawer, setShowPropertyTypeDrawer] = useState(false);
    const [showPriceRangeDrawer, setShowPriceRangeDrawer] = useState(false);
    const [showLocationDrawer, setShowLocationDrawer] = useState(false);
    const [showFeaturesDrawer, setShowFeaturesDrawer] = useState(false);
    const [showSizeRangeDrawer, setShowSizeRangeDrawer] = useState(false);
    const [showRoomsDrawer, setShowRoomsDrawer] = useState(false);
    const [showBathroomsDrawer, setShowBathroomsDrawer] = useState(false);
    const [showAntiquityDrawer, setShowAntiquityDrawer] = useState(false);

    const propertyTypes = [
        { id: 'all', label: 'Todas' },
        { id: 'house', label: 'Casa' },
        { id: 'apartment', label: 'Departamento' },
        { id: 'ph', label: 'PH' }
    ];

    const selectedPropertyType = propertyTypes.find(type => type.id === filters.propertyType);

    const formatPrice = (price) => {
        if (price === null) return 'Sin límite';
        return filters.priceRange.currency === 'ARS'
            ? `$${new Intl.NumberFormat('es-AR').format(price)}`
            : `USD ${new Intl.NumberFormat('en-US').format(price)}`;
    };

    const getPriceRangeLabel = () => {
        const { min, max } = filters.priceRange;
        if (!min && !max) return 'Cualquier precio';
        if (!max) return `Desde ${formatPrice(min)}`;
        if (!min) return `Hasta ${formatPrice(max)}`;
        return `${formatPrice(min)} - ${formatPrice(max)}`;
    };

    const getSizeRangeLabel = () => {
        const { min, max } = filters.sizeRange;
        if (!min && !max) return 'Cualquier tamaño';
        if (!max) return `Desde ${min} m²`;
        if (!min) return `Hasta ${max} m²`;
        return `${min} - ${max} m²`;
    };

    const resetFilters = () => {
        setFilters({
            propertyType: 'all',
            showOnlyWithNotes: false,
            showOnlyFavorites: false,
            priceRange: {
                min: null,
                max: null,
                currency: 'ARS'
            },
            locations: [],
            features: [],
            sizeRange: {
                min: null,
                max: null
            },
            rooms: null,
            bathrooms: null,
            antiquity: null
        });
    };

    // Verificar si hay filtros activos
    const hasActiveFilters = () => {
        return (
            filters.propertyType !== 'all' ||
            filters.showOnlyWithNotes ||
            filters.showOnlyFavorites ||
            filters.priceRange.min !== null ||
            filters.priceRange.max !== null ||
            filters.locations.length > 0 ||
            filters.features.length > 0 ||
            filters.sizeRange.min !== null ||
            filters.sizeRange.max !== null ||
            filters.rooms !== null ||
            filters.bathrooms !== null ||
            filters.antiquity !== null
        );
    };

    return (
        <div className="h-full flex flex-col bg-gray-950 dark:bg-gray-950">
            {/* Header */}
            <div className="p-4 border-b border-gray-800 dark:border-gray-800">
                <div className="flex items-center justify-between">
                    <button 
                        onClick={onClose}
                        className="p-2 hover:bg-gray-800 dark:hover:bg-gray-800 rounded-full"
                    >
                        <svg className="w-6 h-6 text-gray-300 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                    <h1 className="text-xl font-bold text-white dark:text-white">Filtros</h1>
                    <button 
                        onClick={() => onApplyFilters(filters)}
                        className="px-3 py-1 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg"
                    >
                        OK
                    </button>
                </div>
            </div>

            {/* Contenido */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4 pb-16">
                {/* Tipo de propiedad */}
                <button 
                    onClick={() => setShowPropertyTypeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Tipo de propiedad</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{selectedPropertyType?.label || 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Rango de precio */}
                <button 
                    onClick={() => setShowPriceRangeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Precio</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{getPriceRangeLabel()}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Ubicación */}
                <button 
                    onClick={() => setShowLocationDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Ubicación</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.locations.length ? `${filters.locations.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Tamaño */}
                <button 
                    onClick={() => setShowSizeRangeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Tamaño</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{getSizeRangeLabel()}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Ambientes */}
                <button 
                    onClick={() => setShowRoomsDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Ambientes</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.rooms ? `${filters.rooms}+` : 'Todos'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Baños */}
                <button 
                    onClick={() => setShowBathroomsDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Baños</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.bathrooms ? `${filters.bathrooms}+` : 'Todos'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Características */}
                <button 
                    onClick={() => setShowFeaturesDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Características</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.features.length ? `${filters.features.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Antigüedad */}
                <button 
                    onClick={() => setShowAntiquityDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Antigüedad</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.antiquity !== null ? (
                            filters.antiquity === 0 ? 'A estrenar' :
                            filters.antiquity === 100 ? 'Más de 30 años' :
                            `Hasta ${filters.antiquity} años`
                        ) : 'Cualquiera'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Otros filtros */}
                <div className="mt-8">
                    <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase mb-4">
                        OTROS FILTROS
                    </h3>
                    <div className="space-y-4">
                        <div className="flex items-center justify-between p-4 rounded-xl bg-gray-800 dark:bg-gray-800">
                            <span className="text-md text-gray-200 dark:text-gray-200">Solo propiedades con notas</span>
                            <button 
                                onClick={() => setFilters(prev => ({ 
                                    ...prev, 
                                    showOnlyWithNotes: !prev.showOnlyWithNotes 
                                }))}
                                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-300 ${
                                    filters.showOnlyWithNotes ? 'bg-blue-500' : 'bg-gray-600'
                                }`}
                            >
                                <span 
                                    className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-lg transition-transform duration-300 ${
                                        filters.showOnlyWithNotes ? 'translate-x-6' : 'translate-x-1'
                                    }`}
                                />
                            </button>
                        </div>
                        
                        <div className="flex items-center justify-between p-4 rounded-xl bg-gray-800 dark:bg-gray-800">
                            <span className="text-md text-gray-200 dark:text-gray-200">Solo propiedades favoritas</span>
                            <button 
                                onClick={() => setFilters(prev => ({ 
                                    ...prev, 
                                    showOnlyFavorites: !prev.showOnlyFavorites 
                                }))}
                                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-300 ${
                                    filters.showOnlyFavorites ? 'bg-blue-500' : 'bg-gray-600'
                                }`}
                            >
                                <span 
                                    className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-lg transition-transform duration-300 ${
                                        filters.showOnlyFavorites ? 'translate-x-6' : 'translate-x-1'
                                    }`}
                                />
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Drawers */}
            <FilterDrawer 
                isOpen={showPropertyTypeDrawer}
                onClose={() => setShowPropertyTypeDrawer(false)}
                title="TIPO DE PROPIEDAD"
                options={propertyTypes}
                selectedValue={filters.propertyType}
                onSelect={(value) => {
                    setFilters(prev => ({ ...prev, propertyType: value }));
                    setShowPropertyTypeDrawer(false);
                }}
            />

            <FilterDrawer 
                isOpen={showPriceRangeDrawer}
                onClose={() => setShowPriceRangeDrawer(false)}
                title="PRECIO"
                customContent={
                    <PriceRangeFilter
                        initialRange={filters.priceRange}
                        onChange={(range) => setFilters(prev => ({
                            ...prev,
                            priceRange: range
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showLocationDrawer}
                onClose={() => setShowLocationDrawer(false)}
                title="UBICACIÓN"
                customContent={
                    <LocationFilter
                        initialLocations={filters.locations}
                        onChange={(locations) => setFilters(prev => ({
                            ...prev,
                            locations
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showSizeRangeDrawer}
                onClose={() => setShowSizeRangeDrawer(false)}
                title="TAMAÑO"
                customContent={
                    <SizeRangeFilter
                        initialRange={filters.sizeRange}
                        onChange={(range) => setFilters(prev => ({
                            ...prev,
                            sizeRange: range
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showRoomsDrawer}
                onClose={() => setShowRoomsDrawer(false)}
                title="AMBIENTES"
                customContent={
                    <RoomsFilter
                        initialValue={filters.rooms}
                        onChange={(value) => setFilters(prev => ({
                            ...prev,
                            rooms: value
                        }))}
                        title="AMBIENTES"
                    />
                }
            />

            <FilterDrawer 
                isOpen={showBathroomsDrawer}
                onClose={() => setShowBathroomsDrawer(false)}
                title="BAÑOS"
                customContent={
                    <RoomsFilter
                        initialValue={filters.bathrooms}
                        onChange={(value) => setFilters(prev => ({
                            ...prev,
                            bathrooms: value
                        }))}
                        title="BAÑOS"
                    />
                }
            />

            <FilterDrawer 
                isOpen={showFeaturesDrawer}
                onClose={() => setShowFeaturesDrawer(false)}
                title="CARACTERÍSTICAS"
                customContent={
                    <FeaturesFilter
                        initialFeatures={filters.features}
                        onChange={(features) => setFilters(prev => ({
                            ...prev,
                            features
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showAntiquityDrawer}
                onClose={() => setShowAntiquityDrawer(false)}
                title="ANTIGÜEDAD"
                customContent={
                    <AntiquityFilter
                        initialValue={filters.antiquity}
                        onChange={(value) => setFilters(prev => ({
                            ...prev,
                            antiquity: value
                        }))}
                    />
                }
            />
        </div>
    );
}

export default Filters;
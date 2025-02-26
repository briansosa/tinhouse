import { useState } from 'react';
import FilterDrawer from './FilterDrawer';
import PriceRangeFilter from './PriceRangeFilter';
import LocationFilter from './LocationFilter';
import FeaturesFilter from './FeaturesFilter';
import SizeRangeFilter from './SizeRangeFilter';
import RoomsFilter from './RoomsFilter';
import AntiquityFilter from './AntiquityFilter';
import FilterChips from './FilterChips';

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

    const handleRemoveFilter = (filterName, defaultValue) => {
        setFilters(prev => ({
            ...prev,
            [filterName]: defaultValue
        }));
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
        <div className="h-full dark:bg-gray-950 relative flex flex-col">
            {/* Header */}
            <div className="flex items-center justify-between p-4 border-b border-gray-800">
                <div className="flex items-center gap-3">
                    <button 
                        onClick={onClose}
                        className="p-2 hover:bg-gray-800 rounded-full"
                    >
                        <svg className="w-6 h-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                        </svg>
                    </button>
                    <h2 className="text-xl font-bold text-white">Filtros</h2>
                </div>
                {hasActiveFilters() && (
                    <button 
                        onClick={resetFilters}
                        className="text-blue-400 text-sm font-medium"
                    >
                        Limpiar todo
                    </button>
                )}
            </div>

            {/* Chips de filtros activos */}
            <FilterChips 
                filters={filters} 
                onRemove={handleRemoveFilter} 
            />

            {/* Contenido principal */}
            <div className="flex-1 overflow-y-auto">
                <div className="p-4 space-y-8">
                    <div>
                        <button 
                            onClick={() => setShowPropertyTypeDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Tipo de propiedad</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {selectedPropertyType?.label || 'Todas'}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Rango de precio */}
                    <div>
                        <button 
                            onClick={() => setShowPriceRangeDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Precio</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {getPriceRangeLabel()}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Ubicación */}
                    <div>
                        <button 
                            onClick={() => setShowLocationDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Ubicación</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {filters.locations.length ? `${filters.locations.length} seleccionados` : 'Todas'}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Tamaño */}
                    <div>
                        <button 
                            onClick={() => setShowSizeRangeDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Tamaño</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {getSizeRangeLabel()}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Ambientes */}
                    <div>
                        <button 
                            onClick={() => setShowRoomsDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Ambientes</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {filters.rooms ? `${filters.rooms}+` : 'Todos'}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Baños */}
                    <div>
                        <button 
                            onClick={() => setShowBathroomsDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Baños</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {filters.bathrooms ? `${filters.bathrooms}+` : 'Todos'}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Características */}
                    <div>
                        <button 
                            onClick={() => setShowFeaturesDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Características</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {filters.features.length ? `${filters.features.length} seleccionados` : 'Todas'}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Antigüedad */}
                    <div>
                        <button 
                            onClick={() => setShowAntiquityDrawer(true)}
                            className="w-full p-4 rounded-xl bg-gray-800 flex items-center justify-between"
                        >
                            <span className="text-md text-gray-300">Antigüedad</span>
                            <div className="flex items-center gap-2">
                                <span className="text-base text-white font-medium">
                                    {filters.antiquity !== null ? (
                                        filters.antiquity === 0 ? 'A estrenar' :
                                        filters.antiquity === 100 ? 'Más de 30 años' :
                                        `Hasta ${filters.antiquity} años`
                                    ) : 'Cualquiera'}
                                </span>
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>

                    {/* Otros filtros */}
                    <div>
                        <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase mb-4">
                            OTROS FILTROS
                        </h3>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between p-4 rounded-xl bg-gray-800">
                                <span className="text-md text-gray-300">Solo propiedades con notas</span>
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
                            
                            <div className="flex items-center justify-between p-4 rounded-xl bg-gray-800">
                                <span className="text-md text-gray-300">Solo propiedades favoritas</span>
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
            </div>

            {/* Botón de aplicar */}
            <div className="p-4 bg-gray-950 border-t border-gray-800">
                <button 
                    onClick={() => onApplyFilters(filters)}
                    className="w-full p-4 bg-blue-500 hover:bg-blue-600 text-white font-semibold rounded-xl shadow-md transition-colors"
                >
                    Aplicar
                </button>
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
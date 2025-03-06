import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import FilterDrawer from '../components/Filters/FilterDrawer';
import PriceRangeFilter from '../components/Filters/PriceRangeFilter';
import LocationFilter from '../components/Filters/LocationFilter';
import FeaturesFilter from '../components/Filters/FeaturesFilter';
import SurfaceFilter from '../components/Filters/SurfaceFilter';
import RoomsFilter from '../components/Filters/RoomsFilter';
import AntiquityFilter from '../components/Filters/AntiquityFilter';
import DispositionFilter from '../components/Filters/DispositionFilter';
import PropertyTypeFilter from '../components/Filters/PropertyTypeFilter';
import OrientationFilter from '../components/Filters/OrientationFilter';
import { getPropertyTypes } from '../services/api';

export default function Settings({ setShowNavBar }) {
    const navigate = useNavigate();
    const [filters, setFilters] = useState(() => {
        // Intentar cargar los filtros guardados del localStorage
        const savedFilters = localStorage.getItem('globalFilters');
        const defaultFilters = {
            propertyType: ['all'],
            showOnlyWithNotes: false,
            showOnlyFavorites: false,
            priceRange: {
                min: null,
                max: null,
                currency: 'ARS'
            },
            locations: [],
            features: [],
            surface: {
                totalArea: { min: null, max: null },
                coveredArea: { min: null, max: null },
                landArea: { min: null, max: null },
                front: null,
                back: null
            },
            rooms: null,
            bathrooms: null,
            antiquity: null,
            disposition: [],
            orientation: []
        };
        
        if (savedFilters) {
            const parsedFilters = JSON.parse(savedFilters);
            // Convertir propertyType a array si es string (compatibilidad con versión anterior)
            if (parsedFilters.propertyType && !Array.isArray(parsedFilters.propertyType)) {
                parsedFilters.propertyType = [parsedFilters.propertyType];
            }
            
            // Compatibilidad con versión anterior (sizeRange a surface)
            if (parsedFilters.sizeRange && !parsedFilters.surface) {
                parsedFilters.surface = {
                    totalArea: parsedFilters.sizeRange,
                    coveredArea: { min: null, max: null },
                    landArea: { min: null, max: null },
                    front: null,
                    back: null
                };
            }
            
            return parsedFilters;
        }
        
        return defaultFilters;
    });

    // Estado para el tema
    const [themeMode, setThemeMode] = useState(() => {
        const savedTheme = localStorage.getItem('themeMode');
        return savedTheme || 'auto';
    });

    const [showPropertyTypeDrawer, setShowPropertyTypeDrawer] = useState(false);
    const [showPriceRangeDrawer, setShowPriceRangeDrawer] = useState(false);
    const [showLocationDrawer, setShowLocationDrawer] = useState(false);
    const [showFeaturesDrawer, setShowFeaturesDrawer] = useState(false);
    const [showSizeRangeDrawer, setShowSizeRangeDrawer] = useState(false);
    const [showRoomsDrawer, setShowRoomsDrawer] = useState(false);
    const [showBathroomsDrawer, setShowBathroomsDrawer] = useState(false);
    const [showAntiquityDrawer, setShowAntiquityDrawer] = useState(false);
    const [showDispositionDrawer, setShowDispositionDrawer] = useState(false);
    const [showThemeDrawer, setShowThemeDrawer] = useState(false);
    const [showOrientationDrawer, setShowOrientationDrawer] = useState(false);

    const [propertyTypeLabels, setPropertyTypeLabels] = useState({});

    useEffect(() => {
        setShowNavBar(false);
    }, [setShowNavBar]);

    // Cargar tipos de propiedad para mostrar etiquetas correctas
    useEffect(() => {
        const fetchPropertyTypes = async () => {
            try {
                const response = await getPropertyTypes();
                if (response.data && response.data.length > 0) {
                    const typeMap = {};
                    response.data.forEach(type => {
                        typeMap[type.id] = type.name;
                    });
                    setPropertyTypeLabels(typeMap);
                }
            } catch (error) {
                console.error('Error al cargar tipos de propiedad:', error);
            }
        };
        
        fetchPropertyTypes();
    }, []);

    const themeOptions = [
        { id: 'light', label: 'Claro' },
        { id: 'dark', label: 'Oscuro' },
        { id: 'auto', label: 'Automático (según hora)' }
    ];

    const selectedTheme = themeOptions.find(theme => theme.id === themeMode);

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

    const getSurfaceLabel = () => {
        // Verificar si hay algún filtro de superficie activo
        const { totalArea, coveredArea, landArea, front, back } = filters.surface;
        
        if ((!totalArea.min && !totalArea.max) && 
            (!coveredArea.min && !coveredArea.max) && 
            (!landArea.min && !landArea.max) && 
            !front && !back) {
            return 'Cualquier superficie';
        }
        
        const labels = [];
        
        // Superficie Total
        if (totalArea.min || totalArea.max) {
            if (!totalArea.max) labels.push(`Total: desde ${totalArea.min} m²`);
            else if (!totalArea.min) labels.push(`Total: hasta ${totalArea.max} m²`);
            else labels.push(`Total: ${totalArea.min} - ${totalArea.max} m²`);
        }
        
        // Superficie Cubierta
        if (coveredArea.min || coveredArea.max) {
            if (!coveredArea.max) labels.push(`Cubierta: desde ${coveredArea.min} m²`);
            else if (!coveredArea.min) labels.push(`Cubierta: hasta ${coveredArea.max} m²`);
            else labels.push(`Cubierta: ${coveredArea.min} - ${coveredArea.max} m²`);
        }
        
        // Superficie Terreno
        if (landArea.min || landArea.max) {
            if (!landArea.max) labels.push(`Terreno: desde ${landArea.min} m²`);
            else if (!landArea.min) labels.push(`Terreno: hasta ${landArea.max} m²`);
            else labels.push(`Terreno: ${landArea.min} - ${landArea.max} m²`);
        }
        
        // Frente y Fondo
        if (front || back) {
            const frontText = front ? `${front}m` : '-';
            const backText = back ? `${back}m` : '-';
            labels.push(`Frente x Fondo: ${frontText} x ${backText}`);
        }
        
        return labels.join(', ');
    };

    const resetFilters = () => {
        const defaultFilters = {
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
            surface: {
                totalArea: { min: null, max: null },
                coveredArea: { min: null, max: null },
                landArea: { min: null, max: null },
                front: null,
                back: null
            },
            rooms: null,
            bathrooms: null,
            antiquity: null,
            disposition: [],
            orientation: []
        };
        setFilters(defaultFilters);
        localStorage.setItem('globalFilters', JSON.stringify(defaultFilters));
    };

    const saveFilters = () => {
        localStorage.setItem('globalFilters', JSON.stringify(filters));
        localStorage.setItem('themeMode', themeMode);
        
        // Aplicar el tema según la selección
        const hour = new Date().getHours();
        const isDarkHours = hour < 7 || hour > 19; // Oscuro entre 7pm y 7am
        
        // Solo actualizamos el tema actual (theme) basado en el modo seleccionado (themeMode)
        if (themeMode === 'light') {
            localStorage.setItem('theme', 'light');
        } else if (themeMode === 'dark') {
            localStorage.setItem('theme', 'dark');
        } else if (themeMode === 'auto') {
            localStorage.setItem('theme', isDarkHours ? 'dark' : 'light');
        }
        
        // Recargar la página para aplicar los cambios de tema
        window.location.href = '/';
    };

    return (
        <div className="h-full flex flex-col bg-gray-950 dark:bg-gray-950">
            {/* Header */}
            <div className="p-4 border-b border-gray-800 dark:border-gray-800">
                <div className="flex items-center justify-between">
                    <h1 className="text-xl font-bold text-white dark:text-white">Configuración</h1>
                    <button 
                        onClick={saveFilters}
                        className="px-3 py-1 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg"
                    >
                        OK
                    </button>
                </div>
            </div>

            {/* Contenido */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4 pb-28">
                {/* Tipo de propiedad */}
                <button 
                    onClick={() => setShowPropertyTypeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-300 dark:text-gray-300">Tipo de propiedad</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">
                            {filters.propertyType.length > 1 ? 'Múltiples' : filters.propertyType[0] === 'all' ? 'Todas' : 
                             propertyTypeLabels[filters.propertyType[0]] || filters.propertyType[0]}
                        </span>
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
                    <span className="text-gray-300 dark:text-gray-300">Precio</span>
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
                    <span className="text-gray-300 dark:text-gray-300">Ubicación</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.locations.length ? `${filters.locations.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Superficie */}
                <button 
                    onClick={() => setShowSizeRangeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-300 dark:text-gray-300">Superficie</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{getSurfaceLabel()}</span>
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
                    <span className="text-gray-300 dark:text-gray-300">Ambientes</span>
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
                    <span className="text-gray-300 dark:text-gray-300">Baños</span>
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
                    <span className="text-gray-300 dark:text-gray-300">Características</span>
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
                    <span className="text-gray-300 dark:text-gray-300">Antigüedad</span>
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

                {/* Disposición */}
                <button 
                    onClick={() => setShowDispositionDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-300 dark:text-gray-300">Disposición</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.disposition && filters.disposition.length ? `${filters.disposition.length} seleccionados` : 'Cualquiera'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Orientación */}
                <button 
                    onClick={() => setShowOrientationDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-300 dark:text-gray-300">Orientación</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.orientation && filters.orientation.length ? `${filters.orientation.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Ajustes generales */}
                <div className="mt-8">
                    <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase mb-4">
                        AJUSTES GENERALES
                    </h3>
                    <div className="space-y-4">
                        {/* Tema */}
                        <button 
                            onClick={() => setShowThemeDrawer(true)}
                            className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                        >
                            <span className="text-gray-300 dark:text-gray-300">Tema</span>
                            <div className="flex items-center">
                                <span className="text-gray-400 dark:text-gray-400 mr-2">{selectedTheme?.label || 'Automático'}</span>
                                <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </button>
                    </div>
                </div>
            </div>

            {/* Drawers */}
            <FilterDrawer 
                isOpen={showPropertyTypeDrawer}
                onClose={() => setShowPropertyTypeDrawer(false)}
                title="TIPO DE PROPIEDAD"
                customContent={
                    <PropertyTypeFilter
                        initialValue={filters.propertyType}
                        onChange={(value) => setFilters(prev => ({
                            ...prev,
                            propertyType: value
                        }))}
                    />
                }
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
                title="SUPERFICIE"
                customContent={
                    <SurfaceFilter
                        initialValues={filters.surface}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            surface: values
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

            <FilterDrawer 
                isOpen={showDispositionDrawer}
                onClose={() => setShowDispositionDrawer(false)}
                title="DISPOSICIÓN"
                customContent={
                    <DispositionFilter
                        initialValues={filters.disposition}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            disposition: values
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showOrientationDrawer}
                onClose={() => setShowOrientationDrawer(false)}
                title="ORIENTACIÓN"
                customContent={
                    <OrientationFilter
                        initialValues={filters.orientation || []}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            orientation: values
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showThemeDrawer}
                onClose={() => setShowThemeDrawer(false)}
                title="TEMA"
                options={themeOptions}
                selectedValue={themeMode}
                onSelect={(value) => {
                    setThemeMode(value);
                    setShowThemeDrawer(false);
                }}
            />
        </div>
    );
} 
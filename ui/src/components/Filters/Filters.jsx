import { useState, useEffect } from 'react';
import FilterDrawer from './FilterDrawer';
import PriceRangeFilter from './PriceRangeFilter';
import LocationFilter from './LocationFilter';
import FeaturesFilter from './FeaturesFilter';
import SurfaceFilter from './SurfaceFilter';
import RoomsFilter from './RoomsFilter';
import BathroomsFilter from './BathroomsFilter';
import AntiquityFilter from './AntiquityFilter';
import DispositionFilter from './DispositionFilter';
import OrientationFilter from './OrientationFilter';
import ConditionFilter from './ConditionFilter';
import OperationTypeFilter from './OperationTypeFilter';
import SituationFilter from './SituationFilter';
import FilterChips from './FilterChips';
import PropertyTypeFilter from './PropertyTypeFilter';
import { getPropertyTypes } from '../../services/api';

const Filters = ({ onClose, onApplyFilters, initialFilters }) => {
    const [filters, setFilters] = useState({
        propertyType: ['all'],
        showOnlyWithNotes: false,
        showOnlyFavorites: false,
        priceRange: {
            min: null,
            max: null,
            currency: 'USD'
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
        orientation: [],
        condition: [],
        operationType: [],
        situation: []
    });

    const [propertyTypeLabels, setPropertyTypeLabels] = useState({});
    const [showPropertyTypeDrawer, setShowPropertyTypeDrawer] = useState(false);
    const [showPriceRangeDrawer, setShowPriceRangeDrawer] = useState(false);
    const [showLocationDrawer, setShowLocationDrawer] = useState(false);
    const [showFeaturesDrawer, setShowFeaturesDrawer] = useState(false);
    const [showSizeRangeDrawer, setShowSizeRangeDrawer] = useState(false);
    const [showRoomsDrawer, setShowRoomsDrawer] = useState(false);
    const [showBathroomsDrawer, setShowBathroomsDrawer] = useState(false);
    const [showAntiquityDrawer, setShowAntiquityDrawer] = useState(false);
    const [showDispositionDrawer, setShowDispositionDrawer] = useState(false);
    const [showOrientationDrawer, setShowOrientationDrawer] = useState(false);
    const [showConditionDrawer, setShowConditionDrawer] = useState(false);
    const [showOperationTypeDrawer, setShowOperationTypeDrawer] = useState(false);
    const [showSituationDrawer, setShowSituationDrawer] = useState(false);

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

    // Inicializar filtros con los valores iniciales
    useEffect(() => {
        if (initialFilters) {
            // Compatibilidad con versión anterior (sizeRange a surface)
            if (initialFilters.sizeRange && !initialFilters.surface) {
                initialFilters.surface = {
                    totalArea: initialFilters.sizeRange,
                    coveredArea: { min: null, max: null },
                    landArea: { min: null, max: null },
                    front: null,
                    back: null
                };
            }
            
            setFilters(prev => ({
                ...prev,
                ...initialFilters
            }));
        }
    }, [initialFilters]);

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
        setFilters({
            propertyType: ['all'],
            showOnlyWithNotes: false,
            showOnlyFavorites: false,
            priceRange: {
                min: null,
                max: null,
                currency: 'USD'
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
            orientation: [],
            condition: [],
            operationType: [],
            situation: []
        });
    };

    // Verificar si hay filtros activos
    const hasActiveFilters = () => {
        return (
            (Array.isArray(filters.propertyType) ? 
                !filters.propertyType.includes('all') && filters.propertyType.length > 0 : 
                filters.propertyType !== 'all') ||
            filters.showOnlyWithNotes ||
            filters.showOnlyFavorites ||
            filters.priceRange.min !== null ||
            filters.priceRange.max !== null ||
            filters.locations.length > 0 ||
            filters.features.length > 0 ||
            filters.surface.totalArea.min !== null ||
            filters.surface.totalArea.max !== null ||
            filters.rooms !== null ||
            filters.bathrooms !== null ||
            filters.antiquity !== null ||
            (filters.disposition && filters.disposition.length > 0) ||
            (filters.orientation && filters.orientation.length > 0) ||
            (filters.condition && filters.condition.length > 0) ||
            (filters.operationType && filters.operationType.length > 0) ||
            (filters.situation && filters.situation.length > 0)
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
                        <span className="text-gray-400 dark:text-gray-400 mr-2">
                            {Array.isArray(filters.propertyType) ? 
                                (filters.propertyType.includes('all') ? 'Todas' : 
                                 filters.propertyType.length > 1 ? 'Múltiples' : 
                                 propertyTypeLabels[filters.propertyType[0]] || filters.propertyType[0]) : 
                                (filters.propertyType === 'all' ? 'Todas' : 
                                 propertyTypeLabels[filters.propertyType] || filters.propertyType)}
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

                {/* Superficie */}
                <button 
                    onClick={() => setShowSizeRangeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Superficie</span>
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

                {/* Disposición */}
                <button 
                    onClick={() => setShowDispositionDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Disposición</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.disposition && filters.disposition.length ? `${filters.disposition.length} seleccionados` : 'Todas'}</span>
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
                    <span className="text-gray-200 dark:text-gray-200">Orientación</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.orientation && filters.orientation.length ? `${filters.orientation.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Condición */}
                <button 
                    onClick={() => setShowConditionDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Condición</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.condition && filters.condition.length ? `${filters.condition.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Tipo de operación */}
                <button 
                    onClick={() => setShowOperationTypeDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Tipo de Operación</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.operationType && filters.operationType.length ? `${filters.operationType.length} seleccionados` : 'Todas'}</span>
                        <svg className="w-5 h-5 text-gray-400 dark:text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </button>

                {/* Situación */}
                <button 
                    onClick={() => setShowSituationDrawer(true)}
                    className="w-full p-4 bg-gray-800 dark:bg-gray-800 rounded-xl flex justify-between items-center"
                >
                    <span className="text-gray-200 dark:text-gray-200">Situación</span>
                    <div className="flex items-center">
                        <span className="text-gray-400 dark:text-gray-400 mr-2">{filters.situation && filters.situation.length ? `${filters.situation.length} seleccionados` : 'Todas'}</span>
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
                customContent={
                    <PropertyTypeFilter
                        initialValue={filters.propertyType}
                        onChange={(value) => {
                            setFilters(prev => ({ ...prev, propertyType: value }));
                        }}
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
                    <BathroomsFilter
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
                        initialValues={filters.orientation}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            orientation: values
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showConditionDrawer}
                onClose={() => setShowConditionDrawer(false)}
                title="CONDICIÓN"
                customContent={
                    <ConditionFilter
                        initialValues={filters.condition}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            condition: values
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showOperationTypeDrawer}
                onClose={() => setShowOperationTypeDrawer(false)}
                title="TIPO DE OPERACIÓN"
                customContent={
                    <OperationTypeFilter
                        initialValues={filters.operationType}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            operationType: values
                        }))}
                    />
                }
            />

            <FilterDrawer 
                isOpen={showSituationDrawer}
                onClose={() => setShowSituationDrawer(false)}
                title="SITUACIÓN"
                customContent={
                    <SituationFilter
                        initialValues={filters.situation}
                        onChange={(values) => setFilters(prev => ({
                            ...prev,
                            situation: values
                        }))}
                    />
                }
            />
        </div>
    );
}

export default Filters;
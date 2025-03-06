import React, { useState, useEffect } from 'react';
import { getPropertyTypes } from '../../services/api';

const FilterChips = ({ filters, onRemove }) => {
    const [propertyTypeLabels, setPropertyTypeLabels] = useState({
        house: 'Casa',
        apartment: 'Departamento',
        ph: 'PH'
    });

    // Cargar los tipos de propiedad desde el backend
    useEffect(() => {
        const fetchPropertyTypes = async () => {
            try {
                const response = await getPropertyTypes();
                if (response.data && response.data.length > 0) {
                    const typeMap = {};
                    response.data.forEach(type => {
                        // Guardar tanto por código como por ID
                        typeMap[type.code] = type.name;
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

    const getChips = () => {
        const chips = [];

        // Tipo de propiedad
        if (filters.propertyType) {
            if (Array.isArray(filters.propertyType)) {
                // Si es un array y no contiene 'all' y tiene elementos, mostrar chip
                if (!filters.propertyType.includes('all') && filters.propertyType.length > 0) {
                    // Si hay múltiples tipos seleccionados
                    if (filters.propertyType.length > 1) {
                        chips.push({
                            id: 'propertyType',
                            label: `${filters.propertyType.length} tipos de propiedad`,
                            onRemove: () => onRemove('propertyType', ['all'])
                        });
                    } else {
                        // Si solo hay un tipo seleccionado
                        const typeId = filters.propertyType[0];
                        chips.push({
                            id: 'propertyType',
                            label: propertyTypeLabels[typeId] || typeId,
                            onRemove: () => onRemove('propertyType', ['all'])
                        });
                    }
                }
            } else if (filters.propertyType !== 'all') {
                // Compatibilidad con versión anterior (string)
                chips.push({
                    id: 'propertyType',
                    label: propertyTypeLabels[filters.propertyType] || filters.propertyType,
                    onRemove: () => onRemove('propertyType', 'all')
                });
            }
        }

        // Rango de precio
        if (filters.priceRange.min || filters.priceRange.max) {
            const formatPrice = (price) => {
                if (price === null) return '';
                return filters.priceRange.currency === 'ARS' 
                    ? `$${new Intl.NumberFormat('es-AR').format(price)}`
                    : `USD ${new Intl.NumberFormat('en-US').format(price)}`;
            };

            let priceLabel = '';
            if (filters.priceRange.min && filters.priceRange.max) {
                priceLabel = `${formatPrice(filters.priceRange.min)} - ${formatPrice(filters.priceRange.max)}`;
            } else if (filters.priceRange.min) {
                priceLabel = `Desde ${formatPrice(filters.priceRange.min)}`;
            } else if (filters.priceRange.max) {
                priceLabel = `Hasta ${formatPrice(filters.priceRange.max)}`;
            }

            if (priceLabel) {
                chips.push({
                    id: 'priceRange',
                    label: priceLabel,
                    onRemove: () => onRemove('priceRange', { min: null, max: null, currency: filters.priceRange.currency })
                });
            }
        }

        // Ubicaciones
        if (filters.locations.length > 0) {
            chips.push({
                id: 'locations',
                label: `${filters.locations.length} ubicaciones`,
                onRemove: () => onRemove('locations', [])
            });
        }

        // Características
        if (filters.features.length > 0) {
            chips.push({
                id: 'features',
                label: `${filters.features.length} características`,
                onRemove: () => onRemove('features', [])
            });
        }

        // Tamaño
        if (filters.sizeRange.min || filters.sizeRange.max) {
            const formatSize = (size) => {
                if (size === null) return '';
                return `${size} m²`;
            };

            let sizeLabel = '';
            if (filters.sizeRange.min && filters.sizeRange.max) {
                sizeLabel = `${formatSize(filters.sizeRange.min)} - ${formatSize(filters.sizeRange.max)}`;
            } else if (filters.sizeRange.min) {
                sizeLabel = `Desde ${formatSize(filters.sizeRange.min)}`;
            } else if (filters.sizeRange.max) {
                sizeLabel = `Hasta ${formatSize(filters.sizeRange.max)}`;
            }

            if (sizeLabel) {
                chips.push({
                    id: 'sizeRange',
                    label: sizeLabel,
                    onRemove: () => onRemove('sizeRange', { min: null, max: null })
                });
            }
        }

        // Ambientes
        if (filters.rooms) {
            chips.push({
                id: 'rooms',
                label: `${filters.rooms}+ ambientes`,
                onRemove: () => onRemove('rooms', null)
            });
        }

        // Baños
        if (filters.bathrooms) {
            chips.push({
                id: 'bathrooms',
                label: `${filters.bathrooms}+ baños`,
                onRemove: () => onRemove('bathrooms', null)
            });
        }

        // Antigüedad
        if (filters.antiquity !== null) {
            const antiquityLabels = {
                0: 'A estrenar',
                5: 'Hasta 5 años',
                10: 'Hasta 10 años',
                20: 'Hasta 20 años',
                30: 'Hasta 30 años',
                100: 'Más de 30 años'
            };
            chips.push({
                id: 'antiquity',
                label: antiquityLabels[filters.antiquity],
                onRemove: () => onRemove('antiquity', null)
            });
        }

        // Solo con notas
        if (filters.showOnlyWithNotes) {
            chips.push({
                id: 'showOnlyWithNotes',
                label: 'Solo con notas',
                onRemove: () => onRemove('showOnlyWithNotes', false)
            });
        }

        // Solo favoritos
        if (filters.showOnlyFavorites) {
            chips.push({
                id: 'showOnlyFavorites',
                label: 'Solo favoritos',
                onRemove: () => onRemove('showOnlyFavorites', false)
            });
        }

        return chips;
    };

    const chips = getChips();

    if (chips.length === 0) {
        return null;
    }

    return (
        <div className="flex flex-wrap gap-2 p-4">
            {chips.map(chip => (
                <div 
                    key={chip.id}
                    className="flex items-center gap-1 px-3 py-1.5 bg-blue-500/20 text-blue-400 rounded-full text-sm"
                >
                    <span>{chip.label}</span>
                    <button 
                        onClick={chip.onRemove}
                        className="w-4 h-4 flex items-center justify-center rounded-full hover:bg-blue-500/30"
                    >
                        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            ))}
        </div>
    );
};

export default FilterChips; 
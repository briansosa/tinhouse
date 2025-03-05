import { useState, useEffect } from 'react';
import { getPropertyTypes } from '../../services/api';

const PropertyTypeFilter = ({ onChange, initialValue = 'all' }) => {
    const [selectedType, setSelectedType] = useState(initialValue);
    const [propertyTypes, setPropertyTypes] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    
    // Cargar tipos de propiedad desde el backend
    useEffect(() => {
        const fetchPropertyTypes = async () => {
            try {
                setIsLoading(true);
                const response = await getPropertyTypes();                
                if (response.data && response.data.length > 0) {
                    // Convertir los tipos de propiedad al formato esperado
                    const types = [
                        { id: 'all', code: 'all', label: 'Todas' },
                        ...response.data.map(type => ({
                            id: type.id.toString(), // Usar el ID como identificador
                            code: type.code,        // Guardar el código para referencia
                            label: type.name
                        }))
                    ];
                    setPropertyTypes(types);
                } else {
                    console.log("PropertyTypeFilter: No se recibieron tipos de propiedad del servidor");
                    setPropertyTypes(defaultPropertyTypes);
                }
                setError(null);
            } catch (err) {
                console.error('PropertyTypeFilter: Error al cargar tipos de propiedad:', err);
                setError('No se pudieron cargar los tipos de propiedad');
                setPropertyTypes(defaultPropertyTypes);
            } finally {
                setIsLoading(false);
            }
        };
        
        fetchPropertyTypes();
    }, []);
    
    // Tipos de propiedad por defecto en caso de error
    const defaultPropertyTypes = [
        { id: 'all', code: 'all', label: 'Todas' },
        { id: '1', code: 'house', label: 'Casa' },
        { id: '2', code: 'apartment', label: 'Departamento' },
        { id: '3', code: 'ph', label: 'PH' },
        { id: '4', code: 'local', label: 'Local' },
        { id: '5', code: 'office', label: 'Oficina' },
        { id: '6', code: 'land', label: 'Terreno' }
    ];

    const handleTypeSelect = (typeId) => {
        setSelectedType(typeId);
        // Pasar directamente el ID al componente padre
        onChange(typeId);
    };

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4 text-center text-red-500">
                {error}
            </div>
        );
    }

    return (
        <div className="space-y-4 px-4">
            <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase">
                TIPO DE PROPIEDAD
            </h3>
            <div className="grid grid-cols-2 gap-2">
                {propertyTypes.map(type => (
                    <button
                        key={type.id}
                        onClick={() => handleTypeSelect(type.id)}
                        className={`p-3 rounded-xl text-sm font-medium transition-all flex items-center justify-center gap-2 ${
                            selectedType === type.id
                                ? 'bg-blue-500 text-white'
                                : 'bg-gray-800 text-gray-300'
                        }`}
                    >
                        <span className="truncate">{type.label}</span>
                    </button>
                ))}
            </div>
        </div>
    );
};

export default PropertyTypeFilter;
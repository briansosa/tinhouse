import { useState, useEffect } from 'react';
import { getListValues } from '../../services/api';

const ConditionFilter = ({ onChange, initialValues = [] }) => {
    const [conditionOptions, setConditionOptions] = useState([]);
    const [selectedValues, setSelectedValues] = useState(initialValues);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchConditionOptions = async () => {
            try {
                setLoading(true);
                const data = await getListValues('condicion');
                setConditionOptions(data);
                setError(null);
            } catch (err) {
                console.error('Error al cargar opciones de condición:', err);
                setError('No se pudieron cargar las opciones de condición');
            } finally {
                setLoading(false);
            }
        };

        fetchConditionOptions();
    }, []);

    const handleToggleOption = (value) => {
        let newValues;
        if (selectedValues.includes(value)) {
            newValues = selectedValues.filter(v => v !== value);
        } else {
            newValues = [...selectedValues, value];
        }
        setSelectedValues(newValues);
        onChange(newValues);
    };

    const handleReset = () => {
        setSelectedValues([]);
        onChange([]);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center p-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4 text-center">
                <p className="text-red-500">{error}</p>
            </div>
        );
    }

    return (
        <div className="p-4">
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-white font-semibold text-lg">Condición</h3>
                <button
                    onClick={handleReset}
                    className="text-xs text-gray-400 hover:text-gray-300"
                >
                    Restablecer
                </button>
            </div>

            <div className="space-y-3">
                {conditionOptions.map((option) => (
                    <div key={option.id} className="flex items-center">
                        <button
                            onClick={() => handleToggleOption(option.value)}
                            className={`flex-1 p-3 rounded-lg flex justify-between items-center ${
                                selectedValues.includes(option.value)
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-gray-700 text-gray-300'
                            }`}
                        >
                            <span>{option.display_name}</span>
                            {selectedValues.includes(option.value) && (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                            )}
                        </button>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ConditionFilter; 
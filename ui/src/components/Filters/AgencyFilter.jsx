import { useState, useEffect } from 'react';
import { getAgencies } from '../../services/api';

const AgencyFilter = ({ onChange, initialValues = [] }) => {
    const [agencyOptions, setAgencyOptions] = useState([]);
    const [selectedValues, setSelectedValues] = useState(initialValues);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchAgencyOptions = async () => {
            try {
                setLoading(true);
                const data = await getAgencies();
                setAgencyOptions(data || []);
                setError(null);
            } catch (err) {
                console.error('Error al cargar opciones de inmobiliarias:', err);
                setError('No se pudieron cargar las opciones de inmobiliarias');
            } finally {
                setLoading(false);
            }
        };

        fetchAgencyOptions();
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

    return (
        <div className="p-4 pb-16">
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-white font-semibold text-lg">Inmobiliaria</h3>
                <button
                    onClick={handleReset}
                    className="text-xs text-gray-400 hover:text-gray-300"
                >
                    Restablecer
                </button>
            </div>

            <div className="space-y-3 overflow-y-auto">
                {loading ? (
                    <div className="flex justify-center items-center py-8">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                    </div>
                ) : error ? (
                    <div className="p-4 text-center">
                        <p className="text-red-500">{error}</p>
                    </div>
                ) : agencyOptions.length === 0 ? (
                    <div className="text-center text-gray-400 py-8">
                        No hay inmobiliarias disponibles
                    </div>
                ) : (
                    agencyOptions.map((agency) => (
                        <div key={agency.id} className="flex items-center">
                            <button
                                onClick={() => handleToggleOption(agency.id)}
                                className={`flex-1 p-3 rounded-lg flex justify-between items-center ${
                                    selectedValues.includes(agency.id)
                                        ? 'bg-blue-600 text-white'
                                        : 'bg-gray-700 text-gray-300'
                                }`}
                            >
                                <span>{agency.name}</span>
                                {selectedValues.includes(agency.id) && (
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                )}
                            </button>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

export default AgencyFilter; 
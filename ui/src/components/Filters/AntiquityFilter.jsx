import { useState } from 'react';

const AntiquityFilter = ({ onChange, initialValue = null }) => {
    const [selectedValue, setSelectedValue] = useState(initialValue);
    
    const options = [
        { value: 0, label: "A estrenar" },
        { value: 5, label: "Hasta 5 años" },
        { value: 10, label: "Hasta 10 años" },
        { value: 20, label: "Hasta 20 años" },
        { value: 30, label: "Hasta 30 años" },
        { value: 100, label: "Más de 30 años" }
    ];

    const handleSelect = (value) => {
        const newValue = selectedValue === value ? null : value;
        setSelectedValue(newValue);
        onChange(newValue);
    };

    return (
        <div className="space-y-4 px-4 pb-4">
            <div className="flex justify-between items-center">
                <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase">
                    ANTIGÜEDAD
                </h3>
            </div>
            
            <div className="grid grid-cols-2 gap-2">
                {options.map(option => (
                    <button
                        key={option.value}
                        onClick={() => handleSelect(option.value)}
                        className={`p-3 rounded-xl text-sm font-medium transition-all ${
                            selectedValue === option.value
                                ? 'bg-blue-500 text-white'
                                : 'bg-gray-800 text-gray-300'
                        }`}
                    >
                        {option.label}
                    </button>
                ))}
            </div>
        </div>
    );
};

export default AntiquityFilter; 
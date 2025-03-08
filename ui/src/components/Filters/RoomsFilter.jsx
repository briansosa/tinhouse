import { useState } from 'react';

const RoomsFilter = ({ onChange, initialValue = null, title = "Ambientes" }) => {
    const [selectedValue, setSelectedValue] = useState(initialValue);
    
    const options = [
        { value: 1, label: "1+" },
        { value: 2, label: "2+" },
        { value: 3, label: "3+" },
        { value: 4, label: "4+" },
        { value: 5, label: "5+" }
    ];

    const handleSelect = (value) => {
        const newValue = selectedValue === value ? null : value;
        setSelectedValue(newValue);
        onChange(newValue);
    };

    const handleReset = () => {
        setSelectedValue(null);
        onChange(null);
    };

    return (
        <div className="p-4 pb-16">
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-white font-semibold text-lg">{title}</h3>
                <button
                    onClick={handleReset}
                    className="text-xs text-gray-400 hover:text-gray-300"
                >
                    Restablecer
                </button>
            </div>
            
            <div className="space-y-4 min-h-[200px]">
                <div className="grid grid-cols-2 gap-3">
                    {options.map(option => (
                        <button
                            key={option.value}
                            onClick={() => handleSelect(option.value)}
                            className={`p-3 rounded-lg text-base font-medium transition-all ${
                                selectedValue === option.value
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-gray-700 text-gray-300'
                            }`}
                        >
                            {option.label}
                        </button>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default RoomsFilter; 
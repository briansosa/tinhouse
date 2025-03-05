import { useState } from 'react';

const BathroomsFilter = ({ onChange, initialValue = null, title = "Baños" }) => {
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

    return (
        <div className="space-y-4 px-4">
            <h3 className="text-sm font-bold text-gray-400 tracking-wider uppercase">
                {title}
            </h3>
            <div className="flex gap-2 overflow-x-auto pb-2">
                {options.map(option => (
                    <button
                        key={option.value}
                        onClick={() => handleSelect(option.value)}
                        className={`px-4 py-2 rounded-lg text-sm font-medium transition-all flex-shrink-0 ${
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

export default BathroomsFilter; 
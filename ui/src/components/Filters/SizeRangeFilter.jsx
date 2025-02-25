import { useState } from 'react';

const SizeRangeFilter = ({ onChange, initialRange = { min: null, max: null } }) => {
    const [range, setRange] = useState(initialRange);
    const maxSize = 500; // metros cuadrados máximos

    const handleSliderChange = (type) => (e) => {
        const value = parseInt(e.target.value);
        let newRange = { ...range };

        if (type === 'min') {
            newRange.min = Math.min(value, (newRange.max || maxSize) - 1);
        } else {
            newRange.max = Math.max(value, (newRange.min || 0) + 1);
        }

        setRange(newRange);
        onChange(newRange);
    };

    const formatSize = (size) => {
        if (size === null) return 'Sin límite';
        return `${size} m²`;
    };

    return (
        <div className="space-y-12 px-4">
            {/* Valores seleccionados */}
            <div className="flex items-center justify-between text-xl">
                <span className="text-white font-medium">{formatSize(range.min || 0)}</span>
                <span className="text-white font-medium">{formatSize(range.max || maxSize)}</span>
            </div>

            {/* Contenedor de los sliders */}
            <div className="relative h-2 mx-6">
                {/* Barra de fondo */}
                <div className="absolute w-full h-full bg-gray-700 rounded-full" />
                
                {/* Barra de rango seleccionado */}
                <div 
                    className="absolute h-full bg-blue-500 rounded-full"
                    style={{
                        left: `${((range.min || 0) / maxSize) * 100}%`,
                        right: `${100 - ((range.max || maxSize) / maxSize) * 100}%`
                    }}
                />

                {/* Slider máximo */}
                <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                    <input
                        type="range"
                        min="0"
                        max={maxSize}
                        value={range.max || maxSize}
                        onChange={handleSliderChange('max')}
                        className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                    />
                </div>

                {/* Slider mínimo */}
                <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                    <input
                        type="range"
                        min="0"
                        max={maxSize}
                        value={range.min || 0}
                        onChange={handleSliderChange('min')}
                        className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                    />
                </div>
            </div>
        </div>
    );
};

export default SizeRangeFilter;
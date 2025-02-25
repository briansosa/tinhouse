import { useState } from 'react';

const PriceRangeFilter = ({ onChange, initialRange = { min: null, max: null } }) => {
    const [range, setRange] = useState(initialRange);
    const [currency, setCurrency] = useState('ARS'); // 'ARS' o 'USD'
    const maxPrice = currency === 'ARS' ? 500000000 : 500000; // Ajustamos máximo según moneda

    const handleSliderChange = (type) => (e) => {
        const value = parseInt(e.target.value);
        let newRange = { ...range };

        if (type === 'min') {
            newRange.min = Math.min(value, (newRange.max || maxPrice) - 1);
        } else {
            newRange.max = Math.max(value, (newRange.min || 0) + 1);
        }

        setRange(newRange);
        onChange({ ...newRange, currency });
    };

    const handleCurrencyToggle = () => {
        const newCurrency = currency === 'ARS' ? 'USD' : 'ARS';
        setCurrency(newCurrency);
        // Reseteamos el rango al cambiar de moneda
        const newRange = { min: null, max: null };
        setRange(newRange);
        onChange({ ...newRange, currency: newCurrency });
    };

    const formatPrice = (price) => {
        if (price === null) return 'Sin límite';
        return currency === 'ARS' 
            ? `$${new Intl.NumberFormat('es-AR').format(price)}`
            : `USD ${new Intl.NumberFormat('en-US').format(price)}`;
    };

    return (
        <div className="space-y-12 px-4">
            {/* Toggle de moneda */}
            <div className="flex items-center justify-between p-4 rounded-xl bg-gray-800">
                <span className="text-md text-gray-300">$/USD</span>
                <button 
                    onClick={handleCurrencyToggle}
                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-300 ${
                        currency === 'USD' ? 'bg-blue-500' : 'bg-gray-600'
                    }`}
                >
                    <span 
                        className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-lg transition-transform duration-300 ${
                            currency === 'USD' ? 'translate-x-6' : 'translate-x-1'
                        }`}
                    />
                </button>
            </div>

            {/* Valores seleccionados */}
            <div className="flex items-center justify-between text-xl">
                <span className="text-white font-medium">{formatPrice(range.min || 0)}</span>
                <span className="text-white font-medium">{formatPrice(range.max || maxPrice)}</span>
            </div>

            {/* Contenedor de los sliders */}
            <div className="relative h-2 mx-6">
                {/* Barra de fondo */}
                <div className="absolute w-full h-full bg-gray-700 rounded-full" />
                
                {/* Barra de rango seleccionado */}
                <div 
                    className="absolute h-full bg-blue-500 rounded-full"
                    style={{
                        left: `${((range.min || 0) / maxPrice) * 100}%`,
                        right: `${100 - ((range.max || maxPrice) / maxPrice) * 100}%`
                    }}
                />

                {/* Slider máximo */}
                <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                    <input
                        type="range"
                        min="0"
                        max={maxPrice}
                        value={range.max || maxPrice}
                        onChange={handleSliderChange('max')}
                        className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                    />
                </div>

                {/* Slider mínimo */}
                <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                    <input
                        type="range"
                        min="0"
                        max={maxPrice}
                        value={range.min || 0}
                        onChange={handleSliderChange('min')}
                        className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                    />
                </div>
            </div>
        </div>
    );
};

export default PriceRangeFilter; 
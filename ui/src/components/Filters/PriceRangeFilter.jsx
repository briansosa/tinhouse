import { useState } from 'react';

const PriceRangeFilter = ({ onChange, initialRange = { min: null, max: null, currency: 'ARS' } }) => {
    const [range, setRange] = useState({
        min: initialRange.min,
        max: initialRange.max
    });
    const [currency, setCurrency] = useState(initialRange.currency || 'ARS'); // 'ARS' o 'USD'
    
    // Valores máximos según la moneda
    const maxPrices = {
        'ARS': 500000000, // 500 millones de pesos
        'USD': 2000000    // 2 millones de dólares
    };

    const handleSliderChange = (type) => (e) => {
        const value = parseInt(e.target.value);
        let newRange = { ...range };

        if (type === 'min') {
            newRange.min = Math.min(value, (newRange.max || maxPrices[currency]) - 1);
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

    const handlePresetClick = (value) => {
        const newRange = { ...range, max: value };
        setRange(newRange);
        onChange({ ...newRange, currency });
    };

    const formatPrice = (price) => {
        if (price === null) return 'Sin límite';
        return currency === 'ARS' 
            ? `$${new Intl.NumberFormat('es-AR').format(price)}`
            : `USD ${new Intl.NumberFormat('en-US').format(price)}`;
    };

    return (
        <div className="p-4 space-y-6">
            {/* Selector de moneda con explicación */}
            <div className="space-y-2">
                <div className="flex justify-center space-x-2">
                    <button 
                        onClick={() => handleCurrencyToggle()}
                        className={`px-4 py-2 rounded-full text-sm font-medium ${
                            currency === 'ARS' 
                                ? 'bg-blue-500 text-white' 
                                : 'bg-gray-800 dark:bg-gray-800 text-gray-300 dark:text-gray-300'
                        }`}
                    >
                        {currency === 'ARS' ? 'ARS' : 'USD'}
                    </button>
                </div>
                <p className="text-xs text-center text-gray-400">
                    {currency === 'ARS' 
                        ? 'Mostrando solo propiedades en pesos argentinos' 
                        : 'Mostrando solo propiedades en dólares'}
                </p>
            </div>

            {/* Inputs de rango */}
            <div className="flex items-center space-x-4">
                <div className="flex-1">
                    <label className="block text-sm font-medium text-gray-300 dark:text-gray-300 mb-1">Mínimo</label>
                    <div className="relative">
                        <span className="absolute inset-y-0 left-3 flex items-center text-gray-400 dark:text-gray-400">
                            {currency === 'ARS' ? '$' : 'US$'}
                        </span>
                        <input
                            type="number"
                            value={range.min || ''}
                            onChange={(e) => handleSliderChange('min')(e)}
                            placeholder="Mínimo"
                            className="w-full pl-10 pr-3 py-2 border border-gray-600 dark:border-gray-600 rounded-lg bg-gray-800 dark:bg-gray-800 text-white dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
                <div className="flex-1">
                    <label className="block text-sm font-medium text-gray-300 dark:text-gray-300 mb-1">Máximo</label>
                    <div className="relative">
                        <span className="absolute inset-y-0 left-3 flex items-center text-gray-400 dark:text-gray-400">
                            {currency === 'ARS' ? '$' : 'US$'}
                        </span>
                        <input
                            type="number"
                            value={range.max || ''}
                            onChange={(e) => handleSliderChange('max')(e)}
                            placeholder="Máximo"
                            className="w-full pl-10 pr-3 py-2 border border-gray-600 dark:border-gray-600 rounded-lg bg-gray-800 dark:bg-gray-800 text-white dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
            </div>

            {/* Valores seleccionados */}
            <div className="flex items-center justify-between text-xl">
                <span className="text-white font-medium">{formatPrice(range.min || 0)}</span>
                <span className="text-white font-medium">{formatPrice(range.max || maxPrices[currency])}</span>
            </div>

            {/* Contenedor de los sliders */}
            <div className="relative h-2 mx-6">
                {/* Barra de fondo */}
                <div className="absolute w-full h-full bg-gray-700 rounded-full" />
                
                {/* Barra de rango seleccionado */}
                <div 
                    className="absolute h-full bg-blue-500 rounded-full"
                    style={{
                        left: `${((range.min || 0) / maxPrices[currency]) * 100}%`,
                        right: `${100 - ((range.max || maxPrices[currency]) / maxPrices[currency]) * 100}%`
                    }}
                />

                {/* Slider máximo */}
                <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                    <input
                        type="range"
                        min="0"
                        max={maxPrices[currency]}
                        value={range.max || maxPrices[currency]}
                        onChange={handleSliderChange('max')}
                        className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                    />
                </div>

                {/* Slider mínimo */}
                <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                    <input
                        type="range"
                        min="0"
                        max={maxPrices[currency]}
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
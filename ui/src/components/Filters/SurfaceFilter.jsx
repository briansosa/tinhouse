import { useState } from 'react';

const SurfaceFilter = ({ onChange, initialValues = {
    totalArea: { min: null, max: null },
    coveredArea: { min: null, max: null },
    landArea: { min: null, max: null },
    front: null,
    back: null
} }) => {
    const [values, setValues] = useState(initialValues);
    const maxSize = 500; // metros cuadrados máximos
    const maxLength = 100; // metros máximos para frente y fondo

    const handleSliderChange = (type, field) => (e) => {
        const value = parseInt(e.target.value);
        let newValues = { ...values };

        if (type === 'min') {
            newValues[field].min = Math.min(value, (newValues[field].max || maxSize) - 1);
        } else {
            newValues[field].max = Math.max(value, (newValues[field].min || 0) + 1);
        }

        setValues(newValues);
        onChange(newValues);
    };

    const handleInputChange = (field) => (e) => {
        const value = e.target.value === '' ? null : parseInt(e.target.value);
        let newValues = { ...values };
        
        newValues[field] = value;
        
        setValues(newValues);
        onChange(newValues);
    };

    const formatSize = (size) => {
        if (size === null) return 'Sin límite';
        return `${size} m²`;
    };

    const formatLength = (length) => {
        if (length === null) return 'Sin valor';
        return `${length} m`;
    };

    // Función para restablecer todos los valores
    const handleReset = () => {
        const defaultValues = {
            totalArea: { min: null, max: null },
            coveredArea: { min: null, max: null },
            landArea: { min: null, max: null },
            front: null,
            back: null
        };
        
        setValues(defaultValues);
        onChange(defaultValues);
    };

    return (
        <div className="space-y-8 px-4">
            {/* Superficie Total */}
            <div className="space-y-4">
                <h3 className="text-white font-semibold text-lg">Superficie Total</h3>
                
                {/* Valores seleccionados */}
                <div className="flex items-center justify-between text-xl">
                    <span className="text-white font-medium">{formatSize(values.totalArea.min || 0)}</span>
                    <span className="text-white font-medium">{formatSize(values.totalArea.max || maxSize)}</span>
                </div>

                {/* Contenedor de los sliders */}
                <div className="relative h-2 mx-6">
                    {/* Barra de fondo */}
                    <div className="absolute w-full h-full bg-gray-700 rounded-full" />
                    
                    {/* Barra de rango seleccionado */}
                    <div 
                        className="absolute h-full bg-blue-500 rounded-full"
                        style={{
                            left: `${((values.totalArea.min || 0) / maxSize) * 100}%`,
                            right: `${100 - ((values.totalArea.max || maxSize) / maxSize) * 100}%`
                        }}
                    />

                    {/* Slider máximo */}
                    <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                        <input
                            type="range"
                            min="0"
                            max={maxSize}
                            value={values.totalArea.max || maxSize}
                            onChange={handleSliderChange('max', 'totalArea')}
                            className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                        />
                    </div>

                    {/* Slider mínimo */}
                    <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                        <input
                            type="range"
                            min="0"
                            max={maxSize}
                            value={values.totalArea.min || 0}
                            onChange={handleSliderChange('min', 'totalArea')}
                            className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                        />
                    </div>
                </div>
            </div>

            {/* Superficie Cubierta */}
            <div className="space-y-4">
                <h3 className="text-white font-semibold text-lg">Superficie Cubierta</h3>
                
                {/* Valores seleccionados */}
                <div className="flex items-center justify-between text-xl">
                    <span className="text-white font-medium">{formatSize(values.coveredArea.min || 0)}</span>
                    <span className="text-white font-medium">{formatSize(values.coveredArea.max || maxSize)}</span>
                </div>

                {/* Contenedor de los sliders */}
                <div className="relative h-2 mx-6">
                    {/* Barra de fondo */}
                    <div className="absolute w-full h-full bg-gray-700 rounded-full" />
                    
                    {/* Barra de rango seleccionado */}
                    <div 
                        className="absolute h-full bg-blue-500 rounded-full"
                        style={{
                            left: `${((values.coveredArea.min || 0) / maxSize) * 100}%`,
                            right: `${100 - ((values.coveredArea.max || maxSize) / maxSize) * 100}%`
                        }}
                    />

                    {/* Slider máximo */}
                    <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                        <input
                            type="range"
                            min="0"
                            max={maxSize}
                            value={values.coveredArea.max || maxSize}
                            onChange={handleSliderChange('max', 'coveredArea')}
                            className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                        />
                    </div>

                    {/* Slider mínimo */}
                    <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                        <input
                            type="range"
                            min="0"
                            max={maxSize}
                            value={values.coveredArea.min || 0}
                            onChange={handleSliderChange('min', 'coveredArea')}
                            className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                        />
                    </div>
                </div>
            </div>

            {/* Superficie Terreno */}
            <div className="space-y-4">
                <h3 className="text-white font-semibold text-lg">Superficie Terreno</h3>
                
                {/* Valores seleccionados */}
                <div className="flex items-center justify-between text-xl">
                    <span className="text-white font-medium">{formatSize(values.landArea.min || 0)}</span>
                    <span className="text-white font-medium">{formatSize(values.landArea.max || maxSize)}</span>
                </div>

                {/* Contenedor de los sliders */}
                <div className="relative h-2 mx-6">
                    {/* Barra de fondo */}
                    <div className="absolute w-full h-full bg-gray-700 rounded-full" />
                    
                    {/* Barra de rango seleccionado */}
                    <div 
                        className="absolute h-full bg-blue-500 rounded-full"
                        style={{
                            left: `${((values.landArea.min || 0) / maxSize) * 100}%`,
                            right: `${100 - ((values.landArea.max || maxSize) / maxSize) * 100}%`
                        }}
                    />

                    {/* Slider máximo */}
                    <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                        <input
                            type="range"
                            min="0"
                            max={maxSize}
                            value={values.landArea.max || maxSize}
                            onChange={handleSliderChange('max', 'landArea')}
                            className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                        />
                    </div>

                    {/* Slider mínimo */}
                    <div className="absolute inset-x-[-12px] inset-y-0 pointer-events-none">
                        <input
                            type="range"
                            min="0"
                            max={maxSize}
                            value={values.landArea.min || 0}
                            onChange={handleSliderChange('min', 'landArea')}
                            className="w-full h-full appearance-none bg-transparent [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:w-6 [&::-webkit-slider-thumb]:h-6 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:cursor-pointer [&::-webkit-slider-thumb]:shadow-lg [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:w-6 [&::-moz-range-thumb]:h-6 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white [&::-moz-range-thumb]:border-0 [&::-moz-range-thumb]:cursor-pointer [&::-moz-range-thumb]:shadow-lg"
                        />
                    </div>
                </div>
            </div>

            {/* Frente y Fondo */}
            <div className="space-y-4">
                <h3 className="text-white font-semibold text-lg">Frente y Fondo</h3>
                
                <div className="grid grid-cols-2 gap-4">
                    {/* Frente */}
                    <div className="space-y-2">
                        <label className="text-white text-sm">Frente (metros)</label>
                        <input
                            type="number"
                            min="0"
                            max={maxLength}
                            value={values.front || ''}
                            onChange={handleInputChange('front')}
                            placeholder="Metros"
                            className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                    
                    {/* Fondo */}
                    <div className="space-y-2">
                        <label className="text-white text-sm">Fondo (metros)</label>
                        <input
                            type="number"
                            min="0"
                            max={maxLength}
                            value={values.back || ''}
                            onChange={handleInputChange('back')}
                            placeholder="Metros"
                            className="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
            </div>

            {/* Botón de Restablecer */}
            <div className="pt-4">
                <button
                    onClick={handleReset}
                    className="w-full py-3 bg-gray-700 hover:bg-gray-600 text-white font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                    Restablecer
                </button>
            </div>
        </div>
    );
};

export default SurfaceFilter; 
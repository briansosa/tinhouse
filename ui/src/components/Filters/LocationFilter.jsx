import { useState } from 'react';

const LocationFilter = ({ onChange, initialLocations = [] }) => {
    const [selectedLocations, setSelectedLocations] = useState(initialLocations);
    const [selectedZone, setSelectedZone] = useState(null);
    
    const locations = [
        { 
            id: 'caba',
            name: 'Capital Federal',
            neighborhoods: [
                { id: 'palermo', name: 'Palermo' },
                { id: 'belgrano', name: 'Belgrano' },
                { id: 'recoleta', name: 'Recoleta' },
                { id: 'caballito', name: 'Caballito' },
                { id: 'almagro', name: 'Almagro' },
                { id: 'nunez', name: 'Núñez' }
            ]
        },
        {
            id: 'gba_norte',
            name: 'GBA Norte',
            neighborhoods: [
                { id: 'vicente_lopez', name: 'Vicente López' },
                { id: 'san_isidro', name: 'San Isidro' },
                { id: 'tigre', name: 'Tigre' },
                { id: 'san_fernando', name: 'San Fernando' }
            ]
        },
        {
            id: 'gba_sur',
            name: 'GBA Sur',
            neighborhoods: [
                { id: 'avellaneda', name: 'Avellaneda' },
                { id: 'lomas', name: 'Lomas de Zamora' },
                { id: 'quilmes', name: 'Quilmes' },
                { id: 'lanus', name: 'Lanús' }
            ]
        }
    ];

    const handleZoneSelect = (zoneId) => {
        if (selectedZone === zoneId) {
            setSelectedZone(null);
        } else {
            setSelectedZone(zoneId);
        }
    };

    const handleNeighborhoodToggle = (neighborhoodId) => {
        const newLocations = selectedLocations.includes(neighborhoodId)
            ? selectedLocations.filter(id => id !== neighborhoodId)
            : [...selectedLocations, neighborhoodId];
        
        setSelectedLocations(newLocations);
        onChange(newLocations);
    };

    return (
        <div className="space-y-12 px-4">
            {locations.map(zone => (
                <div key={zone.id} className="space-y-3">
                    <button
                        onClick={() => handleZoneSelect(zone.id)}
                        className={`w-full p-4 rounded-xl text-base font-medium transition-all flex items-center justify-between ${
                            selectedZone === zone.id
                                ? 'bg-blue-500 text-white'
                                : 'bg-gray-800 text-gray-300'
                        }`}
                    >
                        <span>{zone.name}</span>
                        <svg 
                            className={`w-5 h-5 transform transition-transform ${
                                selectedZone === zone.id ? 'rotate-180' : ''
                            }`} 
                            fill="none" 
                            viewBox="0 0 24 24" 
                            stroke="currentColor"
                        >
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                    </button>

                    {selectedZone === zone.id && (
                        <div className="pl-4 space-y-2">
                            {zone.neighborhoods.map(neighborhood => (
                                <button
                                    key={neighborhood.id}
                                    onClick={() => handleNeighborhoodToggle(neighborhood.id)}
                                    className={`w-full p-3 rounded-lg text-sm font-medium transition-all flex items-center justify-between ${
                                        selectedLocations.includes(neighborhood.id)
                                            ? 'bg-blue-500/50 text-white'
                                            : 'bg-gray-800/50 text-gray-400'
                                    }`}
                                >
                                    {neighborhood.name}
                                    {selectedLocations.includes(neighborhood.id) && (
                                        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                        </svg>
                                    )}
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            ))}
        </div>
    );
};

export default LocationFilter; 
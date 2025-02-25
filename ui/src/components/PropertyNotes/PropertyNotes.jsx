import { useState, useEffect } from 'react';

export default function PropertyNotes({ property, onClose, onImageClick }) {
    const [notes, setNotes] = useState([]);
    const [newNote, setNewNote] = useState('');

    // Cargar notas del localStorage
    useEffect(() => {
        const savedNotes = localStorage.getItem(`property-notes-${property.id}`);
        if (savedNotes) {
            setNotes(JSON.parse(savedNotes));
        }
    }, [property.id]);

    const saveNote = () => {
        if (!newNote.trim()) return;
        
        const newNoteObj = {
            id: Date.now(),
            text: newNote,
            date: new Date().toISOString(),
        };
        
        const updatedNotes = [...notes, newNoteObj];
        setNotes(updatedNotes);
        localStorage.setItem(`property-notes-${property.id}`, JSON.stringify(updatedNotes));
        setNewNote('');
    };

    return (
        <div className="h-full flex flex-col bg-white dark:bg-gray-950">
            {/* Header */}
            <div className="flex items-center p-4 border-b dark:border-gray-800">
                <button onClick={onClose} className="mr-4">
                    <svg className="w-6 h-6 dark:text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                </button>
                <div 
                    className="w-16 h-16 rounded-full overflow-hidden mr-4 flex-shrink-0 cursor-pointer hover:opacity-80 transition-opacity"
                    onClick={onImageClick}
                >
                    <img 
                        src={property.image_url} 
                        alt={property.title}
                        className="w-full h-full object-cover"
                    />
                </div>
                <div className="flex-1">
                    <h2 className="text-md font-semibold dark:text-white">{property.title}</h2>
                    <p className="text-sm text-gray-500">{property.location}</p>
                </div>
            </div>

            {/* Chat area */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                {notes.length === 0 ? (
                    <div className="text-center text-gray-500 py-8">
                        No hay notas aún
                    </div>
                ) : (
                    notes.map(note => (
                        <div 
                            key={note.id} 
                            className="bg-blue-500 text-white p-3 rounded-lg max-w-[80%] ml-auto"
                        >
                            <p>{note.text}</p>
                            <p className="text-xs opacity-75 mt-1">
                                {new Date(note.date).toLocaleTimeString([], { 
                                    hour: '2-digit', 
                                    minute: '2-digit' 
                                })}
                            </p>
                        </div>
                    ))
                )}
            </div>

            {/* Input area */}
            <div className="p-4 border-t dark:border-gray-800">
                <div className="flex gap-2">
                    <input
                        type="text"
                        value={newNote}
                        onChange={(e) => setNewNote(e.target.value)}
                        placeholder="Escribe una nota..."
                        className="flex-1 px-4 py-2 bg-gray-100 dark:bg-gray-800 rounded-full"
                        onKeyPress={(e) => e.key === 'Enter' && saveNote()}
                    />
                    <button 
                        onClick={saveNote}
                        className="px-4 py-2 bg-blue-500 text-white rounded-full disabled:opacity-50"
                        disabled={!newNote.trim()}
                    >
                        Enviar
                    </button>
                </div>
            </div>
        </div>
    );
} 
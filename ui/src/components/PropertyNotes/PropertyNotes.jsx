import { useState, useEffect } from 'react';
import { getPropertyNotes, addPropertyNote, deletePropertyNote, togglePropertyFavorite, dislikeProperty } from '../../services/api';

export default function PropertyNotes({ property, onClose, onImageClick, onDislike }) {
    const [notes, setNotes] = useState([]);
    const [newNote, setNewNote] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isFavorite, setIsFavorite] = useState(property.is_favorite);
    const [isHeartHovered, setIsHeartHovered] = useState(false);
    const [showMenu, setShowMenu] = useState(false);

    // Cargar notas desde la API
    useEffect(() => {
        const fetchNotes = async () => {
            try {
                setIsLoading(true);
                const response = await getPropertyNotes(property.id);
                setNotes(response.data.notes || []);
                setError(null);
            } catch (err) {
                console.error('Error al cargar las notas:', err);
                setError('No se pudieron cargar las notas');
            } finally {
                setIsLoading(false);
            }
        };

        fetchNotes();
    }, [property.id]);

    const saveNote = async () => {
        if (!newNote.trim()) return;
        
        try {
            const response = await addPropertyNote(property.id, newNote);
            const newNoteObj = response.data.note;
            
            setNotes(prev => [...prev, newNoteObj]);
            setNewNote('');
            setError(null);
        } catch (err) {
            console.error('Error al guardar la nota:', err);
            setError('No se pudo guardar la nota');
        }
    };

    const handleDeleteNote = async (noteId) => {
        try {
            await deletePropertyNote(noteId);
            setNotes(prev => prev.filter(note => note.id !== noteId));
            setError(null);
        } catch (err) {
            console.error('Error al eliminar la nota:', err);
            setError('No se pudo eliminar la nota');
        }
    };

    const handleToggleFavorite = async () => {
        try {
            await togglePropertyFavorite(property.id, !isFavorite);
            setIsFavorite(!isFavorite);
        } catch (err) {
            console.error('Error al cambiar estado de favorito:', err);
        }
    };

    const handleDislike = async () => {
        try {
            if (onDislike) {
                onDislike();
            } else {
                await dislikeProperty(property.id);
                setShowMenu(false);
                if (onClose) onClose();
            }
        } catch (err) {
            console.error('Error al marcar como dislike:', err);
        }
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
                
                {/* Botón de favorito */}
                <div 
                    className="p-2 rounded-full cursor-pointer transition-all duration-200"
                    onClick={handleToggleFavorite}
                    onMouseEnter={() => setIsHeartHovered(true)}
                    onMouseLeave={() => setIsHeartHovered(false)}
                >
                    <svg 
                        className={`w-6 h-6 transition-all duration-200 ${
                            isFavorite 
                                ? 'text-red-500 fill-current' 
                                : isHeartHovered 
                                    ? 'text-red-400 fill-red-300' 
                                    : 'fill-none stroke-gray-300 stroke-1'
                        }`} 
                        viewBox="0 0 24 24"
                    >
                        <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
                    </svg>
                </div>

                {/* Botón de menú (tres puntos) */}
                <div className="relative">
                    <button 
                        className="p-2 rounded-full cursor-pointer transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800 ml-1"
                        onClick={() => setShowMenu(!showMenu)}
                    >
                        <svg className="w-6 h-6 dark:text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                        </svg>
                    </button>

                    {/* Menú desplegable */}
                    {showMenu && (
                        <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg z-10">
                            <div className="py-1">
                                <button 
                                    className="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center"
                                    onClick={handleDislike}
                                >
                                    <svg className="w-5 h-5 mr-2 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.096c.5 0 .905-.405.905-.904 0-.715.211-1.413.608-2.008L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5" />
                                    </svg>
                                    Ya no me gusta
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Chat area */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                {isLoading ? (
                    <div className="flex justify-center items-center h-full">
                        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
                    </div>
                ) : error ? (
                    <div className="text-center text-red-500 py-8">
                        {error}
                    </div>
                ) : notes.length === 0 ? (
                    <div className="text-center text-gray-500 py-8">
                        No hay notas aún
                    </div>
                ) : (
                    notes.map(note => (
                        <div 
                            key={note.id} 
                            className="bg-blue-500 text-white p-3 rounded-lg max-w-[80%] ml-auto relative group"
                        >
                            <p>{note.text}</p>
                            <p className="text-xs opacity-75 mt-1">
                                {new Date(note.created_at).toLocaleTimeString([], { 
                                    hour: '2-digit', 
                                    minute: '2-digit' 
                                })}
                            </p>
                            <button 
                                onClick={() => handleDeleteNote(note.id)}
                                className="absolute -right-2 -top-2 bg-red-500 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                            >
                                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
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
                        disabled={!newNote.trim() || isLoading}
                    >
                        Enviar
                    </button>
                </div>
            </div>
        </div>
    );
} 
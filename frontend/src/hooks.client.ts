import type { HandleClientError } from '@sveltejs/kit';

export const handleError: HandleClientError = ({ error }) => {
	const message = (error as any)?.message ?? '';
	
	// Si falla la carga de un chunk (ej: después de un deploy), recargar SOLO si es un 404 real
	// de lo contrario, intentar de nuevo no ayuda. Dejamos que SvelteKit maneje el error.
	console.warn('Client error:', message);

	return {
		message: 'An unexpected error occurred.',
		error
	};
};

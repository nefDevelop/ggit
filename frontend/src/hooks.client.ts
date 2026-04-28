import type { HandleClientError } from '@sveltejs/kit';

export const handleError: HandleClientError = ({ error }) => {
	const message = (error as any)?.message ?? '';
	
	// If a chunk fails to load (e.g. after a deployment), reload the page to get the latest manifest
	if (message.includes('error loading dynamically imported module') || 
        message.includes('Failed to fetch dynamically imported module')) {
		window.location.reload();
	}

	return {
		message: 'An unexpected error occurred.',
		error
	};
};

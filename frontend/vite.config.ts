import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:3005',
				changeOrigin: true
			},
			'/install': {
				target: 'http://localhost:3005',
				changeOrigin: true
			},
			'/avatars': {
				target: 'http://localhost:3005',
				changeOrigin: true
			},
			'/repo-avatars': {
				target: 'http://localhost:3005',
				changeOrigin: true
			},
		}
	}
});

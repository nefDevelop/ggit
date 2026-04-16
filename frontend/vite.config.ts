import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:3000',
				changeOrigin: true
			},
			'/install': {
				target: 'http://127.0.0.1:3000',
				changeOrigin: true
			},
			'/explore': {
				target: 'http://127.0.0.1:3000',
				changeOrigin: true
			},
			'/admin': {
				target: 'http://127.0.0.1:3000',
				changeOrigin: true
			},
			'/avatars': {
				target: 'http://127.0.0.1:3000',
				changeOrigin: true
			},
			'/repo-avatars': {
				target: 'http://127.0.0.1:3000',
				changeOrigin: true
			},
		}
	}
});

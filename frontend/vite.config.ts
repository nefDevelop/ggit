import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': 'http://localhost:3000',
			'/install': 'http://localhost:3000',
			'/user': 'http://localhost:3000',
			'/explore': 'http://localhost:3000',
			'/admin': 'http://localhost:3000',
			'/avatars': 'http://localhost:3000',
			'/repo-avatars': 'http://localhost:3000',
		}
	}
});

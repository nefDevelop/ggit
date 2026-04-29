import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	const backendPort = env.VITE_BACKEND_PORT || '3005';
	const target = `http://localhost:${backendPort}`;

	return {
		plugins: [sveltekit()],
		server: {
			proxy: {
				'/api': {
					target,
					changeOrigin: true
				},
				'/health': {
					target,
					changeOrigin: true
				},
				'/avatars': {
					target,
					changeOrigin: true
				},
				'/repo-avatars': {
					target,
					changeOrigin: true
				},
			}
		}
	};
});

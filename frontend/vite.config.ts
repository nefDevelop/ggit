import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	const backendPort = env.VITE_BACKEND_PORT || '3005';
	const target = `http://localhost:${backendPort}`;

	return {
		plugins: [sveltekit()],
		build: {
			target: 'esnext',
			cssMinify: true,
			rollupOptions: {
				output: {
					manualChunks: (id) => {
						if (id.includes('lucide-svelte')) return 'icons';
						if (id.includes('highlight.js')) return 'hljs';
						if (id.includes('marked') || id.includes('dompurify')) return 'utils';
						if (id.includes('node_modules')) return 'vendor';
					}
				}
			},
			chunkSizeWarningLimit: 1000
		},
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

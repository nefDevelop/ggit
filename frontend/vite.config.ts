import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import viteCompression from 'vite-plugin-compression';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	const backendPort = env.VITE_BACKEND_PORT || '3005';
	const target = `http://localhost:${backendPort}`;

	return {
		plugins: [
			sveltekit(),
			viteCompression({
				algorithm: 'brotliCompress',
				ext: '.br',
				threshold: 1024,
			}),
		],
		build: {
			target: 'esnext',
			cssMinify: true,
			rollupOptions: {
				output: {
					manualChunks: (id) => {
						if (id.includes('lucide-svelte')) return 'icons';
						if (id.includes('highlight.js')) return 'hljs';
						if (id.includes('marked')) return 'marked';
						if (id.includes('dompurify')) return 'dompurify';
						if (id.includes('node_modules')) return 'vendor';
					}
				}
			},
			chunkSizeWarningLimit: 500
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

<script lang="ts">
	import { page } from '$app/stores';
	import { ArrowLeft, Home, RefreshCw, Terminal } from 'lucide-svelte';

	let status = $derived($page.status);
	let message = $derived($page.error?.message || 'Algo salió mal');

	const errorDetails = {
		404: {
			title: '404 - Ruta no encontrada',
			desc: 'Parece que el código que buscas se ha movido a otra rama o nunca existió.',
			illustration: 'lost'
		},
		500: {
			title: '500 - Error Interno',
			desc: 'Nuestros servidores han encontrado un conflicto imprevisto. Estamos trabajando en el merge.',
			illustration: 'broken'
		}
	} as const;

	let details = $derived(errorDetails[status as keyof typeof errorDetails] || {
		title: `${status} - Error Inesperado`,
		desc: message,
		illustration: 'default'
	});
</script>

<div class="error-container">
	<div class="error-card">
		<div class="error-badge">
			<div class="badge-icon">
				{#if status === 404}
					<div class="glitch-text" data-text="404">404</div>
				{:else if status === 500}
					<RefreshCw size={64} class="animate-spin-slow" />
				{:else}
					<Terminal size={64} />
				{/if}
			</div>
		</div>

		<h1 class="error-title">{details.title}</h1>
		<p class="error-desc">{details.desc}</p>

		{#if $page.error?.stack}
			<div class="error-stack">
				<pre><code>{$page.error.stack}</code></pre>
			</div>
		{/if}

		<div class="error-actions">
			<button class="btn-primary" onclick={() => history.back()}>
				<ArrowLeft size={18} />
				Volver Atrás
			</button>
			<a href="/" class="btn-secondary">
				<Home size={18} />
				Ir al Dashboard
			</a>
		</div>
	</div>

	<div class="decoration-blobs">
		<div class="blob blob-1"></div>
		<div class="blob blob-2"></div>
	</div>
</div>

<style>
	.error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 80vh;
		padding: 2rem;
		position: relative;
		overflow: hidden;
	}

	.error-card {
		background: var(--island-bg);
		border: 1px solid var(--border-color);
		border-radius: 28px;
		padding: 3rem;
		max-width: 600px;
		width: 100%;
		text-align: center;
		box-shadow: var(--shadow-xl);
		z-index: 10;
		backdrop-filter: blur(20px);
	}

	.error-badge {
		margin-bottom: 2rem;
		color: var(--primary);
		display: flex;
		justify-content: center;
	}

	.glitch-text {
		font-size: 5rem;
		font-weight: 900;
		position: relative;
		color: var(--primary);
		letter-spacing: 0.5rem;
	}

	.error-title {
		font-size: 2rem;
		font-weight: 600;
		margin-bottom: 1rem;
		color: var(--text-primary);
	}

	.error-desc {
		font-size: 1.1rem;
		color: var(--text-secondary);
		margin-bottom: 2.5rem;
		line-height: 1.6;
	}

	.error-actions {
		display: flex;
		gap: 1rem;
		justify-content: center;
	}

	.btn-primary, .btn-secondary {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1.5rem;
		border-radius: 12px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
		border: none;
	}

	.btn-primary {
		background: var(--primary);
		color: white;
	}

	.btn-primary:hover {
		background: var(--primary-hover);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(var(--primary-rgb), 0.3);
	}

	.btn-secondary {
		background: var(--surface-variant);
		color: var(--text-primary);
		text-decoration: none;
	}

	.btn-secondary:hover {
		background: var(--surface-variant-hover);
		transform: translateY(-2px);
	}

	.error-stack {
		margin: 2rem 0;
		background: rgba(0, 0, 0, 0.2);
		padding: 1rem;
		border-radius: 12px;
		text-align: left;
		font-size: 0.85rem;
		max-height: 200px;
		overflow-y: auto;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.error-stack pre {
		margin: 0;
		white-space: pre-wrap;
		word-break: break-all;
		color: #ffb4ab;
	}

	/* Animations */
	@keyframes spin-slow {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	.animate-spin-slow {
		animation: spin-slow 8s linear infinite;
	}

	/* Decoration Blobs */
	.decoration-blobs {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 1;
		pointer-events: none;
	}

	.blob {
		position: absolute;
		border-radius: 50%;
		filter: blur(80px);
		opacity: 0.15;
	}

	.blob-1 {
		width: 400px;
		height: 400px;
		background: var(--primary);
		top: -100px;
		right: -100px;
	}

	.blob-2 {
		width: 300px;
		height: 300px;
		background: var(--secondary);
		bottom: -50px;
		left: -50px;
	}

	@media (max-width: 640px) {
		.error-actions {
			flex-direction: column;
		}
		
		.error-card {
			padding: 2rem;
		}
	}
</style>

<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import { authState, checkAuth } from '$lib/auth.svelte';
	import { onMount } from 'svelte';
	import Navigation from '$lib/components/Navigation.svelte';

	let { children } = $props();

	onMount(() => {
		checkAuth();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<script>
		// Prevent FOUC (Flash of Unstyled Content) by setting theme synchronously
		(function() {
			try {
				const theme = localStorage.getItem('theme') || 'default';
				const scheme = localStorage.getItem('color-scheme') || 'auto';
				
				document.documentElement.setAttribute('data-theme', theme);
				
				if (scheme === 'dark' || (scheme === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
					document.documentElement.setAttribute('data-color-scheme', 'dark');
				} else {
					document.documentElement.setAttribute('data-color-scheme', 'light');
				}
			} catch (e) {}
		})();
	</script>
</svelte:head>

<div id="app-layout">
	{#if authState.isAuthenticated}
        <div class="navigation-sidebar">
		    <Navigation />
        </div>
	{/if}

	<main id="main-content">
		{@render children()}
	</main>
</div>

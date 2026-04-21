<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import { authState, checkAuth } from '$lib/auth.svelte';
	import { i18n } from '$lib/i18n.svelte';
	import { onMount } from 'svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import RightSidebar from '$lib/components/RightSidebar.svelte';
	import { page } from '$app/stores';

	let { children } = $props();

	onMount(() => {
		checkAuth();
		i18n.loadLocale(i18n.locale);
	});

	const isDashboard = $derived($page.url.pathname === '/' || $page.url.pathname === '/dashboard');
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

<div id="app-layout" class:has-right-sidebar={authState.isAuthenticated && isDashboard}>
	{#if authState.isAuthenticated}
        <div class="navigation-sidebar">
		    <Navigation />
        </div>
	{/if}

    <div id="main-container-wrapper">
        <main id="main-content">
            {@render children()}
        </main>
    </div>

    {#if authState.isAuthenticated && isDashboard}
        <RightSidebar />
    {/if}
</div>

<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import { authState, checkAuth, checkInstallation } from '$lib/auth.svelte';
	import { i18n } from '$lib/i18n.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Navigation from '$lib/components/Navigation.svelte';
	import RightSidebar from '$lib/components/RightSidebar.svelte';
	import GlobalDialog from '$lib/components/GlobalDialog.svelte';
	import GlobalToast from '$lib/components/GlobalToast.svelte';
	import { page } from '$app/stores';

	let { children } = $props();

	onMount(async () => {
		await checkAuth();
		const installed = await checkInstallation();

		// Redirection Logic
		const path = $page.url.pathname;
		if (!installed && path !== '/install') {
			goto('/install');
		} else if (installed && path === '/install') {
			goto('/');
		}

		// Cargar el idioma guardado o el del navegador
		const savedLocale = localStorage.getItem('ggit-locale') || navigator.language || 'es-ES';
		i18n.init(savedLocale);
	});

	const isDashboard = $derived($page.url.pathname === '/' || $page.url.pathname === '/dashboard');
	const showRightSidebar = $derived(authState.isAuthenticated && isDashboard && !authState.user?.disable_news_feed);
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

<div id="app-layout" class:has-right-sidebar={showRightSidebar}>
	{#if authState.isAuthenticated}
        <div class="navigation-sidebar">
		    <Navigation />
        </div>
	{/if}

    <div id="main-container-wrapper">
        <main id="main-content">
            {#if authState.loading}
                <div class="initial-loading">
                    <div class="spinner"></div>
                </div>
            {:else}
                {@render children()}
            {/if}
        </main>
    </div>

    {#if showRightSidebar}
        <RightSidebar />
    {/if}

    <GlobalDialog />
    <GlobalToast />
</div>

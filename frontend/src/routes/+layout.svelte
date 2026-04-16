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
</svelte:head>

<div id="app-layout">
	{#if authState.isAuthenticated}
		<Navigation />
	{/if}
	<main id="main-content">
		{@render children()}
	</main>
</div>

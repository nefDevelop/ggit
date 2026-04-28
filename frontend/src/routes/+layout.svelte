<script lang="ts">
	import favicon from "$lib/assets/favicon.svg";
	import "../app.css";
	import { authState, checkAuth, checkInstallation } from "$lib/auth.svelte";
	import { i18n } from "$lib/i18n.svelte";
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import Navigation from "$lib/components/Navigation.svelte";
	import RightSidebar from "$lib/components/RightSidebar.svelte";
	import GlobalDialog from "$lib/components/GlobalDialog.svelte";
	import GlobalToast from "$lib/components/GlobalToast.svelte";
	import { page } from "$app/stores";

	let { children } = $props();

	let isSurvivalMode = $state(false);

	async function checkSystemHealth() {
		try {
			const res = await fetch("/health");
			if (res.ok) {
				const health = await res.json();
				isSurvivalMode = health.survival_mode;
			}
		} catch (e) {
			console.error("Health check failed", e);
		}
	}

	onMount(() => {
		const init = async () => {
			await checkAuth();
			const installed = await checkInstallation();

			// Redirection Logic
			const path = $page.url.pathname;
			if (!installed && path !== "/install") {
				goto("/install");
			} else if (installed && path === "/install") {
				goto("/");
			}

			// Cargar el idioma guardado o el del navegador
			const savedLocale =
				localStorage.getItem("ggit-locale") ||
				navigator.language ||
				"es-ES";
			i18n.init(savedLocale);

			checkSystemHealth();
		};

		init();
		
		const healthInterval = setInterval(checkSystemHealth, 30000);
		return () => clearInterval(healthInterval);
	});

	const isDashboard = $derived(
		$page.url.pathname === "/" || $page.url.pathname === "/dashboard",
	);
	const showRightSidebar = $derived(
		authState.isAuthenticated &&
			isDashboard &&
			!authState.user?.disable_news_feed,
	);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<script>
		// Prevent FOUC (Flash of Unstyled Content) by setting theme synchronously
		(function () {
			try {
				const theme = localStorage.getItem("theme") || "default";
				const scheme = localStorage.getItem("color-scheme") || "auto";

				document.documentElement.setAttribute("data-theme", theme);

				if (
					scheme === "dark" ||
					(scheme === "auto" &&
						window.matchMedia("(prefers-color-scheme: dark)")
							.matches)
				) {
					document.documentElement.setAttribute(
						"data-color-scheme",
						"dark",
					);
				} else {
					document.documentElement.setAttribute(
						"data-color-scheme",
						"light",
					);
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
		{#if isSurvivalMode}
			<div class="survival-mode-banner">
				<span class="material-symbols-rounded">warning</span>
				<p>
					<strong>Modo Supervivencia Activo:</strong> El espacio en disco
					es crítico. Las operaciones de escritura están bloqueadas temporalmente.
				</p>
			</div>
		{/if}
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

<style>
	.survival-mode-banner {
		background: var(--md-sys-color-error-container);
		color: var(--md-sys-color-on-error-container);
		padding: 0.75rem 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		border-bottom: 1px solid var(--md-sys-color-error);
		animation: slide-down 0.5s cubic-bezier(0.4, 0, 0.2, 1);
		z-index: 100;
	}
	.survival-mode-banner p {
		margin: 0;
		font-size: 0.9rem;
		font-weight: 600;
	}
	.survival-mode-banner span {
		color: var(--error);
	}

	@keyframes slide-down {
		from {
			transform: translateY(-100%);
		}
		to {
			transform: translateY(0);
		}
	}
</style>

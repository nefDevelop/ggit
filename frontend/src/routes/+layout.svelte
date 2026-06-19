<script lang="ts">
	import favicon from "$lib/assets/favicon.svg";
	import "../app.css";
	import { apiFetch } from "$lib/api";
	import { authState, checkAuth, checkInstallation } from "$lib/auth.svelte";
	import { i18n } from "$lib/i18n.svelte";
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import Navigation from "$lib/components/Navigation.svelte";
	import RightSidebar from "$lib/components/RightSidebar.svelte";
	import GlobalDialog from "$lib/components/GlobalDialog.svelte";
	import GlobalToast from "$lib/components/GlobalToast.svelte";
	import NotificationsPanel from "$lib/components/NotificationsPanel.svelte";
	import { toast } from "$lib/toast.svelte";
	import { page } from "$app/stores";
	
import TriangleAlert from "lucide-svelte/icons/triangle-alert";


	let { children } = $props();

	let isSurvivalMode = $state(false);

	async function checkSystemHealth() {
		try {
			const res = await apiFetch("/health", { silent: true });
			if (res.ok) {
				const health = res.data;
				isSurvivalMode = health?.survival_mode || false;
			}
		} catch (e) {
			// Silent fail for health check
		}
	}

	onMount(() => {
		// Global error handler to prevent full page reload on JS errors
		const handleGlobalError = (event: ErrorEvent) => {
			event.preventDefault();
			console.error('🔴 GLOBAL ERROR:', event.error || event.message);
		};
		const handleRejection = (event: PromiseRejectionEvent) => {
			event.preventDefault();
			console.error('🔴 UNHANDLED REJECTION:', event.reason);
		};
		window.addEventListener('error', handleGlobalError);
		window.addEventListener('unhandledrejection', handleRejection);

		const init = async () => {
			// Run auth and installation checks in parallel to save time
			const [authResult, installResult] = await Promise.all([
				checkAuth(),
				checkInstallation()
			]);

			// Cargar el idioma guardado o el del navegador
			const savedLocale =
				localStorage.getItem("ggit-locale") ||
				navigator.language ||
				"es-ES";
			await i18n.init(savedLocale);

			// Redirection Logic
			const path = $page.url.pathname;
			if (!installResult && path !== "/install") {
				goto("/install");
			} else if (installResult && path === "/install") {
				goto("/");
			}

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
	const hasSubNav = $derived(
		$page.url.pathname.startsWith("/admin") || 
		$page.url.pathname.startsWith("/user/settings")
	);
</script>

<svelte:head>
	<title>Ggit</title>
	<link rel="icon" href={favicon} />
	<script>
		// Prevent FOUC (Flash of Unstyled Content) by setting theme synchronously
		(function () {
			try {
				const theme = localStorage.getItem("theme") || "default";
				const scheme = localStorage.getItem("color-scheme") || "auto";

				document.documentElement.setAttribute("data-theme", theme);

				if (theme === "custom") {
					const customColor = localStorage.getItem("custom-theme-color") || "#0b57d0";
					const isDark = scheme === "dark" || (scheme === "auto" && window.matchMedia("(prefers-color-scheme: dark)").matches);
					
					document.documentElement.style.setProperty("--primary-color", customColor);
					document.documentElement.style.setProperty("--md-sys-color-primary", customColor);
					
					// Basic background/surface fallbacks for custom themes to prevent white flash in dark mode
					if (isDark) {
						document.documentElement.style.setProperty("--md-sys-color-background", "#1a1c1e");
						document.documentElement.style.setProperty("--md-sys-color-surface", "#1e2124");
						document.documentElement.style.setProperty("--md-sys-color-on-surface", "#e2e2e6");
					} else {
						document.documentElement.style.setProperty("--md-sys-color-background", "#f8f9ff");
						document.documentElement.style.setProperty("--md-sys-color-surface", "#ffffff");
						document.documentElement.style.setProperty("--md-sys-color-on-surface", "#1b1b1f");
					}

					if (customColor.startsWith("#")) {
						const r = parseInt(customColor.slice(1, 3), 16);
						const g = parseInt(customColor.slice(3, 5), 16);
						const b = parseInt(customColor.slice(5, 7), 16);
						document.documentElement.style.setProperty("--primary-rgb", `${r}, ${g}, ${b}`);
					}
				}

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

<div id="app-layout" class:has-right-sidebar={showRightSidebar} class:has-sub-nav={hasSubNav}>
	{#if authState.isAuthenticated}
		<div class="navigation-sidebar">
			<Navigation />
		</div>
	{/if}

	<div id="main-container-wrapper">
		{#if isSurvivalMode}
			<div class="survival-mode-banner">
				<TriangleAlert size={20} />
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
		<NotificationsPanel bind:show={toast.showNotices} onclose={() => (toast.showNotices = false)} />
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
	.survival-mode-banner :global(svg) {
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

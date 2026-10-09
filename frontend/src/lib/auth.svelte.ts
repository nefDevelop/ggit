import { apiFetch } from "$lib/api";
export const authState = $state({
	user: null as any,
	isAuthenticated: false,
	isInstalled: true,
	// Tri-estado: hasta que checkInstallation() resuelve no se sabe si hay
	// que mostrar /install. El layout no pinta la ruta hasta entonces para
	// no flashear la pantalla de instalación.
	installChecked: false,
	loading: true,
});

export async function checkAuth() {
	try {
		const res = await apiFetch('/api/v2/user/check');
		if (res.ok) {
			const json = await res.json();
			if (json.status === 'ok' && json.user) {
				authState.user = json.user;
				authState.isAuthenticated = true;
			} else {
				authState.user = null;
				authState.isAuthenticated = false;
			}
		} else {
			authState.user = null;
			authState.isAuthenticated = false;
		}
	} catch (err) {
		console.error('Failed to check auth', err);
		authState.user = null;
		authState.isAuthenticated = false;
	} finally {
		authState.loading = false;
	}
}

export async function checkInstallation() {
	try {
		const res = await apiFetch('/api/v2/install/settings', { silent: true });
		// 403 Forbidden means it's already installed.
		if (res.status === 403) {
			authState.isInstalled = true;
			return true;
		}
		if (res.ok) {
			// Solo un 2xx con el cuerpo esperado significa "sin instalar".
			// Un 200 con otra cosa (HTML del fallback de la SPA, respuesta
			// del proxy, caché) se trata como instalado para no redirigir
			// a /install por error.
			const data = await res.json();
			if (data && typeof data === 'object' && 'default_settings' in data) {
				authState.isInstalled = false;
				return false;
			}
			return true;
		}
		return true; // Assume installed on error to be safe
	} catch (err) {
		console.error('Failed to check installation', err);
		return true;
	} finally {
		authState.installChecked = true;
	}
}

export function logoutUser() {
	authState.user = null;
	authState.isAuthenticated = false;
}

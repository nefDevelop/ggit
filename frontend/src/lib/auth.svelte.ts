import { apiFetch } from "$lib/api";
export const authState = $state({
	user: null as any,
	isAuthenticated: false,
	isInstalled: true,
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
		// 200 OK means it's NOT installed yet.
		if (res.status === 403) {
			authState.isInstalled = true;
			return true;
		} else if (res.ok) {
			authState.isInstalled = false;
			return false;
		}
		return true; // Assume installed on error to be safe
	} catch (err) {
		console.error('Failed to check installation', err);
		return true;
	}
}

export function logoutUser() {
	authState.user = null;
	authState.isAuthenticated = false;
}

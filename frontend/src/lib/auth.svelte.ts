import { apiFetch } from "$lib/api";
export const authState = $state({
	user: null as any,
	isAuthenticated: false,
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

export function logoutUser() {
	authState.user = null;
	authState.isAuthenticated = false;
}

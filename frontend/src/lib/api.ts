/**
 * Get cookie value by name
 */
export function getCookie(name: string): string {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift() || '';
  return '';
}

/**
 * Enhanced fetch that automatically adds CSRF token to mutating requests
 */
export async function apiFetch(url: string, options: RequestInit = {}) {
  const method = (options.method || 'GET').toUpperCase();
  const mutableMethods = ['POST', 'PUT', 'DELETE', 'PATCH'];

  if (mutableMethods.includes(method)) {
    const csrfToken = getCookie('XSRF-TOKEN');
    if (csrfToken) {
      options.headers = {
        ...options.headers,
        'X-CSRF-Token': csrfToken
      };
    }
  }

  return fetch(url, options);
}

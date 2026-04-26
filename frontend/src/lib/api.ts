import { toast } from "./toast.svelte";

/**
 * Get cookie value by name
 */
export function getCookie(name: string): string {
  if (typeof document === 'undefined') return '';
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift() || '';
  return '';
}

interface ApiOptions extends RequestInit {
  silent?: boolean;
}

/**
 * Enhanced fetch that:
 * 1. Automatically adds CSRF token to mutating requests
 * 2. Handles common error status codes globally (Toasts)
 * 3. Provides a cleaner interface for API calls
 */
export async function apiFetch(url: string, options: ApiOptions = {}) {
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

  try {
    const response = await fetch(url, options);

    if (!response.ok && !options.silent) {
      // Intentar extraer mensaje de error del JSON
      let errorMessage = `Error: ${response.statusText || response.status}`;
      try {
        const data = await response.clone().json();
        if (data.message) errorMessage = data.message;
        else if (data.error) errorMessage = data.error;
      } catch (e) {
        // No es JSON, usar el texto por defecto
      }

      if (response.status === 401) {
        toast.error("Sesión expirada. Por favor, vuelve a iniciar sesión.");
        // Podríamos redirigir aquí si fuera necesario
      } else if (response.status === 403) {
        toast.error("No tienes permisos para realizar esta acción.");
      } else if (response.status >= 500) {
        toast.error("Error en el servidor. Inténtalo más tarde.");
      } else {
        toast.error(errorMessage);
      }
    }

    return response;
  } catch (error) {
    if (!options.silent) {
      toast.error("Error de conexión. Verifica tu internet.");
    }
    throw error;
  }
}

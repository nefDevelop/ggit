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
 * 3. Safely parses JSON responses
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
    
    // Intentar extraer JSON de forma segura
    let data: any = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
      try {
        const text = await response.text();
        if (text) data = JSON.parse(text);
      } catch (e) {
        console.error("Error parsing JSON response", e);
      }
    }

    if (!response.ok && !options.silent) {
      let errorMessage = data?.message || data?.error || `Error: ${response.statusText || response.status}`;

      if (response.status === 401) {
        toast.error("Sesión expirada. Por favor, vuelve a iniciar sesión.");
      } else if (response.status === 403) {
        toast.error("No tienes permisos para realizar esta acción.");
      } else if (response.status >= 500) {
        toast.error("Error en el servidor. Inténtalo más tarde.");
      } else {
        toast.error(errorMessage);
      }
    }

    // Devolvemos un proxy de la respuesta que incluye el data ya parseado
    // para mantener compatibilidad con el código existente que usa .json()
    return {
      ...response,
      ok: response.ok,
      status: response.status,
      headers: response.headers,
      data, // El JSON ya parseado
      json: async () => data, // Mock de .json() para compatibilidad
      text: async () => JSON.stringify(data)
    };
  } catch (error) {
    if (!options.silent) {
      toast.error("Error de conexión. Verifica tu internet.");
    }
    throw error;
  }
}

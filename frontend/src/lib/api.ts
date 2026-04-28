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
    
    // Intentar extraer texto y JSON de forma segura
    let data: any = null;
    let textBody: string = "";
    
    try {
      textBody = await response.text();
      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json") && textBody) {
        data = JSON.parse(textBody);
      }
    } catch (e) {
      console.error("Error reading or parsing response", e);
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

    // Devolvemos un objeto que imita la interfaz de Response pero con data pre-parseado
    return {
      ok: response.ok,
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
      url: response.url,
      data, // El JSON ya parseado
      json: async () => data, 
      text: async () => textBody || JSON.stringify(data),
      // Añadimos el original por si alguien necesita acceder a propiedades específicas
      raw: response 
    };
  } catch (error) {
    if (!options.silent) {
      toast.error("Error de conexión. Verifica que el servidor esté encendido.");
    }
    throw error;
  }
}

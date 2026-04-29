import { test, expect } from '@playwright/test';

/**
 * Este test realiza un "crawling" automático de la aplicación.
 * Navega por todos los enlaces internos encontrados, verificando que
 * no haya enlaces rotos (404) ni errores de servidor (500).
 */
test('crawler: navegar por todos los enlaces internos y detectar errores de recursos', async ({ page, baseURL }) => {
  const visited = new Set<string>();
  const toVisit: string[] = ['/'];
  const errors: string[] = [];

  // 1. Capturar excepciones de JavaScript en la página (Uncaught errors)
  page.on('pageerror', (exception) => {
    errors.push(`[JS Error] en ${page.url()}: ${exception.message}`);
  });

  // 2. Capturar errores de consola (Svelte suele loguear fallos críticos aquí)
  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      errors.push(`[Console Error] en ${page.url()}: ${msg.text()}`);
    }
  });

  // 3. Capturar errores de red (CORS, 404 de assets, fallos de conexión)
  page.on('requestfailed', (request) => {
    const failure = request.failure();
    // Filtramos ruidos comunes como extensiones de navegador si fuera el caso
    if (request.url().startsWith(baseURL || 'http://localhost:5173')) {
      errors.push(
        `[Network Error] ${failure?.errorText || 'Unknown'} - ${request.method()} ${request.url()} (Desde: ${page.url()})`
      );
    }
  });

  const excludePaths = ['/user/logout', '/logout', '/delete', '/remove', '#'];

  while (toVisit.length > 0) {
    const currentPath = toVisit.pop()!;
    const normalizedPath = currentPath.endsWith('/') && currentPath !== '/' ? currentPath.slice(0, -1) : currentPath;

    if (visited.has(normalizedPath)) continue;
    visited.add(normalizedPath);

    console.log(`Crawling: ${normalizedPath}`);

    const response = await page.goto(normalizedPath, {
      waitUntil: 'networkidle',
      timeout: 30000,
    });

    // 4. Validar status y "App Panic" (Pantalla en blanco o mensaje de error de SvelteKit)
    if (response) {
      expect(response.status(), `La página devolvió un error: ${response.status()} en ${normalizedPath}`)
        .toBeLessThan(400);
      
      const bodyContent = await page.innerText('body');
      if (bodyContent.trim().length < 10) {
        errors.push(`[White Screen] La página parece estar vacía en: ${normalizedPath}`);
      }
      if (bodyContent.includes('Internal Error') || bodyContent.includes('Internal Server Error')) {
        errors.push(`[App Panic] Mensaje de error crítico detectado en el texto de la página: ${normalizedPath}`);
      }
    }

    // Extraer enlaces (mismo proceso)
    const anchors = await page.locator('a').evaluateAll((elements: HTMLAnchorElement[]) => 
      elements.map(a => ({ href: a.href, host: a.host }))
    );

    for (const link of anchors) {
      const url = new URL(link.href);
      const isInternal = url.origin === (baseURL || 'http://localhost:5173');
      const path = url.pathname + url.search;

      if (isInternal && !visited.has(path) && !excludePaths.some(p => path.includes(p)) && !path.match(/\.(pdf|zip|tar|gz|png|jpg|jpeg)$/i)) {
        toVisit.push(path);
      }
    }
  }

  // 3. Reporte final de errores acumulados
  if (errors.length > 0) {
    console.error('--- ERRORES DETECTADOS DURANTE EL CRAWLING ---');
    errors.forEach(err => console.error(err));
    
    // El test fallará si hay errores de consola o red
    expect(errors, 'Se detectaron errores de red o JavaScript durante la navegación').toHaveLength(0);
  }

  console.log(`Crawling finalizado. Páginas visitadas: ${visited.size}`);
});

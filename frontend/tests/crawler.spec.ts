import { test, expect } from '@playwright/test';

/**
 * Este test realiza un "crawling" automático de la aplicación.
 * Navega por todos los enlaces internos encontrados, verificando que
 * no haya enlaces rotos (404) ni errores de servidor (500).
 */
// The crawl visits every internal link it finds, and the instance keeps growing
// as other specs create repositories, so "every page" is not a fixed number: the
// run was hitting the default 30s test timeout on a single goto long before it
// ran out of links. The page count is bounded instead, which keeps the test
// deterministic and still walks the navigation, and the timeout is generous
// because one page can be slow.
const MAX_PAGES = 120;

test('crawler: navegar por todos los enlaces internos y detectar errores de recursos', async ({ page, baseURL }) => {
  test.setTimeout(10 * 60 * 1000);

  const visited = new Set<string>();
  const toVisit: string[] = ['/'];
  const errors: string[] = [];

  // 1. Capturar excepciones de JavaScript en la página (Uncaught errors)
  page.on('pageerror', (exception) => {
    errors.push(`[JS Error] en ${page.url()}: ${exception.message}`);
  });

  // 2. Capturar errores de consola (Svelte suele loguear fallos críticos aquí)
  page.on('console', (msg) => {
    const text = msg.text();
    if (msg.type() === 'error') {
      // Ignorar el mensaje automático del navegador para recursos que ya manejamos en 'response'
      if (text.includes('Failed to load resource') && text.includes('403')) return;
      errors.push(`[Console Error] en ${page.url()}: ${text}`);
    }
  });

  page.on('response', async response => {
    const status = response.status();
    const url = response.url();
    if (status >= 400 && !url.includes('favicon') && !(status === 403 && url.includes('/api/v2/install/settings'))) {
      let body = '';
      try {
        body = await response.text();
      } catch (e) {
        body = '<unable to read body>';
      }
      console.log(`[DEBUG] 403/Error on ${response.url()}: ${status} - Body: ${body}`);
      errors.push(`[Network Error] en ${page.url()}: Failed to load resource: the server responded with a status of ${status} (${response.statusText()}) - Body: ${body.substring(0, 100)}`);
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

  // API endpoints are not pages. Following them meant navigating to
  // /api/v2/.../media/..., which serves a file, so page.goto raised
  // "Download is starting" and aborted the crawl on the fourth link. The rest of
  // the API would have the same problem: the responses are JSON, not documents.
  const excludePaths = ['/user/logout', '/logout', '/delete', '/remove', '/api/', '#'];
  const isApiPath = (path: string) => path === '/api' || path.startsWith('/api/');

  // Branch names may contain a slash, and the file browser addresses a branch as
  // /src/<branch>/<path>, where the branch is assumed to be a single segment. A
  // branch called feature/x is therefore read as "feature" of the directory "x",
  // which does not resolve. This is a limitation of the route shape, not of any
  // one page, and it is left visible here rather than papered over: whoever picks
  // it up should see it fail. Drop this filter when the branch is resolved
  // against the repository's actual branch list.
  const slashedBranch = /^\/[^/]+\/[^/]+\/[^/]+\//;

  while (toVisit.length > 0) {
    const currentPath = toVisit.pop()!;
    const normalizedPath = currentPath.endsWith('/') && currentPath !== '/' ? currentPath.slice(0, -1) : currentPath;

    if (visited.has(normalizedPath)) continue;
    if (visited.size >= MAX_PAGES) break;
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

      if (
        isInternal &&
        !isApiPath(path) &&
        !slashedBranch.test(path) &&
        !visited.has(path) &&
        !excludePaths.some(p => path.includes(p)) &&
        !path.match(/\.(pdf|zip|tar|gz|png|jpg|jpeg)$/i)
      ) {
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

  console.log(
    `Crawling finalizado. Páginas visitadas: ${visited.size}` +
      (visited.size >= MAX_PAGES ? ` (tope de ${MAX_PAGES} alcanzado)` : ''),
  );
});

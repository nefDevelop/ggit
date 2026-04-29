import { test, expect } from '@playwright/test';

test.describe('Admin Security', () => {
  const regularUser = {
    username: 'regular-user',
    password: 'password123',
    email: 'regular@example.com'
  };

  test('no-admin no debe poder acceder a rutas de administración', async ({ page }) => {
    // 1. Iniciar sesión como usuario normal
    await page.goto('/user/login');
    await page.fill('#username', regularUser.username);
    await page.fill('#password', regularUser.password);
    await page.click('button[type="submit"]');
    await page.waitForURL('/dashboard');

    // 2. Intentar acceder a panel de control de admin
    const adminRoutes = [
      '/admin',
      '/admin/users',
      '/admin/repos',
      '/admin/auths',
      '/admin/notices',
      '/admin/monitor',
      '/admin/repos/adoption'
    ];

    for (const route of adminRoutes) {
      await page.goto(route);
      // Debería redirigir al dashboard o mostrar un 403 (en SvelteKit suele redirigir si no tiene permiso)
      const url = page.url();
      expect(url).not.toContain('/admin');
    }
  });

  test('no-admin no debe poder llamar a la API de administración', async ({ request, page }) => {
    // Intentar llamar a un endpoint de admin directamente
    const response = await request.get('/api/v2/admin/stats');
    expect(response.status()).toBe(403);
  });
});

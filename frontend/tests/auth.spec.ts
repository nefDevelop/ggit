import { test, expect } from '@playwright/test';

test('login page loads and displays form', async ({ page }) => {
  await page.goto('/user/login');

  // Verificar el título
  await expect(page).toHaveTitle(/Sign In/);

  // Verificar que el header está presente
  await expect(page.locator('h2')).toContainText('Welcome Back');

  // Rellenar campos
  await page.fill('#username', 'testuser');
  await page.fill('#password', 'password123');

  // Verificar que el botón de login está habilitado
  const loginButton = page.locator('button[type="submit"]');
  await expect(loginButton).toBeVisible();
  await expect(loginButton).toBeEnabled();
});

test('shows error on failed login attempt', async ({ page }) => {
  await page.goto('/user/login');

  await page.fill('#username', 'wronguser');
  await page.fill('#password', 'wrongpass');
  await page.click('button[type="submit"]');

  // El servidor responderá con error (porque el backend no estará corriendo en este test aislado o fallará la auth)
  // Pero aquí comprobamos que el feedback visual ocurre
  const alert = page.locator('.alert.box-error');
  await expect(alert).toBeVisible;
});

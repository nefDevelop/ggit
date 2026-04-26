import { test, expect } from '@playwright/test';

test.describe('Smoke Tests - Public Areas', () => {
  test('homepage redirects or loads correctly', async ({ page }) => {
    await page.goto('/');
    // Dependiendo de si hay sesión o no, debería haber un título coherente
    await expect(page).toHaveTitle(/Ggit/i);
  });

  test('explore page shows repositories', async ({ page }) => {
    await page.goto('/explore/repos');
    const header = page.locator('h2');
    await expect(header).toBeVisible();
  });

  test('explore page shows users', async ({ page }) => {
    await page.goto('/explore/users');
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();
  });
});

test.describe('Critical Flow - Repository Settings', () => {
  // Nota: Estos tests asumen que el usuario entraría en una URL de repo
  // En un entorno de CI real, usaríamos un "Storage State" para estar logueados
  test('repo general settings structure', async ({ page }) => {
    // Intentamos ir a una ruta de settings (fallará el login, pero probamos la redirección)
    await page.goto('/user/settings/profile');
    await expect(page).toHaveURL(/\/user\/login/);
  });
});

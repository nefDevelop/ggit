import { test, expect } from '@playwright/test';

test.describe('User Settings', () => {
  test('profile settings page loads', async ({ page }) => {
    await page.goto('/user/settings?tab=profile');

    const heading = page.locator('h2:has-text("Profile")').or(
      page.locator('h2:has-text("Perfil")'),
    );
    await expect(heading).toBeVisible({ timeout: 10000 });
  });

  test('SSH keys page is accessible', async ({ page }) => {
    await page.goto('/user/settings?tab=keys');

    const heading = page.locator('h2:has-text("SSH")').or(
      page.locator('h2:has-text("Claves")'),
    );
    await expect(heading).toBeVisible({ timeout: 10000 });
  });

  test('security settings page loads', async ({ page }) => {
    await page.goto('/user/settings?tab=security');

    const heading = page.locator('h2:has-text("Security")').or(
      page.locator('h2:has-text("Seguridad")'),
    );
    await expect(heading).toBeVisible({ timeout: 10000 });
  });

  test('organizations page loads', async ({ page }) => {
    await page.goto('/user/settings?tab=orgs');

    const heading = page.locator('h2:has-text("Organizations")').or(
      page.locator('h2:has-text("Organizaciones")'),
    );
    await expect(heading).toBeVisible({ timeout: 10000 });
  });

  test('repositories page is accessible', async ({ page }) => {
    await page.goto('/user/settings?tab=repos');

    const heading = page.locator('h2:has-text("Repositories")').or(
      page.locator('h2:has-text("Repositorios")'),
    );
    await expect(heading).toBeVisible({ timeout: 10000 });
  });
});

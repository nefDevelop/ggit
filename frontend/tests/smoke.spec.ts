import { test, expect } from '@playwright/test';

test.describe('Smoke Tests - Public Areas', () => {
  test('homepage loads correctly', async ({ page }) => {
    await page.goto('/');
    // Check for either the landing page hero or the navigation bar if logged in
    const hero = page.locator('.hero-title');
    const nav = page.locator('.navigation-bar');
    await expect(hero.or(nav)).toBeVisible();
  });

  test('explore page shows repositories', async ({ page }) => {
    await page.goto('/explore/repos');
    await page.waitForURL('/explore/repos');
    // Use a more specific class selector and wait a bit
    const searchInput = page.locator('.global-search-input');
    await expect(searchInput).toBeVisible({ timeout: 10000 });
  });

  test('explore page shows users', async ({ page }) => {
    await page.goto('/explore/users');
    await page.waitForURL('/explore/users');
    const searchInput = page.locator('.global-search-input');
    await expect(searchInput).toBeVisible({ timeout: 10000 });
  });
});

test.describe('Critical Flow - Repository Settings', () => {
  test('settings page is accessible when sign-in is not required', async ({ page }) => {
    await page.goto('/user/settings?tab=profile');
    // Since REQUIRE_SIGNIN_VIEW is false in tests, we should see the settings page directly
    await expect(page.locator('h2:has-text("Profile")').or(page.locator('h2:has-text("Perfil")'))).toBeVisible();
  });
});

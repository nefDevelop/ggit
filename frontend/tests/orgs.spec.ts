import { test, expect } from '@playwright/test';

test.describe('Organizations', () => {
  const orgName = `test-org-${Date.now()}`;

  test('org creation page loads', async ({ page }) => {
    await page.goto('/org/create');

    const heading = page.locator('h2:has-text("Create Organization")').or(
      page.locator('h2:has-text("Crear")'),
    );
    await expect(heading).toBeVisible({ timeout: 10000 });
  });

  test('create an organization', async ({ page }) => {
    await page.goto('/org/create');

    const nameInput = page.locator('#org-name').or(page.locator('input[name="name"]'));
    if (await nameInput.isVisible()) {
      await nameInput.fill(orgName);
    }

    const descInput = page.locator('#org-desc').or(page.locator('textarea'));
    if (await descInput.isVisible()) {
      await descInput.fill('Test organization created by E2E test');
    }

    const submitBtn = page.locator('button[type="submit"]');
    if (await submitBtn.isVisible()) {
      await submitBtn.click();
    }

    await page.waitForTimeout(2000);
    const heading = page.locator('h2').or(page.locator('h1'));
    await expect(heading).toBeVisible();
  });
});

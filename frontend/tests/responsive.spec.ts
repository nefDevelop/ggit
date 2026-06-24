import { test, expect } from '@playwright/test';

test.describe('Responsive Layout', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('shows sidebar navigation on desktop', async ({ page }) => {
    await page.setViewportSize({ width: 1440, height: 900 });
    await page.waitForTimeout(500);
    const nav = page.locator('.navigation-bar');
    await expect(nav).toBeVisible();
  });

  test('hides sidebar nav and shows bottom bar on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    const nav = page.locator('.navigation-bar');
    await expect(nav).toBeVisible();
    const isFixed = await nav.evaluate(el => window.getComputedStyle(el).position === 'fixed');
    expect(isFixed).toBe(true);
  });

  test('hides nav brand footer on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    const brandSection = page.locator('.brand-section-bottom');
    const isHidden = await brandSection.evaluate(el => window.getComputedStyle(el).display === 'none');
    expect(isHidden).toBe(true);
  });

  test('hides user profile summary on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    const userProfile = page.locator('.user-profile-summary');
    const isHidden = await userProfile.evaluate(el => window.getComputedStyle(el).display === 'none');
    expect(isHidden).toBe(true);
  });

  test('shows mobile avatar on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    const mobileAvatar = page.locator('.mobile-avatar-item');
    const isVisible = await mobileAvatar.evaluate(el => window.getComputedStyle(el).display !== 'none');
    expect(isVisible).toBe(true);
  });

  test('hides right sidebar below 1300px', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.waitForTimeout(500);
    const rightSidebar = page.locator('.right-sidebar');
    if (await rightSidebar.count() > 0) {
      const isHidden = await rightSidebar.evaluate(el => window.getComputedStyle(el).display === 'none');
      expect(isHidden).toBe(true);
    }
  });

  test('hides dashboard sidebar below 1100px', async ({ page }) => {
    await page.setViewportSize({ width: 1024, height: 900 });
    await page.waitForTimeout(500);
    const dashboardSidebar = page.locator('.dashboard-sidebar');
    if (await dashboardSidebar.count() > 0) {
      const isHidden = await dashboardSidebar.evaluate(el => window.getComputedStyle(el).display === 'none');
      expect(isHidden).toBe(true);
    }
  });

  test('main container has correct margin on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    const mainContainer = page.locator('#main-container-wrapper');
    const marginBottom = await mainContainer.evaluate(el => window.getComputedStyle(el).marginBottom);
    expect(marginBottom).toBe('80px');
  });

  test('bottom nav bar is fixed at bottom on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    const nav = page.locator('.navigation-bar');
    const bottom = await nav.evaluate(el => window.getComputedStyle(el).bottom);
    expect(bottom).toBe('0px');
  });
});

import { test, expect } from '@playwright/test';

test.describe('Pull Requests', () => {
  const repoName = `test-prs-${Date.now()}`;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/repo/create');
    await page.fill('#repo-name', repoName);
    await page.fill('#repo-desc', 'Repo for PR testing');
    await page.click('button[type="submit"]');
    await page.waitForURL(new RegExp(`/${repoName}`));
    await page.close();
  });

  test('PR list page loads correctly', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/pulls`);

    // The PR list should have a heading
    const heading = page.locator('h2').or(page.locator('h1'));
    await expect(heading).toBeVisible();
  });

  test('new PR page is accessible', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/pulls/new`);

    // Should show the PR creation page
    await expect(page).toHaveURL(/\/pulls\/new/);
  });
});

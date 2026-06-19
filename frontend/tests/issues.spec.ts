import { test, expect } from '@playwright/test';

test.describe('Issues', () => {
  const repoName = `test-issues-${Date.now()}`;
  const issueTitle = `Bug: login fails on ${Date.now()}`;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/repo/create');
    await page.fill('#repo-name', repoName);
    await page.fill('#repo-desc', 'Repo for issue testing');
    await page.click('button[type="submit"]');
    await page.waitForURL(new RegExp(`/${repoName}`));
    await page.close();
  });

  test('create a new issue', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/issues/new`);

    await page.fill('#title', issueTitle);
    await page.fill('#body', 'Steps to reproduce:\n1. Click login\n2. See error');

    await page.click('button:has-text("Submit Issue")');

    await page.waitForURL(/\/issues\/\d+/);
    const h2 = page.locator('h2');
    await expect(h2).toContainText(issueTitle);
  });

  test('issue appears in issue list', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/issues`);

    await expect(page.locator(`text="${issueTitle}"`)).toBeVisible();
  });

  test('issue detail shows correct data', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/issues`);

    await page.click(`text="${issueTitle}"`);
    await page.waitForURL(/\/issues\/\d+/);

    await expect(page.locator('h2')).toContainText(issueTitle);
    await expect(page.locator('text=Steps to reproduce:')).toBeVisible();
  });
});

import { test, expect, type Page } from '@playwright/test';

// A fixed name, because this suite is not allowed to depend on a fresh instance.
// Playwright discards the worker after a failing test and re-imports this file
// in the next one, so a name built from Date.now() changed on every re-import:
// beforeAll then created a second repository, and the tests that follow looked
// for the issue in a repository that never received it. The cascade made two
// unrelated tests fail behind one real selector mistake.
const REPO = 'test-issues';

// The heading is an h1 next to the repository breadcrumb, and the number is
// rendered inside it, so an exact text match on the title alone does not hold.
const issueHeading = (page: Page, title: string) =>
  page.getByRole('heading', { level: 1 }).filter({ hasText: title });

test.describe('Issues', () => {
  const issueTitle = `Bug: login fails on ${Date.now()}`;

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/repo/create');
    await page.fill('#repo-name', REPO);
    await page.fill('#repo-desc', 'Repo for issue testing');
    await page.click('button[type="submit"]');
    await page.waitForURL(new RegExp(`/${REPO}`));
    await page.close();
  });

  test('create a new issue', async ({ page }) => {
    await page.goto(`/ggit-admin/${REPO}/issues/new`);

    await page.fill('#title', issueTitle);
    await page.fill('#body', 'Steps to reproduce:\n1. Click login\n2. See error');

    // Watching the response is what makes this a real check. The bug that
    // started this file answered 200 with number: 0, the browser navigated to
    // /issues/undefined, and a redirect assertion alone would have missed it.
    const created = page.waitForResponse(
      (r) => r.url().includes('/issues') && r.request().method() === 'POST'
    );
    await page.click('button:has-text("Submit Issue")');

    const response = await created;
    expect(response.status()).toBe(201);
    expect((await response.json()).number).toBeGreaterThan(0);

    await page.waitForURL(/\/issues\/\d+/);
    await expect(issueHeading(page, issueTitle)).toBeVisible();
  });

  test('issue appears in issue list', async ({ page }) => {
    await page.goto(`/ggit-admin/${REPO}/issues`);

    await expect(page.getByRole('link', { name: issueTitle })).toBeVisible();
  });

  test('issue detail shows correct data', async ({ page }) => {
    await page.goto(`/ggit-admin/${REPO}/issues`);

    await page.getByRole('link', { name: issueTitle }).click();
    await page.waitForURL(/\/issues\/\d+/);

    await expect(issueHeading(page, issueTitle)).toBeVisible();
    await expect(page.getByText('Steps to reproduce:')).toBeVisible();
  });
});
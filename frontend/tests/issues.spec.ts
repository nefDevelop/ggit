import { test, expect, type Page, type APIRequestContext } from '@playwright/test';

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

// createRepoIfMissing keeps the suite re-runnable. "make test-e2e" starts from a
// fresh database, but running Playwright twice against the same environment used
// to fail in beforeAll on a name clash, which reads like a product failure and is
// not one.
//
// With fresh, an existing repository is deleted first. The issue specs need that:
// they reuse one repository, and after a few passes it holds a dozen issues, the
// list paginates and the one just created is no longer on the first page.
async function createRepoIfMissing(
  page: Page,
  request: APIRequestContext,
  name: string,
  description: string,
  fresh = false,
) {
  // /view is the repository read endpoint. A bare /api/v2/:owner/:repo does not
  // exist and answers 404 for everything, which made this check always miss and
  // then try to create a name that was already taken.
  const existing = await request.get(`/api/v2/ggit-admin/${name}/view`);

  if (existing.status() === 200) {
    if (!fresh) return;
    await request.delete(`/api/v2/ggit-admin/${name}/settings`, { data: { name } });
  }

  await page.goto('/repo/create');
  await page.fill('#repo-name', name);
  await page.fill('#repo-desc', description);
  await page.click('button[type="submit"]');
  await page.waitForURL(new RegExp(`/${name}`));
}

test.describe('Issues', () => {
  const issueTitle = `Bug: login fails on ${Date.now()}`;

  test.beforeAll(async ({ browser, request }) => {
    const page = await browser.newPage();
    await createRepoIfMissing(page, request, REPO, 'Repo for issue testing', true);
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
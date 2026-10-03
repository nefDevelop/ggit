import { test, expect, type Page, type APIRequestContext } from '@playwright/test';

// A fixed name for the same reason as the other specs: a name built from
// Date.now() changes when Playwright re-imports the file in a new worker, and
// beforeAll then creates a second repository.
const REPO = 'test-prs';

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

test.describe('Pull Requests', () => {
  const repoName = REPO;

  test.beforeAll(async ({ browser, request }) => {
    const page = await browser.newPage();
    await createRepoIfMissing(page, request, REPO, 'Repo for PR testing');
    await page.close();
  });

  test('PR list page loads correctly', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/pulls`);

    // "h2 or h1" matched three elements on this page once the repository header
    // and the empty state both grew a heading, so the title is addressed by role.
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible();
  });

  test('new PR page is accessible', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/pulls/new`);

    // Should show the PR creation page
    await expect(page).toHaveURL(/\/pulls\/new/);
  });
});

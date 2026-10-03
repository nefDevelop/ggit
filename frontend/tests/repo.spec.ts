import { test, expect, type Page, type APIRequestContext } from '@playwright/test';

// A fixed name on purpose. Playwright drops the worker after a failing test and
// re-imports this file in the next one, so a name built from Date.now() changed
// on every re-import: the second test then visited a repository that did not
// exist and reported a missing clone button.
const REPO = 'test-repo-manager';

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

test.describe('Repository Management', () => {
  test.beforeAll(async ({ browser, request }) => {
    const page = await browser.newPage();
    await createRepoIfMissing(page, request, REPO, 'Este es un repo de prueba automatizada');
    await page.close();
  });

  test('should create a new repository', async ({ page }) => {
    await page.goto(`/ggit-admin/${REPO}`);

    // The repository page renders the name twice: once in the breadcrumb and
    // once as the title, so a bare h1 matches two elements.
    await expect(
      page.getByRole('heading', { level: 1, name: REPO }).last()
    ).toBeVisible();
  });

  test('should display correct clone URLs', async ({ page }) => {
    await page.goto(`/ggit-admin/${REPO}`);

    const cloneButton = page.getByRole('button', { name: /Clone|Clonar/i });
    await expect(cloneButton).toBeVisible();

    await cloneButton.click();
    const urlDisplay = page.locator('.url-display');
    await expect(urlDisplay).toBeVisible();
    const url = await urlDisplay.textContent();
    expect(url).not.toBeNull();
    expect(url!).toContain('.git');
    expect(url!).toContain(REPO);
  });
});
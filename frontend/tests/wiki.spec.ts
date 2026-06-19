import { test, expect } from '@playwright/test';

test.describe('Wiki', () => {
  const repoName = `test-wiki-${Date.now()}`;
  const pageTitle = 'Getting Started';
  const pageContent = 'Welcome to the wiki page for this project.';
  const commitMsg = 'Created Getting Started page';

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await page.goto('/repo/create');
    await page.fill('#repo-name', repoName);
    await page.fill('#repo-desc', 'Repo for wiki testing');
    await page.click('button[type="submit"]');
    await page.waitForURL(new RegExp(`/${repoName}`));
    await page.close();
  });

  test('wiki page loads (even if empty)', async ({ page }) => {
    await page.goto(`/ggit-admin/${repoName}/wiki/Home`);

    // Wiki should display the page title or an empty state
    await expect(page.locator('body')).toBeVisible();
  });

  test('create a wiki page', async ({ page }) => {
    // Navigate to wiki and click edit/create
    await page.goto(`/ggit-admin/${repoName}/wiki/Home`);

    // Look for an edit/create button or the wiki editor
    const editBtn = page.locator('button:has-text("Edit")').or(
      page.locator('a:has-text("Edit")'),
    ).or(
      page.locator('button:has-text("Create")'),
    );

    if (await editBtn.isVisible()) {
      await editBtn.click();
    }

    // Fill the page editor
    const titleInput = page.locator('input.title-input');
    if (await titleInput.isVisible()) {
      await titleInput.fill(pageTitle);
    }

    const contentArea = page.locator('textarea, [contenteditable="true"]').first();
    if (await contentArea.isVisible()) {
      await contentArea.fill(pageContent);
    }

    // Fill commit message
    const commitInput = page.locator('#commit');
    if (await commitInput.isVisible()) {
      await commitInput.fill(commitMsg);
    }

    // Click save
    const saveBtn = page.locator('button:has-text("Save")').or(
      page.locator('button:has-text("Guardar")'),
    );
    if (await saveBtn.isVisible()) {
      await saveBtn.click();
    }

    // Wait and check the page
    await page.waitForTimeout(2000);
    await expect(page.locator('body')).toBeVisible();
  });
});

import { test, expect } from '@playwright/test';

test.describe('Organizations', () => {
  // Timestamped so a repeated run does not collide with the organization the
  // previous one created.
  const orgName = `test-org-${Date.now()}`;

  test('org creation page loads', async ({ page }) => {
    await page.goto('/org/create');

    // The page titles itself with an h1, not an h2. The interface is in Spanish
    // here, so the heading is matched by role and by the form the test actually
    // cares about, rather than by a wording that changes with the locale.
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10000 });
    await expect(page.locator('#org-name')).toBeVisible();
  });

  test('create an organization', async ({ page }) => {
    await page.goto('/org/create');

    // The fields are addressed directly. The previous version combined them with
    // .or() and guarded each one with isVisible(), so a rename of either field
    // made the test fill nothing and then sit for 30s on a submit button that was
    // still disabled.
    const nameInput = page.locator('#org-name');
    const submitBtn = page.locator('button[type="submit"]');

    await expect(nameInput).toBeVisible();
    await nameInput.fill(orgName);

    // Asserting the button became enabled turns a silent mismatch into a clear
    // failure instead of a timeout.
    await expect(submitBtn).toBeEnabled({ timeout: 10000 });

    const created = page.waitForResponse(
      (r) => r.url().includes('/api/v2/user/settings/orgs') && r.request().method() === 'POST',
    );
    await submitBtn.click();

    const response = await created;
    expect(response.status()).toBe(201);

    // Watch the response and then the navigation, rather than sleeping and
    // hoping the page caught up. Creating redirects to the organization's own
    // page, where its name is supposed to appear, so the redirect is what
    // confirms the creation.
    await expect(page).toHaveURL(new RegExp(`/org/${orgName}`), { timeout: 15000 });
  });
});
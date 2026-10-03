import { test, expect, type Page } from '@playwright/test';

const regularUser = {
  username: 'regular-user',
  password: 'password123',
  email: 'regular@example.com',
};

async function loginAsRegularUser(page: Page) {
  await page.goto('/user/login');
  await page.fill('#username', regularUser.username);
  await page.fill('#password', regularUser.password);
  await page.click('button[type="submit"]');
  await page.waitForURL('/dashboard');
}

test.describe('Admin Security', () => {
  // The previous version of this file asserted two things that were both wrong.
  //
  // It expected the admin routes to redirect a non-admin somewhere else, but the
  // admin panel is rendered on the client: the server answers 200 with an empty
  // shell for everyone and the decision is made by the API, so there is no
  // redirect to observe and the URL check never had a chance to mean anything.
  //
  // And it used the request fixture to assert 403 on /api/v2/admin/stats, but
  // that fixture inherits this project's storage state, which is the admin
  // session. It was asking the admin for a refusal.
  //
  // What matters is that a non-admin gets no administrative data, so that is
  // what is asserted, and the request is made from the page that is logged in as
  // the non-admin so the session is the one under test.

  test('the API refuses a non-admin', async ({ page }) => {
    await loginAsRegularUser(page);

    for (const endpoint of [
      '/api/v2/admin/stats',
      '/api/v2/admin/users',
      '/api/v2/admin/repos',
    ]) {
      const response = await page.request.get(endpoint);
      expect(response.status(), `${endpoint} must not answer a non-admin`).toBe(403);
    }
  });

  test('the admin pages leak no data to a non-admin', async ({ page }) => {
    await loginAsRegularUser(page);

    const adminRoutes = [
      '/admin',
      '/admin/users',
      '/admin/repos',
      '/admin/auths',
      '/admin/notices',
      '/admin/monitor',
      '/admin/repos/adoption',
    ];

    for (const route of adminRoutes) {
      // "commit" resolves on the response and does not wait for the navigation to
      // settle. The admin panel redirects on the client, /admin to /admin/users,
      // and that second navigation aborts the first one, which surfaced as
      // ERR_ABORTED rather than as anything to do with permissions.
      const response = await page.goto(route, { waitUntil: 'commit' });
      expect(response?.status(), `${route} answered`).toBe(200);
      await page.waitForTimeout(500);

      // Nothing an administrator would see: no totals, no user rows, no
      // repository rows, and no data blob inlined into the page.
      const body = (await page.locator('body').innerText()).toLowerCase();
      for (const leak of ['total users', 'total repos', 'statistics']) {
        expect(body, `${route} exposed "${leak}"`).not.toContain(leak);
      }

      const inlined = await page.content();
      expect(inlined, `${route} inlined administrative data`).not.toMatch(
        /"is_admin"\s*:\s*true/
      );
    }
  });

  test('an admin still gets the data', async ({ page }) => {
    // The counterpart, so the test above cannot pass by the panel being empty
    // for everybody.
    const response = await page.request.get('/api/v2/admin/stats');
    expect(response.status()).toBe(200);
  });
});
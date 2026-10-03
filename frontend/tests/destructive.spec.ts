import { test, expect } from '@playwright/test';
import { existsSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { dirname, resolve } from 'node:path';

// Repository storage layout is <repo root>/<owner>/<lowercase repo>.git, see
// repox.RepositoryPath. The E2E Makefile points REPO_ROOT_PATH at a throwaway
// directory, so asserting on disk here cannot touch real data.
const specDir = dirname(fileURLToPath(import.meta.url));
const repoRootOnDisk = resolve(specDir, '../../ggit-repositories-test');
const diskPath = (owner: string, repo: string) =>
  resolve(repoRootOnDisk, owner, `${repo.toLowerCase()}.git`);

const owner = 'ggit-admin';

test.describe('Destructive paths', () => {
  test('deleting a repository actually removes it', async ({ page, request }) => {
    const repoName = `destructive-${Date.now()}`;

    await page.goto('/repo/create');
    await page.fill('#repo-name', repoName);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(new RegExp(`/${repoName}`));

    // Sanity check. Without it a failure further down could just mean the repo
    // was never created, and the deletion assertions would prove nothing.
    expect(existsSync(diskPath(owner, repoName))).toBe(true);

    await page.goto(`/${owner}/${repoName}/settings`);

    const confirmInput = page.getByPlaceholder(
      /Escribe el nombre del repositorio|Enter the repository name/i
    );
    const deleteButton = page.getByRole('button', {
      name: /Eliminar Repositorio|Delete Repository/i,
    });

    // The button stays disabled until the typed name matches exactly, so a
    // stray click cannot destroy the repository.
    await expect(deleteButton).toBeDisabled();
    await confirmInput.fill('not-the-right-name');
    await expect(deleteButton).toBeDisabled();

    // Still intact: the guard held rather than the request failing later.
    expect((await request.get(`/api/v2/${owner}/${repoName}/view`)).status()).toBe(200);

    await confirmInput.fill(repoName);
    await expect(deleteButton).toBeEnabled();

    // The click is watched rather than assumed. This is the assertion that
    // matters: the original bug had the UI report success while the API
    // answered 400, so the status of this very request is the evidence.
    const deleteResponse = page.waitForResponse(
      (r) => r.url().includes('/settings') && r.request().method() === 'DELETE'
    );
    await deleteButton.click();

    // The handler asks for confirmation through the global dialog before it
    // touches the API, so without this the request is never sent.
    const confirmButton = page.getByRole('button', { name: /^(Aceptar|Accept|Confirmar)$/ });
    await expect(confirmButton).toBeVisible();
    await confirmButton.click();

    const deletion = await deleteResponse;
    expect(
      deletion.status(),
      `DELETE ${deletion.url()} devolvio ${await deletion.text()}`
    ).toBeLessThan(300);

    // The original bug sent no body, the API answered 400 and the repository
    // survived while the UI moved on. Checking the redirect alone would not
    // have caught it, so the assertions below are on the API and on disk.
    expect((await request.get(`/api/v2/${owner}/${repoName}/view`)).status()).toBe(404);
    expect(existsSync(diskPath(owner, repoName))).toBe(false);

    await page.goto(`/${owner}/${repoName}`);
    await expect(page.locator('body')).not.toContainText(repoName);

    // It must be gone from the owner's list too, not merely unreachable by URL.
    await page.goto('/dashboard');
    await expect(page.locator('body')).not.toContainText(repoName);
  });

  test('deleting a webhook removes it', async ({ page, request }) => {
    const repoName = `webhook-target-${Date.now()}`;

    await page.goto('/repo/create');
    await page.fill('#repo-name', repoName);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(new RegExp(`/${repoName}`));

    // Created through the API because the point of this test is the deletion;
    // driving the webhook form would only add ways for it to fail for unrelated
    // reasons.
    const created = await request.post(`/api/v2/${owner}/${repoName}/settings/hooks`, {
      data: {
        url: 'http://127.0.0.1:9/never-called',
        content_type: 'json',
        events: ['push'],
        active: true,
      },
    });
    expect(created.status()).toBe(201);
    const hookId = (await created.json()).id;
    expect(hookId).toBeGreaterThan(0);

    const before = await request.get(`/api/v2/${owner}/${repoName}/settings/hooks`);
    expect(before.status()).toBe(200);
    expect(await before.json()).toHaveLength(1);

    expect((await request.delete(`/api/v2/${owner}/${repoName}/settings/hooks/${hookId}`)).status()).toBe(200);

    const after = await request.get(`/api/v2/${owner}/${repoName}/settings/hooks`);
    expect(after.status()).toBe(200);
    expect(await after.json()).toHaveLength(0);
  });

  test('revoking an access token removes it', async ({ request }) => {
    // The token is created here instead of relying on one already existing, so
    // the test does not depend on the order the suite happens to run in.
    const name = `e2e-token-${Date.now()}`;
    const created = await request.post('/api/v2/user/settings/tokens', {
      data: { name },
    });
    expect(created.status()).toBe(200);

    const listed = await request.get('/api/v2/user/settings/tokens');
    expect(listed.status()).toBe(200);
    const tokens = await listed.json();
    const target = tokens.find((t: { name: string }) => t.name === name);
    expect(target).toBeDefined();

    expect((await request.delete(`/api/v2/user/settings/tokens/${target.id}`)).status()).toBe(200);

    const remaining = await request.get('/api/v2/user/settings/tokens');
    expect((await remaining.json()).some((t: { name: string }) => t.name === name)).toBe(false);
  });
});
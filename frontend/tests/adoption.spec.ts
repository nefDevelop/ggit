import { test, expect } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';
import { fileURLToPath } from 'url';

test.describe('Repository Adoption', () => {
  // The E2E Makefile redirects the repository root under data/test-e2e/, and the
  // fixture has to be created where the server actually looks. It used to point
  // at ../ggit-repositories-test, which is no longer read by anything.
  const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../data/test-e2e/repositories');
  const ownerName = 'ggit-admin';
  const repoName = 'test-adopt';
  const fullPath = path.join(repoRoot, ownerName, `${repoName}.git`);

  test.beforeAll(async ({ request }) => {
    // Adoption moves the repository into the managed tree, so after a first run
    // "test-adopt" exists as a managed repository and no longer shows up in the
    // adoption list. Deleting it first is what makes this spec re-runnable.
    await request.delete(`/api/v2/${ownerName}/${repoName}/settings`, {
      data: { name: repoName },
    });
    fs.rmSync(fullPath, { recursive: true, force: true });

    fs.mkdirSync(fullPath, { recursive: true });
    // Crear un archivo HEAD mínimo para que parezca un repo git bare
    fs.writeFileSync(path.join(fullPath, 'HEAD'), 'ref: refs/heads/master');
  });

  test.beforeEach(async ({ page }) => {
    // Login como administrador
    await page.goto('/user/login');
    await page.fill('#username', 'ggit-admin');
    await page.fill('#password', 'admin1234');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL('/dashboard');
  });

  test('debe listar y adoptar un repositorio no gestionado', async ({ page }) => {
    // Ir a la página de administración de adopción
    await page.goto('/admin/repos/adoption');

    // Verificar que el repositorio aparece en la lista
    await expect(page.getByText(repoName, { exact: true })).toBeVisible();

    // Adoptar pide confirmacion en el dialogo global antes de llamar a la API,
    // asi que hay dos clics: el de la fila y el del "Aceptar". El test solo hacia
    // el primero y se quedaba esperando un aviso que nunca llegaba.
    const adoptButton = page.locator(`button:has-text("Adoptar"), button:has-text("Adopt")`).first();
    await adoptButton.click();
    await page.getByRole('button', { name: /^(Aceptar|Accept|Confirmar)$/ }).click();

    // The confirmation is a dialog.alert whose text includes the owner and the
    // repository name, so it is matched on the stable part of the sentence rather
    // than on a wording that was never what the component rendered.
    await expect(page.getByText(/adoptado con éxito/i)).toBeVisible();
    
    // Verificar que el repositorio ya no aparece en la lista de adopción
    await page.goto('/admin/repos/adoption');
    await expect(page.getByText(repoName)).not.toBeVisible();
  });

  test('Malintención: Recharzar nombres de repositorio maliciosos vía API', async ({ page }) => {
    // Intentar adoptar un repositorio con path traversal vía API directa
    const response = await page.request.post('/api/v2/admin/repos/adopt', {
      data: {
        owner_name: ownerName,
        repo_name: '../../../../etc/passwd'
      }
    });

    expect(response.status()).toBe(400);
    const body = await response.json();
    expect(body.message).toContain('malicious');
  });
});

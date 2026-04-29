import { test, expect } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';

test.describe('Repository Adoption', () => {
  const repoRoot = path.resolve('../ggit-repositories-test');
  const ownerName = 'ggit-admin';
  const repoName = 'test-adopt';
  const fullPath = path.join(repoRoot, ownerName, `${repoName}.git`);

  test.beforeAll(async () => {
    // Crear carpeta de repositorio falso en disco
    if (!fs.existsSync(fullPath)) {
      fs.mkdirSync(fullPath, { recursive: true });
      // Crear un archivo HEAD mínimo para que parezca un repo git bare
      fs.writeFileSync(path.join(fullPath, 'HEAD'), 'ref: refs/heads/master');
    }
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
    // Manejar el diálogo de confirmación
    page.on('dialog', dialog => dialog.accept());

    // Ir a la página de administración de adopción
    await page.goto('/admin/repos/adoption');

    // Verificar que el repositorio aparece en la lista
    await expect(page.getByText(repoName, { exact: true })).toBeVisible();

    // Hacer clic en el botón de adoptar
    const adoptButton = page.locator(`button:has-text("Adoptar"), button:has-text("Adopt")`).first();
    await adoptButton.click();

    // Verificar éxito (debería redirigir o mostrar un mensaje)
    await expect(page.getByText('Repositorio adoptado con éxito').or(page.getByText('Repository adopted successfully'))).toBeVisible();
    
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

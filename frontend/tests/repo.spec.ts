import { test, expect } from '@playwright/test';

test.describe('Repository Management', () => {
  const repoName = `test-repo-${Date.now()}`;

  test('should create a new repository', async ({ page }) => {
    await page.goto('/repo/create');

    // Rellenar formulario de creación
    await page.fill('input[name="name"]', repoName);
    await page.fill('textarea[name="description"]', 'Este es un repo de prueba automatizada');
    
    // Hacer clic en crear
    await page.click('button[type="submit"]');

    // Verificar que nos redirige a la página del repo
    await expect(page).toHaveURL(new RegExp(`/${repoName}`));
    await expect(page.locator('h1')).toContainText(repoName);
  });

  test('should display correct clone URLs', async ({ page }) => {
    // Vamos al repo que acabamos de crear (o uno conocido)
    await page.goto(`/admin/${repoName}`); // Asumiendo que el dueño es admin

    // Verificar que el botón de clonar está presente
    const cloneButton = page.locator('button:has-text("Clone")');
    await expect(cloneButton).toBeVisible();

    // Verificar que la URL de clonado HTTP es correcta
    const cloneInput = page.locator('input[readonly]');
    const url = await cloneInput.inputValue();
    expect(url).toContain('.git');
    expect(url).toContain(repoName);
  });
});

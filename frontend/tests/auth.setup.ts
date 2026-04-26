import { test as setup, expect } from '@playwright/test';

const authFile = 'tests/.auth/user.json';

setup('authenticate', async ({ page }) => {
  await page.goto('/user/login');
  // AQUÍ DEBERÍAS PONER TUS CREDENCIALES DE TEST
  await page.fill('#username', 'admin'); 
  await page.fill('#password', 'admin1234');
  await page.click('button[type="submit"]');

  // Esperar a que lleguemos al dashboard
  await page.waitForURL('/dashboard');
  await expect(page.locator('h2')).toContainText('Dashboard');

  // Guardar el estado de la sesión
  await page.context().storageState({ path: authFile });
});

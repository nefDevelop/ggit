import { test as setup, expect } from '@playwright/test';

const authFile = 'tests/.auth/user.json';

setup('authenticate', async ({ page }) => {
  await page.goto('/user/login');
  // AQUÍ DEBERÍAS PONER TUS CREDENCIALES DE TEST
  await page.fill('#username', 'ggit-admin'); 
  await page.fill('#password', 'admin1234');
  await page.click('button[type="submit"]');

  // Esperar a que lleguemos al dashboard
  await page.waitForURL('/dashboard');
  
  // Verificar que estamos en el dashboard buscando el botón de crear repositorio (idioma flexible)
  await expect(page.getByText('Nuevo Repositorio').or(page.getByText('New Repository'))).toBeVisible();

  // Guardar el estado de la sesión
  await page.context().storageState({ path: authFile });
});

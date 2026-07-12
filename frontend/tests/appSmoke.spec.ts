import { expect, test } from '@playwright/test';

test('opens lobby, room, and match shell', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Crownfall' })).toBeVisible();
  await page.getByRole('link', { name: 'Enter room' }).click();
  await page.getByRole('link', { name: 'Open match table' }).click();
  await expect(page.getByRole('heading', { name: 'Council Table' })).toBeVisible();
  await expect(page.getByLabel(/Animated Crownfall table/)).toBeVisible();
  await expect(page.getByTestId('pixi-board').locator('canvas')).toHaveCount(1);
  await page.goto('/');
  await expect(page.getByTestId('pixi-board')).toHaveCount(0);
});

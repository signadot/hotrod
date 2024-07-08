import { test, expect } from '@playwright/test';

test('request ride and check routing context', async ({ page }) => {
    await page.getByRole('combobox').first().selectOption('1');
    await page.getByRole('combobox').nth(1).selectOption('123');
    await page.getByRole('button', { name: 'Request Ride' }).click();
    await page.locator('div').filter({ hasText: /^17:39:04\.995frontend\(baseline\)Processing dispatch driver request$/ }).getByRole('paragraph').nth(2).click();
    await page.locator('div').filter({ hasText: /^17:39:04\.989browser\(baseline\)Requesting a ride\.$/ }).getByRole('paragraph').nth(2).click();
    await page.locator('div').filter({ hasText: /^17:39:04\.995frontend\(baseline\)Processing dispatch driver request$/ }).getByRole('paragraph').nth(2).click();
    await page.locator('div').filter({ hasText: /^17:39:05\.311location\(baseline\)Resolving locations$/ }).getByRole('paragraph').nth(2).click();
    await page.locator('div').filter({ hasText: /^17:39:05\.734driver\(baseline\)Finding an available driver$/ }).getByRole('paragraph').nth(2).click();
});
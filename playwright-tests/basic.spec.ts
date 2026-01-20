import { test, expect } from "@playwright/test";

test.describe("Request a ride", () => {
	test("check route resolve routes", async ({ page }) => {
		await page.goto("/");
		await page.waitForLoadState();

		await page.getByRole("combobox").first().selectOption("1");
		await page.getByRole("combobox").nth(1).selectOption("123");
		await page.getByRole("button", { name: "Request Ride" }).click();

		// Open the logs drawer to see the route service log
		await page.getByRole("button", { name: "Show Logs" }).click();

		// Wait for "Resolving routes" text to appear in the logs
		await expect(page.getByText("Resolving routes")).toBeVisible({ timeout: 15000 });
	});

	test("check driver series", async ({ page }) => {
		await page.goto("/");
		await page.waitForLoadState();

		await page.getByRole("combobox").first().selectOption("1");
		await page.getByRole("combobox").nth(1).selectOption("123");
		await page.getByRole("button", { name: "Request Ride" }).click();

		await expect(
			page.locator("//div[p[2][contains(text(), 'driver')]]/p[4]").last(),
		).toHaveText(/.*T7\d{5}C.*/, { timeout: 20000 });
	});

	test("check distance display with ETA", async ({ page }) => {
		await page.goto("/");
		await page.waitForLoadState();

		await page.getByRole("combobox").first().selectOption("1");
		await page.getByRole("combobox").nth(1).selectOption("123");
		await page.getByRole("button", { name: "Request Ride" }).click();

		// Wait for the driver countdown text to appear
		// Pattern: "The driver [Name] ([Plate]) will arrive in MM:SS • XX.YYY km away"
		const countdownText = page.getByText(/will arrive in \d{2}:\d{2}/);

		// Wait for element to be visible with timeout
		await expect(countdownText).toBeVisible({ timeout: 20000 });

		// Check that the text contains both ETA and distance
		await expect(countdownText).toContainText(
			/will arrive in \d{2}:\d{2} • \d{1,2}\.\d{1,3} km away/,
		);
	});
});

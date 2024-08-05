import { test, expect } from "@playwright/test";

test.describe("Request a ride", () => {
	test("check route resolve routes", async ({ page }) => {
		await page.goto("/");
		await page.waitForLoadState();

		await page.getByRole("combobox").first().selectOption("1");
		await page.getByRole("combobox").nth(1).selectOption("123");
		await page.getByRole("button", { name: "Request Ride" }).click();

		await expect(
			page.locator("//div[p[2][contains(text(), 'route')]]/p[4]"),
		).toHaveText("Resolving routes")
	});

	test("check driver series", async ({ page }) => {
		await page.goto("/");
		await page.waitForLoadState();

		await page.getByRole("combobox").first().selectOption("1");
		await page.getByRole("combobox").nth(1).selectOption("123");
		await page.getByRole("button", { name: "Request Ride" }).click();

		await expect(
			page.locator("//div[p[2][contains(text(), 'driver')]]/p[4]").last(),
		).toHaveText(/.*T7\d{5}C.*/);
	});
});

import { test, expect } from "@playwright/test";

test("request ride and check routing context", async ({ page }) => {
	const sandboxName = process.env.SIGNADOT_SANDBOX_NAME;

	await page.reload();
	await page.goto("/");
	await page.waitForLoadState();
	await page.reload();


	await page.getByRole("combobox").first().selectOption("1");
	await page.getByRole("combobox").nth(1).selectOption("123");
	await page.getByRole("button", { name: "Request Ride" }).click();

	const getSandboxName = (
		service: "browser" | "frontend" | "location" | "driver" | "route",
	): RegExp => {
		if (!sandboxName) return /baseline/;

		const serviceValue = process.env["SANDBOXED_" + service.toUpperCase()];

		console.log({ serviceValue });
		if (serviceValue !== "1") {
			return /baseline/;
		}

		return new RegExp(`${sandboxName}`);
	};

	await expect(
		page.locator('//*[@id="accordion-panel-:r0:"]/div/div[1]'),
	).toHaveText(getSandboxName("browser"));
	await expect(
		page.locator('//*[@id="accordion-panel-:r0:"]/div/div[2]'),
	).toHaveText(getSandboxName("frontend"));
	await expect(
		page.locator('//*[@id="accordion-panel-:r0:"]/div/div[3]'),
	).toHaveText(getSandboxName("location"));
	await expect(
		page.locator('//*[@id="accordion-panel-:r0:"]/div/div[4]'),
	).toHaveText(getSandboxName("driver"));
	await expect(
		page.locator('//*[@id="accordion-panel-:r0:"]/div/div[5]'),
	).toHaveText(getSandboxName("route"));
});

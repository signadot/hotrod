import { defineConfig, devices } from "@playwright/test";

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// import dotenv from 'dotenv';
// dotenv.config({ path: path.resolve(__dirname, '.env') });

/**
 * See https://playwright.dev/docs/test-configuration.
 */

export const BASE_URL = `http://frontend.${process.env.HOTROD_NAMESPACE}:8080`;

console.log({ ROUTING: process.env.SIGNADOT_ROUTING_KEY });

export default defineConfig({
	testDir: "./playwright-tests",
	/* Run tests in files in parallel */
	fullyParallel: true,
	/* Fail the build on CI if you accidentally left test.only in the source code. */
	forbidOnly: !!process.env.CI,
	/* Retry on CI only */
	retries: process.env.CI ? 0 : 0,
	/* Opt out of parallel tests on CI. */
	workers: process.env.CI ? 1 : undefined,
	/* Reporter to use. See https://playwright.dev/docs/test-reporters */
	reporter: [["html", { open: "never" }]],
	/* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
	use: {
		/* Base URL to use in actions like `await page.goto('/')`. */
		baseURL: BASE_URL,

		/* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
		trace: "on",

		video: "on",
		extraHTTPHeaders: {
			baggage: `sd-routing-key=${process.env.SIGNADOT_ROUTING_KEY}`,
		},
	},

	/* Configure projects for major browsers */
	projects: [
		{
			name: "chromium",
			use: {
				...devices["Desktop Chrome"],
				baseURL: BASE_URL,
				// extraHTTPHeaders: {
				// 	baggage: `sd-routing-key=${process.env.SIGNADOT_ROUTING_KEY}`,
				// },
			},
		},
	],
});

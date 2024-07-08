import { defineConfig, devices } from "@playwright/test";

const BASE_URL = `http://frontend.${process.env.HOTROD_NAMESPACE}:8080`;

console.log({ ROUTING: process.env.SIGNADOT_ROUTING_KEY });

export default defineConfig({
	testDir: "./playwright-tests",
	fullyParallel: true,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 0 : 0,
	workers: process.env.CI ? 1 : undefined,
	reporter: [["html", { open: "never" }]],
	use: {
		baseURL: BASE_URL,
		trace: "on",
		video: "on",
		extraHTTPHeaders: {
			baggage: `sd-routing-key=${process.env.SIGNADOT_ROUTING_KEY}`,
		},
	},
	projects: [
		{
			name: "chromium",
			use: {
				...devices["Desktop Chrome"],
				baseURL: BASE_URL,
			},
		},
	],
});

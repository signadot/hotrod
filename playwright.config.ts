import { defineConfig, devices } from "@playwright/test";

const BASE_URL = `http://frontend.${process.env.HOTROD_NAMESPACE}:8080`;

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

signadot job submit \
  -f .signadot/testing/e2e-playwright-job.yaml \
  --config signadot-prod-config.yaml \
  --set branch=use-playwright-e2e \
  --set namespace="hotrod-istio" \
  --set sandbox=pr-243-driver \
  --set sandboxed_frontend="" \
  --set sandboxed_location="" \
  --set sandboxed_driver="1" \
  --set sandboxed_route="" \
  --attach
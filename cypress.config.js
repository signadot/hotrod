const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    video: true,
    experimentalStudio: true,
    viewportWidth: 1024,
    viewportHeight: 780,
    env: {
      HOTROD_NAMESPACE: 'hotrod',
      SIGNADOT_ROUTING_KEY: '',
      FRONTEND_SANDBOX_NAME: '',
      LOCATION_SANDBOX_NAME: '',
      ROUTE_SANDBOX_NAME: '',
      DRIVER_SANDBOX_NAME: '',
    },
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});

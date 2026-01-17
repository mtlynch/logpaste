// @ts-check
const { devices } = require("@playwright/test");

const config = {
  testDir: "./e2e",
  timeout: 15 * 1000,
  expect: {
    timeout: 5 * 1000,
  },
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  workers: 1,
  reporter: "html",
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL || "http://localhost:3333",
    actionTimeout: 0,
    trace: "on",
    video: "on",
  },
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
      },
    },
  ],
  outputDir: "e2e-results/",
};

module.exports = config;

const { test, expect } = require("@playwright/test");

test("pastes an entry", async ({ page }) => {
  await page.goto("/");

  await page.fill("#upload-textarea", "test upload data");
  await page.click("#upload");

  const resultLink = page.locator("#result a");
  await expect(resultLink).toHaveText(/.+/);

  const href = await resultLink.getAttribute("href");
  expect(href).not.toBeNull();

  const resolvedUrl = new URL(href, page.url()).toString();
  const response = await page.request.get(resolvedUrl);
  expect(response.ok()).toBe(true);
  expect(await response.text()).toBe("test upload data");
});

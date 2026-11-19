import { test, expect } from '@playwright/test';

// --- Mobile (iPhone SE / small phone) ---
test.describe('Mobile (375x667)', () => {
  test.use({ viewport: { width: 375, height: 667 } });

  test('nav bar spans full width', async ({ page }) => {
    await page.goto('/');
    const nav = page.locator('.nav');
    const navBox = await nav.boundingBox();
    expect(navBox).toBeTruthy();
    expect(navBox!.width).toBeGreaterThanOrEqual(375);
  });

  test('nav brand text does not overflow viewport', async ({ page }) => {
    await page.goto('/');
    const brand = page.locator('.nav-brand');
    const brandBox = await brand.boundingBox();
    expect(brandBox).toBeTruthy();
    expect(brandBox!.x + brandBox!.width).toBeLessThanOrEqual(375);
  });

  test('hamburger menu is visible', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('.nav-toggle');
    await expect(toggle).toBeVisible();
  });

  test('nav links are hidden by default', async ({ page }) => {
    await page.goto('/');
    const links = page.locator('.nav-links');
    await expect(links).not.toBeVisible();
  });

  test('hamburger toggles nav links', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('.nav-toggle');
    const links = page.locator('.nav-links');

    await toggle.click();
    await expect(links).toBeVisible();

    await toggle.click();
    await expect(links).not.toBeVisible();
  });

  test('nav link click closes menu and navigates', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('.nav-toggle');
    await toggle.click();

    const docsLink = page.locator('.nav-links a:has-text("Docs")');
    await docsLink.click();

    await expect(page.locator('h1')).toContainText('Documentation');
    const links = page.locator('.nav-links');
    await expect(links).not.toBeVisible();
  });

  test('hero h1 does not overflow viewport', async ({ page }) => {
    await page.goto('/');
    const h1 = page.locator('.hero h1');
    const h1Box = await h1.boundingBox();
    expect(h1Box).toBeTruthy();
    expect(h1Box!.x + h1Box!.width).toBeLessThanOrEqual(375);
  });

  test('hero h1 font size is scaled down', async ({ page }) => {
    await page.goto('/');
    const fontSize = await page.locator('.hero h1').evaluate(
      el => parseFloat(getComputedStyle(el).fontSize)
    );
    expect(fontSize).toBeLessThanOrEqual(24);
  });

  test('features grid is single column', async ({ page }) => {
    await page.goto('/');
    const cards = page.locator('.feature-card');
    const count = await cards.count();
    if (count >= 2) {
      const first = await cards.nth(0).boundingBox();
      const second = await cards.nth(1).boundingBox();
      expect(first).toBeTruthy();
      expect(second).toBeTruthy();
      // Second card should be below first (not side-by-side)
      expect(second!.y).toBeGreaterThan(first!.y);
    }
  });

  test('no horizontal scroll on any page', async ({ page }) => {
    const pages = ['/', '/#/docs', '/#/tutorials', '/#/api'];
    for (const p of pages) {
      await page.goto(p);
      const scrollWidth = await page.evaluate(() => document.documentElement.scrollWidth);
      const clientWidth = await page.evaluate(() => document.documentElement.clientWidth);
      expect(scrollWidth, `Horizontal overflow on ${p}`).toBeLessThanOrEqual(clientWidth + 1);
    }
  });
});

// --- Tablet (iPad) ---
test.describe('Tablet (768x1024)', () => {
  test.use({ viewport: { width: 768, height: 1024 } });

  test('hamburger menu is visible at 768px', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-toggle')).toBeVisible();
  });

  test('hamburger toggles nav links', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('.nav-toggle');
    const links = page.locator('.nav-links');

    await expect(links).not.toBeVisible();
    await toggle.click();
    await expect(links).toBeVisible();
  });

  test('hero h1 does not overflow viewport', async ({ page }) => {
    await page.goto('/');
    const h1 = page.locator('.hero h1');
    const h1Box = await h1.boundingBox();
    expect(h1Box).toBeTruthy();
    expect(h1Box!.x + h1Box!.width).toBeLessThanOrEqual(768);
  });

  test('all pages load without horizontal overflow', async ({ page }) => {
    const pages = ['/', '/#/docs', '/#/tutorials', '/#/api'];
    for (const p of pages) {
      await page.goto(p);
      const scrollWidth = await page.evaluate(() => document.documentElement.scrollWidth);
      const clientWidth = await page.evaluate(() => document.documentElement.clientWidth);
      expect(scrollWidth, `Horizontal overflow on ${p}`).toBeLessThanOrEqual(clientWidth + 1);
    }
  });
});

// --- Desktop ---
test.describe('Desktop (1280x800)', () => {
  test.use({ viewport: { width: 1280, height: 800 } });

  test('hamburger menu is hidden', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-toggle')).not.toBeVisible();
  });

  test('nav links are visible inline', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-links')).toBeVisible();
    const links = page.locator('.nav-links a');
    expect(await links.count()).toBeGreaterThanOrEqual(4);
  });

  test('features grid has multiple columns', async ({ page }) => {
    await page.goto('/');
    const cards = page.locator('.feature-card');
    const count = await cards.count();
    if (count >= 2) {
      const first = await cards.nth(0).boundingBox();
      const second = await cards.nth(1).boundingBox();
      expect(first).toBeTruthy();
      expect(second).toBeTruthy();
      // On desktop, cards should be side-by-side (same Y position)
      expect(Math.abs(second!.y - first!.y)).toBeLessThan(10);
    }
  });
});

// --- Cross-viewport: consistent content ---
test.describe('Cross-viewport consistency', () => {
  for (const { name, width, height } of [
    { name: 'mobile', width: 375, height: 667 },
    { name: 'tablet', width: 768, height: 1024 },
    { name: 'desktop', width: 1280, height: 800 },
  ]) {
    test(`${name}: h1 and footer present on home page`, async ({ page }) => {
      await page.setViewportSize({ width, height });
      await page.goto('/');
      await expect(page.locator('h1')).toContainText('datascience');
      await expect(page.locator('.footer')).toBeVisible();
    });
  }
});

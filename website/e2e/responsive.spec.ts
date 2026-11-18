import { test, expect } from '@playwright/test';

// Viewport sizes for mobile, tablet, and desktop
const viewports = {
  mobile: { width: 375, height: 812 },
  tablet: { width: 768, height: 1024 },
  desktop: { width: 1280, height: 800 },
};

// ----- Mobile viewport tests -----

test.describe('mobile viewport (375x812)', () => {
  test.use({ viewport: viewports.mobile });

  test('page loads and h1 is visible', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1')).toBeVisible();
    await expect(page.locator('h1')).toContainText('datascience');
  });

  test('hamburger menu is visible', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('button.nav-toggle');
    await expect(toggle).toBeVisible();
  });

  test('nav links are hidden by default', async ({ page }) => {
    await page.goto('/');
    const navLinks = page.locator('.nav-links');
    await expect(navLinks).not.toBeVisible();
  });

  test('hamburger toggle opens and closes nav menu', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('button.nav-toggle');
    const navLinks = page.locator('.nav-links');

    // Open menu
    await toggle.click();
    await expect(navLinks).toBeVisible();
    await expect(navLinks).toHaveClass(/open/);

    // Close menu
    await toggle.click();
    await expect(navLinks).not.toBeVisible();
  });

  test('nav links work on mobile', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('button.nav-toggle');

    // Open menu and navigate to Docs
    await toggle.click();
    await page.click('.nav-links a:has-text("Docs")');
    await expect(page.locator('h1')).toContainText('Documentation');

    // Menu should close after navigation
    await expect(page.locator('.nav-links')).not.toBeVisible();
  });

  test('hero section scales properly', async ({ page }) => {
    await page.goto('/');
    const hero = page.locator('.hero');
    await expect(hero).toBeVisible();

    // Hero should fit within viewport width
    const heroBox = await hero.boundingBox();
    expect(heroBox).toBeTruthy();
    expect(heroBox!.width).toBeLessThanOrEqual(viewports.mobile.width);
  });

  test('code blocks do not overflow viewport', async ({ page }) => {
    await page.goto('/');
    const preBlocks = page.locator('pre');
    const count = await preBlocks.count();
    expect(count).toBeGreaterThan(0);

    for (let i = 0; i < count; i++) {
      const box = await preBlocks.nth(i).boundingBox();
      if (box) {
        expect(box.width).toBeLessThanOrEqual(viewports.mobile.width);
      }
    }
  });

  test('features grid stacks to single column', async ({ page }) => {
    await page.goto('/');
    const cards = page.locator('.feature-card');
    const count = await cards.count();
    expect(count).toBeGreaterThan(1);

    // All cards should have roughly the same x position (single column)
    const firstBox = await cards.first().boundingBox();
    expect(firstBox).toBeTruthy();

    for (let i = 1; i < Math.min(count, 3); i++) {
      const box = await cards.nth(i).boundingBox();
      if (box && firstBox) {
        // Same column means same left offset (within rounding tolerance)
        expect(Math.abs(box.x - firstBox.x)).toBeLessThan(2);
      }
    }
  });

  test('footer stacks vertically on mobile', async ({ page }) => {
    await page.goto('/');
    const footerInner = page.locator('.footer-inner');
    await expect(footerInner).toBeVisible();

    // Footer inner should use column layout -- children stack vertically
    const footerBox = await footerInner.boundingBox();
    expect(footerBox).toBeTruthy();
    expect(footerBox!.width).toBeLessThanOrEqual(viewports.mobile.width);
  });

  test('all pages load on mobile', async ({ page }) => {
    const pages = ['/', '/#/docs', '/#/tutorials', '/#/api', '/#/cli'];
    for (const path of pages) {
      await page.goto(path);
      await page.waitForLoadState('networkidle');
      const h1 = await page.$('h1');
      expect(h1, `Page ${path} should have an h1 on mobile`).toBeTruthy();
    }
  });
});

// ----- Tablet viewport tests -----

test.describe('tablet viewport (768x1024)', () => {
  test.use({ viewport: viewports.tablet });

  test('page loads and h1 is visible', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1')).toBeVisible();
    await expect(page.locator('h1')).toContainText('datascience');
  });

  test('hamburger menu is visible at 768px', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('button.nav-toggle');
    // 768px is the breakpoint boundary -- toggle should be visible
    await expect(toggle).toBeVisible();
  });

  test('hamburger toggle opens nav on tablet', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('button.nav-toggle');
    const navLinks = page.locator('.nav-links');

    await toggle.click();
    await expect(navLinks).toBeVisible();
    await expect(navLinks).toHaveClass(/open/);
  });

  test('nav links work on tablet', async ({ page }) => {
    await page.goto('/');
    const toggle = page.locator('button.nav-toggle');

    await toggle.click();
    await page.click('.nav-links a:has-text("Tutorials")');
    await expect(page.locator('h1')).toContainText('Tutorials');
  });

  test('features grid adapts at tablet width', async ({ page }) => {
    await page.goto('/');
    const grid = page.locator('.features-grid');
    await expect(grid).toBeVisible();

    // Grid should be single column at 768px (max-width: 768px breakpoint)
    const cards = page.locator('.feature-card');
    const count = await cards.count();
    expect(count).toBeGreaterThan(1);

    const firstBox = await cards.first().boundingBox();
    const secondBox = await cards.nth(1).boundingBox();
    expect(firstBox).toBeTruthy();
    expect(secondBox).toBeTruthy();

    // Single-column: second card below first (same x, different y)
    expect(Math.abs(secondBox!.x - firstBox!.x)).toBeLessThan(2);
    expect(secondBox!.y).toBeGreaterThan(firstBox!.y);
  });

  test('code blocks are scrollable, not overflowing', async ({ page }) => {
    await page.goto('/');
    const preBlocks = page.locator('pre');
    const count = await preBlocks.count();

    for (let i = 0; i < count; i++) {
      const box = await preBlocks.nth(i).boundingBox();
      if (box) {
        expect(box.width).toBeLessThanOrEqual(viewports.tablet.width);
      }
    }
  });

  test('all pages load on tablet', async ({ page }) => {
    const pages = ['/', '/#/docs', '/#/tutorials', '/#/api', '/#/cli'];
    for (const path of pages) {
      await page.goto(path);
      await page.waitForLoadState('networkidle');
      const h1 = await page.$('h1');
      expect(h1, `Page ${path} should have an h1 on tablet`).toBeTruthy();
    }
  });
});

// ----- Desktop viewport tests -----

test.describe('desktop viewport (1280x800)', () => {
  test.use({ viewport: viewports.desktop });

  test('nav links are visible inline (no hamburger)', async ({ page }) => {
    await page.goto('/');
    const navLinks = page.locator('.nav-links');
    await expect(navLinks).toBeVisible();

    const toggle = page.locator('button.nav-toggle');
    await expect(toggle).not.toBeVisible();
  });

  test('features grid shows multiple columns', async ({ page }) => {
    await page.goto('/');
    const cards = page.locator('.feature-card');
    const count = await cards.count();
    expect(count).toBeGreaterThan(1);

    const firstBox = await cards.first().boundingBox();
    const secondBox = await cards.nth(1).boundingBox();
    expect(firstBox).toBeTruthy();
    expect(secondBox).toBeTruthy();

    // Multi-column: second card should be on the same row (same y, different x)
    expect(secondBox!.x).toBeGreaterThan(firstBox!.x);
  });
});

// ----- Cross-viewport content consistency -----

test.describe('content consistency across viewports', () => {
  for (const [name, viewport] of Object.entries(viewports)) {
    test(`h1 content matches on ${name}`, async ({ page }) => {
      await page.setViewportSize(viewport);
      await page.goto('/');
      await expect(page.locator('h1')).toContainText('datascience');
    });

    test(`footer is present on ${name}`, async ({ page }) => {
      await page.setViewportSize(viewport);
      await page.goto('/');
      await expect(page.locator('footer.footer')).toBeVisible();
    });
  }
});

import { test, expect } from '@playwright/test';

// --- Desktop (>1024px) ---
test.describe('Desktop (1280x800)', () => {
  test.use({ viewport: { width: 1280, height: 800 } });

  test('nav bar has correct height and styling', async ({ page }) => {
    await page.goto('/');
    const nav = page.locator('.nav');
    const inner = page.locator('.nav-inner');
    const innerBox = await inner.boundingBox();
    expect(innerBox).toBeTruthy();
    // Nav inner should be 56px tall (matching SpecifyJS)
    expect(innerBox!.height).toBeGreaterThanOrEqual(54);
    expect(innerBox!.height).toBeLessThanOrEqual(58);
  });

  test('nav brand has correct font styling', async ({ page }) => {
    await page.goto('/');
    const brand = page.locator('.nav-brand');
    const styles = await brand.evaluate(el => {
      const cs = getComputedStyle(el);
      return {
        fontSize: cs.fontSize,
        fontWeight: cs.fontWeight,
        textDecoration: cs.textDecorationLine || cs.textDecoration,
      };
    });
    expect(styles.fontSize).toBe('20px');
    expect(styles.fontWeight).toBe('700');
    expect(styles.textDecoration).toContain('none');
  });

  test('nav brand "Data" is primary color and "Science" is text color', async ({ page }) => {
    await page.goto('/');
    const dataSpan = page.locator('.nav-brand-text span').first();
    const scienceSpan = page.locator('.nav-brand-text span').nth(1);

    const dataColor = await dataSpan.evaluate(el => getComputedStyle(el).color);
    const scienceColor = await scienceSpan.evaluate(el => getComputedStyle(el).color);

    // Data should be primary (teal) — not the same as Science (white)
    expect(dataColor).not.toBe(scienceColor);
    // Science should be the text color (white/light)
    expect(scienceColor).toMatch(/rgb\(241, 245, 249\)/);
  });

  test('nav brand is not styled as a hyperlink', async ({ page }) => {
    await page.goto('/');
    const brand = page.locator('.nav-brand');
    const styles = await brand.evaluate(el => {
      const cs = getComputedStyle(el);
      return {
        textDecoration: cs.textDecorationLine || cs.textDecoration,
        cursor: cs.cursor,
      };
    });
    expect(styles.textDecoration).toContain('none');
  });

  test('hamburger menu is hidden', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-toggle')).not.toBeVisible();
  });

  test('nav links are visible inline with pill styling', async ({ page }) => {
    await page.goto('/');
    const links = page.locator('.nav-links');
    await expect(links).toBeVisible();

    const firstLink = page.locator('.nav-links a').first();
    const styles = await firstLink.evaluate(el => {
      const cs = getComputedStyle(el);
      return {
        fontSize: cs.fontSize,
        fontWeight: cs.fontWeight,
        textDecoration: cs.textDecorationLine || cs.textDecoration,
        padding: cs.padding,
        borderRadius: cs.borderRadius,
      };
    });
    expect(styles.fontSize).toBe('14px');
    expect(styles.fontWeight).toBe('500');
    expect(styles.textDecoration).toContain('none');
    expect(styles.borderRadius).toBe('8px');
  });

  test('active nav link has primary background', async ({ page }) => {
    await page.goto('/');
    const activeLink = page.locator('.nav-links a.active');
    const bg = await activeLink.evaluate(el => getComputedStyle(el).backgroundColor);
    // Should have the primary teal background, not transparent
    expect(bg).not.toBe('rgba(0, 0, 0, 0)');
    expect(bg).not.toBe('transparent');
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

// --- Mobile (<768px) ---
test.describe('Mobile (375x667)', () => {
  test.use({ viewport: { width: 375, height: 667 } });

  test('nav bar spans full width', async ({ page }) => {
    await page.goto('/');
    const nav = page.locator('.nav');
    const navBox = await nav.boundingBox();
    expect(navBox).toBeTruthy();
    expect(navBox!.width).toBeGreaterThanOrEqual(375);
  });

  test('nav brand does not overflow viewport', async ({ page }) => {
    await page.goto('/');
    const brand = page.locator('.nav-brand');
    const brandBox = await brand.boundingBox();
    expect(brandBox).toBeTruthy();
    expect(brandBox!.x + brandBox!.width).toBeLessThanOrEqual(375);
  });

  test('nav brand font scales down on mobile', async ({ page }) => {
    await page.goto('/');
    const fontSize = await page.locator('.nav-brand').evaluate(
      el => parseFloat(getComputedStyle(el).fontSize)
    );
    expect(fontSize).toBeLessThanOrEqual(16);
  });

  test('hamburger menu is visible', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-toggle')).toBeVisible();
  });

  test('nav links are hidden by default', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-links')).not.toBeVisible();
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

  test('mobile nav dropdown is full-width and absolute', async ({ page }) => {
    await page.goto('/');
    await page.locator('.nav-toggle').click();
    const links = page.locator('.nav-links');
    const box = await links.boundingBox();
    expect(box).toBeTruthy();
    // Should span nearly the full viewport width
    expect(box!.width).toBeGreaterThanOrEqual(370);
  });

  test('nav link click closes menu and navigates', async ({ page }) => {
    await page.goto('/');
    await page.locator('.nav-toggle').click();
    await page.locator('.nav-links a:has-text("Docs")').click();
    await expect(page.locator('h1')).toContainText('Documentation');
    await expect(page.locator('.nav-links')).not.toBeVisible();
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
    expect(fontSize).toBeLessThanOrEqual(28);
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

// --- Tablet (768-1024px) ---
test.describe('Tablet (768x1024)', () => {
  test.use({ viewport: { width: 768, height: 1024 } });

  test('hamburger menu is visible', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-toggle')).toBeVisible();
  });

  test('nav links hidden by default, toggle works', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.nav-links')).not.toBeVisible();
    await page.locator('.nav-toggle').click();
    await expect(page.locator('.nav-links')).toBeVisible();
  });

  test('hero h1 does not overflow', async ({ page }) => {
    await page.goto('/');
    const h1Box = await page.locator('.hero h1').boundingBox();
    expect(h1Box).toBeTruthy();
    expect(h1Box!.x + h1Box!.width).toBeLessThanOrEqual(768);
  });

  test('features grid has two columns', async ({ page }) => {
    await page.goto('/');
    const cards = page.locator('.feature-card');
    const count = await cards.count();
    if (count >= 2) {
      const first = await cards.nth(0).boundingBox();
      const second = await cards.nth(1).boundingBox();
      expect(first).toBeTruthy();
      expect(second).toBeTruthy();
      // Two columns: same Y, different X
      expect(Math.abs(second!.y - first!.y)).toBeLessThan(10);
      expect(second!.x).toBeGreaterThan(first!.x);
    }
  });

  test('no horizontal scroll', async ({ page }) => {
    await page.goto('/');
    const scrollWidth = await page.evaluate(() => document.documentElement.scrollWidth);
    const clientWidth = await page.evaluate(() => document.documentElement.clientWidth);
    expect(scrollWidth).toBeLessThanOrEqual(clientWidth + 1);
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

    test(`${name}: nav brand has two-tone styling`, async ({ page }) => {
      await page.setViewportSize({ width, height });
      await page.goto('/');
      const dataSpan = page.locator('.nav-brand-text span').first();
      const scienceSpan = page.locator('.nav-brand-text span').nth(1);
      const dataColor = await dataSpan.evaluate(el => getComputedStyle(el).color);
      const scienceColor = await scienceSpan.evaluate(el => getComputedStyle(el).color);
      expect(dataColor).not.toBe(scienceColor);
    });
  }
});

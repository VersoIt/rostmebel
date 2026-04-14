import { expect, test } from '@playwright/test';

test.describe('public site smoke', () => {
  test('home renders without console errors or horizontal overflow', async ({ page }, testInfo) => {
    const consoleErrors: string[] = [];
    page.on('console', (message) => {
      if (message.type() === 'error') {
        consoleErrors.push(message.text());
      }
    });

    await page.goto('/');
    await expect(page.getByRole('heading', { name: /Кухни и мебель по размеру/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /Посмотреть проекты/i })).toBeVisible();
    await expect(page.getByText(/Расчет за минуту/i).first()).toBeVisible();

    const quickActions = page.getByRole('navigation', { name: /Быстрые действия/i });
    if (testInfo.project.name === 'mobile') {
      await expect(quickActions).toBeVisible();
      await expect(quickActions.getByRole('link', { name: /Расчет/i })).toBeVisible();
    } else {
      await expect(quickActions).toBeHidden();
    }

    const hasHorizontalOverflow = await page.evaluate(() => document.documentElement.scrollWidth > window.innerWidth + 1);
    expect(hasHorizontalOverflow).toBe(false);
    expect(consoleErrors).toEqual([]);

    await page.screenshot({ path: testInfo.outputPath(`home-${testInfo.project.name}.png`), fullPage: true });
  });

  test('contact page keeps the consultation path visible', async ({ page }) => {
    await page.goto('/contact');

    await expect(page.getByRole('heading', { name: /Обсудим мебель/i })).toBeVisible();
    await expect(page.getByRole('heading', { name: /Рассчитать проект/i })).toBeVisible();
    await expect(page.getByText(/Расчет за минуту/i)).toBeVisible();

    await page.getByRole('button', { name: /Дальше/i }).click();
    await page.getByRole('button', { name: /Нужен расчет/i }).click();
    await page.getByRole('button', { name: /Дальше/i }).click();
    await expect(page.getByLabel(/Город или район/i)).toBeVisible();
    await page.getByRole('button', { name: /Дальше/i }).click();
    await expect(page.getByLabel(/Имя/i)).toBeVisible();
    await expect(page.getByLabel(/Телефон/i)).toBeVisible();

    const hasHorizontalOverflow = await page.evaluate(() => document.documentElement.scrollWidth > window.innerWidth + 1);
    expect(hasHorizontalOverflow).toBe(false);
  });
});

test.describe('backend contract smoke', () => {
  test('readiness endpoint reports dependency status', async ({ request }) => {
    const response = await request.get('/readyz');
    expect(response.ok()).toBe(true);

    const body = await response.json();
    expect(body.status).toBe('ready');
    expect(body.checks.postgres).toBe('ok');
    expect(body.checks.redis).toBe('ok');
  });

  test('validation errors expose code and metadata', async ({ request }) => {
    const response = await request.post('/api/v1/orders', { data: {} });
    expect(response.status()).toBe(400);

    const body = await response.json();
    expect(body.error.code).toBe('VALIDATION_FAILED');
    expect(body.error.meta.fields).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ field: 'client_name', rule: 'required' }),
        expect.objectContaining({ field: 'client_phone', rule: 'required' }),
      ]),
    );
  });
});

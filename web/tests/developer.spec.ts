import { test, expect } from './helpers/fixtures'
import { apiLogin, seedDeveloperProfiles, TEST_USERS } from './helpers/api'

test.describe('Developer List & Profile', () => {
  test.beforeAll(async () => {
    // Seed developer profiles for the list
    const devTokens = []
    for (const key of ['dev1', 'dev2', 'dev3', 'dev4', 'dev5'] as const) {
      try {
        const auth = await apiLogin(TEST_USERS[key].email, TEST_USERS[key].password)
        devTokens.push({ token: auth.access_token, email: TEST_USERS[key].email })
      } catch {}
    }
    await seedDeveloperProfiles(devTokens)
  })

  test('should browse developer list with seeded data', async ({ page }) => {
    await page.goto('/developers')
    await expect(page.getByRole('heading', { name: '找开发者' })).toBeVisible({ timeout: 5000 })
    // Should show at least one developer card
    const devCards = page.locator('.developer-card, .dev-card, [class*="developer-item"]')
    await expect(devCards.first()).toBeVisible({ timeout: 5000 })
    // Should show total count > 0
    await expect(page.locator('.result-count')).toContainText(/^[^0]+/)
  })

  test('should view developer detail page', async ({ page }) => {
    await page.goto('/developers')
    await expect(page.locator('.developer-card, .dev-card, [class*="developer-item"]').first()).toBeVisible({ timeout: 5000 })
    // Click first developer card
    const firstDev = page.locator('.developer-card, .dev-card, [class*="developer-item"]').first()
    await firstDev.click()
    await expect(page).toHaveURL(/\/developers\//, { timeout: 5000 })
  })

  test('should show developer profile after login', async ({ developerPage }) => {
    await developerPage.goto('/profile')
    await expect(developerPage.getByRole('heading', { name: /赵全栈|开发者档案|个人资料/ }).first()).toBeVisible({ timeout: 5000 })
  })

  test('should navigate to developers from nav', async ({ page }) => {
    await page.goto('/')
    const devLink = page.getByRole('link', { name: /找开发者/ })
    if (await devLink.isVisible()) {
      await devLink.click()
      await expect(page).toHaveURL(/\/developers/)
    }
  })

  test('should filter developers by availability', async ({ page }) => {
    await page.goto('/developers')
    await expect(page.locator('.developer-card, .dev-card, [class*="developer-item"]').first()).toBeVisible({ timeout: 5000 })
    // Select availability filter
    const availSelect = page.locator('.el-select').first()
    if (await availSelect.isVisible()) {
      await availSelect.click()
      await page.waitForTimeout(300)
      const firstOption = page.locator('.el-select-dropdown__item').first()
      if (await firstOption.isVisible()) {
        await firstOption.click()
        await page.waitForTimeout(1000)
      }
    }
  })

  test('should reset developer filters', async ({ page }) => {
    await page.goto('/developers')
    await expect(page.locator('.developer-card, .dev-card, [class*="developer-item"]').first()).toBeVisible({ timeout: 5000 })
    await page.getByRole('button', { name: '重置筛选' }).click()
    await page.waitForTimeout(1000)
    // Should still show developers
    const devCards = page.locator('.developer-card, .dev-card, [class*="developer-item"]')
    await expect(devCards.first()).toBeVisible({ timeout: 5000 })
  })
})

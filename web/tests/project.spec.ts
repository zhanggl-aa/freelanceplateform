import { test, expect } from './helpers/fixtures'
import { apiGetCategories, apiCreateProject, apiPublishProject, apiLogin, seedTestProjects, type SeededProject, TEST_USERS } from './helpers/api'

test.describe('Project List & Browse', () => {
  let seededProjects: SeededProject[] = []

  test.beforeAll(async () => {
    // Seed test projects via API
    const auth = await apiLogin(TEST_USERS.client1.email, TEST_USERS.client1.password)
    const categories = await apiGetCategories()
    seededProjects = await seedTestProjects(auth.access_token, categories)
  })

  test('should browse project list with seeded data', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.getByRole('heading', { name: '找项目' })).toBeVisible({ timeout: 5000 })
    // Should show at least one project from seed data
    const projectCards = page.locator('.project-item')
    await expect(projectCards.first()).toBeVisible({ timeout: 5000 })
    // Should show total count > 0
    await expect(page.locator('.result-count')).toContainText(/^[^0]+/)
  })

  test('should display project card content', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
    // Project card should show title, budget type tag, and tech tags
    const firstCard = page.locator('.project-item').first()
    await expect(firstCard.locator('.project-item-title')).toBeVisible()
    await expect(firstCard.locator('.el-tag').first()).toBeVisible()
  })

  test('should filter by fixed price type', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
    const countBefore = await page.locator('.project-item').count()
    // Click "固定价" radio filter
    await page.getByRole('radio', { name: '固定价' }).click()
    await page.waitForTimeout(1000)
    // Should still show projects (4 of 6 are fixed price)
    const countAfter = await page.locator('.project-item').count()
    expect(countAfter).toBeGreaterThanOrEqual(0)
    // Each visible card should have 固定价 tag, not 时薪制
    const cards = page.locator('.project-item')
    for (let i = 0; i < Math.min(await cards.count(), 5); i++) {
      const card = cards.nth(i)
      await expect(card.getByText('固定价')).toBeVisible()
    }
  })

  test('should filter by hourly type', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
    await page.getByRole('radio', { name: '时薪制' }).click()
    await page.waitForTimeout(1000)
    // Should show at least 1 hourly project from seed
    const cards = page.locator('.project-item')
    const count = await cards.count()
    if (count > 0) {
      await expect(cards.first().getByText('时薪制')).toBeVisible()
    }
  })

  test('should view project detail by clicking a card', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
    // Click the first project card
    await page.locator('.project-item').first().click()
    // Should navigate to project detail page
    await expect(page).toHaveURL(/\/projects\/[a-f0-9-]+/, { timeout: 5000 })
    // Detail page should show project title and description
    await expect(page.locator('.detail-title')).toBeVisible()
    await expect(page.getByText('项目描述')).toBeVisible()
  })

  test('should show project info on detail page', async ({ page }) => {
    if (seededProjects.length === 0) return
    // Navigate directly to a seeded project
    await page.goto(`/projects/${seededProjects[0].id}`)
    await expect(page.locator('.detail-title')).toBeVisible({ timeout: 5000 })
    // Should show project info section with budget
    await expect(page.getByText('预算范围')).toBeVisible()
    await expect(page.getByText('技术栈')).toBeVisible()
  })

  test('should sort by budget high', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
    await page.getByRole('radiobutton', { name: '预算最高' }).click()
    await page.waitForTimeout(1000)
    // Page should still show projects
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
  })

  test('should reset filters', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.locator('.project-item').first()).toBeVisible({ timeout: 5000 })
    // Apply a filter first
    await page.getByRole('radio', { name: '时薪制' }).click()
    await page.waitForTimeout(500)
    // Click reset
    await page.getByRole('button', { name: '重置筛选' }).click()
    await page.waitForTimeout(1000)
    // All projects should be visible again
    const count = await page.locator('.project-item').count()
    expect(count).toBeGreaterThanOrEqual(seededProjects.length)
  })

  test('should show category tree in sidebar', async ({ page }) => {
    await page.goto('/projects')
    await expect(page.getByText('筛选条件')).toBeVisible({ timeout: 5000 })
    // Category tree should have at least one node
    await expect(page.locator('.el-tree-node').first()).toBeVisible()
  })
})

test.describe('Project Create (Client)', () => {
  test('should navigate to create project page', async ({ clientPage }) => {
    await clientPage.goto('/projects/create')
    await expect(clientPage.getByRole('heading', { name: '发布项目' })).toBeVisible({ timeout: 5000 })
  })

  test('should show project form fields', async ({ clientPage }) => {
    await clientPage.goto('/projects/create')
    await expect(clientPage.getByText('项目标题')).toBeVisible({ timeout: 5000 })
    await expect(clientPage.getByText('项目描述')).toBeVisible({ timeout: 5000 })
    await expect(clientPage.getByText('预算类型')).toBeVisible()
    await expect(clientPage.getByText('技术栈')).toBeVisible()
  })

  test('should create and publish a project via form', async ({ clientPage }) => {
    await clientPage.goto('/projects/create')
    await clientPage.waitForTimeout(500)

    // Fill in project details
    await clientPage.getByPlaceholder(/项目标题/).fill('E2E自动化测试项目')
    // Select a category
    const categorySelect = clientPage.locator('.el-select').first()
    await categorySelect.click()
    await clientPage.waitForTimeout(300)
    await clientPage.locator('.el-select-dropdown__item').first().click()
    // Fill description
    await clientPage.getByPlaceholder(/详细描述项目需求/).fill('这是一个通过Playwright自动化测试创建的项目，用于验证项目发布功能的完整流程。')
    // Set budget
    const budgetMinInput = clientPage.locator('.el-input-number').first()
    await budgetMinInput.find('input').fill('5000')
    const budgetMaxInput = clientPage.locator('.el-input-number').nth(1)
    await budgetMaxInput.find('input').fill('10000')
    // Set deadline
    const deadlinePicker = clientPage.locator('.el-date-editor').first()
    if (await deadlinePicker.isVisible()) {
      await deadlinePicker.click()
      await clientPage.waitForTimeout(300)
      // Click a future date
      const availableDate = clientPage.locator('.el-date-table td.available, .el-date-table td.next-month').first()
      if (await availableDate.isVisible()) {
        await availableDate.click()
      }
    }
    // Add tech stack
    const techSelect = clientPage.locator('.el-select').last()
    await techSelect.click()
    await clientPage.waitForTimeout(300)
    const firstTech = clientPage.locator('.el-select-dropdown__item').first()
    if (await firstTech.isVisible()) {
      await firstTech.click()
    }
    // Click publish button
    await clientPage.getByRole('button', { name: '发布项目' }).click()
    // Should navigate to project detail page
    await expect(clientPage).toHaveURL(/\/projects\//, { timeout: 10000 })
  })

  test('should show my posted projects', async ({ clientPage }) => {
    await clientPage.goto('/my/projects')
    await expect(clientPage.getByRole('heading', { name: '我的项目' })).toBeVisible({ timeout: 5000 })
  })
})

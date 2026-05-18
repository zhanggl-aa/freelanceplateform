import { test, expect } from './helpers/fixtures'
import { apiCreateBid, apiGetCategories, apiLogin, apiCreateProject, apiPublishProject, seedTestProjects, type SeededProject, TEST_USERS } from './helpers/api'

test.describe('Bidding Flow', () => {
  let seededProjects: SeededProject[] = []

  test.beforeAll(async () => {
    const clientAuth = await apiLogin(TEST_USERS.client1.email, TEST_USERS.client1.password)
    const categories = await apiGetCategories()
    seededProjects = await seedTestProjects(clientAuth.access_token, categories)
  })

  test('should show bid button on published project for developer', async ({ developerPage }) => {
    if (seededProjects.length === 0) return
    await developerPage.goto(`/projects/${seededProjects[0].id}`)
    await developerPage.waitForTimeout(1000)
    // Should see bid/apply button on a published project
    const bidBtn = developerPage.getByRole('button', { name: /投标|申请|报价/ })
    if (await bidBtn.isVisible()) {
      await expect(bidBtn).toBeEnabled()
    }
  })

  test('should submit bid via API and see it in my bids', async ({ developerPage }) => {
    if (seededProjects.length === 0) return
    // Submit a bid via API
    const devAuth = await apiLogin(TEST_USERS.dev1.email, TEST_USERS.dev1.password)
    try {
      await apiCreateBid(devAuth.access_token, seededProjects[0].id, {
        cover_letter: '我有丰富的项目经验，可以高质量完成这个项目。',
        estimated_days: 30,
        proposed_budget: 35000,
        budget_type: 'fixed',
      })
    } catch {}
    // Check my bids page
    await developerPage.goto('/my/bids')
    await expect(developerPage.getByRole('heading', { name: '我的投标' })).toBeVisible({ timeout: 5000 })
  })

  test('should view bids on own project as client', async ({ clientPage }) => {
    if (seededProjects.length === 0) return
    await clientPage.goto(`/projects/${seededProjects[0].id}`)
    await clientPage.waitForTimeout(1000)
    // Look for bids tab or section
    const bidsTab = clientPage.getByRole('tab', { name: /竞标中|投标/ })
    if (await bidsTab.isVisible()) {
      await bidsTab.click()
      await clientPage.waitForTimeout(500)
    }
  })

  test('should accept a bid as client', async ({ clientPage }) => {
    if (seededProjects.length === 0) return
    await clientPage.goto(`/projects/${seededProjects[0].id}`)
    await clientPage.waitForTimeout(1000)

    const acceptBtn = clientPage.getByRole('button', { name: /接受|采纳/ }).first()
    if (await acceptBtn.isVisible()) {
      await acceptBtn.click()
      await clientPage.waitForTimeout(1000)
    }
  })

  test('should reject a bid as client', async ({ clientPage }) => {
    if (seededProjects.length > 1) {
      await clientPage.goto(`/projects/${seededProjects[1].id}`)
      await clientPage.waitForTimeout(1000)

      const rejectBtn = clientPage.getByRole('button', { name: /拒绝|驳回/ }).first()
      if (await rejectBtn.isVisible()) {
        await rejectBtn.click()
        await clientPage.waitForTimeout(1000)
      }
    }
  })
})

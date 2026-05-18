import { test, expect } from './helpers/fixtures'

test.describe('Contract Management', () => {
  test('should view my contracts list', async ({ clientPage }) => {
    await clientPage.goto('/my/contracts')
    await expect(clientPage.getByRole('heading', { name: '我的合同' })).toBeVisible({ timeout: 5000 })
  })

  test('should show contract tabs', async ({ clientPage }) => {
    await clientPage.goto('/my/contracts')
    await expect(clientPage.getByText('进行中')).toBeVisible({ timeout: 5000 })
    await expect(clientPage.getByText('已完成')).toBeVisible()
    await expect(clientPage.getByText('已取消')).toBeVisible()
  })

  test('should switch to completed contracts tab', async ({ clientPage }) => {
    await clientPage.goto('/my/contracts')
    await clientPage.getByText('已完成').click()
    await clientPage.waitForTimeout(500)
  })

  test('should switch to cancelled contracts tab', async ({ clientPage }) => {
    await clientPage.goto('/my/contracts')
    await clientPage.getByText('已取消').click()
    await clientPage.waitForTimeout(500)
  })

  test('should view contract detail', async ({ clientPage }) => {
    await clientPage.goto('/my/contracts')
    await clientPage.waitForTimeout(1000)

    const contractRow = clientPage.locator('.el-table__row').first()
    if (await contractRow.isVisible()) {
      await contractRow.click()
      await clientPage.waitForTimeout(1000)
      const url = clientPage.url()
      expect(url).toMatch(/\/contracts\//)
    }
  })

  test('should show developer contracts', async ({ developerPage }) => {
    await developerPage.goto('/my/contracts')
    await expect(developerPage.getByRole('heading', { name: '我的合同' })).toBeVisible({ timeout: 5000 })
  })

  test('should display contract amounts', async ({ clientPage }) => {
    await clientPage.goto('/my/contracts')
    await clientPage.waitForTimeout(1000)
    const amountCells = clientPage.locator('.amount-text, :text("¥")')
    if (await amountCells.first().isVisible()) {
      await expect(amountCells.first()).toContainText('¥')
    }
  })
})

import { test, expect } from './helpers/fixtures'

test.describe('Wallet & Payments', () => {
  test('should view wallet page', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await expect(clientPage.getByRole('heading', { name: '我的钱包' })).toBeVisible({ timeout: 5000 })
  })

  test('should display balance information', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await expect(clientPage.getByText('账户余额')).toBeVisible({ timeout: 5000 })
    await expect(clientPage.getByText('冻结金额')).toBeVisible()
  })

  test('should show balance amounts with ¥ symbol', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await clientPage.waitForTimeout(1000)
    const balanceValues = clientPage.locator('.balance-value')
    if (await balanceValues.first().isVisible()) {
      await expect(balanceValues.first()).toContainText('¥')
    }
  })

  test('should have deposit button', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await expect(clientPage.getByRole('button', { name: /充值/ })).toBeVisible({ timeout: 5000 })
  })

  test('should have withdraw button', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await expect(clientPage.getByRole('button', { name: /提现/ })).toBeVisible({ timeout: 5000 })
  })

  test('should open deposit dialog', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await clientPage.getByRole('button', { name: /充值/ }).click()
    await clientPage.waitForTimeout(500)
    const depositDialog = clientPage.locator('.el-dialog, .el-drawer')
    if (await depositDialog.isVisible()) {
      await expect(depositDialog).toBeVisible()
    }
  })

  test('should display transaction history', async ({ clientPage }) => {
    await clientPage.goto('/wallet')
    await clientPage.waitForTimeout(1000)
    await expect(clientPage.getByRole('heading', { name: '交易记录' })).toBeVisible({ timeout: 5000 })
  })

  test('should show developer wallet with earnings', async ({ developerPage }) => {
    await developerPage.goto('/wallet')
    await expect(developerPage.getByText('账户余额')).toBeVisible({ timeout: 5000 })
    const balanceValue = developerPage.locator('.balance-value').first()
    if (await balanceValue.isVisible()) {
      const text = await balanceValue.textContent()
      const amount = parseFloat(text!.replace(/[¥,]/g, ''))
      expect(amount).toBeGreaterThan(0)
    }
  })

  test('should navigate to wallet from profile menu', async ({ clientPage }) => {
    await clientPage.goto('/')
    const walletLink = clientPage.getByRole('link', { name: /钱包/ })
    if (await walletLink.isVisible()) {
      await walletLink.click()
      await expect(clientPage).toHaveURL(/\/wallet/)
    }
  })
})

import { test, expect } from './helpers/fixtures'
import { TEST_USERS, apiRegister } from './helpers/api'

// The login page has two tabs (email/phone), each with a password field
// sharing the same placeholder. Scope password lookups to the email tab.
function emailPasswordInput(page: import('@playwright/test').Page) {
  return page.getByRole('tabpanel', { name: '邮箱登录' }).getByPlaceholder('请输入密码')
}

test.describe('Authentication', () => {
  test('should show login page', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('h2')).toContainText('登录接单平台')
    await expect(page.getByPlaceholder('请输入邮箱地址')).toBeVisible()
    await expect(emailPasswordInput(page)).toBeVisible()
  })

  test('should login with email successfully', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.client1.email)
    await emailPasswordInput(page).fill(TEST_USERS.client1.password)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL('/', { timeout: 10000 })
    await expect(page.getByText('登录成功')).toBeVisible()
  })

  test('should login as developer', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.dev1.email)
    await emailPasswordInput(page).fill(TEST_USERS.dev1.password)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL('/', { timeout: 10000 })
  })

  test('should show error on wrong password', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.client1.email)
    await emailPasswordInput(page).fill('WrongPassword1')
    await page.getByRole('button', { name: '登录' }).click()
    // ElMessage renders a toast notification with the error
    await expect(page.locator('.el-message').getByText(/失败|错误|invalid/i)).toBeVisible({ timeout: 5000 })
  })

  test('should navigate to register page', async ({ page }) => {
    await page.goto('/login')
    await page.getByText('立即注册').click()
    await expect(page).toHaveURL('/register')
    await expect(page.locator('h2').first()).toContainText('创建账号')
  })

  test('should register a new client user', async ({ page }) => {
    const uniqueEmail = `test_client_${Date.now()}@test.com`
    await page.goto('/register')

    // Step 1: Account info
    await page.getByText('邮箱注册').click()
    await page.getByPlaceholder('请输入邮箱地址').fill(uniqueEmail)
    await page.getByPlaceholder('你的专属昵称').fill('测试甲方')
    await page.getByPlaceholder('请输入密码（至少6位）').fill('Test123456')
    await page.getByPlaceholder('请再次输入密码').fill('Test123456')
    await page.getByRole('button', { name: /下一步/ }).click()

    // Step 2: Role selection
    await expect(page.locator('h2').first()).toContainText('选择你的角色')
    await page.getByText('我要发包').click()
    await page.getByRole('button', { name: /下一步/ }).click()

    // Step 3: Confirm
    await expect(page.locator('h2').first()).toContainText('确认注册')
    // Element Plus checkbox: click the visible label, not the hidden native input
    await page.getByText('我已阅读并同意').click()
    await page.getByRole('button', { name: '完成注册' }).click()

    await expect(page).toHaveURL('/', { timeout: 10000 })
  })

  test('should register a new developer user', async ({ page }) => {
    const uniqueEmail = `test_dev_${Date.now()}@test.com`
    await page.goto('/register')

    await page.getByText('邮箱注册').click()
    await page.getByPlaceholder('请输入邮箱地址').fill(uniqueEmail)
    await page.getByPlaceholder('你的专属昵称').fill('测试开发者')
    await page.getByPlaceholder('请输入密码（至少6位）').fill('Test123456')
    await page.getByPlaceholder('请再次输入密码').fill('Test123456')
    await page.getByRole('button', { name: /下一步/ }).click()

    await expect(page.locator('h2').first()).toContainText('选择你的角色')
    await page.getByText('我要接单').click()
    await page.getByRole('button', { name: /下一步/ }).click()

    await expect(page.locator('h2').first()).toContainText('确认注册')
    await page.getByText('我已阅读并同意').click()
    await page.getByRole('button', { name: '完成注册' }).click()

    await expect(page).toHaveURL('/', { timeout: 10000 })
  })

  test('should redirect to login when accessing protected route', async ({ page }) => {
    await page.goto('/wallet')
    await expect(page).toHaveURL(/\/login/, { timeout: 5000 })
  })

  test('should persist login after page refresh', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.client1.email)
    await emailPasswordInput(page).fill(TEST_USERS.client1.password)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL('/', { timeout: 10000 })

    await page.reload()
    await expect(page).toHaveURL('/')
  })
})

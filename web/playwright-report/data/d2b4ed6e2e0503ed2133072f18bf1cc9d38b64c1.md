# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: auth.spec.ts >> Authentication >> should show error on wrong password
- Location: tests\auth.spec.ts:35:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/失败|错误|invalid/i)
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/失败|错误|invalid/i)

```

```yaml
- img
- heading "登录接单平台" [level=2]
- paragraph: 欢迎回来，请登录您的账号
- tablist:
  - tab "邮箱登录" [selected]
  - tab "手机登录"
- tabpanel "邮箱登录":
  - text: "*邮箱地址"
  - img
  - textbox "*邮箱地址":
    - /placeholder: 请输入邮箱地址
    - text: client1@test.com
  - text: "*密码"
  - img
  - textbox "*密码":
    - /placeholder: 请输入密码
    - text: WrongPassword1
  - img
  - button "登录"
- link "忘记密码？":
  - /url: /forgot-password
- text: "|"
- link "还没有账号？立即注册":
  - /url: /register
```

# Test source

```ts
  1   | import { test, expect } from './helpers/fixtures'
  2   | import { TEST_USERS, apiRegister } from './helpers/api'
  3   | 
  4   | // The login page has two tabs (email/phone), each with a password field
  5   | // sharing the same placeholder. Scope password lookups to the email tab.
  6   | function emailPasswordInput(page: import('@playwright/test').Page) {
  7   |   return page.getByRole('tabpanel', { name: '邮箱登录' }).getByPlaceholder('请输入密码')
  8   | }
  9   | 
  10  | test.describe('Authentication', () => {
  11  |   test('should show login page', async ({ page }) => {
  12  |     await page.goto('/login')
  13  |     await expect(page.locator('h2')).toContainText('登录接单平台')
  14  |     await expect(page.getByPlaceholder('请输入邮箱地址')).toBeVisible()
  15  |     await expect(emailPasswordInput(page)).toBeVisible()
  16  |   })
  17  | 
  18  |   test('should login with email successfully', async ({ page }) => {
  19  |     await page.goto('/login')
  20  |     await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.client1.email)
  21  |     await emailPasswordInput(page).fill(TEST_USERS.client1.password)
  22  |     await page.getByRole('button', { name: '登录' }).click()
  23  |     await expect(page).toHaveURL('/', { timeout: 10000 })
  24  |     await expect(page.getByText('登录成功')).toBeVisible()
  25  |   })
  26  | 
  27  |   test('should login as developer', async ({ page }) => {
  28  |     await page.goto('/login')
  29  |     await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.dev1.email)
  30  |     await emailPasswordInput(page).fill(TEST_USERS.dev1.password)
  31  |     await page.getByRole('button', { name: '登录' }).click()
  32  |     await expect(page).toHaveURL('/', { timeout: 10000 })
  33  |   })
  34  | 
  35  |   test('should show error on wrong password', async ({ page }) => {
  36  |     await page.goto('/login')
  37  |     await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.client1.email)
  38  |     await emailPasswordInput(page).fill('WrongPassword1')
  39  |     await page.getByRole('button', { name: '登录' }).click()
> 40  |     await expect(page.getByText(/失败|错误|invalid/i)).toBeVisible({ timeout: 5000 })
      |                                                    ^ Error: expect(locator).toBeVisible() failed
  41  |   })
  42  | 
  43  |   test('should navigate to register page', async ({ page }) => {
  44  |     await page.goto('/login')
  45  |     await page.getByText('立即注册').click()
  46  |     await expect(page).toHaveURL('/register')
  47  |     await expect(page.locator('h2').first()).toContainText('创建账号')
  48  |   })
  49  | 
  50  |   test('should register a new client user', async ({ page }) => {
  51  |     const uniqueEmail = `test_client_${Date.now()}@test.com`
  52  |     await page.goto('/register')
  53  | 
  54  |     // Step 1: Account info
  55  |     await page.getByText('邮箱注册').click()
  56  |     await page.getByPlaceholder('请输入邮箱地址').fill(uniqueEmail)
  57  |     await page.getByPlaceholder('你的专属昵称').fill('测试甲方')
  58  |     await page.getByPlaceholder('请输入密码（至少6位）').fill('Test123456')
  59  |     await page.getByPlaceholder('请再次输入密码').fill('Test123456')
  60  |     await page.getByRole('button', { name: /下一步/ }).click()
  61  | 
  62  |     // Step 2: Role selection
  63  |     await expect(page.locator('h2').first()).toContainText('选择你的角色')
  64  |     await page.getByText('我要发包').click()
  65  |     await page.getByRole('button', { name: /下一步/ }).click()
  66  | 
  67  |     // Step 3: Confirm
  68  |     await expect(page.locator('h2').first()).toContainText('确认注册')
  69  |     await page.getByRole('checkbox').click()
  70  |     await page.getByRole('button', { name: '完成注册' }).click()
  71  | 
  72  |     await expect(page).toHaveURL('/', { timeout: 10000 })
  73  |   })
  74  | 
  75  |   test('should register a new developer user', async ({ page }) => {
  76  |     const uniqueEmail = `test_dev_${Date.now()}@test.com`
  77  |     await page.goto('/register')
  78  | 
  79  |     await page.getByText('邮箱注册').click()
  80  |     await page.getByPlaceholder('请输入邮箱地址').fill(uniqueEmail)
  81  |     await page.getByPlaceholder('你的专属昵称').fill('测试开发者')
  82  |     await page.getByPlaceholder('请输入密码（至少6位）').fill('Test123456')
  83  |     await page.getByPlaceholder('请再次输入密码').fill('Test123456')
  84  |     await page.getByRole('button', { name: /下一步/ }).click()
  85  | 
  86  |     await expect(page.locator('h2').first()).toContainText('选择你的角色')
  87  |     await page.getByText('我要接单').click()
  88  |     await page.getByRole('button', { name: /下一步/ }).click()
  89  | 
  90  |     await expect(page.locator('h2').first()).toContainText('确认注册')
  91  |     await page.getByRole('checkbox').click()
  92  |     await page.getByRole('button', { name: '完成注册' }).click()
  93  | 
  94  |     await expect(page).toHaveURL('/', { timeout: 10000 })
  95  |   })
  96  | 
  97  |   test('should redirect to login when accessing protected route', async ({ page }) => {
  98  |     await page.goto('/wallet')
  99  |     await expect(page).toHaveURL(/\/login/, { timeout: 5000 })
  100 |   })
  101 | 
  102 |   test('should persist login after page refresh', async ({ page }) => {
  103 |     await page.goto('/login')
  104 |     await page.getByPlaceholder('请输入邮箱地址').fill(TEST_USERS.client1.email)
  105 |     await emailPasswordInput(page).fill(TEST_USERS.client1.password)
  106 |     await page.getByRole('button', { name: '登录' }).click()
  107 |     await expect(page).toHaveURL('/', { timeout: 10000 })
  108 | 
  109 |     await page.reload()
  110 |     await expect(page).toHaveURL('/')
  111 |   })
  112 | })
  113 | 
```
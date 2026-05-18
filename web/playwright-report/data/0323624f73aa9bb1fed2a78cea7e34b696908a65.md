# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: auth.spec.ts >> Authentication >> should register a new client user
- Location: tests\auth.spec.ts:50:3

# Error details

```
Test timeout of 30000ms exceeded.
```

```
Error: locator.click: Test timeout of 30000ms exceeded.
Call log:
  - waiting for getByRole('checkbox')
    - locator resolved to <input type="checkbox" class="el-checkbox__original"/>
  - attempting click action
    2 × waiting for element to be visible, enabled and stable
      - element is not stable
    - retrying click action
    - waiting 20ms
    2 × waiting for element to be visible, enabled and stable
      - element is not stable
    - retrying click action
      - waiting 100ms
    52 × waiting for element to be visible, enabled and stable
       - element is not visible
     - retrying click action
       - waiting 500ms

```

# Page snapshot

```yaml
- generic [ref=e4]:
  - generic [ref=e5]:
    - generic [ref=e6]:
      - generic [ref=e7]: "1"
      - generic [ref=e8]: 账号信息
    - generic [ref=e10]:
      - generic [ref=e11]: "2"
      - generic [ref=e12]: 选择角色
    - generic [ref=e14]:
      - generic [ref=e15]: "3"
      - generic [ref=e16]: 完成注册
  - generic [ref=e18]:
    - generic [ref=e19]:
      - img [ref=e22]
      - heading "确认注册" [level=2] [ref=e25]
      - paragraph [ref=e26]: 请确认以下信息无误
    - generic [ref=e27]:
      - generic [ref=e28]:
        - generic [ref=e29]: 账号
        - generic [ref=e30]: test_client_1778929850213@test.com
      - generic [ref=e31]:
        - generic [ref=e32]: 昵称
        - generic [ref=e33]: 测试甲方
      - generic [ref=e34]:
        - generic [ref=e35]: 角色
        - generic [ref=e37]: 甲方
    - generic [ref=e40] [cursor=pointer]:
      - generic [ref=e41]:
        - checkbox "我已阅读并同意 《用户协议》和《隐私政策》"
      - generic [ref=e43]:
        - text: 我已阅读并同意
        - link "《用户协议》" [ref=e44]:
          - /url: javascript:void(0)
        - text: 和
        - link "《隐私政策》" [ref=e45]:
          - /url: javascript:void(0)
    - generic [ref=e46]:
      - button "上一步" [ref=e47] [cursor=pointer]:
        - generic [ref=e48]:
          - img [ref=e50]
          - text: 上一步
      - button "完成注册" [ref=e52] [cursor=pointer]:
        - generic [ref=e53]: 完成注册
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
  40  |     await expect(page.getByText(/失败|错误|invalid/i)).toBeVisible({ timeout: 5000 })
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
> 69  |     await page.getByRole('checkbox').click()
      |                                      ^ Error: locator.click: Test timeout of 30000ms exceeded.
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
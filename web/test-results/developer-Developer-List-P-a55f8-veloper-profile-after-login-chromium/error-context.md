# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: developer.spec.ts >> Developer List & Profile >> should show developer profile after login
- Location: tests\developer.spec.ts:21:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/赵全栈|开发者|profile/i)
Expected: visible
Error: strict mode violation: getByText(/赵全栈|开发者|profile/i) resolved to 4 elements:
    1) <a class="nav-link" data-v-8da451f1="" href="/developers">找开发者</a> aka getByRole('link', { name: '找开发者' })
    2) <span data-v-8da451f1="" class="user-name hide-on-mobile">赵全栈</span> aka getByRole('button', { name: '赵全栈' })
    3) <div role="tab" tabindex="-1" id="tab-developer" aria-selected="false" class="el-tabs__item is-top" aria-controls="pane-developer">…</div> aka getByRole('tab', { name: '开发者档案' })
    4) <p data-v-8da451f1="">© 2026 接单平台 - 专业的软件开发者自由职业平台</p> aka getByText('© 2026 接单平台 - 专业的软件开发者自由职业平台')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/赵全栈|开发者|profile/i)

```

# Page snapshot

```yaml
- generic [ref=e3]:
  - banner [ref=e4]:
    - generic [ref=e5]:
      - generic [ref=e6]:
        - link "接单平台" [ref=e7] [cursor=pointer]:
          - /url: /
          - img [ref=e9]
          - generic [ref=e11]: 接单平台
        - navigation [ref=e12]:
          - link "首页" [ref=e13] [cursor=pointer]:
            - /url: /
          - link "找项目" [ref=e14] [cursor=pointer]:
            - /url: /projects
          - link "找开发者" [ref=e15] [cursor=pointer]:
            - /url: /developers
      - generic [ref=e16]:
        - link [ref=e18] [cursor=pointer]:
          - /url: /notifications
          - img [ref=e20]
        - link [ref=e25] [cursor=pointer]:
          - /url: /chat
          - img [ref=e27]
        - button "赵全栈" [ref=e31] [cursor=pointer]:
          - img [ref=e34]
          - generic [ref=e36]: 赵全栈
          - img [ref=e38]
  - main [ref=e41]:
    - generic [ref=e42]:
      - generic [ref=e43]:
        - heading "个人中心" [level=2] [ref=e44]
        - paragraph [ref=e45]: 管理你的个人资料和档案信息
      - generic [ref=e46]:
        - tablist [ref=e50]:
          - tab "基本信息" [selected] [ref=e52]
          - tab "开发者档案" [ref=e53]
        - tabpanel "基本信息" [ref=e55]:
          - generic [ref=e58]:
            - group "头像" [ref=e59]:
              - generic [ref=e60]: 头像
              - generic [ref=e62]:
                - img [ref=e65]
                - button "更换头像" [ref=e68] [cursor=pointer]:
                  - button "更换头像" [ref=e69]:
                    - generic [ref=e70]: 更换头像
            - generic [ref=e71]:
              - generic [ref=e72]: "* 昵称"
              - generic [ref=e75]:
                - textbox "* 昵称" [ref=e76]:
                  - /placeholder: 请输入昵称
                  - text: 赵全栈
                - generic [ref=e79]: 3 / 20
            - group "邮箱" [ref=e80]:
              - generic [ref=e81]: 邮箱
              - generic [ref=e83]:
                - generic [ref=e84]: dev1@test.com
                - generic [ref=e86]: 已验证
            - group "手机号" [ref=e87]:
              - generic [ref=e88]: 手机号
              - generic [ref=e90]:
                - generic [ref=e91]: 未绑定
                - generic [ref=e93]: 未验证
            - button "保存修改" [ref=e96] [cursor=pointer]:
              - generic [ref=e97]: 保存修改
  - contentinfo [ref=e98]:
    - paragraph [ref=e100]: © 2026 接单平台 - 专业的软件开发者自由职业平台
```

# Test source

```ts
  1  | import { test, expect } from './helpers/fixtures'
  2  | 
  3  | test.describe('Developer List & Profile', () => {
  4  |   test('should browse developer list as public', async ({ page }) => {
  5  |     await page.goto('/developers')
  6  |     await expect(page.getByText(/开发者|接单/)).toBeVisible({ timeout: 5000 })
  7  |   })
  8  | 
  9  |   test('should view developer detail page', async ({ page }) => {
  10 |     await page.goto('/developers')
  11 |     // Wait for developer list to load
  12 |     await page.waitForTimeout(1000)
  13 |     // Click first developer
  14 |     const firstDev = page.locator('.developer-card, .dev-card, [class*="developer"]').first()
  15 |     if (await firstDev.isVisible()) {
  16 |       await firstDev.click()
  17 |       await expect(page).toHaveURL(/\/developers\//, { timeout: 5000 })
  18 |     }
  19 |   })
  20 | 
  21 |   test('should show developer profile after login', async ({ developerPage }) => {
  22 |     await developerPage.goto('/profile')
> 23 |     await expect(developerPage.getByText(/赵全栈|开发者|profile/i)).toBeVisible({ timeout: 5000 })
     |                                                               ^ Error: expect(locator).toBeVisible() failed
  24 |   })
  25 | 
  26 |   test('should navigate to developers from nav', async ({ page }) => {
  27 |     await page.goto('/')
  28 |     const devLink = page.getByRole('link', { name: /找开发者|开发者/ })
  29 |     if (await devLink.isVisible()) {
  30 |       await devLink.click()
  31 |       await expect(page).toHaveURL(/\/developers/)
  32 |     }
  33 |   })
  34 | })
  35 | 
```
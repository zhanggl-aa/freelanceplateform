# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: contract.spec.ts >> Contract Management >> should show developer contracts
- Location: tests\contract.spec.ts:42:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/我的合同|合同/)
Expected: visible
Error: strict mode violation: getByText(/我的合同|合同/) resolved to 3 elements:
    1) <h2 data-v-fe020488="">我的合同</h2> aka getByRole('heading', { name: '我的合同' })
    2) <p data-v-fe020488="">查看和管理所有合同</p> aka getByText('查看和管理所有合同')
    3) <li tabindex="-1" role="menuitem" aria-disabled="false" data-el-collection-item="" class="el-dropdown-menu__item">…</li> aka getByLabel('赵全栈').getByText('我的合同')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/我的合同|合同/)

```

# Page snapshot

```yaml
- generic [active] [ref=e1]:
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
          - heading "我的合同" [level=2] [ref=e44]
          - paragraph [ref=e45]: 查看和管理所有合同
        - generic [ref=e46]:
          - tablist [ref=e50]:
            - tab "进行中" [selected] [ref=e52]
            - tab "已完成" [ref=e53]
            - tab "已取消" [ref=e54]
            - tab "有争议" [ref=e55]
          - generic:
            - tabpanel "进行中"
        - generic [ref=e56] [cursor=pointer]:
          - generic [ref=e57]:
            - table [ref=e59]:
              - rowgroup [ref=e66]:
                - row "项目名称 对方 金额 状态 开始日期" [ref=e67]:
                  - columnheader "项目名称" [ref=e68]:
                    - generic [ref=e69]: 项目名称
                  - columnheader "对方" [ref=e70]:
                    - generic [ref=e71]: 对方
                  - columnheader "金额" [ref=e72]:
                    - generic [ref=e73]: 金额
                  - columnheader "状态" [ref=e74]:
                    - generic [ref=e75]: 状态
                  - columnheader "开始日期" [ref=e76]:
                    - generic [ref=e77]: 开始日期
            - generic [ref=e81]:
              - table:
                - rowgroup
              - generic [ref=e83]: 暂无数据
          - img [ref=e86]
        - generic [ref=e88]:
          - img [ref=e90]
          - paragraph [ref=e107]: 暂无合同
    - contentinfo [ref=e108]:
      - paragraph [ref=e110]: © 2026 接单平台 - 专业的软件开发者自由职业平台
  - alert [ref=e111]:
    - img [ref=e113]
    - paragraph [ref=e115]: "Key: 'Role' Error:Field validation for 'Role' failed on the 'required' tag"
```

# Test source

```ts
  1  | import { test, expect } from './helpers/fixtures'
  2  | 
  3  | test.describe('Contract Management', () => {
  4  |   test('should view my contracts list', async ({ clientPage }) => {
  5  |     await clientPage.goto('/my/contracts')
  6  |     await expect(clientPage.getByText(/我的合同|合同/)).toBeVisible({ timeout: 5000 })
  7  |   })
  8  | 
  9  |   test('should show contract tabs', async ({ clientPage }) => {
  10 |     await clientPage.goto('/my/contracts')
  11 |     await expect(clientPage.getByText('进行中')).toBeVisible({ timeout: 5000 })
  12 |     await expect(clientPage.getByText('已完成')).toBeVisible()
  13 |     await expect(clientPage.getByText('已取消')).toBeVisible()
  14 |   })
  15 | 
  16 |   test('should switch to completed contracts tab', async ({ clientPage }) => {
  17 |     await clientPage.goto('/my/contracts')
  18 |     await clientPage.getByText('已完成').click()
  19 |     await clientPage.waitForTimeout(500)
  20 |   })
  21 | 
  22 |   test('should switch to cancelled contracts tab', async ({ clientPage }) => {
  23 |     await clientPage.goto('/my/contracts')
  24 |     await clientPage.getByText('已取消').click()
  25 |     await clientPage.waitForTimeout(500)
  26 |   })
  27 | 
  28 |   test('should view contract detail', async ({ clientPage }) => {
  29 |     await clientPage.goto('/my/contracts')
  30 |     await clientPage.waitForTimeout(1000)
  31 | 
  32 |     const contractRow = clientPage.locator('.el-table__row').first()
  33 |     if (await contractRow.isVisible()) {
  34 |       await contractRow.click()
  35 |       await clientPage.waitForTimeout(1000)
  36 |       // Should be on contract detail page
  37 |       const url = clientPage.url()
  38 |       expect(url).toMatch(/\/contracts\//)
  39 |     }
  40 |   })
  41 | 
  42 |   test('should show developer contracts', async ({ developerPage }) => {
  43 |     await developerPage.goto('/my/contracts')
> 44 |     await expect(developerPage.getByText(/我的合同|合同/)).toBeVisible({ timeout: 5000 })
     |                                                      ^ Error: expect(locator).toBeVisible() failed
  45 |   })
  46 | 
  47 |   test('should display contract amounts', async ({ clientPage }) => {
  48 |     await clientPage.goto('/my/contracts')
  49 |     await clientPage.waitForTimeout(1000)
  50 |     // Should show amount column with ¥ symbol
  51 |     const amountCells = clientPage.locator('.amount-text, :text("¥")')
  52 |     if (await amountCells.first().isVisible()) {
  53 |       await expect(amountCells.first()).toContainText('¥')
  54 |     }
  55 |   })
  56 | })
  57 | 
```
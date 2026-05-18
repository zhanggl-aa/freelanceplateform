# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: bid.spec.ts >> Bidding Flow >> should submit bid via API and see it in my bids
- Location: tests\bid.spec.ts:23:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/我的投标|投标/)
Expected: visible
Error: strict mode violation: getByText(/我的投标|投标/) resolved to 3 elements:
    1) <h2 data-v-e524f47e="">我的投标</h2> aka getByRole('heading', { name: '我的投标' })
    2) <p data-v-e524f47e="">查看和管理你提交的所有投标</p> aka getByText('查看和管理你提交的所有投标')
    3) <li tabindex="-1" role="menuitem" aria-disabled="false" data-el-collection-item="" class="el-dropdown-menu__item">…</li> aka getByLabel('赵全栈').getByText('我的投标')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/我的投标|投标/)

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
        - heading "我的投标" [level=2] [ref=e44]
        - paragraph [ref=e45]: 查看和管理你提交的所有投标
      - generic [ref=e46]:
        - tablist [ref=e50]:
          - tab "全部" [selected] [ref=e52]
          - tab "已提交" [ref=e53]
          - tab "已入围" [ref=e54]
          - tab "已接受" [ref=e55]
          - tab "已拒绝" [ref=e56]
        - generic:
          - tabpanel "全部"
      - generic [ref=e58] [cursor=pointer]:
        - table [ref=e60]:
          - rowgroup [ref=e67]:
            - row "项目名称 报价金额 预计天数 状态 提交时间" [ref=e68]:
              - columnheader "项目名称" [ref=e69]:
                - generic [ref=e70]: 项目名称
              - columnheader "报价金额" [ref=e71]:
                - generic [ref=e72]: 报价金额
              - columnheader "预计天数" [ref=e73]:
                - generic [ref=e74]: 预计天数
              - columnheader "状态" [ref=e75]:
                - generic [ref=e76]: 状态
              - columnheader "提交时间" [ref=e77]:
                - generic [ref=e78]: 提交时间
        - table [ref=e83]:
          - rowgroup [ref=e90]:
            - row "未命名项目 ¥42000 45天 已提交 2026-05-16 15:16" [ref=e91]:
              - cell "未命名项目" [ref=e92]:
                - generic [ref=e93]: 未命名项目
              - cell "¥42000" [ref=e94]:
                - generic [ref=e95]: ¥42000
              - cell "45天" [ref=e96]:
                - generic [ref=e97]: 45天
              - cell "已提交" [ref=e98]:
                - generic [ref=e100]: 已提交
              - cell "2026-05-16 15:16" [ref=e101]:
                - generic [ref=e102]: 2026-05-16 15:16
            - row "未命名项目 ¥90000 50天 已提交 2026-05-16 15:16" [ref=e103]:
              - cell "未命名项目" [ref=e104]:
                - generic [ref=e105]: 未命名项目
              - cell "¥90000" [ref=e106]:
                - generic [ref=e107]: ¥90000
              - cell "50天" [ref=e108]:
                - generic [ref=e109]: 50天
              - cell "已提交" [ref=e110]:
                - generic [ref=e112]: 已提交
              - cell "2026-05-16 15:16" [ref=e113]:
                - generic [ref=e114]: 2026-05-16 15:16
            - row "未命名项目 ¥30000 20天 已入围 2026-05-16 15:16" [ref=e115]:
              - cell "未命名项目" [ref=e116]:
                - generic [ref=e117]: 未命名项目
              - cell "¥30000" [ref=e118]:
                - generic [ref=e119]: ¥30000
              - cell "20天" [ref=e120]:
                - generic [ref=e121]: 20天
              - cell "已入围" [ref=e122]:
                - generic [ref=e124]: 已入围
              - cell "2026-05-16 15:16" [ref=e125]:
                - generic [ref=e126]: 2026-05-16 15:16
            - row "未命名项目 ¥120000 90天 已接受 2026-05-16 15:16" [ref=e127]:
              - cell "未命名项目" [ref=e128]:
                - generic [ref=e129]: 未命名项目
              - cell "¥120000" [ref=e130]:
                - generic [ref=e131]: ¥120000
              - cell "90天" [ref=e132]:
                - generic [ref=e133]: 90天
              - cell "已接受" [ref=e134]:
                - generic [ref=e136]: 已接受
              - cell "2026-05-16 15:16" [ref=e137]:
                - generic [ref=e138]: 2026-05-16 15:16
  - contentinfo [ref=e139]:
    - paragraph [ref=e141]: © 2026 接单平台 - 专业的软件开发者自由职业平台
```

# Test source

```ts
  1  | import { test, expect } from './helpers/fixtures'
  2  | import { apiCreateBid, apiGetCategories, apiCreateProject, apiPublishProject, TEST_USERS, apiLogin } from './helpers/api'
  3  | 
  4  | test.describe('Bidding Flow', () => {
  5  |   test('should show bid button on published project for developer', async ({ developerPage }) => {
  6  |     await developerPage.goto('/projects')
  7  |     await developerPage.waitForTimeout(1000)
  8  | 
  9  |     // Click on a published project
  10 |     const projectCards = developerPage.locator('.project-card, [class*="project-item"], .el-card')
  11 |     if (await projectCards.first().isVisible()) {
  12 |       await projectCards.first().click()
  13 |       await developerPage.waitForTimeout(1000)
  14 | 
  15 |       // Should see bid/apply button
  16 |       const bidBtn = developerPage.getByRole('button', { name: /投标|申请|报价/ })
  17 |       if (await bidBtn.isVisible()) {
  18 |         await expect(bidBtn).toBeEnabled()
  19 |       }
  20 |     }
  21 |   })
  22 | 
  23 |   test('should submit bid via API and see it in my bids', async ({ developerPage }) => {
  24 |     // Login as dev1 and check my bids page
  25 |     await developerPage.goto('/my/bids')
  26 |     await developerPage.waitForTimeout(1000)
  27 |     // Page should load without error
> 28 |     await expect(developerPage.getByText(/我的投标|投标/)).toBeVisible({ timeout: 5000 })
     |                                                      ^ Error: expect(locator).toBeVisible() failed
  29 |   })
  30 | 
  31 |   test('should view bids on own project as client', async ({ clientPage }) => {
  32 |     await clientPage.goto('/my/projects')
  33 |     await clientPage.waitForTimeout(1000)
  34 | 
  35 |     // Click on a project that has bids
  36 |     const projectRow = clientPage.locator('.el-table__row, [class*="project"]').first()
  37 |     if (await projectRow.isVisible()) {
  38 |       await projectRow.click()
  39 |       await clientPage.waitForTimeout(1000)
  40 | 
  41 |       // Look for bids section
  42 |       const bidsTab = clientPage.getByText(/投标|竞标/)
  43 |       if (await bidsTab.isVisible()) {
  44 |         await bidsTab.click()
  45 |         await clientPage.waitForTimeout(500)
  46 |       }
  47 |     }
  48 |   })
  49 | 
  50 |   test('should accept a bid as client', async ({ clientPage }) => {
  51 |     // Navigate to a project with bids
  52 |     await clientPage.goto('/projects')
  53 |     await clientPage.waitForTimeout(1000)
  54 | 
  55 |     const projectCards = clientPage.locator('.project-card, [class*="project-item"], .el-card')
  56 |     if (await projectCards.first().isVisible()) {
  57 |       await projectCards.first().click()
  58 |       await clientPage.waitForTimeout(1000)
  59 | 
  60 |       // Look for bid list and accept button
  61 |       const acceptBtn = clientPage.getByRole('button', { name: /接受|采纳/ }).first()
  62 |       if (await acceptBtn.isVisible()) {
  63 |         await acceptBtn.click()
  64 |         await clientPage.waitForTimeout(1000)
  65 |       }
  66 |     }
  67 |   })
  68 | 
  69 |   test('should reject a bid as client', async ({ clientPage }) => {
  70 |     await clientPage.goto('/projects')
  71 |     await clientPage.waitForTimeout(1000)
  72 | 
  73 |     const projectCards = clientPage.locator('.project-card, [class*="project-item"], .el-card')
  74 |     if (await projectCards.first().isVisible()) {
  75 |       await projectCards.first().click()
  76 |       await clientPage.waitForTimeout(1000)
  77 | 
  78 |       const rejectBtn = clientPage.getByRole('button', { name: /拒绝|驳回/ }).first()
  79 |       if (await rejectBtn.isVisible()) {
  80 |         await rejectBtn.click()
  81 |         await clientPage.waitForTimeout(1000)
  82 |       }
  83 |     }
  84 |   })
  85 | })
  86 | 
```
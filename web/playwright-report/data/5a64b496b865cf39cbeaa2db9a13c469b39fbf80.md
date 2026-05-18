# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: bid.spec.ts >> Bidding Flow >> should view bids on own project as client
- Location: tests\bid.spec.ts:31:3

# Error details

```
Error: locator.isVisible: Error: strict mode violation: getByText(/投标|竞标/) resolved to 10 elements:
    1) <div role="tab" tabindex="-1" id="tab-bidding" aria-selected="false" class="el-tabs__item is-top" aria-controls="pane-bidding">…</div> aka getByRole('tab', { name: '竞标中' })
    2) <span data-v-947bc7dc="">2 人投标</span> aka getByText('2 人投标')
    3) <span class="">查看投标</span> aka getByRole('button', { name: '查看投标' }).first()
    4) <span data-v-947bc7dc="">3 人投标</span> aka getByText('3 人投标')
    5) <span class="">查看投标</span> aka getByRole('button', { name: '查看投标' }).nth(1)
    6) <span data-v-947bc7dc="" class="status-badge status-bidding">竞标中</span> aka locator('span').filter({ hasText: '竞标中' })
    7) <span data-v-947bc7dc="">5 人投标</span> aka getByText('5 人投标')
    8) <span class="">查看投标</span> aka getByRole('button', { name: '查看投标' }).nth(2)
    9) <span data-v-947bc7dc="">1 人投标</span> aka getByText('人投标').nth(3)
    10) <span data-v-947bc7dc="">1 人投标</span> aka getByText('人投标').nth(4)

Call log:
    - checking visibility of getByText(/投标|竞标/)

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
        - button "张总" [ref=e31] [cursor=pointer]:
          - img [ref=e34]
          - generic [ref=e36]: 张总
          - img [ref=e38]
  - main [ref=e41]:
    - generic [ref=e42]:
      - generic [ref=e43]:
        - generic [ref=e44]:
          - heading "我的项目" [level=2] [ref=e45]
          - paragraph [ref=e46]: 管理你发布的所有项目
        - link "发布项目" [ref=e47] [cursor=pointer]:
          - /url: /projects/create
          - button "发布项目" [ref=e48]:
            - generic [ref=e49]:
              - img [ref=e51]
              - text: 发布项目
      - generic [ref=e53]:
        - tablist [ref=e57]:
          - tab "已发布" [selected] [ref=e59]
          - tab "竞标中" [ref=e60]
          - tab "进行中" [ref=e61]
          - tab "已完成" [ref=e62]
        - generic:
          - tabpanel "已发布"
      - generic [ref=e63]:
        - generic [ref=e66]:
          - generic [ref=e67]:
            - generic [ref=e68]:
              - heading "微服务API网关开发" [level=3] [ref=e69] [cursor=pointer]
              - generic [ref=e70]: 已发布
            - paragraph [ref=e71]: 开发API网关服务，支持限流、熔断、鉴权、日志等功能，需对接现有微服务集群。
            - generic [ref=e72]:
              - generic [ref=e74]: 后端开发
              - generic [ref=e76]: 固定价
          - generic [ref=e77]:
            - generic [ref=e78]:
              - img [ref=e80]
              - generic [ref=e84]: ¥40000-60000
            - generic [ref=e85]:
              - img [ref=e87]
              - generic [ref=e89]: 2 人投标
            - generic [ref=e90]:
              - img [ref=e92]
              - generic [ref=e96]: 2026-05-13
          - generic [ref=e97]:
            - button "查看投标" [ref=e98] [cursor=pointer]:
              - generic [ref=e99]: 查看投标
            - button "查看详情" [ref=e100] [cursor=pointer]:
              - generic [ref=e101]: 查看详情
        - generic [ref=e104]:
          - generic [ref=e105]:
            - generic [ref=e106]:
              - heading "企业官网全栈重构" [level=3] [ref=e107] [cursor=pointer]
              - generic [ref=e108]: 已发布
            - paragraph [ref=e109]: 现有企业官网需要全面重构，从WordPress迁移到Next.js + Go技术栈，要求SEO友好、响应式设计、后台管理系统。
            - generic [ref=e110]:
              - generic [ref=e112]: 网站开发
              - generic [ref=e114]: 固定价
          - generic [ref=e115]:
            - generic [ref=e116]:
              - img [ref=e118]
              - generic [ref=e122]: ¥30000-50000
            - generic [ref=e123]:
              - img [ref=e125]
              - generic [ref=e127]: 3 人投标
            - generic [ref=e128]:
              - img [ref=e130]
              - generic [ref=e134]: 2026-05-11
          - generic [ref=e135]:
            - button "查看投标" [ref=e136] [cursor=pointer]:
              - generic [ref=e137]: 查看投标
            - button "查看详情" [ref=e138] [cursor=pointer]:
              - generic [ref=e139]: 查看详情
        - generic [ref=e142]:
          - generic [ref=e143]:
            - generic [ref=e144]:
              - heading "后台管理系统前端开发" [level=3] [ref=e145] [cursor=pointer]
              - generic [ref=e146]: 竞标中
            - paragraph [ref=e147]: 使用Vue3 + Element Plus开发后台管理系统，约30个页面，含权限管理、数据可视化。
            - generic [ref=e148]:
              - generic [ref=e150]: 前端开发
              - generic [ref=e152]: 固定价
          - generic [ref=e153]:
            - generic [ref=e154]:
              - img [ref=e156]
              - generic [ref=e160]: ¥20000-35000
            - generic [ref=e161]:
              - img [ref=e163]
              - generic [ref=e165]: 5 人投标
            - generic [ref=e166]:
              - img [ref=e168]
              - generic [ref=e172]: 2026-05-06
          - generic [ref=e173]:
            - button "查看投标" [ref=e174] [cursor=pointer]:
              - generic [ref=e175]: 查看投标
            - button "查看详情" [ref=e176] [cursor=pointer]:
              - generic [ref=e177]: 查看详情
        - generic [ref=e180]:
          - generic [ref=e181]:
            - generic [ref=e182]:
              - heading "在线教育平台开发" [level=3] [ref=e183] [cursor=pointer]
              - generic [ref=e184]: in_progress
            - paragraph [ref=e185]: 开发在线教育平台，支持直播授课、录播回放、作业批改、学习进度跟踪等功能。
            - generic [ref=e186]:
              - generic [ref=e188]: 网站开发
              - generic [ref=e190]: 固定价
          - generic [ref=e191]:
            - generic [ref=e192]:
              - img [ref=e194]
              - generic [ref=e198]: ¥100000-150000
            - generic [ref=e199]:
              - img [ref=e201]
              - generic [ref=e203]: 1 人投标
            - generic [ref=e204]:
              - img [ref=e206]
              - generic [ref=e210]: 2026-04-01
          - button "查看详情" [ref=e212] [cursor=pointer]:
            - generic [ref=e213]: 查看详情
        - generic [ref=e216]:
          - generic [ref=e217]:
            - generic [ref=e218]:
              - heading "社区团购小程序" [level=3] [ref=e219] [cursor=pointer]
              - generic [ref=e220]: 已完成
            - paragraph [ref=e221]: 开发社区团购微信小程序，含商品管理、订单系统、团长管理、支付对接。
            - generic [ref=e222]:
              - generic [ref=e224]: 小程序开发
              - generic [ref=e226]: 固定价
          - generic [ref=e227]:
            - generic [ref=e228]:
              - img [ref=e230]
              - generic [ref=e234]: ¥25000-40000
            - generic [ref=e235]:
              - img [ref=e237]
              - generic [ref=e239]: 1 人投标
            - generic [ref=e240]:
              - img [ref=e242]
              - generic [ref=e246]: 2026-01-16
          - button "查看详情" [ref=e248] [cursor=pointer]:
            - generic [ref=e249]: 查看详情
  - contentinfo [ref=e250]:
    - paragraph [ref=e252]: © 2026 接单平台 - 专业的软件开发者自由职业平台
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
  28 |     await expect(developerPage.getByText(/我的投标|投标/)).toBeVisible({ timeout: 5000 })
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
> 43 |       if (await bidsTab.isVisible()) {
     |                         ^ Error: locator.isVisible: Error: strict mode violation: getByText(/投标|竞标/) resolved to 10 elements:
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
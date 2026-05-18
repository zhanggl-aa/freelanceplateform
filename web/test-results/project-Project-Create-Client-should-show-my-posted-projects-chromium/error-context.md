# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: project.spec.ts >> Project Create (Client) >> should show my posted projects
- Location: tests\project.spec.ts:75:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/我的项目|项目/)
Expected: visible
Error: strict mode violation: getByText(/我的项目|项目/) resolved to 5 elements:
    1) <a href="/projects" class="nav-link" data-v-8da451f1="">找项目</a> aka getByRole('link', { name: '找项目' })
    2) <h2 data-v-947bc7dc="">我的项目</h2> aka getByRole('heading', { name: '我的项目' })
    3) <p data-v-947bc7dc="">管理你发布的所有项目</p> aka getByText('管理你发布的所有项目')
    4) <span class="">…</span> aka getByRole('button', { name: '发布项目' })
    5) <li tabindex="-1" role="menuitem" aria-disabled="false" data-el-collection-item="" class="el-dropdown-menu__item">…</li> aka getByLabel('张总').getByText('我的项目')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/我的项目|项目/)

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
  2  | import { apiGetCategories, apiCreateProject, apiPublishProject } from './helpers/api'
  3  | 
  4  | test.describe('Project List & Browse', () => {
  5  |   test('should browse project list as public', async ({ page }) => {
  6  |     await page.goto('/projects')
  7  |     await expect(page.getByText(/找项目|项目/)).toBeVisible({ timeout: 5000 })
  8  |   })
  9  | 
  10 |   test('should show project cards on list page', async ({ page }) => {
  11 |     await page.goto('/projects')
  12 |     await page.waitForTimeout(1000)
  13 |     // Should have at least some project content
  14 |     const projectCards = page.locator('.project-card, [class*="project-item"], .el-card')
  15 |     await expect(projectCards.first()).toBeVisible({ timeout: 5000 })
  16 |   })
  17 | 
  18 |   test('should filter by project type', async ({ page }) => {
  19 |     await page.goto('/projects')
  20 |     const fixedRadio = page.getByText('固定价')
  21 |     if (await fixedRadio.isVisible()) {
  22 |       await fixedRadio.click()
  23 |       await page.waitForTimeout(500)
  24 |     }
  25 |   })
  26 | 
  27 |   test('should view project detail', async ({ page }) => {
  28 |     await page.goto('/projects')
  29 |     await page.waitForTimeout(1000)
  30 |     const firstProject = page.locator('.project-card, [class*="project-item"], .el-card').first()
  31 |     if (await firstProject.isVisible()) {
  32 |       await firstProject.click()
  33 |       await page.waitForTimeout(1000)
  34 |       // Should be on project detail page
  35 |       const url = page.url()
  36 |       expect(url).toMatch(/\/projects\/[a-f0-9-]+/)
  37 |     }
  38 |   })
  39 | })
  40 | 
  41 | test.describe('Project Create (Client)', () => {
  42 |   test('should navigate to create project page', async ({ clientPage }) => {
  43 |     await clientPage.goto('/projects/create')
  44 |     await expect(clientPage.getByText(/创建项目|发布项目/)).toBeVisible({ timeout: 5000 })
  45 |   })
  46 | 
  47 |   test('should show project form fields', async ({ clientPage }) => {
  48 |     await clientPage.goto('/projects/create')
  49 |     await expect(clientPage.getByText(/项目标题|标题/)).toBeVisible({ timeout: 5000 })
  50 |     await expect(clientPage.getByText(/项目描述|描述/)).toBeVisible({ timeout: 5000 })
  51 |   })
  52 | 
  53 |   test('should create a project via form', async ({ clientPage }) => {
  54 |     await clientPage.goto('/projects/create')
  55 |     await clientPage.waitForTimeout(500)
  56 | 
  57 |     const titleInput = clientPage.getByPlaceholder(/项目标题|标题/)
  58 |     if (await titleInput.isVisible()) {
  59 |       await titleInput.fill('E2E测试项目')
  60 |     }
  61 | 
  62 |     const descInput = clientPage.getByPlaceholder(/项目描述|描述/)
  63 |     if (await descInput.isVisible()) {
  64 |       await descInput.fill('这是一个Playwright自动化测试创建的项目')
  65 |     }
  66 | 
  67 |     // Try to submit
  68 |     const submitBtn = clientPage.getByRole('button', { name: /创建|发布|提交/ })
  69 |     if (await submitBtn.isVisible()) {
  70 |       await submitBtn.click()
  71 |       await clientPage.waitForTimeout(2000)
  72 |     }
  73 |   })
  74 | 
  75 |   test('should show my posted projects', async ({ clientPage }) => {
  76 |     await clientPage.goto('/my/projects')
> 77 |     await expect(clientPage.getByText(/我的项目|项目/)).toBeVisible({ timeout: 5000 })
     |                                                   ^ Error: expect(locator).toBeVisible() failed
  78 |   })
  79 | })
  80 | 
  81 | test.describe('Project Detail', () => {
  82 |   test('should show project info on detail page', async ({ page }) => {
  83 |     await page.goto('/projects')
  84 |     await page.waitForTimeout(1000)
  85 |     const firstProject = page.locator('.project-card, [class*="project-item"], .el-card').first()
  86 |     if (await firstProject.isVisible()) {
  87 |       await firstProject.click()
  88 |       await page.waitForTimeout(1000)
  89 |       // Should show project details
  90 |       await expect(page.getByText(/预算|budget/i)).toBeVisible({ timeout: 5000 })
  91 |     }
  92 |   })
  93 | })
  94 | 
```
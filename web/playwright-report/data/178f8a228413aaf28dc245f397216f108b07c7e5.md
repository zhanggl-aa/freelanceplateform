# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: project.spec.ts >> Project Detail >> should show project info on detail page
- Location: tests\project.spec.ts:82:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/预算|budget/i)
Expected: visible
Error: strict mode violation: getByText(/预算|budget/i) resolved to 2 elements:
    1) <h4 data-v-00cc3a23="">预算范围</h4> aka getByRole('heading', { name: '预算范围' })
    2) <span class="el-radio-button__inner">预算最高</span> aka getByText('预算最高')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/预算|budget/i)

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
        - link "登录" [ref=e17] [cursor=pointer]:
          - /url: /login
          - button "登录" [ref=e18]:
            - generic [ref=e19]: 登录
        - link "注册" [ref=e20] [cursor=pointer]:
          - /url: /register
          - button "注册" [ref=e21]:
            - generic [ref=e22]: 注册
  - main [ref=e24]:
    - generic [ref=e25]:
      - generic [ref=e26]:
        - heading "找项目" [level=2] [ref=e27]
        - paragraph [ref=e28]: 浏览最新项目，找到适合你的机会
      - generic [ref=e29]:
        - complementary [ref=e30]:
          - generic [ref=e32]:
            - heading "筛选条件" [level=3] [ref=e33]
            - generic [ref=e34]:
              - heading "项目分类" [level=4] [ref=e35]
              - tree [ref=e36]:
                - treeitem "网站开发" [expanded] [ref=e37]:
                  - generic [ref=e39] [cursor=pointer]: 网站开发
                  - group
                - treeitem "移动应用" [expanded] [ref=e40]:
                  - generic [ref=e42] [cursor=pointer]: 移动应用
                  - group
                - treeitem "小程序开发" [expanded] [ref=e43]:
                  - generic [ref=e45] [cursor=pointer]: 小程序开发
                  - group
                - treeitem "前端开发" [expanded] [ref=e46]:
                  - generic [ref=e48] [cursor=pointer]: 前端开发
                  - group
                - treeitem "后端开发" [expanded] [ref=e49]:
                  - generic [ref=e51] [cursor=pointer]: 后端开发
                  - group
                - treeitem "AI/机器学习" [expanded] [ref=e52]:
                  - generic [ref=e54] [cursor=pointer]: AI/机器学习
                  - group
                - treeitem "UI/UX设计" [expanded] [ref=e55]:
                  - generic [ref=e57] [cursor=pointer]: UI/UX设计
                  - group
                - treeitem "数据分析" [expanded] [ref=e58]:
                  - generic [ref=e60] [cursor=pointer]: 数据分析
                  - group
                - treeitem "DevOps" [expanded] [ref=e61]:
                  - generic [ref=e63] [cursor=pointer]: DevOps
                  - group
                - treeitem "游戏开发" [expanded] [ref=e64]:
                  - generic [ref=e66] [cursor=pointer]: 游戏开发
                  - group
                - treeitem "区块链" [expanded] [active] [ref=e67]:
                  - generic [ref=e69] [cursor=pointer]: 区块链
                  - group
                - treeitem "其他" [expanded] [ref=e70]:
                  - generic [ref=e72] [cursor=pointer]: 其他
                  - group
            - generic [ref=e73]:
              - heading "预算范围" [level=4] [ref=e74]
              - group "滑块介于 0 至 100000" [ref=e75]:
                - generic [ref=e76] [cursor=pointer]:
                  - slider "选择起始值" [ref=e78]
                  - slider "选择结束值" [ref=e80]
              - generic [ref=e82]:
                - generic [ref=e83]: ¥0
                - generic [ref=e84]: ¥100000
            - generic [ref=e85]:
              - heading "项目类型" [level=4] [ref=e86]
              - radiogroup "radio-group" [ref=e87]:
                - generic [ref=e88] [cursor=pointer]:
                  - radio "全部" [checked] [ref=e90]
                  - generic [ref=e92]: 全部
                - generic [ref=e93] [cursor=pointer]:
                  - radio "固定价" [ref=e95]
                  - generic [ref=e97]: 固定价
                - generic [ref=e98] [cursor=pointer]:
                  - radio "时薪制" [ref=e100]
                  - generic [ref=e102]: 时薪制
            - generic [ref=e103]:
              - heading "技术栈" [level=4] [ref=e104]
              - generic [ref=e106]:
                - generic:
                  - combobox [ref=e108]
                  - generic [ref=e109]: 选择或输入技术
                - img [ref=e112] [cursor=pointer]
            - button "重置筛选" [ref=e114] [cursor=pointer]:
              - generic [ref=e115]: 重置筛选
        - main [ref=e116]:
          - generic [ref=e117]:
            - generic [ref=e118]: 共 0 个项目
            - radiogroup "radio-group" [ref=e119]:
              - generic [ref=e120]:
                - radio "最新发布" [checked] [ref=e121]
                - generic [ref=e122] [cursor=pointer]: 最新发布
              - generic [ref=e123]:
                - radio "预算最高" [ref=e124]
                - generic [ref=e125] [cursor=pointer]: 预算最高
              - generic [ref=e126]:
                - radio "投标最多" [ref=e127]
                - generic [ref=e128] [cursor=pointer]: 投标最多
          - generic [ref=e130]:
            - img [ref=e132]
            - paragraph [ref=e149]: 暂无项目
  - contentinfo [ref=e150]:
    - paragraph [ref=e152]: © 2026 接单平台 - 专业的软件开发者自由职业平台
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
  77 |     await expect(clientPage.getByText(/我的项目|项目/)).toBeVisible({ timeout: 5000 })
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
> 90 |       await expect(page.getByText(/预算|budget/i)).toBeVisible({ timeout: 5000 })
     |                                                  ^ Error: expect(locator).toBeVisible() failed
  91 |     }
  92 |   })
  93 | })
  94 | 
```
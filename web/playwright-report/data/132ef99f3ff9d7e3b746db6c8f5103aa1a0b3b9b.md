# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: project.spec.ts >> Project Create (Client) >> should navigate to create project page
- Location: tests\project.spec.ts:42:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/创建项目|发布项目/)
Expected: visible
Error: strict mode violation: getByText(/创建项目|发布项目/) resolved to 2 elements:
    1) <h2 data-v-36c6fc4c="">发布项目</h2> aka getByRole('heading', { name: '发布项目' })
    2) <span class="">发布项目</span> aka getByRole('button', { name: '发布项目' })

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/创建项目|发布项目/)

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
        - heading "发布项目" [level=2] [ref=e44]
        - paragraph [ref=e45]: 填写项目详情，吸引优秀开发者投标
      - generic [ref=e48]:
        - generic [ref=e49]:
          - generic [ref=e50]: "*项目标题"
          - generic [ref=e53]:
            - textbox "*项目标题" [ref=e54]:
              - /placeholder: 请输入项目标题，简洁明了
            - generic [ref=e57]: 0 / 100
        - generic [ref=e58]:
          - generic [ref=e59]: "*项目分类"
          - generic [ref=e62]:
            - generic:
              - combobox "*项目分类" [ref=e64]
              - generic [ref=e65]: 选择项目分类
            - img [ref=e68] [cursor=pointer]
        - generic [ref=e70]:
          - generic [ref=e71]: "*项目描述"
          - textbox "*项目描述" [ref=e74]:
            - /placeholder: 详细描述项目需求、功能要求、技术要求等
        - generic [ref=e75]:
          - generic [ref=e77]:
            - generic [ref=e78]: "*预算类型"
            - radiogroup "*预算类型" [ref=e80]:
              - generic [ref=e81] [cursor=pointer]:
                - radio "固定价" [checked] [ref=e83]
                - generic [ref=e85]: 固定价
              - generic [ref=e86] [cursor=pointer]:
                - radio "时薪制" [ref=e88]
                - generic [ref=e90]: 时薪制
          - generic [ref=e92]:
            - generic [ref=e93]: 最低预算
            - generic [ref=e95]:
              - button "减少数值" [ref=e96]:
                - img [ref=e98]
              - button "增加数值" [ref=e100] [cursor=pointer]:
                - img [ref=e102]
              - spinbutton "最低预算" [ref=e106]: "0"
          - generic [ref=e108]:
            - generic [ref=e109]: 最高预算
            - generic [ref=e111]:
              - button "减少数值" [ref=e112]:
                - img [ref=e114]
              - button "增加数值" [ref=e116] [cursor=pointer]:
                - img [ref=e118]
              - spinbutton "最高预算" [ref=e122]: "0"
        - generic [ref=e123]:
          - generic [ref=e125]:
            - generic [ref=e126]: "*项目截止日期"
            - generic [ref=e129]:
              - img [ref=e132]
              - combobox "*项目截止日期" [ref=e134]
          - generic [ref=e136]:
            - generic [ref=e137]: 投标截止日期
            - generic [ref=e140]:
              - img [ref=e143]
              - combobox "投标截止日期" [ref=e145]
        - generic [ref=e146]:
          - generic [ref=e147]: 技术栈
          - generic [ref=e150]:
            - generic:
              - combobox "技术栈" [ref=e152]
              - generic [ref=e153]: 选择或输入技术标签
            - img [ref=e156] [cursor=pointer]
        - group "附件上传" [ref=e158]:
          - generic [ref=e159]: 附件上传
          - generic [ref=e161]:
            - button "上传文件" [ref=e162] [cursor=pointer]:
              - button "上传文件" [ref=e163]:
                - generic [ref=e164]:
                  - img [ref=e166]
                  - text: 上传文件
            - generic [ref=e168]: 支持上传文档、图片、压缩包，单个文件不超过20MB
            - list
        - generic [ref=e170]:
          - button "发布项目" [ref=e171] [cursor=pointer]:
            - generic [ref=e172]: 发布项目
          - button "保存草稿" [ref=e173] [cursor=pointer]:
            - generic [ref=e174]: 保存草稿
          - button "取消" [ref=e175] [cursor=pointer]:
            - generic [ref=e176]: 取消
  - contentinfo [ref=e177]:
    - paragraph [ref=e179]: © 2026 接单平台 - 专业的软件开发者自由职业平台
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
> 44 |     await expect(clientPage.getByText(/创建项目|发布项目/)).toBeVisible({ timeout: 5000 })
     |                                                     ^ Error: expect(locator).toBeVisible() failed
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
  90 |       await expect(page.getByText(/预算|budget/i)).toBeVisible({ timeout: 5000 })
  91 |     }
  92 |   })
  93 | })
  94 | 
```
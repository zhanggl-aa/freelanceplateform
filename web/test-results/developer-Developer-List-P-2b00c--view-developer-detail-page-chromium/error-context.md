# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: developer.spec.ts >> Developer List & Profile >> should view developer detail page
- Location: tests\developer.spec.ts:9:3

# Error details

```
Error: expect(page).toHaveURL(expected) failed

Expected pattern: /\/developers\//
Received string:  "http://localhost:3002/developers"
Timeout: 5000ms

Call log:
  - Expect "toHaveURL" with timeout 5000ms
    13 × unexpected value "http://localhost:3002/developers"

```

```yaml
- banner:
  - link "接单平台":
    - /url: /
    - img
    - text: 接单平台
  - navigation:
    - link "首页":
      - /url: /
    - link "找项目":
      - /url: /projects
    - link "找开发者":
      - /url: /developers
  - link "登录":
    - /url: /login
    - button "登录"
  - link "注册":
    - /url: /register
    - button "注册"
- main:
  - heading "找开发者" [level=2]
  - paragraph: 浏览优秀开发者，找到合适的合作伙伴
  - complementary:
    - heading "筛选条件" [level=3]
    - heading "技能搜索" [level=4]
    - img
    - textbox "输入技能关键词"
    - heading "时薪范围（元/时）" [level=4]
    - group "滑块介于 0 至 500":
      - slider "选择起始值"
      - slider "选择结束值"
    - text: ¥0 ¥500
    - heading "可用状态" [level=4]
    - combobox
    - text: 选择可用状态
    - img
    - button "重置筛选"
  - main:
    - text: 共 0 位开发者
    - radiogroup "radio-group":
      - radio "评分最高" [checked]
      - text: 评分最高
      - radio "项目最多"
      - text: 项目最多
      - radio "最新注册"
      - text: 最新注册
    - img
    - heading "吴AI" [level=3]
    - paragraph: AI算法工程师
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000040\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"Python\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000041\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"PyTorch\", \"proficiency\": \"expert\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000042\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"TensorFlow\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000043\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"NLP\", \"proficiency\": \"advanced\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥600/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 6项目
    - button "查看详情"
    - img
    - heading "孙后端" [level=3]
    - paragraph: 后端技术专家
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000020\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Go\", \"proficiency\": \"expert\", \"years_experience\": 8, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000021\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Java\", \"proficiency\": \"advanced\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000022\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Kubernetes\", \"proficiency\": \"expert\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000023\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Redis\", \"proficiency\": \"advanced\", \"years_experience\": 7, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥450/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 22项目
    - button "查看详情"
    - img
    - heading "赵全栈" [level=3]
    - paragraph: 全栈高级工程师
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000001\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"Go\", \"proficiency\": \"expert\", \"years_experience\": 8, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000002\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"React\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000003\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"Vue\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000004\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"Node.js\", \"proficiency\": \"advanced\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥500/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 18项目
    - button "查看详情"
    - img
    - heading "钱前端" [level=3]
    - paragraph: 前端架构师
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000010\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"Vue\", \"proficiency\": \"expert\", \"years_experience\": 7, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000011\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"React\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000012\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"TypeScript\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000013\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"CSS\", \"proficiency\": \"expert\", \"years_experience\": 7, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥400/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 12项目
    - button "查看详情"
    - img
    - heading "冯双修" [level=3]
    - paragraph: DevOps+后端
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000060\", \"developer_id\": \"b0000000-0000-4000-8000-000000000021\", \"skill_name\": \"Docker\", \"proficiency\": \"expert\", \"years_experience\": 8, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000061\", \"developer_id\": \"b0000000-0000-4000-8000-000000000021\", \"skill_name\": \"Kubernetes\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000062\", \"developer_id\": \"b0000000-0000-4000-8000-000000000021\", \"skill_name\": \"Go\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥420/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 10项目
    - button "查看详情"
    - img
    - heading "周移动" [level=3]
    - paragraph: 移动端开发专家
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000030\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"Flutter\", \"proficiency\": \"expert\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000031\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"React Native\", \"proficiency\": \"advanced\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000032\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"Swift\", \"proficiency\": \"intermediate\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000033\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"Kotlin\", \"proficiency\": \"intermediate\", \"years_experience\": 2, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥350/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 8项目
    - button "查看详情"
    - img
    - heading "郑全能" [level=3]
    - paragraph: 全栈+产品经理
    - text: "{ \"id\": \"c0000000-0000-4000-8000-000000000050\", \"developer_id\": \"b0000000-0000-4000-8000-000000000020\", \"skill_name\": \"Vue\", \"proficiency\": \"advanced\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000051\", \"developer_id\": \"b0000000-0000-4000-8000-000000000020\", \"skill_name\": \"Node.js\", \"proficiency\": \"advanced\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } { \"id\": \"c0000000-0000-4000-8000-000000000052\", \"developer_id\": \"b0000000-0000-4000-8000-000000000020\", \"skill_name\": \"Python\", \"proficiency\": \"intermediate\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" } ¥380/时"
    - slider "rating" [disabled]:
      - img
      - img
      - img
      - img
      - img
      - img
    - text: 4项目
    - button "查看详情"
- contentinfo:
  - paragraph: © 2026 接单平台 - 专业的软件开发者自由职业平台
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
> 17 |       await expect(page).toHaveURL(/\/developers\//, { timeout: 5000 })
     |                          ^ Error: expect(page).toHaveURL(expected) failed
  18 |     }
  19 |   })
  20 | 
  21 |   test('should show developer profile after login', async ({ developerPage }) => {
  22 |     await developerPage.goto('/profile')
  23 |     await expect(developerPage.getByText(/赵全栈|开发者|profile/i)).toBeVisible({ timeout: 5000 })
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
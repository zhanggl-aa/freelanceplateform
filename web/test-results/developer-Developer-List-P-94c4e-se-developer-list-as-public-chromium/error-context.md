# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: developer.spec.ts >> Developer List & Profile >> should browse developer list as public
- Location: tests\developer.spec.ts:4:3

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: getByText(/开发者|接单/)
Expected: visible
Error: strict mode violation: getByText(/开发者|接单/) resolved to 6 elements:
    1) <span class="logo-text" data-v-8da451f1="">接单平台</span> aka getByRole('link', { name: '接单平台' })
    2) <a data-v-8da451f1="" href="/developers" aria-current="page" class="router-link-active router-link-exact-active nav-link active">找开发者</a> aka getByRole('link', { name: '找开发者' })
    3) <h2 data-v-f9077165="">找开发者</h2> aka getByRole('heading', { name: '找开发者' })
    4) <p data-v-f9077165="">浏览优秀开发者，找到合适的合作伙伴</p> aka getByText('浏览优秀开发者，找到合适的合作伙伴')
    5) <span data-v-f9077165="" class="result-count">共 0 位开发者</span> aka getByText('共 0 位开发者')
    6) <p data-v-8da451f1="">© 2026 接单平台 - 专业的软件开发者自由职业平台</p> aka getByText('© 2026 接单平台 - 专业的软件开发者自由职业平台')

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for getByText(/开发者|接单/)

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
        - heading "找开发者" [level=2] [ref=e27]
        - paragraph [ref=e28]: 浏览优秀开发者，找到合适的合作伙伴
      - generic [ref=e29]:
        - complementary [ref=e30]:
          - generic [ref=e32]:
            - heading "筛选条件" [level=3] [ref=e33]
            - generic [ref=e34]:
              - heading "技能搜索" [level=4] [ref=e35]
              - generic [ref=e37]:
                - img [ref=e40]
                - textbox "输入技能关键词" [ref=e42]
            - generic [ref=e44]:
              - heading "时薪范围（元/时）" [level=4] [ref=e45]
              - group "滑块介于 0 至 500" [ref=e46]:
                - generic [ref=e47] [cursor=pointer]:
                  - slider "选择起始值" [ref=e49]
                  - slider "选择结束值" [ref=e51]
              - generic [ref=e53]:
                - generic [ref=e54]: ¥0
                - generic [ref=e55]: ¥500
            - generic [ref=e56]:
              - heading "可用状态" [level=4] [ref=e57]
              - generic [ref=e59] [cursor=pointer]:
                - generic:
                  - combobox [ref=e61]
                  - generic [ref=e62]: 选择可用状态
                - img [ref=e65]
            - button "重置筛选" [ref=e67] [cursor=pointer]:
              - generic [ref=e68]: 重置筛选
        - main [ref=e69]:
          - generic [ref=e70]:
            - generic [ref=e71]: 共 0 位开发者
            - radiogroup "radio-group" [ref=e72]:
              - generic [ref=e73]:
                - radio "评分最高" [checked] [ref=e74]
                - generic [ref=e75] [cursor=pointer]: 评分最高
              - generic [ref=e76]:
                - radio "项目最多" [ref=e77]
                - generic [ref=e78] [cursor=pointer]: 项目最多
              - generic [ref=e79]:
                - radio "最新注册" [ref=e80]
                - generic [ref=e81] [cursor=pointer]: 最新注册
          - generic [ref=e82]:
            - generic [ref=e83]:
              - generic [ref=e86] [cursor=pointer]:
                - generic [ref=e87]:
                  - img [ref=e90]
                  - heading "吴AI" [level=3] [ref=e92]
                  - paragraph [ref=e93]: AI算法工程师
                - generic [ref=e94]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000040\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"Python\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000041\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"PyTorch\", \"proficiency\": \"expert\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000042\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"TensorFlow\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000043\", \"developer_id\": \"b0000000-0000-4000-8000-000000000014\", \"skill_name\": \"NLP\", \"proficiency\": \"advanced\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e95]:
                  - generic [ref=e96]: ¥600/时
                  - slider "rating" [disabled] [ref=e98]:
                    - img [ref=e101]
                    - img [ref=e105]
                    - img [ref=e109]
                    - img [ref=e113]
                    - generic [ref=e116]:
                      - img [ref=e117]
                      - img [ref=e120]
                  - generic [ref=e122]: 6项目
                - button "查看详情" [ref=e123]:
                  - generic [ref=e124]: 查看详情
              - generic [ref=e127] [cursor=pointer]:
                - generic [ref=e128]:
                  - img [ref=e131]
                  - heading "孙后端" [level=3] [ref=e133]
                  - paragraph [ref=e134]: 后端技术专家
                - generic [ref=e135]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000020\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Go\", \"proficiency\": \"expert\", \"years_experience\": 8, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000021\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Java\", \"proficiency\": \"advanced\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000022\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Kubernetes\", \"proficiency\": \"expert\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000023\", \"developer_id\": \"b0000000-0000-4000-8000-000000000012\", \"skill_name\": \"Redis\", \"proficiency\": \"advanced\", \"years_experience\": 7, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e136]:
                  - generic [ref=e137]: ¥450/时
                  - slider "rating" [disabled] [ref=e139]:
                    - img [ref=e142]
                    - img [ref=e146]
                    - img [ref=e150]
                    - img [ref=e154]
                    - generic [ref=e157]:
                      - img [ref=e158]
                      - img [ref=e161]
                  - generic [ref=e163]: 22项目
                - button "查看详情" [ref=e164]:
                  - generic [ref=e165]: 查看详情
              - generic [ref=e168] [cursor=pointer]:
                - generic [ref=e169]:
                  - img [ref=e172]
                  - heading "赵全栈" [level=3] [ref=e174]
                  - paragraph [ref=e175]: 全栈高级工程师
                - generic [ref=e176]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000001\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"Go\", \"proficiency\": \"expert\", \"years_experience\": 8, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000002\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"React\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000003\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"Vue\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000004\", \"developer_id\": \"b0000000-0000-4000-8000-000000000010\", \"skill_name\": \"Node.js\", \"proficiency\": \"advanced\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e177]:
                  - generic [ref=e178]: ¥500/时
                  - slider "rating" [disabled] [ref=e180]:
                    - img [ref=e183]
                    - img [ref=e187]
                    - img [ref=e191]
                    - img [ref=e195]
                    - generic [ref=e198]:
                      - img [ref=e199]
                      - img [ref=e202]
                  - generic [ref=e204]: 18项目
                - button "查看详情" [ref=e205]:
                  - generic [ref=e206]: 查看详情
              - generic [ref=e209] [cursor=pointer]:
                - generic [ref=e210]:
                  - img [ref=e213]
                  - heading "钱前端" [level=3] [ref=e215]
                  - paragraph [ref=e216]: 前端架构师
                - generic [ref=e217]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000010\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"Vue\", \"proficiency\": \"expert\", \"years_experience\": 7, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000011\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"React\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000012\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"TypeScript\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000013\", \"developer_id\": \"b0000000-0000-4000-8000-000000000011\", \"skill_name\": \"CSS\", \"proficiency\": \"expert\", \"years_experience\": 7, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e218]:
                  - generic [ref=e219]: ¥400/时
                  - slider "rating" [disabled] [ref=e221]:
                    - img [ref=e224]
                    - img [ref=e228]
                    - img [ref=e232]
                    - img [ref=e236]
                    - generic [ref=e239]:
                      - img [ref=e240]
                      - img [ref=e243]
                  - generic [ref=e245]: 12项目
                - button "查看详情" [ref=e246]:
                  - generic [ref=e247]: 查看详情
              - generic [ref=e250] [cursor=pointer]:
                - generic [ref=e251]:
                  - img [ref=e254]
                  - heading "冯双修" [level=3] [ref=e256]
                  - paragraph [ref=e257]: DevOps+后端
                - generic [ref=e258]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000060\", \"developer_id\": \"b0000000-0000-4000-8000-000000000021\", \"skill_name\": \"Docker\", \"proficiency\": \"expert\", \"years_experience\": 8, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000061\", \"developer_id\": \"b0000000-0000-4000-8000-000000000021\", \"skill_name\": \"Kubernetes\", \"proficiency\": \"expert\", \"years_experience\": 6, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000062\", \"developer_id\": \"b0000000-0000-4000-8000-000000000021\", \"skill_name\": \"Go\", \"proficiency\": \"advanced\", \"years_experience\": 5, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e259]:
                  - generic [ref=e260]: ¥420/时
                  - slider "rating" [disabled] [ref=e262]:
                    - img [ref=e265]
                    - img [ref=e269]
                    - img [ref=e273]
                    - img [ref=e277]
                    - generic [ref=e280]:
                      - img [ref=e281]
                      - img [ref=e284]
                  - generic [ref=e286]: 10项目
                - button "查看详情" [ref=e287]:
                  - generic [ref=e288]: 查看详情
              - generic [ref=e291] [cursor=pointer]:
                - generic [ref=e292]:
                  - img [ref=e295]
                  - heading "周移动" [level=3] [ref=e297]
                  - paragraph [ref=e298]: 移动端开发专家
                - generic [ref=e299]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000030\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"Flutter\", \"proficiency\": \"expert\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000031\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"React Native\", \"proficiency\": \"advanced\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000032\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"Swift\", \"proficiency\": \"intermediate\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000033\", \"developer_id\": \"b0000000-0000-4000-8000-000000000013\", \"skill_name\": \"Kotlin\", \"proficiency\": \"intermediate\", \"years_experience\": 2, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e300]:
                  - generic [ref=e301]: ¥350/时
                  - slider "rating" [disabled] [ref=e303]:
                    - img [ref=e306]
                    - img [ref=e310]
                    - img [ref=e314]
                    - img [ref=e318]
                    - generic [ref=e321]:
                      - img [ref=e322]
                      - img [ref=e325]
                  - generic [ref=e327]: 8项目
                - button "查看详情" [ref=e328]:
                  - generic [ref=e329]: 查看详情
              - generic [ref=e332] [cursor=pointer]:
                - generic [ref=e333]:
                  - img [ref=e336]
                  - heading "郑全能" [level=3] [ref=e338]
                  - paragraph [ref=e339]: 全栈+产品经理
                - generic [ref=e340]:
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000050\", \"developer_id\": \"b0000000-0000-4000-8000-000000000020\", \"skill_name\": \"Vue\", \"proficiency\": \"advanced\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000051\", \"developer_id\": \"b0000000-0000-4000-8000-000000000020\", \"skill_name\": \"Node.js\", \"proficiency\": \"advanced\", \"years_experience\": 4, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                  - generic:
                    - generic: "{ \"id\": \"c0000000-0000-4000-8000-000000000052\", \"developer_id\": \"b0000000-0000-4000-8000-000000000020\", \"skill_name\": \"Python\", \"proficiency\": \"intermediate\", \"years_experience\": 3, \"created_at\": \"2026-05-16T07:16:42.678863Z\" }"
                - generic [ref=e341]:
                  - generic [ref=e342]: ¥380/时
                  - slider "rating" [disabled] [ref=e344]:
                    - img [ref=e347]
                    - img [ref=e351]
                    - img [ref=e355]
                    - img [ref=e359]
                    - generic [ref=e362]:
                      - img [ref=e363]
                      - img [ref=e366]
                  - generic [ref=e368]: 4项目
                - button "查看详情" [ref=e369]:
                  - generic [ref=e370]: 查看详情
            - img [ref=e373]
  - contentinfo [ref=e375]:
    - paragraph [ref=e377]: © 2026 接单平台 - 专业的软件开发者自由职业平台
```

# Test source

```ts
  1  | import { test, expect } from './helpers/fixtures'
  2  | 
  3  | test.describe('Developer List & Profile', () => {
  4  |   test('should browse developer list as public', async ({ page }) => {
  5  |     await page.goto('/developers')
> 6  |     await expect(page.getByText(/开发者|接单/)).toBeVisible({ timeout: 5000 })
     |                                            ^ Error: expect(locator).toBeVisible() failed
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
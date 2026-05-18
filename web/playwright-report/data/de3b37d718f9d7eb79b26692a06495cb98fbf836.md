# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: wallet.spec.ts >> Wallet & Payments >> should display transaction history
- Location: tests\wallet.spec.ts:45:3

# Error details

```
Error: locator.isVisible: Error: strict mode violation: getByText(/交易记录|交易历史|Transaction/) resolved to 2 elements:
    1) <p data-v-d538f058="">管理你的资金和交易记录</p> aka getByText('管理你的资金和交易记录')
    2) <h3 data-v-d538f058="">交易记录</h3> aka getByRole('heading', { name: '交易记录' })

Call log:
    - checking visibility of getByText(/交易记录|交易历史|Transaction/)

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
        - heading "我的钱包" [level=2] [ref=e44]
        - paragraph [ref=e45]: 管理你的资金和交易记录
      - generic [ref=e46]:
        - generic [ref=e50]:
          - img [ref=e53]
          - generic [ref=e57]:
            - generic [ref=e58]: 账户余额
            - generic [ref=e59]: ¥150000
        - generic [ref=e63]:
          - img [ref=e66]
          - generic [ref=e69]:
            - generic [ref=e70]: 冻结金额
            - generic [ref=e71]: ¥20000
      - generic [ref=e72]:
        - button "充值" [ref=e73] [cursor=pointer]:
          - generic [ref=e74]:
            - img [ref=e76]
            - text: 充值
        - button "提现" [ref=e78] [cursor=pointer]:
          - generic [ref=e79]:
            - img [ref=e81]
            - text: 提现
      - generic [ref=e83]:
        - heading "交易记录" [level=3] [ref=e85]
        - generic [ref=e88]:
          - table [ref=e90]:
            - rowgroup [ref=e97]:
              - row "类型 金额 余额 时间 描述" [ref=e98]:
                - columnheader "类型" [ref=e99]:
                  - generic [ref=e100]: 类型
                - columnheader "金额" [ref=e101]:
                  - generic [ref=e102]: 金额
                - columnheader "余额" [ref=e103]:
                  - generic [ref=e104]: 余额
                - columnheader "时间" [ref=e105]:
                  - generic [ref=e106]: 时间
                - columnheader "描述" [ref=e107]:
                  - generic [ref=e108]: 描述
          - table [ref=e113]:
            - rowgroup [ref=e120]:
              - row "托管 +¥200000.00 ¥200000.00 2026-05-16 15:16 充值" [ref=e121]:
                - cell "托管" [ref=e122]:
                  - generic [ref=e125]: 托管
                - cell "+¥200000.00" [ref=e126]:
                  - generic [ref=e127]: +¥200000.00
                - cell "¥200000.00" [ref=e128]:
                  - generic [ref=e129]: ¥200000.00
                - cell "2026-05-16 15:16" [ref=e130]:
                  - generic [ref=e131]: 2026-05-16 15:16
                - cell "充值" [ref=e132]:
                  - generic [ref=e133]: 充值
              - row "escrow_hold ¥15000.00 ¥185000.00 2026-05-16 15:16 托管-在线教育平台-需求确认" [ref=e134]:
                - cell "escrow_hold" [ref=e135]:
                  - generic [ref=e138]: escrow_hold
                - cell "¥15000.00" [ref=e139]:
                  - generic [ref=e140]: ¥15000.00
                - cell "¥185000.00" [ref=e141]:
                  - generic [ref=e142]: ¥185000.00
                - cell "2026-05-16 15:16" [ref=e143]:
                  - generic [ref=e144]: 2026-05-16 15:16
                - cell "托管-在线教育平台-需求确认" [ref=e145]:
                  - generic [ref=e146]: 托管-在线教育平台-需求确认
              - row "escrow_hold ¥45000.00 ¥140000.00 2026-05-16 15:16 托管-在线教育平台-核心功能" [ref=e147]:
                - cell "escrow_hold" [ref=e148]:
                  - generic [ref=e151]: escrow_hold
                - cell "¥45000.00" [ref=e152]:
                  - generic [ref=e153]: ¥45000.00
                - cell "¥140000.00" [ref=e154]:
                  - generic [ref=e155]: ¥140000.00
                - cell "2026-05-16 15:16" [ref=e156]:
                  - generic [ref=e157]: 2026-05-16 15:16
                - cell "托管-在线教育平台-核心功能" [ref=e158]:
                  - generic [ref=e159]: 托管-在线教育平台-核心功能
              - row "托管 +¥100000.00 ¥240000.00 2026-05-16 15:16 充值" [ref=e160]:
                - cell "托管" [ref=e161]:
                  - generic [ref=e164]: 托管
                - cell "+¥100000.00" [ref=e165]:
                  - generic [ref=e166]: +¥100000.00
                - cell "¥240000.00" [ref=e167]:
                  - generic [ref=e168]: ¥240000.00
                - cell "2026-05-16 15:16" [ref=e169]:
                  - generic [ref=e170]: 2026-05-16 15:16
                - cell "充值" [ref=e171]:
                  - generic [ref=e172]: 充值
              - row "escrow_hold ¥20000.00 ¥220000.00 2026-05-16 15:16 托管-在线教育平台-测试优化" [ref=e173]:
                - cell "escrow_hold" [ref=e174]:
                  - generic [ref=e177]: escrow_hold
                - cell "¥20000.00" [ref=e178]:
                  - generic [ref=e179]: ¥20000.00
                - cell "¥220000.00" [ref=e180]:
                  - generic [ref=e181]: ¥220000.00
                - cell "2026-05-16 15:16" [ref=e182]:
                  - generic [ref=e183]: 2026-05-16 15:16
                - cell "托管-在线教育平台-测试优化" [ref=e184]:
                  - generic [ref=e185]: 托管-在线教育平台-测试优化
              - row "提现 ¥70000.00 ¥150000.00 2026-05-16 15:16 提现" [ref=e186]:
                - cell "提现" [ref=e187]:
                  - generic [ref=e190]: 提现
                - cell "¥70000.00" [ref=e191]:
                  - generic [ref=e192]: ¥70000.00
                - cell "¥150000.00" [ref=e193]:
                  - generic [ref=e194]: ¥150000.00
                - cell "2026-05-16 15:16" [ref=e195]:
                  - generic [ref=e196]: 2026-05-16 15:16
                - cell "提现" [ref=e197]:
                  - generic [ref=e198]: 提现
  - contentinfo [ref=e199]:
    - paragraph [ref=e201]: © 2026 接单平台 - 专业的软件开发者自由职业平台
```

# Test source

```ts
  1  | import { test, expect } from './helpers/fixtures'
  2  | 
  3  | test.describe('Wallet & Payments', () => {
  4  |   test('should view wallet page', async ({ clientPage }) => {
  5  |     await clientPage.goto('/wallet')
  6  |     await expect(clientPage.getByText(/我的钱包|钱包/)).toBeVisible({ timeout: 5000 })
  7  |   })
  8  | 
  9  |   test('should display balance information', async ({ clientPage }) => {
  10 |     await clientPage.goto('/wallet')
  11 |     await expect(clientPage.getByText('账户余额')).toBeVisible({ timeout: 5000 })
  12 |     await expect(clientPage.getByText('冻结金额')).toBeVisible()
  13 |   })
  14 | 
  15 |   test('should show balance amounts with ¥ symbol', async ({ clientPage }) => {
  16 |     await clientPage.goto('/wallet')
  17 |     await clientPage.waitForTimeout(1000)
  18 |     const balanceValues = clientPage.locator('.balance-value')
  19 |     if (await balanceValues.first().isVisible()) {
  20 |       await expect(balanceValues.first()).toContainText('¥')
  21 |     }
  22 |   })
  23 | 
  24 |   test('should have deposit button', async ({ clientPage }) => {
  25 |     await clientPage.goto('/wallet')
  26 |     await expect(clientPage.getByRole('button', { name: /充值/ })).toBeVisible({ timeout: 5000 })
  27 |   })
  28 | 
  29 |   test('should have withdraw button', async ({ clientPage }) => {
  30 |     await clientPage.goto('/wallet')
  31 |     await expect(clientPage.getByRole('button', { name: /提现/ })).toBeVisible({ timeout: 5000 })
  32 |   })
  33 | 
  34 |   test('should open deposit dialog', async ({ clientPage }) => {
  35 |     await clientPage.goto('/wallet')
  36 |     await clientPage.getByRole('button', { name: /充值/ }).click()
  37 |     await clientPage.waitForTimeout(500)
  38 |     // Should show deposit dialog or form
  39 |     const depositDialog = clientPage.locator('.el-dialog, .el-drawer')
  40 |     if (await depositDialog.isVisible()) {
  41 |       await expect(depositDialog).toBeVisible()
  42 |     }
  43 |   })
  44 | 
  45 |   test('should display transaction history', async ({ clientPage }) => {
  46 |     await clientPage.goto('/wallet')
  47 |     await clientPage.waitForTimeout(1000)
  48 |     // Should show transaction section
  49 |     const transactionSection = clientPage.getByText(/交易记录|交易历史|Transaction/)
> 50 |     if (await transactionSection.isVisible()) {
     |                                  ^ Error: locator.isVisible: Error: strict mode violation: getByText(/交易记录|交易历史|Transaction/) resolved to 2 elements:
  51 |       await expect(transactionSection).toBeVisible()
  52 |     }
  53 |   })
  54 | 
  55 |   test('should show developer wallet with earnings', async ({ developerPage }) => {
  56 |     await developerPage.goto('/wallet')
  57 |     await expect(developerPage.getByText('账户余额')).toBeVisible({ timeout: 5000 })
  58 |     // Developer should have earnings from completed projects
  59 |     const balanceValue = developerPage.locator('.balance-value').first()
  60 |     if (await balanceValue.isVisible()) {
  61 |       const text = await balanceValue.textContent()
  62 |       const amount = parseFloat(text!.replace(/[¥,]/g, ''))
  63 |       expect(amount).toBeGreaterThan(0)
  64 |     }
  65 |   })
  66 | 
  67 |   test('should navigate to wallet from profile menu', async ({ clientPage }) => {
  68 |     await clientPage.goto('/')
  69 |     // Look for wallet link in navigation
  70 |     const walletLink = clientPage.getByRole('link', { name: /钱包/ })
  71 |     if (await walletLink.isVisible()) {
  72 |       await walletLink.click()
  73 |       await expect(clientPage).toHaveURL(/\/wallet/)
  74 |     }
  75 |   })
  76 | })
  77 | 
```
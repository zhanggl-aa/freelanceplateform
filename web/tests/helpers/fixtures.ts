import { test as base, Page } from '@playwright/test'
import { apiLogin, getStorageState, TEST_USERS, type AuthTokens } from './api'

type Fixtures = {
  clientPage: Page
  developerPage: Page
  _clientAuth: AuthTokens
  _developerAuth: AuthTokens
}

export const test = base.extend<Fixtures>({
  _clientAuth: async ({}, use) => {
    const auth = await apiLogin(TEST_USERS.client1.email, TEST_USERS.client1.password)
    await use(auth)
  },

  _developerAuth: async ({}, use) => {
    const auth = await apiLogin(TEST_USERS.dev1.email, TEST_USERS.dev1.password)
    await use(auth)
  },

  clientPage: async ({ browser, _clientAuth }, use) => {
    const context = await browser.newContext({
      storageState: getStorageState(_clientAuth),
    })
    const page = await context.newPage()
    await page.goto('/')
    await use(page)
    await context.close()
  },

  developerPage: async ({ browser, _developerAuth }, use) => {
    const context = await browser.newContext({
      storageState: getStorageState(_developerAuth),
    })
    const page = await context.newPage()
    await page.goto('/')
    await use(page)
    await context.close()
  },
})

export { expect } from '@playwright/test'

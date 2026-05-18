import Taro from '@tarojs/taro'

const BASE_URL = 'https://your-api-domain.com/api/v1'

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  meta?: {
    page: number
    pageSize: number
    total: number
    totalPages: number
  }
}

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  data?: any
  header?: Record<string, string>
  showLoading?: boolean
  loadingText?: string
}

let isRefreshing = false
let pendingRequests: Array<(token: string) => void> = []

function getAccessToken(): string {
  return Taro.getStorageSync('access_token') || ''
}

function getRefreshToken(): string {
  return Taro.getStorageSync('refresh_token') || ''
}

function setTokens(access: string, refresh: string) {
  Taro.setStorageSync('access_token', access)
  Taro.setStorageSync('refresh_token', refresh)
}

function clearTokens() {
  Taro.removeStorageSync('access_token')
  Taro.removeStorageSync('refresh_token')
}

async function refreshAccessToken(): Promise<string> {
  const refreshToken = getRefreshToken()
  if (!refreshToken) {
    throw new Error('No refresh token')
  }

  const res = await Taro.request({
    url: `${BASE_URL}/auth/refresh`,
    method: 'POST',
    data: { refreshToken },
    header: { 'Content-Type': 'application/json' },
  })

  const { code, data } = res.data as ApiResponse
  if (code === 0 && data) {
    setTokens(data.accessToken, data.refreshToken)
    return data.accessToken
  }

  clearTokens()
  throw new Error('Refresh token expired')
}

async function request<T = any>(options: RequestOptions): Promise<ApiResponse<T>> {
  const {
    url,
    method = 'GET',
    data,
    header = {},
    showLoading = false,
    loadingText = '加载中...',
  } = options

  if (showLoading) {
    Taro.showLoading({ title: loadingText, mask: true })
  }

  const token = getAccessToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...header,
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  try {
    const res = await Taro.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: headers,
    })

    if (showLoading) {
      Taro.hideLoading()
    }

    const responseData = res.data as ApiResponse<T>

    if (res.statusCode === 401) {
      if (!isRefreshing) {
        isRefreshing = true
        try {
          const newToken = await refreshAccessToken()
          isRefreshing = false
          pendingRequests.forEach(cb => cb(newToken))
          pendingRequests = []

          headers['Authorization'] = `Bearer ${newToken}`
          const retryRes = await Taro.request({
            url: `${BASE_URL}${url}`,
            method,
            data,
            header: headers,
          })
          return retryRes.data as ApiResponse<T>
        } catch (err) {
          isRefreshing = false
          pendingRequests = []
          clearTokens()
          Taro.navigateTo({ url: '/pages/auth/login' })
          throw new Error('登录已过期，请重新登录')
        }
      } else {
        return new Promise<ApiResponse<T>>((resolve, reject) => {
          pendingRequests.push(async (newToken: string) => {
            headers['Authorization'] = `Bearer ${newToken}`
            try {
              const retryRes = await Taro.request({
                url: `${BASE_URL}${url}`,
                method,
                data,
                header: headers,
              })
              resolve(retryRes.data as ApiResponse<T>)
            } catch (err) {
              reject(err)
            }
          })
        })
      }
    }

    if (res.statusCode >= 400) {
      const errMsg = responseData?.message || '请求失败'
      Taro.showToast({ title: errMsg, icon: 'none' })
      throw new Error(errMsg)
    }

    return responseData
  } catch (error) {
    if (showLoading) {
      Taro.hideLoading()
    }
    if (error instanceof Error && !error.message.includes('登录已过期')) {
      Taro.showToast({ title: '网络异常，请稍后重试', icon: 'none' })
    }
    throw error
  }
}

export const api = {
  // Auth
  auth: {
    wechatLogin: (code: string) =>
      request({ url: '/auth/wechat/login', method: 'POST', data: { code } }),
    bindPhone: (encryptedData: string, iv: string) =>
      request({ url: '/auth/wechat/bind-phone', method: 'POST', data: { encryptedData, iv } }),
    refresh: (refreshToken: string) =>
      request({ url: '/auth/refresh', method: 'POST', data: { refreshToken } }),
    logout: () => request({ url: '/auth/logout', method: 'POST' }),
  },

  // User
  user: {
    getProfile: () => request({ url: '/user/profile' }),
    updateProfile: (data: any) => request({ url: '/user/profile', method: 'PUT', data }),
    uploadAvatar: (filePath: string) => {
      return new Promise((resolve, reject) => {
        const token = getAccessToken()
        const uploadTask = Taro.uploadFile({
          url: `${BASE_URL}/user/avatar`,
          filePath,
          name: 'avatar',
          header: { Authorization: `Bearer ${token}` },
          success: (res) => {
            const data = JSON.parse(res.data)
            resolve(data)
          },
          fail: (err) => reject(err),
        })
        return uploadTask
      })
    },
  },

  // Developer
  developer: {
    getList: (params?: any) => request({ url: '/developers', data: params }),
    getDetail: (id: string) => request({ url: `/developers/${id}` }),
    getMyProfile: () => request({ url: '/developers/me' }),
    updateMyProfile: (data: any) => request({ url: '/developers/me', method: 'PUT', data }),
    getPortfolio: (id: string) => request({ url: `/developers/${id}/portfolio` }),
  },

  // Client
  client: {
    getMyProfile: () => request({ url: '/clients/me' }),
    updateMyProfile: (data: any) => request({ url: '/clients/me', method: 'PUT', data }),
  },

  // Project
  project: {
    getList: (params?: any) => request({ url: '/projects', data: params }),
    getDetail: (id: string) => request({ url: `/projects/${id}` }),
    create: (data: any) => request({ url: '/projects', method: 'POST', data, showLoading: true, loadingText: '发布中...' }),
    update: (id: string, data: any) => request({ url: `/projects/${id}`, method: 'PUT', data }),
    delete: (id: string) => request({ url: `/projects/${id}`, method: 'DELETE' }),
    saveDraft: (data: any) => request({ url: '/projects/draft', method: 'POST', data }),
    getMyProjects: (params?: any) => request({ url: '/projects/my', data: params }),
    getCategories: () => request({ url: '/projects/categories' }),
  },

  // Bid
  bid: {
    create: (data: any) => request({ url: '/bids', method: 'POST', data, showLoading: true, loadingText: '提交中...' }),
    getList: (projectId: string) => request({ url: `/projects/${projectId}/bids` }),
    getMyBids: (params?: any) => request({ url: '/bids/my', data: params }),
    update: (id: string, data: any) => request({ url: `/bids/${id}`, method: 'PUT', data }),
    withdraw: (id: string) => request({ url: `/bids/${id}/withdraw`, method: 'POST' }),
  },

  // Contract
  contract: {
    getList: (params?: any) => request({ url: '/contracts', data: params }),
    getDetail: (id: string) => request({ url: `/contracts/${id}` }),
    accept: (id: string) => request({ url: `/contracts/${id}/accept`, method: 'POST' }),
    complete: (id: string) => request({ url: `/contracts/${id}/complete`, method: 'POST' }),
  },

  // Milestone
  milestone: {
    getList: (contractId: string) => request({ url: `/contracts/${contractId}/milestones` }),
    update: (contractId: string, id: string, data: any) =>
      request({ url: `/contracts/${contractId}/milestones/${id}`, method: 'PUT', data }),
    submit: (contractId: string, id: string) =>
      request({ url: `/contracts/${contractId}/milestones/${id}/submit`, method: 'POST' }),
    approve: (contractId: string, id: string) =>
      request({ url: `/contracts/${contractId}/milestones/${id}/approve`, method: 'POST' }),
  },

  // Payment
  payment: {
    getBalance: () => request({ url: '/payments/balance' }),
    getTransactions: (params?: any) => request({ url: '/payments/transactions', data: params }),
    deposit: (data: any) => request({ url: '/payments/deposit', method: 'POST', data }),
    withdraw: (data: any) => request({ url: '/payments/withdraw', method: 'POST', data, showLoading: true, loadingText: '提现中...' }),
  },

  // Chat
  chat: {
    getConversations: () => request({ url: '/chat/conversations' }),
    getMessages: (conversationId: string, params?: any) =>
      request({ url: `/chat/conversations/${conversationId}/messages`, data: params }),
    sendMessage: (conversationId: string, data: any) =>
      request({ url: `/chat/conversations/${conversationId}/messages`, method: 'POST', data }),
    createConversation: (data: any) =>
      request({ url: '/chat/conversations', method: 'POST', data }),
    markRead: (conversationId: string) =>
      request({ url: `/chat/conversations/${conversationId}/read`, method: 'POST' }),
  },

  // Review
  review: {
    create: (data: any) => request({ url: '/reviews', method: 'POST', data }),
    getList: (userId: string, params?: any) =>
      request({ url: `/users/${userId}/reviews`, data: params }),
  },

  // Notification
  notification: {
    getList: (params?: any) => request({ url: '/notifications', data: params }),
    markRead: (id: string) => request({ url: `/notifications/${id}/read`, method: 'POST' }),
    markAllRead: () => request({ url: '/notifications/read-all', method: 'POST' }),
    getUnreadCount: () => request({ url: '/notifications/unread-count' }),
  },

  // File
  file: {
    upload: (filePath: string, type: 'image' | 'file' = 'image') => {
      return new Promise((resolve, reject) => {
        const token = getAccessToken()
        Taro.uploadFile({
          url: `${BASE_URL}/files/upload`,
          filePath,
          name: 'file',
          header: { Authorization: `Bearer ${token}` },
          formData: { type },
          success: (res) => {
            const data = JSON.parse(res.data)
            resolve(data)
          },
          fail: (err) => reject(err),
        })
      })
    },
  },
}

export { setTokens, clearTokens, getAccessToken }
export type { ApiResponse }

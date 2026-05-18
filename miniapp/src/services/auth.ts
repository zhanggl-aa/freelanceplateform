import Taro from '@tarojs/taro'
import { api, setTokens, clearTokens } from './api'

export async function wxLogin(): Promise<any> {
  try {
    const loginRes = await Taro.login()
    if (!loginRes.code) {
      throw new Error('微信登录失败')
    }

    const res = await api.auth.wechatLogin(loginRes.code)
    if (res.code === 0 && res.data) {
      setTokens(res.data.accessToken, res.data.refreshToken)
      Taro.setStorageSync('user_info', res.data.user)
      return res.data
    }
    throw new Error(res.message || '登录失败')
  } catch (error) {
    throw error
  }
}

export async function bindPhone(encryptedData: string, iv: string): Promise<any> {
  try {
    const res = await api.auth.bindPhone(encryptedData, iv)
    if (res.code === 0 && res.data) {
      if (res.data.accessToken) {
        setTokens(res.data.accessToken, res.data.refreshToken)
      }
      if (res.data.user) {
        Taro.setStorageSync('user_info', res.data.user)
      }
      return res.data
    }
    throw new Error(res.message || '绑定手机号失败')
  } catch (error) {
    throw error
  }
}

export async function silentLogin(): Promise<any> {
  try {
    const refreshToken = Taro.getStorageSync('refresh_token')
    if (refreshToken) {
      try {
        const res = await api.auth.refresh(refreshToken)
        if (res.code === 0 && res.data) {
          setTokens(res.data.accessToken, res.data.refreshToken)
          if (res.data.user) {
            Taro.setStorageSync('user_info', res.data.user)
          }
          return res.data
        }
      } catch (err) {
        // Refresh failed, fall through to wxLogin
      }
    }

    return await wxLogin()
  } catch (error) {
    clearTokens()
    Taro.removeStorageSync('user_info')
    throw error
  }
}

export async function logout(): Promise<void> {
  try {
    await api.auth.logout()
  } catch (err) {
    // Ignore logout API errors
  } finally {
    clearTokens()
    Taro.removeStorageSync('user_info')
    Taro.reLaunch({ url: '/pages/index/index' })
  }
}

export function isLoggedIn(): boolean {
  return !!Taro.getStorageSync('access_token')
}

export function getStoredUser(): any {
  return Taro.getStorageSync('user_info') || null
}

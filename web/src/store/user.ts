import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '@/api/modules'

interface User {
  id: string
  email?: string
  phone?: string
  nickname: string
  avatar_url?: string
  user_type: string
  status: string
  email_verified: boolean
  phone_verified: boolean
}

interface DeveloperProfile {
  id: string
  title?: string
  bio?: string
  hourly_rate?: number
  availability: string
  experience_years: number
  rating_avg: number
  completed_projects: number
}

interface ClientProfile {
  id: string
  company_name?: string
  industry?: string
  verified: boolean
  total_spent: number
  posted_projects: number
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const developerProfile = ref<DeveloperProfile | null>(null)
  const clientProfile = ref<ClientProfile | null>(null)
  const accessToken = ref(localStorage.getItem('access_token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')

  const isLoggedIn = computed(() => !!accessToken.value)
  const isDeveloper = computed(() =>
    user.value?.user_type === 'developer' || user.value?.user_type === 'both'
  )
  const isClient = computed(() =>
    user.value?.user_type === 'client' || user.value?.user_type === 'both'
  )

  async function login(email: string, password: string) {
    const res: any = await authApi.login({ email, password })
    const { access_token, refresh_token } = res.data
    accessToken.value = access_token
    refreshToken.value = refresh_token
    localStorage.setItem('access_token', access_token)
    localStorage.setItem('refresh_token', refresh_token)
    await fetchUser()
  }

  async function register(data: any) {
    const res: any = await authApi.register(data)
    const { access_token, refresh_token } = res.data
    accessToken.value = access_token
    refreshToken.value = refresh_token
    localStorage.setItem('access_token', access_token)
    localStorage.setItem('refresh_token', refresh_token)
    await fetchUser()
  }

  async function logout() {
    try {
      await authApi.logout(refreshToken.value)
    } catch {}
    clearUser()
  }

  async function fetchUser() {
    try {
      const res: any = await userApi.getMe()
      user.value = res.data.user
      developerProfile.value = res.data.developer_profile || res.data.dev_profile || null
      clientProfile.value = res.data.client_profile || null
    } catch {
      clearUser()
    }
  }

  function clearUser() {
    user.value = null
    developerProfile.value = null
    clientProfile.value = null
    accessToken.value = ''
    refreshToken.value = ''
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  return {
    user,
    developerProfile,
    clientProfile,
    accessToken,
    refreshToken,
    isLoggedIn,
    isDeveloper,
    isClient,
    login,
    register,
    logout,
    fetchUser,
    clearUser,
  }
})

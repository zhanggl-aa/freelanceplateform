import { createContext, useContext, useState, useCallback, useEffect, ReactNode } from 'react'
import Taro from '@tarojs/taro'
import { api } from '@/services/api'
import { silentLogin, logout as authLogout, isLoggedIn, getStoredUser } from '@/services/auth'

interface User {
  id: string
  nickname: string
  avatar: string
  phone: string
  userType: 'developer' | 'client' | null
  createdAt: string
}

interface DeveloperProfile {
  id: string
  title: string
  bio: string
  skills: string[]
  hourlyRate: number
  availability: 'available' | 'busy' | 'unavailable'
  rating: number
  completedProjects: number
  portfolio: any[]
}

interface ClientProfile {
  id: string
  company: string
  industry: string
  postedProjects: number
}

interface UserState {
  user: User | null
  developerProfile: DeveloperProfile | null
  clientProfile: ClientProfile | null
  isLoggedIn: boolean
  loading: boolean
}

interface UserContextType extends UserState {
  loginUser: () => Promise<void>
  logoutUser: () => Promise<void>
  fetchUser: () => Promise<void>
  fetchDeveloperProfile: () => Promise<void>
  fetchClientProfile: () => Promise<void>
  setUser: (user: User | null) => void
}

const UserContext = createContext<UserContextType | null>(null)

export function UserProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<UserState>({
    user: null,
    developerProfile: null,
    clientProfile: null,
    isLoggedIn: false,
    loading: true,
  })

  const fetchUser = useCallback(async () => {
    try {
      const res = await api.user.getProfile()
      if (res.code === 0 && res.data) {
        Taro.setStorageSync('user_info', res.data)
        setState(prev => ({
          ...prev,
          user: res.data,
          isLoggedIn: true,
          loading: false,
        }))
      }
    } catch (err) {
      setState(prev => ({ ...prev, loading: false }))
    }
  }, [])

  const fetchDeveloperProfile = useCallback(async () => {
    try {
      const res = await api.developer.getMyProfile()
      if (res.code === 0 && res.data) {
        setState(prev => ({ ...prev, developerProfile: res.data }))
      }
    } catch (err) {
      // ignore
    }
  }, [])

  const fetchClientProfile = useCallback(async () => {
    try {
      const res = await api.client.getMyProfile()
      if (res.code === 0 && res.data) {
        setState(prev => ({ ...prev, clientProfile: res.data }))
      }
    } catch (err) {
      // ignore
    }
  }, [])

  const loginUser = useCallback(async () => {
    try {
      const data = await silentLogin()
      setState(prev => ({
        ...prev,
        user: data?.user || null,
        isLoggedIn: true,
        loading: false,
      }))
      if (data?.user?.userType === 'developer') {
        fetchDeveloperProfile()
      } else if (data?.user?.userType === 'client') {
        fetchClientProfile()
      }
    } catch (err) {
      setState(prev => ({
        ...prev,
        user: null,
        isLoggedIn: false,
        loading: false,
      }))
    }
  }, [fetchDeveloperProfile, fetchClientProfile])

  const logoutUser = useCallback(async () => {
    await authLogout()
    setState({
      user: null,
      developerProfile: null,
      clientProfile: null,
      isLoggedIn: false,
      loading: false,
    })
  }, [])

  const setUser = useCallback((user: User | null) => {
    setState(prev => ({
      ...prev,
      user,
      isLoggedIn: !!user,
    }))
  }, [])

  useEffect(() => {
    const loggedIn = isLoggedIn()
    if (loggedIn) {
      const storedUser = getStoredUser()
      setState(prev => ({
        ...prev,
        user: storedUser,
        isLoggedIn: true,
        loading: false,
      }))
      fetchUser()
    } else {
      setState(prev => ({ ...prev, loading: false }))
    }
  }, [fetchUser])

  return (
    <UserContext.Provider
      value={{
        ...state,
        loginUser,
        logoutUser,
        fetchUser,
        fetchDeveloperProfile,
        fetchClientProfile,
        setUser,
      }}
    >
      {children}
    </UserContext.Provider>
  )
}

export function useUser(): UserContextType {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useUser must be used within a UserProvider')
  }
  return context
}

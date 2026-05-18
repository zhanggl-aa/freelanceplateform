import { View, Text, Image } from '@tarojs/components'
import { useState, useEffect } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import './user.scss'

export default function UserIndex() {
  const [user, setUser] = useState<any>(null)
  const [developerProfile, setDeveloperProfile] = useState<any>(null)
  const [clientProfile, setClientProfile] = useState<any>(null)

  useEffect(() => {
    loadUser()
  }, [])

  const loadUser = async () => {
    try {
      const res = await api.getMe()
      if (res.code === 0) {
        setUser(res.data.user)
        setDeveloperProfile(res.data.developer_profile)
        setClientProfile(res.data.client_profile)
      }
    } catch {}
  }

  const isLoggedIn = !!Taro.getStorageSync('access_token')

  const handleLogin = () => Taro.navigateTo({ url: '/pages/auth/login' })

  const menuItems = [
    { label: '我的档案', url: '/pages/user/profile', icon: '👤' },
    { label: '我的项目', url: '/pages/user/my-projects', icon: '📋' },
    { label: '我的投标', url: '/pages/user/my-bids', icon: '💼' },
    { label: '我的合同', url: '/pages/user/contracts', icon: '📄' },
    { label: '我的钱包', url: '/pages/user/wallet', icon: '💰' },
  ]

  return (
    <View className='user-page'>
      <View className='user-header'>
        {isLoggedIn && user ? (
          <View className='user-info'>
            <Image className='avatar' src={user.avatar_url || 'https://via.placeholder.com/120'} mode='aspectFill' />
            <View className='info'>
              <Text className='nickname'>{user.nickname}</Text>
              <Text className='type'>{user.user_type === 'developer' ? '开发者' : user.user_type === 'client' ? '需求方' : '双重身份'}</Text>
            </View>
          </View>
        ) : (
          <View className='login-prompt' onClick={handleLogin}>
            <Text className='login-text'>点击登录</Text>
          </View>
        )}
      </View>

      {developerProfile && (
        <View className='stats-row'>
          <View className='stat-item'>
            <Text className='stat-value'>{developerProfile.completed_projects}</Text>
            <Text className='stat-label'>完成项目</Text>
          </View>
          <View className='stat-item'>
            <Text className='stat-value'>{developerProfile.rating_avg}</Text>
            <Text className='stat-label'>评分</Text>
          </View>
          <View className='stat-item'>
            <Text className='stat-value'>¥{developerProfile.total_earnings || 0}</Text>
            <Text className='stat-label'>收入</Text>
          </View>
        </View>
      )}

      <View className='menu-list'>
        {menuItems.map(item => (
          <View key={item.url} className='menu-item' onClick={() => Taro.navigateTo({ url: item.url })}>
            <Text className='menu-icon'>{item.icon}</Text>
            <Text className='menu-label'>{item.label}</Text>
            <Text className='menu-arrow'>›</Text>
          </View>
        ))}
      </View>

      {isLoggedIn && (
        <View className='logout-btn' onClick={() => {
          Taro.removeStorageSync('access_token')
          Taro.removeStorageSync('refresh_token')
          setUser(null)
          Taro.reLaunch({ url: '/pages/index/index' })
        }}>
          <Text>退出登录</Text>
        </View>
      )}
    </View>
  )
}

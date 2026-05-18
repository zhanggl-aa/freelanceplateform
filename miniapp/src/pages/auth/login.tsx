import { View, Text, Input, Button, Image } from '@tarojs/components'
import { useState } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import './auth.scss'

export default function Login() {
  const [phone, setPhone] = useState('')
  const [password, setPassword] = useState('')

  const handleWechatLogin = async () => {
    try {
      const { code } = await Taro.login()
      const res = await api.wechatLogin(code)
      if (res.code === 0) {
        Taro.setStorageSync('access_token', res.data.access_token)
        Taro.setStorageSync('refresh_token', res.data.refresh_token)
        Taro.reLaunch({ url: '/pages/index/index' })
      }
    } catch {
      Taro.showToast({ title: '微信登录失败', icon: 'none' })
    }
  }

  const handlePhoneLogin = async () => {
    if (!phone || !password) {
      Taro.showToast({ title: '请填写手机号和密码', icon: 'none' })
      return
    }
    try {
      const res = await api.phoneLogin(phone, password)
      if (res.code === 0) {
        Taro.setStorageSync('access_token', res.data.access_token)
        Taro.setStorageSync('refresh_token', res.data.refresh_token)
        Taro.reLaunch({ url: '/pages/index/index' })
      }
    } catch {
      Taro.showToast({ title: '登录失败', icon: 'none' })
    }
  }

  return (
    <View className='login-page'>
      <View className='login-header'>
        <Text className='login-title'>接单平台</Text>
        <Text className='login-subtitle'>专业程序员接单服务平台</Text>
      </View>

      <View className='login-card'>
        <Button className='wechat-btn' openType='getUserInfo' onGetUserInfo={handleWechatLogin}>
          微信一键登录
        </Button>

        <View className='divider'>
          <View className='divider-line' />
          <Text className='divider-text'>或</Text>
          <View className='divider-line' />
        </View>

        <Input className='input' type='number' placeholder='手机号' value={phone} onInput={e => setPhone(e.detail.value)} />
        <Input className='input' type='safe-password' placeholder='密码' password value={password} onInput={e => setPassword(e.detail.value)} />
        <Button className='login-btn' onClick={handlePhoneLogin}>登录</Button>
      </View>
    </View>
  )
}

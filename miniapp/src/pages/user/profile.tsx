import { View, Text, Input, Textarea, Button, Picker } from '@tarojs/components'
import { useState, useEffect } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import './user.scss'

export default function Profile() {
  const [profile, setProfile] = useState<any>({})
  const [devProfile, setDevProfile] = useState<any>({})
  const [loading, setLoading] = useState(false)

  useEffect(() => { loadData() }, [])

  const loadData = async () => {
    try {
      const res = await api.getMe()
      if (res.code === 0) {
        setProfile(res.data.user || {})
        setDevProfile(res.data.developer_profile || {})
      }
    } catch {}
  }

  const handleSave = async () => {
    setLoading(true)
    try {
      await api.updateProfile({ nickname: profile.nickname })
      if (devProfile.id) {
        await api.updateDeveloperProfile(devProfile)
      }
      Taro.showToast({ title: '保存成功', icon: 'success' })
    } catch {
      Taro.showToast({ title: '保存失败', icon: 'none' })
    }
    setLoading(false)
  }

  return (
    <View className='page-container'>
      <View className='section'>
        <Text className='section-title'>基本信息</Text>
        <View className='form-item'>
          <Text className='label'>昵称</Text>
          <Input className='input' value={profile.nickname} onInput={e => setProfile({...profile, nickname: e.detail.value})} />
        </View>
      </View>

      <View className='section'>
        <Text className='section-title'>开发者档案</Text>
        <View className='form-item'>
          <Text className='label'>职位头衔</Text>
          <Input className='input' placeholder='如：高级全栈工程师' value={devProfile.title || ''} onInput={e => setDevProfile({...devProfile, title: e.detail.value})} />
        </View>
        <View className='form-item'>
          <Text className='label'>简介</Text>
          <Textarea className='textarea' placeholder='介绍您的技能和经验' value={devProfile.bio || ''} onInput={e => setDevProfile({...devProfile, bio: e.detail.value})} />
        </View>
        <View className='form-item'>
          <Text className='label'>时薪 (元)</Text>
          <Input className='input' type='digit' placeholder='0' value={String(devProfile.hourly_rate || '')} onInput={e => setDevProfile({...devProfile, hourly_rate: Number(e.detail.value)})} />
        </View>
        <View className='form-item'>
          <Text className='label'>工作年限</Text>
          <Input className='input' type='number' placeholder='0' value={String(devProfile.experience_years || 0)} onInput={e => setDevProfile({...devProfile, experience_years: Number(e.detail.value)})} />
        </View>
        <View className='form-item'>
          <Text className='label'>状态</Text>
          <Picker mode='selector' range={['空闲', '忙碌', '不可接单']} onChange={e => setDevProfile({...devProfile, availability: ['available','busy','unavailable'][e.detail.value]})}>
            <Text className='picker-value'>{devProfile.availability === 'available' ? '空闲' : devProfile.availability === 'busy' ? '忙碌' : '不可接单'}</Text>
          </Picker>
        </View>
      </View>

      <Button className='save-btn' loading={loading} onClick={handleSave}>保存</Button>
    </View>
  )
}

import { View, Text } from '@tarojs/components'
import { useState, useEffect } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import '../project/project.scss'

export default function MyProjects() {
  const [projects, setProjects] = useState<any[]>([])
  const [tab, setTab] = useState(0)
  const tabs = ['已发布', '进行中', '已完成']
  const statusMap = ['published', 'in_progress', 'completed']

  useEffect(() => { loadProjects() }, [tab])

  const loadProjects = async () => {
    try {
      const res = await api.myPosted({ status: statusMap[tab] })
      setProjects(res.data || [])
    } catch {}
  }

  return (
    <View className='page-container'>
      <View className='tab-bar'>
        {tabs.map((t, i) => (
          <View key={i} className={`tab-item ${tab === i ? 'active' : ''}`} onClick={() => setTab(i)}>
            <Text>{t}</Text>
          </View>
        ))}
      </View>
      {projects.length === 0 ? (
        <View className='empty'><Text>暂无项目</Text></View>
      ) : (
        projects.map(p => (
          <View key={p.id} className='project-card' onClick={() => Taro.navigateTo({ url: `/pages/project/detail?id=${p.id}` })}>
            <Text className='project-title'>{p.title}</Text>
            <View className='project-meta'>
              <Text className='budget'>¥{p.budget_min || 0} - ¥{p.budget_max || 0}</Text>
              <Text className='bids'>{p.bid_count || 0}个投标</Text>
            </View>
          </View>
        ))
      )}
    </View>
  )
}

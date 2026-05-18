import { View, Text } from '@tarojs/components'
import { useState, useEffect } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import '../project/project.scss'

export default function MyBids() {
  const [bids, setBids] = useState<any[]>([])
  const statusLabels: Record<string, string> = { submitted: '已提交', shortlisted: '已入围', accepted: '已接受', rejected: '已拒绝', withdrawn: '已撤回' }

  useEffect(() => { loadBids() }, [])

  const loadBids = async () => {
    try {
      const res = await api.myBids()
      setBids(res.data || [])
    } catch {}
  }

  return (
    <View className='page-container'>
      <Text className='page-title'>我的投标</Text>
      {bids.length === 0 ? (
        <View className='empty'><Text>暂无投标记录</Text></View>
      ) : (
        bids.map(b => (
          <View key={b.id} className='project-card' onClick={() => Taro.navigateTo({ url: `/pages/project/detail?id=${b.project_id}` })}>
            <View className='flex-between'>
              <Text className='project-title'>{b.project_title || '项目'}</Text>
              <Text className={`status ${b.status}`}>{statusLabels[b.status] || b.status}</Text>
            </View>
            <View className='project-meta'>
              <Text>报价: ¥{b.proposed_budget}</Text>
              <Text>预计 {b.estimated_days} 天</Text>
            </View>
          </View>
        ))
      )}
    </View>
  )
}

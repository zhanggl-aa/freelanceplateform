import { View, Text } from '@tarojs/components'
import { useState, useEffect } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import './user.scss'

export default function Contracts() {
  const [contracts, setContracts] = useState<any[]>([])
  const statusLabels: Record<string, string> = { active: '进行中', completed: '已完成', cancelled: '已取消', disputed: '争议中' }

  useEffect(() => { loadContracts() }, [])

  const loadContracts = async () => {
    try {
      const res = await api.myContracts()
      setContracts(res.data || [])
    } catch {}
  }

  return (
    <View className='page-container'>
      <Text className='page-title'>我的合同</Text>
      {contracts.length === 0 ? (
        <View className='empty'><Text>暂无合同</Text></View>
      ) : (
        contracts.map(c => (
          <View key={c.id} className='menu-item' onClick={() => Taro.navigateTo({ url: `/pages/user/contracts?id=${c.id}` })}>
            <View style={{flex:1}}>
              <Text style={{fontSize:'30px',color:'#303133'}}>合同 #{c.id?.slice(0,8)}</Text>
              <Text style={{fontSize:'24px',color:'#909399',marginTop:'8px',display:'block'}}>¥{c.total_amount} · {statusLabels[c.status] || c.status}</Text>
            </View>
            <Text style={{fontSize:'36px',color:'#c0c4cc'}}>›</Text>
          </View>
        ))
      )}
    </View>
  )
}

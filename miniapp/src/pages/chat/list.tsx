import { View, Text, Image, ScrollView } from '@tarojs/components'
import { useState, useEffect } from 'react'
import Taro from '@tarojs/taro'
import api from '../../services/api'
import './chat.scss'

export default function ChatList() {
  const [conversations, setConversations] = useState<any[]>([])

  useEffect(() => { loadConversations() }, [])

  const loadConversations = async () => {
    try {
      const res = await api.listConversations()
      setConversations(res.data || [])
    } catch {}
  }

  Taro.useDidShow(() => { loadConversations() })

  return (
    <View className='chat-list-page'>
      <View className='page-header-bar'>
        <Text className='page-header-title'>消息</Text>
      </View>
      <ScrollView scrollY className='chat-scroll'>
        {conversations.length === 0 ? (
          <View className='empty'><Text>暂无会话</Text></View>
        ) : (
          conversations.map(c => {
            const other = c.participants?.find((p: any) => p.user_id !== Taro.getStorageSync('user_id'))
            return (
              <View key={c.id} className='chat-item' onClick={() => Taro.navigateTo({ url: `/pages/chat/detail?id=${c.id}` })}>
                <Image className='chat-avatar' src={other?.avatar_url || 'https://via.placeholder.com/80'} mode='aspectFill' />
                <View className='chat-info'>
                  <View className='chat-top'>
                    <Text className='chat-name'>{other?.nickname || '用户'}</Text>
                    <Text className='chat-time'>{c.last_message_at?.slice(5,16) || ''}</Text>
                  </View>
                  <Text className='chat-last-msg'>{c.last_message?.content || ''}</Text>
                </View>
              </View>
            )
          })
        )}
      </ScrollView>
    </View>
  )
}

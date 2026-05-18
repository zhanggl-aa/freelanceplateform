import { View, Text, Input, ScrollView } from '@tarojs/components'
import { useState, useEffect, useRef } from 'react'
import Taro, { useRouter } from '@tarojs/taro'
import api from '../../services/api'
import './chat.scss'

export default function ChatDetail() {
  const router = useRouter()
  const convId = router.params.id
  const [messages, setMessages] = useState<any[]>([])
  const [inputText, setInputText] = useState('')
  const [socketTask, setSocketTask] = useState<any>(null)

  useEffect(() => {
    loadMessages()
    connectWS()
    return () => { socketTask?.close() }
  }, [])

  const loadMessages = async () => {
    if (!convId) return
    try {
      const res = await api.getMessages(convId, { page: 1, page_size: 50 })
      setMessages(res.data || [])
    } catch {}
  }

  const connectWS = () => {
    const token = Taro.getStorageSync('access_token')
    if (!token) return
    const task = Taro.connectSocket({ url: `wss://your-api-domain.com/api/v1/ws?token=${token}` })
    task.onMessage((res) => {
      const msg = JSON.parse(res.data)
      if (msg.conversation_id === convId) {
        setMessages(prev => [...prev, msg])
      }
    })
    setSocketTask(task)
  }

  const handleSend = async () => {
    if (!inputText.trim() || !convId) return
    try {
      await api.sendMessage(convId, { content: inputText, message_type: 'text' })
      setMessages(prev => [...prev, {
        id: Date.now().toString(), content: inputText, sender_id: 'self',
        message_type: 'text', created_at: new Date().toISOString(),
      }])
      setInputText('')
    } catch {}
  }

  const myId = Taro.getStorageSync('user_id')

  return (
    <View className='chat-detail-page'>
      <ScrollView scrollY scrollIntoView={messages.length > 0 ? `msg-${messages[messages.length-1].id}` : ''} className='chat-messages'>
        {messages.map(m => (
          <View key={m.id} id={`msg-${m.id}`} className={`msg-row ${m.sender_id === myId ? 'self' : ''}`}>
            <Text className={`msg-bubble ${m.sender_id === myId ? 'self' : 'other'}`}>
              {m.content}
            </Text>
          </View>
        ))}
      </ScrollView>
      <View className='chat-input-bar'>
        <Input className='chat-input' placeholder='输入消息' value={inputText} onInput={e => setInputText(e.detail.value)} confirmType='send' onConfirm={handleSend} />
        <View className='send-btn' onClick={handleSend}><Text>发送</Text></View>
      </View>
    </View>
  )
}

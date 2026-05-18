import { defineStore } from 'pinia'
import { ref } from 'vue'
import { chatApi } from '@/api/modules'

interface Conversation {
  id: string
  type: string
  project_id?: string
  project_title?: string
  unread_count?: number
  last_message_at: string
  participants: any[]
  last_message?: any
}

interface ChatMessage {
  id: string
  conversation_id: string
  sender_id: string
  content?: string
  message_type: string
  file_url?: string
  file_name?: string
  created_at: string
  sender_name?: string
  sender_avatar?: string
}

export const useChatStore = defineStore('chat', () => {
  const conversations = ref<Conversation[]>([])
  const currentMessages = ref<ChatMessage[]>([])
  const unreadCount = ref(0)
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)

  function connectWebSocket(token: string) {
    if (ws.value?.readyState === WebSocket.OPEN) return

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    ws.value = new WebSocket(`${protocol}//${host}/api/v1/ws?token=${token}`)

    ws.value.onopen = () => { connected.value = true }
    ws.value.onclose = () => {
      connected.value = false
      setTimeout(() => connectWebSocket(token), 3000)
    }
    ws.value.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      if (msg.type === 'chat_message') {
        const convMsg: ChatMessage = {
          id: msg.id || Date.now().toString(),
          conversation_id: msg.conversation_id,
          sender_id: msg.sender_id,
          content: msg.content,
          message_type: msg.message_type || 'text',
          file_url: msg.file_url,
          file_name: msg.file_name,
          created_at: new Date().toISOString(),
          sender_name: msg.sender_name,
          sender_avatar: msg.sender_avatar,
        }
        currentMessages.value.push(convMsg)
        unreadCount.value++
      }
    }
  }

  function disconnect() {
    ws.value?.close()
    ws.value = null
    connected.value = false
  }

  function sendMessage(conversationId: string, content: string) {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({
        type: 'chat_message',
        conversation_id: conversationId,
        content,
        message_type: 'text',
      }))
    }
  }

  async function fetchConversations() {
    const res: any = await chatApi.list()
    conversations.value = res.data || []
  }

  async function fetchMessages(conversationId: string, page = 1) {
    const res: any = await chatApi.getMessages(conversationId, { page, page_size: 50 })
    currentMessages.value = res.data || []
  }

  async function fetchUnreadCount() {
    const res: any = await chatApi.unreadCount()
    unreadCount.value = res.data?.count || 0
  }

  return {
    conversations,
    currentMessages,
    unreadCount,
    ws,
    connected,
    connectWebSocket,
    disconnect,
    sendMessage,
    fetchConversations,
    fetchMessages,
    fetchUnreadCount,
  }
})

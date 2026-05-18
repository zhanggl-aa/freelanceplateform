<template>
  <div class="chat-detail-page">
    <!-- Header -->
    <div class="chat-header">
      <el-button text class="back-btn hide-on-desktop" @click="router.push('/chat')">
        <el-icon :size="20"><ArrowLeft /></el-icon>
      </el-button>
      <el-avatar :size="36" :src="otherUser?.avatar_url" icon="UserFilled" />
      <div class="header-info ml-12">
        <span class="header-name">{{ otherUser?.nickname || '对话' }}</span>
        <span class="header-status text-secondary" v-if="conversation?.project_id">
          项目: {{ conversation.project_title || conversation.project_id }}
        </span>
      </div>
    </div>

    <!-- Messages Area -->
    <div class="messages-area" ref="messagesAreaRef" v-loading="messagesLoading">
      <div class="messages-inner">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message-row', { 'self': msg.sender_id === userId, 'other': msg.sender_id !== userId }]"
        >
          <!-- System Message -->
          <div v-if="msg.message_type === 'system'" class="system-message">
            <span>{{ msg.content }}</span>
          </div>

          <!-- Text / File Message -->
          <template v-else>
            <el-avatar
              v-if="msg.sender_id !== userId"
              :size="32"
              :src="msg.sender_avatar || otherUser?.avatar_url"
              icon="UserFilled"
              class="msg-avatar"
            />
            <div class="message-bubble-wrap">
              <div class="message-bubble">
                <!-- Text -->
                <p v-if="msg.message_type === 'text'" class="message-text">{{ msg.content }}</p>
                <!-- File -->
                <div v-else-if="msg.message_type === 'file'" class="message-file" @click="downloadFile(msg)">
                  <el-icon :size="24"><Document /></el-icon>
                  <div class="file-info">
                    <span class="file-name truncate">{{ msg.file_name || '文件' }}</span>
                    <span class="file-hint text-secondary">点击下载</span>
                  </div>
                </div>
                <!-- Image -->
                <el-image
                  v-else-if="msg.message_type === 'image'"
                  :src="msg.file_url"
                  fit="cover"
                  class="message-image"
                  :preview-src-list="msg.file_url ? [msg.file_url] : []"
                />
              </div>
              <span class="message-time text-secondary">{{ formatMessageTime(msg.created_at) }}</span>
            </div>
            <el-avatar
              v-if="msg.sender_id === userId"
              :size="32"
              :src="userStore.user?.avatar_url"
              icon="UserFilled"
              class="msg-avatar"
            />
          </template>
        </div>

        <div v-if="!messagesLoading && messages.length === 0" class="no-messages text-secondary text-center mt-32">
          还没有消息，开始聊天吧
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="input-area">
      <el-upload
        :show-file-list="false"
        :before-upload="beforeFileUpload"
        :http-request="handleFileUpload"
        class="file-upload-btn"
      >
        <el-button text><el-icon :size="20"><Paperclip /></el-icon></el-button>
      </el-upload>
      <el-input
        v-model="inputText"
        placeholder="输入消息..."
        @keyup.enter="handleSend"
        class="message-input"
        size="large"
      />
      <el-button type="primary" @click="handleSend" :disabled="!inputText.trim()" size="large">
        <el-icon><Promotion /></el-icon>
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useChatStore } from '@/store/chat'
import { chatApi, fileApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import type { UploadRequestOptions } from 'element-plus'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const chatStore = useChatStore()

const conversationId = computed(() => route.params.id as string)
const userId = computed(() => userStore.user?.id || '')
const messagesAreaRef = ref<HTMLElement>()
const messagesLoading = ref(false)
const inputText = ref('')
const otherUser = ref<any>(null)
const conversation = ref<any>(null)

const messages = computed(() => chatStore.currentMessages)

function formatMessageTime(time: string) {
  return dayjs(time).format('HH:mm')
}

async function scrollToBottom() {
  await nextTick()
  if (messagesAreaRef.value) {
    messagesAreaRef.value.scrollTop = messagesAreaRef.value.scrollHeight
  }
}

function beforeFileUpload(file: File) {
  const isLt20M = file.size / 1024 / 1024 < 20
  if (!isLt20M) {
    ElMessage.error('文件大小不能超过20MB')
    return false
  }
  return true
}

async function handleFileUpload(options: UploadRequestOptions) {
  try {
    await chatApi.sendFile(conversationId.value, options.file)
  } catch {
    ElMessage.error('文件发送失败')
  }
}

async function downloadFile(msg: any) {
  if (!msg.file_url) return
  try {
    const res: any = await fileApi.download(msg.file_url)
    const blob = new Blob([res as any])
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = msg.file_name || 'download'
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('下载失败')
  }
}

function handleSend() {
  const text = inputText.value.trim()
  if (!text) return
  chatStore.sendMessage(conversationId.value, text)
  // Also send via REST API for persistence
  chatApi.sendMessage(conversationId.value, { content: text, message_type: 'text' }).catch(() => {})
  inputText.value = ''
  scrollToBottom()
}

async function fetchConversation() {
  try {
    const res: any = await chatApi.getMessages(conversationId.value, { page: 1, page_size: 50 })
    conversation.value = res.data?.conversation || res.data
    // Determine other user from conversation
    const participants = conversation.value?.participants || []
    otherUser.value = participants.find((p: any) => p.id !== userId.value) || participants[0]
  } catch {}
}

async function markAsRead() {
  try {
    await chatApi.markAsRead(conversationId.value)
  } catch {}
}

// Watch for new messages and auto-scroll
watch(() => chatStore.currentMessages.length, () => {
  scrollToBottom()
})

onMounted(async () => {
  messagesLoading.value = true
  try {
    await chatStore.fetchMessages(conversationId.value)
    await fetchConversation()
    await markAsRead()
  } catch {} finally {
    messagesLoading.value = false
  }
  scrollToBottom()
})
</script>

<style scoped lang="scss">
.chat-detail-page {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color-light);
  background: var(--bg-color-white);

  .back-btn {
    margin-right: 4px;
  }

  .header-name {
    font-size: 15px;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  .header-status {
    font-size: 12px;
    display: block;
    margin-top: 2px;
  }
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--bg-color);

  .messages-inner {
    max-width: 800px;
    margin: 0 auto;
  }

  .message-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 16px;

    &.self {
      flex-direction: row-reverse;

      .message-bubble {
        background: var(--color-primary);
        color: #fff;
        border-radius: 12px 4px 12px 12px;
      }

      .message-time {
        text-align: right;
      }

      .message-text {
        color: #fff;
      }
    }

    &.other {
      .message-bubble {
        background: var(--bg-color-white);
        border-radius: 4px 12px 12px 12px;
        border: 1px solid var(--border-color-light);
      }
    }

    .msg-avatar {
      flex-shrink: 0;
      margin-top: 4px;
    }

    .message-bubble-wrap {
      max-width: 70%;
      min-width: 60px;
    }

    .message-bubble {
      padding: 10px 14px;
      word-break: break-word;
    }

    .message-text {
      font-size: 14px;
      line-height: 1.5;
      margin: 0;
    }

    .message-time {
      font-size: 11px;
      display: block;
      margin-top: 4px;
      padding: 0 4px;
    }

    .message-file {
      display: flex;
      align-items: center;
      gap: 10px;
      cursor: pointer;

      .file-info {
        display: flex;
        flex-direction: column;
      }

      .file-name {
        font-size: 13px;
        color: var(--color-primary);
        max-width: 200px;
      }

      .file-hint {
        font-size: 11px;
      }
    }

    .message-image {
      max-width: 240px;
      max-height: 200px;
      border-radius: 8px;
    }
  }

  .system-message {
    text-align: center;
    margin: 12px 0;

    span {
      display: inline-block;
      padding: 4px 12px;
      background: var(--border-color-light);
      border-radius: 4px;
      font-size: 12px;
      color: var(--color-text-secondary);
    }
  }

  .no-messages {
    padding: 60px 0;
  }
}

.input-area {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--border-color-light);
  background: var(--bg-color-white);

  .file-upload-btn {
    flex-shrink: 0;
  }

  .message-input {
    flex: 1;
  }
}
</style>

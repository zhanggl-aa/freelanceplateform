<template>
  <div class="chat-list-page page-container">
    <div class="page-header">
      <h2>消息</h2>
      <p>与项目伙伴沟通交流</p>
    </div>

    <!-- Conversation List -->
    <el-card shadow="never" class="chat-list-card">
      <div class="conversation-search mb-16">
        <el-input v-model="searchKeyword" placeholder="搜索对话" clearable>
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </div>

      <div v-loading="loading" class="conversation-list">
        <div
          v-for="conv in filteredConversations"
          :key="conv.id"
          class="conversation-item"
          @click="goChat(conv.id)"
        >
          <el-badge :value="conv.unread_count" :hidden="!conv.unread_count" :max="99">
            <el-avatar :size="48" :src="otherParticipant(conv)?.avatar_url" icon="UserFilled" />
          </el-badge>
          <div class="conv-info">
            <div class="conv-header flex-between">
              <span class="conv-name truncate">{{ otherParticipant(conv)?.nickname || '对话' }}</span>
              <span class="conv-time text-secondary">{{ formatTime(conv.last_message_at) }}</span>
            </div>
            <p class="conv-last-msg text-secondary truncate">
              {{ conv.last_message?.content || '暂无消息' }}
            </p>
            <span v-if="conv.project_title" class="conv-project text-secondary">
              项目: {{ conv.project_title }}
            </span>
          </div>
          <el-icon class="conv-arrow text-secondary"><ArrowRight /></el-icon>
        </div>

        <el-empty v-if="!loading && conversations.length === 0" description="暂无对话" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '@/store/chat'
import { useUserStore } from '@/store/user'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const chatStore = useChatStore()
const userStore = useUserStore()

const loading = ref(false)
const searchKeyword = ref('')

const conversations = computed(() => chatStore.conversations)

const filteredConversations = computed(() => {
  if (!searchKeyword.value) return conversations.value
  const kw = searchKeyword.value.toLowerCase()
  return conversations.value.filter(conv => {
    const other = otherParticipant(conv)
    return other?.nickname?.toLowerCase().includes(kw)
  })
})

function otherParticipant(conv: any) {
  if (!conv.participants || !userStore.user) return null
  return conv.participants.find((p: any) => p.id !== userStore.user?.id) || conv.participants[0]
}

function formatTime(time: string) {
  if (!time) return ''
  const d = dayjs(time)
  const now = dayjs()
  if (d.isSame(now, 'day')) return d.format('HH:mm')
  if (d.isSame(now.subtract(1, 'day'), 'day')) return '昨天'
  if (d.isSame(now, 'year')) return d.format('MM-DD')
  return d.format('YYYY-MM-DD')
}

function goChat(id: string) {
  router.push(`/chat/${id}`)
}

async function fetchConversations() {
  loading.value = true
  try {
    await chatStore.fetchConversations()
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchConversations()
})
</script>

<style scoped lang="scss">
.chat-list-page {
  .chat-list-card {
    border-radius: 12px;
  }

  .conversation-search {
    max-width: 400px;
  }

  .conversation-list {
    .conversation-item {
      display: flex;
      align-items: center;
      gap: 14px;
      padding: 16px 12px;
      cursor: pointer;
      transition: background 0.2s;
      border-bottom: 1px solid var(--border-color-light);
      border-radius: 8px;

      &:hover {
        background: var(--bg-color);
      }

      &:last-child {
        border-bottom: none;
      }

      .conv-info {
        flex: 1;
        min-width: 0;
      }

      .conv-header {
        margin-bottom: 4px;
      }

      .conv-name {
        font-size: 15px;
        font-weight: 500;
        color: var(--color-text-primary);
        max-width: 200px;
      }

      .conv-time {
        font-size: 12px;
        flex-shrink: 0;
      }

      .conv-last-msg {
        font-size: 13px;
        max-width: 500px;
      }

      .conv-project {
        font-size: 12px;
        display: block;
        margin-top: 4px;
      }

      .conv-arrow {
        flex-shrink: 0;
        font-size: 16px;
      }
    }
  }
}
</style>

<template>
  <div class="notifications-page page-container">
    <div class="page-header flex-between">
      <div>
        <h2>消息通知</h2>
        <p>查看你的所有通知消息</p>
      </div>
      <el-button type="primary" plain @click="handleMarkAllRead" :disabled="notifications.every(n => n.is_read)">
        <el-icon><Check /></el-icon>全部标记已读
      </el-button>
    </div>

    <div v-loading="loading" class="notification-list">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        :class="['notification-item', { unread: !notification.is_read }]"
        @click="handleClickNotification(notification)"
      >
        <div class="notification-icon">
          <el-icon :size="20" :color="iconColor(notification.type)">
            <component :is="iconComponent(notification.type)" />
          </el-icon>
        </div>
        <div class="notification-content">
          <div class="notification-header flex-between">
            <span class="notification-title">{{ notification.title }}</span>
            <span class="notification-time text-secondary">{{ formatTime(notification.created_at) }}</span>
          </div>
          <p class="notification-text">{{ notification.content }}</p>
        </div>
        <div class="notification-dot" v-if="!notification.is_read"></div>
      </div>

      <el-empty v-if="!loading && notifications.length === 0" description="暂无通知" />
    </div>

    <div class="pagination-wrap mt-20" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchNotifications"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { notificationApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()

const loading = ref(false)
const notifications = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

function iconComponent(type: string) {
  const map: Record<string, string> = {
    bid: 'EditPen',
    contract: 'Document',
    payment: 'Money',
    message: 'ChatDotRound',
    system: 'Bell',
    review: 'Star',
    project: 'Folder',
    dispute: 'Warning',
  }
  return map[type] || 'Bell'
}

function iconColor(type: string) {
  const map: Record<string, string> = {
    bid: '#409EFF',
    contract: '#67C23A',
    payment: '#E6A23C',
    message: '#909399',
    system: '#F56C6C',
    review: '#f7ba2a',
    project: '#409EFF',
    dispute: '#F56C6C',
  }
  return map[type] || '#909399'
}

function formatTime(time: string) {
  return dayjs(time).fromNow()
}

async function handleClickNotification(notification: any) {
  if (!notification.is_read) {
    try {
      await notificationApi.markAsRead(notification.id)
      notification.is_read = true
    } catch {}
  }

  // Navigate based on notification type
  if (notification.project_id) {
    router.push(`/projects/${notification.project_id}`)
  } else if (notification.contract_id) {
    router.push(`/contracts/${notification.contract_id}`)
  } else if (notification.conversation_id) {
    router.push(`/chat/${notification.conversation_id}`)
  } else if (notification.link) {
    router.push(notification.link)
  }
}

async function handleMarkAllRead() {
  try {
    await notificationApi.markAllAsRead()
    notifications.value.forEach(n => { n.is_read = true })
    ElMessage.success('已全部标记为已读')
  } catch {
    ElMessage.error('操作失败')
  }
}

async function fetchNotifications() {
  loading.value = true
  try {
    const res: any = await notificationApi.list({ page: currentPage.value, page_size: pageSize.value })
    notifications.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped lang="scss">
.notifications-page {
  .notification-list {
    background: var(--bg-color-white);
    border-radius: 12px;
    border: 1px solid var(--border-color-light);
    overflow: hidden;
  }

  .notification-item {
    display: flex;
    align-items: flex-start;
    gap: 14px;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-color-light);
    cursor: pointer;
    transition: background 0.2s;
    position: relative;

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background: var(--bg-color);
    }

    &.unread {
      background: var(--color-primary-light-9);

      &:hover {
        background: #e0edff;
      }
    }

    .notification-icon {
      width: 40px;
      height: 40px;
      border-radius: 10px;
      background: var(--bg-color);
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .notification-content {
      flex: 1;
      min-width: 0;
    }

    .notification-header {
      margin-bottom: 4px;
    }

    .notification-title {
      font-size: 14px;
      font-weight: 600;
      color: var(--color-text-primary);
    }

    .notification-time {
      font-size: 12px;
      flex-shrink: 0;
    }

    .notification-text {
      font-size: 13px;
      color: var(--color-text-secondary);
      line-height: 1.5;
    }

    .notification-dot {
      position: absolute;
      top: 20px;
      right: 16px;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--color-primary);
    }
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}
</style>

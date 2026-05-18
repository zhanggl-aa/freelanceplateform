<template>
  <div class="admin-users">
    <h2 class="page-title">用户管理</h2>

    <div class="filter-bar flex-between mb-16">
      <div class="flex gap-12">
        <el-input v-model="searchKeyword" placeholder="搜索用户名/邮箱/手机" clearable style="width: 260px;" @keyup.enter="handleSearch">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 140px;" @change="handleSearch">
          <el-option label="正常" value="active" />
          <el-option label="已停用" value="suspended" />
          <el-option label="待验证" value="pending" />
        </el-select>
      </div>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
    </div>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column prop="nickname" label="昵称" width="140">
        <template #default="{ row }">
          <div class="flex gap-8" style="align-items:center">
            <el-avatar :size="28" :src="row.avatar_url" icon="UserFilled" />
            <span>{{ row.nickname }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="phone" label="手机号" width="140">
        <template #default="{ row }">{{ row.phone || '-' }}</template>
      </el-table-column>
      <el-table-column prop="user_type" label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="userTypeTag(row.user_type)" size="small">{{ userTypeLabel(row.user_type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="170">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'active'"
            type="danger"
            text
            size="small"
            @click="handleSuspend(row)"
          >停用</el-button>
          <el-button
            v-else
            type="success"
            text
            size="small"
            @click="handleActivate(row)"
          >启用</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap mt-16">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchUsers"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/api/modules'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const users = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')
const statusFilter = ref('')

function userTypeTag(type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { developer: 'success', client: 'warning', both: 'primary', admin: 'danger' }
  return map[type] || 'info'
}

function userTypeLabel(type: string) {
  const map: Record<string, string> = { developer: '开发者', client: '客户', both: '双重', admin: '管理员' }
  return map[type] || type
}

function statusTag(status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { active: 'success', suspended: 'danger', pending: 'warning' }
  return map[status] || 'info'
}

function statusLabel(status: string) {
  const map: Record<string, string> = { active: '正常', suspended: '已停用', pending: '待验证' }
  return map[status] || status
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function handleSearch() {
  currentPage.value = 1
  fetchUsers()
}

async function handleSuspend(user: any) {
  try {
    await ElMessageBox.confirm(`确定停用用户「${user.nickname}」？`, '确认停用', { type: 'warning' })
    await adminApi.updateUserStatus(user.id, 'suspended')
    ElMessage.success('已停用')
    fetchUsers()
  } catch {}
}

async function handleActivate(user: any) {
  try {
    await ElMessageBox.confirm(`确定启用用户「${user.nickname}」？`, '确认启用')
    await adminApi.updateUserStatus(user.id, 'active')
    ElMessage.success('已启用')
    fetchUsers()
  } catch {}
}

async function fetchUsers() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (statusFilter.value) params.status = statusFilter.value

    const res: any = await adminApi.listUsers(params)
    users.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="scss">
.admin-users {
  .page-title {
    font-size: 22px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 20px;
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}
</style>

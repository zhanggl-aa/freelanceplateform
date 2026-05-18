<template>
  <div class="my-bids-page page-container">
    <div class="page-header">
      <h2>我的投标</h2>
      <p>查看和管理你提交的所有投标</p>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="已提交" name="submitted" />
      <el-tab-pane label="已入围" name="shortlisted" />
      <el-tab-pane label="已接受" name="accepted" />
      <el-tab-pane label="已拒绝" name="rejected" />
    </el-tabs>

    <el-table
      :data="bids"
      v-loading="loading"
      stripe
      class="mt-16"
      @row-click="handleRowClick"
      style="cursor: pointer"
    >
      <el-table-column prop="project_title" label="项目名称" min-width="200">
        <template #default="{ row }">
          <span class="project-name-link">{{ row.project?.title || row.project_title || '未命名项目' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="proposed_budget" label="报价金额" width="130" align="right">
        <template #default="{ row }">
          <span class="budget-text">¥{{ row.proposed_budget }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="estimated_days" label="预计天数" width="110" align="center">
        <template #default="{ row }">
          {{ row.estimated_days }}天
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="110" align="center">
        <template #default="{ row }">
          <span :class="['status-badge', `status-${row.status}`]">{{ statusLabel(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="提交时间" width="160" align="center">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && bids.length === 0" description="暂无投标记录" />

    <div class="pagination-wrap mt-20" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchBids"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { bidApi } from '@/api/modules'
import dayjs from 'dayjs'

const router = useRouter()

const loading = ref(false)
const bids = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const activeTab = ref('all')

function statusLabel(status: string) {
  const map: Record<string, string> = {
    submitted: '已提交',
    shortlisted: '已入围',
    accepted: '已接受',
    rejected: '已拒绝',
    withdrawn: '已撤回',
    counter_offered: '还价中',
  }
  return map[status] || status
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function handleTabChange() {
  currentPage.value = 1
  fetchBids()
}

function handleRowClick(row: any) {
  const projectId = row.project_id || row.project?.id
  if (projectId) {
    router.push(`/projects/${projectId}`)
  }
}

async function fetchBids() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (activeTab.value !== 'all') {
      params.status = activeTab.value
    }

    const res: any = await bidApi.myBids(params)
    bids.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchBids()
})
</script>

<style scoped lang="scss">
.my-bids-page {
  .project-name-link {
    color: var(--color-primary);
    font-weight: 500;
    &:hover {
      text-decoration: underline;
    }
  }

  .budget-text {
    font-weight: 600;
    color: var(--color-danger);
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}
</style>

<template>
  <div class="my-contracts-page page-container">
    <div class="page-header">
      <h2>我的合同</h2>
      <p>查看和管理所有合同</p>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="进行中" name="active" />
      <el-tab-pane label="已完成" name="completed" />
      <el-tab-pane label="已取消" name="cancelled" />
      <el-tab-pane label="有争议" name="disputed" />
    </el-tabs>

    <el-table
      :data="contracts"
      v-loading="loading"
      stripe
      class="mt-16"
      @row-click="handleRowClick"
      style="cursor: pointer"
    >
      <el-table-column prop="project_title" label="项目名称" min-width="200">
        <template #default="{ row }">
          <span class="project-link">{{ row.project?.title || row.project_title || '未命名项目' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="对方" width="140">
        <template #default="{ row }">
          <div class="flex gap-8" style="align-items:center">
            <el-avatar :size="24" :src="row.counterparty?.avatar_url" icon="UserFilled" />
            <span>{{ row.counterparty?.nickname || '用户' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="130" align="right">
        <template #default="{ row }">
          <span class="amount-text">¥{{ row.total_amount || row.amount }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="110" align="center">
        <template #default="{ row }">
          <span :class="['status-badge', `status-${row.status}`]">{{ statusLabel(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="start_date" label="开始日期" width="120" align="center">
        <template #default="{ row }">
          {{ formatDate(row.start_date || row.created_at) }}
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && contracts.length === 0" description="暂无合同" />

    <div class="pagination-wrap mt-20" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchContracts"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { contractApi } from '@/api/modules'
import dayjs from 'dayjs'

const router = useRouter()

const loading = ref(false)
const contracts = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const activeTab = ref('active')

function statusLabel(status: string) {
  const map: Record<string, string> = {
    pending: '待启动',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    disputed: '有争议',
    paused: '已暂停',
  }
  return map[status] || status
}

function formatDate(time: string) {
  return time ? dayjs(time).format('YYYY-MM-DD') : '-'
}

function handleTabChange() {
  currentPage.value = 1
  fetchContracts()
}

function handleRowClick(row: any) {
  router.push(`/contracts/${row.id}`)
}

async function fetchContracts() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (activeTab.value !== 'all') {
      params.status = activeTab.value
    }

    const res: any = await contractApi.myContracts(params)
    contracts.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchContracts()
})
</script>

<style scoped lang="scss">
.my-contracts-page {
  .project-link {
    color: var(--color-primary);
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }

  .amount-text {
    font-weight: 600;
    color: var(--color-danger);
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}
</style>

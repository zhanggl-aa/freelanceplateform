<template>
  <div class="admin-finance" v-loading="loading">
    <h2 class="page-title">财务管理</h2>

    <!-- Revenue Stats -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in revenueStats" :key="stat.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" :style="{ background: stat.bgColor }">
            <el-icon :size="24" :color="stat.color"><component :is="stat.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <span class="stat-label">{{ stat.label }}</span>
            <span class="stat-value">¥{{ stat.value.toLocaleString() }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent Payments Table -->
    <el-card shadow="never" class="payment-card mt-24">
      <template #header>
        <div class="flex-between">
          <h3>近期交易</h3>
          <el-select v-model="paymentTypeFilter" placeholder="交易类型" clearable style="width: 140px;" @change="handleFilterChange">
            <el-option label="托管" value="deposit" />
            <el-option label="释放" value="release" />
            <el-option label="退款" value="refund" />
            <el-option label="服务费" value="fee" />
            <el-option label="提现" value="withdraw" />
          </el-select>
        </div>
      </template>
      <el-table :data="payments" stripe>
        <el-table-column prop="id" label="交易编号" width="120">
          <template #default="{ row }">
            <span class="text-secondary">{{ row.id?.slice(0, 8) || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="用户" width="140">
          <template #default="{ row }">
            {{ row.user?.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="paymentTypeTag(row.type)" size="small">{{ paymentTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="140" align="right">
          <template #default="{ row }">
            <span :class="row.amount >= 0 ? 'amount-positive' : 'amount-negative'">
              {{ row.amount >= 0 ? '+' : '' }}¥{{ Math.abs(row.amount).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'completed' ? 'success' : row.status === 'pending' ? 'warning' : 'danger'" size="small">
              {{ paymentStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="170">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap mt-16">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next, total"
          @current-change="fetchPayments"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { adminApi } from '@/api/modules'
import dayjs from 'dayjs'

const loading = ref(false)
const payments = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const paymentTypeFilter = ref('')

const revenueStats = reactive([
  { label: '总收入', value: 0, icon: 'Money', color: '#409EFF', bgColor: '#ecf5ff' },
  { label: '平台服务费', value: 0, icon: 'Coin', color: '#67C23A', bgColor: '#f0f9eb' },
  { label: '待结算', value: 0, icon: 'Timer', color: '#E6A23C', bgColor: '#fdf6ec' },
  { label: '已退款', value: 0, icon: 'RefreshLeft', color: '#F56C6C', bgColor: '#fef0f0' },
])

function paymentTypeTag(type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { deposit: 'primary', release: 'success', refund: 'warning', fee: 'info', withdraw: 'danger' }
  return map[type] || 'info'
}

function paymentTypeLabel(type: string) {
  const map: Record<string, string> = { deposit: '托管', release: '释放', refund: '退款', fee: '服务费', withdraw: '提现' }
  return map[type] || type
}

function paymentStatusLabel(status: string) {
  const map: Record<string, string> = { completed: '已完成', pending: '处理中', failed: '失败', cancelled: '已取消' }
  return map[status] || status
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function handleFilterChange() {
  currentPage.value = 1
  fetchPayments()
}

async function fetchSummary() {
  try {
    const res: any = await adminApi.financialSummary()
    const data = res.data
    if (data) {
      revenueStats[0].value = data.total_revenue || 0
      revenueStats[1].value = data.platform_fee || 0
      revenueStats[2].value = data.pending_settlement || 0
      revenueStats[3].value = data.total_refund || 0
    }
  } catch {}
}

async function fetchPayments() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (paymentTypeFilter.value) params.type = paymentTypeFilter.value

    const res: any = await adminApi.listPayments(params)
    payments.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchSummary()
  fetchPayments()
})
</script>

<style scoped lang="scss">
.admin-finance {
  .page-title {
    font-size: 22px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 24px;
  }

  .stats-row {
    .stat-card {
      border-radius: 12px;
      border: 1px solid var(--border-color-light);
      margin-bottom: 12px;

      :deep(.el-card__body) {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px;
      }

      .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
      }

      .stat-info {
        display: flex;
        flex-direction: column;
      }

      .stat-label {
        font-size: 12px;
        color: var(--color-text-secondary);
      }

      .stat-value {
        font-size: 22px;
        font-weight: 700;
        color: var(--color-text-primary);
        margin-top: 4px;
      }
    }
  }

  .payment-card {
    border-radius: 12px;

    h3 {
      font-size: 15px;
      font-weight: 600;
    }
  }

  .amount-positive {
    color: var(--color-success);
    font-weight: 600;
  }

  .amount-negative {
    color: var(--color-danger);
    font-weight: 600;
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}
</style>

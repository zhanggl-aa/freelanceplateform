<template>
  <div class="admin-disputes">
    <h2 class="page-title">争议处理</h2>

    <div class="filter-bar flex-between mb-16">
      <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 160px;" @change="handleFilterChange">
        <el-option label="待处理" value="pending" />
        <el-option label="处理中" value="investigating" />
        <el-option label="已解决" value="resolved" />
        <el-option label="已关闭" value="closed" />
      </el-select>
    </div>

    <el-table :data="disputes" v-loading="loading" stripe>
      <el-table-column prop="id" label="争议编号" width="120">
        <template #default="{ row }">
          <span class="text-secondary">{{ row.id?.slice(0, 8) || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="合同编号" width="120">
        <template #default="{ row }">
          <span class="text-secondary">{{ row.contract_id?.slice(0, 8) || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="举报方" width="140">
        <template #default="{ row }">
          <div class="flex gap-8" style="align-items:center">
            <el-avatar :size="24" :src="row.reporter?.avatar_url" icon="UserFilled" />
            <span>{{ row.reporter?.nickname || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="被举报方" width="140">
        <template #default="{ row }">
          <div class="flex gap-8" style="align-items:center">
            <el-avatar :size="24" :src="row.reported?.avatar_url" icon="UserFilled" />
            <span>{{ row.reported?.nickname || '-' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="争议原因" min-width="200" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="disputeStatusTag(row.status)" size="small">{{ disputeStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" text size="small" @click="handleResolve(row)" :disabled="row.status === 'resolved' || row.status === 'closed'">
            处理
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap mt-16">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchDisputes"
      />
    </div>

    <!-- Resolve Dialog -->
    <el-dialog v-model="showResolveDialog" title="处理争议" width="550px" :close-on-click-modal="false">
      <el-descriptions :column="1" border size="small" v-if="resolvingDispute">
        <el-descriptions-item label="举报方">{{ resolvingDispute.reporter?.nickname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="被举报方">{{ resolvingDispute.reported?.nickname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="争议原因">{{ resolvingDispute.reason }}</el-descriptions-item>
      </el-descriptions>
      <el-form class="mt-16" label-position="top">
        <el-form-item label="处理结果">
          <el-radio-group v-model="resolveResult">
            <el-radio value="favor_reporter">支持举报方</el-radio>
            <el-radio value="favor_reported">支持被举报方</el-radio>
            <el-radio value="compromise">双方妥协</el-radio>
            <el-radio value="dismiss">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="处理意见" required>
          <el-input v-model="resolveComment" type="textarea" :rows="4" placeholder="输入处理意见和理由" />
        </el-form-item>
        <el-form-item label="退款金额（选填）" v-if="resolveResult !== 'dismiss'">
          <el-input-number v-model="refundAmount" :min="0" :max="9999999" class="w-full" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showResolveDialog = false">取消</el-button>
        <el-button type="primary" :loading="resolving" @click="handleSubmitResolve">提交处理</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const disputes = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const statusFilter = ref('')

const showResolveDialog = ref(false)
const resolvingDispute = ref<any>(null)
const resolveResult = ref('favor_reporter')
const resolveComment = ref('')
const refundAmount = ref(0)
const resolving = ref(false)

function disputeStatusLabel(status: string) {
  const map: Record<string, string> = { pending: '待处理', investigating: '处理中', resolved: '已解决', closed: '已关闭' }
  return map[status] || status
}

function disputeStatusTag(status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { pending: 'danger', investigating: 'warning', resolved: 'success', closed: 'info' }
  return map[status] || 'info'
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function handleFilterChange() {
  currentPage.value = 1
  fetchDisputes()
}

function handleResolve(dispute: any) {
  resolvingDispute.value = dispute
  resolveResult.value = 'favor_reporter'
  resolveComment.value = ''
  refundAmount.value = 0
  showResolveDialog.value = true
}

async function handleSubmitResolve() {
  if (!resolveComment.value.trim()) {
    ElMessage.warning('请输入处理意见')
    return
  }
  if (!resolvingDispute.value) return

  resolving.value = true
  try {
    await adminApi.resolveDispute(resolvingDispute.value.id, {
      result: resolveResult.value,
      comment: resolveComment.value,
      refund_amount: refundAmount.value || undefined,
    })
    ElMessage.success('争议处理已提交')
    showResolveDialog.value = false
    fetchDisputes()
  } catch {
    ElMessage.error('处理失败')
  } finally {
    resolving.value = false
  }
}

async function fetchDisputes() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (statusFilter.value) params.status = statusFilter.value

    const res: any = await adminApi.listDisputes(params)
    disputes.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDisputes()
})
</script>

<style scoped lang="scss">
.admin-disputes {
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

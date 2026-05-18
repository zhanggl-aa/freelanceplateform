<template>
  <div class="contract-detail-page page-container" v-loading="loading">
    <template v-if="contract">
      <!-- Header -->
      <div class="page-header">
        <div class="flex-between">
          <div>
            <h2>合同详情</h2>
            <p class="mt-4">合同编号: {{ contract.id }}</p>
          </div>
          <span :class="['status-badge', `status-${contract.status}`]">{{ statusLabel(contract.status) }}</span>
        </div>
      </div>

      <!-- Project Info -->
      <el-card shadow="never" class="detail-card">
        <template #header><h3>项目信息</h3></template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目名称" :span="2">
            <span class="project-link" @click="goProject">{{ contract.project?.title || '未命名项目' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="预算类型">{{ contract.project?.budget_type === 'hourly' ? '时薪制' : '固定价' }}</el-descriptions-item>
          <el-descriptions-item label="项目预算">
            <span v-if="contract.project?.budget_type === 'hourly'">¥{{ contract.project?.budget_min }}-{{ contract.project?.budget_max }}/时</span>
            <span v-else>¥{{ contract.project?.budget_min }}-{{ contract.project?.budget_max }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Parties -->
      <el-card shadow="never" class="detail-card">
        <template #header><h3>合同双方</h3></template>
        <el-row :gutter="24">
          <el-col :span="12">
            <div class="party-card">
              <el-avatar :size="48" :src="contract.client?.avatar_url" icon="UserFilled" />
              <div class="party-info">
                <span class="party-name">{{ contract.client?.nickname || '甲方' }}</span>
                <el-tag size="small" type="warning">甲方</el-tag>
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="party-card">
              <el-avatar :size="48" :src="contract.developer?.avatar_url" icon="UserFilled" />
              <div class="party-info">
                <span class="party-name">{{ contract.developer?.nickname || '开发者' }}</span>
                <el-tag size="small" type="success">开发者</el-tag>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- Amount Breakdown -->
      <el-card shadow="never" class="detail-card">
        <template #header><h3>金额明细</h3></template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="合同总额">
            <span class="amount-primary">¥{{ contract.total_amount || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="平台服务费">
            <span>¥{{ contract.platform_fee || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="开发者收入">
            <span class="amount-success">¥{{ contract.developer_payout || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="已支付">
            <span>¥{{ contract.paid_amount || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="已释放">
            <span class="amount-success">¥{{ contract.released_amount || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="待释放">
            <span class="amount-warning">¥{{ (contract.paid_amount || 0) - (contract.released_amount || 0) }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Milestones -->
      <el-card shadow="never" class="detail-card">
        <template #header><h3>里程碑</h3></template>
        <div v-if="milestones.length > 0">
          <div v-for="(ms, idx) in milestones" :key="ms.id" class="milestone-item">
            <div class="milestone-header flex-between">
              <div class="flex gap-12" style="align-items:center">
                <span class="milestone-index">{{ idx + 1 }}</span>
                <div>
                  <h4 class="milestone-title">{{ ms.title }}</h4>
                  <span class="text-secondary">截止 {{ formatDate(ms.deadline) }}</span>
                </div>
              </div>
              <div class="flex gap-12" style="align-items:center">
                <span class="milestone-amount">¥{{ ms.amount }}</span>
                <span :class="['status-badge', `status-${ms.status}`]">{{ milestoneLabel(ms.status) }}</span>
              </div>
            </div>
            <p class="milestone-desc mt-8" v-if="ms.description">{{ ms.description }}</p>
            <div class="milestone-actions mt-8" v-if="contract.status === 'active'">
              <!-- Developer: submit -->
              <el-button
                v-if="isDeveloper && (ms.status === 'pending' || ms.status === 'rejected')"
                size="small"
                type="primary"
                @click="handleSubmitMilestone(ms)"
              >
                提交交付
              </el-button>
              <!-- Client: approve/reject -->
              <template v-if="isClient && ms.status === 'submitted'">
                <el-button size="small" type="success" @click="handleApproveMilestone(ms)">批准</el-button>
                <el-button size="small" type="danger" @click="handleRejectMilestone(ms)">拒绝</el-button>
              </template>
            </div>
          </div>
        </div>
        <el-empty v-else description="暂无里程碑" :image-size="60" />
      </el-card>

      <!-- Payment History -->
      <el-card shadow="never" class="detail-card">
        <template #header><h3>支付记录</h3></template>
        <el-table :data="payments" stripe size="small" v-if="payments.length > 0">
          <el-table-column prop="type" label="类型" width="120">
            <template #default="{ row }">{{ paymentTypeLabel(row.type) }}</template>
          </el-table-column>
          <el-table-column prop="amount" label="金额" width="120" align="right">
            <template #default="{ row }">
              <span :class="row.amount >= 0 ? 'amount-success' : 'amount-danger'">
                {{ row.amount >= 0 ? '+' : '' }}¥{{ row.amount }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag size="small" :type="row.status === 'completed' ? 'success' : 'warning'">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="200" />
          <el-table-column prop="created_at" label="时间" width="160">
            <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="暂无支付记录" :image-size="60" />
      </el-card>

      <!-- Actions -->
      <el-card shadow="never" class="detail-card" v-if="contract.status === 'pending' || contract.status === 'active'">
        <template #header><h3>合同操作</h3></template>
        <div class="flex gap-12">
          <el-button v-if="contract.status === 'pending' && isClient" type="primary" @click="handleStartContract">启动合同</el-button>
          <el-button v-if="contract.status === 'active'" type="danger" plain @click="handleCancelContract">取消合同</el-button>
          <el-button v-if="contract.status === 'active'" type="warning" plain @click="handleOpenDispute">发起争议</el-button>
        </div>
      </el-card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { contractApi, milestoneApi, paymentApi } from '@/api/modules'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const contractId = computed(() => route.params.id as string)
const loading = ref(false)
const contract = ref<any>(null)
const milestones = ref<any[]>([])
const payments = ref<any[]>([])

const isDeveloper = computed(() => userStore.isDeveloper)
const isClient = computed(() => userStore.isClient)

function statusLabel(status: string) {
  const map: Record<string, string> = { pending: '待启动', active: '进行中', completed: '已完成', cancelled: '已取消', disputed: '争议中', paused: '已暂停' }
  return map[status] || status
}

function milestoneLabel(status: string) {
  const map: Record<string, string> = { pending: '待开始', in_progress: '进行中', submitted: '待审批', approved: '已通过', rejected: '已拒绝' }
  return map[status] || status
}

function paymentTypeLabel(type: string) {
  const map: Record<string, string> = { deposit: '托管', release: '释放', refund: '退款', fee: '服务费', withdraw: '提现' }
  return map[type] || type
}

function formatDate(time: string) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

function goProject() {
  if (contract.value?.project_id) {
    router.push(`/projects/${contract.value.project_id}`)
  }
}

async function fetchContract() {
  loading.value = true
  try {
    const res: any = await contractApi.getById(contractId.value)
    contract.value = res.data
  } catch {} finally {
    loading.value = false
  }
}

async function fetchMilestones() {
  try {
    const projectId = contract.value?.project_id
    if (!projectId) return
    const res: any = await milestoneApi.listByProject(projectId)
    milestones.value = res.data || []
  } catch {}
}

async function fetchPayments() {
  try {
    const res: any = await paymentApi.myPayments({ contract_id: contractId.value, page: 1, page_size: 50 })
    payments.value = res.data?.items || res.data || []
  } catch {}
}

async function handleSubmitMilestone(ms: any) {
  try {
    await ElMessageBox.prompt('请输入交付说明', '提交交付', {
      confirmButtonText: '提交',
      cancelButtonText: '取消',
      inputPlaceholder: '描述你的交付内容',
    }).then(async ({ value }) => {
      await milestoneApi.submit(ms.id, { description: value })
      ElMessage.success('已提交交付')
      fetchMilestones()
    })
  } catch {}
}

async function handleApproveMilestone(ms: any) {
  try {
    await ElMessageBox.confirm('确定批准此里程碑？批准后将释放对应款项。', '确认批准')
    await milestoneApi.approve(ms.id)
    ElMessage.success('已批准里程碑')
    fetchMilestones()
    fetchContract()
    fetchPayments()
  } catch {}
}

async function handleRejectMilestone(ms: any) {
  try {
    await ElMessageBox.prompt('请输入拒绝原因', '拒绝交付', {
      confirmButtonText: '拒绝',
      cancelButtonText: '取消',
      inputPlaceholder: '说明拒绝原因',
    }).then(async ({ value }) => {
      await milestoneApi.reject(ms.id, { reason: value })
      ElMessage.success('已拒绝')
      fetchMilestones()
    })
  } catch {}
}

async function handleStartContract() {
  try {
    await ElMessageBox.confirm('确定启动此合同？', '确认启动')
    await contractApi.start(contractId.value)
    ElMessage.success('合同已启动')
    fetchContract()
  } catch {}
}

async function handleCancelContract() {
  try {
    await ElMessageBox.confirm('确定取消此合同？此操作不可撤销。', '确认取消', { type: 'warning' })
    await contractApi.cancel(contractId.value)
    ElMessage.success('合同已取消')
    fetchContract()
  } catch {}
}

async function handleOpenDispute() {
  try {
    await ElMessageBox.prompt('请描述争议原因', '发起争议', {
      confirmButtonText: '提交',
      cancelButtonText: '取消',
      inputType: 'textarea',
      inputPlaceholder: '详细说明争议原因',
    }).then(async ({ value }) => {
      await contractApi.dispute(contractId.value, { reason: value })
      ElMessage.success('争议已提交')
      fetchContract()
    })
  } catch {}
}

onMounted(async () => {
  await fetchContract()
  await fetchMilestones()
  await fetchPayments()
})
</script>

<style scoped lang="scss">
.contract-detail-page {
  .detail-card {
    margin-bottom: 20px;
    border-radius: 12px;

    h3 {
      font-size: 16px;
      font-weight: 600;
    }
  }

  .project-link {
    color: var(--color-primary);
    cursor: pointer;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }

  .party-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    border: 1px solid var(--border-color-light);
    border-radius: 8px;

    .party-info {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .party-name {
      font-size: 15px;
      font-weight: 500;
    }
  }

  .amount-primary {
    font-size: 18px;
    font-weight: 700;
    color: var(--color-text-primary);
  }

  .amount-success {
    color: var(--color-success);
    font-weight: 600;
  }

  .amount-warning {
    color: var(--color-warning);
    font-weight: 600;
  }

  .amount-danger {
    color: var(--color-danger);
  }

  .milestone-item {
    padding: 16px 0;
    border-bottom: 1px solid var(--border-color-light);

    &:last-child {
      border-bottom: none;
    }

    .milestone-index {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 28px;
      height: 28px;
      border-radius: 50%;
      background: var(--color-primary-light-9);
      color: var(--color-primary);
      font-weight: 600;
      font-size: 14px;
      flex-shrink: 0;
    }

    .milestone-title {
      font-size: 15px;
      font-weight: 500;
      color: var(--color-text-primary);
    }

    .milestone-amount {
      font-size: 16px;
      font-weight: 600;
      color: var(--color-danger);
    }

    .milestone-desc {
      font-size: 13px;
      color: var(--color-text-secondary);
      padding-left: 40px;
    }

    .milestone-actions {
      padding-left: 40px;
    }
  }
}
</style>

<template>
  <div class="wallet-page page-container">
    <div class="page-header">
      <h2>我的钱包</h2>
      <p>管理你的资金和交易记录</p>
    </div>

    <!-- Balance Cards -->
    <el-row :gutter="20" class="balance-row">
      <el-col :xs="24" :sm="12">
        <el-card shadow="never" class="balance-card main-balance">
          <div class="balance-content">
            <div class="balance-icon">
              <el-icon :size="40" color="#409EFF"><Wallet /></el-icon>
            </div>
            <div class="balance-info">
              <span class="balance-label">账户余额</span>
              <span class="balance-value">¥{{ balance }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12">
        <el-card shadow="never" class="balance-card frozen-balance">
          <div class="balance-content">
            <div class="balance-icon">
              <el-icon :size="40" color="#E6A23C"><Lock /></el-icon>
            </div>
            <div class="balance-info">
              <span class="balance-label">冻结金额</span>
              <span class="balance-value frozen">¥{{ frozenAmount }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Actions -->
    <div class="action-bar mt-20 flex gap-12">
      <el-button type="primary" @click="showDepositDialog = true">
        <el-icon><Plus /></el-icon>充值
      </el-button>
      <el-button @click="showWithdrawDialog = true">
        <el-icon><Minus /></el-icon>提现
      </el-button>
    </div>

    <!-- Transaction History -->
    <el-card shadow="never" class="transaction-card mt-24">
      <template #header>
        <h3>交易记录</h3>
      </template>
      <el-table :data="transactions" v-loading="loading" stripe>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[row.type] || 'info'" size="small">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="140" align="right">
          <template #default="{ row }">
            <span :class="row.amount >= 0 ? 'amount-positive' : 'amount-negative'">
              {{ row.amount >= 0 ? '+' : '' }}¥{{ Math.abs(row.amount).toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="余额" width="140" align="right">
          <template #default="{ row }">
            ¥{{ (row.balance_after || 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" />
      </el-table>

      <el-empty v-if="!loading && transactions.length === 0" description="暂无交易记录" />

      <div class="pagination-wrap mt-16" v-if="txTotal > txPageSize">
        <el-pagination
          v-model:current-page="txCurrentPage"
          :page-size="txPageSize"
          :total="txTotal"
          layout="prev, pager, next"
          @current-change="fetchTransactions"
        />
      </div>
    </el-card>

    <!-- Deposit Dialog -->
    <el-dialog v-model="showDepositDialog" title="充值" width="400px" :close-on-click-modal="false">
      <el-form label-position="top">
        <el-form-item label="充值金额">
          <el-input-number v-model="depositAmount" :min="1" :max="999999" :step="100" class="w-full" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDepositDialog = false">取消</el-button>
        <el-button type="primary" :loading="depositing" @click="handleDeposit">确认充值</el-button>
      </template>
    </el-dialog>

    <!-- Withdraw Dialog -->
    <el-dialog v-model="showWithdrawDialog" title="提现" width="400px" :close-on-click-modal="false">
      <el-form label-position="top">
        <el-form-item label="提现金额">
          <el-input-number v-model="withdrawAmount" :min="1" :max="balance" :step="100" class="w-full" />
        </el-form-item>
        <el-form-item label="提现方式">
          <el-select v-model="withdrawMethod" class="w-full">
            <el-option label="银行卡" value="bank" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showWithdrawDialog = false">取消</el-button>
        <el-button type="primary" :loading="withdrawing" @click="handleWithdraw">确认提现</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { paymentApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const balance = ref(0)
const frozenAmount = ref(0)
const transactions = ref<any[]>([])
const txTotal = ref(0)
const txCurrentPage = ref(1)
const txPageSize = ref(20)

const showDepositDialog = ref(false)
const depositAmount = ref(100)
const depositing = ref(false)

const showWithdrawDialog = ref(false)
const withdrawAmount = ref(100)
const withdrawMethod = ref('bank')
const withdrawing = ref(false)

const typeTagMap: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
  deposit: 'primary',
  release: 'success',
  refund: 'warning',
  fee: 'info',
  withdraw: 'danger',
  income: 'success',
}

function typeLabel(type: string) {
  const map: Record<string, string> = { deposit: '托管', release: '释放', refund: '退款', fee: '服务费', withdraw: '提现', income: '收入' }
  return map[type] || type
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

async function fetchBalance() {
  try {
    const res: any = await paymentApi.walletBalance()
    balance.value = res.data?.balance || res.data?.available || 0
    frozenAmount.value = res.data?.frozen || res.data?.frozen_amount || 0
  } catch {}
}

async function fetchTransactions() {
  loading.value = true
  try {
    const res: any = await paymentApi.walletTransactions({ page: txCurrentPage.value, page_size: txPageSize.value })
    transactions.value = res.data?.items || res.data || []
    txTotal.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

async function handleDeposit() {
  if (depositAmount.value <= 0) {
    ElMessage.warning('请输入充值金额')
    return
  }
  depositing.value = true
  try {
    await paymentApi.deposit({ amount: depositAmount.value })
    ElMessage.success('充值请求已提交')
    showDepositDialog.value = false
    fetchBalance()
    fetchTransactions()
  } catch {
    ElMessage.error('充值失败')
  } finally {
    depositing.value = false
  }
}

async function handleWithdraw() {
  if (withdrawAmount.value <= 0) {
    ElMessage.warning('请输入提现金额')
    return
  }
  if (withdrawAmount.value > balance.value) {
    ElMessage.warning('提现金额不能超过余额')
    return
  }
  withdrawing.value = true
  try {
    await paymentApi.withdraw({ amount: withdrawAmount.value, method: withdrawMethod.value })
    ElMessage.success('提现请求已提交')
    showWithdrawDialog.value = false
    fetchBalance()
    fetchTransactions()
  } catch {
    ElMessage.error('提现失败')
  } finally {
    withdrawing.value = false
  }
}

onMounted(() => {
  fetchBalance()
  fetchTransactions()
})
</script>

<style scoped lang="scss">
.wallet-page {
  .balance-row {
    .balance-card {
      border-radius: 12px;
      border: 1px solid var(--border-color-light);

      .balance-content {
        display: flex;
        align-items: center;
        gap: 16px;
      }

      .balance-icon {
        width: 72px;
        height: 72px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: var(--color-primary-light-9);
      }

      .frozen-balance .balance-icon {
        background: var(--color-warning-light-9);
      }

      .balance-info {
        display: flex;
        flex-direction: column;
      }

      .balance-label {
        font-size: 14px;
        color: var(--color-text-secondary);
      }

      .balance-value {
        font-size: 28px;
        font-weight: 700;
        color: var(--color-text-primary);
        margin-top: 4px;

        &.frozen {
          color: var(--color-warning);
        }
      }
    }
  }

  .transaction-card {
    border-radius: 12px;

    h3 {
      font-size: 16px;
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

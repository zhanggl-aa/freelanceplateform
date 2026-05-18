<template>
  <div class="admin-dashboard" v-loading="loading">
    <h2 class="page-title">仪表盘</h2>

    <!-- Stats Cards -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="8" :md="4" v-for="stat in statsCards" :key="stat.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" :style="{ background: stat.bgColor }">
            <el-icon :size="24" :color="stat.color"><component :is="stat.icon" /></el-icon>
          </div>
          <el-statistic :title="stat.label" :value="stat.value" class="stat-value" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Chart Placeholders -->
    <el-row :gutter="20" class="mt-24">
      <el-col :xs="24" :sm="12">
        <el-card shadow="never" class="chart-card">
          <template #header><h3>用户增长趋势</h3></template>
          <div class="chart-placeholder">
            <el-icon :size="48" color="#c0c4cc"><TrendCharts /></el-icon>
            <p class="text-secondary mt-8">图表区域 - 用户增长趋势</p>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12">
        <el-card shadow="never" class="chart-card">
          <template #header><h3>收入趋势</h3></template>
          <div class="chart-placeholder">
            <el-icon :size="48" color="#c0c4cc"><Money /></el-icon>
            <p class="text-secondary mt-8">图表区域 - 收入趋势</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent Activity -->
    <el-row :gutter="20" class="mt-24">
      <el-col :span="24">
        <el-card shadow="never" class="activity-card">
          <template #header><h3>近期动态</h3></template>
          <el-table :data="recentActivities" stripe size="small">
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="activityTagType(row.type)" size="small">{{ row.type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="300" />
            <el-table-column prop="time" label="时间" width="160" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { adminApi } from '@/api/modules'

const loading = ref(false)
const recentActivities = ref<any[]>([])

const statsCards = reactive([
  { label: '总用户数', value: 0, icon: 'User', color: '#409EFF', bgColor: '#ecf5ff' },
  { label: '项目总数', value: 0, icon: 'Folder', color: '#67C23A', bgColor: '#f0f9eb' },
  { label: '活跃合同', value: 0, icon: 'Document', color: '#E6A23C', bgColor: '#fdf6ec' },
  { label: '总收入', value: 0, icon: 'Money', color: '#F56C6C', bgColor: '#fef0f0' },
  { label: '待处理争议', value: 0, icon: 'Warning', color: '#9b59b6', bgColor: '#f5eef8' },
  { label: '今日新增', value: 0, icon: 'TrendCharts', color: '#34495e', bgColor: '#eef1f5' },
])

function activityTagType(type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { 用户: 'primary', 项目: 'success', 合同: 'warning', 支付: 'danger', 争议: 'info' }
  return map[type] || 'info'
}

async function fetchDashboard() {
  loading.value = true
  try {
    const res: any = await adminApi.dashboard()
    const data = res.data
    if (data) {
      statsCards[0].value = data.total_users || 0
      statsCards[1].value = data.total_projects || 0
      statsCards[2].value = data.active_contracts || 0
      statsCards[3].value = data.total_revenue || 0
      statsCards[4].value = data.pending_disputes || 0
      statsCards[5].value = data.today_new_users || 0
      recentActivities.value = data.recent_activities || []
    }
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboard()
})
</script>

<style scoped lang="scss">
.admin-dashboard {
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

      .stat-value {
        :deep(.el-statistic__head) {
          font-size: 12px;
          color: var(--color-text-secondary);
          margin-bottom: 4px;
        }

        :deep(.el-statistic__content) {
          font-size: 24px;
          font-weight: 700;
        }
      }
    }
  }

  .chart-card {
    border-radius: 12px;
    margin-bottom: 12px;

    h3 {
      font-size: 15px;
      font-weight: 600;
    }

    .chart-placeholder {
      height: 240px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      background: var(--bg-color);
      border-radius: 8px;
    }
  }

  .activity-card {
    border-radius: 12px;

    h3 {
      font-size: 15px;
      font-weight: 600;
    }
  }
}
</style>

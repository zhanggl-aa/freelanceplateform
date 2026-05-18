<template>
  <div class="admin-projects">
    <h2 class="page-title">项目管理</h2>

    <div class="filter-bar flex-between mb-16">
      <div class="flex gap-12">
        <el-input v-model="searchKeyword" placeholder="搜索项目名称" clearable style="width: 260px;" @keyup.enter="handleSearch">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable style="width: 140px;" @change="handleSearch">
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="published" />
          <el-option label="进行中" value="ongoing" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </div>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
    </div>

    <el-table :data="projects" v-loading="loading" stripe>
      <el-table-column prop="title" label="项目名称" min-width="220">
        <template #default="{ row }">
          <span class="project-name">{{ row.title }}</span>
        </template>
      </el-table-column>
      <el-table-column label="发布方" width="140">
        <template #default="{ row }">
          {{ row.client?.nickname || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="分类" width="120">
        <template #default="{ row }">
          {{ row.category_name || row.category?.name || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="预算" width="140" align="right">
        <template #default="{ row }">
          <span v-if="row.budget_type === 'hourly'">¥{{ row.budget_min }}-{{ row.budget_max }}/时</span>
          <span v-else>¥{{ row.budget_min }}-{{ row.budget_max }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <span :class="['status-badge', `status-${row.status}`]">{{ statusLabel(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" text size="small" @click="handleReview(row)">审核</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap mt-16">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchProjects"
      />
    </div>

    <!-- Review Dialog -->
    <el-dialog v-model="showReviewDialog" title="项目审核" width="500px">
      <el-descriptions :column="1" border size="small" v-if="reviewingProject">
        <el-descriptions-item label="项目名称">{{ reviewingProject.title }}</el-descriptions-item>
        <el-descriptions-item label="发布方">{{ reviewingProject.client?.nickname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ statusLabel(reviewingProject.status) }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ reviewingProject.description }}</el-descriptions-item>
      </el-descriptions>
      <el-form class="mt-16" label-position="top">
        <el-form-item label="审核操作">
          <el-radio-group v-model="reviewAction">
            <el-radio value="approve">通过</el-radio>
            <el-radio value="reject">拒绝</el-radio>
            <el-radio value="suspend">暂停</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审核意见">
          <el-input v-model="reviewComment" type="textarea" :rows="3" placeholder="输入审核意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" :loading="reviewing" @click="handleSubmitReview">提交审核</el-button>
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
const projects = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')
const statusFilter = ref('')

const showReviewDialog = ref(false)
const reviewingProject = ref<any>(null)
const reviewAction = ref('approve')
const reviewComment = ref('')
const reviewing = ref(false)

function statusLabel(status: string) {
  const map: Record<string, string> = { draft: '草稿', published: '已发布', ongoing: '进行中', completed: '已完成', cancelled: '已取消', disputed: '争议中' }
  return map[status] || status
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function handleSearch() {
  currentPage.value = 1
  fetchProjects()
}

function handleReview(project: any) {
  reviewingProject.value = project
  reviewAction.value = 'approve'
  reviewComment.value = ''
  showReviewDialog.value = true
}

async function handleSubmitReview() {
  if (!reviewingProject.value) return
  reviewing.value = true
  try {
    const statusMap: Record<string, string> = { approve: 'published', reject: 'rejected', suspend: 'suspended' }
    await adminApi.moderateProject(reviewingProject.value.id, {
      status: statusMap[reviewAction.value],
      comment: reviewComment.value,
    })
    ElMessage.success('审核已提交')
    showReviewDialog.value = false
    fetchProjects()
  } catch {
    ElMessage.error('审核失败')
  } finally {
    reviewing.value = false
  }
}

async function fetchProjects() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (statusFilter.value) params.status = statusFilter.value

    const res: any = await adminApi.listProjects(params)
    projects.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProjects()
})
</script>

<style scoped lang="scss">
.admin-projects {
  .page-title {
    font-size: 22px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 20px;
  }

  .project-name {
    color: var(--color-primary);
    font-weight: 500;
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}
</style>

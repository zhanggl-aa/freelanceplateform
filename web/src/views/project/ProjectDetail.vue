<template>
  <div class="project-detail-page page-container" v-loading="loading">
    <div class="detail-layout" v-if="project">
      <!-- Main Content -->
      <main class="detail-main">
        <!-- Header -->
        <div class="detail-header">
          <div class="detail-header-top">
            <h1 class="detail-title">{{ project.title }}</h1>
            <div class="header-actions">
              <el-button :type="isBookmarked ? 'warning' : 'default'" @click="toggleBookmark">
                <el-icon><Star /></el-icon>
                {{ isBookmarked ? '已收藏' : '收藏' }}
              </el-button>
            </div>
          </div>
          <div class="detail-meta">
            <el-tag :type="statusTagType(project.status)">{{ statusLabel(project.status) }}</el-tag>
            <el-tag type="primary" effect="plain">{{ project.category_name || project.category?.name || '未分类' }}</el-tag>
            <el-tag :type="project.budget_type === 'hourly' ? 'warning' : 'primary'">
              {{ project.budget_type === 'hourly' ? '时薪制' : '固定价' }}
            </el-tag>
            <span class="meta-item"><el-icon><View /></el-icon> {{ project.view_count || 0 }} 次浏览</span>
            <span class="meta-item"><el-icon><Clock /></el-icon> {{ formatTime(project.created_at) }}</span>
          </div>
        </div>

        <!-- Description -->
        <el-card shadow="never" class="detail-card">
          <template #header><h3>项目描述</h3></template>
          <div class="project-description" v-html="project.description"></div>
        </el-card>

        <!-- Project Info -->
        <el-card shadow="never" class="detail-card">
          <template #header><h3>项目信息</h3></template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="预算范围">
              <span v-if="project.budget_type === 'hourly'">¥{{ project.budget_min }} - {{ project.budget_max }}/时</span>
              <span v-else>¥{{ project.budget_min }} - {{ project.budget_max }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="截止日期">{{ project.deadline ? formatDate(project.deadline) : '无期限' }}</el-descriptions-item>
            <el-descriptions-item label="投标截止">{{ project.bid_deadline ? formatDate(project.bid_deadline) : '无限制' }}</el-descriptions-item>
            <el-descriptions-item label="投标数量">{{ project.bid_count || 0 }} 人</el-descriptions-item>
            <el-descriptions-item label="技术栈" :span="2">
              <el-tag v-for="tech in (project.tech_stack || [])" :key="tech" size="small" effect="plain" class="tech-tag">{{ tech }}</el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- Milestones -->
        <el-card shadow="never" class="detail-card" v-if="milestones.length > 0">
          <template #header><h3>里程碑</h3></template>
          <el-timeline>
            <el-timeline-item
              v-for="(ms, idx) in milestones"
              :key="ms.id"
              :type="milestoneColor(ms.status)"
              :timestamp="`第${idx + 1}阶段 - 截止 ${formatDate(ms.deadline)}`"
              placement="top"
            >
              <h4>{{ ms.title }}</h4>
              <p class="text-secondary">{{ ms.description }}</p>
              <p>金额: ¥{{ ms.amount }}</p>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- Attachments -->
        <el-card shadow="never" class="detail-card" v-if="project.attachments && project.attachments.length > 0">
          <template #header><h3>附件</h3></template>
          <div class="attachment-list">
            <div v-for="file in project.attachments" :key="file.id" class="attachment-item" @click="downloadFile(file)">
              <el-icon size="20"><Document /></el-icon>
              <span class="attachment-name">{{ file.name || file.filename }}</span>
              <el-button text type="primary" size="small">下载</el-button>
            </div>
          </div>
        </el-card>

        <!-- Bids Received (for client) -->
        <el-card shadow="never" class="detail-card" v-if="isProjectOwner && project.status === 'published'">
          <template #header>
            <div class="flex-between">
              <h3>收到的投标 ({{ bids.length }})</h3>
            </div>
          </template>
          <div v-loading="bidsLoading">
            <div v-for="bid in bids" :key="bid.id" class="bid-item">
              <div class="bid-header">
                <div class="bid-developer flex" style="gap:10px">
                  <el-avatar :size="36" :src="bid.developer?.avatar_url" icon="UserFilled" />
                  <div>
                    <h4>{{ bid.developer?.nickname || '开发者' }}</h4>
                    <span class="text-secondary">{{ bid.developer?.title || '' }}</span>
                  </div>
                </div>
                <span :class="['status-badge', `status-${bid.status}`]">{{ bidStatusLabel(bid.status) }}</span>
              </div>
              <p class="bid-cover mt-8">{{ bid.cover_letter }}</p>
              <div class="bid-details mt-8 flex gap-20">
                <span>报价: ¥{{ bid.proposed_budget }}</span>
                <span>工期: {{ bid.estimated_days }}天</span>
                <span>类型: {{ bid.budget_type === 'hourly' ? '时薪制' : '固定价' }}</span>
              </div>
              <div class="bid-actions mt-12" v-if="bid.status === 'submitted' || bid.status === 'shortlisted'">
                <el-button type="success" size="small" @click="handleAcceptBid(bid.id)">接受</el-button>
                <el-button type="danger" size="small" @click="handleRejectBid(bid.id)">拒绝</el-button>
                <el-button size="small" @click="handleShortlistBid(bid.id)" v-if="bid.status === 'submitted'">加入候选</el-button>
              </div>
            </div>
            <el-empty v-if="bids.length === 0" description="暂无投标" />
          </div>
        </el-card>
      </main>

      <!-- Sidebar -->
      <aside class="detail-sidebar">
        <!-- Client Info -->
        <el-card shadow="never" class="sidebar-card">
          <template #header><h3>发布方</h3></template>
          <div class="client-info">
            <el-avatar :size="56" :src="project.client?.avatar_url" icon="UserFilled" />
            <h4 class="mt-8">{{ project.client?.nickname || '甲方' }}</h4>
            <div class="client-stats mt-8 text-secondary">
              <span>已发布 {{ project.client?.posted_projects || 0 }} 个项目</span>
            </div>
            <el-button type="primary" plain class="w-full mt-12" @click="startChat(project.client?.id)" v-if="!isProjectOwner">联系甲方</el-button>
          </div>
        </el-card>

        <!-- Bid Action (for developer) -->
        <el-card shadow="never" class="sidebar-card" v-if="userStore.isDeveloper && !isProjectOwner && project.status === 'published'">
          <el-button type="primary" size="large" class="w-full" @click="showBidDialog = true">提交投标</el-button>
        </el-card>

        <!-- Quick Stats -->
        <el-card shadow="never" class="sidebar-card">
          <template #header><h3>项目统计</h3></template>
          <div class="stats-grid">
            <div class="stat-cell">
              <el-icon><EditPen /></el-icon>
              <span>{{ project.bid_count || 0 }}</span>
              <label>投标</label>
            </div>
            <div class="stat-cell">
              <el-icon><View /></el-icon>
              <span>{{ project.view_count || 0 }}</span>
              <label>浏览</label>
            </div>
          </div>
        </el-card>
      </aside>
    </div>

    <!-- Bid Dialog -->
    <el-dialog v-model="showBidDialog" title="提交投标" width="600px" :close-on-click-modal="false">
      <el-form ref="bidFormRef" :model="bidForm" :rules="bidRules" label-position="top">
        <el-form-item label="求职信" prop="cover_letter">
          <el-input v-model="bidForm.cover_letter" type="textarea" :rows="4" placeholder="请介绍你的经验和方案" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="报价金额" prop="proposed_budget">
              <el-input-number v-model="bidForm.proposed_budget" :min="1" :max="9999999" class="w-full" placeholder="请输入报价" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="预计工期(天)" prop="estimated_days">
              <el-input-number v-model="bidForm.estimated_days" :min="1" :max="365" class="w-full" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="报价类型" prop="budget_type">
          <el-radio-group v-model="bidForm.budget_type">
            <el-radio value="fixed">固定价</el-radio>
            <el-radio value="hourly">时薪制</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="里程碑计划（选填）">
          <el-input v-model="bidForm.milestone_plan" type="textarea" :rows="3" placeholder="描述你计划的里程碑和交付节点" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBidDialog = false">取消</el-button>
        <el-button type="primary" :loading="bidSubmitting" @click="handleSubmitBid">提交投标</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { projectApi, bidApi, milestoneApi, chatApi, fileApi } from '@/api/modules'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const projectId = computed(() => route.params.id as string)
const loading = ref(false)
const project = ref<any>(null)
const milestones = ref<any[]>([])
const bids = ref<any[]>([])
const bidsLoading = ref(false)
const isBookmarked = ref(false)
const showBidDialog = ref(false)
const bidSubmitting = ref(false)
const bidFormRef = ref<FormInstance>()

const bidForm = reactive({
  cover_letter: '',
  proposed_budget: 0,
  estimated_days: 7,
  budget_type: 'fixed',
  milestone_plan: '',
})

const bidRules: FormRules = {
  cover_letter: [{ required: true, message: '请填写求职信', trigger: 'blur' }],
  proposed_budget: [{ required: true, message: '请输入报价', trigger: 'blur' }],
  estimated_days: [{ required: true, message: '请输入预计工期', trigger: 'blur' }],
  budget_type: [{ required: true, message: '请选择报价类型', trigger: 'change' }],
}

const isProjectOwner = computed(() => {
  return userStore.user && project.value?.client_id === userStore.user.id
})

function statusTagType(status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { published: 'success', draft: 'warning', ongoing: 'primary', completed: 'success', cancelled: 'danger', disputed: 'danger' }
  return map[status] || 'info'
}

function statusLabel(status: string) {
  const map: Record<string, string> = { published: '招标中', draft: '草稿', ongoing: '进行中', completed: '已完成', cancelled: '已取消', disputed: '争议中' }
  return map[status] || status
}

function bidStatusLabel(status: string) {
  const map: Record<string, string> = { submitted: '已提交', shortlisted: '候选中', accepted: '已接受', rejected: '已拒绝', withdrawn: '已撤回' }
  return map[status] || status
}

function milestoneColor(status: string) {
  const map: Record<string, string> = { pending: 'warning', in_progress: 'primary', submitted: 'primary', approved: 'success', rejected: 'danger' }
  return map[status] || 'info'
}

function formatTime(time: string) { return dayjs(time).fromNow() }
function formatDate(time: string) { return dayjs(time).format('YYYY-MM-DD') }

async function fetchProject() {
  loading.value = true
  try {
    const res: any = await projectApi.getById(projectId.value)
    project.value = res.data
    isBookmarked.value = res.data?.is_bookmarked || false
  } catch {} finally {
    loading.value = false
  }
}

async function fetchMilestones() {
  try {
    const res: any = await milestoneApi.listByProject(projectId.value)
    milestones.value = res.data || []
  } catch {}
}

async function fetchBids() {
  if (!isProjectOwner.value) return
  bidsLoading.value = true
  try {
    const res: any = await bidApi.listByProject(projectId.value)
    bids.value = res.data?.items || res.data || []
  } catch {} finally {
    bidsLoading.value = false
  }
}

async function toggleBookmark() {
  try {
    if (isBookmarked.value) {
      await projectApi.removeBookmark(projectId.value)
      isBookmarked.value = false
      ElMessage.success('已取消收藏')
    } else {
      await projectApi.bookmark(projectId.value)
      isBookmarked.value = true
      ElMessage.success('已收藏')
    }
  } catch {}
}

async function handleSubmitBid() {
  const valid = await bidFormRef.value?.validate().catch(() => false)
  if (!valid) return

  bidSubmitting.value = true
  try {
    await bidApi.create(projectId.value, bidForm)
    ElMessage.success('投标提交成功')
    showBidDialog.value = false
    bidFormRef.value?.resetFields()
  } catch (err: any) {
    ElMessage.error(err.message || '投标提交失败')
  } finally {
    bidSubmitting.value = false
  }
}

async function handleAcceptBid(bidId: string) {
  try {
    await ElMessageBox.confirm('确定接受此投标？接受后将创建合同。', '确认接受')
    await bidApi.accept(bidId)
    ElMessage.success('已接受投标')
    fetchBids()
    fetchProject()
  } catch {}
}

async function handleRejectBid(bidId: string) {
  try {
    await ElMessageBox.confirm('确定拒绝此投标？', '确认拒绝')
    await bidApi.reject(bidId)
    ElMessage.success('已拒绝投标')
    fetchBids()
  } catch {}
}

async function handleShortlistBid(bidId: string) {
  try {
    await bidApi.shortlist(bidId)
    ElMessage.success('已加入候选')
    fetchBids()
  } catch {}
}

async function startChat(userId: string) {
  if (!userId) return
  try {
    const res: any = await chatApi.create({ recipient_id: userId, project_id: projectId.value })
    router.push(`/chat/${res.data?.id}`)
  } catch {}
}

async function downloadFile(file: any) {
  try {
    const res: any = await fileApi.download(file.id || file.file_id)
    const blob = new Blob([res as any])
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = file.name || file.filename || 'download'
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('下载失败')
  }
}

onMounted(() => {
  fetchProject()
  fetchMilestones()
  fetchBids()
})
</script>

<style scoped lang="scss">
.project-detail-page {
  .detail-layout {
    display: flex;
    gap: 24px;
  }

  .detail-main {
    flex: 1;
    min-width: 0;
  }

  .detail-header {
    margin-bottom: 24px;

    .detail-header-top {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 12px;
    }

    .detail-title {
      font-size: 26px;
      font-weight: 700;
      color: var(--color-text-primary);
      line-height: 1.3;
    }

    .detail-meta {
      display: flex;
      align-items: center;
      gap: 10px;
      flex-wrap: wrap;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: var(--color-text-secondary);
      }
    }
  }

  .detail-card {
    margin-bottom: 20px;
    border-radius: 12px;

    h3 {
      font-size: 16px;
      font-weight: 600;
    }

    .project-description {
      font-size: 14px;
      line-height: 1.8;
      color: var(--color-text-regular);
      word-break: break-word;
    }

    .tech-tag {
      margin-right: 6px;
      margin-bottom: 4px;
    }
  }

  .attachment-list {
    .attachment-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 10px 0;
      border-bottom: 1px solid var(--border-color-light);
      cursor: pointer;

      &:last-child { border-bottom: none; }

      .attachment-name {
        flex: 1;
        color: var(--color-primary);
      }
    }
  }

  .bid-item {
    padding: 16px 0;
    border-bottom: 1px solid var(--border-color-light);

    &:last-child { border-bottom: none; }

    .bid-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .bid-cover {
      font-size: 14px;
      color: var(--color-text-regular);
      line-height: 1.6;
    }

    .bid-details {
      font-size: 14px;
      color: var(--color-text-secondary);
    }

    .bid-actions {
      display: flex;
      gap: 8px;
    }
  }

  .detail-sidebar {
    width: 300px;
    flex-shrink: 0;

    .sidebar-card {
      margin-bottom: 16px;
      border-radius: 12px;

      h3 {
        font-size: 15px;
        font-weight: 600;
      }
    }

    .client-info {
      text-align: center;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 16px;

      .stat-cell {
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 4px;

        span {
          font-size: 20px;
          font-weight: 600;
          color: var(--color-text-primary);
        }

        label {
          font-size: 12px;
          color: var(--color-text-secondary);
        }
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .project-detail-page {
    .detail-layout {
      flex-direction: column;
    }

    .detail-sidebar {
      width: 100%;
      order: -1;
    }
  }
}
</style>

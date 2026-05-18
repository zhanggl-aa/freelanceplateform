<template>
  <div class="my-projects-page page-container">
    <div class="page-header flex-between">
      <div>
        <h2>我的项目</h2>
        <p>管理你发布的所有项目</p>
      </div>
      <router-link to="/projects/create">
        <el-button type="primary"><el-icon><Plus /></el-icon>发布项目</el-button>
      </router-link>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="已发布" name="published" />
      <el-tab-pane label="竞标中" name="bidding" />
      <el-tab-pane label="进行中" name="ongoing" />
      <el-tab-pane label="已完成" name="completed" />
    </el-tabs>

    <div v-loading="loading" class="project-list mt-16">
      <el-card v-for="project in projects" :key="project.id" shadow="hover" class="project-card card-hover">
        <div class="project-card-inner">
          <div class="project-main">
            <div class="project-header flex-between">
              <h3 class="project-title" @click="goDetail(project.id)">{{ project.title }}</h3>
              <span :class="['status-badge', `status-${project.status}`]">{{ statusLabel(project.status) }}</span>
            </div>
            <p class="project-desc line-clamp-2 mt-8">{{ project.description }}</p>
            <div class="project-tags mt-8">
              <el-tag size="small" effect="plain">{{ project.category_name || project.category?.name || '未分类' }}</el-tag>
              <el-tag size="small" :type="project.budget_type === 'hourly' ? 'warning' : 'primary'" effect="plain">
                {{ project.budget_type === 'hourly' ? '时薪制' : '固定价' }}
              </el-tag>
            </div>
          </div>
          <div class="project-meta">
            <div class="meta-row">
              <el-icon><Money /></el-icon>
              <span v-if="project.budget_type === 'hourly'">¥{{ project.budget_min }}-{{ project.budget_max }}/时</span>
              <span v-else>¥{{ project.budget_min }}-{{ project.budget_max }}</span>
            </div>
            <div class="meta-row">
              <el-icon><EditPen /></el-icon>
              <span>{{ project.bid_count || 0 }} 人投标</span>
            </div>
            <div class="meta-row text-secondary">
              <el-icon><Clock /></el-icon>
              <span>{{ formatDate(project.created_at) }}</span>
            </div>
          </div>
          <div class="project-actions">
            <el-button v-if="project.status === 'draft'" size="small" type="primary" @click="goEdit(project.id)">编辑</el-button>
            <el-button v-if="project.status === 'published' || project.status === 'bidding'" size="small" @click="goBids(project.id)">查看投标</el-button>
            <el-button v-if="project.status === 'ongoing'" size="small" type="primary" @click="goMilestones(project.id)">里程碑</el-button>
            <el-button size="small" @click="goDetail(project.id)">查看详情</el-button>
          </div>
        </div>
      </el-card>

      <el-empty v-if="!loading && projects.length === 0" description="暂无项目" />
    </div>

    <div class="pagination-wrap mt-20" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchProjects"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { projectApi } from '@/api/modules'
import dayjs from 'dayjs'

const router = useRouter()

const loading = ref(false)
const projects = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const activeTab = ref('published')

function statusLabel(status: string) {
  const map: Record<string, string> = {
    draft: '草稿',
    published: '已发布',
    bidding: '竞标中',
    ongoing: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    disputed: '争议中',
  }
  return map[status] || status
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD')
}

function handleTabChange() {
  currentPage.value = 1
  fetchProjects()
}

function goDetail(id: string) {
  router.push(`/projects/${id}`)
}

function goEdit(id: string) {
  router.push(`/projects/${id}/edit`)
}

function goBids(id: string) {
  router.push(`/projects/${id}`)
}

function goMilestones(id: string) {
  router.push(`/projects/${id}`)
}

const apiMap: Record<string, any> = {
  published: projectApi.myPosted,
  bidding: projectApi.myBidding,
  ongoing: projectApi.myWorking,
  completed: projectApi.myCompleted,
}

async function fetchProjects() {
  loading.value = true
  try {
    const apiFn = apiMap[activeTab.value] || projectApi.myPosted
    const res: any = await apiFn({ page: currentPage.value, page_size: pageSize.value })
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
.my-projects-page {
  .project-card {
    margin-bottom: 16px;
    border-radius: 12px;
    border: 1px solid var(--border-color-light);

    .project-card-inner {
      display: flex;
      gap: 24px;
      align-items: flex-start;
    }

    .project-main {
      flex: 1;
      min-width: 0;
    }

    .project-title {
      font-size: 17px;
      font-weight: 600;
      color: var(--color-text-primary);
      cursor: pointer;

      &:hover {
        color: var(--color-primary);
      }
    }

    .project-desc {
      font-size: 14px;
      color: var(--color-text-secondary);
      line-height: 1.6;
    }

    .project-tags {
      display: flex;
      gap: 6px;
    }

    .project-meta {
      width: 160px;
      flex-shrink: 0;
      display: flex;
      flex-direction: column;
      gap: 8px;
      font-size: 13px;
      color: var(--color-text-regular);
      border-left: 1px solid var(--border-color-light);
      padding-left: 20px;

      .meta-row {
        display: flex;
        align-items: center;
        gap: 6px;
      }
    }

    .project-actions {
      display: flex;
      flex-direction: column;
      gap: 8px;
      flex-shrink: 0;
    }
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}

@media screen and (max-width: 768px) {
  .my-projects-page {
    .project-card .project-card-inner {
      flex-direction: column;

      .project-meta {
        width: 100%;
        flex-direction: row;
        flex-wrap: wrap;
        border-left: none;
        padding-left: 0;
        border-top: 1px solid var(--border-color-light);
        padding-top: 12px;
      }

      .project-actions {
        flex-direction: row;
        width: 100%;
      }
    }
  }
}
</style>

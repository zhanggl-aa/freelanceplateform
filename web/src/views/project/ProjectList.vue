<template>
  <div class="project-list-page page-container">
    <div class="page-header">
      <h2>找项目</h2>
      <p>浏览最新项目，找到适合你的机会</p>
    </div>

    <div class="content-layout">
      <!-- Filter Sidebar -->
      <aside class="filter-sidebar">
        <el-card shadow="never" class="filter-card">
          <h3 class="filter-title">筛选条件</h3>

          <div class="filter-section">
            <h4>项目分类</h4>
            <el-tree
              ref="treeRef"
              :data="categoryTree"
              :props="{ label: 'name', children: 'children' }"
              node-key="id"
              highlight-current
              default-expand-all
              @node-click="handleCategoryClick"
            />
          </div>

          <div class="filter-section">
            <h4>预算范围</h4>
            <el-slider
              v-model="budgetRange"
              range
              :min="0"
              :max="100000"
              :step="1000"
              :format-tooltip="(val: number) => `¥${val}`"
              @change="handleFilterChange"
            />
            <div class="range-labels">
              <span>¥{{ budgetRange[0] }}</span>
              <span>¥{{ budgetRange[1] }}</span>
            </div>
          </div>

          <div class="filter-section">
            <h4>项目类型</h4>
            <el-radio-group v-model="filters.budget_type" @change="handleFilterChange">
              <el-radio value="">全部</el-radio>
              <el-radio value="fixed">固定价</el-radio>
              <el-radio value="hourly">时薪制</el-radio>
            </el-radio-group>
          </div>

          <div class="filter-section">
            <h4>技术栈</h4>
            <el-select v-model="filters.tech_stack" multiple filterable allow-create placeholder="选择或输入技术" class="w-full" @change="handleFilterChange">
              <el-option v-for="tech in techOptions" :key="tech" :label="tech" :value="tech" />
            </el-select>
          </div>

          <el-button type="default" class="w-full mt-16" @click="resetFilters">重置筛选</el-button>
        </el-card>
      </aside>

      <!-- Project List -->
      <main class="project-main">
        <div class="sort-bar flex-between mb-16">
          <span class="result-count">共 {{ total }} 个项目</span>
          <el-radio-group v-model="sortBy" size="small" @change="handleSortChange">
            <el-radio-button value="newest">最新发布</el-radio-button>
            <el-radio-button value="budget_high">预算最高</el-radio-button>
            <el-radio-button value="most_bids">投标最多</el-radio-button>
          </el-radio-group>
        </div>

        <div v-loading="loading" class="project-list">
          <el-card v-for="project in projects" :key="project.id" shadow="hover" class="project-item card-hover" @click="goDetail(project.id)">
            <div class="project-item-inner">
              <div class="project-item-main">
                <div class="project-item-header">
                  <h3 class="project-item-title">{{ project.title }}</h3>
                  <el-tag size="small" :type="project.budget_type === 'hourly' ? 'warning' : 'primary'">
                    {{ project.budget_type === 'hourly' ? '时薪制' : '固定价' }}
                  </el-tag>
                </div>
                <p class="project-item-desc line-clamp-2">{{ project.description }}</p>
                <div class="project-item-tags">
                  <el-tag v-for="tech in (project.tech_stack || []).slice(0, 5)" :key="tech" size="small" effect="plain" class="tech-tag">{{ tech }}</el-tag>
                </div>
              </div>
              <div class="project-item-meta">
                <div class="meta-row">
                  <el-icon><Money /></el-icon>
                  <span v-if="project.budget_type === 'hourly'">¥{{ project.budget_min }}-{{ project.budget_max }}/时</span>
                  <span v-else>¥{{ project.budget_min }}-{{ project.budget_max }}</span>
                </div>
                <div class="meta-row">
                  <el-icon><EditPen /></el-icon>
                  <span>{{ project.bid_count || 0 }} 人投标</span>
                </div>
                <div class="meta-row">
                  <el-icon><Calendar /></el-icon>
                  <span>截止 {{ formatDate(project.bid_deadline || project.deadline) }}</span>
                </div>
                <div class="meta-row text-secondary">
                  <el-icon><Clock /></el-icon>
                  <span>{{ formatTime(project.created_at) }}</span>
                </div>
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
            layout="prev, pager, next, jumper"
            @current-change="fetchProjects"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { projectApi, categoryApi } from '@/api/modules'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const projects = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const categoryTree = ref<any[]>([])
const treeRef = ref()
const budgetRange = ref([0, 100000])
const sortBy = ref('newest')

const techOptions = ['Vue.js', 'React', 'Angular', 'Node.js', 'Python', 'Java', 'Go', 'PHP', 'TypeScript', 'Flutter', 'Swift', 'Kotlin', 'C++', 'Rust', 'Docker', 'AWS', 'AI/ML']

const filters = reactive({
  category_id: '',
  budget_type: '',
  tech_stack: [] as string[],
})

function formatTime(time: string) {
  return dayjs(time).fromNow()
}

function formatDate(time: string) {
  return time ? dayjs(time).format('YYYY-MM-DD') : '无期限'
}

function handleCategoryClick(data: any) {
  filters.category_id = data.id
  currentPage.value = 1
  fetchProjects()
}

function handleFilterChange() {
  currentPage.value = 1
  fetchProjects()
}

function handleSortChange() {
  currentPage.value = 1
  fetchProjects()
}

function resetFilters() {
  filters.category_id = ''
  filters.budget_type = ''
  filters.tech_stack = []
  budgetRange.value = [0, 100000]
  sortBy.value = 'newest'
  currentPage.value = 1
  treeRef.value?.setCurrentKey(null)
  fetchProjects()
}

function goDetail(id: string) {
  router.push(`/projects/${id}`)
}

async function fetchCategories() {
  try {
    const res: any = await categoryApi.getTree()
    categoryTree.value = res.data || []
  } catch {}
}

async function fetchProjects() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      budget_min: budgetRange.value[0] || undefined,
      budget_max: budgetRange.value[1] < 100000 ? budgetRange.value[1] : undefined,
      sort: sortBy.value,
    }
    if (filters.category_id) params.category_id = filters.category_id
    if (filters.budget_type) params.budget_type = filters.budget_type
    if (filters.tech_stack.length) params.tech_stack = filters.tech_stack.join(',')

    const keyword = route.query.keyword as string
    if (keyword) params.keyword = keyword
    const catId = route.query.category_id as string
    if (catId) params.category_id = catId

    const res: any = await projectApi.search(params)
    projects.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchProjects()
})
</script>

<style scoped lang="scss">
.project-list-page {
  .content-layout {
    display: flex;
    gap: 24px;
  }

  .filter-sidebar {
    width: 260px;
    flex-shrink: 0;

    .filter-card {
      position: sticky;
      top: calc(var(--header-height) + 24px);
      border-radius: 12px;
    }

    .filter-title {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 16px;
      color: var(--color-text-primary);
    }

    .filter-section {
      margin-bottom: 20px;
      padding-bottom: 16px;
      border-bottom: 1px solid var(--border-color-light);

      &:last-child {
        border-bottom: none;
      }

      h4 {
        font-size: 13px;
        font-weight: 600;
        color: var(--color-text-primary);
        margin-bottom: 10px;
      }
    }

    .range-labels {
      display: flex;
      justify-content: space-between;
      font-size: 12px;
      color: var(--color-text-secondary);
      margin-top: 4px;
    }
  }

  .project-main {
    flex: 1;
    min-width: 0;
  }

  .sort-bar {
    .result-count {
      font-size: 14px;
      color: var(--color-text-secondary);
    }
  }

  .project-item {
    margin-bottom: 16px;
    cursor: pointer;
    border-radius: 12px;
    border: 1px solid var(--border-color-light);

    .project-item-inner {
      display: flex;
      gap: 24px;
    }

    .project-item-main {
      flex: 1;
      min-width: 0;
    }

    .project-item-header {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 8px;
    }

    .project-item-title {
      font-size: 17px;
      font-weight: 600;
      color: var(--color-text-primary);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .project-item-desc {
      font-size: 14px;
      color: var(--color-text-secondary);
      line-height: 1.6;
      margin-bottom: 10px;
    }

    .project-item-tags {
      .tech-tag {
        margin-right: 4px;
        margin-bottom: 4px;
      }
    }

    .project-item-meta {
      width: 180px;
      flex-shrink: 0;
      display: flex;
      flex-direction: column;
      gap: 8px;
      font-size: 13px;
      color: var(--color-text-regular);
      border-left: 1px solid var(--border-color-light);
      padding-left: 24px;

      .meta-row {
        display: flex;
        align-items: center;
        gap: 6px;
      }
    }
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}

@media screen and (max-width: 768px) {
  .project-list-page {
    .content-layout {
      flex-direction: column;
    }

    .filter-sidebar {
      width: 100%;
    }

    .project-item {
      .project-item-inner {
        flex-direction: column;
      }

      .project-item-meta {
        width: 100%;
        flex-direction: row;
        flex-wrap: wrap;
        border-left: none;
        padding-left: 0;
        border-top: 1px solid var(--border-color-light);
        padding-top: 12px;
      }
    }
  }
}
</style>

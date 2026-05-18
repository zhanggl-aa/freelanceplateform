<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">专业开发者自由职业平台</h1>
        <p class="hero-subtitle">连接优质项目与顶尖开发者，让技术创造更大价值</p>
        <div class="hero-search">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索项目关键词..."
            size="large"
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="searchCategory"
            placeholder="选择分类"
            size="large"
            class="search-category"
            clearable
          >
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
          <el-button type="primary" size="large" @click="handleSearch" class="search-btn">
            <el-icon><Search /></el-icon>搜索项目
          </el-button>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-num">{{ stats.projectCount }}+</span>
            <span class="stat-label">在线项目</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">{{ stats.developerCount }}+</span>
            <span class="stat-label">注册开发者</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">{{ stats.totalPaid }}万</span>
            <span class="stat-label">累计交易</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Category Section -->
    <section class="section category-section">
      <div class="section-inner">
        <h2 class="section-title">热门分类</h2>
        <el-row :gutter="20">
          <el-col :xs="12" :sm="8" :md="6" v-for="cat in displayCategories" :key="cat.id">
            <el-card shadow="hover" class="category-card card-hover" @click="goCategory(cat.id)">
              <div class="category-icon" :style="{ background: cat.bgColor }">
                <el-icon :size="32" :color="cat.color"><component :is="cat.icon" /></el-icon>
              </div>
              <h3 class="category-name">{{ cat.name }}</h3>
              <p class="category-desc">{{ cat.desc }}</p>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </section>

    <!-- Featured Projects -->
    <section class="section featured-section">
      <div class="section-inner">
        <div class="flex-between mb-20">
          <h2 class="section-title" style="margin-bottom:0">精选项目</h2>
          <router-link to="/projects">
            <el-button text type="primary">查看更多 <el-icon><ArrowRight /></el-icon></el-button>
          </router-link>
        </div>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" v-for="project in featuredProjects" :key="project.id">
            <el-card shadow="hover" class="project-card card-hover" @click="goProject(project.id)">
              <div class="project-header">
                <el-tag size="small" :type="budgetTypeTag(project.budget_type)">{{ project.budget_type === 'hourly' ? '时薪制' : '固定价' }}</el-tag>
                <span class="project-time">{{ formatTime(project.created_at) }}</span>
              </div>
              <h3 class="project-title">{{ project.title }}</h3>
              <p class="project-desc line-clamp-2">{{ project.description }}</p>
              <div class="project-tags">
                <el-tag v-for="tech in (project.tech_stack || []).slice(0, 3)" :key="tech" size="small" effect="plain" class="tech-tag">{{ tech }}</el-tag>
              </div>
              <div class="project-footer">
                <span class="project-budget">
                  <template v-if="project.budget_type === 'hourly'">
                    ¥{{ project.budget_min }}-{{ project.budget_max }}/时
                  </template>
                  <template v-else>
                    ¥{{ project.budget_min }}-{{ project.budget_max }}
                  </template>
                </span>
                <span class="project-bids">{{ project.bid_count || 0 }}人投标</span>
              </div>
            </el-card>
          </el-col>
        </el-row>
        <el-empty v-if="!loading && featuredProjects.length === 0" description="暂无精选项目" />
      </div>
    </section>

    <!-- Top Developers -->
    <section class="section developer-section">
      <div class="section-inner">
        <div class="flex-between mb-20">
          <h2 class="section-title" style="margin-bottom:0">优秀开发者</h2>
          <router-link to="/developers">
            <el-button text type="primary">查看更多 <el-icon><ArrowRight /></el-icon></el-button>
          </router-link>
        </div>
        <el-row :gutter="20">
          <el-col :xs="12" :sm="8" :md="6" v-for="dev in topDevelopers" :key="dev.id">
            <el-card shadow="hover" class="developer-card card-hover" @click="goDeveloper(dev.id)">
              <div class="developer-avatar-wrap">
                <el-avatar :size="64" :src="dev.avatar_url" icon="UserFilled" />
              </div>
              <h3 class="developer-name">{{ dev.nickname }}</h3>
              <p class="developer-title">{{ dev.title || '全栈开发者' }}</p>
              <div class="developer-skills">
                <el-tag v-for="skill in (dev.skills || []).slice(0, 3)" :key="skill.name || skill" size="small" effect="plain" class="skill-tag">{{ skill.name || skill }}</el-tag>
              </div>
              <div class="developer-stats">
                <div class="dev-stat">
                  <el-icon color="#f7ba2a"><StarFilled /></el-icon>
                  <span>{{ dev.rating_avg?.toFixed(1) || '5.0' }}</span>
                </div>
                <div class="dev-stat">
                  <el-icon color="#67C23A"><CircleCheck /></el-icon>
                  <span>{{ dev.completed_projects || 0 }}项目</span>
                </div>
                <div class="dev-stat" v-if="dev.hourly_rate">
                  <span class="rate">¥{{ dev.hourly_rate }}/时</span>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
        <el-empty v-if="!devLoading && topDevelopers.length === 0" description="暂无开发者" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { projectApi, developerApi, categoryApi } from '@/api/modules'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()

const searchKeyword = ref('')
const searchCategory = ref('')
const categories = ref<any[]>([])
const featuredProjects = ref<any[]>([])
const topDevelopers = ref<any[]>([])
const loading = ref(false)
const devLoading = ref(false)

const stats = ref({
  projectCount: 5000,
  developerCount: 12000,
  totalPaid: 3800,
})

const defaultCategories = [
  { id: 'web', name: '网站开发', icon: 'Monitor', color: '#409EFF', bgColor: '#ecf5ff', desc: '企业官网/商城/平台' },
  { id: 'mobile', name: '移动开发', icon: 'Iphone', color: '#67C23A', bgColor: '#f0f9eb', desc: 'iOS/Android/小程序' },
  { id: 'backend', name: '后端开发', icon: 'Cpu', color: '#E6A23C', bgColor: '#fdf6ec', desc: 'API/微服务/架构' },
  { id: 'ai', name: 'AI/数据', icon: 'MagicStick', color: '#9b59b6', bgColor: '#f5eef8', desc: '机器学习/NLP/数据分析' },
  { id: 'design', name: 'UI/UX设计', icon: 'Brush', color: '#F56C6C', bgColor: '#fef0f0', desc: '界面设计/交互设计' },
  { id: 'devops', name: '运维部署', icon: 'Setting', color: '#34495e', bgColor: '#eef1f5', desc: '云服务/容器/自动化' },
  { id: 'game', name: '游戏开发', icon: 'VideoPlay', color: '#e74c3c', bgColor: '#fde8e8', desc: 'Unity/UE/小游戏' },
  { id: 'security', name: '安全测试', icon: 'Lock', color: '#1abc9c', bgColor: '#e8f8f5', desc: '渗透测试/代码审计' },
]

const displayCategories = computed(() => categories.value.length > 0 ? categories.value : defaultCategories)

function budgetTypeTag(type: string) {
  return type === 'hourly' ? 'warning' : 'primary'
}

function formatTime(time: string) {
  return dayjs(time).fromNow()
}

function handleSearch() {
  router.push({
    path: '/projects',
    query: {
      keyword: searchKeyword.value || undefined,
      category_id: searchCategory.value || undefined,
    },
  })
}

function goCategory(id: string) {
  router.push({ path: '/projects', query: { category_id: id } })
}

function goProject(id: string) {
  router.push(`/projects/${id}`)
}

function goDeveloper(id: string) {
  router.push(`/developers/${id}`)
}

async function fetchCategories() {
  try {
    const res: any = await categoryApi.getTree()
    categories.value = res.data || []
  } catch {}
}

async function fetchFeaturedProjects() {
  loading.value = true
  try {
    const res: any = await projectApi.search({ page: 1, page_size: 6, status: 'published' })
    featuredProjects.value = res.data?.items || res.data || []
  } catch {} finally {
    loading.value = false
  }
}

async function fetchTopDevelopers() {
  devLoading.value = true
  try {
    const res: any = await developerApi.search({ page: 1, page_size: 8, sort: 'rating' })
    topDevelopers.value = res.data?.items || res.data || []
  } catch {} finally {
    devLoading.value = false
  }
}

onMounted(() => {
  fetchCategories()
  fetchFeaturedProjects()
  fetchTopDevelopers()
})
</script>

<style scoped lang="scss">
.home-page {
  background: var(--bg-color);
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #409EFF 50%, #764ba2 100%);
  padding: 80px 24px 60px;
  color: #fff;

  .hero-content {
    max-width: 800px;
    margin: 0 auto;
    text-align: center;
  }

  .hero-title {
    font-size: 40px;
    font-weight: 700;
    margin-bottom: 12px;
    letter-spacing: 2px;
  }

  .hero-subtitle {
    font-size: 18px;
    opacity: 0.9;
    margin-bottom: 40px;
  }

  .hero-search {
    display: flex;
    gap: 0;
    max-width: 680px;
    margin: 0 auto 40px;
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);

    .search-input {
      flex: 1;
      :deep(.el-input__wrapper) {
        box-shadow: none;
        border-radius: 0;
        font-size: 16px;
      }
    }

    .search-category {
      width: 160px;
      :deep(.el-input__wrapper) {
        box-shadow: none;
        border-radius: 0;
        border-left: 1px solid var(--border-color-light);
      }
    }

    .search-btn {
      border-radius: 0;
      padding: 0 28px;
      font-size: 16px;
    }
  }

  .hero-stats {
    display: flex;
    justify-content: center;
    gap: 60px;

    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: center;

      .stat-num {
        font-size: 32px;
        font-weight: 700;
      }

      .stat-label {
        font-size: 14px;
        opacity: 0.8;
        margin-top: 4px;
      }
    }
  }
}

.section {
  padding: 48px 24px;

  .section-inner {
    max-width: 1200px;
    margin: 0 auto;
  }

  .section-title {
    font-size: 24px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 24px;
    position: relative;
    padding-left: 16px;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 4px;
      bottom: 4px;
      width: 4px;
      background: var(--color-primary);
      border-radius: 2px;
    }
  }
}

.category-section {
  background: var(--bg-color-white);

  .category-card {
    text-align: center;
    cursor: pointer;
    margin-bottom: 20px;
    border-radius: 12px;
    border: 1px solid var(--border-color-light);

    :deep(.el-card__body) {
      padding: 24px 16px;
    }

    .category-icon {
      width: 72px;
      height: 72px;
      border-radius: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 12px;
    }

    .category-name {
      font-size: 16px;
      font-weight: 600;
      color: var(--color-text-primary);
      margin-bottom: 4px;
    }

    .category-desc {
      font-size: 12px;
      color: var(--color-text-secondary);
    }
  }
}

.featured-section {
  .project-card {
    cursor: pointer;
    margin-bottom: 20px;
    border-radius: 12px;
    border: 1px solid var(--border-color-light);

    .project-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;

      .project-time {
        font-size: 12px;
        color: var(--color-text-secondary);
      }
    }

    .project-title {
      font-size: 16px;
      font-weight: 600;
      color: var(--color-text-primary);
      margin-bottom: 8px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .project-desc {
      font-size: 13px;
      color: var(--color-text-secondary);
      margin-bottom: 12px;
      line-height: 1.5;
      min-height: 39px;
    }

    .project-tags {
      margin-bottom: 12px;

      .tech-tag {
        margin-right: 4px;
        margin-bottom: 4px;
      }
    }

    .project-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-top: 12px;
      border-top: 1px solid var(--border-color-light);

      .project-budget {
        font-size: 16px;
        font-weight: 600;
        color: var(--color-danger);
      }

      .project-bids {
        font-size: 13px;
        color: var(--color-text-secondary);
      }
    }
  }
}

.developer-section {
  background: var(--bg-color-white);

  .developer-card {
    text-align: center;
    cursor: pointer;
    margin-bottom: 20px;
    border-radius: 12px;
    border: 1px solid var(--border-color-light);

    :deep(.el-card__body) {
      padding: 24px 16px;
    }

    .developer-avatar-wrap {
      margin-bottom: 12px;
    }

    .developer-name {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 4px;
    }

    .developer-title {
      font-size: 13px;
      color: var(--color-text-secondary);
      margin-bottom: 10px;
    }

    .developer-skills {
      margin-bottom: 12px;
      min-height: 28px;

      .skill-tag {
        margin: 0 2px 4px;
      }
    }

    .developer-stats {
      display: flex;
      justify-content: center;
      gap: 16px;
      font-size: 13px;
      color: var(--color-text-secondary);

      .dev-stat {
        display: flex;
        align-items: center;
        gap: 4px;
      }

      .rate {
        color: var(--color-danger);
        font-weight: 600;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .hero-section {
    padding: 48px 16px 40px;

    .hero-title {
      font-size: 24px;
    }

    .hero-subtitle {
      font-size: 14px;
      margin-bottom: 24px;
    }

    .hero-search {
      flex-direction: column;
      border-radius: 8px;

      .search-input,
      .search-category {
        width: 100%;
      }

      .search-btn {
        border-radius: 0 0 8px 8px;
      }
    }

    .hero-stats {
      gap: 24px;

      .stat-num {
        font-size: 22px;
      }
    }
  }
}
</style>

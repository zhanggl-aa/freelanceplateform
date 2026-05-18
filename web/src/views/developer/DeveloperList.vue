<template>
  <div class="developer-list-page page-container">
    <div class="page-header">
      <h2>找开发者</h2>
      <p>浏览优秀开发者，找到合适的合作伙伴</p>
    </div>

    <div class="content-layout">
      <!-- Filter Sidebar -->
      <aside class="filter-sidebar">
        <el-card shadow="never" class="filter-card">
          <h3 class="filter-title">筛选条件</h3>

          <div class="filter-section">
            <h4>技能搜索</h4>
            <el-input v-model="filters.skill" placeholder="输入技能关键词" clearable @clear="handleFilterChange" @keyup.enter="handleFilterChange">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
          </div>

          <div class="filter-section">
            <h4>时薪范围（元/时）</h4>
            <el-slider
              v-model="rateRange"
              range
              :min="0"
              :max="500"
              :step="10"
              :format-tooltip="(val: number) => `¥${val}/时`"
              @change="handleFilterChange"
            />
            <div class="range-labels">
              <span>¥{{ rateRange[0] }}</span>
              <span>¥{{ rateRange[1] }}</span>
            </div>
          </div>

          <div class="filter-section">
            <h4>可用状态</h4>
            <el-select v-model="filters.availability" placeholder="选择可用状态" clearable class="w-full" @change="handleFilterChange">
              <el-option label="全职可用" value="full_time" />
              <el-option label="兼职可用" value="part_time" />
              <el-option label="随时可用" value="available" />
              <el-option label="暂不可用" value="unavailable" />
            </el-select>
          </div>

          <el-button type="default" class="w-full mt-16" @click="resetFilters">重置筛选</el-button>
        </el-card>
      </aside>

      <!-- Developer List -->
      <main class="developer-main">
        <div class="sort-bar flex-between mb-16">
          <span class="result-count">共 {{ total }} 位开发者</span>
          <el-radio-group v-model="sortBy" size="small" @change="handleSortChange">
            <el-radio-button value="rating">评分最高</el-radio-button>
            <el-radio-button value="projects">项目最多</el-radio-button>
            <el-radio-button value="newest">最新注册</el-radio-button>
          </el-radio-group>
        </div>

        <div v-loading="loading">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" v-for="dev in developers" :key="dev.id">
              <el-card shadow="hover" class="developer-card card-hover" @click="goDetail(dev.id)">
                <div class="card-top">
                  <el-avatar :size="64" :src="dev.avatar_url" icon="UserFilled" />
                  <h3 class="dev-name">{{ dev.nickname }}</h3>
                  <p class="dev-title">{{ dev.title || '开发者' }}</p>
                </div>
                <div class="dev-skills">
                  <el-tag v-for="skill in (dev.skills || []).slice(0, 4)" :key="skill.name || skill" size="small" effect="plain" class="skill-tag">
                    {{ skill.name || skill }}
                  </el-tag>
                </div>
                <div class="dev-stats flex-between">
                  <div class="stat-item">
                    <span class="stat-value">¥{{ dev.hourly_rate || 0 }}/时</span>
                  </div>
                  <div class="stat-item">
                    <el-rate v-model="dev.rating_avg" disabled :colors="['#f7ba2a', '#f7ba2a', '#f7ba2a']" size="small" />
                  </div>
                  <div class="stat-item">
                    <span class="stat-label">{{ dev.completed_projects || 0 }}项目</span>
                  </div>
                </div>
                <el-button type="primary" plain class="w-full mt-12" @click.stop="goDetail(dev.id)">查看详情</el-button>
              </el-card>
            </el-col>
          </el-row>

          <el-empty v-if="!loading && developers.length === 0" description="暂无开发者" />
        </div>

        <div class="pagination-wrap mt-20" v-if="total > pageSize">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="total"
            layout="prev, pager, next, jumper"
            @current-change="fetchDevelopers"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { developerApi } from '@/api/modules'

const router = useRouter()

const loading = ref(false)
const developers = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const rateRange = ref([0, 500])
const sortBy = ref('rating')

const filters = reactive({
  skill: '',
  availability: '',
})

function handleFilterChange() {
  currentPage.value = 1
  fetchDevelopers()
}

function handleSortChange() {
  currentPage.value = 1
  fetchDevelopers()
}

function resetFilters() {
  filters.skill = ''
  filters.availability = ''
  rateRange.value = [0, 500]
  sortBy.value = 'rating'
  currentPage.value = 1
  fetchDevelopers()
}

function goDetail(id: string) {
  // The search result id is developer_profiles.id, but the route expects user_id
  // The enriched search result also has user_id field
  router.push(`/developers/${id}`)
}

async function fetchDevelopers() {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    if (filters.skill) params.skills = filters.skill
    if (filters.availability) params.availability = filters.availability
    if (rateRange.value[0] > 0) params.min_rate = rateRange.value[0]
    if (rateRange.value[1] < 500) params.max_rate = rateRange.value[1]

    const res: any = await developerApi.search(params)
    developers.value = res.data?.items || res.data || []
    total.value = res.data?.total || 0
  } catch {} finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDevelopers()
})
</script>

<style scoped lang="scss">
.developer-list-page {
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

  .developer-main {
    flex: 1;
    min-width: 0;
  }

  .sort-bar {
    .result-count {
      font-size: 14px;
      color: var(--color-text-secondary);
    }
  }

  .developer-card {
    margin-bottom: 20px;
    cursor: pointer;
    border-radius: 12px;
    border: 1px solid var(--border-color-light);

    :deep(.el-card__body) {
      padding: 20px;
    }

    .card-top {
      text-align: center;
      margin-bottom: 12px;
    }

    .dev-name {
      font-size: 16px;
      font-weight: 600;
      color: var(--color-text-primary);
      margin-top: 10px;
      margin-bottom: 4px;
    }

    .dev-title {
      font-size: 13px;
      color: var(--color-text-secondary);
    }

    .dev-skills {
      display: flex;
      flex-wrap: wrap;
      justify-content: center;
      gap: 4px;
      margin-bottom: 12px;
      min-height: 28px;

      .skill-tag {
        margin: 0;
      }
    }

    .dev-stats {
      padding-top: 12px;
      border-top: 1px solid var(--border-color-light);
      align-items: center;

      .stat-value {
        font-size: 15px;
        font-weight: 600;
        color: var(--color-danger);
      }

      .stat-label {
        font-size: 13px;
        color: var(--color-text-secondary);
      }
    }
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
  }
}

@media screen and (max-width: 768px) {
  .developer-list-page {
    .content-layout {
      flex-direction: column;
    }

    .filter-sidebar {
      width: 100%;
    }
  }
}
</style>

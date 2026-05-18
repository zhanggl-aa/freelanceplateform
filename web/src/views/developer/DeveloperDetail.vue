<template>
  <div class="developer-detail-page page-container" v-loading="loading">
    <template v-if="developer">
      <!-- Hero Header -->
      <div class="detail-hero">
        <div class="hero-bg"></div>
        <div class="hero-content">
          <el-avatar :size="100" :src="developer.avatar_url" icon="UserFilled" class="hero-avatar" />
          <div class="hero-info">
            <div class="hero-top">
              <h1 class="dev-name">{{ developer.nickname }}</h1>
              <el-tag v-if="developer.verified" type="success" effect="dark" size="small" class="verified-tag">
                <el-icon><CircleCheckFilled /></el-icon> 已认证
              </el-tag>
            </div>
            <p class="dev-title">{{ developer.title || '全栈开发者' }}</p>
            <p class="dev-bio">{{ developer.bio || '这位开发者很低调，还没有填写简介' }}</p>
            <div class="hero-actions mt-12">
              <el-button v-if="userStore.isClient" type="primary" size="large" round @click="handleInvite">
                <el-icon><Promotion /></el-icon> 邀请合作
              </el-button>
              <el-button size="large" round @click="handleMessage">
                <el-icon><ChatDotRound /></el-icon> 发消息
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- Stats Row -->
      <el-row :gutter="16" class="stats-row">
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #f7ba2a, #f5a623);">
              <el-icon :size="22" color="#fff"><StarFilled /></el-icon>
            </div>
            <div class="stat-text">
              <span class="stat-value">{{ developer.rating_avg?.toFixed(1) || '0.0' }}</span>
              <span class="stat-label">评分</span>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #67C23A, #43e97b);">
              <el-icon :size="22" color="#fff"><CircleCheck /></el-icon>
            </div>
            <div class="stat-text">
              <span class="stat-value">{{ developer.completed_projects || 0 }}</span>
              <span class="stat-label">已完成</span>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #409EFF, #667eea);">
              <el-icon :size="22" color="#fff"><Timer /></el-icon>
            </div>
            <div class="stat-text">
              <span class="stat-value">{{ developer.experience_years || 0 }}年</span>
              <span class="stat-label">经验</span>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(135deg, #F56C6C, #fa709a);">
              <el-icon :size="22" color="#fff"><Money /></el-icon>
            </div>
            <div class="stat-text">
              <span class="stat-value">¥{{ developer.hourly_rate || 0 }}/时</span>
              <span class="stat-label">时薪</span>
            </div>
          </div>
        </el-col>
      </el-row>

      <div class="detail-layout">
        <!-- Main Content -->
        <main class="detail-main">
          <!-- Skills -->
          <div class="section-card">
            <div class="section-header">
              <el-icon :size="20" color="#667eea"><CollectionTag /></el-icon>
              <h3>技能标签</h3>
            </div>
            <div class="skills-cloud" v-if="developer.skills && developer.skills.length > 0">
              <el-tag
                v-for="skill in developer.skills"
                :key="skill.id || skill.name || skill"
                size="large"
                effect="plain"
                class="skill-tag"
                round
              >
                {{ skill.name || skill }}
              </el-tag>
            </div>
            <div v-else class="empty-state">
              <el-icon :size="40" color="#c0c4cc"><CollectionTag /></el-icon>
              <p>暂无技能标签</p>
            </div>
          </div>

          <!-- Portfolio - Highlight Section -->
          <div class="section-card portfolio-section">
            <div class="section-header">
              <el-icon :size="20" color="#667eea"><PictureFilled /></el-icon>
              <h3>个人作品</h3>
              <span class="section-count" v-if="developer.portfolio">{{ developer.portfolio.length }} 个作品</span>
            </div>

            <div class="portfolio-grid" v-if="developer.portfolio && developer.portfolio.length > 0">
              <div
                v-for="item in developer.portfolio"
                :key="item.id"
                class="portfolio-card"
                @click="handlePortfolioClick(item)"
              >
                <div class="portfolio-cover">
                  <el-image
                    v-if="item.image_url || (item.image_urls && item.image_urls.length > 0)"
                    :src="item.image_url || item.image_urls[0]"
                    fit="cover"
                    class="portfolio-img"
                  >
                    <template #error>
                      <div class="img-placeholder">
                        <el-icon :size="36" color="#c0c4cc"><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                  <div v-else class="img-placeholder">
                    <el-icon :size="36" color="#c0c4cc"><Picture /></el-icon>
                  </div>
                  <div class="portfolio-overlay">
                    <el-icon :size="24" color="#fff"><ZoomIn /></el-icon>
                  </div>
                </div>
                <div class="portfolio-body">
                  <h4 class="portfolio-title">{{ item.title }}</h4>
                  <p class="portfolio-desc line-clamp-2">{{ item.description || '暂无描述' }}</p>
                  <div class="portfolio-tech" v-if="item.tech_stack && item.tech_stack.length > 0">
                    <el-tag v-for="tech in item.tech_stack.slice(0, 3)" :key="tech" size="small" effect="light" type="info">
                      {{ tech }}
                    </el-tag>
                    <span v-if="item.tech_stack.length > 3" class="more-tech">+{{ item.tech_stack.length - 3 }}</span>
                  </div>
                  <a v-if="item.project_url || item.url" :href="item.project_url || item.url" target="_blank" class="portfolio-link" @click.stop>
                    <el-icon><Link /></el-icon> 查看项目
                  </a>
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              <el-icon :size="40" color="#c0c4cc"><PictureFilled /></el-icon>
              <p>暂无作品展示</p>
            </div>
          </div>

          <!-- Reviews -->
          <div class="section-card">
            <div class="section-header">
              <el-icon :size="20" color="#667eea"><ChatLineSquare /></el-icon>
              <h3>客户评价</h3>
              <span class="section-count" v-if="reviews">{{ reviews.length }} 条</span>
            </div>
            <div v-loading="reviewsLoading">
              <div v-for="review in reviews" :key="review.id" class="review-item">
                <div class="review-header">
                  <div class="reviewer">
                    <el-avatar :size="36" :src="review.reviewer?.avatar_url" icon="UserFilled" />
                    <div class="reviewer-info">
                      <span class="reviewer-name">{{ review.reviewer?.nickname || '用户' }}</span>
                      <el-rate v-model="review.rating" disabled size="small" />
                    </div>
                  </div>
                  <span class="review-time">{{ formatDate(review.created_at) }}</span>
                </div>
                <p class="review-content">{{ review.comment }}</p>
              </div>
              <div v-if="!reviewsLoading && reviews.length === 0" class="empty-state">
                <el-icon :size="40" color="#c0c4cc"><ChatLineSquare /></el-icon>
                <p>暂无评价</p>
              </div>
            </div>
          </div>
        </main>

        <!-- Sidebar -->
        <aside class="detail-sidebar">
          <div class="sidebar-card">
            <h3 class="sidebar-title">个人信息</h3>
            <div class="info-list">
              <div class="info-item">
                <span class="info-label">可用状态</span>
                <el-tag :type="availabilityType(developer.availability)" size="small" effect="dark" round>
                  {{ availabilityLabel(developer.availability) }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="info-label">经验年限</span>
                <span class="info-value">{{ developer.experience_years || 0 }} 年</span>
              </div>
              <div class="info-item">
                <span class="info-label">时薪</span>
                <span class="info-value highlight">¥{{ developer.hourly_rate || 0 }}/时</span>
              </div>
              <div class="info-item" v-if="developer.github_url">
                <span class="info-label">GitHub</span>
                <a :href="developer.github_url" target="_blank" class="social-link">
                  <el-icon><Link /></el-icon> 访问
                </a>
              </div>
              <div class="info-item" v-if="developer.linkedin_url">
                <span class="info-label">LinkedIn</span>
                <a :href="developer.linkedin_url" target="_blank" class="social-link">
                  <el-icon><Link /></el-icon> 访问
                </a>
              </div>
              <div class="info-item" v-if="developer.website_url">
                <span class="info-label">个人网站</span>
                <a :href="developer.website_url" target="_blank" class="social-link">
                  <el-icon><Link /></el-icon> 访问
                </a>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </template>

    <!-- Portfolio Lightbox -->
    <el-dialog v-model="showLightbox" :title="lightboxItem?.title || ''" width="600px" class="portfolio-dialog">
      <div class="lightbox-content">
        <el-image
          v-if="lightboxItem?.image_url || (lightboxItem?.image_urls && lightboxItem.image_urls.length > 0)"
          :src="lightboxItem?.image_url || lightboxItem?.image_urls?.[0]"
          fit="contain"
          style="width: 100%; max-height: 400px;"
        />
        <p class="lightbox-desc">{{ lightboxItem?.description }}</p>
        <div class="lightbox-tech" v-if="lightboxItem?.tech_stack && lightboxItem.tech_stack.length > 0">
          <el-tag v-for="tech in lightboxItem.tech_stack" :key="tech" size="small">{{ tech }}</el-tag>
        </div>
        <a v-if="lightboxItem?.project_url || lightboxItem?.url" :href="lightboxItem.project_url || lightboxItem.url" target="_blank" class="portfolio-link">
          <el-icon><Link /></el-icon> 查看项目
        </a>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { developerApi, reviewApi, chatApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const developerId = computed(() => route.params.id as string)
const loading = ref(false)
const reviewsLoading = ref(false)
const developer = ref<any>(null)
const reviews = ref<any[]>([])
const showLightbox = ref(false)
const lightboxItem = ref<any>(null)

function availabilityLabel(status: string) {
  const map: Record<string, string> = { full_time: '全职可用', part_time: '兼职可用', available: '随时可用', unavailable: '暂不可用' }
  return map[status] || status || '未知'
}

function availabilityType(status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { full_time: 'success', part_time: 'primary', available: 'success', unavailable: 'danger' }
  return map[status] || 'info'
}

function formatDate(time: string) {
  return dayjs(time).format('YYYY-MM-DD')
}

function handlePortfolioClick(item: any) {
  lightboxItem.value = item
  showLightbox.value = true
}

function handleMessage() {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  handleInvite()
}

async function handleInvite() {
  try {
    const res: any = await chatApi.create({ recipient_id: developer.value.user_id || developerId.value })
    router.push(`/chat/${res.data?.id}`)
  } catch {
    ElMessage.error('发起对话失败')
  }
}

async function fetchDeveloper() {
  loading.value = true
  try {
    const res: any = await developerApi.getDeveloper(developerId.value)
    developer.value = res.data
  } catch {} finally {
    loading.value = false
  }
}

async function fetchReviews() {
  reviewsLoading.value = true
  try {
    const res: any = await reviewApi.getByUser(developerId.value)
    reviews.value = res.data?.items || res.data || []
  } catch {} finally {
    reviewsLoading.value = false
  }
}

onMounted(() => {
  fetchDeveloper()
  fetchReviews()
})
</script>

<style scoped lang="scss">
.developer-detail-page {
  .detail-hero {
    position: relative;
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 24px;

    .hero-bg {
      position: absolute;
      inset: 0;
      background: linear-gradient(135deg, #0f0c29, #302b63, #24243e);
      z-index: 0;
    }

    .hero-content {
      position: relative;
      z-index: 1;
      padding: 32px;
      display: flex;
      gap: 24px;
      align-items: flex-start;
    }

    .hero-avatar {
      border: 4px solid rgba(255, 255, 255, 0.2);
      box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
      flex-shrink: 0;
    }

    .hero-info {
      color: #fff;
    }

    .hero-top {
      display: flex;
      align-items: center;
      gap: 10px;
    }

    .dev-name {
      font-size: 28px;
      font-weight: 700;
      margin: 0;
    }

    .verified-tag {
      display: inline-flex;
      align-items: center;
      gap: 4px;
    }

    .dev-title {
      font-size: 16px;
      color: rgba(255, 255, 255, 0.7);
      margin-top: 4px;
    }

    .dev-bio {
      font-size: 14px;
      color: rgba(255, 255, 255, 0.6);
      line-height: 1.6;
      max-width: 600px;
      margin-top: 8px;
    }
  }

  .stats-row {
    margin-bottom: 24px;
  }

  .stat-card {
    background: #fff;
    border-radius: 14px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 14px;
    border: 1px solid #f0f0f5;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
    transition: transform 0.2s;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
    }

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .stat-text {
      display: flex;
      flex-direction: column;
    }

    .stat-value {
      font-size: 20px;
      font-weight: 700;
      color: #1a1a2e;
    }

    .stat-label {
      font-size: 12px;
      color: #8c8c9a;
    }
  }

  .detail-layout {
    display: flex;
    gap: 24px;
  }

  .detail-main {
    flex: 1;
    min-width: 0;
  }

  .section-card {
    background: #fff;
    border-radius: 14px;
    padding: 24px;
    margin-bottom: 20px;
    border: 1px solid #f0f0f5;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 20px;

    h3 {
      font-size: 17px;
      font-weight: 700;
      color: #1a1a2e;
      margin: 0;
      flex: 1;
    }

    .section-count {
      font-size: 13px;
      color: #8c8c9a;
      background: #f7f8fa;
      padding: 2px 10px;
      border-radius: 12px;
    }
  }

  .skills-cloud {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;

    .skill-tag {
      font-size: 14px;
      padding: 4px 14px;
      border-radius: 20px;
      border-color: #e0e0e8;
      color: #5a5a6e;

      &:hover {
        color: #667eea;
        border-color: #667eea;
        background: rgba(102, 126, 234, 0.05);
      }
    }
  }

  .empty-state {
    text-align: center;
    padding: 32px;
    color: #c0c4cc;

    p {
      margin-top: 8px;
      font-size: 14px;
    }
  }

  // Portfolio grid
  .portfolio-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .portfolio-card {
    border: 1px solid #f0f0f5;
    border-radius: 14px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s;
    background: #fff;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1);
      border-color: #667eea;
    }

    .portfolio-cover {
      position: relative;
      height: 180px;
      overflow: hidden;
      background: #f7f8fa;

      .portfolio-img {
        width: 100%;
        height: 100%;
      }

      .img-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #f7f8fa, #e8e8ed);
      }

      .portfolio-overlay {
        position: absolute;
        inset: 0;
        background: rgba(0, 0, 0, 0.4);
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transition: opacity 0.3s;
      }

      &:hover .portfolio-overlay {
        opacity: 1;
      }
    }

    .portfolio-body {
      padding: 16px;
    }

    .portfolio-title {
      font-size: 16px;
      font-weight: 700;
      color: #1a1a2e;
      margin: 0 0 6px;
    }

    .portfolio-desc {
      font-size: 13px;
      color: #8c8c9a;
      line-height: 1.5;
      margin: 0 0 10px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .portfolio-tech {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
      margin-bottom: 8px;

      .more-tech {
        font-size: 12px;
        color: #8c8c9a;
        line-height: 22px;
      }
    }

    .portfolio-link {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      font-size: 13px;
      color: #667eea;
      text-decoration: none;
      font-weight: 500;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .review-item {
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f5;

    &:last-child {
      border-bottom: none;
    }

    .review-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .reviewer {
      display: flex;
      align-items: center;
      gap: 10px;
    }

    .reviewer-info {
      display: flex;
      flex-direction: column;
    }

    .reviewer-name {
      font-size: 14px;
      font-weight: 600;
      color: #1a1a2e;
    }

    .review-time {
      font-size: 12px;
      color: #c0c4cc;
    }

    .review-content {
      font-size: 14px;
      color: #5a5a6e;
      line-height: 1.6;
      padding-left: 46px;
      margin-top: 8px;
    }
  }

  .detail-sidebar {
    width: 300px;
    flex-shrink: 0;
  }

  .sidebar-card {
    background: #fff;
    border-radius: 14px;
    padding: 24px;
    border: 1px solid #f0f0f5;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
    position: sticky;
    top: calc(64px + 24px);
  }

  .sidebar-title {
    font-size: 16px;
    font-weight: 700;
    color: #1a1a2e;
    margin: 0 0 16px;
  }

  .info-list {
    .info-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 10px 0;
      border-bottom: 1px solid #f7f8fa;

      &:last-child {
        border-bottom: none;
      }

      .info-label {
        font-size: 13px;
        color: #8c8c9a;
      }

      .info-value {
        font-size: 14px;
        font-weight: 600;
        color: #1a1a2e;

        &.highlight {
          color: #F56C6C;
        }
      }

      .social-link {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        color: #667eea;
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}

.lightbox-content {
  .lightbox-desc {
    font-size: 14px;
    color: #5a5a6e;
    line-height: 1.6;
    margin-top: 16px;
  }

  .lightbox-tech {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-top: 12px;
  }

  .portfolio-link {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 14px;
    color: #667eea;
    text-decoration: none;
    margin-top: 12px;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }
}

@media screen and (max-width: 768px) {
  .developer-detail-page {
    .detail-hero .hero-content {
      flex-direction: column;
      align-items: center;
      text-align: center;
    }

    .hero-top {
      justify-content: center !important;
    }

    .hero-actions {
      justify-content: center;
    }

    .detail-layout {
      flex-direction: column;
    }

    .detail-sidebar {
      width: 100%;
    }

    .portfolio-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>

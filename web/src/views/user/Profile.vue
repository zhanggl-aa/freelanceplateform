<template>
  <div class="profile-page page-container">
    <div class="page-header">
      <h2>个人中心</h2>
      <p>管理你的个人资料和档案信息</p>
    </div>

    <el-tabs v-model="activeTab" class="profile-tabs">
      <!-- 基本信息 -->
      <el-tab-pane label="基本信息" name="basic">
        <el-card shadow="never" class="profile-card" v-loading="basicLoading">
          <el-form ref="basicFormRef" :model="basicForm" :rules="basicRules" label-width="100px" style="max-width: 600px;">
            <el-form-item label="头像">
              <div class="avatar-upload">
                <el-avatar :size="80" :src="basicForm.avatar_url" icon="UserFilled" />
                <el-upload
                  :show-file-list="false"
                  :before-upload="beforeAvatarUpload"
                  :http-request="handleAvatarUpload"
                  accept="image/*"
                >
                  <el-button size="small" type="primary" plain class="ml-12">更换头像</el-button>
                </el-upload>
              </div>
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="basicForm.nickname" placeholder="请输入昵称" maxlength="20" show-word-limit />
            </el-form-item>
            <el-form-item label="邮箱">
              <div class="flex-between w-full">
                <span>{{ userStore.user?.email || '未绑定' }}</span>
                <el-tag :type="userStore.user?.email_verified ? 'success' : 'warning'" size="small">
                  {{ userStore.user?.email_verified ? '已验证' : '未验证' }}
                </el-tag>
              </div>
            </el-form-item>
            <el-form-item label="手机号">
              <div class="flex-between w-full">
                <span>{{ userStore.user?.phone || '未绑定' }}</span>
                <el-tag :type="userStore.user?.phone_verified ? 'success' : 'warning'" size="small">
                  {{ userStore.user?.phone_verified ? '已验证' : '未验证' }}
                </el-tag>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="basicSaving" @click="handleSaveBasic">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 开发者档案 -->
      <el-tab-pane label="开发者档案" name="developer" v-if="userStore.isDeveloper">
        <div class="dev-profile-section" v-loading="devLoading">
          <!-- Portfolio Showcase -->
          <div class="portfolio-showcase">
            <div class="showcase-header">
              <h3><el-icon><PictureFilled /></el-icon> 作品展示</h3>
              <el-button type="primary" size="small" @click="openAddPortfolio">
                <el-icon><Plus /></el-icon> 添加作品
              </el-button>
            </div>

            <div class="portfolio-grid" v-if="devForm.portfolio.length > 0">
              <div v-for="(item, idx) in devForm.portfolio" :key="item.id || idx" class="portfolio-card">
                <div class="portfolio-cover">
                  <el-image
                    v-if="item.image_url"
                    :src="item.image_url"
                    fit="cover"
                    class="portfolio-img"
                  >
                    <template #error>
                      <div class="img-placeholder">
                        <el-icon :size="32" color="#c0c4cc"><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                  <div v-else class="img-placeholder">
                    <el-icon :size="32" color="#c0c4cc"><Picture /></el-icon>
                  </div>
                  <div class="portfolio-actions-overlay">
                    <el-button circle type="primary" size="small" @click="handleEditPortfolio(idx)">
                      <el-icon><Edit /></el-icon>
                    </el-button>
                    <el-button circle type="danger" size="small" @click="handleDeletePortfolio(idx, item)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                </div>
                <div class="portfolio-body">
                  <h4>{{ item.title }}</h4>
                  <p class="line-clamp-2">{{ item.description }}</p>
                  <div class="portfolio-meta" v-if="item.url">
                    <a :href="item.url" target="_blank" class="portfolio-link">
                      <el-icon><Link /></el-icon> 查看链接
                    </a>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="portfolio-empty">
              <el-icon :size="48" color="#dcdfe6"><PictureFilled /></el-icon>
              <p>还没有作品，添加你的第一个作品吧</p>
              <el-button type="primary" plain @click="openAddPortfolio">
                <el-icon><Plus /></el-icon> 添加作品
              </el-button>
            </div>
          </div>

          <!-- Profile Form -->
          <el-card shadow="never" class="profile-card">
            <el-form ref="devFormRef" :model="devForm" label-width="100px" style="max-width: 700px;">
              <el-form-item label="职业头衔">
                <el-input v-model="devForm.title" placeholder="如：全栈开发工程师" />
              </el-form-item>
              <el-form-item label="个人简介">
                <el-input v-model="devForm.bio" type="textarea" :rows="4" placeholder="介绍你的专业能力和经验" />
              </el-form-item>
              <el-row :gutter="16">
                <el-col :span="12">
                  <el-form-item label="时薪（元）">
                    <el-input-number v-model="devForm.hourly_rate" :min="0" :max="9999" class="w-full" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="可用状态">
                    <el-select v-model="devForm.availability" class="w-full">
                      <el-option label="全职可用" value="full_time" />
                      <el-option label="兼职可用" value="part_time" />
                      <el-option label="随时可用" value="available" />
                      <el-option label="暂不可用" value="unavailable" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="经验年限">
                <el-input-number v-model="devForm.experience_years" :min="0" :max="50" class="w-full" />
              </el-form-item>
              <el-form-item label="GitHub">
                <el-input v-model="devForm.github_url" placeholder="https://github.com/yourname" />
              </el-form-item>
              <el-form-item label="LinkedIn">
                <el-input v-model="devForm.linkedin_url" placeholder="https://linkedin.com/in/yourname" />
              </el-form-item>
              <el-form-item label="个人网站">
                <el-input v-model="devForm.website_url" placeholder="https://yoursite.com" />
              </el-form-item>

              <!-- Skills -->
              <el-form-item label="技能标签">
                <div class="w-full">
                  <div class="skills-wrap mb-8">
                    <el-tag
                      v-for="skill in devForm.skills"
                      :key="skill.id || skill.name"
                      closable
                      size="large"
                      class="skill-tag"
                      effect="plain"
                      round
                      @close="handleRemoveSkill(skill)"
                    >
                      {{ skill.name }}
                    </el-tag>
                  </div>
                  <div class="flex gap-8">
                    <el-input v-model="newSkill" placeholder="输入技能名称" size="small" @keyup.enter="handleAddSkill" />
                    <el-button size="small" type="primary" @click="handleAddSkill">添加</el-button>
                  </div>
                </div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" :loading="devSaving" @click="handleSaveDevProfile">保存档案</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </div>
      </el-tab-pane>

      <!-- 客户档案 -->
      <el-tab-pane label="客户档案" name="client" v-if="userStore.isClient">
        <el-card shadow="never" class="profile-card" v-loading="clientLoading">
          <el-form ref="clientFormRef" :model="clientForm" label-width="100px" style="max-width: 600px;">
            <el-form-item label="公司名称">
              <el-input v-model="clientForm.company_name" placeholder="请输入公司名称" />
            </el-form-item>
            <el-form-item label="所属行业">
              <el-select v-model="clientForm.industry" placeholder="选择行业" class="w-full" clearable>
                <el-option v-for="ind in industries" :key="ind" :label="ind" :value="ind" />
              </el-select>
            </el-form-item>
            <el-form-item label="公司规模">
              <el-select v-model="clientForm.company_size" placeholder="选择公司规模" class="w-full" clearable>
                <el-option label="1-10人" value="1-10" />
                <el-option label="11-50人" value="11-50" />
                <el-option label="51-200人" value="51-200" />
                <el-option label="201-500人" value="201-500" />
                <el-option label="500人以上" value="500+" />
              </el-select>
            </el-form-item>
            <el-form-item label="公司网站">
              <el-input v-model="clientForm.website" placeholder="https://company.com" />
            </el-form-item>
            <el-form-item label="认证状态">
              <el-tag :type="clientForm.verified ? 'success' : 'warning'">
                {{ clientForm.verified ? '已认证' : '未认证' }}
              </el-tag>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="clientSaving" @click="handleSaveClientProfile">保存档案</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- Portfolio Dialog -->
    <el-dialog v-model="showPortfolioDialog" :title="editingPortfolioIdx >= 0 ? '编辑作品' : '添加作品'" width="560px" :close-on-click-modal="false">
      <el-form ref="portfolioFormRef" :model="portfolioForm" label-position="top">
        <el-form-item label="作品标题" required>
          <el-input v-model="portfolioForm.title" placeholder="请输入作品标题" />
        </el-form-item>
        <el-form-item label="作品描述">
          <el-input v-model="portfolioForm.description" type="textarea" :rows="4" placeholder="描述作品内容和你的职责" />
        </el-form-item>
        <el-form-item label="技术栈">
          <el-select v-model="portfolioForm.tech_stack" multiple filterable allow-create placeholder="添加技术栈标签" class="w-full">
            <el-option v-for="t in commonTechStack" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="作品链接">
          <el-input v-model="portfolioForm.url" placeholder="https://..." />
        </el-form-item>
        <el-form-item label="封面图片">
          <el-upload
            :show-file-list="false"
            :before-upload="beforeAvatarUpload"
            :http-request="(opt: any) => handlePortfolioImageUpload(opt)"
            accept="image/*"
          >
            <div class="upload-area" v-if="!portfolioForm.image_url">
              <el-icon :size="32" color="#c0c4cc"><Plus /></el-icon>
              <p>点击上传图片</p>
            </div>
            <el-image v-else :src="portfolioForm.image_url" fit="cover" style="width: 100%; max-height: 200px; border-radius: 12px;" />
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPortfolioDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSavePortfolio">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { userApi, developerApi, clientApi, fileApi } from '@/api/modules'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules, UploadRequestOptions } from 'element-plus'

const userStore = useUserStore()

const activeTab = ref('basic')
const basicLoading = ref(false)
const basicSaving = ref(false)
const devLoading = ref(false)
const devSaving = ref(false)
const clientLoading = ref(false)
const clientSaving = ref(false)
const basicFormRef = ref<FormInstance>()
const newSkill = ref('')
const showPortfolioDialog = ref(false)
const editingPortfolioIdx = ref(-1)
const portfolioFormRef = ref<FormInstance>()

const industries = ['互联网/IT', '金融', '教育', '医疗', '电商', '游戏', '制造业', '房地产', '餐饮', '物流', '其他']
const commonTechStack = ['Vue', 'React', 'Angular', 'TypeScript', 'JavaScript', 'Go', 'Python', 'Java', 'Node.js', 'Docker', 'Kubernetes', 'PostgreSQL', 'MySQL', 'Redis', 'MongoDB', 'AWS', 'GCP']

const basicForm = reactive({
  nickname: '',
  avatar_url: '',
})

const basicRules: FormRules = {
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
}

const devForm = reactive({
  title: '',
  bio: '',
  hourly_rate: 0,
  availability: 'available',
  experience_years: 0,
  github_url: '',
  linkedin_url: '',
  website_url: '',
  skills: [] as any[],
  portfolio: [] as any[],
})

const clientForm = reactive({
  company_name: '',
  industry: '',
  company_size: '',
  website: '',
  verified: false,
})

const portfolioForm = reactive({
  title: '',
  description: '',
  url: '',
  image_url: '',
  tech_stack: [] as string[],
})

function beforeAvatarUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过2MB')
    return false
  }
  return true
}

async function handleAvatarUpload(options: UploadRequestOptions) {
  try {
    const res: any = await fileApi.upload(options.file, 'avatar', userStore.user?.id)
    basicForm.avatar_url = res.data?.url || res.data
    ElMessage.success('头像上传成功')
  } catch {
    ElMessage.error('头像上传失败')
  }
}

async function handlePortfolioImageUpload(options: UploadRequestOptions) {
  try {
    const res: any = await fileApi.upload(options.file, 'portfolio')
    portfolioForm.image_url = res.data?.url || res.data
    ElMessage.success('图片上传成功')
  } catch {
    ElMessage.error('图片上传失败')
  }
}

function handleAddSkill() {
  const name = newSkill.value.trim()
  if (!name) return
  if (devForm.skills.some(s => s.name === name)) {
    ElMessage.warning('该技能已存在')
    return
  }
  devForm.skills.push({ name })
  newSkill.value = ''
}

function handleRemoveSkill(skill: any) {
  const idx = devForm.skills.findIndex(s => s.name === skill.name || s.id === skill.id)
  if (idx >= 0) devForm.skills.splice(idx, 1)
}

function openAddPortfolio() {
  editingPortfolioIdx.value = -1
  portfolioForm.title = ''
  portfolioForm.description = ''
  portfolioForm.url = ''
  portfolioForm.image_url = ''
  portfolioForm.tech_stack = []
  showPortfolioDialog.value = true
}

function handleEditPortfolio(idx: number) {
  editingPortfolioIdx.value = idx
  const item = devForm.portfolio[idx]
  portfolioForm.title = item.title || ''
  portfolioForm.description = item.description || ''
  portfolioForm.url = item.url || ''
  portfolioForm.image_url = item.image_url || ''
  portfolioForm.tech_stack = item.tech_stack || []
  showPortfolioDialog.value = true
}

async function handleDeletePortfolio(idx: number, item: any) {
  try {
    await ElMessageBox.confirm('确定删除该作品？', '确认删除')
    if (item.id) {
      await developerApi.deletePortfolio(item.id)
    }
    devForm.portfolio.splice(idx, 1)
    ElMessage.success('已删除')
  } catch {}
}

function handleSavePortfolio() {
  if (!portfolioForm.title.trim()) {
    ElMessage.warning('请输入作品标题')
    return
  }
  const data = { ...portfolioForm }
  if (editingPortfolioIdx.value >= 0) {
    devForm.portfolio[editingPortfolioIdx.value] = { ...devForm.portfolio[editingPortfolioIdx.value], ...data }
  } else {
    devForm.portfolio.push(data)
  }
  showPortfolioDialog.value = false
  editingPortfolioIdx.value = -1
  portfolioForm.title = ''
  portfolioForm.description = ''
  portfolioForm.url = ''
  portfolioForm.image_url = ''
  portfolioForm.tech_stack = []
}

async function handleSaveBasic() {
  const valid = await basicFormRef.value?.validate().catch(() => false)
  if (!valid) return

  basicSaving.value = true
  try {
    await userApi.updateProfile({ nickname: basicForm.nickname, avatar_url: basicForm.avatar_url })
    await userStore.fetchUser()
    ElMessage.success('基本信息已保存')
  } catch {
    ElMessage.error('保存失败')
  } finally {
    basicSaving.value = false
  }
}

async function handleSaveDevProfile() {
  devSaving.value = true
  try {
    const data = { ...devForm }
    await developerApi.updateProfile(data)
    ElMessage.success('开发者档案已保存')
  } catch {
    ElMessage.error('保存失败')
  } finally {
    devSaving.value = false
  }
}

async function handleSaveClientProfile() {
  clientSaving.value = true
  try {
    await clientApi.updateProfile(clientForm)
    ElMessage.success('客户档案已保存')
  } catch {
    ElMessage.error('保存失败')
  } finally {
    clientSaving.value = false
  }
}

async function fetchDevProfile() {
  devLoading.value = true
  try {
    const res: any = await developerApi.getProfile()
    const data = res.data
    if (data) {
      devForm.title = data.title || ''
      devForm.bio = data.bio || ''
      devForm.hourly_rate = data.hourly_rate || 0
      devForm.availability = data.availability || 'available'
      devForm.experience_years = data.experience_years || 0
      devForm.github_url = data.github_url || ''
      devForm.linkedin_url = data.linkedin_url || ''
      devForm.website_url = data.website_url || ''
      devForm.skills = data.skills || []
      devForm.portfolio = data.portfolio || []
    }
  } catch {} finally {
    devLoading.value = false
  }
}

async function fetchClientProfile() {
  clientLoading.value = true
  try {
    const res: any = await clientApi.getProfile()
    const data = res.data
    if (data) {
      clientForm.company_name = data.company_name || ''
      clientForm.industry = data.industry || ''
      clientForm.company_size = data.company_size || ''
      clientForm.website = data.website || ''
      clientForm.verified = data.verified || false
    }
  } catch {} finally {
    clientLoading.value = false
  }
}

onMounted(() => {
  if (userStore.user) {
    basicForm.nickname = userStore.user.nickname || ''
    basicForm.avatar_url = userStore.user.avatar_url || ''
  }
  if (userStore.isDeveloper) fetchDevProfile()
  if (userStore.isClient) fetchClientProfile()
})
</script>

<style scoped lang="scss">
.profile-page {
  .profile-tabs {
    :deep(.el-tabs__item) {
      font-size: 15px;
      font-weight: 500;
    }
  }

  .profile-card {
    border-radius: 14px;
    border: 1px solid #f0f0f5;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  }

  .avatar-upload {
    display: flex;
    align-items: center;
  }

  .dev-profile-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  // Portfolio showcase
  .portfolio-showcase {
    background: #fff;
    border-radius: 14px;
    padding: 24px;
    border: 1px solid #f0f0f5;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  }

  .showcase-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;

    h3 {
      font-size: 17px;
      font-weight: 700;
      color: #1a1a2e;
      margin: 0;
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }

  .portfolio-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 16px;
  }

  .portfolio-card {
    border: 1px solid #f0f0f5;
    border-radius: 12px;
    overflow: hidden;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
      border-color: #667eea;
    }

    .portfolio-cover {
      position: relative;
      height: 160px;
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

      .portfolio-actions-overlay {
        position: absolute;
        top: 8px;
        right: 8px;
        display: flex;
        gap: 6px;
        opacity: 0;
        transition: opacity 0.2s;
      }

      &:hover .portfolio-actions-overlay {
        opacity: 1;
      }
    }

    .portfolio-body {
      padding: 14px;

      h4 {
        font-size: 15px;
        font-weight: 600;
        color: #1a1a2e;
        margin: 0 0 6px;
      }

      p {
        font-size: 13px;
        color: #8c8c9a;
        line-height: 1.5;
        margin: 0 0 8px;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .portfolio-meta {
        .portfolio-link {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          font-size: 12px;
          color: #667eea;
          text-decoration: none;
        }
      }
    }
  }

  .portfolio-empty {
    text-align: center;
    padding: 48px 20px;
    color: #c0c4cc;

    p {
      margin: 12px 0;
      font-size: 14px;
      color: #8c8c9a;
    }
  }

  .skills-wrap {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    .skill-tag {
      font-size: 14px;
    }
  }

  .upload-area {
    width: 100%;
    height: 120px;
    border: 2px dashed #dcdfe6;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.3s;

    &:hover {
      border-color: #667eea;
    }

    p {
      margin-top: 8px;
      font-size: 13px;
      color: #8c8c9a;
    }
  }
}
</style>

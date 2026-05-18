<template>
  <div class="project-edit-page page-container" v-loading="loading">
    <div class="page-header">
      <h2>编辑项目</h2>
      <p>修改项目信息（仅草稿状态可编辑）</p>
    </div>

    <el-alert v-if="project && project.status !== 'draft'" type="warning" :closable="false" class="mb-20" show-icon>
      当前项目状态为「{{ statusLabel(project?.status) }}」，仅草稿状态可以编辑
    </el-alert>

    <el-card shadow="never" class="form-card" v-if="project">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" size="large" :disabled="project.status !== 'draft'">
        <el-form-item label="项目标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入项目标题" maxlength="100" show-word-limit />
        </el-form-item>

        <el-form-item label="项目分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="选择项目分类" class="w-full" filterable>
            <el-option v-for="cat in flatCategories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="项目描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="8" placeholder="详细描述项目需求" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="预算类型" prop="budget_type">
              <el-radio-group v-model="form.budget_type">
                <el-radio value="fixed">固定价</el-radio>
                <el-radio value="hourly">时薪制</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="最低预算" prop="budget_min">
              <el-input-number v-model="form.budget_min" :min="0" :max="9999999" class="w-full" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="最高预算" prop="budget_max">
              <el-input-number v-model="form.budget_max" :min="form.budget_min || 0" :max="9999999" class="w-full" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="项目截止日期" prop="deadline">
              <el-date-picker v-model="form.deadline" type="date" placeholder="选择截止日期" class="w-full" value-format="YYYY-MM-DD" :disabled-date="(d: Date) => d < new Date()" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="投标截止日期" prop="bid_deadline">
              <el-date-picker v-model="form.bid_deadline" type="date" placeholder="选择投标截止日期" class="w-full" value-format="YYYY-MM-DD" :disabled-date="(d: Date) => d < new Date()" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="技术栈" prop="tech_stack">
          <el-select v-model="form.tech_stack" multiple filterable allow-create default-first-option placeholder="选择或输入技术标签" class="w-full">
            <el-option v-for="tech in techOptions" :key="tech" :label="tech" :value="tech" />
          </el-select>
        </el-form-item>

        <el-form-item label="附件上传">
          <el-upload
            :action="uploadUrl"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-remove="handleRemoveAttachment"
            :file-list="fileList"
            multiple
          >
            <el-button type="primary" plain><el-icon><Upload /></el-icon> 上传文件</el-button>
            <template #tip>
              <div class="upload-tip">支持上传文档、图片、压缩包，单个文件不超过20MB</div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item class="form-actions" v-if="project.status === 'draft'">
          <el-button type="primary" size="large" @click="handleSave" :loading="submitting">保存修改</el-button>
          <el-button type="success" size="large" @click="handlePublish" :loading="submitting">发布项目</el-button>
          <el-button size="large" @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { projectApi, categoryApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const projectId = computed(() => route.params.id as string)

const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const project = ref<any>(null)
const categories = ref<any[]>([])
const fileList = ref<any[]>([])
const attachmentIds = ref<string[]>([])

const uploadUrl = '/api/v1/files/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.accessToken}`,
}))

const techOptions = ['Vue.js', 'React', 'Angular', 'Node.js', 'Python', 'Java', 'Go', 'PHP', 'TypeScript', 'Flutter', 'Swift', 'Kotlin', 'C++', 'Rust', 'Docker', 'AWS', 'AI/ML', 'MySQL', 'PostgreSQL', 'MongoDB', 'Redis']

const form = reactive({
  title: '',
  category_id: '',
  description: '',
  budget_type: 'fixed',
  budget_min: 0,
  budget_max: 0,
  deadline: '',
  bid_deadline: '',
  tech_stack: [] as string[],
})

const validateBudget = (_rule: any, _value: any, callback: Function) => {
  if (form.budget_max > 0 && form.budget_max < form.budget_min) {
    callback(new Error('最高预算不能低于最低预算'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  title: [{ required: true, message: '请输入项目标题', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择项目分类', trigger: 'change' }],
  description: [{ required: true, message: '请输入项目描述', trigger: 'blur' }],
  budget_type: [{ required: true, message: '请选择预算类型', trigger: 'change' }],
  budget_min: [{ validator: validateBudget, trigger: 'blur' }],
  budget_max: [{ validator: validateBudget, trigger: 'blur' }],
  deadline: [{ required: true, message: '请选择截止日期', trigger: 'change' }],
}

const flatCategories = computed(() => {
  const result: any[] = []
  function flatten(items: any[]) {
    for (const item of items) {
      result.push(item)
      if (item.children) flatten(item.children)
    }
  }
  flatten(categories.value)
  return result
})

function statusLabel(status: string) {
  const map: Record<string, string> = { published: '招标中', draft: '草稿', ongoing: '进行中', completed: '已完成', cancelled: '已取消' }
  return map[status] || status
}

function handleUploadSuccess(response: any) {
  const fileId = response.data?.id || response.data?.file_id
  if (fileId) attachmentIds.value.push(fileId)
}

function handleRemoveAttachment(file: any) {
  const fileId = file.response?.data?.id || file.response?.data?.file_id
  if (fileId) {
    attachmentIds.value = attachmentIds.value.filter(id => id !== fileId)
  }
}

function populateForm(data: any) {
  form.title = data.title || ''
  form.category_id = data.category_id || ''
  form.description = data.description || ''
  form.budget_type = data.budget_type || 'fixed'
  form.budget_min = data.budget_min || 0
  form.budget_max = data.budget_max || 0
  form.deadline = data.deadline ? data.deadline.substring(0, 10) : ''
  form.bid_deadline = data.bid_deadline ? data.bid_deadline.substring(0, 10) : ''
  form.tech_stack = data.tech_stack || []
  if (data.attachments) {
    fileList.value = data.attachments.map((a: any) => ({
      name: a.name || a.filename,
      url: a.url,
      response: { data: a },
    }))
  }
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    await projectApi.update(projectId.value, { ...form, attachment_ids: attachmentIds.value.length > 0 ? attachmentIds.value : undefined })
    ElMessage.success('保存成功')
    fetchProject()
  } catch (err: any) {
    ElMessage.error(err.message || '保存失败')
  } finally {
    submitting.value = false
  }
}

async function handlePublish() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    await projectApi.update(projectId.value, { ...form, attachment_ids: attachmentIds.value.length > 0 ? attachmentIds.value : undefined })
    await projectApi.publish(projectId.value)
    ElMessage.success('项目已发布')
    router.push(`/projects/${projectId.value}`)
  } catch (err: any) {
    ElMessage.error(err.message || '发布失败')
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.back()
}

async function fetchProject() {
  loading.value = true
  try {
    const res: any = await projectApi.getById(projectId.value)
    project.value = res.data
    populateForm(res.data)
  } catch {} finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const res: any = await categoryApi.getTree()
    categories.value = res.data || []
  } catch {}
}

onMounted(() => {
  fetchProject()
  fetchCategories()
})
</script>

<style scoped lang="scss">
.project-edit-page {
  max-width: 900px;
  margin: 0 auto;

  .form-card {
    border-radius: 12px;
    padding: 12px;
  }

  .upload-tip {
    font-size: 12px;
    color: var(--color-text-secondary);
    margin-top: 4px;
  }

  .form-actions {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid var(--border-color-light);

    .el-button {
      min-width: 120px;
    }
  }
}
</style>

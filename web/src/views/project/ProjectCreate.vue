<template>
  <div class="project-create-page page-container">
    <div class="page-header">
      <h2>发布项目</h2>
      <p>填写项目详情，吸引优秀开发者投标</p>
    </div>

    <el-card shadow="never" class="form-card">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" size="large" v-loading="submitting">
        <el-form-item label="项目标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入项目标题，简洁明了" maxlength="100" show-word-limit />
        </el-form-item>

        <el-form-item label="项目分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="选择项目分类" class="w-full" filterable>
            <el-option v-for="cat in flatCategories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>

        <el-form-item label="项目描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="8" placeholder="详细描述项目需求、功能要求、技术要求等" />
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
              <el-input-number v-model="form.budget_min" :min="0" :max="9999999" class="w-full" placeholder="最低" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="最高预算" prop="budget_max">
              <el-input-number v-model="form.budget_max" :min="form.budget_min || 0" :max="9999999" class="w-full" placeholder="最高" />
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
            :on-error="handleUploadError"
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

        <el-form-item class="form-actions">
          <el-button type="primary" size="large" @click="handleSubmit('published')" :loading="submitting">发布项目</el-button>
          <el-button size="large" @click="handleSubmit('draft')" :loading="submitting">保存草稿</el-button>
          <el-button size="large" @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { projectApi, categoryApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const submitting = ref(false)
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
  if (form.budget_min < 0 || form.budget_max < 0) {
    callback(new Error('预算不能为负数'))
  } else if (form.budget_max > 0 && form.budget_max < form.budget_min) {
    callback(new Error('最高预算不能低于最低预算'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  title: [
    { required: true, message: '请输入项目标题', trigger: 'blur' },
    { min: 5, max: 100, message: '标题长度在5到100个字符之间', trigger: 'blur' },
  ],
  category_id: [{ required: true, message: '请选择项目分类', trigger: 'change' }],
  description: [
    { required: true, message: '请输入项目描述', trigger: 'blur' },
    { min: 20, message: '项目描述不少于20个字符', trigger: 'blur' },
  ],
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

function handleUploadSuccess(response: any, file: any) {
  const fileId = response.data?.id || response.data?.file_id
  if (fileId) attachmentIds.value.push(fileId)
  ElMessage.success('文件上传成功')
}

function handleUploadError() {
  ElMessage.error('文件上传失败')
}

function handleRemoveAttachment(file: any) {
  const fileId = file.response?.data?.id || file.response?.data?.file_id
  if (fileId) {
    attachmentIds.value = attachmentIds.value.filter(id => id !== fileId)
  }
}

async function handleSubmit(status: string) {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data: any = {
      ...form,
      status,
      attachment_ids: attachmentIds.value.length > 0 ? attachmentIds.value : undefined,
    }
    const res: any = await projectApi.create(data)
    ElMessage.success(status === 'published' ? '项目发布成功' : '草稿保存成功')
    router.push(`/projects/${res.data?.id}`)
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.back()
}

async function fetchCategories() {
  try {
    const res: any = await categoryApi.getTree()
    categories.value = res.data || []
  } catch {}
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped lang="scss">
.project-create-page {
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

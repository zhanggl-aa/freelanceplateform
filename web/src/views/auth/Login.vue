<template>
  <div class="login-page">
    <div class="login-container">
      <el-card class="login-card" shadow="always">
        <div class="login-header">
          <el-icon :size="36" color="#409EFF"><Monitor /></el-icon>
          <h2>登录接单平台</h2>
          <p>欢迎回来，请登录您的账号</p>
        </div>

        <el-tabs v-model="loginType" class="login-tabs" stretch>
          <el-tab-pane label="邮箱登录" name="email">
            <el-form ref="emailFormRef" :model="emailForm" :rules="emailRules" label-position="top" size="large" @submit.prevent="handleEmailLogin">
              <el-form-item label="邮箱地址" prop="email">
                <el-input v-model="emailForm.email" placeholder="请输入邮箱地址" prefix-icon="Message" />
              </el-form-item>
              <el-form-item label="密码" prop="password">
                <el-input v-model="emailForm.password" type="password" placeholder="请输入密码" prefix-icon="Lock" show-password @keyup.enter="handleEmailLogin" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" class="login-btn" :loading="loading" @click="handleEmailLogin">登录</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <el-tab-pane label="手机登录" name="phone">
            <el-form ref="phoneFormRef" :model="phoneForm" :rules="phoneRules" label-position="top" size="large" @submit.prevent="handlePhoneLogin">
              <el-form-item label="手机号码" prop="phone">
                <el-input v-model="phoneForm.phone" placeholder="请输入手机号码" prefix-icon="Iphone">
                  <template #prepend>+86</template>
                </el-input>
              </el-form-item>
              <el-form-item label="密码" prop="password">
                <el-input v-model="phoneForm.password" type="password" placeholder="请输入密码" prefix-icon="Lock" show-password @keyup.enter="handlePhoneLogin" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" class="login-btn" :loading="loading" @click="handlePhoneLogin">登录</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>

        <div class="login-footer">
          <router-link to="/forgot-password" class="forgot-link">忘记密码？</router-link>
          <span class="divider">|</span>
          <router-link to="/register" class="register-link">还没有账号？立即注册</router-link>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginType = ref('email')
const loading = ref(false)

const emailFormRef = ref<FormInstance>()
const phoneFormRef = ref<FormInstance>()

const emailForm = reactive({
  email: '',
  password: '',
})

const phoneForm = reactive({
  phone: '',
  password: '',
})

const emailRules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不少于6位', trigger: 'blur' },
  ],
}

const phoneRules: FormRules = {
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不少于6位', trigger: 'blur' },
  ],
}

async function handleEmailLogin() {
  const valid = await emailFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login(emailForm.email, emailForm.password)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    ElMessage.error(err.message || '登录失败，请检查账号密码')
  } finally {
    loading.value = false
  }
}

async function handlePhoneLogin() {
  const valid = await phoneFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login(phoneForm.phone, phoneForm.password)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    ElMessage.error(err.message || '登录失败，请检查账号密码')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #409EFF 50%, #764ba2 100%);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 440px;
  border-radius: 16px;
  border: none;

  :deep(.el-card__body) {
    padding: 40px 36px;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 28px;

  h2 {
    font-size: 22px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin: 12px 0 4px;
  }

  p {
    font-size: 14px;
    color: var(--color-text-secondary);
  }
}

.login-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 24px;
  }
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
}

.login-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;

  .forgot-link,
  .register-link {
    color: var(--color-primary);
    text-decoration: none;
    &:hover {
      text-decoration: underline;
    }
  }

  .divider {
    color: var(--color-text-placeholder);
    margin: 0 12px;
  }
}
</style>

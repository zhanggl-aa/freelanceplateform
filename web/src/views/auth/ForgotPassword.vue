<template>
  <div class="forgot-page">
    <div class="forgot-container">
      <el-card class="forgot-card" shadow="always">
        <div class="forgot-header">
          <el-icon :size="36" color="#409EFF"><Lock /></el-icon>
          <h2>找回密码</h2>
          <p>请按照步骤重置您的密码</p>
        </div>

        <el-steps :active="step" align-center class="forgot-steps" finish-status="success">
          <el-step title="验证账号" />
          <el-step title="输入验证码" />
          <el-step title="重置密码" />
        </el-steps>

        <!-- Step 1: Enter account -->
        <div v-if="step === 0" class="step-content">
          <el-form ref="step1FormRef" :model="step1Form" :rules="step1Rules" label-position="top" size="large">
            <el-form-item label="邮箱或手机号" prop="account">
              <el-input v-model="step1Form.account" placeholder="请输入注册时使用的邮箱或手机号" prefix-icon="User" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="step-btn" :loading="loading" @click="handleSendCode">发送验证码</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- Step 2: Enter verification code -->
        <div v-if="step === 1" class="step-content">
          <el-form ref="step2FormRef" :model="step2Form" :rules="step2Rules" label-position="top" size="large">
            <el-alert type="info" :closable="false" class="code-tip" show-icon>
              验证码已发送至 {{ maskedAccount }}，请查收
            </el-alert>
            <el-form-item label="验证码" prop="code">
              <div class="code-input-wrap">
                <el-input v-model="step2Form.code" placeholder="请输入6位验证码" maxlength="6" prefix-icon="Key" />
                <el-button :disabled="countdown > 0" @click="handleResendCode" class="resend-btn">
                  {{ countdown > 0 ? `${countdown}s后重发` : '重新发送' }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="step-btn" :loading="loading" @click="handleVerifyCode">验证</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- Step 3: Set new password -->
        <div v-if="step === 2" class="step-content">
          <el-form ref="step3FormRef" :model="step3Form" :rules="step3Rules" label-position="top" size="large">
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="step3Form.newPassword" type="password" placeholder="请输入新密码（至少6位）" prefix-icon="Lock" show-password />
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirmPassword">
              <el-input v-model="step3Form.confirmPassword" type="password" placeholder="请再次输入新密码" prefix-icon="Lock" show-password @keyup.enter="handleResetPassword" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="step-btn" :loading="loading" @click="handleResetPassword">重置密码</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- Success -->
        <div v-if="step === 3" class="step-content step-success">
          <el-icon :size="64" color="#67C23A"><CircleCheckFilled /></el-icon>
          <h3>密码重置成功！</h3>
          <p>您现在可以使用新密码登录了</p>
          <el-button type="primary" size="large" @click="goLogin" class="step-btn">去登录</el-button>
        </div>

        <div class="forgot-footer" v-if="step < 3">
          <el-button text @click="goLogin">
            <el-icon><ArrowLeft /></el-icon>返回登录
          </el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/modules'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const step = ref(0)
const loading = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const step1FormRef = ref<FormInstance>()
const step2FormRef = ref<FormInstance>()
const step3FormRef = ref<FormInstance>()

const step1Form = reactive({ account: '' })
const step2Form = reactive({ code: '' })
const step3Form = reactive({ newPassword: '', confirmPassword: '' })

const maskedAccount = computed(() => {
  const acc = step1Form.account
  if (acc.includes('@')) {
    const [local, domain] = acc.split('@')
    return local.substring(0, 2) + '***@' + domain
  }
  if (/^1\d{10}$/.test(acc)) {
    return acc.substring(0, 3) + '****' + acc.substring(7)
  }
  return acc
})

const step1Rules: FormRules = {
  account: [
    { required: true, message: '请输入邮箱或手机号', trigger: 'blur' },
  ],
}

const step2Rules: FormRules = {
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '验证码只能为数字', trigger: 'blur' },
  ],
}

const validateConfirmPassword = (_rule: any, value: string, callback: Function) => {
  if (value !== step3Form.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const step3Rules: FormRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 32, message: '密码长度在6到32个字符之间', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

function startCountdown() {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      if (timer) clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function handleSendCode() {
  const valid = await step1FormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authApi.forgotPassword(step1Form.account)
    ElMessage.success('验证码已发送')
    step.value = 1
    startCountdown()
  } catch (err: any) {
    ElMessage.error(err.message || '发送验证码失败')
  } finally {
    loading.value = false
  }
}

async function handleResendCode() {
  loading.value = true
  try {
    await authApi.forgotPassword(step1Form.account)
    ElMessage.success('验证码已重新发送')
    startCountdown()
  } catch (err: any) {
    ElMessage.error(err.message || '发送验证码失败')
  } finally {
    loading.value = false
  }
}

async function handleVerifyCode() {
  const valid = await step2FormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    step.value = 2
  } finally {
    loading.value = false
  }
}

async function handleResetPassword() {
  const valid = await step3FormRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authApi.resetPassword({
      account: step1Form.account,
      code: step2Form.code,
      new_password: step3Form.newPassword,
    })
    step.value = 3
    ElMessage.success('密码重置成功')
  } catch (err: any) {
    ElMessage.error(err.message || '密码重置失败')
  } finally {
    loading.value = false
  }
}

function goLogin() {
  router.push('/login')
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped lang="scss">
.forgot-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #409EFF 50%, #764ba2 100%);
  padding: 20px;
}

.forgot-card {
  width: 100%;
  max-width: 480px;
  border-radius: 16px;
  border: none;

  :deep(.el-card__body) {
    padding: 36px;
  }
}

.forgot-header {
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

.forgot-steps {
  margin-bottom: 32px;
}

.step-content {
  min-height: 180px;
}

.code-tip {
  margin-bottom: 20px;
}

.code-input-wrap {
  display: flex;
  gap: 12px;
  width: 100%;

  .el-input {
    flex: 1;
  }

  .resend-btn {
    white-space: nowrap;
  }
}

.step-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
}

.step-success {
  text-align: center;
  padding: 20px 0;

  h3 {
    font-size: 20px;
    color: var(--color-text-primary);
    margin: 16px 0 8px;
  }

  p {
    color: var(--color-text-secondary);
    margin-bottom: 24px;
  }
}

.forgot-footer {
  text-align: center;
  margin-top: 16px;
}
</style>

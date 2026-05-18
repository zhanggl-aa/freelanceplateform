<template>
  <div class="register-page">
    <!-- Animated background -->
    <div class="bg-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>

    <div class="register-container">
      <!-- Step indicators -->
      <div class="steps-indicator">
        <div class="step-dot" :class="{ active: step >= 1, done: step > 1 }">
          <span class="step-num">1</span>
          <span class="step-label">账号信息</span>
        </div>
        <div class="step-line" :class="{ active: step > 1 }"></div>
        <div class="step-dot" :class="{ active: step >= 2, done: step > 2 }">
          <span class="step-num">2</span>
          <span class="step-label">选择角色</span>
        </div>
        <div class="step-line" :class="{ active: step > 2 }"></div>
        <div class="step-dot" :class="{ active: step >= 3 }">
          <span class="step-num">3</span>
          <span class="step-label">完成注册</span>
        </div>
      </div>

      <transition name="slide-fade" mode="out-in">
        <!-- Step 1: Account Info -->
        <el-card v-if="step === 1" key="step1" class="register-card" shadow="always">
          <div class="register-header">
            <div class="logo-circle">
              <el-icon :size="32" color="#fff"><Monitor /></el-icon>
            </div>
            <h2>创建账号</h2>
            <p>填写基本信息，开启你的自由职业之旅</p>
          </div>

          <el-form ref="formRef" :model="form" :rules="rules" label-position="top" size="large" @submit.prevent>
            <el-form-item label="注册方式" prop="accountType">
              <el-segmented v-model="form.accountType" :options="accountTypeOptions" block />
            </el-form-item>

            <el-form-item v-if="form.accountType === 'email'" label="邮箱地址" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱地址" prefix-icon="Message" />
            </el-form-item>

            <el-form-item v-if="form.accountType === 'phone'" label="手机号码" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号码" prefix-icon="Iphone">
                <template #prepend>+86</template>
              </el-input>
            </el-form-item>

            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="form.nickname" placeholder="你的专属昵称" prefix-icon="User" maxlength="20" show-word-limit />
            </el-form-item>

            <el-form-item label="密码" prop="password">
              <el-input v-model="form.password" type="password" placeholder="请输入密码（至少6位）" prefix-icon="Lock" show-password />
              <div class="password-strength" v-if="form.password">
                <div class="strength-bar">
                  <div class="strength-fill" :style="{ width: passwordStrength.percent + '%' }" :class="passwordStrength.level"></div>
                </div>
                <span class="strength-text" :class="passwordStrength.level">{{ passwordStrength.text }}</span>
              </div>
            </el-form-item>

            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="form.confirmPassword" type="password" placeholder="请再次输入密码" prefix-icon="Lock" show-password />
            </el-form-item>

            <el-button type="primary" class="next-btn" @click="goToStep2">
              下一步<el-icon class="ml-8"><ArrowRight /></el-icon>
            </el-button>
          </el-form>

          <div class="register-footer">
            已有账号？<router-link to="/login" class="login-link">立即登录</router-link>
          </div>
        </el-card>

        <!-- Step 2: Role Selection -->
        <el-card v-else-if="step === 2" key="step2" class="register-card role-card" shadow="always">
          <div class="register-header">
            <div class="logo-circle">
              <el-icon :size="32" color="#fff"><UserFilled /></el-icon>
            </div>
            <h2>选择你的角色</h2>
            <p>选择最适合你的身份，随时可以切换</p>
          </div>

          <div class="role-cards">
            <div
              class="role-card-item"
              :class="{ selected: form.user_type === 'client' }"
              @click="form.user_type = 'client'"
            >
              <div class="role-icon-wrap client-icon">
                <el-icon :size="36"><Briefcase /></el-icon>
              </div>
              <h3>我要发包</h3>
              <p class="role-desc">作为甲方发布项目需求，寻找优秀开发者合作</p>
              <ul class="role-features">
                <li><el-icon><Check /></el-icon> 发布项目需求</li>
                <li><el-icon><Check /></el-icon> 筛选开发者</li>
                <li><el-icon><Check /></el-icon> 托管支付</li>
              </ul>
              <div class="role-check" v-if="form.user_type === 'client'">
                <el-icon :size="20"><CircleCheckFilled /></el-icon>
              </div>
            </div>

            <div
              class="role-card-item"
              :class="{ selected: form.user_type === 'developer' }"
              @click="form.user_type = 'developer'"
            >
              <div class="role-icon-wrap dev-icon">
                <el-icon :size="36"><EditPen /></el-icon>
              </div>
              <h3>我要接单</h3>
              <p class="role-desc">作为开发者展示技能，承接项目赚取收入</p>
              <ul class="role-features">
                <li><el-icon><Check /></el-icon> 展示技能作品</li>
                <li><el-icon><Check /></el-icon> 投标接单</li>
                <li><el-icon><Check /></el-icon> 安全收款</li>
              </ul>
              <div class="role-check" v-if="form.user_type === 'developer'">
                <el-icon :size="20"><CircleCheckFilled /></el-icon>
              </div>
            </div>
          </div>

          <div class="step2-actions">
            <el-button size="large" @click="step = 1">
              <el-icon><ArrowLeft /></el-icon>上一步
            </el-button>
            <el-button type="primary" size="large" class="next-btn" @click="goToStep3">
              下一步<el-icon class="ml-8"><ArrowRight /></el-icon>
            </el-button>
          </div>
        </el-card>

        <!-- Step 3: Confirm -->
        <el-card v-else-if="step === 3" key="step3" class="register-card" shadow="always">
          <div class="register-header">
            <div class="logo-circle success-icon">
              <el-icon :size="32" color="#fff"><CircleCheck /></el-icon>
            </div>
            <h2>确认注册</h2>
            <p>请确认以下信息无误</p>
          </div>

          <div class="confirm-info">
            <div class="info-row">
              <span class="info-label">账号</span>
              <span class="info-value">{{ form.accountType === 'email' ? form.email : '+86 ' + form.phone }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">昵称</span>
              <span class="info-value">{{ form.nickname }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">角色</span>
              <el-tag :type="form.user_type === 'developer' ? 'success' : 'primary'" size="large" effect="dark">
                {{ form.user_type === 'developer' ? '开发者' : '甲方' }}
              </el-tag>
            </div>
          </div>

          <el-form-item prop="agree" class="agree-item">
            <el-checkbox v-model="form.agree">
              我已阅读并同意
              <a href="javascript:void(0)" class="terms-link">《用户协议》</a>和<a href="javascript:void(0)" class="terms-link">《隐私政策》</a>
            </el-checkbox>
          </el-form-item>

          <div class="step2-actions">
            <el-button size="large" @click="step = 2">
              <el-icon><ArrowLeft /></el-icon>上一步
            </el-button>
            <el-button type="primary" size="large" class="register-btn" :loading="loading" @click="handleRegister">
              完成注册
            </el-button>
          </div>
        </el-card>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const step = ref(1)

const accountTypeOptions = [
  { label: '邮箱注册', value: 'email' },
  { label: '手机注册', value: 'phone' },
]

const form = reactive({
  accountType: 'email',
  email: '',
  phone: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  user_type: 'developer',
  agree: false,
})

const passwordStrength = computed(() => {
  const pwd = form.password
  if (!pwd) return { percent: 0, level: '', text: '' }
  let score = 0
  if (pwd.length >= 6) score += 20
  if (pwd.length >= 10) score += 20
  if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score += 20
  if (/\d/.test(pwd)) score += 20
  if (/[^a-zA-Z0-9]/.test(pwd)) score += 20
  if (score <= 40) return { percent: score, level: 'weak', text: '弱' }
  if (score <= 60) return { percent: score, level: 'medium', text: '中' }
  return { percent: score, level: 'strong', text: '强' }
})

const validateConfirmPassword = (_rule: any, value: string, callback: Function) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const validateAgree = (_rule: any, value: boolean, callback: Function) => {
  if (!value) {
    callback(new Error('请同意用户协议和隐私政策'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' },
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在2到20个字符之间', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 32, message: '密码长度在6到32个字符之间', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
  user_type: [
    { required: true, message: '请选择角色类型', trigger: 'change' },
  ],
  agree: [
    { validator: validateAgree, trigger: 'change' },
  ],
}

async function goToStep2() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  step.value = 2
}

function goToStep3() {
  step.value = 3
}

async function handleRegister() {
  if (!form.agree) {
    ElMessage.warning('请同意用户协议和隐私政策')
    return
  }

  loading.value = true
  try {
    const data: any = {
      password: form.password,
      nickname: form.nickname,
      user_type: form.user_type,
    }
    if (form.accountType === 'email') {
      data.email = form.email
    } else {
      data.phone = form.phone
    }
    await userStore.register(data)
    ElMessage.success('注册成功')
    router.push('/')
  } catch (err: any) {
    ElMessage.error(err.message || '注册失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.bg-shapes {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: float 20s infinite ease-in-out;
}

.shape-1 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  top: -200px;
  right: -100px;
  animation-delay: 0s;
}

.shape-2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #43e97b, #38f9d7);
  bottom: -150px;
  left: -100px;
  animation-delay: -7s;
}

.shape-3 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, #fa709a, #fee140);
  top: 50%;
  left: 50%;
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.05); }
  66% { transform: translate(-20px, 20px) scale(0.95); }
}

.register-container {
  width: 100%;
  max-width: 520px;
  position: relative;
  z-index: 1;
}

.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 28px;
  gap: 0;
}

.step-dot {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  opacity: 0.4;
  transition: all 0.3s;

  &.active {
    opacity: 1;
  }

  .step-num {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.15);
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    transition: all 0.3s;
  }

  &.active .step-num {
    background: linear-gradient(135deg, #667eea, #764ba2);
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  }

  &.done .step-num {
    background: #67C23A;
  }

  .step-label {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.7);
    white-space: nowrap;
  }

  &.active .step-label {
    color: #fff;
  }
}

.step-line {
  width: 60px;
  height: 2px;
  background: rgba(255, 255, 255, 0.15);
  margin: 0 8px;
  margin-bottom: 20px;
  transition: background 0.3s;

  &.active {
    background: linear-gradient(90deg, #67C23A, #667eea);
  }
}

.register-card {
  width: 100%;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  background: rgba(255, 255, 255, 0.95);

  :deep(.el-card__body) {
    padding: 40px;
  }
}

.register-header {
  text-align: center;
  margin-bottom: 28px;

  .logo-circle {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea, #764ba2);
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
  }

  .success-icon {
    background: linear-gradient(135deg, #67C23A, #43e97b);
    box-shadow: 0 8px 25px rgba(103, 194, 58, 0.3);
  }

  h2 {
    font-size: 24px;
    font-weight: 700;
    color: #1a1a2e;
    margin: 0 0 6px;
  }

  p {
    font-size: 14px;
    color: #8c8c9a;
  }
}

.password-strength {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 8px;

  .strength-bar {
    flex: 1;
    height: 4px;
    background: #e8e8ed;
    border-radius: 2px;
    overflow: hidden;
  }

  .strength-fill {
    height: 100%;
    border-radius: 2px;
    transition: width 0.3s, background 0.3s;

    &.weak { background: #F56C6C; }
    &.medium { background: #E6A23C; }
    &.strong { background: #67C23A; }
  }

  .strength-text {
    font-size: 12px;
    font-weight: 500;
    min-width: 20px;

    &.weak { color: #F56C6C; }
    &.medium { color: #E6A23C; }
    &.strong { color: #67C23A; }
  }
}

.next-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  margin-top: 8px;

  &:hover {
    background: linear-gradient(135deg, #7c93f0, #8b5fb8);
  }
}

// Role selection cards
.role-cards {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.role-card-item {
  flex: 1;
  padding: 24px 20px;
  border: 2px solid #e8e8ed;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  text-align: center;

  &:hover {
    border-color: #c0c0cc;
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
  }

  &.selected {
    border-color: #667eea;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.05), rgba(118, 75, 162, 0.05));
    box-shadow: 0 8px 25px rgba(102, 126, 234, 0.15);
  }

  .role-icon-wrap {
    width: 72px;
    height: 72px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    transition: all 0.3s;
  }

  .client-icon {
    background: linear-gradient(135deg, #e8f4fd, #d6eaf8);
    color: #409EFF;
  }

  .dev-icon {
    background: linear-gradient(135deg, #e8fdf5, #d6f5e8);
    color: #67C23A;
  }

  &.selected .client-icon {
    background: linear-gradient(135deg, #409EFF, #337ecc);
    color: #fff;
  }

  &.selected .dev-icon {
    background: linear-gradient(135deg, #67C23A, #4ea82a);
    color: #fff;
  }

  h3 {
    font-size: 18px;
    font-weight: 700;
    color: #1a1a2e;
    margin: 0 0 8px;
  }

  .role-desc {
    font-size: 13px;
    color: #8c8c9a;
    line-height: 1.5;
    margin: 0 0 16px;
  }

  .role-features {
    list-style: none;
    padding: 0;
    margin: 0;
    text-align: left;

    li {
      font-size: 13px;
      color: #5a5a6e;
      padding: 4px 0;
      display: flex;
      align-items: center;
      gap: 6px;

      .el-icon {
        color: #67C23A;
        font-size: 14px;
      }
    }
  }

  .role-check {
    position: absolute;
    top: 12px;
    right: 12px;
    color: #667eea;
    animation: scaleIn 0.3s ease;
  }
}

@keyframes scaleIn {
  from { transform: scale(0); }
  to { transform: scale(1); }
}

.step2-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;

  .el-button {
    flex: 1;
    height: 48px;
    border-radius: 12px;
    font-size: 16px;
  }

  .next-btn {
    flex: 2;
  }
}

.confirm-info {
  background: #f7f8fa;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;

  .info-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid #e8e8ed;

    &:last-child {
      border-bottom: none;
    }

    .info-label {
      font-size: 14px;
      color: #8c8c9a;
    }

    .info-value {
      font-size: 14px;
      font-weight: 500;
      color: #1a1a2e;
    }
  }
}

.agree-item {
  margin-bottom: 0;
}

.terms-link {
  color: #667eea;
  text-decoration: none;
  &:hover {
    text-decoration: underline;
  }
}

.register-btn {
  background: linear-gradient(135deg, #67C23A, #43e97b) !important;
  border: none !important;

  &:hover {
    background: linear-gradient(135deg, #7dd355, #5cf0af) !important;
  }
}

.register-footer {
  text-align: center;
  font-size: 14px;
  color: #8c8c9a;
  margin-top: 20px;

  .login-link {
    color: #667eea;
    text-decoration: none;
    font-weight: 600;
    &:hover {
      text-decoration: underline;
    }
  }
}

.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s ease-in;
}

.slide-fade-enter-from {
  transform: translateX(20px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateX(-20px);
  opacity: 0;
}

@media screen and (max-width: 600px) {
  .role-cards {
    flex-direction: column;
  }

  .register-card :deep(.el-card__body) {
    padding: 24px;
  }
}
</style>

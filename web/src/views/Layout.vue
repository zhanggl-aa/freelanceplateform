<template>
  <div class="layout">
    <!-- Top Navbar -->
    <header class="navbar">
      <div class="navbar-inner">
        <div class="navbar-left">
          <router-link to="/" class="logo">
            <el-icon :size="24"><Monitor /></el-icon>
            <span class="logo-text">接单平台</span>
          </router-link>
          <nav class="nav-links">
            <router-link to="/" class="nav-link" :class="{ active: $route.path === '/' }">首页</router-link>
            <router-link to="/projects" class="nav-link" :class="{ active: $route.path.startsWith('/projects') }">找项目</router-link>
            <router-link to="/developers" class="nav-link" :class="{ active: $route.path.startsWith('/developers') }">找开发者</router-link>
          </nav>
        </div>

        <div class="navbar-right">
          <template v-if="userStore.isLoggedIn">
            <el-badge :value="notificationCount" :hidden="notificationCount === 0" :max="99" class="nav-icon-badge">
              <router-link to="/notifications" class="nav-icon-btn">
                <el-icon :size="20"><Bell /></el-icon>
              </router-link>
            </el-badge>

            <el-badge :value="chatStore.unreadCount" :hidden="chatStore.unreadCount === 0" :max="99" class="nav-icon-badge">
              <router-link to="/chat" class="nav-icon-btn">
                <el-icon :size="20"><ChatDotRound /></el-icon>
              </router-link>
            </el-badge>

            <el-dropdown trigger="click" @command="handleCommand">
              <div class="user-avatar-wrap">
                <el-avatar :size="32" :src="userStore.user?.avatar_url" icon="UserFilled" />
                <span class="user-name hide-on-mobile">{{ userStore.user?.nickname }}</span>
                <el-icon class="arrow"><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>个人中心
                  </el-dropdown-item>
                  <el-dropdown-item command="myProjects">
                    <el-icon><Folder /></el-icon>我的项目
                  </el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isDeveloper" command="myBids">
                    <el-icon><EditPen /></el-icon>我的投标
                  </el-dropdown-item>
                  <el-dropdown-item command="myContracts">
                    <el-icon><Document /></el-icon>我的合同
                  </el-dropdown-item>
                  <el-dropdown-item command="wallet">
                    <el-icon><Wallet /></el-icon>我的钱包
                  </el-dropdown-item>
                  <el-dropdown-item v-if="isAdmin" command="admin" divided>
                    <el-icon><Setting /></el-icon>后台管理
                  </el-dropdown-item>
                  <el-dropdown-item command="logout" divided>
                    <el-icon><SwitchButton /></el-icon>退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>

          <template v-else>
            <router-link to="/login">
              <el-button type="primary" plain size="small">登录</el-button>
            </router-link>
            <router-link to="/register" class="ml-12">
              <el-button type="primary" size="small">注册</el-button>
            </router-link>
          </template>
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <div class="main-wrapper">
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- Footer -->
    <footer class="footer">
      <div class="footer-inner">
        <p>&copy; {{ new Date().getFullYear() }} 接单平台 - 专业的软件开发者自由职业平台</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useChatStore } from '@/store/chat'
import { notificationApi } from '@/api/modules'

const router = useRouter()
const userStore = useUserStore()
const chatStore = useChatStore()

const notificationCount = ref(0)
let pollTimer: ReturnType<typeof setInterval> | null = null

const isAdmin = computed(() => userStore.user?.user_type === 'admin')

function handleCommand(command: string) {
  const routeMap: Record<string, string> = {
    profile: '/profile',
    myProjects: '/my/projects',
    myBids: '/my/bids',
    myContracts: '/my/contracts',
    wallet: '/wallet',
    admin: '/admin',
  }
  if (command === 'logout') {
    userStore.logout()
    chatStore.disconnect()
    router.push('/login')
  } else if (routeMap[command]) {
    router.push(routeMap[command])
  }
}

async function fetchNotificationCount() {
  if (!userStore.isLoggedIn) return
  try {
    const res: any = await notificationApi.unreadCount()
    notificationCount.value = res.data?.count || 0
  } catch {}
}

function startPolling() {
  fetchNotificationCount()
  pollTimer = setInterval(fetchNotificationCount, 60000)
}

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  if (userStore.isLoggedIn) {
    chatStore.fetchUnreadCount()
    chatStore.connectWebSocket(userStore.accessToken)
    startPolling()
  }
})

onUnmounted(() => {
  stopPolling()
  chatStore.disconnect()
})
</script>

<style scoped lang="scss">
.layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: var(--header-height);
  background: var(--bg-color-white);
  border-bottom: 1px solid var(--border-color-light);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  z-index: 1000;

  .navbar-inner {
    max-width: 1400px;
    margin: 0 auto;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
  }

  .navbar-left {
    display: flex;
    align-items: center;
    gap: 32px;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--color-primary);
    font-size: 20px;
    font-weight: 700;
    text-decoration: none;

    .logo-text {
      background: linear-gradient(135deg, var(--color-primary), #667eea);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }

  .nav-links {
    display: flex;
    gap: 4px;

    .nav-link {
      padding: 8px 16px;
      border-radius: 6px;
      color: var(--color-text-regular);
      font-size: 15px;
      transition: all 0.2s;
      text-decoration: none;

      &:hover {
        color: var(--color-primary);
        background: var(--color-primary-light-9);
      }

      &.active {
        color: var(--color-primary);
        font-weight: 600;
        background: var(--color-primary-light-9);
      }
    }
  }

  .navbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .nav-icon-badge {
    margin-right: 4px;
  }

  .nav-icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 8px;
    color: var(--color-text-regular);
    transition: all 0.2s;
    text-decoration: none;

    &:hover {
      color: var(--color-primary);
      background: var(--color-primary-light-9);
    }
  }

  .user-avatar-wrap {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 8px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;

    &:hover {
      background: var(--bg-color);
    }

    .user-name {
      font-size: 14px;
      color: var(--color-text-primary);
      max-width: 100px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .arrow {
      font-size: 12px;
      color: var(--color-text-secondary);
    }
  }
}

.main-wrapper {
  margin-top: var(--header-height);
  flex: 1;
}

.main-content {
  min-height: calc(100vh - var(--header-height) - 60px);
}

.footer {
  background: var(--bg-color-white);
  border-top: 1px solid var(--border-color-light);
  padding: 20px 0;

  .footer-inner {
    max-width: 1400px;
    margin: 0 auto;
    text-align: center;
    color: var(--color-text-secondary);
    font-size: 13px;
  }
}

@media screen and (max-width: 768px) {
  .navbar {
    .navbar-inner {
      padding: 0 12px;
    }
    .nav-links {
      gap: 0;
      .nav-link {
        padding: 8px 10px;
        font-size: 14px;
      }
    }
  }
}
</style>

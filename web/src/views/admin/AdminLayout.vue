<template>
  <div class="admin-layout">
    <el-container>
      <!-- Sidebar -->
      <el-aside :width="isCollapsed ? '64px' : '220px'" class="admin-aside">
        <div class="aside-header">
          <h2 v-if="!isCollapsed" class="aside-title">后台管理</h2>
          <el-icon v-else :size="20" color="#409EFF"><Setting /></el-icon>
        </div>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapsed"
          router
          background-color="#1d1e2c"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
          class="aside-menu"
        >
          <el-menu-item index="/admin">
            <el-icon><DataAnalysis /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/projects">
            <el-icon><Folder /></el-icon>
            <template #title>项目管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/disputes">
            <el-icon><Warning /></el-icon>
            <template #title>争议处理</template>
          </el-menu-item>
          <el-menu-item index="/admin/finance">
            <el-icon><Money /></el-icon>
            <template #title>财务管理</template>
          </el-menu-item>
        </el-menu>

        <div class="aside-footer">
          <el-button text @click="isCollapsed = !isCollapsed" class="collapse-btn">
            <el-icon :size="18">
              <component :is="isCollapsed ? 'Expand' : 'Fold'" />
            </el-icon>
          </el-button>
        </div>
      </el-aside>

      <!-- Main Content -->
      <el-container>
        <el-header class="admin-header">
          <div class="flex-between w-full">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/admin' }">后台管理</el-breadcrumb-item>
              <el-breadcrumb-item v-if="currentPageTitle">{{ currentPageTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="flex gap-12" style="align-items:center">
              <router-link to="/">
                <el-button text type="primary" size="small">返回前台</el-button>
              </router-link>
              <el-dropdown trigger="click" @command="handleCommand">
                <span class="flex gap-8 cursor-pointer" style="align-items:center">
                  <el-avatar :size="28" :src="userStore.user?.avatar_url" icon="UserFilled" />
                  {{ userStore.user?.nickname }}
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <el-main class="admin-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapsed = ref(false)

const activeMenu = computed(() => route.path)

const pageTitles: Record<string, string> = {
  '/admin': '',
  '/admin/users': '用户管理',
  '/admin/projects': '项目管理',
  '/admin/disputes': '争议处理',
  '/admin/finance': '财务管理',
}

const currentPageTitle = computed(() => pageTitles[route.path] || '')

function handleCommand(command: string) {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else if (command === 'profile') {
    router.push('/profile')
  }
}
</script>

<style scoped lang="scss">
.admin-layout {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.admin-aside {
  background: #1d1e2c;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  overflow: hidden;

  .aside-header {
    height: var(--header-height);
    display: flex;
    align-items: center;
    justify-content: center;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);

    .aside-title {
      color: #fff;
      font-size: 18px;
      font-weight: 600;
      white-space: nowrap;
    }
  }

  .aside-menu {
    flex: 1;
    border-right: none;
    overflow-y: auto;
  }

  .aside-footer {
    padding: 12px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    text-align: center;

    .collapse-btn {
      color: #bfcbd9;
    }
  }
}

.admin-header {
  background: var(--bg-color-white);
  border-bottom: 1px solid var(--border-color-light);
  display: flex;
  align-items: center;
  height: 50px;
  padding: 0 20px;
}

.admin-main {
  background: var(--bg-color);
  overflow-y: auto;
}
</style>

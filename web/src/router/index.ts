import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { guest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { guest: true },
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: { guest: true },
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
      },
      {
        path: 'projects',
        name: 'ProjectList',
        component: () => import('@/views/project/ProjectList.vue'),
      },
      {
        path: 'projects/:id',
        name: 'ProjectDetail',
        component: () => import('@/views/project/ProjectDetail.vue'),
      },
      {
        path: 'projects/create',
        name: 'ProjectCreate',
        component: () => import('@/views/project/ProjectCreate.vue'),
        meta: { auth: true, role: ['client', 'both'] },
      },
      {
        path: 'projects/:id/edit',
        name: 'ProjectEdit',
        component: () => import('@/views/project/ProjectEdit.vue'),
        meta: { auth: true },
      },
      {
        path: 'developers',
        name: 'DeveloperList',
        component: () => import('@/views/developer/DeveloperList.vue'),
      },
      {
        path: 'developers/:id',
        name: 'DeveloperDetail',
        component: () => import('@/views/developer/DeveloperDetail.vue'),
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { auth: true },
      },
      {
        path: 'my/projects',
        name: 'MyProjects',
        component: () => import('@/views/user/MyProjects.vue'),
        meta: { auth: true },
      },
      {
        path: 'my/bids',
        name: 'MyBids',
        component: () => import('@/views/developer/MyBids.vue'),
        meta: { auth: true, role: ['developer', 'both'] },
      },
      {
        path: 'my/contracts',
        name: 'MyContracts',
        component: () => import('@/views/user/MyContracts.vue'),
        meta: { auth: true },
      },
      {
        path: 'contracts/:id',
        name: 'ContractDetail',
        component: () => import('@/views/user/ContractDetail.vue'),
        meta: { auth: true },
      },
      {
        path: 'chat',
        name: 'Chat',
        component: () => import('@/views/chat/ChatList.vue'),
        meta: { auth: true },
      },
      {
        path: 'chat/:id',
        name: 'ChatDetail',
        component: () => import('@/views/chat/ChatDetail.vue'),
        meta: { auth: true },
      },
      {
        path: 'wallet',
        name: 'Wallet',
        component: () => import('@/views/user/Wallet.vue'),
        meta: { auth: true },
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/user/Notifications.vue'),
        meta: { auth: true },
      },
    ],
  },
  {
    path: '/admin',
    component: () => import('@/views/admin/AdminLayout.vue'),
    meta: { auth: true, role: ['admin'] },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
      },
      {
        path: 'projects',
        name: 'AdminProjects',
        component: () => import('@/views/admin/Projects.vue'),
      },
      {
        path: 'disputes',
        name: 'AdminDisputes',
        component: () => import('@/views/admin/Disputes.vue'),
      },
      {
        path: 'finance',
        name: 'AdminFinance',
        component: () => import('@/views/admin/Finance.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()

  if (to.meta.auth && !userStore.isLoggedIn) {
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }

  if (to.meta.guest && userStore.isLoggedIn) {
    return next({ name: 'Home' })
  }

  if (to.meta.role && userStore.user) {
    const roles = to.meta.role as string[]
    if (!roles.includes(userStore.user.user_type) && userStore.user.user_type !== 'both') {
      return next({ name: 'Home' })
    }
  }

  if (userStore.isLoggedIn && !userStore.user) {
    await userStore.fetchUser()
  }

  next()
})

export default router

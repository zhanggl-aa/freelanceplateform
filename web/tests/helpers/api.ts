import { APIRequestContext } from '@playwright/test'

const API_BASE = 'http://localhost:8080/api/v1'

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface UserInfo {
  id: string
  email?: string
  nickname: string
  user_type: string
}

export async function apiLogin(email: string, password: string): Promise<AuthTokens & { user: UserInfo }> {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Login failed')
  return json.data
}

export async function apiRegister(data: {
  email?: string
  phone?: string
  password: string
  nickname: string
  user_type: string
}): Promise<AuthTokens & { user: UserInfo }> {
  const res = await fetch(`${API_BASE}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Register failed')
  return json.data
}

export async function apiCreateProject(token: string, data: {
  category_id: string
  title: string
  description: string
  budget_min: number
  budget_max: number
  budget_type: string
  tech_stack: string[]
  deadline?: string
}) {
  const res = await fetch(`${API_BASE}/projects`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify(data),
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Create project failed')
  return json.data
}

export async function apiPublishProject(token: string, projectId: string) {
  const res = await fetch(`${API_BASE}/projects/${projectId}/publish`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Publish project failed')
  return json.data
}

export async function apiCreateBid(token: string, projectId: string, data: {
  cover_letter: string
  estimated_days: number
  proposed_budget: number
  budget_type: string
}) {
  const res = await fetch(`${API_BASE}/projects/${projectId}/bids`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify(data),
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Create bid failed')
  return json.data
}

export async function apiAcceptBid(token: string, bidId: string) {
  const res = await fetch(`${API_BASE}/bids/${bidId}/accept`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Accept bid failed')
  return json.data
}

export async function apiRejectBid(token: string, bidId: string, message?: string) {
  const res = await fetch(`${API_BASE}/bids/${bidId}/reject`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify({ client_message: message || '' }),
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Reject bid failed')
  return json.data
}

export async function apiGetCategories(): Promise<Array<{ id: string; name: string; slug: string }>> {
  const res = await fetch(`${API_BASE}/categories`)
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Get categories failed')
  return json.data
}

export async function apiGetWalletBalance(token: string) {
  const res = await fetch(`${API_BASE}/wallet/balance`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const json = await res.json()
  if (json.code && json.code !== 0) throw new Error(json.message || 'Get wallet failed')
  return json.data
}

export function getStorageState(tokens: AuthTokens) {
  return {
    cookies: [],
    origins: [{
      origin: 'http://localhost:3002',
      localStorage: [
        { name: 'access_token', value: tokens.access_token },
        { name: 'refresh_token', value: tokens.refresh_token },
      ],
    }],
  }
}

export const TEST_USERS = {
  client1: { email: 'client1@test.com', password: 'Test123456' },
  client2: { email: 'client2@test.com', password: 'Test123456' },
  client3: { email: 'client3@test.com', password: 'Test123456' },
  dev1: { email: 'dev1@test.com', password: 'Test123456' },
  dev2: { email: 'dev2@test.com', password: 'Test123456' },
  dev3: { email: 'dev3@test.com', password: 'Test123456' },
  dev4: { email: 'dev4@test.com', password: 'Test123456' },
  dev5: { email: 'dev5@test.com', password: 'Test123456' },
  both1: { email: 'both1@test.com', password: 'Test123456' },
  both2: { email: 'both2@test.com', password: 'Test123456' },
}

// ─── Seed: create and publish test projects via API ───

const SEED_PROJECTS = [
  {
    title: '电商后台管理系统开发',
    description: '需要开发一个完整的电商后台管理系统，包含商品管理、订单管理、用户管理、数据统计等核心功能模块。要求使用Vue3 + TypeScript + Element Plus技术栈，后端对接现有Go微服务API。',
    budget_min: 30000,
    budget_max: 50000,
    budget_type: 'fixed',
    tech_stack: ['Vue.js', 'TypeScript', 'Element Plus', 'Go'],
  },
  {
    title: '微信小程序商城开发',
    description: '开发一个微信小程序商城，支持商品浏览、购物车、下单支付、订单跟踪等完整购物流程。需要对接微信支付，支持优惠券和拼团活动。UI设计稿已准备好。',
    budget_min: 20000,
    budget_max: 35000,
    budget_type: 'fixed',
    tech_stack: ['小程序', 'TypeScript', 'Node.js'],
  },
  {
    title: 'SaaS平台后端架构设计',
    description: '为即将上线的SaaS产品搭建后端架构，包括多租户隔离、RBAC权限系统、计费系统、API网关等。要求高可用、可扩展，支持水平扩容。',
    budget_min: 500,
    budget_max: 800,
    budget_type: 'hourly',
    tech_stack: ['Go', 'PostgreSQL', 'Redis', 'Docker', 'AWS'],
  },
  {
    title: 'AI智能客服系统集成',
    description: '将AI大模型能力集成到现有客服系统中，实现智能问答、工单自动分类、对话摘要生成等功能。需要对接OpenAI API，开发RAG检索增强生成模块。',
    budget_min: 40000,
    budget_max: 60000,
    budget_type: 'fixed',
    tech_stack: ['Python', 'AI/ML', 'FastAPI', 'PostgreSQL'],
  },
  {
    title: 'Flutter跨平台App开发',
    description: '开发一款跨平台移动应用（iOS + Android），核心功能包括实时通讯、地图定位、支付集成。已有UI设计稿和API文档，需要完整实现并上架应用商店。',
    budget_min: 50000,
    budget_max: 80000,
    budget_type: 'fixed',
    tech_stack: ['Flutter', 'Dart', 'Firebase'],
  },
  {
    title: '数据可视化大屏开发',
    description: '为企业运营中心开发数据可视化大屏，展示实时业务数据、KPI指标、趋势图表等。需要支持多种数据源接入，刷新频率不低于30秒。ECharts + WebSocket方案。',
    budget_min: 15000,
    budget_max: 25000,
    budget_type: 'fixed',
    tech_stack: ['Vue.js', 'ECharts', 'WebSocket', 'Node.js'],
  },
]

export interface SeededProject {
  id: string
  title: string
  category_id: string
}

export async function seedTestProjects(
  clientToken: string,
  categories: Array<{ id: string; name: string; slug: string }>
): Promise<SeededProject[]> {
  const seeded: SeededProject[] = []

  for (let i = 0; i < SEED_PROJECTS.length; i++) {
    const p = SEED_PROJECTS[i]
    const categoryId = categories[i % categories.length]?.id
    if (!categoryId) continue

    const deadline = new Date()
    deadline.setDate(deadline.getDate() + 30 + i * 7)

    try {
      const created = await apiCreateProject(clientToken, {
        category_id: categoryId,
        title: p.title,
        description: p.description,
        budget_min: p.budget_min,
        budget_max: p.budget_max,
        budget_type: p.budget_type,
        tech_stack: p.tech_stack,
        deadline: deadline.toISOString(),
      })

      const projectId = created.id
      if (projectId) {
        await apiPublishProject(clientToken, projectId)
        seeded.push({ id: projectId, title: p.title, category_id: categoryId })
      }
    } catch (e) {
      console.warn(`Seed: failed to create project "${p.title}":`, e)
    }
  }

  return seeded
}

// ─── Seed: create developer profiles via API ───

const SEED_DEVELOPERS = [
  { title: '全栈开发工程师', bio: '8年全栈开发经验，擅长Vue/React + Go/Node.js技术栈。', hourly_rate: 300, availability: 'full_time', skills: ['Vue.js', 'React', 'Go', 'Node.js', 'TypeScript'] },
  { title: '高级前端工程师', bio: '5年前端开发经验，精通React生态系统，熟悉性能优化。', hourly_rate: 250, availability: 'part_time', skills: ['React', 'TypeScript', 'Angular', 'CSS'] },
  { title: 'Python后端专家', bio: '6年Python开发经验，专注AI/ML和数据分析领域。', hourly_rate: 350, availability: 'available', skills: ['Python', 'AI/ML', 'FastAPI', 'PostgreSQL'] },
  { title: '移动端开发工程师', bio: '4年Flutter和原生开发经验，已上线10+款App。', hourly_rate: 280, availability: 'full_time', skills: ['Flutter', 'Swift', 'Kotlin'] },
  { title: 'DevOps架构师', bio: '7年运维和架构经验，精通云原生和容器化部署。', hourly_rate: 400, availability: 'part_time', skills: ['Docker', 'AWS', 'Go', 'DevOps'] },
]

export async function seedDeveloperProfiles(devTokens: Array<{ token: string; email: string }>): Promise<void> {
  for (let i = 0; i < devTokens.length && i < SEED_DEVELOPERS.length; i++) {
    const { token } = devTokens[i]
    const d = SEED_DEVELOPERS[i]

    try {
      // Create developer profile
      const res = await fetch(`${API_BASE}/developers/profile`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({
          title: d.title,
          bio: d.bio,
          hourly_rate: d.hourly_rate,
          availability: d.availability,
          experience_years: 3 + i,
        }),
      })
      const json = await res.json()
      if (json.code && json.code !== 0) {
        // Profile may already exist, try updating
        await fetch(`${API_BASE}/developers/profile`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
          body: JSON.stringify({
            title: d.title,
            bio: d.bio,
            hourly_rate: d.hourly_rate,
            availability: d.availability,
          }),
        })
      }

      // Add skills
      for (const skill of d.skills) {
        await fetch(`${API_BASE}/developers/skills`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
          body: JSON.stringify({ skill_name: skill, proficiency: 'expert', years_experience: 3 + i }),
        }).catch(() => {})
      }
    } catch (e) {
      console.warn(`Seed: failed to create developer profile for ${devTokens[i].email}:`, e)
    }
  }
}

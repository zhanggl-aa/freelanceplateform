import request from './index'

// ─── Auth API ───

export const authApi = {
  register: (data: { email?: string; phone?: string; password: string; nickname: string; user_type: string }) =>
    request.post('/auth/register', data),
  login: (data: { email?: string; phone?: string; password: string }) =>
    request.post('/auth/login', data),
  logout: (refresh_token: string) =>
    request.post('/auth/logout', { refresh_token }),
  refresh: (refresh_token: string) =>
    request.post('/auth/refresh', { refresh_token }),
  forgotPassword: (account: string) =>
    request.post('/auth/forgot-password', { account }),
  resetPassword: (data: { account: string; code: string; new_password: string }) =>
    request.post('/auth/reset-password', data),
  changePassword: (data: { old_password: string; new_password: string }) =>
    request.put('/auth/change-password', data),
  verifyEmail: () => request.post('/auth/verify-email'),
  verifyPhone: () => request.post('/auth/verify-phone'),
}

// ─── User API ───

export const userApi = {
  getMe: () => request.get('/users/me'),
  updateProfile: (data: { nickname?: string; avatar_url?: string }) =>
    request.put('/users/me', data),
  deleteAccount: () => request.delete('/users/me'),
  getUser: (id: string) => request.get(`/users/${id}`),
}

// ─── Developer API ───

export const developerApi = {
  createProfile: (data: any) => request.post('/developers/profile', data),
  getProfile: () => request.get('/developers/profile'),
  updateProfile: (data: any) => request.put('/developers/profile', data),
  getDeveloper: (id: string) => request.get(`/developers/${id}`),
  search: (params: any) => request.get('/developers', { params }),
  addSkill: (data: any) => request.post('/developers/skills', data),
  updateSkill: (id: string, data: any) => request.put(`/developers/skills/${id}`, data),
  deleteSkill: (id: string) => request.delete(`/developers/skills/${id}`),
  addPortfolio: (data: any) => request.post('/developers/portfolio', data),
  updatePortfolio: (id: string, data: any) => request.put(`/developers/portfolio/${id}`, data),
  deletePortfolio: (id: string) => request.delete(`/developers/portfolio/${id}`),
}

// ─── Client API ───

export const clientApi = {
  createProfile: (data: any) => request.post('/clients/profile', data),
  getProfile: () => request.get('/clients/profile'),
  updateProfile: (data: any) => request.put('/clients/profile', data),
  verify: (data: any) => request.post('/clients/verify', data),
}

// ─── Category API ───

export const categoryApi = {
  getTree: () => request.get('/categories'),
  getById: (id: string) => request.get(`/categories/${id}`),
  create: (data: any) => request.post('/categories', data),
  update: (id: string, data: any) => request.put(`/categories/${id}`, data),
  delete: (id: string) => request.delete(`/categories/${id}`),
}

// ─── Project API ───

export const projectApi = {
  create: (data: any) => request.post('/projects', data),
  search: (params: any) => request.get('/projects', { params }),
  getById: (id: string) => request.get(`/projects/${id}`),
  update: (id: string, data: any) => request.put(`/projects/${id}`, data),
  delete: (id: string) => request.delete(`/projects/${id}`),
  publish: (id: string) => request.post(`/projects/${id}/publish`),
  close: (id: string) => request.post(`/projects/${id}/close`),
  uploadAttachment: (id: string, file: File) => {
    const fd = new FormData()
    fd.append('file', file)
    return request.post(`/projects/${id}/attachments`, fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  bookmark: (id: string) => request.post(`/projects/${id}/bookmark`),
  removeBookmark: (id: string) => request.delete(`/projects/${id}/bookmark`),
  myPosted: (params?: any) => request.get('/projects/my/posted', { params }),
  myBidding: (params?: any) => request.get('/projects/my/bidding', { params }),
  myWorking: (params?: any) => request.get('/projects/my/working', { params }),
  myCompleted: (params?: any) => request.get('/projects/my/completed', { params }),
}

// ─── Bid API ───

export const bidApi = {
  create: (projectId: string, data: any) =>
    request.post(`/projects/${projectId}/bids`, data),
  listByProject: (projectId: string, params?: any) =>
    request.get(`/projects/${projectId}/bids`, { params }),
  getById: (id: string) => request.get(`/bids/${id}`),
  update: (id: string, data: any) => request.put(`/bids/${id}`, data),
  withdraw: (id: string) => request.delete(`/bids/${id}`),
  accept: (id: string) => request.post(`/bids/${id}/accept`),
  reject: (id: string, data?: any) => request.post(`/bids/${id}/reject`, data),
  shortlist: (id: string) => request.post(`/bids/${id}/shortlist`),
  counterOffer: (id: string, data: any) =>
    request.post(`/bids/${id}/counter-offer`, data),
  myBids: (params?: any) => request.get('/developers/me/bids', { params }),
}

// ─── Contract API ───

export const contractApi = {
  getById: (id: string) => request.get(`/contracts/${id}`),
  myContracts: (params?: any) => request.get('/contracts/my', { params }),
  update: (id: string, data: any) => request.put(`/contracts/${id}`, data),
  start: (id: string) => request.post(`/contracts/${id}/start`),
  cancel: (id: string) => request.post(`/contracts/${id}/cancel`),
  dispute: (id: string, data: any) => request.post(`/contracts/${id}/dispute`, data),
}

// ─── Milestone API ───

export const milestoneApi = {
  create: (projectId: string, data: any) =>
    request.post(`/projects/${projectId}/milestones`, data),
  listByProject: (projectId: string) =>
    request.get(`/projects/${projectId}/milestones`),
  update: (id: string, data: any) => request.put(`/milestones/${id}`, data),
  delete: (id: string) => request.delete(`/milestones/${id}`),
  submit: (id: string, data: any) => request.post(`/milestones/${id}/submit`, data),
  approve: (id: string) => request.post(`/milestones/${id}/approve`),
  reject: (id: string, data: any) => request.post(`/milestones/${id}/reject`, data),
  dispute: (id: string, data: any) => request.post(`/milestones/${id}/dispute`, data),
}

// ─── Payment API ───

export const paymentApi = {
  deposit: (data: any) => request.post('/payments/deposit', data),
  release: (data: any) => request.post('/payments/release', data),
  refund: (data: any) => request.post('/payments/refund', data),
  getById: (id: string) => request.get(`/payments/${id}`),
  myPayments: (params?: any) => request.get('/payments/my', { params }),
  walletBalance: () => request.get('/wallet/balance'),
  walletTransactions: (params?: any) =>
    request.get('/wallet/transactions', { params }),
  withdraw: (data: any) => request.post('/wallet/withdraw', data),
}

// ─── Chat API ───

export const chatApi = {
  list: (params?: any) => request.get('/conversations', { params }),
  getMessages: (id: string, params?: any) =>
    request.get(`/conversations/${id}`, { params }),
  create: (data: any) => request.post('/conversations', data),
  sendMessage: (id: string, data: any) =>
    request.post(`/conversations/${id}/messages`, data),
  sendFile: (id: string, file: File) => {
    const fd = new FormData()
    fd.append('file', file)
    return request.post(`/conversations/${id}/files`, fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  markAsRead: (id: string) => request.put(`/conversations/${id}/read`),
  unreadCount: () => request.get('/conversations/unread-count'),
}

// ─── Review API ───

export const reviewApi = {
  create: (projectId: string, data: any) =>
    request.post(`/projects/${projectId}/reviews`, data),
  getByProject: (projectId: string) =>
    request.get(`/projects/${projectId}/reviews`),
  getByUser: (userId: string, params?: any) =>
    request.get(`/users/${userId}/reviews`, { params }),
  update: (id: string, data: any) => request.put(`/reviews/${id}`, data),
  delete: (id: string) => request.delete(`/reviews/${id}`),
}

// ─── Notification API ───

export const notificationApi = {
  list: (params?: any) => request.get('/notifications', { params }),
  unreadCount: () => request.get('/notifications/unread-count'),
  markAsRead: (id: string) => request.put(`/notifications/${id}/read`),
  markAllAsRead: () => request.put('/notifications/read-all'),
  getSettings: () => request.get('/notifications/settings'),
  updateSettings: (data: any) => request.put('/notifications/settings', data),
}

// ─── File API ───

export const fileApi = {
  upload: (file: File, entityType?: string, entityId?: string) => {
    const fd = new FormData()
    fd.append('file', file)
    if (entityType) fd.append('entity_type', entityType)
    if (entityId) fd.append('entity_id', entityId)
    return request.post('/files/upload', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  getById: (id: string) => request.get(`/files/${id}`),
  download: (id: string) => request.get(`/files/${id}/download`, { responseType: 'blob' }),
  delete: (id: string) => request.delete(`/files/${id}`),
}

// ─── Admin API ───

export const adminApi = {
  dashboard: () => request.get('/admin/dashboard'),
  listUsers: (params?: any) => request.get('/admin/users', { params }),
  updateUserStatus: (id: string, status: string) =>
    request.put(`/admin/users/${id}/status`, { status }),
  listProjects: (params?: any) => request.get('/admin/projects', { params }),
  moderateProject: (id: string, data: any) =>
    request.put(`/admin/projects/${id}`, data),
  listDisputes: (params?: any) => request.get('/admin/disputes', { params }),
  resolveDispute: (id: string, data: any) =>
    request.put(`/admin/disputes/${id}`, data),
  listPayments: (params?: any) => request.get('/admin/payments', { params }),
  financialSummary: () => request.get('/admin/finance/summary'),
}

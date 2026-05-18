export interface Project {
  id: string
  title: string
  description: string
  status: 'open' | 'in_progress' | 'completed' | 'cancelled' | 'draft'
  category: string
  categoryName: string
  budgetType: 'fixed' | 'range' | 'hourly'
  budgetMin: number
  budgetMax: number
  deadline: string
  techStack: string[]
  bidCount: number
  clientId: string
  clientName: string
  clientAvatar: string
  createdAt: string
  updatedAt: string
}

export interface Developer {
  id: string
  userId: string
  nickname: string
  avatar: string
  title: string
  bio: string
  skills: string[]
  hourlyRate: number
  availability: 'available' | 'busy' | 'unavailable'
  rating: number
  completedProjects: number
  portfolio: PortfolioItem[]
}

export interface PortfolioItem {
  id: string
  title: string
  description: string
  images: string[]
  url: string
}

export interface Bid {
  id: string
  projectId: string
  projectTitle: string
  developerId: string
  developerName: string
  developerAvatar: string
  coverLetter: string
  estimatedDays: number
  proposedBudget: number
  status: 'pending' | 'accepted' | 'rejected' | 'withdrawn'
  createdAt: string
}

export interface Contract {
  id: string
  projectId: string
  projectTitle: string
  clientId: string
  clientName: string
  developerId: string
  developerName: string
  amount: number
  status: 'pending' | 'active' | 'completed' | 'cancelled' | 'disputed'
  startDate: string
  endDate: string
  milestones: Milestone[]
}

export interface Milestone {
  id: string
  title: string
  description: string
  amount: number
  dueDate: string
  status: 'pending' | 'in_progress' | 'submitted' | 'approved' | 'rejected'
}

export interface Conversation {
  id: string
  participantId: string
  participantName: string
  participantAvatar: string
  lastMessage: string
  lastMessageTime: string
  unreadCount: number
}

export interface Message {
  id: string
  conversationId: string
  senderId: string
  content: string
  type: 'text' | 'image' | 'file'
  fileUrl?: string
  createdAt: string
}

export interface Transaction {
  id: string
  type: 'deposit' | 'withdraw' | 'payment' | 'refund' | 'income'
  amount: number
  description: string
  status: 'pending' | 'completed' | 'failed'
  createdAt: string
}

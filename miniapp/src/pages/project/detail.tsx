import { useState, useEffect } from 'react'
import { View, Text, Image, ScrollView, Input, Textarea, Button } from '@tarojs/components'
import Taro, { useRouter } from '@tarojs/taro'
import { api } from '@/services/api'
import { useUser } from '@/store/user'
import { Project } from '@/components/types'
import Empty from '@/components/Empty'
import './detail.scss'

export default function ProjectDetail() {
  const router = useRouter()
  const { user, isLoggedIn } = useUser()
  const [project, setProject] = useState<Project | null>(null)
  const [loading, setLoading] = useState(true)
  const [showBidForm, setShowBidForm] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [bidForm, setBidForm] = useState({
    coverLetter: '',
    estimatedDays: '',
    proposedBudget: '',
  })

  useEffect(() => {
    if (router.params.id) {
      loadProject(router.params.id)
    }
  }, [router.params.id])

  async function loadProject(id: string) {
    setLoading(true)
    try {
      const res = await api.project.getDetail(id)
      if (res.code === 0) {
        setProject(res.data)
        Taro.setNavigationBarTitle({ title: res.data?.title || '项目详情' })
      }
    } catch (err) {
      // handle error
    } finally {
      setLoading(false)
    }
  }

  async function handleSubmitBid() {
    if (!isLoggedIn) {
      Taro.navigateTo({ url: '/pages/auth/login' })
      return
    }
    if (!bidForm.coverLetter.trim()) {
      Taro.showToast({ title: '请填写求职信', icon: 'none' })
      return
    }
    if (!bidForm.estimatedDays || Number(bidForm.estimatedDays) <= 0) {
      Taro.showToast({ title: '请填写预计天数', icon: 'none' })
      return
    }
    if (!bidForm.proposedBudget || Number(bidForm.proposedBudget) <= 0) {
      Taro.showToast({ title: '请填写报价', icon: 'none' })
      return
    }

    setSubmitting(true)
    try {
      const res = await api.bid.create({
        projectId: project!.id,
        coverLetter: bidForm.coverLetter,
        estimatedDays: Number(bidForm.estimatedDays),
        proposedBudget: Number(bidForm.proposedBudget),
      })
      if (res.code === 0) {
        Taro.showToast({ title: '投标成功', icon: 'success' })
        setShowBidForm(false)
        setBidForm({ coverLetter: '', estimatedDays: '', proposedBudget: '' })
        loadProject(project!.id)
      }
    } catch (err) {
      // handled in api
    } finally {
      setSubmitting(false)
    }
  }

  const statusMap: Record<string, string> = {
    open: '竞标中',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    draft: '草稿',
  }

  if (loading) {
    return (
      <View className="project-detail-page__loading">
        <Text>加载中...</Text>
      </View>
    )
  }

  if (!project) {
    return <Empty type="nodata" description="项目不存在" />
  }

  return (
    <View className="project-detail-page">
      <ScrollView scrollY className="project-detail-page__content">
        {/* Header */}
        <View className="project-detail-page__header">
          <Text className="project-detail-page__title">{project.title}</Text>
          <View className="project-detail-page__meta">
            <Text className="project-detail-page__status">{statusMap[project.status] || project.status}</Text>
            <Text className="project-detail-page__category">{project.categoryName}</Text>
          </View>
        </View>

        {/* Client Info */}
        <View className="project-detail-page__client">
          <Image
            className="project-detail-page__client-avatar"
            src={project.clientAvatar || 'https://via.placeholder.com/64'}
            mode="aspectFill"
          />
          <View className="project-detail-page__client-info">
            <Text className="project-detail-page__client-name">{project.clientName}</Text>
            <Text className="project-detail-page__client-label">发布者</Text>
          </View>
          <View className="project-detail-page__client-chat" onClick={() => {
            if (!isLoggedIn) {
              Taro.navigateTo({ url: '/pages/auth/login' })
              return
            }
            Taro.navigateTo({ url: `/pages/chat/detail?userId=${project.clientId}` })
          }}>
            <Text className="project-detail-page__client-chat-text">联系</Text>
          </View>
        </View>

        {/* Budget & Deadline */}
        <View className="project-detail-page__info-row">
          <View className="project-detail-page__info-item">
            <Text className="project-detail-page__info-label">预算</Text>
            <Text className="project-detail-page__info-value project-detail-page__info-value--price">
              {project.budgetType === 'fixed' ? `¥${project.budgetMin}` : `¥${project.budgetMin}-${project.budgetMax}`}
            </Text>
          </View>
          <View className="project-detail-page__info-item">
            <Text className="project-detail-page__info-label">截止日期</Text>
            <Text className="project-detail-page__info-value">{project.deadline}</Text>
          </View>
          <View className="project-detail-page__info-item">
            <Text className="project-detail-page__info-label">竞标</Text>
            <Text className="project-detail-page__info-value">{project.bidCount || 0}人</Text>
          </View>
        </View>

        {/* Description */}
        <View className="project-detail-page__section">
          <Text className="project-detail-page__section-title">项目描述</Text>
          <Text className="project-detail-page__desc">{project.description}</Text>
        </View>

        {/* Tech Stack */}
        {project.techStack && project.techStack.length > 0 && (
          <View className="project-detail-page__section">
            <Text className="project-detail-page__section-title">技术要求</Text>
            <View className="project-detail-page__tech-list">
              {project.techStack.map((tech, index) => (
                <Text key={index} className="project-detail-page__tech-tag">{tech}</Text>
              ))}
            </View>
          </View>
        )}
      </ScrollView>

      {/* Bid Form Modal */}
      {showBidForm && (
        <View className="project-detail-page__bid-modal">
          <View className="project-detail-page__bid-form">
            <View className="project-detail-page__bid-form-header">
              <Text className="project-detail-page__bid-form-title">提交投标</Text>
              <Text className="project-detail-page__bid-form-close" onClick={() => setShowBidForm(false)}>✕</Text>
            </View>
            <View className="project-detail-page__bid-form-body">
              <View className="project-detail-page__form-group">
                <Text className="project-detail-page__form-label">求职信</Text>
                <Textarea
                  className="project-detail-page__form-textarea"
                  placeholder="介绍你的经验和优势，说明你如何完成该项目"
                  value={bidForm.coverLetter}
                  onInput={(e) => setBidForm(prev => ({ ...prev, coverLetter: e.detail.value }))}
                  maxlength={500}
                />
              </View>
              <View className="project-detail-page__form-row">
                <View className="project-detail-page__form-group project-detail-page__form-group--half">
                  <Text className="project-detail-page__form-label">预计天数</Text>
                  <Input
                    className="project-detail-page__form-input"
                    type="number"
                    placeholder="预计完成天数"
                    value={bidForm.estimatedDays}
                    onInput={(e) => setBidForm(prev => ({ ...prev, estimatedDays: e.detail.value }))}
                  />
                </View>
                <View className="project-detail-page__form-group project-detail-page__form-group--half">
                  <Text className="project-detail-page__form-label">报价(元)</Text>
                  <Input
                    className="project-detail-page__form-input"
                    type="digit"
                    placeholder="你的报价"
                    value={bidForm.proposedBudget}
                    onInput={(e) => setBidForm(prev => ({ ...prev, proposedBudget: e.detail.value }))}
                  />
                </View>
              </View>
            </View>
            <View className="project-detail-page__bid-form-footer">
              <Button
                className="project-detail-page__bid-btn project-detail-page__bid-btn--cancel"
                onClick={() => setShowBidForm(false)}
              >
                取消
              </Button>
              <Button
                className="project-detail-page__bid-btn project-detail-page__bid-btn--submit"
                loading={submitting}
                onClick={handleSubmitBid}
              >
                提交投标
              </Button>
            </View>
          </View>
        </View>
      )}

      {/* Bottom Action */}
      {project.status === 'open' && (
        <View className="project-detail-page__bottom">
          <Button
            className="project-detail-page__bottom-btn"
            onClick={() => {
              if (!isLoggedIn) {
                Taro.navigateTo({ url: '/pages/auth/login' })
                return
              }
              setShowBidForm(true)
            }}
          >
            立即投标
          </Button>
        </View>
      )}
    </View>
  )
}

detail.config = {
  navigationBarTitleText: '项目详情',
}

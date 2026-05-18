import { useState, useEffect } from 'react'
import { View, Text, Input, Textarea, Button, Picker, ScrollView } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { api } from '@/services/api'
import { useUser } from '@/store/user'
import './create.scss'

const CATEGORY_OPTIONS = [
  '请选择分类', '网站开发', '移动应用', '小程序', '前端', '后端', 'AI/ML', 'UI设计', '数据分析', '运维部署'
]

const BUDGET_TYPE_OPTIONS = ['固定预算', '预算范围', '时薪预算']

export default function ProjectCreate() {
  const { isLoggedIn } = useUser()
  const [submitting, setSubmitting] = useState(false)
  const [savingDraft, setSavingDraft] = useState(false)
  const [categoryIndex, setCategoryIndex] = useState(0)
  const [budgetTypeIndex, setBudgetTypeIndex] = useState(0)
  const [form, setForm] = useState({
    title: '',
    description: '',
    budgetMin: '',
    budgetMax: '',
    deadline: '',
    techStack: '',
  })

  useEffect(() => {
    if (!isLoggedIn) {
      Taro.navigateTo({ url: '/pages/auth/login' })
    }
  }, [isLoggedIn])

  function updateForm(field: string, value: string) {
    setForm(prev => ({ ...prev, [field]: value }))
  }

  function validateForm(): boolean {
    if (!form.title.trim()) {
      Taro.showToast({ title: '请输入项目标题', icon: 'none' })
      return false
    }
    if (categoryIndex === 0) {
      Taro.showToast({ title: '请选择项目分类', icon: 'none' })
      return false
    }
    if (!form.description.trim()) {
      Taro.showToast({ title: '请输入项目描述', icon: 'none' })
      return false
    }
    if (!form.budgetMin || Number(form.budgetMin) <= 0) {
      Taro.showToast({ title: '请输入预算金额', icon: 'none' })
      return false
    }
    if (budgetTypeIndex === 1 && (!form.budgetMax || Number(form.budgetMax) <= Number(form.budgetMin))) {
      Taro.showToast({ title: '请输入正确的预算上限', icon: 'none' })
      return false
    }
    if (!form.deadline) {
      Taro.showToast({ title: '请选择截止日期', icon: 'none' })
      return false
    }
    return true
  }

  async function handlePublish() {
    if (!validateForm()) return

    setSubmitting(true)
    try {
      const budgetType = ['fixed', 'range', 'hourly'][budgetTypeIndex]
      const techStack = form.techStack
        ? form.techStack.split(/[,，、]/).map(s => s.trim()).filter(Boolean)
        : []

      const res = await api.project.create({
        title: form.title,
        category: CATEGORY_OPTIONS[categoryIndex],
        description: form.description,
        budgetType,
        budgetMin: Number(form.budgetMin),
        budgetMax: form.budgetMax ? Number(form.budgetMax) : Number(form.budgetMin),
        deadline: form.deadline,
        techStack,
      })

      if (res.code === 0) {
        Taro.showToast({ title: '发布成功', icon: 'success' })
        setTimeout(() => {
          Taro.navigateBack()
        }, 1500)
      }
    } catch (err) {
      // handled in api
    } finally {
      setSubmitting(false)
    }
  }

  async function handleSaveDraft() {
    if (!form.title.trim()) {
      Taro.showToast({ title: '请至少填写项目标题', icon: 'none' })
      return
    }

    setSavingDraft(true)
    try {
      const budgetType = ['fixed', 'range', 'hourly'][budgetTypeIndex]
      const techStack = form.techStack
        ? form.techStack.split(/[,，、]/).map(s => s.trim()).filter(Boolean)
        : []

      const res = await api.project.saveDraft({
        title: form.title,
        category: categoryIndex > 0 ? CATEGORY_OPTIONS[categoryIndex] : '',
        description: form.description,
        budgetType,
        budgetMin: form.budgetMin ? Number(form.budgetMin) : 0,
        budgetMax: form.budgetMax ? Number(form.budgetMax) : 0,
        deadline: form.deadline,
        techStack,
      })

      if (res.code === 0) {
        Taro.showToast({ title: '草稿已保存', icon: 'success' })
      }
    } catch (err) {
      // handled
    } finally {
      setSavingDraft(false)
    }
  }

  return (
    <View className="project-create-page">
      <ScrollView scrollY className="project-create-page__content">
        {/* Title */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">项目标题</Text>
          <Input
            className="project-create-page__input"
            placeholder="输入项目标题"
            value={form.title}
            onInput={(e) => updateForm('title', e.detail.value)}
            maxlength={50}
          />
        </View>

        {/* Category */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">项目分类</Text>
          <Picker mode="selector" range={CATEGORY_OPTIONS} value={categoryIndex} onChange={(e) => setCategoryIndex(Number(e.detail.value))}>
            <View className="project-create-page__picker">
              <Text className={`project-create-page__picker-text ${categoryIndex === 0 ? 'project-create-page__picker-text--placeholder' : ''}`}>
                {CATEGORY_OPTIONS[categoryIndex]}
              </Text>
              <Text className="project-create-page__picker-arrow">▼</Text>
            </View>
          </Picker>
        </View>

        {/* Description */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">项目描述</Text>
          <Textarea
            className="project-create-page__textarea"
            placeholder="详细描述你的项目需求、功能要求、交付标准等"
            value={form.description}
            onInput={(e) => updateForm('description', e.detail.value)}
            maxlength={2000}
          />
        </View>

        {/* Budget Type */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">预算方式</Text>
          <Picker mode="selector" range={BUDGET_TYPE_OPTIONS} value={budgetTypeIndex} onChange={(e) => setBudgetTypeIndex(Number(e.detail.value))}>
            <View className="project-create-page__picker">
              <Text className="project-create-page__picker-text">{BUDGET_TYPE_OPTIONS[budgetTypeIndex]}</Text>
              <Text className="project-create-page__picker-arrow">▼</Text>
            </View>
          </Picker>
        </View>

        {/* Budget Range */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">预算金额(元)</Text>
          <View className="project-create-page__budget-row">
            <Input
              className="project-create-page__input project-create-page__input--budget"
              type="digit"
              placeholder="最低预算"
              value={form.budgetMin}
              onInput={(e) => updateForm('budgetMin', e.detail.value)}
            />
            {budgetTypeIndex === 1 && (
              <>
                <Text className="project-create-page__budget-sep">—</Text>
                <Input
                  className="project-create-page__input project-create-page__input--budget"
                  type="digit"
                  placeholder="最高预算"
                  value={form.budgetMax}
                  onInput={(e) => updateForm('budgetMax', e.detail.value)}
                />
              </>
            )}
          </View>
        </View>

        {/* Deadline */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">截止日期</Text>
          <Picker mode="date" value={form.deadline || ''} onChange={(e) => updateForm('deadline', e.detail.value)}>
            <View className="project-create-page__picker">
              <Text className={`project-create-page__picker-text ${!form.deadline ? 'project-create-page__picker-text--placeholder' : ''}`}>
                {form.deadline || '请选择截止日期'}
              </Text>
              <Text className="project-create-page__picker-arrow">▼</Text>
            </View>
          </Picker>
        </View>

        {/* Tech Stack */}
        <View className="project-create-page__form-group">
          <Text className="project-create-page__label">技术栈</Text>
          <Input
            className="project-create-page__input"
            placeholder="用逗号分隔，如: React, Node.js, MySQL"
            value={form.techStack}
            onInput={(e) => updateForm('techStack', e.detail.value)}
          />
        </View>
      </ScrollView>

      {/* Bottom Buttons */}
      <View className="project-create-page__bottom">
        <Button
          className="project-create-page__btn project-create-page__btn--draft"
          loading={savingDraft}
          onClick={handleSaveDraft}
        >
          存草稿
        </Button>
        <Button
          className="project-create-page__btn project-create-page__btn--publish"
          loading={submitting}
          onClick={handlePublish}
        >
          发布项目
        </Button>
      </View>
    </View>
  )
}

create.config = {
  navigationBarTitleText: '发布项目',
}

import { useState, useEffect } from 'react'
import { View, Text, Input, ScrollView, Image, Navigator } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { api } from '@/services/api'
import ProjectCard from '@/components/ProjectCard'
import DeveloperCard from '@/components/DeveloperCard'
import Empty from '@/components/Empty'
import { Project, Developer } from '@/components/types'
import './index.scss'

const CATEGORIES = [
  { name: '网站开发', icon: '🌐' },
  { name: '移动应用', icon: '📱' },
  { name: '小程序', icon: '🔧' },
  { name: '前端', icon: '💻' },
  { name: '后端', icon: '⚙️' },
  { name: 'AI/ML', icon: '🤖' },
  { name: 'UI设计', icon: '🎨' },
  { name: '更多', icon: '📋' },
]

export default function Index() {
  const [keyword, setKeyword] = useState('')
  const [featuredProjects, setFeaturedProjects] = useState<Project[]>([])
  const [topDevelopers, setTopDevelopers] = useState<Developer[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  async function loadData() {
    setLoading(true)
    try {
      const [projectRes, devRes] = await Promise.all([
        api.project.getList({ page: 1, pageSize: 6, sort: 'featured' }),
        api.developer.getList({ page: 1, pageSize: 8, sort: 'rating' }),
      ])
      if (projectRes.code === 0) {
        setFeaturedProjects(projectRes.data?.list || projectRes.data || [])
      }
      if (devRes.code === 0) {
        setTopDevelopers(devRes.data?.list || devRes.data || [])
      }
    } catch (err) {
      // handle error
    } finally {
      setLoading(false)
    }
  }

  function handleSearch() {
    if (!keyword.trim()) return
    Taro.navigateTo({ url: `/pages/project/list?keyword=${encodeURIComponent(keyword)}` })
  }

  function handleCategoryTap(name: string) {
    if (name === '更多') {
      Taro.switchTab({ url: '/pages/project/list' })
      return
    }
    Taro.navigateTo({ url: `/pages/project/list?category=${encodeURIComponent(name)}` })
  }

  return (
    <View className="index-page">
      {/* Search Bar */}
      <View className="index-page__search">
        <Input
          className="index-page__search-input"
          placeholder="搜索项目或开发者"
          placeholderClass="index-page__search-placeholder"
          value={keyword}
          onInput={(e) => setKeyword(e.detail.value)}
          onConfirm={handleSearch}
          confirmType="search"
        />
        <View className="index-page__search-btn" onClick={handleSearch}>
          <Text className="index-page__search-btn-text">搜索</Text>
        </View>
      </View>

      {/* Category Grid */}
      <View className="index-page__categories">
        {CATEGORIES.map((cat, index) => (
          <View
            key={index}
            className="index-page__category-item"
            onClick={() => handleCategoryTap(cat.name)}
          >
            <Text className="index-page__category-icon">{cat.icon}</Text>
            <Text className="index-page__category-name">{cat.name}</Text>
          </View>
        ))}
      </View>

      {/* Featured Projects */}
      <View className="index-page__section">
        <View className="index-page__section-header">
          <Text className="index-page__section-title">精选项目</Text>
          <Navigator url="/pages/project/list" className="index-page__section-more">
            <Text>查看更多</Text>
          </Navigator>
        </View>
        {loading ? (
          <View className="index-page__loading">
            <Text>加载中...</Text>
          </View>
        ) : featuredProjects.length > 0 ? (
          featuredProjects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))
        ) : (
          <Empty type="nodata" description="暂无精选项目" />
        )}
      </View>

      {/* Top Developers */}
      <View className="index-page__section">
        <View className="index-page__section-header">
          <Text className="index-page__section-title">优秀开发者</Text>
          <Navigator url="/pages/developer/list" className="index-page__section-more">
            <Text>查看更多</Text>
          </Navigator>
        </View>
        {loading ? (
          <View className="index-page__loading">
            <Text>加载中...</Text>
          </View>
        ) : topDevelopers.length > 0 ? (
          <ScrollView scrollX className="index-page__dev-scroll" enhanced showScrollbar={false}>
            <View className="index-page__dev-list">
              {topDevelopers.map((dev) => (
                <DeveloperCard key={dev.id} developer={dev} horizontal />
              ))}
            </View>
          </ScrollView>
        ) : (
          <Empty type="nodata" description="暂无开发者" />
        )}
      </View>
    </View>
  )
}

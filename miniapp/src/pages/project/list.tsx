import { useState, useEffect, useCallback } from 'react'
import { View, Text, ScrollView, Picker } from '@tarojs/components'
import Taro, { usePullDownRefresh, useReachBottom, useRouter } from '@tarojs/taro'
import { api } from '@/services/api'
import ProjectCard from '@/components/ProjectCard'
import Empty from '@/components/Empty'
import { Project } from '@/components/types'
import './list.scss'

const SORT_TABS = [
  { key: 'all', label: '全部' },
  { key: 'latest', label: '最新' },
  { key: 'budget_high', label: '预算最高' },
  { key: 'bids_most', label: '竞标最多' },
]

const CATEGORY_OPTIONS = [
  '全部分类', '网站开发', '移动应用', '小程序', '前端', '后端', 'AI/ML', 'UI设计', '数据分析', '运维部署'
]

export default function ProjectList() {
  const router = useRouter()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(false)
  const [refreshing, setRefreshing] = useState(false)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  const [activeSort, setActiveSort] = useState('all')
  const [categoryIndex, setCategoryIndex] = useState(0)
  const [keyword, setKeyword] = useState('')

  useEffect(() => {
    const { category, keyword: kw } = router.params
    if (category) {
      const idx = CATEGORY_OPTIONS.indexOf(decodeURIComponent(category))
      if (idx > -1) setCategoryIndex(idx)
    }
    if (kw) {
      setKeyword(decodeURIComponent(kw))
    }
  }, [router.params])

  useEffect(() => {
    loadProjects(true)
  }, [activeSort, categoryIndex])

  const loadProjects = useCallback(async (reset = false) => {
    if (loading) return
    const currentPage = reset ? 1 : page
    setLoading(true)

    try {
      const params: any = {
        page: currentPage,
        pageSize: 10,
        sort: activeSort,
      }
      if (categoryIndex > 0) {
        params.category = CATEGORY_OPTIONS[categoryIndex]
      }
      if (keyword) {
        params.keyword = keyword
      }

      const res = await api.project.getList(params)
      if (res.code === 0) {
        const list = res.data?.list || res.data || []
        const total = res.meta?.total || 0
        setProjects(prev => reset ? list : [...prev, ...list])
        setHasMore(currentPage * 10 < total)
        setPage(currentPage + 1)
      }
    } catch (err) {
      // handle error
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [page, activeSort, categoryIndex, keyword, loading])

  usePullDownRefresh(() => {
    setRefreshing(true)
    loadProjects(true)
    Taro.stopPullDownRefresh()
  })

  useReachBottom(() => {
    if (hasMore && !loading) {
      loadProjects()
    }
  })

  function handleSortChange(key: string) {
    setActiveSort(key)
  }

  function handleCategoryChange(e) {
    setCategoryIndex(Number(e.detail.value))
  }

  return (
    <View className="project-list-page">
      {/* Sort Tabs */}
      <View className="project-list-page__tabs">
        {SORT_TABS.map((tab) => (
          <View
            key={tab.key}
            className={`project-list-page__tab ${activeSort === tab.key ? 'project-list-page__tab--active' : ''}`}
            onClick={() => handleSortChange(tab.key)}
          >
            <Text className="project-list-page__tab-text">{tab.label}</Text>
          </View>
        ))}
      </View>

      {/* Category Filter */}
      <View className="project-list-page__filter">
        <Picker mode="selector" range={CATEGORY_OPTIONS} value={categoryIndex} onChange={handleCategoryChange}>
          <View className="project-list-page__filter-picker">
            <Text className="project-list-page__filter-text">
              {CATEGORY_OPTIONS[categoryIndex]}
            </Text>
            <Text className="project-list-page__filter-arrow">▼</Text>
          </View>
        </Picker>
      </View>

      {/* Project List */}
      <ScrollView scrollY className="project-list-page__content">
        {projects.length > 0 ? (
          projects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))
        ) : !loading ? (
          <Empty type="search" description="暂无匹配的项目" />
        ) : null}

        {loading && (
          <View className="project-list-page__loading">
            <Text>加载中...</Text>
          </View>
        )}

        {!hasMore && projects.length > 0 && (
          <View className="project-list-page__no-more">
            <Text>没有更多了</Text>
          </View>
        )}
      </ScrollView>
    </View>
  )
}

list.config = {
  navigationBarTitleText: '项目列表',
  enablePullDownRefresh: true,
}

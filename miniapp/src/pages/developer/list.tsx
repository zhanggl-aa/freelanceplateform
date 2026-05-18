import { useState, useEffect, useCallback } from 'react'
import { View, Text, Input, ScrollView, Picker } from '@tarojs/components'
import Taro, { usePullDownRefresh, useReachBottom } from '@tarojs/taro'
import { api } from '@/services/api'
import DeveloperCard from '@/components/DeveloperCard'
import Empty from '@/components/Empty'
import { Developer } from '@/components/types'
import './list.scss'

const RATE_OPTIONS = ['不限时薪', '¥0-100/时', '¥100-200/时', '¥200-300/时', '¥300+/时']
const AVAILABILITY_OPTIONS = ['不限状态', '空闲', '忙碌']

export default function DeveloperList() {
  const [developers, setDevelopers] = useState<Developer[]>([])
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  const [keyword, setKeyword] = useState('')
  const [rateIndex, setRateIndex] = useState(0)
  const [availabilityIndex, setAvailabilityIndex] = useState(0)

  useEffect(() => {
    loadDevelopers(true)
  }, [rateIndex, availabilityIndex])

  const loadDevelopers = useCallback(async (reset = false) => {
    if (loading) return
    const currentPage = reset ? 1 : page
    setLoading(true)

    try {
      const params: any = {
        page: currentPage,
        pageSize: 10,
      }
      if (keyword) {
        params.skill = keyword
      }
      if (rateIndex > 0) {
        const rateRanges = [[0, 100], [100, 200], [200, 300], [300, 99999]]
        const [min, max] = rateRanges[rateIndex - 1]
        params.hourlyRateMin = min
        params.hourlyRateMax = max
      }
      if (availabilityIndex > 0) {
        const availMap = ['available', 'busy']
        params.availability = availMap[availabilityIndex - 1]
      }

      const res = await api.developer.getList(params)
      if (res.code === 0) {
        const list = res.data?.list || res.data || []
        const total = res.meta?.total || 0
        setDevelopers(prev => reset ? list : [...prev, ...list])
        setHasMore(currentPage * 10 < total)
        setPage(currentPage + 1)
      }
    } catch (err) {
      // handle error
    } finally {
      setLoading(false)
    }
  }, [page, keyword, rateIndex, availabilityIndex, loading])

  usePullDownRefresh(() => {
    loadDevelopers(true)
    Taro.stopPullDownRefresh()
  })

  useReachBottom(() => {
    if (hasMore && !loading) {
      loadDevelopers()
    }
  })

  function handleSearch() {
    loadDevelopers(true)
  }

  return (
    <View className="dev-list-page">
      {/* Search Bar */}
      <View className="dev-list-page__search">
        <Input
          className="dev-list-page__search-input"
          placeholder="按技能搜索开发者"
          placeholderClass="dev-list-page__search-placeholder"
          value={keyword}
          onInput={(e) => setKeyword(e.detail.value)}
          onConfirm={handleSearch}
          confirmType="search"
        />
      </View>

      {/* Filters */}
      <View className="dev-list-page__filters">
        <Picker mode="selector" range={RATE_OPTIONS} value={rateIndex} onChange={(e) => setRateIndex(Number(e.detail.value))}>
          <View className="dev-list-page__filter-item">
            <Text className="dev-list-page__filter-text">{RATE_OPTIONS[rateIndex]}</Text>
            <Text className="dev-list-page__filter-arrow">▼</Text>
          </View>
        </Picker>
        <Picker mode="selector" range={AVAILABILITY_OPTIONS} value={availabilityIndex} onChange={(e) => setAvailabilityIndex(Number(e.detail.value))}>
          <View className="dev-list-page__filter-item">
            <Text className="dev-list-page__filter-text">{AVAILABILITY_OPTIONS[availabilityIndex]}</Text>
            <Text className="dev-list-page__filter-arrow">▼</Text>
          </View>
        </Picker>
      </View>

      {/* Developer List */}
      <ScrollView scrollY className="dev-list-page__content">
        {developers.length > 0 ? (
          developers.map((dev) => (
            <DeveloperCard key={dev.id} developer={dev} />
          ))
        ) : !loading ? (
          <Empty type="search" description="未找到匹配的开发者" />
        ) : null}

        {loading && (
          <View className="dev-list-page__loading">
            <Text>加载中...</Text>
          </View>
        )}

        {!hasMore && developers.length > 0 && (
          <View className="dev-list-page__no-more">
            <Text>没有更多了</Text>
          </View>
        )}
      </ScrollView>
    </View>
  )
}

list.config = {
  navigationBarTitleText: '找开发者',
  enablePullDownRefresh: true,
}

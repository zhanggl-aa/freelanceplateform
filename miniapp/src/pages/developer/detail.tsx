import { useState, useEffect } from 'react'
import { View, Text, Image, ScrollView } from '@tarojs/components'
import Taro, { useRouter } from '@tarojs/taro'
import { api } from '@/services/api'
import { useUser } from '@/store/user'
import { Developer, PortfolioItem } from '@/components/types'
import Empty from '@/components/Empty'
import './detail.scss'

export default function DeveloperDetail() {
  const router = useRouter()
  const { isLoggedIn } = useUser()
  const [developer, setDeveloper] = useState<Developer | null>(null)
  const [portfolio, setPortfolio] = useState<PortfolioItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (router.params.id) {
      loadDeveloper(router.params.id)
    }
  }, [router.params.id])

  async function loadDeveloper(id: string) {
    setLoading(true)
    try {
      const [devRes, portRes] = await Promise.all([
        api.developer.getDetail(id),
        api.developer.getPortfolio(id),
      ])
      if (devRes.code === 0) {
        setDeveloper(devRes.data)
        Taro.setNavigationBarTitle({ title: devRes.data?.nickname || '开发者详情' })
      }
      if (portRes.code === 0) {
        setPortfolio(portRes.data?.list || portRes.data || [])
      }
    } catch (err) {
      // handle error
    } finally {
      setLoading(false)
    }
  }

  function handleChat() {
    if (!isLoggedIn) {
      Taro.navigateTo({ url: '/pages/auth/login' })
      return
    }
    Taro.navigateTo({ url: `/pages/chat/detail?userId=${developer!.userId}` })
  }

  const availabilityMap: Record<string, string> = {
    available: '空闲',
    busy: '忙碌',
    unavailable: '不可用',
  }

  if (loading) {
    return (
      <View className="dev-detail-page__loading">
        <Text>加载中...</Text>
      </View>
    )
  }

  if (!developer) {
    return <Empty type="nodata" description="开发者不存在" />
  }

  return (
    <View className="dev-detail-page">
      <ScrollView scrollY className="dev-detail-page__content">
        {/* Profile Header */}
        <View className="dev-detail-page__header">
          <Image
            className="dev-detail-page__avatar"
            src={developer.avatar || 'https://via.placeholder.com/120'}
            mode="aspectFill"
          />
          <View className="dev-detail-page__header-info">
            <Text className="dev-detail-page__name">{developer.nickname}</Text>
            <Text className="dev-detail-page__title">{developer.title}</Text>
            <View className="dev-detail-page__header-row">
              <Text className="dev-detail-page__availability">
                {availabilityMap[developer.availability] || '未知'}
              </Text>
              <Text className="dev-detail-page__rate">¥{developer.hourlyRate}/小时</Text>
            </View>
          </View>
        </View>

        {/* Stats */}
        <View className="dev-detail-page__stats">
          <View className="dev-detail-page__stat-item">
            <Text className="dev-detail-page__stat-value">{developer.rating?.toFixed(1) || '0.0'}</Text>
            <Text className="dev-detail-page__stat-label">评分</Text>
          </View>
          <View className="dev-detail-page__stat-item">
            <Text className="dev-detail-page__stat-value">{developer.completedProjects || 0}</Text>
            <Text className="dev-detail-page__stat-label">完成项目</Text>
          </View>
          <View className="dev-detail-page__stat-item">
            <Text className="dev-detail-page__stat-value">{developer.skills?.length || 0}</Text>
            <Text className="dev-detail-page__stat-label">技能</Text>
          </View>
        </View>

        {/* Bio */}
        <View className="dev-detail-page__section">
          <Text className="dev-detail-page__section-title">个人简介</Text>
          <Text className="dev-detail-page__bio">{developer.bio || '暂无简介'}</Text>
        </View>

        {/* Skills */}
        {developer.skills && developer.skills.length > 0 && (
          <View className="dev-detail-page__section">
            <Text className="dev-detail-page__section-title">技能标签</Text>
            <View className="dev-detail-page__skills">
              {developer.skills.map((skill, index) => (
                <Text key={index} className="dev-detail-page__skill-tag">{skill}</Text>
              ))}
            </View>
          </View>
        )}

        {/* Portfolio */}
        {portfolio.length > 0 && (
          <View className="dev-detail-page__section">
            <Text className="dev-detail-page__section-title">作品集</Text>
            {portfolio.map((item) => (
              <View key={item.id} className="dev-detail-page__portfolio-item">
                <Text className="dev-detail-page__portfolio-title">{item.title}</Text>
                <Text className="dev-detail-page__portfolio-desc">{item.description}</Text>
                {item.images && item.images.length > 0 && (
                  <View className="dev-detail-page__portfolio-images">
                    {item.images.slice(0, 3).map((img, idx) => (
                      <Image
                        key={idx}
                        className="dev-detail-page__portfolio-img"
                        src={img}
                        mode="aspectFill"
                        onClick={() => Taro.previewImage({ urls: item.images, current: img })}
                      />
                    ))}
                  </View>
                )}
              </View>
            ))}
          </View>
        )}
      </ScrollView>

      {/* Bottom Action */}
      <View className="dev-detail-page__bottom">
        <View className="dev-detail-page__bottom-btn" onClick={handleChat}>
          <Text className="dev-detail-page__bottom-btn-text">联系TA</Text>
        </View>
      </View>
    </View>
  )
}

detail.config = {
  navigationBarTitleText: '开发者详情',
}

import { View, Text, Image, Navigator } from '@tarojs/components'
import { Developer } from './types'
import './DeveloperCard.scss'

interface DeveloperCardProps {
  developer: Developer
  horizontal?: boolean
}

export default function DeveloperCard({ developer, horizontal = false }: DeveloperCardProps) {
  const availabilityMap: Record<string, string> = {
    available: '空闲',
    busy: '忙碌',
    unavailable: '不可用',
  }

  const availabilityColorMap: Record<string, string> = {
    available: '#67C23A',
    busy: '#E6A23C',
    unavailable: '#909399',
  }

  if (horizontal) {
    return (
      <Navigator url={`/pages/developer/detail?id=${developer.id}`} className="developer-card--horizontal">
        <Image
          className="developer-card--horizontal__avatar"
          src={developer.avatar || 'https://via.placeholder.com/80'}
          mode="aspectFill"
        />
        <Text className="developer-card--horizontal__name">{developer.nickname}</Text>
        <View
          className="developer-card--horizontal__availability"
          style={{ backgroundColor: availabilityColorMap[developer.availability] || '#909399' }}
        />
        <Text className="developer-card--horizontal__title">{developer.title}</Text>
        <Text className="developer-card--horizontal__rate">¥{developer.hourlyRate}/时</Text>
      </Navigator>
    )
  }

  return (
    <Navigator url={`/pages/developer/detail?id=${developer.id}`} className="developer-card">
      <View className="developer-card__header">
        <Image
          className="developer-card__avatar"
          src={developer.avatar || 'https://via.placeholder.com/80'}
          mode="aspectFill"
        />
        <View className="developer-card__info">
          <View className="developer-card__name-row">
            <Text className="developer-card__name">{developer.nickname}</Text>
            <View
              className="developer-card__availability"
              style={{ backgroundColor: availabilityColorMap[developer.availability] || '#909399' }}
            >
              <Text className="developer-card__availability-text">
                {availabilityMap[developer.availability] || '未知'}
              </Text>
            </View>
          </View>
          <Text className="developer-card__title">{developer.title}</Text>
        </View>
      </View>

      <View className="developer-card__skills">
        {developer.skills && developer.skills.slice(0, 5).map((skill, index) => (
          <Text key={index} className="developer-card__skill-tag">{skill}</Text>
        ))}
      </View>

      <View className="developer-card__footer">
        <View className="developer-card__rate">
          <Text className="developer-card__rate-value">¥{developer.hourlyRate}</Text>
          <Text className="developer-card__rate-unit">/小时</Text>
        </View>
        <View className="developer-card__stats">
          <Text className="developer-card__rating">★ {developer.rating?.toFixed(1) || '0.0'}</Text>
          <Text className="developer-card__completed">{developer.completedProjects || 0}个项目</Text>
        </View>
      </View>
    </Navigator>
  )
}

import { View, Text, Image } from '@tarojs/components'
import './Empty.scss'

interface EmptyProps {
  type?: 'default' | 'search' | 'network' | 'nodata' | 'login'
  title?: string
  description?: string
  image?: string
  actionText?: string
  onAction?: () => void
}

const defaultImages: Record<string, string> = {
  default: 'https://via.placeholder.com/200x200?text=Empty',
  search: 'https://via.placeholder.com/200x200?text=No+Results',
  network: 'https://via.placeholder.com/200x200?text=Network+Error',
  nodata: 'https://via.placeholder.com/200x200?text=No+Data',
  login: 'https://via.placeholder.com/200x200?text=Login',
}

const defaultTitles: Record<string, string> = {
  default: '暂无数据',
  search: '未找到相关结果',
  network: '网络异常',
  nodata: '暂无内容',
  login: '请先登录',
}

const defaultDescriptions: Record<string, string> = {
  default: '当前没有相关数据',
  search: '换个关键词试试吧',
  network: '请检查网络连接后重试',
  nodata: '这里暂时还没有内容',
  login: '登录后即可查看',
}

export default function Empty({
  type = 'default',
  title,
  description,
  image,
  actionText,
  onAction,
}: EmptyProps) {
  const displayTitle = title || defaultTitles[type]
  const displayDesc = description || defaultDescriptions[type]
  const displayImage = image || defaultImages[type]

  return (
    <View className="empty">
      <Image className="empty__image" src={displayImage} mode="aspectFit" />
      <Text className="empty__title">{displayTitle}</Text>
      <Text className="empty__description">{displayDesc}</Text>
      {actionText && (
        <View className="empty__action" onClick={onAction}>
          <Text className="empty__action-text">{actionText}</Text>
        </View>
      )}
    </View>
  )
}

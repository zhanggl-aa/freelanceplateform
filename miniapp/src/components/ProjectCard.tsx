import { View, Text, Navigator } from '@tarojs/components'
import { Project } from './types'
import './ProjectCard.scss'

interface ProjectCardProps {
  project: Project
}

export default function ProjectCard({ project }: ProjectCardProps) {
  const statusMap: Record<string, string> = {
    open: '竞标中',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
  }

  const statusColorMap: Record<string, string> = {
    open: '#409EFF',
    in_progress: '#67C23A',
    completed: '#909399',
    cancelled: '#F56C6C',
  }

  return (
    <Navigator url={`/pages/project/detail?id=${project.id}`} className="project-card">
      <View className="project-card__header">
        <Text className="project-card__title">{project.title}</Text>
        <View
          className="project-card__status"
          style={{ color: statusColorMap[project.status] || '#909399' }}
        >
          {statusMap[project.status] || project.status}
        </View>
      </View>

      <View className="project-card__category">
        <Text className="project-card__category-tag">{project.categoryName}</Text>
      </View>

      <Text className="project-card__desc">{project.description}</Text>

      {project.techStack && project.techStack.length > 0 && (
        <View className="project-card__tech">
          {project.techStack.slice(0, 4).map((tech, index) => (
            <Text key={index} className="project-card__tech-tag">{tech}</Text>
          ))}
          {project.techStack.length > 4 && (
            <Text className="project-card__tech-tag project-card__tech-tag--more">
              +{project.techStack.length - 4}
            </Text>
          )}
        </View>
      )}

      <View className="project-card__footer">
        <View className="project-card__budget">
          <Text className="project-card__budget-label">预算</Text>
          <Text className="project-card__budget-value">
            {project.budgetType === 'fixed' ? `¥${project.budgetMin}` : `¥${project.budgetMin}-${project.budgetMax}`}
          </Text>
        </View>
        <View className="project-card__info">
          <Text className="project-card__bids">{project.bidCount || 0}人竞标</Text>
          <Text className="project-card__deadline">截止 {project.deadline}</Text>
        </View>
      </View>
    </Navigator>
  )
}

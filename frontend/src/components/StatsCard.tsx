import { HeroIcon } from '../types/heroicons'

interface StatsCardProps {
  name: string
  value: number | string
  icon: HeroIcon
  color: 'blue' | 'green' | 'purple' | 'red' | 'yellow' | 'indigo'
  change?: {
    value: number
    type: 'increase' | 'decrease'
  }
}

const colorClasses = {
  blue: {
    icon: 'text-blue-600',
    bg: 'bg-blue-50',
  },
  green: {
    icon: 'text-green-600',
    bg: 'bg-green-50',
  },
  purple: {
    icon: 'text-purple-600',
    bg: 'bg-purple-50',
  },
  red: {
    icon: 'text-red-600',
    bg: 'bg-red-50',
  },
  yellow: {
    icon: 'text-yellow-600',
    bg: 'bg-yellow-50',
  },
  indigo: {
    icon: 'text-indigo-600',
    bg: 'bg-indigo-50',
  },
}

export default function StatsCard({ name, value, icon: Icon, color, change }: StatsCardProps) {
  const colors = colorClasses[color]

  return (
    <div className="bg-white overflow-hidden shadow rounded-lg">
      <div className="p-5">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <div className={`${colors.bg} rounded-md p-3`}>
              <Icon className={`h-6 w-6 ${colors.icon}`} />
            </div>
          </div>
          <div className="ml-5 w-0 flex-1">
            <dl>
              <dt className="text-sm font-medium text-gray-500 truncate">{name}</dt>
              <dd className="flex items-baseline">
                <div className="text-2xl font-semibold text-gray-900">{value}</div>
                {change && (
                  <div className={`ml-2 flex items-baseline text-sm font-semibold ${
                    change.type === 'increase' ? 'text-green-600' : 'text-red-600'
                  }`}>
                    {change.type === 'increase' ? '+' : '-'}{Math.abs(change.value)}%
                  </div>
                )}
              </dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  )
}
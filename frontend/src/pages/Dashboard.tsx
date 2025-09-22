import { useQuery } from '@apollo/client'
import { GET_FLEETS, GET_VEHICLES, GET_DRIVERS, GET_ALERTS } from '../graphql/queries'
import {
  TruckIcon,
  UserGroupIcon,
  ExclamationTriangleIcon,
  BuildingOfficeIcon,
} from '@heroicons/react/24/outline'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import StatsCard from '../components/StatsCard'

export default function Dashboard() {
  const { data: fleetsData, loading: fleetsLoading } = useQuery(GET_FLEETS)
  const { data: vehiclesData, loading: vehiclesLoading } = useQuery(GET_VEHICLES)
  const { data: driversData, loading: driversLoading } = useQuery(GET_DRIVERS)

  // For alerts, we'll need to get all fleets first, then query alerts for each
  // For now, let's just show a placeholder
  const activeAlerts = 0

  if (fleetsLoading || vehiclesLoading || driversLoading) {
    return <LoadingSpinner />
  }

  const fleets = fleetsData?.fleets || []
  const vehicles = vehiclesData?.vehicles || []
  const drivers = driversData?.drivers || []

  const stats = [
    {
      name: 'Total Fleets',
      value: fleets.length,
      icon: BuildingOfficeIcon,
      color: 'blue',
    },
    {
      name: 'Total Vehicles',
      value: vehicles.length,
      icon: TruckIcon,
      color: 'green',
    },
    {
      name: 'Total Drivers',
      value: drivers.length,
      icon: UserGroupIcon,
      color: 'purple',
    },
    {
      name: 'Active Alerts',
      value: activeAlerts,
      icon: ExclamationTriangleIcon,
      color: 'red',
    },
  ]

  // Calculate risk metrics
  const highRiskVehicles = vehicles.filter((v: any) => v.riskScore > 7).length
  const highRiskDrivers = drivers.filter((d: any) => d.riskScore > 7).length

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-1 text-sm text-gray-500">
          Overview of your fleet risk intelligence system
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <StatsCard
            key={stat.name}
            name={stat.name}
            value={stat.value}
            icon={stat.icon}
            color={stat.color}
          />
        ))}
      </div>

      {/* Risk Overview */}
      <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <ExclamationTriangleIcon className="h-6 w-6 text-red-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    High Risk Vehicles
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {highRiskVehicles} of {vehicles.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div className="bg-gray-50 px-5 py-3">
            <div className="text-sm">
              <span className="font-medium text-red-600">
                {vehicles.length > 0 ? Math.round((highRiskVehicles / vehicles.length) * 100) : 0}%
              </span>
              <span className="text-gray-500"> of fleet at high risk</span>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <UserGroupIcon className="h-6 w-6 text-orange-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    High Risk Drivers
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {highRiskDrivers} of {drivers.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div className="bg-gray-50 px-5 py-3">
            <div className="text-sm">
              <span className="font-medium text-orange-600">
                {drivers.length > 0 ? Math.round((highRiskDrivers / drivers.length) * 100) : 0}%
              </span>
              <span className="text-gray-500"> of drivers need attention</span>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">
            Recent Activity
          </h3>
          <div className="mt-5">
            <div className="text-sm text-gray-500">
              Real-time activity feed will be displayed here once WebSocket integration is complete.
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
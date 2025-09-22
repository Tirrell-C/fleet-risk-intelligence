import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@apollo/client'
import { GET_FLEET } from '../graphql/queries'
import {
  TruckIcon,
  UserGroupIcon,
  ArrowLeftIcon,
} from '@heroicons/react/24/outline'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import StatsCard from '../components/StatsCard'
import { Fleet } from '../types/graphql'

export default function FleetDetail() {
  const { id } = useParams<{ id: string }>()
  const { data, loading, error } = useQuery(GET_FLEET, {
    variables: { id },
    skip: !id,
  })

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage message={error.message} />
  if (!data?.fleet) return <ErrorMessage message="Fleet not found" />

  const fleet: Fleet = data.fleet

  const stats = [
    {
      name: 'Total Vehicles',
      value: fleet.vehicles?.length || 0,
      icon: TruckIcon,
      color: 'blue' as const,
    },
    {
      name: 'Total Drivers',
      value: fleet.drivers?.length || 0,
      icon: UserGroupIcon,
      color: 'green' as const,
    },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center space-x-4">
        <Link
          to="/fleets"
          className="inline-flex items-center text-sm font-medium text-gray-500 hover:text-gray-700"
        >
          <ArrowLeftIcon className="mr-2 h-4 w-4" />
          Back to Fleets
        </Link>
      </div>

      <div>
        <h1 className="text-2xl font-bold text-gray-900">{fleet.name}</h1>
        <p className="mt-1 text-sm text-gray-500">{fleet.companyName}</p>
        <p className="text-sm text-gray-500">{fleet.contactEmail}</p>
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

      {/* Vehicles Section */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg leading-6 font-medium text-gray-900">Vehicles</h3>
            <Link
              to={`/vehicles?fleet=${fleet.id}`}
              className="text-sm font-medium text-primary-600 hover:text-primary-500"
            >
              View all
            </Link>
          </div>
          {fleet.vehicles && fleet.vehicles.length > 0 ? (
            <div className="overflow-hidden">
              <ul className="divide-y divide-gray-200">
                {fleet.vehicles.slice(0, 5).map((vehicle) => (
                  <li key={vehicle.id} className="py-4">
                    <div className="flex items-center space-x-4">
                      <div className="flex-shrink-0">
                        <TruckIcon className="h-6 w-6 text-gray-400" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-gray-900 truncate">
                          {vehicle.make} {vehicle.model} ({vehicle.year})
                        </p>
                        <p className="text-sm text-gray-500">{vehicle.licensePlate}</p>
                      </div>
                      <div className="flex-shrink-0">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          vehicle.status === 'ACTIVE' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                        }`}>
                          {vehicle.status}
                        </span>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <p className="text-sm text-gray-500">No vehicles assigned to this fleet.</p>
          )}
        </div>
      </div>

      {/* Drivers Section */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg leading-6 font-medium text-gray-900">Drivers</h3>
            <Link
              to={`/drivers?fleet=${fleet.id}`}
              className="text-sm font-medium text-primary-600 hover:text-primary-500"
            >
              View all
            </Link>
          </div>
          {fleet.drivers && fleet.drivers.length > 0 ? (
            <div className="overflow-hidden">
              <ul className="divide-y divide-gray-200">
                {fleet.drivers.slice(0, 5).map((driver) => (
                  <li key={driver.id} className="py-4">
                    <div className="flex items-center space-x-4">
                      <div className="flex-shrink-0">
                        <UserGroupIcon className="h-6 w-6 text-gray-400" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-gray-900 truncate">
                          {driver.firstName} {driver.lastName}
                        </p>
                        <p className="text-sm text-gray-500">{driver.email}</p>
                      </div>
                      <div className="flex-shrink-0">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          driver.status === 'ACTIVE' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                        }`}>
                          {driver.status}
                        </span>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <p className="text-sm text-gray-500">No drivers assigned to this fleet.</p>
          )}
        </div>
      </div>
    </div>
  )
}
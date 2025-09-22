import { useQuery } from '@apollo/client'
import { Link, useSearchParams } from 'react-router-dom'
import { GET_DRIVERS } from '../graphql/queries'
import { UserGroupIcon } from '@heroicons/react/24/outline'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import { Driver } from '../types/graphql'

export default function Drivers() {
  const [searchParams] = useSearchParams()
  const fleetId = searchParams.get('fleet')

  const { data, loading, error } = useQuery(GET_DRIVERS, {
    variables: fleetId ? { fleetId } : {},
  })

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage message={error.message} />

  const drivers: Driver[] = data?.drivers || []

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Drivers</h1>
        <p className="mt-1 text-sm text-gray-500">
          Monitor and manage your fleet drivers.
        </p>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {drivers.map((driver) => (
            <li key={driver.id}>
              <Link
                to={`/drivers/${driver.id}`}
                className="block hover:bg-gray-50 px-4 py-4 sm:px-6"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <UserGroupIcon className="h-6 w-6 text-gray-400" />
                    </div>
                    <div className="ml-4">
                      <div className="text-sm font-medium text-gray-900">
                        {driver.firstName} {driver.lastName}
                      </div>
                      <div className="text-sm text-gray-500">
                        {driver.email} • {driver.employeeId}
                      </div>
                      {driver.fleet && (
                        <div className="text-sm text-gray-500">
                          Fleet: {driver.fleet.name}
                        </div>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center space-x-4">
                    <div className="text-right">
                      <div className="text-sm font-medium text-gray-900">
                        Risk Score: {driver.riskScore}/10
                      </div>
                      <div className={`text-sm ${
                        driver.riskScore > 7 ? 'text-red-600' :
                        driver.riskScore > 4 ? 'text-yellow-600' : 'text-green-600'
                      }`}>
                        {driver.riskScore > 7 ? 'High Risk' :
                         driver.riskScore > 4 ? 'Medium Risk' : 'Low Risk'}
                      </div>
                    </div>
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      driver.status === 'ACTIVE' ? 'bg-green-100 text-green-800' :
                      driver.status === 'ON_DUTY' ? 'bg-blue-100 text-blue-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {driver.status.replace('_', ' ')}
                    </span>
                  </div>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      </div>

      {drivers.length === 0 && (
        <div className="text-center">
          <UserGroupIcon className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">No drivers</h3>
          <p className="mt-1 text-sm text-gray-500">
            No drivers found{fleetId ? ' for this fleet' : ''}.
          </p>
        </div>
      )}
    </div>
  )
}
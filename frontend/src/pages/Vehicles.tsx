import { useQuery } from '@apollo/client'
import { Link, useSearchParams } from 'react-router-dom'
import { GET_VEHICLES } from '../graphql/queries'
import { TruckIcon } from '@heroicons/react/24/outline'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import { Vehicle } from '../types/graphql'

export default function Vehicles() {
  const [searchParams] = useSearchParams()
  const fleetId = searchParams.get('fleet')

  const { data, loading, error } = useQuery(GET_VEHICLES, {
    variables: fleetId ? { fleetId } : {},
  })

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage message={error.message} />

  const vehicles: Vehicle[] = data?.vehicles || []

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Vehicles</h1>
        <p className="mt-1 text-sm text-gray-500">
          Monitor and manage your fleet vehicles.
        </p>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {vehicles.map((vehicle) => (
            <li key={vehicle.id}>
              <Link
                to={`/vehicles/${vehicle.id}`}
                className="block hover:bg-gray-50 px-4 py-4 sm:px-6"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <TruckIcon className="h-6 w-6 text-gray-400" />
                    </div>
                    <div className="ml-4">
                      <div className="text-sm font-medium text-gray-900">
                        {vehicle.make} {vehicle.model} ({vehicle.year})
                      </div>
                      <div className="text-sm text-gray-500">
                        {vehicle.licensePlate} • VIN: {vehicle.vin}
                      </div>
                      {vehicle.fleet && (
                        <div className="text-sm text-gray-500">
                          Fleet: {vehicle.fleet.name}
                        </div>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center space-x-4">
                    <div className="text-right">
                      <div className="text-sm font-medium text-gray-900">
                        Risk Score: {vehicle.riskScore}/10
                      </div>
                      <div className={`text-sm ${
                        vehicle.riskScore > 7 ? 'text-red-600' :
                        vehicle.riskScore > 4 ? 'text-yellow-600' : 'text-green-600'
                      }`}>
                        {vehicle.riskScore > 7 ? 'High Risk' :
                         vehicle.riskScore > 4 ? 'Medium Risk' : 'Low Risk'}
                      </div>
                    </div>
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      vehicle.status === 'ACTIVE' ? 'bg-green-100 text-green-800' :
                      vehicle.status === 'MAINTENANCE' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {vehicle.status}
                    </span>
                  </div>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      </div>

      {vehicles.length === 0 && (
        <div className="text-center">
          <TruckIcon className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">No vehicles</h3>
          <p className="mt-1 text-sm text-gray-500">
            No vehicles found{fleetId ? ' for this fleet' : ''}.
          </p>
        </div>
      )}
    </div>
  )
}
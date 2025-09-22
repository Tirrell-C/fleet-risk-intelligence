import { useQuery } from '@apollo/client'
import { GET_RISK_EVENTS } from '../graphql/queries'
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import { RiskEvent } from '../types/graphql'

export default function RiskEvents() {
  const { data, loading, error } = useQuery(GET_RISK_EVENTS, {
    variables: { limit: 50 }
  })

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage message={error.message} />

  const riskEvents: RiskEvent[] = data?.riskEvents || []

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'CRITICAL': return 'bg-red-100 text-red-800'
      case 'HIGH': return 'bg-orange-100 text-orange-800'
      case 'MEDIUM': return 'bg-yellow-100 text-yellow-800'
      case 'LOW': return 'bg-green-100 text-green-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Risk Events</h1>
        <p className="mt-1 text-sm text-gray-500">
          Monitor driving behavior and safety incidents across your fleet.
        </p>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {riskEvents.map((event) => (
            <li key={event.id} className="px-4 py-4 sm:px-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <ExclamationTriangleIcon className="h-6 w-6 text-red-400" />
                  </div>
                  <div className="ml-4">
                    <div className="text-sm font-medium text-gray-900">
                      {event.eventType.replace('_', ' ')}
                    </div>
                    <div className="text-sm text-gray-500">
                      {event.description}
                    </div>
                    <div className="text-sm text-gray-500">
                      {event.vehicle && `${event.vehicle.make} ${event.vehicle.model} (${event.vehicle.licensePlate})`}
                      {event.driver && ` • ${event.driver.firstName} ${event.driver.lastName}`}
                    </div>
                    <div className="text-xs text-gray-400">
                      {new Date(event.timestamp).toLocaleString()}
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  <div className="text-right">
                    <div className="text-sm font-medium text-gray-900">
                      Risk Score: {event.riskScore}/10
                    </div>
                  </div>
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getSeverityColor(event.severity)}`}>
                    {event.severity}
                  </span>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>

      {riskEvents.length === 0 && (
        <div className="text-center">
          <ExclamationTriangleIcon className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">No risk events</h3>
          <p className="mt-1 text-sm text-gray-500">
            No risk events have been recorded yet.
          </p>
        </div>
      )}
    </div>
  )
}
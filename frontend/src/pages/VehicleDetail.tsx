import { useParams } from 'react-router-dom'

export default function VehicleDetail() {
  const { id } = useParams<{ id: string }>()

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Vehicle Details</h1>
        <p className="mt-1 text-sm text-gray-500">Vehicle ID: {id}</p>
      </div>

      <div className="bg-white shadow rounded-lg p-6">
        <p className="text-gray-500">Vehicle detail page will be implemented here.</p>
      </div>
    </div>
  )
}
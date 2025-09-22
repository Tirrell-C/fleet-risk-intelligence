import { useState } from 'react'
import { useQuery, useMutation } from '@apollo/client'
import { Link } from 'react-router-dom'
import { GET_FLEETS } from '../graphql/queries'
import { CREATE_FLEET } from '../graphql/mutations'
import { PlusIcon, BuildingOfficeIcon } from '@heroicons/react/24/outline'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import { Fleet, CreateFleetInput } from '../types/graphql'

export default function Fleets() {
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [formData, setFormData] = useState<CreateFleetInput>({
    name: '',
    companyName: '',
    contactEmail: '',
  })

  const { data, loading, error, refetch } = useQuery(GET_FLEETS)
  const [createFleet, { loading: createLoading }] = useMutation(CREATE_FLEET, {
    onCompleted: () => {
      setShowCreateForm(false)
      setFormData({ name: '', companyName: '', contactEmail: '' })
      refetch()
    },
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createFleet({
        variables: { input: formData }
      })
    } catch (error) {
      console.error('Error creating fleet:', error)
    }
  }

  if (loading) return <LoadingSpinner />
  if (error) return <ErrorMessage message={error.message} onRetry={() => refetch()} />

  const fleets: Fleet[] = data?.fleets || []

  return (
    <div className="space-y-6">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-2xl font-bold text-gray-900">Fleets</h1>
          <p className="mt-1 text-sm text-gray-500">
            Manage your fleet organizations and their vehicles and drivers.
          </p>
        </div>
        <div className="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
          <button
            type="button"
            onClick={() => setShowCreateForm(true)}
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 sm:w-auto"
          >
            <PlusIcon className="-ml-1 mr-2 h-4 w-4" />
            Add Fleet
          </button>
        </div>
      </div>

      {/* Create Fleet Form */}
      {showCreateForm && (
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">Create New Fleet</h3>
            <form onSubmit={handleSubmit} className="mt-5 space-y-4">
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                  Fleet Name
                </label>
                <input
                  type="text"
                  id="name"
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                />
              </div>
              <div>
                <label htmlFor="companyName" className="block text-sm font-medium text-gray-700">
                  Company Name
                </label>
                <input
                  type="text"
                  id="companyName"
                  required
                  value={formData.companyName}
                  onChange={(e) => setFormData({ ...formData, companyName: e.target.value })}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                />
              </div>
              <div>
                <label htmlFor="contactEmail" className="block text-sm font-medium text-gray-700">
                  Contact Email
                </label>
                <input
                  type="email"
                  id="contactEmail"
                  required
                  value={formData.contactEmail}
                  onChange={(e) => setFormData({ ...formData, contactEmail: e.target.value })}
                  className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                />
              </div>
              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={createLoading}
                  className="inline-flex justify-center rounded-md border border-transparent bg-primary-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50"
                >
                  {createLoading ? <LoadingSpinner size="sm" /> : 'Create Fleet'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Fleets Grid */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {fleets.map((fleet) => (
          <Link
            key={fleet.id}
            to={`/fleets/${fleet.id}`}
            className="relative rounded-lg border border-gray-300 bg-white px-6 py-5 shadow-sm flex items-center space-x-3 hover:border-gray-400 focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-primary-500"
          >
            <div className="flex-shrink-0">
              <BuildingOfficeIcon className="h-10 w-10 text-gray-400" />
            </div>
            <div className="flex-1 min-w-0">
              <span className="absolute inset-0" aria-hidden="true" />
              <p className="text-sm font-medium text-gray-900 truncate">{fleet.name}</p>
              <p className="text-sm text-gray-500 truncate">{fleet.companyName}</p>
              <p className="text-xs text-gray-400">{fleet.contactEmail}</p>
            </div>
          </Link>
        ))}
      </div>

      {fleets.length === 0 && (
        <div className="text-center">
          <BuildingOfficeIcon className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">No fleets</h3>
          <p className="mt-1 text-sm text-gray-500">Get started by creating a new fleet.</p>
          <div className="mt-6">
            <button
              type="button"
              onClick={() => setShowCreateForm(true)}
              className="inline-flex items-center rounded-md border border-transparent bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            >
              <PlusIcon className="-ml-1 mr-2 h-4 w-4" />
              New Fleet
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
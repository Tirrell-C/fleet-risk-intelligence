import { gql } from '@apollo/client'

export const VEHICLE_UPDATES_SUBSCRIPTION = gql`
  subscription VehicleUpdates($vehicleId: ID!) {
    vehicleUpdates(vehicleId: $vehicleId) {
      vehicle {
        id
        make
        model
        licensePlate
      }
      location {
        latitude
        longitude
        address
      }
      speed
      heading
      fuelLevel
      engineStatus
      lastUpdate
    }
  }
`

export const RISK_EVENT_NOTIFICATIONS_SUBSCRIPTION = gql`
  subscription RiskEventNotifications($fleetId: ID!) {
    riskEventNotifications(fleetId: $fleetId) {
      id
      vehicleId
      driverId
      eventType
      severity
      description
      latitude
      longitude
      riskScore
      status
      vehicle {
        id
        make
        model
        licensePlate
      }
      driver {
        id
        firstName
        lastName
      }
      timestamp
      createdAt
    }
  }
`

export const ALERT_NOTIFICATIONS_SUBSCRIPTION = gql`
  subscription AlertNotifications($fleetId: ID!) {
    alertNotifications(fleetId: $fleetId) {
      id
      fleetId
      vehicleId
      driverId
      riskEventId
      type
      priority
      title
      message
      status
      vehicle {
        id
        make
        model
        licensePlate
      }
      driver {
        id
        firstName
        lastName
      }
      createdAt
    }
  }
`
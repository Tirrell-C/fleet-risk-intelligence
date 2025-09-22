import { gql } from '@apollo/client'

export const GET_FLEETS = gql`
  query GetFleets {
    fleets {
      id
      name
      companyName
      contactEmail
      status
      createdAt
      updatedAt
    }
  }
`

export const GET_FLEET = gql`
  query GetFleet($id: ID!) {
    fleet(id: $id) {
      id
      name
      companyName
      contactEmail
      status
      vehicles {
        id
        make
        model
        year
        licensePlate
        status
        riskScore
      }
      drivers {
        id
        firstName
        lastName
        email
        status
        riskScore
      }
      createdAt
      updatedAt
    }
  }
`

export const GET_VEHICLES = gql`
  query GetVehicles($fleetId: ID) {
    vehicles(fleetId: $fleetId) {
      id
      fleetId
      driverId
      make
      model
      year
      vin
      licensePlate
      status
      riskScore
      fleet {
        id
        name
      }
      driver {
        id
        firstName
        lastName
      }
      createdAt
      updatedAt
    }
  }
`

export const GET_VEHICLE = gql`
  query GetVehicle($id: ID!) {
    vehicle(id: $id) {
      id
      fleetId
      driverId
      make
      model
      year
      vin
      licensePlate
      status
      riskScore
      fleet {
        id
        name
        companyName
      }
      driver {
        id
        firstName
        lastName
        email
        phone
      }
      currentLocation {
        latitude
        longitude
        address
      }
      lastTelemetry {
        id
        speed
        timestamp
      }
      createdAt
      updatedAt
    }
  }
`

export const GET_DRIVERS = gql`
  query GetDrivers($fleetId: ID) {
    drivers(fleetId: $fleetId) {
      id
      fleetId
      employeeId
      firstName
      lastName
      email
      phone
      licenseNumber
      status
      riskScore
      fleet {
        id
        name
      }
      currentVehicle {
        id
        make
        model
        licensePlate
      }
      createdAt
      updatedAt
    }
  }
`

export const GET_DRIVER = gql`
  query GetDriver($id: ID!) {
    driver(id: $id) {
      id
      fleetId
      employeeId
      firstName
      lastName
      email
      phone
      licenseNumber
      status
      riskScore
      fleet {
        id
        name
        companyName
      }
      currentVehicle {
        id
        make
        model
        licensePlate
        vin
      }
      driverScore {
        id
        overallScore
        safetyScore
        efficiencyScore
        totalMiles
        totalTrips
        riskEvents
        lastUpdated
      }
      createdAt
      updatedAt
    }
  }
`

export const GET_RISK_EVENTS = gql`
  query GetRiskEvents($vehicleId: ID, $driverId: ID, $limit: Int) {
    riskEvents(vehicleId: $vehicleId, driverId: $driverId, limit: $limit) {
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

export const GET_ALERTS = gql`
  query GetAlerts($fleetId: ID!, $status: AlertStatus) {
    alerts(fleetId: $fleetId, status: $status) {
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
      fleet {
        id
        name
      }
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
      updatedAt
    }
  }
`

export const GET_DRIVER_SCORES = gql`
  query GetDriverScores($fleetId: ID!) {
    driverScores(fleetId: $fleetId) {
      id
      driverId
      overallScore
      safetyScore
      efficiencyScore
      totalMiles
      totalTrips
      riskEvents
      driver {
        id
        firstName
        lastName
        employeeId
      }
      lastUpdated
    }
  }
`

export const GET_LIVE_VEHICLE_DATA = gql`
  query GetLiveVehicleData($vehicleId: ID!) {
    liveVehicleData(vehicleId: $vehicleId) {
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
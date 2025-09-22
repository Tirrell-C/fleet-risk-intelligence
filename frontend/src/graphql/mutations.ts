import { gql } from '@apollo/client'

export const CREATE_FLEET = gql`
  mutation CreateFleet($input: CreateFleetInput!) {
    createFleet(input: $input) {
      id
      name
      companyName
      contactEmail
      status
      createdAt
    }
  }
`

export const UPDATE_FLEET = gql`
  mutation UpdateFleet($id: ID!, $input: UpdateFleetInput!) {
    updateFleet(id: $id, input: $input) {
      id
      name
      companyName
      contactEmail
      status
      updatedAt
    }
  }
`

export const CREATE_VEHICLE = gql`
  mutation CreateVehicle($input: CreateVehicleInput!) {
    createVehicle(input: $input) {
      id
      fleetId
      make
      model
      year
      vin
      licensePlate
      status
      fleet {
        id
        name
      }
      createdAt
    }
  }
`

export const UPDATE_VEHICLE = gql`
  mutation UpdateVehicle($id: ID!, $input: UpdateVehicleInput!) {
    updateVehicle(id: $id, input: $input) {
      id
      make
      model
      year
      vin
      licensePlate
      status
      updatedAt
    }
  }
`

export const ASSIGN_DRIVER = gql`
  mutation AssignDriver($vehicleId: ID!, $driverId: ID!) {
    assignDriver(vehicleId: $vehicleId, driverId: $driverId) {
      id
      driverId
      driver {
        id
        firstName
        lastName
      }
      updatedAt
    }
  }
`

export const CREATE_DRIVER = gql`
  mutation CreateDriver($input: CreateDriverInput!) {
    createDriver(input: $input) {
      id
      fleetId
      employeeId
      firstName
      lastName
      email
      phone
      licenseNumber
      status
      fleet {
        id
        name
      }
      createdAt
    }
  }
`

export const UPDATE_DRIVER = gql`
  mutation UpdateDriver($id: ID!, $input: UpdateDriverInput!) {
    updateDriver(id: $id, input: $input) {
      id
      employeeId
      firstName
      lastName
      email
      phone
      licenseNumber
      status
      updatedAt
    }
  }
`

export const ACKNOWLEDGE_ALERT = gql`
  mutation AcknowledgeAlert($id: ID!) {
    acknowledgeAlert(id: $id) {
      id
      status
      updatedAt
    }
  }
`

export const DISMISS_ALERT = gql`
  mutation DismissAlert($id: ID!) {
    dismissAlert(id: $id) {
      id
      status
      updatedAt
    }
  }
`
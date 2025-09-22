// Generated GraphQL types - would normally be auto-generated from schema
export interface Fleet {
  id: string
  name: string
  companyName: string
  contactEmail: string
  status: string
  vehicles: Vehicle[]
  drivers: Driver[]
  createdAt: string
  updatedAt: string
}

export interface Vehicle {
  id: string
  fleetId: string
  driverId?: string
  make: string
  model: string
  year: number
  vin: string
  licensePlate: string
  status: VehicleStatus
  riskScore: number
  fleet: Fleet
  driver?: Driver
  currentLocation?: Location
  lastTelemetry?: TelemetryEvent
  createdAt: string
  updatedAt: string
}

export interface Driver {
  id: string
  fleetId: string
  employeeId: string
  firstName: string
  lastName: string
  email: string
  phone: string
  licenseNumber: string
  status: DriverStatus
  riskScore: number
  fleet: Fleet
  currentVehicle?: Vehicle
  driverScore?: DriverScore
  createdAt: string
  updatedAt: string
}

export interface RiskEvent {
  id: string
  vehicleId: string
  driverId?: string
  eventType: RiskEventType
  severity: RiskSeverity
  description: string
  latitude: number
  longitude: number
  riskScore: number
  status: RiskEventStatus
  data: any
  vehicle: Vehicle
  driver?: Driver
  timestamp: string
  createdAt: string
  updatedAt: string
}

export interface Alert {
  id: string
  fleetId: string
  vehicleId?: string
  driverId?: string
  riskEventId?: string
  type: AlertType
  priority: AlertPriority
  title: string
  message: string
  status: AlertStatus
  fleet: Fleet
  vehicle?: Vehicle
  driver?: Driver
  riskEvent?: RiskEvent
  createdAt: string
  updatedAt: string
}

export interface TelemetryEvent {
  id: string
  vehicleId: string
  eventType: string
  latitude: number
  longitude: number
  speed: number
  acceleration: number
  data: any
  vehicle: Vehicle
  timestamp: string
  processedAt?: string
  createdAt: string
}

export interface DriverScore {
  id: string
  driverId: string
  overallScore: number
  safetyScore: number
  efficiencyScore: number
  totalMiles: number
  totalTrips: number
  riskEvents: number
  driver: Driver
  lastUpdated: string
  createdAt: string
  updatedAt: string
}

export interface Location {
  latitude: number
  longitude: number
  address?: string
}

export interface VehicleData {
  vehicle: Vehicle
  location: Location
  speed: number
  heading: number
  fuelLevel: number
  engineStatus: string
  lastUpdate: string
}

// Enums
export enum VehicleStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  MAINTENANCE = 'MAINTENANCE',
  OUT_OF_SERVICE = 'OUT_OF_SERVICE'
}

export enum DriverStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  ON_DUTY = 'ON_DUTY',
  OFF_DUTY = 'OFF_DUTY',
  SUSPENDED = 'SUSPENDED'
}

export enum RiskEventType {
  HARD_BRAKING = 'HARD_BRAKING',
  RAPID_ACCELERATION = 'RAPID_ACCELERATION',
  SHARP_TURN = 'SHARP_TURN',
  SPEEDING = 'SPEEDING',
  DISTRACTED_DRIVING = 'DISTRACTED_DRIVING',
  COLLISION = 'COLLISION',
  HARSH_CORNERING = 'HARSH_CORNERING'
}

export enum RiskSeverity {
  LOW = 'LOW',
  MEDIUM = 'MEDIUM',
  HIGH = 'HIGH',
  CRITICAL = 'CRITICAL'
}

export enum RiskEventStatus {
  ACTIVE = 'ACTIVE',
  RESOLVED = 'RESOLVED',
  IGNORED = 'IGNORED'
}

export enum AlertType {
  RISK_EVENT = 'RISK_EVENT',
  VEHICLE_MAINTENANCE = 'VEHICLE_MAINTENANCE',
  DRIVER_BEHAVIOR = 'DRIVER_BEHAVIOR',
  SYSTEM = 'SYSTEM'
}

export enum AlertPriority {
  LOW = 'LOW',
  MEDIUM = 'MEDIUM',
  HIGH = 'HIGH',
  URGENT = 'URGENT'
}

export enum AlertStatus {
  ACTIVE = 'ACTIVE',
  ACKNOWLEDGED = 'ACKNOWLEDGED',
  RESOLVED = 'RESOLVED',
  DISMISSED = 'DISMISSED'
}

// Input types
export interface CreateFleetInput {
  name: string
  companyName: string
  contactEmail: string
}

export interface UpdateFleetInput {
  name?: string
  companyName?: string
  contactEmail?: string
  status?: string
}

export interface CreateVehicleInput {
  fleetId: string
  make: string
  model: string
  year: number
  vin: string
  licensePlate: string
}

export interface UpdateVehicleInput {
  make?: string
  model?: string
  year?: number
  vin?: string
  licensePlate?: string
  status?: VehicleStatus
}

export interface CreateDriverInput {
  fleetId: string
  employeeId: string
  firstName: string
  lastName: string
  email: string
  phone: string
  licenseNumber: string
}

export interface UpdateDriverInput {
  employeeId?: string
  firstName?: string
  lastName?: string
  email?: string
  phone?: string
  licenseNumber?: string
  status?: DriverStatus
}
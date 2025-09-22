import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Fleets from './pages/Fleets'
import FleetDetail from './pages/FleetDetail'
import Vehicles from './pages/Vehicles'
import VehicleDetail from './pages/VehicleDetail'
import Drivers from './pages/Drivers'
import DriverDetail from './pages/DriverDetail'
import RiskEvents from './pages/RiskEvents'
import Alerts from './pages/Alerts'

function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/fleets" element={<Fleets />} />
        <Route path="/fleets/:id" element={<FleetDetail />} />
        <Route path="/vehicles" element={<Vehicles />} />
        <Route path="/vehicles/:id" element={<VehicleDetail />} />
        <Route path="/drivers" element={<Drivers />} />
        <Route path="/drivers/:id" element={<DriverDetail />} />
        <Route path="/risk-events" element={<RiskEvents />} />
        <Route path="/alerts" element={<Alerts />} />
      </Routes>
    </Layout>
  )
}

export default App
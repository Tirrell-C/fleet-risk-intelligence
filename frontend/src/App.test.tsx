import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { MockedProvider } from '@apollo/client/testing'
import App from './App'

const AppWithProviders = () => (
  <MockedProvider mocks={[]} addTypename={false}>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </MockedProvider>
)

describe('App', () => {
  it('renders without crashing', () => {
    render(<AppWithProviders />)
    expect(screen.getByText('Fleet Risk Intelligence')).toBeInTheDocument()
  })
})
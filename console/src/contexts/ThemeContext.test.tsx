import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { ThemeProvider, useTheme } from './ThemeContext'
import { resolveTheme } from '../themeMode'

function Probe() {
  const { mode, resolved, setMode } = useTheme()
  return (
    <div>
      <span data-testid="mode">{mode}</span>
      <span data-testid="resolved">{resolved}</span>
      <button onClick={() => setMode('dark')}>dark</button>
      <button onClick={() => setMode('light')}>light</button>
    </div>
  )
}

describe('ThemeContext', () => {
  beforeEach(() => {
    localStorage.clear()
    delete document.documentElement.dataset.theme
  })

  it('answers light without a provider', () => {
    render(<Probe />)
    expect(screen.getByTestId('resolved').textContent).toBe('light')
  })

  it('stores the choice and stamps <html>', () => {
    render(
      <ThemeProvider>
        <Probe />
      </ThemeProvider>
    )
    act(() => screen.getByText('dark').click())
    expect(screen.getByTestId('mode').textContent).toBe('dark')
    expect(screen.getByTestId('resolved').textContent).toBe('dark')
    expect(document.documentElement.dataset.theme).toBe('dark')
    expect(localStorage.getItem('notifuse_theme')).toBe('dark')

    act(() => screen.getByText('light').click())
    expect(document.documentElement.dataset.theme).toBe('light')
  })

  it('reads a stored choice on mount', () => {
    localStorage.setItem('notifuse_theme', 'dark')
    render(
      <ThemeProvider>
        <Probe />
      </ThemeProvider>
    )
    expect(screen.getByTestId('resolved').textContent).toBe('dark')
  })

  it('resolves explicit modes without touching matchMedia', () => {
    expect(resolveTheme('light')).toBe('light')
    expect(resolveTheme('dark')).toBe('dark')
  })
})

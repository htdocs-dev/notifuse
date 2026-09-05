import { createContext, useCallback, useContext, useEffect, useMemo, useState, useSyncExternalStore, ReactNode } from 'react'
import {
  applyTheme,
  readStoredThemeMode,
  resolveTheme,
  storeThemeMode,
  subscribeSystemTheme,
  systemTheme,
  type ResolvedTheme,
  type ThemeMode
} from '../themeMode'

// Apply before the first paint so a dark user never sees a light flash.
applyTheme(resolveTheme(readStoredThemeMode()))

interface ThemeContextType {
  mode: ThemeMode
  resolved: ResolvedTheme
  setMode: (mode: ThemeMode) => void
}

const ThemeContext = createContext<ThemeContextType | null>(null)

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setModeState] = useState<ThemeMode>(readStoredThemeMode)
  const system = useSyncExternalStore(subscribeSystemTheme, systemTheme, () => 'light' as const)
  const resolved = resolveTheme(mode, system)

  useEffect(() => {
    applyTheme(resolved)
  }, [resolved])

  const setMode = useCallback((next: ThemeMode) => {
    setModeState(next)
    storeThemeMode(next)
  }, [])

  const value = useMemo(() => ({ mode, resolved, setMode }), [mode, resolved, setMode])

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
}

const fallback: ThemeContextType = { mode: 'light', resolved: 'light', setMode: () => {} }

/** Usable without a provider (tests, isolated renders): answers light. */
// eslint-disable-next-line react-refresh/only-export-components
export function useTheme(): ThemeContextType {
  return useContext(ThemeContext) ?? fallback
}

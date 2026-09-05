/**
 * Console colour theme, the part with no React in it.
 *
 * The choice is stored per browser (localStorage) and applied as
 * `data-theme="light" | "dark"` on <html>. Everything else keys off that
 * attribute: Ant Design switches algorithm in App.tsx, and index.css swaps the
 * Tailwind grey scale and the surface variables. "system" follows the OS.
 */
export type ThemeMode = 'light' | 'dark' | 'system'
export type ResolvedTheme = 'light' | 'dark'

export const THEME_STORAGE_KEY = 'notifuse_theme'
const DEFAULT_MODE: ThemeMode = 'system'
const DARK_QUERY = '(prefers-color-scheme: dark)'

export function readStoredThemeMode(): ThemeMode {
  try {
    const v = localStorage.getItem(THEME_STORAGE_KEY)
    if (v === 'light' || v === 'dark' || v === 'system') return v
  } catch {
    // storage unavailable (private mode, tests): fall through
  }
  return DEFAULT_MODE
}

export function storeThemeMode(mode: ThemeMode) {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, mode)
  } catch {
    // storage unavailable: the choice lasts for this page only
  }
}

function darkQuery(): MediaQueryList | null {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return null
  return window.matchMedia(DARK_QUERY)
}

export function systemTheme(): ResolvedTheme {
  return darkQuery()?.matches ? 'dark' : 'light'
}

/** Calls onChange whenever the OS preference flips; returns the unsubscribe. */
export function subscribeSystemTheme(onChange: () => void): () => void {
  const mq = darkQuery()
  if (!mq) return () => {}
  mq.addEventListener('change', onChange)
  return () => mq.removeEventListener('change', onChange)
}

export function resolveTheme(mode: ThemeMode, system: ResolvedTheme = systemTheme()): ResolvedTheme {
  return mode === 'system' ? system : mode
}

export function applyTheme(resolved: ResolvedTheme) {
  if (typeof document === 'undefined') return
  document.documentElement.dataset.theme = resolved
}

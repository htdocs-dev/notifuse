import { Dropdown, Button } from 'antd'
import type { MenuProps } from 'antd'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faSun, faMoon } from '@fortawesome/free-regular-svg-icons'
import { faDesktop } from '@fortawesome/free-solid-svg-icons'
import { useLingui } from '@lingui/react/macro'
import { useTheme } from '../contexts/ThemeContext'
import type { ThemeMode } from '../themeMode'

export function ThemeSwitcher() {
  const { t } = useLingui()
  const { mode, resolved, setMode } = useTheme()

  const items: MenuProps['items'] = (
    [
      { key: 'light', icon: faSun, label: t`Light` },
      { key: 'dark', icon: faMoon, label: t`Dark` },
      { key: 'system', icon: faDesktop, label: t`Match system` }
    ] as const
  ).map((item) => ({
    key: item.key,
    label: item.label,
    icon: <FontAwesomeIcon icon={item.icon} fixedWidth />,
    onClick: () => setMode(item.key as ThemeMode)
  }))

  return (
    <Dropdown trigger={['click']} menu={{ items, selectedKeys: [mode] }} placement="bottomRight">
      <Button
        color="default"
        variant="filled"
        aria-label={t`Colour theme`}
        icon={<FontAwesomeIcon icon={resolved === 'dark' ? faMoon : faSun} />}
      />
    </Dropdown>
  )
}

/**
 * 主题模式：日间 / 夜间 / 跟随系统（设备）
 * auto = 跟随操作系统（prefers-color-scheme），默认值
 */
export type ThemeMode = 'auto' | 'light' | 'dark'

/** 主题模式持久化键（供 index.html 首屏脚本直接读取） */
export const THEME_MODE_STORAGE_KEY = 'themeMode'

/** 旧版本只存了布尔值 isDark，保留读取以完成平滑迁移 */
export const LEGACY_IS_DARK_STORAGE_KEY = 'isDark'

export const THEME_MODE_OPTIONS: ThemeMode[] = ['auto', 'light', 'dark']

export const DEFAULT_THEME_MODE: ThemeMode = 'auto'

export const isThemeMode = (value: unknown): value is ThemeMode =>
  typeof value === 'string' && (THEME_MODE_OPTIONS as string[]).includes(value)

/** 系统当前是否处于深色模式（SSR / 老浏览器安全降级） */
export const getSystemDark = (): boolean => {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return false
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

/** 由模式推导最终是否应用深色 */
export const resolveIsDark = (mode: ThemeMode, systemDark: boolean): boolean =>
  mode === 'auto' ? systemDark : mode === 'dark'

/** 把 dark/light 类与 color-scheme 写到 <html>，原生控件、滚动条也会跟随 */
export const applyThemeClass = (isDark: boolean): void => {
  if (typeof document === 'undefined') return
  const root = document.documentElement
  root.classList.toggle('dark', isDark)
  root.classList.toggle('light', !isDark)
  root.style.colorScheme = isDark ? 'dark' : 'light'
}

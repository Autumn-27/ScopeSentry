import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useAppStore } from './app'
import { LEGACY_IS_DARK_STORAGE_KEY, THEME_MODE_STORAGE_KEY } from '@/constants/theme'

type MediaHandler = (event: { matches: boolean }) => void

const mediaHandlers: MediaHandler[] = []

/**
 * jsdom 没有实现 matchMedia，这里按真实浏览器的最小契约装一个桩：
 * 记录 change 回调，便于测试“系统深浅色变化”。
 */
const installMatchMedia = (options: { dark?: boolean; reduceMotion?: boolean } = {}) => {
  const { dark = false, reduceMotion = false } = options
  mediaHandlers.length = 0
  Object.defineProperty(window, 'matchMedia', {
    configurable: true,
    writable: true,
    value: (query: string) => ({
      matches: query.includes('prefers-reduced-motion') ? reduceMotion : dark,
      media: query,
      onchange: null,
      addEventListener: (_type: string, handler: MediaHandler) => {
        mediaHandlers.push(handler)
      },
      removeEventListener: () => {},
      addListener: (handler: MediaHandler) => {
        mediaHandlers.push(handler)
      },
      removeListener: () => {},
      dispatchEvent: () => false
    })
  })
}

/** 模拟操作系统切换深浅色 */
const fireSystemChange = (matches: boolean) => {
  for (const handler of [...mediaHandlers]) handler({ matches })
}

const hasClass = (className: string) => document.documentElement.classList.contains(className)

const storedMode = () => JSON.parse(window.localStorage.getItem(THEME_MODE_STORAGE_KEY) as string).value

describe('app store theme mode', () => {
  beforeEach(() => {
    window.localStorage.clear()
    document.documentElement.className = ''
    document.documentElement.style.colorScheme = ''
    installMatchMedia()
    setActivePinia(createPinia())
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('mirrors a theme change made in another tab and ignores corrupt payloads', () => {
    const store = useAppStore()
    store.initTheme()
    expect(store.themeMode).toBe('auto')

    // 其它标签页切换成夜间后，storage 事件把本标签页带过去
    window.dispatchEvent(
      new StorageEvent('storage', {
        key: THEME_MODE_STORAGE_KEY,
        newValue: JSON.stringify({ type: 'String', value: 'dark' })
      })
    )
    expect(store.themeMode).toBe('dark')
    expect(store.isDark).toBe(true)
    expect(hasClass('dark')).toBe(true)
    expect(storedMode()).toBe('dark')

    // 损坏的值不能让页面白屏或回退成浅色
    window.dispatchEvent(
      new StorageEvent('storage', { key: THEME_MODE_STORAGE_KEY, newValue: '{not-json' })
    )
    expect(store.themeMode).toBe('dark')
    expect(hasClass('dark')).toBe(true)
  })

  it('restores a persisted mode on boot instead of blindly following the system', () => {
    window.localStorage.setItem(THEME_MODE_STORAGE_KEY, JSON.stringify({ type: 'String', value: 'light' }))
    installMatchMedia({ dark: true }) // 系统是深色

    const store = useAppStore()
    store.initTheme()

    expect(store.systemDark).toBe(true)
    expect(store.themeMode).toBe('light')
    expect(store.isDark).toBe(false)
    expect(hasClass('light')).toBe(true)
  })

  it('never turns a legacy isDark flag into a fixed mode', () => {
    // 旧版本每次启动都按系统覆盖写 isDark，它不能代表用户的选择
    window.localStorage.setItem(LEGACY_IS_DARK_STORAGE_KEY, JSON.stringify({ type: 'Boolean', value: true }))
    installMatchMedia({ dark: false })

    const store = useAppStore()
    store.initTheme()

    expect(store.themeMode).toBe('auto')
    expect(store.isDark).toBe(false)
    expect(storedMode()).toBe('auto')
  })

  it('keeps state, storage and the html class consistent across all three modes', () => {
    installMatchMedia({ dark: true })
    const store = useAppStore()
    store.initTheme()
    // auto 下跟随系统（深色）
    expect(store.isDark).toBe(true)

    store.setThemeMode('light')
    expect(store.themeMode).toBe('light')
    expect(store.isDark).toBe(false)
    expect(hasClass('light')).toBe(true)
    expect(document.documentElement.style.colorScheme).toBe('light')
    expect(storedMode()).toBe('light')

    store.setThemeMode('dark')
    expect(store.themeMode).toBe('dark')
    expect(store.isDark).toBe(true)
    expect(hasClass('dark')).toBe(true)
    expect(document.documentElement.style.colorScheme).toBe('dark')
    expect(storedMode()).toBe('dark')

    // 回到 auto 立刻对齐当前系统状态，不是简单地取反
    store.setThemeMode('auto')
    expect(store.themeMode).toBe('auto')
    expect(store.isDark).toBe(true)
    expect(storedMode()).toBe('auto')
  })

  it('follows system changes only while in auto mode', () => {
    installMatchMedia({ dark: false })
    const store = useAppStore()
    store.initTheme()
    expect(store.isDark).toBe(false)

    fireSystemChange(true)
    expect(store.systemDark).toBe(true)
    expect(store.isDark).toBe(true)
    expect(hasClass('dark')).toBe(true)

    // 手动选浅色后，系统再变也不应该把页面拖走
    store.setThemeMode('light')
    fireSystemChange(true)
    expect(store.systemDark).toBe(true)
    expect(store.isDark).toBe(false)
    expect(hasClass('light')).toBe(true)

    store.setThemeMode('auto')
    expect(store.isDark).toBe(true)
  })

  it('does not write storage when the selected mode is unchanged', () => {
    const store = useAppStore()
    store.initTheme()
    const spy = vi.spyOn(window.localStorage, 'setItem')
    store.setThemeMode('auto')
    expect(spy).not.toHaveBeenCalled()
    spy.mockRestore()
  })

  it('adds a one-shot transition class on manual switching', () => {
    vi.useFakeTimers()
    installMatchMedia({ dark: false })
    const store = useAppStore()
    store.setThemeMode('dark')

    expect(document.documentElement.classList.contains('theme-transition')).toBe(true)
    vi.advanceTimersByTime(300)
    expect(document.documentElement.classList.contains('theme-transition')).toBe(false)
  })

  it('skips the transition when the user asked for reduced motion', () => {
    installMatchMedia({ dark: false, reduceMotion: true })
    const store = useAppStore()
    store.setThemeMode('dark')

    expect(store.isDark).toBe(true)
    expect(document.documentElement.classList.contains('theme-transition')).toBe(false)
  })

  it('survives an environment without matchMedia', () => {
    Object.defineProperty(window, 'matchMedia', { configurable: true, writable: true, value: undefined })
    const store = useAppStore()

    expect(() => store.setThemeMode('dark')).not.toThrow()
    expect(store.isDark).toBe(true)
    expect(hasClass('dark')).toBe(true)
  })
})

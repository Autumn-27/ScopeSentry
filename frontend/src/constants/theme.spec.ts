import { describe, expect, it, beforeEach } from 'vitest'
import {
  DEFAULT_THEME_MODE,
  applyThemeClass,
  isThemeMode,
  resolveInitialThemeMode,
  resolveIsDark
} from './theme'

describe('theme constants', () => {
  it('defaults to follow-system', () => {
    expect(DEFAULT_THEME_MODE).toBe('auto')
  })

  it('validates modes', () => {
    expect(isThemeMode('auto')).toBe(true)
    expect(isThemeMode('light')).toBe(true)
    expect(isThemeMode('dark')).toBe(true)
    expect(isThemeMode('system')).toBe(false)
    expect(isThemeMode(null)).toBe(false)
    expect(isThemeMode(1)).toBe(false)
  })

  it('never migrates the legacy boolean into a fixed mode', () => {
    // 旧版的 isDark 是机器写入的，不能当成用户显式选择
    expect(resolveInitialThemeMode(true)).toBe('auto')
    expect(resolveInitialThemeMode(false)).toBe('auto')
    expect(resolveInitialThemeMode('dark')).toBe('dark')
    expect(resolveInitialThemeMode('garbage')).toBe('auto')
  })

  it('resolves the final dark state', () => {
    expect(resolveIsDark('auto', true)).toBe(true)
    expect(resolveIsDark('auto', false)).toBe(false)
    expect(resolveIsDark('dark', false)).toBe(true)
    expect(resolveIsDark('light', true)).toBe(false)
  })

  describe('applyThemeClass', () => {
    beforeEach(() => {
      document.documentElement.className = ''
      document.documentElement.style.colorScheme = ''
    })

    it('applies the dark class and color scheme', () => {
      applyThemeClass(true)
      expect(document.documentElement.classList.contains('dark')).toBe(true)
      expect(document.documentElement.classList.contains('light')).toBe(false)
      expect(document.documentElement.style.colorScheme).toBe('dark')
    })

    it('switches back to light', () => {
      applyThemeClass(true)
      applyThemeClass(false)
      expect(document.documentElement.classList.contains('light')).toBe(true)
      expect(document.documentElement.classList.contains('dark')).toBe(false)
      expect(document.documentElement.style.colorScheme).toBe('light')
    })
  })
})

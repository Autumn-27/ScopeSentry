import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useStorage } from './useStorage'

describe('useStorage', () => {
  beforeEach(() => {
    window.localStorage.clear()
  })

  it('round-trips values with their type wrapper', () => {
    const { setStorage, getStorage } = useStorage('localStorage')
    setStorage('flag', true)
    setStorage('count', 3)
    setStorage('name', 'scope')

    expect(getStorage('flag')).toBe(true)
    expect(getStorage('count')).toBe(3)
    expect(getStorage('name')).toBe('scope')
  })

  it('returns the raw value for hand-written JSON', () => {
    window.localStorage.setItem('bare', JSON.stringify('light'))
    const { getStorage } = useStorage('localStorage')
    expect(getStorage('bare')).toBe('light')
  })

  it('returns null instead of throwing on corrupted JSON', () => {
    window.localStorage.setItem('broken', '{not-json')
    const { getStorage } = useStorage('localStorage')
    expect(() => getStorage('broken')).not.toThrow()
    expect(getStorage('broken')).toBeNull()
  })

  it('returns null when the storage backend throws (private mode)', () => {
    const spy = vi.spyOn(window.localStorage, 'getItem').mockImplementation(() => {
      throw new Error('SecurityError')
    })
    const { getStorage } = useStorage('localStorage')
    expect(getStorage('anything')).toBeNull()
    spy.mockRestore()
  })

  it('returns the raw empty value for missing keys', () => {
    const { getStorage } = useStorage('localStorage')
    expect(getStorage('missing')).toBeNull()
  })

  it('clear keeps excluded keys', () => {
    const { setStorage, clear } = useStorage('localStorage')
    setStorage('keep-me', 1)
    setStorage('drop-me', 2)
    setStorage('dynamicRouter', 3)

    clear(['keep-me'])

    expect(window.localStorage.getItem('keep-me')).not.toBeNull()
    expect(window.localStorage.getItem('drop-me')).toBeNull()
    // 默认排除项不会被清掉
    expect(window.localStorage.getItem('dynamicRouter')).not.toBeNull()
  })
})

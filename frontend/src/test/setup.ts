/**
 * 测试环境补齐：
 * Node 26 自带一个未启用的 localStorage 全局（需要 --localstorage-file），
 * 它会挡住 jsdom 的实现，导致 window.localStorage 为 undefined。
 * 这里在缺失时装一个内存版 Storage，并像真实实现一样把键暴露为自有属性
 * （项目里的 clear() 依赖 Object.keys(localStorage)）。
 */
const backing = new WeakMap<object, Map<string, string>>()

const storeOf = (target: object): Map<string, string> => {
  let m = backing.get(target)
  if (!m) {
    m = new Map()
    backing.set(target, m)
  }
  return m
}

class MemoryStorage implements Storage {
  [key: string]: unknown

  get length(): number {
    return storeOf(this).size
  }

  clear(): void {
    for (const key of Array.from(storeOf(this).keys())) {
      this.removeItem(key)
    }
  }

  getItem(key: string): string | null {
    const m = storeOf(this)
    return m.has(key) ? (m.get(key) as string) : null
  }

  key(index: number): string | null {
    return Array.from(storeOf(this).keys())[index] ?? null
  }

  removeItem(key: string): void {
    storeOf(this).delete(key)
    if (Object.prototype.hasOwnProperty.call(this, key)) {
      delete this[key]
    }
  }

  setItem(key: string, value: string): void {
    storeOf(this).set(key, String(value))
    this[key] = String(value)
  }
}

const install = (target: Window & typeof globalThis) => {
  if (!target.localStorage) {
    Object.defineProperty(target, 'localStorage', {
      value: new MemoryStorage(),
      configurable: true
    })
  }
  if (!target.sessionStorage) {
    Object.defineProperty(target, 'sessionStorage', {
      value: new MemoryStorage(),
      configurable: true
    })
  }
}

if (typeof window !== 'undefined') {
  install(window as unknown as Window & typeof globalThis)
}

export {}

/**
 * 主题逻辑自检
 * 运行：pnpm theme:check  ->  esno ./scripts/check-theme.ts
 *
 * 前半部分是**行为断言**：直接调用 src/constants/theme.ts 的真实函数，
 * 用最小的假 document / matchMedia 驱动，验证三态决策与类切换。
 * 后半部分只是源码接线 smoke（证明入口确实挂上了），不作为逻辑证明。
 */
import { readFileSync } from 'node:fs'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'
import {
  DEFAULT_THEME_MODE,
  THEME_MODE_STORAGE_KEY,
  THEME_MODE_OPTIONS,
  isThemeMode,
  resolveIsDark,
  resolveInitialThemeMode,
  applyThemeClass,
  getSystemDark
} from '../src/constants/theme.ts'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..')

let failures = 0
const check = (name: string, ok: boolean, detail = '') => {
  if (ok) {
    console.log(`PASS ${name}`)
  } else {
    failures++
    console.log(`FAIL ${name}${detail ? ' -> ' + detail : ''}`)
  }
}

// ---------------------------------------------------------------- 假 DOM
type FakeClassList = {
  value: Set<string>
  add: (c: string) => void
  remove: (c: string) => void
  toggle: (c: string, force?: boolean) => void
  contains: (c: string) => boolean
}
const createClassList = (): FakeClassList => {
  const value = new Set<string>()
  return {
    value,
    add: (c) => void value.add(c),
    remove: (c) => void value.delete(c),
    toggle: (c, force) => {
      const next = force === undefined ? !value.has(c) : force
      if (next) value.add(c)
      else value.delete(c)
    },
    contains: (c) => value.has(c)
  }
}
const fakeRoot = { classList: createClassList(), style: { colorScheme: '' } }
const g = globalThis as any
g.document = { documentElement: fakeRoot }

// ---------------------------------------------------------------- 行为断言
check('默认模式为 auto（跟随系统）', DEFAULT_THEME_MODE === 'auto', DEFAULT_THEME_MODE)
check('持久化键名稳定', THEME_MODE_STORAGE_KEY === 'themeMode', THEME_MODE_STORAGE_KEY)
check('只暴露三种模式', THEME_MODE_OPTIONS.join(',') === 'auto,light,dark', THEME_MODE_OPTIONS.join(','))

for (const mode of ['auto', 'light', 'dark']) {
  check(`isThemeMode('${mode}')`, isThemeMode(mode) === true)
}
for (const bad of ['system', 'Dark', '', null, undefined, 1, {}, true, false]) {
  check(`isThemeMode(${JSON.stringify(bad)}) 为 false`, isThemeMode(bad) === false)
}

// 启动模式：缺失/非法 -> auto；旧版布尔 isDark 不能把用户钉死
check('未持久化 -> auto', resolveInitialThemeMode(null) === 'auto')
check('非法值 -> auto', resolveInitialThemeMode('system') === 'auto')
check('旧版 isDark=true 不迁移为 dark', resolveInitialThemeMode(true) === 'auto')
check('旧版 isDark=false 不迁移为 light', resolveInitialThemeMode(false) === 'auto')
check('合法 light 保留', resolveInitialThemeMode('light') === 'light')
check('合法 dark 保留', resolveInitialThemeMode('dark') === 'dark')
check('合法 auto 保留', resolveInitialThemeMode('auto') === 'auto')

// 模式矩阵
check('auto + 系统深色 = 深色', resolveIsDark('auto', true) === true)
check('auto + 系统浅色 = 浅色', resolveIsDark('auto', false) === false)
check('dark 覆盖系统浅色', resolveIsDark('dark', false) === true)
check('light 覆盖系统深色', resolveIsDark('light', true) === false)

// 类切换真实执行
applyThemeClass(true)
check('深色：html 有 dark 类', fakeRoot.classList.contains('dark'))
check('深色：html 没有 light 类', !fakeRoot.classList.contains('light'))
check('深色：color-scheme=dark', fakeRoot.style.colorScheme === 'dark')
applyThemeClass(false)
check('浅色：html 有 light 类', fakeRoot.classList.contains('light'))
check('浅色：html 没有 dark 类', !fakeRoot.classList.contains('dark'))
check('浅色：color-scheme=light', fakeRoot.style.colorScheme === 'light')

// 系统深浅色读取
g.window = {
  matchMedia: (q: string) => ({ matches: q.includes('dark'), media: q })
}
check('getSystemDark 读 matchMedia', getSystemDark() === true)
g.window = { matchMedia: () => ({ matches: false, media: '' }) }
check('getSystemDark 浅色系统', getSystemDark() === false)
g.window = {}
check('无 matchMedia 时安全降级为浅色', getSystemDark() === false)
delete g.window

// 无 document 环境（SSR）不能崩
const savedDoc = g.document
delete g.document
let threw = false
try {
  applyThemeClass(true)
} catch {
  threw = true
}
check('无 document 时 applyThemeClass 安全返回', threw === false)
g.document = savedDoc

// ---------------------------------------------------------------- 接线 smoke
const indexHtml = readFileSync(resolve(root, 'index.html'), 'utf8')
check('[smoke] index.html 首屏读取 themeMode', indexHtml.includes("readValue('themeMode')"))
check('[smoke] index.html 解析 useStorage 包装格式', indexHtml.includes("'value' in parsed"))
check('[smoke] index.html 不再把旧版 isDark 当选择', !indexHtml.includes("readValue('isDark')"))
check('[smoke] index.html 深色加载背景', indexHtml.includes('html.dark .app-loading'))

const store = readFileSync(resolve(root, 'src/store/modules/app.ts'), 'utf8')
check('[smoke] store 监听系统变化', store.includes("addEventListener('change'"))
check('[smoke] store 兼容旧 Safari 的 addListener', store.includes('addListener'))
check('[smoke] store 仅 auto 模式跟随系统', store.includes("this.themeMode === 'auto'"))
check('[smoke] store 有 SSR 守卫', store.includes("typeof window === 'undefined'"))
check('[smoke] store 调用统一应用函数', store.includes('applyThemeClass(isDark)'))

const themeSwitch = readFileSync(resolve(root, 'src/components/ThemeSwitch/src/ThemeSwitch.vue'), 'utf8')
for (const mode of ["'light'", "'dark'", "'auto'"]) {
  check(`[smoke] ThemeSwitch 提供 ${mode} 选项`, themeSwitch.includes(`value: ${mode}`))
}

const header = readFileSync(resolve(root, 'src/layout/components/ToolHeader.vue'), 'utf8')
check('[smoke] 顶栏挂载 ThemeSwitch', header.includes('<ThemeSwitch'))

const storage = readFileSync(resolve(root, 'src/hooks/web/useStorage.ts'), 'utf8')
check('[smoke] getStorage 有 try/catch 容错', /const getStorage[\s\S]{0,200}try \{/.test(storage))

console.log(`\n${failures === 0 ? 'ALL PASS' : failures + ' FAILED'}`)
process.exit(failures === 0 ? 0 : 1)

/**
 * 主题逻辑自检（无浏览器依赖）
 * 运行：pnpm theme:check  ->  esno ./scripts/check-theme.ts
 */
import { readFileSync } from 'node:fs'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'
import {
  DEFAULT_THEME_MODE,
  THEME_MODE_STORAGE_KEY,
  LEGACY_IS_DARK_STORAGE_KEY,
  isThemeMode,
  resolveIsDark,
  applyThemeClass
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

// 1. 默认值必须跟随系统
check('默认模式为 auto（跟随系统）', DEFAULT_THEME_MODE === 'auto', DEFAULT_THEME_MODE)
check('持久化键名稳定', THEME_MODE_STORAGE_KEY === 'themeMode', THEME_MODE_STORAGE_KEY)
check('旧版键名保留', LEGACY_IS_DARK_STORAGE_KEY === 'isDark', LEGACY_IS_DARK_STORAGE_KEY)

// 2. 模式校验
for (const mode of ['auto', 'light', 'dark']) {
  check(`isThemeMode('${mode}')`, isThemeMode(mode) === true)
}
for (const bad of ['system', 'Dark', '', null, undefined, 1, {}]) {
  check(`isThemeMode(${JSON.stringify(bad)}) 为 false`, isThemeMode(bad) === false)
}

// 3. 模式 -> 最终深浅色
check('auto + 系统深色 = 深色', resolveIsDark('auto', true) === true)
check('auto + 系统浅色 = 浅色', resolveIsDark('auto', false) === false)
check('dark 覆盖系统浅色', resolveIsDark('dark', false) === true)
check('light 覆盖系统深色', resolveIsDark('light', true) === false)

// 4. 无 DOM 环境不能崩（SSR / 脚本环境）
let threw = false
try {
  applyThemeClass(true)
} catch {
  threw = true
}
check('无 document 时 applyThemeClass 安全返回', threw === false)

// 5. 关键接线：首屏脚本、store 监听、UI 三选项
const indexHtml = readFileSync(resolve(root, 'index.html'), 'utf8')
check('index.html 含首屏主题脚本', indexHtml.includes("readValue('themeMode')"))
check('index.html 解析 useStorage 包装格式', indexHtml.includes("'value' in parsed"))
check('index.html 兼容旧版 isDark', indexHtml.includes("readValue('isDark')"))
check('index.html 深色加载背景', indexHtml.includes('html.dark .app-loading'))

const store = readFileSync(resolve(root, 'src/store/modules/app.ts'), 'utf8')
check('store 读取系统深浅色', store.includes('prefers-color-scheme: dark'))
check('store 监听系统变化', store.includes("addEventListener('change'"))
check('store 在 auto 模式才跟随系统', store.includes("this.themeMode === 'auto'"))
check('store 迁移旧版 isDark', store.includes('legacyIsDark === true'))
check('store 调用统一的主题应用函数', store.includes('applyThemeClass(isDark)'))

const themeConstants = readFileSync(resolve(root, 'src/constants/theme.ts'), 'utf8')
check('应用主题时同步 color-scheme', themeConstants.includes('colorScheme'))
check('应用主题时切换 dark/light 类', themeConstants.includes("classList.toggle('dark'"))

const themeSwitch = readFileSync(resolve(root, 'src/components/ThemeSwitch/src/ThemeSwitch.vue'), 'utf8')
for (const mode of ["'light'", "'dark'", "'auto'"]) {
  check(`ThemeSwitch 提供 ${mode} 选项`, themeSwitch.includes(`value: ${mode}`))
}

const header = readFileSync(resolve(root, 'src/layout/components/ToolHeader.vue'), 'utf8')
check('顶栏挂载 ThemeSwitch', header.includes('<ThemeSwitch'))

console.log(`\n${failures === 0 ? 'ALL PASS' : failures + ' FAILED'}`)
process.exit(failures === 0 ? 0 : 1)

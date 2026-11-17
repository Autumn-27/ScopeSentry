import { defineStore } from 'pinia'
import { store } from '../index'
import { setCssVar, humpToUnderline, getCssVar } from '@/utils'
import { ElMessage, ComponentSize } from 'element-plus'
import { colorIsDark, hexToRGB, lighten, mix } from '@/utils/color'
import { unref } from 'vue'
import { useCssVar } from '@vueuse/core'
import { useStorage } from '@/hooks/web/useStorage'
import {
  ThemeMode,
  THEME_MODE_STORAGE_KEY,
  LEGACY_IS_DARK_STORAGE_KEY,
  DEFAULT_THEME_MODE,
  getSystemDark,
  isThemeMode,
  resolveInitialThemeMode,
  resolveIsDark,
  applyThemeClass
} from '@/constants/theme'

const { setStorage, getStorage } = useStorage('localStorage')

interface AppState {
  breadcrumb: boolean
  breadcrumbIcon: boolean
  collapse: boolean
  uniqueOpened: boolean
  hamburger: boolean
  screenfull: boolean
  size: boolean
  locale: boolean
  tagsView: boolean
  tagsViewIcon: boolean
  logo: boolean
  fixedHeader: boolean
  greyMode: boolean
  dynamicRouter: boolean
  serverDynamicRouter: boolean
  pageLoading: boolean
  layout: LayoutType
  title: string
  isDark: boolean
  // 主题模式：auto(跟随系统) / light / dark
  themeMode: ThemeMode
  // 系统当前是否深色，由 prefers-color-scheme 同步
  systemDark: boolean
  currentSize: ComponentSize
  sizeMap: ComponentSize[]
  mobile: boolean
  footer: boolean
  theme: ThemeTypes
  fixedMenu: boolean
}

let themeStorageBound = false

/**
 * 多标签页主题同步：A 标签切换主题后，B 标签立即跟随
 * （storage 事件只在其它标签触发，所以不会形成回环）
 */
const bindThemeStorageSync = (storeInstance: {
  setThemeMode: (mode: ThemeMode) => void
  themeMode: ThemeMode
}) => {
  if (themeStorageBound || typeof window === 'undefined') return
  themeStorageBound = true
  window.addEventListener('storage', (event) => {
    if (event.key !== THEME_MODE_STORAGE_KEY || !event.newValue) return
    try {
      const parsed = JSON.parse(event.newValue)
      const mode = parsed && typeof parsed === 'object' && 'value' in parsed ? parsed.value : parsed
      // 注意：Pinia 的 getter 是属性而不是方法，这里读 state 上的 themeMode 值
      if (isThemeMode(mode) && mode !== storeInstance.themeMode) {
        storeInstance.setThemeMode(mode)
      }
    } catch {
      /* 忽略损坏值 */
    }
  })
}

let mediaQueryList: MediaQueryList | null = null
let mediaQueryHandler: ((event: MediaQueryListEvent) => void) | null = null
let themeTransitionTimer: number | null = null

/**
 * 监听系统深浅色变化，只在“跟随系统”模式下实时生效
 * 单例绑定，避免 HMR / 重复初始化堆积监听
 */
const bindSystemThemeListener = (storeInstance: { setSystemDark: (value: boolean) => void }) => {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  // matchMedia 每次返回新对象，不能靠引用判断是否已绑定，统一先解绑再绑
  if (mediaQueryList && mediaQueryHandler) {
    if (typeof mediaQueryList.removeEventListener === 'function') {
      mediaQueryList.removeEventListener('change', mediaQueryHandler)
    } else if (typeof mediaQueryList.removeListener === 'function') {
      // Safari <= 13 / 旧 WKWebView 只有废弃的 add/removeListener
      mediaQueryList.removeListener(mediaQueryHandler)
    }
  }
  mediaQueryList = mediaQuery
  mediaQueryHandler = (event: MediaQueryListEvent) => storeInstance.setSystemDark(event.matches)
  if (typeof mediaQuery.addEventListener === 'function') {
    mediaQuery.addEventListener('change', mediaQueryHandler)
  } else if (typeof mediaQuery.addListener === 'function') {
    mediaQuery.addListener(mediaQueryHandler)
  }
}

export const useAppStore = defineStore('app', {
  state: (): AppState => {
    return {
      sizeMap: ['default', 'large', 'small'],
      mobile: false, // 是否是移动端
      title: import.meta.env.VITE_APP_TITLE, // 标题
      pageLoading: false, // 路由跳转loading
      breadcrumb: true, // 面包屑
      breadcrumbIcon: true, // 面包屑图标
      collapse: false, // 折叠菜单
      uniqueOpened: false, // 是否只保持一个子菜单的展开
      hamburger: true, // 折叠图标
      screenfull: true, // 全屏图标
      size: true, // 尺寸图标
      locale: true, // 多语言图标
      tagsView: true, // 标签页
      tagsViewIcon: false, // 是否显示标签图标
      logo: true, // logo
      fixedHeader: true, // 固定toolheader
      footer: true, // 显示页脚
      greyMode: false, // 是否开始灰色模式，用于特殊悼念日
      dynamicRouter: false, // 是否动态路由
      serverDynamicRouter: false, // 是否服务端渲染动态路由
      fixedMenu: false, // 是否固定菜单

      layout: 'classic', // layout布局
      isDark: false, // 是否应用暗黑模式（最终态，由 themeMode + 系统决定）
      themeMode: DEFAULT_THEME_MODE, // 主题模式，默认跟随系统
      systemDark: false, // 系统当前是否深色
      currentSize: 'default', // 组件尺寸
      theme: {
        // 主题色
        elColorPrimary: '#2f6feb',
        // 左侧菜单边框颜色
        leftMenuBorderColor: 'inherit',
        // 左侧菜单背景颜色
        leftMenuBgColor: '#fff',
        // 左侧菜单浅色背景颜色
        leftMenuBgLightColor: '#0f2438',
        // 左侧菜单选中背景颜色
        leftMenuBgActiveColor: 'var(--el-color-primary)',
        // 左侧菜单收起选中背景颜色
        leftMenuCollapseBgActiveColor: 'var(--el-color-primary)',
        // 左侧菜单字体颜色
        leftMenuTextColor: '#bfcbd9',
        // 左侧菜单选中字体颜色
        leftMenuTextActiveColor: '#fff',
        // logo字体颜色
        logoTitleTextColor: '#fff',
        // logo边框颜色
        logoBorderColor: 'inherit',
        // 头部背景颜色
        topHeaderBgColor: '#fff',
        // 头部字体颜色
        topHeaderTextColor: 'inherit',
        // 头部悬停颜色
        topHeaderHoverColor: '#f6f6f6',
        // 头部边框颜色
        topToolBorderColor: '#eee'
      }
    }
  },
  getters: {
    getBreadcrumb(): boolean {
      return this.breadcrumb
    },
    getBreadcrumbIcon(): boolean {
      return this.breadcrumbIcon
    },
    getCollapse(): boolean {
      return this.collapse
    },
    getUniqueOpened(): boolean {
      return this.uniqueOpened
    },
    getHamburger(): boolean {
      return this.hamburger
    },
    getScreenfull(): boolean {
      return this.screenfull
    },
    getSize(): boolean {
      return this.size
    },
    getLocale(): boolean {
      return this.locale
    },
    getTagsView(): boolean {
      return this.tagsView
    },
    getTagsViewIcon(): boolean {
      return this.tagsViewIcon
    },
    getLogo(): boolean {
      return this.logo
    },
    getFixedHeader(): boolean {
      return this.fixedHeader
    },
    getGreyMode(): boolean {
      return this.greyMode
    },
    getDynamicRouter(): boolean {
      return this.dynamicRouter
    },
    getServerDynamicRouter(): boolean {
      return this.serverDynamicRouter
    },
    getFixedMenu(): boolean {
      return this.fixedMenu
    },
    getPageLoading(): boolean {
      return this.pageLoading
    },
    getLayout(): LayoutType {
      return this.layout
    },
    getTitle(): string {
      return 'Scope Sentry'
    },
    getIsDark(): boolean {
      return this.isDark
    },
    getThemeMode(): ThemeMode {
      return this.themeMode
    },
    getSystemDark(): boolean {
      return this.systemDark
    },
    getCurrentSize(): ComponentSize {
      return this.currentSize
    },
    getSizeMap(): ComponentSize[] {
      return this.sizeMap
    },
    getMobile(): boolean {
      return this.mobile
    },
    getTheme(): ThemeTypes {
      return this.theme
    },
    getFooter(): boolean {
      return this.footer
    }
  },
  actions: {
    setBreadcrumb(breadcrumb: boolean) {
      this.breadcrumb = breadcrumb
    },
    setBreadcrumbIcon(breadcrumbIcon: boolean) {
      this.breadcrumbIcon = breadcrumbIcon
    },
    setCollapse(collapse: boolean) {
      this.collapse = collapse
    },
    setUniqueOpened(uniqueOpened: boolean) {
      this.uniqueOpened = uniqueOpened
    },
    setHamburger(hamburger: boolean) {
      this.hamburger = hamburger
    },
    setScreenfull(screenfull: boolean) {
      this.screenfull = screenfull
    },
    setSize(size: boolean) {
      this.size = size
    },
    setLocale(locale: boolean) {
      this.locale = locale
    },
    setTagsView(tagsView: boolean) {
      this.tagsView = tagsView
    },
    setTagsViewIcon(tagsViewIcon: boolean) {
      this.tagsViewIcon = tagsViewIcon
    },
    setLogo(logo: boolean) {
      this.logo = logo
    },
    setFixedHeader(fixedHeader: boolean) {
      this.fixedHeader = fixedHeader
    },
    setGreyMode(greyMode: boolean) {
      this.greyMode = greyMode
    },
    setDynamicRouter(dynamicRouter: boolean) {
      this.dynamicRouter = dynamicRouter
    },
    setServerDynamicRouter(serverDynamicRouter: boolean) {
      this.serverDynamicRouter = serverDynamicRouter
    },
    setFixedMenu(fixedMenu: boolean) {
      this.fixedMenu = fixedMenu
    },
    setPageLoading(pageLoading: boolean) {
      this.pageLoading = pageLoading
    },
    setLayout(layout: LayoutType) {
      if (this.mobile && layout !== 'classic') {
        ElMessage.warning('移动端模式下不支持切换其它布局')
        return
      }
      this.layout = layout
    },
    setTitle(title: string) {
      this.title = title
    },
    /**
     * 应用最终的主题状态（dark 类、color-scheme、菜单/头部配色）
     * 一般不直接调用，请用 setThemeMode / setSystemDark
     */
    setIsDark(isDark: boolean) {
      this.isDark = isDark
      // 兼容旧版本：老版本只认 isDark 布尔值
      setStorage(LEGACY_IS_DARK_STORAGE_KEY, isDark)
      applyThemeClass(isDark)
      this.setPrimaryLight()
      // 深浅切换后让菜单与头部跟随新的基底背景色
      const baseBgColor = getCssVar('--el-bg-color').trim()
      if (baseBgColor) {
        this.setMenuTheme(baseBgColor)
        this.setHeaderTheme(baseBgColor)
      }
    },
    /** 手动切换：日间 / 夜间 / 跟随系统 */
    setThemeMode(themeMode: ThemeMode) {
      if (themeMode === this.themeMode) return
      this.themeMode = themeMode
      setStorage(THEME_MODE_STORAGE_KEY, themeMode)
      // 手动切换时给整页一次短过渡，避免颜色硬切；首次加载/跟随系统变化不加
      this.playThemeTransition()
      this.setIsDark(resolveIsDark(themeMode, this.systemDark))
    },
    /** 主题切换的 260ms 全局过渡（系统级变化与人机切换的手感差异） */
    playThemeTransition() {
      if (typeof window === 'undefined' || typeof document === 'undefined') return
      const root = document.documentElement
      // 老 WebView / 非浏览器环境可能没有 matchMedia，取不到就按“不减弱动效”处理
      if (window.matchMedia?.('(prefers-reduced-motion: reduce)')?.matches) return
      root.classList.add('theme-transition')
      if (themeTransitionTimer) window.clearTimeout(themeTransitionTimer)
      themeTransitionTimer = window.setTimeout(() => {
        root.classList.remove('theme-transition')
        themeTransitionTimer = null
      }, 260)
    },
    /** 系统主题变化（prefers-color-scheme） */
    setSystemDark(systemDark: boolean) {
      this.systemDark = systemDark
      if (this.themeMode === 'auto') {
        this.setIsDark(systemDark)
      }
    },
    /** 启动时初始化主题：读取持久化模式（缺省/非法 -> 跟随系统）并监听系统变化 */
    initTheme() {
      // SSR/预渲染环境没有 window，直接跳过（其余主题逻辑都有 typeof window 守卫）
      if (typeof window === 'undefined') return
      this.systemDark = getSystemDark()
      const themeMode = resolveInitialThemeMode(getStorage(THEME_MODE_STORAGE_KEY))
      this.themeMode = themeMode
      setStorage(THEME_MODE_STORAGE_KEY, themeMode)
      this.setIsDark(resolveIsDark(themeMode, this.systemDark))
      bindSystemThemeListener(this)
      bindThemeStorageSync(this)
    },
    setCurrentSize(currentSize: ComponentSize) {
      this.currentSize = currentSize
    },
    setMobile(mobile: boolean) {
      this.mobile = mobile
    },
    setTheme(theme: ThemeTypes) {
      this.theme = Object.assign(this.theme, theme)
    },
    setCssVarTheme() {
      for (const key in this.theme) {
        setCssVar(`--${humpToUnderline(key)}`, this.theme[key])
      }
      this.setPrimaryLight()
    },
    setPrimaryLight() {
      if (this.theme.elColorPrimary) {
        const elColorPrimary = this.theme.elColorPrimary
        const color = this.isDark ? '#000000' : '#ffffff'
        const lightList = [3, 5, 7, 8, 9]
        lightList.forEach((v) => {
          setCssVar(`--el-color-primary-light-${v}`, mix(color, elColorPrimary, v / 10))
        })
        setCssVar(`--el-color-primary-dark-2`, mix(color, elColorPrimary, 0.2))
      }
    },
    setMenuTheme(color: string) {
      const primaryColor = useCssVar('--el-color-primary', document.documentElement)
      const isDarkColor = colorIsDark(color)
      const theme: Recordable = {
        // 左侧菜单边框颜色
        leftMenuBorderColor: isDarkColor ? 'inherit' : '#eee',
        // 左侧菜单背景颜色
        leftMenuBgColor: color,
        // 左侧菜单浅色背景颜色
        leftMenuBgLightColor: isDarkColor ? lighten(color!, 6) : color,
        // 左侧菜单选中背景颜色
        leftMenuBgActiveColor: isDarkColor
          ? 'var(--el-color-primary)'
          : hexToRGB(unref(primaryColor), 0.1),
        // 左侧菜单收起选中背景颜色
        leftMenuCollapseBgActiveColor: isDarkColor
          ? 'var(--el-color-primary)'
          : hexToRGB(unref(primaryColor), 0.1),
        // 左侧菜单字体颜色
        leftMenuTextColor: isDarkColor ? '#bfcbd9' : '#333',
        // 左侧菜单选中字体颜色
        leftMenuTextActiveColor: isDarkColor ? '#fff' : 'var(--el-color-primary)',
        // logo字体颜色
        logoTitleTextColor: isDarkColor ? '#fff' : 'inherit',
        // logo边框颜色
        logoBorderColor: isDarkColor ? color : '#eee'
      }
      this.setTheme(theme)
      this.setCssVarTheme()
    },
    setHeaderTheme(color: string) {
      const isDarkColor = colorIsDark(color)
      const textColor = isDarkColor ? '#fff' : 'inherit'
      const textHoverColor = isDarkColor ? lighten(color!, 6) : '#f6f6f6'
      const topToolBorderColor = isDarkColor ? color : '#eee'
      setCssVar('--top-header-bg-color', color)
      setCssVar('--top-header-text-color', textColor)
      setCssVar('--top-header-hover-color', textHoverColor)
      this.setTheme({
        topHeaderBgColor: color,
        topHeaderTextColor: textColor,
        topHeaderHoverColor: textHoverColor,
        topToolBorderColor
      })
      if (this.getLayout === 'top') {
        this.setMenuTheme(color)
      }
    },
    setFooter(footer: boolean) {
      this.footer = footer
    }
  },
  persist: true
})

export const useAppStoreWithOut = () => {
  return useAppStore(store)
}

import type { Composer } from 'vue-i18n'
import { i18n } from '@/plugins/vueI18n'

// Resolve only registered keys. Server prose and unknown keys remain unchanged.
export const translateApiMessage = (message: unknown): string => {
  if (typeof message !== 'string' || !i18n) return typeof message === 'string' ? message : ''
  const { te, t, locale } = i18n.global as unknown as Composer
  const lang = locale.value
  if (te(message, lang) || te(message, 'en')) return String(t(message))
  // 后端新加的 key 还没翻：成功类回退到通用成功文案，避免把 "api.user.login.success" 原样弹给用户
  if (message.endsWith('.success')) return String(t('api.success'))
  return message
}

import type { Composer } from 'vue-i18n'
import { i18n } from '@/plugins/vueI18n'

// Resolve only registered keys. Server prose and unknown keys remain unchanged.
export const translateApiMessage = (message: unknown): string => {
  if (typeof message !== 'string' || !i18n) return typeof message === 'string' ? message : ''
  const { te, t, locale } = i18n.global as unknown as Composer
  const lang = locale.value
  return te(message, lang) || te(message, 'en') ? String(t(message)) : message
}

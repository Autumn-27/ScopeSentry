import fs from 'node:fs'
import path from 'node:path'
import vm from 'node:vm'
import assert from 'node:assert/strict'
import { fileURLToPath } from 'node:url'
import { createRequire } from 'node:module'

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const require = createRequire(path.join(root, 'package.json'))
const { createI18n } = require('vue-i18n')
const load = (lang) =>
  vm.runInNewContext(
    `(${fs.readFileSync(path.join(root, 'src/locales', `${lang}.ts`), 'utf8').replace(/^export default/, '')})`
  )
const flatten = (obj, prefix = '') =>
  Object.fromEntries(
    Object.entries(obj).flatMap(([key, value]) =>
      typeof value === 'string'
        ? [[prefix + key, value]]
        : Object.entries(flatten(value, `${prefix}${key}.`))
    )
  )
const messages = Object.fromEntries(['en', 'zh-CN', 'ko-KR'].map((lang) => [lang, load(lang)]))
const locales = Object.fromEntries(
  Object.entries(messages).map(([lang, value]) => [lang, flatten(value)])
)
const placeholders = (value) => [...value.matchAll(/\{(\w+)\}/g)].map((m) => m[1]).sort()
const sourceKeys = new Set([...Object.keys(locales.en), ...Object.keys(locales['zh-CN'])])
for (const key of sourceKeys) {
  assert.ok(key in locales['ko-KR'], `Missing Korean key: ${key}`)
  const source = locales.en[key] ?? locales['zh-CN'][key]
  assert.deepEqual(
    placeholders(locales['ko-KR'][key]),
    placeholders(source),
    `Interpolation mismatch: ${key}`
  )
  assert.ok(locales['ko-KR'][key].length || source === '', `Empty Korean value: ${key}`)
}
assert.equal(Object.keys(locales['ko-KR']).length, sourceKeys.size, 'Unexpected Korean keys')
const i18n = createI18n({
  legacy: false,
  locale: 'ko-KR',
  fallbackLocale: 'en',
  messages,
  missingWarn: false,
  fallbackWarn: false
})
assert.equal(i18n.global.t('common.lengthRange', { min: 2, max: 10 }), '길이는 2~10자여야 합니다.')
assert.ok(i18n.global.t('apiKey.deleteConfirm', { name: 'test-key' }).includes('test-key'))
assert.ok(
  i18n.global.t('plugin.confirmInstall', { action: '설치', name: 'Nuclei' }).includes('Nuclei')
)
// Remove one key from an isolated copy to prove English fallback is available.
const fallbackMessages = structuredClone(messages)
delete fallbackMessages['ko-KR'].common.save
const fallback = createI18n({
  legacy: false,
  locale: 'ko-KR',
  fallbackLocale: 'en',
  messages: fallbackMessages,
  missingWarn: false,
  fallbackWarn: false
})
assert.equal(fallback.global.t('common.save'), 'Save')
console.log(
  JSON.stringify(
    {
      keys: Object.fromEntries(
        Object.entries(locales).map(([lang, values]) => [lang, Object.keys(values).length])
      ),
      missingKoreanKeys: 0,
      interpolation: 'PASS',
      englishFallback: 'PASS'
    },
    null,
    2
  )
)

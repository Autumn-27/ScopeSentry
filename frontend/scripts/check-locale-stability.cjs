// Local component regression checks; no backend, credentials or network requests.
const assert = require('node:assert/strict')
const fs = require('node:fs')
const path = require('node:path')
const vm = require('node:vm')
const { createRequire } = require('node:module')
const root = path.resolve(__dirname, '..')
const fromFrontend = createRequire(path.join(root, 'package.json'))
const ts = fromFrontend('typescript')
const { parse, compileScript } = fromFrontend('@vue/compiler-sfc')
const vue = fromFrontend('vue')

function evaluate(source, mocks, filename) {
  const output = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020 }
  }).outputText
  const module = { exports: {} }
  vm.runInNewContext(output, {
    module,
    exports: module.exports,
    require: (name) => name in mocks ? mocks[name] : fromFrontend(name),
    console,
    document: mocks.document
  }, { filename })
  return module.exports
}

function component(relative, mocks, inlineTemplate = false) {
  const filename = path.join(root, relative)
  const { descriptor } = parse(fs.readFileSync(filename, 'utf8'), { filename })
  const script = compileScript(descriptor, { id: 'locale-regression', inlineTemplate })
  return evaluate(script.content, mocks, filename).default
}

const design = { useDesign: () => ({ getPrefixCls: (s) => s, variables: { elNamespace: 'el' } }) }
const props = { propTypes: { string: { def: (value) => ({ type: String, default: value }) }, oneOf: () => ({ def: (value) => ({ default: value }) }) } }
const titleSource = fs.readFileSync(path.join(root, 'src/hooks/web/useTitle.ts'), 'utf8')
let watcherCount = 0
const doc = { title: '' }
let language = 'en'
const title = evaluate(titleSource, {
  vue: { ...vue, watch: (...args) => { watcherCount++; return vue.watch(...args) } },
  '@/utils/is': { isString: (s) => typeof s === 'string' },
  '@/store/modules/app': { useAppStoreWithOut: () => ({ getTitle: 'Scope Sentry' }) },
  '@/hooks/web/useI18n': { useI18n: () => ({ t: (key) => `${language}:${key}` }) },
  document: doc
}, 'useTitle.ts')
for (let i = 0; i < 100; i++) {
  language = ['en', 'ko-KR', 'zh-CN'][i % 3]
  title.setPageTitle('router.login')
  assert.equal(doc.title, `Scope Sentry - ${language}:router.login`)
}
title.setPageTitle()
assert.equal(doc.title, 'Scope Sentry')
assert.equal(watcherCount, 0, 'One-off title updates must not create watchers')
for (const relative of ['src/permission.ts', 'src/hooks/web/useLocale.ts']) {
  const source = fs.readFileSync(path.join(root, relative), 'utf8')
  assert.ok(source.includes('setPageTitle('))
  assert.ok(!/\buseTitle\(/.test(source), `${relative} creates an event-owned watcher`)
}

async function main() {
  const locale = vue.reactive({ lang: 'ko-KR', elLocale: {} })
  let resolveConfirmation, rejectConfirmation
  let confirmations = 0, changes = 0, failImport = false
  const dropdown = component('src/components/LocaleDropdown/src/LocaleDropdown.vue', {
    'element-plus': { ElMessageBox: { confirm: () => { confirmations++; return new Promise((resolve, reject) => { resolveConfirmation = resolve; rejectConfirmation = reject }) } } },
    '@/store/modules/locale': { useLocaleStore: () => ({ getCurrentLocale: locale, getLocaleMap: [] }) },
    '@/hooks/web/useI18n': { useI18n: () => ({ t: (key) => key }) },
    '@/hooks/web/useLocale': { useLocale: () => ({ changeLocale: async (lang) => { if (failImport) throw new Error('Import failed'); changes++; locale.lang = lang } }) },
    '@/utils/propTypes': props,
    '@/hooks/web/useDesign': design
  })
  const setup = dropdown.setup({}, { expose() {} })
  await setup.setLang('ko-KR')
  assert.equal(confirmations, 0)
  const canceled = setup.setLang('en')
  assert.equal(changes, 0, 'No locale/storage mutation before confirmation')
  await setup.setLang('zh-CN')
  assert.equal(confirmations, 1, 'Overlapping choices must not create another operation')
  rejectConfirmation('cancel')
  await canceled
  assert.equal(locale.lang, 'ko-KR')
  assert.equal(setup.isChanging.value, false)
  const accepted = setup.setLang('en')
  resolveConfirmation()
  await accepted
  assert.equal(locale.lang, 'en')
  assert.equal(changes, 1)
  const closed = setup.setLang('zh-CN')
  rejectConfirmation('close')
  await closed
  assert.equal(locale.lang, 'en')
  failImport = true
  const failed = setup.setLang('zh-CN')
  resolveConfirmation()
  await assert.rejects(failed, /Import failed/)
  assert.equal(locale.lang, 'en')
  assert.equal(setup.isChanging.value, false, 'A failed load must unlock the menu')

  // Mount the real ConfigGlobal with a minimal renderer and a synthetic child.
  // This proves the reset boundary and async callback lifetime, not real API behavior.
  let providerMounts = 0, providerUnmounts = 0
  const config = component('src/components/ConfigGlobal/src/ConfigGlobal.vue', {
    '@/utils/propTypes': props,
    'element-plus': { ElConfigProvider: { inheritAttrs: false, setup: (_, { slots }) => {
      providerMounts++
      vue.onUnmounted(() => { providerUnmounts++ })
      return () => slots.default?.()[0]
    } } },
    '@/store/modules/locale': { useLocaleStore: () => ({ currentLocale: locale }) },
    '@vueuse/core': { useWindowSize: () => ({ width: vue.ref(1280) }) },
    '@/store/modules/app': { useAppStore: () => ({ setCssVarTheme() {}, setMobile() {}, getMobile: false }) },
    '@/utils': { setCssVar() {} },
    '@/hooks/web/useDesign': design
  }, true)
  const renderer = vue.createRenderer({
    createElement: () => ({ children: [] }), createText: () => ({}), createComment: () => ({}),
    setText() {}, setElementText() {}, parentNode: (node) => node.parent,
    nextSibling: (node) => node.parent?.children[node.parent.children.indexOf(node) + 1] || null,
    insert(node, parent, anchor) {
      if (node.parent) node.parent.children.splice(node.parent.children.indexOf(node), 1)
      const at = anchor ? parent.children.indexOf(anchor) : -1
      parent.children.splice(at < 0 ? parent.children.length : at, 0, node)
      node.parent = parent
    },
    remove(node) { if (node.parent) node.parent.children.splice(node.parent.children.indexOf(node), 1); node.parent = null },
    patchProp() {}
  })
  let mounts = 0, unmounts = 0, draft, asyncCompleted = false
  let finishPending
  const pending = new Promise((resolve) => { finishPending = resolve })
  const probe = { setup() {
    const instance = ++mounts
    draft = vue.ref('')
    if (instance === 1) pending.then(() => { asyncCompleted = true })
    vue.onUnmounted(() => { unmounts++ })
    return () => vue.h('div', draft.value)
  } }
  const app = renderer.createApp({ setup: () => () => vue.h(config, {}, { default: () => vue.h(probe) }) })
  app.mount({ children: [] })
  draft.value = 'unsaved test input'
  const cancelRemount = setup.setLang('ko-KR')
  rejectConfirmation('cancel')
  await cancelRemount
  await vue.nextTick()
  assert.equal(mounts, 1)
  assert.equal(draft.value, 'unsaved test input')
  failImport = false
  const confirmedRemount = setup.setLang('ko-KR')
  resolveConfirmation()
  await confirmedRemount
  await vue.nextTick()
  assert.equal(mounts, 2)
  assert.equal(unmounts, 1)
  assert.equal(draft.value, '')
  assert.equal(providerMounts, 1, 'Element Plus provider must survive locale changes')
  assert.equal(providerUnmounts, 0)
  finishPending()
  await pending
  await vue.nextTick()
  assert.equal(asyncCompleted, true, 'Unmount does not cancel an existing async operation')
  app.unmount()
  console.log(JSON.stringify({
    titleUpdatesWithoutWatchers: 'PASS (101 calls, 0 watchers)',
    sameLanguageNoOp: 'PASS', confirmationCancelAndClose: 'PASS',
    overlappingSwitchGuard: 'PASS', confirmedSwitch: 'PASS',
    failedLoadUnlocksMenu: 'PASS', canceledSwitchKeepsDraftAndInstance: 'PASS',
    providerInstancePreserved: 'PASS',
    confirmedLocaleRecreatesSubtree: 'PASS (known reset behavior)',
    unmountDoesNotCancelPendingAsync: 'PASS (synthetic callback; no real API)'
  }, null, 2))
}
main().catch((error) => { console.error(error); process.exitCode = 1 })

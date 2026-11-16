<script setup lang="tsx">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
// 只引用类型；运行时通过动态 import 拉 monaco，避免资产列表等页面被动加载 ~1MB
import type * as Monaco from 'monaco-editor'
import { useAppStore } from '@/store/modules/app'

const props = defineProps<{ original: string; modified: string }>()

const appStore = useAppStore()

const isDark = computed(() => appStore.getIsDark)

const editorContainer = ref<HTMLDivElement | null>(null)

let monaco: typeof Monaco | null = null

let diffEditor: Monaco.editor.IStandaloneDiffEditor | null = null

const monacoTheme = () => (isDark.value ? 'vs-dark' : 'vs')

const loadMonaco = async () => {
  if (!monaco) {
    monaco = await import('monaco-editor')
  }
  return monaco
}

onMounted(async () => {
  const m = await loadMonaco()
  if (!editorContainer.value) return
  diffEditor = m.editor.createDiffEditor(editorContainer.value, {
    theme: monacoTheme(),
    originalEditable: true,
    automaticLayout: true
  })

  diffEditor.setModel({
    original: m.editor.createModel(props.original, 'javascript'),
    modified: m.editor.createModel(props.modified, 'javascript')
  })
})

// monaco 的主题是全局设置，切换后所有实例一起变
watch(isDark, () => {
  monaco?.editor.setTheme(monacoTheme())
})

onBeforeUnmount(() => {
  diffEditor?.dispose()
  diffEditor = null
})
</script>
<template>
  <div ref="editorContainer" class="monaco-diff-editor"></div>
</template>
<style scoped>
.monaco-diff-editor {
  width: 100%;
  height: 100%;
}
</style>

/**
 * Monaco 差异编辑器：跟随全局主题（vs / vs-dark）
 */
<script setup lang="tsx">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as monaco from 'monaco-editor'
import { useAppStore } from '@/store/modules/app'

const props = defineProps<{ original: string; modified: string }>()

const appStore = useAppStore()

const isDark = computed(() => appStore.getIsDark)

const editorContainer = ref<HTMLDivElement | null>(null)
let diffEditor: monaco.editor.IStandaloneDiffEditor | null = null

const monacoTheme = () => (isDark.value ? 'vs-dark' : 'vs')

onMounted(() => {
  if (editorContainer.value) {
    diffEditor = monaco.editor.createDiffEditor(editorContainer.value, {
      theme: monacoTheme(),
      originalEditable: true,
      automaticLayout: true
    })

    diffEditor.setModel({
      original: monaco.editor.createModel(props.original, 'javascript'),
      modified: monaco.editor.createModel(props.modified, 'javascript')
    })
  }
})

// monaco 的主题是全局设置，切换后所有实例一起变
watch(isDark, () => {
  monaco.editor.setTheme(monacoTheme())
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

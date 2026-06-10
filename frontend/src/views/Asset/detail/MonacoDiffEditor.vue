<script setup lang="tsx">
import { onMounted, ref } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps<{ original: string; modified: string }>()

const editorContainer = ref<HTMLDivElement | null>(null)
let diffEditor: monaco.editor.IStandaloneDiffEditor | null = null

onMounted(() => {
  if (editorContainer.value) {
    diffEditor = monaco.editor.createDiffEditor(editorContainer.value, {
      theme: 'vs-dark',
      originalEditable: true,
      automaticLayout: true
    })

    diffEditor.setModel({
      original: monaco.editor.createModel(props.original, 'javascript'),
      modified: monaco.editor.createModel(props.modified, 'javascript')
    })
  }
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

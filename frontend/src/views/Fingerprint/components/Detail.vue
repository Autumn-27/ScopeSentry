<script setup lang="ts">
import { ElButton } from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { toRefs } from '@vueuse/core'
import { updateFingerprintDataApi, addFingerprintDataApi } from '@/api/Fingerprint'
import * as monaco from 'monaco-editor'

const { t } = useI18n()
const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  fingerprintForm: {
    id: string
    content: string
  }
}>()
const { fingerprintForm } = toRefs(props)
const localFingerprintForm = ref({ ...fingerprintForm.value })

const editorContainer = ref<HTMLDivElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

// 默认 YAML 模板
const defaultYamlTemplate = `fingerprint:
  name: 
  category: 
  parent_category: 
  company: 
  tags:
    - "demo"
  rules:
  - logic: AND
    conditions:
    - location: body
      match_type: contains
      pattern: "demo"
`

onMounted(() => {
  if (editorContainer.value) {
    // 如果是新建且内容为空，使用默认模板
    let initialValue = localFingerprintForm.value.content || ''
    if (!localFingerprintForm.value.id && !initialValue.trim()) {
      initialValue = defaultYamlTemplate
      localFingerprintForm.value.content = initialValue
    }

    editor = monaco.editor.create(editorContainer.value, {
      value: initialValue,
      language: 'yaml',
      theme: 'vs-dark',
      automaticLayout: true,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      fontSize: 14,
      lineNumbers: 'on',
      wordWrap: 'on'
    })

    // 监听编辑器内容变化
    editor.onDidChangeModelContent(() => {
      const value = editor?.getValue() || ''
      localFingerprintForm.value.content = value
    })
  }
})

// 监听外部数据变化
watch(
  () => [props.fingerprintForm.content, props.fingerprintForm.id],
  ([newValue, id]) => {
    if (editor) {
      // 如果是新建且内容为空，使用默认模板
      if (!id && (!newValue || !newValue.trim())) {
        const defaultValue = defaultYamlTemplate
        if (editor.getValue() !== defaultValue) {
          editor.setValue(defaultValue)
          localFingerprintForm.value.content = defaultValue
        }
      } else if (newValue !== editor.getValue()) {
        editor.setValue(newValue || '')
        localFingerprintForm.value.content = newValue || ''
      }
    }
  }
)

onBeforeUnmount(() => {
  if (editor) {
    editor.dispose()
  }
})

const saveLoading = ref(false)
const submitForm = async () => {
  saveLoading.value = true
  try {
    const content = editor?.getValue() || ''
    if (!content.trim()) {
      saveLoading.value = false
      return
    }

    let res
    if (localFingerprintForm.value.id != '') {
      res = await updateFingerprintDataApi(localFingerprintForm.value.id, content)
    } else {
      res = await addFingerprintDataApi(content)
    }
    if (res.code === 200) {
      props.getList()
      props.closeDialog()
    }
  } finally {
    saveLoading.value = false
  }
}
</script>
<template>
  <div class="fingerprint-detail-container">
    <div class="editor-wrapper">
      <div ref="editorContainer" class="yaml-editor"></div>
    </div>
    <div class="action-bar">
      <ElButton @click="props.closeDialog">{{ t('common.cancel') || '取消' }}</ElButton>
      <ElButton type="primary" @click="submitForm" :loading="saveLoading">{{
        t('task.save') || '保存'
      }}</ElButton>
    </div>
  </div>
</template>
<style scoped>
.fingerprint-detail-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.editor-wrapper {
  flex: 1;
  min-height: 0;
  margin-bottom: 0;
  display: flex;
  flex-direction: column;
}

.yaml-editor {
  width: 100%;
  height: 100%;
  min-height: 350px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid var(--el-border-color);
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 4px 0;
  border-top: 1px solid var(--el-border-color-light);
  margin-top: 8px;
  flex-shrink: 0;
}

/* 减少 Dialog body 的 padding */
:deep(.el-dialog__body) {
  padding: 10px 15px !important;
}

/* Monaco Editor 样式优化 */
:deep(.monaco-editor) {
  border-radius: 8px;
}

:deep(.monaco-editor .margin) {
  background-color: #1e1e1e !important;
}

:deep(.monaco-editor .monaco-editor-background) {
  background-color: #1e1e1e !important;
}

:deep(.monaco-scrollable-element > .scrollbar) {
  background-color: transparent !important;
}
</style>

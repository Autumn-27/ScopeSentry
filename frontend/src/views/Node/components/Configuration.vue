<script setup lang="ts">
import {
  ElFormItem,
  ElInput,
  ElRow,
  ElCol,
  FormInstance,
  ElForm,
  ElButton,
  ElSwitch,
  ElDivider,
  ElText
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { onMounted, ref } from 'vue'
import { toRefs } from '@vueuse/core'
import { updateNodeConfigDataApi } from '@/api/node'
import { Codemirror } from 'vue-codemirror'
import { javascript } from '@codemirror/lang-javascript'
import { oneDark } from '@codemirror/theme-one-dark'
const extensions = [javascript(), oneDark]
const { t } = useI18n()
const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  nodeConfForm: {
    name: string
    maxTaskNum: string
    state: string
    ModulesConfig: string
  }
}>()
const { nodeConfForm } = toRefs(props)
const localForm = ref({ ...nodeConfForm.value })
const saveLoading = ref(false)
const isDisabled = ref(false)
const switchValue = ref(false)
onMounted(() => {
  if (localForm.value.state === '1') {
    switchValue.value = true
    isDisabled.value = false
  } else if (localForm.value.state === '2') {
    switchValue.value = false
    isDisabled.value = false
  } else if (localForm.value.state === '3') {
    switchValue.value = false
    isDisabled.value = true
  }
})
const ruleFormRef = ref<FormInstance>()

const oldName = localForm.value.name
const submitForm = async (formEl: FormInstance | undefined) => {
  saveLoading.value = true
  if (!formEl) return
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      let res
      try {
        res = await updateNodeConfigDataApi(
          oldName,
          localForm.value.name,
          localForm.value.ModulesConfig,
          switchValue.value
        )
        if (res.code === 200) {
          props.getList()
          props.closeDialog()
        }
      } finally {
        saveLoading.value = false
      }
    } else {
      console.log('error submit!', fields)
      saveLoading.value = false
    }
  })
}
</script>
<template>
  <ElForm :model="localForm" label-width="auto" status-icon ref="ruleFormRef">
    <ElFormItem :label="t('node.nodeName')" prop="name">
      <ElInput v-model="localForm.name" />
    </ElFormItem>
    <ElFormItem label="Module Config">
      <Codemirror
        v-model="localForm.ModulesConfig"
        :extensions="extensions"
        :autofocus="true"
        :indent-with-tab="true"
        :tab-size="2"
        :style="{ height: '550px', width: '100%' }"
      />
    </ElFormItem>
    <ElFormItem :label="t('common.state')">
      <ElSwitch
        v-model="switchValue"
        inline-prompt
        :active-text="t('common.switchAction')"
        :inactive-text="t('common.switchInactive')"
        :disabled="isDisabled"
      />
    </ElFormItem>
    <ElRow>
      <ElCol :span="16" :offset="8">
        <ElFormItem>
          <ElButton type="primary" @click="submitForm(ruleFormRef)" :loading="saveLoading">
            {{ t('task.save') }}
          </ElButton>
          <ElDivider direction="vertical" />
          <ElText size="small" type="danger">{{ t('configuration.threadMsg') }}</ElText>
        </ElFormItem>
      </ElCol>
    </ElRow>
  </ElForm>
</template>

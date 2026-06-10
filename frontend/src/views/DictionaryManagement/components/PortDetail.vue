<script setup lang="ts">
import {
  ElFormItem,
  ElInput,
  ElRow,
  ElCol,
  FormRules,
  FormInstance,
  ElForm,
  ElButton,
  ElDivider
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { reactive, ref } from 'vue'
import { toRefs } from '@vueuse/core'
import { upgradePortDictDataApi, addPortDictDataApi } from '@/api/DictionaryManagement'
const { t } = useI18n()
const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  portDictForm: {
    id: string
    name: string
    value: string
  }
}>()
const { portDictForm } = toRefs(props)
const localSensitiveForm = ref({ ...portDictForm.value })

interface RuleForm {
  name: string
  regular: string
}
const rules = reactive<FormRules<RuleForm>>({
  name: [{ required: true, message: t('portDict.nameMsg'), trigger: 'blur' }],
  regular: [{ required: true, message: t('portDict.valueMsg'), trigger: 'blur' }]
})
const saveLoading = ref(false)
const ruleFormRef = ref<FormInstance>()
const submitForm = async (formEl: FormInstance | undefined) => {
  saveLoading.value = true
  if (!formEl) return
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      let res
      console.log('submit!')
      if (localSensitiveForm.value.id != '') {
        res = await upgradePortDictDataApi(
          localSensitiveForm.value.id,
          localSensitiveForm.value.name,
          localSensitiveForm.value.value
        )
      } else {
        res = await addPortDictDataApi(
          localSensitiveForm.value.name,
          localSensitiveForm.value.value
        )
      }
      if (res.code === 200) {
        props.getList()
        props.closeDialog()
      }
      saveLoading.value = false
    } else {
      console.log('error submit!', fields)
      saveLoading.value = false
    }
  })
}
</script>
<template>
  <ElForm
    :model="localSensitiveForm"
    label-width="auto"
    :rules="rules"
    status-icon
    ref="ruleFormRef"
  >
    <ElFormItem :label="t('portDict.name')" prop="name">
      <ElInput v-model="localSensitiveForm.name" :placeholder="t('portDict.nameMsg')" />
    </ElFormItem>
    <ElFormItem :label="t('portDict.value')">
      <ElInput
        v-model="localSensitiveForm.value"
        type="textarea"
        prop="value"
        :placeholder="t('portDict.valueMsg')"
        :autosize="{ minRows: 11 }"
      />
    </ElFormItem>
    <ElDivider />
    <ElRow>
      <ElCol :span="2" :offset="8">
        <ElFormItem>
          <ElButton type="primary" @click="submitForm(ruleFormRef)" :loading="saveLoading">{{
            t('task.save')
          }}</ElButton>
        </ElFormItem>
      </ElCol>
    </ElRow>
  </ElForm>
</template>

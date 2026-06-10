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
  ElDivider,
  ElSelectV2,
  ElSwitch
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { reactive, ref } from 'vue'
import { toRefs } from '@vueuse/core'
import { updateSensitiveDataApi, addSensitiveDataApi } from '@/api/sensitive'
const { t } = useI18n()
const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  sensitiveForm: {
    id: string
    name: string
    regular: string
    color: string
    state: boolean
  }
}>()
const { sensitiveForm } = toRefs(props)
const localSensitiveForm = ref({ ...sensitiveForm.value })

interface RuleForm {
  name: string
  regular: string
}
const rules = reactive<FormRules<RuleForm>>({
  name: [{ required: true, message: t('sensitiveInformation.sensitiveNameMsg'), trigger: 'blur' }],
  regular: [
    { required: true, message: t('sensitiveInformation.sensitiveRegularMsg'), trigger: 'blur' }
  ]
})
const colorOptions = [
  {
    value: 'null',
    label: 'null'
  },
  {
    value: 'green',
    label: 'green'
  },
  {
    value: 'red',
    label: 'red'
  },
  {
    value: 'cyan',
    label: 'cyan'
  },
  {
    value: 'yellow',
    label: 'yellow'
  },
  {
    value: 'orange',
    label: 'orange'
  },
  {
    value: 'gray',
    label: 'gray'
  },
  {
    value: 'pink',
    label: 'pink'
  }
]
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
        res = await updateSensitiveDataApi(
          localSensitiveForm.value.id,
          localSensitiveForm.value.name,
          localSensitiveForm.value.regular,
          localSensitiveForm.value.color,
          localSensitiveForm.value.state
        )
      } else {
        res = await addSensitiveDataApi(
          localSensitiveForm.value.name,
          localSensitiveForm.value.regular,
          localSensitiveForm.value.color,
          localSensitiveForm.value.state
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
    <ElFormItem :label="t('sensitiveInformation.sensitiveName')" prop="name">
      <ElInput
        v-model="localSensitiveForm.name"
        :placeholder="t('sensitiveInformation.sensitiveNameMsg')"
      />
    </ElFormItem>
    <ElFormItem :label="t('sensitiveInformation.sensitiveRegular')" prop="regular">
      <ElInput
        v-model="localSensitiveForm.regular"
        :placeholder="t('sensitiveInformation.sensitiveRegularMsg')"
      />
    </ElFormItem>
    <ElFormItem :label="t('sensitiveInformation.sensitiveColor')">
      <ElSelectV2
        v-model="localSensitiveForm.color"
        placeholder="Please select color"
        :options="colorOptions"
      />
    </ElFormItem>
    <ElFormItem :label="t('common.state')">
      <ElSwitch
        v-model="localSensitiveForm.state"
        inline-prompt
        :active-text="t('common.switchAction')"
        :inactive-text="t('common.switchInactive')"
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

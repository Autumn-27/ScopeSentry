<script setup lang="ts">
import { useI18n } from '@/hooks/web/useI18n'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Dialog } from '@/components/Dialog'
import { Form } from '@/components/Form'
import { useForm } from '@/hooks/web/useForm'
import { reactive, computed } from 'vue'
import { useValidator } from '@/hooks/web/useValidator'
import { FormSchema } from '@/components/Form'
import { useDesign } from '@/hooks/web/useDesign'
import { changePasswordApi } from '@/api/login'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const { getPrefixCls } = useDesign()
const prefixCls = getPrefixCls('lock-dialog')
const { required } = useValidator()

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Boolean
  }
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  }
})

const dialogTitle = ref(t('common.changePassword'))

const rules = reactive({
  newPassword: [required()]
})

const schema: FormSchema[] = reactive([
  {
    label: t('common.newPassword'),
    field: 'newPassword',
    component: 'Input',
    componentProps: {
      type: 'password'
    }
  }
])

const { formRegister, formMethods } = useForm()

const { getFormData, getElFormExpose } = formMethods

export interface changePassword {
  newPassword: string
}
const changePassword = async () => {
  const formExpose = await getElFormExpose()
  formExpose?.validate(async (valid) => {
    if (valid) {
      dialogVisible.value = false
      const formData = await getFormData()
      console.log(formData)
      const changePasswordData: changePassword = {
        newPassword: formData.newPassword
      }
      const res = await changePasswordApi(changePasswordData)
      console.log(res)
      if (res.code == 200) {
        ElMessage.success(res.data.message)
      }
    }
  })
}
const userInfo = userStore.getUserInfo
const username = userInfo!.username
</script>

<template>
  <Dialog
    v-model="dialogVisible"
    width="500px"
    max-height="170px"
    :class="prefixCls"
    :title="dialogTitle"
  >
    <div class="flex flex-col items-center">
      <img src="@/assets/imgs/avatar.jpg" alt="" class="w-70px h-70px rounded-[50%]" />
      <span class="text-14px my-10px text-[var(--top-header-text-color)]"> {{ username }} </span>
    </div>
    <Form :is-col="false" :schema="schema" :rules="rules" @register="formRegister" />
    <template #footer>
      <BaseButton type="primary" @click="changePassword">{{ t('common.submit') }}</BaseButton>
    </template>
  </Dialog>
</template>

<style lang="less" scoped>
:global(.v-lock-dialog) {
  @media (width <= 767px) {
    max-width: calc(100vw - 16px);
  }
}
</style>

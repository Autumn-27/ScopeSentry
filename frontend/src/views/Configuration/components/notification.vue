<script setup lang="tsx">
import {
  ElRow,
  ElCol,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElRadioGroup,
  ElRadio,
  ElText,
  ElSwitch,
  ElDivider
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { ElCard } from 'element-plus'
import { h, reactive, ref } from 'vue'
import { Icon } from '@/components/Icon'
import { Dialog } from '@/components/Dialog'
import { Table, TableColumn } from '@/components/Table'
import { BaseButton } from '@/components/Button'
import { useTable } from '@/hooks/web/useTable'
import {
  addNotificationApi,
  deletePocDataApi,
  getNotificationApi,
  getNotificationConfigApi,
  updateNotificationApi,
  updateNotificationConfigApi
} from '@/api/Configuration'
const { t } = useI18n()
const taskColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'name',
    label: 'Name',
    minWidth: 20
  },
  {
    field: 'method',
    label: 'Method',
    minWidth: 20
  },
  {
    field: 'url',
    label: 'URL'
  },
  {
    field: 'contentType',
    label: 'Content Type',
    minWidth: 25
  },
  {
    field: 'data',
    label: 'POST DATA'
  },
  {
    field: 'state',
    label: t('common.state'),
    minWidth: 25,
    formatter: (_: Recordable, __: TableColumn, stateValue: boolean) => {
      let color = ''
      let flag = ''
      if (stateValue == true) {
        color = '#2eb98a'
        flag = t('common.on')
      } else {
        color = 'red'
        flag = t('common.off')
      }
      return h(ElRow, { gutter: 20 }, [
        h(ElCol, { span: 1 }, [h(Icon, { icon: 'clarity:circle-solid', color, size: 10 })]),
        h(ElCol, { span: 5 }, [h(ElText, { type: 'info' }, flag)])
      ])
    }
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    formatter: (row, __: TableColumn, _: number) => {
      return h('div', [
        h(BaseButton, { type: 'primary', onClick: () => edit(row) }, t('common.edit')),
        h(BaseButton, { type: 'danger', onClick: () => del(row) }, t('common.delete'))
      ])
    }
  }
])

const { tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getNotificationApi()
    return {
      list: res.data.list
    }
  }
})
const { dataList } = tableState
const { getList, getElTableExpose } = tableMethods
const notificationForm = reactive({
  name: '',
  url: '',
  method: 'GET',
  contentType: 'raw',
  data: '',
  state: true
})
const notificationConfigForm = reactive({
  dirScanNotification: true,
  portScanNotification: true,
  sensitiveNotification: true,
  subdomainNotification: true,
  subdomainTakeoverNotification: true,
  pageMonNotification: true,
  vulNotification: true
})
const getNotificationConfig = async () => {
  const res = await getNotificationConfigApi()
  console.log(res)
  notificationConfigForm.dirScanNotification = res.data.dirScanNotification
  notificationConfigForm.portScanNotification = res.data.portScanNotification
  notificationConfigForm.sensitiveNotification = res.data.sensitiveNotification
  notificationConfigForm.subdomainNotification = res.data.subdomainNotification
  notificationConfigForm.subdomainTakeoverNotification = res.data.subdomainTakeoverNotification
  notificationConfigForm.pageMonNotification = res.data.pageMonNotification
  notificationConfigForm.vulNotification = res.data.vulNotification
}
getNotificationConfig()
const updateNotificationSaveLoading = ref(false)
const updateNotificationConfig = async () => {
  updateNotificationSaveLoading.value = true
  await updateNotificationConfigApi(
    notificationConfigForm.dirScanNotification,
    notificationConfigForm.portScanNotification,
    notificationConfigForm.sensitiveNotification,
    notificationConfigForm.subdomainNotification,
    notificationConfigForm.subdomainTakeoverNotification,
    notificationConfigForm.pageMonNotification,
    notificationConfigForm.vulNotification
  )
  updateNotificationSaveLoading.value = false
}
const dialogVisible = ref(false)
const addNotification = async () => {
  notificationId.value = ''
  notificationForm.name = ''
  notificationForm.url = ''
  notificationForm.method = 'GET'
  notificationForm.contentType = 'raw'
  notificationForm.data = ''
  notificationForm.state = true
  dialogVisible.value = true
}
const notificationId = ref('')
const addNotificationSaveLoading = ref(false)
const submitAddPageMonitForm = async () => {
  console.log(notificationId.value)
  addNotificationSaveLoading.value = true
  // await addScheduledTaskPageMonitApi(pageMontForm.url)
  try {
    if (notificationId.value == '') {
      await addNotificationApi(
        notificationForm.name,
        notificationForm.url,
        notificationForm.method,
        notificationForm.contentType,
        notificationForm.data,
        notificationForm.state
      )
    } else {
      await updateNotificationApi(
        notificationId.value,
        notificationForm.name,
        notificationForm.url,
        notificationForm.method,
        notificationForm.contentType,
        notificationForm.data,
        notificationForm.state
      )
    }
    getList()
    addNotificationSaveLoading.value = false
  } finally {
    dialogVisible.value = false
  }
}
const edit = (data) => {
  notificationId.value = data.id
  notificationForm.name = data.name
  notificationForm.url = data.url
  notificationForm.method = data.method
  notificationForm.contentType = data.contentType
  notificationForm.data = data.data
  notificationForm.state = data.state
  dialogVisible.value = true
  console.log(notificationId.value)
}
const delLoading = ref(false)
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await deletePocDataApi([data.id])
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const ids = ref<string[]>([])

const confirmDelete = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}

const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await deletePocDataApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
</script>

<template>
  <ElCard shadow="never" class="mb-20px">
    <template #header>
      <ElRow>
        <ElCol :span="3" style="height: 100%">
          <span>{{ t('configuration.noticeConfig') }}</span>
        </ElCol>
      </ElRow>
    </template>
    <ElRow>
      <ElCol style="position: relative; top: 16px">
        <div class="mb-10px">
          <BaseButton type="primary" @click="addNotification">{{
            t('configuration.newWebhookConfig')
          }}</BaseButton>
          <BaseButton type="danger" :loading="delLoading" @click="confirmDelete">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </ElCol>
    </ElRow>
    <div style="position: relative; top: 12px">
      <Table
        :data="dataList"
        :columns="taskColums"
        stripe
        :border="true"
        :resizable="true"
        maxHeight="200"
        :style="{
          fontFamily:
            '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
        }"
      />
    </div>
    <ElDivider />
    <ElForm
      :model="notificationConfigForm"
      label-width="auto"
      status-icon
      ref="ruleFormRef"
      style="position: relative; top: 1rem"
    >
      <ElRow>
        <ElCol :span="5">
          <ElFormItem :label="t('subdomain.subdomainName')">
            <ElSwitch
              v-model="notificationConfigForm.subdomainNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="5">
          <ElFormItem :label="t('task.subdomainTakeover')">
            <ElSwitch
              v-model="notificationConfigForm.subdomainTakeoverNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="5">
          <ElFormItem :label="t('dirScan.dirScanName')">
            <ElSwitch
              v-model="notificationConfigForm.dirScanNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
      </ElRow>
      <ElRow>
        <ElCol :span="5">
          <ElFormItem :label="t('task.portScan')">
            <ElSwitch
              v-model="notificationConfigForm.portScanNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="5">
          <ElFormItem :label="t('sensitiveInformation.sensitiveInformationName')">
            <ElSwitch
              v-model="notificationConfigForm.sensitiveNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="5">
          <ElFormItem :label="t('PageMonitoring.pageMonitoringName')">
            <ElSwitch
              v-model="notificationConfigForm.pageMonNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
      </ElRow>
      <ElRow>
        <ElCol :span="5">
          <ElFormItem :label="t('vulnerability.vulnerabilityName')">
            <ElSwitch
              v-model="notificationConfigForm.vulNotification"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElFormItem>
        </ElCol>
      </ElRow>
      <ElRow>
        <ElCol :span="2" :offset="8">
          <ElFormItem>
            <ElButton
              type="primary"
              @click="updateNotificationConfig()"
              :loading="updateNotificationSaveLoading"
              >{{ t('common.submit') }}</ElButton
            >
          </ElFormItem>
        </ElCol>
      </ElRow>
    </ElForm>
  </ElCard>
  <Dialog
    v-model="dialogVisible"
    :title="t('configuration.newWebhookConfig')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    maxHeight="100"
  >
    <ElText class="mx-2" type="danger" size="small" style="position: relative; left: 2rem">{{
      t('configuration.noticeHelp')
    }}</ElText>
    <ElForm
      :model="notificationForm"
      label-width="auto"
      status-icon
      ref="ruleFormRef"
      style="position: relative; top: 1rem"
    >
      <ElFormItem label="Name" prop="name">
        <ElInput v-model="notificationForm.name" placeholder="Input name." />
      </ElFormItem>
      <ElFormItem label="Method" prop="method">
        <ElRadioGroup v-model="notificationForm.method">
          <ElRadio value="GET">GET</ElRadio>
          <ElRadio value="POST">POST</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="URL" prop="url">
        <ElInput v-model="notificationForm.url" placeholder="Input URL." />
      </ElFormItem>
      <ElFormItem label="Data Type" prop="contentType" v-if="notificationForm.method == 'POST'">
        <ElRadioGroup v-model="notificationForm.contentType">
          <ElRadio value="raw">Raw</ElRadio>
          <ElRadio value="json">Json</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem label="Data" prop="Data" v-if="notificationForm.method == 'POST'">
        <ElInput v-model="notificationForm.data" placeholder="Input POST Data." />
      </ElFormItem>
      <ElFormItem :label="t('common.state')">
        <ElSwitch
          v-model="notificationForm.state"
          inline-prompt
          :active-text="t('common.switchAction')"
          :inactive-text="t('common.switchInactive')"
        />
      </ElFormItem>
      <ElRow>
        <ElCol :span="2" :offset="8">
          <ElFormItem>
            <ElButton
              type="primary"
              @click="submitAddPageMonitForm()"
              :loading="addNotificationSaveLoading"
              >{{ t('common.submit') }}</ElButton
            >
          </ElFormItem>
        </ElCol>
      </ElRow>
    </ElForm>
  </Dialog>
</template>

<style scoped>
.header-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>

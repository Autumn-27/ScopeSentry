<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, onMounted } from 'vue'
import { ElButton, ElCol, ElInput, ElRow, ElText, ElForm, ElFormItem } from 'element-plus'
import { Dialog } from '@/components/Dialog'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { useIcon } from '@/hooks/web/useIcon'
import {
  getScheduledTaskPageMonitDataApi,
  deleteScheduledTaskPageMonitApi,
  addScheduledTaskPageMonitApi
} from '@/api/task'
import { BaseButton } from '@/components/Button'

const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const search = ref('')
const handleSearch = () => {
  getList()
}
const taskColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'url',
    label: 'URL'
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 20,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          <BaseButton type="danger" onClick={() => confirmDelete(row)}>
            {t('common.delete')}
          </BaseButton>
        </>
      )
    }
  }
])

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const res = await getScheduledTaskPageMonitDataApi(
      search.value,
      currentPage.value,
      pageSize.value
    )
    return {
      list: res.data.list,
      total: res.data.total
    }
  },
  immediate: true
})
const { loading, dataList, total, currentPage, pageSize } = tableState
pageSize.value = 20
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}

const confirmDeleteSelect = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}

const confirmDelete = async (data) => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await del(data)
  }
}
const delLoading = ref(false)
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await deleteScheduledTaskPageMonitApi([data.id])
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
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await deleteScheduledTaskPageMonitApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})
const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.75
}
const dialogVisible = ref(false)
const addPageMonit = async () => {
  dialogVisible.value = true
}
const pageMontForm = reactive({
  url: ''
})
const addPageMonitSaveLoading = ref(false)
const submitAddPageMonitForm = async () => {
  addPageMonitSaveLoading.value = true
  await addScheduledTaskPageMonitApi(pageMontForm.url)
  addPageMonitSaveLoading.value = false
  dialogVisible.value = false
}
</script>

<template>
  <ContentWrap>
    <ElRow>
      <ElCol :span="1">
        <ElText class="mx-1" style="position: relative; top: 8px">URL:</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="search" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5" style="position: relative; left: 16px">
        <ElButton type="primary" :icon="searchicon" style="height: 100%" @click="handleSearch"
          >Search</ElButton
        >
      </ElCol>
    </ElRow>
    <ElRow>
      <ElCol style="position: relative; top: 16px">
        <div class="mb-10px">
          <BaseButton type="primary" @click="addPageMonit">{{ t('task.addURL') }}</BaseButton>
          <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">
            {{ t('task.delURL') }}
          </BaseButton>
        </div>
      </ElCol>
    </ElRow>
    <div style="position: relative; top: 12px">
      <Table
        :tooltip-options="{
          offset: 1,
          showArrow: false,
          effect: 'dark',
          enterable: false,
          showAfter: 0,
          popperOptions: {},
          popperClass: 'test',
          placement: 'bottom',
          hideAfter: 0,
          disabled: true
        }"
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
        :columns="taskColums"
        :data="dataList"
        stripe
        :border="true"
        :loading="loading"
        :max-height="maxHeight"
        :resizable="true"
        :pagination="{
          total: total,
          pageSizes: [20, 30, 50, 100, 200, 500, 1000]
        }"
        @register="tableRegister"
        :headerCellStyle="tableHeaderColor"
        :style="{
          fontFamily:
            '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
        }"
      />
    </div>
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="t('task.addPageMonitTask')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    maxHeight="100"
  >
    <ElForm :model="pageMontForm" label-width="auto" status-icon ref="ruleFormRef">
      <ElFormItem label="URL" prop="url">
        <ElInput v-model="pageMontForm.url" placeholder="Input URL. Eg: http(s)://xxx.com" />
      </ElFormItem>
      <ElRow>
        <ElCol :span="2" :offset="8">
          <ElFormItem>
            <ElButton
              type="primary"
              @click="submitAddPageMonitForm()"
              :loading="addPageMonitSaveLoading"
              >{{ t('common.submit') }}</ElButton
            >
          </ElFormItem>
        </ElCol>
      </ElRow>
    </ElForm></Dialog
  >
</template>

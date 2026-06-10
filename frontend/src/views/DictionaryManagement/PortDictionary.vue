<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive } from 'vue'
import { ElButton, ElCol, ElRow, ElInput } from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { Dialog } from '@/components/Dialog'
import { useIcon } from '@/hooks/web/useIcon'
import { useTable } from '@/hooks/web/useTable'
import { BaseButton } from '@/components/Button'
import { getPortDictDataApi, deletePortDictDataApi } from '@/api/DictionaryManagement'
import PortDetail from './components/PortDetail.vue'
const { t } = useI18n()
const searchicon = useIcon({ icon: 'iconoir:search' })
const search = ref('')
const handleSearch = () => {
  getList()
}
const nodeColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'id',
    hidden: true
  },
  {
    field: 'name',
    label: t('portDict.name'),
    minWidth: 40
  },
  {
    field: 'value',
    label: t('portDict.value')
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 40,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          <BaseButton type="primary" onClick={() => edit(row)}>
            {t('common.edit')}
          </BaseButton>
          <BaseButton type="danger" onClick={() => del(row)}>
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
    const res = await getPortDictDataApi(search.value, currentPage.value, pageSize.value)
    return {
      list: res.data.list,
      total: res.data.total
    }
  }
})
const { loading, dataList, total, currentPage, pageSize } = tableState
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
const dialogVisible = ref(false)
const addPortDict = async () => {
  portDictForm.id = ''
  portDictForm.value = ''
  portDictForm.name = ''
  dialogVisible.value = true
}
const closeDialog = () => {
  dialogVisible.value = false
}
let portDictForm = reactive({
  id: '',
  name: '',
  value: ''
})
const edit = (data) => {
  portDictForm.id = data.id
  portDictForm.value = data.value
  portDictForm.name = data.name
  dialogVisible.value = true
}
const delLoading = ref(false)
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await deletePortDictDataApi([data.id])
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
    const res = await deletePortDictDataApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmDelete = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}
</script>

<template>
  <ContentWrap>
    <ElRow :gutter="20" style="margin-bottom: 15px">
      <ElCol :span="1.5">
        <ElText class="mx-1" style="position: relative; top: 8px">Search :</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="search" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5">
        <ElButton type="primary" :icon="searchicon" style="height: 38px" @click="handleSearch"
          >Search</ElButton
        >
      </ElCol>
    </ElRow>
    <ElRow :gutter="60">
      <ElCol :span="1">
        <div class="mb-10px">
          <ElButton type="primary" @click="addPortDict">{{ t('common.new') }}</ElButton>
        </div>
      </ElCol>
      <ElCol :span="1">
        <div class="mb-10px">
          <BaseButton type="danger" :loading="delLoading" @click="confirmDelete">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </ElCol>
    </ElRow>
    <Table
      v-model:pageSize="pageSize"
      v-model:currentPage="currentPage"
      :columns="nodeColums"
      :data="dataList"
      stripe
      :border="true"
      :loading="loading"
      :resizable="true"
      :pagination="{
        total: total,
        pageSizes: [10, 20, 50, 100, 200, 500, 1000]
      }"
      @register="tableRegister"
      :headerCellStyle="tableHeaderColor"
      :style="{
        fontFamily:
          '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
      }"
    />
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="t('common.new')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="400"
  >
    <PortDetail :closeDialog="closeDialog" :portDictForm="portDictForm" :getList="getList" />
  </Dialog>
</template>

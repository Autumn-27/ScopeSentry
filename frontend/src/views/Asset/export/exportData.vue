<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import {
  ElTabs,
  ElTabPane,
  ElFormItem,
  ElForm,
  ElRadio,
  ElRadioGroup,
  ElInput,
  ElButton,
  ElTag,
  ElSpace,
  ElCheckbox,
  CheckboxValueType,
  ElCheckboxGroup
} from 'element-plus'
import { onMounted, reactive, Ref, ref } from 'vue'
import {
  exportApi,
  getExportRecordApi,
  delExportApi,
  getFieldApi,
  downloadExportApi
} from '@/api/export'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { BaseButton } from '@/components/Button'
const { t } = useI18n()

const props = defineProps<{
  index: string
  searchParams: string
  getFilter: () => { [key: string]: any }
}>()
const exportForm = reactive({
  type: 'all',
  quantity: 0
})
const create = async () => {
  createLoading.value = true
  const searchFilter = props.getFilter()
  await exportApi(
    props.index,
    exportForm.quantity,
    exportForm.type,
    props.searchParams,
    searchFilter,
    checkedField.value,
    filetype.value
  )
  createLoading.value = false
}

const exportColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection'
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    fixed: 'left',
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          <BaseButton type="success" onClick={() => download(row.file_name)}>
            {t('export.download')}
          </BaseButton>
          <BaseButton type="danger" onClick={() => confirmDelete(row)}>
            {t('common.delete')}
          </BaseButton>
        </>
      )
    }
  },
  {
    field: 'file_name',
    label: t('export.fileName')
  },
  {
    field: 'state',
    label: t('export.state'),
    formatter: (_: Recordable, __: TableColumn, value: number) => {
      if (value == 0) {
        return <ElTag type="info">{t('export.run')}</ElTag>
      } else if (value == 1) {
        return <ElTag type="success">{t('export.success')}</ElTag>
      } else {
        return <ElTag type="danger">{t('export.fail')}</ElTag>
      }
    }
  },
  {
    field: 'create_time',
    label: t('export.createTime')
  },
  {
    field: 'end_time',
    label: t('export.endTime'),
    formatter: (_: Recordable, __: TableColumn, value: string) => {
      if (value == '') {
        return '-'
      } else {
        return value
      }
    }
  },
  {
    field: 'data_type',
    label: t('export.type')
  },
  {
    field: 'file_size',
    label: t('export.fileSize'),
    formatter: (_: Recordable, __: TableColumn, value: string) => {
      if (value == '') {
        return '-'
      } else {
        return value + ' MB'
      }
    }
  }
])
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getExportRecordApi()
    return {
      list: res.data.list
    }
  },
  immediate: false
})
const { dataList, loading } = tableState
const { getList, getElTableExpose } = tableMethods

const download = async (id) => {
  try {
    const blob = await downloadExportApi(id)
    const url = window.URL.createObjectURL(blob)

    const a = document.createElement('a')
    a.href = url
    a.download = `${id}`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error('下载失败:', err)
  }
}
const createLoading = ref(false)
const onClick = (name) => {
  if (name == 'exportRecords') {
    getList()
  }
}
const delLoading = ref(false)
const confirmDelete = async (data) => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await del(data)
  }
}
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await delExportApi([data.file_name])
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmDeleteSelect = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}
const ids = ref<string[]>([])
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.file_name)
  delLoading.value = true
  try {
    const res = await delExportApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const fields = ref([]) as Ref<string[]>
const getField = async () => {
  const res = await getFieldApi(props.index)
  fields.value = res.data.field
}

onMounted(() => {
  getField()
})
const isIndeterminate = ref(true)
const checkAll = ref(false)
const checkedField = ref([]) as Ref<string[]>
const handleCheckAllChange = (val: CheckboxValueType) => {
  checkedField.value = val ? fields.value : []
  isIndeterminate.value = false
}
const handleCheckedCitiesChange = (value: CheckboxValueType[]) => {
  const checkedCount = value.length
  checkAll.value = checkedCount === fields.value.length
  isIndeterminate.value = checkedCount > 0 && checkedCount < fields.value.length
}
const filetype = ref('csv')
</script>

<template>
  <ElTabs tabPosition="left" @tab-change="onClick" model-value="export">
    <ElTabPane :label="t('asset.export')" name="export">
      <ElForm :model="exportForm" label-width="auto" style="position: relative">
        <ElFormItem :label="t('export.exportType')">
          <ElRadioGroup v-model="exportForm.type">
            <ElRadio value="all">{{ t('export.exportTypeAll') }}</ElRadio>
            <ElRadio value="search">{{ t('export.exportTypeSearch') }}</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem :label="t('export.exportQuantity')">
          <ElInput v-model.number="exportForm.quantity" />
        </ElFormItem>
        <ElFormItem :label="t('export.field')">
          <ElCheckbox
            v-model="checkAll"
            :indeterminate="isIndeterminate"
            @change="handleCheckAllChange"
          >
            All
          </ElCheckbox>
        </ElFormItem>
        <ElFormItem label=" ">
          <ElCheckboxGroup v-model="checkedField" @change="handleCheckedCitiesChange">
            <ElCheckbox v-for="field in fields" :key="field" :label="field" :value="field">
              {{ field }}
            </ElCheckbox>
          </ElCheckboxGroup>
        </ElFormItem>
        <ElFormItem :label="t('export.fileType')">
          <ElRadioGroup v-model="filetype">
            <ElRadio value="xlsx">csv</ElRadio>
            <ElRadio value="json">json</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem>
          <ElButton
            type="primary"
            @click="create"
            style="left: 40%; position: relative"
            :loading="createLoading"
          >
            Create
          </ElButton>
        </ElFormItem>
      </ElForm>
    </ElTabPane>
    <ElTabPane :label="t('export.exportRecords')" name="exportRecords">
      <ElSpace direction="vertical" alignment="flex-start" :style="{ width: '100%' }">
        <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">
          {{ t('common.delete') }}
        </BaseButton>
      </ElSpace>
      <Table
        @register="tableRegister"
        :columns="exportColums"
        :data="dataList"
        :loading="loading"
        max-height="500"
        :style="{
          fontFamily:
            '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji',
          width: '100%'
        }"
      />
    </ElTabPane>
  </ElTabs>
</template>

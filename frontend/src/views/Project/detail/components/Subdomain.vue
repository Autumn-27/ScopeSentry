<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { reactive, ref } from 'vue'
import { onMounted } from 'vue'
import { useTable } from '@/hooks/web/useTable'
import {
  ElCard,
  ElRow,
  ElCol,
  ElInput,
  ElMessageBox,
  ElMessage,
  ElButton,
  ElDivider,
  ElText
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import { useRoute } from 'vue-router'
import { getProjectSubdomainDataApi } from '@/api/ProjectAggregation'
import { delDataApi } from '@/api/asset'
const { t } = useI18n()
const { query } = useRoute()
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})

const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.8
}

const searchParams = ref('')
const handleSearch = (data: any) => {
  searchParams.value = data
  getList()
}
const filter = reactive<{ [key: string]: any }>({})
filter.project = [query.id as string]
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  getList()
}
const crudSchemas = reactive<CrudSchema[]>([
  {
    field: 'selection',
    type: 'selection',
    minWidth: '55'
  },
  {
    field: 'index',
    label: t('tableDemo.index'),
    type: 'index',
    minWidth: '30'
  },
  {
    field: 'host',
    label: t('subdomain.subdomainName'),
    minWidth: '200',
    formatter: (row, __: TableColumn, hostValue: string) => {
      if (!row.count) {
        return <ElText>{hostValue}</ElText>
      }
      return (
        <>
          <ElText>{hostValue}</ElText>
          <ElText type="info">({row.count})</ElText>
        </>
      )
    },
    slots: {
      header: () => {
        return (
          <div>
            <span>{t('subdomain.subdomainName')}</span>
            <ElInput
              v-model={hostValue.value}
              placeholder="Search"
              style="width: 200px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('sub_host')}
            />
          </div>
        )
      }
    }
  },
  {
    field: 'type',
    label: t('subdomain.recordType'),
    minWidth: '200',
    columnKey: 'type',
    filters: [
      { text: 'A', value: 'A' },
      { text: 'NS', value: 'NS' },
      { text: 'CNAME', value: 'CNAME' },
      { text: 'PTR', value: 'PTR' },
      { text: 'TXT', value: 'TXT' }
    ]
  },
  {
    field: 'value',
    label: t('subdomain.recordValue'),
    minWidth: '250',
    formatter: (_: Recordable, __: TableColumn, RecordValue: string[]) => {
      let content = ''
      RecordValue.forEach((item, _) => {
        content += `${item}\r\n`
      })
      return content
    },
    slots: {
      header: () => {
        return (
          <div>
            <span>{t('subdomain.recordValue')}</span>
            <ElInput
              v-model={valueValue.value}
              placeholder="Search"
              style="width: 200px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('sub_value')}
            />
          </div>
        )
      }
    }
  },
  {
    field: 'ip',
    label: 'IP',
    minWidth: '150',
    formatter: (_: Recordable, __: TableColumn, IPValue: string[]) => {
      let content = ''
      IPValue.forEach((item, _) => {
        content += `${item}\r\n`
      })
      return content
    },
    slots: {
      header: () => {
        return (
          <div>
            <span>IP</span>
            <ElInput
              v-model={ipValue.value}
              placeholder="Search"
              style="width: 180px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('sub_ip')}
            />
          </div>
        )
      }
    }
  },
  {
    field: 'time',
    label: t('asset.time'),
    minWidth: '200'
  }
])

const { allSchemas } = useCrudSchemas(crudSchemas)
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getProjectSubdomainDataApi('', filter, fq)
    return {
      list: res.data.list
    }
  },
  immediate: true
})
const { loading, dataList } = tableState
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}

const hostValue = ref('')
const valueValue = ref('')
const ipValue = ref('')
const fq = reactive<{ [key: string]: any }>({})
const filterSearchChange = async (type: string) => {
  let value = ''
  if (type == 'sub_host') {
    value = hostValue.value
  }
  if (type == 'sub_value') {
    value = valueValue.value
  }
  if (type == 'sub_ip') {
    value = ipValue.value
  }
  fq[type] = value
  getList()
}
const ids = ref<string[]>([])
const delSelect = async () => {
  console.log('dwa')
  ElMessageBox.confirm('Whether to delete?', 'Warning', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    type: 'warning'
  })
    .then(async () => {
      const elTableExpose = await getElTableExpose()
      const selectedRows = elTableExpose?.getSelectionRows() || []
      ids.value = selectedRows.map((row) => row.id)
      await delDataApi(ids.value, 'subdomain')
      getList()
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Delete canceled'
      })
    })
}
let deleteDisplay = ref(false)
const changeDeleteDisplay = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  if (ids.value.length != 0) {
    deleteDisplay.value = true
  } else {
    deleteDisplay.value = false
  }
}
</script>

<template>
  <ElRow>
    <ElCol>
      <ElCard style="height: min-content">
        <ElButton v-if="deleteDisplay" @click="delSelect" type="danger" size="small"
          >Dlete</ElButton
        >
        <Table
          :columns="allSchemas.tableColumns"
          :data="dataList"
          stripe
          :max-height="maxHeight"
          :border="true"
          :loading="loading"
          @selection-change="changeDeleteDisplay"
          rowKey="id"
          :resizable="true"
          @register="tableRegister"
          @filter-change="filterChange"
          :headerCellStyle="tableHeaderColor"
          :style="{
            fontFamily:
              '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
          }"
        />
      </ElCard>
    </ElCol>
  </ElRow>
</template>

<style lang="less" scoped>
.el-button {
  margin-top: 10px;
}
:deep(.el-table .cell.el-tooltip) {
  white-space: pre-line;
}
</style>

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
import { getProjectServiceDataApi } from '@/api/ProjectAggregation'
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
    field: 'service',
    label: t('asset.service'),
    minWidth: '100',
    formatter: (row, __: TableColumn, domainValue: string) => {
      if (!row.count) {
        return <ElText>{domainValue}</ElText>
      }
      return (
        <>
          <ElText>{domainValue}</ElText>
          <ElText type="info">({row.count})</ElText>
        </>
      )
    },
    slots: {
      header: () => {
        return (
          <div>
            <span>{t('asset.service')}</span>
            <ElInput
              v-model={serviceValue.value}
              placeholder="Search"
              style="width: 80px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('service_service')}
            />
          </div>
        )
      }
    }
  },
  {
    field: 'host',
    label: t('asset.domain'),
    minWidth: '200',
    slots: {
      header: () => {
        return (
          <div>
            <span>{t('asset.domain')}</span>
            <ElInput
              v-model={domainValue.value}
              placeholder="Search"
              style="width: 80px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('service_domain')}
            />
          </div>
        )
      }
    }
  },
  {
    field: 'ip',
    label: 'IP',
    minWidth: '250',
    slots: {
      header: () => {
        return (
          <div>
            <span>IP</span>
            <ElInput
              v-model={ipValue.value}
              placeholder="Search"
              style="width: 200px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('service_ip')}
            />
          </div>
        )
      }
    }
  },
  {
    field: 'port',
    label: t('asset.port'),
    minWidth: '250',
    slots: {
      header: () => {
        return (
          <div>
            <span>{t('asset.port')}</span>
            <ElInput
              v-model={portValue.value}
              placeholder="Search"
              style="width: 200px; margin-left: 10px;"
              size="small"
              onChange={() => filterSearchChange('service_port')}
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
    const res = await getProjectServiceDataApi('', filter, fq)
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
const portValue = ref('')
const domainValue = ref('')
const ipValue = ref('')
const serviceValue = ref('')
const fq = reactive<{ [key: string]: any }>({})
const filterSearchChange = async (type: string) => {
  let value = ''
  if (type == 'service_port') {
    value = portValue.value
  }
  if (type == 'service_domain') {
    value = domainValue.value
  }
  if (type == 'service_ip') {
    value = ipValue.value
  }
  if (type == 'service_service') {
    value = serviceValue.value
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
      await delDataApi(ids.value, 'asset')
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

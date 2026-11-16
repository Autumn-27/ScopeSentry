<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { h, nextTick, reactive, Ref, ref } from 'vue'
import { onMounted, onUnmounted } from 'vue'
import { useTable } from '@/hooks/web/useTable'
import { ElCard, ElPagination, ElRow, ElCol, ElTag } from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import { addTagApi, deleteTagApi, getIPAssetApi } from '@/api/asset'
import Csearch from '../search/Csearch.vue'
const { t } = useI18n()
interface Project {
  value: string
  label: string
  children?: Project[]
}
const props = defineProps<{
  projectList: Project[]
  taskList: { id: string; name: string }[]
}>()
const searchKeywordsData = [
  {
    keyword: 'domain',
    example: 'domain="baidu.com"',
    explain: t('searchHelp.name')
  },
  {
    keyword: 'ip',
    example: 'ip="192.168.2.1"',
    explain: t('searchHelp.ip')
  },
  {
    keyword: 'port',
    example: 'port="3306"',
    explain: t('searchHelp.port')
  },
  {
    keyword: 'app',
    example: 'app="Nginx"',
    explain: t('searchHelp.app')
  },
  {
    keyword: 'service',
    example: 'service="ssh"',
    explain: t('searchHelp.protocol')
  },
  {
    keyword: 'webServer',
    example: 'webServer="Flask"',
    explain: 'Search Web Server'
  },
  {
    keyword: 'project',
    example: 'project="Hackerone"',
    explain: t('searchHelp.project')
  }
]
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', setMaxHeight)
})

const maxHeight = ref(0)
const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.7
}

const searchParams = ref('')
const handleSearch = (data: any) => {
  searchParams.value = data
  getList()
}
const crudSchemas = reactive<CrudSchema[]>([
  {
    field: 'selection',
    type: 'selection',
    minWidth: '35'
  },
  {
    field: 'ip',
    label: t('asset.IP'),
    minWidth: '120'
  },
  {
    field: 'port',
    label: t('asset.port'),
    minWidth: '100'
  },
  {
    field: 'domain',
    label: t('asset.domain'),
    minWidth: '100'
  },
  {
    field: 'service',
    label: t('asset.service'),
    minWidth: '100'
  },
  {
    field: 'products',
    label: t('asset.products'),
    minWidth: '180',
    formatter: (_: Recordable, __: TableColumn, ProductsValue: string[] | null) => {
      if (!ProductsValue || ProductsValue.length === 0) return
      if (ProductsValue.length != 0) {
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
            {ProductsValue.map((product) => (
              <div key={product} style={{ cursor: 'pointer' }}>
                <ElTag type={'success'}>{product}</ElTag>
              </div>
            ))}
          </div>
        )
      }
    }
  },
  {
    field: 'time',
    label: t('asset.time'),
    minWidth: '180'
  }
])
let index = 'IPAsset'
crudSchemas.forEach((schema) => {
  schema.hidden = schema.hidden ?? false // 如果没有 hidden 属性，添加并设置为 false
})
let statisticsHidden = ref(false)
// 从localStorage读取配置并更新列的显示状态
const loadColumnConfig = () => {
  const savedConfig = JSON.parse(localStorage.getItem(`columnConfig_${index}`) || '{}')
  console.log(savedConfig)
  crudSchemas.forEach((col) => {
    if (savedConfig[col.field] !== undefined && col.field != 'select') {
      col.hidden = savedConfig[col.field] // 恢复列的显示状态
    }
  })
  statisticsHidden.value = savedConfig['statisticsHidden']
}
// 保存配置到localStorage
const saveColumnConfig = () => {
  const config = crudSchemas.reduce((acc, col) => {
    acc[col.field] = col.hidden // 保存每列的显示状态
    return acc
  }, {})
  config['statisticsHidden'] = statisticsHidden.value
  localStorage.setItem(`columnConfig_${index}`, JSON.stringify(config)) // 按index保存配置
}

// 处理列显示状态变化
const handleColumnVisibilityChange = ({ field, hidden }) => {
  console.log(field, hidden)
  const columnIndex = crudSchemas.findIndex((col) => col.field === field)
  if (columnIndex !== -1) {
    // 使用对象的展开运算符来创建一个新的对象，并更新隐藏属性
    crudSchemas[columnIndex].hidden = hidden
  }
  saveColumnConfig()
}
loadColumnConfig()
const { allSchemas } = useCrudSchemas(crudSchemas)
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const res = await getIPAssetApi(searchParams.value, currentPage.value, pageSize.value, filter)
    return {
      list: res.data.list,
      total: res.data.total
    }
  },
  immediate: false
})
const { loading, dataList, total, currentPage, pageSize } = tableState
// 单独设置 IP 组件的 pageSize 默认值
pageSize.value = 10 // 可以根据需要修改为你想要的默认值
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--app-surface-2)' }
}
const filter = reactive<{ [key: string]: any }>({})
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  getList()
}
const handleFilterSearch = (data: any, newFilters: any) => {
  Object.assign(filter, newFilters)
  searchParams.value = data
  getList()
}
const dynamicTags = ref<string[]>([])

const handleClose = (tag: string) => {
  if (dynamicTags.value) {
    const [key, value] = tag.split('=')
    if (key in filter && Array.isArray(filter[key])) {
      filter[key] = filter[key].filter((item: string) => item !== value)
      if (filter[key].length === 0) {
        delete filter[key]
      }
    }
    dynamicTags.value = dynamicTags.value.filter((item) => item !== tag)
  }
}
const getFilter = () => {
  return filter
}
const spanMethod = ({ row, column, rowIndex, columnIndex }) => {
  // columnIndex:
  // 0 = selection (复选框列)
  // 1 = IP
  // 2 = Port

  // 📌 Selection 列合并（与 IP 列同步）
  if (columnIndex === 0) {
    if (row.ipRowSpan > 0) {
      return [row.ipRowSpan, 1] // 显示并合并
    } else {
      return [0, 0] // 隐藏单元格
    }
  }

  // 📌 IP 合并
  if (columnIndex === 1) {
    if (row.ipRowSpan > 0) {
      return [row.ipRowSpan, 1] // 显示并合并
    } else {
      return [0, 0] // 隐藏单元格
    }
  }

  // 📌 Port 合并
  if (columnIndex === 2) {
    if (row.portRowSpan > 0) {
      return [row.portRowSpan, 1]
    } else {
      return [0, 0]
    }
  }

  // 📌 默认正常显示
  return [1, 1]
}
</script>

<template>
  <Csearch
    :getList="getList"
    :handleSearch="handleSearch"
    :searchKeywordsData="searchKeywordsData"
    :index="index"
    :getElTableExpose="getElTableExpose"
    :projectList="$props.projectList"
    :taskList="$props.taskList"
    :handleFilterSearch="handleFilterSearch"
    :crudSchemas="crudSchemas"
    :dynamicTags="dynamicTags"
    :handleClose="handleClose"
    @update-column-visibility="handleColumnVisibilityChange"
    :searchResultCount="total"
    :getFilter="getFilter"
  />
  <ElRow>
    <ElCol>
      <ElCard style="height: min-content">
        <Table
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
          :columns="allSchemas.tableColumns"
          :data="dataList"
          rowKey="datakey"
          stripe
          :border="true"
          :loading="loading"
          :span-method="spanMethod"
          :resizable="true"
          :max-height="maxHeight"
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
    <ElCol ::span="24">
      <ElCard>
        <ElPagination
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
          :page-sizes="[10, 20, 50, 100, 200, 500, 1000]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
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

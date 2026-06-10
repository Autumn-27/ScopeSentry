<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { h, nextTick, onMounted, reactive, Ref, ref } from 'vue'
import { Icon } from '@/components/Icon'
import { useTable } from '@/hooks/web/useTable'
import {
  ElCard,
  ElPagination,
  ElCol,
  ElRow,
  ElText,
  ElButton,
  ElInput,
  ElTag,
  InputInstance,
  ElLink
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import { addTagApi, deleteTagApi, getDirScanApi, totalDataApi } from '@/api/asset'
import Csearch from '../search/Csearch.vue'
import { RowState } from '@/api/asset/types'
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
    keyword: 'url',
    example: 'url="http://example.com"',
    explain: t('searchHelp.url')
  },
  {
    keyword: 'statuscode',
    example: 'statuscode=="200"',
    explain: t('searchHelp.statuscode')
  },
  {
    keyword: 'redirect',
    example: 'redirect="https://example.com"',
    explain: t('searchHelp.redirect')
  },
  {
    keyword: 'project',
    example: 'project="Hackerone"',
    explain: t('searchHelp.project')
  },
  {
    keyword: 'length',
    example: 'length="1234"',
    explain: t('searchHelp.length')
  }
]
const searchParams = ref('')
const handleSearch = (data: any) => {
  searchParams.value = data
  getList()
}
const rowStateMap = reactive<Record<string, RowState>>({})
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
    minWidth: 55
  },
  {
    field: 'url',
    label: 'URL',
    minWidth: 200,
    formatter: (_, __: TableColumn, url: string) => {
      return (
        <ElLink href={url} underline={false} target="_blank">
          {url}
        </ElLink>
      )
    }
  },
  {
    field: 'status',
    label: t('dirScan.status'),
    columnKey: 'status',
    minWidth: 120,
    formatter: (_: Recordable, __: TableColumn, statusValue: number) => {
      if (statusValue == null) {
        return <div>-</div>
      }
      let color = ''
      if (statusValue < 300) {
        color = '#2eb98a'
      } else if (statusValue < 400) {
        color = '#ff5252'
      } else {
        color = '#ff5252'
      }
      return (
        <ElRow gutter={1}>
          <ElCol span={1}>
            <Icon
              icon="clarity:circle-solid"
              color={color}
              size={10}
              style={'transform: translateY(8%)'}
            />
          </ElCol>
          <ElCol span={2}>
            <ElText>{statusValue}</ElText>
          </ElCol>
        </ElRow>
      )
    },
    filters: [
      { text: '200', value: 200 },
      { text: '201', value: 201 },
      { text: '204', value: 204 },
      { text: '301', value: 301 },
      { text: '302', value: 302 },
      { text: '304', value: 304 },
      { text: '400', value: 400 },
      { text: '401', value: 401 },
      { text: '403', value: 403 },
      { text: '404', value: 404 },
      { text: '405', value: 405 },
      { text: '500', value: 500 },
      { text: '502', value: 502 },
      { text: '503', value: 503 },
      { text: '504', value: 504 }
    ]
  },
  {
    field: 'length',
    label: 'Length',
    minWidth: 120,
    sortable: 'custom'
  },
  {
    field: 'msg',
    label: 'Redirect',
    minWidth: 200
  },
  {
    field: 'tags',
    label: 'TAG',
    fit: 'true',
    showOverflowTooltip: false,
    formatter: (row: Recordable, __: TableColumn, tags: string[]) => {
      if (tags == null) {
        tags = []
      }
      // 初始化状态
      if (!rowStateMap[row.id]) {
        rowStateMap[row.id] = {
          inputVisible: false,
          inputValue: '',
          inputRef: ref(null) as Ref<InputInstance | null>
        }
      }
      const rowState = rowStateMap[row.id]
      const handleInputConfirm = async () => {
        if (rowState.inputValue) {
          tags.push(rowState.inputValue) // 将输入值添加到 tags 中
          addTagApi(row.id, index, rowState.inputValue)
        }
        rowState.inputVisible = false // 隐藏输入框
        rowState.inputValue = '' // 清空输入框的值
      }
      const deleteTag = async (row: Recordable, tag: string) => {
        if (!row.tags) {
          row.tags = []
        }
        const indexT = row.tags.indexOf(tag)
        if (indexT > -1) {
          row.tags.splice(indexT, 1) // 从数组中移除指定的元素
        }
        await deleteTagApi(row.id, index, tag)
        // 如需强制刷新表格数据，可取消下行注释
        // await getList()
      }
      const showInput = () => {
        rowState.inputVisible = true
        nextTick(() => {
          // console.log('inputRef:', rowState.inputRef)
          // if (rowState.inputRef.value?.input) {
          //   rowState.inputRef.value.input.focus()
          // }
        })
      }
      // 标签点击处理函数
      const handleTagClick = (event: MouseEvent, tag: string) => {
        if ((event.target as HTMLElement).classList.contains('el-tag__close')) {
          // 点击关闭按钮时不处理
          return
        }
        // 这里可以添加处理点击事件的逻辑
        console.log('Tag clicked:', tag)
        changeTags('tags', tag)
      }

      return h(ElRow, {}, () => [
        // 渲染标签
        ...tags.map((tag) =>
          h(ElCol, { span: 24, key: tag }, () => [
            h('div', { onClick: (event: MouseEvent) => handleTagClick(event, tag) }, [
              h(ElTag, { closable: true, onClose: () => deleteTag(row, tag) }, () => tag)
            ])
          ])
        ),

        // 输入框或按钮
        h(
          ElCol,
          { span: 24 },
          rowState.inputVisible
            ? () =>
                h(ElInput, {
                  ref: rowState.inputRef,
                  modelValue: rowState.inputValue, // 双向绑定输入框值
                  'onUpdate:modelValue': (value: string) => (rowState.inputValue = value),
                  class: 'w-20',
                  size: 'small',
                  onKeyup: (event: KeyboardEvent) => {
                    if (event.key === 'Enter') {
                      handleInputConfirm() // 只在回车键被按下时触发
                    }
                  },
                  onBlur: handleInputConfirm // 失去焦点时调用 handleInputConfirm
                })
            : () =>
                h(
                  ElButton,
                  { class: 'button-new-tag', size: 'small', onClick: () => showInput() },
                  () => '+ New Tag'
                )
        )
      ])
    },
    minWidth: '130'
  }
])

let index = 'DirScanResult'
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

const lastSearchParams = ref('')
let lastFilter = reactive<{ [key: string]: any }>({})

const { allSchemas } = useCrudSchemas(crudSchemas)
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const searchParamsChanged = searchParams.value !== lastSearchParams.value
    const filterChanged = JSON.stringify(filter) !== JSON.stringify(lastFilter)
    const { currentPage, pageSize } = tableState

    if (searchParamsChanged || filterChanged) {
      currentPage.value = 1
      getTotal(searchParams.value, currentPage.value, pageSize.value, filter)
      lastSearchParams.value = searchParams.value
      lastFilter = { ...filter }
    }
    const res = await getDirScanApi(
      searchParams.value,
      currentPage.value,
      pageSize.value,
      filter,
      sortBy
    )
    return {
      list: res.data.list,
      flag: true
    }
  },
  immediate: false
})
const { loading, dataList, total, currentPage, pageSize } = tableState
const { getList, getElTableExpose } = tableMethods

const getTotal = async (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
) => {
  let res = await totalDataApi(search, pageIndex, pageSize, filter, index)
  total.value = res.data.total
}

pageSize.value = 20
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})

const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.7
}
const filter = reactive<{ [key: string]: any }>({})
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  getList()
}
const sortBy = reactive<{ [key: string]: any }>({})
const sortChange = async (column: any) => {
  const key = column.prop
  const value = column.order
  sortBy[key] = value
  getList()
}
const handleFilterSearch = (data: any, newFilters: any) => {
  Object.assign(filter, newFilters)
  searchParams.value = data
  getList()
}

const dynamicTags = ref<string[]>([])
const changeTags = (type, value) => {
  const key = `${type}=${value}`
  console.log(key)
  dynamicTags.value = [...dynamicTags.value, key]
}
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
</script>

<template>
  <Csearch
    :getList="getList"
    :handleSearch="handleSearch"
    :searchKeywordsData="searchKeywordsData"
    :index="index"
    :dynamicTags="dynamicTags"
    :handleClose="handleClose"
    :getElTableExpose="getElTableExpose"
    :handleFilterSearch="handleFilterSearch"
    :projectList="$props.projectList"
    :taskList="$props.taskList"
    :crudSchemas="crudSchemas"
    @update-column-visibility="handleColumnVisibilityChange"
    :searchResultCount="total"
    :getFilter="getFilter"
  />
  <ElRow>
    <ElCol>
      <ElCard>
        <Table
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
          :columns="allSchemas.tableColumns"
          :data="dataList"
          stripe
          :border="true"
          :loading="loading"
          :resizable="true"
          :max-height="maxHeight"
          @register="tableRegister"
          @filter-change="filterChange"
          @sort-change="sortChange"
          :headerCellStyle="tableHeaderColor"
          :tooltip-options="{
            offset: 1,
            showArrow: false,
            effect: 'dark',
            enterable: true,
            showAfter: 0,
            popperOptions: {},
            popperClass: 'test',
            placement: 'top',
            hideAfter: 0,
            disabled: true
          }"
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
</style>

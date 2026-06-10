<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { computed, h, nextTick, reactive, Ref, ref } from 'vue'
import { onMounted } from 'vue'
import { useTable } from '@/hooks/web/useTable'
import {
  ElCard,
  ElPagination,
  ElCol,
  ElRow,
  ElScrollbar,
  ElDrawer,
  ElTable,
  ElTableColumn,
  ElTag,
  ElInput,
  ElButton,
  InputInstance,
  ElSelect,
  ElOption,
  ElLink
} from 'element-plus'
import { Dialog } from '@/components/Dialog'
import { Table, TableColumn } from '@/components/Table'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import {
  addTagApi,
  deleteTagApi,
  getSensitiveInfoApi,
  getSensitiveNamesApi,
  getSensitiveResultApi,
  getSensitiveResultBodyApi,
  getSensitiveResultNumberApi,
  updateStatusApi
} from '@/api/asset'
import { Icon } from '@iconify/vue'
import { BaseButton } from '@/components/Button'
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
    keyword: 'sname',
    example: 'sname="twilio_account_sid"',
    explain: t('searchHelp.sname')
  },
  {
    keyword: 'info',
    example: 'info="api-key-example"',
    explain: t('searchHelp.sinfo')
  },
  {
    keyword: 'project',
    example: 'project="Hackerone"',
    explain: t('searchHelp.project')
  },
  {
    keyword: 'md5',
    example: 'md5=="1d49e5e190f7a38ab498e28e6578f64f"',
    explain: t('searchHelp.sensMd5')
  },
  {
    keyword: 'level',
    example: 'level=="rad"',
    explain: t('searchHelp.sensLevel')
  }
]
const searchParams = ref('')
const handleSearch = (data: any) => {
  searchParams.value = data
  callFunctionsConcurrently()
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
    label: '',
    type: 'index',
    minWidth: 50
  },
  {
    field: 'url',
    label: 'URL',
    minWidth: 200
  },
  {
    field: 'name',
    label: t('sensitiveInformation.sensitiveName'),
    minWidth: 150
  },
  // {
  //   field: 'color',
  //   label: 'Level',
  //   minWidth: 50,
  //   columnKey: 'color',
  //   formatter: (_, __: TableColumn, cellValue: string) => {
  //     if (!cellValue) {
  //       return
  //     }
  //     return (
  //       <Icon icon="clarity:circle-solid" color={cellValue} style={'transform: translateY(-35%)'} />
  //     )
  //   },
  //   filters: [
  //     { text: 'Red', value: 'red' },
  //     { text: 'Green', value: 'green' },
  //     { text: 'Cyan', value: 'cyan' },
  //     { text: 'Yellow', value: 'yellow' },
  //     { text: 'Orange', value: 'orange' },
  //     { text: 'Gray', value: 'gray' },
  //     { text: 'Pink', value: 'pink' },
  //     { text: 'Null', value: 'null' }
  //   ]
  // },
  {
    field: 'match',
    label: 'Info',
    minWidth: 150,
    formatter: (_, __: TableColumn, cellValue: string[]) => {
      if (!cellValue) {
        return
      }
      const elements = cellValue.map((line, index) => <div key={index}>{line}</div>)
      return (
        <ElScrollbar height="100px">
          <div class="scrollbar-demo-item">{elements}</div>
        </ElScrollbar>
      )
    }
  },
  {
    field: 'status',
    label: t('common.state'),
    minWidth: 100,
    columnKey: 'status',
    formatter: (row: Recordable, __: TableColumn, _: number) => {
      if (row.id.includes('//')) {
        return
      }
      if (row.id.includes('APP')) {
        return
      }
      if (row.status == null) {
        row.status = 1
      }

      if (row.status == 0) {
        row.status = 1
      }

      const options = [
        { value: 1, label: t('common.unprocessed'), color: '#909399' },
        { value: 2, label: t('common.processing'), color: '#409EFF' },
        { value: 3, label: t('common.ignored'), color: '#C0C4CC' },
        { value: 4, label: t('common.suspected'), color: '#E6A23C' },
        { value: 5, label: t('common.confirmed'), color: '#F56C6C' },
        { value: 6, label: t('common.processed'), color: '#67C23A' }
      ]
      const selected = options.find((opt) => opt.value === row.status)
      const selectedColor = selected?.color || '#000'
      return (
        <ElSelect
          modelValue={row.status}
          class="colored-select"
          popper-class="colored-select-popper"
          style={{ '--select-text-color': selectedColor }}
          onUpdate:modelValue={async (newValue) => {
            try {
              row.status = newValue
              await updateStatusApi(row.id, 'SensitiveResult', newValue)
            } catch (error) {
              console.error(error)
            }
          }}
        >
          {options.map((item) => (
            <ElOption
              key={item.value}
              label={item.label}
              value={item.value}
              style={{ color: item.color }}
            />
          ))}
        </ElSelect>
      )
    },
    filters: [
      { text: t('common.unprocessed'), value: 1 },
      { text: t('common.processing'), value: 2 },
      { text: t('common.ignored'), value: 3 },
      { text: t('common.suspected'), value: 4 },
      { text: t('common.confirmed'), value: 5 },
      { text: t('common.processed'), value: 5 }
    ]
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
  },
  {
    field: 'time',
    label: t('asset.time'),
    minWidth: 200
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    formatter: (row, __: TableColumn, _: number) => {
      if (row.body_id) {
        return (
          <>
            <BaseButton type="primary" onClick={() => action(row.body_id)}>
              {t('asset.detail')}
            </BaseButton>
          </>
        )
      } else {
        return
      }
    },
    minWidth: 100
  }
])
let index = 'SensitiveResult'
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

const filter = reactive<{ [key: string]: any }>({})
const { allSchemas } = useCrudSchemas(crudSchemas)
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const searchParamsChanged = searchParams.value !== lastSearchParams.value
    const filterChanged = JSON.stringify(filter) !== JSON.stringify(lastFilter)
    if (searchParamsChanged || filterChanged) {
      currentPage.value = 1
      getTotal(searchParams.value, currentPage.value, pageSize.value, filter)
      lastSearchParams.value = searchParams.value
      lastFilter = { ...filter }
    }
    const res = await getSensitiveResultApi(
      searchParams.value,
      currentPage.value,
      pageSize.value,
      filter
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

let allNumber = ref(0)
const getTotal = async (
  search: string,
  pageIndex: number,
  pageSize: number,
  filter: Record<string, any>
) => {
  let res = await getSensitiveResultNumberApi(search, pageIndex, pageSize, filter)
  total.value = res.data.total
  allNumber.value = res.data.all
}

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
const DialogVisible = ref(false)
const body = ref('')
const action = async (id) => {
  id.replace('md5==', '')
  const res = await getSensitiveResultBodyApi(id)
  body.value = res.data.body
  DialogVisible.value = true
}
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  callFunctionsConcurrently()
}
const handleFilterSearch = (data: any, newFilters: any) => {
  Object.assign(filter, newFilters)
  searchParams.value = data
  callFunctionsConcurrently()
}
const drawerFlag = ref(false)

const openAggregation = () => {
  drawerFlag.value = true
}
const aggregationData = ref<
  {
    name: string
    color: string
    count: number
  }[]
>([])

const getSensNames = async () => {
  try {
    const res = await getSensitiveNamesApi(searchParams.value, filter)
    if (res && res.data && Array.isArray(res.data.list)) {
      aggregationData.value = res.data.list.map((data) => ({
        name: data.name,
        color: data.color,
        count: data.count
      }))
    }
  } catch (error) {
    console.error('Error fetching node data:', error)
  } finally {
  }
}
const callFunctionsConcurrently = async () => {
  try {
    await Promise.all([getList(), getSensNames()])
    console.log('Both getList and getSensNames have completed')
  } catch (error) {
    console.error('Error occurred while executing getList and getSensNames:', error)
  }
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
const ViewInfoDialogVisible = ref(false)

const infoArray = ref<string[]>([])
const infNmae = ref('')
const OpenViewInfoDialogVisible = async (sid) => {
  infNmae.value = sid
  const res = await getSensitiveInfoApi(sid, searchParams.value, filter)
  infoArray.value = res.data.list
  ViewInfoDialogVisible.value = true
}
const colorMap = {
  null: '#f4f4f5',
  green: '#D0F5EA',
  red: '#FFD6D6',
  cyan: '#D4F6FF',
  yellow: '#FFF9C4',
  orange: '#FFD9B3',
  gray: '#ECECEC',
  pink: '#FFE0F0',
  default: '#E9F3FF'
}

const parseTags = (tags) => {
  if (!tags || !tags.length) return []
  return tags.map((tagStr) => {
    const [type, value] = tagStr.split('=')
    return { type, value }
  })
}
interface SensitiveAggregation {
  name: string
  color: string
  count: number
}
const drawerAggregationData = computed<SensitiveAggregation[]>(() => {
  if (!dynamicTags.value?.length) return aggregationData.value
  const parsedTags = parseTags(dynamicTags.value)
  const snames = parsedTags.filter((tag) => tag.type === 'sname').map((tag) => tag.value)
  return aggregationData.value.filter((item) => snames.includes(item.name))
})
</script>

<template>
  <Csearch
    :getList="getList"
    :handleSearch="handleSearch"
    :searchKeywordsData="searchKeywordsData"
    :index="index"
    :getElTableExpose="getElTableExpose"
    :handleFilterSearch="handleFilterSearch"
    :projectList="$props.projectList"
    :taskList="$props.taskList"
    :openAggregation="openAggregation"
    :dynamicTags="dynamicTags"
    :handleClose="handleClose"
    :crudSchemas="crudSchemas"
    @update-column-visibility="handleColumnVisibilityChange"
    :searchResultCount="total"
    :sensitiveAllNumber="allNumber"
    :getFilter="getFilter"
  />
  <ElRow>
    <div style="display: flex; width: 100%">
      <!-- 左侧目录栏，高度与表格一致，内容可滚动 -->
      <div style="width: 220px; margin-right: 18px">
        <ElCard style="height: 100%">
          <ElScrollbar :height="`${maxHeight}px`" class="agg-scroll">
            <div
              v-for="item in aggregationData"
              :key="item.name"
              style="display: flex; align-items: center; cursor: pointer; margin-bottom: 10px"
              @click="changeTags('sname', item.name)"
            >
              <ElTag
                :style="{
                  background: colorMap[item.color || 'null'] || colorMap.default,
                  color: '#333',
                  border: 'none'
                }"
                effect="plain"
              >
                {{ item.name }}
              </ElTag>
              <span style="margin-left: 8px; color: #888">({{ item.count }})</span>
            </div>
          </ElScrollbar>
        </ElCard>
      </div>
      <!-- 右侧表格 -->
      <div style="flex: 1">
        <ElCol>
          <ElCard>
            <Table
              v-model:pageSize="pageSize"
              v-model:currentPage="currentPage"
              :columns="allSchemas.tableColumns"
              :data="dataList"
              stripe
              rowKey="id"
              :border="true"
              :max-height="maxHeight"
              :loading="loading"
              :resizable="true"
              @register="tableRegister"
              @filter-change="filterChange"
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
                disabled: false
              }"
              :style="{
                fontFamily:
                  '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
              }"
            />
          </ElCard>
        </ElCol>
      </div>
    </div>
    <ElCol ::span="22">
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
  <Dialog
    v-model="DialogVisible"
    :title="t('asset.detail')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    width="70%"
    :max-height="maxHeight"
  >
    <ElScrollbar :max-height="maxHeight">
      <div :style="{ whiteSpace: 'pre-line' }">{{ body }}</div>
    </ElScrollbar>
  </Dialog>
  <ElDrawer
    v-model="drawerFlag"
    :title="t('sensitiveInformation.sensAggre')"
    direction="rtl"
    size="30%"
  >
    <ElTable :data="drawerAggregationData">
      <ElTableColumn prop="name" :label="t('sensitiveInformation.sensitiveName')" width="180">
        <template #default="scope">
          <div
            style="display: flex; align-items: center"
            @click="changeTags('sname', scope.row.name)"
          >
            <ElTag>{{ scope.row.name }}</ElTag>
          </div>
        </template>
      </ElTableColumn>
      <ElTableColumn prop="color" label="color" width="100">
        <template #default="scope">
          <ElTag :color="scope.row.color" round effect="plain" size="small" style="width: 20px" />
        </template>
      </ElTableColumn>
      <ElTableColumn prop="count" :label="t('common.quantity')" width="130" />
      <ElTableColumn :label="t('tableDemo.operate')" width="180">
        <template #default="scope">
          <div
            style="display: flex; align-items: center"
            @click="OpenViewInfoDialogVisible(scope.row.name)"
          >
            <ElButton type="success">info</ElButton>
          </div>
        </template>
      </ElTableColumn>
    </ElTable>
  </ElDrawer>
  <Dialog v-model="ViewInfoDialogVisible" :title="infNmae">
    <div v-for="(item, index) in infoArray" :key="index">
      <p>{{ item }}</p>
    </div>
  </Dialog>
</template>

<style lang="less" scoped>
.el-button {
  margin-top: 10px;
}

::v-deep(.colored-select .el-select__selected-item) {
  color: var(--select-text-color) !important;
}

/* 生效于 popper 中的下拉选项 */
.colored-select-popper .el-select-dropdown__item {
  color: inherit;
}
.agg-scroll {
  // WebKit 美化
  &::-webkit-scrollbar {
    width: 8px;
  }
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  &::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 6px;
  }
  &::-webkit-scrollbar-thumb:hover {
    background-color: rgba(0, 0, 0, 0.3);
  }

  // Firefox
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.3) transparent;
}
</style>

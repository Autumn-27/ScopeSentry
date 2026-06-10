<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { h, nextTick, reactive, Ref, ref } from 'vue'
import { onMounted } from 'vue'
import { useTable } from '@/hooks/web/useTable'
import {
  ElCard,
  ElPagination,
  ElCol,
  ElRow,
  ElDescriptions,
  ElDescriptionsItem,
  ElText,
  ElButton,
  ElInput,
  ElTag,
  InputInstance,
  ElSelect,
  ElOption,
  ElMessage
} from 'element-plus'
import { Dialog } from '@/components/Dialog'
import { Table, TableColumn } from '@/components/Table'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import { getVulDetailApi, getVulResultDataApi } from '@/api/vul'
import { Icon } from '@iconify/vue'
import Csearch from '../search/Csearch.vue'
import { BaseButton } from '@/components/Button'
import { addTagApi, deleteTagApi, updateStatusApi } from '@/api/asset'
import { RowState } from '@/api/asset/types'
import Detail from '../../Poc/components/Detail.vue'
import { getPocContentApi, getPocDetailApi } from '@/api/poc'
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
    keyword: 'vulname',
    example: 'vulname="nginxwebui-runcmd-rce"',
    explain: t('searchHelp.vulname')
  },
  {
    keyword: 'level',
    example: 'level="info"',
    explain: t('searchHelp.level')
  },
  {
    keyword: 'matched',
    example: 'matched="https://example.com"',
    explain: t('searchHelp.matched')
  },
  {
    keyword: 'request',
    example: 'request="cmd=whoami"',
    explain: t('searchHelp.vulRequest')
  },
  {
    keyword: 'response',
    example: 'response="root"',
    explain: t('searchHelp.response')
  },
  {
    keyword: 'project',
    example: 'project="Hackerone"',
    explain: t('searchHelp.project')
  }
]
const searchParams = ref('')
const indexName = 'vulnerability'
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
    minWidth: 100
  },
  {
    field: 'vulnerability',
    label: 'Vulnerability',
    minWidth: 120
  },
  {
    field: 'level',
    label: 'Level',
    minWidth: 100,
    columnKey: 'level',
    formatter: (_: Recordable, __: TableColumn, levelValue: string) => {
      if (levelValue == null) {
        return <div></div>
      }
      let color = ''
      let flag = ''
      if (levelValue === 'critical') {
        color = 'red'
        flag = t('poc.critical')
      } else if (levelValue === 'high') {
        color = 'orange'
        flag = t('poc.high')
      } else if (levelValue === 'medium') {
        color = 'yellow'
        flag = t('poc.medium')
      } else if (levelValue === 'low') {
        color = 'blue'
        flag = t('poc.low')
      } else if (levelValue === 'info') {
        color = 'green'
        flag = t('poc.info')
      } else if (levelValue === 'unknown') {
        color = 'gray'
        flag = t('poc.unknown')
      }
      return (
        <ElRow gutter={20} style="width: 80%">
          <ElCol span={1}>
            <Icon icon="clarity:circle-solid" color={color} />
          </ElCol>
          <ElCol span={5}>
            <ElText type="info">{flag}</ElText>
          </ElCol>
        </ElRow>
      )
    },
    filters: [
      { text: t('poc.critical'), value: 'critical' },
      { text: t('poc.high'), value: 'high' },
      { text: t('poc.medium'), value: 'medium' },
      { text: t('poc.low'), value: 'low' },
      { text: t('poc.info'), value: 'info' },
      { text: t('poc.unknown'), value: 'unknown' }
    ]
  },
  {
    field: 'matched',
    label: 'Matched',
    minWidth: 150
  },
  {
    field: 'status',
    label: t('common.state'),
    minWidth: 100,
    columnKey: 'status',
    formatter: (row: Recordable, __: TableColumn, _: number) => {
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
              await updateStatusApi(row.id, index, newValue)
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
      { text: t('common.confirmed'), value: 5 }
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
      return (
        <>
          <BaseButton type="primary" onClick={() => action(row)}>
            {t('asset.detail')}
          </BaseButton>
          <BaseButton type="success" onClick={() => openPoc(row.vulnid)}>
            POC
          </BaseButton>
        </>
      )
    },
    minWidth: 100
  }
])
let index = 'vulnerability'
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
    const res = await getVulResultDataApi(
      searchParams.value,
      currentPage.value,
      pageSize.value,
      filter
    )
    return {
      list: res.data.list,
      total: res.data.total
    }
  },
  immediate: false
})
const { loading, dataList, total, currentPage, pageSize } = tableState
const { getList, getElTableExpose } = tableMethods
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
interface RowData {
  URL: string
  Vulnerability: string
  Level: string
  Matched: string
  Time: string
  Request: string
  Response: string
}
const DialogData = reactive<RowData>({
  URL: '',
  Vulnerability: '',
  Level: '',
  Matched: '',
  Time: '',
  Request: '',
  Response: ''
})
const color = ref('')
const DialogVisible = ref(false)
const action = async (data: any) => {
  const levelValue = data.level
  color.value = ''
  let flag = ''
  if (levelValue === 'critical') {
    color.value = 'red'
    flag = t('poc.critical')
  } else if (levelValue === 'high') {
    color.value = 'orange'
    flag = t('poc.high')
  } else if (levelValue === 'medium') {
    color.value = 'yellow'
    flag = t('poc.medium')
  } else if (levelValue === 'low') {
    color.value = 'blue'
    flag = t('poc.low')
  } else if (levelValue === 'info') {
    color.value = 'green'
    flag = t('poc.info')
  } else if (levelValue === 'unknown') {
    color.value = 'gray'
    flag = t('poc.unknown')
  }
  const res = await getVulDetailApi(data.hash)
  DialogData.Level = flag
  DialogData.Vulnerability = data.vulnerability
  DialogData.Matched = data.matched
  DialogData.Time = data.time
  DialogData.URL = data.url
  DialogData.Request = res.data.req
  DialogData.Response = res.data.res
  DialogVisible.value = true
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
let pocForm = reactive({
  id: '',
  name: '',
  level: '',
  content: '',
  tags: []
})

const dialogVisible = ref(false)
const openPoc = async (id) => {
  pocForm.id = ''
  pocForm.name = ''
  pocForm.level = ''
  pocForm.content = ''
  pocForm.tags = []
  const res = await getPocDetailApi(id)
  pocForm.id = res.data.id
  pocForm.name = res.data.name
  pocForm.level = res.data.level
  pocForm.content = res.data.content
  pocForm.tags = res.data.tags
  dialogVisible.value = true
}
const closeDialog = () => {
  dialogVisible.value = false
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
    :index="indexName"
    :getElTableExpose="getElTableExpose"
    :handleFilterSearch="handleFilterSearch"
    :projectList="$props.projectList"
    :taskList="$props.taskList"
    :crudSchemas="crudSchemas"
    :dynamicTags="dynamicTags"
    :handleClose="handleClose"
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
          :max-height="maxHeight"
          :loading="loading"
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
    <ElCol ::span="24">
      <ElCard>
        <ElPagination
          v-model:pageSize="pageSize"
          v-model:currentPage="currentPage"
          :page-sizes="[20, 50, 100, 200, 500, 1000]"
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
    <ElDescriptions :border="true" :column="2" style="word-break: break-word">
      <ElDescriptionsItem label="URL" :span="2">{{ DialogData.URL }}</ElDescriptionsItem>
      <ElDescriptionsItem label="Level">
        <Icon icon="clarity:circle-solid" :color="color" />
        <ElText type="info">{{ DialogData.Level }}</ElText>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="Vulnerability">{{ DialogData.Vulnerability }}</ElDescriptionsItem>
      <ElDescriptionsItem label="Matched">{{ DialogData.Matched }}</ElDescriptionsItem>
      <ElDescriptionsItem label="Time">{{ DialogData.Time }}</ElDescriptionsItem>
      <ElDescriptionsItem label="Request">
        <ElScrollbar :max-height="maxHeight" max-width="maxHeight">
          <div :style="{ whiteSpace: 'pre-line', width: '500px' }"> {{ DialogData.Request }}</div>
        </ElScrollbar>
      </ElDescriptionsItem>
      <ElDescriptionsItem label="Response">
        <ElScrollbar :max-height="maxHeight">
          <div :style="{ whiteSpace: 'pre-line' }">{{ DialogData.Response }}</div>
        </ElScrollbar>
      </ElDescriptionsItem>
    </ElDescriptions>
  </Dialog>
  <Dialog
    v-model="dialogVisible"
    :title="pocForm.id ? $t('common.edit') : $t('common.new')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="800"
  >
    <Detail :closeDialog="closeDialog" :pocForm="pocForm" :getList="getList" />
  </Dialog>
</template>

<style lang="less" scoped>
.el-button {
  margin-top: 10px;
}
.el-descriptions {
  margin-top: 20px;
}
::v-deep(.colored-select .el-select__selected-item) {
  color: var(--select-text-color) !important;
}

.colored-select-popper .el-select-dropdown__item {
  color: inherit;
}
.cell-item {
  display: flex;
  align-items: center;
}
.margin-top {
  margin-top: 20px;
}
</style>

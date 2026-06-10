<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { Ref, h, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useTable } from '@/hooks/web/useTable'
import { Dialog } from '@/components/Dialog'
import {
  ElRow,
  ElCol,
  ElCard,
  ElScrollbar,
  ElTag,
  ElTooltip,
  ElCollapse,
  ElCollapseItem,
  ElBadge,
  ElPagination,
  ElLink,
  ElText,
  ElButton,
  InputInstance,
  ElInput
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { CrudSchema, useCrudSchemas } from '@/hooks/web/useCrudSchemas'
import {
  addTagApi,
  deleteTagApi,
  getAssetApi,
  getAssetCardApi,
  getAssetScreenshotApi,
  getAssetStatisticsPortApi,
  getAssetStatisticsTypeApi,
  getAssetStatisticsappApi,
  getAssetStatisticsiconApi,
  totalDataApi
} from '@/api/asset'
import { Icon } from '@/components/Icon'
import { BaseButton } from '@/components/Button'
import { useRouter } from 'vue-router'
import Csearch from '../search/Csearch.vue'
import { createImageViewer } from '@/components/ImageViewer'
import { AssetData, RowState } from '@/api/asset/types'
import AssetDetail2 from '../detail/AssetDetail2.vue'
import { url } from 'inspector'
const { push } = useRouter()
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
    keyword: 'app',
    example: 'app="Nginx"',
    explain: t('searchHelp.app')
  },
  {
    keyword: 'body',
    example: 'body="bootstrap.min.css"',
    explain: t('searchHelp.body')
  },
  {
    keyword: 'header',
    example: 'header="rememberMe"',
    explain: t('searchHelp.header')
  },
  {
    keyword: 'title',
    example: 'title="admin console"',
    explain: t('searchHelp.title')
  },
  {
    keyword: 'statuscode',
    example: 'statuscode=="403"',
    explain: t('searchHelp.statuscode')
  },
  {
    keyword: 'icon',
    example: 'icon="54256234"',
    explain: t('searchHelp.icon')
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
    keyword: 'domain',
    example: 'domain="example.com"',
    explain: t('searchHelp.domain')
  },
  {
    keyword: 'service',
    example: 'service="ssh"',
    explain: t('searchHelp.protocol')
  },
  {
    keyword: 'banner',
    example: 'banner="SSH-2.0-OpenSSH"',
    explain: t('searchHelp.banner')
  },
  {
    keyword: 'project',
    example: 'project="Hackerone"',
    explain: t('searchHelp.project')
  },
  {
    keyword: 'type',
    example: 'type="http"',
    explain: t('searchHelp.protocol')
  }
]
const staticLoading = ref(true)
const searchParams = ref('')
const handleSearch = (data: any) => {
  searchParams.value = data
  if (activeSegment.value == 'cardSegment') {
    getAssetCardData()
    return
  }
  staticLoading.value = true
  AssetstatisticsData.value.Icon = []
  getList()
  staticLoading.value = false
}

let AssetstatisticsData: Ref<{
  Port: { value: number; number: number }[]
  Service: { value: string; number: number }[]
  Product: { value: string; number: number }[]
  Icon: { value: string; number: number; icon_hash: string }[]
}> = ref({
  Port: [],
  Service: [],
  Product: [],
  Icon: []
})
let iconPage = 1 // 当前页
let iconPageSize = 50 // 每次加载的数据量

const getAssetstatistics = async () => {
  if (statisticsHidden.value) {
    return
  }
  AssetstatisticsData.value.Port = []
  AssetstatisticsData.value.Service = []
  AssetstatisticsData.value.Product = []
  staticLoading.value = true
  AssetstatisticsData.value.Icon = []
  const [portRes, serviceRes, productRes] = await Promise.all([
    getAssetStatisticsPortApi(searchParams.value, filter),
    getAssetStatisticsTypeApi(searchParams.value, filter),
    getAssetStatisticsappApi(searchParams.value, filter)
  ])

  AssetstatisticsData.value.Port = portRes.data.Port
  AssetstatisticsData.value.Service = serviceRes.data.Service
  AssetstatisticsData.value.Product = productRes.data.Product
  staticLoading.value = false
  iconPage = 1
  iconPageSize = 50
  let iconRes = await getAssetStatisticsiconApi(searchParams.value, filter, iconPage, iconPageSize)
  AssetstatisticsData.value.Icon = iconRes.data.Icon
}

const rowStateMap = reactive<Record<string, RowState>>({})

let crudSchemas = reactive<CrudSchema[]>([
  {
    field: 'selection',
    type: 'selection',
    minWidth: '55'
  },
  {
    field: 'index',
    label: t('tableDemo.index'),
    type: 'index',
    minWidth: '15'
  },
  {
    field: 'domain',
    label: t('asset.domain'),
    minWidth: '200',
    formatter: (row, __: TableColumn, domainValue: string) => {
      return (
        <div class="flex">
          <Icon
            icon="material-symbols-light:bring-your-own-ip"
            style={'transform: translateY(35%)'}
            size={16}
            color="#409eff"
          />
          <ElLink
            href={row.type === 'http' ? row.url : `${row.service}://${domainValue}`}
            underline={false}
            target="_blank"
          >
            {domainValue}
          </ElLink>
        </div>
      )
    }
  },
  {
    field: 'ip',
    label: t('asset.IP'),
    minWidth: '130',
    formatter: (row, __: TableColumn, ipValue: string) => {
      return (
        <div class="flex">
          <Icon
            icon="arcticons:ip-tools"
            style={'transform: translateY(30%)'}
            size={15}
            color="red"
          />
          <ElLink href={row.url} underline={false}>
            {ipValue}
          </ElLink>
        </div>
      )
    }
  },
  {
    field: 'port',
    label: t('asset.port') + '/' + t('asset.service'),
    minWidth: '110',
    formatter: (raw, __: TableColumn, statusValue: number) => {
      if (raw.service == '') {
        return <div>{statusValue}</div>
      } else {
        return (
          <div class="flex">
            <div>{statusValue}</div>
            <ElTag
              type="info"
              effect="dark"
              round
              size="small"
              style={'top: 2px; left:6px; position:relative'}
            >
              {raw.service}
            </ElTag>
          </div>
        )
      }
    }
  },
  {
    field: 'status',
    label: t('asset.status'),
    minWidth: '85',
    columnKey: 'statuscode',
    formatter: (row: Recordable, __: TableColumn, statusValue: number) => {
      if (row.type == 'other') {
        return <div>-</div>
      }
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
        <ElRow gutter={10}>
          <ElCol span={2}>
            <Icon
              icon="clarity:circle-solid"
              color={color}
              size={6}
              style={'transform: translateY(-35%)'}
            />
          </ElCol>
          <ElCol span={18}>
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
      { text: '500', value: 500 },
      { text: '502', value: 502 },
      { text: '503', value: 503 },
      { text: '504', value: 504 }
    ]
  },
  {
    field: 'title',
    label: t('asset.title'),
    minWidth: '150',
    formatter: (row: Recordable, __: TableColumn, title: string) => {
      if (title == null || title == '') {
        title = ''
      }
      if (row.faviconmmh3 == '' || row.faviconmmh3 == null) {
        return (
          <ElRow gutter={10}>
            <ElCol span={24}>
              <ElText size="small" class="w-200px mb-2" truncated>
                {title}
              </ElText>
            </ElCol>
          </ElRow>
        )
      }
      const st = '/images/icon/' + row.faviconmmh3 + '.png'
      return (
        <ElRow gutter={20}>
          <ElCol span={2}>
            <img src={st} alt="Icon" style="width: 20px; height: 20px" />
          </ElCol>
          <ElCol span={18}>
            <ElText size="small" class="w-200px mb-2" truncated>
              {title}
            </ElText>
          </ElCol>
        </ElRow>
      )
    }
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
    field: 'banner',
    label: t('asset.banner'),
    fit: 'true',
    formatter: (row: Recordable, __: TableColumn, bannerValue: string) => {
      if (row.type == 'other') {
        return (
          <ElScrollbar>
            <div class="scrollbar-demo-item">{row.metadata}</div>
          </ElScrollbar>
        )
      }
      const lines = bannerValue.split('\n')
      const elements = lines.map((line, index) => <div key={index}>{line}</div>)
      return (
        <ElScrollbar height="150px">
          <div class="scrollbar-demo-item">{elements}</div>
        </ElScrollbar>
      )
    },
    minWidth: '190'
  },
  {
    field: 'products',
    label: t('asset.products'),
    minWidth: '110',
    formatter: (_: Recordable, __: TableColumn, ProductsValue: string[] | null) => {
      if (!ProductsValue || ProductsValue.length === 0) return
      if (ProductsValue.length != 0) {
        return (
          <ElRow style={{ flexWrap: 'wrap' }}>
            {ProductsValue.map((product) => (
              <ElCol span={24} key={product}>
                <div
                  onClick={() => changeTags('app', product)}
                  style={'display: inline-block; cursor: pointer'}
                >
                  <ElTag type={'success'}>{product}</ElTag>
                </div>
              </ElCol>
            ))}
          </ElRow>
        )
        // if (ProductsValue.length > 1) {
        //   let contentTool = ''
        //   if (Array.isArray(ProductsValue)) {
        //     // It's an array, you can use forEach
        //     ProductsValue.forEach((item, _) => {
        //       contentTool += `<div>${item}</div>`
        //     })
        //   } else {
        //     console.error('ProductsValue is not an array')
        //   }
        //   return (
        //     <div class="flex">
        //       <ElTag type="success" effect="light" round>
        //         {ProductsValue[0]}
        //       </ElTag>
        //       <ElTooltip
        //         class="box-item"
        //         effect="dark"
        //         placement="top-start"
        //         content={contentTool}
        //         popper-class="tagtooltip"
        //         rawContent
        //       >
        //         <ElTag type="info" effect="plain" round style={'left:3px; position:relative'}>
        //           {t('asset.total')} {ProductsValue.length} {t('asset.p')}
        //         </ElTag>
        //       </ElTooltip>
        //     </div>
        //   )
        // } else {
        //   return (
        //     <div class="flex">
        //       <ElTag type="success" effect="light">
        //         {ProductsValue[0]}
        //       </ElTag>
        //     </div>
        //   )
        // }
      }
    }
  },
  {
    field: 'screenshot',
    label: t('asset.screenshot'),
    minWidth: '170',
    formatter: (row) => {
      if (row.ResponseBodyHash == undefined) {
        return
      }
      if (row.ResponseBodyHash == '') {
        return
      }
      if (row.ResponseBodyHash != '') {
        const imageSrc = `/images/screenshots/${row.ResponseBodyHash}.png`
        return (
          <img
            key={`${row.id}-${row.ResponseBodyHash}`}
            src={imageSrc}
            alt="screenshot"
            style={{
              width: '100%',
              height: 'auto',
              maxHeight: '250px'
            }}
            onError={(e) => {
              // 404 或其他错误时，隐藏图片（不显示裂开的图片）
              ;(e.target as HTMLImageElement).style.display = 'none'
            }}
            onClick={() => handleImageClick(imageSrc)}
          />
        )
      }
    }
  },
  {
    field: 'time',
    label: t('asset.time'),
    minWidth: '170'
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    fixed: 'right',
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          <BaseButton
            type="primary"
            onClick={() => openDetail(row.id, row.service + '://' + row.domain, row.ip, row.port)}
          >
            {t('asset.detail')}
          </BaseButton>
        </>
      )
    },
    minWidth: '100'
  }
])

const getScreenshot = async (id) => {
  const response = await getAssetScreenshotApi(id)
  return response.data.screenshot
}
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  getList()
}
const action = (id: string) => {
  push(`/asset-information/asset-detail?id=${id}`)
}

const handleImageClick = (screenshot: string) => {
  createImageViewer({
    urlList: [screenshot]
  })
}
let index = 'asset'
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
      col.hidden = savedConfig[col.field] // 复列的显示状态
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
  console.log('statisticsHidden.value', statisticsHidden.value)
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

const filter = reactive<{ [key: string]: any }>({})

const lastSearchParams = ref('')
let lastFilter = reactive<{ [key: string]: any }>({})

const { allSchemas } = useCrudSchemas(crudSchemas)
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    if (activeSegment.value == 'cardSegment') {
      await getAssetCardData()
      return {
        list: [],
        flag: true
      }
    }
    const searchParamsChanged = searchParams.value !== lastSearchParams.value
    const filterChanged = JSON.stringify(filter) !== JSON.stringify(lastFilter)
    const { currentPage, pageSize } = tableState

    if (searchParamsChanged || filterChanged) {
      currentPage.value = 1
      getTotal(searchParams.value, currentPage.value, pageSize.value, filter)
      getAssetstatistics()
      lastSearchParams.value = searchParams.value
      lastFilter = { ...filter }
    }
    const res = await getAssetApi(searchParams.value, currentPage.value, pageSize.value, filter)
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
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
function rowstyle() {
  return { maxheight: '10px' }
}
const activeNames = ref(['1', '2', '3', '4', '5'])
const handleFilterSearch = (data: any, newFilters: any) => {
  Object.assign(filter, newFilters)
  searchParams.value = data
  if (activeSegment.value == 'cardSegment') {
    getAssetCardData()
    return
  }
  getList()
}
const dynamicTags = ref<string[]>([])
const changeTags = (type, value) => {
  const key = `${type}=${value}`
  console.log(key)
  dynamicTags.value = [...dynamicTags.value, key]
}
const handleClose = (tag: string) => {
  if (tag == 'close') {
    dynamicTags.value = []
  } else {
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
}
const changeStatisticsHidden = (value: boolean) => {
  statisticsHidden.value = value
  saveColumnConfig()
}
const detailVisible = ref(false)
const detailId = ref('')
const detailhost = ref('')
const detailip = ref('')
const detailport = ref()
const openDetail = (id: string, host: string, ip: string, port: number) => {
  detailhost.value = ''
  detailip.value = ''
  detailport.value = null
  detailId.value = id
  detailhost.value = host
  detailip.value = ip
  detailport.value = port
  detailVisible.value = true
}

const isLoading = ref(false)
// 滚动监听，触发加载更多
const scrollbarRef = ref<InstanceType<typeof ElScrollbar>>()

const handleScroll = ({ scrollTop }: { scrollTop: number }) => {
  const wrap = scrollbarRef.value?.wrapRef
  if (!wrap || isLoading.value) return
  const { scrollHeight, clientHeight } = wrap
  // 当滚动到距离底部 20px 时触发加载
  if (scrollHeight - (scrollTop + clientHeight) < 20) {
    loadMoreIcons()
  }
}

const loadMoreIcons = async () => {
  try {
    isLoading.value = true
    iconPage++
    const iconRes = await getAssetStatisticsiconApi(
      searchParams.value,
      filter,
      iconPage,
      iconPageSize
    )
    if (iconRes.data.Icon?.length) {
      AssetstatisticsData.value.Icon.push(...iconRes.data.Icon)
    }
  } finally {
    isLoading.value = false
  }
}
const activeSegment = ref<'tableSegment' | 'cardSegment'>('tableSegment')

const setActiveSegment = (segment: 'tableSegment' | 'cardSegment', flag: boolean) => {
  activeSegment.value = segment
  // 将配置存储到 localStorage
  localStorage.setItem(`assetActiveSegment`, JSON.stringify({ activeSegment: segment }))
  if (flag) {
    getList()
  }
}
const websites = ref<AssetData[]>([])
// 跟踪每个图片的加载状态：'loading' | 'loaded' | 'error'
const imageLoadStates = ref<Record<string, 'loading' | 'loaded' | 'error'>>({})
const handleImageLoad = (hash: string) => {
  if (imageLoadStates.value[hash] !== 'loaded') {
    imageLoadStates.value[hash] = 'loaded'
  }
}
const handleImageError = (hash: string) => {
  if (imageLoadStates.value[hash] !== 'error') {
    imageLoadStates.value[hash] = 'error'
  }
}
const getAssetCardData = async () => {
  websites.value = []
  // 重置图片加载状态
  imageLoadStates.value = {}
  getTotal(searchParams.value, currentPage.value, pageSize.value, filter)
  const res = await getAssetCardApi(searchParams.value, currentPage.value, pageSize.value, filter)
  websites.value = res.data.list
  // 初始化所有图片为加载中状态
  res.data.list.forEach((site) => {
    if (site.ResponseBodyHash) {
      imageLoadStates.value[site.ResponseBodyHash] = 'loading'
    }
  })
}
const getStatusColor = (statusValue) => {
  if (statusValue < 300) {
    return '#2eb98a' // 绿色，表示成功
  } else if (statusValue < 400) {
    return '#ff9800' // 橙色，表示重定向
  } else {
    return '#ff5252' // 红色，表示错误
  }
}
const getFilter = () => {
  return filter
}
const selectedIcons = ref<{ value: string; number: number; icon_hash: string }[]>([])

const toggleSelectIcon = (iconItem: { value: string; number: number; icon_hash: string }) => {
  const idx = selectedIcons.value.findIndex((i) => i.icon_hash === iconItem.icon_hash)
  if (idx > -1) {
    selectedIcons.value.splice(idx, 1) // 已选再点取消
  } else {
    selectedIcons.value.push(iconItem)
  }
}

const confirmSelectedIcons = () => {
  dynamicTags.value = [
    ...dynamicTags.value,
    ...selectedIcons.value.map((i) => `icon=${i.icon_hash}`)
  ]
  selectedIcons.value = []
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
    :dynamicTags="dynamicTags"
    :handleClose="handleClose"
    :crudSchemas="crudSchemas"
    @update-column-visibility="handleColumnVisibilityChange"
    :statisticsHidden="statisticsHidden"
    :changeStatisticsHidden="changeStatisticsHidden"
    :searchResultCount="total"
    :activeSegment="activeSegment"
    :setActiveSegment="setActiveSegment"
    :getFilter="getFilter"
    :iconData="AssetstatisticsData.Icon"
  />
  <ElRow :gutter="3" v-if="activeSegment == 'tableSegment'">
    <ElCol :span="statisticsHidden ? 0 : 3">
      <ElCard v-loading="staticLoading">
        <div>
          <ElRow>
            <ElCol :span="12">
              <ElText tag="b" size="small">{{ t('asset.assetTotalNum') }}:</ElText>
            </ElCol>
            <ElCol :span="12" style="text-align: end">
              <ElText size="small">{{ total }}</ElText>
            </ElCol>
          </ElRow>
        </div>
        <ElCollapse v-model="activeNames" style="position: relative">
          <!-- <ElCollapseItem name="1">
            <template #title>
              <ElText tag="b" size="small">{{ t('asset.assetTotalNum') }}</ElText>
            </template>
            xxxxx
          </ElCollapseItem> -->
          <ElCollapseItem name="2">
            <template #title>
              <ElText tag="b" size="small">{{ t('asset.port') }}</ElText>
            </template>
            <ElScrollbar height="20rem">
              <ElRow v-for="portItem in AssetstatisticsData.Port" :key="portItem.value">
                <ElCol :span="12">
                  <div
                    @click="changeTags('port', portItem.value)"
                    style="display: inline-block; cursor: pointer"
                  >
                    <ElTag effect="light" round size="small">
                      {{ portItem.value }}
                    </ElTag>
                  </div>
                </ElCol>
                <ElCol :span="12" style="text-align: end">
                  <ElText size="small">{{ portItem.number }}</ElText>
                </ElCol>
              </ElRow>
            </ElScrollbar>
          </ElCollapseItem>
          <ElCollapseItem name="3">
            <template #title>
              <ElText tag="b" size="small">{{ t('asset.service') }}</ElText>
            </template>
            <ElScrollbar height="13rem">
              <ElRow v-for="serviceItem in AssetstatisticsData.Service" :key="serviceItem.value">
                <ElCol :span="12">
                  <div
                    @click="changeTags('service', serviceItem.value)"
                    style="display: inline-block; cursor: pointer"
                  >
                    <ElTag effect="light" round size="small">
                      {{ serviceItem.value }}
                    </ElTag>
                  </div>
                </ElCol>
                <ElCol :span="12" style="text-align: end">
                  <ElText size="small">{{ serviceItem.number }}</ElText>
                </ElCol>
              </ElRow>
            </ElScrollbar>
          </ElCollapseItem>
          <ElCollapseItem name="4">
            <template #title>
              <ElText tag="b" size="small">{{ t('asset.products') }}</ElText>
            </template>
            <ElScrollbar height="20rem">
              <ElRow v-for="productItem in AssetstatisticsData.Product" :key="productItem.value">
                <ElCol :span="12">
                  <div
                    @click="changeTags('app', productItem.value)"
                    style="display: inline-block; cursor: pointer"
                  >
                    <ElTag effect="light" round size="small">
                      {{ productItem.value }}
                    </ElTag>
                  </div>
                </ElCol>
                <ElCol :span="12" style="text-align: end">
                  <ElText size="small">{{ productItem.number }}</ElText>
                </ElCol>
              </ElRow>
            </ElScrollbar>
          </ElCollapseItem>
          <!-- 吸附于icon块上方的浮层 -->
          <template v-if="activeNames.includes('5') && selectedIcons.length">
            <div class="icon-selection-float-abs">
              <div class="float-header">
                <span>{{ t('asset.iconSelected') }}</span>
                <span class="icon-selection-close" @click="selectedIcons = []">×</span>
              </div>
              <div class="float-body">
                <span
                  v-for="icon in selectedIcons"
                  :key="icon.icon_hash"
                  class="icon-selection-img"
                >
                  <img :src="'/images/icon/' + icon.icon_hash + '.png'" />
                </span>
                <ElButton
                  class="float-confirm"
                  type="primary"
                  size="small"
                  @click="confirmSelectedIcons"
                >
                  {{ t('asset.confirm') }}
                </ElButton>
              </div>
            </div>
          </template>

          <ElCollapseItem name="5">
            <template #title>
              <div
                class="icon-collapse-title"
                ref="iconTitleRef"
                style="display: inline-block; position: relative"
              >
                <ElText tag="b" size="small">Icon</ElText>
              </div>
            </template>
            <div class="collapse-item-icon-area">
              <ElScrollbar ref="scrollbarRef" height="25rem" @scroll="handleScroll">
                <ElRow style="margin-top: 10px; margin-left: 10px">
                  <ElCol
                    :span="8"
                    v-for="iconItem in AssetstatisticsData.Icon"
                    :key="iconItem.value"
                  >
                    <ElBadge :value="iconItem.number" :max="99" style="font-size: 8px">
                      <ElTooltip :content="iconItem.icon_hash" placement="top-start">
                        <img
                          :src="'/images/icon/' + iconItem.icon_hash + '.png'"
                          alt="Icon"
                          style="width: 30px; height: 30px"
                          :class="{
                            'selected-icon': selectedIcons.some(
                              (i) => i.icon_hash === iconItem.icon_hash
                            )
                          }"
                          @click="toggleSelectIcon(iconItem)"
                        />
                      </ElTooltip>
                    </ElBadge>
                  </ElCol>
                </ElRow>
              </ElScrollbar>
            </div>
          </ElCollapseItem>
        </ElCollapse>
      </ElCard>
    </ElCol>
    <ElCol :span="statisticsHidden ? 24 : 21">
      <ElRow>
        <ElCol :span="24">
          <ElCard>
            <Table
              :columns="allSchemas.tableColumns"
              :data="dataList"
              stripe
              :border="true"
              :loading="loading"
              @filter-change="filterChange"
              :rowStyle="rowstyle"
              :resizable="true"
              @register="tableRegister"
              :headerCellStyle="tableHeaderColor"
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
              :style="{
                fontFamily:
                  '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
              }"
              class="asset-table"
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
    </ElCol>
  </ElRow>
  <div v-if="activeSegment != 'tableSegment'" v-loading="loading">
    <ElRow :gutter="20" type="flex" justify="start" wrap>
      <ElCol
        :span="6"
        v-for="(site, index) in websites"
        :key="index"
        style="margin-bottom: 20px; display: flex; flex-direction: column; height: 350px"
      >
        <ElCard
          :body-style="{ padding: '0px', height: '100%', width: '100%' }"
          style="display: flex; flex-direction: column; height: 100%"
        >
          <!-- 占位符，在图片未加载成功时显示 -->
          <div
            v-if="imageLoadStates[site.ResponseBodyHash] !== 'loaded'"
            style="
              width: 100%;
              height: 100%;
              max-height: 270px;
              background-color: #f0f0f0;
              display: flex;
              justify-content: center;
              align-items: center;
              color: #ccc;
            "
          >
            <span v-if="site.type == 'http'">No pictures available</span>
            <span v-else>{{ site.service }}</span>
          </div>
          <!-- 图片元素，先隐藏，加载成功后再显示 -->
          <img
            v-show="imageLoadStates[site.ResponseBodyHash] === 'loaded'"
            :src="`/images/screenshots/${site.ResponseBodyHash}.png`"
            alt="screenshot"
            style="width: 100%; height: 100%; max-height: 270px"
            @click="handleImageClick(`/images/screenshots/${site.ResponseBodyHash}.png`)"
            @load="handleImageLoad(site.ResponseBodyHash)"
            @error="handleImageError(site.ResponseBodyHash)"
          />
          <!-- 隐藏的预加载图片，用于在不显示时也能触发加载事件 -->
          <img
            v-if="imageLoadStates[site.ResponseBodyHash] !== 'loaded'"
            v-show="false"
            :src="`/images/screenshots/${site.ResponseBodyHash}.png`"
            @load="handleImageLoad(site.ResponseBodyHash)"
            @error="handleImageError(site.ResponseBodyHash)"
          />
          <template #footer>
            <ElRow>
              <ElCol
                style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap"
                :span="20"
              >
                <ElText>
                  {{ site.type == 'http' ? site.title : site.service }}
                </ElText>
              </ElCol>
              <ElCol :span="4" v-if="site.type == 'http'">
                <Icon
                  icon="clarity:circle-solid"
                  :color="getStatusColor(site.statuscode)"
                  :size="6"
                  style="transform: translateY(-30%)"
                />
                <ElText style="margin-left: 5px">{{ site.statuscode }}</ElText>
              </ElCol>
            </ElRow>
            <ElRow>
              <ElCol
                :span="20"
                style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap"
              >
                <ElLink
                  :underline="false"
                  target="_blank"
                  :href="site.type == 'http' ? site.url : site.host"
                  style="font-weight: bold; color: #60a0ef"
                >
                  <Icon
                    icon="carbon:link"
                    :size="16"
                    style="margin-right: 5px"
                    v-if="site.type == 'http'"
                  />
                  {{ site.type == 'http' ? site.url : site.host }}
                </ElLink>
              </ElCol>
              <ElCol :span="4">
                <ElTag>{{ site.port }}</ElTag>
              </ElCol>
            </ElRow>
          </template>
        </ElCard>
      </ElCol>
      <ElCol ::span="24">
        <ElCard>
          <ElPagination
            :loading="loading"
            v-model:pageSize="pageSize"
            v-model:currentPage="currentPage"
            :page-sizes="[20, 40, 60, 100, 200, 400, 600, 1000]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
          />
        </ElCard>
      </ElCol>
    </ElRow>
  </div>

  <Dialog
    v-model="detailVisible"
    :title="t('asset.detail')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    width="50%"
  >
    <AssetDetail2 :id="detailId" :host="detailhost" :ip="detailip" :port="detailport" />
  </Dialog>
</template>

<style lang="less">
.el-popper.is-dark.tagtooltip {
  max-width: 50% !important;
  line-height: 24px;
}
</style>
<style lang="less" scoped>
.icon-selection-float-abs {
  left: 0; // 靠左对齐icon块
  bottom: 100%; // 下边界紧贴icon块上边界
  margin-bottom: 4px; // 微小间隔
  z-index: 220;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px #aaa;
  padding: 8px 14px 10px 14px;
  min-width: 120px;
  border: 1px solid #e3e4e6;
}
.icon-collapse-title {
  position: relative;
  z-index: 10;
}
.selected-icon {
  border: 2px solid #409eff;
  box-shadow: 0 0 4px #409eff;
  border-radius: 5px;
}
.icon-selection-float {
  position: absolute;
  right: 8px;
  top: 2px;
  z-index: 100;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px #aaa;
  padding: 8px 14px 10px 14px;
  min-width: 120px;
  border: 1px solid #e3e4e6;
  /* 可根据实际位置微调 */
}
.float-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #555;
  font-size: 12px;
  margin-bottom: 4px;
}
.icon-selection-close {
  color: #bbb;
  font-size: 18px;
  cursor: pointer;
  margin-left: 12px;
  transition: color 0.2s;
}
.icon-selection-close:hover {
  color: #ff4949;
}
.float-body {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}
.icon-selection-img {
  margin: 0 2px 2px 0;
}
.icon-selection-img img {
  width: 22px;
  height: 22px;
  border-radius: 4px;
  border: 1px solid #e3e3e3;
  margin-right: 2px;
  background: #fafafa;
}
.float-confirm {
  margin-left: 10px;
}

// 优化表格滚动性能
:deep(.asset-table) {
  .el-table__body-wrapper {
    // 启用硬件加速
    transform: translateZ(0);
    -webkit-transform: translateZ(0);
    will-change: scroll-position;
    // 优化滚动性能
    -webkit-overflow-scrolling: touch;
    overflow-anchor: none;

    // 防止滚动时重排
    .el-table__body {
      transform: translateZ(0);
      -webkit-transform: translateZ(0);
    }
  }

  // 优化单元格滚动条性能
  .el-scrollbar {
    .el-scrollbar__wrap {
      // 启用硬件加速
      transform: translateZ(0);
      -webkit-transform: translateZ(0);
      will-change: scroll-position;
      -webkit-overflow-scrolling: touch;
    }
  }

  // 优化滚动条
  .el-scrollbar__bar {
    // 启用硬件加速
    transform: translateZ(0);
    -webkit-transform: translateZ(0);
  }
}

// 优化表格单元格内容
:deep(.asset-table .el-table__cell) {
  // 防止内容变化导致重排
  contain: layout style paint;

  .scrollbar-demo-item {
    // 启用硬件加速
    transform: translateZ(0);
    -webkit-transform: translateZ(0);
  }
}
</style>

<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, onMounted, h, inject, computed, type Ref } from 'vue'
import { ArrowDown, Search } from '@element-plus/icons-vue'
import {
  ElButton,
  ElCol,
  ElInput,
  ElRow,
  ElText,
  ElMessageBox,
  ElMessage,
  ElUpload,
  ElTooltip,
  ElScrollbar,
  ElDropdownItem,
  ElDropdownMenu,
  ElDropdown,
  ElIcon,
  ElBadge,
  ElSwitch,
  ElDrawer,
  ElSpace,
  UploadInstance,
  UploadProps,
  UploadRawFile
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { useIcon } from '@/hooks/web/useIcon'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import {
  checkKeyApi,
  cleanPluginLogApi,
  deletePluginDataApi,
  getPluginDataApi,
  getPluginLogApi,
  reCheckPluginApi,
  reInstallPluginApi,
  runPluginOnceApi,
  uninstallPluginApi,
  updatePluginStatusApi
} from '@/api/plugins'
import detail from './detail.vue'
import { useUserStore } from '@/store/modules/user'

const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const search = ref('')
const handleSearch = () => {
  getList()
}

// 从父组件注入插件市场相关方法
const openMarketDialog = inject<() => void>('openMarketDialog', () => {})
const pendingPluginsCount = inject<Ref<number>>('pendingPluginsCount', ref(0))
const pendingPluginsCountValue = computed(() => pendingPluginsCount.value)

const taskColums = reactive<TableColumn[]>([
  {
    field: 'index',
    label: t('tableDemo.index'),
    type: 'index',
    minWidth: '15'
  },
  {
    field: 'selection',
    type: 'selection',
    minWidth: 55
  },
  {
    field: 'name',
    label: t('plugin.name'),
    formatter: (row, __: TableColumn, value: string) => {
      return (
        <a
          href={`https://plugin.scope-sentry.top/plugin/${row.hash}`}
          style="color: #409EFF; text-decoration: none;"
          target="_blank"
        >
          {value}
        </a>
      )
    }
  },
  {
    field: 'status',
    label: t('plugin.status'),
    minWidth: 100,
    formatter: (row, __: TableColumn, value: boolean | undefined) => {
      const handleStatusChange = async (checked: boolean) => {
        const status = checked
        try {
          await updatePluginStatusApi(row.hash, status)
          // 更新本地数据
          row.status = status
          getList()
        } catch (error) {
          console.error('Error updating plugin status:', error)
        }
      }
      return (
        <ElSwitch
          modelValue={value ?? false}
          onChange={handleStatusChange}
          activeText=""
          inactiveText=""
        />
      )
    }
  },
  {
    field: 'lastTime',
    label: t('task.lastTime'),
    minWidth: 100,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      if (cellValue == '') {
        return '-'
      }
      return cellValue
    }
  },
  {
    field: 'nextTime',
    label: t('task.nextTime'),
    minWidth: 100,
    formatter: (row: Recordable, __: TableColumn, cellValue: string) => {
      if (cellValue == '') {
        return '-'
      }
      if (row.state == false) {
        return '-'
      }
      return cellValue
    }
  },
  {
    field: 'version',
    label: t('plugin.version'),
    minWidth: 50
  },
  {
    field: 'introduction',
    label: t('plugin.introduction'),
    minWidth: 200
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: '300',
    fixed: 'right',
    formatter: (row, __: TableColumn, _: number) => {
      const handleCommand = (command) => {
        switch (command) {
          case 'reinstall':
            reInstallPluginApi('all', row.hash, row.module)
            break
          case 'runOnce':
            confirmRunPluginOnce(row.hash)
            break
          // case 'recheck':
          //   reCheckPluginApi('all', row.hash, row.module)
          //   break
          // case 'uninstall':
          //   uninstallPluginApi('all', row.hash, row.module)
          //   break
        }
      }
      const retestAndDeleteDropdown = h(
        ElDropdown,
        {
          onCommand: handleCommand
        },
        {
          default: () =>
            h(
              ElButton,
              {
                style: { outline: 'none', boxShadow: 'none' }
              },
              () => [
                t('common.operation'), // 下拉菜单触发按钮文字
                h(
                  ElIcon,
                  {},
                  () => h(ArrowDown) // 向下箭头图标
                )
              ]
            ),
          dropdown: () =>
            h(ElDropdownMenu, null, () => {
              return [
                h(ElDropdownItem, { command: 'reinstall' }, () => t('plugin.reInstall')),
                h(ElDropdownItem, { command: 'runOnce' }, () => t('plugin.runOnce'))
              ]
            })
        }
      )
      return (
        <>
          {retestAndDeleteDropdown}
          <BaseButton
            type="warning"
            style={{ marginLeft: '10px' }}
            onClick={() => openLogDialogVisible(row)}
          >
            {t('common.log')}
          </BaseButton>
          <BaseButton
            type="info"
            style={{ marginLeft: '10px' }}
            onClick={() => confirmCleanLog(row.hash, row.module)}
          >
            {t('common.cleanLog')}
          </BaseButton>
          <BaseButton type="success" onClick={() => editPlugin(row.id, row.hash)}>
            {t('common.edit')}
          </BaseButton>
          <BaseButton
            type="danger"
            onClick={() => confirmDelete(row.hash, row.module)}
            disabled={row.isSystem}
          >
            {t('common.delete')}
          </BaseButton>
        </>
      )
    }
  }
])

const moduleColorMap = {
  TargetHandler: '#2243dda6', // 浅红色
  SubdomainScan: '#FF9B85', // 更深的浅橙色
  SubdomainSecurity: '#FFFFBA', // 浅黄色
  PortScanPreparation: '#BAFFB3', // 浅绿色
  PortScan: '#BAE1FF', // 浅蓝色
  AssetMapping: '#e3ffba', // 浅粉红色
  URLScan: '#D1BAFF', // 浅紫色
  WebCrawler: '#FFABAB', // 浅红
  DirScan: '#3ccde6', // 选择浅桃色
  VulnerabilityScan: '#FF677D', // 浅粉色
  AssetHandle: '#B2E1FF', // 浅青色
  PortFingerprint: '#ffb5e4', // 更亮的浅橙色
  URLSecurity: '#FFE4BA', // 浅米色
  PassiveScan: '#A2DFF7'
}

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    // 传递 type="server" 参数
    const res = await getPluginDataApi(search.value, currentPage.value, pageSize.value, 'server')
    return {
      list: res.data.list,
      total: res.data.total
    }
  },
  immediate: false
})
const { loading, dataList, total, currentPage, pageSize } = tableState
pageSize.value = 20
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
const dialogVisible = ref(false)

let DialogTitle = t('plugin.new')
const closeDialog = () => {
  dialogVisible.value = false
}

const confirmDeleteSelect = async () => {
  ElMessageBox({
    title: 'Delete',
    draggable: true
  }).then(async () => {
    await delSelect()
  })
}

const confirmDelete = async (hash: string, module: string) => {
  ElMessageBox({
    title: 'Delete',
    draggable: true
  }).then(async () => {
    await del(hash, module)
  })
}

const confirmCleanLog = async (hash: string, module: string) => {
  ElMessageBox({
    title: 'Clean Log',
    message: 'Are you sure you want to clean the logs?',
    draggable: true
  }).then(async () => {
    await cleanPluginLogApi(module, hash, 'server')
  })
}

const confirmRunPluginOnce = async (hash: string) => {
  ElMessageBox({
    title: t('plugin.runOnce'),
    message: t('plugin.runOnceConfirm'),
    draggable: true
  })
    .then(async () => {
      try {
        await runPluginOnceApi(hash)
        ElMessage.success(t('plugin.runOnceSuccess'))
      } catch (error) {
        console.error('运行插件失败:', error)
        ElMessage.error(t('plugin.runOnceFailed'))
      }
    })
    .catch(() => {
      // 用户取消操作
    })
}

const delLoading = ref(false)
const del = async (hash: string, module: string) => {
  delLoading.value = true
  try {
    const res = await deletePluginDataApi([{ hash, module }])
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  const deleteData = selectedRows.map((row) => ({
    hash: row.hash,
    module: row.module
  }))

  delLoading.value = true
  try {
    const res = await deletePluginDataApi(deleteData)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}

const addPlugin = async () => {
  id.value = ''
  dialogVisible.value = true
}

const id = ref('')
const hash = ref('')
const editPlugin = async (data, h) => {
  id.value = data
  hash.value = h
  DialogTitle = t('common.edit')
  dialogVisible.value = true
}

const keyDialogVisible = ref(false)
const pluginKey = ref('')

const LoadPluginKey = () => {
  const key = localStorage.getItem(`plugin_key`) as string
  if (!key) {
    keyDialogVisible.value = true
  }
  pluginKey.value = key || ''
}

const savePluginKey = async () => {
  if (pluginKey.value) {
    const res = await checkKeyApi(pluginKey.value)
    if (res.code == 200) {
      localStorage.setItem('plugin_key', pluginKey.value)
      keyDialogVisible.value = false
    }
  }
}

const handlePluginKeyChange = () => {
  if (pluginKey.value) {
    localStorage.setItem('plugin_key', pluginKey.value)
  } else {
    localStorage.removeItem('plugin_key')
  }
}

onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
  LoadPluginKey()
  // 当组件挂载时（切换到服务端插件 tab 时）加载数据
  getList()
})

const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.7
}
const logDialogVisible = ref(false)
const closeLogDialogVisible = () => {
  logDialogVisible.value = false
  logSearchText.value = ''
}
const logContent = ref('')
const logSearchText = ref('')
const logScrollbarRef = ref<InstanceType<typeof ElScrollbar>>()
const autoScroll = ref(true)

const logModule = ref('')
const logHash = ref('')
const openLogDialogVisible = async (data) => {
  logModule.value = data.module
  logHash.value = data.hash
  await refreshLog()
  logDialogVisible.value = true
  // 等待DOM更新后滚动到底部
  setTimeout(() => {
    scrollToBottom()
  }, 100)
}

const refreshLog = async () => {
  const res = await getPluginLogApi('', logHash.value, 'server')
  logContent.value = res.data.data
  if (autoScroll.value) {
    setTimeout(() => {
      scrollToBottom()
    }, 50)
  }
}

const cleanLog = async () => {
  await cleanPluginLogApi(logModule.value, logHash.value, 'server')
  logContent.value = ''
  ElMessage.success(t('common.cleanLog') + ' ' + t('common.success'))
}

const scrollToBottom = () => {
  if (logScrollbarRef.value) {
    const scrollbar = logScrollbarRef.value
    const scrollbarEl = scrollbar.$el
    const wrapEl = scrollbarEl?.querySelector('.el-scrollbar__wrap')
    if (wrapEl) {
      wrapEl.scrollTop = wrapEl.scrollHeight
    }
  }
}

const copyLog = async () => {
  if (logContent.value) {
    try {
      await navigator.clipboard.writeText(logContent.value)
      ElMessage.success(t('common.copySuccess'))
    } catch (error) {
      ElMessage.error(t('common.copyFailed'))
    }
  }
}

// 过滤后的日志内容（用于搜索高亮）
const filteredLogContent = computed(() => {
  if (!logSearchText.value || !logContent.value) {
    return logContent.value
  }
  const searchText = logSearchText.value
  const lines = logContent.value.split('\n')
  return lines.filter((line) => line.toLowerCase().includes(searchText.toLowerCase())).join('\n')
})

/**
 * HTML转义函数 - 防止XSS攻击
 */
const escapeHtml = (text: string): string => {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

/**
 * 验证颜色值是否安全
 * 只允许十六进制颜色 (#rrggbb) 或 rgb/rgba 格式
 */
const isValidColor = (color: string): boolean => {
  if (!color) return false
  color = color.trim()

  // 检查是否包含危险字符
  if (/[<>'"`]/.test(color) || /javascript:/i.test(color) || /on\w+=/i.test(color)) {
    return false
  }

  // 验证十六进制颜色格式 #rrggbb 或 #rgb
  const hexPattern = /^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$/
  if (hexPattern.test(color)) {
    return true
  }

  // 验证 rgb/rgba 格式
  const rgbPattern = /^rgba?\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*(,\s*[\d.]+\s*)?\)$/
  if (rgbPattern.test(color)) {
    return true
  }

  return false
}

/**
 * 解析日志中的 <color> 标签
 * 格式: <color value="#ff6b6b">文字</color> 或 <color value="#51cf66" bold>文字</color>
 *
 * @param text 原始日志文本
 * @returns 解析后的HTML字符串（已转义，安全）
 */
const parseColorTags = (text: string): string => {
  if (!text) return ''

  // 匹配 <color value="颜色值" bold>内容</color> 标签
  // 支持格式:
  //   <color value="#ff6b6b">文字</color>
  //   <color value="#51cf66" bold>文字</color>
  const colorTagRegex =
    /<color\s+value=["']([^"']+)["'](?:\s+bold)?\s*>((?:[^<]|<(?!\/color>))*?)<\/color>/gi

  return text.replace(colorTagRegex, (match, colorValue, content) => {
    // 验证颜色值
    if (!isValidColor(colorValue)) {
      // 如果颜色值不安全，返回转义后的原始文本
      return escapeHtml(match)
    }

    // 检查是否有 bold 属性
    const isBold = /\s+bold\s*>/i.test(match)

    // 转义内容，防止XSS
    const safeContent = escapeHtml(content)

    // 构建样式
    const styles = [`color: ${colorValue}`]
    if (isBold) {
      styles.push('font-weight: bold')
    }

    // 返回安全的HTML标签
    return `<span style="${styles.join('; ')}">${safeContent}</span>`
  })
}

/**
 * 渲染日志内容（解析颜色标签 + 搜索高亮）
 */
const renderLogContent = (text: string): string => {
  if (!text) return ''

  // 1. 先解析颜色标签
  let result = parseColorTags(text)

  // 2. 转义其他所有HTML标签（除了我们已经生成的span标签）
  // 先标记我们生成的span标签
  const tempMarker = '___COLOR_SPAN___'
  const spanMatches: string[] = []
  result = result.replace(/<span style="[^"]+">[^<]*<\/span>/g, (match) => {
    spanMatches.push(match)
    return tempMarker + (spanMatches.length - 1) + tempMarker
  })

  // 转义所有剩余的HTML
  result = escapeHtml(result)

  // 恢复我们的span标签
  spanMatches.forEach((match, index) => {
    result = result.replace(tempMarker + index + tempMarker, match)
  })

  // 3. 如果有搜索文本，再应用搜索高亮
  if (logSearchText.value) {
    const searchText = logSearchText.value
    const escapedSearchText = escapeHtml(searchText)
    const regex = new RegExp(`(${escapedSearchText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')

    result = result.replace(regex, (match) => {
      // 如果匹配的内容已经在span标签内，在span内部添加mark
      if (match.includes('<span')) {
        return match.replace(/(<span[^>]*>)(.*?)(<\/span>)/g, (_, openTag, content, closeTag) => {
          const contentRegex = new RegExp(
            `(${escapedSearchText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`,
            'gi'
          )
          const highlightedContent = content.replace(contentRegex, '<mark>$1</mark>')
          return `${openTag}${highlightedContent}${closeTag}`
        })
      }
      return `<mark>${match}</mark>`
    })
  }

  return result
}

const userStore = useUserStore()
const uploadHeaders = { Authorization: `${userStore.getToken}` }
const upload = ref<UploadInstance>()
const handleExceed: UploadProps['onExceed'] = (files) => {
  upload.value!.clearFiles()
  const file = files[0] as UploadRawFile
  upload.value!.handleStart(file)
}

const handleUploadSuccess = (response) => {
  console.log(response)
  if (response.code === 200) {
    ElMessage.success('Upload succes')
  } else {
    ElMessage.error(response.message)
  }
  if (response.code == 505) {
    localStorage.removeItem('plugin_key')
  }
  getList()
  upload.value?.clearFiles()
}
const handleFileChange = (_file, fileList) => {
  if (fileList.length > 0) {
    upload.value!.submit()
  }
}
</script>

<template>
  <ContentWrap>
    <ElRow>
      <ElCol :span="1">
        <ElText class="mx-1" style="position: relative; top: 8px">{{ t('plugin.name') }}:</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="search" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5" style="position: relative; left: 16px">
        <ElButton type="primary" :icon="searchicon" style="height: 100%" @click="handleSearch"
          >Search</ElButton
        >
      </ElCol>
      <ElCol :span="1" style="position: relative; left: 32px">
        <ElText class="mx-1" style="position: relative; top: 8px">{{ t('plugin.key') }}:</ElText>
      </ElCol>
      <ElCol :span="5" style="position: relative; left: 32px">
        <ElInput
          v-model="pluginKey"
          :placeholder="t('plugin.key')"
          style="height: 38px"
          @blur="handlePluginKeyChange"
        />
      </ElCol>
    </ElRow>
    <ElRow :gutter="16" class="mt-4">
      <ElCol :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <div class="flex flex-wrap gap-3 items-center">
          <BaseButton type="primary" @click="addPlugin">
            {{ t('plugin.new') }}
          </BaseButton>

          <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">
            {{ t('plugin.delete') }}
          </BaseButton>

          <ElBadge
            :value="pendingPluginsCountValue"
            :hidden="pendingPluginsCountValue === 0"
            :max="99"
          >
            <BaseButton type="success" @click="openMarketDialog">
              {{ t('plugin.market') }}
            </BaseButton>
          </ElBadge>
          <ElUpload
            ref="upload"
            class="flex items-center"
            :action="'/api/plugin/import?key=' + (pluginKey || '')"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :limit="1"
            :on-exceed="handleExceed"
            :auto-upload="false"
            @change="handleFileChange"
          >
            <template #trigger>
              <BaseButton>
                <template #icon>
                  <Icon icon="iconoir:upload" />
                </template>
                {{ t('plugin.import') }}
              </BaseButton>
            </template>
          </ElUpload>
        </div>
      </ElCol>
    </ElRow>
    <div style="position: relative; top: 12px">
      <Table
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
        :columns="taskColums"
        :data="dataList"
        stripe
        :border="true"
        :loading="loading"
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
  <ElDrawer v-model="dialogVisible" :title="DialogTitle" size="50%" direction="rtl">
    <detail :closeDialog="closeDialog" :getList="getList" :id="id" tp="server" :hash="hash" />
  </ElDrawer>
  <ElDrawer
    v-model="logDialogVisible"
    :title="t('node.log')"
    size="80%"
    direction="rtl"
    :with-header="true"
  >
    <template #header>
      <div style="display: flex; align-items: center; justify-content: space-between; width: 100%">
        <span style="font-weight: 500">{{ t('node.log') }}</span>
        <ElSpace>
          <ElInput
            v-model="logSearchText"
            :placeholder="t('common.search')"
            clearable
            style="width: 200px"
          >
            <template #prefix>
              <ElIcon><Search /></ElIcon>
            </template>
          </ElInput>
          <ElSwitch v-model="autoScroll" :active-text="t('common.autoScroll')" inactive-text="" />
        </ElSpace>
      </div>
    </template>
    <div style="display: flex; flex-direction: column; height: 100%">
      <ElScrollbar ref="logScrollbarRef" style="flex: 1; height: 0">
        <div
          style="
            padding: 16px;
            background: #1e1e1e;
            color: #d4d4d4;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 13px;
            line-height: 1.6;
            min-height: 100%;
            white-space: pre-wrap;
            word-wrap: break-word;
          "
        >
          <div v-if="logSearchText && filteredLogContent">
            <div
              v-for="(line, index) in filteredLogContent.split('\n')"
              :key="index"
              v-html="renderLogContent(line)"
            ></div>
          </div>
          <div v-else-if="logContent">
            <div
              v-for="(line, index) in logContent.split('\n')"
              :key="index"
              v-html="renderLogContent(line)"
            ></div>
          </div>
          <div v-else style="color: #888; text-align: center; padding: 40px">
            {{ t('common.noData') }}
          </div>
        </div>
      </ElScrollbar>
      <div
        style="
          padding: 16px;
          border-top: 1px solid var(--el-border-color);
          display: flex;
          justify-content: space-between;
          align-items: center;
        "
      >
        <div style="color: var(--el-text-color-secondary); font-size: 12px">
          <span v-if="logSearchText">
            {{ t('common.searchResult') }}:
            {{ filteredLogContent.split('\n').filter((l) => l).length }}
            {{ t('common.lines') }}
          </span>
          <span v-else-if="logContent">
            {{ logContent.split('\n').filter((l) => l).length }} {{ t('common.lines') }}
          </span>
        </div>
        <ElSpace>
          <BaseButton @click="refreshLog" type="primary">
            {{ t('common.refresh') }}
          </BaseButton>
          <BaseButton @click="copyLog" type="info">
            {{ t('common.copy') }}
          </BaseButton>
          <BaseButton @click="cleanLog" type="danger">{{ t('common.cleanLog') }}</BaseButton>
          <BaseButton @click="closeLogDialogVisible">{{ t('common.off') }}</BaseButton>
        </ElSpace>
      </div>
    </div>
  </ElDrawer>
  <Dialog
    v-model="keyDialogVisible"
    :title="t('plugin.key')"
    center
    width="30%"
    style="max-width: 400px; height: 200px"
  >
    <div class="flex flex-col gap-2">
      <el-tooltip class="item" effect="dark" :content="t('plugin.keyMsg')" placement="top">
        <ElInput v-model="pluginKey" />
      </el-tooltip>
      <BaseButton @click="savePluginKey" type="primary" class="w-full">Save</BaseButton>
    </div>
  </Dialog>
</template>

<style scoped lang="less">
// 确保搜索高亮的mark标签在深色背景下可见
:deep(mark) {
  background-color: #ffd700;
  color: #000;
  padding: 2px 4px;
  border-radius: 2px;
}
</style>

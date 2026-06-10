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
  ElTag,
  ElTooltip,
  ElScrollbar,
  UploadInstance,
  UploadProps,
  UploadRawFile,
  ElUpload,
  ElMessage,
  ElDropdownItem,
  ElDropdownMenu,
  ElDropdown,
  ElIcon,
  ElBadge,
  ElDrawer,
  ElSpace,
  ElSwitch
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { useIcon } from '@/hooks/web/useIcon'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import {
  checkKeyApi,
  cleanAllPluginLogApi,
  cleanPluginLogApi,
  deletePluginDataApi,
  getPluginDataApi,
  getPluginLogApi,
  reCheckPluginApi,
  reInstallPluginApi,
  uninstallPluginApi
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
    minWidth: 55,
    selectable: (row) => !row.isSystem
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
    field: 'module',
    label: t('plugin.module'),
    formatter: (_, __: TableColumn, value: string) => {
      const color = moduleColorMap[value] || '#FFFFFF' // 默认白色
      return <ElTag style={{ backgroundColor: color, color: '#000' }}>{value}</ElTag>
    }
  },
  {
    field: 'version',
    label: t('plugin.version'),
    minWidth: 100
  },
  {
    field: 'parameter',
    label: t('plugin.parameter'),
    formatter: (row, __: TableColumn, value: string) => {
      return (
        <ElTooltip content={row.help} placement="top" effect="light">
          <span style="cursor: pointer;">{value}</span>
        </ElTooltip>
      )
    }
  },
  {
    field: 'introduction',
    label: t('plugin.introduction'),
    minWidth: 200
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    width: 450,
    fixed: 'right',
    formatter: (row, __: TableColumn, _: number) => {
      const handleCommand = (command) => {
        switch (command) {
          case 'reinstall':
            reInstallPluginApi('all', row.hash, row.module)
            break
          case 'recheck':
            reCheckPluginApi('all', row.hash, row.module)
            break
          case 'uninstall':
            uninstallPluginApi('all', row.hash, row.module)
            break
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
                h(ElDropdownItem, { command: 'recheck' }, () => t('plugin.reCheck')),
                h(ElDropdownItem, { command: 'uninstall' }, () => t('plugin.uninstall'))
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
          <BaseButton type="success" onClick={() => editPlugin(row.id)}>
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
    // 不传递 type 参数，默认为扫描端插件
    const res = await getPluginDataApi(search.value, currentPage.value, pageSize.value, 'scan')
    return {
      list: res.data.list,
      total: res.data.total
    }
  },
  immediate: true
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
    await cleanPluginLogApi(module, hash)
  })
}

const confirmCleanAllLog = async () => {
  ElMessageBox({
    title: 'Clean All Plugin Logs',
    message: 'Are you sure you want to clean all plugin logs?',
    draggable: true
  }).then(async () => {
    await cleanAllPluginLogApi()
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

const editPlugin = async (data) => {
  id.value = data
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
  const res = await getPluginLogApi(logModule.value, logHash.value)
  logContent.value = res.data.data
  if (autoScroll.value) {
    setTimeout(() => {
      scrollToBottom()
    }, 50)
  }
}

const cleanLog = async () => {
  await cleanPluginLogApi(logModule.value, logHash.value)
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

// 高亮搜索关键词
const highlightLog = (text: string) => {
  if (!logSearchText.value || !text) {
    return text
  }
  const searchText = logSearchText.value
  const regex = new RegExp(`(${searchText})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
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

// 暴露刷新列表方法给父组件
defineExpose({
  refreshList: getList
})
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

          <BaseButton type="warning" @click="confirmCleanAllLog">
            {{ t('common.cleanAllLog') }}
          </BaseButton>

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
  <ElDrawer v-model="dialogVisible" size="50%" direction="rtl">
    <template #title>
      <div style="display: flex; align-items: center; gap: 16px; width: 100%">
        <span style="font-weight: 500; white-space: nowrap">{{ DialogTitle }}</span>
        <span style="color: #f56c6c; font-size: 12px; font-weight: normal; line-height: 1.4">
          {{ t('plugin.parameterConfigTip') }}
        </span>
      </div>
    </template>
    <detail :closeDialog="closeDialog" :getList="getList" :id="id" tp="scan" />
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
              v-html="highlightLog(line)"
            ></div>
          </div>
          <div v-else-if="logContent">
            {{ logContent }}
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

<style scoped lang="less"></style>

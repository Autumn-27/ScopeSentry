<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, h, watch } from 'vue'
import { ElCol, ElRow, ElScrollbar, ElTag, ElTooltip } from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
// import { useIcon } from '@/hooks/web/useIcon'
import { onMounted } from 'vue'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import Configuration from './components/Configuration.vue'
import plugin from './components/plugin.vue'
import { getNodeDataApi, deleteNodeApi, getNodeLogApi, restartNodeApi } from '@/api/node'
const { t } = useI18n()
const nodeColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'name',
    label: t('node.nodeName'),
    minWidth: 20
  },
  {
    field: 'maxTaskNum',
    label: t('configuration.maxTaskNum'),
    minWidth: 10,
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(
        ElTag,
        {
          type: 'info'
        },
        () => cellValue
      )
    }
  },
  {
    field: 'running',
    label: t('node.taskCount'),
    minWidth: 10,
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(
        ElTag,
        {
          round: true,
          effect: 'plain',
          hit: true
        },
        () => cellValue
      )
    }
  },
  {
    field: 'finished',
    label: t('node.finished'),
    minWidth: 10,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      return h(
        ElTag,
        {
          round: true,
          effect: 'plain',
          hit: true
        },
        () => cellValue
      )
    }
  },
  {
    field: 'cpuNum',
    label: t('node.nodeUsageCpu'),
    minWidth: 20,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      let numericValue = parseFloat(cellValue)
      numericValue = parseFloat(numericValue.toFixed(2))
      return h(
        ElTag,
        {
          round: true,
          effect: 'plain',
          hit: true,
          type: numericValue < 50 ? '' : numericValue < 80 ? 'warning' : 'danger'
        },
        () => numericValue + '%'
      )
    }
  },
  {
    field: 'memNum',
    label: t('node.nodeUsageMemory'),
    minWidth: 20,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      let numericValue = parseFloat(cellValue)
      numericValue = parseFloat(numericValue.toFixed(2))
      return h(
        ElTag,
        {
          round: true,
          effect: 'plain',
          hit: true,
          type: numericValue < 50 ? '' : numericValue < 80 ? 'warning' : 'danger'
        },
        () => numericValue + '%'
      )
    }
  },
  {
    field: 'state',
    label: t('node.nodeStatus'),
    minWidth: 20,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      return h(
        ElTag,
        {
          type: cellValue === '1' ? 'success' : cellValue === '2' ? 'warning' : 'danger',
          effect: 'light',
          hit: true
        },
        () =>
          cellValue === '1'
            ? t('node.statusRun')
            : cellValue === '2'
              ? t('node.statusStop')
              : t('node.statusError')
      )
    }
  },
  {
    field: 'updateTime',
    label: t('node.updateTime'),
    minWidth: 20
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 30,
    formatter: (row, __: TableColumn, _: number) => {
      console.log(row)
      return (
        <>
          <BaseButton type="warning" size="small" onClick={() => openPlugin(row.name)}>
            {t('node.plugin')}
          </BaseButton>
          <BaseButton type="success" size="small" onClick={() => openLogDialogVisible(row)}>
            {t('node.log')}
          </BaseButton>
          <BaseButton type="primary" size="small" onClick={() => openConfig(row)}>
            {t('common.config')}
          </BaseButton>
          <ElTooltip content={t('node.restartMsg')}>
            <BaseButton type="danger" size="small" onClick={() => restartNode(row.name)}>
              {t('node.restart')}
            </BaseButton>
          </ElTooltip>
        </>
      )
    }
  }
])
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getNodeDataApi()
    return {
      list: res.data.list
    }
  }
})
const { loading, dataList, currentPage, pageSize } = tableState
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
const dialogVisible = ref(false)
const closeDialog = () => {
  dialogVisible.value = false
}
const detailData = reactive({
  name: '',
  maxTaskNum: '',
  state: '',
  ModulesConfig: ''
})
const openConfig = async (data) => {
  detailData.name = data.name
  detailData.maxTaskNum = data.maxTaskNum
  detailData.ModulesConfig = data.modulesConfig
  detailData.state = data.state
  dialogVisible.value = true
}
const restartNode = async (name) => {
  await restartNodeApi(name)
}
const delLoading = ref(false)
const names = ref<string[]>([])
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  names.value = selectedRows.map((row) => row.name)
  delLoading.value = true
  try {
    const res = await deleteNodeApi(names.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmDelete = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}
const logDialogVisible = ref(false)
const closeLogDialogVisible = () => {
  logDialogVisible.value = false
}
const logContent = ref('')
const scrollbarRef = ref<InstanceType<typeof ElScrollbar>>()
const openLogDialogVisible = async (data) => {
  const res = await getNodeLogApi(data.name)
  logContent.value = res.data.logs
  logDialogVisible.value = true
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

const nodeName = ref('')
const pluginDialogVisible = ref(false)
const openPlugin = async (data) => {
  nodeName.value = data
  pluginDialogVisible.value = true
}
const closepluginDialogVisible = () => {
  pluginDialogVisible.value = false
}
</script>

<template>
  <ContentWrap>
    <!-- <ElRow :gutter="20" style="margin-bottom: 15px">
      <ElCol :span="1.5">
        <ElText class="mx-1" style="position: relative; top: 8px">{{ t('node.nodeName') }}:</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="search" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5" style="position: relative; left: 16px">
        <ElButton type="primary" :icon="searchicon" style="height: 100%" @click="handleSearch"
          >Search</ElButton
        >
      </ElCol>
    </ElRow> -->
    <ElRow>
      <ElCol style="position: relative; top: 16px">
        <div class="mb-10px">
          <BaseButton type="danger" :loading="delLoading" @click="confirmDelete">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </ElCol>
    </ElRow>
    <div style="position: relative; top: 12px">
      <Table
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
        :columns="nodeColums"
        :data="dataList"
        stripe
        :border="true"
        :loading="loading"
        :resizable="true"
        @register="tableRegister"
        :headerCellStyle="tableHeaderColor"
        :style="{
          fontFamily:
            '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
        }"
      />
    </div>
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="$t('common.config')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="maxHeight"
  >
    <Configuration :closeDialog="closeDialog" :nodeConfForm="detailData" :getList="getList" />
  </Dialog>
  <Dialog
    v-model="logDialogVisible"
    :title="t('node.log')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="maxHeight"
  >
    <ElScrollbar ref="scrollbarRef">
      <pre v-if="logContent">{{ logContent }}</pre>
    </ElScrollbar>
    <template #footer>
      <BaseButton @click="closeLogDialogVisible">{{ t('common.off') }}</BaseButton>
    </template>
  </Dialog>

  <Dialog
    v-model="pluginDialogVisible"
    :title="t('node.plugin')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="maxHeight"
  >
    <plugin :closeDialog="closepluginDialogVisible" :name="nodeName" />
  </Dialog>
</template>

<script setup lang="ts">
import PanelGroup from './components/PanelGroup.vue'
import {
  ElRow,
  ElCol,
  ElCard,
  ElProgress,
  ElText,
  ElTooltip,
  ElButton,
  ElPopconfirm,
  ElForm,
  ElFormItem,
  ElInput
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { ElTag } from 'element-plus'
import { Dialog } from '@/components/Dialog'
import { ref, reactive, h, onBeforeUnmount, Ref } from 'vue'
import { getNodeDataApi } from '@/api/node'
import { getTaskDataApi } from '@/api/task'
import { useI18n } from '@/hooks/web/useI18n'
import { UPDATEsYSTEMApi, getVersionDataApi } from '@/api/dashboard/analysis'
import { checkKeyApi } from '@/api/plugins'

const { t } = useI18n()

const loading = ref(true)

const nodeColumns = reactive<TableColumn[]>([
  {
    field: 'name',
    label: t('node.nodeName')
  },
  {
    field: 'running',
    label: t('node.taskCount'),
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(
        ElTag,
        {
          round: true,
          effect: 'dark'
        },
        () => cellValue
      )
    }
  },
  {
    field: 'finished',
    label: t('node.finished'),
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(
        ElTag,
        {
          round: true,
          effect: 'dark'
        },
        () => cellValue
      )
    }
  },
  {
    field: 'state',
    label: t('node.nodeStatus'),
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      return h(
        ElTag,
        {
          type: cellValue === '1' ? 'success' : cellValue === '2' ? 'warning' : 'danger',
          effect: 'dark'
        },
        () =>
          cellValue == '1'
            ? t('node.statusRun')
            : cellValue == '2'
              ? t('node.statusStop')
              : t('node.statusError')
      )
    }
  }
])

const taskColums = reactive<TableColumn[]>([
  {
    field: 'name',
    label: t('task.taskName')
  },
  {
    field: 'taskNum',
    label: t('task.taskCount'),
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(
        ElTag,
        {
          round: true,
          effect: 'dark'
        },
        () => cellValue
      )
    }
  },
  {
    field: 'progress',
    label: t('task.taskProgress'),
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      return h(ElProgress, {
        percentage: cellValue,
        type: 'line',
        striped: true,
        status: cellValue < 100 ? '' : 'success',
        stripedFlow: cellValue < 100 ? true : false
      })
    }
  },
  {
    field: 'creatTime',
    label: t('task.createTime')
  }
])

const nodeUsageColumns = reactive<TableColumn[]>([
  {
    field: 'name',
    label: t('node.nodeName')
  },
  {
    field: 'cpuNum',
    label: t('node.nodeUsageCpu'),
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      let numericValue = parseFloat(cellValue)
      numericValue = parseFloat(numericValue.toFixed(2))
      return h(ElProgress, {
        percentage: numericValue,
        type: 'dashboard',
        color: numericValue < 50 ? '#26a33f' : numericValue <= 80 ? '#fe9900' : '#df2800'
      })
    }
  },
  {
    field: 'memNum',
    label: t('node.nodeUsageMemory'),
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      let numericValue = parseFloat(cellValue)
      numericValue = parseFloat(numericValue.toFixed(2))
      return h(ElProgress, {
        percentage: numericValue,
        type: 'dashboard',
        color: numericValue < 50 ? '#26a33f' : numericValue < 80 ? '#fe9900' : '#df2800'
      })
    }
  }
])

const versionColumns = reactive<TableColumn[]>([
  {
    field: 'name',
    label: t('common.name')
  },
  {
    field: 'cversion',
    label: t('common.cversion')
  },
  {
    field: 'lversion',
    label: t('common.lversion'),
    formatter: (row: Recordable, __: TableColumn, cellValue: string) => {
      if (row.cversion != row.lversion) {
        updateFlag.value = true
        const msgArray = row.msg.split('\\n')
        let content = ''
        msgArray.forEach((line) => {
          content += `<div>${line}</div>`
        })
        return h(
          ElTooltip,
          {
            placement: 'top',
            content: content,
            rawContent: true
          },
          [
            h(
              ElText,
              {
                type: 'danger'
              },
              cellValue
            )
          ]
        )
      } else {
        return h(ElText, cellValue)
      }
    }
  }
])

let nodeUsageData: Ref<{ name: string; cpuNum: number; memNum: number }[]> = ref([])

const nodeData = ref<
  {
    name: string
    running: number
    finished: number
    state: number
    cpuNum: number
    memNum: number
  }[]
>([])

const getNodeState = async () => {
  try {
    const res = await getNodeDataApi()
    if (res && res.data && Array.isArray(res.data.list)) {
      nodeData.value = res.data.list.map((node) => ({
        name: node.name,
        running: node.running,
        state: node.state,
        finished: node.finished,
        cpuNum: node.cpuNum,
        memNum: node.memNum
      }))
      nodeUsageData.value = reactive(
        res.data.list.map((node) => ({
          name: node.name,
          cpuNum: node.cpuNum,
          memNum: node.memNum
        }))
      )
    }
  } catch (error) {
    console.error('Error fetching node data:', error)
  } finally {
    // 不论请求成功或失败，都会执行的代码块
    loading.value = false
  }
}

const taskData = ref<
  {
    name: string
    taskNum: string
    progress: string
    creatTime: string
  }[]
>([])

const getTaskData = async () => {
  const res = await getTaskDataApi('', 1, 10)
  console.log(res)
  taskData.value = reactive(
    res.data.list.map((task) => ({
      name: task.name,
      taskNum: task.taskNum,
      progress: task.progress,
      creatTime: task.creatTime
    }))
  )
}

const versionData = ref<
  {
    name: string
    cversion: string
    lversion: string
    msg: string
  }[]
>([])

const getVersionData = async () => {
  const res = await getVersionDataApi()
  console.log(res)
  versionData.value = reactive(
    res.data.list.map((v) => ({
      name: v.name,
      cversion: v.cversion,
      lversion: v.lversion,
      msg: v.msg
    }))
  )
}

const getAllApi = async () => {
  await Promise.all([getNodeState(), getTaskData()])
  loading.value = false
}
getVersionData()
getAllApi()
const refreshInterval = setInterval(getAllApi, 10000)

onBeforeUnmount(() => {
  clearInterval(refreshInterval)
})

const pluginKey = ref('')
const updateSystem = async () => {
  const key = localStorage.getItem(`plugin_key`) as string
  if (!key) {
    keyDialogVisible.value = true
  } else {
    UpdatedialogVisible.value = true
  }
}
const savePluginKey = async () => {
  if (pluginKey.value) {
    const res = await checkKeyApi(pluginKey.value)
    if (res.code == 200) {
      localStorage.setItem('plugin_key', pluginKey.value)
      keyDialogVisible.value = false
      UpdatedialogVisible.value = true
    }
  }
}
const updateFlag = ref(false)

const form = ref({
  server: 'https://github.com/Autumn-27/ScopeSentry/archive/refs/heads/main.zip',
  scan: ''
})
const UpdatedialogVisible = ref(false)
const keyDialogVisible = ref(false)
async function handleSubmit() {
  const key = localStorage.getItem(`plugin_key`) as string
  const res = await UPDATEsYSTEMApi(form.value.server, form.value.scan, key)
  if (res.code == 505) {
    localStorage.removeItem('plugin_key')
  }
}
</script>

<template>
  <PanelGroup />
  <ElRow :gutter="20" justify="space-between">
    <ElCol :xl="12" :lg="12" :md="24" :sm="24" :xs="24">
      <ElCard shadow="hover" class="mb-20px">
        <template #header>
          <span>{{ t('dashboard.nodeInfo') }}</span>
        </template>
        <Table :columns="nodeColumns" :data="nodeData" stripe :border="false" :height="250" />
      </ElCard>
    </ElCol>
    <ElCol :xl="12" :lg="12" :md="24" :sm="24" :xs="24">
      <ElCard shadow="hover" class="mb-20px">
        <template #header>
          <span>{{ t('dashboard.taskInfo') }}</span>
        </template>
        <Table
          :columns="taskColums"
          :data="taskData"
          stripe
          :border="false"
          :height="250"
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
        />
      </ElCard>
    </ElCol>
    <ElCol :span="12">
      <ElCard shadow="hover" class="mb-25px">
        <template #header>
          <div>
            <span>{{ t('node.nodeUsageStatus') }}</span>
          </div>
        </template>
        <Table
          :columns="nodeUsageColumns"
          :data="nodeUsageData"
          :highlightCurrentRow="false"
          stripe
          :border="false"
          :height="600"
        />
      </ElCard>
    </ElCol>
    <ElCol :span="12">
      <ElCard shadow="hover" class="mb-25px">
        <template #header>
          <ElRow>
            <ElCol :span="12">
              <div>
                <span>{{ t('common.version') }}</span>
                <ElText
                  v-if="updateFlag"
                  type="danger"
                  size="small"
                  style="position: relative; left: 1rem"
                  >*{{ t('common.updatemsg') }}</ElText
                >
              </div>
            </ElCol>
            <!-- <ElCol :span="3" :offset="8" v-if="updateFlag">
              <ElPopconfirm title="Are you sure?" @confirm="updateSystem">
                <template #reference>
                  <ElButton color="#626aef">
                    {{ t('common.update') }}
                  </ElButton>
                </template>
              </ElPopconfirm>
            </ElCol> -->
          </ElRow>
        </template>
        <Table :columns="versionColumns" :data="versionData" stripe :border="false" :height="600" />
      </ElCard>
    </ElCol>
  </ElRow>
  <Dialog
    v-model="UpdatedialogVisible"
    title="Update"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="430"
  >
    <ElText type="danger" size="small" style="position: relative; left: 1rem"
      >*更新目前只支持docker容器搭建的程序，输入的url地址确保docker内可访问，节点最新版在github中releases的linux版本</ElText
    >
    <ElForm :model="form" label-width="120px" class="upload-form">
      <ElFormItem label="server url">
        <ElInput v-model="form.server" placeholder="server url" />
      </ElFormItem>
      <ElFormItem label="scan url">
        <ElInput
          v-model="form.scan"
          placeholder="scan url(https://github.com/Autumn-27/ScopeSentry-Scan/releases/download/vx.x.x/ScopeSentry-Scan_linux_amd64_vx.x.x.zip)"
        />
      </ElFormItem>
      <ElFormItem>
        <ElButton type="primary" @click="handleSubmit">Submit</ElButton>
      </ElFormItem>
    </ElForm>
  </Dialog>
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
      <BaseButton @click="savePluginKey" type="primary" class="w-full">确定</BaseButton>
    </div>
  </Dialog>
</template>

<style scoped>
.header-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
.tooltip-content {
  white-space: pre-line !important;
}
</style>

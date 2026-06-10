<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, onMounted } from 'vue'
import {
  ElButton,
  ElCol,
  ElInput,
  ElRow,
  ElText,
  ElTabs,
  ElTabPane,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  CheckboxValueType,
  ElSelectV2,
  ElCheckbox,
  ElTooltip,
  ElSwitch
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { useIcon } from '@/hooks/web/useIcon'
import {
  getScheduledTaskDataApi,
  scheduledDeleteTaskApi,
  updateScheduledTaskPageMonitApi
} from '@/api/task'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import AddTask from './components/AddTask.vue'
import AddProject from '../Project/components/AddProject.vue'
import { Icon } from '@iconify/vue'
import PageMonit from './components/PageMonit.vue'
import { getNodeDataOnlineApi } from '@/api/node'

const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const search = ref('')
const handleSearch = () => {
  console.log('as')
  getList()
}
const taskColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'name',
    label: t('task.taskName'),
    minWidth: 30
  },
  {
    field: 'cycle',
    label: t('task.taskCycle'),
    minWidth: 20
  },
  {
    field: 'type',
    label: t('task.typeTask'),
    minWidth: 20
  },
  {
    field: 'lastTime',
    label: t('task.lastTime'),
    minWidth: 40,
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
    minWidth: 40,
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
    field: 'state',
    label: t('common.state'),
    minWidth: 20,
    formatter: (_: Recordable, __: TableColumn, stateValue: boolean) => {
      if (stateValue == null) {
        return <div></div>
      }
      let color = ''
      let flag = ''
      if (stateValue == true) {
        color = '#2eb98a'
        flag = t('common.on')
      } else {
        color = 'red'
        flag = t('common.statusStop')
      }
      return (
        <ElRow gutter={20}>
          <ElCol span={1}>
            <Icon icon="clarity:circle-solid" color={color} />
          </ElCol>
          <ElCol span={5}>
            <ElText type="info">{flag}</ElText>
          </ElCol>
        </ElRow>
      )
    }
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 40,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          {row.TaskID === 'page_monitoring' ? (
            <BaseButton type="success" onClick={() => getPageMonitContent(row)}>
              {t('common.edit')}
            </BaseButton>
          ) : (
            <>
              <BaseButton type="success" onClick={() => getTaskContent(row)}>
                {t('common.edit')}
              </BaseButton>
              <BaseButton type="danger" onClick={() => confirmDelete(row)}>
                {t('common.delete')}
              </BaseButton>
            </>
          )}
          {/* <BaseButton type="warning" onClick={() => taskRunNow(row.id)}>
            {t('task.runNow')}
          </BaseButton> */}
        </>
      )
    }
  }
])

// const taskRunNow = async (id) => {
//   await taskRunApi(id)
// }

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const res = await getScheduledTaskDataApi(search.value, currentPage.value, pageSize.value)
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
let DialogTitle = t('task.addTask')
const closeDialog = () => {
  dialogVisible.value = false
}
let ProjectId = ''
let TaskId = ref('')
let Create = ref(true)
const getTaskContent = async (data) => {
  TaskId.value = data.id
  DialogTitle = t('common.edit')
  dialogVisible.value = true
}

const confirmDeleteSelect = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}

const confirmDelete = async (data) => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await del(data)
  }
}
const delLoading = ref(false)
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await scheduledDeleteTaskApi([data.id])
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const ids = ref<string[]>([])
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await scheduledDeleteTaskApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
onMounted(() => {
  setMaxHeight()
  window.addEventListener('resize', setMaxHeight)
})
const maxHeight = ref(0)

const setMaxHeight = () => {
  const screenHeight = window.innerHeight || document.documentElement.clientHeight
  maxHeight.value = screenHeight * 0.75
}
const projectDialogVisible = ref(false)
const closeProjectDialog = () => {
  projectDialogVisible.value = false
}
const pageMontDialogVisible = ref(false)
const ConfigPageMonitSaveLoading = ref(false)
const pageMontForm = reactive({
  hour: 24,
  allNode: true,
  node: [] as string[],
  scheduledTasks: true
})
const submitConfigPageMonitForm = async () => {
  ConfigPageMonitSaveLoading.value = true
  try {
    await updateScheduledTaskPageMonitApi(
      pageMontForm.hour,
      pageMontForm.node,
      pageMontForm.allNode,
      pageMontForm.scheduledTasks
    )
    getList()
  } finally {
    ConfigPageMonitSaveLoading.value = false
  }
}
const getPageMonitContent = async (data) => {
  pageMontForm.hour = data.hour
  pageMontForm.allNode = data.allNode
  pageMontForm.node = data.node
  pageMontForm.scheduledTasks = data.scheduledTasks
  pageMontDialogVisible.value = true
}
const nodeOptions = reactive<{ value: string; label: string }[]>([])
const indeterminate = ref(false)
const isCheckboxDisabledNode = ref(false)
const handleCheckAll = (val: CheckboxValueType) => {
  indeterminate.value = false
  if (val) {
    pageMontForm.allNode = true
    pageMontForm.node = []
    nodeOptions.forEach((option) => {
      return pageMontForm.node.push(option.value)
    })
  } else {
    pageMontForm.allNode = false
    pageMontForm.node = []
  }
}
const getNodeList = async () => {
  const res = await getNodeDataOnlineApi()
  if (res.data.list.length > 0) {
    isCheckboxDisabledNode.value = false
    res.data.list.forEach((item) => {
      nodeOptions.push({ value: item, label: item })
    })
  } else {
    isCheckboxDisabledNode.value = true
    ElMessage.warning(t('node.onlineNodeMsg'))
  }
}
getNodeList()
const addTask = async () => {
  DialogTitle = t('task.addScheduled')
  Create.value = true
  dialogVisible.value = true
}
</script>

<template>
  <ContentWrap>
    <ElRow>
      <ElCol :span="1">
        <ElText class="mx-1" style="position: relative; top: 8px">{{ t('task.taskName') }}:</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="search" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5" style="position: relative; left: 16px">
        <ElButton type="primary" :icon="searchicon" style="height: 100%" @click="handleSearch"
          >Search</ElButton
        >
      </ElCol>
    </ElRow>
    <ElRow>
      <ElCol style="position: relative; top: 16px">
        <div class="mb-10px">
          <BaseButton type="primary" @click="addTask">{{ t('task.addScheduled') }}</BaseButton>
          <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">
            {{ t('task.delTask') }}
          </BaseButton>
        </div>
      </ElCol>
    </ElRow>
    <div style="position: relative; top: 12px">
      <Table
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
        v-model:pageSize="pageSize"
        v-model:currentPage="currentPage"
        :columns="taskColums"
        :data="dataList"
        stripe
        :border="true"
        :loading="loading"
        :max-height="maxHeight"
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
    <Dialog
      v-model="dialogVisible"
      :title="DialogTitle"
      center
      style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    >
      <AddTask
        :closeDialog="closeDialog"
        :getList="getList"
        :create="Create"
        :taskid="TaskId"
        :schedule="true"
        tp="scan"
        :target-ids="[]"
      />
    </Dialog>
    <Dialog
      v-model="projectDialogVisible"
      :title="t('common.edit')"
      center
      style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    >
      <AddProject
        :closeDialog="closeProjectDialog"
        :projectid="ProjectId"
        :getProjectData="getList"
        :schedule="false"
      />
    </Dialog>
    <Dialog
      v-model="pageMontDialogVisible"
      :title="t('common.edit')"
      center
      style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    >
      <ElTabs type="card">
        <ElTabPane :label="t('router.configuration')">
          <ElForm :model="pageMontForm" label-width="100px" status-icon ref="ruleFormRef">
            <ElTooltip :content="t('task.selectNodeMsg')" placement="top">
              <ElFormItem :label="t('task.nodeSelect')" prop="node">
                <ElSelectV2
                  v-model="pageMontForm.node"
                  filterable
                  :options="nodeOptions"
                  placeholder="Please select node"
                  style="width: 80%"
                  multiple
                  tag-type="success"
                  collapse-tags
                  collapse-tags-tooltip
                  :max-collapse-tags="7"
                >
                  <template #header>
                    <ElCheckbox
                      v-model="pageMontForm.allNode"
                      :disabled="isCheckboxDisabledNode"
                      :indeterminate="indeterminate"
                      @change="handleCheckAll"
                    >
                      All
                    </ElCheckbox>
                  </template>
                </ElSelectV2>
              </ElFormItem>
            </ElTooltip>
            <ElFormItem :label="t('project.cycle')" prop="type">
              <ElInputNumber
                v-model="pageMontForm.hour"
                :min="1"
                controls-position="right"
                size="small"
              /><ElText style="position: relative; left: 16px">Hour</ElText>
            </ElFormItem>
            <ElFormItem :label="t('common.state')">
              <ElSwitch
                v-model="pageMontForm.scheduledTasks"
                inline-prompt
                :active-text="t('common.switchAction')"
                :inactive-text="t('common.switchInactive')"
              />
            </ElFormItem>
            <ElRow>
              <ElCol :span="2" :offset="8">
                <ElFormItem>
                  <ElButton
                    type="primary"
                    @click="submitConfigPageMonitForm()"
                    :loading="ConfigPageMonitSaveLoading"
                    >{{ t('task.save') }}</ElButton
                  >
                </ElFormItem>
              </ElCol>
            </ElRow>
          </ElForm>
        </ElTabPane>
        <ElTabPane :label="t('task.data')"><PageMonit /></ElTabPane>
      </ElTabs>
    </Dialog>
  </ContentWrap>
</template>

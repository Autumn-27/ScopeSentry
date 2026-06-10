<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, h, onMounted, resolveComponent } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import {
  ElButton,
  ElCol,
  ElInput,
  ElRow,
  ElText,
  ElProgress,
  ElTag,
  ElMessageBox,
  ElSwitch,
  ElDropdown,
  ElDropdownMenu,
  ElDropdownItem,
  ElIcon,
  ElRadioGroup,
  ElRadioButton,
  ElSelect,
  ElOption,
  ElTreeSelect
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { useIcon } from '@/hooks/web/useIcon'
import {
  getTaskDataApi,
  deleteTaskApi,
  retestTaskApi,
  stopTaskApi,
  starTaskApi,
  syancProjectApi
} from '@/api/task'
import { Dialog } from '@/components/Dialog'
import { BaseButton } from '@/components/Button'
import AddTask from './components/AddTask.vue'
import ProgressInfo from './components/ProgressInfo.vue'
import { useRouter } from 'vue-router'
import { getProjectAllApi } from '@/api/project'
const { push } = useRouter()
const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const search = ref('')
const handleSearch = () => {
  getList()
}
const taskColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    minWidth: 55
  },
  {
    field: 'name',
    label: t('task.taskName'),
    minWidth: 100
  },
  {
    field: 'taskNum',
    label: t('task.taskCount'),
    minWidth: 70,
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
    field: 'progress',
    label: t('task.taskProgress'),
    minWidth: 200,
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
    field: 'status',
    label: t('common.state'),
    minWidth: 200,
    formatter: (_: Recordable, __: TableColumn, cellValue: number) => {
      // 1为运行中，2为暂停，3为运行完成
      let tagType, tagText
      switch (cellValue) {
        case 1:
          tagType = 'info'
          tagText = t('task.running') // 运行中
          break
        case 2:
          tagType = 'warning'
          tagText = t('task.stop') // 暂停
          break
        case 3:
          tagType = 'success'
          tagText = t('task.finish') // 运行完成
          break
        default:
          tagType = 'default'
          tagText = '' // 未知状态
      }
      return h(ElTag, { type: tagType }, () => tagText)
    }
  },
  {
    field: 'creatTime',
    minWidth: 200,
    label: t('task.createTime')
  },
  {
    field: 'endTime',
    label: t('task.endTime'),
    minWidth: 200,
    formatter: (_: Recordable, __: TableColumn, cellValue: string) => {
      if (cellValue == '') {
        return '-'
      }
      return cellValue
    }
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: '420',
    fixed: 'right',
    formatter: (row, __: TableColumn, _: number) => {
      const handleCommand = (command) => {
        ids.value = []
        switch (command) {
          case 'retest':
            confirmRetest(row)
            break
          case 'delete':
            confirmDelete(row)
            break
          case 'stop':
            ids.value.push(row.id)
            stopTask(ids.value)
            break
          case 'start':
            ids.value.push(row.id)
            startTask(ids.value)
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
              // 根据 row.status 渲染不同的菜单项
              if (row.status === 3) {
                return [
                  h(ElDropdownItem, { command: 'retest' }, () => t('task.retest')),
                  h(ElDropdownItem, { command: 'delete' }, () => t('common.delete'))
                ]
              } else {
                return [
                  h(ElDropdownItem, { command: 'start' }, () => t('task.start')),
                  h(ElDropdownItem, { command: 'stop' }, () => t('task.stop')), // 如果是运行中，显示“停止”按钮
                  h(ElDropdownItem, { command: 'retest' }, () => t('task.retest')),
                  h(ElDropdownItem, { command: 'delete' }, () => t('common.delete'))
                ]
              }
            })
        }
      )
      return (
        <>
          {retestAndDeleteDropdown}
          <BaseButton
            type="primary"
            onClick={() => getTaskResult(row.name)}
            style={{ marginLeft: '10px' }}
          >
            {t('task.result')}
          </BaseButton>
          <BaseButton
            type="success"
            onClick={() => getTaskContent(row)}
            style={{ marginLeft: '10px' }}
          >
            {t('common.view')}
          </BaseButton>
          <ElButton type="warning" onClick={() => getProgressInfo(row.id)}>
            {t('task.taskProgress')}
          </ElButton>
        </>
      )
    }
  }
])

const progressDialogVisible = ref(false)
let getProgressInfoID = ''
const getProgressInfo = async (id) => {
  getProgressInfoID = id
  progressDialogVisible.value = true
}

const getTaskResult = async (id) => {
  push(`/asset-information/index?task=${id}`)
}

const stopTask = async (ids) => {
  console.log('begin stop')
  await stopTaskApi(ids)
}

const startTask = async (ids) => {
  console.log('begin start')
  await starTaskApi(ids)
}

const progresscloseDialog = () => {
  progressDialogVisible.value = false
}
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const res = await getTaskDataApi(search.value, currentPage.value, pageSize.value)
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
const addTask = async () => {
  taskid.value = ''
  DialogTitle = t('task.addTask')
  Create.value = true
  dialogVisible.value = true
}

let DialogTitle = t('task.addTask')
const closeDialog = () => {
  dialogVisible.value = false
}

let Create = ref(true)
const taskid = ref('')
const getTaskContent = async (data) => {
  taskid.value = data.id
  dialogVisible.value = true
  Create.value = false
  DialogTitle = t('common.view')
}
const confirmDeleteSelect = async () => {
  const deleteAssetS = ref<boolean | string | number>(false)
  ElMessageBox({
    title: 'Delete',
    draggable: true,
    // Should pass a function if VNode contains dynamic props
    message: () =>
      h('div', { style: { display: 'flex', alignItems: 'center' } }, [
        h('p', { style: { margin: '0 10px 0 0' } }, t('task.delAsset')),
        h(ElSwitch, {
          modelValue: deleteAssetS.value,
          'onUpdate:modelValue': (val: boolean | string | number) => {
            deleteAssetS.value = val
          }
        })
      ])
  }).then(async () => {
    await delSelect(deleteAssetS.value)
  })
}

const confirmStopSelect = async () => {
  ElMessageBox({
    title: 'Stop Task',
    draggable: true
  }).then(async () => {
    await stopTaskSelect()
  })
}
const stopTaskSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    await stopTask(ids.value)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error Stop data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmStartSelect = async () => {
  ElMessageBox({
    title: 'Start Task',
    draggable: true
  }).then(async () => {
    await startTaskSelect()
  })
}
const startTaskSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    await startTask(ids.value)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error Stop data:', error)
    delLoading.value = false
    getList()
  }
}
interface Project {
  value: string
  label: string
  children?: Project[]
}
const projectList = reactive<Project[]>([])
const getProjectList = async () => {
  const res = await getProjectAllApi()
  res.data.list.forEach((item: Project) => {
    projectList.push({
      label: item.label,
      value: item.value || `parent-${item.label}`, // 避免空字符串
      children: item.children || []
    })
  })
}

const confirmSyncToProjectSelect = async () => {
  const option = ref<'existing' | 'new'>('existing') // 选项：已有 or 新建
  const selectedProjectId = ref<string>('')
  const newProjectName = ref('')
  const newProjectTag = ref('')
  await getProjectList()
  ElMessageBox({
    title: t('task.syncToProject'),
    draggable: true,
    message: () =>
      h('div', { style: { display: 'flex', flexDirection: 'column', gap: '10px' } }, [
        // 选择类型
        h('div', [
          h('label', { style: { marginRight: '8px' } }),
          h(
            ElRadioGroup,
            {
              modelValue: option.value,
              'onUpdate:modelValue': (val: 'existing' | 'new') => {
                option.value = val
              }
            },
            {
              default: () => [
                h(ElRadioButton, { label: 'existing' }, () => t('task.syncToExisting')),
                h(ElRadioButton, { label: 'new' }, () => t('task.createNewProject'))
              ]
            }
          )
        ]),

        // 如果是同步到已有项目，显示下拉框
        option.value === 'existing'
          ? h(ElTreeSelect, {
              modelValue: selectedProjectId.value,
              'onUpdate:modelValue': (val: string) => {
                selectedProjectId.value = val
              },
              data: projectList,
              showCheckbox: true,
              placeholder: t('project.project'),
              filterable: true,
              style: { width: '100%' }
            })
          : null,

        // 如果是创建新项目，显示输入框
        option.value === 'new'
          ? h('div', { style: { display: 'flex', flexDirection: 'column', gap: '8px' } }, [
              h(ElInput, {
                modelValue: newProjectName.value,
                placeholder: t('project.msgProject'),
                'onUpdate:modelValue': (val: string) => (newProjectName.value = val)
              }),
              h(ElInput, {
                modelValue: newProjectTag.value,
                placeholder: t('project.msgProjectTag'),
                'onUpdate:modelValue': (val: string) => (newProjectTag.value = val)
              })
            ])
          : null
      ])
  }).then(async () => {
    if (option.value === 'existing') {
      console.log('同步到已有项目ID:', selectedProjectId.value)
    } else {
      console.log('创建新项目:', newProjectName.value, newProjectTag.value)
    }
    const elTableExpose = await getElTableExpose()
    const selectedRows = elTableExpose?.getSelectionRows() || []
    ids.value = selectedRows.map((row) => row.id)
    await syancProjectApi(
      ids.value,
      option.value,
      selectedProjectId.value,
      newProjectTag.value,
      newProjectName.value
    )
  })
}

const confirmDelete = async (data) => {
  const deleteAsset = ref<boolean | string | number>(false)
  ElMessageBox({
    title: 'Delete',
    draggable: true,
    // Should pass a function if VNode contains dynamic props
    message: () =>
      h('div', { style: { display: 'flex', alignItems: 'center' } }, [
        h('p', { style: { margin: '0 10px 0 0' } }, t('task.delAsset')),
        h(ElSwitch, {
          modelValue: deleteAsset.value,
          'onUpdate:modelValue': (val: boolean | string | number) => {
            deleteAsset.value = val
          }
        })
      ])
  }).then(async () => {
    await del(data, deleteAsset.value)
  })
}
const delLoading = ref(false)
const del = async (data, delA) => {
  delLoading.value = true
  try {
    const res = await deleteTaskApi([data.id], delA)
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
const delSelect = async (delA) => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await deleteTaskApi(ids.value, delA)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmRetest = async (data) => {
  const confirmed = window.confirm('Are you sure you want to retest?')
  if (confirmed) {
    await retestTask(data)
  }
}
const retestTask = async (data) => {
  try {
    await retestTaskApi(data.id)
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
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
          <BaseButton type="primary" @click="addTask">{{ t('task.addTask') }}</BaseButton>
          <BaseButton type="danger" :loading="delLoading" @click="confirmDeleteSelect">
            {{ t('task.delTask') }}
          </BaseButton>
          <BaseButton type="warning" :loading="delLoading" @click="confirmStopSelect">
            {{ t('task.stop') }}
          </BaseButton>
          <BaseButton type="success" :loading="delLoading" @click="confirmStartSelect">
            {{ t('task.start') }}
          </BaseButton>
          <BaseButton type="info" :loading="delLoading" @click="confirmSyncToProjectSelect">
            {{ t('task.syncToProject') }}
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
  </ContentWrap>
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
      :taskid="taskid"
      :schedule="false"
      tp="scan"
      :targetIds="[]"
    />
  </Dialog>
  <Dialog
    v-model="progressDialogVisible"
    :title="t('task.taskProgress')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    width="70%"
    max-height="700"
  >
    <ProgressInfo
      :closeDialog="progresscloseDialog"
      :getProgressInfoID="getProgressInfoID"
      getProgressInfotype="scan"
      getProgressInforunnerid=""
  /></Dialog>
</template>

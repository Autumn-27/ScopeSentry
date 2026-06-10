<script setup lang="ts">
import {
  ElCheckbox,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElRow,
  ElCol,
  ElSwitch,
  FormRules,
  ElTooltip,
  ElRadioGroup,
  ElRadio,
  ElSelectV2,
  ElOption,
  ElButton,
  FormInstance,
  ElSelect,
  ElMessage,
  ElInputNumber,
  CheckboxValueType,
  ElText,
  ElTreeSelect,
  ElDrawer
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { onMounted, reactive, ref, toRefs, watch } from 'vue'
import { getNodeDataOnlineApi } from '@/api/node'
import {
  addScheduledTaskApi,
  addTaskApi,
  getScheduleDetailApi,
  getTaskDetailApi,
  getTemplateDataApi,
  updateScheduleApi
} from '@/api/task'
import DetailTemplate from './DetailTemplate.vue'
import { getProjectAllApi } from '@/api/project'
const { t } = useI18n()

const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  create: boolean
  schedule: boolean
  taskid: string
  tp: string
  targetIds: string[]
  getFilter?: () => { [key: string]: any }
  searchParams?: string
}>()

interface RuleForm {
  name: string
  target: string
  node: []
  template: string
  day: number
  hour: number
  minute: number
}
const sourceTp = props.tp.includes('Source') ? true : false

const rules = reactive<FormRules<RuleForm>>({
  name: [{ required: true, message: t('task.msgTaskName'), trigger: 'blur' }],
  target: [
    {
      required: sourceTp ? false : true,
      message: t('task.msgTarget'),
      trigger: 'blur'
    }
  ],
  node: [{ required: true, message: t('task.nodeMsg'), trigger: 'blur' }],
  template: [{ required: true, message: 'Please select template', trigger: 'blur' }],
  day: [
    {
      message: '1-31',
      trigger: 'change',
      validator: (_, value, callback) => {
        console.log(value)
        if (!value) {
          callback(new Error('1-31'))
        } else if (!/^\d+$/.test(value)) {
          callback(new Error('1-31'))
        } else if (value < 1 || value > 31) {
          callback(new Error('1-31'))
        } else {
          callback() // 验证通过
        }
      }
    }
  ],
  hour: [
    {
      message: '0-24',
      trigger: 'change',
      validator: (rule, value, callback) => {
        console.log(value)
        if (!value) {
          callback(new Error('0-24'))
        } else if (!/^\d+$/.test(value)) {
          callback(new Error('0-24'))
        } else if (value < 0 || value > 24) {
          callback(new Error('0-24'))
        } else {
          callback() // 验证通过
        }
      }
    }
  ],
  minute: [
    {
      message: '0-60',
      trigger: 'change',
      validator: (rule, value, callback) => {
        console.log(value)
        if (!value) {
          callback(new Error('0-60'))
        } else if (!/^\d+$/.test(value)) {
          callback(new Error('0-60'))
        } else if (value < 0 || value > 60) {
          callback(new Error('0-60'))
        } else {
          callback() // 验证通过
        }
      }
    }
  ]
})

const saveLoading = ref(false)
const ruleFormRef = ref<FormInstance>()
const submitForm = async (formEl: FormInstance | undefined) => {
  saveLoading.value = true
  try {
    if (!formEl) return

    const valid = await formEl.validate() // 使用 Promise 风格的 validate
    if (valid) {
      let res
      let searchFilter = reactive<{ [key: string]: any }>({})
      if (targetTp.value == 'search') {
        if (props.getFilter) {
          searchFilter = props.getFilter()
        }
        if (props.searchParams) {
          taskData.search = props.searchParams
        }
      }
      if (sourceTp) {
        taskData.targetSource = props.tp
      }
      if (props.taskid) {
        // 修改计划任务
        res = await updateScheduleApi(
          props.taskid,
          taskData.name,
          taskData.target,
          taskData.ignore,
          taskData.node,
          taskData.allNode,
          taskData.duplicates,
          taskData.scheduledTasks,
          taskData.hour,
          taskData.template,
          targetTp.value,
          taskData.search,
          searchFilter,
          targetNumber.value,
          props.targetIds,
          taskData.project,
          taskData.targetSource,
          taskData.day,
          taskData.minute,
          taskData.week,
          taskData.bindProject,
          taskData.cycleType
        )
      } else {
        // 创建新任务
        if (props.schedule) {
          res = await addScheduledTaskApi(
            taskData.name,
            taskData.target,
            taskData.ignore,
            taskData.node,
            taskData.allNode,
            taskData.duplicates,
            taskData.scheduledTasks,
            taskData.hour,
            taskData.template,
            targetTp.value,
            taskData.search,
            searchFilter,
            targetNumber.value,
            props.targetIds,
            taskData.project,
            taskData.targetSource,
            taskData.day,
            taskData.minute,
            taskData.week,
            taskData.bindProject,
            taskData.cycleType
          )
        } else {
          res = await addTaskApi(
            taskData.name,
            taskData.target,
            taskData.ignore,
            taskData.node,
            taskData.allNode,
            taskData.duplicates,
            taskData.scheduledTasks,
            taskData.hour,
            taskData.template,
            targetTp.value,
            taskData.search,
            searchFilter,
            targetNumber.value,
            props.targetIds,
            taskData.project,
            taskData.targetSource,
            taskData.day,
            taskData.minute,
            taskData.week,
            taskData.bindProject,
            taskData.cycleType
          )
        }
      }
      if (res.code === 200) {
        props.closeDialog()
        props.getList()
      }
    }
  } catch (error) {
    console.error('提交表单时发生错误:', error)
  } finally {
    saveLoading.value = false // 确保无论成功或失败都重置加载状态
  }
}
const nodeOptions = reactive<{ value: string; label: string }[]>([])
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

const templateOptions = reactive<{ value: string; label: string }[]>([])

const getTemplateList = async () => {
  templateOptions.length = 0
  const res = await getTemplateDataApi('', 1, 1000)
  if (res.data.list.length > 0) {
    res.data.list.forEach((item) => {
      templateOptions.push({ value: item.id, label: item.name })
    })
  }
}

onMounted(() => {
  getNodeList()
  getTemplateList()
})
const indeterminate = ref(false)
const isCheckboxDisabledNode = ref(false)
const taskData = reactive({
  name: '',
  target: '',
  ignore: '',
  node: [] as string[],
  allNode: true,
  scheduledTasks: false,
  duplicates: 'None',
  template: '',
  cycleType: 'daily',
  search: '',
  project: [] as string[],
  targetSource: 'general',
  day: 1,
  hour: 1,
  minute: 30,
  week: 1,
  bindProject: null
})
const handleCheckAll = (val: CheckboxValueType) => {
  indeterminate.value = false
  if (val) {
    taskData.node = []
    nodeOptions.forEach((option) => {
      return taskData.node.push(option.value)
    })
  } else {
    taskData.node = []
  }
}
const templateId = ref('')
const dialogVisible = ref(false)

let DialogTitle = t('task.addTemplate')
const editTemplate = async (data) => {
  templateId.value = data
  if (data != '') {
    DialogTitle = t('task.editTemplate')
  }
  dialogVisible.value = true
}
const closeTemplateDialog = () => {
  dialogVisible.value = false
}

const loadTaskData = async (id) => {
  const res = await getTaskDetailApi(id)
  taskData.name = res.data.name
  taskData.target = res.data.target
  taskData.ignore = res.data.ignore
  taskData.node = res.data.node
  taskData.allNode = res.data.allNode
  taskData.scheduledTasks = res.data.scheduledTasks
  taskData.hour = res.data.hour
  taskData.duplicates = res.data.duplicates
  taskData.template = res.data.template
  taskData.day = res.data.day
  taskData.minute = res.data.minute
  taskData.week = res.data.week
  taskData.cycleType = res.data.cycleType
}

const loadScheduleData = async (id) => {
  const res = await getScheduleDetailApi(id)
  taskData.name = res.data.name
  taskData.target = res.data.target
  taskData.ignore = res.data.ignore
  taskData.node = res.data.node
  taskData.allNode = res.data.allNode
  taskData.scheduledTasks = res.data.scheduledTasks
  taskData.hour = res.data.hour
  taskData.duplicates = res.data.duplicates
  taskData.template = res.data.template
  taskData.project = res.data.project
  taskData.targetSource = res.data.targetSource
  taskData.day = res.data.day
  taskData.minute = res.data.minute
  taskData.week = res.data.week
  taskData.cycleType = res.data.cycleType
}

watch(
  () => props.taskid, // 监听 props.taskid 的变化
  async (newId) => {
    if (newId) {
      // 如果传入了 ID，则加载已有数据
      if (props.schedule) {
        // 如果是计划任务则从计划任务中加载数据
        await loadScheduleData(newId)
      } else {
        // 从任务中加载数据
        await loadTaskData(newId)
      }
    } else {
      if (props.schedule) {
        taskData.scheduledTasks = true
      } else {
        taskData.scheduledTasks = false
      }
      taskData.name = ''
      taskData.target = ''
      taskData.ignore = ''
      taskData.node = []
      taskData.allNode = true
      taskData.duplicates = 'None'
      taskData.template = ''
    }
  },
  { immediate: true } // 确保组件挂载时立即触发
)
const targetTp = ref('select')
const targetNumber = ref(0)
const targetSourceOptions = [
  {
    label: t('task.general'),
    value: 'general'
  },
  {
    label: t('task.fromProject'),
    value: 'project'
  },
  {
    label: t('task.fromAsset'),
    value: 'asset'
  },
  {
    label: t('task.fromRootDomain'),
    value: 'RootDomain'
  },
  {
    label: t('task.fromSubdomain'),
    value: 'subdomain'
  }
]
interface Project {
  value: string
  label: string
  children?: Project[]
}
const projectList = reactive<Project[]>([])
const getProjectList = async () => {
  const res = await getProjectAllApi()
  res.data.list.forEach((item: Project) => {
    projectList.push(item)
  })
}
getProjectList()
</script>
<template>
  <ElForm
    :model="taskData"
    label-width="auto"
    :rules="rules"
    status-icon
    ref="ruleFormRef"
    :disabled="create ? false : true"
  >
    <ElFormItem :label="t('task.taskName')" prop="name">
      <ElInput v-model="taskData.name" :placeholder="t('task.msgTaskName')" />
    </ElFormItem>
    <ElFormItem :label="t('task.targetSource')" v-if="!sourceTp">
      <ElSelectV2
        style="width: 50%"
        v-model="taskData.targetSource"
        :options="targetSourceOptions"
      />
    </ElFormItem>
    <ElRow
      v-if="taskData.targetSource != 'project' && taskData.targetSource != 'general' && !sourceTp"
    >
      <ElCol :span="12">
        <ElFormItem :label="t('task.search')">
          <ElInput v-model="taskData.search" :placeholder="t('form.input')" v-if="!sourceTp" />
        </ElFormItem>
      </ElCol>
      <ElCol :span="12">
        <ElFormItem :label="t('task.targetProject')">
          <ElTreeSelect
            v-model="taskData.project"
            :data="projectList"
            :placeholder="t('project.project')"
            multiple
            filterable
            show-checkbox
            collapse-tags
            :max-collapse-tags="1"
          />
        </ElFormItem>
      </ElCol>
    </ElRow>
    <ElFormItem :label="t('task.targetProject')" v-if="taskData.targetSource == 'project'">
      <ElTreeSelect
        v-model="taskData.project"
        :data="projectList"
        :placeholder="t('project.project')"
        multiple
        filterable
        show-checkbox
        collapse-tags
        :max-collapse-tags="1"
      />
    </ElFormItem>
    <ElFormItem
      :label="t('task.taskTarget')"
      prop="target"
      v-if="taskData.targetSource == 'general'"
    >
      <ElInput
        v-model="taskData.target"
        :placeholder="t('task.msgTarget')"
        type="textarea"
        rows="10"
        v-if="!sourceTp"
      />
      <ElRadioGroup v-model="targetTp" v-if="sourceTp">
        <ElRadio value="select">{{ t('task.select') }}</ElRadio>
        <ElRadio value="search">{{ t('export.exportTypeSearch') }}</ElRadio>
      </ElRadioGroup>
    </ElFormItem>
    <ElFormItem :label="t('task.targetNumber')" v-if="targetTp == 'search' && sourceTp">
      <ElInput v-model.number="targetNumber" />
    </ElFormItem>
    <ElFormItem :label="t('task.ignore')" prop="ignore" v-if="taskData.targetSource != 'project'">
      <ElInput
        v-model="taskData.ignore"
        :placeholder="t('task.ignoreMsg')"
        type="textarea"
        rows="5"
      />
    </ElFormItem>
    <!-- <ElFormItem :label="t('task.bindProject')" v-if="taskData.targetSource != 'project'">
      <ElTreeSelect
        v-model="taskData.bindProject"
        :data="projectList"
        :placeholder="t('project.project')"
        filterable
        show-checkbox
      />
    </ElFormItem> -->
    <ElRow>
      <ElCol :span="12">
        <ElFormItem :label="t('task.nodeSelect')" prop="node">
          <ElSelectV2
            v-model="taskData.node"
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
                :disabled="isCheckboxDisabledNode"
                :indeterminate="indeterminate"
                @change="handleCheckAll"
              >
                All
              </ElCheckbox>
            </template>
          </ElSelectV2>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12">
        <ElFormItem :label="t('task.autoNode')">
          <ElTooltip effect="dark" :content="t('task.selectNodeMsg')" placement="top">
            <ElSwitch
              v-model="taskData.allNode"
              inline-prompt
              :active-text="t('common.switchAction')"
              :inactive-text="t('common.switchInactive')"
            />
          </ElTooltip>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElFormItem :label="t('project.scheduledTasks')">
      <ElTooltip effect="dark" :content="t('project.msgScheduledTasks')" placement="top">
        <ElSwitch
          v-model="taskData.scheduledTasks"
          inline-prompt
          :active-text="t('common.switchAction')"
          :inactive-text="t('common.switchInactive')"
        />
      </ElTooltip>
    </ElFormItem>
    <ElFormItem :label="t('project.cycle')" prop="type" v-if="taskData.scheduledTasks">
      <ElRow :gutter="10" style="width: 100%">
        <ElCol :span="5">
          <ElSelect v-model="taskData.cycleType" style="width: 100%">
            <ElOption :label="t('task.daily')" value="daily" />
            <ElOption :label="t('task.ndays')" value="ndays" />
            <ElOption :label="t('task.nhours')" value="nhours" />
            <ElOption :label="t('task.weekly')" value="weekly" />
            <ElOption :label="t('task.monthly')" value="monthly" />
          </ElSelect>
        </ElCol>
        <ElCol :span="5" v-if="taskData.cycleType == 'weekly'">
          <ElSelect v-model.number="taskData.week" style="width: 100%">
            <ElOption :label="t('task.monday')" :value="0" />
            <ElOption :label="t('task.tuesday')" :value="1" />
            <ElOption :label="t('task.wednesday')" :value="2" />
            <ElOption :label="t('task.thursday')" :value="3" />
            <ElOption :label="t('task.friday')" :value="4" />
            <ElOption :label="t('task.saturday')" :value="5" />
            <ElOption :label="t('task.sunday')" :value="6" />
          </ElSelect>
        </ElCol>
        <ElCol :span="5" v-if="taskData.cycleType === 'ndays' || taskData.cycleType == 'monthly'">
          <ElFormItem prop="day">
            <ElInput style="width: 100%" v-model.number="taskData.day">
              <template #append>{{ t('task.day') }}</template>
            </ElInput>
          </ElFormItem>
        </ElCol>
        <ElCol
          :span="5"
          v-if="
            taskData.cycleType === 'daily' ||
            taskData.cycleType === 'ndays' ||
            taskData.cycleType == 'nhours' ||
            taskData.cycleType == 'weekly' ||
            taskData.cycleType == 'monthly'
          "
        >
          <ElFormItem prop="hour">
            <ElInput style="width: 100%" v-model.number="taskData.hour">
              <template #append>{{ t('task.hour') }}</template>
            </ElInput>
          </ElFormItem>
        </ElCol>
        <ElCol
          :span="5"
          v-if="
            taskData.cycleType === 'daily' ||
            taskData.cycleType === 'ndays' ||
            taskData.cycleType == 'nhours' ||
            taskData.cycleType == 'weekly' ||
            taskData.cycleType == 'monthly'
          "
        >
          <ElFormItem prop="minute">
            <ElInput style="width: 100%" v-model.number="taskData.minute">
              <template #append>{{ t('task.minute') }}</template>
            </ElInput>
          </ElFormItem>
        </ElCol>
      </ElRow>
    </ElFormItem>
    <ElDivider content-position="center" style="width: 60%; left: 20%">{{
      t('task.duplicates')
    }}</ElDivider>
    <ElRow>
      <ElCol :span="24">
        <ElFormItem :label="t('task.duplicates')" prop="type">
          <ElRadioGroup v-model="taskData.duplicates">
            <ElRadio label="None" name="duplicates" :checked="true" value="None" />
            <ElTooltip effect="dark" :content="t('task.duplicatesMsg')" placement="top">
              <ElRadio :label="t('task.duplicatesSubdomain')" name="duplicates" value="subdomain" />
            </ElTooltip>
          </ElRadioGroup>
        </ElFormItem>
      </ElCol>
    </ElRow>
    <ElDivider content-position="center" style="width: 60%; left: 20%">{{
      t('router.scanTemplate')
    }}</ElDivider>
    <ElFormItem :label="t('router.scanTemplate')" prop="template">
      <!-- <ElSelectV2 v-model="taskData.template" placeholder="Please select node" style="width: 50%" /> -->
      <ElSelect v-model="taskData.template" placeholder="Please select template" style="width: 30%">
        <ElOption
          v-for="item in templateOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        >
          <ElRow>
            <ElCol :span="16">
              <ElTooltip :content="item.label" placement="top">
                <span
                  style="
                    float: left;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    width: 100%;
                  "
                >
                  {{ item.label }}
                </span>
              </ElTooltip>
            </ElCol>
            <ElCol :span="8">
              <ElButton
                type="primary"
                size="small"
                style="margin-left: 15px"
                @click.stop="editTemplate(item.value)"
              >
                {{ t('common.edit') }}
              </ElButton>
            </ElCol>
          </ElRow>
        </ElOption>
        <template #footer>
          <ElButton
            type="success"
            size="small"
            plain
            style="margin-left: 15px"
            @click.stop="editTemplate('')"
          >
            {{ t('common.new') }}
          </ElButton>
        </template>
      </ElSelect>
    </ElFormItem>
    <ElDivider />
    <ElRow>
      <ElCol :span="2" :offset="10">
        <ElFormItem>
          <ElButton type="primary" @click="submitForm(ruleFormRef)" :loading="saveLoading">
            {{ t('task.save') }}
          </ElButton>
        </ElFormItem>
      </ElCol>
    </ElRow>
  </ElForm>
  <ElDrawer v-model="dialogVisible" :title="DialogTitle" direction="rtl" size="80%">
    <DetailTemplate
      :closeDialog="closeTemplateDialog"
      :getList="getTemplateList"
      :id="templateId"
    />
  </ElDrawer>
</template>

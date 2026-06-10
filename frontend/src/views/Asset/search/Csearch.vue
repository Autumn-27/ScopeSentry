<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { nextTick, onMounted, reactive, ref, watch } from 'vue'
import {
  ElCol,
  ElRow,
  ElButton,
  ElTable,
  ElTableColumn,
  ElText,
  ElDivider,
  ElAutocomplete,
  ElIcon,
  ElDropdown,
  ElDropdownMenu,
  ElDropdownItem,
  ElMessageBox,
  ElMessage,
  ElTreeSelect,
  ElTag,
  ElSwitch,
  ElSelect,
  ElOption
} from 'element-plus'
import { Dialog } from '@/components/Dialog'
import { useIcon } from '@/hooks/web/useIcon'
import { Icon } from '@/components/Icon'
import exportData from '../export/exportData.vue'
import { delDataApi } from '@/api/asset'
import { useRoute } from 'vue-router'
import { defineProps, defineEmits } from 'vue'
import { CrudSchema } from '@/hooks/web/useCrudSchemas'
import AddTask from '../../Task/components/AddTask.vue'
const { t } = useI18n()
const { query } = useRoute()
const props = defineProps<{
  getList: () => void
  handleSearch: (string) => void
  searchKeywordsData: {
    keyword: string
    example: string
    explain: string
  }[]
  index: string
  getElTableExpose: () => void
  handleFilterSearch: (string, any) => void
  projectList: Project[]
  taskList: { id: string; name: string }[]
  dynamicTags?: string[]
  handleClose?: (string) => void
  openAggregation?: () => void
  crudSchemas: Array<CrudSchema>
  statisticsHidden?: boolean
  changeStatisticsHidden?: (boolean) => void
  searchResultCount: number
  sensitiveAllNumber?: number
  activeSegment?: 'tableSegment' | 'cardSegment' // 可选属性
  setActiveSegment?: (segment: 'tableSegment' | 'cardSegment', flag: boolean) => void // 可选方法
  getFilter: () => { [key: string]: any }
  iconData?: { value: string; number: number; icon_hash: string }[]
}>()
const localSearchKeywordsData = reactive([...props.searchKeywordsData])
const newKeyword = {
  keyword: 'task',
  example: 'task=="test"',
  explain: t('searchHelp.taskName')
}
const tagKeyword = {
  keyword: 'tag',
  example: 'tag=="test"',
  explain: 'find tags'
}
localSearchKeywordsData.push(newKeyword)
localSearchKeywordsData.push(tagKeyword)
const AssignmentHelp = [
  {
    operator: '=',
    meaning: t('searchHelp.like'),
    value: '=""'
  },
  {
    operator: '!=',
    meaning: t('searchHelp.notIn'),
    value: '!=""'
  },
  {
    operator: '==',
    meaning: t('searchHelp.equal'),
    value: '==""'
  }
]
const logicHelp = [
  {
    operator: '&&',
    meaning: t('searchHelp.and'),
    value: '&&',
    logic: '1'
  },
  {
    operator: '||',
    meaning: t('searchHelp.or'),
    value: '||',
    logic: '1'
  },
  {
    operator: '()',
    meaning: t('searchHelp.brackets'),
    value: '()',
    logic: '1'
  }
]
const searchHelpData = AssignmentHelp.concat(logicHelp)
const dialogVisible = ref(false)

// 保存列显示配置到localStorage
const saveColumnConfig = () => {
  const config = props.crudSchemas.reduce((acc, column) => {
    if (column.field != 'select') {
      acc[column.field] = column.hidden
    }
    return acc
  }, {})
  localStorage.setItem(`columnConfig_${props.index}`, JSON.stringify(config))
}

// 从localStorage加载配置
const loadColumnConfig = () => {
  const savedConfig = JSON.parse(localStorage.getItem(`columnConfig_${props.index}`) || '{}')
  props.crudSchemas.forEach((col) => {
    if (savedConfig[col.field] !== undefined) {
      col.hidden = savedConfig[col.field]
    }
  })
}

// 初始化加载配置
loadColumnConfig()
watch(
  () => props.crudSchemas,
  () => {
    saveColumnConfig()
  },
  { deep: true }
)
const getHelp = () => {
  dialogVisible.value = true
}
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
const searchParams = ref('')
const searchicon = useIcon({ icon: 'iconoir:search' })
const help = useIcon({ icon: 'tdesign:chat-bubble-help' })
const elDropdownicon = useIcon({ icon: 'ri:arrow-drop-down-line' })
const exporticon = useIcon({ icon: 'ph:export-light' })
const aggregationIcon = useIcon({ icon: 'carbon:data-vis-1' })
const deleteicon = useIcon({ icon: 'openmoji:delete' })
const TASKicon = useIcon({ icon: 'carbon:task-complete' })
const exportDialogVisible = ref(false)
const openExport = () => {
  exportDialogVisible.value = true
}
const ids = ref<string[]>([])

const getIds = async () => {
  const elTableExpose = await props.getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  return ids.value
}
const delSelect = async () => {
  ElMessageBox.confirm('Whether to delete?', 'Warning', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    type: 'warning'
  })
    .then(async () => {
      await getIds()
      await delDataApi(ids.value, props.index)
      props.getList()
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Delete canceled'
      })
    })
}

const selectedKeyword = ref('')
const opSelect = ref(false)
const opSelect2 = ref(false)
let keyword = ref(true)
let assignmen = ref(false)
let logic = ref(false)
const querySearch = (queryString, cb) => {
  selectedKeyword.value = queryString
  console.log(queryString)
  if (queryString == '') {
    keyword.value = true
    logic.value = false
    assignmen.value = false
  }
  if (keyword.value) {
    console.log('第一')
    if (logic.value) {
      queryString = queryString.replace(selectedKeyword.value, '').trim()
    }
    const results = localSearchKeywordsData.filter((item) => {
      return item.keyword.toLowerCase().includes(queryString.toLowerCase())
    })
    cb(results)
    return
  }
  if (assignmen.value) {
    console.log('第2')
    const searchStr = queryString.replace(selectedKeyword.value, '').trim()
    const results = AssignmentHelp.filter((item) => item.operator.includes(searchStr))
    cb(results)
    return
  }
  if (logic.value && queryString.endsWith(' ')) {
    console.log('第3')
    const searchStr = queryString.replace(searchParams.value, '').trim()
    const results = logicHelp.filter((item) => item.operator.includes(searchStr))
    cb(results)
    return
  }
  cb([])
  return
}

const handleSelect = (item) => {
  console.log(item)
  if (item.keyword) {
    console.log('handleSelect 1')
    let key = ''
    if (logic.value) {
      key = selectedKeyword.value + item.keyword
    } else {
      key = item.keyword
    }
    selectedKeyword.value = key
    searchParams.value = key
    opSelect.value = true
    keyword.value = false
    assignmen.value = true
  } else if (item.logic) {
    console.log('handleSelect 2')
    searchParams.value = `${selectedKeyword.value}${item.value}`
    selectedKeyword.value = searchParams.value
    keyword.value = true
  } else {
    console.log('handleSelect 3')
    searchParams.value = `${selectedKeyword.value}${item.value}`
    selectedKeyword.value = searchParams.value
    opSelect2.value = true
    assignmen.value = false
    logic.value = true
  }
}
interface Project {
  value: string
  label: string
  children?: Project[]
}
const projectLoading = ref(false)
const projectValue = ref([])
const taskValue = ref([])
const filterChange = async () => {
  console.log(projectValue.value)
  console.log(taskValue.value)
  const filterData: { [key: string]: any } = { project: projectValue.value }
  filterData.task = taskValue.value
  props.handleFilterSearch(searchParams.value, filterData)
}
watch(
  () => projectValue.value,
  (newValue) => {
    filterChange()
  }
)

watch(
  () => taskValue.value,
  (newValue) => {
    filterChange()
  }
)
const localDynamicTags = ref<string[]>(props.dynamicTags ? [...props.dynamicTags] : [])
const tagClickFilterSearch = () => {
  const dictionary: { [key: string]: string[] } = {}
  console.log(localDynamicTags.value)
  console.log('Updated dictionary:', dictionary)
  localDynamicTags.value.forEach((tag) => {
    const [key, value] = tag.split('=')
    if (key in dictionary) {
      dictionary[key].push(value)
    } else {
      dictionary[key] = [value]
    }
  })
  dictionary['project'] = projectValue.value
  console.log('Updated dictionary:', dictionary)
  props.handleFilterSearch(searchParams.value, dictionary)
}
let task = query.task as string
if (task !== undefined && task !== '') {
  localDynamicTags.value.push(`task=${task}`)
}
const savedActiveSegmentConfig = JSON.parse(localStorage.getItem('assetActiveSegment') || '{}')

// 如果配置中有 activeSegment，则使用它，否则使用默认值
if (savedActiveSegmentConfig && savedActiveSegmentConfig.activeSegment) {
  if (props.setActiveSegment) {
    props.setActiveSegment(savedActiveSegmentConfig.activeSegment, false)
  }
}
tagClickFilterSearch()
watch(
  () => props.dynamicTags,
  (newTags) => {
    if (newTags) {
      localDynamicTags.value = [...newTags]
      if (task !== undefined && task !== '') {
        localDynamicTags.value.push(`task=${task}`)
      }
    } else {
      localDynamicTags.value = []
    }
    tagClickFilterSearch()
  },
  { immediate: false }
)
function handleCloseTag(tag: string) {
  if (tag.includes('task=')) {
    task = ''
  }
  if (props.handleClose) {
    props.handleClose(tag)
  } else {
    console.warn('handleClose function is not defined')
  }
}
function getIconByHash(hash: string) {
  return props.iconData?.find((item) => item.icon_hash === hash)
}
function clearAllTags() {
  localDynamicTags.value = []
  if (props.handleClose) {
    props.handleClose('close')
  }
}
const emit = defineEmits<{
  (event: 'update-column-visibility', payload: { field: string; hidden: boolean }): void
}>()
// 处理开关变化，通知父组件
const handleSwitchChange = (field) => {
  emit('update-column-visibility', { field: field.field, hidden: field.hidden })
}
const localStatisticsHidden = ref(props.statisticsHidden)
const refreshPage = () => {
  location.reload()
}
// const activeSegment = ref<'tableSegment' | 'cardSegment'>('tableSegment')

// const setActiveSegment = (segment: 'tableSegment' | 'cardSegment') => {
//   activeSegment.value = segment
// }

const tableSegmentIcon = useIcon({ icon: 'icons8:insert-table' })
const cardSegmentIcon = useIcon({ icon: 'flowbite:grid-solid' })
function handleSetActiveSegment(segment: 'tableSegment' | 'cardSegment') {
  if (props.setActiveSegment) {
    props.setActiveSegment(segment, true)
  }
}
const taskDialogVisible = ref(false)
let DialogTitle = t('task.addTask')
const taskCloseDialog = () => {
  taskDialogVisible.value = false
}
const openCreateTask = async () => {
  await getIds()
  taskDialogVisible.value = true
}
</script>

<template>
  <ContentWrap>
    <!-- <ElRow justify="start">
      <ElCol :span="1">
        <ElText>{{ t('form.input') }}:</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="searchParams" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5">
        <ElButton type="primary" :icon="searchicon" @click="props.handleSearch(searchParams)">
          Search
        </ElButton>
      </ElCol>
    </ElRow> -->
    <div class="search-toolbar">
      <div class="search-input-container">
        <ElAutocomplete
          v-model="searchParams"
          :fetch-suggestions="querySearch"
          :placeholder="t('form.input')"
          popperClass="my-autocomplete"
          @select="handleSelect"
          class="search-autocomplete"
          style="width: 100%"
        >
          <template #append>
            <ElButton @click="getHelp" text style="display: contents" :icon="help">
              <!-- <template #default>
                <ElIcon :icon="help" style="color: black" />
              </template> -->
            </ElButton>
          </template>
          <template #default="{ item }">
            <span style="float: left">{{ item.keyword || item.operator }}</span>
            <span style="float: right; color: var(--el-text-color-secondary); font-size: 13px">
              {{ item.explain || item.meaning }}
            </span>
          </template>
        </ElAutocomplete>
      </div>
      <ElButton
        type="primary"
        :icon="searchicon"
        @click="$props.handleSearch(searchParams)"
        class="toolbar-btn"
      >
        {{ t('form.input') }}
      </ElButton>
      <ElButton type="primary" @click="openExport" :icon="exporticon" class="toolbar-btn">
        {{ t('asset.export') }}
      </ElButton>
      <div class="task-select-container">
        <ElSelect
          v-model="taskValue"
          :placeholder="t('task.taskName')"
          multiple
          filterable
          collapse-tags
          :max-collapse-tags="1"
          class="task-select"
          clearable
        >
          <ElOption
            v-for="taskItem in $props.taskList"
            :key="taskItem.id"
            :label="taskItem.name"
            :value="taskItem.name"
          />
        </ElSelect>
      </div>
      <div class="project-select-container">
        <ElTreeSelect
          :loading="projectLoading"
          v-model="projectValue"
          :data="$props.projectList"
          :placeholder="t('project.project')"
          multiple
          filterable
          show-checkbox
          collapse-tags
          :max-collapse-tags="1"
          class="project-select"
        />
      </div>
      <ElDropdown trigger="click" class="toolbar-dropdown">
        <ElButton plain class="custom-button toolbar-btn">
          {{ t('common.operation') }}
          <ElIcon class="el-icon--right"><elDropdownicon /></ElIcon>
        </ElButton>
        <template #dropdown>
          <ElDropdownMenu>
            <ElDropdownItem :icon="deleteicon" @click="delSelect">{{
              t('common.delete')
            }}</ElDropdownItem>
            <ElDropdownItem :icon="TASKicon" @click="openCreateTask">{{
              t('task.addTask')
            }}</ElDropdownItem>
          </ElDropdownMenu>
        </template>
      </ElDropdown>
      <ElDropdown class="toolbar-dropdown">
        <div class="custom-dropdown">
          <Icon icon="ant-design:setting-outlined" class="cursor-pointer" />
        </div>
        <template #dropdown>
          <ElDropdownMenu>
            <ElDropdownItem v-for="(field, i) in crudSchemas" :key="i">
              <div class="dropdown-item" v-if="field.field != 'selection'">
                <span class="label-text">{{ field.label }}</span>
                <ElSwitch
                  size="small"
                  v-model="field.hidden"
                  :active-value="false"
                  :inactive-value="true"
                  @change="handleSwitchChange(field)"
                />
              </div>
            </ElDropdownItem>
            <ElDropdownItem v-if="$props.index == 'asset'">
              <span class="label-text">{{ t('asset.Chart') }}</span>
              <ElSwitch
                size="small"
                v-model="localStatisticsHidden"
                :active-value="false"
                :inactive-value="true"
                @change="changeStatisticsHidden(localStatisticsHidden)"
              />
            </ElDropdownItem>
            <ElDropdownItem divided>
              <ElButton style="width: 100%" type="primary" @click="refreshPage">Save</ElButton>
            </ElDropdownItem>
          </ElDropdownMenu>
        </template>
      </ElDropdown>
      <div class="segment-container" v-if="index == 'asset'">
        <div class="segment-control">
          <div
            class="segment"
            :class="{ active: props.activeSegment === 'tableSegment' }"
            @click="handleSetActiveSegment('tableSegment')"
          >
            <ElIcon>
              <tableSegmentIcon />
            </ElIcon>
          </div>
          <div
            class="segment"
            :class="{ active: props.activeSegment === 'cardSegment' }"
            @click="handleSetActiveSegment('cardSegment')"
          >
            <ElIcon>
              <cardSegmentIcon />
            </ElIcon>
          </div>
        </div>
      </div>
      <ElButton
        type="success"
        @click="$props.openAggregation"
        :icon="aggregationIcon"
        v-if="index == 'SensitiveResult'"
        class="toolbar-btn"
      >
        {{ t('project.aggregation') }}
      </ElButton>
    </div>
    <ElRow class="result-info-row">
      <ElCol :span="24">
        <div class="flex gap-2" style="flex-wrap: wrap">
          <span style="color: #888">{{ t('asset.total') }}</span>
          <span style="font-weight: bold; color: #333333">{{ props.searchResultCount }}</span>
          <span style="color: #888">{{ t('asset.result') }}</span>
          <!-- <div v-if="index == 'SensitiveResult'">
            <span style="color: #888">{{ t('asset.total') }}</span>
            <span style="font-weight: bold; color: #333333">{{ props?.sensitiveAllNumber }}</span>
            <span style="color: #888">{{ t('asset.sensitiveNumber') }}</span>
          </div> -->
          <ElTag
            v-for="tag in localDynamicTags"
            :key="tag"
            closable
            :disable-transitions="false"
            type="info"
            size="small"
            @close="handleCloseTag(tag)"
          >
            <template v-if="tag.startsWith('icon=')">
              <img
                :src="'/images/icon/' + tag.split('=')[1] + '.png'"
                :alt="tag"
                style="width: 20px; height: 20px; vertical-align: middle"
              />
            </template>
            <template v-else>
              {{ tag }}
            </template>
          </ElTag>
          <!-- <ElButton
            v-if="localDynamicTags.length > 0 && index == 'asset'"
            size="small"
            type="danger"
            plain
            @click="clearAllTags"
            style="margin-left: 8px; align-self: flex-start"
            >{{ t('common.clearAll') }}</ElButton
          > -->
        </div>
      </ElCol>
    </ElRow>
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="t('common.querysyntax')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <ElRow>
      <ElCol>
        <ElText tag="b" size="small">{{ t('searchHelp.operator') }}</ElText>
        <ElDivider direction="vertical" />
        <ElText size="small" type="danger">{{ t('searchHelp.notice') }}</ElText>
      </ElCol>
      <ElCol style="margin-top: 10px">
        <ElTable :headerCellStyle="tableHeaderColor" :data="searchHelpData">
          <ElTableColumn prop="operator" :label="t('searchHelp.operator')" width="300" />
          <ElTableColumn prop="meaning" :label="t('searchHelp.meaning')" />
        </ElTable>
      </ElCol>
      <ElCol style="margin-top: 15px">
        <ElText tag="b" size="small">{{ t('searchHelp.keywords') }}</ElText>
      </ElCol>
      <ElCol style="margin-top: 10px">
        <ElTable :headerCellStyle="tableHeaderColor" :data="localSearchKeywordsData">
          <ElTableColumn prop="keyword" :label="t('searchHelp.keywords')" />
          <ElTableColumn prop="example" :label="t('searchHelp.example')" />
          <ElTableColumn prop="explain" :label="t('searchHelp.explain')" />
        </ElTable>
      </ElCol>
    </ElRow>
  </Dialog>
  <Dialog
    v-model="exportDialogVisible"
    :title="t('asset.export')"
    center
    max-height="300"
    width="70%"
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <exportData :index="$props.index" :searchParams="searchParams" :getFilter="$props.getFilter" />
  </Dialog>
  <Dialog
    v-model="taskDialogVisible"
    :title="DialogTitle"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <AddTask
      :closeDialog="taskCloseDialog"
      :create="true"
      taskid=""
      :schedule="false"
      :getList="function () {}"
      :tp="$props.index + 'Source'"
      :target-ids="ids"
      :getFilter="$props.getFilter"
      :searchParams="searchParams"
    />
  </Dialog>
</template>
<style scoped>
/* 搜索工具栏 - 使用 flexbox 布局 */
.search-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  width: 100%;
  padding: 0;
  min-height: 40px;
  margin-bottom: 10px;
}

/* 搜索框容器 - 占据主要空间但有限制 */
.search-input-container {
  flex: 1 1 600px;
  min-width: 500px;
  max-width: 700px;
  height: 40px;
  display: flex;
  align-items: center;
  margin-right: 6px;
}

.search-autocomplete {
  width: 100% !important;
  min-width: 100% !important;
  max-width: 100% !important;
}

/* 确保 ElAutocomplete 输入框高度一致 */
.search-autocomplete :deep(.el-input__wrapper) {
  height: 40px;
  box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
  width: 100% !important;
}

/* 确保 ElAutocomplete 容器宽度 */
.search-autocomplete :deep(.el-autocomplete) {
  width: 100% !important;
  display: block;
}

.search-autocomplete :deep(.el-input) {
  width: 100% !important;
}

.search-autocomplete :deep(.el-input__inner) {
  height: 38px;
  line-height: 38px;
}

/* 工具栏按钮 */
.toolbar-btn {
  flex-shrink: 1;
  white-space: nowrap;
  height: 40px;
  padding: 0 20px;
  margin: 0;
  margin-right: 6px;
  min-width: fit-content;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 任务选择器 */
.task-select-container {
  flex-shrink: 1;
  width: 200px;
  min-width: 150px;
  max-width: 200px;
  height: 40px;
  display: flex;
  align-items: center;
  margin-right: 6px;
}

.task-select {
  width: 100%;
}

/* 确保 ElSelect 高度一致 */
.task-select :deep(.el-select__wrapper) {
  height: 40px;
  box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
}

.task-select :deep(.el-select__placeholder) {
  line-height: 38px;
}

/* 项目选择器 */
.project-select-container {
  flex-shrink: 1;
  width: 200px;
  min-width: 150px;
  max-width: 200px;
  height: 40px;
  display: flex;
  align-items: center;
  margin-right: 6px;
}

.project-select {
  width: 100%;
}

/* 确保 ElTreeSelect 高度一致 */
.project-select :deep(.el-select__wrapper) {
  height: 40px;
  box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
}

.project-select :deep(.el-select__placeholder) {
  line-height: 38px;
}

/* 下拉菜单 */
.toolbar-dropdown {
  flex-shrink: 1;
  height: 40px;
  display: flex;
  align-items: center;
  margin-right: 6px;
  min-width: fit-content;
}

/* 确保下拉菜单中的按钮高度一致 */
.toolbar-dropdown :deep(.el-button) {
  height: 40px;
  padding: 0 15px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.toolbar-dropdown .custom-button {
  height: 40px;
  padding: 0 15px;
}

.toolbar-dropdown .custom-dropdown {
  height: 40px;
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.3s;
  border: 1px solid var(--el-border-color);
}

.toolbar-dropdown .custom-dropdown:hover {
  background-color: var(--el-fill-color-light);
  border-color: var(--el-border-color-hover);
}

/* 分段控制器 */
.segment-container {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  height: 40px;
  margin-right: 6px;
  min-width: 80px;
  order: 999;
}

.segment-control {
  display: flex;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  height: 40px;
}

.segment {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 12px;
  cursor: pointer;
  transition: all 0.3s;
  background-color: #fff;
  min-width: 40px;
}

.segment:hover {
  background-color: #f5f7fa;
}

.segment.active {
  background-color: var(--el-color-primary);
  color: #fff;
}

.segment.active:hover {
  background-color: var(--el-color-primary);
}

/* 响应式设计 */
@media (max-width: 1600px) {
  .search-input-container {
    flex: 1 1 550px;
    min-width: 500px;
    max-width: 650px;
  }

  .task-select-container {
    width: 180px;
    min-width: 160px;
  }

  .project-select-container {
    width: 180px;
    min-width: 160px;
  }

  .segment-container {
    flex-shrink: 0;
    min-width: 80px;
  }
}

@media (max-width: 1500px) {
  .search-input-container {
    flex: 1 1 450px;
    min-width: 350px;
    max-width: 550px;
  }

  .toolbar-btn {
    padding: 0 15px;
    font-size: 14px;
  }

  .task-select-container {
    width: 160px;
    min-width: 140px;
  }

  .project-select-container {
    width: 160px;
    min-width: 140px;
  }

  .segment-container {
    flex-shrink: 0;
    min-width: 80px;
  }
}

@media (max-width: 1400px) {
  .search-input-container {
    flex: 1 1 380px;
    min-width: 300px;
    max-width: 450px;
  }

  .toolbar-btn {
    padding: 0 12px;
    font-size: 13px;
  }

  .task-select-container {
    width: 150px;
    min-width: 120px;
  }

  .project-select-container {
    width: 150px;
    min-width: 120px;
  }

  .segment-container {
    flex-shrink: 0;
    min-width: 80px;
  }
}

@media (max-width: 1450px) {
  .search-toolbar {
    flex-wrap: wrap;
    gap: 8px;
  }

  .search-input-container {
    flex: 1 1 100%;
    min-width: 100%;
    max-width: 100%;
    order: 1;
    margin-bottom: 8px;
  }

  .toolbar-btn,
  .task-select-container,
  .project-select-container,
  .toolbar-dropdown {
    order: 2;
  }

  .segment-container {
    order: 3;
    flex-shrink: 0;
    min-width: 80px;
  }
}

@media (max-width: 1300px) {
  .search-toolbar {
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 12px;
  }

  .search-input-container {
    flex: 1 1 100%;
    min-width: 100%;
    max-width: 100%;
    order: 1;
    margin-bottom: 8px;
  }

  .toolbar-btn,
  .task-select-container,
  .project-select-container,
  .toolbar-dropdown,
  .segment-container {
    order: 2;
    flex-shrink: 1;
  }

  .toolbar-btn {
    padding: 0 12px;
    font-size: 13px;
  }

  .segment-container {
    min-width: 80px;
    flex-shrink: 0;
  }
}

@media (max-width: 1200px) {
  .search-toolbar {
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
    min-height: auto;
    padding-bottom: 8px;
  }

  .search-input-container {
    flex: 1 1 100%;
    min-width: 100%;
    order: 1;
    margin-bottom: 8px;
  }

  .toolbar-btn,
  .task-select-container,
  .project-select-container,
  .toolbar-dropdown,
  .segment-container {
    order: 2;
  }
}

@media (max-width: 768px) {
  .search-toolbar {
    gap: 6px;
    margin-bottom: 16px;
    padding-bottom: 8px;
  }

  .toolbar-btn {
    padding: 8px 12px;
    font-size: 12px;
  }

  .task-select-container {
    width: 100%;
    min-width: 100%;
    order: 2;
    margin-bottom: 6px;
  }

  .project-select-container {
    width: 100%;
    min-width: 100%;
    order: 2;
    margin-bottom: 6px;
  }
}

/* 结果信息行 */
.result-info-row {
  margin-top: 10px;
  clear: both;
}

@media (max-width: 1200px) {
  .result-info-row {
    margin-top: 16px;
  }
}

/* 原有样式 */
.custom-button:hover {
  background-color: transparent !important;
  color: inherit !important;
  box-shadow: none !important;
  border-color: inherit !important;
  border-width: 1px !important;
}

.my-autocomplete .el-scrollbar__view {
  max-height: 300px;
  overflow-y: auto;
}

.my-autocomplete li {
  line-height: normal;
  padding: 7px;
}

.my-autocomplete li .name {
  text-overflow: ellipsis;
  overflow: hidden;
}

.my-autocomplete li .addr {
  font-size: 12px;
  color: #b4b4b4;
}

.my-autocomplete li .highlighted .addr {
  color: #ddd;
}

.custom-dropdown:focus-visible {
  outline: unset;
}

.dropdown-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.label-text {
  margin-right: 10px;
}
</style>

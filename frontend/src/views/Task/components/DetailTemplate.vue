<script setup lang="ts">
import { ref, reactive, watch, nextTick } from 'vue'
import {
  ElMessage,
  ElTooltip,
  ElTag,
  ElCard,
  ElInput,
  ElSwitch,
  ElButton,
  ElForm,
  ElFormItem,
  ElRow,
  ElCol,
  ElSelectV2,
  ElTreeV2,
  ElSelect,
  ElOption,
  ElDrawer
} from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import { Dialog } from '@/components/Dialog'
import { useI18n } from '@/hooks/web/useI18n'
import { getPluginDataByModuleApi } from '@/api/plugins'
import { getPocDataAllApi } from '@/api/poc'
import { getTemplateDetailApi, saveTemplateDetailApi } from '@/api/task'
import { pocData } from '@/api/poc/types'
import { getManagetListApi, getPortDictDataApi } from '@/api/DictionaryManagement'
import type { fileData, portDictData } from '@/api/DictionaryManagement/types'

const { t } = useI18n()

// 接收父组件传递的 props
const props = defineProps<{
  closeDialog: () => void
  getList: () => void
  id: string
}>()

// 模块数组
const modules = [
  'TargetHandler',
  'SubdomainScan',
  'SubdomainSecurity',
  'PortScanPreparation',
  'PortScan',
  'PortFingerprint',
  'AssetMapping',
  'AssetHandle',
  'URLScan',
  'WebCrawler',
  'URLSecurity',
  'DirScan',
  'VulnerabilityScan',
  'PassiveScan'
]

// 参数类型定义
interface ParameterItem {
  name: string
  type: 'string' | 'bool' | 'dict' | 'port'
  defaultValue?: string
  dictCategory?: string
  dictName?: string
  portName?: string
}

// 存储每个模块的插件和参数数据
const plugins = reactive<
  Record<
    string,
    {
      name: string
      hash: string
      parameter: string
      parameterList?: string
      help: string
      introduction: string
      enabled: boolean
    }[]
  >
>({})
// 存储每个插件的参数列表（结构化数据）
const parameterLists = reactive<Record<string, Record<string, ParameterItem[]>>>({})
// 存储生成的参数字符串（用于显示和保存）
const parameters = reactive<Record<string, Record<string, string>>>({})
const selectPlugin = reactive<Record<string, string[]>>({})

// 字典和端口数据
const dictList = ref<fileData[]>([])
const portList = ref<portDictData[]>([])

// 加载字典列表
const loadDictList = async () => {
  try {
    const res = await getManagetListApi()
    if (res.code === 200) {
      dictList.value = res.data.list
    }
  } catch (error) {
    console.error('加载字典列表失败:', error)
  }
}

// 加载端口列表
const loadPortList = async () => {
  try {
    const res = await getPortDictDataApi('', 1, 1000)
    if (res.code === 200) {
      portList.value = res.data.list
    }
  } catch (error) {
    console.error('加载端口列表失败:', error)
  }
}

// 格式化参数列表为参数字符串
const formatParameters = (params: ParameterItem[]): string => {
  if (!params || params.length === 0) {
    return ''
  }
  return params
    .filter((param) => param && param.name && param.type)
    .map((param) => {
      if (param.type === 'dict') {
        // dict 类型：-name {dict.分类.名称}
        if (param.dictCategory && param.dictName) {
          return `-${param.name} {dict.${param.dictCategory}.${param.dictName}}`
        }
        return null
      }
      if (param.type === 'port') {
        // port 类型：-name {port.name}
        if (param.portName) {
          return `-${param.name} {port.${param.portName}}`
        }
        return null
      }
      // string 和 bool 类型：只有当默认值有值时才格式化
      if (
        param.defaultValue !== undefined &&
        param.defaultValue !== null &&
        param.defaultValue !== ''
      ) {
        return `-${param.name} ${param.defaultValue}`
      }
      return null
    })
    .filter((item) => item !== null)
    .join(' ')
}

// 解析参数字符串为参数列表（用于向后兼容，如果parameterList不存在时）
const parseParameters = (paramStr: string): ParameterItem[] => {
  if (!paramStr || paramStr.trim() === '') {
    return []
  }
  // 简单的解析逻辑，将 -name value 格式解析为参数列表
  const parts = paramStr.match(/-(\w+)\s+(\S+)/g)
  if (!parts) {
    return []
  }
  return parts
    .map((part) => {
      const match = part.match(/-(\w+)\s+(.+)/)
      if (match) {
        const name = match[1]
        const value = match[2]
        // 判断是否为 bool 类型
        if (value === 'true' || value === 'false') {
          return { name, type: 'bool', defaultValue: value }
        }
        // 判断是否为 dict 类型
        if (value.startsWith('{dict.')) {
          const dictMatch = value.match(/{dict\.(.+)\.(.+)}/)
          if (dictMatch) {
            return { name, type: 'dict', dictCategory: dictMatch[1], dictName: dictMatch[2] }
          }
        }
        // 判断是否为 port 类型
        if (value.startsWith('{port.')) {
          const portMatch = value.match(/{port\.(.+)}/)
          if (portMatch) {
            return { name, type: 'port', portName: portMatch[1] }
          }
        }
        // 默认 string 类型
        return { name, type: 'string', defaultValue: value }
      }
      return null
    })
    .filter((item) => item !== null) as ParameterItem[]
}

// 处理参数列表变化，更新参数字符串
const handleParameterListChange = (module: string, hash: string) => {
  setTimeout(() => {
    const params = parameterLists[module]?.[hash] || []
    const formattedParams = formatParameters(params)
    if (!parameters[module]) {
      parameters[module] = {}
    }
    parameters[module][hash] = formattedParams
  }, 0)
}

// 根据分类筛选字典名称
const getFilteredDictNames = (category: string | undefined) => {
  if (!category) {
    return []
  }
  return dictList.value.filter((dict) => (dict.category || dict.name) === category)
}

// 处理分类改变，清空名称
const handleDictCategoryChange = (module: string, hash: string, index: number) => {
  const param = parameterLists[module][hash][index]
  param.dictName = undefined
  handleParameterListChange(module, hash)
}

// 添加参数（只允许添加string类型）
const addParameter = (module: string, hash: string) => {
  if (!parameterLists[module]) {
    parameterLists[module] = {}
  }
  if (!parameterLists[module][hash]) {
    parameterLists[module][hash] = []
  }
  parameterLists[module][hash].push({
    name: '',
    type: 'string',
    defaultValue: ''
  })
  handleParameterListChange(module, hash)
}

// 删除参数
const removeParameter = (module: string, hash: string, index: number) => {
  parameterLists[module][hash].splice(index, 1)
  handleParameterListChange(module, hash)
}

// 初始化插件数据
const initPlugins = async () => {
  for (const module of modules) {
    const modulePlugins = await getPluginDataByModuleApi(module) // 调用实际接口获取插件数据
    parameters[module] = {} // 初始化空的参数对象
    parameterLists[module] = {} // 初始化空的参数列表对象
    selectPlugin[module] = []

    plugins[module] = modulePlugins.data.list.map((plugin) => {
      // 如果有 parameterList，解析它；否则解析 parameter 字符串
      let paramList: ParameterItem[] = []
      if (plugin.parameterList) {
        try {
          paramList = JSON.parse(plugin.parameterList) as ParameterItem[]
        } catch (error) {
          console.error('解析 parameterList 失败:', error)
          // 如果解析失败，尝试从 parameter 字符串解析
          paramList = parseParameters(plugin.parameter || '')
        }
      } else {
        paramList = parseParameters(plugin.parameter || '')
      }

      // 存储参数列表
      if (!parameterLists[module]) {
        parameterLists[module] = {}
      }
      parameterLists[module][plugin.hash] = paramList

      // 生成参数字符串
      const formattedParams = formatParameters(paramList)
      parameters[module][plugin.hash] = formattedParams

      return {
        ...plugin,
        enabled: false // 初始化时开关为关闭状态
      }
    })
  }
}

const vulList = ref<string[]>([])
// 根据 ID 加载模板数据
const loadTemplate = async (id: string) => {
  const template = await getTemplateDetailApi(id) // 调用实际接口获取已有数据
  console.log(template)
  templateName.value = template.data.name
  vulList.value = template.data.vullist
  for (const module of modules) {
    parameters[module] = {}
    parameterLists[module] = {}

    const moduleData = template.data?.[module] || []
    const modulePlugins = await getPluginDataByModuleApi(module) // 获取模块对应的插件

    plugins[module] = modulePlugins.data.list.map((plugin) => {
      // 优先使用保存的 parameterList，如果没有则尝试老版本的 Parameters，最后使用插件的默认 parameterList
      const templateData = template.data as any
      const savedParameterList = templateData.ParameterLists?.[module]?.[plugin.hash]
      const savedParameterString = templateData.Parameters?.[module]?.[plugin.hash] // 老版本存储的参数字符串
      let paramList: ParameterItem[] = []

      if (savedParameterList) {
        // 新版本：如果保存的是 parameterList（JSON字符串），直接解析
        if (typeof savedParameterList === 'string') {
          try {
            paramList = JSON.parse(savedParameterList) as ParameterItem[]
          } catch (error) {
            console.error('解析保存的 parameterList 失败:', error)
            // 解析失败，尝试从保存的参数字符串解析
            if (savedParameterString && typeof savedParameterString === 'string') {
              paramList = parseParameters(savedParameterString)
            } else {
              // 回退到使用插件默认 parameterList
              paramList = plugin.parameterList
                ? (JSON.parse(plugin.parameterList) as ParameterItem[])
                : []
            }
          }
        } else if (Array.isArray(savedParameterList)) {
          // 如果保存的已经是数组格式，直接使用
          paramList = savedParameterList as ParameterItem[]
        } else {
          // 其他格式，尝试从保存的参数字符串解析
          if (savedParameterString && typeof savedParameterString === 'string') {
            paramList = parseParameters(savedParameterString)
          } else {
            // 回退到使用插件默认 parameterList
            paramList = plugin.parameterList
              ? (JSON.parse(plugin.parameterList) as ParameterItem[])
              : []
          }
        }
      } else if (savedParameterString && typeof savedParameterString === 'string') {
        // 老版本兼容：如果只有参数字符串，解析它
        // 但需要与插件的 parameterList 合并，确保不丢失未配置的参数
        const parsedParams = parseParameters(savedParameterString)

        if (plugin.parameterList) {
          // 如果有插件的 parameterList，合并两者
          try {
            const pluginParamList = JSON.parse(plugin.parameterList) as ParameterItem[]
            // 将解析出的参数值更新到插件的参数列表中
            paramList = pluginParamList.map((pluginParam) => {
              const savedParam = parsedParams.find((p) => p.name === pluginParam.name)
              if (savedParam) {
                // 如果模板中有配置，使用模板的值，但保留插件的类型和其他信息
                return {
                  ...pluginParam,
                  defaultValue: savedParam.defaultValue,
                  dictCategory: savedParam.dictCategory,
                  dictName: savedParam.dictName,
                  portName: savedParam.portName
                }
              }
              return pluginParam
            })
          } catch (error) {
            console.error('解析插件 parameterList 失败:', error)
            // 解析失败，直接使用从参数字符串解析的结果
            paramList = parsedParams
          }
        } else {
          // 插件没有 parameterList，直接使用解析结果
          paramList = parsedParams
        }
      } else if (plugin.parameterList) {
        // 没有保存的参数，使用插件的默认 parameterList
        try {
          paramList = JSON.parse(plugin.parameterList) as ParameterItem[]
        } catch (error) {
          console.error('解析插件 parameterList 失败:', error)
          paramList = []
        }
      } else {
        // 插件也没有 parameterList，尝试从 parameter 字符串解析（最后的后向兼容）
        paramList = parseParameters(plugin.parameter || '')
      }

      // 存储参数列表
      parameterLists[module][plugin.hash] = paramList

      // 生成参数字符串（用于显示和向后兼容）
      const formattedParams = formatParameters(paramList)
      parameters[module][plugin.hash] = formattedParams

      return {
        ...plugin,
        enabled: moduleData.includes(plugin.hash) || false // 判断插件是否启用
      }
    })
  }
}

// 监听 id 的变化来判断是创建还是编辑模式
watch(
  () => props.id,
  async (newId) => {
    if (newId) {
      await loadTemplate(newId) // 如果传入了 ID，则进入编辑模式，加载已有数据
    } else {
      await initPlugins() // 否则初始化空白表单，表示创建新项
    }
  },
  { immediate: true } // 确保组件挂载时立即触发
)
const saveLoading = ref(false)

// 提交表单数据
const onSubmit = async () => {
  saveLoading.value = true
  const result: Record<string, any> = {}
  if (templateName.value == '') {
    ElMessage.error('name 不能为空')
    saveLoading.value = false
    return
  }
  // 收集每个模块启用的插件和对应的参数
  for (const module of modules) {
    const enabledPlugins = plugins[module].filter((plugin) => plugin.enabled)
    result[module] = enabledPlugins.map((plugin) => plugin.hash)

    // 保存参数字符串（用于向后兼容和执行）
    result.Parameters = result.Parameters || {}
    result.Parameters[module] = {}

    // 保存 parameterList（JSON字符串，用于完整恢复参数配置）
    result.ParameterLists = result.ParameterLists || {}
    result.ParameterLists[module] = {}

    // 只收集已启用插件的参数
    for (const plugin of enabledPlugins) {
      // 保存参数字符串
      if (parameters[module]?.[plugin.hash]) {
        result.Parameters[module][plugin.hash] = parameters[module][plugin.hash]
      }

      // 保存 parameterList（JSON字符串）
      if (parameterLists[module]?.[plugin.hash]) {
        result.ParameterLists[module][plugin.hash] = JSON.stringify(
          parameterLists[module][plugin.hash]
        )
      }
    }
  }
  result['name'] = templateName.value
  result['vullist'] = vulList.value
  try {
    const res = await saveTemplateDetailApi(result, props.id)
    console.log(result)

    if (res.code === 200) {
      ElMessage.success('success')
      // 提交成功才执行父组件逻辑
      props.closeDialog()
      props.getList()
    } else {
      // 提交失败
      ElMessage.error(res.message || 'error')
    }
  } catch (error) {
    ElMessage.error('error')
    console.error(error)
  } finally {
    saveLoading.value = false
  }
}
const moduleColorMap = {
  TargetHandler: '#409EFF',
  SubdomainScan: '#E6A23C',
  SubdomainSecurity: '#F56C6C',
  PortScanPreparation: '#67C23A',
  PortScan: '#00CED1',
  AssetMapping: '#8A2BE2',
  URLScan: '#C71585',
  WebCrawler: '#FF4500',
  DirScan: '#20B2AA',
  VulnerabilityScan: '#DC143C',
  AssetHandle: '#4682B4',
  PortFingerprint: '#DAA520',
  URLSecurity: '#9370DB',
  PassiveScan: '#5F9EA0'
}
const templateName = ref('')
interface TreeNode {
  value: string
  label: string
  children: TreeNode[]
}

const vulOptions = reactive<TreeNode[]>([])

const buildTree = (data: pocData[]): TreeNode[] => {
  const tree: TreeNode[] = []

  data.forEach((item) => {
    let currentLevel = tree
    item.tags.forEach((tag) => {
      // 查找当前层级是否有该标签
      const existingNode = currentLevel.find((node) => node.label === tag)
      if (!existingNode) {
        // 如果没有找到，则生成一个分类节点
        const randomString = Math.random().toString(36).substring(2, 8)
        const newNode: TreeNode = { value: randomString, label: tag, children: [] }
        currentLevel.push(newNode)
        currentLevel = newNode.children // 进入下一层
      } else {
        currentLevel = existingNode.children // 如果找到了，继续向下
      }
    })

    // 添加实际数据节点
    currentLevel.push({ value: item.id, label: item.name, children: [] })
  })

  return tree
}

const vulSelectOptions = reactive<{ value: string; label: string }[]>([])
// 获取数据并生成树形结构
const getPocList = async () => {
  const res = await getPocDataAllApi() // 调用后端API
  if (res.data.list.length > 0) {
    vulSelectOptions.push({ value: 'All Poc', label: 'All Poc' })
    res.data.list.forEach((item) => {
      vulSelectOptions.push({ value: item.id, label: item.name })
    })
  }

  vulOptions.push({ value: 'All Poc', label: 'All Poc', children: [] })
  const tree = buildTree(res.data.list)
  vulOptions.push(...tree)
}

getPocList()
// 初始化时加载字典和端口列表
loadDictList()
loadPortList()

const dialogVisible = ref(false)
const openPocList = async () => {
  dialogVisible.value = true
}

// 参数配置抽屉
const parameterDrawerVisible = ref(false)
const currentPluginModule = ref('')
const currentPluginHash = ref('')
const currentPluginName = ref('')
const currentPluginHelp = ref('')
const drawerSize = ref('50%')

// 计算抽屉尺寸
const calculateDrawerSize = () => {
  const screenWidth = window.innerWidth
  if (screenWidth < 768) {
    drawerSize.value = '90%'
  } else if (screenWidth < 1024) {
    drawerSize.value = '60%'
  } else {
    drawerSize.value = '50%'
  }
}

// 打开参数配置抽屉
const openParameterDialog = (module: string, hash: string, name: string) => {
  currentPluginModule.value = module
  currentPluginHash.value = hash
  currentPluginName.value = name
  // 获取当前插件的 help 信息
  const plugin = plugins[module]?.find((p) => p.hash === hash)
  currentPluginHelp.value = plugin?.help || ''
  calculateDrawerSize()
  parameterDrawerVisible.value = true
}

// 关闭参数配置抽屉
const closeParameterDialog = () => {
  parameterDrawerVisible.value = false
  currentPluginModule.value = ''
  currentPluginHash.value = ''
  currentPluginName.value = ''
  currentPluginHelp.value = ''
}
watch(dialogVisible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      const tree = treeRef.value
      if (tree) {
        console.log('treeRef 已经获取到实例:', tree)
        setTreeCheckedNodes(vulList.value) // 初始化选中的节点
      } else {
        console.log('treeRef 未能正确获取到实例')
      }
    })
  }
})
const treeRef = ref<InstanceType<typeof ElTreeV2> | null>(null)
const setTreeCheckedNodes = (checkedKeys) => {
  nextTick(() => {
    const tree = treeRef.value
    if (tree) {
      tree.setCheckedKeys(checkedKeys) // 使用 setCheckedKeys 方法来批量选中指定的节点
    } else {
      console.log('treeRef 未能正确获取到实例')
    }
  })
}
const propss = {
  value: 'value',
  label: 'label',
  children: 'children'
}

const handleCheckChange = (data, checked) => {
  const nodeValue = data.value // 当前节点的 value

  // 判断当前节点是否是叶子节点
  const isLeafNode = !data.children || data.children.length === 0

  // 如果是叶子节点，直接处理它
  if (isLeafNode) {
    if (checked && !vulList.value.includes(nodeValue)) {
      vulList.value.push(nodeValue) // 选中叶子节点，添加到 vulList
    } else if (!checked) {
      const index = vulList.value.indexOf(nodeValue)
      if (index > -1) {
        vulList.value.splice(index, 1) // 取消选中叶子节点，移除它
      }
    }
  } else {
    // 如果是父节点，遍历它的所有子节点并处理
    const addLeafNodes = (node) => {
      if (node.children && node.children.length > 0) {
        node.children.forEach((child) => {
          const isChildLeafNode = !child.children || child.children.length === 0 // 判断子节点是否为叶子节点
          if (isChildLeafNode) {
            // 只处理叶子节点
            if (checked && !vulList.value.includes(child.value)) {
              vulList.value.push(child.value) // 添加叶子节点到 vulList
            } else if (!checked) {
              const index = vulList.value.indexOf(child.value)
              if (index > -1) {
                vulList.value.splice(index, 1) // 移除取消选中的叶子节点
              }
            }
          } else {
            // 如果是父节点，则递归处理其子节点
            addLeafNodes(child)
          }
        })
      }
    }

    // 如果是父节点，递归地添加或移除其叶子节点
    addLeafNodes(data)
  }

  // 打印当前选中的叶子节点的值
  console.log('当前选中的叶子节点的值:', vulList.value)
}
</script>

<template>
  <ElForm @submit.prevent="onSubmit" label-width="auto">
    <ElFormItem :label="t('task.templateName')">
      <ElInput v-model="templateName" />
    </ElFormItem>
    <div v-for="module in modules" :key="module">
      <ElCard class="module-card" :body-style="{ padding: '20px' }" shadow="always">
        <div
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            border-bottom: 1px solid #f0f2f5;
            padding-bottom: 15px;
          "
        >
          <div style="display: flex; align-items: center; gap: 10px">
            <div
              :style="{
                width: '4px',
                height: '16px',
                backgroundColor: moduleColorMap[module],
                borderRadius: '2px'
              }"
            ></div>
            <span style="font-weight: 600; font-size: 16px; color: #303133">
              {{ t(`scanTemplate.${module}`) }}
            </span>
          </div>
        </div>

        <div class="plugins-container">
          <ElCard
            v-for="plugin in plugins[module]"
            :key="plugin.hash"
            :class="['plugin-card', { 'plugin-card-enabled': plugin.enabled }]"
            :style="{ '--card-accent-color': moduleColorMap[module] }"
            :body-style="{ padding: '0', height: '100%', display: 'flex', flexDirection: 'column' }"
            shadow="hover"
          >
            <div class="plugin-card-header">
              <div class="plugin-title">
                <div class="plugin-name-wrapper">
                  <span class="plugin-name" :title="plugin.name">{{ plugin.name }}</span>
                </div>
                <ElTooltip
                  placement="top"
                  effect="dark"
                  :content="plugin.enabled ? t('common.enabled') : t('common.disabled')"
                >
                  <ElSwitch
                    v-model="plugin.enabled"
                    class="plugin-switch"
                    style="--el-switch-on-color: var(--card-accent-color)"
                  />
                </ElTooltip>
              </div>
              <div class="plugin-desc" v-if="plugin.introduction">
                <ElTooltip placement="top" effect="light" :content="plugin.introduction">
                  <span class="text-truncate">{{ plugin.introduction }}</span>
                </ElTooltip>
              </div>
            </div>
            <div class="plugin-card-body">
              <ElFormItem
                :label="t('task.vulList')"
                prop="type"
                v-if="plugin.enabled && plugin.hash === 'ed93b8af6b72fe54a60efdb932cf6fbc'"
              >
                <ElSelectV2
                  v-model="vulList"
                  filterable
                  :options="vulSelectOptions"
                  placeholder="Please select vul"
                  style="width: 80%; margin-right: 10px"
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  tag-type="info"
                  :max-collapse-tags="3"
                />
                <ElButton type="primary" @click="openPocList" :loading="saveLoading">
                  {{ t('common.selectCategory') }}
                </ElButton>
              </ElFormItem>
              <!-- 参数预览和配置按钮 -->
              <ElFormItem v-if="plugin.enabled" :label="t('plugin.parameter')">
                <div style="display: flex; gap: 8px; align-items: flex-start">
                  <ElInput
                    :model-value="parameters[module]?.[plugin.hash] || ''"
                    readonly
                    style="background-color: #f5f5f5; flex: 1"
                    type="textarea"
                    :rows="2"
                  />
                  <ElButton
                    type="primary"
                    size="small"
                    @click="openParameterDialog(module, plugin.hash, plugin.name)"
                  >
                    {{ t('plugin.parameterConfig') }}
                  </ElButton>
                </div>
              </ElFormItem>
            </div>
          </ElCard>
        </div>
      </ElCard>
    </div>
    <ElRow>
      <ElCol :span="12" style="text-align: right">
        <ElButton type="primary" @click="onSubmit" :loading="saveLoading"> 保存 </ElButton>
      </ElCol>
    </ElRow>
  </ElForm>
  <Dialog
    v-model="dialogVisible"
    title="POC"
    center
    fullscreen
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <ElTreeV2
      ref="treeRef"
      style="max-width: 100%"
      :data="vulOptions"
      :props="propss"
      show-checkbox
      :height="600"
      @check-change="handleCheckChange"
    />
  </Dialog>

  <!-- 参数配置抽屉 -->
  <ElDrawer
    v-model="parameterDrawerVisible"
    :title="
      currentPluginName
        ? `${currentPluginName} - ${t('plugin.parameterConfig')}`
        : t('plugin.parameterConfig')
    "
    :size="drawerSize"
    direction="rtl"
    :close-on-click-modal="false"
  >
    <ElForm label-width="100px">
      <!-- 插件帮助信息 -->
      <ElRow v-if="currentPluginHelp" style="margin-bottom: 20px">
        <ElCol :span="24">
          <ElCard shadow="never" style="background-color: #f0f9ff; border: 1px solid #b3d8ff">
            <div style="display: flex; align-items: flex-start; gap: 8px">
              <span style="font-weight: 600; color: #409eff; flex-shrink: 0">
                {{ t('plugin.help') }}:
              </span>
              <span style="color: #606266; line-height: 1.6">{{ currentPluginHelp }}</span>
            </div>
          </ElCard>
        </ElCol>
      </ElRow>
      <ElRow :gutter="20" v-if="parameterLists[currentPluginModule]?.[currentPluginHash]?.length">
        <template
          v-for="(param, index) in parameterLists[currentPluginModule]?.[currentPluginHash] || []"
          :key="index"
        >
          <ElCol :span="24" style="margin-bottom: 16px">
            <div style="padding: 12px; border: 1px solid #dcdfe6; border-radius: 4px">
              <ElRow :gutter="10">
                <!-- 参数名称 -->
                <ElCol :span="24">
                  <ElFormItem
                    :prop="`parameterList.${currentPluginModule}.${currentPluginHash}.${index}.name`"
                    style="margin-bottom: 8px"
                    :label="t('plugin.parameterName')"
                  >
                    <div style="display: flex; align-items: center; gap: 8px">
                      <ElInput
                        v-model="param.name"
                        :placeholder="t('plugin.parameterName')"
                        :readonly="param.type !== 'string'"
                        :disabled="param.type !== 'string'"
                        style="flex: 1"
                        @input="handleParameterListChange(currentPluginModule, currentPluginHash)"
                      />
                      <ElButton
                        :icon="Delete"
                        type="danger"
                        circle
                        size="small"
                        @click="removeParameter(currentPluginModule, currentPluginHash, index)"
                      />
                    </div>
                  </ElFormItem>
                </ElCol>

                <!-- string 和 bool 类型的默认值 -->
                <template v-if="param.type === 'string' || param.type === 'bool'">
                  <ElCol :span="24">
                    <ElFormItem
                      :prop="`parameterList.${currentPluginModule}.${currentPluginHash}.${index}.defaultValue`"
                      style="margin-bottom: 0"
                      :label="t('plugin.defaultValue')"
                    >
                      <ElSelect
                        v-if="param.type === 'bool'"
                        v-model="param.defaultValue"
                        :placeholder="t('plugin.defaultValue')"
                        @change="handleParameterListChange(currentPluginModule, currentPluginHash)"
                      >
                        <ElOption label="true" value="true" />
                        <ElOption label="false" value="false" />
                      </ElSelect>
                      <ElInput
                        v-else
                        v-model="param.defaultValue"
                        :placeholder="t('plugin.defaultValue')"
                        @input="handleParameterListChange(currentPluginModule, currentPluginHash)"
                      />
                    </ElFormItem>
                  </ElCol>
                </template>

                <!-- dict 类型的分类和名称 -->
                <template v-if="param.type === 'dict'">
                  <ElCol :span="24" style="margin-bottom: 8px">
                    <ElFormItem
                      :prop="`parameterList.${currentPluginModule}.${currentPluginHash}.${index}.dictCategory`"
                      style="margin-bottom: 0"
                      :label="t('plugin.category')"
                    >
                      <ElSelect
                        v-model="param.dictCategory"
                        :placeholder="t('plugin.category')"
                        filterable
                        @change="
                          handleDictCategoryChange(currentPluginModule, currentPluginHash, index)
                        "
                      >
                        <ElOption
                          v-for="dict in dictList"
                          :key="dict.id"
                          :label="dict.category || dict.name"
                          :value="dict.category || dict.name"
                        />
                      </ElSelect>
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="24">
                    <ElFormItem
                      :prop="`parameterList.${currentPluginModule}.${currentPluginHash}.${index}.dictName`"
                      style="margin-bottom: 0"
                      :label="t('plugin.dictName')"
                    >
                      <ElSelect
                        v-model="param.dictName"
                        :placeholder="t('plugin.dictName')"
                        filterable
                        :disabled="!param.dictCategory"
                        @change="handleParameterListChange(currentPluginModule, currentPluginHash)"
                      >
                        <ElOption
                          v-for="dict in getFilteredDictNames(param.dictCategory)"
                          :key="dict.id"
                          :label="dict.name"
                          :value="dict.name"
                        />
                      </ElSelect>
                    </ElFormItem>
                  </ElCol>
                </template>

                <!-- port 类型的 name -->
                <template v-if="param.type === 'port'">
                  <ElCol :span="24">
                    <ElFormItem
                      :prop="`parameterList.${currentPluginModule}.${currentPluginHash}.${index}.portName`"
                      style="margin-bottom: 0"
                      :label="t('plugin.portName')"
                    >
                      <ElSelect
                        v-model="param.portName"
                        :placeholder="t('plugin.portName')"
                        filterable
                        @change="handleParameterListChange(currentPluginModule, currentPluginHash)"
                      >
                        <ElOption
                          v-for="port in portList"
                          :key="port.id"
                          :label="port.name"
                          :value="port.name"
                        />
                      </ElSelect>
                    </ElFormItem>
                  </ElCol>
                </template>
              </ElRow>
            </div>
          </ElCol>
        </template>
      </ElRow>
      <div v-else style="text-align: center; padding: 40px; color: #909399">
        {{ t('plugin.noParameters') }}
      </div>

      <!-- 添加参数按钮 -->
      <ElRow style="margin-top: 20px">
        <ElCol :span="24">
          <ElButton
            :icon="Plus"
            style="width: 100%"
            @click="addParameter(currentPluginModule, currentPluginHash)"
          >
            {{ t('plugin.addParameter') }}
          </ElButton>
        </ElCol>
      </ElRow>

      <!-- 参数预览 -->
      <ElFormItem :label="t('plugin.parameter')" style="margin-top: 20px">
        <ElInput
          :model-value="parameters[currentPluginModule]?.[currentPluginHash] || ''"
          readonly
          style="background-color: #f5f5f5"
          type="textarea"
          :rows="3"
        />
        <div style="font-size: 12px; color: #909399; margin-top: 4px">
          {{ t('plugin.parameterTip') }}
        </div>
      </ElFormItem>
    </ElForm>

    <template #footer>
      <div style="text-align: right">
        <ElButton @click="closeParameterDialog">{{ t('common.cancel') }}</ElButton>
        <ElButton type="primary" @click="closeParameterDialog" style="margin-left: 10px">
          {{ t('common.confirmed') }}
        </ElButton>
      </div>
    </template>
  </ElDrawer>
</template>

<style scoped>
/* 样式部分 */
.ElFormItem {
  margin-bottom: 20px;
}

.module-card {
  margin-bottom: 24px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  background: #ffffff;
  box-shadow: none;
}

.module-card:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.plugins-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.plugin-card {
  border-radius: 4px;
  border: 1px solid #e4e7ed;
  background-color: #ffffff;
  transition: all 0.2s ease-in-out;
  position: relative;
  border-left: 4px solid var(--card-accent-color, #409eff);
}

.plugin-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #c0c4cc;
}

.plugin-card-enabled {
  background: #fcfcfc;
}

.plugin-card-header {
  padding: 12px 16px;
  border-bottom: 1px solid #f2f6fc;
  background: #fff;
  border-radius: 4px 4px 0 0;
}

.plugin-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.plugin-name-wrapper {
  flex: 1;
  overflow: hidden;
  margin-right: 12px;
}

.plugin-name {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

.plugin-desc {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}

.text-truncate {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.plugin-card-body {
  padding: 16px;
  background-color: #ffffff;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  border-radius: 0 0 4px 4px;
}

.plugin-switch {
  --el-switch-on-color: var(--card-accent-color);
}
/* Customizing the parameter input area */
.plugin-card-body :deep(.el-textarea__inner) {
  background-color: #f8fafc !important;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  color: #475569;
  box-shadow: none;
}

.plugin-card-body :deep(.el-textarea__inner:hover) {
  border-color: #cbd5e1;
}

.plugin-card-body :deep(.el-button--small) {
  padding: 8px 12px;
  font-weight: 500;
}
</style>

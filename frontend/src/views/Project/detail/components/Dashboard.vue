<script setup lang="tsx">
import { useI18n } from '@/hooks/web/useI18n'
import { reactive, Ref, ref, watch } from 'vue'
import { Echart } from '@/components/Echart'
import {
  ElRow,
  ElCol,
  ElCard,
  ElStatistic,
  ElContainer,
  ElFooter,
  ElDescriptions,
  ElDescriptionsItem,
  ElTag,
  ElDivider,
  ElSkeleton,
  ElTable,
  ElTableColumn,
  ElText,
  ElSpace,
  ElCollapseItem,
  ElCollapse
} from 'element-plus'
import { useRoute } from 'vue-router'
import { useIcon } from '@/hooks/web/useIcon'
import {
  getProjectAssetCountApi,
  getProjectInfoApi,
  getProjectVulDataApi,
  getProjectVulLevelCountApi
} from '@/api/ProjectAggregation'
import { EChartsOption } from 'echarts'
import {
  getAssetStatisticsappApi,
  getAssetStatisticsPortApi,
  getAssetStatisticsTitleApi,
  getAssetStatisticsTypeApi
} from '@/api/asset'
const { t } = useI18n()
const { query } = useRoute()
let projectInfo = reactive({
  name: '',
  tag: '',
  scheduledTasks: false,
  hour: 0,
  AssetCount: 0,
  domain: [''],
  nextTime: ''
})

let projectAssetCount = reactive({
  subdomainCount: 0,
  vulCount: 0
})
const projectInfoLoading = ref(false)
const getProjectInfo = async () => {
  try {
    projectInfoLoading.value = true

    const [res1, res2] = await Promise.all([
      getProjectInfoApi(query.id as string),
      getProjectAssetCountApi(query.id as string)
    ])

    projectInfo.name = res1.data.name
    projectInfo.tag = res1.data.tag
    projectInfo.scheduledTasks = res1.data.scheduledTasks
    projectInfo.hour = res1.data.hour
    projectInfo.AssetCount = res1.data.AssetCount
    projectInfo.domain = res1.data.root_domains
    projectInfo.nextTime = res1.data.next_time

    projectAssetCount.subdomainCount = res2.data.subdomainCount
    projectAssetCount.vulCount = res2.data.vulCount
  } catch (error) {
    console.error('Error fetching project info:', error)
  } finally {
    projectInfoLoading.value = false
  }
}

const projectNameIcon = useIcon({ icon: 'icon-park:edit-name' })
const projectTagIcon = useIcon({ icon: 'icon-park:tag' })
const projectScopeIcon = useIcon({ icon: 'zondicons:network', color: '#40c9c6' })
const projectTaskIcon = useIcon({ icon: 'hugeicons:task-done-01', color: '#36a3f7' })
const projectCycleIcon = useIcon({ icon: 'icon-park-outline:cycle', color: '#36a3f7' })
const projectNextTimeIcon = useIcon({ icon: 'tdesign:time', color: '#f4516c' })
const levelMap = {
  critical: { color: '#E74C3C', flag: t('poc.critical') },
  high: { color: '#F39C12', flag: t('poc.high') },
  medium: { color: '#F1C40F', flag: t('poc.medium') },
  low: { color: '#3498DB', flag: t('poc.low') },
  info: { color: '#2ECC71', flag: t('poc.info') },
  unknown: { color: '#95A5A6', flag: t('poc.unknown') }
}
const getVulTagType = (level) => {
  switch (level) {
    case 'critical':
      return 'danger'
    case 'high':
      return 'warning'
    case 'medium':
      return 'primary'
    case 'low':
      return 'success'
    case 'info':
      return 'info'
    case 'unknown':
      return 'info'
    default:
      return 'info'
  }
}
const getVulLevelName = (level) => {
  return levelMap[level].flag
}
const vulLevelloading = ref(false)
let chartData: any[] = []
const getProjectVulLevelInfo = async () => {
  vulLevelloading.value = true
  const res = await getProjectVulLevelCountApi(query.id as string)
  if (Array.isArray(res.data)) {
    chartData = res.data
      .map((item) => {
        if (item._id && item.count && levelMap[item._id]) {
          const { color, flag } = levelMap[item._id]
          return {
            value: item.count,
            name: item._id,
            itemStyle: { color }
          }
        }
        return null
      })
      .filter(Boolean)
  }
  chartOptions.series![0].data = chartData
  vulLevelloading.value = false
}
const chartOptions: EChartsOption = reactive({
  tooltip: {
    trigger: 'item'
  },
  legend: {
    orient: 'horizontal',
    bottom: '20%',
    textStyle: {
      fontSize: 11,
      fontWeight: 'bold',
      color: '#333'
    }
  },
  series: [
    {
      name: 'Vulnerability',
      type: 'pie',
      center: ['50%', '30%'],
      radius: ['30%', '60%'],
      avoidLabelOverlap: false,
      padAngle: 4,
      itemStyle: {
        borderRadius: 10
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 10,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: []
    }
  ]
})

interface VulData {
  time: string
  level: string
  vulname: string
  url: string
  matched: string
}
const projectVulDataLoading = ref(false)
let vulTableData = reactive<VulData[]>([])
const getProjectVulData = async () => {
  projectVulDataLoading.value = true
  const res = await getProjectVulDataApi(query.id as string)
  if (Array.isArray(res.data.list)) {
    // 更新 vulTableData 的值，使用 .value 属性访问 ref 的值
    vulTableData = res.data.list
  }
  projectVulDataLoading.value = false
}

let AssetstatisticsData: Ref<{
  Port: { value: number; number: number }[]
  Service: { value: string; number: number }[]
  Product: { value: string; number: number }[]
  Icon: { value: string; number: number; icon_hash: string }[]
  Title: { value: string; number: number }[]
}> = ref({
  Port: [],
  Service: [],
  Product: [],
  Icon: [],
  Title: []
})
const projectStatisLoading = ref(false)
const filter = reactive<{ [key: string]: any }>({})
const getAssetstatistics = async () => {
  projectStatisLoading.value = true
  filter.project = [query.id as string]
  AssetstatisticsData.value.Port = []
  AssetstatisticsData.value.Service = []
  AssetstatisticsData.value.Product = []
  AssetstatisticsData.value.Icon = []
  AssetstatisticsData.value.Title = []
  const [portRes, serviceRes, productRes, titleRes] = await Promise.all([
    getAssetStatisticsPortApi('', filter),
    getAssetStatisticsTypeApi('', filter),
    getAssetStatisticsappApi('', filter),
    getAssetStatisticsTitleApi('', filter)
  ])

  AssetstatisticsData.value.Port = portRes.data.Port
  AssetstatisticsData.value.Service = serviceRes.data.Service
  AssetstatisticsData.value.Product = productRes.data.Product
  AssetstatisticsData.value.Title = titleRes.data.Title
  projectStatisLoading.value = false
}

const getAllApi = async () => {
  await Promise.all([
    getProjectInfo(),
    getProjectVulLevelInfo(),
    getProjectVulData(),
    getAssetstatistics()
  ])
}
getAllApi()
</script>

<template>
  <ElRow :gutter="20">
    <ElCol :xl="8" :lg="8" :md="8" :sm="8" :xs="8">
      <ElCard shadow="hover" class="mb-20px">
        <ElSkeleton :loading="projectInfoLoading" animated>
          <ElContainer>
            <ElMain>
              <ElDescriptions :column="2" direction="vertical" :border="true">
                <ElDescriptionsItem labelClassName="projectinfoclass">
                  <template #label>
                    <div class="cell-item">
                      <ElIcon style="position: relative; top: 3px">
                        <projectNameIcon />
                      </ElIcon>
                      {{ t('project.projectName') }}
                    </div>
                  </template>
                  {{ projectInfo.name }}
                </ElDescriptionsItem>
                <ElDescriptionsItem labelClassName="projectinfoclass">
                  <template #label>
                    <div class="cell-item">
                      <ElIcon style="position: relative; top: 3px">
                        <projectTagIcon />
                      </ElIcon>
                      TAG
                    </div>
                  </template>
                  <ElTag> {{ projectInfo.tag }}</ElTag>
                </ElDescriptionsItem>
                <ElDescriptionsItem labelClassName="projectinfoclass">
                  <template #label>
                    <div class="cell-item">
                      <ElIcon style="position: relative; top: 3px">
                        <projectScopeIcon />
                      </ElIcon>
                      {{ t('project.projectScope') }}
                    </div>
                  </template>
                  <ElScrollbar max-height="50px" always>
                    <template v-for="(o, index) in projectInfo.domain" :key="index">
                      <div class="text item">{{ o }}</div>
                    </template>
                  </ElScrollbar>
                </ElDescriptionsItem>
                <ElDescriptionsItem labelClassName="projectinfoclass">
                  <template #label>
                    <div class="cell-item">
                      <ElIcon style="position: relative; top: 3px">
                        <projectTaskIcon />
                      </ElIcon>
                      {{ t('project.scheduledTasks') }}
                    </div>
                  </template>
                  <ElTag :type="projectInfo.scheduledTasks ? 'success' : 'info'">
                    {{
                      projectInfo.scheduledTasks
                        ? t('common.switchAction')
                        : t('common.switchInactive')
                    }}
                  </ElTag>
                </ElDescriptionsItem>
                <ElDescriptionsItem
                  v-if="projectInfo.scheduledTasks"
                  labelClassName="projectinfoclass"
                >
                  <template #label>
                    <div class="cell-item">
                      <ElIcon style="position: relative; top: 3px">
                        <projectCycleIcon />
                      </ElIcon>
                      {{ t('project.cycle') }}(h)
                    </div>
                  </template>
                  {{ projectInfo.hour }}
                </ElDescriptionsItem>
                <ElDescriptionsItem
                  v-if="projectInfo.scheduledTasks"
                  labelClassName="projectinfoclass"
                >
                  <template #label>
                    <div class="cell-item">
                      <ElIcon style="position: relative; top: 3px">
                        <projectNextTimeIcon />
                      </ElIcon>
                      {{ t('task.nextTime') }}
                    </div>
                  </template>
                  {{ projectInfo.nextTime }}
                </ElDescriptionsItem>
              </ElDescriptions>
            </ElMain>
            <ElDivider />
            <ElFooter heigh="50%">
              <ElRow>
                <ElCol :span="8">
                  <ElStatistic
                    :title="t('dashboard.totalAssets')"
                    :value="projectInfo.AssetCount"
                  />
                </ElCol>
                <ElCol :span="8">
                  <ElStatistic
                    :title="t('dashboard.subDomain')"
                    :value="projectAssetCount.subdomainCount"
                  />
                </ElCol>
                <ElCol :span="8">
                  <ElStatistic
                    :title="t('vulnerability.vulnerabilityName')"
                    :value="projectAssetCount.vulCount"
                  />
                </ElCol>
              </ElRow>
            </ElFooter>
          </ElContainer>
        </ElSkeleton>
      </ElCard>
    </ElCol>
    <ElCol :xl="16" :lg="16" :md="16" :sm="16" :xs="16">
      <ElCard shadow="hover" class="mb-20px">
        <ElRow>
          <ElCol :span="5">
            <ElSkeleton :loading="vulLevelloading" animated>
              <Echart :options="chartOptions" :height="353" :width="250" />
            </ElSkeleton>
          </ElCol>

          <ElDivider direction="vertical" />
          <ElCol :span="17">
            <ElSkeleton :loading="projectVulDataLoading" animated>
              <ElTable
                :data="vulTableData"
                :showHeader="false"
                :scrollbarAlwaysOn="true"
                :fit="true"
                :max-height="353"
              >
                <ElTableColumn minWidth="50">
                  <template #default="scope">
                    <ElTag :type="getVulTagType(scope.row.level)" size="small">{{
                      getVulLevelName(scope.row.level)
                    }}</ElTag>
                  </template>
                </ElTableColumn>
                <ElTableColumn :showOverflowTooltip="true" minWidth="300">
                  <template #default="scope">
                    <ElSpace>
                      <ElText>{{ scope.row.vulname }}</ElText>
                      <ElText type="info" size="small">{{
                        scope.row.url ? scope.row.url : scope.row.matched
                      }}</ElText>
                    </ElSpace>
                  </template>
                </ElTableColumn>
                <ElTableColumn :showOverflowTooltip="true" minWidth="200">
                  <template #default="scope">
                    {{ scope.row.time }}
                  </template>
                </ElTableColumn>
              </ElTable>
            </ElSkeleton>
          </ElCol>
        </ElRow>
      </ElCard>
    </ElCol>
  </ElRow>
  <ElRow :gutter="20">
    <ElSkeleton :loading="projectStatisLoading" animated>
      <ElCol :span="4">
        <ElCard shadow="hover" class="mb-20px">
          <template #header>
            <ElText tag="b">{{ t('asset.port') }}</ElText>
          </template>
          <ElRow v-for="portItem in AssetstatisticsData.Port" :key="portItem.value">
            <ElCol :span="12">
              <ElTag effect="light" round size="small">
                {{ portItem.value }}
              </ElTag>
            </ElCol>
            <ElCol :span="12" style="text-align: end">
              <ElText size="small">{{ portItem.number }}</ElText>
            </ElCol>
          </ElRow>
        </ElCard>
      </ElCol>
      <ElCol :span="5">
        <ElCard shadow="hover" class="mb-20px">
          <template #header>
            <ElText tag="b">{{ t('asset.service') }}</ElText>
          </template>
          <ElRow v-for="serviceItem in AssetstatisticsData.Service" :key="serviceItem.value">
            <ElCol :span="12">
              <ElTag effect="light" round size="small">
                {{ serviceItem.value }}
              </ElTag>
            </ElCol>
            <ElCol :span="12" style="text-align: end">
              <ElText size="small">{{ serviceItem.number }}</ElText>
            </ElCol>
          </ElRow>
        </ElCard>
      </ElCol>
      <ElCol :span="6">
        <ElCard shadow="hover" class="mb-20px">
          <template #header>
            <ElText tag="b">{{ t('asset.products') }}</ElText>
          </template>
          <ElRow v-for="productItem in AssetstatisticsData.Product" :key="productItem.value">
            <ElCol :span="12">
              <ElTag effect="light" round size="small">
                {{ productItem.value }}
              </ElTag>
            </ElCol>
            <ElCol :span="12" style="text-align: end">
              <ElText size="small">{{ productItem.number }}</ElText>
            </ElCol>
          </ElRow>
        </ElCard>
      </ElCol>
      <ElCol :span="9">
        <ElCard shadow="hover" class="mb-20px">
          <template #header>
            <ElText tag="b">{{ t('asset.title') }}</ElText>
          </template>
          <ElRow v-for="titleItem in AssetstatisticsData.Title" :key="titleItem.value">
            <ElCol :span="18">
              <ElTag effect="light" round size="small">
                {{ titleItem.value }}
              </ElTag>
            </ElCol>
            <ElCol :span="5" style="text-align: end">
              <ElText size="small">{{ titleItem.number }}</ElText>
            </ElCol>
          </ElRow>
        </ElCard>
      </ElCol>
    </ElSkeleton>
  </ElRow>
</template>
<style scoped>
:deep(.projectinfoclass) {
  background: transparent !important;
  color: var(--el-text-color-primary) !important;
}
.el-divider--vertical {
  height: 13em;
  margin-left: 20px;
}
.el-table__row:hover {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); /* 阴影效果，可以根据需求调整 */
}
</style>

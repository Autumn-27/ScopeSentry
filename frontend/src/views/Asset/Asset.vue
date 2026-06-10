<script setup lang="ts">
import { useI18n } from '@/hooks/web/useI18n'
import { ElTabs } from 'element-plus'
import { ElTabPane } from 'element-plus'
import AssetInfo2 from './components/AssetInfo2.vue'
import Subdomain from './components/Subdomain.vue'
import URL from './components/URL.vue'
import Crawler from './components/Crawler.vue'
import SensitiveInformation from './components/SensitiveInformation.vue'
import DirScan from './components/DirScan.vue'
import PageMonitoring from './components/PageMonitoring.vue'
import vul from './components/vul.vue'
import SubdomainTakeover from './components/SubdomainTakeover.vue'
import { reactive } from 'vue'
import { getProjectAllApi } from '@/api/project'
import { getTaskNamesApi } from '@/api/task'
import RootDomain from './components/RootDomain.vue'
import APP from './components/APP.vue'
import MP from './components/MP.vue'
import IP from './components/IP.vue'
const { t } = useI18n()
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

interface TaskNameData {
  id: string
  name: string
}
const taskList = reactive<TaskNameData[]>([])
const getTaskList = async () => {
  const res = await getTaskNamesApi()
  taskList.splice(0, taskList.length, ...(res.data || []))
}
getTaskList()

const handleTabClick = (tab: any) => {
  if (tab.paneName === 'map') {
    // 跳转到同一网站的map项目
    window.location.href = '/map'
  }
}
</script>

<template>
  <ElTabs type="border-card" @tab-click="handleTabClick">
    <ElTabPane :label="t('asset.assetName')"
      ><AssetInfo2 :projectList="projectList" :taskList="taskList"
    /></ElTabPane>
    <ElTabPane label="IP"><IP :projectList="projectList" :taskList="taskList" /></ElTabPane>
    <ElTabPane :label="t('rootDomain.rootDomainName')">
      <RootDomain :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('subdomain.subdomainName')">
      <Subdomain :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('task.subdomainTakeover')">
      <SubdomainTakeover :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('app.appName')">
      <APP :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('miniProgram.miniProgramName')">
      <MP :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('URL.URLName')"
      ><URL :projectList="projectList" :taskList="taskList"
    /></ElTabPane>
    <ElTabPane :label="t('crawler.crawlerName')"
      ><Crawler :projectList="projectList" :taskList="taskList"
    /></ElTabPane>
    <ElTabPane :label="t('sensitiveInformation.sensitiveInformationName')">
      <SensitiveInformation :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('dirScan.dirScanName')"
      ><DirScan :projectList="projectList" :taskList="taskList"
    /></ElTabPane>
    <ElTabPane :label="t('vulnerability.vulnerabilityName')">
      <vul :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane :label="t('PageMonitoring.pageMonitoringName')">
      <PageMonitoring :projectList="projectList" :taskList="taskList" />
    </ElTabPane>
    <ElTabPane label="Map" name="map" />
  </ElTabs>
</template>

<script setup lang="ts">
import { useI18n } from '@/hooks/web/useI18n'
import { ElButton, ElIcon, ElTableV2, ElTag, ElTooltip } from 'element-plus'
import type { Column } from 'element-plus'
import { TableV2FixedDir } from 'element-plus'
import { ref } from 'vue'
const { t } = useI18n()

const props = defineProps<{
  closeDialog: () => void
  getProgressInfoID: string
  getProgressInfotype: string
  getProgressInforunnerid: string
}>()

let id = 0
const dataGenerator = () => ({
  id: `random-id-${++id}`,
  name: 'Tom',
  date: '2020-10-1'
})

const columns: Column<any>[] = [
  {
    key: 'target',
    title: t('scanTemplate.TargetHandler'),
    dataKey: 'target',
    width: 150
  },
  {
    key: 'TargetHandler',
    title: t('scanTemplate.TargetHandler'),
    dataKey: 'TargetHandler',
    width: 100
  },
  {
    key: 'SubdomainScan',
    title: t('scanTemplate.SubdomainScan'),
    dataKey: 'name',
    width: 100,
    align: 'center'
  },
  {
    key: 'SubdomainSecurity',
    title: t('scanTemplate.SubdomainSecurity'),
    width: 100
  },
  {
    key: 'PortScanPreparation',
    title: t('scanTemplate.PortScanPreparation'),
    width: 100
  },
  {
    key: 'PortScan',
    title: t('scanTemplate.PortScan'),
    width: 100
  },
  {
    key: 'PortFingerprint',
    title: t('scanTemplate.PortFingerprint'),
    width: 100
  },
  {
    key: 'AssetMapping',
    title: t('scanTemplate.AssetMapping'),
    width: 100
  },
  {
    key: 'AssetHandle',
    title: t('scanTemplate.AssetHandle'),
    width: 100
  },
  {
    key: 'URLScan',
    title: t('scanTemplate.URLScan'),
    width: 100
  },
  {
    key: 'WebCrawler',
    title: t('scanTemplate.WebCrawler'),
    width: 100
  },
  {
    key: 'URLSecurity',
    title: t('scanTemplate.URLSecurity'),
    width: 100
  },
  {
    key: 'DirScan',
    title: t('scanTemplate.DirScan'),
    width: 100
  },
  {
    key: 'VulnerabilityScan',
    title: t('scanTemplate.VulnerabilityScan'),
    width: 100
  },
  {
    key: 'All',
    title: 'All',
    width: 100
  }
]
columns[0].fixed = true
const data = ref(Array.from({ length: 200 }).map(dataGenerator))
</script>
<template>
  <ElTableV2 :columns="columns" :data="data" :width="900" :height="700" fixed />
</template>

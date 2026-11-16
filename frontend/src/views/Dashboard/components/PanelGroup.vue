<script setup lang="ts">
import { ElRow, ElCol, ElCard, ElSkeleton } from 'element-plus'
import { CountTo } from '@/components/CountTo'
import { useDesign } from '@/hooks/web/useDesign'
import { useI18n } from '@/hooks/web/useI18n'
import { computed, ref, reactive } from 'vue'
import { getAssetStatisticsApi } from '@/api/dashboard/analysis'
import type { DashboardTotalTypes } from '@/api/dashboard/analysis/types'

const { t } = useI18n()

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('panel')

const loading = ref(true)

let totalState = reactive<DashboardTotalTypes>({
  asetCount: 0,
  subdomainCount: 0,
  sensitiveCount: 0,
  urlCount: 0,
  vulnerabilityCount: 0
})

const getCount = async () => {
  try {
    const res = await getAssetStatisticsApi()
    // 统计接口失败不应该把首页拖成白屏，保持 0 即可
    if (!res?.data) return
    totalState.asetCount = res.data.assetCount
    totalState.subdomainCount = res.data.subdomainCount
    totalState.sensitiveCount = res.data.sensitiveCount
    totalState.urlCount = res.data.urlCount
    totalState.vulnerabilityCount = res.data.vulnerabilityCount
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

getCount()

// 五个指标卡：颜色只做“识别”，不做装饰
const cards = computed(() => [
  {
    key: 'asset',
    icon: 'zondicons:network',
    label: t('dashboard.totalAssets'),
    value: totalState.asetCount || 0,
    tone: '#0ea5a4',
    soft: 'rgb(14 165 164 / 12%)'
  },
  {
    key: 'subdomain',
    icon: 'formkit:url',
    label: t('dashboard.subDomain'),
    value: totalState.subdomainCount || 0,
    tone: '#2f6feb',
    soft: 'rgb(47 111 235 / 12%)'
  },
  {
    key: 'sensitive',
    icon: 'lucide:info',
    label: t('dashboard.informationLeakage'),
    value: totalState.sensitiveCount || 0,
    tone: '#e11d48',
    soft: 'rgb(225 29 72 / 12%)'
  },
  {
    key: 'url',
    icon: 'carbon:url',
    label: t('dashboard.URL'),
    value: totalState.urlCount || 0,
    tone: '#6366f1',
    soft: 'rgb(99 102 241 / 12%)'
  },
  {
    key: 'vul',
    icon: 'ant-design:bug-filled',
    label: t('vulnerability.vulnerabilityName'),
    value: totalState.vulnerabilityCount || 0,
    tone: '#d97706',
    soft: 'rgb(217 119 6 / 12%)'
  }
])
</script>

<template>
  <ElRow :class="prefixCls" :gutter="16" justify="space-between">
    <ElCol v-for="card in cards" :key="card.key" :xs="24" :sm="12" :md="8" :lg="card.key === 'vul' ? 4 : 5">
      <ElCard shadow="never" :class="`${prefixCls}__card`" class="mb-16px">
        <ElSkeleton :loading="loading" animated :rows="2">
          <template #default>
            <div
              :class="`${prefixCls}__item`"
              :style="{ '--tone': card.tone, '--tone-soft': card.soft }"
            >
              <div :class="`${prefixCls}__icon`">
                <Icon :icon="card.icon" :size="22" />
              </div>
              <div :class="`${prefixCls}__meta`">
                <div :class="`${prefixCls}__label`">{{ card.label }}</div>
                <CountTo
                  :class="`${prefixCls}__value`"
                  :start-val="0"
                  :end-val="card.value"
                  :duration="1600"
                />
              </div>
            </div>
          </template>
        </ElSkeleton>
      </ElCard>
    </ElCol>
  </ElRow>
</template>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-panel';

.@{prefix-cls} {
  &__card {
    border: 1px solid var(--app-border);
    border-radius: var(--app-radius-md);
    transition:
      border-color 0.18s ease,
      transform 0.18s ease,
      box-shadow 0.18s ease;

    &:hover {
      border-color: var(--tone, var(--app-border-strong));
      box-shadow: var(--app-shadow-sm);
      transform: translateY(-1px);
    }

    :deep(.el-card__body) {
      padding: 16px 18px;
    }
  }

  &__item {
    display: flex;
    align-items: center;
    gap: 14px;
    min-height: 52px;
  }

  &__icon {
    display: flex;
    flex: none;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    color: var(--tone);
    background: var(--tone-soft);
    border-radius: 12px;
  }

  &__meta {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  &__label {
    overflow: hidden;
    font-size: 13px;
    color: var(--el-text-color-secondary);
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__value {
    font-size: 26px;
    font-weight: 600;
    font-variant-numeric: tabular-nums;
    line-height: 1.15;
    letter-spacing: -0.02em;
    color: var(--el-text-color-primary);
  }
}
</style>

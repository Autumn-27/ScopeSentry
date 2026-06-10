<script setup lang="ts">
import { ContentDetailWrap } from '@/components/ContentDetailWrap'
import { h, reactive, ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useRouter, useRoute } from 'vue-router'
import { DescriptionsSchema } from '@/components/Descriptions'
import { Descriptions } from '@/components/Descriptions'
import { getAssetDetailApi } from '@/api/asset'
import { ElCol, ElRow, ElScrollbar, ElTag, ElText } from 'element-plus'
import { Icon } from '@/components/Icon'

const { push, go } = useRouter()

const { t } = useI18n()
const { query } = useRoute()
const schema = reactive<DescriptionsSchema[]>([
  {
    field: 'host',
    label: t('asset.domain'),
    slots: {
      default: (data) => {
        const host = data.host
        if (host == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, host)
      }
    }
  },
  {
    field: 'IP',
    label: t('asset.IP'),
    slots: {
      default: (data) => {
        const IP = data.IP
        if (IP == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, IP)
      }
    }
  },
  {
    field: 'URL',
    label: 'URL',
    slots: {
      default: (data) => {
        const URL = data.URL
        if (URL == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, URL)
      }
    }
  },
  {
    field: 'port',
    label: t('asset.port'),
    slots: {
      default: (data) => {
        const statusValue = data.port
        if (statusValue == '') {
          return h('div', '-')
        }
        return h('div', [h(ElTag, statusValue)])
      }
    }
  },
  {
    field: 'service',
    label: t('asset.service'),
    slots: {
      default: (data) => {
        const statusValue = data.service
        if (statusValue == '') {
          return h('div', '-')
        }
        return h('div', [h(ElTag, { type: 'info', effect: 'light', round: true }, statusValue)])
      }
    }
  },
  {
    field: 'title',
    label: t('asset.title'),
    slots: {
      default: (data) => {
        const title = data.title
        if (title == '') {
          return h('div', '-')
        }
        return h('div', title)
      }
    }
  },
  {
    field: 'status',
    label: t('asset.status'),
    slots: {
      default: (data) => {
        const statusValue = data.status
        if (statusValue == '') {
          return h('div', '-')
        }
        const getColor = (value) => {
          return value < 300 ? '#2eb98a' : '#ff5252'
        }
        const color = getColor(statusValue)
        return h('div', [
          h(ElRow, { gutter: 2 }, [
            h(ElCol, { span: 0.99999 }, [
              h(Icon, {
                icon: 'clarity:circle-solid',
                color: color,
                size: 6,
                style: { transform: 'translateY(-35%)' }
              })
            ]),
            h(ElCol, { span: 2 }, [h(ElText, statusValue)])
          ])
        ])
      }
    }
  },
  {
    field: 'FaviconHash',
    label: 'Favicon Hash',
    slots: {
      default: (data) => {
        const FaviconHash = data.FaviconHash
        if (FaviconHash == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, FaviconHash)
      }
    }
  },
  {
    field: 'jarm',
    label: 'Jarm',
    slots: {
      default: (data) => {
        const jarm = data.jarm
        if (jarm == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, jarm)
      }
    }
  },
  {
    field: 'time',
    label: t('asset.time'),
    slots: {
      default: (data) => {
        const time = data.time
        if (time == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, time)
      }
    }
  },
  {
    field: 'products',
    label: t('asset.products'),
    span: 12,
    slots: {
      default: (data) => {
        const products: any[] = data.products
        if (!Array.isArray(products) || products.length === 0) {
          return h('div', '-')
        }
        const rows = []
        for (let i = 0; i < products.length; i += 6) {
          const sliceStart = i
          const sliceEnd = i + 6
          const row: any[] = products.slice(sliceStart, sliceEnd) as any[]
          rows.push(row)
        }

        const tags = rows.map((row, rowIndex) => {
          const rowTags = row.map((product, colIndex) => {
            return h(
              ElCol,
              { span: 3 },
              h(ElTag, { key: rowIndex * 6 + colIndex, type: 'success' }, product)
            )
          })
          return h(ElRow, { gutter: 1 }, rowTags)
        })

        return h('div', tags)
      }
    }
  },
  {
    field: 'project',
    label: t('project.projectName'),
    span: 12,
    slots: {
      default: (data) => {
        const statusValue = data.project
        if (statusValue == '') {
          return h('div', '-')
        }
        return h('div', [h(ElTag, statusValue)])
      }
    }
  },
  {
    field: 'TLSData',
    label: 'TLS',
    span: 24,
    slots: {
      default: (data) => {
        const TLSData = data.TLSData
        if (TLSData == '') {
          return h('div', '-')
        }
        return h(
          ElScrollbar,
          { maxHeight: '100px' },
          {
            default: () => h('div', { style: { whiteSpace: 'pre-line' } }, TLSData)
          }
        )
      }
    }
  },
  {
    field: 'hash',
    label: 'Hash',
    span: 24,
    slots: {
      default: (data) => {
        const hashValue = data.hash
        if (hashValue == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, hashValue)
      }
    }
  },
  {
    field: 'banner',
    label: 'Banner',
    span: 24,
    slots: {
      default: (data) => {
        const banner = data.banner
        if (banner == '') {
          return h('div', '-')
        }
        return h('div', { style: { whiteSpace: 'pre-line' } }, banner)
      }
    }
  },
  {
    field: 'ResponseBody',
    label: t('asset.responseBody'),
    span: 24,
    slots: {
      default: (data) => {
        const ResponseBody = data.ResponseBody
        if (ResponseBody == '') {
          return h('div', '-')
        }
        return h(
          ElScrollbar,
          { maxHeight: '100px' },
          {
            default: () => h('div', { style: { whiteSpace: 'pre-line' } }, ResponseBody)
          }
        )
      }
    }
  }
])
const descriptionsDoading = ref(true)
let assetData = reactive({})
const getTableDet = async () => {
  const res = await getAssetDetailApi(query.id as string)
  if (res) {
    assetData = Object.assign(assetData, res?.data || {})
    descriptionsDoading.value = false
  }
}
getTableDet()
</script>

<template>
  <ContentDetailWrap
    :title="t('exampleDemo.detail')"
    @back="push('/asset-information/index')"
    v-loading="descriptionsDoading"
  >
    <template #header>
      <BaseButton @click="go(-1)">
        {{ t('common.back') }}
      </BaseButton>
    </template>
    <Descriptions
      :title="t('asset.assetDetail')"
      :schema="schema"
      :data="assetData"
      :collapse="false"
    />
  </ContentDetailWrap>
</template>

<style lang="less">
.el-row {
  margin-bottom: 10px;
}
.el-row:last-child {
  margin-bottom: 0;
}
.el-col {
  border-radius: 4px;
}
</style>

<script setup lang="ts">
import {
  ElCard,
  ElTabPane,
  ElTabs,
  ElTag,
  ElTimeline,
  ElTimelineItem,
  ElIcon,
  ElSpace,
  ElText,
  ElRow,
  ElCol
} from 'element-plus'
import { Codemirror } from 'vue-codemirror'
import { ArrowDown, ArrowRight } from '@element-plus/icons-vue'
import { javascript } from '@codemirror/lang-javascript'
import { oneDark } from '@codemirror/theme-one-dark'
import { getAssetChangeLogApi, getAssetDetailApi } from '@/api/asset'
import { ref } from 'vue'
import { AssetChangeLog } from '@/api/asset/types'
import { createImageViewer } from '@/components/ImageViewer'
const extensions = [javascript(), oneDark]
const props = defineProps<{
  id: string
  host: string
  ip: string
  port: number
}>()

const detailJson = ref('')
const getAssetDetail = async () => {
  const res = await getAssetDetailApi(props.id)
  if (res.code == 200) {
    // 格式化 JSON，使用 2 个空格缩进
    detailJson.value = JSON.stringify(res.data, null, 2)
  }
}
getAssetDetail()

const assetChangeLog = ref<AssetChangeLog[]>([])
// 获取资产变更日志
const getAssetChangeLog = async () => {
  const res = await getAssetChangeLogApi(props.id)
  if (res.code === 200) {
    // 在获取到数据后，添加 isExpanded 属性
    assetChangeLog.value = res.data.map((log) => ({
      ...log, // 保留原始的 log 数据
      isExpanded: false // 添加 isExpanded 属性，默认是 false
    }))
  }
}
getAssetChangeLog()
const handleImageClick = (screenshot: string) => {
  createImageViewer({
    urlList: [screenshot],
    zIndex: 999999
  })
}
</script>

<template>
  <ElTabs type="border-card" tab-position="left">
    <ElTabPane label="原始数据">
      <Codemirror
        v-model="detailJson"
        :extensions="extensions"
        :autofocus="true"
        :indent-with-tab="true"
        :tab-size="2"
        :style="{ height: '550px', width: '100%' }"
      />
    </ElTabPane>
    <ElTabPane label="资产变更">
      <ElRow style="margin-bottom: 20px">
        <ElCol :offset="2">
          <ElSpace>
            <ElText>{{ props.host }}</ElText>
            <ElText type="info" size="small">{{ props.ip }}</ElText>
            <ElTag type="success">{{ props.port }}</ElTag>
          </ElSpace>
        </ElCol>
      </ElRow>
      <ElRow>
        <ElCol>
          <ElTimeline>
            <!-- 动态渲染每个变更 -->
            <ElTimelineItem
              v-for="(log, index) in assetChangeLog"
              :key="index"
              :timestamp="log.timestamp"
              :type="index % 2 === 0 ? 'primary' : 'danger'"
              placement="top"
            >
              <ElCard>
                <div
                  style="display: flex; flex-wrap: wrap; gap: 10px; flex-grow: 1"
                  @click="log.isExpanded = !log.isExpanded"
                >
                  <!-- 添加箭头图标，控制展开收起 -->
                  <ElIcon :style="{ marginLeft: '10px' }">
                    <ArrowRight v-if="!log.isExpanded" />
                    <ArrowDown v-if="log.isExpanded" />
                  </ElIcon>
                  <ElTag
                    v-for="(change, index) in log.change"
                    :key="index"
                    type="danger"
                    :style="{ marginBottom: '10px' }"
                  >
                    {{ change.fieldname }}
                  </ElTag>
                </div>

                <div class="p-6" v-if="log.isExpanded">
                  <div class="grid grid-cols-2 gap-6">
                    <!-- 旧值部分 -->
                    <div class="space-y-2">
                      <div class="el-card border-gray-200" style="border-radius: 12px">
                        <div
                          class="px-4 py-2 bg-gray-100 border-b border-gray-200 font-medium text-sm"
                        >
                          旧值
                        </div>
                        <div class="p-4 text-sm whitespace-pre-wrap">
                          <!-- 遍历并输出 fieldname: old -->
                          <div v-for="(change, index) in log.change" :key="'old-' + index">
                            <strong>{{ change.fieldname }}:</strong>
                            <template v-if="change.fieldname === 'Screenshot'">
                              <img
                                :src="change.old"
                                alt="screenshot"
                                style="width: 100%; height: auto; max-height: 250px"
                                @click="handleImageClick(change.old)"
                              />
                            </template>
                            <template v-else>
                              {{ change.old }}
                            </template>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- 新值部分 -->
                    <div class="space-y-2">
                      <div class="el-card border-gray-200" style="border-radius: 12px">
                        <div
                          class="px-4 py-2 bg-blue-100 border-b border-blue-200 font-medium text-sm"
                        >
                          新值
                        </div>
                        <div class="p-4 text-sm whitespace-pre-wrap">
                          <!-- 遍历并输出 fieldname: new -->
                          <div v-for="(change, index) in log.change" :key="'new-' + index">
                            <strong>{{ change.fieldname }}:</strong>
                            <template v-if="change.fieldname === 'Screenshot'">
                              <img
                                :src="change.new"
                                alt="screenshot"
                                style="width: 100%; height: auto; max-height: 250px"
                                @click="handleImageClick(change.new)"
                              />
                            </template>
                            <template v-else>
                              {{ change.new }}
                            </template>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </ElCard>
            </ElTimelineItem>
          </ElTimeline>
        </ElCol>
      </ElRow>
    </ElTabPane>
  </ElTabs>
</template>
<style lang="less"></style>

/**
 * App 根组件：初始化主题（日间/夜间/跟随系统）
 */
<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store/modules/app'
import { ConfigGlobal } from '@/components/ConfigGlobal'
import { useDesign } from '@/hooks/web/useDesign'
import { setCssVar } from './utils'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('app')

const appStore = useAppStore()

const currentSize = computed(() => appStore.getCurrentSize)

const greyMode = computed(() => appStore.getGreyMode)

// 初始化主题：读取持久化模式（旧版本 isDark 自动迁移），默认跟随系统，
// 并在“跟随系统”模式下监听 prefers-color-scheme 的实时变化
appStore.initTheme()

// 初始化标签页高度
const tagsViewEnabled = appStore.getTagsView
setCssVar('--tags-view-height', tagsViewEnabled ? '42px' : '0px')
</script>

<template>
  <ConfigGlobal :size="currentSize">
    <RouterView :class="greyMode ? `${prefixCls}-grey-mode` : ''" />
  </ConfigGlobal>
</template>

<style lang="less">
@prefix-cls: ~'@{namespace}-app';

.size {
  width: 100%;
  height: 100%;
}

html,
body {
  padding: 0 !important;
  margin: 0;
  overflow: hidden;
  .size;

  #app {
    .size;
  }
}

.@{prefix-cls}-grey-mode {
  filter: grayscale(100%);
}
</style>

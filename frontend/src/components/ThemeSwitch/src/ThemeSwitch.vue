/**
 * 主题切换：日间 / 夜间 / 跟随系统（设备）
 * 只负责 UI，状态与持久化在 app store（setThemeMode）
 */
<script setup lang="ts">
import { computed } from 'vue'
import { ElDropdown, ElDropdownItem, ElDropdownMenu, ElTooltip } from 'element-plus'
import { useAppStore } from '@/store/modules/app'
import { useI18n } from '@/hooks/web/useI18n'
import { useDesign } from '@/hooks/web/useDesign'
import { propTypes } from '@/utils/propTypes'
import type { ThemeMode } from '@/constants/theme'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('theme-switch')

defineProps({
  color: propTypes.string.def('')
})

const appStore = useAppStore()

const { t } = useI18n()

const mode = computed<ThemeMode>(() => appStore.getThemeMode)

const isDark = computed(() => appStore.getIsDark)

const options = computed(() => [
  {
    value: 'light' as ThemeMode,
    icon: 'emojione-monotone:sun',
    label: t('setting.themeLight')
  },
  {
    value: 'dark' as ThemeMode,
    icon: 'emojione-monotone:crescent-moon',
    label: t('setting.themeDark')
  },
  {
    value: 'auto' as ThemeMode,
    icon: 'ant-design:desktop-outlined',
    label: t('setting.themeAuto')
  }
])

const triggerIcon = computed(() => {
  if (mode.value === 'auto') return 'ant-design:desktop-outlined'
  return isDark.value ? 'emojione-monotone:crescent-moon' : 'emojione-monotone:sun'
})

const handleCommand = (command: ThemeMode) => {
  appStore.setThemeMode(command)
}
</script>

<template>
  <ElTooltip :content="t('setting.themeMode')" placement="bottom">
    <ElDropdown :class="prefixCls" trigger="click" @command="handleCommand">
      <span :class="`${prefixCls}__trigger`" class="cursor-pointer flex items-center">
        <Icon :icon="triggerIcon" :color="color || undefined" :size="18" />
      </span>
      <template #dropdown>
        <ElDropdownMenu>
          <ElDropdownItem v-for="item in options" :key="item.value" :command="item.value">
            <span :class="`${prefixCls}__item`">
              <Icon :icon="item.icon" :size="16" />
              <span :class="`${prefixCls}__label`">{{ item.label }}</span>
              <Icon
                v-if="mode === item.value"
                icon="ant-design:check-outlined"
                :size="14"
                color="var(--el-color-primary)"
              />
            </span>
          </ElDropdownItem>
        </ElDropdownMenu>
      </template>
    </ElDropdown>
  </ElTooltip>
</template>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-theme-switch';

.@{prefix-cls} {
  &__trigger {
    padding: 0 6px;
  }

  &__item {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 120px;
  }

  &__label {
    flex: 1;
  }
}
</style>

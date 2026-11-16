<script setup lang="ts">
import { computed, ref, unref } from 'vue'
import { ElDropdown, ElDropdownMenu, ElDropdownItem, ElMessageBox } from 'element-plus'
import { useLocaleStore } from '@/store/modules/locale'
import { useI18n } from '@/hooks/web/useI18n'
import { useLocale } from '@/hooks/web/useLocale'
import { propTypes } from '@/utils/propTypes'
import { useDesign } from '@/hooks/web/useDesign'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('locale-dropdown')

defineProps({
  color: propTypes.string.def('')
})

const localeStore = useLocaleStore()

const langMap = computed(() => localeStore.getLocaleMap)

const currentLang = computed(() => localeStore.getCurrentLocale)

const { t } = useI18n()
const { changeLocale } = useLocale()
const isChanging = ref(false)

const setLang = async (lang: LocaleType) => {
  if (isChanging.value || lang === unref(currentLang).lang) return
  isChanging.value = true
  try {
    // Setup-time translated schemas currently require recreating the view.
    // Ask before discarding input; do not cancel or replay requests here.
    try {
      await ElMessageBox.confirm(t('common.changeLanguageMessage'), t('common.reminder'), {
        confirmButtonText: t('common.ok'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
    } catch (action) {
      if (action === 'cancel' || action === 'close') return
      throw action
    }
    await changeLocale(lang)
  } finally {
    isChanging.value = false
  }
}
</script>

<template>
  <ElDropdown :class="prefixCls" :disabled="isChanging" trigger="click" @command="setLang">
    <Icon
      :size="18"
      icon="ion:language-sharp"
      class="cursor-pointer !p-0"
      :class="$attrs.class"
      :color="color"
    />
    <template #dropdown>
      <ElDropdownMenu>
        <ElDropdownItem v-for="item in langMap" :key="item.lang" :command="item.lang">
          {{ item.name }}
        </ElDropdownItem>
      </ElDropdownMenu>
    </template>
  </ElDropdown>
</template>

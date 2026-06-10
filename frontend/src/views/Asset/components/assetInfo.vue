<script setup lang="ts">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { Search } from '@/components/Search'
import { reactive, ref } from 'vue'
import { FormSchema } from '@/components/Form'
import { useSearch } from '@/hooks/web/useSearch'
import { onMounted } from 'vue'
import { ElPagination } from 'element-plus'
import { getAssetApi } from '@/api/Asset'
import { ElRow, ElCol, ElCard } from 'element-plus'
import { ElTabs } from 'element-plus'
import { ElTabPane, ElInput, ElLink, ElTag, ElText } from 'element-plus'
import type { TagProps } from 'element-plus'

const { t } = useI18n()
const loading = ref(true)
const { searchRegister } = useSearch()
const restaurants = ref<Recordable[]>([])
const querySearch = (queryString: string, cb: Fn) => {
  const results = queryString
    ? restaurants.value.filter(createFilter(queryString))
    : restaurants.value
  // call callback function to return suggestions
  cb(results)
}

const handleSelect = (item: Recordable) => {
  console.log(item)
}

const schema = reactive<FormSchema[]>([
  {
    field: 'search',
    label: t('form.input'),
    component: 'Autocomplete',
    componentProps: {
      fetchSuggestions: querySearch,
      on: {
        select: handleSelect
      }
    },
    formItemProps: {
      size: 'large',
      style: { width: '100%' }
    }
  }
])
const createFilter = (queryString: string) => {
  return (restaurant: Recordable) => {
    return restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
  }
}
const loadAll = () => {
  return [
    { value: 'body' },
    { value: 'header' },
    { value: 'title' },
    { value: 'icon_hash' },
    { value: 'ip' },
    { value: 'host' },
    { value: 'and' },
    { value: 'or' },
    { value: '=' }
  ]
}

onMounted(() => {
  restaurants.value = loadAll()
})

const isGrid = ref(true)
const layout = ref('inline')

const buttonPosition = ref('left')

const searchParams = ref('')
const handleSearch = (data: any) => {
  searchParams.value = data.search
  getAssetData()
}

const searchLoading = ref(false)

const currentPage = ref(1)
const pageSize = ref(10)
const small = ref(false)
const background = ref(false)
const disabled = ref(false)
const handleSizeChange = () => {
  getAssetData()
}
const handleCurrentChange = () => {
  getAssetData()
}
const getAssetData = async () => {
  console.log(pageSize.value)
  console.log(currentPage.value)
  console.log(searchParams.value)
  try {
    const res = await getAssetApi(searchParams.value, currentPage.value, pageSize.value)
    console.log(res)
  } catch (error) {
    console.error('Error fetching node data:', error)
  } finally {
    // 不论请求成功或失败，都会执行的代码块
    loading.value = false
  }
}
getAssetData()
const textarea = ref(
  '111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111111111111111xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx11111xxxxxxxxxxxxxxxxxxx11111'
)

type Item = { type: TagProps['type']; label: string }

const items = ref<Array<Item>>([
  { type: '', label: 'Tag 1' },
  { type: 'success', label: 'Tag 2' },
  { type: 'info', label: 'Tag 3' },
  { type: 'danger', label: 'Tag 4' },
  { type: 'warning', label: 'Tag 5' }
])
</script>

<template>
  <ContentWrap style="height: 80px">
    <Search
      :schema="schema"
      :is-col="isGrid"
      :layout="layout"
      :show-reset="false"
      :expand="false"
      :button-position="buttonPosition"
      :search-loading="searchLoading"
      @search="handleSearch"
      @reset="handleSearch"
      @register="searchRegister"
    />
  </ContentWrap>
  <ElRow>
    <ElCol :span="6">
      <div class="grid-content ep-bg-purple">xxxxxxxxxxxxx</div>
    </ElCol>
  </ElRow>
  <ElRow :gutter="20" justify="space-between">
    <ElCol :span="16" :offset="4">
      <ElCard shadow="never" class="mb-25px" style="background-color: var(--el-fill-color-light)">
        <template #header>
          <div class="header-container" style="height: 10%">
            <span class="header-content">
              <ElLink
                :underline="false"
                href="https://xxxxxxxx.top"
                style="color: var(--el-color-primary)"
              >
                https://xxxxxx.top
              </ElLink>
            </span>
          </div>
        </template>
        <ElRow :gutter="20">
          <ElCol :span="10">
            <ElRow>
              <ElCol :span="24">
                <ElText class="mx-1">XXX XXX</ElText>
              </ElCol>
              <ElCol :span="24">
                <ElLink
                  :underline="false"
                  href="https://xxxxxxxx.top"
                  style="color: var(--el-color-primary); font-size: large"
                >
                  127.0.0.1
                </ElLink>
              </ElCol>
              <ElCol :span="24">
                <ElText class="mx-1">2023-12-24</ElText>
              </ElCol>
            </ElRow>
          </ElCol>
          <ElCol :span="13">
            <ElTabs type="border-card">
              <ElTabPane :label="t('asset.banner')">
                <ElInput
                  v-model="textarea"
                  type="textarea"
                  :rows="5"
                  readonly
                  autosize
                  resize="none"
                />
              </ElTabPane>
              <ElTabPane :label="t('asset.products')">
                <div class="tag-group my-2 flex flex-wrap gap-1 items-center">
                  <span class="tag-group__title m-1 line-height-2">Dark</span>
                  <ElTag
                    v-for="item in items"
                    :key="item.label"
                    :type="item.type"
                    class="mx-1"
                    effect="dark"
                  >
                    {{ item.label }}
                  </ElTag>
                </div>
              </ElTabPane>
            </ElTabs>
          </ElCol>
        </ElRow>
      </ElCard>
    </ElCol>
  </ElRow>
  <div class="demo-pagination-block">
    <ElPagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[10, 100, 200, 500, 1000]"
      :small="small"
      :disabled="disabled"
      :background="background"
      layout="total, sizes, prev, pager, next, jumper"
      :total="400"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<style lang="less" scoped>
.el-button {
  margin-top: 10px;
}
.el-textarea {
  --el-input-focus-border-color: #a8abb2;
}
.el-tabs__item {
  flex: 1;
  text-align: center;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>

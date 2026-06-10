<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, h, nextTick, Ref } from 'vue'
import {
  ElButton,
  ElCol,
  ElInput,
  ElRow,
  ElText,
  ElUpload,
  ElTooltip,
  ElMessage,
  UploadProps,
  UploadRawFile,
  UploadInstance,
  ElTag,
  InputInstance
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { useTable } from '@/hooks/web/useTable'
import { Icon } from '@/components/Icon'
import { useIcon } from '@/hooks/web/useIcon'
import { BaseButton } from '@/components/Button'
import { getPocDataApi, getPocContentApi, deletePocDataApi, getPocDetailApi } from '@/api/poc'
import Detail from './components/Detail.vue'
import { Dialog } from '@/components/Dialog'
import { useUserStore } from '@/store/modules/user'
import { RowState } from '@/api/asset/types'
import { addTagApi, deleteTagApi } from '@/api/asset'
const searchicon = useIcon({ icon: 'iconoir:search' })
const { t } = useI18n()
const dialogVisible = ref(false)
const search = ref('')
const handleSearch = () => {
  getList()
}
const rowStateMap = reactive<Record<string, RowState>>({})
const nodeColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'name',
    label: t('poc.pocName'),
    minWidth: 70
  },
  {
    field: 'level',
    label: t('poc.level'),
    minWidth: 50,
    columnKey: 'level',
    formatter: (_: Recordable, __: TableColumn, levelValue: string) => {
      if (levelValue == null) {
        return <div></div>
      }
      let color = ''
      let flag = ''
      if (levelValue === 'critical') {
        color = 'red'
        flag = t('poc.critical')
      } else if (levelValue === 'high') {
        color = 'orange'
        flag = t('poc.high')
      } else if (levelValue === 'medium') {
        color = 'yellow'
        flag = t('poc.medium')
      } else if (levelValue === 'low') {
        color = 'blue'
        flag = t('poc.low')
      } else if (levelValue === 'info') {
        color = 'green'
        flag = t('poc.info')
      } else if (levelValue === 'unknown') {
        color = 'gray'
        flag = t('poc.unknown')
      }
      return (
        <ElRow gutter={20} style="width: 80%">
          <ElCol span={1}>
            <Icon icon="clarity:circle-solid" color={color} size={10} />
          </ElCol>
          <ElCol span={5}>
            <ElText type="info">{flag}</ElText>
          </ElCol>
        </ElRow>
      )
    },
    filters: [
      { text: t('poc.critical'), value: 'critical' },
      { text: t('poc.high'), value: 'high' },
      { text: t('poc.medium'), value: 'medium' },
      { text: t('poc.low'), value: 'low' },
      { text: t('poc.info'), value: 'info' },
      { text: t('poc.unknown'), value: 'unknown' }
    ]
  },
  {
    field: 'tags',
    label: 'TAG',
    fit: 'true',
    formatter: (row: Recordable, __: TableColumn, tags: string[]) => {
      if (tags.length != 0) {
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
            {tags.map((product) => (
              <div
                key={product}
                onClick={() => changeTags('app', product)}
                style={{ cursor: 'pointer' }}
              >
                <ElTag type={'success'}>{product}</ElTag>
              </div>
            ))}
          </div>
        )
        // if (ProductsValue.length > 1) {
        //   let contentTool = ''
        //   if (Array.isArray(ProductsValue)) {
        //     // It's an array, you can use forEach
        //     ProductsValue.forEach((item, _) => {
        //       contentTool += `<div>${item}</div>`
        //     })
        //   } else {
        //     console.error('ProductsValue is not an array')
        //   }
        //   return (
        //     <div class="flex">
        //       <ElTag type="success" effect="light" round>
        //         {ProductsValue[0]}
        //       </ElTag>
        //       <ElTooltip
        //         class="box-item"
        //         effect="dark"
        //         placement="top-start"
        //         content={contentTool}
        //         popper-class="tagtooltip"
        //         rawContent
        //       >
        //         <ElTag type="info" effect="plain" round style={'left:3px; position:relative'}>
        //           {t('asset.total')} {ProductsValue.length} {t('asset.p')}
        //         </ElTag>
        //       </ElTooltip>
        //     </div>
        //   )
        // } else {
        //   return (
        //     <div class="flex">
        //       <ElTag type="success" effect="light">
        //         {ProductsValue[0]}
        //       </ElTag>
        //     </div>
        //   )
        // }
      }
    }
  },
  {
    field: 'time',
    label: t('node.createTime'),
    minWidth: 50
  },

  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 30,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
          <BaseButton type="primary" onClick={() => edit(row)}>
            {t('common.edit')}
          </BaseButton>
          <BaseButton type="danger" onClick={() => del(row)}>
            {t('common.delete')}
          </BaseButton>
        </>
      )
    }
  }
])
const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const { currentPage, pageSize } = tableState
    const res = await getPocDataApi(search.value, currentPage.value, pageSize.value, filter)
    return {
      list: res.data.list,
      total: res.data.total
    }
  }
})
const { loading, dataList, total, currentPage, pageSize } = tableState
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}

const changeTags = (type, value) => {
  const key = `${type}=${value}`
  console.log(key)
  // dynamicTags.value = [...dynamicTags.value, key]
}

let pocForm = reactive({
  id: '',
  name: '',
  level: '',
  content: '',
  tags: []
})
const addPoc = async () => {
  pocForm.id = ''
  pocForm.name = ''
  pocForm.level = ''
  pocForm.content = ''
  pocForm.tags = []
  dialogVisible.value = true
}
const edit = async (data) => {
  pocForm.id = data.id
  pocForm.name = data.name
  pocForm.level = data.level
  pocForm.tags = data.tags
  const res = await getPocDetailApi(pocForm.id)
  pocForm.content = res.data.content
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}
const delLoading = ref(false)
const del = async (data) => {
  delLoading.value = true
  try {
    const res = await deletePocDataApi([data.id])
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const ids = ref<string[]>([])
const delSelect = async () => {
  const elTableExpose = await getElTableExpose()
  const selectedRows = elTableExpose?.getSelectionRows() || []
  ids.value = selectedRows.map((row) => row.id)
  delLoading.value = true
  try {
    const res = await deletePocDataApi(ids.value)
    console.log('Data deleted successfully:', res)
    delLoading.value = false
    getList()
  } catch (error) {
    console.error('Error deleting data:', error)
    delLoading.value = false
    getList()
  }
}
const confirmDelete = async () => {
  const confirmed = window.confirm('Are you sure you want to delete the selected data?')
  if (confirmed) {
    await delSelect()
  }
}
const userStore = useUserStore()
const uploadHeaders = ref({ Authorization: `${userStore.getToken}` })

const upload = ref<UploadInstance>()
const uploadSuccess = async () => {
  console.log('导入中')
  ElMessage.success('导入中')
}

const handleExceed: UploadProps['onExceed'] = (files) => {
  upload.value!.clearFiles()
  const file = files[0] as UploadRawFile
  upload.value!.handleStart(file)
}

const handleUploadSuccess = (response) => {
  console.log(response)
  if (response.code === 200) {
    ElMessage.success('Upload succes')
  } else {
    ElMessage.error(response.message)
  }
  if (response.code == 505) {
    localStorage.removeItem('plugin_key')
  }
  getList()
  upload.value?.clearFiles()
}
const handleFileChange = (file, fileList) => {
  if (fileList.length > 0) {
    upload.value!.submit()
  }
}
const filter = reactive<{ [key: string]: any }>({})
const filterChange = async (newFilters: any) => {
  Object.assign(filter, newFilters)
  getList()
}
</script>

<template>
  <ContentWrap>
    <ElRow :gutter="20" style="margin-bottom: 15px">
      <ElCol :span="1.5">
        <ElText class="mx-1" style="position: relative; top: 8px">{{ t('poc.pocName') }}:</ElText>
      </ElCol>
      <ElCol :span="5">
        <ElInput v-model="search" :placeholder="t('common.inputText')" style="height: 38px" />
      </ElCol>
      <ElCol :span="5" style="position: relative; left: 16px">
        <ElButton type="primary" :icon="searchicon" style="height: 100%" @click="handleSearch"
          >Search</ElButton
        >
      </ElCol>
    </ElRow>
    <ElRow :gutter="60">
      <ElCol :span="1">
        <div class="mb-10px">
          <ElButton type="primary" @click="addPoc">{{ t('common.new') }}</ElButton>
        </div>
      </ElCol>
      <ElCol :span="1">
        <div class="mb-10px">
          <BaseButton type="danger" :loading="delLoading" @click="confirmDelete">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </ElCol>
      <ElCol :span="3">
        <ElTooltip :content="t('common.uploadMsg')" placement="top">
          <!-- <ElUpload class="upload-demo" action="/api/poc/data/import" :headers="uploadHeaders">
            <ElButton :icon="uploadicon">{{ t('common.import') }}</ElButton>
          </ElUpload> -->
          <ElUpload
            ref="upload"
            class="flex items-center"
            action="/api/poc/data/import"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :limit="1"
            :on-exceed="handleExceed"
            :auto-upload="false"
            @change="handleFileChange"
          >
            <template #trigger>
              <BaseButton>
                <template #icon>
                  <Icon icon="iconoir:upload" />
                </template>
                {{ t('plugin.import') }}
              </BaseButton>
            </template>
          </ElUpload>
        </ElTooltip>
      </ElCol>
    </ElRow>
    <Table
      v-model:pageSize="pageSize"
      v-model:currentPage="currentPage"
      @filter-change="filterChange"
      :columns="nodeColums"
      :data="dataList"
      stripe
      :border="true"
      :loading="loading"
      :resizable="true"
      :pagination="{
        total: total,
        pageSizes: [10, 20, 50, 100, 200, 500, 1000]
      }"
      @register="tableRegister"
      :headerCellStyle="tableHeaderColor"
      :style="{
        fontFamily:
          '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji'
      }"
    />
  </ContentWrap>
  <Dialog
    v-model="dialogVisible"
    :title="pocForm.id ? $t('common.edit') : $t('common.new')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="800"
  >
    <Detail :closeDialog="closeDialog" :pocForm="pocForm" :getList="getList" />
  </Dialog>
</template>

<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, h } from 'vue'
import {
  ElButton,
  ElCol,
  ElRow,
  ElInput,
  ElUpload,
  ElMessage,
  UploadInstance,
  ElForm,
  ElFormItem,
  ElIcon,
  ElTag
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { Dialog } from '@/components/Dialog'
import { Plus } from '@element-plus/icons-vue'
import { useTable } from '@/hooks/web/useTable'
import { BaseButton } from '@/components/Button'
import {
  createDictApi,
  deleteDictApi,
  downloadDictApi,
  getManagetListApi
} from '@/api/DictionaryManagement'
import { useUserStore } from '@/store/modules/user'
const { t } = useI18n()
const userStore = useUserStore()
const uploadHeaders = { Authorization: `${userStore.getToken}` }
const nodeColums = reactive<TableColumn[]>([
  {
    field: 'selection',
    type: 'selection',
    width: '55'
  },
  {
    field: 'id',
    hidden: true
  },
  {
    field: 'name',
    label: t('common.name'),
    minWidth: 40
  },
  {
    field: 'category',
    label: t('common.category'),
    minWidth: 40,
    formatter: (row, __: TableColumn, value: string) => {
      return (
        <ElTag type="success" effect="light" size="large">
          {value}
        </ElTag>
      )
    }
  },
  {
    field: 'size',
    label: t('common.filesize'),
    formatter: (_, __: TableColumn, value: string) => {
      return value + 'MB'
    }
  },
  {
    field: 'upload',
    label: t('common.upload'),
    formatter: (row, __: TableColumn, _: string) => {
      const uploadAction = `/api/dictionary/manage/save?id=${row.id}`
      return h('div', [
        h(
          ElUpload,
          {
            class: 'upload-demo',
            action: uploadAction,
            headers: uploadHeaders,
            onSuccess: (response) => {
              console.log(response.code)
              if (response.code === 200) {
                ElMessage.success('Upload succes')
              } else {
                ElMessage.error('Upload failed')
              }
              getList()
            },
            onError: (err) => {
              // 如果上传失败，触发该回调
              ElMessage.error(`Upload failed: ${err}`)
            }
          },
          () => h(ElButton, { type: 'warning' }, t('common.cover'))
        )
      ])
    }
  },
  {
    field: 'download',
    label: t('common.download'),
    formatter: (row, __: TableColumn, _: string) => {
      return (
        <>
          <BaseButton type="primary" onClick={() => downloadFile(row.id, row.name)}>
            {t('common.download')}
          </BaseButton>
        </>
      )
    }
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 40,
    formatter: (row, __: TableColumn, _: number) => {
      return (
        <>
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
    const res = await getManagetListApi()
    return {
      list: res.data.list
    }
  }
})
const { loading, dataList } = tableState
const { getList, getElTableExpose } = tableMethods
function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}
const dialogVisible = ref(false)
const form = ref({
  name: '',
  category: ''
})
const delLoading = ref(false)
const del = async (data) => {
  const confirmed = window.confirm('Are you sure you want to delete this data?')

  if (!confirmed) {
    return
  }
  delLoading.value = true
  try {
    const res = await deleteDictApi([data.id])
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
    const res = await deleteDictApi(ids.value)
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

const upload = ref<UploadInstance>()

async function customRequest(option) {
  try {
    const file = option.file
    const formData = new FormData()
    formData.append('file', file)
    formData.append('name', form.value.name)
    formData.append('category', form.value.category)
    // 在这里可以处理实际的文件数据
    await createDictApi(formData) // 传递文件和其他参数
    getList() // 假设 getList 是一个用于获取更新后的数据的函数
    dialogVisible.value = false
    ElMessage.success('Upload success')
  } catch (error) {
    ElMessage.error(`Upload failed: ${error}`)
  }
}
async function handleSubmit() {
  upload.value!.submit()
}
const openCreateDialog = async () => {
  form.value.name = ''
  dialogVisible.value = true
}

const downloadFile = async (id, name) => {
  const response = await downloadDictApi(id)
  console.log(response)
  const url = window.URL.createObjectURL(new Blob([response.data]))
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', name) // 设置文件名
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
</script>

<template>
  <ContentWrap>
    <ElRow :gutter="60">
      <ElCol :span="1">
        <div class="mb-10px">
          <ElButton type="primary" @click="openCreateDialog">{{ t('common.new') }}</ElButton>
        </div>
      </ElCol>
      <ElCol :span="1">
        <div class="mb-10px">
          <BaseButton type="danger" :loading="delLoading" @click="confirmDelete">
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </ElCol>
    </ElRow>
    <Table
      :columns="nodeColums"
      :data="dataList"
      stripe
      :border="true"
      :loading="loading"
      :resizable="true"
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
    :title="t('common.new')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="430"
  >
    <ElForm :model="form" label-width="120px" class="upload-form">
      <ElFormItem :label="t('common.name')">
        <el-input v-model="form.name" placeholder="Enter name" />
      </ElFormItem>
      <ElFormItem :label="t('common.category')">
        <el-input v-model="form.category" placeholder="Enter category" />
      </ElFormItem>
      <ElFormItem label="">
        <ElUpload
          class="upload-demo"
          drag
          ref="upload"
          :auto-upload="false"
          :limit="1"
          :httpRequest="customRequest"
        >
          <el-icon class="avatar-uploader-icon"><Plus /></el-icon>
        </ElUpload>
      </ElFormItem>
      <ElFormItem>
        <ElButton type="primary" @click="handleSubmit">Submit</ElButton>
      </ElFormItem>
    </ElForm>
  </Dialog>
</template>
<style scoped>
.avatar-uploader .avatar {
  width: 178px;
  height: 130px;
  display: block;
}
</style>

<style>
.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 130px;
  text-align: center;
}
</style>

<script setup lang="tsx">
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive } from 'vue'
import {
  ElButton,
  ElCol,
  ElRow,
  ElInput,
  ElForm,
  ElFormItem,
  ElAlert,
  ElTag,
  ElMessageBox
} from 'element-plus'
import { Table, TableColumn } from '@/components/Table'
import { Dialog } from '@/components/Dialog'
import { useTable } from '@/hooks/web/useTable'
import { BaseButton } from '@/components/Button'
import { createApiKeyApi, deleteApiKeyApi, getApiKeyListApi } from '@/api/apikey'
import { useClipboard } from '@/hooks/web/useClipboard'
import { formatTime } from '@/utils'

const { t } = useI18n()
const { copy, copied } = useClipboard()

const mcpConfigExample = `{
  "mcpServers": {
    "scopesentry": {
      "url": "http://your-host:8082/mcp",
      "headers": {
        "X-API-Key": "ssk_你的密钥"
      }
    }
  }
}`

const formatDateTime = (value?: string) => {
  if (!value) return '-'
  return formatTime(value, 'yyyy-MM-dd HH:mm:ss')
}

const columns = reactive<TableColumn[]>([
  {
    field: 'name',
    label: t('common.name'),
    minWidth: 120
  },
  {
    field: 'keyPrefix',
    label: t('apiKey.keyPrefix'),
    minWidth: 140
  },
  {
    field: 'createdBy',
    label: t('apiKey.createdBy'),
    minWidth: 100
  },
  {
    field: 'createdAt',
    label: t('apiKey.createdAt'),
    minWidth: 160,
    formatter: (_row, __: TableColumn, value: string) => formatDateTime(value)
  },
  {
    field: 'lastUsedAt',
    label: t('apiKey.lastUsedAt'),
    minWidth: 160,
    formatter: (_row, __: TableColumn, value: string) => formatDateTime(value)
  },
  {
    field: 'enabled',
    label: t('apiKey.status'),
    minWidth: 80,
    formatter: (_row, __: TableColumn, value: boolean) => (
      <ElTag type={value ? 'success' : 'info'} effect="light">
        {value ? t('apiKey.enabled') : t('apiKey.disabled')}
      </ElTag>
    )
  },
  {
    field: 'action',
    label: t('tableDemo.action'),
    minWidth: 100,
    fixed: 'right',
    formatter: (row) => (
      <BaseButton type="danger" onClick={() => handleDelete(row)}>
        {t('common.delete')}
      </BaseButton>
    )
  }
])

const { tableRegister, tableState, tableMethods } = useTable({
  fetchDataApi: async () => {
    const res = await getApiKeyListApi()
    return {
      list: res.data.list || []
    }
  }
})

const { loading, dataList } = tableState
const { getList } = tableMethods

function tableHeaderColor() {
  return { background: 'var(--el-fill-color-light)' }
}

const createDialogVisible = ref(false)
const keyDialogVisible = ref(false)
const createLoading = ref(false)
const form = ref({ name: '' })
const createdKey = ref('')

const openCreateDialog = () => {
  form.value.name = ''
  createDialogVisible.value = true
}

const handleCreate = async () => {
  const name = form.value.name.trim()
  if (!name) {
    return
  }
  createLoading.value = true
  try {
    const res = await createApiKeyApi(name)
    createDialogVisible.value = false
    createdKey.value = res.data.key
    keyDialogVisible.value = true
    getList()
  } finally {
    createLoading.value = false
  }
}

const handleCopyKey = () => {
  copy(createdKey.value)
}

const handleDelete = async (row: { id: string; name: string }) => {
  await ElMessageBox.confirm(
    t('apiKey.deleteConfirm', { name: row.name }),
    t('common.reminder'),
    {
      confirmButtonText: t('common.ok'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }
  )
  await deleteApiKeyApi(row.id)
  getList()
}

const closeKeyDialog = () => {
  keyDialogVisible.value = false
  createdKey.value = ''
}
</script>

<template>
  <ContentWrap>
    <ElAlert
      :title="t('apiKey.usageTitle')"
      type="info"
      :closable="false"
      show-icon
      class="mb-16px"
    >
      <template #default>
        <div class="api-key-usage">
          <p>{{ t('apiKey.usageDesc') }}</p>
          <pre class="api-key-example">{{ mcpConfigExample }}</pre>
        </div>
      </template>
    </ElAlert>

    <ElRow :gutter="16" class="mb-10px">
      <ElCol :span="4">
        <ElButton type="primary" @click="openCreateDialog">{{ t('apiKey.createTitle') }}</ElButton>
      </ElCol>
    </ElRow>

    <Table
      :columns="columns"
      :data="dataList"
      stripe
      :border="true"
      :loading="loading"
      :resizable="true"
      @register="tableRegister"
      :headerCellStyle="tableHeaderColor"
    />
  </ContentWrap>

  <Dialog
    v-model="createDialogVisible"
    :title="t('apiKey.createTitle')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="280"
  >
    <ElForm :model="form" label-width="80px" @submit.prevent>
      <ElFormItem :label="t('common.name')" required>
        <ElInput
          v-model="form.name"
          :placeholder="t('apiKey.namePlaceholder')"
          maxlength="64"
          show-word-limit
        />
      </ElFormItem>
      <ElFormItem>
        <ElButton type="primary" :loading="createLoading" @click="handleCreate">
          {{ t('common.ok') }}
        </ElButton>
        <ElButton @click="createDialogVisible = false">{{ t('common.cancel') }}</ElButton>
      </ElFormItem>
    </ElForm>
  </Dialog>

  <Dialog
    v-model="keyDialogVisible"
    :title="t('apiKey.keyCreatedTitle')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
    :maxHeight="360"
    @close="closeKeyDialog"
  >
    <ElAlert :title="t('apiKey.keyShowOnce')" type="warning" :closable="false" show-icon class="mb-16px" />
    <ElForm label-width="80px">
      <ElFormItem :label="t('apiKey.fullKey')">
        <ElInput v-model="createdKey" readonly>
          <template #append>
            <ElButton @click="handleCopyKey">
              {{ copied ? t('apiKey.copied') : t('apiKey.copyKey') }}
            </ElButton>
          </template>
        </ElInput>
      </ElFormItem>
    </ElForm>
    <div class="dialog-footer">
      <ElButton type="primary" @click="closeKeyDialog">{{ t('common.ok') }}</ElButton>
    </div>
  </Dialog>
</template>

<style scoped>
.api-key-usage p {
  margin: 0 0 8px;
  line-height: 1.6;
}

.api-key-example {
  margin: 0;
  padding: 12px;
  border-radius: 8px;
  background: var(--el-fill-color-light);
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  line-height: 1.5;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}

.mb-16px {
  margin-bottom: 16px;
}

.mb-10px {
  margin-bottom: 10px;
}
</style>

<script setup lang="ts">
import { Table } from '@/components/Table'
import { h, ref, watch } from 'vue'
import {
  ElAvatar,
  ElDropdown,
  ElDropdownMenu,
  ElDropdownItem,
  ElPagination,
  ElCheckbox,
  ElSwitch,
  ElForm,
  ElFormItem,
  ElCheckboxGroup
} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { Dialog } from '@/components/Dialog'
import { defineProps } from 'vue'
import { ElMessageBox } from 'element-plus'
import { deleteProjectApi } from '@/api/project'
import AddProject from './AddProject.vue'
import { useIcon } from '@/hooks/web/useIcon'
import { useRouter } from 'vue-router'
const { t } = useI18n()
const { push } = useRouter()
interface Recordable {
  id: string
  name: string
  logo: string
  AssetCount: number
  tag: string
}
const props = defineProps({
  tableDataList: {
    type: Array as () => Recordable[],
    default: () => []
  },
  getProjectTag: {
    type: Function as unknown as () => (pageIndex: number, pageSize: number) => Promise<void>,
    required: true
  },
  total: {
    type: Number,
    default: 0
  },
  multipleSelection: {
    type: Boolean
  },
  selectedRows: {
    type: Array as () => (string | number)[],
    default: () => []
  }
})
const loading = ref(false)
let ProjectId = ''
const dialogVisible = ref(false)
const closeDialog = () => {
  dialogVisible.value = false
}
const edit = async (id: string) => {
  ProjectId = id
  dialogVisible.value = true
}
const del = (id: string) => {
  // ElMessageBox.alert('Are you sure you want to delete the selected data?', '', {
  //   confirmButtonText: 'YES',
  //   callback: async () => {
  //     await deleteProjectApi([id])
  //     props.getProjectTag(currentPage.value, pageSize.value)
  //   }
  // })
  const deleteAsset = ref<boolean>(false)
  ElMessageBox({
    title: 'Delete',
    draggable: true,
    // Should pass a function if VNode contains dynamic props
    message: () =>
      h('div', { style: { display: 'flex', alignItems: 'center' } }, [
        h('p', { style: { margin: '0 10px 0 0' } }, t('task.delAsset')),
        h(ElSwitch, {
          modelValue: deleteAsset.value,
          'onUpdate:modelValue': (val: boolean) => {
            deleteAsset.value = val
          }
        })
      ])
  }).then(async () => {
    await deleteProjectApi([id], deleteAsset.value)
  })
}
const editIcon = useIcon({ icon: 'uil:edit' })
const delIcon = useIcon({ icon: 'material-symbols:delete-outline' })
const data = useIcon({ icon: 'carbon:data-vis-1' })
const handleCommand = (command: string | number | object) => {
  if (command['type'] == 'edit') {
    edit(command['id'])
  } else if (command['type'] == 'del') {
    del(command['id'])
  } else {
    action(command['id'])
  }
}
const handlePageChange = () => {
  props.getProjectTag(currentPage.value, pageSize.value)
}
const currentPage = ref(1)
const pageSize = ref(50)
const small = ref(false)
const background = ref(false)
const disabled = ref(false)
const emit = defineEmits<{
  (event: 'update:selectedRows', value: (string | number)[]): void
}>()

const localSelectedRows = ref([...props.selectedRows])

watch(localSelectedRows, (newVal) => {
  if (JSON.stringify(newVal) !== JSON.stringify(props.selectedRows)) {
    emit('update:selectedRows', newVal)
  }
  isAllSelected.value = newVal.length === props.tableDataList.length
})

watch(
  () => props.selectedRows,
  (newVal) => {
    localSelectedRows.value = [...newVal]
  }
)
watch(
  () => props.tableDataList,
  (newVal) => {
    isAllSelected.value = localSelectedRows.value.length === newVal.length
  }
)
const isAllSelected = ref(false)
const toggleAllSelection = () => {
  if (isAllSelected.value) {
    localSelectedRows.value = props.tableDataList.map((item) => item.id)
  } else {
    localSelectedRows.value = []
  }
}
const action = (id: string) => {
  push(`/project-management/project-detail?id=${id}`)
}
</script>

<template>
  <ElCheckbox v-if="multipleSelection" v-model="isAllSelected" @change="toggleAllSelection">
    {{ t('common.selectAll') }}
  </ElCheckbox>
  <ElCheckboxGroup v-model="localSelectedRows">
    <Table
      :columns="[]"
      :data="tableDataList"
      :loading="loading"
      custom-content
      :card-wrap-style="{
        width: '210px',
        marginBottom: '20px',
        marginRight: '20px'
      }"
    >
      <template #content="row">
        <ElDropdown trigger="contextmenu" @command="handleCommand">
          <div class="flex cursor-pointer">
            <ElCheckbox :value="row.id" class="pr-16px" v-if="multipleSelection" />
            <div class="pr-16px">
              <template v-if="row.logo != ''">
                <ElAvatar :src="row.logo" class="avatar" fit="cover" />
              </template>
              <template v-else>
                <ElAvatar class="avatar avatar-placeholder"> {{ row.name.charAt(0) }} </ElAvatar>
              </template>
            </div>
            <div>
              <div class="name">{{ row.name }}</div>
              <div class="assets-info">{{ t('project.totalAssets') }} : {{ row.AssetCount }}</div>
            </div>
          </div>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem :icon="editIcon" :command="{ type: 'edit', id: row.id }">{{
                t('common.edit')
              }}</ElDropdownItem>
              <ElDropdownItem :icon="delIcon" :command="{ type: 'del', id: row.id }">{{
                t('common.delete')
              }}</ElDropdownItem>
              <ElDropdownItem :icon="data" :command="{ type: 'aggregation', id: row.id }">{{
                t('project.aggregation')
              }}</ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Table>
  </ElCheckboxGroup>
  <ElPagination
    v-model:current-page="currentPage"
    v-model:page-size="pageSize"
    :page-sizes="[50, 70, 100, 200, 400]"
    :small="small"
    :disabled="disabled"
    :background="background"
    layout="total, sizes, prev, pager, next, jumper"
    :total="total"
    @size-change="handlePageChange"
    @current-change="handlePageChange"
  />
  <Dialog
    v-model="dialogVisible"
    :title="t('common.edit')"
    center
    style="border-radius: 15px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3)"
  >
    <AddProject
      :closeDialog="closeDialog"
      :projectid="ProjectId"
      :getProjectData="$props.getProjectTag"
      :schedule="false"
    />
  </Dialog>
</template>
<style>
.avatar {
  width: 45px;
  height: 45px;
  line-height: 45px;
  font-size: 24px;
}

.avatar-placeholder {
  background-color: cornflowerblue;
}
.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
.name {
  margin-bottom: 12px;
  font-weight: 700;
  font-size: 16px;
}

.assets-info {
  color: #b1b3b8;
  font-size: 11px;
  position: relative;
  top: -6px;
}
</style>

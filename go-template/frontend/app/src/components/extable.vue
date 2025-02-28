<script setup lang="ts">
import type { IApiQuery, IApiQueryPagination } from '@/model/api'
import { SORT_COLUMN_TYPES } from '@/model/api'
import type { ITableColumn } from '@/model/table'
import ApiService from '@/services/ApiService'
import { useUserStore } from '@/store/user'
import { pathJoin } from '@/utils/Path'
import { ErrorPopup, SuccessPopup, WarningPopup } from '@/utils/Popup'
import type { PropType } from 'vue'
import { VForm } from 'vuetify/components'

const props = defineProps({
  columns: {
    type: Array as () => Array<ITableColumn>,
    required: true,
  },
  apiUrl: {
    type: String,
    required: true,
  },
  createButton: {
    type: Boolean,
    required: true,
  },
  actionBar: {
    type: Boolean,
  },
  form: {
    type: Object,
    required: true,
  },
  resetForm: {
    required: false,
    type: Function,
  },
  tableActions: {
    type: Boolean,
    default: true,
  },
  emptyTableText: {
    type: String,
    default: '.....',
  },
  indexColumn: {
    type: Boolean,
    default: true,
  },
  actionsColumn: {
    type: Boolean,
    default: true,
  },
  defaultSortColumn: {
    type: String,
    default: 'id',
  },
  defaultSortColumnType: {
    type: Number as PropType<SORT_COLUMN_TYPES>,
    default: SORT_COLUMN_TYPES.NORMAL,
  },
  defaultSortOrder: {
    type: String as PropType<'asc' | 'desc'>,
    default: 'desc',
  },
  query: {
    type: Object as PropType<IApiQuery>,
    default: () => ({ query: [], columns: [], columnTypes: [] }),
  },
})

const emits = defineEmits(['update:form'])
const slots = useSlots()

const loading = ref(false)
const formLoading = ref(false)
const query = ref<IApiQuery>(props.query)

// 👉 Watching query
watch(
  () => props.query,
  newVal => {
    query.value = newVal
    fetchData()
  },
  { deep: true },
)

// 👉 Table refs
const rows = ref<any[]>([])
const sortColumn = ref(props.defaultSortColumn)
const sortColumnType = ref(props.defaultSortColumnType)
const sortOrder = ref(props.defaultSortOrder)

const pagination = ref({
  currentPage: 1,
  rowPerPage: 10,
  totalPage: 1,
  totalRows: 0,
})

// 👉 Computing pagination data
const paginationText = computed(() => {
  if (!rows.value.length)
    return ''

  const firstIndex = rows.value.length ? ((pagination.value.currentPage - 1) * pagination.value.rowPerPage) + 1 : 0
  const lastIndex = rows.value.length + ((pagination.value.currentPage - 1) * pagination.value.rowPerPage)

  return `${pagination.value.totalRows} satırdan ${firstIndex} - ${lastIndex} arası gösteriliyor`
})

// 👉 watching current page
watchEffect(() => {
  if (pagination.value.currentPage > pagination.value.totalPage)
    pagination.value.currentPage = pagination.value.totalPage
})

// 👉 compute columns
const columns = computed(() => {
  const userStore = useUserStore()
  const userRole = userStore.user.role

  if (!props.columns)
    return [] as ITableColumn[]

  const cols = [...props.columns]

  /*  // Index kolonu ekle
   if (props.indexColumn && (userRole === 10))
     cols.unshift({ key: 'i', name: '#', max_width: '10px' })
 
   // Role göre actions kolonunu kontrol et
   if (props.actionsColumn && (userRole === 10)) {
     cols.push({ key: 'actions', name: 'İŞLEMLER' })
   } */

  return cols
})
const sort = (sortable: boolean | undefined, columnKey: string, columnType: SORT_COLUMN_TYPES | undefined) => {
  if (!sortable)
    return

  if (sortOrder.value === 'asc')
    sortOrder.value = 'desc'
  else
    sortOrder.value = 'asc'

  // defaultta asc olsun
  if (sortColumn.value !== columnKey)
    sortOrder.value = 'asc'

  sortColumn.value = columnKey
  sortColumnType.value = columnType ?? SORT_COLUMN_TYPES.NORMAL
  // eslint-disable-next-line @typescript-eslint/no-use-before-define
  fetchData()
}

// 👉 Form refs
const formRef = ref<VForm>()
const isCreateForm = ref(false)
const isModalOpen = ref(false)

const closeModal = () => {
  isModalOpen.value = false
  formRef.value?.reset()
  props.resetForm?.()
}

const openModal = () => {
  isModalOpen.value = true
}

const fetchData = async () => {
  // if (sortColumn.value === '')
  //   sortColumn.value = 'created_at'
  //
  // if (sortOrder.value === '')
  //   sortOrder.value = 'desc'

  const q = {
    ...query.value,
  } as IApiQueryPagination

  // const q = { ...query.value }
  q.sortColumns = [sortColumn.value]
  q.sortOrders = [sortOrder.value]
  q.sortColumnTypes = [sortColumnType.value]

  q.page = pagination.value.currentPage
  q.perPage = pagination.value.rowPerPage

  loading.value = true

  const [error, resp] = await ApiService.get<any[]>(props.apiUrl, q)

  loading.value = false
  if (error)
    return ErrorPopup(error)

  rows.value = resp.data || []
  pagination.value.totalRows = resp.data_count
  pagination.value.totalPage = Math.ceil(resp.data_count / pagination.value.rowPerPage)

  // emit('update:modelValue', items)
  // currentSort.value = sortColumn.value + sortOrder.value
}

const setItemsPerPage = (val: number) => {
  if ((val * pagination.value.currentPage) > pagination.value.totalRows)
    pagination.value.currentPage = Math.ceil(pagination.value.totalRows / val)

  pagination.value.rowPerPage = val
  fetchData()
}

onMounted(() => {
  fetchData()
})

const onSubmit = async () => {
  const { valid } = await formRef.value!.validate()
  if (!valid)
    return
  formLoading.value = true
  let error
  if (isCreateForm.value)
    [error] = await ApiService.post<null>(props.apiUrl, props.form)
  else
    [error] = await ApiService.put<null>(pathJoin(props.apiUrl, props.form.id), props.form)
  formLoading.value = false
  if (error)
    return ErrorPopup(error)

  isModalOpen.value = false
  fetchData()
  SuccessPopup('İşlem Başarılı')
}

const onCreate = () => {
  isCreateForm.value = true
  isModalOpen.value = true
  nextTick(() => {
    formRef.value?.reset()
    props.resetForm?.()
  })
}

const onEdit = async (id: number) => {
  isCreateForm.value = false

  formLoading.value = true

  const [error, resp] = await ApiService.get<typeof props.form>(pathJoin(props.apiUrl, id))

  formLoading.value = false
  if (error)
    return ErrorPopup(error)

  emits('update:form', resp.data)
  isModalOpen.value = true
}

const onDelete = async (id: number) => {
  const confirm = await WarningPopup(
    'İşlemi Onaylıyor musunuz?',
    'Evet',
    'Hayır',
  )

  if (!confirm.isConfirmed)
    return

  loading.value = true

  const [error] = await ApiService.delete<null>(pathJoin(props.apiUrl, id))

  loading.value = false
  if (error)
    return ErrorPopup(error)

  fetchData()
  SuccessPopup('İşlem Başarılı')
}

const refresh = async (q?: IApiQuery) => {
  if (q)
    query.value = { ...q }

  await fetchData()
}

defineExpose({
  refresh,
})
</script>

<template>
  <div>
    <VCard flat>
      <VCardText v-if="createButton || actionBar" class="d-flex justify-end">
        <slot name="actionBar" :actionloading="loading" />
        <VBtn v-if="createButton" color="primary" prepend-icon="tabler-edit" @click="onCreate">
          Ekle
        </VBtn>
      </VCardText>
      <VCardText>
        <VTable class="text-no-wrap">
          <!-- 👉 table head -->
          <thead>
            <tr>
              <th v-for="(column, i) in columns" :key="i" scope="col" :style="{ 'max-width': column.max_width }" @click="
                sort(
                  column.sortable,
                  column.sortField ? column.sortField : column.key,
                  column.sortFieldType,
                )
                ">
                <div class="d-flex align-center" :class="{ 'justify-end': column.key === 'actions' }">
                  <span style="cursor: pointer" class="text-center">
                    {{ column.name }}
                  </span>
                  <VIcon v-if="sortColumn === column.key"
                    :icon="sortOrder === 'asc' ? 'material-symbols:keyboard-arrow-up-rounded' : 'material-symbols:keyboard-arrow-down-rounded'" />
                </div>
              </th>
            </tr>
          </thead>

          <!-- 👉 table body -->
          <tbody>
            <template v-for="(row, i) in rows" :key="i">
              <tr style="height: 3.75rem;">
                <template v-for="(column, j) in columns" :key="j">
                  <td :class="{ 'text-end': column.key === 'actions' }" :style="{ 'max-width': column.max_width }">
                    <span v-if="column.key === 'i'">{{ i + 1 }}</span>
                    <span v-else-if="column.key === 'actions' && tableActions">
                      <slot name="actions" :row="row" />
                      <VBtn icon size="x-small" color="default" variant="text" @click="onEdit(row.id)">
                        <VIcon size="22" icon="tabler-edit" />
                      </VBtn>

                      <VBtn icon size="x-small" color="default" variant="text" @click="onDelete(row.id)">
                        <VIcon size="22" color="error" icon="tabler-trash" />
                      </VBtn>
                    </span>
                    <slot v-else-if="slots[column.key]" :name="column.key" :row="row" />
                    <!-- column slot olarak verilmediyse direk row içinden ilgili kolonu span olarak yaz -->
                    <span v-else>{{ row[column.key] }}</span>
                  </td>
                </template>
              </tr>
            </template>
          </tbody>

          <!-- 👉 table footer  -->
          <tfoot v-show="!rows.length">
            <tr>
              <td colspan="7" class="text-center">
                {{ emptyTableText }}
              </td>
            </tr>
          </tfoot>
        </VTable>
      </VCardText>
      <VDivider />

      <VCardText class="d-flex align-center flex-wrap justify-space-between gap-4 py-3 px-5">
        <div class="d-flex align-center">
          <div class="me-3" style="width: 80px;">
            <VSelect v-model="pagination.rowPerPage" density="compact" variant="outlined" :items="[10, 20, 30, 50]"
              @update:model-value="setItemsPerPage" />
          </div>
          <span class="text-sm text-disabled">
            {{ paginationText }}
          </span>
        </div>

        <VPagination v-model="pagination.currentPage" size="small" :total-visible="5" :length="pagination.totalPage"
          @update:model-value="fetchData" />
      </VCardText>
      <VOverlay v-model="loading" class="align-center justify-center" contained persistent>
        <VProgressCircular :size="50" color="primary" indeterminate />
      </VOverlay>
    </VCard>
    <VDialog v-model="isModalOpen" :width="$vuetify.display.smAndDown ? 'auto' : 700" persistent>
      <DialogCloseBtn @click="closeModal" />
      <VCard :loading="formLoading" class="pa-sm-8 pa-5">
        <!--        <VCardItem class="text-center"> -->
        <!--          <VCardTitle class="text-h5 mb-3"> -->
        <!--            {{ formCardTitle }} -->
        <!--          </VCardTitle> -->
        <!--        </VCardItem> -->

        <VCardText>
          <VForm ref="formRef" @submit.prevent="onSubmit">
            <slot name="modalBody" :formref="formRef" :closemodal="closeModal" :openmodal="openModal"
              :formloading="formLoading" :iscreateform="isCreateForm" />
          </VForm>
        </VCardText>
      </VCard>
    </VDialog>
  </div>
</template>

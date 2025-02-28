<script setup lang="ts">
import { ref } from 'vue'
import type { ITableColumn } from '@/model/table'
import { ApiQuery } from '@/model/api'
import { useUserStore } from '@/store/user'
import Extable from '@/components/extable.vue'
import tarihFormat from '@/utils/ExDate'
import ApiService from '@/services/ApiService'
import { ErrorPopup, WarningPopup } from '@/utils/Popup'
import { router } from '@/plugins/1.router'

interface Like {
  id: number
  user_id: number
}

interface Comment {
  id: number
  post_id: number
  user_id: number
  content: string
  user?: {
    name: string
    surname: string
  }
}

interface Post {
  id: number
  title: string
  content: string
  main_content: string
  user_id: number
  like_count: number
  comment_count: number
  created_at: string
  likes: Like[]
  comments: Comment[]
  user: {
    name: string
    surname: string
  }
}

// Interface ekleyelim
interface MyUser {
  id: number
  name: string
  surname: string
  role: number
}

const apiUrl = 'post/'
const userStore = useUserStore()
const users = ref<MyUser[]>([])
const searchQuery = ref('')
const statusFilter = ref('all')
const deleteDialog = ref(false)
const postIdToDelete = ref<number | null>(null)

const form = ref<Post>({
  id: 0,
  title: '',
  content: '',
  main_content: '',
  user_id: 0,
  like_count: 0,
  comment_count: 0,
  created_at: '',
  likes: [],
  comments: [],
  user: {
    name: '',
    surname: ''
  }
})

// Tüm içeriği tek kolonda göstereceğiz
const columns = ref<ITableColumn[]>([
  {
    key: 'title',
    name: 'BAŞLIK',
    sortable: true,
  },

  {
    key: 'user',
    name: 'YAZAR',
  },
  {
    key: 'like_count',
    name: 'BEĞENİ',
  },
  {
    key: 'comment_count',
    name: 'YORUM',
  },
  {
    key: 'created_at',
    name: 'TARİH',
  },
])

const tableQuery = ref<ApiQuery>(new ApiQuery())
form.value.user_id = userStore.user.id
form.value.user.name = userStore.user.name
form.value.user.surname = userStore.user.surname

const navigateToPost = (postId: number) => {
  router.push(`third_page/${postId}`)
}


const confirmDelete = (postId: number) => {
  postIdToDelete.value = postId
  deleteDialog.value = true
}

const deletePost = async () => {
  if (!postIdToDelete.value) return

  try {
    const [error] = await ApiService.delete(`post/${postIdToDelete.value}`)

    if (error) {
      ErrorPopup('Gönderi silinirken bir hata oluştu')
      return
    }

    // Tabloyu yenile
    // Muhtemelen Extable bu işi kendisi yapar
  } catch (error) {
    ErrorPopup('Bir hata oluştu')
  } finally {
    deleteDialog.value = false
    postIdToDelete.value = null
  }
}

// Kullanıcıları yükleyecek fonksiyon
const loadUsers = async () => {
  const [error, response] = await ApiService.get<MyUser[]>('user')
  if (!error) {
    users.value = response.data
  }
}

// Post veya yorum sahibinin adını getiren fonksiyon
const getUserFullName = (userId: number) => {
  const user = users.value.find(u => u.id === userId)
  return user ? `${user.name} ${user.surname}` : 'Bilinmeyen Kullanıcı'
}

const getUserNameForIcon = (userId: number) => {
  const user = users.value.find(u => u.id === userId)
  return user?.name[0]
}

// Component mount olduğunda kullanıcıları yükle
onMounted(() => {
  loadUsers()
})

const createNewPost = () => {
  form.value = {
    id: 0,
    title: '',
    content: '',
    main_content: '',
    user_id: userStore.user.id,
    like_count: 0,
    comment_count: 0,
    created_at: '',
    likes: [],
    comments: [],
    user: {
      name: userStore.user.name,
      surname: userStore.user.surname
    }
  }
  // Extable bileşeni create modunu tetikleyecek
}
</script>

<template>
  <div class="admin-posts-container">
    <VCard>
      <VCardText>
        <VRow>
          <VCol
            sm="6"
            md="4"
            lg="3"
          >
            <VTextField
              v-model="searchQuery"
              label="Arama"
              @update:model-value="tableQuery.append($event, 'title;content')"
            />
          </VCol>
        </VRow>
      </VCardText>
      <VDivider />

      <Extable 
        v-model:form="form" 
        :api-url="apiUrl" 
        :columns="columns" 
        :query="tableQuery" 
        create-button 
        table-actions
      >
        <!-- User column template -->
        <template #user="{ row }">
          <div class="d-flex align-center">
            <span>{{ getUserFullName(row.user_id) }}</span>
          </div>
        </template>

        <!-- Created at column template -->
        <template #created_at="{ row }">
          {{ tarihFormat(row.created_at) }}
        </template>

        <!-- Modal template -->
        <template #modalBody="{ closemodal, formloading, iscreateform }">
          <div v-if="iscreateform || userStore.getRole() === 10 || form.user_id === userStore.user.id"
            class="d-flex gap-2">
            <VRow>
              <VCol cols="12">
                <VTextField v-model="form.title" label="Başlık" :rules="[v => !!v || 'Başlık zorunludur']" />
              </VCol>



              <VCol cols="12">
                <VTextarea v-model="form.main_content" label="Ne düşünüyorsun?" rows="6" counter auto-grow no-resize
                  clearable :max-rows="8" :rules="[v => !!v || 'İçerik zorunludur']" />
              </VCol>
              <VCol cols="12" class="d-flex justify-end gap-4">
                <VBtn type="submit" :loading="formloading" color="primary">
                  {{ iscreateform ? 'Paylaş' : 'Güncelle' }}
                </VBtn>
                <VBtn type="reset" color="secondary" variant="tonal" :loading="formloading" @click="closemodal">
                  Vazgeç
                </VBtn>
              </VCol>
            </VRow>
          </div>
          <div v-else>
            Bu gönderiyi düzenleme yetkiniz yok.
            <VBtn variant="tonal" :loading="formloading" @click="closemodal">
              Vazgeç
            </VBtn>
          </div>
        </template>
      </Extable>
    </VCard>
  </div>
</template>

<style scoped>
.admin-posts-container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
}
</style>

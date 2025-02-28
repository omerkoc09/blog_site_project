<script setup lang="ts">
import tarihFormat from '@/utils/ExDate'
import { useUserStore } from '@/store/user'
import { ErrorPopup } from '@/utils/Popup'
import ApiService from '@/services/ApiService'
import ImageUploader from '@/components/ImageUploader.vue'
import { ref } from "vue"
import AppBar from '@/components/AppBar.vue'


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
  image: string
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

interface MyUser {
  id: number
  name: string
  surname: string
  role: number
}

const newComment = ref('')
const showCommentInput = ref<number | null>(null)
const editingCommentId = ref<number | null>(null)
const editingCommentPostId = ref<number | null>(null)
const editingCommentUserId = ref<number | null>(null)
const editedComment = ref('')
const users = ref<MyUser[]>([])


const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)
const postId = computed(() => (route.params as { id: string }).id)
console.log('Post ID:', postId.value)


const form = ref<Post>({
  id: 0,
  title: '',
  content: '',
  main_content: '',
  image: '',
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


form.value.user_id = userStore.user.id //veya backend tarafında tokendan alıp orada da doldurulabilir.
form.value.user.name = userStore.user.name
form.value.user.surname = userStore.user.surname

const isLikedByUser = (likes: Like[]) => {
  return likes.some((like: Like) => like.user_id === userStore.user.id)
}

const handleLike = async (postId: number, post: any) => {
  const [error] = await ApiService.post(`like/${postId}`, {})

  if (!userStore.isAuthenticated) {
    ErrorPopup('Bu işlemi yapabilmek için giriş yapmanız gerekmektedir.')
    return
  }

  if (!error) {
    // Yerel state'i güncelle
    if (isLikedByUser(post.likes)) {
      // Eğer kullanıcı zaten beğenmişse, beğeniyi kaldır
      post.likes = post.likes.filter((like: { user_id: any }) => like.user_id !== userStore.user.id)
      post.like_count--
    } else {
      // Beğeni ekle
      post.likes.push({ user_id: userStore.user.id })
      post.like_count++
    }
  }
}

const handleComment = async (postId: number, post: any) => {
  console.log('Mevcut yorumlar:', post.comments) // Debug için mevcut yorumları görelim

  if (!userStore.isAuthenticated) {
    ErrorPopup('Bu işlemi yapabilmek için giriş yapmanız gerekmektedir.')
    return
  }
  if (!newComment.value.trim()) return
  try {

    const [error, response] = await ApiService.post('comment', {
      post_id: postId,
      content: newComment.value

    })

    console.log('API yanıtı:', response) // Debug için API yanıtını görelim

    if (!error) {
      // Backend'den gelen yanıtı kullan
      post.comments.push({
        id: (response.data as any).id, // Backend'den gelen id'yi kullan
        content: newComment.value,
        post_id: postId,
        user_id: userStore.user.id,
        user: {
          name: userStore.user.name,
          surname: userStore.user.surname
        }
      })

      post.comment_count++

      newComment.value = ''

      showCommentInput.value = null

    }

  } catch (error) {
    console.error('Yorum hatası:', error) // Hata varsa görelim

    ErrorPopup('Yorum gönderilirken bir hata oluştu')

  }

}

// Yorum düzenleme işlemleri
const startEdit = (comment: Comment) => {
  editingCommentId.value = comment.id
  editedComment.value = comment.content
  editingCommentPostId.value = comment.post_id
  editingCommentUserId.value = comment.user_id

}

const cancelEdit = () => {
  editingCommentId.value = null
  editingCommentPostId.value = null
  editingCommentUserId.value = null
  editedComment.value = ''
}

const updateComment = async (commentId: number, post: any) => {  // post parametresini ekledik
  if (!editedComment.value.trim()) return

  try {
    const [error] = await ApiService.put(`comment/${commentId}`, {
      id: editingCommentId.value,
      content: editedComment.value,
      post_id: editingCommentPostId.value,
      user_id: editingCommentUserId.value
    })

    if (!error) {
      const comment = post.comments.find((c: { id: number }) => c.id === commentId)
      if (comment) {
        comment.content = editedComment.value
      }
      editingCommentId.value = null
      editedComment.value = ''
      editingCommentPostId.value = null
      editingCommentUserId.value = null
    }
  } catch (error) {
    ErrorPopup('Yorum güncellenirken bir hata oluştu')
  }
}

// Yorum silme işlemi
const deleteComment = async (commentId: number, post: any) => {  // post parametresini ekledik
  const result = await WarningPopup(
    'Yorumu silmek istediğinize emin misiniz?',
    'Sil',
    'İptal'
  )

  if (!result.isConfirmed) return

  try {
    const [error] = await ApiService.delete(`comment/${commentId}`)
    if (!error) {
      // Direkt gelen post üzerinde güncelleme yap
      post.comments = post.comments.filter((c: { id: number }) => c.id !== commentId)
      post.comment_count--
    }
  } catch (error) {
    ErrorPopup('Yorum silinirken bir hata oluştu')
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

// Script bölümüne eklenecek fonksiyonlar ve değişkenler
const showEditModal = ref(false)
const deleteConfirmDialog = ref(false)
const updateLoading = ref(false)
const deleteLoading = ref(false)
const showCreateModal = ref(false)


// Kullanıcı post sahibi mi kontrolü
const isPostOwner = computed(() => {
  console.log('post id', form.value.id)
  console.log('user id', userStore.user.id)

  return form.value.user_id === userStore.user.id
})


// Post silme işlevi
const deletePost = async () => {
  if (!isPostOwner.value) {
    ErrorPopup('Bu gönderiyi silme yetkiniz yok')
    return
  }

  deleteLoading.value = true
  try {
    const [error] = await ApiService.delete(`post/${form.value.id}`)

    if (error) {
      ErrorPopup('Gönderi silinirken bir hata oluştu')
      return
    }

    // Ana sayfaya yönlendir
    router.push('/')
  } catch (error) {
    ErrorPopup('Bir hata oluştu')
  } finally {
    deleteLoading.value = false
    deleteConfirmDialog.value = false
  }
}



onMounted(async () => {

  // Test için loglar
  console.log('Route params:', route.params)
  console.log('Full route:', route)

  loading.value = true
  const [err, data] = await ApiService.get<any>(`post/${postId.value}`)
  loading.value = false
  console.log('Post Verisi:', data)  // Gelen veriyi kontrol et    
  console.log('Post ID:', postId.value)
  if (err)
    return ErrorPopup(err)

  form.value = data.data
  loadUsers()

  console.log('form value', form.value)


})


// Resim yükleme için değişkenler
const selectedFile = ref<File | null>(null)

const profileMenu = ref(false)



const updatePost = async () => {
  if (!isPostOwner.value) {
    ErrorPopup('Bu gönderiyi düzenleme yetkiniz yok')
    return
  }

  updateLoading.value = true
  try {
    // Create FormData object
    const formData = new FormData()
    formData.append('id', form.value.id.toString())
    formData.append('title', form.value.title)
    formData.append('content', form.value.content)
    formData.append('main_content', form.value.main_content)
    formData.append('user_id', String(userStore.user.id))

    // If there's a selected file, append it
    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    // FormData ile gönderirken Content-Type header'ı otomatik olarak ayarlanacak
    const [error] = await ApiService.put(`post/${form.value.id}`, formData)

    if (error) {
      ErrorPopup('Gönderi güncellenirken bir hata oluştu')
      return
    }

    showEditModal.value = false
    // Başarı mesajı göster
    const successMessage = document.createElement('div')
    successMessage.className = 'success-popup'
    successMessage.textContent = 'Gönderi başarıyla güncellendi'
    document.body.appendChild(successMessage)
    setTimeout(() => {
      document.body.removeChild(successMessage)
    }, 3000)

    // Sayfayı yenile
    const [err, data] = await ApiService.get<any>(`post/${form.value.id}`)
    if (!err) {
      form.value = data.data
      selectedFile.value = null // Reset selected file after successful update
    }
  } catch (error) {
    ErrorPopup('Bir hata oluştu')
  } finally {
    updateLoading.value = false
  }
}

const bottomRef = ref<HTMLElement | null>(null)

const isAuthenticated = () => {
  if (!userStore.isAuthenticated) {
    ErrorPopup('Bu işlemi yapabilmek için giriş yapmanız gerekmektedir.')
    return
  }
}

const navigateToAuthor = (userId: number) => {
  router.push(`../author/${userId}`)
}

const scrollToBottom = () => {

  if (!userStore.isAuthenticated) {
    ErrorPopup('Bu işlemi yapabilmek için giriş yapmanız gerekmektedir.')
    return
  }
  else
  bottomRef.value?.scrollIntoView({ behavior: "smooth" })
}

// Çıkış yapma fonksiyonu
const logout =  () => {
  try {
    useUserStore().logout()
    router.push('/')
  } catch (error) {
    ErrorPopup('Çıkış yapılırken bir hata oluştu')
  }
}

const navigateToLogin = () => {
  router.push('../auth/login')
}

const navigateToRegister = () => {
  router.push('../auth/register')
}

</script>

<template>

  <!-- Site Başlığı ve Profil Menüsü -->
  <VAppBar color="white" elevation="1">
    <VContainer class="d-flex align-center pa-0">
      <!-- Sol taraf: Site Adı -->
      <div class="d-flex align-center" >
        <RouterLink to="/" >
          <VRow to="/" color="black">

            <VAppBarTitle class="font-weight-bold" >BLOGGER.COM</VAppBarTitle>
          </VRow>
        </RouterLink>

      </div>

      <VSpacer />

      <!-- Sağ taraf: Profil Menüsü -->
      <div class="d-flex align-center" >

        <!-- Giriş yapmamış kullanıcılar için butonlar -->
        <div v-if="!userStore.user.name" class="d-flex align-center">
          <VBtn variant="outlined" color="primary" class="me-2" @click="navigateToLogin">
            Giriş Yap
          </VBtn>
          <VBtn variant="outlined" color="primary" @click="navigateToRegister">
            Kaydol
          </VBtn>
        </div>


        <!-- Profil Menüsü -->
        <VMenu v-model="profileMenu" location="bottom end" v-if="userStore.user.name">
          <template v-slot:activator="{ props }">
            <VBtn
                icon
                v-bind="props"
            >
              <VAvatar color="primary" size="40">
                <span class="text-h6 text-white">{{ userStore.user.name[0] }}</span>
              </VAvatar>
            </VBtn>
          </template>

          <VList width="220">
            <VListItem to="/profile">
              <template #prepend>
                <VIcon icon="tabler-user" size="small" class="me-2" />
              </template>
              <VListItemTitle>Profil</VListItemTitle>
            </VListItem>

            <VDivider class="my-2" />

            <VListItem  color="error" @click="logout">
              <template #prepend>
                <VIcon icon="tabler-logout" size="small" class="me-2" color="error"/>
              </template>
              <VListItemTitle>Çıkış Yap</VListItemTitle>
            </VListItem>
          </VList>
        </VMenu>
      </div>
    </VContainer>
  </VAppBar>

    <div class="posts-container">
    <!-- Loading state -->
    <div v-if="loading" class="d-flex justify-center my-4">
      <VProgressCircular indeterminate />
    </div>

    <!-- Post content -->
    <template v-else>
      <VCard class="mb-6">
        <!-- İşlem butonları - sadece post sahibi görebilir -->
        <VCardItem v-if="isPostOwner" class="d-flex justify-end pa-2">

          <VBtn text="Düzenle" variant="outlined" color="primary" class="me-2" @click="showEditModal = true">
            <VIcon start>tabler-edit</VIcon>
            Düzenle
          </VBtn>

          <VBtn text="Sil" variant="outlined" color="error" @click="deleteConfirmDialog = true">
            <VIcon start>tabler-trash</VIcon>
            Sil
          </VBtn>
        </VCardItem>

        <!-- Post title -->
        <VCardTitle class="text-h3 px-4 text-center text-wrap" :style="{
          wordBreak: 'break-word',
          whiteSpace: 'normal',
          lineHeight: '1.2',
          overflow: 'visible'
        }">
          {{ form.title }}
        </VCardTitle>

        <!-- User info and date -->
        <VCardItem>
          <template #prepend>
            <VAvatar color="primary" size="40">
              {{ getUserNameForIcon(form.user_id) }}
            </VAvatar>
          </template>

          <div>
            <VCardTitle class="user-name">
              <a class="author-link" @click="navigateToAuthor(form.user_id)" target="_blank"> {{ getUserFullName(form.user_id) }}</a>
            </VCardTitle>
            <VCardSubtitle class="post-date">
              {{ tarihFormat(form.created_at) }}
            </VCardSubtitle>
          </div>
        </VCardItem>

        <!-- Like and comment buttons -->
        <VCardActions>
          <VBtn variant="text" prepend-icon="tabler-heart" :color="isLikedByUser(form.likes) ? 'error' : 'default'"
            @click="handleLike(form.id, form)">
            {{ form.like_count }}
          </VBtn>

          <VBtn variant="text" prepend-icon="tabler-message-circle" color="default"
            @click="showCommentInput = showCommentInput === form.id ? null : form.id ; scrollToBottom()">
            {{ form.comment_count }} 
          </VBtn>
          <VSpacer />
        </VCardActions>



        <VCardText v-if="form.image" class="text-center">
          <img :src="`http://localhost:3001/${form.image}`" :alt="form.title" class="post-image" />
        </VCardText>

        <VCardText class="text-body-1 text-left" :style="{
          fontSize: '1.5rem !important',
          lineHeight: '1.5',
          color: 'black',
          whiteSpace: 'pre-wrap',
          fontWeight: '400',
          letterSpacing: '0.5px',
          wordSpacing: '2px',
          textTransform: 'none',
          padding: '1rem',
          margin: '1rem 0',
          maxWidth: '100%',
          overflowWrap: 'break-word'
        }">
          {{ form.main_content }}
        </VCardText>

        <VDivider />
      </VCard>
    </template>

    <!-- Düzenleme Modal -->
    <VDialog v-model="showEditModal" max-width="600px">
      <VCard>
        <VCardTitle class="text-h5 bg-primary text-black pa-4">
          Gönderi Düzenle
        </VCardTitle>

        <VCardText class="pa-4">
          <VForm @submit.prevent="updatePost">
            <VRow>
              <ImageUploader 
                :post-id="form.id" 
                :current-image="form.image"
                v-model:selectedFile="selectedFile"
              />
              
              <VCol cols="12">
                <VTextField v-model="form.title" label="Başlık" :rules="[v => !!v || 'Başlık zorunludur']" required />
              </VCol>

              <VCol cols="12">
                <VTextField v-model="form.content" label="Giriş" :rules="[v => !!v || 'Giriş zorunludur']" required />
              </VCol>

              <VCol cols="12">
                <VTextarea v-model="form.main_content" label="Ana İçerik" :rules="[v => !!v || 'İçerik zorunludur']"
                  rows="6" auto-grow counter required />
              </VCol>
            </VRow>

            <VCardActions class="pa-0 mt-4">
              <VSpacer />
              <VBtn variant="text" @click="showEditModal = false">
                İptal
              </VBtn>
              <VBtn color="primary" type="submit" :loading="updateLoading">
                Güncelle
              </VBtn>
            </VCardActions>
          </VForm>
        </VCardText>
      </VCard>
    </VDialog>

    <!-- Silme Onay Dialog -->
    <VDialog v-model="deleteConfirmDialog" max-width="500px">
      <VCard>
        <VCardTitle class="text-h5 bg-error text-white pa-4">
          Gönderiyi Sil
        </VCardTitle>

        <VCardText class="pa-4 pt-5">
          <p>Bu gönderiyi silmek istediğinize emin misiniz? Bu işlem geri alınamaz.</p>
        </VCardText>

        <VCardActions>
          <VSpacer />
          <VBtn variant="text" @click="deleteConfirmDialog = false">
            İptal
          </VBtn>
          <VBtn color="error" @click="deletePost" :loading="deleteLoading">
            Sil
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>

      <!-- Comments section -->
      <VExpandTransition>
        <div class="pa-4 bg-grey-lighten-4">
          <!-- Existing comments -->
          <div v-if="form.comments.length" class="mb-4">
            <div v-for="comment in form.comments" :key="comment.id" class="mb-3">
              <v-card>
              <div class="d-flex align-center justify-space-between">
                <div class="d-flex align-center">
                  <VAvatar color="primary" size="32" class="me-2">
                    {{ getUserNameForIcon(comment.user_id) }}
                  </VAvatar>
                  <div>
                    <div class="font-weight-medium">
                      {{ getUserFullName(comment.user_id) }}
                    </div>
                    <!-- Edit mode -->
                    <div v-if="editingCommentId === comment.id"  >
                      <VTextField v-model="editedComment" density="compact" variant="outlined"
                                  :rules="[v => !!v || 'Yorum boş olamaz']" @keyup.enter="updateComment(comment.id, form)" />
                      <div class="d-flex gap-2 mt-1">
                        <VBtn size="small" color="primary" @click="updateComment(comment.id, form)">
                          Kaydet
                        </VBtn>
                        <VBtn size="small" variant="outlined" @click="cancelEdit">
                          İptal
                        </VBtn>
                      </div>
                    </div>
                    <!-- Normal view -->
                    <div v-else class="text-body-2">
                      {{ comment.content }}
                    </div>
                  </div>
                </div>
                <!-- Edit/Delete buttons -->
                <div v-if="comment.user_id === userStore.user.id" class="d-flex gap-2">
                  <VBtn icon="tabler-edit" size="x-small" variant="text" @click="startEdit(comment)" />
                  <VBtn icon="tabler-trash" size="x-small" variant="text" color="error"
                        @click="deleteComment(comment.id, form)" />
                </div>
              </div>
              </v-card>
            </div>

          </div>
          <div ref="bottomRef"></div>

          <!-- Comment input -->
          <div  class="mt-3">
            <VTextField v-model="newComment" label="Yorum yaz..." :rules="[v => !!v || 'Yorum boş olamaz']"
                        append-inner-icon="tabler-send" @click:append-inner="handleComment(form.id, form)"
                        @keyup.enter="handleComment(form.id, form)" />
          </div>
        </div>
      </VExpandTransition>
  </div>

</template>

<style scoped>
.posts-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
  /* viewport height */
  overflow-y: auto;
  /* dikey scroll ekler */
}

.text-h3 {
  font-size: 3rem !important;
  /* Apply a larger font size with higher priority */
  font-weight: bold;
  color: #000000;
  margin-bottom: 10px;
}

.user-name {
  font-size: 1rem;
  /* Smaller font size */
}

.post-date {
  font-size: 1rem;
  /* Even smaller font size */
  color: #6c757d;
}

.posts-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 80px;
  min-height: 100vh;
  /* viewport height */
  overflow-y: auto;
  /* dikey scroll ekler */
}

.text-h3 {
  font-size: 3rem !important;
  /* Apply a larger font size with higher priority */
  font-weight: bold;
  color: #000000;
  margin-bottom: 10px;
}

.user-name {
  font-size: 1rem;
  /* Smaller font size */
}

.post-date {
  font-size: 1rem;
  /* Even smaller font size */
  color: #6c757d;
}

/* Başarılı güncelleme bildirimi */
.success-popup {
  position: fixed;
  top: 20px;
  right: 20px;
  background-color: #4caf50;
  color: white;
  padding: 16px 24px;
  border-radius: 4px;
  z-index: 9999;
  box-shadow: 0 3px 5px rgba(0, 0, 0, 0.2);
  animation: fadeInOut 3s ease;
}

@keyframes fadeInOut {

  0%,
  100% {
    opacity: 0;
    transform: translateY(-20px);
  }

  10%,
  90% {
    opacity: 1;
    transform: translateY(0);
  }
}

.image-upload-area {
  border: 2px dashed #ccc;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}


.preview-container {
  position: relative;
  width: 100%;
  max-height: 300px;
  overflow: hidden;
  border-radius: 4px;
}

.preview-image {
  width: 100%;
  height: auto;
  object-fit: contain;
  max-height: 300px;
}

.remove-image {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0.8;
}

.remove-image:hover {
  opacity: 1;
}

.post-image {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin: 1rem 0;
  object-fit: contain;
  /* Resmi orantılı şekilde sığdırır */
  max-height: 500px;
  /* Maksimum yükseklik */
}

.post-image:hover {
  transform: scale(1.01);
}

.author-link {
  color: #555;
  font-size: 15px;
  text-decoration: none;
  position: relative;
  cursor: pointer;
}

.author-link:hover {
  text-decoration: underline;
}
</style>

<route lang="yaml">
  meta:
    layout: blank  # Blank layout kullanılıyor
  </route>

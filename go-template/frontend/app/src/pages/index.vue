<script setup lang="ts">
import { ref } from 'vue'
import type { ITableColumn } from '@/model/table'
import { ApiQuery } from '@/model/api'
import { useUserStore } from '@/store/user'
import tarihFormat from '@/utils/ExDate'
import ApiService from '@/services/ApiService'
import { ErrorPopup, WarningPopup } from '@/utils/Popup'
import { router } from '@/plugins/1.router'
import apiService from "@/services/ApiService";
import ImageUploader from "@/components/ImageUploader.vue";
import Pagination from '@/components/Pagination.vue'





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

// Interface ekleyelim
interface MyUser {
  id: number
  name: string
  surname: string
  role: number
}

const apiUrl = 'post/'
const userStore = useUserStore()
const newComment = ref('')
const showCommentInput = ref<number | null>(null)
const users = ref<MyUser[]>([])
const profileMenu = ref(false)

// Form ve tablo ayarları
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



console.log('form value2', form.value)

const navigateToPost = (postId: number) => {
  router.push(`third_page/${postId}`)
}

const navigateToAuthor = (userId: number) => {
  router.push(`author/${userId}`)
}

const isLikedByUser = (likes: Like[]) => {
  return likes.some((like: Like) => like.user_id === userStore.user.id)
}

const handleLike = async (postId: number, post: any) => {
  const [error] = await ApiService.post(`like/${postId}`, {})


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

// Kullanıcıları yükleyecek fonksiyon
const loadUsers = async () => {
  const [error, response] = await ApiService.get<MyUser[]>('user')
  if (!error) {
    users.value = response.data
  }

  console.log('users', users.value)
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

// Çıkış yapma fonksiyonu
const logout =  () => {
  try {
    useUserStore().logout()
    router.push('/')
  } catch (error) {
    ErrorPopup('Çıkış yapılırken bir hata oluştu')
  }
}

// Resim yükleme için gerekli değişkenler
const fileInput = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const imagePreview = ref<string | null>(null)
const selectedFile = ref<File | null>(null)
const showEditModal = ref(false)
const loading = ref(false)
const updateLoading = ref(false)

const createImage = async (file: File) => {
  if (!file) return
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('image', file)
    const [error, response] = await ApiService.post(
        `post/image`,
        formData
    )
    if (error) {
      ErrorPopup(error)
      return null
    }
    SuccessPopup('Resim başarıyla yüklendi')
    return response.data.imageUrl // Backendden dönen resim URL'sini varsayıyorum
  } catch (err) {
    ErrorPopup('Bir hata oluştu')
    return null
  } finally {
    loading.value = false
  }
}


// Resim yükleme fonksiyonları
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files && input.files[0]) {
    handleImage(input.files[0])
  }
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  if (event.dataTransfer?.files[0]) {
    handleImage(event.dataTransfer.files[0])
  }
}

const handleImage = (file: File) => {
  if (file.size > 5 * 1024 * 1024) {
    ErrorPopup('Dosya boyutu 5MB\'dan küçük olmalıdır')
    return
  }


  selectedFile.value = file
  console.log('selected', selectedFile.value)
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

const removeImage = () => {
  selectedFile.value = null
  imagePreview.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}
const refreshPosts = () => {
  console.log("Gönderi başarıyla eklendi/güncellendi, listeyi güncelle!")
}

const posts = ref<Post[]>([]);
 onMounted(async () => {
   try {
     const  [err, response] = await apiService.get<any>(`post`); // Backend URL'ni ekle
     if(err) return
     posts.value = response.data;
     console.log('posts value:', posts.value)
   } catch (error) {
     console.error("Veri alınırken hata oluştu:", error);
   }

   await loadUsers()
 })

const onSubmit = async () => {
  if (!form.value.title || !form.value.content || !form.value.main_content) {
    ErrorPopup('Lütfen tüm alanları doldurun')
    return
  }

  updateLoading.value = true
  loading.value = true

  try {
    // FormData oluştur
    const formData = new FormData()

    // Form verilerini ekle
    formData.append('title', form.value.title)
    formData.append('content', form.value.content)
    formData.append('main_content', form.value.main_content)
    formData.append('user_id', String(userStore.user.id))

    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    // Form verilerini ve resmi tek istekte gönder
    const [error, resp] = await ApiService.post(
        'post/image',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        },
        undefined
    )

    if (error) {
      ErrorPopup(error)
      return
    }

    showEditModal.value = false
    SuccessPopup('Gönderi başarıyla eklendi')

    // Formu sıfırla
    form.value = {
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
    }

    removeImage()

    // Postları yeniden yükle
    const [err, response] = await apiService.get<any>(`post`);
    if(!err) {
      posts.value = response.data;
    }
  } catch (error) {
    ErrorPopup('Bir hata oluştu')
  } finally {
    updateLoading.value = false
    loading.value = false
  }
}
const rows = ref<any[]>([])

const onPageChanged = (pageData: { page: number, items: Post[] }) => {
  console.log('Sayfa değişti:', pageData.page)
  console.log('Sayfa öğeleri:', pageData.items)
}

const navigateToLogin = () => {
  router.push('auth/login')
}

const navigateToRegister = () => {
  router.push('auth/register')
}

</script>


<template>



  <!-- Site Başlığı ve Profil Menüsü -->
  <VAppBar color="white" elevation="1">
    <VContainer class="d-flex align-center pa-0">
      <!-- Sol taraf: Site Adı -->
      <div class="d-flex align-center">
        <RouterLink to="/">
          <VRow to="/" color="black">
            <VAppBarTitle class="font-weight-bold">BLOGGER.COM</VAppBarTitle>
          </VRow>
        </RouterLink>
      </div>

      <VSpacer />

      <!-- Sağ taraf: Ekle Butonu ve Profil Menüsü -->
      <div class="d-flex align-center">
        <!-- Giriş yapmamış kullanıcılar için butonlar -->
        <div v-if="!userStore.user.name" class="d-flex align-center">
          <VBtn variant="outlined" color="primary" class="me-2" @click="navigateToLogin">
            Giriş Yap
          </VBtn>
          <VBtn variant="outlined" color="primary" @click="navigateToRegister">
            Kaydol
          </VBtn>
        </div>

        <!-- Giriş yapmış kullanıcılar için Ekle butonu ve profil menüsü -->
        <template v-if="userStore.user.name">
          <VBtn variant="outlined" color="primary" class="me-3" @click="showEditModal = true">
            <VIcon start>tabler-plus</VIcon>
            Ekle
          </VBtn>

          <VMenu v-model="profileMenu" location="bottom end">
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

              <VListItem color="error" @click="logout">
                <template #prepend>
                  <VIcon icon="tabler-logout" size="small" class="me-2" color="error"/>
                </template>
                <VListItemTitle>Çıkış Yap</VListItemTitle>
              </VListItem>
            </VList>
          </VMenu>
        </template>
      </div>
    </VContainer>
  </VAppBar>


  <div class="post-list">
    <v-container>
      <Pagination :items="posts" :items-per-page="12" @page-changed="onPageChanged">
        <template #default="{ items: paginatedPosts }">
          <v-row>
            <v-col v-for="form in paginatedPosts" :key="form.id" cols="12" md="6" lg="4">
              <v-card class="mx-auto post-card" max-width="400" @click="navigateToPost(form.id)">
                <div class="image-container">
                  <v-img
                      v-if="form.image"
                      :src="`http://localhost:3001/${form.image}`"
                      :alt="form.title"
                      height="200"
                      cover
                  />
                  <v-img
                      v-else
                      src="@/assets/placeholder.jpg"
                      height="200"
                      cover
                      class="grey darken-3"
                  >
                    <div class="d-flex justify-center align-center fill-height">
                      <v-icon size="48" color="grey">tabler-photo</v-icon>
                    </div>
                  </v-img>
                </div>

                <v-card-title class="pa-2 text--primary title-container">{{ form.title }}</v-card-title>

                <v-card-text class="pa-2 text--primary content-container">
                  {{ form.content }}
                </v-card-text>

                <VCardItem class="pa-2 text--primary">
                  <VCardTitle class="user-name" @click.stop>
                    <a class="author-link" @click="navigateToAuthor(form.user_id)" target="_blank">
                      {{ getUserFullName(form.user_id) }}
                    </a>
                  </VCardTitle>
                  <VCardSubtitle class="post-date">
                    {{ tarihFormat(form.created_at) }}
                  </VCardSubtitle>
                </VCardItem>

                <v-card-actions @click.stop>
                  <VBtn
                      variant="text"
                      prepend-icon="tabler-heart"
                      :color="isLikedByUser(form.likes) ? 'error' : 'default'"
                      @click="handleLike(form.id, form)"
                  >
                    {{ form.like_count }}
                  </VBtn>
                  <VBtn
                      variant="text"
                      prepend-icon="tabler-message-circle"
                      @click="navigateToPost(form.id)"
                      color="default"
                  >
                    {{ form.comment_count }}
                  </VBtn>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </template>
      </Pagination>
    </v-container>
  </div>

  <VDialog v-model="showEditModal" max-width="600px">
    <VCard>
      <VCardTitle class="text-h5 bg-primary text-black pa-4">
        Gönderi Ekle
      </VCardTitle>

      <VCardText class="pa-4">
        <VForm @submit.prevent="onSubmit">
          <VRow>
            <!-- Resim yükleme alanı -->
            <VCol cols="12">
              <div
                  class="upload-area"
                  :class="{ 'dragging': isDragging }"
                  @click="triggerFileInput"
                  @dragover.prevent="isDragging = true"
                  @dragleave.prevent="isDragging = false"
                  @drop.prevent="handleDrop"
                  v-if="!imagePreview"
              >
                <div class="text-center">
                  <VIcon icon="tabler-upload" size="48" class="mb-2" />
                  <div>Resim yüklemek için tıklayın veya sürükleyin</div>
                  <div class="text-caption">Maksimum dosya boyutu: 5MB</div>
                </div>
                <input
                    type="file"
                    ref="fileInput"
                    style="display: none"
                    accept="image/*"
                    @change="handleFileSelect"
                />
              </div>

              <!-- Resim önizleme -->
              <div class="preview-container" v-if="imagePreview">
                <img :src="imagePreview" class="preview-image" />
                <VBtn
                    icon
                    class="remove-btn"
                    size="small"
                    color="error"
                    @click.stop="removeImage"
                >
                  <VIcon>tabler-x</VIcon>
                </VBtn>
              </div>
            </VCol>

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
              Ekle
            </VBtn>
          </VCardActions>
        </VForm>
      </VCardText>
    </VCard>
  </VDialog>
</template>




<style scoped>
.posts-container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
}

.user-name {
  font-size: 0.8rem;  /* Smaller font size */
}

.post-date {
  font-size: 0.7rem;  /* Even smaller font size */
  color: #6c757d;
}

.profile-btn {
  border-radius: 24px;
  padding-left: 12px;
  padding-right: 12px;
}

.upload-area {
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
}

.preview-image {
  width: 100%;
  height: auto;
  object-fit: contain;
  max-height: 300px;
}

.remove-btn {
  position: absolute;
  top: 8px;
  right: 8px;
}
.post-list {
  padding: 80px;
}

.user-info {
  padding: 3px 15px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-link {
  color: #555;
  font-size: 13px;
  text-decoration: none;
  position: relative;
  cursor: pointer;
}

.author-link:hover {
  text-decoration: underline;
}

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

.post-card {
  height: 400px;
  display: flex;
  flex-direction: column;
}

.image-container {
  height: 350px;
  overflow: hidden;
}

.title-container {
  height: 60px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  font-size: 1.1rem;
  line-height: 1.4;
}

.content-container {
  height: 100px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  font-size: 0.9rem;
  line-height: 1.4;
  color: #666;
}

</style>

<route lang="yaml">
meta:
  auth: false
  layout: blank
</route>

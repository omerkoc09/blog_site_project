<script setup lang="ts">
import tarihFormat from '@/utils/ExDate'
import { useUserStore } from '@/store/user'
import { ErrorPopup } from '@/utils/Popup'
import ApiService from '@/services/ApiService'



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


</script>

<template>
  <div class="posts-container">
    <!-- Loading state -->
    <div v-if="loading" class="d-flex justify-center my-4">
      <VProgressCircular indeterminate />
    </div>

    <!-- Post content -->
    <template v-else>
      <VCard class="mb-6">
        <!-- Post title -->
        <VCardTitle class="text-h3 px-4 text-center text-wrap" :style="{
          wordBreak: 'break-word',
          whiteSpace: 'normal',
          lineHeight: '1.2',  // satır yüksekliğini ayarlar
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
              {{ getUserFullName(form.user_id) }}
            </VCardTitle>
            <VCardSubtitle class="post-date">
              {{ tarihFormat(form.created_at) }}
            </VCardSubtitle>
          </div>
        </VCardItem>

        <!-- Like and comment buttons -->
        <VCardActions>
          <VBtn variant="text" prepend-icon="tabler-heart">
            {{ form.like_count }} Beğeni
          </VBtn>

          <VBtn variant="text" prepend-icon="tabler-message-circle"
            @click="showCommentInput = showCommentInput === form.id ? null : form.id">
            {{ form.comment_count }} Yorum
          </VBtn>
          <VSpacer />
        </VCardActions>
        <!-- Comments section -->
        <VExpandTransition>
          <div class="pa-4 bg-grey-lighten-4">
            <!-- Existing comments -->
            <div v-if="form.comments.length" class="mb-4">
              <div v-for="comment in form.comments" :key="comment.id" class="mb-3">
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
                      <div v-if="editingCommentId === comment.id">
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
                  <div v-if="comment.user_id === userStore.user.id || userStore.getRole() === 10" class="d-flex gap-2">
                    <VBtn icon="tabler-edit" size="x-small" variant="text" @click="startEdit(comment)" />
                    <VBtn icon="tabler-trash" size="x-small" variant="text" color="error"
                      @click="deleteComment(comment.id, form)" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </VExpandTransition>

        <VCardText class="text-body-1 text-left" :style="{
          fontSize: '1.5rem !important',
          lineHeight: '1.5',
          color: 'black',
          whiteSpace: 'pre-wrap',  // Tüm boşlukları ve satır sonlarını korur
          fontWeight: '400',       // Normal yazı kalınlığı
          letterSpacing: '0.5px',  // Harfler arası mesafe
          wordSpacing: '2px',      // Kelimeler arası mesafe
          textTransform: 'none',   // Metni olduğu gibi gösterir
          padding: '1rem',         // Metin etrafında boşluk
          margin: '1rem 0',        // Üst ve altında boşluk
          maxWidth: '100%',        // Container'a göre maksimum genişlik
          overflowWrap: 'break-word'  // Uzun kelimeleri alt satıra geçirir
        }">
          {{ form.main_content }}
        </VCardText>

        <VDivider />
      </VCard>
    </template>
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
</style>

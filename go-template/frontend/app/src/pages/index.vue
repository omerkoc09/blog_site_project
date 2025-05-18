<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useUserStore } from '@/store/user'
import tarihFormat from '@/utils/ExDate'
import ApiService from '@/services/ApiService'
import { ErrorPopup, WarningPopup, SuccessPopup } from '@/utils/Popup'
import { router } from '@/plugins/1.router'
import { useRoute } from 'vue-router'
import apiService from "@/services/ApiService";
import Pagination from '@/components/Pagination.vue'
import { useFollow } from '@/composables/useFollowService'
const { followers, fetchFollowings } = useFollow()
const route = useRoute()

// Prop olarak modal durumunu alabiliriz (layout tarafında yönetiliyor)
const props = defineProps({
  showAddModal: {
    type: Boolean,
    default: false
  }
})

// Emit tanımlayalım ki parent komponenti güncelleyebilelim
const emit = defineEmits(['update:showAddModal'])

// showEditModal değişkenini prop ile senkronize edelim
const showEditModal = computed({
  get: () => props.showAddModal,
  set: (value) => emit('update:showAddModal', value)
})

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
  topics?: {id: number, name: string}[]  // Topic bilgisini ekleyelim
}

// Interface ekleyelim
interface MyUser {
  id: number
  name: string
  surname: string
  role: number
}

const userStore = useUserStore()
const users = ref<MyUser[]>([])

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


const navigateToPost = (postId: number) => {
  router.push(`third_page/${postId}`)
}

const navigateToAuthor = (userId: number) => {
  router.push(`author/${userId}`)
}

const isLikedByUser = (likes: Like[]) => {
  if (!likes || !Array.isArray(likes) || !userStore.isAuthenticated) {
    return false;
  }
  return likes.some((like: Like) => like && like.user_id === userStore.user.id);
}

const handleLike = async (postId: number, post: any) => {
  if (!userStore.isAuthenticated) {
    ErrorPopup('Beğenmek için giriş yapmalısınız');
    return;
  }
  
  const [error] = await ApiService.post(`like/${postId}`, {})

  if (!error) {
    // Initialize post.likes if it doesn't exist
    if (!post.likes || !Array.isArray(post.likes)) {
      post.likes = [];
    }
    
    // Yerel state'i güncelle
    if (isLikedByUser(post.likes)) {
      // Eğer kullanıcı zaten beğenmişse, beğeniyi kaldır
      post.likes = post.likes.filter((like: { user_id: any }) => like.user_id !== userStore.user.id)
      post.like_count = (post.like_count || 0) - 1;
    } else {
      // Beğeni ekle
      post.likes.push({ user_id: userStore.user.id })
      post.like_count = (post.like_count || 0) + 1;
    }
  } else {
    ErrorPopup('Beğeni işlemi sırasında bir hata oluştu: ' + error);
  }
}

const handleShowAddModal = () => {
  showEditModal.value = true
}

// Kullanıcıları yükleyecek fonksiyon
const loadUsers = async () => {
  try {
    const [error, response] = await ApiService.get<MyUser[]>('user');
    if (!error && response?.data) {
      users.value = response.data || [];
    } else if (error) {
      console.error('Error loading users:', error);
      users.value = [];
    }
  } catch (err) {
    console.error('Exception loading users:', err);
    users.value = [];
  }
  
  console.log('users', users.value);
}

// Post veya yorum sahibinin adını getiren fonksiyon
const getUserFullName = (userId: number) => {
  if (!users.value || !Array.isArray(users.value)) {
    return 'Bilinmeyen Kullanıcı';
  }
  
  const user = users.value.find(u => u && u.id === userId);
  return user ? `${user.name || ''} ${user.surname || ''}`.trim() : 'Bilinmeyen Kullanıcı';
}


// Resim yükleme için gerekli değişkenler
const fileInput = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const imagePreview = ref<string | null>(null)
const selectedFile = ref<File | null>(null)
const loading = ref(false)
const updateLoading = ref(false)
const tab = ref<number>(1)

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
const posts = ref<Post[]>([]);
const savedPosts = ref<Post[]>([]);
const postSavedStatus = ref<{[key: number]: boolean}>({});

// Topic için gerekli state değişkenlerini ekleyelim
const selectedTopicId = ref<number | null>(null)
const topics = ref<{ id: number, name: string }[]>([])
const selectedTopicIds = ref<number[]>([])

// Topic yükleme fonksiyonu
const loadTopics = async () => {
  try {
    const [error, response] = await ApiService.get<any>('topic')
    if (!error && response?.data) {
      topics.value = response.data || response
    } else if (error) {
      console.error('Error loading topics:', error)
      topics.value = []
    }
  } catch (err) {
    console.error('Exception loading topics:', err)
    topics.value = []
  }
}

// Post'ların topic'lerini yükle
const loadPostTopics = async () => {
  console.log('load post topics CALISTI')
  console.log('posts.value', posts.value)
  for (const post of posts.value) {
    try {
      const [error, response] = await ApiService.get<any>(`post/${post.id}/topics`)
      if (!error && response) {
        // API'nın response formatına göre doğru veriyi alın
        const topicsData = response.data || response
        console.log('topic data', topicsData)
        post.topics = topicsData
        
        console.log('Topics:')
        if (Array.isArray(topicsData)) {
          topicsData.forEach((topic: any, index: number) => {
            console.log(`Topic ${index + 1}: ID=${topic.id}, Name=${topic.name}`)
          })
        }
      }
    } catch (err) {
      console.error(`Error loading topics for post ${post.id}:`, err)
      post.topics = []
    }
  }
}

// Topic'leri temizle ve filtreyi kaldır
const clearTopicFilter = () => {
  selectedTopicId.value = null
}

const filteredPosts = computed(() => {
  if (!posts.value) {
    return tab.value === 4 ? savedPosts.value || [] : [];
  }

  for (const post of posts.value) {
    console.log('post.topics', post.topics)
  }
   // Kullanıcının ilgi alanı ile eşleşenler
   const interestedPosts = posts.value.filter(post =>
    post.topics?.some((topic) => userStore.user?.interests?.includes(topic.id))
  );
  
  // Eşleşmeyenler
  const otherPosts = posts.value.filter(post =>
    !post.topics?.some(topic => userStore.user?.interests?.includes(topic.id))
  );
  
  let result: any[] = [];
  
  if(tab.value == 1) {
    result = [...interestedPosts, ...otherPosts];
  }
  else if(tab.value == 2) {
    result = userStore.isAuthenticated && userStore.user?.id ? 
      posts.value.filter(post => post.user_id == userStore.user.id) : [];
  }
  else if(tab.value == 3) {
    result = posts.value.filter(post => followers.value?.includes(post.user_id) || false);
  }
  else if(tab.value == 4) {
    result = savedPosts.value || [];
  }

  
  // Ardından topic filtresi uygulayalım (eğer bir topic seçildiyse)
  if (selectedTopicId.value) {
    console.log('SELECTED TOPIC ID', selectedTopicId.value)
    result = result.filter(post => 
      post.topics?.some((topic: any) => topic.id === selectedTopicId.value)
    )
  }
  
  return result;
})

const isSavedByUser = (postId: number) => {
  return !!postSavedStatus.value[postId];
}

const handleSave = async (postId: number, event: MouseEvent) => {
  // Stop event propagation to prevent navigating to post detail
  event.stopPropagation();
  
  if (!userStore.isAuthenticated) {
    ErrorPopup('Kaydetmek için giriş yapmalısınız');
    return;
  }
  
  const [error] = await ApiService.post(`saved/${postId}`, {});
  
  if (!error) {
    // Update local state
    postSavedStatus.value[postId] = !postSavedStatus.value[postId];
    
    // If we're in saved tab, we might need to refresh the list
    if (tab.value === 4 && !postSavedStatus.value[postId]) {
      loadSavedPosts();
    }
  } else {
    ErrorPopup('Bir hata oluştu: ' + error);
  }
}

const loadSavedPosts = async () => {
  if (!userStore.isAuthenticated) {
    savedPosts.value = [];
    return;
  }
  
  try {
    const [error, response] = await ApiService.get<any>('saved');
    if (!error && response?.data) {
      savedPosts.value = response.data || [];
      
      // Update the saved status for each post
      savedPosts.value.forEach(post => {
        if (post && post.id) {
          postSavedStatus.value[post.id] = true;
        }
      });
    } else if (error) {
      console.error('Error loading saved posts:', error);
      ErrorPopup('Kaydedilen gönderiler yüklenirken hata oluştu: ' + error);
      savedPosts.value = [];
    }
  } catch (err) {
    console.error('Exception loading saved posts:', err);
    savedPosts.value = [];
  }
}

onMounted(async () => {
  try {
    const [err, response] = await apiService.get<any>(`post`);
    if (err) {
      console.error("Error loading posts:", err);
      posts.value = [];
      return;
    }
    
    if (response?.data) {
      posts.value = response.data;
      console.log('Posts loaded:', posts.value?.length || 0);
      // Post'ların topic'lerini yükleyelim
      await loadPostTopics();
    } else {
      console.warn("No posts data returned from API");
      posts.value = [];
    }
  } catch (error) {
    console.error("Veri alınırken hata oluştu:", error);
    posts.value = [];
  }

  await loadUsers();
  // Topic'leri yükleyelim
  await loadTopics();
  
  if (userStore.isAuthenticated && userStore.user?.id) {
    console.log('Kullanıcı giriş yapmış, takip edilenleri yüklüyorum...')
    await fetchFollowings(userStore.user.id)
    await loadSavedPosts()
  } else {
    console.log('Kullanıcı giriş yapmamış, takip edilenler yüklenmiyor')
  }

  if (route.query.topic) {
    selectedTopicId.value = Number(route.query.topic)
  }
})

watch(tab, async (newValue) => {
  console.log('Tab değişti:', newValue)
  
  if (newValue === 3) {
    if (userStore.isAuthenticated && userStore.user?.id) {
      console.log('Following tabına geçildi, takip edilenleri yüklüyorum...')
      await fetchFollowings(userStore.user.id)
    } else {
      console.warn('Kullanıcı giriş yapmamış, takip edilenler yüklenemiyor')
      ErrorPopup('Bu özelliği kullanmak için giriş yapmalısınız')
      tab.value = 1 // Redirect to Explore tab
    }
  } else if (newValue === 4) {
    if (userStore.isAuthenticated && userStore.user?.id) {
      console.log('Kaydedilenler tabına geçildi, kaydedilen gönderileri yüklüyorum...')
      await loadSavedPosts()
    } else {
      console.warn('Kullanıcı giriş yapmamış, kaydedilen gönderiler yüklenemiyor')
      ErrorPopup('Bu özelliği kullanmak için giriş yapmalısınız')
      tab.value = 1 // Redirect to Explore tab
    }
  }
})

watch(() => userStore.isAuthenticated, async (isLoggedIn) => {
  if (isLoggedIn && userStore.user?.id) {
    console.log('Kullanıcı girişi tespit edildi, veri yükleniyor...')
    if (tab.value === 3) {
      await fetchFollowings(userStore.user.id)
    } else if (tab.value === 4) {
      await loadSavedPosts()
    }
  } else {
    savedPosts.value = []
    postSavedStatus.value = {}
  }
})

const onSubmit = async () => {
  if (!form.value.title || !form.value.content || !form.value.main_content) {
    ErrorPopup('Lütfen tüm alanları doldurun')
    return
  }

  updateLoading.value = true
  loading.value = true

  try {
    const formData = new FormData()

    formData.append('title', form.value.title)
    formData.append('content', form.value.content)
    formData.append('main_content', form.value.main_content)
    formData.append('user_id', String(userStore.user.id))
    if (selectedTopicIds.value.length > 0) {
  formData.append('topic_ids', selectedTopicIds.value.join(','))
}

    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    const [error, resp] = await ApiService.post(
        'post/image',
        formData
    )

    if (error) {
      ErrorPopup(error)
      return
    }

    showEditModal.value = false
    SuccessPopup('Gönderi başarıyla eklendi')

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

const onPageChanged = (pageData: { page: number, items: Post[] }) => {
  console.log('Sayfa değişti:', pageData.page)
  console.log('Sayfa öğeleri:', pageData.items)
}


</script>


<template>
  <div>
    <v-tabs
        v-model="tab"
        align-tabs="center"
        color="deep-purple-accent-4"
    >
      <v-tab :value="1">For you</v-tab>
      <v-tab :value="2">Your Posts</v-tab>
      <v-tab :value="3">Following</v-tab>
      <v-tab :value="4">Saved</v-tab>
    </v-tabs>
    
    <!-- Topic filtreleme alanı ekleyelim -->
    <div class="topic-filter-container">
      <v-select
        v-model="selectedTopicId"
        :items="topics"
        item-title="name"
        item-value="id"
        label="Filter by Topic"
        clearable
        @click:clear="clearTopicFilter"
        class="topic-filter"
      >
        <template v-slot:prepend>
          <v-icon>tabler-tag</v-icon>
        </template>
      </v-select>
    </div>

    <div class="post-list">
      <v-container>
        <Pagination :items="filteredPosts" :items-per-page="12" @page-changed="onPageChanged">
          <template #default="{ items: paginatedPosts }">
            <v-row>
              <v-col v-for="post in paginatedPosts as Post[]" :key="post.id" cols="12" md="6" lg="4">
                <v-card class="mx-auto post-card" max-width="400" @click="navigateToPost(post.id)">
                  <div class="image-container">
                    <v-img
                        v-if="post.image"
                        :src="`http://localhost:3001/${post.image}`"
                        :alt="post.title"
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

                  <v-card-title class="pa-2 text--primary title-container">{{ post.title }}</v-card-title>

                  <v-card-text class="pa-2 text--primary content-container">
                    {{ post.content }}
                  </v-card-text>

                  <VCardItem class="pa-2 text--primary">
                    <VCardTitle class="user-name" @click.stop>
                      <a class="author-link" @click="navigateToAuthor(post.user_id)" target="_blank">
                        {{ getUserFullName(post.user_id) }}
                      </a>
                    </VCardTitle>
                    <VCardSubtitle class="post-date">
                      {{ tarihFormat(post.created_at) }}
                    </VCardSubtitle>
                  </VCardItem>

                  <v-card-actions @click.stop>
                    <VBtn
                        variant="text"
                        prepend-icon="tabler-heart"
                        :color="isLikedByUser(post.likes) ? 'error' : 'default'"
                        @click="handleLike(post.id, post)"
                    >
                      
                    </VBtn>
                    <VBtn
                        variant="text"
                        prepend-icon="tabler-message-circle"
                        @click="navigateToPost(post.id)"
                        color="default"
                    >
                      {{ post.comment_count }}
                    </VBtn>
                    <VSpacer />
                    <VBtn
                        variant="text"
                        :prepend-icon="isSavedByUser(post.id) ? 'tabler-bookmark-filled' : 'tabler-bookmark'"
                        :color="isSavedByUser(post.id) ? 'primary' : 'default'"
                        @click="(e: MouseEvent) => handleSave(post.id, e)"
                    />
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </Pagination>
      </v-container>
    </div>
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

            <VCol cols="12">
  <v-select
    v-model="selectedTopicIds"
    :items="topics"
    item-title="name"
    item-value="id"
    label="Topic Seç"
    multiple
    chips
    clearable
    :loading="topics.length === 0"
  />
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
  padding: 20px;
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
  line-clamp: 2;
  -webkit-box-orient: vertical;
  font-size: 1.1rem;
  line-height: 1.4;
}

.content-container {
  height: 100px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  line-clamp: 4;
  -webkit-box-orient: vertical;
  font-size: 0.9rem;
  line-height: 1.4;
  color: #666;
}

.topic-filter-container {
  padding: 0 20px;
  max-width: 300px;
  margin: 10px auto;
}

.topic-filter {
  margin-bottom: 10px;
}

</style>

<route lang="yaml">
meta:
  auth: false
  layout: default
</route>

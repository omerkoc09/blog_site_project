<script setup lang="ts">
import apiService from "@/services/ApiService";
import {ref} from "vue";
import {useUserStore} from "@/store/user";
import {router} from "@/plugins/1.router";
import {ErrorPopup} from "@/utils/Popup";
import * as https from "node:https";

interface Post {
  id: number;
  title: string;
  content: string;
  main_content: string;
  created_at: string;
}

const form = ref({
  name: '',
  surname: '',
  email: '',
  phone: '',
  about: '',
  posts: [] as Post[],
})

const userStore = useUserStore()
const profileMenu = ref(false)
const route = useRoute()
const tab = ref(0)

const logout =  () => {
  try {
    useUserStore().logout()
    router.push('/')
  } catch (error) {
    ErrorPopup('Çıkış yapılırken bir hata oluştu')
  }
}

const navigateToPost = (postId: number) => {
  router.push(`/third_page/${postId}`)
}

const userId = computed(() => (route.params as { id: string }).id)
console.log('Post ID:', userId.value)

const navigateToLogin = () => {
  router.push('../auth/login')
}

const navigateToRegister = () => {
  router.push('../auth/register')
}


onMounted(async () => {
  try {
    const [err, response] = await apiService.get<any>(`user_guest/${userId.value}`);
    if(err) return
    form.value = response.data;
    console.log('User data:', form.value)
  } catch (error) {
    console.error("Veri alınırken hata oluştu:", error);
  }
})
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

  <div class="author-container">
    <!-- Author Header -->
    <div class="author-header">
      <h1 class="author-name">{{ form.name }} {{ form.surname }}</h1>
    </div>

    <!-- Tabs -->
    <VTabs v-model="tab" color="primary" grow>
      <VTab value="0">
        <VIcon start icon="tabler-article"/>
        Gönderiler
      </VTab>
      <VTab value="1">
        <VIcon start icon="tabler-user"/>
        Profil Bilgileri
      </VTab>
    </VTabs>

    <VWindow v-model="tab" class="mt-6">
      <!-- Posts Tab -->
      <VWindowItem value="0">
        <div v-if="form.posts.length === 0" class="text-center pa-4">
          <VIcon icon="tabler-article-off" size="48" color="grey"/>
          <p class="text-h6 text-grey mt-2">Henüz gönderi bulunmamaktadır.</p>
        </div>
        
        <div v-else class="posts-grid">
          <VCard
            v-for="post in form.posts"
            :key="post.id"
            class="post-card"
            @click="navigateToPost(post.id)"
          >
            <VCardTitle class="text-h5">{{ post.title }}</VCardTitle>
            <VCardText>
              {{ post.content }}
            </VCardText>
            <VCardActions>
              <VSpacer />
              <VBtn
                color="primary"
                variant="text"
                @click.stop="navigateToPost(post.id)"
              >
                Devamını Oku
                <VIcon end icon="tabler-chevron-right"/>
              </VBtn>
            </VCardActions>
          </VCard>
        </div>
      </VWindowItem>

      <!-- Profile Info Tab -->
      <VWindowItem value="1">
        <VCard class="profile-info-card">
          <VCardText>
            <div class="info-section">
              <VIcon icon="tabler-mail" color="primary"/>
              <div class="info-content">
                <h3>E-posta</h3>
                <p>{{ form.email }}</p>
              </div>
            </div>

            <VDivider class="my-4"/>

            <div class="info-section">
              <VIcon icon="tabler-phone" color="primary"/>
              <div class="info-content">
                <h3>Telefon</h3>
                <p>{{ form.phone }}</p>
              </div>
            </div>

            <VDivider class="my-4"/>

            <div class="info-section">
              <VIcon icon="tabler-user" color="primary"/>
              <div class="info-content">
                <h3>Hakkında</h3>
                <p>{{ form.about }}</p>
              </div>
            </div>
          </VCardText>
        </VCard>
      </VWindowItem>
    </VWindow>
  </div>
</template>

<style scoped>
.author-container {
  max-width: 1200px;
  margin: 80px auto 0;
  padding: 24px;
}

.author-header {
  text-align: center;
  margin-bottom: 32px;
}

.author-name {
  font-size: 2.5rem;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
  padding: 16px 0;
}

.post-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.info-section {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 16px;
}

.info-content {
  flex: 1;
}

.info-content h3 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #666;
  margin: 0 0 8px 0;
}

.info-content p {
  font-size: 1.2rem;
  color: #333;
  margin: 0;
  line-height: 1.5;
}

.profile-info-card {
  max-width: 800px;
  margin: 0 auto;
}
</style>

<route lang="yaml">
meta:
  auth: false
  layout: blank
</route>
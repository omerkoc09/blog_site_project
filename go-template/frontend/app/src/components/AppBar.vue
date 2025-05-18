<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { router } from '@/plugins/1.router'
import { Notification } from '@/services/notification'
import ApiService from '@/services/ApiService';

// Props tanımlama
interface AppBarProps {
  showAddButton?: boolean; // Opsiyonel prop, varsayılan değeri aşağıda belirlenecek
}

// Prop'ları tanımla ve varsayılan değerlerini belirle
const props = withDefaults(defineProps<AppBarProps>(), {
  showAddButton: true, // Varsayılan olarak "Ekle" butonu gösterilsin
})

const userStore = useUserStore()
const profileMenu = ref(false)
const notifications = ref<Notification[]>([])
const notificationMenu = ref(false)
const loadingNotifications = ref(false)

const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

const openAddModal = () => {
  emit('showAddModal')
}

async function loadNotifications() {
  loadingNotifications.value = true
  const [error, response] = await ApiService.get<Notification[]>('notifications')
  if (!error && response) {
    notifications.value = response.data || [] || response
  }
  loadingNotifications.value = false
}

async function handleNotificationClick(notification: Notification) {
  if (!notification.is_read) {
    const [error, response] = await ApiService.put<Notification>(`notifications/${notification.id}/read`)
    if (!error && response) {
      notification.is_read = true
    }
  }
  if (notification.post_id) {
    router.push(`/third_page/${notification.post_id}`)
  }
  if (notification.follow_id) {
    router.push(`/author/${notification.sender_id}`)
  }
  
  notificationMenu.value = false
}

onMounted(() => {
  if (userStore.user.name) {
    loadNotifications()
  }
})

// Çıkış yapma fonksiyonu
const logout = () => {
  try {
    useUserStore().logout()
    router.push('/')
  } catch (error) {
    ErrorPopup('Çıkış yapılırken bir hata oluştu')
  }
}

const navigateToLogin = () => {
  router.push('auth/login')
}

const navigateToRegister = () => {
  router.push('auth/register')
}

// Yeni post/gönderi eklemeyi tetikleyen emit
const emit = defineEmits(['showAddModal'])
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
          <!-- showAddButton prop'u true ise Ekle butonunu göster -->
          <VBtn 
            v-if="showAddButton" 
            variant="outlined" 
            color="primary" 
            class="me-3" 
            @click="openAddModal"
          >
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

          <!-- Bildirim ikonu -->
          <VMenu v-model="notificationMenu" location="bottom end">
            <template v-slot:activator="{ props }">
              <VBtn icon v-bind="props" @click="loadNotifications">
                <VBadge :content="unreadCount" color="red" v-if="unreadCount > 0">
                  <VIcon>tabler-bell</VIcon>
                </VBadge>
                <VIcon v-else>tabler-bell</VIcon>
              </VBtn>
            </template>
            <VList width="350" style="max-height:400px;overflow:auto;">
              <VListItem v-if="loadingNotifications">
                <VListItemTitle>Yükleniyor...</VListItemTitle>
              </VListItem>
              <VListItem v-for="n in notifications" :key="n.id" @click="handleNotificationClick(n)" :class="{'bg-grey-lighten-4': !n.is_read}">
                <VListItemTitle>
                  <span v-if="n.type === 1">Birisi gönderini beğendi.</span>
                  <span v-else-if="n.type === 2">Birisi gönderine yorum yaptı.</span>
                  <span v-else-if="n.type === 3">Biri seni takip etti.</span>
                  <span v-else-if="n.type === 4">Takip ettiğin biri gönderi paylaştı.</span>
                  <span v-else>Yeni bildirim.</span>
                </VListItemTitle>
                <VListItemSubtitle v-if="n.created_at">{{ new Date(n.created_at).toLocaleString('tr-TR') }}</VListItemSubtitle>
              </VListItem>
              <VListItem v-if="!loadingNotifications && notifications.length === 0">
                <VListItemTitle>Hiç bildirimin yok.</VListItemTitle>
              </VListItem>
            </VList>
          </VMenu>
        </template>
      </div>
    </VContainer>
  </VAppBar>
</template>

<style scoped>
.profile-btn {
  border-radius: 24px;
  padding-left: 12px;
  padding-right: 12px;
}
</style>
<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/store/user'
import { router } from '@/plugins/1.router'

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

const openAddModal = () => {
  emit('showAddModal')
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
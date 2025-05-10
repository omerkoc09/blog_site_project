<script lang="ts" setup>
// Sadece AppBar'ı import ediyoruz, nav yapıları kaldırıldı
import AppBar from '@/components/AppBar.vue'
import { ref } from 'vue'
import { useUserStore } from '@/store/user'
import { ErrorPopup, SuccessPopup } from '@/utils/Popup'
import apiService from '@/services/ApiService'

// Modal kontrolü için
const showAddModal = ref(false)

// Post ekleme için gerekli değişkenler
const userStore = useUserStore()
const fileInput = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const imagePreview = ref<string | null>(null)
const selectedFile = ref<File | null>(null)
const loading = ref(false)
const updateLoading = ref(false)

// Formlar için temel yapı
const form = ref({
  title: '',
  content: '',
  main_content: '',
  user_id: 0
})

// Modal açma fonksiyonu
const handleShowAddModal = () => {
  showAddModal.value = true
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

// Form submit
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
    const [error, resp] = await apiService.post(
        'post/image',
        formData
    )

    if (error) {
      ErrorPopup(error)
      return
    }

    showAddModal.value = false
    SuccessPopup('Gönderi başarıyla eklendi')

    // Formu sıfırla
    form.value = {
      title: '',
      content: '',
      main_content: '',
      user_id: 0
    }

    removeImage()

    // Event fonksiyonu ekleyin, sayfa yenilemesi için
    window.dispatchEvent(new CustomEvent('post-added'))
  } catch (error) {
    ErrorPopup('Bir hata oluştu')
  } finally {
    updateLoading.value = false
    loading.value = false
  }
}
</script>

<template>
  <div class="app-layout">
    <!-- AppBar eklendi -->
    <AppBar @showAddModal="handleShowAddModal" />
    
    <!-- İçerik alanı -->
    <div class="content-area">
      <RouterView v-slot="{ Component }">
        <Suspense>
          <Component :is="Component" :showAddModal="showAddModal" @update:showAddModal="showAddModal = $event" />
        </Suspense>
      </RouterView>
    </div>

    <!-- Post ekleme modalı tüm sayfalarda erişilebilir -->
    <VDialog v-model="showAddModal" max-width="600px">
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
              <VBtn variant="text" @click="showAddModal = false">
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
  </div>
</template>

<style lang="scss">
// Basit layout stili
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  width: 100%;
}

.content-area {
  flex: 1;
  padding-top: 64px; /* AppBar yüksekliği kadar boşluk */
}

// Modal stil tanımlamaları
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
</style>

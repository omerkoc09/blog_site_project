<template>
    <div class="image-upload-container">
      <VCard class="upload-card">
        <VCardTitle class="text-center">Resim Yükle</VCardTitle>
        <VCardText class="text-center">
          <VFileInput
            v-model="selectedImage"
            label="Resim Seç"
            accept="image/*"
            show-size
            variant="outlined"
            density="compact"
            prepend-icon="mdi-camera"
            class="custom-file-input"
            @update:model-value="onFileSelect"
          />
        </VCardText>
  
        <!-- Mevcut resim veya yeni seçilen resmin önizlemesi -->
        <VCardText v-if="currentImage || imagePreview" class="text-center">
          <img
            :src="imagePreview || (currentImage ? 'http://localhost:3001/' + currentImage : '')"
            class="image-preview"
            alt="Post Resmi"
          />
        </VCardText>

        <VCardActions class="d-flex justify-center">
          <VBtn color="error" variant="outlined" @click="resetImage">
            Sıfırla
          </VBtn>
        </VCardActions>
      </VCard>
    </div>
  </template>

  <script setup lang="ts">
  import { ref, watch } from 'vue'
  import { ErrorPopup } from '@/utils/Popup'

  const props = defineProps<{
    postId: number,
    currentImage?: string
  }>()

  const emit = defineEmits<{
    (e: 'update:selectedFile', file: File | null): void
  }>()

  const selectedImage = ref<File[]>([])
  const imagePreview = ref<string | null>(null)

  // Dosya seçildiğinde
  const onFileSelect = (files: File[]) => {
    if (!files || files.length === 0) {
      resetImage()
      return
    }

    const file = files[0]

    // Dosya boyutu kontrolü (5MB)
    if (file.size > 5 * 1024 * 1024) {
      ErrorPopup('Dosya boyutu 5MB\'dan küçük olmalıdır')
      resetImage()
      return
    }

    // Önizleme oluştur
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)

    // Seçilen dosyayı parent komponente ilet
    emit('update:selectedFile', file)
  }

  const resetImage = () => {
    selectedImage.value = []
    imagePreview.value = null
    emit('update:selectedFile', null)
  }
  </script>

<style scoped>
.image-upload-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  text-align: center;
  padding: 20px;
  cursor: pointer;
  width: 1000px;
}

.upload-card {
  width: 350px;
  padding: 15px;
  border-radius: 12px;
}

.custom-file-input {
  margin-bottom: 10px;
}

.image-preview {
  width: 100%;
  max-height: 200px;
  object-fit: contain;
  border-radius: 8px;
  margin-top: 10px;
}
</style>
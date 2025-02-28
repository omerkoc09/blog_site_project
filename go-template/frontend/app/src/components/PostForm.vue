<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import ApiService from '@/services/ApiService'
import { useUserStore } from '@/store/user'
import { ErrorPopup, SuccessPopup } from '@/utils/Popup'


const userStore = useUserStore()

const props = defineProps({
  iscreateform: Boolean,
    formloading: Boolean,
})

const emit = defineEmits(['closemodal', 'postSuccess'])

const form = ref({
  title: '',
  content: '',
  main_content: '',
  image: null,
  user_id: userStore.user.id
})

const fileInput = ref(null)
const selectedFile = ref(null)
const imagePreview = ref(null)
const isDragging = ref(false)

const triggerFileInput = () => {
  fileInput.value.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
    form.value.image = file
    imagePreview.value = URL.createObjectURL(file)
  }
}

const handleDrop = (event) => {
  isDragging.value = false
  const file = event.dataTransfer.files[0]
  if (file) {
    selectedFile.value = file
    form.value.image = file
    imagePreview.value = URL.createObjectURL(file)
  }
}

const removeImage = () => {
  selectedFile.value = null
  form.value.image = null
  imagePreview.value = null
}

const submitPost = async () => {
  if (!form.value.title || !form.value.content || !form.value.main_content) {
    ErrorPopup('Tüm alanları doldurun!')
    return
  }

  try {
    const formData = new FormData()
    formData.append('title', form.value.title)
    formData.append('content', form.value.content)
    formData.append('main_content', form.value.main_content)
    formData.append('user_id', form.value.user_id)

    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    const endpoint = props.iscreateform ? 'post' : `post/${form.value.id}`
    const method = props.iscreateform ? 'post' : 'put'

    const [error, response] = await ApiService.post('post/image', formData)

    if (error) {
      ErrorPopup(error)
      return
    }

    SuccessPopup(props.iscreateform ? 'Gönderi başarıyla paylaşıldı!' : 'Gönderi güncellendi!')

    emit('postSuccess')
    emit('closemodal')

  } catch (err) {
    ErrorPopup('Bir hata oluştu!')
  }
}
</script>


<template>
    <VForm @submit.prevent="submitPost">
      <VRow>
        <!-- Resim yükleme alanı -->
        <VCol cols="12">
          <VCard
            class="upload-area pa-4"
            :class="{ 'dragover': isDragging }"
            @dragenter.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @dragover.prevent
            @drop.prevent="handleDrop"
            @click="triggerFileInput"
          >
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              class="d-none"
              @change="handleFileSelect"
            />
            <div v-if="!imagePreview" class="text-center">
              <VIcon icon="tabler-upload" size="48" color="primary" class="mb-2" />
              <div class="text-h6">Resim Yükle</div>
              <div class="text-body-2 text-medium-emphasis">
                Sürükle bırak veya tıklayarak seç
              </div>
            </div>
            <div v-else class="preview-container">
              <img :src="imagePreview" class="preview-image">
              <VBtn
                icon="tabler-x"
                size="x-small"
                variant="text"
                color="error"
                class="remove-btn"
                @click.stop="removeImage"
              />
            </div>
          </VCard>
        </VCol>
  
        <!-- Form Alanları -->
        <VCol cols="12">
          <VTextField v-model="form.title" label="Başlık" required />
        </VCol>
  
        <VCol cols="12">
          <VTextField v-model="form.content" label="Giriş" required />
        </VCol>
  
        <VCol cols="12">
          <VTextarea
            v-model="form.main_content"
            label="Ne düşünüyorsun?"
            rows="6"
            counter
            auto-grow
            no-resize
            clearable
            :max-rows="8"
            required
          />
        </VCol>
  
        <VCol cols="12" class="d-flex justify-end gap-4">
          <VBtn type="submit" :loading="formloading" color="primary">
            {{ iscreateform ? 'Paylaş' : 'Güncelle' }}
          </VBtn>
          <VBtn
          variant="tonal"
          :loading="formloading"
          @click="emit('closemodal')"
        >
          Vazgeç
        </VBtn>
        </VCol>
      </VRow>
    </VForm>
  </template>
  
 
  <style scoped>
  .upload-area {
    border: 2px dashed #ccc;
    border-radius: 8px;
    text-align: center;
    padding: 20px;
    cursor: pointer;
  }
  
  .upload-area.dragover {
    border-color: #1976d2;
  }
  
  .preview-container {
    position: relative;
    text-align: center;
  }
  
  .preview-image {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
  }
  
  .remove-btn {
    position: absolute;
    top: 10px;
    right: 10px;
  }
  </style>
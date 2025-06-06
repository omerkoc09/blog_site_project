<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { emailValidator, phoneValidator, requiredValidator } from '@validators'
import ApiService from '@/services/ApiService'
import { ErrorPopup } from '@/utils/Popup'

const router = useRouter()
const loading = ref(false)
const topics = ref<{ id: number, name: string }[]>([])
const selectedTopicIds = ref<number[]>([])


const form = ref({
  name: '',
  surname: '',
  phone: '',
  email: '',
  password: '',
  role: 1,
})

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

const onSubmit = async () => {
  loading.value = true

  // form verisini kopyala ve interests alanını ekle
  const payload = {
    ...form.value,
    interests: selectedTopicIds.value, // topic id dizisi olarak gönder
  }

  const [error, resp] = await ApiService.post(
    'user',
    payload,
    undefined,
    undefined,
  )

  loading.value = false

  if (error)
    return ErrorPopup(error)

  router.push('/auth/login')
}

onMounted(async () => {
  await loadTopics()
})
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <VCard class="auth-card pa-4" max-width="448">
      <VCardText class="pt-2">
        <h5 class="text-h5 font-weight-semibold mb-1">
          Register Account 🚀
        </h5>
        <p class="mb-0">
          Create your account
        </p>
      </VCardText>

      <VCardText>
        <VForm @submit.prevent="onSubmit">
          <VRow>
            <VCol cols="12" md="6">
              <VTextField v-model="form.name" label="Name" :rules="[requiredValidator]" />
            </VCol>

            <VCol cols="12" md="6">
              <VTextField v-model="form.surname" label="Surname" :rules="[requiredValidator]" />
            </VCol>

            <VCol cols="12">
              <VTextField v-model="form.email" label="Email" type="email"
                :rules="[requiredValidator, emailValidator]" />
            </VCol>

            <VCol cols="12">
              <VTextField v-model="form.phone" label="Phone" :rules="[requiredValidator, phoneValidator]" />
            </VCol>

            <VCol cols="12">
              <VTextField v-model="form.password" label="Password" type="password" :rules="[requiredValidator]" />
            </VCol>
            <VCol cols="12">
              <v-select v-model="selectedTopicIds" :items="topics" item-title="name" item-value="id" label="Topic Seç"
                multiple chips clearable :loading="topics.length === 0" />
            </VCol>

            <VCol cols="12"></VCol>


            <VCol cols="12">
              <VBtn block type="submit" :loading="loading">
                Register
              </VBtn>
            </VCol>

            <VCol cols="12" class="text-center">
              <RouterLink class="text-primary ms-2" to="/auth/login">
                Already have an account?
              </RouterLink>
            </VCol>
          </VRow>
        </VForm>
      </VCardText>
    </VCard>
  </div>
</template>

<route lang="yaml">
meta:
  layout: blank
  auth: false
</route>
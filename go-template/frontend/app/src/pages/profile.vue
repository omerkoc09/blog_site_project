<script setup lang="ts">
import ApiService from '@/services/ApiService'
import { ErrorPopup, SuccessPopup } from '@/utils/Popup'
import { emailValidator, phoneValidator, requiredValidator } from '@validators'
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'


const loading = ref(false)
const form = ref({
  name: '',
  surname: '',
  email: '',
  phone: '',
  password: '',
  about: '',
  likes: [],
  comments: [],
  posts:[],
})

onMounted(async () => {
  loading.value = true
  const [err, data] = await ApiService.get<any>('user/me')
  loading.value = false
  if (err)
    return ErrorPopup(err)
  form.value = data.data
})

const router = useRouter()

const onSubmit = async () => {
  loading.value = true
  const [err] = await ApiService.put('user/me', form.value)
  console.log('profile value', form.value.about)
  loading.value = false
  if (err)
    return ErrorPopup(err)
  SuccessPopup('Profil başarıyla güncellendi')
  await useUserStore().updateUser()
  router.push('/')
}
</script>

<template>

<div>
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
      </VContainer>
    </VAppBar>

        
  <VCard>
    <div class="d-flex justify-center pa-10 profile-container" >
      <VCard max-width="400px">
        <VCardTitle>
          <h2 class="text-center pa-5">
            Profil
          </h2>
        </VCardTitle>
        <VCardText>
          <VForm
            ref="formRef"
            @submit.prevent="onSubmit"
          >
            <VRow>
              <!-- Name -->
              <VCol cols="12">
                <VTextField
                  v-model="form.name"
                  label="İsim"
                  :rules="[requiredValidator]"
                />
              </VCol>
              <!-- Surname -->
              <VCol cols="12">
                <VTextField
                  v-model="form.surname"
                  label="Soyisim"
                  :rules="[requiredValidator]"
                />
              </VCol>
              <!-- email -->
              <VCol cols="12">
                <VTextField
                  v-model="form.email"
                  label="Email"
                  type="email"
                  :rules="[requiredValidator, emailValidator]"
                />
              </VCol>

              <!-- phone -->
              <VCol cols="12">
                <VTextField
                  v-model="form.phone"
                  name="telefon"
                  label="Telefon"
                  type="phone"
                  :rules="[requiredValidator, phoneValidator]"
                />
              </VCol>

              <!-- password -->
              <VCol cols="12">
                <VTextField
                  v-model="form.password"
                  label="Parola"
                  type="text"
                />
              </VCol>
              <VCol cols="12">
                <VTextarea
                    v-model="form.about"
                    label="Hakkımda"
                    :rules="[requiredValidator]"
                />
              </VCol>
            </VRow>
            <VRow>
              <VCol cols="12">
                <VBtn
                  block
                  type="submit"
                  :loading="loading"
                >
                  Kaydet
                </VBtn>
              </VCol>
            </VRow>
          </VForm>
        </VCardText>
      </VCard>
    </div>
  </VCard>
</div>
</template>

<style scoped lang="scss">
 .profile-container {
  margin-top: 150px;
 }
</style>

<route lang="yaml">
  meta:
    auth: false
    layout: blank
  </route>

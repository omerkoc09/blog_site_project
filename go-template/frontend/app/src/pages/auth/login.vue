<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { VForm } from 'vuetify/components'
import { emailValidator, requiredValidator } from '@validators'
import ApiService from '@/services/ApiService'
import { ErrorPopup } from '@/utils/Popup'
import { useUserStore } from '@/store/user'
import JwtService from '@/services/JwtService'

const loading = ref(false)
const isPasswordVisible = ref(false)
const formRef = ref<VForm>()

interface LoginForm {
  email: string
  password: string
}

const form = ref<LoginForm>({
  email: '',
  password: '',
})

const router = useRouter()

const onSubmit = async () => {
  const { valid } = await formRef.value!.validate()
  if (!valid)
    return

  console.log('Form values:', form.value)
  console.log(import.meta.env.VITE_API_BASE_URL)
  loading.value = true
  const [error, resp] = await ApiService.post<any>('auth/login', form.value)
  loading.value = false
  if (error)
    return ErrorPopup(error)
  await useUserStore().login(resp.data.access_token, resp.data.refresh_token)

  console.log('Access Token:', JwtService.getAccessToken())
  console.log('User Role:', useUserStore().user.role)
  console.log('Entire User:', useUserStore().user)

  await router.push('/')
}
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <div class="position-relative my-sm-16">


      <!-- 👉 Auth Card -->
      <VCard
        class="auth-card pa-4"
        max-width="448"
        min-width="448"
      >
        <VCardItem class="justify-center">
          <template #prepend>
            <div class="d-flex flex-column align-center">
              <h1 class="text-h2 font-weight-bold primary--text mb-2">BLOGGER.COM</h1>
              <div class="text-subtitle-1">Welcome to your blogging platform</div>
            </div>
          </template>
        </VCardItem>
        <VCardText>
          <VForm
            ref="formRef"
            @submit.prevent="onSubmit"
          >
            <VRow>
              <!-- email -->
              <VCol cols="12">
                <VTextField
                  v-model="form.email"
                  label="Email"
                  type="email"
                  :rules="[requiredValidator, emailValidator]"
                />
              </VCol>

              <!-- password -->
              <VCol cols="12">
                <VTextField
                  v-model="form.password"
                  label="Parola"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'tabler-eye-off' : 'tabler-eye'"
                  :rules="[requiredValidator]"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />

                <!-- remember me checkbox -->
                <div class="d-flex align-center justify-space-between flex-wrap mt-2 mb-4">
                  <!--                  <VCheckbox -->
                  <!--                    v-model="form.remember" -->
                  <!--                    label="Remember me" -->
                  <!--                  /> -->

                  <!--                  <RouterLink -->
                  <!--                    class="text-primary ms-2 mb-1" -->
                  <!--                    :to="{ name: 'auth-forgot-password' }" -->
                  <!--                  > -->
                  <!--                    Parolamı Unuttum? -->
                  <!--                  </RouterLink> -->
                </div>

                <div class="forgot-password-link">
                 <router-link to="/auth/forgot-password">
               Forgot Password?
               </router-link>
                </div>

                <!-- login button -->
                <VBtn
                  block
                  type="submit"
                  :loading="loading"
                >
                  GİRİŞ
                </VBtn>
              </VCol>

              <!-- create account -->
                           <VCol 
                             cols="12" 
                             class="text-center text-base" 
                           > 
                             <span>Hesabınız yok mu?</span> 
                             <RouterLink 
                               class="text-primary ms-2" 
                               to="/auth/register"
                             > 
                               Hesap oluştur 
                             </RouterLink> 
                           </VCol> 
               <!--  <VCol
                cols="12"
                class="d-flex align-center"
                >
                <VDivider />
                <span class="mx-4">or</span>
                <VDivider />
                </VCol>

                &lt;!&ndash; auth providers &ndash;&gt;
                <VCol
                cols="12"
                class="text-center"
                >
                <AuthProvider />
                </VCol> -->
            
            </VRow>
          </VForm>
        </VCardText>
      </VCard>
    </div>
  </div>
</template>

<style lang="scss">
@use "@core/scss/template/pages/page-auth.scss";
</style>

<route lang="yaml">
meta:
  layout: blank
  redirectIfLoggedIn: true
</route>

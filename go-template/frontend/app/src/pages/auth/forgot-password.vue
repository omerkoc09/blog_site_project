<script setup lang="ts">
import authV1BottomShape from '@images/svg/auth-v1-bottom-shape.svg'
import authV1TopShape from '@images/svg/auth-v1-top-shape.svg'
import { VNodeRenderer } from '@layouts/components/VNodeRenderer'
import { themeConfig } from '@themeConfig'
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

const form = ref({
  email: '',
})

const router = useRouter()

const onSubmit = async () => {
  const { valid } = await formRef.value!.validate()
  if (!valid)
    return

  console.log('Form value:', form.value)
  console.log(import.meta.env.VITE_API_BASE_URL)
  loading.value = true
  const [error, resp] = await ApiService.post<any>('auth/forgot-password', form.value)
  loading.value = false
  if (error)
    return ErrorPopup(error)
   await useUserStore().login(resp.data.access_token, resp.data.refresh_token)

  console.log('Access Token:', resp.data.access_token)
  console.log('User Role:', useUserStore().user.role)
  console.log('Entire User:', useUserStore().user)

 // await router.push('/reset-password')
}
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <div class="position-relative my-sm-16">
      <!-- 👉 Top shape -->
      <VImg
        :src="authV1TopShape"
        class="auth-v1-top-shape d-none d-sm-block"
      />

      <!-- 👉 Bottom shape -->
      <VImg
        :src="authV1BottomShape"
        class="auth-v1-bottom-shape d-none d-sm-block"
      />

      <!-- 👉 Auth card -->
      <VCard
        class="auth-card pa-4"
        max-width="448"
      >
       <!--  <VCardItem class="justify-center">                            IT SHOWS ROW BINARY PNG IMAGE FILE?
          <template #prepend>
            <div class="d-flex">
              <VNodeRenderer :nodes="themeConfig.app.logo" />
            </div>
          </template>

          <VCardTitle class="font-weight-bold text-h5 py-1">
            {{ themeConfig.app.title }}
          </VCardTitle>
        </VCardItem> -->

        <VCardText class="pt-2">
          <h5 class="text-h5 font-weight-semibold mb-1">
            Forgot Password? 🔒
          </h5>
          <p class="mb-0">
            Enter your email and we'll send you instructions to reset your password
          </p>
        </VCardText>

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
                />
              </VCol>

              <!-- reset password -->
              <VCol cols="12">
                <VBtn
                  block
                  type="submit"
                >
                  Send Reset Link
                </VBtn>
              </VCol>

              <!-- back to login -->
              <VCol cols="12">
                <RouterLink
                  class="d-flex align-center justify-center"
                  :to="{ name: 'auth-login' }"
                >
                  <VIcon
                    icon="tabler-chevron-left"
                    class="flip-in-rtl"
                  />
                  <span>Back to login</span>
                </RouterLink>
              </VCol>
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


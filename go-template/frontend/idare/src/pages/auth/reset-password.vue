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
import ForgotPassword from './auth/forgot-password.vue'

const loading = ref(false)
const isPasswordVisible = ref(false)
const formRef = ref<VForm>()

const form = ref({
  password: '',
})



const router = useRoute()
const token = router.query.token as string
console.log("access token from link", token)
console.log('Access Token from forgot:', JwtService.getAccessToken())

const onSubmit = async () => {
  const resetData = {
    token: token,
    password: form.value.password
  }
  
  const [error, resp] = await ApiService.put('auth/reset-password', resetData)
  if (error) {
    return ErrorPopup(error)
  }

  // Success - redirect to login
  const route = useRouter()
  route.push('/login')
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


        <VCardText class="pt-2">
          <h5 class="text-h5 font-weight-semibold mb-1">
            Reset Password 
          </h5>
          <p class="mb-0">
            Enter your new password
          </p>
        </VCardText>

        <VCardText>
          <VForm
            ref="formRef"
            @submit.prevent="onSubmit"
          >
            <VRow>
              <!-- new password -->
              <VCol cols="12">
                <VTextField
                  v-model="form.password"
                  label="New password"
                />
              </VCol>

              <!-- reset password -->
              <VCol cols="12">
                <VBtn
                  block
                  type="submit"
                >
                  Submit
                </VBtn>
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
</route>



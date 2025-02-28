<script setup lang="ts">
import { ref } from 'vue'

const isLoading = ref(false)

const fallbackHandle = () => {
  isLoading.value = true
}

const resolveHandle = () => {
  isLoading.value = false
}

defineExpose({
  fallbackHandle,
  resolveHandle
})
</script>

<template>
  <VOverlay
    v-model="isLoading"
    persistent
    class="align-center justify-center"
  >
    <VProgressCircular
      indeterminate
      color="primary"
      size="50"
    />
  </VOverlay>
</template>
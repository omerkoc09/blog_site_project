<script setup lang="ts">
import { computed, ref, watch } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  itemsPerPage: {
    type: Number,
    default: 6
  }
})

const emit = defineEmits(['pageChanged'])

const currentPage = ref(1)
const totalPages = computed(() => Math.ceil(props.items.length / props.itemsPerPage))
const paginatedItems = computed(() => {
  const start = (currentPage.value - 1) * props.itemsPerPage
  const end = start + props.itemsPerPage
  return props.items.slice(start, end)
})

function goToPage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    emit('pageChanged', {
      page: currentPage.value,
      items: paginatedItems.value
    })
  }
}

// Reset to page 1 when items change
watch(() => props.items, () => {
  currentPage.value = 1
})
</script>

<template>
  <div class="pagination-wrapper">
    <!-- Sayfalanmış öğeler -->
    <slot :items="paginatedItems" />

    <!-- Pagination kontrolleri -->
    <div class="d-flex justify-center mt-6" v-if="totalPages > 1">
      <v-pagination
          v-model="currentPage"
          :length="totalPages"
          :total-visible="5"
          rounded="circle"
          @update:model-value="goToPage"
      ></v-pagination>
    </div>
  </div>
</template>
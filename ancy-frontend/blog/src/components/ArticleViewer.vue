<template>
  <div class="mx-auto max-w-3xl px-4">
    <slot name="header"></slot>
    <MdViewer :content="content" />
    <slot name="footer"></slot>
  </div>
</template>

<script setup lang="ts">
import MdViewer from '@/components/MdViewer.vue'
import { computed } from 'vue'

const props = defineProps<{
  article?: { content: string }
}>()
const content = computed(() => props.article?.content)
</script>

<style scoped lang="scss"></style>

<template>
  <div class="mx-auto max-w-3xl px-4">
    <slot name="header"></slot>
    <MdViewer :content="content" />
    <slot name="footer"></slot>
    <DisplayComment
      :id="props.article.id"
      type="Article"
      v-if="props.article?.id && props.showComment"
    />
  </div>
</template>

<script setup lang="ts">
import MdViewer from '@/components/MdViewer.vue'
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    article?: { id: number; content: string }
    showComment?: boolean
  }>(),
  {
    showComment: true,
  },
)
const content = computed(() => props.article?.content)
</script>

<style scoped lang="scss"></style>

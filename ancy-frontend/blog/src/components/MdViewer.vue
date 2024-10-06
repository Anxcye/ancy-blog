<template>
  <div>
    <MdPreview :model-value="props.content" previewTheme="ancy" />
  </div>
</template>

<script setup lang="ts">
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

const props = defineProps<{
  content?: string
}>()
</script>

<style lang="scss">
.md-editor {
  background-color: transparent !important;
}
.md-editor-preview-wrapper {
  padding: 0 !important;
}
</style>

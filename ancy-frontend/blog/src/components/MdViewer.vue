<template>
  <div>
    <MdPreview :model-value="props.content" previewTheme="ancy" />
  </div>
</template>

<script setup lang="ts">
import { config, MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
config({
  markdownItPlugins(plugins) {
    return plugins.filter((p) => p.type !== 'xss')
  },
})
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

<template>
  <div
    class="absolute top-0 left-0 w-full bg-gradient-to-b from-primary-bg to-transparent -z-10"
    :style="{ height: props.height || '100%' }"
  ></div>
</template>

<script setup lang="ts">
const props = defineProps<{
  height?: string
}>()
</script>

<style scoped lang="scss"></style>

<template>
  <a-tooltip>
    <template #title>
      {{ props.time }}
    </template>
    {{ timeAgo(new Date(props.time)) }}
  </a-tooltip>
</template>

<script setup lang="ts">
import timeAgo from '@/utils/timeAgo'

const props = defineProps<{
  time: string
}>()
</script>

<style scoped></style>

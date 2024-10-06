<template>
  <div>
    <a-timeline class="mt-8">
      <a-timeline-item
        :color="colorStore.getPrimaryColor()"
        v-for="item in props.list"
        :key="item.id"
      >
        <slot name="item" :item="item" />
      </a-timeline-item>
    </a-timeline>
  </div>
</template>

<script setup lang="ts">
import { useColorStore } from '@/stores/color'

const colorStore = useColorStore()

const props = defineProps<{
  list: { id: number; [key: string]: any }[]
}>()
</script>

<style scoped lang="scss"></style>

<template>
  <div class="timeline_container">
    <a-timeline class="mt-8" mode="left">
      <a-timeline-item
        :color="colorStore.getPrimaryColor()"
        v-for="item in props.list"
        :key="item.id"
      >
        <template #label>
          <div class="text-sm text-gray">{{ timeAgo(new Date(item.createTime)) }}</div>
        </template>
        <slot name="item" :item="item" />
      </a-timeline-item>
    </a-timeline>
  </div>
</template>

<script setup lang="ts">
import { useColorStore } from '@/stores/color'
import timeAgo from '@/utils/timeAgo'

const colorStore = useColorStore()

const props = defineProps<{
  list: { id: number; [key: string]: any }[]
}>()
</script>

<style scoped lang="scss"></style>

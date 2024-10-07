<template>
  <div class="timeline_container">
    <a-timeline class="mt-8" mode="left">
      <a-timeline-item
        :color="colorStore.getPrimaryColor()"
        v-for="item in props.list"
        :key="item.id"
      >
        <template #label>
          <TimeTip :time="time(item)" />
        </template>
        <slot name="item" :item="item" />
      </a-timeline-item>
    </a-timeline>
  </div>
</template>

<script setup lang="ts">
import { useColorStore } from '@/stores/color'
import TimeTip from '@/components/TimeTip.vue'
const colorStore = useColorStore()

const props = defineProps<{
  list: { id: number; [key: string]: any }[]
  timeField?: string
}>()

const timeField = props.timeField || 'createTime'

const time = (item: any) => {
  return item[timeField]
}
</script>

<style scoped lang="scss"></style>

<template>
  <div class="timeline_container">
    <a-timeline class="mt-8" mode="left" :pending="isPending">
      <template #pendingDot>
        <LoadingOutlined v-if="hasMore" />
        <CheckCircleFilled v-else />
      </template>
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
import { LoadingOutlined, CheckCircleFilled } from '@ant-design/icons-vue'
import { computed } from 'vue'

const colorStore = useColorStore()
const props = withDefaults(
  defineProps<{
    list: { id: number; [key: string]: any }[]
    timeField?: string
    pending?: boolean
    total?: number
  }>(),
  {
    timeField: 'createTime',
    pending: true,
  },
)

const hasMore = computed(() => {
  return props.total !== props.list.length
})

const isPending = computed(() => {
  return props.pending ? (hasMore.value ? '加载中...' : '已经是全部了') : false
})

const time = (item: any) => {
  return item[props.timeField]
}
</script>

<style scoped lang="scss"></style>

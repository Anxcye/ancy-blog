<template>
  <div class="hover-card" ref="card">
    <div class="content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const ANGLE = 10
const card = ref<HTMLElement | null>(null)

const handleMouseMove = (e: MouseEvent) => {
  if (!card.value) return
  const content = card.value.querySelector('.content') as HTMLElement
  const rect = card.value.getBoundingClientRect()
  const w = rect.width
  const h = rect.height
  const y = ((e.clientX - rect.left - w * 0.5) / w) * ANGLE
  const x = ((1 - (e.clientY - rect.top - h * 0.5)) / h) * ANGLE

  content.style.transform = `
    perspective(1000px)
    rotateX(${x}deg)
    rotateY(${y}deg)
    scale3d(1, 1, 1)
  `
}

const handleMouseOut = () => {
  if (!card.value) return
  const content = card.value.querySelector('.content') as HTMLElement
  content.style.transform = `
    perspective(1000px)
    rotateX(0deg)
    rotateY(0deg)
    scale3d(1, 1, 1)
  `
  content.style.transition = 'all 0.2s ease'
}

onMounted(() => {
  if (card.value) {
    card.value.addEventListener('mousemove', handleMouseMove)
    card.value.addEventListener('mouseout', handleMouseOut)
  }
})

onUnmounted(() => {
  if (card.value) {
    card.value.removeEventListener('mousemove', handleMouseMove)
    card.value.removeEventListener('mouseout', handleMouseOut)
  }
})
</script>

<style scoped lang="scss">
.hover-card {
  perspective: 1000px;
  .content {
    transition: transform 0.2s;
    transform-style: preserve-3d;
    transform: translateZ(0);
  }
}
</style>

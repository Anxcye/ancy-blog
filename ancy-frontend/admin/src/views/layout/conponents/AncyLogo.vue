<template>
  <div class="ancy-logo">
    Ancy
    <div class="ancy-logo-text">{{ hello }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// const hello = ref('hello')
const hello = computed(() => {
  if (new Date().getHours() < 12) {
    return '上午好'
  } else if (new Date().getHours() < 18) {
    return '下午好'
  } else {
    return '晚上好'
  }
  // return
})
</script>

<style scoped lang="scss">
.ancy-logo {
  font-size: 20rem;
  font-weight: bold;
  color: #00000022;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;

  .ancy-logo-text {
    font-size: 5rem;
    font-weight: lighter;
    color: #00000022;
    align-self: flex-end;
  }
}
</style>

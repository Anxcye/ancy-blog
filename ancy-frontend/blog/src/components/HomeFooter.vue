<template>
  <div class="border-t border-gray-bg bg-primary-bg-1 h-40 p-3">
    <div class="content">
      <div class="flex flex-col">
        <div class="flex flex-row items-center gap-3 text-sm">
          <span>
            Powered by
            <a href="https://github.com/Anxcye/ancy-blog" target="_blank">Ancy</a>
          </span>
          |
          <a @click="openLogin" class="hover:cursor-pointer">管理员登录</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const openLogin = () => {
  window.open('https://github.com/Anxcye/anci-blog', '_blank')
}
</script>

<style scoped lang="scss"></style>

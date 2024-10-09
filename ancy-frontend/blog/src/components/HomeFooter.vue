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
          <router-link
            v-for="item in footer1"
            :key="item.index"
            :to="item.url"
            class="hover:cursor-pointer"
          >
            {{ item.text }}
          </router-link>
          |
          <a @click="openLogin" class="hover:cursor-pointer">管理员登录</a>
        </div>
        <div class="flex flex-row items-center gap-3 text-sm">
          <a
            v-for="item in footer2"
            :key="item.index"
            :href="item.url"
            class="hover:cursor-pointer"
          >
            {{ item.text }}
          </a>
        </div>
      </div>
    </div>
  </div>
  <a-modal v-model:open="open" :footer="null">
    <template #closeIcon></template>
    <div class="flex flex-col items-center w-full">
      <div class="bg-bg-color text-tx-color p-5 rounded-xl">
        <div class="text-xl font-medium mb-4 text-center">登录到 {{ baseInfoStore.getName() }}</div>
        <a-button type="primary" @click="handleLogout" class="w-full mb-4" v-if="adminStore.token">
          当前已登录,点击以登出
        </a-button>
        <a-form v-else>
          <a-form-item>
            <a-input v-model:value="form.userName" placeholder="用户名" />
          </a-form-item>
          <a-form-item>
            <a-input v-model:value="form.password" type="password" placeholder="密码" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="handleLogin" class="w-full" :loading="loginLoading">
              登录
            </a-button>
          </a-form-item>
        </a-form>
        <a-button @click="handleCancel" class="w-full">取消</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import type { LoginParams } from '@/api/admin/type'
import { useAdminStore } from '@/stores/admin'
import { useBaseInfoStore } from '@/stores/baseInfo'
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'

const baseInfoStore = useBaseInfoStore()
const open = ref(false)
const form = ref<LoginParams>({})
const adminStore = useAdminStore()
const footer1 = computed(() => baseInfoStore.getFooter()?.filter((item) => item.position === 1))
const footer2 = computed(() => baseInfoStore.getFooter()?.filter((item) => item.position === 2))
const loginLoading = ref(false)

const openLogin = () => {
  open.value = true
  form.value = {}
}

const handleLogin = () => {
  loginLoading.value = true
  adminStore.login(form.value).finally(() => {
    loginLoading.value = false
  })
}

const handleCancel = () => {
  open.value = false
}

const handleLogout = () => {
  adminStore.logout()
}
</script>

<style scoped lang="scss"></style>

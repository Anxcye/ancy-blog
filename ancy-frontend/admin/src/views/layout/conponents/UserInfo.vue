<template>
  <div>
    <el-dropdown>
      <el-avatar :src="userInfo.avatar" />
      <template #dropdown>
        <el-dropdown-menu class="header-user">
          <div class="user-info">
            <div class="user-info-name">{{ userInfo.nickName }}</div>
            <div class="user-info-email">{{ userInfo.email }}</div>
            <div class="user-info-role">
              <div v-for="role in roles" :key="role">
                <el-tag type="primary">{{ role }}</el-tag>
              </div>
            </div>
          </div>
          <el-dropdown-item divided class="logout">
            <el-button type="primary" @click="logout">退出登录</el-button>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()
const userInfo = userStore.userInfo!.userInfoVo
const roles = userStore.userInfo!.role

const logout = () => {
  userStore.logout()
}
</script>

<style scoped lang="scss">
.header-user {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px;

  .user-info-name {
    margin: 10px 0;
    font-size: 16px;
    font-weight: bold;
  }

  .user-info-email {
    margin: 10px 0;
    font-size: 14px;
    color: #999;
  }

  .user-info-role {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }

  .logout {
    display: flex;
    justify-content: center;
    text-align: center;
  }
}
</style>

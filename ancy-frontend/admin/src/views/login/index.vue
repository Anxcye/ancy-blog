<template>
  <div class="login-container">
    <div class="gradient-background"></div>
    <HoverCard>
      <div class="login-card">
        <div class="login-left">
          <h2 class="login-title">Ancy</h2>
        </div>

        <div class="login-right">
          <el-form :model="loginForm" :rules="rules" ref="loginFormRef" class="login-form">
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                :prefix-icon="UserFilled"
                placeholder="用户名"
              ></el-input>
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                :prefix-icon="Lock"
                type="password"
                placeholder="密码"
                show-password
              ></el-input>
            </el-form-item>
            <el-form-item>
              <span class="login-button-container">
                <el-button
                  type="primary"
                  @click="handleVisitor"
                  :loading="loading"
                  class="visitor-button"
                >
                  访客参观
                </el-button>
                <el-button
                  type="primary"
                  @click="handleLogin"
                  :loading="loading"
                  class="login-button"
                >
                  登录
                </el-button>
              </span>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </HoverCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { UserFilled, Lock } from '@element-plus/icons-vue'
import HoverCard from '@/components/HoverCard.vue'
import { useUserStore } from '@/stores/modules/user'
import router from '@/router'

const userStore = useUserStore()

const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  loading.value = true
  await loginFormRef.value.validate()
  try {
    await userStore.userLogin(loginForm)
    router.push('/ancy')
  } finally {
    loading.value = false
  }
}

const handleVisitor = async () => {
  loginForm.username = 'visitor'
  loginForm.password = '1234'
  await handleLogin()
}
</script>

<style scoped lang="scss">
.login-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.gradient-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, #ff5ccc, #dc52bf, #b947b1, #973da4, #743296, #512889);
  background-size: 400% 400%;
  animation: gradientFlow 15s ease infinite;
}

@keyframes gradientFlow {
  0% {
    background-position: 0% 50%;
  }

  50% {
    background-position: 100% 0%;
  }

  100% {
    background-position: 0% 50%;
  }
}

.login-card {
  max-width: 550px;
  width: 550px;
  background-color: rgba(255, 255, 255);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  height: 350px;
  display: flex;
  overflow: hidden;
}

@media (max-width: 768px) {
  .login-card {
    flex-direction: column;
    width: 80vw;
    margin: 0 auto;
  }
  .login-button-container {
    flex-direction: column;
    .visitor-button {
      margin-bottom: 10px;
      width: 100%;
    }
    .el-button + .el-button {
      margin-left: 0px;
    }
  }
}

:deep(.el-input__wrapper) {
  background-color: rgba(255, 255, 255);
}

:deep(.el-card__body) {
  padding: 0px;
  height: 100%;
}

.login-left {
  flex: 1;
  text-align: left;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: $ac-primary-color;
  color: #fff;
  height: 100%;

  .login-title {
    text-align: left;
    font-weight: 600;
    font-size: 1.2em;
  }
}

.login-right {
  flex: 1;
  padding: 15px;
  display: flex;
  align-items: center;
  justify-content: center;

  .login-button-container {
    display: flex;
    align-items: center;
    width: 100%;
    gap: 10px;

    .login-button {
      flex: 1;
      width: 100%;
    }

    .visitor-button {
      background-color: #fff;
      color: #000;
    }
  }
}

.login-form {
  width: 100%;
}
</style>

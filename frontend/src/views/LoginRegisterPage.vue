<template>
  <div class="auth-page">
    <ArtPlaceholder label="游戏背景海报" sub="全屏背景图" width="100vw" height="100vh" bgColor="#0a0a1e" class="bg" />
    
    <div class="card">
      <button class="card-close" title="关闭程序" @click="handleExit">✕</button>
      
      <ArtPlaceholder label="游戏 LOGO" width="180" height="60" bgColor="#1a1a3e" style="margin:0 auto" />
      
      <div class="tabs">
        <button :class="['tab', { active: mode === 'login' }]" @click="switchTab('login')">登录</button>
        <button :class="['tab', { active: mode === 'register' }]" @click="switchTab('register')">注册</button>
      </div>

      <!-- 登录表单 -->
      <form v-if="mode === 'login'" @submit.prevent="handleLogin" class="form">
        <input v-model="loginForm.username" class="input" placeholder="请输入账号" />
        <input v-model="loginForm.password" type="password" class="input" placeholder="请输入密码" />
        <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
        <button type="submit" class="btn-submit" :disabled="loginLoading">
          {{ loginLoading ? '登录中...' : '登录' }}
        </button>
      </form>

      <!-- 注册表单 -->
      <form v-if="mode === 'register'" @submit.prevent="handleRegister" class="form">
        <input v-model="registerForm.username" class="input" placeholder="请输入账号" />
        <input v-model="registerForm.password" type="password" class="input" placeholder="请输入密码" />
        <input v-model="registerForm.confirmPassword" type="password" class="input" placeholder="请确认密码" />
        <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
        <button type="submit" class="btn-submit" :disabled="registerLoading">
          {{ registerLoading ? '注册中...' : '注册' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useUserStore } from '@/store/user'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'

const router = useRouter()
const userStore = useUserStore()
const mode = ref<'login' | 'register'>('login')
const errorMsg = ref('')
const loginLoading = ref(false)
const registerLoading = ref(false)

const loginForm = reactive({ username: '', password: '' })
const registerForm = reactive({ username: '', password: '', confirmPassword: '' })

function switchTab(tab: 'login' | 'register') {
  mode.value = tab
  errorMsg.value = ''
}

async function handleLogin() {
  errorMsg.value = ''
  if (!loginForm.username || !loginForm.password) {
    errorMsg.value = '请填写完整信息'
    return
  }
  loginLoading.value = true
  try {
    const res: any = await api.post('/login', { username: loginForm.username, password: loginForm.password })
    const data = res.data || res
    localStorage.setItem('token', data.token)
    await userStore.fetchUserInfo()
    router.replace('/lobby')
  } catch (err: any) {
    const msg = err?.response?.data?.message || err?.data?.message || err?.message || '登录失败'
    errorMsg.value = msg
  } finally {
    loginLoading.value = false
  }
}

async function handleRegister() {
  errorMsg.value = ''
  if (!registerForm.username || !registerForm.password || !registerForm.confirmPassword) {
    errorMsg.value = '请填写完整信息'
    return
  }
  if (registerForm.password !== registerForm.confirmPassword) {
    errorMsg.value = '两次密码不一致'
    return
  }
  registerLoading.value = true
  try {
    const res: any = await api.post('/register', { username: registerForm.username, password: registerForm.password })
    const data = res.data || res
    localStorage.setItem('token', data.token)
    await userStore.fetchUserInfo()
    router.replace('/lobby')
  } catch (err: any) {
    const msg = err?.response?.data?.message || err?.data?.message || err?.message || '注册失败'
    errorMsg.value = msg
  } finally {
    registerLoading.value = false
  }
}

function handleExit() {
  if (window.close) window.close()
}
</script>

<style scoped>
.auth-page {
  width: 100vw; height: 100vh;
  display: flex; align-items: center; justify-content: center;
  position: relative;
}
.bg { position: absolute; inset: 0; z-index: 0; }
.card {
  position: relative; z-index: 1;
  width: 400px;
  background: rgba(20, 20, 40, 0.92);
  backdrop-filter: blur(12px);
  border: 2px solid #3a3a5e;
  border-radius: 20px;
  padding: 36px 32px 32px;
  display: flex; flex-direction: column; gap: 24px;
}
.card-close {
  position: absolute; top: 12px; right: 14px;
  background: none; border: none; color: #666; font-size: 20px;
  cursor: pointer; width: 30px; height: 30px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
}
.card-close:hover { background: rgba(255,255,255,0.1); color: #fff; }
.tabs { display: flex; gap: 8px; }
.tab {
  flex: 1; padding: 10px;
  background: rgba(255,255,255,0.05);
  border: 1px solid #3a3a5e; border-radius: 8px;
  color: #888; font-size: 15px; cursor: pointer; transition: all 0.2s;
}
.tab.active { background: linear-gradient(135deg, #ffd70033, #ff8c0033); border-color: #ffd700; color: #ffd700; font-weight: bold; }
.form { display: flex; flex-direction: column; gap: 14px; }
.input {
  background: rgba(255,255,255,0.08); border: 1px solid #3a3a5e; border-radius: 8px;
  padding: 12px 16px; color: #e0e0e0; font-size: 14px; outline: none;
}
.input:focus { border-color: #ffd700; }
.input::placeholder { color: #666; }
.error-text { color: #ef4444; font-size: 13px; margin: 0; text-align: center; }
.btn-submit {
  background: linear-gradient(135deg, #ffd700, #ff8c00);
  border: none; border-radius: 10px; padding: 13px;
  color: #1a1a1a; font-size: 16px; font-weight: bold; cursor: pointer; transition: all 0.2s;
}
.btn-submit:hover:not(:disabled) { filter: brightness(1.1); transform: translateY(-1px); }
.btn-submit:disabled { opacity: 0.5; cursor: not-allowed; }
</style>

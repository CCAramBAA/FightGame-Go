<template>
  <div class="login-page">
    <div class="login-card">
      <div class="logo-section">
        <h1 class="game-title">FightGame</h1>
        <p class="subtitle">格斗对战平台</p>
      </div>

      <div class="form">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            class="form-input"
            type="text"
            placeholder="请输入用户名"
            @keyup.enter="handleLogin"
          />
          <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            class="form-input"
            type="password"
            placeholder="请输入密码"
            @keyup.enter="handleLogin"
          />
          <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
        </div>

        <button class="btn-primary" :disabled="loading" @click="handleLogin">
          {{ loading ? '登录中...' : '登录' }}
        </button>

        <div v-if="serverError" class="server-error">{{ serverError }}</div>
      </div>

      <div class="bottom-actions">
        <button class="btn-link" @click="$router.push('/register')">没有账号？去注册</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const serverError = ref('')
const errors = reactive({ username: '', password: '' })

const form = reactive({
  username: '',
  password: '',
})

function validate(): boolean {
  errors.username = ''
  errors.password = ''
  if (!form.username.trim()) {
    errors.username = '请输入用户名'
    return false
  }
  if (!form.password) {
    errors.password = '请输入密码'
    return false
  }
  if (form.password.length < 6) {
    errors.password = '密码至少6位'
    return false
  }
  return true
}

async function handleLogin() {
  if (!validate()) return

  serverError.value = ''
  loading.value = true
  try {
    const res: any = await api.post('/login', {
      username: form.username,
      password: form.password,
    })

    const d = res.data
    userStore.setAuth(d.token, d.id, d.username, d.nickname || d.username)
    router.push('/lobby')
  } catch (err: any) {
    const msg = err?.response?.data?.message || err?.message || '登录失败，请检查后端服务是否启动'
    serverError.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

.login-card {
  width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

.logo-section {
  text-align: center;
  margin-bottom: 32px;
}

.game-title {
  font-size: 2.5rem;
  font-weight: 800;
  color: #e94560;
  margin: 0;
  text-shadow: 0 0 20px rgba(233, 69, 96, 0.5);
}

.subtitle {
  color: #a0a0a0;
  margin: 8px 0 0;
  font-size: 0.9rem;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  color: #ccc;
  font-size: 0.9rem;
}

.form-input {
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  color: #fff;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: #e94560;
}

.form-input::placeholder {
  color: #555;
}

.form-error {
  color: #ff6b6b;
  font-size: 0.8rem;
}

.btn-primary {
  padding: 12px;
  background: #e94560;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: #d63851;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.server-error {
  padding: 10px 14px;
  background: rgba(255, 107, 107, 0.15);
  border: 1px solid rgba(255, 107, 107, 0.3);
  border-radius: 8px;
  color: #ff6b6b;
  font-size: 0.85rem;
  text-align: center;
}

.bottom-actions {
  text-align: center;
  margin-top: 20px;
}

.btn-link {
  background: none;
  border: none;
  color: #a0a0a0;
  font-size: 0.85rem;
  cursor: pointer;
  text-decoration: underline;
}

.btn-link:hover {
  color: #e94560;
}
</style>

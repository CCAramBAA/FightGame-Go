<template>
  <div class="register-page">
    <div class="register-card">
      <div class="logo-section">
        <h1 class="game-title">FightGame</h1>
        <p class="subtitle">创建你的账号</p>
      </div>

      <div class="form">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            class="form-input"
            type="text"
            placeholder="4-20位字母数字"
            @keyup.enter="handleRegister"
          />
          <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">昵称</label>
          <input
            v-model="form.nickname"
            class="form-input"
            type="text"
            placeholder="显示名称"
            @keyup.enter="handleRegister"
          />
          <span v-if="errors.nickname" class="form-error">{{ errors.nickname }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            class="form-input"
            type="password"
            placeholder="至少6位"
            @keyup.enter="handleRegister"
          />
          <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
        </div>

        <div class="form-group">
          <label class="form-label">确认密码</label>
          <input
            v-model="form.confirmPassword"
            class="form-input"
            type="password"
            placeholder="再次输入密码"
            @keyup.enter="handleRegister"
          />
          <span v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</span>
        </div>

        <button class="btn-primary" :disabled="loading" @click="handleRegister">
          {{ loading ? '注册中...' : '注册' }}
        </button>

        <div v-if="serverError" class="server-error">{{ serverError }}</div>
        <div v-if="serverSuccess" class="server-success">{{ serverSuccess }}</div>
      </div>

      <div class="bottom-actions">
        <button class="btn-link" @click="$router.push('/login')">已有账号？去登录</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
const loading = ref(false)
const serverError = ref('')
const serverSuccess = ref('')
const errors = reactive({ username: '', nickname: '', password: '', confirmPassword: '' })

const form = reactive({
  username: '',
  nickname: '',
  password: '',
  confirmPassword: '',
})

function validate(): boolean {
  errors.username = ''
  errors.nickname = ''
  errors.password = ''
  errors.confirmPassword = ''

  if (!form.username.trim()) {
    errors.username = '请输入用户名'
    return false
  }
  if (form.username.length < 4 || form.username.length > 20) {
    errors.username = '用户名4-20位'
    return false
  }
  if (!/^[a-zA-Z0-9_]+$/.test(form.username)) {
    errors.username = '只允许字母、数字和下划线'
    return false
  }
  if (!form.nickname.trim()) {
    errors.nickname = '请输入昵称'
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
  if (form.password !== form.confirmPassword) {
    errors.confirmPassword = '两次输入的密码不一致'
    return false
  }
  return true
}

async function handleRegister() {
  if (!validate()) return

  serverError.value = ''
  serverSuccess.value = ''
  loading.value = true
  try {
    await api.post('/register', {
      username: form.username,
      nickname: form.nickname,
      password: form.password,
    })

    serverSuccess.value = '注册成功，正在跳转...'
    setTimeout(() => router.push('/login'), 800)
  } catch (err: any) {
    const msg = err?.response?.data?.message || err?.message || '注册失败，请检查后端服务是否启动'
    serverError.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

.register-card {
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

.server-success {
  padding: 10px 14px;
  background: rgba(46, 204, 113, 0.15);
  border: 1px solid rgba(46, 204, 113, 0.3);
  border-radius: 8px;
  color: #2ecc71;
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

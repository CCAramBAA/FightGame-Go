<template>
  <div class="home-page">
    <header class="header">
      <h1>FightGame</h1>
      <p class="subtitle">在线格斗对战平台</p>
    </header>

    <main class="main">
      <!-- 未登录 -->
      <div v-if="!userStore.isLoggedIn" class="actions">
        <button class="btn btn-primary" @click="$router.push('/login')">登录</button>
        <button class="btn btn-secondary" @click="$router.push('/register')">注册</button>
      </div>

      <!-- 已登录 -->
      <div v-else class="actions">
        <button class="btn btn-primary" @click="$router.push('/lobby')">进入大厅</button>
        <button class="btn btn-secondary" @click="$router.push('/game')">快速对战</button>
      </div>

      <!-- 角色展示 -->
      <div v-if="userStore.isLoggedIn && userStore.myCharacters.length" class="character-preview">
        <h3>我的角色</h3>
        <div class="character-grid">
          <div
            v-for="char in userStore.myCharacters"
            :key="char.id"
            class="character-card"
          >
            <div class="char-avatar">{{ char.name[0] }}</div>
            <span class="char-name">{{ char.name }}</span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

onMounted(async () => {
  if (userStore.isLoggedIn) {
    await userStore.fetchUserInfo()
    await userStore.fetchMyCharacters()
  }
})
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  color: #fff;
}
.header {
  text-align: center;
  margin-bottom: 3rem;
}
.header h1 {
  font-size: 3.5rem;
  margin: 0;
  background: linear-gradient(90deg, #e94560, #f39c12);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.subtitle {
  font-size: 1.2rem;
  color: #a0a0a0;
  margin-top: 0.5rem;
}
.actions {
  display: flex;
  gap: 1.5rem;
}
.btn {
  padding: 0.8rem 2.5rem;
  font-size: 1.1rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 600;
}
.btn-primary {
  background: #e94560;
  color: #fff;
}
.btn-primary:hover {
  background: #c43a52;
  transform: translateY(-2px);
}
.btn-secondary {
  background: transparent;
  color: #e94560;
  border: 2px solid #e94560;
}
.btn-secondary:hover {
  background: rgba(233, 69, 96, 0.1);
  transform: translateY(-2px);
}
.character-preview {
  margin-top: 3rem;
  text-align: center;
}
.character-preview h3 {
  color: #a0a0a0;
  margin-bottom: 1rem;
}
.character-grid {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}
.character-card {
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 12px;
  padding: 1rem;
  width: 100px;
  transition: all 0.3s;
}
.character-card:hover {
  border-color: #e94560;
  transform: translateY(-2px);
}
.char-avatar {
  width: 60px;
  height: 60px;
  background: #e94560;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0 auto 8px;
}
.char-name {
  color: #ddd;
  font-size: 0.85rem;
}
</style>

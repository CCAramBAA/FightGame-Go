<template>
  <header class="global-nav">
    <div class="nav-left">
      <slot name="left">
        <button v-if="showBack" class="btn-back" @click="$router.back()">← 返回</button>
        <ArtPlaceholder label="LOGO" width="120" height="36" />
      </slot>
    </div>
    <div class="nav-center">
      <slot name="center">
        <span v-if="userStore.nickname" class="nickname">{{ userStore.nickname }}</span>
        <span class="gold">💰 {{ userStore.gold }}</span>
        <span class="rank">🏆 {{ userStore.rankScore }} 分</span>
      </slot>
    </div>
    <div class="nav-right">
      <slot name="right">
        <button class="btn-icon" title="好友列表" @click="$emit('openFriends')">👥</button>
        <button class="btn-icon" title="设置" @click="$emit('openSettings')">⚙️</button>
        <button class="btn-icon" title="退出" @click="handleLogout">🚪</button>
      </slot>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'
import ArtPlaceholder from './ArtPlaceholder.vue'

defineProps<{ showBack?: boolean }>()
defineEmits<{
  openSettings: []
  openFriends: []
}>()

const userStore = useUserStore()
const router = useRouter()

function handleLogout() {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.global-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 20px;
  background: linear-gradient(180deg, #0f0f23, #1a1a2e);
  border-bottom: 1px solid #2a2a4a;
  position: relative;
  z-index: 100;
  flex-shrink: 0;
}
.nav-left, .nav-center, .nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.nav-center {
  gap: 20px;
}
.nickname {
  color: #ffd700;
  font-weight: bold;
  font-size: 16px;
}
.gold, .rank {
  color: #ccc;
  font-size: 14px;
}
.btn-back {
  background: rgba(255,255,255,0.1);
  border: 1px solid #3a3a5e;
  color: #ccc;
  padding: 6px 14px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}
.btn-back:hover { background: rgba(255,255,255,0.2); }
.btn-icon {
  background: rgba(255,255,255,0.08);
  border: 1px solid #3a3a5e;
  color: #ccc;
  width: 38px;
  height: 38px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.btn-icon:hover { background: rgba(255,255,255,0.2); }
</style>

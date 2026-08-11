<template>
  <div class="lobby">
    <!-- 顶部导航 -->
    <GlobalNav
      showBack="false"
      @openSettings="showSettings = true"
      @openFriends="showFriends = true"
    />

    <!-- 主体：5大圆形入口 -->
    <main class="lobby-main">
      <ArtPlaceholder label="大厅背景图" sub="全屏装饰" width="100%" height="100%" bgColor="#0f0f1e" class="lobby-bg" />
      
      <div class="entries">
        <button v-for="(entry, idx) in entries" :key="idx" class="entry-btn" @click="entry.action">
          <div class="entry-circle">
            <ArtPlaceholder :label="entry.icon" width="100" height="100" bgColor="#1a1a3e" />
          </div>
          <span class="entry-label">{{ entry.label }}</span>
          <span class="entry-sub">{{ entry.sub }}</span>
        </button>
      </div>
    </main>

    <!-- 设置弹窗 -->
    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    
    <!-- 好友列表弹窗 -->
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const router = useRouter()
const userStore = useUserStore()
const showSettings = ref(false)
const showFriends = ref(false)

const entries = [
  { label: 'PVP 房间', sub: '在线对战', icon: 'PVP', action: () => router.push('/pvp-rooms') },
  { label: 'PVE 闯关', sub: '副本挑战', icon: 'PVE', action: () => router.push('/pve-stages') },
  { label: '商城', sub: '英雄/皮肤', icon: '商城', action: () => router.push('/shop') },
  { label: '教程', sub: '玩法教学', icon: '教程', action: () => router.push('/tutorial') },
  { label: '英雄图鉴', sub: '全英雄', icon: '图鉴', action: () => router.push('/herodex') },
]

onMounted(async () => {
  try {
    await userStore.fetchUserInfo()
  } catch { /* ignore */ }
})
</script>

<style scoped>
.lobby {
  width: 100vw; height: 100vh;
  display: flex; flex-direction: column;
  overflow: hidden;
}
.lobby-main {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
.lobby-bg {
  position: absolute; inset: 0; z-index: 0;
}
.entries {
  position: relative; z-index: 1;
  display: flex; gap: 36px;
}
.entry-btn {
  background: none; border: none; cursor: pointer;
  display: flex; flex-direction: column; align-items: center;
  gap: 8px; padding: 16px; border-radius: 16px;
  transition: all 0.25s;
}
.entry-btn:hover {
  background: rgba(255,255,255,0.06);
  transform: translateY(-4px);
}
.entry-circle {
  width: 100px; height: 100px; border-radius: 50%;
  overflow: hidden; border: 3px solid #3a3a5e;
  transition: border-color 0.3s;
}
.entry-btn:hover .entry-circle { border-color: #ffd700; }
.entry-label {
  color: #e0e0e0; font-size: 16px; font-weight: bold;
}
.entry-sub {
  color: #777; font-size: 12px;
}
</style>

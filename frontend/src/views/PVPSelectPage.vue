<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/pvp-rooms')">← 退出房间</button>
      </template>
      <template #center>
        <span class="countdown">{{ countdown > 0 ? `选人中: ${countdown}s` : '已锁定' }}</span>
      </template>
    </GlobalNav>

    <main class="content">
      <!-- 左侧英雄列表 -->
      <div class="hero-grid">
        <div v-for="h in heroes" :key="h.id" 
          :class="['hero-cell', { selected: selectedHero?.id === h.id }]"
          @click="selectHero(h)"
        >
          <ArtPlaceholder :label="h.name" width="80" height="80" bgColor="#1a1a3e" />
        </div>
      </div>

      <!-- 右侧预览 -->
      <div class="preview">
        <ArtPlaceholder :label="selectedHero ? selectedHero.name : '选择英雄'" sub="高清立绘" height="260" bgColor="#1a1a2e" />
        
        <!-- 皮肤列表 -->
        <div class="skins" v-if="selectedHero">
          <span class="section-title">皮肤</span>
          <div class="skin-grid">
            <div v-for="s in heroSkins" :key="s.id"
              :class="['skin-cell', { owned: s.owned, active: selectedSkin?.id === s.id }]"
              @click="selectSkin(s)"
            >
              <ArtPlaceholder :label="s.name" width="64" height="64" bgColor="#1a1a3e" />
              <span v-if="!s.owned" class="lock-icon">🔒</span>
            </div>
          </div>
        </div>

        <!-- 技能预览 -->
        <div class="skills" v-if="selectedHero">
          <span class="section-title">技能预览</span>
          <div v-for="sk in selectedHero.skills || []" :key="sk.id" class="skill-row">
            <ArtPlaceholder :label="sk.name" width="40" height="40" bgColor="#2a2a4a" />
            <span>{{ sk.name }}</span>
          </div>
        </div>
      </div>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '@/api'
import { useWebSocket } from '@/utils/websocket'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const router = useRouter()
const route = useRoute()
const showSettings = ref(false)
const showFriends = ref(false)
const countdown = ref(10)
const heroes = ref<any[]>([])
const selectedHero = ref<any>(null)
const selectedSkin = ref<any>(null)
const heroSkins = ref<any[]>([])
let timer: any = null

const { send, on, off } = useWebSocket()

onMounted(async () => {
  // 加载全部英雄
  try {
    const res: any = await api.get('/characters')
    heroes.value = (res.data || res) || []
  } catch { /* 使用假数据 */ }

  // 倒计时
  timer = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--
    } else {
      clearInterval(timer)
      lockIn()
    }
  }, 1000)

  on('game_start', handleGameStart)
})

onUnmounted(() => {
  clearInterval(timer)
  off('game_start', handleGameStart)
})

function selectHero(h: any) {
  selectedHero.value = h
  selectedSkin.value = null
  // 加载该英雄的皮肤
  heroSkins.value = [
    { id: 1, name: '默认', owned: true },
    { id: 2, name: h.name + '皮肤1', owned: false },
    { id: 3, name: h.name + '皮肤2', owned: false },
  ]
}

function selectSkin(s: any) {
  if (!s.owned) {
    alert('请前往商城购买此皮肤')
    return
  }
  selectedSkin.value = s
}

function lockIn() {
  if (!selectedHero.value) {
    // 随机选
    if (heroes.value.length > 0) selectHero(heroes.value[0])
  }
  send({ type: 'hero_selected', heroId: selectedHero.value?.id, skinId: selectedSkin.value?.id })
}

function handleGameStart(data: any) {
  router.push(`/game?mode=pvp&battleId=${data.battleId || ''}`)
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; display: flex; gap: 16px; padding: 16px; background: #0f0f1e; }
.hero-grid {
  flex: 1; display: grid; grid-template-columns: repeat(5, 1fr);
  gap: 8px; overflow-y: auto; align-content: start;
}
.hero-cell { cursor: pointer; border-radius: 10px; border: 2px solid transparent; transition: all .2s; }
.hero-cell.selected { border-color: #ffd700; }
.preview {
  width: 320px; display: flex; flex-direction: column; gap: 14px;
  background: rgba(20,20,40,.8); border: 1px solid #3a3a5e; border-radius: 12px; padding: 16px;
}
.countdown { color: #ffd700; font-size: 18px; font-weight: bold; }
.section-title { color: #ffd700; font-size: 13px; display: block; margin-bottom: 6px; }
.skin-grid { display: flex; gap: 8px; }
.skin-cell { position: relative; cursor: pointer; border: 2px solid transparent; border-radius: 8px; }
.skin-cell.owned { border-color: #3a3a5e; }
.skin-cell.active { border-color: #ffd700; }
.lock-icon { position: absolute; top: 50%; left: 50%; transform: translate(-50%,-50%); font-size: 20px; }
.skill-row { display: flex; align-items: center; gap: 8px; color: #888; font-size: 13px; }
.btn-back { background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc; padding: 6px 16px; border-radius: 6px; cursor: pointer; }
</style>

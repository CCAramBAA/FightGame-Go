<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/lobby')">← 返回大厅</button>
      </template>
      <template #center>
        <div class="replay-input">
          <input v-model="replayId" class="input-sm" placeholder="输入回放 ID" />
          <button class="btn-load" @click="loadReplay">加载回放</button>
        </div>
      </template>
    </GlobalNav>

    <main class="content">
      <ArtPlaceholder label="战斗回放画布" sub="Phaser 回放画布" width="100%" height="100%" bgColor="#0a0a1a" />
      
      <!-- 回放控制栏 -->
      <div class="controls">
        <button @click="togglePlay">{{ playing ? '⏸ 暂停' : '▶ 播放' }}</button>
        <button @click="cycleSpeed">{{ speedLabel }}</button>
        <button @click="prevFrame">⏮ 上一帧</button>
        <button @click="nextFrame">⏭ 下一帧</button>
        <span class="frame-info">帧: {{ currentFrame }} / {{ totalFrames }}</span>
      </div>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import api from '@/api'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const showSettings = ref(false)
const showFriends = ref(false)
const replayId = ref('')
const playing = ref(false)
const currentFrame = ref(0)
const totalFrames = ref(0)
const speed = ref(1)

const speedLabel = ref('1x')

function cycleSpeed() {
  if (speed.value === 1) { speed.value = 2; speedLabel.value = '2x' }
  else if (speed.value === 2) { speed.value = 4; speedLabel.value = '4x' }
  else { speed.value = 1; speedLabel.value = '1x' }
}

function togglePlay() { playing.value = !playing.value }
function prevFrame() { if (currentFrame.value > 0) currentFrame.value-- }
function nextFrame() { if (currentFrame.value < totalFrames.value) currentFrame.value++ }

async function loadReplay() {
  if (!replayId.value) { alert('请输入回放 ID'); return }
  try {
    const res: any = await api.get(`/battle/replay/${replayId.value}`)
    const data = (res as any).data || res
    totalFrames.value = data.frames?.length || 0
    currentFrame.value = 0
  } catch {
    alert('加载回放失败')
  }
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; position: relative; }
.controls {
  position: absolute; bottom: 0; left: 0; right: 0;
  display: flex; gap: 10px; padding: 12px 20px; justify-content: center; align-items: center;
  background: rgba(0,0,0,.8); border-top: 1px solid #2a2a4a;
}
.controls button {
  padding: 8px 16px; background: #2a2a4a; border: 1px solid #3a3a5e;
  border-radius: 6px; color: #ccc; cursor: pointer; font-size: 13px;
}
.frame-info { color: #888; font-size: 13px; }
.replay-input { display: flex; gap: 8px; }
.input-sm { padding: 6px 12px; background: #2a2a4a; border: 1px solid #3a3a5e; border-radius: 6px; color: #ccc; width: 160px; }
.btn-load { padding: 6px 14px; background: #ffd70033; border: 1px solid #ffd700; border-radius: 6px; color: #ffd700; cursor: pointer; }
.btn-back { background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc; padding: 6px 16px; border-radius: 6px; cursor: pointer; }
</style>

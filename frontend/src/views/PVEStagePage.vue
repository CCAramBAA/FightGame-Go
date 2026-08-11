<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/lobby')">← 返回大厅</button>
      </template>
    </GlobalNav>

    <main class="content">
      <ArtPlaceholder label="闯关大地图背景" width="100%" height="100%" bgColor="#0a1a0e" class="bg" />

      <!-- 关卡区域 -->
      <div class="stages">
        <!-- 简单区 -->
        <div class="stage-zone zone-easy">
          <span class="zone-title">简单关卡</span>
          <div class="stage-dots">
            <button v-for="s in easyStages" :key="s.id" class="stage-dot easy" @click="enterStage(s)">
              {{ s.level }}
            </button>
          </div>
        </div>
        
        <!-- 普通区 -->
        <div class="stage-zone zone-normal">
          <span class="zone-title">普通关卡</span>
          <div class="stage-dots">
            <button v-for="s in normalStages" :key="s.id" class="stage-dot normal" @click="enterStage(s)">
              {{ s.level }}
            </button>
          </div>
        </div>

        <!-- 困难区 -->
        <div class="stage-zone zone-hard">
          <span class="zone-title">困难关卡</span>
          <div class="stage-dots">
            <button v-for="s in hardStages" :key="s.id" :class="['stage-dot', 'hard', { locked: s.locked }]" @click="enterStage(s)">
              {{ s.level }}
              <span v-if="s.locked" class="lock">🔒</span>
            </button>
          </div>
        </div>

        <!-- BOSS 区 -->
        <div class="stage-zone zone-boss">
          <span class="zone-title">BOSS 关卡</span>
          <div class="stage-dots">
            <button v-for="s in bossStages" :key="s.id" :class="['stage-dot', 'boss', { locked: s.locked }]" @click="enterStage(s)">
              {{ s.level }}
              <span v-if="s.locked" class="lock">🔒</span>
            </button>
          </div>
        </div>
      </div>

      <button class="btn-train" @click="startTraining">训练模式</button>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useUserStore } from '@/store/user'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const router = useRouter()
const userStore = useUserStore()
const showSettings = ref(false)
const showFriends = ref(false)

const stages = ref<any[]>([])

const easyStages = computed(() => stages.value.filter((s: any) => s.difficulty === 'easy'))
const normalStages = computed(() => stages.value.filter((s: any) => s.difficulty === 'normal'))
const hardStages = computed(() => stages.value.filter((s: any) => s.difficulty === 'hard'))
const bossStages = computed(() => stages.value.filter((s: any) => s.difficulty === 'boss'))

onMounted(async () => {
  try {
    const res: any = await api.get('/pve/stages')
    stages.value = (res.data || res) || []
  } catch {
    // 假数据
    stages.value = Array.from({ length: 10 }, (_, i) => ({
      id: i + 1, level: i + 1,
      difficulty: i < 4 ? 'easy' : i < 7 ? 'normal' : i < 9 ? 'hard' : 'boss',
      locked: false
    }))
  }
})

function enterStage(stage: any) {
  if (stage.locked) { alert('此关卡未解锁'); return }
  router.push(`/pve-select?stageId=${stage.id}`)
}

function startTraining() {
  router.push('/game?mode=training')
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; position: relative; padding: 20px; }
.bg { position: absolute; inset: 0; z-index: 0; }
.stages { position: relative; z-index: 1; display: flex; flex-direction: column; gap: 24px; }
.stage-zone { display: flex; align-items: center; gap: 16px; }
.zone-title { color: #ccc; font-size: 14px; min-width: 80px; }
.stage-dots { display: flex; gap: 10px; flex-wrap: wrap; }
.stage-dot {
  width: 48px; height: 48px; border-radius: 12px; border: 2px solid transparent;
  cursor: pointer; font-size: 14px; font-weight: bold; position: relative;
  display: flex; align-items: center; justify-content: center; transition: all .2s;
}
.stage-dot.easy { background: #22c55e44; border-color: #22c55e; color: #22c55e; }
.stage-dot.normal { background: #f59e0b44; border-color: #f59e0b; color: #f59e0b; }
.stage-dot.hard { background: #dc262644; border-color: #dc2626; color: #dc2626; }
.stage-dot.boss { background: #9333ea44; border-color: #9333ea; color: #9333ea; }
.stage-dot.locked { background: #2a2a4a; border-color: #666; color: #666; cursor: not-allowed; }
.stage-dot:hover:not(.locked) { transform: scale(1.1); }
.lock { font-size: 12px; position: absolute; }
.btn-train {
  position: fixed; bottom: 30px; left: 30px; z-index: 10;
  padding: 12px 22px; border-radius: 50px; border: 2px solid #3a3a5e;
  background: rgba(20,20,40,.9); color: #ccc; cursor: pointer; font-size: 14px;
}
.btn-back { background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc; padding: 6px 16px; border-radius: 6px; cursor: pointer; }
</style>

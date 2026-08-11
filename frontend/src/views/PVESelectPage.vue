<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/pve-stages')">← 返回关卡地图</button>
      </template>
    </GlobalNav>

    <main class="content">
      <!-- 左侧英雄列表（有解锁限制） -->
      <div class="hero-grid">
        <div v-for="h in heroes" :key="h.id" 
          :class="['hero-cell', { selected: selectedHero?.id === h.id, locked: !h.unlocked }]"
          @click="h.unlocked ? selectHero(h) : alert('请前往商城购买或通关指定关卡解锁')"
        >
          <ArtPlaceholder :label="h.name" width="80" height="80" bgColor="#1a1a3e" />
          <span v-if="!h.unlocked" class="lock-icon">🔒</span>
        </div>
      </div>

      <!-- 右侧预览 -->
      <div class="preview">
        <ArtPlaceholder :label="selectedHero ? selectedHero.name : '选择英雄'" sub="高清立绘" height="260" bgColor="#1a1a2e" />

        <div class="skins" v-if="selectedHero">
          <span class="section-title">皮肤</span>
          <div v-for="s in heroSkins" :key="s.id"
            :class="['skin-cell', { owned: s.owned, active: selectedSkin?.id === s.id }]"
            @click="selectSkin(s)"
          >
            <ArtPlaceholder :label="s.name" width="64" height="64" bgColor="#1a1a3e" />
            <span v-if="!s.owned">🔒</span>
          </div>
        </div>

        <div class="skills" v-if="selectedHero">
          <span class="section-title">技能预览</span>
          <div v-for="sk in selectedHero.skills || []" :key="sk.id" class="skill-row">
            <ArtPlaceholder :label="sk.name" width="40" height="40" bgColor="#2a2a4a" />
            <span>{{ sk.name }}</span>
          </div>
        </div>

        <button v-if="selectedHero" class="btn-start" @click="startBattle">开始战斗</button>
      </div>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '@/api'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const router = useRouter()
const route = useRoute()
const showSettings = ref(false)
const showFriends = ref(false)
const heroes = ref<any[]>([])
const selectedHero = ref<any>(null)
const selectedSkin = ref<any>(null)
const heroSkins = ref<any[]>([])

onMounted(async () => {
  try {
    const res: any = await api.get('/profile/characters')
    heroes.value = (res.data || res) || []
  } catch {
    // 假数据：前3个英雄已解锁
    heroes.value = [
      { id: 1, name: '火焰战士', unlocked: true, skills: [{ id: 1, name: '火球术' }, { id: 2, name: '烈焰斩' }] },
      { id: 2, name: '寒冰法师', unlocked: true, skills: [{ id: 3, name: '冰锥术' }, { id: 4, name: '暴风雪' }] },
      { id: 3, name: '暗影刺客', unlocked: true, skills: [{ id: 5, name: '暗影步' }, { id: 6, name: '致命一击' }] },
      { id: 4, name: '圣光骑士', unlocked: false, skills: [{ id: 7, name: '圣光斩' }] },
      { id: 5, name: '风行者', unlocked: false, skills: [{ id: 8, name: '飓风箭' }] },
    ]
  }
})

function selectHero(h: any) {
  selectedHero.value = h
  selectedSkin.value = null
  heroSkins.value = [
    { id: 1, name: '默认', owned: true },
    { id: 2, name: h.name + '皮肤1', owned: false },
  ]
}

function selectSkin(s: any) {
  if (!s.owned) { alert('请前往商城购买此皮肤'); return }
  selectedSkin.value = s
}

function startBattle() {
  router.push(`/game?mode=pve&stageId=${route.query.stageId}`)
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; display: flex; gap: 16px; padding: 16px; background: #0f0f1e; }
.hero-grid { flex: 1; display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; overflow-y: auto; align-content: start; }
.hero-cell { position: relative; cursor: pointer; border-radius: 10px; border: 2px solid transparent; }
.hero-cell.selected { border-color: #ffd700; }
.hero-cell.locked { opacity: 0.4; cursor: not-allowed; }
.lock-icon { position: absolute; top: 50%; left: 50%; transform: translate(-50%,-50%); font-size: 24px; }
.preview { width: 320px; display: flex; flex-direction: column; gap: 14px; background: rgba(20,20,40,.8); border: 1px solid #3a3a5e; border-radius: 12px; padding: 16px; }
.section-title { color: #ffd700; font-size: 13px; margin-bottom: 6px; }
.skins { display: flex; gap: 8px; flex-wrap: wrap; }
.skin-cell { position: relative; cursor: pointer; border: 2px solid #3a3a5e; border-radius: 8px; }
.skin-cell.active { border-color: #ffd700; }
.skill-row { display: flex; align-items: center; gap: 8px; color: #888; font-size: 13px; }
.btn-start { padding: 12px; border-radius: 10px; border: none; background: linear-gradient(135deg,#ffd700,#ff8c00); color: #1a1a1a; font-size: 16px; font-weight: bold; cursor: pointer; }
.btn-back { background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc; padding: 6px 16px; border-radius: 6px; cursor: pointer; }
</style>

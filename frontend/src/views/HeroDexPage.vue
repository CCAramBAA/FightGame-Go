<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/lobby')">← 返回大厅</button>
      </template>
    </GlobalNav>

    <main class="content">
      <ArtPlaceholder label="图鉴背景" width="100%" height="100%" bgColor="#0f0f1e" class="bg" />
      
      <!-- 左侧英雄列表 -->
      <div class="hero-grid">
        <div v-for="h in heroes" :key="h.id" 
          :class="['hero-cell', { active: selectedHero?.id === h.id }]"
          @click="selectHero(h)"
        >
          <ArtPlaceholder :label="h.name" width="90" height="90" bgColor="#1a1a3e" />
        </div>
      </div>

      <!-- 右侧详情面板 -->
      <div class="detail-panel" v-if="selectedHero">
        <ArtPlaceholder :label="selectedHero.name" sub="高清立绘" height="200" bgColor="#1a1a2e" />
        
        <section>
          <h4>技能参数</h4>
          <div v-for="sk in selectedHero.skills || []" :key="sk.id" class="skill-detail">
            <ArtPlaceholder :label="sk.name" width="40" height="40" bgColor="#2a2a4a" />
            <div>
              <div class="sk-name">{{ sk.name }}</div>
              <div class="sk-stats">伤害: {{ sk.damage || '-' }} | CD: {{ sk.cooldown || '-' }}s | 消耗: {{ sk.cost || '-' }}</div>
            </div>
          </div>
        </section>

        <section>
          <h4>技能优先级</h4>
          <div class="priority">{{ selectedHero.priority || '无特殊优先级' }}</div>
        </section>

        <section>
          <h4>特殊标签</h4>
          <div class="tags">
            <span v-for="t in selectedHero.tags || []" :key="t" class="tag">{{ t }}</span>
            <span v-if="!selectedHero.tags?.length" class="no-tag">无</span>
          </div>
        </section>

        <section>
          <h4>背景故事</h4>
          <p class="lore">{{ selectedHero.lore || '暂无背景故事' }}</p>
        </section>
      </div>

      <div v-else class="no-select">
        <p>← 请选择一位英雄查看详情</p>
      </div>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/api'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const showSettings = ref(false)
const showFriends = ref(false)
const heroes = ref<any[]>([])
const selectedHero = ref<any>(null)

onMounted(async () => {
  try {
    const res: any = await api.get('/characters')
    heroes.value = (res.data || res) || []
  } catch {
    heroes.value = [
      { id: 1, name: '火焰战士', tags: ['近战', '高爆发'], lore: '来自火山之地的勇者...', priority: '攻击力 > 暴击', skills: [{ id: 1, name: '火球术', damage: 80, cooldown: 3, cost: 20 }, { id: 2, name: '烈焰斩', damage: 150, cooldown: 8, cost: 50 }] },
      { id: 2, name: '寒冰法师', tags: ['远程', '控制'], lore: '冰雪王国的守护者...', priority: '闪避 > 生命值', skills: [{ id: 3, name: '冰锥术', damage: 60, cooldown: 2, cost: 15 }, { id: 4, name: '暴风雪', damage: 120, cooldown: 10, cost: 60 }] },
      { id: 3, name: '暗影刺客', tags: ['近战', '刺客'], lore: '来自暗影教团的杀手...', priority: '暴击 > 攻击速度', skills: [{ id: 5, name: '暗影步', damage: 40, cooldown: 4, cost: 25 }, { id: 6, name: '致命一击', damage: 200, cooldown: 15, cost: 80 }] },
    ]
  }
})

function selectHero(h: any) {
  selectedHero.value = h
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; position: relative; display: flex; gap: 16px; padding: 16px; }
.bg { position: absolute; inset: 0; z-index: 0; }
.hero-grid {
  position: relative; z-index: 1;
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; align-content: start;
  width: 300px; height: 100%; overflow-y: auto;
}
.hero-cell { cursor: pointer; border: 2px solid transparent; border-radius: 10px; transition: all .2s; }
.hero-cell.active { border-color: #ffd700; }
.detail-panel {
  flex: 1; position: relative; z-index: 1;
  background: rgba(20,20,40,.85); border: 1px solid #3a3a5e; border-radius: 12px;
  padding: 16px; overflow-y: auto; display: flex; flex-direction: column; gap: 16px;
}
.detail-panel h4 { color: #ffd700; margin: 0 0 8px; font-size: 14px; }
.skill-detail { display: flex; gap: 10px; margin-bottom: 8px; }
.sk-name { color: #e0e0e0; font-size: 14px; }
.sk-stats { color: #888; font-size: 12px; }
.priority { color: #4fc3f7; font-size: 14px; }
.tags { display: flex; gap: 6px; }
.tag { padding: 3px 10px; background: #3a3a5e; border-radius: 12px; color: #ccc; font-size: 12px; }
.no-tag { color: #666; }
.lore { color: #999; font-size: 13px; line-height: 1.6; }
.no-select { flex: 1; display: flex; align-items: center; justify-content: center; z-index: 1; }
.no-select p { color: #666; font-size: 16px; }
.btn-back { background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc; padding: 6px 16px; border-radius: 6px; cursor: pointer; }
</style>

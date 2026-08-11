<template>
  <ModalOverlay :visible="visible" :showClose="false" maxWidth="480px">
    <div class="result">
      <!-- 上半区 -->
      <div class="result-top">
        <div class="winner">{{ resultText }}</div>
        <div class="hp-compare">
          <div class="player">
            <span class="label">你</span>
            <div class="hp-bar-bg"><div class="hp-bar" :style="{ width: playerHp + '%' }"></div></div>
            <span>{{ playerHp }}%</span>
          </div>
          <span class="vs">VS</span>
          <div class="enemy">
            <span class="label">对手</span>
            <div class="hp-bar-bg"><div class="hp-bar enemy-bar" :style="{ width: enemyHp + '%' }"></div></div>
            <span>{{ enemyHp }}%</span>
          </div>
        </div>
      </div>

      <!-- 下半区 -->
      <div class="result-bottom">
        <div v-if="mode === 'pvp'" class="stat">
          <span>段位积分</span>
          <span :class="['score', rankChange >= 0 ? 'up' : 'down']">
            {{ rankChange >= 0 ? '+' : '' }}{{ rankChange }}
          </span>
        </div>
        <div v-if="mode === 'pve'" class="stat">
          <span>金币奖励</span>
          <span class="up">+{{ goldReward }}</span>
        </div>
        <div v-if="mode === 'pve'" class="stat">
          <span>星级评价</span>
          <span class="stars">{{ '★'.repeat(stars) }}{{ '☆'.repeat(3 - stars) }}</span>
        </div>
      </div>
    </div>

    <template #footer>
      <button class="btn-replay" @click="$emit('replay')">查看本局回放</button>
      <button class="btn-lobby" @click="$emit('lobby')">返回大厅</button>
    </template>
  </ModalOverlay>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import ModalOverlay from './ModalOverlay.vue'

const props = defineProps<{
  visible: boolean
  mode: 'pvp' | 'pve' | 'training'
  won?: boolean
  playerHp: number
  enemyHp: number
  rankChange?: number
  goldReward?: number
  stars?: number
}>()

defineEmits<{ replay: []; lobby: [] }>()

const resultText = computed(() => {
  if (props.playerHp === props.enemyHp) return '平局'
  return props.won ? '胜利!' : '失败'
})
</script>

<style scoped>
.result { display: flex; flex-direction: column; gap: 20px; }
.result-top { text-align: center; display: flex; flex-direction: column; gap: 16px; }
.winner { font-size: 28px; font-weight: bold; color: #ffd700; }
.hp-compare { display: flex; align-items: center; gap: 12px; }
.player, .enemy { flex: 1; display: flex; align-items: center; gap: 6px; }
.label { color: #888; font-size: 12px; }
.hp-bar-bg { flex: 1; height: 12px; background: #2a2a4a; border-radius: 6px; overflow: hidden; }
.hp-bar { height: 100%; background: #22c55e; border-radius: 6px; }
.enemy-bar { background: #dc2626; }
.vs { color: #666; font-size: 14px; }
.player span, .enemy span { color: #ccc; font-size: 13px; min-width: 36px; }
.result-bottom { display: flex; flex-direction: column; gap: 8px; }
.stat { display: flex; justify-content: space-between; padding: 6px 0; border-bottom: 1px solid #2a2a4a; color: #ccc; }
.stat .up { color: #22c55e; font-weight: bold; }
.stat .down { color: #dc2626; font-weight: bold; }
.stars { color: #ffd700; letter-spacing: 2px; }
.btn-replay {
  padding: 10px 20px; background: #2a2a4a; border: 1px solid #3a3a5e;
  border-radius: 8px; color: #ccc; cursor: pointer; font-size: 14px;
}
.btn-lobby {
  padding: 10px 20px; background: linear-gradient(135deg,#ffd700,#ff8c00);
  border: none; border-radius: 8px; color: #1a1a1a; cursor: pointer; font-size: 14px; font-weight: bold;
}
</style>

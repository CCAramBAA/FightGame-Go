<template>
  <div class="game-page">
    <div ref="gameContainer" class="game-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { FightGame } from '@/game/FightGame'

const gameContainer = ref<HTMLDivElement>()
let game: FightGame | null = null

onMounted(() => {
  if (gameContainer.value) {
    game = new FightGame(gameContainer.value.id || 'game-container')
  }
})

onUnmounted(() => {
  if (game) {
    game.destroy(true)
    game = null
  }
})
</script>

<style scoped>
.game-page {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background: #1a1a2e;
}
.game-container {
  width: 100%;
  height: 100%;
}
</style>

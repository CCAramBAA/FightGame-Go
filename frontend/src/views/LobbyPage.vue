<template>
  <div class="lobby-page">
    <h2>游戏大厅</h2>
    <div class="room-list">
      <p>暂无可用房间</p>
    </div>
    <button class="btn" @click="createRoom">创建房间</button>
  </div>
</template>

<script setup lang="ts">
import { useWebSocket } from '@/utils/websocket'

const { send } = useWebSocket()

function createRoom() {
  send({ type: 'create_room' })
}
</script>

<style scoped>
.lobby-page {
  padding: 2rem;
  color: #fff;
  background: #1a1a2e;
  min-height: 100vh;
}
.room-list {
  margin: 2rem 0;
  padding: 2rem;
  background: #16213e;
  border-radius: 8px;
  min-height: 200px;
}
.btn {
  padding: 0.8rem 2rem;
  background: #e94560;
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}
</style>

<template>
  <div class="home-page">
    <header class="header">
      <h1>FightGame</h1>
      <p class="subtitle">在线对战游戏平台</p>
    </header>
    <main class="main">
      <div class="actions">
        <button class="btn btn-primary" @click="$router.push('/lobby')">进入大厅</button>
        <button class="btn btn-secondary" @click="$router.push('/game')">快速对战</button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  color: #fff;
}
.header {
  text-align: center;
  margin-bottom: 3rem;
}
.header h1 {
  font-size: 3.5rem;
  margin: 0;
  background: linear-gradient(90deg, #e94560, #f39c12);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.subtitle {
  font-size: 1.2rem;
  color: #a0a0a0;
  margin-top: 0.5rem;
}
.actions {
  display: flex;
  gap: 1.5rem;
}
.btn {
  padding: 0.8rem 2.5rem;
  font-size: 1.1rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 600;
}
.btn-primary {
  background: #e94560;
  color: #fff;
}
.btn-primary:hover {
  background: #c43a52;
  transform: translateY(-2px);
}
.btn-secondary {
  background: transparent;
  color: #e94560;
  border: 2px solid #e94560;
}
.btn-secondary:hover {
  background: rgba(233, 69, 96, 0.1);
  transform: translateY(-2px);
}
</style>

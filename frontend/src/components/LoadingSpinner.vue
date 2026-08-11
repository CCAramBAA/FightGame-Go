<template>
  <div class="loading-spinner" :class="{ fullscreen }">
    <div class="spinner-ring"></div>
    <p v-if="text" class="spinner-text">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{ text?: string; fullscreen?: boolean }>()
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 40px;
}
.loading-spinner.fullscreen {
  position: fixed;
  inset: 0;
  background: #0a0a1a;
  z-index: 9999;
}
.spinner-ring {
  width: 48px;
  height: 48px;
  border: 4px solid #2a2a4a;
  border-top-color: #ffd700;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
.spinner-text {
  color: #888;
  font-size: 14px;
  margin: 0;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>

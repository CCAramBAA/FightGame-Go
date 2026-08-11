<template>
  <ModalOverlay :visible="true" :showClose="false" maxWidth="800px">
    <div class="tutorial">
      <ArtPlaceholder label="教程动画播放窗口" sub="分段引导动画" height="400" bgColor="#1a1a2e" />
      
      <!-- 进度条 -->
      <input type="range" :min="0" :max="totalSteps - 1" v-model.number="currentStep" class="slider" />

      <div class="controls">
        <button class="btn-ctrl" @click="prev" :disabled="currentStep <= 0">上一页</button>
        <button class="btn-ctrl btn-skip" @click="skip" :disabled="!canSkip">跳过</button>
        <button class="btn-ctrl" @click="next" :disabled="currentStep >= totalSteps - 1">下一页</button>
        <button v-if="currentStep >= totalSteps - 1" class="btn-ctrl btn-done" @click="finish">完成</button>
      </div>
    </div>
  </ModalOverlay>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import ModalOverlay from '@/components/ModalOverlay.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'

const router = useRouter()
const totalSteps = 5
const currentStep = ref(0)
const isFirstTime = ref(!localStorage.getItem('tutorial_done'))

// 首次登录不可跳过
const canSkip = computed(() => !isFirstTime.value)

function prev() { if (currentStep.value > 0) currentStep.value-- }
function next() { if (currentStep.value < totalSteps - 1) currentStep.value++ }
function skip() { if (canSkip.value) finish() }
function finish() {
  localStorage.setItem('tutorial_done', '1')
  router.replace('/lobby')
}
</script>

<style scoped>
.tutorial { display: flex; flex-direction: column; gap: 16px; }
.slider { width: 100%; accent-color: #ffd700; }
.controls { display: flex; gap: 10px; justify-content: center; }
.btn-ctrl {
  padding: 8px 20px; border-radius: 8px; border: 1px solid #3a3a5e;
  background: rgba(255,255,255,0.08); color: #ccc; cursor: pointer; font-size: 14px;
}
.btn-ctrl:hover:not(:disabled) { background: rgba(255,255,255,0.15); }
.btn-ctrl:disabled { opacity: 0.35; cursor: not-allowed; }
.btn-skip { border-color: #888; }
.btn-done { background: linear-gradient(135deg, #ffd700, #ff8c00); border-color: #ffd700; color: #1a1a1a; font-weight: bold; }
</style>

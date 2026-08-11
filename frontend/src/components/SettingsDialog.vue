<template>
  <ModalOverlay :visible="visible" title="系统设置" maxWidth="520px" @close="saveAndClose">
    <div class="settings">
      <!-- 按键自定义 -->
      <section>
        <h4>按键自定义</h4>
        <div class="key-row" v-for="k in keyBindings" :key="k.action">
          <span>{{ k.label }}</span>
          <button class="btn-key" @click="bindKey(k)">{{ k.key }}</button>
        </div>
      </section>
      
      <!-- 画质 -->
      <section>
        <h4>画质</h4>
        <select v-model="quality" class="select">
          <option value="low">低</option>
          <option value="medium" selected>中</option>
          <option value="high">高</option>
        </select>
      </section>
      
      <!-- 音量 -->
      <section>
        <h4>音量</h4>
        <label>BGM <input type="range" v-model.number="bgmVolume" min="0" max="100" /> {{ bgmVolume }}%</label>
        <label>音效 <input type="range" v-model.number="sfxVolume" min="0" max="100" /> {{ sfxVolume }}%</label>
      </section>
      
      <!-- 手柄 -->
      <section>
        <h4>游戏手柄适配</h4>
        <button class="btn-small" @click="detectGamepad">检测手柄</button>
        <span v-if="gamepadDetected" class="ok">已检测到手柄</span>
      </section>
      
      <!-- 快捷按钮 -->
      <section>
        <button class="btn-link" @click="goToHeroDex">前往英雄图鉴 →</button>
      </section>
    </div>
    <template #footer>
      <button class="btn-ok" @click="saveAndClose">关闭</button>
    </template>
  </ModalOverlay>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import ModalOverlay from './ModalOverlay.vue'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()
const router = useRouter()

const quality = ref('medium')
const bgmVolume = ref(80)
const sfxVolume = ref(80)
const gamepadDetected = ref(false)

const keyBindings = reactive([
  { action: 'attack1', label: '轻攻击', key: 'J' },
  { action: 'attack2', label: '重攻击', key: 'K' },
  { action: 'skill1', label: '技能1', key: 'U' },
  { action: 'skill2', label: '技能2', key: 'I' },
  { action: 'skill3', label: '技能3', key: 'O' },
  { action: 'skill4', label: '技能4', key: 'P' },
])

function bindKey(k: typeof keyBindings[number]) {
  const newKey = prompt(`按下新按键替换 ${k.key}:`, k.key)
  if (newKey && newKey.length === 1) k.key = newKey.toUpperCase()
  saveSettings()
}

function detectGamepad() {
  gamepadDetected.value = navigator.getGamepads ? true : false
}

function goToHeroDex() {
  emit('close')
  router.push('/herodex')
}

function saveSettings() {
  localStorage.setItem('settings', JSON.stringify({ quality: quality.value, bgmVolume: bgmVolume.value, sfxVolume: sfxVolume.value, keyBindings }))
}

function saveAndClose() {
  saveSettings()
  emit('close')
}
</script>

<style scoped>
.settings { display: flex; flex-direction: column; gap: 18px; }
.settings h4 { color: #ffd700; margin: 0 0 8px; font-size: 14px; }
section { display: flex; flex-direction: column; gap: 6px; }
.key-row { display: flex; justify-content: space-between; align-items: center; color: #ccc; }
.btn-key {
  background: #2a2a4a; border: 1px solid #3a3a5e; border-radius: 6px;
  color: #fff; padding: 6px 16px; cursor: pointer; font-size: 14px;
}
.select {
  padding: 8px; background: #2a2a4a; border: 1px solid #3a3a5e; color: #ccc; border-radius: 6px;
}
.btn-small {
  padding: 6px 14px; background: #2a2a4a; border: 1px solid #3a3a5e; color: #ccc; border-radius: 6px; cursor: pointer;
}
.btn-link { background: none; border: none; color: #4fc3f7; cursor: pointer; font-size: 14px; }
.btn-ok {
  background: linear-gradient(135deg, #ffd700, #ff8c00); border: none;
  padding: 10px 32px; border-radius: 8px; cursor: pointer; font-weight: bold;
}
.ok { color: #22c55e; font-size: 13px; }
</style>

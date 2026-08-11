<template>
  <div class="settings-page">
    <header class="page-header">
      <button class="btn-back" @click="$router.back()">← 返回</button>
      <h1>游戏设置</h1>
      <span class="hint">所有设置自动保存</span>
    </header>

    <div class="settings-container">
      <!-- 键盘按键设置 -->
      <section class="setting-section">
        <h2>操作按键</h2>
        <div class="key-grid">
          <div v-for="key in keyBindings" :key="key.action" class="key-item">
            <span class="key-label">{{ key.label }}</span>
            <button class="key-btn" :class="{ recording: recordingKey === key.action }" @click="startRecord(key.action)">
              {{ recordingKey === key.action ? '按下按键...' : key.displayKey }}
            </button>
            <button v-if="key.binding !== key.defaultKey" class="key-reset" @click="resetKey(key.action)">重置</button>
          </div>
        </div>
      </section>

      <!-- 画质设置 -->
      <section class="setting-section">
        <h2>画质</h2>
        <div class="radio-group">
          <label v-for="q in qualityOptions" :key="q.value" class="radio-label" :class="{ active: quality === q.value }">
            <input type="radio" v-model="quality" :value="q.value" @change="saveSettings" />
            <span>{{ q.label }}</span>
            <small>{{ q.desc }}</small>
          </label>
        </div>
      </section>

      <!-- 音量设置 -->
      <section class="setting-section">
        <h2>音量</h2>
        <div class="slider-group">
          <div class="slider-item">
            <span>BGM 音量</span>
            <input type="range" v-model.number="bgmVolume" min="0" max="100" @input="saveSettings" />
            <span class="val">{{ bgmVolume }}%</span>
          </div>
          <div class="slider-item">
            <span>技能音效</span>
            <input type="range" v-model.number="sfxVolume" min="0" max="100" @input="saveSettings" />
            <span class="val">{{ sfxVolume }}%</span>
          </div>
          <div class="slider-item">
            <span>UI 音效</span>
            <input type="range" v-model.number="uiVolume" min="0" max="100" @input="saveSettings" />
            <span class="val">{{ uiVolume }}%</span>
          </div>
        </div>
      </section>

      <!-- 手柄适配 -->
      <section class="setting-section">
        <h2>手柄</h2>
        <div class="gamepad-status">
          <div class="status-indicator" :class="{ connected: gamepadConnected }"></div>
          <span>{{ gamepadConnected ? '已连接: ' + gamepadName : '未检测到手柄' }}</span>
          <button v-if="!gamepadConnected" class="btn-detect" @click="detectGamepad">检测手柄</button>
        </div>
        <p class="hint-text">支持 Xbox / PlayStation 手柄，无需额外驱动</p>
      </section>

      <!-- 快速跳转 -->
      <section class="setting-section">
        <h2>快捷入口</h2>
        <div class="quick-links">
          <router-link to="/herodex" class="quick-link">英雄图鉴 →</router-link>
          <router-link to="/tutorial" class="quick-link">新手教程 →</router-link>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface KeyBinding {
  action: string
  label: string
  binding: string
  defaultKey: string
  displayKey: string
}

const keyBindings = ref<KeyBinding[]>([
  { action: 'move_up', label: '上移', binding: 'W', defaultKey: 'W', displayKey: 'W' },
  { action: 'move_down', label: '下移', binding: 'S', defaultKey: 'S', displayKey: 'S' },
  { action: 'move_left', label: '左移', binding: 'A', defaultKey: 'A', displayKey: 'A' },
  { action: 'move_right', label: '右移', binding: 'D', defaultKey: 'D', displayKey: 'D' },
  { action: 'attack', label: '普攻', binding: 'J', defaultKey: 'J', displayKey: 'J' },
  { action: 'skill_1', label: '技能1', binding: 'U', defaultKey: 'U', displayKey: 'U' },
  { action: 'skill_2', label: '技能2', binding: 'I', defaultKey: 'I', displayKey: 'I' },
  { action: 'skill_3', label: '技能3', binding: 'O', defaultKey: 'O', displayKey: 'O' },
  { action: 'dash', label: '闪避', binding: 'L', defaultKey: 'L', displayKey: 'L' },
  { action: 'skill_panel', label: '技能面板', binding: 'Tab', defaultKey: 'Tab', displayKey: 'Tab' },
])

const recordingKey = ref<string | null>(null)

const quality = ref('high')
const qualityOptions = [
  { value: 'low', label: '低', desc: '流畅优先' },
  { value: 'medium', label: '中', desc: '均衡' },
  { value: 'high', label: '高', desc: '画质优先' },
]

const bgmVolume = ref(80)
const sfxVolume = ref(100)
const uiVolume = ref(60)

const gamepadConnected = ref(false)
const gamepadName = ref('')

function startRecord(action: string) {
  recordingKey.value = action
}

function resetKey(action: string) {
  const binding = keyBindings.value.find(k => k.action === action)
  if (binding) {
    binding.binding = binding.defaultKey
    binding.displayKey = binding.defaultKey
  }
  saveKeySettings()
}

function saveSettings() {
  const settings = {
    keyBindings: keyBindings.value.map(k => ({ action: k.action, binding: k.binding })),
    quality: quality.value,
    bgmVolume: bgmVolume.value,
    sfxVolume: sfxVolume.value,
    uiVolume: uiVolume.value,
  }
  localStorage.setItem('game_settings', JSON.stringify(settings))
}

function saveKeySettings() {
  const keyMap: Record<string, string> = {}
  keyBindings.value.forEach(k => { keyMap[k.action] = k.binding })
  localStorage.setItem('game_keybindings', JSON.stringify(keyMap))
}

function detectGamepad() {
  const gamepads = navigator.getGamepads?.()
  if (gamepads) {
    for (const gp of gamepads) {
      if (gp) {
        gamepadConnected.value = true
        gamepadName.value = gp.id
        return
      }
    }
  }
  alert('未检测到手柄，请确保手柄已连接并按下任意按钮激活')
}

// 键盘监听
function onKeyDown(e: KeyboardEvent) {
  if (!recordingKey.value) return
  e.preventDefault()
  const binding = keyBindings.value.find(k => k.action === recordingKey.value)
  if (binding) {
    binding.binding = e.key
    binding.displayKey = e.key === ' ' ? 'Space' : e.key.length === 1 ? e.key.toUpperCase() : e.key
  }
  recordingKey.value = null
  saveKeySettings()
}

// 手柄事件
function onGamepadConnected(e: GamepadEvent) {
  gamepadConnected.value = true
  gamepadName.value = e.gamepad.id
}
function onGamepadDisconnected() {
  gamepadConnected.value = false
  gamepadName.value = ''
}

onMounted(() => {
  // 加载设置
  const saved = localStorage.getItem('game_settings')
  if (saved) {
    try {
      const s = JSON.parse(saved)
      quality.value = s.quality || 'high'
      bgmVolume.value = s.bgmVolume ?? 80
      sfxVolume.value = s.sfxVolume ?? 100
      uiVolume.value = s.uiVolume ?? 60
    } catch {}
  }

  const savedKeys = localStorage.getItem('game_keybindings')
  if (savedKeys) {
    try {
      const map = JSON.parse(savedKeys)
      keyBindings.value.forEach(k => {
        if (map[k.action]) {
          k.binding = map[k.action]
          k.displayKey = map[k.action].length === 1 ? map[k.action].toUpperCase() : map[k.action]
        }
      })
    } catch {}
  }

  document.addEventListener('keydown', onKeyDown)
  window.addEventListener('gamepadconnected', onGamepadConnected as any)
  window.addEventListener('gamepaddisconnected', onGamepadDisconnected)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeyDown)
  window.removeEventListener('gamepadconnected', onGamepadConnected as any)
  window.removeEventListener('gamepaddisconnected', onGamepadDisconnected)
})
</script>

<style scoped>
.settings-page {
  min-height: 100vh;
  background: #0a0a1a;
  color: #e0e0e0;
  padding: 20px;
  max-width: 700px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 30px;
}
.btn-back {
  background: none;
  border: none;
  color: #888;
  font-size: 1rem;
  cursor: pointer;
}
.btn-back:hover { color: #fff; }
.page-header h1 { margin: 0; font-size: 1.5rem; }
.hint { color: #666; font-size: 0.8rem; margin-left: auto; }

.settings-container {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.setting-section {
  background: #111128;
  border-radius: 12px;
  padding: 20px;
}
.setting-section h2 {
  font-size: 1.1rem;
  color: #6c63ff;
  margin: 0 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #222;
}

/* 按键设置 */
.key-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}
.key-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.key-label {
  color: #aaa;
  font-size: 0.85rem;
  width: 60px;
  flex-shrink: 0;
}
.key-btn {
  background: #1a1a38;
  border: 1px solid #333;
  color: #fff;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  min-width: 64px;
  font-size: 0.9rem;
  transition: all 0.15s;
}
.key-btn:hover { border-color: #6c63ff; }
.key-btn.recording {
  border-color: #f59e0b;
  background: #2a2a10;
  animation: pulse 0.8s infinite;
}
@keyframes pulse { 50% { opacity: 0.6; } }

.key-reset {
  background: none;
  border: none;
  color: #f66;
  font-size: 0.75rem;
  cursor: pointer;
  padding: 4px;
}

/* 画质 */
.radio-group {
  display: flex;
  gap: 12px;
}
.radio-label {
  flex: 1;
  background: #1a1a38;
  border: 2px solid #2a2a4a;
  border-radius: 10px;
  padding: 12px;
  cursor: pointer;
  text-align: center;
  transition: all 0.2s;
}
.radio-label input { display: none; }
.radio-label.active { border-color: #6c63ff; background: #1f1f3e; }
.radio-label span { display: block; font-weight: 600; font-size: 1rem; }
.radio-label small { display: block; color: #888; font-size: 0.75rem; margin-top: 4px; }

/* 音量 */
.slider-group {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.slider-item {
  display: flex;
  align-items: center;
  gap: 12px;
}
.slider-item span:first-child { color: #aaa; width: 80px; font-size: 0.9rem; }
.slider-item input[type="range"] {
  flex: 1;
  accent-color: #6c63ff;
  height: 6px;
}
.val { color: #fff; font-size: 0.85rem; width: 40px; text-align: right; }

/* 手柄 */
.gamepad-status {
  display: flex;
  align-items: center;
  gap: 10px;
}
.status-indicator {
  width: 10px; height: 10px;
  border-radius: 50%;
  background: #f44;
}
.status-indicator.connected { background: #22c55e; }
.gamepad-status span { color: #aaa; font-size: 0.9rem; }
.btn-detect {
  background: #333;
  border: 1px solid #555;
  color: #fff;
  padding: 6px 14px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
}
.hint-text { color: #666; font-size: 0.8rem; margin-top: 8px; }

/* 快捷入口 */
.quick-links {
  display: flex;
  gap: 12px;
}
.quick-link {
  background: #1a1a38;
  border: 1px solid #333;
  color: #6c63ff;
  padding: 10px 20px;
  border-radius: 8px;
  text-decoration: none;
  font-size: 0.9rem;
  transition: all 0.2s;
}
.quick-link:hover { border-color: #6c63ff; background: #1f1f3e; }
</style>

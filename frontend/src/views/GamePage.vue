<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { FightGame } from '@/game/FightGame'
import { GameScene } from '@/game/scenes/GameScene'
import { PVEBattleScene } from '@/game/scenes/PVEBattleScene'
import type { CharacterConfig, SkillConfig } from '@/game/BattleEngine'
import type { AIConfig } from '@/game/AIEngine'
import { useWebSocket } from '@/utils/websocket'
import { useGameStore } from '@/store/game'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import BattleResultModal from '@/components/BattleResultModal.vue'

declare global { interface Window { fightGame: FightGame | undefined } }

const router = useRouter()
const route = useRoute()
const gameStore = useGameStore()
const { send: wsSend } = useWebSocket()

const gameCanvas = ref<HTMLDivElement>()
const isReady = ref(false)
const isGameOver = ref(false)
const isPVP = ref(false)
const gameResult = ref<{ winner: string; playerHp: number; enemyHp: number } | null>(null)
const playerName = ref('You')
const enemyName = ref('Enemy')
const isPveWinner = ref(false)
const stageName = ref('')
const showResult = ref(false)

// 战斗 UI 状态
const battleTime = ref(99)
const playerHpPercent = ref(100)
const playerEnergyPercent = ref(100)
const enemyHpPercent = ref(100)
const enemyEnergyPercent = ref(100)
const playerMaxHp = ref(1200)
const enemyMaxHp = ref(1200)
const playerCurrentHp = ref(1200)
const enemyCurrentHp = ref(1200)
const mode = ref('pvp')

const skills = ref([
  { name: '技能1', key: 'U', cd: 0, maxCd: 3, cost: 30 },
  { name: '技能2', key: 'I', cd: 0, maxCd: 5, cost: 40 },
  { name: '技能3', key: 'O', cd: 0, maxCd: 8, cost: 60 },
  { name: '技能4', key: 'P', cd: 0, maxCd: 12, cost: 80 },
])

const isTraining = computed(() => mode.value === 'training')
const countdownDisplay = computed(() => battleTime.value > 0 ? `${battleTime.value}s` : '对战结束')

let gameScene: GameScene | null = null
let pveScene: PVEBattleScene | null = null
let remoteFrameHandler: ((msg: any) => void) | null = null
let battleTimer: any = null

onMounted(async () => {
  const mode = (route.query.mode as string) || 'pvp'
  isPVP.value = mode === 'pvp'

  if (!gameCanvas.value) return

  // 销毁旧实例
  if (window.fightGame) {
    window.fightGame.destroy(true)
    window.fightGame = undefined
  }

  const game = new FightGame('game-canvas')
  window.fightGame = game

  // 等待场景就绪
  game.events.once('ready', async () => {
    if (isPVP.value) {
      await loadAndStartPVPBattle()
    } else {
      await loadAndStartPVEBattle()
    }
  })
})

// ===================== PVE 战斗 =====================

async function loadAndStartPVEBattle(): Promise<void> {
  try {
    const { default: api } = await import('@/api')
    const stageId = Number(route.query.stageId) || 1

    // 获取关卡数据
    let stageData: any = null
    try {
      const res: any = await api.get(`/pve/stages/${stageId}`)
      stageData = res.data?.data || res.data
    } catch (_) { /* fallback */ }

    const bossConfig = stageData?.boss_config
      ? (typeof stageData.boss_config === 'string' ? JSON.parse(stageData.boss_config) : stageData.boss_config)
      : { name: '火焰守护者', phases: [], enemy_hp: 1500, enemy_energy: 80, enemy_speed: 160, enemy_attack: 80, enemy_defense: 40, enemy_energy_regen: 3 }

    stageName.value = stageData?.name || '火焰试炼'
    playerName.value = localStorage.getItem('nickname') || localStorage.getItem('username') || 'You'
    enemyName.value = bossConfig.name || '火焰守护者'

    // 玩家配置
    const playerConfig: CharacterConfig = { id: 1, name: '烈焰战士', hp: 1200, energy: 100, energy_regen: 3, speed: 180, attack: 120, defense: 60 }
    const playerSkills: SkillConfig[] = [
      { id: 1, character_id: 1, name: '烈焰斩', skill_type: 'active', energy_cost: 30, cool_down: 3, damage: 80, range: 70, priority_level: 2, tags: [], effect_params: {} },
      { id: 2, character_id: 1, name: '炎爆', skill_type: 'active', energy_cost: 50, cool_down: 8, damage: 150, range: 90, priority_level: 1, tags: [], effect_params: {} },
      { id: 3, character_id: 1, name: '战意', skill_type: 'passive', energy_cost: 0, cool_down: 0, damage: 0, range: 0, priority_level: 0, tags: [], effect_params: {} },
    ]

    // 敌人配置
    const enemyConfig: CharacterConfig = {
      id: 99, name: bossConfig.name || '火焰守护者',
      hp: bossConfig.enemy_hp || 1500,
      energy: bossConfig.enemy_energy || 80,
      energy_regen: bossConfig.enemy_energy_regen || 3,
      speed: bossConfig.enemy_speed || 160,
      attack: bossConfig.enemy_attack || 80,
      defense: bossConfig.enemy_defense || 40,
    }
    const enemySkills: SkillConfig[] = [
      { id: 101, character_id: 99, name: '火焰拳', skill_type: 'active', energy_cost: 25, cool_down: 4, damage: 60, range: 70, priority_level: 2, tags: [], effect_params: {} },
      { id: 102, character_id: 99, name: '烈焰爆发', skill_type: 'active', energy_cost: 45, cool_down: 7, damage: 120, range: 80, priority_level: 1, tags: [], effect_params: {} },
    ]

    // AI 配置
    const aiConfig: AIConfig = {
      stages: (bossConfig.phases || []).map((p: any, i: number) => ({
        stage: i + 1,
        hp_threshold: p.hp_pct ?? 1,
        aggression: p.behavior === 'aggressive' ? 0.9 : 0.6,
        skill_ids: [101, 102],
        skill_freq: p.behavior === 'aggressive' ? 0.3 : 0.15,
        move_pattern: (p.behavior === 'aggressive' ? 'chase' : 'random') as any,
        description: p.description || '',
      })),
      reaction_time_ms: 300,
      block_chance: 0.1,
      dash_chance: 0.2,
    }
    if (aiConfig.stages.length === 0) {
      aiConfig.stages.push({ stage: 1, hp_threshold: 1, aggression: 0.6, skill_ids: [101, 102], skill_freq: 0.15, move_pattern: 'random', description: '默认' })
    }

    const game = window.fightGame
    if (!game) return

    // 启动 PVEBattleScene，通过 scene.start(data) 直接传入战斗配置
    game.scene.start('PVEBattleScene', {
      playerConfig, playerSkills, playerRules: [],
      enemyConfig, enemySkills, enemyRules: [],
      aiConfig,
    })

    pveScene = game.scene.getScene('PVEBattleScene') as PVEBattleScene
    if (!pveScene) {
      console.error('PVEBattleScene not found')
      return
    }

    // 监听 PVE 结束
    pveScene.onBattleEnd = (result) => {
      gameResult.value = { winner: result.winner, playerHp: result.playerHp, enemyHp: result.enemyHp }
      isGameOver.value = true
      isPveWinner.value = result.winner === 'player'
      showResult.value = true
      clearInterval(battleTimer)

      // 保存 PVE 进度
      const token = localStorage.getItem('token')
      if (token) {
        import('@/api').then(({ default: api }) => {
          api.post('/pve/progress', {
            stage_id: stageId,
            won: result.winner === 'player',
            hp_remaining: result.playerHp,
          }).catch(() => {})
        })
      }
    }

    // 设置 HP 显示
    playerMaxHp.value = playerConfig.hp
    enemyMaxHp.value = enemyConfig.hp
    playerCurrentHp.value = playerConfig.hp
    enemyCurrentHp.value = enemyConfig.hp

    // 启动倒计时
    battleTime.value = 99
    battleTimer = setInterval(() => {
      if (battleTime.value > 0 && !isGameOver.value) {
        battleTime.value--
      }
    }, 1000)

    window.addEventListener('pve-battle-end', handlePveBattleEnd)
    isReady.value = true
  } catch (err) {
    console.error('PVE init error:', err)
  }
}

function handlePveBattleEnd(e: Event): void {
  window.removeEventListener('pve-battle-end', handlePveBattleEnd)
}

// ===================== PVP 战斗 =====================

async function loadAndStartPVPBattle(): Promise<void> {
  try {
    const { default: api } = await import('@/api')
    const token = localStorage.getItem('token')

    let characters: CharacterConfig[] = [{ id: 1, name: '烈焰战士', hp: 1200, energy: 100, energy_regen: 3, speed: 180, attack: 120, defense: 60 }]
    let player1Char = characters[0]
    let player2Char = { ...characters[0], name: gameStore.opponentCharName || '对手' }

    if (token) {
      try {
        const res: any = await api.get('/characters')
        const list = res.data?.data || res.data || []
        if (list.length > 0) {
          characters = list.map((c: any) => ({ id: c.id, name: c.name, hp: c.hp, energy: c.energy, energy_regen: c.energy_regen, speed: c.speed, attack: c.attack, defense: c.defense }))
          player1Char = characters[0]
          player2Char = { ...characters[0], name: gameStore.opponentCharName || '对手' }
        }
      } catch (_) { /* 使用默认角色 */ }
    }

    const game = window.fightGame
    if (!game) return

    // 激活 GameScene
    game.scene.start('GameScene', { player1Char, player2Char })

    await new Promise(r => setTimeout(r, 200))
    gameScene = game.scene.getScene('GameScene') as GameScene

    if (!gameScene) {
      console.error('GameScene not found')
      return
    }

    // 设置角色
    playerName.value = localStorage.getItem('nickname') || localStorage.getItem('username') || 'You'
    enemyName.value = player2Char.name
    gameScene.setPlayerNames(playerName.value, enemyName.value)
    gameScene.setPlayerHP(player1Char.hp, player1Char.hp, player2Char.hp, player2Char.hp)

    // PVP 模式
    gameScene.setRemoteMode(true)

    // 每帧发送本地快照
    gameScene.setOnFrameSnapshot((snapshot) => {
      wsSend({ type: 'frame_input', data: snapshot })
    })

    // 接收远程输入
    remoteFrameHandler = (msg: any) => {
      const data = msg.data || msg
      if (gameScene && !gameResult.value) {
        gameScene.applyRemoteInput({
          keys: data.keys || {},
          x: data.x, y: data.y,
          hp: data.hp, facing: data.facing,
        })
      }
    }

    const { on } = useWebSocket()
    on('frame_input', remoteFrameHandler)

    // 战斗结束
    gameScene.setOnBattleEnd((result) => {
      gameResult.value = result
      isGameOver.value = true
      showResult.value = true
      clearInterval(battleTimer)
      wsSend({ type: 'battle_over', data: { winner: result.winner, playerHp: result.playerHp, enemyHp: result.enemyHp } })
    })

    // 设置 HP 显示
    playerMaxHp.value = player1Char.hp
    enemyMaxHp.value = player2Char.hp
    playerCurrentHp.value = player1Char.hp
    enemyCurrentHp.value = player2Char.hp

    // 启动战斗计时
    battleTime.value = 99
    battleTimer = setInterval(() => {
      if (battleTime.value > 0 && !isGameOver.value) {
        battleTime.value--
      }
    }, 1000)

    // 倒计时
    setTimeout(() => gameScene?.showCountdown(3), 500)
    setTimeout(() => gameScene?.showCountdown(2), 1500)
    setTimeout(() => gameScene?.showCountdown(1), 2500)
    setTimeout(() => gameScene?.showCountdown(0), 3500)

    isReady.value = true
  } catch (err) {
    console.error('PVP init error:', err)
  }
}

function handleExitBattle(): void {
  if (confirm('确认退出对局?')) {
    cleanup()
    if (isPVP.value) {
      wsSend({ type: 'leave_room' })
    }
    router.push(isPVP.value ? '/pvp-rooms' : isPVP.value ? '/lobby' : '/pve-stages')
  }
}

function castSkill(index: number): void {
  const sk = skills.value[index]
  if (sk.cd > 0) return
  sk.cd = sk.maxCd
  // 通知 Phaser 场景释放技能
  if (isPVP.value) {
    gameScene?.castSkill(index)
  } else {
    pveScene?.castSkill(index)
  }
  // CD 倒计时
  const timer = setInterval(() => {
    if (sk.cd > 0) sk.cd--
    else clearInterval(timer)
  }, 1000)
}

function goToReplay(): void {
  cleanup()
  router.push('/replay')
}

function goToLobby(): void {
  cleanup()
  router.push('/lobby')
}

function cleanup(): void {
  window.removeEventListener('pve-battle-end', handlePveBattleEnd)
  if (battleTimer) clearInterval(battleTimer)
  if (remoteFrameHandler) {
    const { off } = useWebSocket()
    off('frame_input', remoteFrameHandler)
    remoteFrameHandler = null
  }
}

onBeforeUnmount(() => {
  cleanup()
  if (window.fightGame) {
    window.fightGame.destroy(true)
    window.fightGame = undefined
  }
})
</script>

<template>
  <div class="battle-page">
    <div v-if="!isReady" class="loading-overlay">
      <div class="loading-text">加载战斗中...</div>
    </div>

    <!-- 战斗 HUD 层 -->
    <div v-show="isReady" class="battle-hud">
      <!-- 退出按钮 -->
      <button class="btn-exit" @click="handleExitBattle">退出对局</button>

      <!-- PVE 关卡名 -->
      <div v-if="!isPVP && stageName" class="stage-tag">关卡: {{ stageName }}</div>

      <!-- 顶部计时器 -->
      <div class="timer-bar" :class="{ warning: battleTime <= 10 }">
        <span class="timer-text">{{ countdownDisplay }}</span>
      </div>

      <!-- 左玩家血条 -->
      <div class="hp-panel left">
        <div class="avatar-box">
          <ArtPlaceholder label="头像" width="48" height="48" bgColor="#1a1a3e" />
        </div>
        <div class="bars">
          <div class="hp-row">
            <span class="label">{{ playerName }}</span>
            <div class="bar-bg hp"><div class="bar-fill hp-fill" :style="{ width: playerHpPercent + '%' }"></div></div>
            <span class="num">{{ playerCurrentHp }} / {{ playerMaxHp }}</span>
          </div>
          <div class="en-row">
            <span class="label en-label">能量</span>
            <div class="bar-bg en"><div class="bar-fill en-fill" :style="{ width: playerEnergyPercent + '%' }"></div></div>
          </div>
        </div>
      </div>

      <!-- 右敌血条 -->
      <div class="hp-panel right">
        <div class="bars">
          <div class="hp-row right-align">
            <span class="num">{{ enemyCurrentHp }} / {{ enemyMaxHp }}</span>
            <div class="bar-bg hp"><div class="bar-fill hp-fill enemy-fill" :style="{ width: enemyHpPercent + '%' }"></div></div>
            <span class="label">{{ enemyName }}</span>
          </div>
          <div class="en-row right-align">
            <div class="bar-bg en"><div class="bar-fill en-fill" :style="{ width: enemyEnergyPercent + '%' }"></div></div>
            <span class="label en-label">能量</span>
          </div>
        </div>
        <div class="avatar-box">
          <ArtPlaceholder label="头像" width="48" height="48" bgColor="#1a1a3e" />
        </div>
      </div>

      <!-- 底部技能栏 -->
      <div class="skill-bar">
        <div v-for="(sk, i) in skills" :key="i" class="skill-slot" @click="castSkill(i)">
          <ArtPlaceholder :label="sk.name" width="52" height="52" bgColor="#2a2a4a" />
          <span class="skill-key">{{ sk.key }}</span>
          <div v-if="sk.cd > 0" class="cd-overlay">{{ sk.cd }}</div>
        </div>
      </div>
    </div>

    <div ref="gameCanvas" id="game-canvas" class="game-canvas" />

    <!-- 对局结算弹窗 -->
    <BattleResultModal
      :visible="showResult"
      :mode="isPVP ? 'pvp' : isTraining ? 'training' : 'pve'"
      :won="gameResult?.winner === 'player'"
      :playerHp="gameResult?.playerHp ?? 100"
      :enemyHp="gameResult?.enemyHp ?? 100"
      :goldReward="isPveWinner ? Math.floor(Math.random() * 200 + 100) : 0"
      :stars="isPveWinner ? Math.ceil(gameResult?.playerHp ? gameResult.playerHp / 33 : 3) : 1"
      @replay="goToReplay"
      @lobby="goToLobby"
    />
  </div>
</template>

<style scoped>
.battle-page {
  width: 100vw; height: 100vh;
  background: #0f0f1a;
  display: flex; align-items: center; justify-content: center;
  overflow: hidden; position: relative;
}
.game-canvas { width: 1024px; height: 768px; }
.loading-overlay {
  position: absolute; inset: 0;
  display: flex; align-items: center; justify-content: center;
  background: #1a1a2e; z-index: 50;
}
.loading-text { color: #666; font-size: 20px; }

/* ====== Battle HUD ====== */
.battle-hud {
  position: absolute; inset: 0;
  z-index: 10; pointer-events: none;
}
.battle-hud > * { pointer-events: auto; }
.btn-exit {
  position: absolute; top: 10px; left: 10px;
  padding: 6px 14px; background: rgba(220,38,38,.2); border: 1px solid #dc2626;
  color: #dc2626; border-radius: 6px; cursor: pointer; font-size: 12px; z-index: 20;
}
.btn-exit:hover { background: rgba(220,38,38,.4); }
.stage-tag {
  position: absolute; top: 10px; left: 50%; transform: translateX(-50%);
  background: rgba(251,191,36,.2); color: #fbbf24;
  padding: 4px 16px; border-radius: 12px; font-size: 13px;
}

/* Timer */
.timer-bar {
  position: absolute; top: 10px; left: 50%; transform: translateX(-50%);
  padding: 4px 20px; background: rgba(0,0,0,.6); border: 1px solid #ffd700;
  border-radius: 16px; margin-top: 28px;
}
.timer-bar.warning { border-color: #dc2626; animation: pulse 0.5s infinite; }
.timer-text { color: #ffd700; font-size: 14px; font-weight: bold; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: .5; } }

/* HP panels */
.hp-panel {
  position: absolute; top: 10px;
  display: flex; gap: 8px; align-items: center;
}
.hp-panel.left { left: 16px; }
.hp-panel.right { right: 16px; flex-direction: row-reverse; }
.avatar-box { border-radius: 8px; overflow: hidden; border: 2px solid #3a3a5e; }
.bars { display: flex; flex-direction: column; gap: 4px; min-width: 200px; }
.hp-row, .en-row { display: flex; align-items: center; gap: 6px; }
.right-align { flex-direction: row-reverse; }
.hp-row .label, .en-row .label { color: #ccc; font-size: 11px; min-width: 36px; white-space: nowrap; }
.en-label { color: #4fc3f7 !important; }
.hp-row .num { color: #aaa; font-size: 10px; white-space: nowrap; }
.bar-bg { flex: 1; height: 14px; background: #1a1a2e; border-radius: 7px; overflow: hidden; border: 1px solid #3a3a5e; }
.bar-bg.en { height: 8px; }
.bar-fill { height: 100%; border-radius: 7px; transition: width .3s; }
.hp-fill { background: linear-gradient(90deg, #22c55e, #16a34a); }
.enemy-fill { background: linear-gradient(270deg, #dc2626, #b91c1c); }
.en-fill { background: linear-gradient(90deg, #4fc3f7, #0288d1); }

/* Skill bar */
.skill-bar {
  position: absolute; bottom: 16px; left: 50%; transform: translateX(-50%);
  display: flex; gap: 12px;
}
.skill-slot {
  position: relative; cursor: pointer; border: 2px solid #3a3a5e;
  border-radius: 10px; overflow: hidden; transition: all .15s;
}
.skill-slot:hover { border-color: #ffd700; transform: translateY(-2px); }
.skill-key {
  position: absolute; bottom: 2px; right: 4px;
  background: rgba(0,0,0,.7); color: #888; font-size: 10px;
  padding: 1px 5px; border-radius: 4px;
}
.cd-overlay {
  position: absolute; inset: 0;
  background: rgba(0,0,0,.7); color: #fff;
  display: flex; align-items: center; justify-content: center;
  font-size: 20px; font-weight: bold;
}
</style>

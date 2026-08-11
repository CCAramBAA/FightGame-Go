/**
 * PVE 战斗场景 — 集成 AI 引擎与战斗引擎
 * 支持 99s 倒计时、技能冷却 UI、训练模式
 */
import Phaser from 'phaser'
import { AIEngine, AIConfig } from '../AIEngine'
import {
  BattleEngine,
  BattleEntity,
  CharacterConfig,
  SkillConfig,
  SpecialRule,
  SkillTag,
} from '../BattleEngine'

interface SkillSlot {
  key: string
  skill: SkillConfig
  cooldownRemaining: number
}

export class PVEBattleScene extends Phaser.Scene {
  // 战斗实体
  private player!: BattleEntity
  private enemy!: BattleEntity

  // AI 引擎
  private aiEngine!: AIEngine

  // 渲染对象
  private playerSprite!: Phaser.GameObjects.Image
  private enemySprite!: Phaser.GameObjects.Image
  private playerNameText!: Phaser.GameObjects.Text
  private enemyNameText!: Phaser.GameObjects.Text
  private playerHpBar!: Phaser.GameObjects.Rectangle
  private enemyHpBar!: Phaser.GameObjects.Rectangle
  private playerHpText!: Phaser.GameObjects.Text
  private enemyHpText!: Phaser.GameObjects.Text
  private playerEnergyBar!: Phaser.GameObjects.Rectangle
  private enemyEnergyBar!: Phaser.GameObjects.Rectangle
  private timerText!: Phaser.GameObjects.Text
  private infoText!: Phaser.GameObjects.Text
  private shieldText!: Phaser.GameObjects.Text
  private enemyShieldText!: Phaser.GameObjects.Text
  private readyCountdownText!: Phaser.GameObjects.Text
  private resultOverlay!: Phaser.GameObjects.Container

  // 技能栏
  private skillSlots: SkillSlot[] = []
  private skillIconBg: Phaser.GameObjects.Rectangle[] = []
  private skillIcons: Phaser.GameObjects.Text[] = []
  private skillCooldownTexts: Phaser.GameObjects.Text[] = []

  // 输入
  private keys!: {
    W: Phaser.Input.Keyboard.Key; A: Phaser.Input.Keyboard.Key
    S: Phaser.Input.Keyboard.Key; D: Phaser.Input.Keyboard.Key
    J: Phaser.Input.Keyboard.Key; U: Phaser.Input.Keyboard.Key
    I: Phaser.Input.Keyboard.Key; O: Phaser.Input.Keyboard.Key
    L: Phaser.Input.Keyboard.Key; TAB: Phaser.Input.Keyboard.Key
  }

  // 状态
  private battleTimeLeft: number = 99
  private totalBattleTime: number = 99
  private isFighting: boolean = false
  private isPaused: boolean = false
  private isOver: boolean = false
  private frameCount: number = 0
  private energyRegenAccum: number = 0

  // 训练模式
  private isTrainingMode: boolean = false

  // 回调
  public onBattleEnd?: (result: { winner: 'player' | 'enemy' | 'draw'; playerHp: number; enemyHp: number }) => void

  constructor() {
    super({ key: 'PVEBattleScene' })
  }

  // ===================== 外部初始化接口 =====================

  initBattle(data: {
    playerConfig: CharacterConfig
    playerSkills: SkillConfig[]
    playerRules: SpecialRule[]
    enemyConfig: CharacterConfig
    enemySkills: SkillConfig[]
    enemyRules: SpecialRule[]
    aiConfig: AIConfig
    isTraining?: boolean
  }): void {
    this.isTrainingMode = data.isTraining || false

    const W = 1024, H = 768
    this.player = BattleEngine.createEntity(
      data.playerConfig, data.playerSkills, data.playerRules,
      250, H - 150, 1,
    )
    this.enemy = BattleEngine.createEntity(
      data.enemyConfig, data.enemySkills, data.enemyRules,
      W - 250, H - 150, -1,
    )

    this.aiEngine = new AIEngine(data.aiConfig)
    this.battleTimeLeft = this.totalBattleTime
    this.isFighting = false
    this.isOver = false
    this.frameCount = 0
  }

  // ===================== Phaser 生命周期 =====================

  create(data?: any): void {
    // 如果有传入数据，先初始化战斗实体
    if (data) {
      this.initBattle(data)
    }

    const W = 1024, H = 768

    // 背景
    this.add.rectangle(W / 2, H / 2, W, H, 0x0a0a1a)
    this.add.rectangle(W / 2, H - 20, W, 40, 0x111128)
    this.add.rectangle(W / 2, H - 100, W - 200, 2, 0x333356)
    this.add.rectangle(W / 2, 0, W, 50, 0x111128).setDepth(10)

    // 玩家精灵
    this.playerSprite = this.add.image(this.player.x, this.player.y, 'flame_warrior_sprite')
    this.playerSprite.setScale(0.9)
    this.playerSprite.setFlipX(true)
    this.enemySprite = this.add.image(this.enemy.x, this.enemy.y, 'flame_warrior_sprite')
    this.enemySprite.setScale(0.9)
    this.enemySprite.setTint(0x3344cc)

    // 名字
    this.playerNameText = this.add.text(this.player.x, this.player.y - 60, 'YOU', {
      fontSize: '14px', color: '#e94560', fontStyle: 'bold',
    }).setOrigin(0.5)
    this.enemyNameText = this.add.text(this.enemy.x, this.enemy.y - 60, 'ENEMY', {
      fontSize: '14px', color: '#3b82f6', fontStyle: 'bold',
    }).setOrigin(0.5)

    // 顶部 UI
    // 玩家 HP
    this.add.text(20, 16, 'YOU', { fontSize: '13px', color: '#e94560', fontStyle: 'bold' }).setOrigin(0, 0.5).setDepth(11)
    let hpBg = this.add.rectangle(56, 16, 200, 14, 0x333355).setOrigin(0, 0.5).setDepth(11)
    this.playerHpBar = this.add.rectangle(56, 16, 200, 14, 0xe94560).setOrigin(0, 0.5).setDepth(12)
    this.playerHpText = this.add.text(260, 16, '', { fontSize: '11px', color: '#fff' }).setOrigin(0, 0.5).setDepth(12)
    this.playerEnergyBar = this.add.rectangle(56, 34, 0, 6, 0xfbbf24).setOrigin(0, 0.5).setDepth(12)
    this.add.rectangle(56, 34, 200, 6, 0x333355).setOrigin(0, 0.5).setDepth(11)

    // 护盾文字
    this.shieldText = this.add.text(10, 55, '', { fontSize: '11px', color: '#a78bfa' }).setDepth(12)

    // 敌人 HP
    this.add.text(W - 20, 16, 'ENEMY', { fontSize: '13px', color: '#3b82f6', fontStyle: 'bold' }).setOrigin(1, 0.5).setDepth(11)
    let eHpBg = this.add.rectangle(W - 56, 16, 200, 14, 0x333355).setOrigin(1, 0.5).setDepth(11)
    this.enemyHpBar = this.add.rectangle(W - 56, 16, 200, 14, 0x3b82f6).setOrigin(1, 0.5).setDepth(12)
    this.enemyHpText = this.add.text(W - 260, 16, '', { fontSize: '11px', color: '#fff' }).setOrigin(0, 0.5).setDepth(12)
    this.enemyEnergyBar = this.add.rectangle(W - 56, 34, 0, 6, 0xfbbf24).setOrigin(1, 0.5).setDepth(12)
    this.add.rectangle(W - 56, 34, 200, 6, 0x333355).setOrigin(1, 0.5).setDepth(11)
    this.enemyShieldText = this.add.text(W - 10, 55, '', { fontSize: '11px', color: '#a78bfa' }).setOrigin(1, 0).setDepth(12)

    // 倒计时
    this.timerText = this.add.text(W / 2, 16, `${this.totalBattleTime}`.padStart(2, '0'), {
      fontSize: '22px', color: '#fbbf24', fontStyle: 'bold',
    }).setOrigin(0.5, 0.5).setDepth(12)

    // 底部技能栏
    this.createSkillBar()

    // 底部提示
    this.infoText = this.add.text(W / 2, H - 10, 'WASD 移动 | J 普攻 | U/I/O 技能 | L 闪避 | Tab 技能详情', {
      fontSize: '11px', color: '#555',
    }).setOrigin(0.5)

    // 3 秒倒计时
    this.readyCountdownText = this.add.text(W / 2, H / 2, '', {
      fontSize: '64px', color: '#fff', stroke: '#000', strokeThickness: 6,
    }).setOrigin(0.5).setDepth(100)

    // 结果遮罩层（初始隐藏）
    this.resultOverlay = this.add.container(0, 0).setDepth(200).setAlpha(0)

    // 键盘
    this.keys = {
      W: this.input.keyboard!.addKey('W'), A: this.input.keyboard!.addKey('A'),
      S: this.input.keyboard!.addKey('S'), D: this.input.keyboard!.addKey('D'),
      J: this.input.keyboard!.addKey('J'), U: this.input.keyboard!.addKey('U'),
      I: this.input.keyboard!.addKey('I'), O: this.input.keyboard!.addKey('O'),
      L: this.input.keyboard!.addKey('L'), TAB: this.input.keyboard!.addKey('TAB'),
    }

    // 启动 3 秒准备倒计时
    this.startReadyCountdown()
  }

  private createSkillBar(): void {
    const W = 1024, H = 768
    const activeSkills = this.player.skills.filter((s: SkillConfig) => s.skill_type === 'active')
    this.skillSlots = activeSkills.slice(0, 4).map((s: SkillConfig) => ({ key: String(s.id), skill: s, cooldownRemaining: 0 }))

    const startX = W / 2 - this.skillSlots.length * 35
    this.skillSlots.forEach((slot, i) => {
      const x = startX + i * 70
      const y = H - 45

      const bg = this.add.rectangle(x, y, 50, 50, 0x1a1a38)
      bg.setStrokeStyle(2, 0x444466).setDepth(20)
      this.skillIconBg.push(bg)

      const shortcut = ['U', 'I', 'O', 'P'][i]
      const text = this.add.text(x, y, shortcut, {
        fontSize: '18px', color: '#fff', fontStyle: 'bold',
      }).setOrigin(0.5).setDepth(21)
      this.skillIcons.push(text)

      // 冷却文字
      const cd = this.add.text(x, y + 30, '', {
        fontSize: '10px', color: '#f66',
      }).setOrigin(0.5).setDepth(21)
      this.skillCooldownTexts.push(cd)

      // 能量消耗
      this.add.text(x + 18, y - 18, `${slot.skill.energy_cost}`, {
        fontSize: '9px', color: '#fbbf24',
      }).setDepth(21)
    })
  }

  // ===================== 准备倒计时 =====================
  private readyCount: number = 3
  private startReadyCountdown(): void {
    this.readyCount = 3
    const doCount = () => {
      if (this.readyCount > 0) {
        this.readyCountdownText.setText(String(this.readyCount))
        this.tweens.add({
          targets: this.readyCountdownText,
          scaleX: 1.5, scaleY: 1.5, alpha: 1,
          duration: 300,
          yoyo: false,
          onComplete: () => {
            this.readyCountdownText.setScale(1).setAlpha(1)
          },
        })
        this.readyCount--
        this.time.delayedCall(800, doCount)
      } else {
        this.readyCountdownText.setText('FIGHT!').setColor('#fbbf24')
        this.tweens.add({
          targets: this.readyCountdownText,
          alpha: 0, duration: 500,
          onComplete: () => {
            this.readyCountdownText.setVisible(false)
            this.isFighting = true
          },
        })
      }
    }
    doCount()
  }

  // ===================== 每帧更新 =====================
  update(_time: number, delta: number): void {
    if (!this.isFighting || this.isOver) return

    this.frameCount++
    this.handlePlayerInput(delta)
    this.handleAI(delta)
    this.updateBuffs(delta)
    this.updateTimers(delta)
    this.updateSpritePositions()
    this.updateUI()
    this.checkBattleEnd()
  }

  // ===================== 玩家输入 =====================
  private handlePlayerInput(delta: number): void {
    if (BattleEngine.isControlled(this.player)) return

    const p = this.player
    const dt = delta / 1000
    const speed = BattleEngine.getEffectiveSpeed(p)

    // 移动
    let moveX = 0
    if (this.keys.A.isDown) { moveX = -speed * dt; p.facing = -1 }
    if (this.keys.D.isDown) { moveX = speed * dt; p.facing = 1 }
    p.x = Phaser.Math.Clamp(p.x + moveX, 50, 974)

    // 普攻
    if (this.keys.J.isDown && p.attackCooldown <= 0) {
      p.attackCooldown = 400
      this.executeMeleeAttack(p, this.enemy)
    }

    // 技能键绑定到 U/I/O
    const skillKeys = [this.keys.U, this.keys.I, this.keys.O]
    skillKeys.forEach((key, i) => {
      if (key.isDown && this.skillSlots[i] && this.skillSlots[i].cooldownRemaining <= 0) {
        const slot = this.skillSlots[i]
        if (p.energy >= slot.skill.energy_cost) {
          p.energy -= slot.skill.energy_cost
          p.attackCooldown = 300
          this.skillSlots[i].cooldownRemaining = slot.skill.cool_down * 1000
          this.executeSkill(p, this.enemy, slot.skill)
        }
      }
    })

    // 闪避
    if (this.keys.L.isDown && p.attackCooldown <= 0) {
      p.attackCooldown = 300
      p.x += p.facing * 100
      p.x = Phaser.Math.Clamp(p.x, 50, 974)
    }

    // 技能面板
    if (Phaser.Input.Keyboard.JustDown(this.keys.TAB)) {
      this.toggleSkillPanel()
    }
  }

  // ===================== AI 控制 =====================
  private handleAI(delta: number): void {
    if (BattleEngine.isControlled(this.enemy)) return

    const p = this.player
    const e = this.enemy

    const decision = this.aiEngine.tick(
      { x: e.x, y: e.y, hp: e.hp, maxHp: e.maxHp, energy: e.energy, maxEnergy: e.maxEnergy, speed: e.config.speed, facing: e.facing, attackCooldown: e.attackCooldown, isHit: false, hitTimer: 0, skillCooldowns: e.skillCooldowns },
      { x: p.x, y: p.y, hp: p.hp, maxHp: p.maxHp, energy: p.energy, maxEnergy: p.maxEnergy, speed: p.config.speed, facing: p.facing, attackCooldown: p.attackCooldown, isHit: false, hitTimer: 0, skillCooldowns: p.skillCooldowns },
      delta,
    )

    e.x = Phaser.Math.Clamp(e.x + decision.moveX, 50, 974)

    // AI 普攻
    if (decision.attack && e.attackCooldown <= 0) {
      e.attackCooldown = 400
      this.executeMeleeAttack(e, p)
    }

    // AI 技能
    if (decision.skillIndex >= 0 && e.skills[decision.skillIndex]) {
      const skill = e.skills[decision.skillIndex]
      if (e.energy >= skill.energy_cost && (e.skillCooldowns[skill.id] || 0) <= 0) {
        e.energy -= skill.energy_cost
        e.attackCooldown = 300
        e.skillCooldowns[skill.id] = skill.cool_down * 1000
        this.executeSkill(e, p, skill)
      }
    }
  }

  // ===================== 普攻 =====================
  private executeMeleeAttack(attacker: BattleEntity, defender: BattleEntity): void {
    const dist = Math.abs(attacker.x - defender.x)
    if (dist > 100) return

    // 普攻无技能标签，固定伤害
    const rawDamage = attacker.config.attack * 0.3
    const defenseReduction = defender.config.defense / (defender.config.defense + 100)
    const dmg = Math.max(1, rawDamage * (1 - defenseReduction))

    defender.hp = Math.max(0, defender.hp - dmg)

    // 攻击动画
    const sprite = attacker === this.player ? this.playerSprite : this.enemySprite
    this.tweens.add({
      targets: sprite, scaleX: 1.3, scaleY: 0.7, duration: 80, yoyo: true,
    })

    // 受击闪烁
    const defSprite = defender === this.player ? this.playerSprite : this.enemySprite
    defSprite.setTint(0xffffff)
    this.time.delayedCall(150, () => {
      if (defender === this.player) {
        defSprite.clearTint()
      } else {
        defSprite.setTint(0x3344cc)
      }
    })

    // 能量恢复
    attacker.energy = Math.min(attacker.maxEnergy, attacker.energy + 5)
    defender.energy = Math.min(defender.maxEnergy, defender.energy + 8)
  }

  // ===================== 技能释放 =====================
  private executeSkill(attacker: BattleEntity, defender: BattleEntity, skill: SkillConfig): void {
    const dist = Math.abs(attacker.x - defender.x)
    if (dist > skill.range) return

    // 优先级判定（如果对方正在攻击）
    const result = BattleEngine.calcDamage(attacker, defender, skill)

    // 动画
    const sprite = attacker === this.player ? this.playerSprite : this.enemySprite
    this.tweens.add({
      targets: sprite, scaleX: 1.4, scaleY: 0.6, duration: 100, yoyo: true,
    })

    const defSprite = defender === this.player ? this.playerSprite : this.enemySprite
    defSprite.setTint(0xffffff)
    this.time.delayedCall(150, () => {
      if (defender === this.player) {
        defSprite.clearTint()
      } else {
        defSprite.setTint(0x3344cc)
      }
    })

    // 受击能量恢复
    defender.energy = Math.min(defender.maxEnergy, defender.energy + 10)
  }

  // ===================== Buff 更新 =====================
  private updateBuffs(delta: number): void {
    BattleEngine.updateBuffs(this.player, delta)
    BattleEngine.updateBuffs(this.enemy, delta)

    // 更新技能冷却
    const updateCooldowns = (entity: BattleEntity) => {
      entity.attackCooldown = Math.max(0, entity.attackCooldown - delta)
      for (const key in entity.skillCooldowns) {
        entity.skillCooldowns[key] = Math.max(0, entity.skillCooldowns[key] - delta)
      }
    }
    updateCooldowns(this.player)
    updateCooldowns(this.enemy)

    // 更新技能栏冷却
    this.skillSlots.forEach((s, i) => {
      if (s.cooldownRemaining > 0) {
        s.cooldownRemaining = Math.max(0, s.cooldownRemaining - delta)
      }
    })
  }

  // ===================== 计时器 =====================
  private updateTimers(delta: number): void {
    // 战斗时间
    this.battleTimeLeft = Math.max(0, this.battleTimeLeft - delta / 1000)

    // 能量恢复
    this.energyRegenAccum += delta
    if (this.energyRegenAccum >= 1000) {
      this.energyRegenAccum -= 1000
      this.player.energy = Math.min(this.player.maxEnergy, this.player.energy + this.player.config.energy_regen)
      this.enemy.energy = Math.min(this.enemy.maxEnergy, this.enemy.energy + this.enemy.config.energy_regen)
    }
  }

  // ===================== 战斗结束判定 =====================
  private checkBattleEnd(): void {
    if (this.isOver) return

    let winner: 'player' | 'enemy' | 'draw' | null = null

    if (this.player.hp <= 0) winner = 'enemy'
    if (this.enemy.hp <= 0) winner = 'player'
    if (this.battleTimeLeft <= 0) {
      if (this.player.hp > this.enemy.hp) winner = 'player'
      else if (this.enemy.hp > this.player.hp) winner = 'enemy'
      else winner = 'draw'
    }

    if (!winner) return

    this.isOver = true
    this.showResult(winner)

    if (this.onBattleEnd) {
      this.onBattleEnd({
        winner,
        playerHp: this.player.hp,
        enemyHp: this.enemy.hp,
      })
    }
  }

  private showResult(winner: 'player' | 'enemy' | 'draw'): void {
    const W = 1024, H = 768
    this.resultOverlay.setAlpha(1)

    const bg = this.add.rectangle(W / 2, H / 2, W, H, 0x000000, 0.7).setDepth(199)
    this.resultOverlay.add(bg)

    let text = winner === 'draw' ? '平局!' : winner === 'player' ? '胜利!' : '失败'
    let color = winner === 'draw' ? '#fbbf24' : winner === 'player' ? '#22c55e' : '#ef4444'

    const title = this.add.text(W / 2, H / 2 - 60, text, {
      fontSize: '56px', color, fontStyle: 'bold',
    }).setOrigin(0.5)
    this.resultOverlay.add(title)

    const details = this.add.text(W / 2, H / 2 + 10, `剩余血量: ${Math.ceil(this.player.hp)} | 敌人血量: ${Math.ceil(this.enemy.hp)}`, {
      fontSize: '16px', color: '#aaa',
    }).setOrigin(0.5)
    this.resultOverlay.add(details)

    // 按钮
    const btnReturn = this.add.text(W / 2, H / 2 + 70, '[ 返回大厅 ]', {
      fontSize: '20px', color: '#6c63ff', fontStyle: 'bold',
      padding: { x: 20, y: 10 },
      backgroundColor: '#1f1f3e',
    }).setOrigin(0.5).setInteractive({ useHandCursor: true })
    btnReturn.on('pointerdown', () => {
      this.scene.stop()
      // 通过事件通知 Vue 层
      window.dispatchEvent(new CustomEvent('pve-battle-end', {
        detail: { winner, playerHp: this.player.hp, enemyHp: this.enemy.hp },
      }))
    })
    this.resultOverlay.add(btnReturn)
  }

  // ===================== 技能面板 =====================
  private toggleSkillPanel(): void {
    // TODO: 侧边面板展示技能描述
  }

  // ===================== 精灵位置与 UI 更新 =====================
  private updateSpritePositions(): void {
    this.playerSprite.setPosition(this.player.x, this.player.y)
    this.enemySprite.setPosition(this.enemy.x, this.enemy.y)
    this.playerNameText.setPosition(this.player.x, this.player.y - 60)
    this.enemyNameText.setPosition(this.enemy.x, this.enemy.y - 60)
  }

  private updateUI(): void {
    // HP
    const pHpR = this.player.hp / this.player.maxHp
    this.playerHpBar.setSize(200 * pHpR, 14)
    this.playerHpText.setText(`${Math.ceil(this.player.hp)}/${this.player.maxHp}`)

    const eHpR = this.enemy.hp / this.enemy.maxHp
    this.enemyHpBar.setSize(200 * eHpR, 14)
    this.enemyHpText.setText(`${Math.ceil(this.enemy.hp)}/${this.enemy.maxHp}`)

    // 能量
    this.playerEnergyBar.setSize(200 * (this.player.energy / this.player.maxEnergy), 6)
    this.enemyEnergyBar.setSize(200 * (this.enemy.energy / this.enemy.maxEnergy), 6)

    // 护盾
    this.shieldText.setText(this.player.shield > 0 ? `护盾 ${Math.ceil(this.player.shield)}` : '')
    this.enemyShieldText.setText(this.enemy.shield > 0 ? `护盾 ${Math.ceil(this.enemy.shield)}` : '')

    // 倒计时
    const secs = Math.ceil(this.battleTimeLeft)
    this.timerText.setText(secs.toString().padStart(2, '0'))
    if (secs <= 10) this.timerText.setColor('#ef4444')
    else if (secs <= 30) this.timerText.setColor('#fbbf24')

    // 技能栏冷却
    this.skillCooldownTexts.forEach((t, i) => {
      const remaining = this.skillSlots[i]?.cooldownRemaining || 0
      if (remaining > 0) {
        t.setText(`${(remaining / 1000).toFixed(1)}s`)
        this.skillIconBg[i].setFillStyle(0x2a1a1a)
      } else {
        t.setText('')
        this.skillIconBg[i].setFillStyle(0x1a1a38)
      }
    })
  }

  // ===================== 外部控制 =====================
  setTrainingMode(enabled: boolean): void {
    this.isTrainingMode = enabled
    if (this.aiEngine) {
      this.aiEngine.infiniteHP = enabled
      this.aiEngine.infiniteEnergy = enabled
    }
  }

  toggleAIPause(): void {
    if (this.aiEngine) {
      this.aiEngine.paused = !this.aiEngine.paused
    }
  }
}

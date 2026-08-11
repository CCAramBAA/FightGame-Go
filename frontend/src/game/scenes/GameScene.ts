import Phaser from 'phaser'

interface CharacterState {
  x: number; y: number
  hp: number; maxHp: number
  energy: number; maxEnergy: number
  speed: number; name: string
  velocityX: number; velocityY: number
  isAttacking: boolean; isDefending: boolean; isHit: boolean
  hitTimer: number; attackCooldown: number
  facing: number // 1=right, -1=left
}

export interface RemoteInput {
  keys: { left: boolean; right: boolean; up: boolean; down: boolean; attack: boolean; defense: boolean; skill1: boolean; skill2: boolean; skill3: boolean }
  x?: number; y?: number; hp?: number; facing?: number
}

export interface FrameSnapshot {
  frame: number
  keys: { left: boolean; right: boolean; up: boolean; down: boolean; attack: boolean; defense: boolean; skill1: boolean; skill2: boolean; skill3: boolean }
  x: number; y: number; hp: number; facing: number
}

export class GameScene extends Phaser.Scene {
  private player1!: CharacterState
  private player2!: CharacterState

  private p1Body!: Phaser.GameObjects.Image
  private p2Body!: Phaser.GameObjects.Image
  private p1NameText!: Phaser.GameObjects.Text
  private p2NameText!: Phaser.GameObjects.Text

  private p1HpBar!: Phaser.GameObjects.Rectangle
  private p1HpBg!: Phaser.GameObjects.Rectangle
  private p2HpBar!: Phaser.GameObjects.Rectangle
  private p2HpBg!: Phaser.GameObjects.Rectangle
  private p1EnergyBar!: Phaser.GameObjects.Rectangle
  private p2EnergyBar!: Phaser.GameObjects.Rectangle
  private p1HpText!: Phaser.GameObjects.Text
  private p2HpText!: Phaser.GameObjects.Text
  private countdownText!: Phaser.GameObjects.Text
  private winnerText!: Phaser.GameObjects.Text
  private infoText!: Phaser.GameObjects.Text

  private keys!: {
    W: Phaser.Input.Keyboard.Key; A: Phaser.Input.Keyboard.Key
    S: Phaser.Input.Keyboard.Key; D: Phaser.Input.Keyboard.Key
    J: Phaser.Input.Keyboard.Key; K: Phaser.Input.Keyboard.Key
    L: Phaser.Input.Keyboard.Key; ENTER: Phaser.Input.Keyboard.Key
  }

  private onFrameSnapshot?: (snapshot: FrameSnapshot) => void
  private onAttack?: (skillIndex: number) => void
  private onBattleEnd?: (result: { winner: string; playerHp: number; enemyHp: number }) => void

  private energyTimer = 0
  private readonly ENERGY_REGEN_RATE = 2
  private readonly ENERGY_REGEN_INTERVAL = 60

  // 远程控制
  private remoteMode = false
  private remoteInputBuffer: RemoteInput[] = []
  private frameCounter = 0
  private battleOver = false

  constructor() { super({ key: 'GameScene' }) }

  create(): void {
    const W = 1024; const H = 768
    this.battleOver = false

    this.add.rectangle(W / 2, H / 2, W, H, 0x1a1a2e)
    this.add.rectangle(W / 2, H - 20, W, 40, 0x16213e)
    this.add.rectangle(W / 2, 20, W, 40, 0x16213e)
    this.add.rectangle(W / 2, H - 100, W, 4, 0x334155)

    this.player1 = makeCharState(250, H - 150, 100, 100, 200, 'Player 1', 1)
    this.player2 = makeCharState(W - 250, H - 150, 100, 100, 200, 'Player 2', -1)

    this.p1Body = this.add.image(this.player1.x, this.player1.y, 'flame_warrior_sprite').setScale(0.9).setFlipX(true)
    this.p2Body = this.add.image(this.player2.x, this.player2.y, 'flame_warrior_sprite').setScale(0.9).setTint(0x3344cc)

    this.p1NameText = this.add.text(this.player1.x, this.player1.y - 60, 'YOU', { fontSize: '14px', color: '#e94560' }).setOrigin(0.5)
    this.p2NameText = this.add.text(this.player2.x, this.player2.y - 60, 'ENEMY', { fontSize: '14px', color: '#3b82f6' }).setOrigin(0.5)

    this.add.text(20, 8, 'YOU', { fontSize: '14px', color: '#e94560' }).setOrigin(0, 0.5)
    this.p1HpBg = this.add.rectangle(180, 14, 200, 14, 0x333333)
    this.p1HpBar = this.add.rectangle(180, 14, 200, 14, 0xe94560).setOrigin(0, 0.5)
    this.p1HpText = this.add.text(290, 14, '100/100', { fontSize: '11px', color: '#fff' }).setOrigin(1, 0.5)
    this.p1EnergyBar = this.add.rectangle(180, 32, 0, 8, 0xfbbf24).setOrigin(0, 0.5)
    this.add.rectangle(180, 32, 200, 8, 0x333333).setOrigin(0, 0.5).setDepth(0)

    this.add.text(W - 20, 8, 'ENEMY', { fontSize: '14px', color: '#3b82f6' }).setOrigin(1, 0.5)
    this.p2HpBg = this.add.rectangle(W - 180, 14, 200, 14, 0x333333)
    this.p2HpBar = this.add.rectangle(W - 180, 14, 200, 14, 0x3b82f6).setOrigin(1, 0.5)
    this.p2HpText = this.add.text(W - 290, 14, '100/100', { fontSize: '11px', color: '#fff' }).setOrigin(0, 0.5)
    this.p2EnergyBar = this.add.rectangle(W - 180, 32, 0, 8, 0xfbbf24).setOrigin(1, 0.5)
    this.add.rectangle(W - 180, 32, 200, 8, 0x333333).setOrigin(1, 0.5).setDepth(0)

    this.countdownText = this.add.text(W / 2, H / 2 - 50, '', {
      fontSize: '72px', color: '#fff', stroke: '#000', strokeThickness: 6,
    }).setOrigin(0.5).setDepth(100).setAlpha(0)

    this.winnerText = this.add.text(W / 2, H / 2, '', {
      fontSize: '48px', color: '#fbbf24', stroke: '#000', strokeThickness: 4,
    }).setOrigin(0.5).setDepth(100).setAlpha(0)

    this.infoText = this.add.text(W / 2, H - 30, 'WASD 移动 | J 普攻 | K 技能1 | L 技能2', {
      fontSize: '12px', color: '#666',
    }).setOrigin(0.5)

    this.keys = {
      W: this.input.keyboard!.addKey('W'), A: this.input.keyboard!.addKey('A'),
      S: this.input.keyboard!.addKey('S'), D: this.input.keyboard!.addKey('D'),
      J: this.input.keyboard!.addKey('J'), K: this.input.keyboard!.addKey('K'),
      L: this.input.keyboard!.addKey('L'), ENTER: this.input.keyboard!.addKey('ENTER'),
    }
  }

  // ========= 外部 API =========
  setRemoteMode(enabled: boolean): void { this.remoteMode = enabled }
  setRemotePlayerName(name: string): void {
    this.player2.name = name
    this.p2NameText.setText(name)
  }

  setPlayerNames(p1Name: string, p2Name: string): void {
    this.player1.name = p1Name
    this.player2.name = p2Name
    this.p1NameText.setText(p1Name)
    this.p2NameText.setText(p2Name)
  }

  setPlayerHP(p1HP: number, p1MaxHP: number, p2HP: number, p2MaxHP: number): void {
    this.player1.hp = Math.max(0, p1HP); this.player1.maxHp = p1MaxHP
    this.player2.hp = Math.max(0, p2HP); this.player2.maxHp = p2MaxHP
  }

  showCountdown(count: number): void {
    this.countdownText.setText(String(count)).setAlpha(1)
    this.tweens.add({ targets: this.countdownText, alpha: 0, scaleX: 1.5, scaleY: 1.5, duration: 800 })
  }

  showWinner(winnerName: string): void {
    this.winnerText.setText(`${winnerName} 获胜!`).setAlpha(1)
  }

  resetBattle(): void {
    this.player1.hp = this.player1.maxHp; this.player2.hp = this.player2.maxHp
    this.player1.energy = 0; this.player2.energy = 0
    this.player1.x = 250; this.player2.x = 1024 - 250
    this.winnerText.setAlpha(0); this.countdownText.setAlpha(0)
    this.battleOver = false
  }

  setOnFrameSnapshot(cb: (snapshot: FrameSnapshot) => void): void { this.onFrameSnapshot = cb }
  setOnAttack(cb: (skillIndex: number) => void): void { this.onAttack = cb }
  setOnBattleEnd(cb: (result: { winner: string; playerHp: number; enemyHp: number }) => void): void { this.onBattleEnd = cb }

  /** 应用远程玩家的输入 */
  applyRemoteInput(input: RemoteInput): void {
    this.remoteInputBuffer.push(input)
    if (this.remoteInputBuffer.length > 5) this.remoteInputBuffer.shift()
  }

  // ========= 主循环 =========
  update(_time: number, delta: number): void {
    if (this.battleOver) return
    this.handleLocalInput(delta)
    this.handleRemoteInput(delta)
    this.updatePhysics(delta)
    this.updateUI()
    this.sendFrameSnapshot()
    this.checkBattleEnd()
  }

  private sendFrameSnapshot(): void {
    this.frameCounter++
    // 每 3 帧发送一次快照
    if (this.frameCounter % 3 !== 0) return
    if (!this.remoteMode) return

    const p1 = this.player1
    const snapshot: FrameSnapshot = {
      frame: this.frameCounter,
      keys: getKeyState(this.keys),
      x: Math.round(p1.x), y: Math.round(p1.y),
      hp: Math.round(p1.hp),
      facing: p1.facing,
    }
    this.onFrameSnapshot?.(snapshot)
  }

  private handleLocalInput(delta: number): void {
    const p1 = this.player1
    const dt = delta / 1000

    let moveX = 0
    if (this.keys.A.isDown) { moveX = -p1.speed * dt; p1.facing = -1 }
    if (this.keys.D.isDown) { moveX = p1.speed * dt; p1.facing = 1 }
    p1.x = Phaser.Math.Clamp(p1.x + moveX, 60, 1024 - 60)
    p1.isDefending = this.keys.S.isDown

    if (p1.attackCooldown > 0) p1.attackCooldown -= delta
    if (p1.isHit) { p1.hitTimer -= delta; if (p1.hitTimer <= 0) p1.isHit = false }

    // 攻击
    if (this.keys.J.isDown && p1.attackCooldown <= 0) {
      this.executeAttack(p1, this.player2, 10)
      this.onAttack?.(0)
    }
    if (this.keys.K.isDown && p1.attackCooldown <= 0 && p1.energy >= 30) {
      p1.energy -= 30
      this.executeAttack(p1, this.player2, 20)
      this.onAttack?.(1)
    }
    if (this.keys.L.isDown && p1.attackCooldown <= 0 && p1.energy >= 50) {
      p1.energy -= 50
      this.executeAttack(p1, this.player2, 35)
      this.onAttack?.(2)
    }
  }

  private handleRemoteInput(delta: number): void {
    if (!this.remoteMode) return
    const input = this.remoteInputBuffer.shift()
    if (!input) return

    const p2 = this.player2
    const dt = delta / 1000

    // 状态校正
    if (input.x !== undefined && input.y !== undefined) {
      p2.x = Phaser.Math.Linear(p2.x, input.x, 0.3)
    }
    if (input.hp !== undefined) {
      p2.hp = Phaser.Math.Linear(p2.hp, input.hp, 0.5)
    }
    if (input.facing !== undefined) p2.facing = input.facing

    // 按键→移动
    const keys = input.keys
    let moveX = 0
    if (keys.left) moveX = -p2.speed * dt
    if (keys.right) moveX = p2.speed * dt
    p2.x = Phaser.Math.Clamp(p2.x + moveX, 60, 1024 - 60)
    p2.isDefending = keys.defense

    if (p2.attackCooldown > 0) p2.attackCooldown -= delta
    if (p2.isHit) { p2.hitTimer -= delta; if (p2.hitTimer <= 0) p2.isHit = false }

    if (keys.attack && p2.attackCooldown <= 0) this.executeAttack(p2, this.player1, 10)
    if (keys.skill1 && p2.attackCooldown <= 0 && p2.energy >= 30) { p2.energy -= 30; this.executeAttack(p2, this.player1, 20) }
    if (keys.skill2 && p2.attackCooldown <= 0 && p2.energy >= 50) { p2.energy -= 50; this.executeAttack(p2, this.player1, 35) }
  }

  private executeAttack(attacker: CharacterState, defender: CharacterState, damage: number): void {
    attacker.isAttacking = true
    attacker.attackCooldown = 400
    const dist = Math.abs(attacker.x - defender.x)
    if (dist < 100) {
      const multiplier = defender.isDefending ? 0.3 : 1
      defender.hp = Math.max(0, defender.hp - damage * multiplier)
      defender.isHit = true
      defender.hitTimer = 200
    }
    const body = attacker === this.player1 ? this.p1Body : this.p2Body
    this.tweens.add({ targets: body, scaleX: 1.2, scaleY: 0.8, duration: 100, yoyo: true })
  }

  private updatePhysics(_delta: number): void {
    // 非远程模式下 P2 使用简单 AI
    if (!this.remoteMode) {
      const p2 = this.player2
      p2.x += p2.velocityX * (1 / 60)
      p2.velocityX += (Math.sin(Date.now() / 1000) * 2 - p2.velocityX) * 0.02
      p2.facing = this.player1.x > p2.x ? 1 : -1
      p2.x = Phaser.Math.Clamp(p2.x, 60, 1024 - 60)
    }

    this.energyTimer++
    if (this.energyTimer >= this.ENERGY_REGEN_INTERVAL) {
      this.energyTimer = 0
      this.player1.energy = Math.min(this.player1.maxEnergy, this.player1.energy + this.ENERGY_REGEN_RATE)
      this.player2.energy = Math.min(this.player2.maxEnergy, this.player2.energy + this.ENERGY_REGEN_RATE)
    }
  }

  private checkBattleEnd(): void {
    if (this.player1.hp <= 0) {
      this.battleOver = true
      this.showWinner(this.player2.name)
      this.onBattleEnd?.({ winner: 'enemy', playerHp: this.player1.hp, enemyHp: this.player2.hp })
    } else if (this.player2.hp <= 0) {
      this.battleOver = true
      this.showWinner(this.player1.name)
      this.onBattleEnd?.({ winner: 'player', playerHp: this.player1.hp, enemyHp: this.player2.hp })
    }
  }

  private updateUI(): void {
    this.p1Body.setPosition(this.player1.x, this.player1.y)
    this.p1NameText.setPosition(this.player1.x, this.player1.y - 60)
    this.p2Body.setPosition(this.player2.x, this.player2.y)
    this.p2NameText.setPosition(this.player2.x, this.player2.y - 60)

    const p1R = Math.max(0, this.player1.hp / this.player1.maxHp)
    this.p1HpBar.setSize(200 * p1R, 14)
    this.p1HpText.setText(`${Math.ceil(this.player1.hp)}/${this.player1.maxHp}`)

    const p2R = Math.max(0, this.player2.hp / this.player2.maxHp)
    this.p2HpBar.setSize(200 * p2R, 14)
    this.p2HpText.setText(`${Math.ceil(this.player2.hp)}/${this.player2.maxHp}`)

    this.p1EnergyBar.setSize(200 * (this.player1.energy / this.player1.maxEnergy), 8)
    this.p2EnergyBar.setSize(200 * (this.player2.energy / this.player2.maxEnergy), 8)

    if (this.player1.isHit) this.p1Body.setTint(0xffffff)
    else this.p1Body.clearTint()
    if (this.player2.isHit) this.p2Body.setTint(0xffffff)
    else if (this.remoteMode) this.p2Body.setTint(0x3344cc)
  }
}

function makeCharState(x: number, y: number, hp: number, maxHp: number, speed: number, name: string, facing: number): CharacterState {
  return { x, y, hp, maxHp, energy: 0, maxEnergy: 100, speed, name, velocityX: 0, velocityY: 0, isAttacking: false, isDefending: false, isHit: false, hitTimer: 0, attackCooldown: 0, facing }
}

function getKeyState(keys: GameScene['keys']): RemoteInput['keys'] {
  return {
    left: keys.A.isDown, right: keys.D.isDown,
    up: keys.W.isDown, down: keys.S.isDown,
    attack: keys.J.isDown, defense: keys.S.isDown,
    skill1: keys.K.isDown, skill2: keys.L.isDown,
    skill3: false,
  }
}

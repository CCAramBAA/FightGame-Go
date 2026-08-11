/**
 * AI 行为引擎 — PVE 模式专用
 * 支持多阶段 BOSS AI、技能释放频率控制、目标选择
 * 所有行为参数由后端配置下发，无硬编码
 */
export interface AIStageConfig {
  stage: number
  hp_threshold: number // 触发该阶段的 HP 百分比阈值 (0-1)
  skill_ids: number[] // 该阶段可用技能 ID 列表
  skill_freq: number // 技能释放频率 (每秒技能释放次数)
  aggression: number // 激进程度 0-1，影响走位和攻击倾向
  move_pattern: 'chase' | 'keep_distance' | 'random' | 'idle'
  description: string
}

export interface AIConfig {
  stages: AIStageConfig[]
  reaction_time_ms: number // 反应延迟 (ms)
  block_chance: number // 格挡概率 0-1
  dash_chance: number // 闪避概率 0-1
}

interface AIEntity {
  x: number
  y: number
  hp: number
  maxHp: number
  energy: number
  maxEnergy: number
  speed: number
  facing: number
  attackCooldown: number
  isHit: boolean
  hitTimer: number
  // 技能冷却独立
  skillCooldowns: Record<number, number>
}

interface AIDecision {
  moveX: number
  moveY: number
  attack: boolean
  skillIndex: number // -1 = no skill
  block: boolean
  dash: boolean
}

export class AIEngine {
  private config: AIConfig
  private currentStage: number = 0
  private lastSkillTime: number = 0
  private frameCount: number = 0
  private reactionTimer: number = 0

  // 训练模式参数
  public paused: boolean = false
  public infiniteHP: boolean = false
  public infiniteEnergy: boolean = false

  constructor(config: AIConfig) {
    this.config = config
  }

  /**
   * 根据敌人当前 HP 百分比切换 AI 阶段
   */
  getCurrentStage(hpRatio: number): AIStageConfig {
    let selected = this.config.stages[0]
    for (const stage of this.config.stages) {
      if (hpRatio <= stage.hp_threshold) {
        selected = stage
      }
    }
    return selected
  }

  /**
   * 每帧调用，返回 AI 决策
   */
  tick(ai: AIEntity, target: AIEntity, delta: number): AIDecision {
    if (this.paused) {
      return { moveX: 0, moveY: 0, attack: false, skillIndex: -1, block: false, dash: false }
    }

    this.frameCount++
    const hpRatio = ai.hp / ai.maxHp
    const stage = this.getCurrentStage(hpRatio)
    const dt = delta / 1000

    // 反应延迟
    this.reactionTimer -= delta
    if (this.reactionTimer > 0) {
      return { moveX: 0, moveY: 0, attack: false, skillIndex: -1, block: false, dash: false }
    }
    this.reactionTimer = this.config.reaction_time_ms

    const decision: AIDecision = {
      moveX: 0,
      moveY: 0,
      attack: false,
      skillIndex: -1,
      block: false,
      dash: false,
    }

    // 移动逻辑
    const dist = target.x - ai.x
    const absDist = Math.abs(dist)

    switch (stage.move_pattern) {
      case 'chase':
        decision.moveX = Math.sign(dist) * ai.speed * dt * stage.aggression
        break
      case 'keep_distance':
        if (absDist < 200) {
          decision.moveX = -Math.sign(dist) * ai.speed * dt * 0.8
        } else if (absDist > 300) {
          decision.moveX = Math.sign(dist) * ai.speed * dt * 0.5
        }
        break
      case 'random':
        decision.moveX = (Math.sin(this.frameCount * 0.05) * ai.speed * dt * 0.6)
        break
      case 'idle':
        decision.moveX = 0
        break
    }

    // 保持朝向面对目标
    ai.facing = dist > 0 ? 1 : -1

    // 攻击决策
    const timeSinceLastSkill = (Date.now() - this.lastSkillTime) / 1000
    const skillInterval = 1 / Math.max(stage.skill_freq, 0.1)

    // 普攻：距离近且激进程度高
    if (absDist < 120 && Math.random() < stage.aggression * 0.3 && ai.attackCooldown <= 0) {
      decision.attack = true
    }

    // 技能释放
    if (timeSinceLastSkill >= skillInterval && stage.skill_ids.length > 0) {
      // 按概率选择一个技能
      const skillIdx = stage.skill_ids[Math.floor(Math.random() * stage.skill_ids.length)]
      const cdKey = `skill_${skillIdx}` as any
      const cooldown = ai.skillCooldowns[skillIdx] || 0
      if (cooldown <= 0) {
        decision.skillIndex = skillIdx
        this.lastSkillTime = Date.now()
      }
    }

    // 格挡：受到攻击时有概率格挡
    if (target.attackCooldown > 0 && absDist < 150 && Math.random() < this.config.block_chance) {
      decision.block = true
    }

    // 闪避
    if (target.attackCooldown > 0 && absDist < 150 && Math.random() < this.config.dash_chance) {
      decision.dash = true
    }

    return decision
  }

  /** 重置引擎状态 */
  reset(): void {
    this.currentStage = 0
    this.lastSkillTime = 0
    this.frameCount = 0
    this.reactionTimer = 0
    this.paused = false
  }
}

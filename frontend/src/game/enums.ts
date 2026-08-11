/**
 * 战斗系统统一枚举定义
 * 前端与后端保持同步
 */

// ===================== 技能标签枚举 =====================
export enum SkillTag {
  /** 免控 - 抵消低层级控制技能 */
  UNSTOPPABLE = 'unstoppable',
  /** 护盾 - 吸收固定伤害，被破甲穿透 */
  SHIELD = 'shield',
  /** 偷取技能 - 临时复制敌方技能替换自身技能栏 */
  STEAL = 'steal',
  /** 吸血 - 造成伤害时同步恢复自身血量 */
  LIFESTEAL = 'lifesteal',
  /** 沉默 - 禁止释放技能 */
  SILENCE = 'silence',
  /** 破甲 - 穿透护盾直接造成伤害 */
  ARMOR_BREAK = 'armor_break',
  /** 眩晕 - 无法移动/攻击/释放技能 */
  STUN = 'stun',
  /** 击飞 - 强制位移 */
  KNOCKUP = 'knockup',
  /** 减速 - 降低移动速度 */
  SLOW = 'slow',
  /** 变身 - 切换全套角色属性/技能/特效 */
  TRANSFORM = 'transform',
  /** 位移 - 冲刺/闪烁 */
  DASH = 'dash',
  /** 持续伤害 */
  DOT = 'dot',
  /** 治疗 */
  HEAL = 'heal',
}

// ===================== 技能类型枚举 =====================
export enum SkillType {
  ACTIVE = 'active',
  PASSIVE = 'passive',
  TRANSFORM = 'transform',
}

// ===================== 交互效果类型枚举 =====================
export enum EffectType {
  /** 偷取技能 */
  STEAL_SKILL = 'steal_skill',
  /** 穿透护盾 */
  PIERCE_SHIELD = 'pierce_shield',
  /** 免疫伤害 */
  IMMUNE_DAMAGE = 'immune_damage',
  /** 优先级覆盖 */
  PRIORITY_OVERRIDE = 'priority_override',
}

// ===================== 游戏模式枚举 =====================
export enum GameMode {
  PVE = 'pve',
  PVP = 'pvp',
  TRAINING = 'training',
}

// ===================== 关卡难度枚举 =====================
export enum StageDifficulty {
  EASY = 'easy',
  NORMAL = 'normal',
  HARD = 'hard',
  BOSS = 'boss',
}

// ===================== 房间状态枚举 =====================
export enum RoomStatus {
  WAITING = 'waiting',
  SELECTING = 'selecting',
  PLAYING = 'playing',
  FINISHED = 'finished',
}

// ===================== 对局结果枚举 =====================
export enum BattleResult {
  WIN = 'win',
  LOSE = 'lose',
  DRAW = 'draw',
}

// ===================== 金币来源枚举 =====================
export enum GoldSourceType {
  PVP_WIN = 'pvp_win',
  PVE_STAR = 'pve_star',
  SHOP_BUY = 'shop_buy',
}

// ===================== 英雄解锁类型 =====================
export enum UnlockType {
  GOLD = 'gold',
  PVE_STAGE = 'pve_stage',
}

// ===================== 好友关系状态 =====================
export enum FriendStatus {
  PENDING = 'pending',
  ACCEPTED = 'accepted',
  BLOCKED = 'blocked',
}

/**
 * 战斗底层通用执行函数接口
 * 所有效果必须实现 execute()
 */
export interface IBattleEffect {
  tag: SkillTag
  /** 技能优先级层级（越大越高） */
  priorityLevel: number
  /** 执行效果 */
  execute(context: BattleContext): void
}

/**
 * 战斗上下文 - 包含当前战斗所有状态
 */
export interface BattleContext {
  /** 施放者 */
  caster: BattleEntity
  /** 目标 */
  target: BattleEntity
  /** 技能基础伤害 */
  baseDamage: number
  /** 当前帧时间 */
  currentTime: number
  /** 全场所有实体 */
  allEntities: BattleEntity[]
}

/**
 * 战斗实体基础属性
 */
export interface BattleEntity {
  id: number
  characterId: number
  hp: number
  maxHp: number
  energy: number
  maxEnergy: number
  energyRegen: number
  speed: number
  attack: number
  defense: number
  /** 当前激活的 buff/debuff */
  buffs: Map<SkillTag, BuffState>
  /** 是否己方 */
  isPlayer: boolean
  /** 当前皮肤ID */
  skinId: number
}

/**
 * Buff/Debuff 状态
 */
export interface BuffState {
  tag: SkillTag
  duration: number  // 剩余帧数
  value: number     // 效果数值
  sourceEntityId: number
}

/**
 * 技能定义（后端配置下发）
 */
export interface SkillDef {
  id: number
  characterId: number
  name: string
  skillType: SkillType
  energyCost: number
  cooldown: number    // 秒
  damage: number
  range: number       // 像素
  priorityLevel: number
  tags: SkillTag[]
  description: string
  effectData: Record<string, number>
}

/**
 * 英雄特殊交互规则（后端配置下发）
 */
export interface HeroSpecialRuleDef {
  id: number
  characterId: number
  name: string
  description: string
  priorityOrder: number
  effectType: EffectType
  condition: object   // JSON 触发条件
  effectData: object  // JSON 效果参数
}

/**
 * 帧数据 - 用于帧同步和回放
 */
export interface FrameData {
  frame: number
  timestamp: number
  actions: PlayerAction[]
  entities: EntitySnapshot[]
}

export interface PlayerAction {
  playerId: number
  actionType: string       // move / attack / skill_1 / skill_2 / skill_3 / skill_4 / idle
  params: Record<string, number | string>
}

export interface EntitySnapshot {
  id: number
  x: number
  y: number
  hp: number
  energy: number
  facing: number
  state: string           // idle / moving / attacking / casting / hit / dead
  animation: string
}

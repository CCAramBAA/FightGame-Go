// 战斗引擎
export enum SkillTag {
  NONE = 'none',
  IMMUNE_CONTROL = 'immune_control',
  SHIELD = 'shield',
  STEAL = 'steal',
  SILENCE = 'silence',
  ARMOR_BREAK = 'armor_break',
  STUN = 'stun',
  KNOCK_UP = 'knock_up',
  SLOW = 'slow',
  LIFE_STEAL = 'life_steal',
  TRANSFORM = 'transform',
}

export interface CharacterConfig {
  id: number; name: string; hp: number; energy: number; energy_regen: number
  speed: number; attack: number; defense: number
}

export interface SkillConfig {
  id: number; character_id: number; name: string; skill_type: 'active' | 'passive' | 'ultimate'
  energy_cost: number; cool_down: number; damage: number; range: number; priority_level: number
  tags: SkillTag[]; effect_params: Record<string, any>
}

export interface SpecialRule {
  id: number; character_id: number; priority_order: number; rule_tag: string
  attacker_skill_tag?: SkillTag; defender_state?: string
  effect: string; effect_type: string; effect_value: number
}

export interface ActiveBuff { tag: SkillTag; remaining_ms: number; value: number; source_skill_id: number }

export interface BattleEntity {
  config: CharacterConfig; skills: SkillConfig[]; specialRules: SpecialRule[]
  hp: number; maxHp: number; energy: number; maxEnergy: number; shield: number
  x: number; y: number; facing: number; isAlive: boolean; buffs: ActiveBuff[]
  skillCooldowns: Record<number, number>; attackCooldown: number
}

export interface BattleResult {
  attacker_hp_change: number; defender_hp_change: number; shield_damage: number
  effects: string[]; priority_used: number
}

export class BattleEngine {
  static calcDamage(attacker: BattleEntity, defender: BattleEntity, skill: SkillConfig): BattleResult {
    const result: BattleResult = { attacker_hp_change: 0, defender_hp_change: 0, shield_damage: 0, effects: [], priority_used: skill.priority_level }
    let rawDamage = skill.damage + attacker.config.attack * 0.5
    let pierceShield = false
    if (skill.tags.includes(SkillTag.ARMOR_BREAK)) { pierceShield = true; result.effects.push('破甲：无视护盾'); rawDamage *= 1.2 }
    if (!pierceShield && defender.shield > 0) {
      if (defender.shield >= rawDamage) { defender.shield -= rawDamage; result.shield_damage = rawDamage; rawDamage = 0 }
      else { rawDamage -= defender.shield; result.shield_damage = defender.shield; defender.shield = 0 }
    }
    const defReduction = defender.config.defense / (defender.config.defense + 100)
    const finalDamage = Math.max(1, rawDamage * (1 - defReduction))
    defender.hp = Math.max(0, defender.hp - finalDamage)
    result.defender_hp_change = -finalDamage
    if (skill.tags.includes(SkillTag.LIFE_STEAL)) { const s = finalDamage * 0.3; attacker.hp = Math.min(attacker.maxHp, attacker.hp + s); result.attacker_hp_change = s }
    if (skill.tags.includes(SkillTag.STUN)) BattleEngine.applyBuff(defender, SkillTag.STUN, 1500, 0, skill.id)
    if (skill.tags.includes(SkillTag.SLOW)) BattleEngine.applyBuff(defender, SkillTag.SLOW, 3000, 0.4, skill.id)
    if (skill.tags.includes(SkillTag.SILENCE)) BattleEngine.applyBuff(defender, SkillTag.SILENCE, 2000, 0, skill.id)
    if (skill.tags.includes(SkillTag.SHIELD)) { const sv = (skill.effect_params && skill.effect_params.shield_value) || (attacker.config.hp * 0.15); attacker.shield += sv }
    return result
  }

  static resolvePriority(as: SkillConfig, ds: SkillConfig | null, _a: BattleEntity, _d: BattleEntity): 'attacker_wins' | 'defender_wins' | 'both_execute' {
    if (!ds) return 'attacker_wins'
    if (as.priority_level > ds.priority_level) return 'attacker_wins'
    if (as.priority_level < ds.priority_level) return 'defender_wins'
    return 'both_execute'
  }

  static applyBuff(entity: BattleEntity, tag: SkillTag, ms: number, v: number, sid: number): void {
    entity.buffs = entity.buffs.filter(b => b.tag !== tag)
    entity.buffs.push({ tag, remaining_ms: ms, value: v, source_skill_id: sid })
  }

  static updateBuffs(entity: BattleEntity, deltaMs: number): string[] {
    const expired: string[] = []
    entity.buffs = entity.buffs.filter(b => { b.remaining_ms -= deltaMs; if (b.remaining_ms <= 0) { expired.push(b.tag); return false } return true })
    return expired
  }

  static getEffectiveSpeed(entity: BattleEntity): number {
    let spd = entity.config.speed
    for (const b of entity.buffs) { if (b.tag === SkillTag.SLOW) spd *= (1 - b.value); if (b.tag === SkillTag.STUN) spd = 0 }
    return spd
  }

  static isControlled(entity: BattleEntity): boolean {
    return entity.buffs.some(b => b.tag === SkillTag.STUN || b.tag === SkillTag.SILENCE)
  }

  static createEntity(cfg: CharacterConfig, skills: SkillConfig[], rules: SpecialRule[], x: number, y: number, f: number): BattleEntity {
    return { config: cfg, skills, specialRules: rules, hp: cfg.hp, maxHp: cfg.hp, energy: 0, maxEnergy: cfg.energy, shield: 0, x, y, facing: f, isAlive: true, buffs: [], skillCooldowns: {}, attackCooldown: 0 }
  }
}

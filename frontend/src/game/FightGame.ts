import Phaser from 'phaser'
import { BootScene } from './scenes/BootScene'
import { GameScene } from './scenes/GameScene'
import { PVEBattleScene } from './scenes/PVEBattleScene'

export class FightGame extends Phaser.Game {
  constructor(parent: string) {
    const config: Phaser.Types.Core.GameConfig = {
      type: Phaser.AUTO,
      width: 1024,
      height: 768,
      parent,
      backgroundColor: '#1a1a2e',
      physics: {
        default: 'arcade',
        arcade: {
          gravity: { x: 0, y: 0 },
          debug: false,
        },
      },
      scene: [BootScene, PVEBattleScene, GameScene],
      scale: {
        mode: Phaser.Scale.FIT,
        autoCenter: Phaser.Scale.CENTER_BOTH,
      },
    }

    super(config)
  }
}

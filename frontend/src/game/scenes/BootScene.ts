import Phaser from 'phaser'

export class BootScene extends Phaser.Scene {
  constructor() {
    super({ key: 'BootScene' })
  }

  preload(): void {
    // 显示加载进度
    const progressBar = this.add.graphics()
    const progressBox = this.add.graphics()
    progressBox.fillStyle(0x222222, 0.8)
    progressBox.fillRect(262, 370, 500, 50)

    this.load.on('progress', (value: number) => {
      progressBar.clear()
      progressBar.fillStyle(0xe94560, 1)
      progressBar.fillRect(272, 380, 480 * value, 30)
    })

    this.load.on('complete', () => {
      progressBar.destroy()
      progressBox.destroy()
    })

    // TODO: 预加载游戏资源
    // this.load.image('player', 'assets/player.png')
    // this.load.spritesheet('fighter', 'assets/fighter.png', { frameWidth: 64, frameHeight: 64 })
  }

  create(): void {
    this.scene.start('GameScene')
  }
}

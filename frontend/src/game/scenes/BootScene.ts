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

    // 加载像素画角色资源
    this.load.image('flame_warrior_sprite', 'assets/flame_warrior_sprite.png')
    this.load.image('flame_warrior_portrait', 'assets/flame_warrior_portrait.png')
    this.load.spritesheet('flame_warrior_spritesheet', 'assets/flame_warrior_spritesheet.png', { frameWidth: 128, frameHeight: 128 })

    this.load.on('complete', () => {
      // 创建 spritesheet 动画
      if (this.anims.exists('flame_idle')) return
      this.anims.create({
        key: 'flame_idle',
        frames: this.anims.generateFrameNumbers('flame_warrior_spritesheet', { start: 0, end: 0 }),
        frameRate: 4,
        repeat: -1,
      })
      this.anims.create({
        key: 'flame_attack',
        frames: this.anims.generateFrameNumbers('flame_warrior_spritesheet', { start: 1, end: 2 }),
        frameRate: 8,
        repeat: 0,
      })
      this.anims.create({
        key: 'flame_hurt',
        frames: [{ key: 'flame_warrior_spritesheet', frame: 3 }],
        frameRate: 4,
        repeat: 0,
      })
    })
  }

  create(): void {
    this.scene.start('GameScene')
    // 通知外部场景已就绪
    this.game.events.emit('ready')
  }
}

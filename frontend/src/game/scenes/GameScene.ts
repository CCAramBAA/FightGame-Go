import Phaser from 'phaser'

export class GameScene extends Phaser.Scene {
  private player1!: Phaser.Physics.Arcade.Sprite
  private player2!: Phaser.Physics.Arcade.Sprite
  private cursors!: Phaser.Types.Input.Keyboard.CursorKeys
  private wasd!: { W: Phaser.Input.Keyboard.Key; A: Phaser.Input.Keyboard.Key; S: Phaser.Input.Keyboard.Key; D: Phaser.Input.Keyboard.Key }

  constructor() {
    super({ key: 'GameScene' })
  }

  create(): void {
    // 游戏场景背景
    this.add.text(512, 384, 'Fight Game - 对战游戏', {
      fontSize: '32px',
      color: '#ffffff',
    }).setOrigin(0.5)

    // TODO: 创建玩家角色
    // this.player1 = this.physics.add.sprite(200, 384, 'fighter')
    // this.player2 = this.physics.add.sprite(824, 384, 'fighter')

    // 键盘输入
    this.cursors = this.input.keyboard!.createCursorKeys()
    this.wasd = {
      W: this.input.keyboard!.addKey(Phaser.Input.Keyboard.KeyCodes.W),
      A: this.input.keyboard!.addKey(Phaser.Input.Keyboard.KeyCodes.A),
      S: this.input.keyboard!.addKey(Phaser.Input.Keyboard.KeyCodes.S),
      D: this.input.keyboard!.addKey(Phaser.Input.Keyboard.KeyCodes.D),
    }
  }

  update(): void {
    // TODO: 游戏逻辑更新
    // 玩家1 方向键控制
    // 玩家2 WASD 控制
  }
}

import { Howl } from 'howler'

// 音效管理器
class AudioManager {
  private sounds: Map<string, Howl> = new Map()
  private musicVolume: number = 0.5
  private sfxVolume: number = 0.7
  private currentMusic: Howl | null = null

  load(key: string, src: string | string[]): void {
    const sound = new Howl({
      src: Array.isArray(src) ? src : [src],
      preload: true,
      volume: this.sfxVolume,
    })
    this.sounds.set(key, sound)
  }

  play(key: string): void {
    const sound = this.sounds.get(key)
    if (sound) sound.play()
  }

  playMusic(key: string): void {
    if (this.currentMusic) {
      this.currentMusic.stop()
    }
    const music = new Howl({
      src: [key],
      loop: true,
      volume: this.musicVolume,
      preload: true,
    })
    this.currentMusic = music
    music.play()
  }

  stopMusic(): void {
    if (this.currentMusic) {
      this.currentMusic.stop()
      this.currentMusic = null
    }
  }

  setMusicVolume(vol: number): void {
    this.musicVolume = vol
    if (this.currentMusic) {
      this.currentMusic.volume(vol)
    }
  }

  setSfxVolume(vol: number): void {
    this.sfxVolume = vol
  }
}

export const audioManager = new AudioManager()

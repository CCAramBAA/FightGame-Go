import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'
import router from '@/router'

export interface CharacterInfo {
  id: number
  name: string
  description: string
  hp: number
  energy: number
  speed: number
  attack: number
  defense: number
  avatar_path: string
  skills: any[]
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string>((() => {
    const t = localStorage.getItem('token')
    return t && t !== 'undefined' && t !== 'null' ? t : ''
  })())
  const username = ref<string>(localStorage.getItem('username') || '')
  const userId = ref<string>(localStorage.getItem('userId') || '')
  const nickname = ref<string>(localStorage.getItem('nickname') || '')
  const gold = ref<number>(0)
  const rankScore = ref<number>(1000)
  const winCount = ref<number>(0)
  const loseCount = ref<number>(0)
  const myCharacters = ref<CharacterInfo[]>([])
  const initialized = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(t: string, id: string, name: string, nick?: string) {
    if (!t || t === 'undefined' || t === 'null') {
      console.error('setAuth: invalid token, clearing auth')
      logout()
      return
    }
    token.value = t
    userId.value = id || ''
    username.value = name
    nickname.value = nick || name
    localStorage.setItem('token', t)
    localStorage.setItem('userId', id || '')
    localStorage.setItem('username', name)
    localStorage.setItem('nickname', nickname.value)
  }

  function logout() {
    // 调用后端退出接口
    api.post('/logout').catch(() => {})
    token.value = ''
    username.value = ''
    userId.value = ''
    nickname.value = ''
    gold.value = 0
    rankScore.value = 1000
    winCount.value = 0
    loseCount.value = 0
    myCharacters.value = []
    initialized.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('userId')
    localStorage.removeItem('nickname')
    router.push('/')
  }

  // 从服务器同步用户信息
  async function fetchUserInfo() {
    if (!token.value) return
    try {
      const [profileRes, rankRes]: any[] = await Promise.all([
        api.get('/profile'),
        api.get('/profile/rank'),
      ])
      const profile = profileRes.data || profileRes
      const rank = rankRes.data || rankRes
      nickname.value = profile.nickname || username.value
      gold.value = profile.gold || 0
      rankScore.value = rank.score || 1000
      winCount.value = rank.win_count || 0
      loseCount.value = rank.lose_count || 0
      initialized.value = true
    } catch (err) {
      console.error('获取用户信息失败', err)
    }
  }

  // 加载我的角色列表
  async function fetchMyCharacters() {
    if (!token.value) return
    try {
      const res: any = await api.get('/profile/characters')
      const chars = res.data || res
      myCharacters.value = chars || []
    } catch (err) {
      console.error('获取角色列表失败', err)
    }
  }

  return {
    token, username, userId, nickname,
    gold, rankScore, winCount, loseCount,
    myCharacters, initialized, isLoggedIn,
    setAuth, logout, fetchUserInfo, fetchMyCharacters,
  }
})

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface RoomInfo {
  id: string
  host_id: number
  guest_id: number
  host_char_id: number
  guest_char_id: number
  host_skin_id: number
  guest_skin_id: number
  host_ready: boolean
  guest_ready: boolean
  status: string
  player_count: number
  created_at: number
}

export interface CharacterInfo {
  id: number
  name: string
  hp: number
  energy: number
  speed: number
  attack: number
  defense: number
  avatar_path: string
  skills: Array<{
    id: number
    name: string
    damage: number
    energy_cost: number
    cooldown: number
    tags: string[]
  }>
}

export interface BattleState {
  myHP: number
  myMaxHP: number
  myEnergy: number
  myMaxEnergy: number
  enemyHP: number
  enemyMaxHP: number
  enemyEnergy: number
  enemyMaxEnergy: number
  myCharName: string
  enemyCharName: string
  round: number
  isMyTurn: boolean
  status: 'waiting' | 'countdown' | 'playing' | 'finished'
  winner: string
}

export const useGameStore = defineStore('game', () => {
  // 房间状态
  const currentRoom = ref<RoomInfo | null>(null)
  const rooms = ref<RoomInfo[]>([])
  const inRoom = computed(() => !!currentRoom.value)
  const isHost = computed(() => {
    if (!currentRoom.value) return false
    return currentRoom.value.host_id === userID.value
  })

  // 战斗状态
  const battle = ref<BattleState>({
    myHP: 100, myMaxHP: 100,
    myEnergy: 0, myMaxEnergy: 100,
    enemyHP: 100, enemyMaxHP: 100,
    enemyEnergy: 0, enemyMaxEnergy: 100,
    myCharName: '', enemyCharName: '',
    round: 0, isMyTurn: false,
    status: 'waiting',
    winner: '',
  })

  // 本地数据
  const selectedCharID = ref<number>(0)
  const selectedSkinID = ref<number>(0)
  const selectedChar = ref<CharacterInfo | null>(null)
  const opponentCharID = ref<number>(0)
  const opponentChar = ref<CharacterInfo | null>(null)
  const opponentCharName = ref<string>('')
  const opponentID = ref<number>(0)
  const userID = ref<number>(parseInt(localStorage.getItem('userId') || '0'))
  const isReady = ref(false)
  const gameCountdown = ref(0)
  const mySkins = ref<Array<{ id: number; character_id: number; name: string }>>([])

  // 帧同步
  const currentFrame = ref(0)
  const localFrame = ref(0)

  function setCurrentRoom(room: RoomInfo | null) {
    currentRoom.value = room
    isReady.value = false
    gameCountdown.value = 0
  }

  function setRooms(list: RoomInfo[]) {
    rooms.value = list
  }

  function updateRoomState(data: Partial<RoomInfo>) {
    if (currentRoom.value) {
      Object.assign(currentRoom.value, data)
    }
  }

  function setIsReady(ready: boolean) {
    isReady.value = ready
  }

  function setGameCountdown(count: number) {
    gameCountdown.value = count
    battle.value.status = 'countdown'
  }

  function startBattle(
    hostCharID: number, guestCharID: number,
    hostName: string, guestName: string,
    hostHP: number, guestHP: number,
    hostMaxHP: number, guestMaxHP: number,
  ) {
    const isMeHost = currentRoom.value?.host_id === userID.value
    battle.value = {
      myHP: isMeHost ? hostHP : guestHP,
      myMaxHP: isMeHost ? hostMaxHP : guestMaxHP,
      myEnergy: 0,
      myMaxEnergy: 100,
      enemyHP: isMeHost ? guestHP : hostHP,
      enemyMaxHP: isMeHost ? guestMaxHP : hostMaxHP,
      enemyEnergy: 0,
      enemyMaxEnergy: 100,
      myCharName: isMeHost ? hostName : guestName,
      enemyCharName: isMeHost ? guestName : hostName,
      round: 0, isMyTurn: true,
      status: 'playing',
      winner: '',
    }
    currentFrame.value = 0
    localFrame.value = 0
  }

  function updateBattle(data: Partial<BattleState>) {
    Object.assign(battle.value, data)
  }

  function updateHP(myHP: number, enemyHP: number) {
    battle.value.myHP = Math.max(0, myHP)
    battle.value.enemyHP = Math.max(0, enemyHP)
  }

  function setWinner(winner: string) {
    battle.value.winner = winner
    battle.value.status = 'finished'
  }

  function endBattle() {
    setCurrentRoom(null)
    battle.value.status = 'finished'
  }

  function incrementFrame() {
    localFrame.value++
  }

  // 重置
  function reset() {
    setCurrentRoom(null)
    selectedCharID.value = 0
    selectedSkinID.value = 0
    selectedChar.value = null
    opponentCharID.value = 0
    opponentChar.value = null
    isReady.value = false
    gameCountdown.value = 0
    battle.value.status = 'waiting'
    battle.value.winner = ''
  }

  return {
    currentRoom, rooms, inRoom, isHost,
    battle,
    selectedCharID, selectedSkinID, selectedChar, opponentCharID, opponentChar, opponentCharName, opponentID,
    userID, isReady, gameCountdown,
    currentFrame, localFrame, mySkins,
    setCurrentRoom, setRooms, updateRoomState,
    setIsReady, setGameCountdown, startBattle, updateBattle,
    updateHP, setWinner, endBattle, incrementFrame, reset,
  }
})

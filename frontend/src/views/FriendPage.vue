<template>
  <div class="friend-page">
    <div class="friend-header">
      <button class="btn-back" @click="$router.push('/lobby')">&larr; 返回大厅</button>
      <h2>好友</h2>
    </div>

    <div class="add-friend">
      <input
        v-model="friendIdInput"
        placeholder="输入好友ID添加"
        type="number"
        class="friend-input"
      />
      <button class="btn-add" :disabled="adding" @click="addFriend">
        {{ adding ? '添加中...' : '添加好友' }}
      </button>
    </div>

    <div class="friend-list-card">
      <div class="card-header">好友列表 ({{ friends.length }})</div>
      <div class="friend-list">
        <div v-if="loading" class="loading">加载中...</div>
        <div v-for="f in friends" :key="f.id" class="friend-item">
          <div class="friend-avatar">{{ f.nickname?.[0] || f.username?.[0] || '?' }}</div>
          <div class="friend-info">
            <div class="friend-name">{{ f.nickname || f.username }}</div>
            <div class="friend-status">
              <span class="status-dot" :class="{ online: f.online }" />
              {{ f.online ? '在线' : '离线' }}
            </div>
          </div>
          <div class="friend-actions">
            <button
              class="btn-invite"
              :disabled="!f.online"
              @click="inviteFriend(f)"
            >邀请对战</button>
            <button class="btn-remove" @click="removeFriend(f)">删除</button>
          </div>
        </div>
        <div v-if="!loading && friends.length === 0" class="empty">还没有好友，快去添加吧</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useWebSocket } from '@/utils/websocket'
import { useGameStore } from '@/store/game'

const router = useRouter()
const gameStore = useGameStore()
const { on, off } = useWebSocket()

interface Friend {
  id: number
  username: string
  nickname: string
  online: boolean
}

const friends = ref<Friend[]>([])
const friendIdInput = ref('')
const loading = ref(false)
const adding = ref(false)

async function loadFriends() {
  loading.value = true
  try {
    const res: any = await api.get('/friends')
    friends.value = (res.data?.data || res.data || res) || []
  } catch (err) {
    console.error('加载好友列表失败', err)
  } finally {
    loading.value = false
  }
}

async function addFriend() {
  const id = parseInt(friendIdInput.value)
  if (!id) {
    alert('请输入有效的好友ID')
    return
  }
  adding.value = true
  try {
    await api.post('/friends/add', { friend_id: id })
    alert('添加好友请求已发送')
    friendIdInput.value = ''
    loadFriends()
  } catch (err: any) {
    alert(err?.response?.data?.message || '添加失败')
  } finally {
    adding.value = false
  }
}

async function removeFriend(f: Friend) {
  if (!confirm(`确定删除好友 ${f.nickname || f.username}？`)) return
  try {
    await api.delete(`/friends/${f.id}`)
    alert('已删除好友')
    loadFriends()
  } catch (err) { /* cancelled */ }
}

function inviteFriend(f: Friend) {
  const { send } = useWebSocket()
  send({ type: 'invite_friend', data: { friend_id: f.id } })
  alert(`已向 ${f.nickname || f.username} 发送对战邀请`)
}

onMounted(() => {
  loadFriends()
  on('invite_result', handleInviteResult)
  on('game_start', handleInviteGameStart)
})

onUnmounted(() => {
  off('invite_result', handleInviteResult)
  off('game_start', handleInviteGameStart)
})

function handleInviteResult(msg: any) {
  const data = msg.data || msg
  if (data.accepted) {
    gameStore.opponentCharName = data.opponent_name || '对手'
    router.push('/game')
  } else {
    alert('对方拒绝了邀请')
  }
}

function handleInviteGameStart(msg: any) {
  const data = msg.data || msg
  gameStore.opponentCharName = data.opponent_name || '对手'
  router.push('/game')
}
</script>

<style scoped>
.friend-page { background: #0a1628; min-height: 100vh; color: #fff; padding: 20px; }
.friend-header { display: flex; align-items: center; gap: 20px; margin-bottom: 20px; }
.friend-header h2 { font-size: 1.5rem; color: #e94560; margin: 0; }
.btn-back { background: transparent; border: 1px solid #666; color: #aaa; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
.btn-back:hover { border-color: #fff; color: #fff; }
.add-friend { display: flex; gap: 10px; margin-bottom: 20px; }
.friend-input { padding: 8px 12px; border: 1px solid #1f4287; border-radius: 6px; background: #162447; color: #fff; width: 220px; font-size: 0.9rem; }
.friend-input::placeholder { color: #666; }
.btn-add {
  padding: 8px 16px; border: none; border-radius: 6px; cursor: pointer;
  background: #e94560; color: #fff; font-size: 0.9rem;
}
.btn-add:disabled { opacity: 0.4; cursor: not-allowed; }
.friend-list-card { background: #162447; border: 1px solid #1f4287; border-radius: 8px; }
.card-header { padding: 12px 16px; border-bottom: 1px solid #1f4287; font-size: 0.95rem; color: #ddd; }
.friend-item { display: flex; align-items: center; padding: 12px 16px; border-bottom: 1px solid #1f4287; }
.friend-item:last-child { border-bottom: none; }
.friend-avatar {
  width: 40px; height: 40px; border-radius: 50%; background: #e94560;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.1rem; font-weight: bold; margin-right: 14px; flex-shrink: 0;
}
.friend-info { flex: 1; }
.friend-name { font-size: 0.95rem; font-weight: 500; }
.friend-status { display: flex; align-items: center; gap: 6px; font-size: 0.8rem; color: #888; margin-top: 2px; }
.status-dot { width: 8px; height: 8px; border-radius: 50%; background: #666; display: inline-block; }
.status-dot.online { background: #67c23a; }
.friend-actions { display: flex; gap: 8px; }
.btn-invite {
  padding: 4px 12px; border: none; border-radius: 4px; cursor: pointer;
  background: #3b82f6; color: #fff; font-size: 0.8rem;
}
.btn-invite:disabled { opacity: 0.4; cursor: not-allowed; }
.btn-invite:not(:disabled):hover { background: #2563eb; }
.btn-remove {
  padding: 4px 12px; border: none; border-radius: 4px; cursor: pointer;
  background: #e94560; color: #fff; font-size: 0.8rem;
}
.btn-remove:hover { background: #d63852; }
.loading { text-align: center; color: #aaa; padding: 30px; }
.empty { text-align: center; color: #666; padding: 30px; }
</style>

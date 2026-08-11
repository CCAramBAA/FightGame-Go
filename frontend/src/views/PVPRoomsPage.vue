<template>
  <div class="page">
    <GlobalNav showBack @openSettings="showSettings = true" @openFriends="showFriends = true">
      <template #left>
        <button class="btn-back" @click="$router.push('/lobby')">← 返回大厅</button>
      </template>
    </GlobalNav>
    
    <main class="content">
      <ArtPlaceholder label="房间列表背景" width="100%" height="100%" bgColor="#0f0f1e" class="bg" />
      
      <div class="room-grid" v-if="rooms.length > 0">
        <div v-for="r in rooms" :key="r.id" class="room-card">
          <ArtPlaceholder label="房间" width="100%" height="80" bgColor="#1a1a3e" />
          <div class="room-info">
            <span>房主: ID{{ r.host_id }}</span>
            <span>人数: {{ r.player_count || 1 }}/2</span>
          </div>
          <button 
            class="btn-join" 
            :disabled="(r.player_count || 1) >= 2"
            @click="joinRoom(r)"
          >
            {{ (r.player_count || 1) >= 2 ? '已满' : '进入房间' }}
          </button>
        </div>
      </div>
      <p v-else class="empty">暂无房间，快来创建一个吧！</p>
      
      <button class="btn-create" @click="createRoom">+ 创建房间</button>
    </main>

    <SettingsDialog :visible="showSettings" @close="showSettings = false" />
    <FriendsDialog :visible="showFriends" @close="showFriends = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useWebSocket } from '@/utils/websocket'
import GlobalNav from '@/components/GlobalNav.vue'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import FriendsDialog from '@/components/FriendsDialog.vue'

const router = useRouter()
const rooms = ref<any[]>([])
const showSettings = ref(false)
const showFriends = ref(false)

const { send, isConnected } = useWebSocket()

function handleRoomList(msg: any) {
  const list = msg.data || msg || []
  if (Array.isArray(list)) rooms.value = list
}

onMounted(async () => {
  // HTTP 获取
  try {
    const res: any = await api.get('/rooms')
    const list = res.data || res
    if (Array.isArray(list)) rooms.value = list
  } catch { /* fallback */ }
  // WS 请求
  if (isConnected.value) send({ type: 'get_room_list' })
})

async function joinRoom(room: any) {
  try {
    await api.post(`/rooms/join/${room.id}`)
    router.push(`/pvp-select?roomId=${room.id}`)
  } catch (err: any) {
    alert(err?.data?.message || '加入房间失败')
  }
}

async function createRoom() {
  try {
    // 通过 WebSocket 创建房间
    send({ type: 'create_room' })
    alert('房间已创建')
    // 等待 room_created 消息
  } catch {
    alert('创建房间失败')
  }
}
</script>

<style scoped>
.page { width: 100vw; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
.content { flex: 1; position: relative; padding: 24px; }
.bg { position: absolute; inset: 0; z-index: 0; }
.room-grid {
  position: relative; z-index: 1;
  display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px; max-height: calc(100vh - 160px); overflow-y: auto;
}
.room-card {
  background: rgba(20,20,40,0.9); border: 1px solid #3a3a5e; border-radius: 12px;
  padding: 12px; display: flex; flex-direction: column; gap: 10px;
}
.room-info { display: flex; justify-content: space-between; color: #999; font-size: 13px; }
.btn-join {
  padding: 8px; border-radius: 8px; border: none; cursor: pointer;
  background: linear-gradient(135deg, #4fc3f7, #0288d1); color: #fff; font-size: 14px;
}
.btn-join:disabled { background: #2a2a4a; color: #666; cursor: not-allowed; }
.btn-create {
  position: fixed; bottom: 30px; right: 30px; z-index: 10;
  padding: 14px 24px; border-radius: 50px; border: none; cursor: pointer;
  background: linear-gradient(135deg, #ffd700, #ff8c00); color: #1a1a1a;
  font-size: 16px; font-weight: bold; box-shadow: 0 4px 20px rgba(255,140,0,.4);
}
.empty { color: #777; text-align: center; margin-top: 100px; font-size: 16px; }
.btn-back {
  background: rgba(255,255,255,.1); border: 1px solid #3a3a5e; color: #ccc;
  padding: 6px 16px; border-radius: 6px; cursor: pointer;
}
</style>

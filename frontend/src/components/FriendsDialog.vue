<template>
  <ModalOverlay :visible="visible" title="好友列表" maxWidth="480px" @close="$emit('close')">
    <div class="friends">
      <!-- 添加好友 -->
      <div class="add-friend">
        <input v-model="searchId" class="input-sm" placeholder="输入玩家 ID 搜索" />
        <button class="btn-sm" @click="searchFriend">搜索</button>
      </div>
      
      <!-- 好友列表 -->
      <div class="list">
        <p v-if="friends.length === 0" class="empty">暂无好友</p>
        <div v-for="f in friends" :key="f.id" class="friend-row" :class="{ online: f.online }">
          <span class="status-dot" :style="{ background: f.online ? '#22c55e' : '#666' }"></span>
          <span class="name">{{ f.nickname || f.username }}</span>
          <span class="tag" v-if="f.online">在线</span>
          <span class="tag off" v-else>离线</span>
          <button class="btn-sm" v-if="f.online" @click="inviteToBattle(f)">邀请对战</button>
          <button class="btn-sm btn-del" @click="deleteFriend(f)">删除</button>
        </div>
      </div>
    </div>
    <template #footer>
      <button class="btn-ok" @click="$emit('close')">关闭</button>
    </template>
  </ModalOverlay>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import api from '@/api'
import ModalOverlay from './ModalOverlay.vue'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()

const friends = ref<any[]>([])
const searchId = ref('')

onMounted(fetchFriends)
watch(() => props.visible, (v) => { if (v) fetchFriends() })

async function fetchFriends() {
  try {
    const res: any = await api.get('/friends')
    friends.value = (res.data || res) || []
  } catch { friends.value = [] }
}

async function searchFriend() {
  if (!searchId.value) return
  try {
    await api.post('/friends/add', { target_id: Number(searchId.value) })
    alert('好友申请已发送')
    searchId.value = ''
  } catch (err: any) {
    alert(err?.data?.message || '搜索失败')
  }
}

async function deleteFriend(f: any) {
  if (!confirm(`确认删除好友 ${f.nickname || f.username}?`)) return
  try {
    await api.delete(`/friends/${f.id}`)
    friends.value = friends.value.filter(x => x.id !== f.id)
  } catch { alert('删除失败') }
}

async function inviteToBattle(f: any) {
  try {
    await api.post('/friends/invite', { friend_id: f.id })
    alert('对战邀请已发送')
  } catch (err: any) {
    alert(err?.data?.message || '邀请失败')
  }
}
</script>

<style scoped>
.friends { display: flex; flex-direction: column; gap: 14px; }
.add-friend { display: flex; gap: 8px; }
.input-sm {
  flex: 1; padding: 8px 12px; background: #2a2a4a; border: 1px solid #3a3a5e;
  border-radius: 6px; color: #ccc; font-size: 13px; outline: none;
}
.btn-sm {
  padding: 8px 14px; background: #2a2a4a; border: 1px solid #3a3a5e;
  border-radius: 6px; color: #ccc; cursor: pointer; font-size: 13px;
}
.btn-sm:hover { background: #3a3a5e; }
.btn-del { border-color: #dc262633; color: #dc2626; }
.btn-del:hover { background: #dc262622; }
.empty { color: #666; text-align: center; padding: 20px; }
.list { max-height: 300px; overflow-y: auto; }
.friend-row { display: flex; align-items: center; gap: 10px; padding: 10px; border-bottom: 1px solid #2a2a4a; }
.friend-row .name { flex: 1; color: #ccc; font-size: 14px; }
.friend-row .tag { font-size: 11px; padding: 2px 8px; border-radius: 10px; background: #22c55e22; color: #22c55e; }
.friend-row .tag.off { background: #66622; color: #666; }
.status-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.btn-ok {
  background: linear-gradient(135deg, #ffd700, #ff8c00); border: none;
  padding: 10px 32px; border-radius: 8px; cursor: pointer; font-weight: bold;
}
</style>

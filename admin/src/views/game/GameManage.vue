<template>
  <div class="game-manage">
    <h2>游戏管理</h2>
    <el-table :data="rooms" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="房间ID" />
      <el-table-column prop="host_id" label="房主ID" />
      <el-table-column prop="guest_id" label="玩家ID">
        <template #default="{ row }">
          <span v-if="row.guest_id">{{ row.guest_id }}</span>
          <el-tag v-else size="small" type="info">等待中</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="tagType(row.status)">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="closeRoom(row)">关闭房间</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const rooms = ref<any[]>([])
const loading = ref(false)

function tagType(status: string) {
  switch (status) {
    case 'waiting': return 'warning'
    case 'playing': return 'success'
    case 'selecting': return 'primary'
    default: return 'info'
  }
}

async function fetchRooms() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/rooms')
    rooms.value = (res.data || res) || []
  } catch (err) {
    console.error('获取房间列表失败', err)
  } finally {
    loading.value = false
  }
}

async function closeRoom(row: any) {
  try {
    await ElMessageBox.confirm(`确定要关闭房间 ${row.id} 吗？`, '提示')
    await api.delete(`/admin/rooms/${row.id}`)
    ElMessage.success('房间已关闭')
    fetchRooms()
  } catch (err) {
    // cancelled or error
  }
}

onMounted(() => fetchRooms())
</script>

<style scoped>
.game-manage h2 { margin-bottom: 20px; }
</style>

<template>
  <div class="dashboard">
    <h2>数据概览</h2>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总用户数</template>
          <div class="card-value">{{ stats.totalUsers }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>今日对局</template>
          <div class="card-value">{{ stats.battlesToday }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>活跃房间</template>
          <div class="card-value">{{ stats.activeRooms }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>在线玩家</template>
          <div class="card-value">{{ stats.onlineUsers }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总角色数</template>
          <div class="card-value">{{ stats.totalCharacters }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总皮肤数</template>
          <div class="card-value">{{ stats.totalSkins }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import api from '@/api'

const stats = reactive({
  totalUsers: 0,
  battlesToday: 0,
  activeRooms: 0,
  onlineUsers: 0,
  totalCharacters: 0,
  totalSkins: 0,
})

onMounted(async () => {
  try {
    const res: any = await api.get('/admin/dashboard')
    const d = res.data || res
    stats.totalUsers = d.total_users || 0
    stats.battlesToday = d.battles_today || 0
    stats.activeRooms = d.active_rooms || 0
    stats.onlineUsers = d.online_users || 0
    stats.totalCharacters = d.total_characters || 0
    stats.totalSkins = d.total_skins || 0
  } catch (err) {
    console.error('获取仪表盘数据失败', err)
  }
})
</script>

<style scoped>
.dashboard h2 { margin-bottom: 20px; }
.card-value { font-size: 2rem; font-weight: bold; color: #e94560; }
</style>

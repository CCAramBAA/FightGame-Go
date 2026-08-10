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
          <template #header>今日活跃</template>
          <div class="card-value">{{ stats.todayActive }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>总对局数</template>
          <div class="card-value">{{ stats.totalGames }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>在线玩家</template>
          <div class="card-value">{{ stats.onlinePlayers }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'

const stats = reactive({
  totalUsers: 0,
  todayActive: 0,
  totalGames: 0,
  onlinePlayers: 0,
})
</script>

<style scoped>
.dashboard h2 { margin-bottom: 20px; }
.card-value {
  font-size: 2rem;
  font-weight: bold;
  color: #e94560;
}
</style>

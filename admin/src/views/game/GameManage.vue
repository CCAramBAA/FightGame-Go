<template>
  <div class="game-manage">
    <h2>游戏管理</h2>
    <el-table :data="games" border style="width: 100%">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="room_name" label="房间名" />
      <el-table-column prop="host_id" label="房主" />
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 'waiting' ? 'warning' : 'success'">
            {{ row.status === 'waiting' ? '等待中' : '对战中' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default>
          <el-button size="small" type="danger">关闭房间</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const games = ref<any[]>([])
</script>

<style scoped>
.game-manage h2 { margin-bottom: 20px; }
</style>

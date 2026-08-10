<template>
  <div class="user-list">
    <h2>用户管理</h2>
    <el-table :data="users" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="win_count" label="胜场" />
      <el-table-column prop="lose_count" label="败场" />
      <el-table-column prop="created_at" label="注册时间" />
      <el-table-column label="操作" width="180">
        <template #default>
          <el-button size="small" type="primary">编辑</el-button>
          <el-button size="small" type="danger">封禁</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface User {
  id: number
  username: string
  nickname: string
  win_count: number
  lose_count: number
  created_at: string
}

const users = ref<User[]>([])
</script>

<style scoped>
.user-list h2 { margin-bottom: 20px; }
</style>

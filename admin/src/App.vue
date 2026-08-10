<template>
  <el-container class="app-container">
    <el-aside width="220px">
      <div class="logo">
        <h2>FightGame 管理</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#1a1a2e"
        text-color="#a0a0a0"
        active-text-color="#e94560"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据概览</span>
        </el-menu-item>
        <el-menu-item index="/user">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/game">
          <el-icon><VideoGame /></el-icon>
          <span>游戏管理</span>
        </el-menu-item>
        <el-menu-item index="/log">
          <el-icon><Document /></el-icon>
          <span>日志查看</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <span>管理后台</span>
        <el-button type="danger" size="small" @click="logout">退出登录</el-button>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const activeMenu = computed(() => route.path)

function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>

<style scoped>
.app-container {
  height: 100vh;
}
.el-aside {
  background: #1a1a2e;
  overflow: hidden;
}
.logo {
  padding: 20px;
  text-align: center;
  color: #e94560;
  border-bottom: 1px solid #333;
}
.el-aside,
.el-menu {
  border-right: none !important;
}
.el-header {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #eee;
  padding: 0 20px;
}
.el-main {
  background: #f5f7fa;
}
</style>

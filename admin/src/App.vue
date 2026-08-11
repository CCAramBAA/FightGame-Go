<template>
  <el-container class="layout">
    <el-aside width="220px">
      <div class="logo">FightGame Admin</div>
      <el-menu :default-active="activeMenu" router background-color="#1a1a2e" text-color="#ccc" active-text-color="#e94560">
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon><span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon><span>用户管理</span>
        </el-menu-item>
        <el-sub-menu index="config">
          <template #title><el-icon><Setting /></el-icon><span>配置管理</span></template>
          <el-menu-item index="/characters">角色管理</el-menu-item>
          <el-menu-item index="/skills">技能管理</el-menu-item>
          <el-menu-item index="/skins">皮肤管理</el-menu-item>
          <el-menu-item index="/special-rules">特殊交互规则</el-menu-item>
          <el-menu-item index="/stages">PVE关卡</el-menu-item>
          <el-menu-item index="/shop">商城商品</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/game">
          <el-icon><Monitor /></el-icon><span>游戏管理</span>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Document /></el-icon><span>系统日志</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { DataAnalysis, User, Setting, Monitor, Document } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
.layout { height: 100vh; }
.el-aside { background: #1a1a2e; overflow-y: auto; }
.logo { color: #e94560; font-size: 1.2rem; font-weight: bold; padding: 20px 16px; text-align: center; border-bottom: 1px solid #333; }
.el-main { background: #f0f2f5; min-height: 100vh; }
</style>

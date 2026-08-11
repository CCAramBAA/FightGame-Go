<template>
  <div class="connection-page">
    <ArtPlaceholder label="游戏主视觉海报" sub="全屏背景图" width="100vw" height="100vh" bgColor="#0a0a1e" class="bg" />
    
    <div class="content">
      <ArtPlaceholder label="游戏 LOGO" sub="大图居中" width="260" height="100" bgColor="#1a1a3e" />
      <LoadingSpinner text="正在连接服务器..." />
    </div>

    <NetworkErrorDialog 
      :visible="showError" 
      message="无法连接服务器，请检查网络"
      @close="handleExit"
      @retry="checkConnection"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { useUserStore } from '@/store/user'
import ArtPlaceholder from '@/components/ArtPlaceholder.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import NetworkErrorDialog from '@/components/NetworkErrorDialog.vue'

const router = useRouter()
const userStore = useUserStore()
const showError = ref(false)

async function checkConnection() {
  showError.value = false
  try {
    await api.get('/health')
    const token = localStorage.getItem('token')
    if (token) {
      try {
        await userStore.fetchUserInfo()
        router.replace('/lobby')
        return
      } catch {
        localStorage.removeItem('token')
      }
    }
    router.replace('/login')
  } catch {
    showError.value = true
  }
}

function handleExit() {
  if (window.close) window.close()
}
</script>

<style scoped>
.connection-page {
  width: 100vw;
  height: 100vh;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
.bg {
  position: absolute;
  inset: 0;
  z-index: 0;
}
.content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 40px;
}
</style>

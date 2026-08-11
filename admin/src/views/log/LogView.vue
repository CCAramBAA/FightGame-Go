<template>
  <div class="log-view">
    <h2>系统日志</h2>
    <div class="log-toolbar">
      <el-input v-model="searchText" placeholder="搜索日志..." clearable style="width: 300px" />
      <el-button type="primary" @click="fetchLogs" style="margin-left: 10px">刷新</el-button>
    </div>
    <el-table :data="filteredLogs" border style="width: 100%; margin-top: 16px" max-height="600" v-loading="loading">
      <el-table-column prop="time" label="时间" width="200" />
      <el-table-column prop="level" label="级别" width="100">
        <template #default="{ row }">
          <el-tag :type="row.level === 'ERROR' ? 'danger' : row.level === 'WARN' ? 'warning' : 'info'" size="small">
            {{ row.level }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="message" label="消息" show-overflow-tooltip />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import api from '@/api'

interface LogEntry {
  time: string
  level: string
  message: string
}

const logs = ref<LogEntry[]>([])
const searchText = ref('')
const loading = ref(false)

const filteredLogs = computed(() => {
  if (!searchText.value) return logs.value
  return logs.value.filter((l) => l.message.toLowerCase().includes(searchText.value.toLowerCase()))
})

function parseLogLine(line: string): LogEntry {
  // Format: "[LEVEL] 2026-08-11 12:00:00 message..."
  const match = line.match(/^\[(\w+)\]\s*(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s*(.*)$/)
  if (match) {
    return { level: match[1], time: match[2], message: match[3] }
  }
  return { level: 'INFO', time: '', message: line }
}

async function fetchLogs() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/logs')
    const lines = (res.data || res) || []
    logs.value = lines.map((l: string) => parseLogLine(l)).reverse()
  } catch (err) {
    console.error('获取日志失败', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchLogs())
</script>

<style scoped>
.log-view h2 { margin-bottom: 20px; }
.log-toolbar { display: flex; align-items: center; }
</style>

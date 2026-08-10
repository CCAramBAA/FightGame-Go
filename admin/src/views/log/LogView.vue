<template>
  <div class="log-view">
    <h2>系统日志</h2>
    <el-input
      v-model="searchText"
      placeholder="搜索日志..."
      clearable
      style="width: 300px; margin-bottom: 20px"
    />
    <el-table :data="filteredLogs" border style="width: 100%" max-height="600">
      <el-table-column prop="time" label="时间" width="180" />
      <el-table-column prop="level" label="级别" width="100">
        <template #default="{ row }">
          <el-tag
            :type="row.level === 'ERROR' ? 'danger' : row.level === 'WARN' ? 'warning' : 'info'"
            size="small"
          >
            {{ row.level }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="message" label="消息" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface LogEntry {
  time: string
  level: string
  message: string
}

const logs = ref<LogEntry[]>([
  { time: '2026-08-10 17:00:00', level: 'INFO', message: 'Server started' },
])

const searchText = ref('')

const filteredLogs = computed(() =>
  logs.value.filter((l) => l.message.includes(searchText.value))
)
</script>

<style scoped>
.log-view h2 { margin-bottom: 20px; }
</style>

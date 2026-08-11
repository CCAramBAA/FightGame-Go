<template>
  <div class="stage-manage">
    <h2>PVE关卡管理</h2>
    <el-button type="primary" @click="openAdd">新增关卡</el-button>
    <el-table :data="stages" border style="width: 100%; margin-top: 16px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="difficulty" label="难度" width="100">
        <template #default="{ row }">
          <el-tag :type="diffTag(row.difficulty)" size="small">{{ row.difficulty }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="unlock_condition" label="解锁条件" />
      <el-table-column prop="gold_reward" label="金币奖励" width="100" />
      <el-table-column prop="boss_config" label="BOSS配置" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.boss_config" type="warning" size="small">BOSS</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delStage(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="isEdit ? '编辑关卡' : '新增关卡'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="难度">
          <el-select v-model="form.difficulty">
            <el-option label="简单" value="easy" />
            <el-option label="普通" value="normal" />
            <el-option label="困难" value="hard" />
            <el-option label="BOSS" value="boss" />
          </el-select>
        </el-form-item>
        <el-form-item label="解锁条件"><el-input v-model="form.unlock_condition" placeholder="前一关卡ID或角色ID" /></el-form-item>
        <el-form-item label="金币奖励"><el-input-number v-model="form.gold_reward" :min="0" /></el-form-item>
        <el-form-item label="关卡排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="BOSS配置(JSON)">
          <el-input v-model="form.boss_config" type="textarea" rows="4" placeholder="BOSS多阶段AI行为配置JSON" />
        </el-form-item>
        <el-form-item label="AI配置(JSON)">
          <el-input v-model="form.ai_config" type="textarea" rows="4" placeholder="普通AI行为配置JSON" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const stages = ref<any[]>([])
const loading = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const form = ref<any>({})

function diffTag(d: string) {
  switch (d) {
    case 'easy': return 'success'
    case 'normal': return 'primary'
    case 'hard': return 'danger'
    case 'boss': return 'warning'
    default: return 'info'
  }
}

async function fetchStages() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/pve/stages')
    stages.value = (res.data || res) || []
  } finally { loading.value = false }
}

function openAdd() {
  isEdit.value = false
  form.value = { name: '', difficulty: 'easy', unlock_condition: '', gold_reward: 50, sort_order: stages.value.length + 1, boss_config: '', ai_config: '' }
  visible.value = true
}

function openEdit(row: any) {
  isEdit.value = true
  form.value = { ...row }
  visible.value = true
}

async function save() {
  try {
    const url = isEdit.value ? `/admin/pve/stages/${form.value.id}` : '/admin/pve/stages'
    const method = isEdit.value ? api.put : api.post
    await method(url, form.value)
    ElMessage.success('保存成功')
    visible.value = false
    fetchStages()
  } catch (err) {
    ElMessage.error('操作失败')
  }
}

async function delStage(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除关卡 "${row.name}"？`, '警告', { type: 'warning' })
    await api.delete(`/admin/pve/stages/${row.id}`)
    ElMessage.success('删除成功')
    fetchStages()
  } catch (err) { /* cancelled */ }
}

onMounted(() => fetchStages())
</script>

<style scoped>
.stage-manage h2 { margin-bottom: 20px; }
</style>

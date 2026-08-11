<template>
  <div class="character-manage">
    <h2>角色配置</h2>
    <el-button type="primary" @click="openAdd">新增角色</el-button>
    <el-table :data="characters" border style="width: 100%; margin-top: 16px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="hp" label="血量" width="80" />
      <el-table-column prop="energy" label="能量" width="80" />
      <el-table-column prop="speed" label="速度" width="80" />
      <el-table-column prop="attack" label="攻击" width="80" />
      <el-table-column prop="defense" label="防御" width="80" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delCharacter(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="isEdit ? '编辑角色' : '新增角色'" width="600px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
        <el-row :gutter="12">
          <el-col :span="8"><el-form-item label="血量"><el-input-number v-model="form.hp" :min="1" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="能量"><el-input-number v-model="form.energy" :min="0" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="速度"><el-input-number v-model="form.speed" :min="1" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="8"><el-form-item label="攻击"><el-input-number v-model="form.attack" :min="1" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="防御"><el-input-number v-model="form.defense" :min="0" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="立绘路径"><el-input v-model="form.avatar_path" /></el-form-item>
        <el-form-item label="背景故事"><el-input v-model="form.story" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="save" :loading="saving">保存并刷新缓存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const characters = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const form = ref<any>({})

async function fetchCharacters() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/characters')
    characters.value = (res.data || res) || []
  } finally { loading.value = false }
}

function openAdd() {
  isEdit.value = false
  form.value = { name: '', description: '', hp: 100, energy: 100, speed: 5, attack: 20, defense: 10, avatar_path: '', story: '' }
  visible.value = true
}

function openEdit(row: any) {
  isEdit.value = true
  form.value = { ...row }
  visible.value = true
}

async function save() {
  saving.value = true
  try {
    const url = isEdit.value ? `/admin/characters/${form.value.id}` : '/admin/characters'
    const method = isEdit.value ? api.put : api.post
    await method(url, form.value)
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    visible.value = false
    fetchCharacters()
    await api.post('/admin/cache/refresh', null, { params: { type: 'characters' } })
  } catch (err) {
    ElMessage.error('操作失败')
  } finally { saving.value = false }
}

async function delCharacter(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除角色 "${row.name}"？`, '警告', { type: 'warning' })
    await api.delete(`/admin/characters/${row.id}`)
    ElMessage.success('删除成功')
    fetchCharacters()
  } catch (err) { /* cancelled */ }
}

onMounted(() => fetchCharacters())
</script>

<style scoped>
.character-manage h2 { margin-bottom: 20px; }
</style>

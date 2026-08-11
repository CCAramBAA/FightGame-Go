<template>
  <div class="skill-manage">
    <h2>技能配置</h2>
    <div class="toolbar">
      <el-select v-model="filterCharId" placeholder="筛选角色" clearable style="width: 200px" @change="fetchSkills">
        <el-option v-for="c in characters" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
      <el-button type="primary" @click="openAdd">新增技能</el-button>
    </div>
    <el-table :data="skills" border style="width: 100%; margin-top: 16px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="character_id" label="角色ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="priority" label="优先级" width="80" />
      <el-table-column prop="energy_cost" label="能耗" width="80" />
      <el-table-column prop="cooldown" label="冷却" width="80" />
      <el-table-column prop="damage" label="伤害" width="80" />
      <el-table-column prop="tags" label="标签" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delSkill(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="isEdit ? '编辑技能' : '新增技能'" width="600px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="角色">
          <el-select v-model="form.character_id">
            <el-option v-for="c in characters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
        <el-row :gutter="12">
          <el-col :span="8"><el-form-item label="优先级"><el-input-number v-model="form.priority" :min="0" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="能耗"><el-input-number v-model="form.energy_cost" :min="0" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="冷却(秒)"><el-input-number v-model="form.cooldown" :min="0" :step="0.1" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="8"><el-form-item label="伤害"><el-input-number v-model="form.damage" :min="0" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="范围"><el-input-number v-model="form.range" :min="0" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="标签"><el-input v-model="form.tags" placeholder="逗号分隔，如: stun,knockback" /></el-form-item>
        <el-form-item label="特效路径"><el-input v-model="form.effect_path" /></el-form-item>
        <el-form-item label="音效路径"><el-input v-model="form.sound_path" /></el-form-item>
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

const skills = ref<any[]>([])
const characters = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const filterCharId = ref('')
const form = ref<any>({})

async function fetchCharacters() {
  const res: any = await api.get('/admin/characters')
  characters.value = (res.data || res) || []
}

async function fetchSkills() {
  loading.value = true
  try {
    const params: any = {}
    if (filterCharId.value) params.character_id = filterCharId.value
    const res: any = await api.get('/admin/skills', { params })
    skills.value = (res.data || res) || []
  } finally { loading.value = false }
}

function openAdd() {
  isEdit.value = false
  form.value = { character_id: '', name: '', description: '', priority: 0, energy_cost: 10, cooldown: 1, damage: 10, range: 1, tags: '', effect_path: '', sound_path: '' }
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
    const url = isEdit.value ? `/admin/skills/${form.value.id}` : '/admin/skills'
    const method = isEdit.value ? api.put : api.post
    await method(url, form.value)
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    visible.value = false
    fetchSkills()
    await api.post('/admin/cache/refresh', null, { params: { type: 'skills' } })
  } catch (err) {
    ElMessage.error('操作失败')
  } finally { saving.value = false }
}

async function delSkill(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除技能 "${row.name}"？`, '警告', { type: 'warning' })
    await api.delete(`/admin/skills/${row.id}`)
    ElMessage.success('删除成功')
    fetchSkills()
  } catch (err) { /* cancelled */ }
}

onMounted(async () => {
  await fetchCharacters()
  fetchSkills()
})
</script>

<style scoped>
.skill-manage h2 { margin-bottom: 20px; }
.toolbar { display: flex; align-items: center; gap: 12px; }
</style>

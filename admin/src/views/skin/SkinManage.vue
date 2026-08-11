<template>
  <div class="skin-manage">
    <h2>皮肤管理</h2>
    <el-button type="primary" @click="openAdd">新增皮肤</el-button>
    <el-table :data="skins" border style="width: 100%; margin-top: 16px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="character_id" label="角色ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="preview_path" label="预览路径" />
      <el-table-column prop="price" label="价格" width="100" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delSkin(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="isEdit ? '编辑皮肤' : '新增皮肤'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="角色">
          <el-select v-model="form.character_id">
            <el-option v-for="c in characters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="预览路径"><el-input v-model="form.preview_path" /></el-form-item>
        <el-form-item label="立绘路径"><el-input v-model="form.avatar_path" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
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

const skins = ref<any[]>([])
const characters = ref<any[]>([])
const loading = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const form = ref<any>({})

async function fetchData() {
  loading.value = true
  try {
    const [skinRes, charRes]: any[] = await Promise.all([
      api.get('/admin/skins'),
      api.get('/admin/characters'),
    ])
    skins.value = (skinRes.data || skinRes) || []
    characters.value = (charRes.data || charRes) || []
  } finally { loading.value = false }
}

function openAdd() {
  isEdit.value = false
  form.value = { character_id: '', name: '', preview_path: '', avatar_path: '', description: '' }
  visible.value = true
}

function openEdit(row: any) {
  isEdit.value = true
  form.value = { ...row }
  visible.value = true
}

async function save() {
  try {
    const url = isEdit.value ? `/admin/skins/${form.value.id}` : '/admin/skins'
    const method = isEdit.value ? api.put : api.post
    await method(url, form.value)
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    visible.value = false
    fetchData()
  } catch (err) {
    ElMessage.error('操作失败')
  }
}

async function delSkin(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除皮肤 "${row.name}"？`, '警告', { type: 'warning' })
    await api.delete(`/admin/skins/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch (err) { /* cancelled */ }
}

onMounted(() => fetchData())
</script>

<style scoped>
.skin-manage h2 { margin-bottom: 20px; }
</style>

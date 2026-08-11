<template>
  <div class="shop-manage">
    <h2>商城配置</h2>
    <el-button type="primary" @click="openAdd">新增商品</el-button>
    <el-table :data="items" border style="width: 100%; margin-top: 16px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="商品名" />
      <el-table-column prop="item_type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.item_type === 'character' ? 'primary' : 'success'" size="small">{{ row.item_type === 'character' ? '角色' : '皮肤' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="item_id" label="关联ID" width="100" />
      <el-table-column prop="price" label="价格(金币)" width="120" />
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delItem(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="isEdit ? '编辑商品' : '新增商品'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="商品名"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.item_type">
            <el-option label="角色" value="character" />
            <el-option label="皮肤" value="skin" />
          </el-select>
        </el-form-item>
        <el-form-item label="关联ID"><el-input-number v-model="form.item_id" :min="1" /></el-form-item>
        <el-form-item label="价格(金币)"><el-input-number v-model="form.price" :min="1" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="图片路径"><el-input v-model="form.image_path" /></el-form-item>
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

const items = ref<any[]>([])
const loading = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const form = ref<any>({})

async function fetchItems() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/shop/items')
    items.value = (res.data || res) || []
  } finally { loading.value = false }
}

function openAdd() {
  isEdit.value = false
  form.value = { name: '', item_type: 'character', item_id: 1, price: 100, sort_order: 0, image_path: '', description: '' }
  visible.value = true
}

function openEdit(row: any) {
  isEdit.value = true
  form.value = { ...row }
  visible.value = true
}

async function save() {
  try {
    const url = isEdit.value ? `/admin/shop/items/${form.value.id}` : '/admin/shop/items'
    const method = isEdit.value ? api.put : api.post
    await method(url, form.value)
    ElMessage.success('保存成功')
    visible.value = false
    fetchItems()
  } catch (err) {
    ElMessage.error('操作失败')
  }
}

async function delItem(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除商品 "${row.name}"？`, '警告', { type: 'warning' })
    await api.delete(`/admin/shop/items/${row.id}`)
    ElMessage.success('删除成功')
    fetchItems()
  } catch (err) { /* cancelled */ }
}

onMounted(() => fetchItems())
</script>

<style scoped>
.shop-manage h2 { margin-bottom: 20px; }
</style>

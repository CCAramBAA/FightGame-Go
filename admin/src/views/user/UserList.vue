<template>
  <div class="user-list">
    <h2>用户管理</h2>
    <div class="search-bar">
      <el-input v-model="searchUsername" placeholder="搜索用户名" clearable style="width: 240px" @clear="fetchUsers" @keyup.enter="fetchUsers" />
      <el-button type="primary" @click="fetchUsers" style="margin-left: 10px">搜索</el-button>
    </div>
    <el-table :data="users" border style="width: 100%; margin-top: 16px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="gold" label="金币" width="100" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag size="small" :type="row.role === 'admin' ? 'danger' : 'info'">{{ row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag size="small" :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '正常' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="180" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="editUser(row)">编辑</el-button>
          <el-button size="small" :type="row.status === 1 ? 'danger' : 'success'" @click="toggleUser(row)">
            {{ row.status === 1 ? '封禁' : '解封' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      style="margin-top: 20px; justify-content: flex-end"
      @current-change="fetchUsers"
    />

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑用户" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称"><el-input v-model="editForm.nickname" /></el-form-item>
        <el-form-item label="角色">
          <el-select v-model="editForm.role">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="金币"><el-input-number v-model="editForm.gold" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUser">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const users = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchUsername = ref('')
const loading = ref(false)

const editVisible = ref(false)
const editForm = ref({ id: 0, nickname: '', role: '', gold: 0 })

async function fetchUsers() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/users', {
      params: { page: page.value, page_size: pageSize.value, username: searchUsername.value },
    })
    const d = res.data || res
    users.value = d.list || []
    total.value = d.total || 0
  } catch (err) {
    console.error('获取用户列表失败', err)
  } finally {
    loading.value = false
  }
}

function editUser(row: any) {
  editForm.value = { id: row.id, nickname: row.nickname, role: row.role, gold: row.gold }
  editVisible.value = true
}

async function saveUser() {
  try {
    await api.put(`/admin/users/${editForm.value.id}`, {
      nickname: editForm.value.nickname,
      role: editForm.value.role,
      gold: editForm.value.gold,
    })
    ElMessage.success('保存成功')
    editVisible.value = false
    fetchUsers()
  } catch (err) {
    ElMessage.error('保存失败')
  }
}

async function toggleUser(row: any) {
  const action = row.status === 1 ? '封禁' : '解封'
  try {
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示')
    await api.put(`/admin/users/${row.id}`, { status: row.status === 1 ? 0 : 1 })
    ElMessage.success(`${action}成功`)
    fetchUsers()
  } catch (err) {
    // cancelled
  }
}

onMounted(() => fetchUsers())
</script>

<style scoped>
.user-list h2 { margin-bottom: 20px; }
.search-bar { display: flex; align-items: center; }
</style>

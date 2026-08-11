<template>
  <div class="rule-manage">
    <h2>英雄特殊交互规则</h2>
    <el-button type="primary" @click="openAdd">新增规则</el-button>
    <el-alert type="info" :closable="false" style="margin-top: 10px">
      执行优先级：英雄特殊交互规则 > 全局技能优先级。同表格从上至下依次执行，互斥效果靠前规则优先生效。
    </el-alert>
    <el-table :data="rules" border style="width: 100%; margin-top: 16px" v-loading="loading" row-key="id">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="priority" label="排序" width="80" />
      <el-table-column prop="hero_a_id" label="英雄A ID" width="100" />
      <el-table-column prop="hero_b_id" label="英雄B ID" width="100" />
      <el-table-column prop="skill_a_id" label="技能A ID" width="100" />
      <el-table-column prop="skill_b_id" label="技能B ID" width="100" />
      <el-table-column prop="rule_type" label="规则类型" width="120" />
      <el-table-column prop="result_type" label="结果" width="100">
        <template #default="{ row }">
          <el-tag :type="row.result_type === 'override' ? 'danger' : 'info'" size="small">{{ row.result_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="delRule(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="isEdit ? '编辑规则' : '新增规则'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="排序编号"><el-input-number v-model="form.priority" :min="0" /></el-form-item>
        <el-form-item label="英雄A ID"><el-input-number v-model="form.hero_a_id" :min="0" /></el-form-item>
        <el-form-item label="英雄B ID"><el-input-number v-model="form.hero_b_id" :min="0" /></el-form-item>
        <el-form-item label="技能A ID"><el-input-number v-model="form.skill_a_id" :min="0" /></el-form-item>
        <el-form-item label="技能B ID"><el-input-number v-model="form.skill_b_id" :min="0" /></el-form-item>
        <el-form-item label="规则类型">
          <el-select v-model="form.rule_type">
            <el-option label="技能对抗" value="skill_clash" />
            <el-option label="特殊覆盖" value="special_override" />
            <el-option label="免疫" value="immune" />
          </el-select>
        </el-form-item>
        <el-form-item label="结果">
          <el-select v-model="form.result_type">
            <el-option label="覆盖" value="override" />
            <el-option label="抵消" value="negate" />
            <el-option label="保留" value="keep" />
          </el-select>
        </el-form-item>
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

const rules = ref<any[]>([])
const loading = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const form = ref<any>({})

async function fetchRules() {
  loading.value = true
  try {
    const res: any = await api.get('/admin/special-rules')
    rules.value = (res.data || res) || []
  } finally { loading.value = false }
}

function openAdd() {
  isEdit.value = false
  form.value = { priority: rules.value.length + 1, hero_a_id: 0, hero_b_id: 0, skill_a_id: 0, skill_b_id: 0, rule_type: 'skill_clash', result_type: 'override', description: '' }
  visible.value = true
}

function openEdit(row: any) {
  isEdit.value = true
  form.value = { ...row }
  visible.value = true
}

async function save() {
  try {
    const url = isEdit.value ? `/admin/special-rules/${form.value.id}` : '/admin/special-rules'
    const method = isEdit.value ? api.put : api.post
    await method(url, form.value)
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    visible.value = false
    fetchRules()
    api.post('/admin/cache/refresh', null, { params: { type: 'rules' } }).catch(() => {})
  } catch (err) {
    ElMessage.error('操作失败')
  }
}

async function delRule(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除此规则？`, '警告', { type: 'warning' })
    await api.delete(`/admin/special-rules/${row.id}`)
    ElMessage.success('删除成功')
    fetchRules()
  } catch (err) { /* cancelled */ }
}

onMounted(() => fetchRules())
</script>

<style scoped>
.rule-manage h2 { margin-bottom: 20px; }
</style>

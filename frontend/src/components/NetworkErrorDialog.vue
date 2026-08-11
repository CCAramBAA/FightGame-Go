<template>
  <ModalOverlay :visible="visible" title="网络异常" :showClose="false" maxWidth="420px">
    <p class="error-msg">{{ message }}</p>
    <template #footer>
      <button class="btn-secondary" @click="$emit('close')">退出程序</button>
      <button class="btn-primary" @click="$emit('retry')">重新连接</button>
    </template>
  </ModalOverlay>
</template>

<script setup lang="ts">
import ModalOverlay from './ModalOverlay.vue'

defineProps<{ visible: boolean; message?: string }>()
defineEmits<{ close: []; retry: [] }>()
</script>

<style scoped>
.error-msg {
  color: #ccc;
  text-align: center;
  font-size: 15px;
  margin: 8px 0;
}
.btn-primary {
  background: linear-gradient(135deg, #4fc3f7, #0288d1);
  border: none;
  color: #fff;
  padding: 10px 24px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
}
.btn-primary:hover { filter: brightness(1.1); }
.btn-secondary {
  background: rgba(255,255,255,0.1);
  border: 1px solid #555;
  color: #ccc;
  padding: 10px 24px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}
.btn-secondary:hover { background: rgba(255,255,255,0.2); }
</style>

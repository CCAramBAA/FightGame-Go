<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="onOverlayClick">
      <div class="modal-panel" :style="{ maxWidth: maxWidth }">
        <div class="modal-header" v-if="title || showClose">
          <h3>{{ title }}</h3>
          <button v-if="showClose" class="modal-close" @click="$emit('close')">×</button>
        </div>
        <div class="modal-body">
          <slot />
        </div>
        <div class="modal-footer" v-if="$slots.footer">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineProps<{
  visible: boolean
  title?: string
  showClose?: boolean
  maxWidth?: string
  closeOnOverlay?: boolean
}>()

const emit = defineEmits<{ close: [] }>()

function onOverlayClick() {
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}
.modal-panel {
  background: linear-gradient(135deg, #1a1a2e, #16213e);
  border: 2px solid #3a3a5e;
  border-radius: 16px;
  min-width: 360px;
  max-width: 90vw;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 0 40px rgba(0, 0, 0, 0.5);
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #2a2a4a;
}
.modal-header h3 {
  margin: 0;
  color: #e0e0e0;
  font-size: 18px;
}
.modal-close {
  background: none;
  border: none;
  color: #888;
  font-size: 24px;
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.modal-close:hover {
  background: rgba(255,255,255,0.1);
  color: #fff;
}
.modal-body {
  padding: 20px 24px;
}
.modal-footer {
  padding: 12px 24px 20px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
</style>

<template>
  <div class="art-placeholder" :style="placeholderStyle">
    <span class="art-label">{{ label }}</span>
    <span v-if="sub" class="art-sub">{{ sub }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  label: string
  sub?: string
  width?: string | number
  height?: string | number
  bgColor?: string
}>()

const placeholderStyle = computed(() => ({
  width: typeof props.width === 'number' ? `${props.width}px` : (props.width || '100%'),
  height: typeof props.height === 'number' ? `${props.height}px` : (props.height || '100%'),
  backgroundColor: props.bgColor || '#1a1a2e',
}))
</script>

<style scoped>
.art-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 2px dashed #444;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}
.art-placeholder::before {
  content: '';
  position: absolute;
  inset: 0;
  background: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 10px,
    rgba(255,255,255,0.02) 10px,
    rgba(255,255,255,0.02) 20px
  );
}
.art-label {
  color: #888;
  font-size: 14px;
  font-weight: bold;
  z-index: 1;
}
.art-sub {
  color: #666;
  font-size: 11px;
  margin-top: 4px;
  z-index: 1;
}
</style>

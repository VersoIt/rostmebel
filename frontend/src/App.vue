<script setup lang="ts">
import Toast from '@/components/common/Toast.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
</script>

<template>
  <Toast />
  <ConfirmDialog />
  <router-view />
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>

<script setup lang="ts">
import { computed } from 'vue';
import { useNotificationStore } from '@/stores/notifications';
import { LucideAlertCircle, LucideCheckCircle, LucideInfo, LucideX } from 'lucide-vue-next';

const store = useNotificationStore();

const toastClass = computed(() => {
  if (store.type === 'error') return 'bg-red-50/95 border-red-100 text-red-800';
  if (store.type === 'success') return 'bg-green-50/95 border-green-100 text-green-800';
  return 'bg-brand-cream/95 border-brand-brown/10 text-brand-brown';
});

const toastIcon = computed(() => {
  if (store.type === 'error') return LucideAlertCircle;
  if (store.type === 'success') return LucideCheckCircle;
  return LucideInfo;
});
</script>

<template>
  <transition 
    enter-active-class="transform transition duration-300 ease-out"
    enter-from-class="translate-y-[-100%] opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transform transition duration-200 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-[-100%] opacity-0"
  >
    <div v-if="store.visible" class="fixed top-6 left-1/2 -translate-x-1/2 z-[300] w-full max-w-md px-4">
      <div :class="[
        'flex items-center gap-4 p-4 rounded-lg shadow-2xl border backdrop-blur-md',
        toastClass
      ]">
        <component :is="toastIcon" :size="24" />
        <p class="flex-1 font-medium text-sm">{{ store.message }}</p>
        <button @click="store.visible = false" class="hover:opacity-50 transition-opacity">
          <LucideX :size="20" />
        </button>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { useNotificationStore } from '@/stores/notifications';
import { LucideAlertCircle, LucideCheckCircle, LucideX } from 'lucide-vue-next';

const store = useNotificationStore();
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
    <div v-if="store.visible" class="fixed top-6 left-1/2 -translate-x-1/2 z-[100] w-full max-w-md px-4">
      <div :class="[
        'flex items-center gap-4 p-4 rounded-2xl shadow-2xl border backdrop-blur-md',
        store.type === 'error' ? 'bg-red-50/90 border-red-100 text-red-800' : 'bg-green-50/90 border-green-100 text-green-800'
      ]">
        <component :is="store.type === 'error' ? LucideAlertCircle : LucideCheckCircle" :size="24" />
        <p class="flex-1 font-medium text-sm">{{ store.message }}</p>
        <button @click="store.visible = false" class="hover:opacity-50 transition-opacity">
          <LucideX :size="20" />
        </button>
      </div>
    </div>
  </transition>
</template>

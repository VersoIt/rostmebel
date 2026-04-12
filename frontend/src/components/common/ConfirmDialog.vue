<script setup lang="ts">
import { useConfirmStore } from '@/stores/confirm';
import { LucideAlertTriangle } from 'lucide-vue-next';

const store = useConfirmStore();
</script>

<template>
  <Teleport to="body">
    <transition name="confirm-fade">
      <div v-if="store.visible" class="fixed inset-0 z-[220] flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="store.resolve(false)"></div>
        <section class="relative w-full max-w-md rounded-lg border border-brand-brown/10 bg-white p-6 shadow-2xl">
          <div class="mb-5 flex items-start gap-4">
            <div
              :class="[
                'flex h-11 w-11 shrink-0 items-center justify-center rounded-lg',
                store.tone === 'danger' ? 'bg-red-50 text-red-600' : 'bg-brand-cream text-brand-gold'
              ]"
            >
              <LucideAlertTriangle :size="22" />
            </div>
            <div>
              <h2 class="font-serif text-2xl text-brand-brown">{{ store.title }}</h2>
              <p class="mt-2 leading-7 text-brand-brown/65">{{ store.message }}</p>
            </div>
          </div>

          <div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
            <button
              class="rounded-lg border border-brand-brown/10 px-5 py-3 font-semibold text-brand-brown/70 transition hover:bg-brand-gray"
              @click="store.resolve(false)"
            >
              {{ store.cancelLabel }}
            </button>
            <button
              :class="[
                'rounded-lg px-5 py-3 font-semibold text-white transition',
                store.tone === 'danger' ? 'bg-red-600 hover:bg-red-700' : 'bg-brand-brown hover:bg-brand-gold'
              ]"
              @click="store.resolve(true)"
            >
              {{ store.confirmLabel }}
            </button>
          </div>
        </section>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
.confirm-fade-enter-active,
.confirm-fade-leave-active {
  transition: opacity 0.2s ease;
}

.confirm-fade-enter-from,
.confirm-fade-leave-to {
  opacity: 0;
}
</style>

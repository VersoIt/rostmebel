<script setup lang="ts">
import { useConfirmStore } from '@/stores/confirm';
import { LucideAlertTriangle } from 'lucide-vue-next';

const store = useConfirmStore();
</script>

<template>
  <Teleport to="body">
    <transition name="confirm-fade">
      <div v-if="store.visible" class="ui-modal-backdrop">
        <div class="absolute inset-0" @click="store.resolve(false)"></div>
        <section class="ui-modal-panel max-w-md p-6">
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
              class="ui-button ui-button-secondary"
              @click="store.resolve(false)"
            >
              {{ store.cancelLabel }}
            </button>
            <button
              :class="[
                'ui-button',
                store.tone === 'danger' ? 'ui-button-danger' : 'ui-button-primary'
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
  transition: opacity 0.18s ease;
}

.confirm-fade-enter-from,
.confirm-fade-leave-to {
  opacity: 0;
}
</style>

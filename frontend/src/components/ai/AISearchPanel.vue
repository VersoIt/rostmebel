<script setup lang="ts">
import { ref } from 'vue';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideLoader2, LucideSearch } from 'lucide-vue-next';
import type { Product } from '@/types';
import { getApiErrorMessage } from '@/api/errors';
import { useNotificationStore } from '@/stores/notifications';

const productStore = useProductStore();
const notificationStore = useNotificationStore();
const query = ref('');
const results = ref<Product[]>([]);
const isSearching = ref(false);
const hasSearched = ref(false);
const searchFailed = ref(false);

const handleSearch = async () => {
  if (!query.value.trim() || isSearching.value) return;

  isSearching.value = true;
  hasSearched.value = true;
  searchFailed.value = false;

  try {
    results.value = await productStore.aiSearch(query.value);
    if (results.value.length === 0) {
      notificationStore.show('Подходящих проектов пока нет. Попробуйте описать задачу иначе.', 'info');
    }
  } catch (err) {
    searchFailed.value = true;
    notificationStore.show(getApiErrorMessage(err), 'error');
  } finally {
    isSearching.value = false;
  }

  if (results.value.length > 0) {
    window.setTimeout(() => {
      document.getElementById('ai-results')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, 100);
  }
};
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <div class="ui-card p-3 sm:p-4">
      <div class="relative flex flex-col gap-3 sm:block">
        <input
          v-model="query"
          type="text"
          placeholder="Например: светлая кухня с техникой до 250 000 ₽"
          class="ui-input min-h-14 pl-12 pr-4 text-base sm:pr-32"
          @keyup.enter="handleSearch"
        >
        <LucideSearch class="absolute left-4 top-7 -translate-y-1/2 text-brand-brown/35 sm:top-1/2" :size="23" />
        <button
          type="button"
          class="ui-button ui-button-primary sm:absolute sm:bottom-2 sm:right-2 sm:top-2"
          :disabled="isSearching"
          @click="handleSearch"
        >
          <LucideLoader2 v-if="isSearching" class="animate-spin" :size="19" />
          {{ isSearching ? 'Ищем...' : 'Найти' }}
        </button>
      </div>
    </div>

    <div v-if="results.length > 0" id="ai-results" class="mt-12 motion-fade-up">
      <div class="mb-7 flex items-center justify-center gap-3">
        <div class="h-px flex-1 bg-brand-gold/20"></div>
        <span class="rounded-lg bg-brand-gold px-3 py-1 text-[10px] font-black uppercase tracking-widest text-white">Подбор по описанию</span>
        <div class="h-px flex-1 bg-brand-gold/20"></div>
      </div>
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <ProductCard v-for="product in results" :key="product.id" :product="product" />
      </div>
    </div>

    <div v-else-if="hasSearched && !searchFailed && !isSearching && query" class="ui-empty mt-8">
      Проектов с таким описанием пока нет в портфолио.
    </div>
  </div>
</template>

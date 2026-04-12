<script setup lang="ts">
import { ref } from 'vue';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideSearch, LucideLoader2 } from 'lucide-vue-next';
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
      notificationStore.show('Мы не нашли подходящих проектов по вашему описанию. Попробуйте изменить запрос.', 'info');
    }
  } catch (err) {
    searchFailed.value = true;
    notificationStore.show(getApiErrorMessage(err), 'error');
  } finally {
    isSearching.value = false;
  }
  
  if (results.value.length > 0) {
    setTimeout(() => {
      document.getElementById('ai-results')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, 100);
  }
};
</script>

<template>
  <div class="max-w-4xl mx-auto px-4">
    <div class="relative group flex flex-col gap-3 sm:block">
      <input 
        v-model="query"
        type="text"
        placeholder="Например: светлая кухня в стиле лофт до 150 000 руб..."
        class="w-full pl-14 pr-6 py-5 rounded-lg border-2 border-brand-brown/5 focus:border-brand-gold outline-none text-lg transition-all shadow-lg group-hover:shadow-xl sm:pr-32 sm:py-6 sm:text-xl"
        @keyup.enter="handleSearch"
      >
      <LucideSearch class="absolute left-6 top-7 -translate-y-1/2 text-brand-brown/40 sm:top-1/2" :size="28" />
      <button 
        @click="handleSearch"
        :disabled="isSearching"
        class="flex items-center justify-center gap-2 rounded-lg bg-brand-brown px-8 py-4 font-medium text-white transition-colors hover:bg-brand-gold disabled:cursor-not-allowed disabled:opacity-50 sm:absolute sm:bottom-3 sm:right-3 sm:top-3 sm:py-0"
      >
        <LucideLoader2 v-if="isSearching" class="animate-spin" :size="20" />
        {{ isSearching ? 'Ищу проекты...' : 'Найти' }}
      </button>
    </div>

    <div v-if="results.length > 0" id="ai-results" class="mt-20">
      <div class="flex items-center justify-center gap-2 mb-8">
        <div class="h-px bg-brand-gold/20 flex-1"></div>
        <span class="bg-brand-gold text-white text-[10px] font-bold px-3 py-1 rounded-lg uppercase">Персональный подбор ИИ</span>
        <div class="h-px bg-brand-gold/20 flex-1"></div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <ProductCard v-for="p in results" :key="p.id" :product="p" />
      </div>
    </div>
    
    <div v-else-if="hasSearched && !searchFailed && !isSearching && query && results.length === 0" class="mt-8 text-center text-brand-brown/40 italic">
      Проектов с таким описанием пока нет в нашем портфолио
    </div>
  </div>
</template>

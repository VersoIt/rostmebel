<script setup lang="ts">
import { ref } from 'vue';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideSearch, LucideLoader2 } from 'lucide-vue-next';
import type { Product } from '@/types';

import { useNotificationStore } from '@/stores/notifications';

const productStore = useProductStore();
const notificationStore = useNotificationStore();
const query = ref('');
const results = ref<Product[]>([]);
const isSearching = ref(false);

const handleSearch = async () => {
  if (!query.value.trim() || isSearching.value) return;
  
  isSearching.value = true;
  try {
    results.value = await productStore.aiSearch(query.value);
    if (results.value.length === 0) {
      notificationStore.show('Gemma 4 не нашла подходящих товаров. Попробуйте другой запрос.', 'info');
    }
  } catch (err) {
    notificationStore.show('Ошибка ИИ-агента. Повторите попытку позже.', 'error');
  } finally {
    isSearching.value = false;
  }
  
  if (results.value.length > 0) {
    document.getElementById('ai-results')?.scrollIntoView({ behavior: 'smooth' });
  }
};
</script>

<template>
  <div class="max-w-4xl mx-auto px-4">
    <div class="relative group">
      <input 
        v-model="query"
        type="text"
        placeholder="Например: хочу светлую скандинавскую спальню до 50000 рублей..."
        class="w-full pl-14 pr-32 py-6 rounded-2xl border-2 border-brand-brown/5 focus:border-brand-gold outline-none text-xl transition-all shadow-lg group-hover:shadow-xl"
        @keyup.enter="handleSearch"
      >
      <LucideSearch class="absolute left-6 top-1/2 -translate-y-1/2 text-brand-brown/40" :size="28" />
      <button 
        @click="handleSearch"
        :disabled="isSearching"
        class="absolute right-3 top-3 bottom-3 bg-brand-brown text-white px-8 rounded-xl font-medium hover:bg-brand-gold disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
      >
        <LucideLoader2 v-if="isSearching" class="animate-spin" :size="20" />
        {{ isSearching ? 'Думаю...' : 'Найти' }}
      </button>
    </div>

    <div v-if="results.length > 0" id="ai-results" class="mt-20">
      <div class="flex items-center justify-center gap-2 mb-8">
        <div class="h-px bg-brand-gold/20 flex-1"></div>
        <span class="bg-brand-gold text-white text-[10px] font-bold px-3 py-1 rounded-full uppercase tracking-widest">Персональный подбор Gemma 4</span>
        <div class="h-px bg-brand-gold/20 flex-1"></div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <ProductCard v-for="p in results" :key="p.id" :product="p" />
      </div>
    </div>
    
    <div v-else-if="!isSearching && query && results.length === 0" class="mt-8 text-center text-brand-brown/40 italic">
      Ничего не нашлось, попробуйте изменить запрос
    </div>
  </div>
</template>

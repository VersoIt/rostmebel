<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { 
  LucidePackage, 
  LucideClipboardList, 
} from 'lucide-vue-next';
import api from '@/api/client';
import { useProductStore } from '@/stores/products';
import ProductModal from '@/components/admin/ProductModal.vue';
import type { Stats, Product } from '@/types';

const stats = ref<Stats | null>(null);
const productStore = useProductStore();
const isModalOpen = ref(false);
const editingProduct = ref<Product | null>(null);

onMounted(async () => {
  try {
    const { data } = await api.get('/admin/stats');
    stats.value = data;
    await productStore.fetchCategories();
  } catch (err) {
    console.error(err);
  }
});

const openEdit = async (id: number) => {
  const p = await productStore.fetchProduct(id);
  if (p) {
    editingProduct.value = p;
    isModalOpen.value = true;
  }
};

const handleSaved = async () => {
  isModalOpen.value = false;
  const { data } = await api.get('/admin/stats');
  stats.value = data;
};

// SVG Chart logic
const getChartPath = () => {
  if (!stats.value?.orders_by_day || stats.value.orders_by_day.length < 2) return '';
  const data = stats.value.orders_by_day.map(d => d.count);
  const max = Math.max(...data, 5);
  const width = 1000;
  const height = 200;
  const step = width / (data.length - 1);
  
  return data.map((v, i) => `${i === 0 ? 'M' : 'L'} ${i * step} ${height - (v / max * height)}`).join(' ');
};
</script>

<template>
  <div>
    <h1 class="font-serif text-4xl text-brand-brown mb-12">Обзор системы</h1>
    
    <div v-if="stats" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-12">
      <div class="bg-white p-8 rounded-3xl border border-brand-brown/5 shadow-sm">
        <span class="text-brand-brown/40 uppercase tracking-widest text-xs font-semibold mb-2 block">Всего проектов</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.products_count }}</div>
      </div>
      <div class="bg-brand-gold p-8 rounded-3xl shadow-xl shadow-brand-gold/20">
        <span class="text-brand-brown/60 uppercase tracking-widest text-xs font-semibold mb-2 block">Заявки сегодня</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.new_orders_today }}</div>
      </div>
      <div class="bg-white p-8 rounded-3xl border border-brand-brown/5 shadow-sm">
        <span class="text-brand-brown/40 uppercase tracking-widest text-xs font-semibold mb-2 block">Всего заявок</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.total_orders }}</div>
      </div>
      <div class="bg-brand-brown p-8 rounded-3xl shadow-xl">
        <span class="text-white/40 uppercase tracking-widest text-xs font-semibold mb-2 block">Успешных сделок</span>
        <div class="text-4xl font-serif text-brand-gold">{{ stats.success_rate.toFixed(0) }}%</div>
      </div>
    </div>

    <!-- Chart -->
    <div v-if="stats?.orders_by_day && stats.orders_by_day.length > 0" class="bg-white p-8 rounded-3xl border border-brand-brown/5 shadow-sm mb-12">
      <h3 class="font-serif text-xl mb-8">Активность заявок (30 дней)</h3>
      <div class="h-48 w-full">
        <svg viewBox="0 0 1000 200" class="w-full h-full preserve-aspect-ratio-none">
          <path 
            :d="getChartPath()" 
            fill="none" 
            stroke="#C9A84C" 
            stroke-width="4" 
            stroke-linecap="round" 
            stroke-linejoin="round"
            class="drop-shadow-lg"
          />
        </svg>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <div class="bg-white rounded-3xl border border-brand-brown/5 shadow-sm overflow-hidden">
        <div class="p-8 border-b border-brand-brown/5 bg-brand-gray/20">
          <h3 class="font-serif text-xl">Популярные проекты</h3>
        </div>
        <div class="p-8 space-y-6">
          <div 
            v-for="p in stats?.top_products" 
            :key="p.id" 
            @click="openEdit(p.id)"
            class="flex items-center justify-between group cursor-pointer hover:bg-brand-gray/5 p-2 rounded-xl transition-all"
          >
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 bg-brand-gray rounded-xl flex items-center justify-center font-bold text-brand-brown/30">{{ p.id }}</div>
              <span class="font-medium group-hover:text-brand-gold transition-colors">{{ p.name }}</span>
            </div>
            <span class="bg-brand-gray px-4 py-2 rounded-lg text-sm font-semibold">{{ p.count }} заявок</span>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-8">
        <div class="bg-white p-8 rounded-3xl border border-brand-brown/5 shadow-sm">
          <h3 class="font-serif text-xl mb-6">Быстрые действия</h3>
          <div class="grid grid-cols-2 gap-4">
            <router-link to="/admin/products" class="bg-brand-brown text-white p-6 rounded-2xl hover:bg-brand-gold transition-all text-left group">
              <LucidePackage class="mb-4 group-hover:scale-110 transition-transform" />
              <span class="font-medium">Управление проектами</span>
            </router-link>
            <router-link to="/admin/orders" class="bg-brand-gray p-6 rounded-2xl hover:bg-brand-brown hover:text-white transition-all text-left group">
              <LucideClipboardList class="mb-4 group-hover:scale-110 transition-transform text-brand-gold" />
              <span class="font-medium">Проверить заявки</span>
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <ProductModal 
      v-if="isModalOpen" 
      :product="editingProduct" 
      :categories="productStore.categories"
      @close="isModalOpen = false" 
      @saved="handleSaved" 
    />
  </div>
</template>

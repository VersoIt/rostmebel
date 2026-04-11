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

// Advanced SVG Chart logic
const chartWidth = 1000;
const chartHeight = 250;
const padding = { top: 20, right: 20, bottom: 40, left: 40 };

const getChartPoints = () => {
  if (!stats.value?.orders_by_day || stats.value.orders_by_day.length < 2) return [];
  const data = stats.value.orders_by_day.map(d => d.count);
  const max = Math.max(...data, 5);
  const stepX = (chartWidth - padding.left - padding.right) / (data.length - 1);
  const drawHeight = chartHeight - padding.top - padding.bottom;
  
  return data.map((v, i) => ({
    x: padding.left + i * stepX,
    y: padding.top + drawHeight - (v / max * drawHeight),
    val: v,
    date: stats.value!.orders_by_day[i].date
  }));
};

const getLinePath = () => {
  const points = getChartPoints();
  if (points.length === 0) return '';
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ');
};

const getAreaPath = () => {
  const points = getChartPoints();
  if (points.length === 0) return '';
  const path = getLinePath();
  const lastPoint = points[points.length - 1];
  const firstPoint = points[0];
  const bottom = chartHeight - padding.bottom;
  return `${path} L ${lastPoint.x} ${bottom} L ${firstPoint.x} ${bottom} Z`;
};
</script>

<template>
  <div>
    <h1 class="font-serif text-4xl text-brand-brown mb-12">Обзор системы</h1>
    
    <!-- Top Stats -->
    <div v-if="stats" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-12">
      <div class="bg-white p-8 rounded-3xl border border-brand-brown/5 shadow-sm">
        <span class="text-brand-brown/40 uppercase tracking-widest text-xs font-semibold mb-2 block">Всего проектов</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.projects_count }}</div>
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

    <!-- Chart Section -->
    <div v-if="stats?.orders_by_day && stats.orders_by_day.length > 0" class="bg-white p-10 rounded-[2.5rem] border border-brand-brown/5 shadow-sm mb-12">
      <div class="flex items-center justify-between mb-10">
        <div>
          <h3 class="font-serif text-2xl text-brand-brown">Активность заявок</h3>
          <p class="text-brand-brown/40 text-sm font-medium">Статистика за последние 30 дней</p>
        </div>
        <div class="flex gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-brand-gold rounded-full"></div>
            <span class="text-xs font-bold uppercase tracking-widest text-brand-brown/60">Новые заявки</span>
          </div>
        </div>
      </div>

      <div class="relative w-full overflow-hidden">
        <svg :viewBox="`0 0 ${chartWidth} ${chartHeight}`" class="w-full h-auto overflow-visible">
          <!-- Definitions -->
          <defs>
            <linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#C9A84C" stop-opacity="0.3" />
              <stop offset="100%" stop-color="#C9A84C" stop-opacity="0" />
            </linearGradient>
          </defs>

          <!-- Horizontal Grid Lines -->
          <g class="grid-lines">
            <line v-for="i in 5" :key="i"
              :x1="padding.left" 
              :y1="padding.top + (chartHeight - padding.top - padding.bottom) * (i-1) / 4"
              :x2="chartWidth - padding.right"
              :y2="padding.top + (chartHeight - padding.top - padding.bottom) * (i-1) / 4"
              stroke="#F5F5F0"
              stroke-width="1"
            />
          </g>

          <!-- Area -->
          <path :d="getAreaPath()" fill="url(#chartGradient)" />

          <!-- Line -->
          <path 
            :d="getLinePath()" 
            fill="none" 
            stroke="#C9A84C" 
            stroke-width="4" 
            stroke-linecap="round" 
            stroke-linejoin="round"
            class="drop-shadow-lg"
          />

          <!-- Points and Tooltips -->
          <g v-for="(p, i) in getChartPoints()" :key="i" class="chart-point group/point">
            <!-- Invisible hit area -->
            <circle :cx="p.x" :cy="p.y" r="15" fill="transparent" class="cursor-pointer" />
            <!-- Visual point -->
            <circle 
              :cx="p.x" :cy="p.y" r="5" 
              fill="white" 
              stroke="#C9A84C" 
              stroke-width="3"
              class="transition-all duration-300 group-hover/point:r-8 group-hover/point:stroke-brand-brown"
            />
            
            <!-- Value on hover -->
            <g class="opacity-0 group-hover/point:opacity-100 transition-opacity duration-300 pointer-events-none">
              <rect :x="p.x - 20" :y="p.y - 45" width="40" height="30" rx="8" fill="#2C1810" />
              <text :x="p.x" :y="p.y - 25" text-anchor="middle" fill="white" font-size="14" font-weight="bold">{{ p.val }}</text>
            </g>

            <!-- Date Labels (show every 5th to avoid mess) -->
            <text 
              v-if="i % 5 === 0 || i === getChartPoints().length - 1"
              :x="p.x" 
              :y="chartHeight - 10" 
              text-anchor="middle" 
              fill="#2C1810" 
              opacity="0.3" 
              font-size="10" 
              font-weight="bold"
              class="uppercase tracking-widest"
            >
              {{ p.date.split('-').slice(2).join('.') }}
            </text>
          </g>
        </svg>
      </div>
    </div>

    <!-- Bottom Lists -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <div class="bg-white rounded-3xl border border-brand-brown/5 shadow-sm overflow-hidden">
        <div class="p-8 border-b border-brand-brown/5 bg-brand-gray/20">
          <h3 class="font-serif text-xl">Популярные проекты</h3>
        </div>
        <div class="p-8 space-y-6">
          <div 
            v-for="p in stats?.top_projects" 
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
            <router-link to="/admin/projects" class="bg-brand-brown text-white p-6 rounded-2xl hover:bg-brand-gold transition-all text-left group">
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

<style scoped>
.preserve-aspect-ratio-none {
  preserveAspectRatio: none;
}
.chart-point circle {
  transition: all 0.3s ease;
}
</style>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { 
  LucidePackage, 
  LucideClipboardList,
  LucideMessageSquare,
  LucideChevronRight,
  LucideClock
} from 'lucide-vue-next';
import api from '@/api/client';
import { useProductStore } from '@/stores/products';
import ProductModal from '@/components/admin/ProductModal.vue';
import type { Stats, Product } from '@/types';

const stats = ref<Stats | null>(null);
const productStore = useProductStore();
const isModalOpen = ref(false);
const editingProduct = ref<Product | null>(null);

const fetchStats = async () => {
  try {
    const { data } = await api.get('/admin/stats');
    stats.value = data;
  } catch (err) {
    console.error(err);
  }
};

onMounted(async () => {
  await fetchStats();
  await productStore.fetchCategories();
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
  await fetchStats();
};

// Advanced SVG Chart logic
const chartWidth = 1000;
const chartHeight = 300;
const padding = { top: 60, right: 20, bottom: 40, left: 40 };

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
  <div class="space-y-12">
    <div class="flex items-center justify-between">
      <h1 class="font-serif text-4xl text-brand-brown">Обзор системы</h1>
      <div class="text-sm font-bold text-brand-brown/40 uppercase tracking-widest">
        {{ new Date().toLocaleDateString('ru-RU', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
      </div>
    </div>
    
    <!-- Top Stats -->
    <div v-if="stats" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
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
    <div v-if="stats?.orders_by_day && stats.orders_by_day.length > 0" class="bg-white p-10 rounded-[2.5rem] border border-brand-brown/5 shadow-sm">
      <div class="flex items-center justify-between mb-10">
        <div>
          <h3 class="font-serif text-2xl text-brand-brown">Активность заявок</h3>
          <p class="text-brand-brown/40 text-sm font-medium">Динамика за последние 30 дней</p>
        </div>
        <div class="flex items-center gap-2">
          <div class="w-3 h-3 bg-brand-gold rounded-full"></div>
          <span class="text-xs font-bold uppercase tracking-widest text-brand-brown/60">Новые заявки</span>
        </div>
      </div>

      <div class="relative w-full">
        <svg :viewBox="`0 0 ${chartWidth} ${chartHeight}`" class="w-full h-auto overflow-visible">
          <defs>
            <linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#C9A84C" stop-opacity="0.3" />
              <stop offset="100%" stop-color="#C9A84C" stop-opacity="0" />
            </linearGradient>
          </defs>
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
          <path :d="getAreaPath()" fill="url(#chartGradient)" />
          <path :d="getLinePath()" fill="none" stroke="#C9A84C" stroke-width="4" stroke-linecap="round" stroke-linejoin="round" class="drop-shadow-lg" />
          <g v-for="(p, i) in getChartPoints()" :key="i" class="chart-point group/point">
            <circle :cx="p.x" :cy="p.y" r="15" fill="transparent" class="cursor-pointer" />
            <circle :cx="p.x" :cy="p.y" r="5" fill="white" stroke="#C9A84C" stroke-width="3" class="transition-all duration-300 group-hover/point:r-8 group-hover/point:stroke-brand-brown" />
            <g class="opacity-0 group-hover/point:opacity-100 transition-opacity duration-300 pointer-events-none">
              <rect :x="p.x - 20" :y="p.y < 60 ? p.y + 15 : p.y - 45" width="40" height="30" rx="8" fill="#2C1810" />
              <text :x="p.x" :y="p.y < 60 ? p.y + 35 : p.y - 25" text-anchor="middle" fill="white" font-size="14" font-weight="bold">{{ p.val }}</text>
            </g>
            <text v-if="i % 5 === 0 || i === getChartPoints().length - 1" :x="p.x" :y="chartHeight - 10" text-anchor="middle" fill="#2C1810" opacity="0.3" font-size="10" font-weight="bold" class="uppercase tracking-widest">
              {{ p.date.split('-').slice(2).join('.') }}
            </text>
          </g>
        </svg>
      </div>
    </div>

    <!-- MAIN OPERATIONAL SECTION -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      
      <!-- 1. Latest Leads (Left - 8 columns) -->
      <div class="lg:col-span-8 bg-white rounded-[2.5rem] border border-brand-brown/5 shadow-sm overflow-hidden flex flex-col">
        <div class="p-8 border-b border-brand-brown/5 bg-brand-gray/10 flex items-center justify-between">
          <div>
            <h3 class="font-serif text-2xl">Свежие заявки</h3>
            <p class="text-brand-brown/40 text-sm font-medium">Требуют первичной обработки</p>
          </div>
          <router-link to="/admin/orders" class="text-brand-gold font-bold text-xs uppercase tracking-widest hover:underline flex items-center gap-2">
            Все заявки <LucideChevronRight :size="16" />
          </router-link>
        </div>
        
        <div class="flex-1">
          <div v-if="stats?.recent_orders.length" class="divide-y divide-brand-brown/5">
            <div 
              v-for="order in stats.recent_orders" :key="order.id"
              class="p-6 hover:bg-brand-gray/5 transition-colors flex items-center justify-between group"
            >
              <div class="flex items-center gap-6">
                <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-2xl flex items-center justify-center">
                  <LucideClock :size="24" />
                </div>
                <div>
                  <div class="font-bold text-brand-brown">{{ order.client_name }}</div>
                  <div class="text-xs text-brand-gold font-bold uppercase tracking-widest">{{ order.project_name }}</div>
                </div>
              </div>
              <div class="flex items-center gap-8">
                <div class="text-right hidden md:block">
                  <div class="text-sm font-medium text-brand-brown/60">{{ new Date(order.created_at).toLocaleDateString() }}</div>
                  <div class="text-[10px] text-brand-brown/30 font-bold uppercase">{{ new Date(order.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</div>
                </div>
                <router-link to="/admin/orders" class="p-3 bg-brand-brown text-white rounded-xl opacity-0 group-hover:opacity-100 transition-all hover:bg-brand-gold">
                  <LucideChevronRight :size="20" />
                </router-link>
              </div>
            </div>
          </div>
          <div v-else class="py-20 text-center text-brand-brown/20 italic">
            Новых заявок пока нет
          </div>
        </div>
      </div>

      <!-- 2. Priority Tasks (Right - 4 columns) -->
      <div class="lg:col-span-4 space-y-8">
        <!-- Review Task -->
        <div 
          :class="[
            'p-8 rounded-[2.5rem] border transition-all duration-500 flex flex-col justify-between h-full',
            stats?.pending_reviews_count ? 'bg-brand-gold border-brand-gold text-brand-brown shadow-xl shadow-brand-gold/20' : 'bg-white border-brand-brown/5 text-brand-brown opacity-60'
          ]"
        >
          <div>
            <LucideMessageSquare class="mb-6" :size="40" />
            <h3 class="font-serif text-3xl mb-4 leading-tight">Модерация отзывов</h3>
            <p v-if="stats?.pending_reviews_count" class="font-bold mb-8">
              У вас {{ stats.pending_reviews_count }} новых отзывов, ждущих публикации на сайте.
            </p>
            <p v-else class="font-medium opacity-60 mb-8">
              Все отзывы проверены. Новых публикаций нет.
            </p>
          </div>
          <router-link 
            to="/admin/reviews" 
            :class="[
              'w-full py-4 rounded-2xl font-black uppercase tracking-widest text-center transition-all',
              stats?.pending_reviews_count ? 'bg-brand-brown text-white hover:bg-black shadow-lg' : 'bg-brand-gray text-brand-brown/40 pointer-events-none'
            ]"
          >
            Проверить
          </router-link>
        </div>

        <!-- Popular Project Quick Look -->
        <div class="bg-white p-8 rounded-[2.5rem] border border-brand-brown/5 shadow-sm">
          <div class="flex items-center gap-4 mb-6">
            <LucidePackage class="text-brand-gold" :size="24" />
            <h4 class="font-bold text-brand-brown">ТОП проект</h4>
          </div>
          <div v-if="stats?.top_projects.length" class="space-y-2">
            <div class="text-xl font-serif text-brand-brown line-clamp-1">{{ stats.top_projects[0].name }}</div>
            <div class="text-brand-gold font-bold text-sm">{{ stats.top_projects[0].count }} заявок</div>
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
.chart-point circle {
  transition: all 0.3s ease;
}
</style>

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
import { useNotificationStore } from '@/stores/notifications';
import { getApiErrorMessage } from '@/api/errors';

const stats = ref<Stats | null>(null);
const productStore = useProductStore();
const notificationStore = useNotificationStore();
const isModalOpen = ref(false);
const editingProduct = ref<Product | null>(null);

const fetchStats = async () => {
  try {
    const { data } = await api.get('/admin/stats');
    stats.value = data;
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
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
  <div class="space-y-6 lg:space-y-8">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-xs font-black uppercase tracking-widest text-brand-gold">Админ-панель</p>
        <h1 class="mt-1 font-serif text-3xl font-bold text-brand-brown sm:text-4xl">Обзор системы</h1>
      </div>
      <div class="text-sm font-bold uppercase tracking-widest text-brand-brown/40">
        {{ new Date().toLocaleDateString('ru-RU', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
      </div>
    </div>
    
    <!-- Top Stats -->
    <div v-if="stats" class="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-lg border border-brand-brown/10 bg-white p-5 shadow-sm sm:p-6">
        <span class="mb-2 block text-xs font-semibold uppercase tracking-widest text-brand-brown/40">Всего проектов</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.projects_count }}</div>
      </div>
      <div class="rounded-lg bg-brand-gold p-5 shadow-xl shadow-brand-gold/20 sm:p-6">
        <span class="mb-2 block text-xs font-semibold uppercase tracking-widest text-brand-brown/60">Заявки сегодня</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.new_orders_today }}</div>
      </div>
      <div class="rounded-lg border border-brand-brown/10 bg-white p-5 shadow-sm sm:p-6">
        <span class="mb-2 block text-xs font-semibold uppercase tracking-widest text-brand-brown/40">Всего заявок</span>
        <div class="text-4xl font-serif text-brand-brown">{{ stats.total_orders }}</div>
      </div>
      <div class="rounded-lg bg-brand-brown p-5 shadow-xl sm:p-6">
        <span class="mb-2 block text-xs font-semibold uppercase tracking-widest text-white/40">Успешных сделок</span>
        <div class="text-4xl font-serif text-brand-gold">{{ stats.success_rate.toFixed(0) }}%</div>
      </div>
    </div>

    <!-- Chart Section -->
    <div v-if="stats?.orders_by_day && stats.orders_by_day.length > 0" class="rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm sm:p-6">
      <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 class="font-serif text-2xl text-brand-brown">Активность заявок</h3>
          <p class="text-brand-brown/40 text-sm font-medium">Динамика за последние 30 дней</p>
        </div>
        <div class="flex items-center gap-2">
          <div class="w-3 h-3 bg-brand-gold rounded-full"></div>
          <span class="text-xs font-bold uppercase tracking-widest text-brand-brown/60">Новые заявки</span>
        </div>
      </div>

      <div class="relative w-full overflow-x-auto">
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
    <div class="grid grid-cols-1 gap-4 xl:grid-cols-12 xl:gap-6">
      
      <!-- 1. Latest Leads (Left - 8 columns) -->
      <div class="flex flex-col overflow-hidden rounded-lg border border-brand-brown/10 bg-white shadow-sm xl:col-span-8">
        <div class="flex items-center justify-between gap-4 border-b border-brand-brown/10 bg-brand-gray/20 p-4 sm:p-6">
          <div>
            <h3 class="font-serif text-2xl">Свежие заявки</h3>
            <p class="text-sm font-medium text-brand-brown/40">Требуют первичной обработки</p>
          </div>
          <router-link to="/admin/orders" class="flex items-center gap-1 text-xs font-bold uppercase tracking-widest text-brand-gold hover:underline">
            Все заявки <LucideChevronRight :size="16" />
          </router-link>
        </div>
        
        <div class="flex-1">
          <div v-if="stats?.recent_orders.length" class="divide-y divide-brand-brown/5">
            <div 
              v-for="order in stats.recent_orders" :key="order.id"
              class="group flex items-center justify-between gap-4 p-4 transition-colors hover:bg-brand-gray/30 sm:p-5"
            >
              <div class="flex min-w-0 items-center gap-4">
                <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                  <LucideClock :size="24" />
                </div>
                <div class="min-w-0">
                  <div class="truncate font-bold text-brand-brown">{{ order.client_name }}</div>
                  <div class="truncate text-xs font-bold uppercase tracking-widest text-brand-gold">{{ order.project_name }}</div>
                </div>
              </div>
              <div class="flex shrink-0 items-center gap-3">
                <div class="text-right hidden md:block">
                  <div class="text-sm font-medium text-brand-brown/60">{{ new Date(order.created_at).toLocaleDateString() }}</div>
                  <div class="text-[10px] text-brand-brown/30 font-bold uppercase">{{ new Date(order.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</div>
                </div>
                <router-link to="/admin/orders" class="flex h-10 w-10 items-center justify-center rounded-lg bg-brand-brown text-white transition-all hover:bg-brand-gold lg:opacity-0 lg:group-hover:opacity-100">
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
      <div class="space-y-4 xl:col-span-4">
        <!-- Review Task -->
        <div 
          :class="[
            'flex h-full flex-col justify-between rounded-lg border p-5 transition-all duration-500 sm:p-6',
            stats?.pending_reviews_count ? 'bg-brand-gold border-brand-gold text-brand-brown shadow-xl shadow-brand-gold/20' : 'bg-white border-brand-brown/5 text-brand-brown opacity-60'
          ]"
        >
          <div>
            <LucideMessageSquare class="mb-5" :size="36" />
            <h3 class="mb-4 font-serif text-2xl leading-tight sm:text-3xl">Модерация отзывов</h3>
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
              'w-full rounded-lg py-3 text-center font-black uppercase tracking-widest transition-all',
              stats?.pending_reviews_count ? 'bg-brand-brown text-white hover:bg-black shadow-lg' : 'bg-brand-gray text-brand-brown/40 pointer-events-none'
            ]"
          >
            Проверить
          </router-link>
        </div>

        <!-- Popular Project Quick Look -->
        <div class="rounded-lg border border-brand-brown/10 bg-white p-5 shadow-sm sm:p-6">
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

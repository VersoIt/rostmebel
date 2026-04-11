<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideFilter, LucideChevronDown, LucideSearch, LucideRotateCcw } from 'lucide-vue-next';

const productStore = useProductStore();
const route = useRoute();
const router = useRouter();

// Filter States (Internal for form)
const search = ref(route.query.search?.toString() || '');
const minPrice = ref(route.query.min_price?.toString() || '');
const maxPrice = ref(route.query.max_price?.toString() || '');
const selectedCategory = ref(route.query.category?.toString() || '');
const sortBy = ref(route.query.sort_by?.toString() || 'id');

onMounted(async () => {
  await productStore.fetchCategories();
  loadProducts();
});

const loadProducts = () => {
  const params: any = {
    sort_by: sortBy.value,
    status: 'published',
    search: search.value,
  };
  
  if (selectedCategory.value) params.category_id = selectedCategory.value;
  if (minPrice.value) params.min_price = minPrice.value;
  if (maxPrice.value) params.max_price = maxPrice.value;
  
  // Sync URL
  router.push({ query: params });
  productStore.fetchProducts(params);
};

const resetFilters = () => {
  search.value = '';
  minPrice.value = '';
  maxPrice.value = '';
  selectedCategory.value = '';
  sortBy.value = 'id';
  loadProducts();
};

const selectCategory = (id: string) => {
  selectedCategory.value = id;
  loadProducts(); // Category click still triggers immediate load for UX
};
</script>

<template>
  <div class="bg-brand-gray/30 min-h-screen pt-32 pb-24 px-4">
    <div class="max-w-7xl mx-auto">
      <header class="mb-12">
        <h1 class="font-serif text-5xl text-brand-brown mb-6">Каталог мебели</h1>
        
        <div class="flex gap-4 overflow-x-auto pb-4 no-scrollbar">
          <button 
            @click="selectCategory('')"
            :class="[
              'px-6 py-3 rounded-full font-medium transition-all whitespace-nowrap border-2',
              !selectedCategory ? 'bg-brand-brown border-brand-brown text-white shadow-lg' : 'bg-white border-transparent text-brand-brown hover:bg-brand-brown/5'
            ]"
          >
            Все товары
          </button>
          <button 
            v-for="cat in productStore.categories" 
            :key="cat.id"
            @click="selectCategory(cat.id.toString())"
            :class="[
              'px-6 py-3 rounded-full font-medium transition-all whitespace-nowrap border-2',
              selectedCategory === cat.id.toString() ? 'bg-brand-brown border-brand-brown text-white shadow-lg' : 'bg-white border-transparent text-brand-brown hover:bg-brand-brown/5'
            ]"
          >
            {{ cat.name }}
          </button>
        </div>
      </header>

      <div class="flex flex-col md:flex-row gap-8">
        <aside class="w-full md:w-72 shrink-0">
          <div class="bg-white p-8 rounded-3xl border border-brand-brown/5 shadow-sm sticky top-32 space-y-8">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2">
                <LucideFilter :size="20" class="text-brand-gold" />
                <h3 class="font-serif text-xl">Фильтры</h3>
              </div>
              <button @click="resetFilters" class="text-brand-brown/40 hover:text-brand-gold transition-colors p-1">
                <LucideRotateCcw :size="18" />
              </button>
            </div>

            <div>
              <label class="block text-xs font-bold text-brand-brown/40 mb-3 uppercase tracking-widest">Поиск</label>
              <div class="relative">
                <input 
                  v-model="search"
                  type="text"
                  placeholder="Название..."
                  class="w-full pl-10 pr-4 py-3 rounded-xl bg-brand-gray/50 border-none outline-none focus:ring-2 ring-brand-gold/20 text-sm"
                >
                <LucideSearch class="absolute left-3 top-1/2 -translate-y-1/2 text-brand-brown/30" :size="16" />
              </div>
            </div>
            
            <div>
              <label class="block text-xs font-bold text-brand-brown/40 mb-3 uppercase tracking-widest">Цена (₽)</label>
              <div class="flex items-center gap-3">
                <input v-model="minPrice" type="number" placeholder="От" class="w-full px-3 py-3 rounded-xl bg-brand-gray/50 border-none outline-none focus:ring-2 ring-brand-gold/20 text-sm">
                <div class="w-2 h-px bg-brand-brown/20"></div>
                <input v-model="maxPrice" type="number" placeholder="До" class="w-full px-3 py-3 rounded-xl bg-brand-gray/50 border-none outline-none focus:ring-2 ring-brand-gold/20 text-sm">
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-brand-brown/40 mb-3 uppercase tracking-widest">Сортировка</label>
              <div class="relative">
                <select 
                  v-model="sortBy"
                  class="w-full appearance-none bg-brand-gray/50 border-none px-4 py-3 rounded-xl outline-none focus:ring-2 ring-brand-gold/20 text-sm cursor-pointer"
                >
                  <option value="id">По умолчанию</option>
                  <option value="price">По цене</option>
                  <option value="views_count">Популярные</option>
                  <option value="created_at">Новинки</option>
                </select>
                <LucideChevronDown class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-brand-brown/40" :size="16" />
              </div>
            </div>

            <button 
              @click="loadProducts"
              class="w-full bg-brand-brown text-white py-4 rounded-xl font-bold hover:bg-brand-gold transition-all shadow-lg shadow-brand-brown/10 active:scale-[0.98]"
            >
              ПРИМЕНИТЬ
            </button>
          </div>
        </aside>

        <main class="flex-1">
          <div v-if="productStore.loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <div v-for="i in 6" :key="i" class="animate-pulse bg-white rounded-3xl h-96"></div>
          </div>
          
          <div v-else-if="productStore.products.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <ProductCard v-for="p in productStore.products" :key="p.id" :product="p" />
          </div>

          <div v-else class="bg-white py-32 text-center rounded-[3rem] border border-dashed border-brand-brown/10 flex flex-col items-center">
            <div class="w-20 h-20 bg-brand-gray rounded-full flex items-center justify-center text-brand-brown/20 mb-6">
              <LucideSearch :size="40" />
            </div>
            <h3 class="text-2xl font-serif text-brand-brown mb-2">Ничего не нашли</h3>
            <p class="text-brand-brown/40 mb-8">Измените параметры фильтров</p>
            <button @click="resetFilters" class="text-brand-gold font-bold hover:underline">Сбросить всё</button>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

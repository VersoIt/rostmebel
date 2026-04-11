<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideSearch, LucideFilterX, LucideLoader2 } from 'lucide-vue-next';

const productStore = useProductStore();
const route = useRoute();
const router = useRouter();

const selectedCategory = ref(route.query.category?.toString() || '');
const searchQuery = ref('');

const fetch = async () => {
  const params: any = { status: 'published' };
  if (selectedCategory.value) {
    const cat = productStore.categories.find(c => c.slug === selectedCategory.value);
    if (cat) params.project_category_id = cat.id;
  }
  if (searchQuery.value) params.search = searchQuery.value;
  
  await productStore.fetchProducts(params);
};

onMounted(async () => {
  await productStore.fetchCategories();
  await fetch();
});

const selectCategory = (slug: string) => {
  selectedCategory.value = slug;
  router.push({ query: { ...route.query, category: slug || undefined } });
  fetch();
};

const resetFilters = () => {
  selectedCategory.value = '';
  searchQuery.value = '';
  router.push({ query: {} });
  fetch();
};
</script>

<template>
  <div class="bg-brand-cream/20 min-h-screen pt-32 pb-24">
    <div class="max-w-7xl mx-auto px-6">
      <header class="mb-16">
        <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block">Портфолио</span>
        <h1 class="font-serif text-5xl md:text-7xl text-brand-brown mb-8">Наши проекты</h1>
        
        <!-- Category Filters -->
        <div class="flex flex-wrap gap-3">
          <button 
            @click="selectCategory('')"
            :class="[
              'px-8 py-3.5 rounded-full font-bold text-xs uppercase tracking-widest transition-all border-2',
              !selectedCategory ? 'bg-brand-brown border-brand-brown text-white shadow-xl scale-105' : 'bg-white border-brand-brown/5 text-brand-brown/60 hover:border-brand-gold hover:text-brand-gold'
            ]"
          >
            Все работы
          </button>
          <button 
            v-for="cat in productStore.categories" 
            :key="cat.id"
            @click="selectCategory(cat.slug)"
            :class="[
              'px-8 py-3.5 rounded-full font-bold text-xs uppercase tracking-widest transition-all border-2',
              selectedCategory === cat.slug ? 'bg-brand-brown border-brand-brown text-white shadow-xl scale-105' : 'bg-white border-brand-brown/5 text-brand-brown/60 hover:border-brand-gold hover:text-brand-gold'
            ]"
          >
            {{ cat.name }}
          </button>
        </div>
      </header>

      <!-- Search & Results -->
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-12">
        <!-- Sidebar Search -->
        <aside class="lg:col-span-1">
          <div class="sticky top-32 space-y-8">
            <div class="relative group">
              <input 
                v-model="searchQuery"
                @input="fetch"
                type="text" 
                placeholder="Поиск проекта..." 
                class="w-full pl-12 pr-4 py-4 rounded-2xl bg-white border-2 border-transparent focus:border-brand-gold outline-none shadow-sm transition-all group-hover:shadow-md font-medium"
              >
              <LucideSearch class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/20" :size="20" />
            </div>

            <div class="p-8 bg-brand-brown rounded-[2rem] text-white shadow-2xl relative overflow-hidden">
              <div class="relative z-10">
                <h4 class="font-serif text-xl mb-4 text-brand-gold">Нужна помощь?</h4>
                <p class="text-white/60 text-sm leading-relaxed mb-6">Наш ИИ поможет подобрать проект под ваши требования.</p>
                <router-link to="/#ai-search" class="inline-flex items-center gap-2 text-xs font-black uppercase tracking-widest hover:text-brand-gold transition-colors">
                  Запустить ИИ-поиск
                  <LucideArrowRight :size="14" />
                </router-link>
              </div>
              <div class="absolute -bottom-10 -right-10 w-32 h-32 bg-white/5 rounded-full blur-3xl"></div>
            </div>
          </div>
        </aside>

        <!-- Main Grid -->
        <main class="lg:col-span-3">
          <div v-if="productStore.loading" class="flex flex-col items-center justify-center py-40 gap-4">
            <LucideLoader2 class="animate-spin text-brand-gold" :size="48" />
            <span class="text-brand-brown/40 font-bold text-xs uppercase tracking-widest">Загружаем шедевры...</span>
          </div>

          <div v-else-if="productStore.products.length > 0">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 animate-in">
              <ProductCard v-for="p in productStore.products" :key="p.id" :product="p" />
            </div>
            
            <div class="mt-16 pt-8 border-t border-brand-brown/5 flex justify-between items-center">
              <span class="text-sm text-brand-brown/40 font-medium">Показано {{ productStore.products.length }} проектов</span>
            </div>
          </div>

          <div v-else class="text-center py-40 bg-white rounded-[3rem] border-2 border-dashed border-brand-brown/5">
            <LucideFilterX :size="64" class="mx-auto text-brand-brown/10 mb-6" />
            <h3 class="text-2xl font-serif text-brand-brown mb-2">Ничего не найдено</h3>
            <p class="text-brand-brown/40 mb-8">Попробуйте изменить параметры поиска или категорию</p>
            <button @click="resetFilters" class="text-brand-gold font-bold hover:underline uppercase text-xs tracking-widest">Сбросить всё</button>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-in {
  animation: fade-up 0.6s cubic-bezier(0.2, 1, 0.3, 1) forwards;
}

@keyframes fade-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useProductStore } from '@/stores/products';
import AISearchPanel from '@/components/ai/AISearchPanel.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { gsap } from 'gsap';
import { LucideChevronRight } from 'lucide-vue-next';
import type { Product } from '@/types';

const productStore = useProductStore();
const hits = ref<Product[]>([]);

onMounted(async () => {
  await productStore.fetchProducts({ limit: 4, sort_by: 'views_count', sort_order: 'desc', status: 'published' });
  hits.value = productStore.products;
  
  gsap.from('.hero-text', { 
    y: 50, 
    opacity: 0, 
    duration: 1, 
    stagger: 0.2,
    ease: 'power3.out' 
  });
});
</script>

<template>
  <div class="bg-brand-cream min-h-screen">
    <!-- Hero Section -->
    <section class="relative h-screen flex items-center justify-center overflow-hidden">
      <div class="absolute inset-0 z-0">
        <img src="https://images.unsplash.com/photo-1555041469-a586c61ea9bc?auto=format&fit=crop&q=80" 
             class="w-full h-full object-cover opacity-20" alt="Furniture">
      </div>
      
      <div class="relative z-10 text-center px-4">
        <h1 class="hero-text font-serif text-5xl md:text-7xl text-brand-brown mb-6">
          Эстетика РОСТа <br> в вашем доме
        </h1>
        <p class="hero-text text-xl text-brand-brown/80 mb-10 max-w-2xl mx-auto">
          Создаем мебель, которая вдохновляет на новую жизнь. Качество, проверенное временем.
        </p>
        <div class="hero-text flex flex-col md:flex-row gap-4 justify-center">
          <router-link to="/catalog" 
            class="bg-brand-brown text-brand-white px-8 py-4 rounded-lg font-medium hover:bg-brand-gold transition-colors flex items-center justify-center gap-2">
            Смотреть каталог
            <LucideChevronRight :size="20" />
          </router-link>
          <a href="#ai-search" 
            class="bg-white text-brand-brown border border-brand-brown/10 px-8 py-4 rounded-lg font-medium hover:bg-brand-gray transition-colors">
            Подобрать с ИИ
          </a>
        </div>
      </div>
    </section>

    <!-- AI Search Section -->
    <section id="ai-search" class="py-24 px-4 bg-white">
      <div class="max-w-4xl mx-auto text-center mb-16">
        <h2 class="font-serif text-4xl text-brand-brown mb-4">Умный поиск мебели</h2>
        <p class="text-brand-brown/60">Опишите ваши предпочтения, и наш ИИ подберет идеальные варианты</p>
      </div>
      <AISearchPanel />
    </section>

    <!-- Hits Section -->
    <section class="py-24 px-4 max-w-7xl mx-auto">
      <div class="flex items-end justify-between mb-12">
        <div>
          <h2 class="font-serif text-4xl text-brand-brown mb-2">Хиты продаж</h2>
          <p class="text-brand-brown/60">Мебель, которую выбирают чаще всего</p>
        </div>
        <router-link to="/catalog" class="text-brand-gold font-medium hover:underline">
          Весь каталог
        </router-link>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <ProductCard v-for="product in hits" :key="product.id" :product="product" />
      </div>
    </section>

    <!-- Why Us -->
    <section class="py-24 px-4 bg-brand-brown text-brand-white">
      <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-12 text-center">
        <div>
          <div class="text-brand-gold text-4xl mb-4 font-serif">15+</div>
          <h3 class="text-xl mb-2 font-serif">Лет опыта</h3>
          <p class="text-brand-white/60">Мы знаем о мебели всё и даже больше</p>
        </div>
        <div>
          <div class="text-brand-gold text-4xl mb-4 font-serif">24/7</div>
          <h3 class="text-xl mb-2 font-serif">Поддержка</h3>
          <p class="text-brand-white/60">Всегда на связи для решения ваших вопросов</p>
        </div>
        <div>
          <div class="text-brand-gold text-4xl mb-4 font-serif">100%</div>
          <h3 class="text-xl mb-2 font-serif">Гарантия</h3>
          <p class="text-brand-white/60">Контроль качества на каждом этапе производства</p>
        </div>
      </div>
    </section>
  </div>
</template>

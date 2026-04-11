<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useProductStore } from '@/stores/products';
import AISearchPanel from '@/components/ai/AISearchPanel.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { gsap } from 'gsap';
import { 
  LucideChevronRight, 
  LucideArrowRight, 
  LucideRuler, 
  LucidePenTool, 
  LucideHammer, 
  LucideTruck 
} from 'lucide-vue-next';
import type { Product } from '@/types';

const productStore = useProductStore();
const hits = ref<Product[]>([]);

onMounted(async () => {
  await productStore.fetchProducts({ limit: 6, sort_by: 'views_count', sort_order: 'desc', status: 'published' });
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
  <div class="bg-brand-cream min-h-screen text-brand-brown">
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
            Наши проекты
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
        <h2 class="font-serif text-4xl text-brand-brown mb-4">Умный поиск по проектам</h2>
        <p class="text-brand-brown/60">Опишите ваши пожелания, и наш ИИ подберет похожие реализованные проекты</p>
      </div>
      <AISearchPanel />
    </section>

    <!-- Process Section -->
    <section class="py-32 px-4 bg-brand-cream/50 relative overflow-hidden">
      <div class="max-w-7xl mx-auto relative z-10">
        <div class="text-center mb-20">
          <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block">Путь к идеалу</span>
          <h2 class="font-serif text-5xl text-brand-brown">Этапы реализации проекта</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-12 relative">
          <!-- Connector line (hidden on mobile) -->
          <div class="hidden md:block absolute top-10 left-0 right-0 h-0.5 bg-brand-gold/10 -z-0"></div>

          <!-- Step 1 -->
          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucideRuler :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">01</div>
            </div>
            <div class="text-center">
              <h3 class="font-serif text-xl mb-4">Бесплатный замер</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Выезжаем на объект, делаем точные замеры и обсуждаем ваши пожелания.</p>
            </div>
          </div>

          <!-- Step 2 -->
          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucidePenTool :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">02</div>
            </div>
            <div class="text-center">
              <h3 class="font-serif text-xl mb-4">3D-проектирование</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Создаем фотореалистичный проект и подбираем материалы под ваш бюджет.</p>
            </div>
          </div>

          <!-- Step 3 -->
          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucideHammer :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">03</div>
            </div>
            <div class="text-center">
              <h3 class="font-serif text-xl mb-4">Производство</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Изготавливаем мебель на собственном производстве с контролем качества.</p>
            </div>
          </div>

          <!-- Step 4 -->
          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucideTruck :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">04</div>
            </div>
            <div class="text-center">
              <h3 class="font-serif text-xl mb-4">Монтаж и сборка</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Бережно доставляем и профессионально устанавливаем готовую мебель.</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Hits Section -->
    <section class="py-24 px-4 max-w-7xl mx-auto">
      <div class="flex items-end justify-between mb-12">
        <div>
          <h2 class="font-serif text-4xl text-brand-brown mb-2">Популярные проекты</h2>
          <p class="text-brand-brown/60">Решения, которые вдохновляют чаще всего</p>
        </div>
        <router-link to="/catalog" class="text-brand-gold font-medium hover:underline">
          Все проекты
        </router-link>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8">
        <ProductCard v-for="product in hits.slice(0, 5)" :key="product.id" :product="product" />
        
        <!-- More Card -->
        <router-link 
          v-if="hits.length > 5"
          to="/catalog" 
          class="bg-brand-brown rounded-3xl p-8 flex flex-col items-center justify-center text-center group transition-all hover:bg-brand-gold shadow-xl"
        >
          <div class="w-16 h-16 bg-white/10 rounded-full flex items-center justify-center text-white mb-6 group-hover:scale-110 transition-transform">
            <LucideArrowRight :size="32" />
          </div>
          <h3 class="text-white font-serif text-xl mb-2">Смотреть всё</h3>
          <p class="text-white/60 text-sm">Полное портфолио работ</p>
        </router-link>
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

    <!-- SEO Content Section -->
    <section class="py-24 px-4 bg-white">
      <div class="max-w-4xl mx-auto prose prose-brand">
        <h2 class="font-serif text-3xl text-brand-brown mb-8 text-center">Мебель по индивидуальным проектам в Симферополе и Крыму</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-12 text-brand-brown/70 leading-relaxed text-sm">
          <p>
            Компания <strong>РОСТ Мебель</strong> специализируется на проектировании и изготовлении премиальной мебели по индивидуальным размерам. Наше производство оснащено современным оборудованием, что позволяет создавать изделия любой сложности: от классических кухонь из массива до современных шкафов-купе в стиле минимализм.
          </p>
          <p>
            Мы предлагаем комплексный подход: от замера и создания 3D-проекта до профессионального монтажа. Используем только проверенные материалы и надежную фурнитуру от ведущих мировых производителей. Наша мебель на заказ — это сочетание эргономики, стиля и долговечности для вашего дома или офиса.
          </p>
        </div>
      </div>
    </section>
  </div>
</template>

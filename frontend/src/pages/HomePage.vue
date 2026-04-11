<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
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

// Hero Slider
const heroImages = [
  'https://images.unsplash.com/photo-1556911223-05345a39365e?q=80&w=1600&auto=format&fit=crop', // Modern Dark
  'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?q=80&w=1600&auto=format&fit=crop', // White Classic
  'https://images.unsplash.com/photo-1600585152220-90363fe7e115?q=80&w=1600&auto=format&fit=crop', // Marble Luxury
  'https://images.unsplash.com/photo-1565183275290-0b174804bb3d?q=80&w=1600&auto=format&fit=crop'  // Minimalist
];
const currentHeroIndex = ref(0);
let heroInterval: any = null;

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

  heroInterval = setInterval(() => {
    currentHeroIndex.value = (currentHeroIndex.value + 1) % heroImages.length;
  }, 5000);
});

onUnmounted(() => {
  if (heroInterval) clearInterval(heroInterval);
});
</script>

<template>
  <div class="bg-brand-cream min-h-screen text-brand-brown">
    <!-- Hero Section -->
    <section class="relative h-screen flex items-center justify-center overflow-hidden bg-brand-brown">
      <!-- Background Slider -->
      <div class="absolute inset-0 z-0">
        <transition name="fade-bg">
          <div :key="currentHeroIndex" class="absolute inset-0">
            <div 
              class="absolute inset-0 bg-cover bg-center transition-transform duration-[10000ms] ease-linear scale-100 hover:scale-110"
              :style="{ backgroundImage: `url(${heroImages[currentHeroIndex]})`, opacity: 0.6 }"
            ></div>
          </div>
        </transition>
        <!-- Overlay for readability -->
        <div class="absolute inset-0 bg-black/20"></div>
        <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-transparent to-brand-cream"></div>
      </div>
      
      <div class="relative z-10 text-center px-4">
        <h1 class="hero-text font-serif text-5xl md:text-8xl text-white mb-6 drop-shadow-[0_10px_10px_rgba(0,0,0,0.5)] leading-tight">
          Эстетика РОСТа <br> <span class="text-brand-gold">в вашем доме</span>
        </h1>
        <p class="hero-text text-xl md:text-2xl text-white/90 mb-12 max-w-3xl mx-auto font-light tracking-wide drop-shadow-lg">
          Создаем кухни и авторские интерьеры, которые вдохновляют на новую жизнь.
        </p>
        <div class="hero-text flex flex-col md:flex-row gap-6 justify-center items-center">
          <router-link to="/catalog" 
            class="bg-brand-gold text-brand-brown px-12 py-5 rounded-full font-black uppercase tracking-widest hover:bg-white transition-all shadow-2xl hover:scale-105 flex items-center gap-3">
            Наши проекты
            <LucideChevronRight :size="20" />
          </router-link>
          <a href="#ai-search" 
            class="bg-white/10 backdrop-blur-xl text-white border-2 border-white/30 px-12 py-5 rounded-full font-black uppercase tracking-widest hover:bg-white/20 transition-all hover:scale-105">
            Подобрать с ИИ
          </a>
        </div>
      </div>

      <!-- Bottom indicator -->
      <div class="absolute bottom-12 left-1/2 -translate-x-1/2 z-10 flex gap-3">
        <div 
          v-for="(_, idx) in heroImages" 
          :key="idx"
          class="h-1 transition-all duration-500 rounded-full"
          :class="currentHeroIndex === idx ? 'w-12 bg-brand-gold' : 'w-6 bg-white/30'"
        ></div>
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
    <section class="py-32 px-4 bg-brand-cream/50 relative overflow-hidden border-y border-brand-brown/5">
      <div class="max-w-7xl mx-auto relative z-10">
        <div class="text-center mb-20">
          <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block">Путь к идеалу</span>
          <h2 class="font-serif text-5xl text-brand-brown">Этапы реализации проекта</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-12 relative">
          <div class="hidden md:block absolute top-10 left-0 right-0 h-0.5 bg-brand-gold/10 -z-0"></div>

          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucideRuler :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">01</div>
            </div>
            <div class="text-center px-4">
              <h3 class="font-serif text-xl mb-4">Бесплатный замер</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Выезжаем на объект, делаем точные замеры и обсуждаем ваши пожелания.</p>
            </div>
          </div>

          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucidePenTool :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">02</div>
            </div>
            <div class="text-center px-4">
              <h3 class="font-serif text-xl mb-4">3D-проектирование</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Создаем фотореалистичный проект и подбираем материалы под ваш бюджет.</p>
            </div>
          </div>

          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucideHammer :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">03</div>
            </div>
            <div class="text-center px-4">
              <h3 class="font-serif text-xl mb-4">Производство</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed">Изготавливаем мебель на собственном производстве с контролем качества.</p>
            </div>
          </div>

          <div class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <LucideTruck :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">04</div>
            </div>
            <div class="text-center px-4">
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
        <h2 class="font-serif text-3xl text-brand-brown mb-8 text-center">Мебель по индивидуальным проектам в Севастополе и Крыму</h2>
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

<style scoped>
.fade-bg-enter-active, .fade-bg-leave-active { 
  transition: opacity 2s ease-in-out; 
}
.fade-bg-enter-from, .fade-bg-leave-to { 
  opacity: 0; 
}

/* Эффект плавного зума фона */
.bg-cover {
  animation: ken-burns 20s infinite alternate;
}

@keyframes ken-burns {
  from { transform: scale(1); }
  to { transform: scale(1.1); }
}
</style>

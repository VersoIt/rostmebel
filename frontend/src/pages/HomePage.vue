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
  LucideTruck,
  LucideLayout,
  LucideCrown,
  LucideMapPin,
  LucideZap,
  LucideCompass,
  LucideLayers,
  LucideShieldCheck,
  LucideMessageSquare
} from 'lucide-vue-next';
import type { Product } from '@/types';

const productStore = useProductStore();
const hits = ref<Product[]>([]);

// Hero Slider - Using LOCAL images
const heroImages = [
  '/assets/images/hero-1.jpg',
  '/assets/images/hero-2.jpg',
  '/assets/images/hero-3.jpg'
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
  <div class="bg-brand-cream min-h-screen text-brand-brown font-sans selection:bg-brand-gold selection:text-white">
    <!-- 1. Hero Section -->
    <section class="relative h-screen flex items-center justify-center overflow-hidden bg-neutral-900">
      <div class="absolute inset-0 z-0">
        <div 
          v-for="(img, idx) in heroImages" 
          :key="idx"
          class="absolute inset-0 transition-opacity duration-[2000ms] ease-in-out"
          :style="{ 
            opacity: currentHeroIndex === idx ? 1 : 0,
            zIndex: currentHeroIndex === idx ? 10 : 0
          }"
        >
          <img 
            :src="img" 
            class="w-full h-full object-cover"
            :class="{ 'animate-ken-burns': currentHeroIndex === idx }"
            alt="Premium Kitchen"
          >
        </div>
        
        <div class="absolute inset-0 bg-black/30 z-[11]"></div>
        <div class="absolute inset-0 bg-gradient-to-b from-black/20 via-transparent to-brand-cream/40 z-[12]"></div>
      </div>
      
      <div class="relative z-20 text-center px-4">
        <div class="hero-text mb-6 inline-flex items-center gap-2 bg-white/10 backdrop-blur-md px-4 py-2 rounded-full border border-white/20 text-white text-xs font-bold uppercase tracking-widest">
          <LucideMapPin :size="14" class="text-brand-gold" />
          Работаем по всему Крыму
        </div>
        <h1 class="hero-text font-serif text-5xl md:text-8xl text-white mb-6 drop-shadow-[0_10px_10px_rgba(0,0,0,0.5)] leading-tight">
          Идеальные кухни <br> <span class="text-brand-gold">на заказ</span>
        </h1>
        <p class="hero-text text-xl md:text-2xl text-white/90 mb-12 max-w-3xl mx-auto font-light tracking-wide drop-shadow-lg text-center">
          Создаем авторские кухни и корпусную мебель любой сложности. Индивидуальный подход к каждому сантиметру вашего пространства.
        </p>
        <div class="hero-text flex flex-col md:flex-row gap-6 justify-center items-center">
          <router-link to="/catalog" 
            class="bg-brand-gold text-brand-brown px-12 py-5 rounded-full font-black uppercase tracking-widest hover:bg-white transition-all shadow-2xl hover:scale-105 flex items-center gap-3">
            Наши работы
            <LucideChevronRight :size="20" />
          </router-link>
          <a href="#projects-grid" 
            class="bg-white/10 backdrop-blur-xl text-white border-2 border-white/30 px-12 py-5 rounded-full font-black uppercase tracking-widest hover:bg-white/20 transition-all hover:scale-105">
            Узнать больше
          </a>
        </div>
      </div>

      <div class="absolute bottom-12 left-1/2 -translate-x-1/2 z-30 flex gap-3">
        <div v-for="(_, idx) in heroImages" :key="idx"
          class="h-1 transition-all duration-500 rounded-full"
          :class="currentHeroIndex === idx ? 'w-12 bg-brand-gold' : 'w-6 bg-white/30'">
        </div>
      </div>
    </section>

    <!-- 2. Hits Section -->
    <section id="projects-grid" class="py-32 px-4 max-w-7xl mx-auto">
      <div class="flex items-end justify-between mb-16">
        <div>
          <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block text-left">Портфолио</span>
          <h2 class="font-serif text-5xl text-brand-brown mb-2 text-left">Наши последние проекты</h2>
          <p class="text-brand-brown/60 text-lg text-left">Решения, которые вдохновляют на перемены</p>
        </div>
        <router-link to="/catalog" class="hidden md:flex items-center gap-2 text-brand-gold font-bold hover:underline group">
          Смотреть все работы
          <LucideArrowRight :size="20" class="group-hover:translate-x-2 transition-transform" />
        </router-link>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
        <ProductCard v-for="product in hits.slice(0, 6)" :key="product.id" :product="product" />
      </div>

      <div class="mt-16 text-center md:hidden">
        <router-link to="/catalog" class="inline-flex items-center gap-2 bg-brand-brown text-white px-8 py-4 rounded-xl font-bold">
          Все проекты
          <LucideArrowRight :size="20" />
        </router-link>
      </div>
    </section>

    <!-- 3. Supervision Section -->
    <section class="py-32 px-4 bg-[#1a1410] text-white relative overflow-hidden">
      <div class="absolute top-0 right-0 w-1/3 h-full bg-brand-gold/5 skew-x-12 translate-x-20"></div>
      <div class="max-w-7xl mx-auto relative z-10">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-24 items-center text-left">
          <div>
            <div class="inline-flex items-center gap-2 text-brand-gold font-bold text-xs uppercase tracking-[0.4em] mb-6">
              <span class="w-10 h-px bg-brand-gold"></span>
              Expert Supervision
            </div>
            <h2 class="font-serif text-5xl md:text-6xl mb-10 leading-tight">Мы не просто рисуем, <br> <span class="text-brand-gold italic">мы строим.</span></h2>
            <p class="text-white/60 text-lg mb-12 leading-relaxed max-w-xl text-left">
              Чтобы мебель встала идеально, помещение должно быть подготовлено безупречно. Мы берем на себя авторское сопровождение и выдаем строителям точные технические карты.
            </p>
            
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-8 text-left">
              <div class="space-y-4 p-6 bg-white/5 rounded-2xl border border-white/5 hover:border-brand-gold/30 transition-colors group">
                <LucideZap class="text-brand-gold group-hover:scale-110 transition-transform" :size="32" />
                <h4 class="font-bold text-xl text-white">Электрика</h4>
                <p class="text-sm text-white/40">План розеток под вашу технику и подсветку.</p>
              </div>
              <div class="space-y-4 p-6 bg-white/5 rounded-2xl border border-white/5 hover:border-brand-gold/30 transition-colors group">
                <LucideCompass class="text-brand-gold group-hover:scale-110 transition-transform" :size="32" />
                <h4 class="font-bold text-xl text-white">Плитка</h4>
                <p class="text-sm text-white/40">Раскладка фартука для идеального примыкания.</p>
              </div>
              <div class="space-y-4 p-6 bg-white/5 rounded-2xl border border-white/5 hover:border-brand-gold/30 transition-colors group">
                <LucideLayers class="text-brand-gold group-hover:scale-110 transition-transform" :size="32" />
                <h4 class="font-bold text-xl text-white">Инженерия</h4>
                <p class="text-sm text-white/40">Разметка уровней потолков и сантехники.</p>
              </div>
              <div class="space-y-4 p-6 bg-white/5 rounded-2xl border border-white/5 hover:border-brand-gold/30 transition-colors group">
                <LucideShieldCheck class="text-brand-gold group-hover:scale-110 transition-transform" :size="32" />
                <h4 class="font-bold text-xl text-white">Технадзор</h4>
                <p class="text-sm text-white/40">Контроль подготовки стен и углов 90°.</p>
              </div>
            </div>
          </div>
          
          <div class="relative text-left">
            <div class="aspect-[4/5] rounded-[3rem] overflow-hidden shadow-2xl border border-white/10">
              <img src="/assets/images/hero-1.jpg" class="w-full h-full object-cover grayscale-[0.5] hover:grayscale-0 transition-all duration-700" alt="Technical Drawing">
            </div>
            <div class="absolute -bottom-10 -left-10 bg-brand-gold p-10 rounded-3xl shadow-2xl hidden md:block text-left">
              <div class="text-brand-brown font-serif text-4xl mb-2">0%</div>
              <div class="text-brand-brown/80 text-xs font-bold uppercase tracking-widest text-left leading-tight">ошибок при <br> монтаже</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 4. Services Section -->
    <section class="py-32 px-4 bg-white">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-20 items-center">
          <div class="grid grid-cols-2 gap-4">
            <img src="/assets/images/hero-2.jpg" class="rounded-3xl aspect-[3/4] object-cover mt-12 shadow-2xl" alt="Interior 1">
            <img src="/assets/images/hero-3.jpg" class="rounded-3xl aspect-[3/4] object-cover shadow-2xl" alt="Interior 2">
          </div>
          <div class="text-left">
            <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block">All-in-one</span>
            <h2 class="font-serif text-5xl text-brand-brown mb-8 leading-tight">Комплексное меблирование дома</h2>
            <p class="text-lg text-brand-brown/70 mb-10 leading-relaxed text-left">
              Мы создаем единую экосистему вашего интерьера. Никакого «разнобоя» в материалах и стилях. Наша команда закроет все потребности в мебели: от кухни до прикроватной тумбочки.
            </p>
            <div class="space-y-8 text-left">
              <div class="flex gap-6 items-start">
                <div class="w-14 h-14 bg-brand-cream rounded-2xl flex items-center justify-center text-brand-gold shrink-0 border border-brand-gold/10">
                  <LucideCrown :size="28" />
                </div>
                <div class="text-left">
                  <h4 class="font-bold text-xl text-brand-brown mb-2">Фурнитура ТОП-уровня</h4>
                  <p class="text-brand-brown/60">Blum, Hettich, Grass. Петли и ящики, которые работают бесшумно десятилетиями.</p>
                </div>
              </div>
              <div class="flex gap-6 items-start text-left">
                <div class="w-14 h-14 bg-brand-cream rounded-2xl flex items-center justify-center text-brand-gold shrink-0 border border-brand-gold/10">
                  <LucideLayout :size="28" />
                </div>
                <div class="text-left">
                  <h4 class="font-bold text-xl text-brand-brown mb-2">Эксклюзивные фасады</h4>
                  <p class="text-brand-brown/60">Итальянский пластик Fenix, матовая эмаль, натуральный шпон и компакт-плиты.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 5. Process Section -->
    <section class="py-32 px-4 bg-brand-cream/50 relative overflow-hidden border-y border-brand-brown/5">
      <div class="max-w-7xl mx-auto relative z-10 text-center">
        <div class="text-center mb-20 text-center">
          <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block">Workflow</span>
          <h2 class="font-serif text-5xl text-brand-brown text-center">Как мы работаем</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-12 relative">
          <div class="hidden md:block absolute top-10 left-0 right-0 h-0.5 bg-brand-gold/10 -z-0"></div>

          <div v-for="(step, idx) in [
            { icon: LucideRuler, title: 'Замер', desc: 'Бесплатный выезд дизайнера по всему Крыму для точного обмера.' },
            { icon: LucidePenTool, title: 'Проект', desc: 'Создание фотореалистичного 3D-проекта и подбор материалов.' },
            { icon: LucideHammer, title: 'Фабрика', desc: 'Изготовление на современном производстве с контролем ОТК.' },
            { icon: LucideTruck, title: 'Монтаж', desc: 'Профессиональная установка «под ключ» в любой точке полуострова.' }
          ]" :key="idx" class="relative group">
            <div class="w-20 h-20 bg-white rounded-3xl shadow-xl flex items-center justify-center mb-8 mx-auto group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 relative z-10 border border-brand-brown/5">
              <component :is="step.icon" :size="32" />
              <div class="absolute -top-3 -right-3 w-8 h-8 bg-brand-brown text-white text-xs font-bold rounded-full flex items-center justify-center">0{{ idx + 1 }}</div>
            </div>
            <div class="text-center px-4">
              <h3 class="font-serif text-xl mb-4 text-center">{{ step.title }}</h3>
              <p class="text-brand-brown/60 text-sm leading-relaxed text-center">{{ step.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 6. AI Search Section -->
    <section id="ai-search" class="py-32 px-4 bg-white text-center">
      <div class="max-w-4xl mx-auto text-center mb-16">
        <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block">Инновации</span>
        <h2 class="font-serif text-4xl text-brand-brown mb-4 text-center">Не нашли то, что искали?</h2>
        <p class="text-brand-brown/60 text-lg text-center">Опишите вашу идеальную кухню, и наш ИИ подберет похожие проекты из нашего опыта</p>
      </div>
      <AISearchPanel />
    </section>

    <!-- 7. Why Us -->
    <section class="py-24 px-4 bg-brand-brown text-brand-white text-center">
      <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-12 text-center">
        <div>
          <div class="text-brand-gold text-5xl mb-4 font-serif">15+</div>
          <h3 class="text-2xl mb-2 font-serif text-center">Лет на рынке</h3>
          <p class="text-brand-white/60 text-center">Реализовали более 1000 проектов разной сложности</p>
        </div>
        <div>
          <div class="text-brand-gold text-5xl mb-4 font-serif">2 года</div>
          <h3 class="text-2xl mb-2 font-serif text-center">Прямая гарантия</h3>
          <p class="text-brand-white/60 text-center">Честная гарантия на мебель и фурнитуру по договору</p>
        </div>
        <div>
          <div class="text-brand-gold text-5xl mb-4 font-serif">100%</div>
          <h3 class="text-2xl mb-2 font-serif text-center">Своя фабрика</h3>
          <p class="text-brand-white/60 text-center">Никаких посредников — полный цикл производства</p>
        </div>
      </div>
    </section>

    <!-- 8. Final CTA -->
    <section class="py-32 px-4 bg-brand-gold relative overflow-hidden text-center">
      <div class="absolute inset-0 bg-black/5"></div>
      <div class="max-w-4xl mx-auto text-center relative z-10">
        <h2 class="font-serif text-5xl text-brand-brown mb-8 text-center">Готовы создать интерьер <br> своей мечты?</h2>
        <p class="text-brand-brown/80 text-xl mb-12 text-center">Запишитесь на бесплатный замер и получите 3D-проектирование в подарок!</p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <router-link to="/contact" class="bg-brand-brown text-white px-12 py-5 rounded-full font-black uppercase tracking-widest hover:bg-black transition-all shadow-2xl flex items-center justify-center gap-3">
            <LucideMessageSquare :size="20" />
            Обсудить проект
          </router-link>
          <a href="tel:+79787631603" class="bg-white text-brand-brown px-12 py-5 rounded-full font-black uppercase tracking-widest hover:bg-brand-cream transition-all shadow-xl flex items-center justify-center gap-3">
            Позвонить нам
          </a>
        </div>
      </div>
    </section>

    <!-- 9. About / SEO Section -->
    <section class="py-32 px-4 bg-white border-t border-brand-brown/5">
      <div class="max-w-5xl mx-auto">
        <h2 class="font-serif text-4xl text-brand-brown mb-12 text-center">Искусство создания <span class="text-brand-gold italic">вашего пространства</span></h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-16 text-brand-brown/80 leading-relaxed text-lg text-left">
          <p class="font-light">
            Компания <strong class="font-bold text-brand-brown">РОСТ Мебель</strong> — это не просто производство, это команда профессионалов, влюбленных в свое дело. Мы специализируемся на проектировании и изготовлении премиальной мебели по индивидуальным размерам в Севастополе и по всему Крыму, превращая сложные технические задачи в элегантные интерьерные решения.
          </p>
          <p class="font-light">
            Наше производство оснащено передовым оборудованием, что позволяет нам работать с любыми материалами: от натурального массива и шпона до итальянского нано-пластика Fenix. Мы обеспечиваем полный цикл авторского сопровождения: от первого наброска и разметки розеток до финальной установки «под ключ».
          </p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.animate-ken-burns {
  animation: ken-burns 20s infinite alternate ease-in-out;
}

@keyframes ken-burns {
  from { transform: scale(1); }
  to { transform: scale(1.15); }
}
</style>

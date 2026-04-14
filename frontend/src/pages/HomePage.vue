<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useProductStore } from '@/stores/products';
import AISearchPanel from '@/components/ai/AISearchPanel.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import QuoteQuiz from '@/components/order/QuoteQuiz.vue';
import {
  LucideArrowRight,
  LucideCheckCircle2,
  LucideMapPin,
  LucideShieldCheck,
} from 'lucide-vue-next';
import type { Product } from '@/types';
import { buildBusinessSchema, removeJsonLd, setJsonLd } from '@/utils/seo';

const productStore = useProductStore();
const hits = ref<Product[]>([]);

const heroImages = [
  '/assets/images/hero-1.jpg',
  '/assets/images/hero-2.jpg',
  '/assets/images/hero-3.jpg',
];

const currentHeroIndex = ref(0);
let heroInterval: number | undefined;

const proof = [
  { value: '15+', label: 'лет в мебели' },
  { value: '2 года', label: 'гарантия по договору' },
  { value: '1 смета', label: 'мебель, техника, монтаж' },
];

const highlights = [
  'Смета до запуска',
  'Карты розеток и выводов',
  'Подбор техники',
  'Производство и монтаж',
];

const applianceItems = [
  'варочные поверхности',
  'духовые шкафы',
  'вытяжки',
  'посудомоечные машины',
  'холодильники',
  'мойки и смесители',
];

onMounted(async () => {
  setJsonLd('schema-business', buildBusinessSchema());

  await productStore.fetchProducts({
    limit: 3,
    sort_by: 'views_count',
    sort_order: 'desc',
    status: 'published',
  });
  hits.value = productStore.products;

  heroInterval = window.setInterval(() => {
    currentHeroIndex.value = (currentHeroIndex.value + 1) % heroImages.length;
  }, 5600);
});

onUnmounted(() => {
  if (heroInterval) window.clearInterval(heroInterval);
  removeJsonLd('schema-business');
});
</script>

<template>
  <div class="bg-brand-cream text-brand-brown">
    <section class="relative flex min-h-[86svh] items-end overflow-hidden bg-neutral-950 pb-10 pt-28 lg:pb-12">
      <div class="absolute inset-0">
        <div
          v-for="(img, idx) in heroImages"
          :key="img"
          class="absolute inset-0 transition-opacity duration-700"
          :style="{ opacity: currentHeroIndex === idx ? 1 : 0 }"
        >
          <img :src="img" class="h-full w-full object-cover" alt="Кухня и корпусная мебель РОСТ Мебель">
        </div>
        <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-black/48 to-black/18"></div>
        <div class="absolute inset-x-0 bottom-0 h-28 bg-gradient-to-t from-brand-cream to-transparent"></div>
      </div>

      <div class="ui-container relative z-10">
        <div class="max-w-4xl text-white motion-fade-up">
          <div class="mb-5 inline-flex items-center gap-2 rounded-lg border border-white/20 bg-black/22 px-3 py-2 text-sm backdrop-blur">
            <LucideMapPin :size="16" class="text-brand-gold" />
            Работаем по Крыму: замер, доставка, монтаж
          </div>

          <h1 class="max-w-3xl font-serif text-4xl font-bold leading-tight sm:text-5xl lg:text-6xl">
            Кухни и корпусная мебель по размеру
          </h1>

          <p class="mt-5 max-w-2xl text-lg leading-8 text-white/82">
            Проект, производство, техника и установка в одной понятной смете. Сразу проверяем размеры, розетки, материалы и сроки.
          </p>

          <div class="mt-8 flex flex-col gap-3 sm:flex-row">
            <router-link to="/catalog" class="ui-button ui-button-accent">
              Посмотреть проекты
              <LucideArrowRight :size="18" />
            </router-link>
            <a href="#quote-quiz" class="ui-button border border-white/30 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
              Получить расчет
            </a>
          </div>
        </div>

        <div class="mt-9 grid max-w-4xl grid-cols-3 gap-2 sm:gap-3">
          <div v-for="item in proof" :key="item.label" class="rounded-lg border border-white/14 bg-black/24 p-3 text-white backdrop-blur sm:p-4">
            <div class="font-serif text-2xl leading-none sm:text-3xl">{{ item.value }}</div>
            <div class="mt-2 text-xs leading-4 text-white/70 sm:text-sm">{{ item.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <section class="border-b border-brand-brown/10 bg-white">
      <div class="ui-container grid grid-cols-1 gap-3 py-7 sm:grid-cols-2 lg:grid-cols-4">
        <div v-for="item in highlights" :key="item" class="flex items-center gap-3 text-sm font-bold text-brand-brown">
          <span class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-brand-gold/10 text-brand-gold">
            <LucideCheckCircle2 :size="17" />
          </span>
          {{ item }}
        </div>
      </div>
    </section>

    <section id="projects-grid" class="ui-section">
      <div class="ui-container">
        <div class="mb-8 flex flex-col justify-between gap-4 md:flex-row md:items-end">
          <div>
            <p class="ui-eyebrow mb-3">Портфолио</p>
            <h2 class="ui-title-lg">Несколько проектов, чтобы быстро понять уровень</h2>
          </div>
          <router-link to="/catalog" class="ui-button ui-button-secondary">
            Все проекты
            <LucideArrowRight :size="18" />
          </router-link>
        </div>

        <div v-if="hits.length" class="grid grid-cols-1 gap-6 md:grid-cols-3">
          <ProductCard v-for="product in hits" :key="product.id" :product="product" />
        </div>
        <div v-else class="ui-empty">
          Проекты появятся после публикации в админке.
        </div>
      </div>
    </section>

    <section class="bg-white">
      <div class="ui-container grid grid-cols-1 gap-8 py-12 lg:grid-cols-[0.95fr_1.05fr] lg:items-center">
        <div>
          <p class="ui-eyebrow mb-3">Техника для кухни</p>
          <h2 class="ui-title-lg">Комплектуем кухню техникой сразу в проекте</h2>
          <p class="ui-copy-lg mt-4">
            Подбираем технику вместе с мебелью, чтобы размеры, розетки, вентиляция, фасады и посадочные места сошлись до запуска производства.
          </p>
          <a href="#quote-quiz" class="ui-button ui-button-primary mt-6">
            Посчитать кухню
            <LucideArrowRight :size="18" />
          </a>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-[0.8fr_1fr]">
          <img src="/assets/images/tech-drawing.jpg" class="h-full min-h-64 rounded-lg object-cover" alt="Технический проект кухни">
          <div class="grid content-center gap-2">
            <div
              v-for="item in applianceItems"
              :key="item"
              class="flex items-center gap-3 border-b border-brand-brown/10 py-3 text-sm font-semibold text-brand-brown/75"
            >
              <LucideCheckCircle2 :size="18" class="text-brand-gold" />
              {{ item }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="ai-search" class="bg-white">
      <div class="ui-container grid grid-cols-1 gap-8 py-12 lg:grid-cols-[0.8fr_1.2fr] lg:items-center">
        <div>
          <p class="ui-eyebrow mb-3">Быстрый подбор</p>
          <h2 class="ui-title-lg">Не хочется листать каталог?</h2>
          <p class="ui-copy-lg mt-4">
            Опишите задачу одной фразой. Например: светлая кухня с техникой до 250 000 ₽.
          </p>
        </div>
        <AISearchPanel />
      </div>
    </section>

    <section id="quote-quiz" class="scroll-mt-28 py-10 sm:py-12">
      <div class="ui-container">
        <div class="grid grid-cols-1 gap-7 rounded-lg bg-brand-brown p-5 text-white lg:grid-cols-[0.9fr_1.1fr] lg:items-center lg:p-8">
          <div class="lg:pr-4">
            <div class="mb-4 flex items-center gap-3 text-brand-gold">
              <LucideShieldCheck :size="22" />
              <span class="font-semibold">Посчитаем без долгих созвонов</span>
            </div>
            <h2 class="font-serif text-3xl font-bold leading-tight sm:text-4xl">Ответьте на 4 вопроса и получите следующий шаг по проекту</h2>
            <p class="mt-4 max-w-2xl leading-8 text-white/72">
              Подскажем реалистичный бюджет, сроки и слабые места планировки до того, как вы вложитесь в материалы, технику и ремонт.
            </p>

            <div class="mt-6 grid gap-3 text-sm font-semibold text-white/72 sm:grid-cols-3">
              <div class="rounded-lg border border-white/10 bg-white/5 p-3">Размеры и техника</div>
              <div class="rounded-lg border border-white/10 bg-white/5 p-3">Бюджет и сроки</div>
              <div class="rounded-lg border border-white/10 bg-white/5 p-3">Замер по Крыму</div>
            </div>

            <a href="tel:+79787631603" class="ui-button mt-6 border border-white/20 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
              Позвонить: +7 (978) 763-16-03
            </a>
          </div>

          <div class="rounded-lg bg-white p-4 text-brand-brown shadow-2xl shadow-black/20 sm:p-5">
            <QuoteQuiz initial-project-type="Кухня с техникой" />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

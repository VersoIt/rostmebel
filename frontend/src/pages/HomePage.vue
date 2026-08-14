<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { useProductStore } from '@/stores/products';
import AISearchPanel from '@/components/ai/AISearchPanel.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import QuoteQuiz from '@/components/order/QuoteQuiz.vue';
import {
  LucideArrowRight,
  LucideCheckCircle2,
  LucideMapPin,
  LucideMessageSquare,
  LucideRuler,
  LucideShieldCheck,
  LucideWrench,
} from 'lucide-vue-next';
import type { Product } from '@/types';
import { buildBusinessSchema, buildWebsiteSchema, removeJsonLd, setJsonLd } from '@/utils/seo';
import { PHONE_DISPLAY, PHONE_HREF } from '@/constants/contacts';

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
  { value: '15+', label: 'лет в мебели на заказ' },
  { value: '2 года', label: 'гарантия по договору' },
  { value: '1 смета', label: 'мебель, техника и монтаж в одном проекте' },
];

const stages = [
  {
    title: 'Замер и проект',
    text: 'Сначала смотрим помещение, размеры, розетки, технику и важные детали.',
    icon: LucideRuler,
  },
  {
    title: 'Понятная смета',
    text: 'Сразу раскладываем стоимость по позициям, чтобы бюджет был понятен заранее.',
    icon: LucideMessageSquare,
  },
  {
    title: 'Производство и монтаж',
    text: 'Изготавливаем мебель и доводим проект до готового результата.',
    icon: LucideWrench,
  },
];

const applianceItems = [
  'варочные поверхности и духовые шкафы',
  'вытяжки и вентиляция',
  'посудомоечные машины',
  'холодильники и встроенные колонны',
  'мойки, смесители и подсветка',
  'схемы розеток и выводов до монтажа',
];

onMounted(async () => {
  setJsonLd('schema-website', buildWebsiteSchema());
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
  removeJsonLd('schema-website');
  removeJsonLd('schema-business');
});
</script>

<template>
  <div class="bg-transparent text-brand-brown">
    <section class="relative isolate overflow-hidden bg-[#08110f] pb-12 pt-28 text-white lg:pb-16">
      <div class="absolute inset-0">
        <div
          v-for="(img, idx) in heroImages"
          :key="img"
          class="absolute inset-0 transition-opacity duration-700"
          :style="{ opacity: currentHeroIndex === idx ? 1 : 0 }"
        >
          <img :src="img" class="h-full w-full object-cover" alt="Мебель на заказ РОСТ Мебель">
        </div>
        <div class="absolute inset-0 bg-[linear-gradient(110deg,rgba(8,17,15,0.94)_8%,rgba(8,17,15,0.74)_46%,rgba(8,17,15,0.38)_100%)]"></div>
        <div class="absolute -left-16 top-10 h-56 w-56 rounded-full bg-brand-gold/20 blur-3xl"></div>
        <div class="absolute right-0 top-1/3 h-72 w-72 rounded-full bg-white/10 blur-3xl"></div>
        <div class="absolute inset-x-0 bottom-0 h-36 bg-gradient-to-t from-[#f2f5f3] to-transparent"></div>
      </div>

      <div class="ui-container relative z-10">
        <div class="max-w-4xl text-white motion-fade-up">
          <div class="mb-5 inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-4 py-2 text-sm backdrop-blur">
            <LucideMapPin :size="16" class="text-brand-gold" />
            Работаем по Крыму: замер, доставка, монтаж
          </div>

          <h1 class="max-w-3xl font-serif text-4xl font-bold leading-[1.02] sm:text-5xl lg:text-[4.35rem]">
            Кухни и корпусная мебель на заказ
          </h1>

          <p class="mt-6 max-w-2xl text-lg leading-8 text-white/80">
            Проектируем, производим и устанавливаем мебель по вашим размерам. Сразу учитываем технику, розетки,
            бюджет и монтаж.
          </p>

          <div class="mt-8 flex flex-col gap-3 sm:flex-row">
            <router-link to="/catalog" class="ui-button ui-button-accent">
              Посмотреть проекты
              <LucideArrowRight :size="18" />
            </router-link>
            <a href="#quote-quiz" class="ui-button border border-white/20 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
              Получить расчет
            </a>
          </div>

          <div class="mt-10 grid max-w-4xl grid-cols-1 gap-3 sm:grid-cols-3">
            <div
              v-for="item in proof"
              :key="item.label"
              class="rounded-2xl border border-white/12 bg-white/8 p-4 backdrop-blur-xl"
            >
              <div class="font-serif text-3xl leading-none text-white">{{ item.value }}</div>
              <div class="mt-2 text-sm leading-5 text-white/68">{{ item.label }}</div>
            </div>
          </div>

          <div class="mt-6 flex items-center gap-2">
            <button
              v-for="(img, idx) in heroImages"
              :key="`hero-dot-${img}`"
              type="button"
              :class="[
                'h-2.5 rounded-full transition-all duration-300',
                currentHeroIndex === idx ? 'w-10 bg-brand-gold' : 'w-2.5 bg-white/35 hover:bg-white/70'
              ]"
              :aria-label="`Показать фото ${idx + 1}`"
              @click="currentHeroIndex = idx"
            />
          </div>
        </div>
      </div>
    </section>

    <section class="ui-section pt-0">
      <div class="ui-container">
        <div class="grid gap-4 md:grid-cols-3">
          <article
            v-for="stage in stages"
            :key="stage.title"
            class="ui-card ui-card-hover p-6"
          >
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-brand-gold/10 text-brand-gold">
              <component :is="stage.icon" :size="20" />
            </div>
            <h2 class="font-serif text-2xl font-bold text-brand-brown">{{ stage.title }}</h2>
            <p class="mt-3 leading-7 text-brand-brown/62">{{ stage.text }}</p>
          </article>
        </div>
      </div>
    </section>

    <section id="projects-grid" class="ui-section pt-4">
      <div class="ui-container">
        <div class="ui-surface overflow-hidden p-5 sm:p-7 lg:p-8">
          <div class="mb-8 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
            <div>
              <p class="ui-eyebrow mb-3">Портфолио</p>
              <h2 class="ui-title-lg">Реальные проекты из нашего портфолио</h2>
              <p class="ui-copy mt-4 max-w-2xl">
                Кухни, шкафы и другие проекты с фотографиями, бюджетом и деталями исполнения.
              </p>
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
      </div>
    </section>

    <section class="ui-section pb-0">
      <div class="ui-container grid grid-cols-1 gap-8 lg:grid-cols-[0.96fr_1.04fr] lg:items-center">
        <div>
          <p class="ui-eyebrow mb-3">Техника для кухни</p>
          <h2 class="ui-title-lg">Продумываем кухню вместе с техникой</h2>
          <p class="ui-copy-lg mt-4">
            Подбираем не только фасады и корпуса. Сразу увязываем мебель с техникой, розетками, вентиляцией,
            колоннами и логикой хранения.
          </p>
          <a href="#quote-quiz" class="ui-button ui-button-primary mt-6">
            Посчитать кухню
            <LucideArrowRight :size="18" />
          </a>
        </div>

        <div class="ui-surface grid gap-4 p-4 sm:grid-cols-[0.84fr_1fr] sm:p-5">
          <img src="/assets/images/tech-drawing.jpg" class="h-full min-h-72 rounded-[1.6rem] object-cover" alt="Технический проект кухни">
          <div class="grid content-center gap-2">
            <div
              v-for="item in applianceItems"
              :key="item"
              class="flex items-center gap-3 rounded-2xl border border-brand-brown/8 bg-white/85 px-4 py-3 text-sm font-semibold text-brand-brown/78"
            >
              <LucideCheckCircle2 :size="18" class="shrink-0 text-brand-gold" />
              {{ item }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="ai-search" class="ui-section">
      <div class="ui-container">
        <div class="mb-8 max-w-2xl">
          <p class="ui-eyebrow mb-3">Быстрый подбор</p>
          <h2 class="ui-title-lg">Не хочется листать каталог вручную?</h2>
          <p class="ui-copy-lg mt-4">
            Опишите, что именно вам нужно, и мы сразу покажем похожие проекты из портфолио.
          </p>
        </div>

        <AISearchPanel />
      </div>
    </section>

    <section id="quote-quiz" class="scroll-mt-28 pb-12 sm:pb-16">
      <div class="ui-container">
        <div class="overflow-hidden rounded-[2rem] bg-brand-brown p-5 text-white shadow-[0_28px_80px_rgba(23,33,29,0.2)] lg:p-8">
          <div class="grid grid-cols-1 gap-7 lg:grid-cols-[0.9fr_1.1fr] lg:items-center">
            <div class="lg:pr-4">
              <div class="mb-4 inline-flex items-center gap-3 rounded-full bg-white/8 px-4 py-2 text-brand-gold">
                <LucideShieldCheck :size="20" />
                <span class="font-semibold">Посчитаем без долгих созвонов</span>
              </div>
              <h2 class="font-serif text-3xl font-bold leading-tight sm:text-4xl">
                Ответьте на 4 вопроса и получите понятный следующий шаг по проекту
              </h2>
              <p class="mt-4 max-w-2xl leading-8 text-white/72">
                Подскажем реалистичный бюджет, сроки и важные детали до старта проекта.
              </p>

              <div class="mt-6 grid gap-3 text-sm font-semibold text-white/74 sm:grid-cols-3">
                <div class="rounded-2xl border border-white/10 bg-white/5 p-3">Размеры и техника</div>
                <div class="rounded-2xl border border-white/10 bg-white/5 p-3">Бюджет и сроки</div>
                <div class="rounded-2xl border border-white/10 bg-white/5 p-3">Замер по Крыму</div>
              </div>

              <a :href="PHONE_HREF" class="ui-button mt-6 border border-white/18 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
                Позвонить: {{ PHONE_DISPLAY }}
              </a>
            </div>

            <div class="rounded-[1.8rem] bg-white p-4 text-brand-brown shadow-2xl shadow-black/20 sm:p-5">
              <QuoteQuiz initial-project-type="Кухня с техникой" />
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

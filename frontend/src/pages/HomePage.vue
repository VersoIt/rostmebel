<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useProductStore } from '@/stores/products';
import AISearchPanel from '@/components/ai/AISearchPanel.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import {
  LucideArrowRight,
  LucideCheckCircle2,
  LucideHammer,
  LucideLayers,
  LucideMapPin,
  LucideMessageSquare,
  LucidePenTool,
  LucideRuler,
  LucideShieldCheck,
  LucideTruck,
} from 'lucide-vue-next';
import type { Product } from '@/types';

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
  { value: '15+', label: 'лет в корпусной мебели' },
  { value: '2 года', label: 'гарантия по договору' },
  { value: 'Крым', label: 'замер, доставка, монтаж' },
];

const services = [
  {
    icon: LucideRuler,
    title: 'Точный замер',
    text: 'Проверяем стены, углы, выводы воды и электрики до сметы, а не после монтажа.',
  },
  {
    icon: LucidePenTool,
    title: 'Проект и смета',
    text: 'Фиксируем материалы, фурнитуру, наполнение и сроки до запуска производства.',
  },
  {
    icon: LucideHammer,
    title: 'Свое производство',
    text: 'Делаем кухни, шкафы, гардеробные и мебель для всего объекта без лишних посредников.',
  },
  {
    icon: LucideTruck,
    title: 'Монтаж под ключ',
    text: 'Привозим, собираем, регулируем фасады и сдаем готовую мебель в рабочем состоянии.',
  },
];

const scopes = [
  {
    title: 'Кухня под размер',
    text: 'Планировка, фасады, столешница, фурнитура, подсветка, техника и монтаж в одной смете.',
    items: ['3D-проект', 'карта розеток', 'сборка и регулировка'],
  },
  {
    title: 'Системы хранения',
    text: 'Шкафы, гардеробные, прихожие и постирочные, рассчитанные под реальные вещи и проходы.',
    items: ['наполнение', 'раздвижные системы', 'защита стен и пола'],
  },
  {
    title: 'Мебель для объекта',
    text: 'Комплектуем квартиру, дом или коммерческое помещение в едином стиле и бюджете.',
    items: ['единый менеджер', 'график поставок', 'контроль качества'],
  },
];

const process = [
  { title: 'Бриф и планировка', text: 'Уточняем задачи, бюджет, технику и сценарии хранения.' },
  { title: 'Замер и проект', text: 'Снимаем размеры, готовим визуализацию и технические карты.' },
  { title: 'Производство', text: 'Согласованный проект уходит в цех, закупки и контроль комплектации.' },
  { title: 'Монтаж и сдача', text: 'Устанавливаем, проверяем геометрию, фурнитуру и чистовую готовность.' },
];

onMounted(async () => {
  await productStore.fetchProducts({
    limit: 6,
    sort_by: 'views_count',
    sort_order: 'desc',
    status: 'published',
  });
  hits.value = productStore.products;

  heroInterval = window.setInterval(() => {
    currentHeroIndex.value = (currentHeroIndex.value + 1) % heroImages.length;
  }, 5200);
});

onUnmounted(() => {
  if (heroInterval) window.clearInterval(heroInterval);
});
</script>

<template>
  <div class="bg-brand-cream min-h-screen text-brand-brown font-sans selection:bg-brand-gold selection:text-white">
    <section class="relative min-h-[calc(100svh-88px)] pt-28 pb-12 flex items-end overflow-hidden bg-neutral-950">
      <div class="absolute inset-0">
        <div
          v-for="(img, idx) in heroImages"
          :key="img"
          class="absolute inset-0 transition-opacity duration-1000"
          :style="{ opacity: currentHeroIndex === idx ? 1 : 0 }"
        >
          <img :src="img" class="h-full w-full object-cover" alt="Кухня и корпусная мебель РОСТ Мебель">
        </div>
        <div class="absolute inset-0 bg-gradient-to-r from-black/75 via-black/35 to-black/10"></div>
        <div class="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-t from-brand-cream to-transparent"></div>
      </div>

      <div class="relative z-10 w-full px-5">
        <div class="mx-auto max-w-7xl">
          <div class="max-w-3xl text-white">
            <div class="mb-5 inline-flex items-center gap-2 border border-white/20 bg-black/20 px-3 py-2 text-sm">
              <LucideMapPin :size="16" class="text-brand-gold" />
              Работаем по Крыму: замер, доставка, монтаж
            </div>

            <h1 class="font-serif text-4xl leading-tight md:text-6xl">
              Кухни и корпусная мебель под размер
            </h1>

            <p class="mt-5 max-w-2xl text-lg leading-8 text-white/82 md:text-xl">
              Проектируем, производим и устанавливаем мебель без разрыва между красивой картинкой и реальным монтажом. Смета, материалы и сроки фиксируются до запуска.
            </p>

            <div class="mt-8 flex flex-col gap-3 sm:flex-row">
              <router-link
                to="/catalog"
                class="inline-flex items-center justify-center gap-2 rounded-lg bg-brand-gold px-7 py-4 font-semibold text-white transition hover:bg-white hover:text-brand-brown"
              >
                Посмотреть проекты
                <LucideArrowRight :size="18" />
              </router-link>
              <router-link
                to="/contact"
                class="inline-flex items-center justify-center gap-2 rounded-lg border border-white/30 px-7 py-4 font-semibold text-white transition hover:bg-white hover:text-brand-brown"
              >
                Получить расчет
              </router-link>
            </div>
          </div>

          <div class="mt-10 grid max-w-3xl grid-cols-1 gap-3 sm:grid-cols-3">
            <div v-for="item in proof" :key="item.label" class="border border-white/14 bg-black/25 p-4 text-white backdrop-blur">
              <div class="font-serif text-3xl">{{ item.value }}</div>
              <div class="mt-1 text-sm text-white/70">{{ item.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="px-5 py-14">
      <div class="mx-auto grid max-w-7xl grid-cols-1 gap-4 md:grid-cols-4">
        <div v-for="service in services" :key="service.title" class="border border-brand-brown/10 bg-white p-6">
          <component :is="service.icon" class="mb-5 text-brand-gold" :size="28" />
          <h2 class="font-serif text-2xl">{{ service.title }}</h2>
          <p class="mt-3 leading-7 text-brand-brown/68">{{ service.text }}</p>
        </div>
      </div>
    </section>

    <section id="projects-grid" class="px-5 py-16">
      <div class="mx-auto max-w-7xl">
        <div class="mb-9 flex flex-col justify-between gap-5 md:flex-row md:items-end">
          <div>
            <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">Портфолио</p>
            <h2 class="font-serif text-4xl leading-tight md:text-5xl">Проекты, которые уже собраны и работают</h2>
          </div>
          <router-link to="/catalog" class="inline-flex items-center gap-2 font-semibold text-brand-gold hover:text-brand-brown">
            Все проекты
            <LucideArrowRight :size="18" />
          </router-link>
        </div>

        <div v-if="hits.length" class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          <ProductCard v-for="product in hits.slice(0, 6)" :key="product.id" :product="product" />
        </div>
        <div v-else class="border border-brand-brown/10 bg-white p-8 text-brand-brown/60">
          Проекты появятся после публикации в админке.
        </div>
      </div>
    </section>

    <section class="bg-white px-5 py-16">
      <div class="mx-auto grid max-w-7xl grid-cols-1 gap-10 lg:grid-cols-[0.85fr_1.15fr] lg:items-start">
        <div>
          <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">Под ключ</p>
          <h2 class="font-serif text-4xl leading-tight md:text-5xl">Не набор обещаний, а понятный состав работ</h2>
          <p class="mt-5 leading-8 text-brand-brown/68">
            Заранее собираем состав проекта, зоны ответственности и точки контроля, чтобы смета не расползалась после запуска.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <article v-for="scope in scopes" :key="scope.title" class="border border-brand-brown/10 bg-brand-cream p-6">
            <LucideLayers class="mb-5 text-brand-gold" :size="28" />
            <h3 class="font-serif text-2xl">{{ scope.title }}</h3>
            <p class="mt-3 min-h-24 leading-7 text-brand-brown/68">{{ scope.text }}</p>
            <ul class="mt-5 space-y-3">
              <li v-for="item in scope.items" :key="item" class="flex items-center gap-2 text-sm font-medium text-brand-brown/78">
                <LucideCheckCircle2 :size="16" class="text-brand-gold" />
                {{ item }}
              </li>
            </ul>
          </article>
        </div>
      </div>
    </section>

    <section class="px-5 py-16">
      <div class="mx-auto max-w-7xl">
        <div class="grid grid-cols-1 gap-10 lg:grid-cols-[1fr_1fr] lg:items-center">
          <div class="grid grid-cols-2 gap-3">
            <img src="/assets/images/hero-2.jpg" class="aspect-[4/5] w-full object-cover" alt="Проект кухни с деревянными фасадами">
            <img src="/assets/images/interior-1.jpg" class="mt-8 aspect-[4/5] w-full object-cover" alt="Интерьер с корпусной мебелью">
          </div>

          <div>
            <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">Процесс</p>
            <h2 class="font-serif text-4xl leading-tight md:text-5xl">Четыре шага, где каждый следующий опирается на проверенные данные</h2>
            <div class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
              <article v-for="(step, idx) in process" :key="step.title" class="border border-brand-brown/10 bg-white p-6">
                <div class="mb-4 flex h-10 w-10 items-center justify-center rounded-lg bg-brand-brown text-white">
                  {{ idx + 1 }}
                </div>
                <h3 class="font-serif text-2xl">{{ step.title }}</h3>
                <p class="mt-3 leading-7 text-brand-brown/68">{{ step.text }}</p>
              </article>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="ai-search" class="bg-white px-5 py-16">
      <div class="mx-auto mb-9 max-w-3xl text-center">
        <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">Подбор по описанию</p>
        <h2 class="font-serif text-4xl leading-tight md:text-5xl">Опишите задачу обычными словами</h2>
        <p class="mt-4 leading-8 text-brand-brown/68">
          Например: светлая кухня с техникой, высокий шкаф под хранение, бюджет до 250 000 ₽.
        </p>
      </div>
      <AISearchPanel />
    </section>

    <section class="px-5 py-16">
      <div class="mx-auto grid max-w-7xl grid-cols-1 gap-8 bg-brand-brown p-8 text-white md:grid-cols-[1fr_auto] md:items-center">
        <div>
          <div class="mb-4 flex items-center gap-3 text-brand-gold">
            <LucideShieldCheck :size="22" />
            <span class="font-semibold">Готовы посчитать ваш объект</span>
          </div>
          <h2 class="font-serif text-4xl leading-tight">Пришлите планировку или вызовите замерщика</h2>
          <p class="mt-4 max-w-2xl leading-8 text-white/70">
            Подскажем реалистичный бюджет, сроки и слабые места планировки до того, как вы вложитесь в материалы.
          </p>
        </div>
        <div class="flex flex-col gap-3 sm:flex-row md:flex-col">
          <router-link to="/contact" class="inline-flex items-center justify-center gap-2 rounded-lg bg-brand-gold px-7 py-4 font-semibold text-white transition hover:bg-white hover:text-brand-brown">
            <LucideMessageSquare :size="18" />
            Оставить заявку
          </router-link>
          <a href="tel:+79787631603" class="inline-flex items-center justify-center gap-2 rounded-lg border border-white/20 px-7 py-4 font-semibold text-white transition hover:bg-white hover:text-brand-brown">
            Позвонить
          </a>
        </div>
      </div>
    </section>
  </div>
</template>

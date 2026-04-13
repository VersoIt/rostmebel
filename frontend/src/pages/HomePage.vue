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
  { value: '15+', label: 'лет делаем корпусную мебель' },
  { value: '2 года', label: 'гарантия по договору' },
  { value: 'Крым', label: 'замер, доставка и монтаж' },
];

const services = [
  {
    icon: LucideRuler,
    title: 'Замер без догадок',
    text: 'Проверяем стены, углы, выводы воды, розетки и технику до сметы, чтобы проект не пришлось переделывать на монтаже.',
  },
  {
    icon: LucidePenTool,
    title: 'Проект и смета',
    text: 'Фиксируем материалы, фурнитуру, наполнение, сроки и состав работ до запуска производства.',
  },
  {
    icon: LucideHammer,
    title: 'Собственное производство',
    text: 'Делаем кухни, шкафы, гардеробные и мебель для всего объекта без лишних посредников.',
  },
  {
    icon: LucideTruck,
    title: 'Монтаж под ключ',
    text: 'Привозим, собираем, регулируем фасады и сдаем мебель в рабочем состоянии.',
  },
];

const scopes = [
  {
    title: 'Кухни',
    text: 'Планировка, фасады, столешница, фурнитура, подсветка, техника и монтаж в одной смете.',
    items: ['3D-проект', 'карта розеток', 'сборка и регулировка'],
  },
  {
    title: 'Системы хранения',
    text: 'Шкафы, гардеробные, прихожие и постирочные под реальные вещи, проходы и привычки семьи.',
    items: ['наполнение', 'раздвижные системы', 'защита стен и пола'],
  },
  {
    title: 'Мебель для объекта',
    text: 'Комплектуем квартиру, дом или коммерческое помещение в едином стиле и бюджете.',
    items: ['единый менеджер', 'график поставок', 'контроль качества'],
  },
];

const process = [
  { title: 'Бриф', text: 'Уточняем задачу, бюджет, технику, сроки и сценарии хранения.' },
  { title: 'Замер', text: 'Снимаем размеры, проверяем ограничения и готовим техническую основу.' },
  { title: 'Производство', text: 'Согласованный проект уходит в цех, закупку и контроль комплектации.' },
  { title: 'Монтаж', text: 'Устанавливаем, регулируем фурнитуру и сдаем готовый результат.' },
];

onMounted(async () => {
  setJsonLd('schema-business', buildBusinessSchema());

  await productStore.fetchProducts({
    limit: 6,
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
    <section class="relative flex min-h-[calc(100svh-96px)] items-end overflow-hidden bg-neutral-950 pb-12 pt-28">
      <div class="absolute inset-0">
        <div
          v-for="(img, idx) in heroImages"
          :key="img"
          class="absolute inset-0 transition-opacity duration-700"
          :style="{ opacity: currentHeroIndex === idx ? 1 : 0 }"
        >
          <img :src="img" class="h-full w-full object-cover" alt="Кухня и корпусная мебель РОСТ Мебель">
        </div>
        <div class="absolute inset-0 bg-gradient-to-r from-black/78 via-black/42 to-black/10"></div>
        <div class="absolute inset-x-0 bottom-0 h-28 bg-gradient-to-t from-brand-cream to-transparent"></div>
      </div>

      <div class="ui-container relative z-10">
        <div class="max-w-3xl text-white motion-fade-up">
          <div class="mb-5 inline-flex items-center gap-2 rounded-lg border border-white/20 bg-black/22 px-3 py-2 text-sm">
            <LucideMapPin :size="16" class="text-brand-gold" />
            Работаем по Крыму: замер, доставка, монтаж
          </div>

          <h1 class="font-serif text-4xl font-bold leading-tight sm:text-5xl lg:text-6xl">
            Кухни и корпусная мебель по размеру
          </h1>

          <p class="mt-5 max-w-2xl text-lg leading-8 text-white/82">
            Проектируем, производим и устанавливаем мебель без разрыва между красивой картинкой и реальным монтажом. Смета, материалы и сроки фиксируются до запуска.
          </p>

          <div class="mt-8 flex flex-col gap-3 sm:flex-row">
            <router-link to="/catalog" class="ui-button ui-button-accent">
              Посмотреть проекты
              <LucideArrowRight :size="18" />
            </router-link>
            <router-link to="/contact" class="ui-button border border-white/30 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
              Получить расчет
            </router-link>
          </div>
        </div>

        <div class="mt-10 grid max-w-3xl grid-cols-1 gap-3 sm:grid-cols-3">
          <div v-for="item in proof" :key="item.label" class="rounded-lg border border-white/14 bg-black/24 p-4 text-white backdrop-blur">
            <div class="font-serif text-3xl">{{ item.value }}</div>
            <div class="mt-1 text-sm text-white/70">{{ item.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <section class="ui-section-tight">
      <div class="ui-container grid grid-cols-1 gap-4 md:grid-cols-4">
        <article v-for="service in services" :key="service.title" class="ui-card ui-card-hover p-5 motion-fade-up">
          <component :is="service.icon" class="mb-5 text-brand-gold" :size="28" />
          <h2 class="ui-title-md">{{ service.title }}</h2>
          <p class="ui-copy mt-3">{{ service.text }}</p>
        </article>
      </div>
    </section>

    <section id="projects-grid" class="ui-section">
      <div class="ui-container">
        <div class="mb-9 flex flex-col justify-between gap-5 md:flex-row md:items-end">
          <div>
            <p class="ui-eyebrow mb-3">Портфолио</p>
            <h2 class="ui-title-lg">Проекты, которые уже собраны и работают</h2>
          </div>
          <router-link to="/catalog" class="ui-button ui-button-secondary">
            Все проекты
            <LucideArrowRight :size="18" />
          </router-link>
        </div>

        <div v-if="hits.length" class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          <ProductCard v-for="product in hits.slice(0, 6)" :key="product.id" :product="product" />
        </div>
        <div v-else class="ui-empty">
          Проекты появятся после публикации в админке.
        </div>
      </div>
    </section>

    <section class="bg-white">
      <div class="ui-container ui-section grid grid-cols-1 gap-10 lg:grid-cols-[0.85fr_1.15fr] lg:items-start">
        <div>
          <p class="ui-eyebrow mb-3">Под ключ</p>
          <h2 class="ui-title-lg">Понятный состав работ вместо общих обещаний</h2>
          <p class="ui-copy-lg mt-5">
            Сразу собираем состав проекта, зоны ответственности и точки контроля, чтобы смета не расползалась после запуска.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <article v-for="scope in scopes" :key="scope.title" class="ui-card p-5">
            <LucideLayers class="mb-5 text-brand-gold" :size="28" />
            <h3 class="ui-title-md">{{ scope.title }}</h3>
            <p class="ui-copy mt-3">{{ scope.text }}</p>
            <ul class="mt-5 space-y-3">
              <li v-for="item in scope.items" :key="item" class="flex items-center gap-2 text-sm font-semibold text-brand-brown/78">
                <LucideCheckCircle2 :size="16" class="text-brand-gold" />
                {{ item }}
              </li>
            </ul>
          </article>
        </div>
      </div>
    </section>

    <section class="ui-section">
      <div class="ui-container grid grid-cols-1 gap-10 lg:grid-cols-[1fr_1fr] lg:items-center">
        <div class="grid grid-cols-2 gap-3">
          <img src="/assets/images/hero-2.jpg" class="aspect-[4/5] w-full rounded-lg object-cover" alt="Проект кухни с деревянными фасадами">
          <img src="/assets/images/interior-1.jpg" class="mt-8 aspect-[4/5] w-full rounded-lg object-cover" alt="Интерьер с корпусной мебелью">
        </div>

        <div>
          <p class="ui-eyebrow mb-3">Процесс</p>
          <h2 class="ui-title-lg">Четыре шага, где каждый следующий опирается на проверенные данные</h2>
          <div class="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
            <article v-for="(step, idx) in process" :key="step.title" class="ui-card p-5">
              <div class="mb-4 flex h-10 w-10 items-center justify-center rounded-lg bg-brand-brown font-bold text-white">
                {{ idx + 1 }}
              </div>
              <h3 class="ui-title-md">{{ step.title }}</h3>
              <p class="ui-copy mt-3">{{ step.text }}</p>
            </article>
          </div>
        </div>
      </div>
    </section>

    <section id="ai-search" class="bg-white">
      <div class="ui-container ui-section">
        <div class="mx-auto mb-9 max-w-3xl text-center">
          <p class="ui-eyebrow mb-3">Подбор по описанию</p>
          <h2 class="ui-title-lg">Опишите задачу обычными словами</h2>
          <p class="ui-copy-lg mt-4">
            Например: светлая кухня с техникой, высокий шкаф под хранение, бюджет до 250 000 ₽.
          </p>
        </div>
        <AISearchPanel />
      </div>
    </section>

    <section class="ui-section">
      <div class="ui-container">
        <div class="grid grid-cols-1 gap-8 rounded-lg bg-brand-brown p-6 text-white md:grid-cols-[1fr_auto] md:items-center lg:p-8">
          <div>
            <div class="mb-4 flex items-center gap-3 text-brand-gold">
              <LucideShieldCheck :size="22" />
              <span class="font-semibold">Посчитаем ваш объект</span>
            </div>
            <h2 class="font-serif text-3xl font-bold leading-tight sm:text-4xl">Пришлите планировку или вызовите замерщика</h2>
            <p class="mt-4 max-w-2xl leading-8 text-white/72">
              Подскажем реалистичный бюджет, сроки и слабые места планировки до того, как вы вложитесь в материалы.
            </p>
          </div>
          <div class="flex flex-col gap-3 sm:flex-row md:flex-col">
            <router-link to="/contact" class="ui-button ui-button-accent">
              <LucideMessageSquare :size="18" />
              Оставить заявку
            </router-link>
            <a href="tel:+79787631603" class="ui-button border border-white/20 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
              Позвонить
            </a>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

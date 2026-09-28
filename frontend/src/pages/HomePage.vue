<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import QuoteQuiz from '@/components/order/QuoteQuiz.vue';
import {
  LucideArrowRight,
  LucideClock,
  LucideMapPin,
  LucideMessageSquare,
  LucidePhoneCall,
  LucideRuler,
  LucideShieldCheck,
  LucideWrench,
} from 'lucide-vue-next';
import type { Product } from '@/types';
import { buildBusinessSchema, buildWebsiteSchema, removeJsonLd, setJsonLd } from '@/utils/seo';
import { MESSENGER_LINKS, PHONE_DISPLAY, PHONE_HREF } from '@/constants/contacts';

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
  { value: 'Крым', label: 'замер и монтаж' },
];

const stages = [
  {
    title: '1. Заявка',
    text: 'Фото, размеры или просто идея: кухня, шкаф, прихожая, гардеробная.',
    icon: LucideRuler,
  },
  {
    title: '2. Проект',
    text: 'Подберем материалы, компоновку и согласуем детали.',
    icon: LucideMessageSquare,
  },
  {
    title: '3. Монтаж',
    text: 'Изготовим, доставим и установим мебель на объекте.',
    icon: LucideWrench,
  },
];

const heroBullets = [
  'кухни, шкафы, гардеробные',
  'замер, производство, монтаж',
  'гарантия 2 года по договору',
];

onMounted(async () => {
  setJsonLd('schema-website', buildWebsiteSchema());
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
        <div class="grid gap-8 lg:grid-cols-[minmax(0,1fr)_410px] lg:items-center">
          <div class="max-w-4xl text-white motion-fade-up">
            <div class="mb-5 inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-4 py-2 text-sm backdrop-blur">
              <LucideMapPin :size="16" class="text-brand-gold" />
              Работаем по Крыму: замер, доставка, монтаж
            </div>

            <h1 class="max-w-3xl font-serif text-4xl font-bold leading-[1.02] sm:text-5xl lg:text-[4.35rem]">
              Кухни и мебель на заказ в Крыму
            </h1>

            <p class="mt-6 max-w-2xl text-lg leading-8 text-white/80">
              Проектируем, производим и устанавливаем мебель по вашим размерам.
            </p>

            <div class="mt-7 grid max-w-2xl gap-2 sm:grid-cols-3">
              <div
                v-for="item in heroBullets"
                :key="item"
                class="rounded-2xl border border-white/12 bg-white/8 px-4 py-3 text-sm font-semibold leading-5 text-white/78 backdrop-blur"
              >
                {{ item }}
              </div>
            </div>

            <div class="mt-8 flex flex-col gap-3 sm:flex-row">
              <a href="#quote-quiz" class="ui-button ui-button-accent min-h-14 text-base">
                Получить расчет
                <LucideArrowRight :size="18" />
              </a>
              <router-link to="/catalog" class="ui-button border border-white/20 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
                Посмотреть проекты
              </router-link>
            </div>

            <div class="mt-4 flex flex-wrap items-center gap-2">
              <span class="mr-1 text-xs font-semibold uppercase tracking-widest text-white/42">Написать напрямую</span>
              <a :href="MESSENGER_LINKS.whatsapp" target="_blank" rel="noopener noreferrer" class="inline-flex min-h-9 items-center gap-1.5 rounded-lg border border-white/14 bg-white/10 px-3 py-1.5 text-xs font-bold text-white/86 transition hover:-translate-y-0.5 hover:bg-white hover:text-brand-brown">
                <LucideMessageSquare :size="17" />
                WhatsApp
              </a>
              <a :href="MESSENGER_LINKS.telegram" target="_blank" rel="noopener noreferrer" class="inline-flex min-h-9 items-center gap-1.5 rounded-lg border border-white/14 bg-white/10 px-3 py-1.5 text-xs font-bold text-white/86 transition hover:-translate-y-0.5 hover:bg-white hover:text-brand-brown">
                <LucideMessageSquare :size="17" />
                Telegram
              </a>
              <a :href="PHONE_HREF" class="inline-flex min-h-9 items-center gap-1.5 rounded-lg border border-white/14 bg-white/10 px-3 py-1.5 text-xs font-bold text-white/86 transition hover:-translate-y-0.5 hover:bg-white hover:text-brand-brown">
                <LucidePhoneCall :size="17" />
                {{ PHONE_DISPLAY }}
              </a>
            </div>

            <div class="mt-7 flex max-w-3xl flex-wrap gap-x-7 gap-y-3 border-t border-white/14 pt-5">
              <div
                v-for="item in proof"
                :key="item.label"
                class="flex items-baseline gap-2"
              >
                <div class="font-serif text-2xl leading-none text-white">{{ item.value }}</div>
                <div class="text-sm leading-5 text-white/58">{{ item.label }}</div>
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

          <aside id="quote-quiz" class="motion-scale-in scroll-mt-28 rounded-[2rem] border border-white/14 bg-white p-4 text-brand-brown shadow-[0_28px_90px_rgba(0,0,0,0.28)] sm:p-5">
            <div class="mb-4 rounded-2xl bg-brand-gray px-4 py-3">
              <div class="flex items-center gap-2 text-[11px] font-black uppercase tracking-widest text-brand-gold">
                <LucideClock :size="16" />
                Быстрый старт
              </div>
              <h2 class="mt-2 font-serif text-2xl font-bold leading-tight text-brand-brown">
                Заявка на расчет
              </h2>
              <p class="mt-2 text-sm leading-6 text-brand-brown/62">
                Ответьте на 4 вопроса, чтобы мы поняли задачу.
              </p>
            </div>
            <QuoteQuiz initial-project-type="Кухня с техникой" />
          </aside>
        </div>
      </div>
    </section>

    <section id="projects-grid" class="ui-section pt-8">
      <div class="ui-container">
        <div class="ui-surface overflow-hidden p-5 sm:p-7 lg:p-8">
          <div class="mb-8 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
            <div>
              <p class="ui-eyebrow mb-3">Портфолио</p>
              <h2 class="ui-title-lg">Проекты</h2>
              <p class="ui-copy mt-4 max-w-2xl">
                Кухни, шкафы и другие решения, которые уже установлены у клиентов.
              </p>
            </div>
            <router-link to="/catalog" class="ui-button ui-button-secondary">
              Все проекты
              <LucideArrowRight :size="18" />
            </router-link>
          </div>

          <div v-if="hits.length" class="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3">
            <ProductCard v-for="product in hits" :key="product.id" :product="product" />
          </div>
          <div v-else class="ui-empty">
            Проекты появятся после публикации в админке.
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

    <section class="pb-12 sm:pb-16">
      <div class="ui-container">
        <div class="overflow-hidden rounded-[2rem] bg-brand-brown p-5 text-white shadow-[0_28px_80px_rgba(23,33,29,0.2)] lg:p-8">
          <div class="grid grid-cols-1 gap-7 lg:grid-cols-[minmax(0,1fr)_360px] lg:items-center">
            <div>
              <div class="mb-4 inline-flex items-center gap-3 rounded-full bg-white/8 px-4 py-2 text-brand-gold">
                <LucideShieldCheck :size="20" />
                <span class="font-semibold">Связаться</span>
              </div>
              <h2 class="font-serif text-3xl font-bold leading-tight sm:text-4xl">
                Обсудим ваш проект
              </h2>
              <p class="mt-4 max-w-2xl leading-8 text-white/72">
                Напишите в мессенджер, позвоните или заполните короткую форму выше.
              </p>

              <div class="mt-6 grid gap-3 text-sm font-semibold text-white/74 sm:grid-cols-3">
                <div class="rounded-2xl border border-white/10 bg-white/5 p-3">Фото помещения</div>
                <div class="rounded-2xl border border-white/10 bg-white/5 p-3">Размеры или идея</div>
                <div class="rounded-2xl border border-white/10 bg-white/5 p-3">Удобный способ связи</div>
              </div>

              <div class="mt-6 flex flex-wrap gap-2">
                <a :href="MESSENGER_LINKS.whatsapp" target="_blank" rel="noopener noreferrer" class="ui-button bg-white text-brand-brown hover:bg-brand-gold hover:text-white">
                  WhatsApp
                </a>
                <a :href="MESSENGER_LINKS.telegram" target="_blank" rel="noopener noreferrer" class="ui-button border border-white/18 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
                  Telegram
                </a>
                <a :href="PHONE_HREF" class="ui-button border border-white/18 bg-white/8 text-white hover:bg-white hover:text-brand-brown">
                  Позвонить: {{ PHONE_DISPLAY }}
                </a>
              </div>
            </div>

            <div class="rounded-[1.8rem] border border-white/10 bg-white/8 p-5">
              <div class="font-serif text-3xl font-bold leading-tight text-white">Начать расчет</div>
              <p class="mt-3 leading-7 text-white/70">
                Достаточно описать мебель и указать удобный способ связи.
              </p>
              <a href="#quote-quiz" class="ui-button mt-5 bg-white text-brand-brown hover:bg-brand-gold hover:text-white">
                Заполнить короткую форму
                <LucideArrowRight :size="18" />
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

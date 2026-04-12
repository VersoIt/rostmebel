<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProductStore } from '@/stores/products';
import {
  LucideArrowRight,
  LucideCheckCircle,
  LucideChevronLeft,
  LucideMessageSquare,
  LucideSearch,
  LucideShieldCheck,
  LucideTruck,
  LucideX,
} from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import ReviewForm from '@/components/catalog/ReviewForm.vue';
import ReviewList from '@/components/catalog/ReviewList.vue';
import type { Product } from '@/types';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';

const route = useRoute();
const router = useRouter();
const productStore = useProductStore();
const product = ref<Product | null>(null);
const relatedProjects = ref<Product[]>([]);
const activeImage = ref('');
const isOrderModalOpen = ref(false);
const isReviewModalOpen = ref(false);
const isLightboxOpen = ref(false);
const reviewListRef = ref<any>(null);

const handleImageError = (event: Event) => {
  (event.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
};

const openLightbox = (url: string) => {
  activeImage.value = url;
  isLightboxOpen.value = true;
};

const updateSchema = (item: Product) => {
  const schema = {
    '@context': 'https://schema.org/',
    '@type': 'Product',
    name: item.name,
    image: item.images.map((image) => image.url),
    description: item.description,
    offers: {
      '@type': 'Offer',
      url: window.location.href,
      priceCurrency: 'RUB',
      price: item.price,
      availability: 'https://schema.org/InStock',
      seller: { '@type': 'Organization', name: 'РОСТ Мебель' },
    },
  };

  const scriptId = 'schema-product';
  let script = document.getElementById(scriptId) as HTMLScriptElement;
  if (!script) {
    script = document.createElement('script');
    script.id = scriptId;
    script.type = 'application/ld+json';
    document.head.appendChild(script);
  }
  script.textContent = JSON.stringify(schema);
};

const loadProjectData = async () => {
  const id = route.params.id as string;
  const loadedProduct = await productStore.fetchProduct(id);

  if (loadedProduct) {
    product.value = loadedProduct;
    activeImage.value = loadedProduct.images[0]?.url || PLACEHOLDER_IMAGE;
    updateSchema(loadedProduct);
    document.title = `${loadedProduct.name} — РОСТ Мебель`;

    await productStore.fetchProducts({
      project_category_id: loadedProduct.project_category_id,
      limit: 4,
      status: 'published',
    });
    relatedProjects.value = productStore.products.filter((item) => item.id !== loadedProduct.id).slice(0, 3);
  }

  window.scrollTo({ top: 0, behavior: 'smooth' });
};

watch(() => route.params.id, loadProjectData);
onMounted(loadProjectData);

onUnmounted(() => {
  document.getElementById('schema-product')?.remove();
});

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(price);
};

const handleReviewSuccess = () => {
  isReviewModalOpen.value = false;
  reviewListRef.value?.refresh();
};
</script>

<template>
  <div v-if="product" class="min-h-screen bg-white">
    <section class="bg-brand-cream pt-28">
      <div class="ui-container ui-section-tight">
        <button type="button" class="mb-6 inline-flex items-center gap-2 text-sm font-bold text-brand-brown/50 transition-colors hover:text-brand-gold" @click="router.push('/catalog')">
          <LucideChevronLeft :size="17" />
          Назад к проектам
        </button>
        <h1 class="ui-title-xl">{{ product.name }}</h1>
        <div class="mt-4 flex flex-wrap items-center gap-2">
          <span class="ui-status bg-brand-gold/10 text-brand-gold ring-brand-gold/20">
            {{ productStore.categories.find((category) => category.id === product?.project_category_id)?.name || 'Проект' }}
          </span>
        </div>
      </div>
    </section>

    <section class="ui-container ui-section grid grid-cols-1 gap-10 lg:grid-cols-2 lg:gap-14">
      <div class="space-y-5">
        <button
          type="button"
          class="group relative aspect-square w-full overflow-hidden rounded-lg bg-brand-gray"
          @click="openLightbox(activeImage)"
        >
          <img :src="activeImage" :alt="product.name" class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-[1.035]" @error="handleImageError">
          <span class="absolute inset-0 flex items-center justify-center bg-black/0 transition-colors group-hover:bg-black/12">
            <LucideSearch :size="40" class="text-white opacity-0 transition-opacity group-hover:opacity-100" />
          </span>
        </button>

        <div class="flex gap-3 overflow-x-auto pb-2 no-scrollbar">
          <button
            v-for="image in product.images"
            :key="image.url"
            type="button"
            :class="[
              'h-20 w-20 shrink-0 overflow-hidden rounded-lg border-2 transition-colors',
              activeImage === image.url ? 'border-brand-gold' : 'border-transparent opacity-70 hover:opacity-100'
            ]"
            @click="activeImage = image.url"
          >
            <img :src="image.url" class="h-full w-full object-cover" alt="" @error="handleImageError">
          </button>
        </div>
      </div>

      <div class="flex flex-col">
        <div class="mb-8">
          <div class="mb-6 flex flex-wrap items-end gap-6">
            <div>
              <div class="ui-label-compact">Бюджет реализации</div>
              <div class="font-serif text-4xl font-bold text-brand-gold">{{ formatPrice(product.price) }}</div>
            </div>
            <div v-if="product.price_old">
              <div class="ui-label-compact">Ориентир</div>
              <div class="text-xl text-brand-brown/25 line-through">{{ formatPrice(product.price_old) }}</div>
            </div>
          </div>
          <p class="ui-copy-lg">{{ product.description }}</p>
        </div>

        <div class="ui-card-muted mb-8 p-5 sm:p-6">
          <h2 class="ui-title-md mb-6">Детали проекта</h2>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div v-for="(value, key) in product.specs" :key="key" class="border-b border-brand-brown/10 pb-3">
              <div class="ui-label-compact">{{ key }}</div>
              <div class="font-bold text-brand-brown">{{ value }}</div>
            </div>
          </div>
        </div>

        <div class="mt-auto space-y-6">
          <button type="button" class="ui-button ui-button-primary w-full min-h-14 text-base" @click="isOrderModalOpen = true">
            Рассчитать похожий проект
          </button>

          <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
            <div
              v-for="item in [
                { icon: LucideShieldCheck, text: 'Гарантия 2 года' },
                { icon: LucideTruck, text: 'Монтаж в Крыму' },
                { icon: LucideCheckCircle, text: 'Контроль сборки' }
              ]"
              :key="item.text"
              class="flex items-center gap-2 text-xs font-black uppercase tracking-widest text-brand-brown/45"
            >
              <component :is="item.icon" class="text-brand-gold" :size="18" />
              {{ item.text }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="ui-container ui-section border-t border-brand-brown/10">
      <div class="mb-10 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="ui-eyebrow mb-3">Отзывы</p>
          <h2 class="ui-title-lg">Что говорят клиенты</h2>
          <p class="ui-copy mt-3">Публикуем отзывы после модерации и проверки заказа.</p>
        </div>
        <button type="button" class="ui-button ui-button-secondary" @click="isReviewModalOpen = true">
          <LucideMessageSquare :size="19" />
          Оставить отзыв
        </button>
      </div>

      <ReviewList ref="reviewListRef" :project-id="product.id" />
    </section>

    <section v-if="relatedProjects.length" class="ui-container ui-section border-t border-brand-brown/10">
      <div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <h2 class="ui-title-lg">Похожие проекты</h2>
        <router-link to="/catalog" class="ui-button ui-button-secondary">
          Смотреть все
          <LucideArrowRight :size="18" />
        </router-link>
      </div>
      <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
        <ProductCard v-for="related in relatedProjects" :key="related.id" :product="related" />
      </div>
    </section>

    <Teleport to="body">
      <transition name="fade">
        <div v-if="isOrderModalOpen" class="ui-modal-backdrop" @click.self="isOrderModalOpen = false">
          <section class="ui-modal-panel max-w-lg p-5 sm:p-8">
            <button type="button" class="absolute right-4 top-4 rounded-lg p-2 text-brand-brown/35 transition-colors hover:bg-brand-gray hover:text-brand-brown" @click="isOrderModalOpen = false">
              <LucideX :size="24" />
            </button>
            <h2 class="ui-title-md mb-2">Заявка на расчет</h2>
            <p class="ui-copy mb-6">Обсудим похожий проект и подскажем реалистичный бюджет.</p>
            <OrderForm :project-id="product.id" @success="isOrderModalOpen = false" />
          </section>
        </div>
      </transition>
    </Teleport>

    <Teleport to="body">
      <transition name="fade">
        <div v-if="isReviewModalOpen" class="ui-modal-backdrop" @click.self="isReviewModalOpen = false">
          <section class="ui-modal-panel max-w-2xl">
            <button type="button" class="absolute right-4 top-4 z-10 rounded-lg p-2 text-brand-brown/35 transition-colors hover:bg-brand-gray hover:text-brand-brown" @click="isReviewModalOpen = false">
              <LucideX :size="24" />
            </button>
            <ReviewForm :project-id="product.id" @success="handleReviewSuccess" />
          </section>
        </div>
      </transition>
    </Teleport>

    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLightboxOpen" class="ui-modal-backdrop" @click="isLightboxOpen = false">
          <button type="button" class="absolute right-5 top-5 rounded-lg bg-white/10 p-3 text-white transition-colors hover:bg-white hover:text-brand-brown">
            <LucideX :size="28" />
          </button>
          <img :src="activeImage" class="z-10 max-h-full max-w-full rounded-lg object-contain shadow-2xl" alt="">
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

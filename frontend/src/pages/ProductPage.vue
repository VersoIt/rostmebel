<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProductStore } from '@/stores/products';
import {
  LucideArrowRight,
  LucideCheckCircle,
  LucideChevronLeft,
  LucideChevronRight,
  LucideMessageSquare,
  LucideSearch,
  LucideShieldCheck,
  LucideTruck,
  LucideX,
} from 'lucide-vue-next';
import QuoteQuiz from '@/components/order/QuoteQuiz.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import ReviewForm from '@/components/catalog/ReviewForm.vue';
import ReviewList from '@/components/catalog/ReviewList.vue';
import type { Product } from '@/types';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';
import { buildImageVariantUrl, buildResponsiveImageSet, preloadResponsiveImage } from '@/utils/images';
import { absoluteUrl, compactDescription, removeJsonLd, setJsonLd, setPageSeo } from '@/utils/seo';

const route = useRoute();
const router = useRouter();
const productStore = useProductStore();
const product = ref<Product | null>(null);
const relatedProjects = ref<Product[]>([]);
const activeImage = ref('');
const displayedMainImage = ref(PLACEHOLDER_IMAGE);
const isOrderModalOpen = ref(false);
const isReviewModalOpen = ref(false);
const isLightboxOpen = ref(false);
const isMainImageLoading = ref(false);
const isLightboxImageLoading = ref(false);
const pendingLightboxIndex = ref<number | null>(null);
const reviewListRef = ref<any>(null);
const lightboxTouchStartX = ref(0);
const lightboxTouchStartY = ref(0);
let mainImageRequestId = 0;
let lightboxNavigationRequestId = 0;
let lightboxBackgroundPreloadTimer: number | undefined;

const productImages = computed(() => product.value?.images || []);
const hasMultipleImages = computed(() => productImages.value.length > 1);
const mainImageSource = computed(() => getMainImageSource(displayedMainImage.value));
const lightboxImageSource = computed(() => getLightboxImageSource(activeImage.value || PLACEHOLDER_IMAGE));
const currentImageIndex = computed(() => {
  if (!productImages.value.length) return -1;

  const foundIndex = productImages.value.findIndex((image) => image.url === activeImage.value);
  return foundIndex >= 0 ? foundIndex : 0;
});
const visibleLightboxIndex = computed(() => pendingLightboxIndex.value ?? currentImageIndex.value);

const quoteProjectType = computed(() => {
  if (!product.value) return 'Пока не знаю';

  const categorySlug = productStore.categories.find((category) => category.id === product.value?.project_category_id)?.slug || '';

  const haystack = [
    product.value.name,
    product.value.description,
    product.value.ai_tags,
    ...Object.values(product.value.specs || {}),
  ].join(' ').toLowerCase();

  if (categorySlug === 'kitchens' || haystack.includes('кух')) {
    return haystack.includes('техник') ? 'Кухня с техникой' : 'Кухня';
  }

  if (categorySlug === 'wardrobes' || categorySlug === 'dressing-rooms' || haystack.includes('шкаф') || haystack.includes('гардероб')) {
    return 'Шкаф или гардеробная';
  }

  if (categorySlug === 'hallways' || haystack.includes('прихож')) {
    return 'Прихожая';
  }

  if (categorySlug === 'tables' || haystack.includes('стол') || haystack.includes('рабоч')) {
    return 'Стол или рабочая зона';
  }

  if (categorySlug === 'living-rooms' || haystack.includes('гостин') || haystack.includes('тв-тумб') || haystack.includes('тумб')) {
    return 'ТВ-тумба или гостиная';
  }

  if (categorySlug === 'children-rooms' || haystack.includes('детск')) {
    return 'Детская';
  }

  if (categorySlug === 'commercial-furniture' || haystack.includes('коммер') || haystack.includes('офис') || haystack.includes('салон')) {
    return 'Коммерческий объект';
  }

  return 'Пока не знаю';
});

const handleImageError = (event: Event) => {
  (event.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
};

const thumbnailImageUrl = (url: string) => buildImageVariantUrl(url, 180, 72);

const getMainImageSource = (url: string) =>
  buildResponsiveImageSet(
    url || PLACEHOLDER_IMAGE,
    [640, 960, 1280, 1600],
    '(min-width: 1024px) 50vw, 100vw',
    82,
  );

const getLightboxImageSource = (url: string) =>
  buildResponsiveImageSet(
    url || PLACEHOLDER_IMAGE,
    [960, 1440, 1920, 2560],
    '100vw',
    86,
  );

const normalizeImageIndex = (index: number) => {
  if (!productImages.value.length) return -1;
  return ((index % productImages.value.length) + productImages.value.length) % productImages.value.length;
};

const imageUrlByIndex = (index: number) => {
  const normalizedIndex = normalizeImageIndex(index);
  if (normalizedIndex < 0) {
    return '';
  }
  return productImages.value[normalizedIndex]?.url || '';
};

const clearLightboxBackgroundPreloadTimer = () => {
  if (lightboxBackgroundPreloadTimer !== undefined) {
    window.clearTimeout(lightboxBackgroundPreloadTimer);
    lightboxBackgroundPreloadTimer = undefined;
  }
};

const syncDisplayedMainImage = async (nextImage: string, immediate = false) => {
  const normalizedImage = nextImage || PLACEHOLDER_IMAGE;
  const requestId = ++mainImageRequestId;

  if (immediate || displayedMainImage.value === normalizedImage) {
    displayedMainImage.value = normalizedImage;
    isMainImageLoading.value = false;
    void preloadResponsiveImage(getMainImageSource(normalizedImage));
    return;
  }

  isMainImageLoading.value = true;
  await preloadResponsiveImage(getMainImageSource(normalizedImage));

  if (requestId !== mainImageRequestId) {
    return;
  }

  displayedMainImage.value = normalizedImage;
  isMainImageLoading.value = false;
};

const preloadLightboxImage = (url: string) => preloadResponsiveImage(getLightboxImageSource(url || PLACEHOLDER_IMAGE));

const buildLightboxPreloadOrder = (centerIndex: number) => {
  if (!productImages.value.length) {
    return [] as number[];
  }

  const visited = new Set<number>();
  const ordered: number[] = [];

  for (let offset = 0; ordered.length < productImages.value.length; offset += 1) {
    const nextIndex = normalizeImageIndex(centerIndex + offset);
    if (nextIndex >= 0 && !visited.has(nextIndex)) {
      visited.add(nextIndex);
      ordered.push(nextIndex);
    }

    if (offset === 0) {
      continue;
    }

    const previousIndex = normalizeImageIndex(centerIndex - offset);
    if (previousIndex >= 0 && !visited.has(previousIndex)) {
      visited.add(previousIndex);
      ordered.push(previousIndex);
    }
  }

  return ordered;
};

const scheduleLightboxPreload = (centerIndex: number) => {
  if (!productImages.value.length) {
    return;
  }

  clearLightboxBackgroundPreloadTimer();

  const preloadOrder = buildLightboxPreloadOrder(centerIndex);
  preloadOrder.slice(0, 3).forEach((index) => {
    const url = imageUrlByIndex(index);
    if (url) {
      void preloadLightboxImage(url);
    }
  });

  const remainingIndexes = preloadOrder.slice(3);
  if (!remainingIndexes.length) {
    return;
  }

  lightboxBackgroundPreloadTimer = window.setTimeout(() => {
    void (async () => {
      for (const index of remainingIndexes) {
        const url = imageUrlByIndex(index);
        if (url) {
          await preloadLightboxImage(url);
        }
      }
    })();
    lightboxBackgroundPreloadTimer = undefined;
  }, 90);
};

const setActiveImageByIndex = (index: number) => {
  const nextUrl = imageUrlByIndex(index);
  activeImage.value = nextUrl || PLACEHOLDER_IMAGE;
};

const navigateLightboxToIndex = async (index: number) => {
  const normalizedIndex = normalizeImageIndex(index);
  const nextUrl = imageUrlByIndex(normalizedIndex);

  if (normalizedIndex < 0 || !nextUrl) {
    return;
  }

  pendingLightboxIndex.value = normalizedIndex;
  isLightboxImageLoading.value = true;
  scheduleLightboxPreload(normalizedIndex);

  const requestId = ++lightboxNavigationRequestId;
  await preloadLightboxImage(nextUrl);

  if (requestId !== lightboxNavigationRequestId) {
    return;
  }

  activeImage.value = nextUrl;
  pendingLightboxIndex.value = null;
  isLightboxImageLoading.value = false;
};

const openLightbox = (url: string) => {
  activeImage.value = url || productImages.value[0]?.url || PLACEHOLDER_IMAGE;
  pendingLightboxIndex.value = null;
  isLightboxImageLoading.value = false;
  isLightboxOpen.value = true;
  scheduleLightboxPreload(currentImageIndex.value);
};

const showPreviousImage = () => {
  if (!hasMultipleImages.value) return;
  void navigateLightboxToIndex(visibleLightboxIndex.value - 1);
};

const showNextImage = () => {
  if (!hasMultipleImages.value) return;
  void navigateLightboxToIndex(visibleLightboxIndex.value + 1);
};

const handleLightboxTouchStart = (event: TouchEvent) => {
  if (!hasMultipleImages.value) return;

  const touch = event.changedTouches[0];
  lightboxTouchStartX.value = touch.clientX;
  lightboxTouchStartY.value = touch.clientY;
};

const handleLightboxTouchEnd = (event: TouchEvent) => {
  if (!hasMultipleImages.value) return;

  const touch = event.changedTouches[0];
  const deltaX = touch.clientX - lightboxTouchStartX.value;
  const deltaY = touch.clientY - lightboxTouchStartY.value;

  if (Math.abs(deltaX) < 44 || Math.abs(deltaX) < Math.abs(deltaY) * 1.15) {
    return;
  }

  if (deltaX < 0) {
    showNextImage();
  } else {
    showPreviousImage();
  }
};

const handleLightboxKeydown = (event: KeyboardEvent) => {
  if (!isLightboxOpen.value) return;

  if (event.key === 'Escape') {
    isLightboxOpen.value = false;
    return;
  }

  if (event.key === 'ArrowLeft') {
    event.preventDefault();
    showPreviousImage();
    return;
  }

  if (event.key === 'ArrowRight') {
    event.preventDefault();
    showNextImage();
  }
};

const updateSchema = (item: Product) => {
  const productPath = `/product/${item.slug || item.id}`;
  const categoryName = productStore.categories.find((category) => category.id === item.project_category_id)?.name || 'Мебель по размеру';
  const image = item.images[0]?.url || PLACEHOLDER_IMAGE;
  const additionalProperty = Object.entries(item.specs || {})
    .filter(([, value]) => value)
    .map(([name, value]) => ({
      '@type': 'PropertyValue',
      name,
      value,
    }));

  setPageSeo({
    title: `${item.name} — РОСТ Мебель`,
    description: compactDescription(item.description || `Проект ${item.name}: фотографии, бюджет, материалы и заявка на расчет.`),
    path: productPath,
    image,
    imageAlt: item.name,
    type: 'product',
  });

  const schema = {
    '@context': 'https://schema.org/',
    '@type': 'Product',
    '@id': absoluteUrl(`${productPath}#product`),
    name: item.name,
    sku: String(item.id),
    url: absoluteUrl(productPath),
    mainEntityOfPage: absoluteUrl(productPath),
    brand: {
      '@type': 'Brand',
      name: 'РОСТ Мебель',
    },
    image: (item.images.length ? item.images : [{ url: PLACEHOLDER_IMAGE }]).map((imageItem) => absoluteUrl(imageItem.url)),
    description: item.description,
    category: categoryName,
    additionalProperty: additionalProperty.length ? additionalProperty : undefined,
    offers: {
      '@type': 'Offer',
      url: absoluteUrl(productPath),
      priceCurrency: 'RUB',
      price: item.price,
      itemCondition: 'https://schema.org/NewCondition',
      availability: 'https://schema.org/InStock',
      seller: { '@type': 'Organization', name: 'РОСТ Мебель' },
    },
  };

  setJsonLd('schema-product', schema);
  setJsonLd('schema-product-breadcrumbs', {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement: [
      {
        '@type': 'ListItem',
        position: 1,
        name: 'Главная',
        item: absoluteUrl('/'),
      },
      {
        '@type': 'ListItem',
        position: 2,
        name: 'Проекты',
        item: absoluteUrl('/catalog'),
      },
      {
        '@type': 'ListItem',
        position: 3,
        name: item.name,
        item: absoluteUrl(productPath),
      },
    ],
  });
};

const loadProjectData = async () => {
  const id = route.params.id as string;
  if (!productStore.categories.length) {
    await productStore.fetchCategories();
  }
  const loadedProduct = await productStore.fetchProduct(id);

  if (loadedProduct) {
    product.value = loadedProduct;
    activeImage.value = loadedProduct.images[0]?.url || PLACEHOLDER_IMAGE;
    displayedMainImage.value = activeImage.value;
    isMainImageLoading.value = false;
    pendingLightboxIndex.value = null;
    updateSchema(loadedProduct);

    await productStore.fetchProducts({
      project_category_id: loadedProduct.project_category_id,
      limit: 4,
      status: 'published',
    });
    relatedProjects.value = productStore.products.filter((item) => item.id !== loadedProduct.id).slice(0, 3);
  }

  window.scrollTo({ top: 0, behavior: 'smooth' });
};

watch(() => route.params.id, () => {
  void loadProjectData();
});

watch(activeImage, (nextImage, previousImage) => {
  void syncDisplayedMainImage(nextImage || PLACEHOLDER_IMAGE, !previousImage);
});

watch(isLightboxOpen, (isOpen) => {
  if (isOpen && currentImageIndex.value >= 0) {
    scheduleLightboxPreload(currentImageIndex.value);
    return;
  }

  clearLightboxBackgroundPreloadTimer();
  pendingLightboxIndex.value = null;
  isLightboxImageLoading.value = false;
  lightboxNavigationRequestId += 1;
});

onMounted(() => {
  void loadProjectData();
  window.addEventListener('keydown', handleLightboxKeydown);
});

onUnmounted(() => {
  clearLightboxBackgroundPreloadTimer();
  window.removeEventListener('keydown', handleLightboxKeydown);
  removeJsonLd('schema-product');
  removeJsonLd('schema-product-breadcrumbs');
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
          <img
            :src="mainImageSource.src"
            :srcset="mainImageSource.srcset"
            :sizes="mainImageSource.sizes"
            :alt="product.name"
            :class="[
              'absolute inset-0 h-full w-full object-cover transition-all duration-300 ease-out group-hover:scale-[1.035]',
              isMainImageLoading ? 'opacity-90' : 'opacity-100'
            ]"
            decoding="async"
            fetchpriority="high"
            @error="handleImageError"
          >
          <div
            v-if="isMainImageLoading"
            class="absolute inset-0 flex items-center justify-center bg-brand-brown/10 backdrop-blur-[1px]"
          >
            <div class="h-9 w-9 animate-spin rounded-full border-2 border-white/85 border-t-transparent"></div>
          </div>
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
            <img :src="thumbnailImageUrl(image.url)" class="h-full w-full object-cover" alt="" loading="lazy" decoding="async" @error="handleImageError">
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
          <section class="ui-modal-panel max-w-2xl p-5 sm:p-8">
            <button type="button" class="absolute right-4 top-4 rounded-lg p-2 text-brand-brown/35 transition-colors hover:bg-brand-gray hover:text-brand-brown" @click="isOrderModalOpen = false">
              <LucideX :size="24" />
            </button>
            <h2 class="ui-title-md mb-2">Заявка на расчет</h2>
            <p class="ui-copy mb-6">Обсудим похожий проект и подскажем реалистичный бюджет.</p>
            <QuoteQuiz :project-id="product.id" :initial-project-type="quoteProjectType" @success="isOrderModalOpen = false" />
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
          <div
            class="relative z-10 flex h-full w-full items-center justify-center px-4 py-8 sm:px-8"
            @click.stop
            @touchstart.passive="handleLightboxTouchStart"
            @touchend.passive="handleLightboxTouchEnd"
          >
            <button
              type="button"
              class="absolute right-5 top-5 rounded-lg bg-white/10 p-3 text-white transition-colors hover:bg-white hover:text-brand-brown"
              @click="isLightboxOpen = false"
            >
              <LucideX :size="28" />
            </button>

            <button
              v-if="hasMultipleImages"
              type="button"
              class="absolute left-3 top-1/2 z-20 flex h-11 w-11 -translate-y-1/2 items-center justify-center rounded-full bg-black/35 text-white backdrop-blur transition-colors hover:bg-black/55 sm:left-5 sm:h-12 sm:w-12"
              aria-label="Предыдущее изображение"
              @click="showPreviousImage"
            >
              <LucideChevronLeft :size="28" />
            </button>

            <div class="relative flex max-h-full max-w-full items-center justify-center">
              <img
                :src="lightboxImageSource.src"
                :srcset="lightboxImageSource.srcset"
                :sizes="lightboxImageSource.sizes"
                :alt="product?.name || ''"
                :class="[
                  'max-h-full max-w-full rounded-lg object-contain shadow-2xl transition-opacity duration-300 ease-out',
                  isLightboxImageLoading ? 'opacity-75' : 'opacity-100'
                ]"
                decoding="async"
                fetchpriority="high"
                @error="handleImageError"
              >

              <div
                v-if="isLightboxImageLoading"
                class="absolute inset-0 flex items-center justify-center"
              >
                <div class="h-11 w-11 animate-spin rounded-full border-2 border-white/90 border-t-transparent"></div>
              </div>

              <div
                v-if="productImages.length"
                class="absolute bottom-4 left-1/2 -translate-x-1/2 rounded-full bg-black/55 px-3 py-1 text-sm font-semibold text-white backdrop-blur"
              >
                {{ visibleLightboxIndex + 1 }} / {{ productImages.length }}
              </div>
            </div>

            <button
              v-if="hasMultipleImages"
              type="button"
              class="absolute right-3 top-1/2 z-20 flex h-11 w-11 -translate-y-1/2 items-center justify-center rounded-full bg-black/35 text-white backdrop-blur transition-colors hover:bg-black/55 sm:right-5 sm:h-12 sm:w-12"
              aria-label="Следующее изображение"
              @click="showNextImage"
            >
              <LucideChevronRight :size="28" />
            </button>
          </div>
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

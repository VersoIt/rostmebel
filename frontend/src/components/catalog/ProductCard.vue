<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import {
  LucideArrowRight,
  LucideChevronLeft,
  LucideChevronRight,
  LucideEye,
  LucideHeart,
  LucideImages,
  LucideX,
} from 'lucide-vue-next';
import type { Product } from '@/types';
import { useFavorites } from '@/composables/useFavorites';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';
import { buildResponsiveImageSet } from '@/utils/images';
import { useProductStore } from '@/stores/products';

const props = defineProps<{
  product: Product;
}>();

const router = useRouter();
const productStore = useProductStore();
const { toggleFavorite, isFavorite } = useFavorites();

const isQuickViewOpen = ref(false);
const activeSlideIdx = ref(0);
const slideDirection = ref<'next' | 'prev'>('next');
const touchStartX = ref(0);
const touchStartY = ref(0);

const productTags = computed(() => {
  return (props.product.ai_tags || '')
    .split(',')
    .map((tag) => tag.trim())
    .filter(Boolean);
});

const categoryName = computed(() => {
  return productStore.categories.find((category) => category.id === props.product.project_category_id)?.name || 'Проект';
});

const sellingPoints = computed(() => {
  const specEntries = Object.entries(props.product.specs || {});
  const points: string[] = [];
  const material = specEntries.find(([key]) => /материал|фасад/i.test(key))?.[1];
  const style = specEntries.find(([key]) => /стиль/i.test(key))?.[1];
  const place = specEntries.find(([key]) => /город|локац|объект|адрес/i.test(key))?.[1];

  if (material) points.push(material);
  if (style) points.push(style);
  if (place) points.push(place);
  if (productTags.value[0]) points.push(productTags.value[0]);

  return Array.from(new Set(points)).slice(0, 3);
});

const priceLabel = computed(() => {
  return props.product.price > 0 ? 'Ориентир бюджета' : 'После расчета';
});

const imageCount = () => props.product.images.length;

const primaryImageSource = computed(() =>
  buildResponsiveImageSet(
    props.product.images[0]?.url || PLACEHOLDER_IMAGE,
    [320, 480, 640],
    '(min-width: 1280px) 400px, (min-width: 768px) 50vw, 100vw',
    74,
  ),
);

const quickViewImageSource = computed(() =>
  buildResponsiveImageSet(
    props.product.images[activeSlideIdx.value]?.url || PLACEHOLDER_IMAGE,
    [640, 960, 1280],
    '(min-width: 1024px) 60vw, 100vw',
    82,
  ),
);

const nextSlide = () => {
  if (imageCount() > 1) {
    slideDirection.value = 'next';
    activeSlideIdx.value = (activeSlideIdx.value + 1) % imageCount();
  }
};

const prevSlide = () => {
  if (imageCount() > 1) {
    slideDirection.value = 'prev';
    activeSlideIdx.value = (activeSlideIdx.value - 1 + imageCount()) % imageCount();
  }
};

const handleTouchStart = (event: TouchEvent) => {
  if (imageCount() < 2) return;
  const touch = event.changedTouches[0];
  touchStartX.value = touch.clientX;
  touchStartY.value = touch.clientY;
};

const handleTouchEnd = (event: TouchEvent) => {
  if (imageCount() < 2) return;
  const touch = event.changedTouches[0];
  const deltaX = touch.clientX - touchStartX.value;
  const deltaY = touch.clientY - touchStartY.value;

  if (Math.abs(deltaX) < 44 || Math.abs(deltaX) < Math.abs(deltaY) * 1.15) {
    return;
  }

  if (deltaX < 0) {
    nextSlide();
  } else {
    prevSlide();
  }
};

const handleImageError = (event: Event) => {
  (event.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
};

const goToProduct = () => {
  router.push(`/product/${props.product.slug || props.product.id}`);
};

const openQuickView = () => {
  isQuickViewOpen.value = true;
  activeSlideIdx.value = 0;
};

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(price);
};
</script>

<template>
  <article
    class="ui-card ui-card-hover group flex h-full cursor-pointer flex-col overflow-hidden"
    @click="goToProduct"
  >
    <div class="relative aspect-square shrink-0 overflow-hidden bg-brand-gray">
      <img
        :src="primaryImageSource.src"
        :srcset="primaryImageSource.srcset"
        :sizes="primaryImageSource.sizes"
        :alt="product.name"
        class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-[1.04]"
        loading="lazy"
        decoding="async"
        @error="handleImageError"
      >

      <div class="absolute inset-x-0 top-0 flex items-start justify-between p-3">
        <span class="rounded-full border border-white/15 bg-black/35 px-3 py-1 text-[10px] font-black uppercase tracking-[0.18em] text-white backdrop-blur">
          {{ categoryName }}
        </span>

        <div class="flex gap-2 sm:flex-col">
          <button
            type="button"
            :class="[
              'flex h-10 w-10 items-center justify-center rounded-xl shadow-md backdrop-blur transition-colors',
              isFavorite(product.id) ? 'bg-brand-gold text-white' : 'bg-white/92 text-brand-brown hover:text-brand-gold'
            ]"
            aria-label="Добавить в избранное"
            @click.stop="toggleFavorite(product)"
          >
            <LucideHeart :size="20" :fill="isFavorite(product.id) ? 'currentColor' : 'none'" />
          </button>
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-xl bg-white/92 text-brand-brown shadow-md backdrop-blur transition-colors hover:text-brand-gold"
            aria-label="Быстрый просмотр"
            @click.stop="openQuickView"
          >
            <LucideEye :size="20" />
          </button>
        </div>
      </div>

      <div class="absolute inset-x-0 bottom-0 p-3">
        <div class="rounded-2xl border border-white/12 bg-[linear-gradient(180deg,rgba(8,17,15,0.08),rgba(8,17,15,0.62))] p-3 text-white backdrop-blur-md">
          <div class="flex items-center justify-between gap-3">
            <div class="min-w-0">
              <div class="text-[10px] font-black uppercase tracking-[0.18em] text-white/66">{{ priceLabel }}</div>
              <div class="mt-1 truncate text-lg font-semibold text-white">{{ formatPrice(product.price) }}</div>
            </div>
            <div class="inline-flex items-center gap-1 rounded-full bg-white/10 px-2.5 py-1 text-xs font-bold text-white/78">
              <LucideImages :size="14" />
              {{ imageCount() }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex flex-1 flex-col p-5">
      <div class="mb-2 text-xs font-black uppercase tracking-[0.18em] text-brand-brown/35">
        {{ productTags[0] || 'Проект по размеру' }}
      </div>
      <h3 class="line-clamp-2 font-serif text-2xl font-bold text-brand-brown transition-colors group-hover:text-brand-gold">
        {{ product.name }}
      </h3>

      <p class="mt-3 line-clamp-3 leading-7 text-brand-brown/62">
        {{ product.description }}
      </p>

      <div v-if="sellingPoints.length" class="mt-4 flex flex-wrap gap-2">
        <span
          v-for="point in sellingPoints"
          :key="point"
          class="rounded-full bg-brand-gray px-3 py-1.5 text-[11px] font-bold text-brand-brown/62"
        >
          {{ point }}
        </span>
      </div>

      <div class="mt-auto pt-6">
        <div class="flex items-center justify-between gap-3 border-t border-brand-brown/10 pt-4">
          <div class="min-w-0">
            <div class="text-[11px] font-black uppercase tracking-[0.18em] text-brand-brown/30">Смотреть проект</div>
            <div class="mt-1 text-sm text-brand-brown/58">Фото, бюджет, детали и похожие решения</div>
          </div>
          <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl bg-brand-brown text-white transition-colors group-hover:bg-brand-gold">
            <LucideArrowRight :size="20" />
          </span>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <transition name="modal-fade">
        <div v-if="isQuickViewOpen" class="ui-modal-backdrop" @click.stop>
          <div class="absolute inset-0" @click="isQuickViewOpen = false"></div>

          <section class="ui-modal-panel z-10 max-w-6xl overflow-hidden bg-white">
            <div class="grid grid-cols-1 lg:min-h-[560px] lg:grid-cols-[minmax(0,1fr)_360px]">
              <div
                class="relative aspect-[4/3] min-h-[320px] touch-pan-y overflow-hidden bg-brand-gray sm:aspect-video lg:aspect-auto lg:min-h-[560px]"
                @touchstart.passive="handleTouchStart"
                @touchend.passive="handleTouchEnd"
              >
                <transition :name="slideDirection === 'next' ? 'slide-next' : 'slide-prev'">
                  <img
                    :key="activeSlideIdx"
                    :src="quickViewImageSource.src"
                    :srcset="quickViewImageSource.srcset"
                    :sizes="quickViewImageSource.sizes"
                    class="absolute inset-0 h-full w-full object-cover"
                    alt=""
                    decoding="async"
                    @error="handleImageError"
                  >
                </transition>

                <div v-if="product.images.length > 1" class="absolute inset-x-0 top-1/2 flex -translate-y-1/2 items-center justify-between px-3 sm:px-5">
                  <button type="button" class="flex h-11 w-11 items-center justify-center rounded-xl bg-black/35 text-white backdrop-blur transition-colors hover:bg-black/55" @click.stop="prevSlide">
                    <LucideChevronLeft :size="28" />
                  </button>
                  <button type="button" class="flex h-11 w-11 items-center justify-center rounded-xl bg-black/35 text-white backdrop-blur transition-colors hover:bg-black/55" @click.stop="nextSlide">
                    <LucideChevronRight :size="28" />
                  </button>
                </div>

                <div v-if="product.images.length > 1" class="absolute bottom-4 left-1/2 flex -translate-x-1/2 items-center gap-2 rounded-full bg-black/45 px-3 py-1.5 text-xs font-bold text-white backdrop-blur">
                  <LucideImages :size="14" />
                  {{ activeSlideIdx + 1 }} / {{ product.images.length }}
                </div>
              </div>

              <aside class="flex flex-col bg-white p-6">
                <div class="mb-auto">
                  <div class="mb-3 inline-flex rounded-full bg-brand-gold/10 px-3 py-1 text-[10px] font-black uppercase tracking-[0.18em] text-brand-gold">
                    {{ categoryName }}
                  </div>
                  <h3 class="ui-title-md mb-4">{{ product.name }}</h3>
                  <div class="mb-1 text-[11px] font-black uppercase tracking-[0.18em] text-brand-brown/35">{{ priceLabel }}</div>
                  <div class="mb-4 text-2xl font-semibold text-brand-brown">{{ formatPrice(product.price) }}</div>
                  <div v-if="sellingPoints.length" class="mb-5 flex flex-wrap gap-2">
                    <span
                      v-for="point in sellingPoints"
                      :key="point"
                      class="rounded-full bg-brand-gray px-3 py-1.5 text-[11px] font-bold text-brand-brown/62"
                    >
                      {{ point }}
                    </span>
                  </div>
                  <p class="line-clamp-[8] leading-7 text-brand-brown/62">
                    {{ product.description }}
                  </p>
                </div>
                <button type="button" class="ui-button ui-button-primary mt-6 w-full" @click="goToProduct">
                  Подробнее
                </button>
              </aside>
            </div>

            <button
              type="button"
              class="absolute right-4 top-4 flex h-11 w-11 items-center justify-center rounded-xl bg-black/35 text-white backdrop-blur transition-colors hover:bg-red-600"
              aria-label="Закрыть быстрый просмотр"
              @click.stop="isQuickViewOpen = false"
            >
              <LucideX :size="23" />
            </button>
          </section>
        </div>
      </transition>
    </Teleport>
  </article>
</template>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.slide-next-enter-active,
.slide-next-leave-active,
.slide-prev-enter-active,
.slide-prev-leave-active {
  transition: opacity 0.26s ease, transform 0.26s ease;
}

.slide-next-enter-from,
.slide-prev-leave-to {
  opacity: 0;
  transform: translateX(16px);
}

.slide-prev-enter-from,
.slide-next-leave-to {
  opacity: 0;
  transform: translateX(-16px);
}
</style>

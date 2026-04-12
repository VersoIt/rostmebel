<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { LucideChevronLeft, LucideChevronRight, LucideEye, LucideHeart, LucideX } from 'lucide-vue-next';
import type { Product } from '@/types';
import { useFavorites } from '@/composables/useFavorites';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';

const props = defineProps<{
  product: Product;
}>();

const router = useRouter();
const { toggleFavorite, isFavorite } = useFavorites();

const isQuickViewOpen = ref(false);
const activeSlideIdx = ref(0);
const slideDirection = ref<'next' | 'prev'>('next');

const imageCount = () => props.product.images.length;

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

const handleImageError = (event: Event) => {
  (event.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
};

const goToProduct = () => {
  router.push(`/product/${props.product.id}`);
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
        :src="product.images[0]?.url || PLACEHOLDER_IMAGE"
        :alt="product.name"
        class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-[1.035]"
        @error="handleImageError"
      >

      <div class="absolute right-3 top-3 flex gap-2 sm:flex-col">
        <button
          type="button"
          :class="[
            'flex h-10 w-10 items-center justify-center rounded-lg shadow-md transition-colors',
            isFavorite(product.id) ? 'bg-brand-gold text-white' : 'bg-white text-brand-brown hover:text-brand-gold'
          ]"
          aria-label="Добавить в избранное"
          @click.stop="toggleFavorite(product)"
        >
          <LucideHeart :size="20" :fill="isFavorite(product.id) ? 'currentColor' : 'none'" />
        </button>
        <button
          type="button"
          class="flex h-10 w-10 items-center justify-center rounded-lg bg-white text-brand-brown shadow-md transition-colors hover:text-brand-gold"
          aria-label="Быстрый просмотр"
          @click.stop="openQuickView"
        >
          <LucideEye :size="20" />
        </button>
      </div>
    </div>

    <div class="flex flex-1 flex-col p-5">
      <div class="mb-2 text-xs font-bold uppercase tracking-widest text-brand-brown/40">
        {{ product.ai_tags?.split(',')[0]?.trim() || 'Проект' }}
      </div>
      <h3 class="line-clamp-2 font-serif text-xl font-bold text-brand-brown transition-colors group-hover:text-brand-gold">
        {{ product.name }}
      </h3>
      <div class="mt-auto flex flex-wrap items-baseline gap-3 pt-5">
        <span class="text-xl font-semibold text-brand-brown">{{ formatPrice(product.price) }}</span>
        <span v-if="product.price_old" class="text-sm text-brand-brown/30 line-through">
          {{ formatPrice(product.price_old) }}
        </span>
      </div>
    </div>

    <Teleport to="body">
      <transition name="modal-fade">
        <div v-if="isQuickViewOpen" class="ui-modal-backdrop" @click.stop>
          <div class="absolute inset-0" @click="isQuickViewOpen = false"></div>

          <section class="ui-modal-panel z-10 max-w-6xl overflow-hidden bg-black lg:aspect-video">
            <div class="grid min-h-[70vh] grid-cols-1 lg:grid-cols-[1fr_320px]">
              <div class="relative min-h-[58vh] overflow-hidden bg-neutral-900">
                <transition :name="slideDirection === 'next' ? 'slide-next' : 'slide-prev'">
                  <img
                    :key="activeSlideIdx"
                    :src="product.images[activeSlideIdx]?.url || PLACEHOLDER_IMAGE"
                    class="absolute inset-0 h-full w-full object-cover"
                    alt=""
                    @error="handleImageError"
                  >
                </transition>

                <div v-if="product.images.length > 1" class="absolute inset-x-0 top-1/2 flex -translate-y-1/2 items-center justify-between px-3 sm:px-5">
                  <button type="button" class="flex h-11 w-11 items-center justify-center rounded-lg bg-black/35 text-white backdrop-blur transition-colors hover:bg-black/55" @click.stop="prevSlide">
                    <LucideChevronLeft :size="28" />
                  </button>
                  <button type="button" class="flex h-11 w-11 items-center justify-center rounded-lg bg-black/35 text-white backdrop-blur transition-colors hover:bg-black/55" @click.stop="nextSlide">
                    <LucideChevronRight :size="28" />
                  </button>
                </div>
              </div>

              <aside class="flex flex-col bg-white p-6">
                <div class="mb-auto">
                  <div class="ui-eyebrow mb-3">Быстрый просмотр</div>
                  <h3 class="ui-title-md mb-4">{{ product.name }}</h3>
                  <div class="mb-5 text-2xl font-semibold text-brand-brown">{{ formatPrice(product.price) }}</div>
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
              class="absolute right-4 top-4 flex h-11 w-11 items-center justify-center rounded-lg bg-black/35 text-white backdrop-blur transition-colors hover:bg-red-600"
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

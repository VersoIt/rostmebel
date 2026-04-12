<script setup lang="ts">
import { LucideHeart, LucideEye, LucideChevronLeft, LucideChevronRight, LucideX } from 'lucide-vue-next';
import type { Product } from '@/types';
import { useFavorites } from '@/composables/useFavorites';
import { useRouter } from 'vue-router';
import { ref } from 'vue';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';

const props = defineProps<{
  product: Product;
}>();

const router = useRouter();
const { toggleFavorite, isFavorite } = useFavorites();

const isQuickViewOpen = ref(false);
const activeSlideIdx = ref(0);
const slideDirection = ref('next');

const nextSlide = () => {
  if (props.product.images.length > 1) {
    slideDirection.value = 'next';
    activeSlideIdx.value = (activeSlideIdx.value + 1) % props.product.images.length;
  }
};

const prevSlide = () => {
  if (props.product.images.length > 1) {
    slideDirection.value = 'prev';
    activeSlideIdx.value = (activeSlideIdx.value - 1 + props.product.images.length) % props.product.images.length;
  }
};

const handleImageError = (e: Event) => {
  (e.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
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
  <div 
    @click="goToProduct"
    class="group bg-white rounded-lg overflow-hidden hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 border border-brand-brown/5 cursor-pointer h-full flex flex-col"
  >
    <!-- Image Area -->
    <div class="relative aspect-square overflow-hidden bg-brand-gray shrink-0">
      <img 
        :src="product.images[0]?.url || PLACEHOLDER_IMAGE" 
        @error="handleImageError"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        :alt="product.name"
      >
      
      <div class="absolute top-4 right-4 flex flex-col gap-2 translate-x-12 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 transition-all duration-300">
        <button 
          @click.stop="toggleFavorite(product)"
          :class="['w-10 h-10 rounded-lg flex items-center justify-center shadow-lg transition-colors', isFavorite(product.id) ? 'bg-brand-gold text-white' : 'bg-white text-brand-brown hover:text-brand-gold']"
        >
          <LucideHeart :size="20" :fill="isFavorite(product.id) ? 'currentColor' : 'none'" />
        </button>
        <button 
          @click.stop="openQuickView" 
          class="w-10 h-10 bg-white rounded-lg flex items-center justify-center text-brand-brown hover:text-brand-gold shadow-lg"
        >
          <LucideEye :size="20" />
        </button>
      </div>
    </div>

    <!-- Content Area -->
    <div class="p-6 flex-1 flex flex-col">
      <div class="text-xs text-brand-brown/40 uppercase tracking-widest mb-2">
        {{ product.ai_tags?.split(',')[0] || 'Проект' }}
      </div>
      <h3 class="font-serif text-lg text-brand-brown mb-2 group-hover:text-brand-gold transition-colors line-clamp-2">
        {{ product.name }}
      </h3>
      <div class="mt-auto flex items-center gap-3">
        <span class="text-xl font-medium text-brand-brown">{{ formatPrice(product.price) }}</span>
        <span v-if="product.price_old" class="text-brand-brown/30 line-through text-sm">
          {{ formatPrice(product.price_old) }}
        </span>
      </div>
    </div>

    <!-- Quick View Modal -->
    <Teleport to="body">
      <transition name="modal-fade">
        <div v-if="isQuickViewOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-12">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-md" @click.stop="isQuickViewOpen = false"></div>
          
          <div class="relative w-full max-w-6xl aspect-video bg-black rounded-lg overflow-hidden shadow-2xl flex group/modal">
            
            <!-- Slider Area -->
            <div class="flex-1 relative bg-neutral-900 overflow-hidden">
              <transition :name="slideDirection === 'next' ? 'slide-next' : 'slide-prev'">
                <img 
                  :key="activeSlideIdx"
                  :src="product.images[activeSlideIdx]?.url || PLACEHOLDER_IMAGE" 
                  class="absolute inset-0 w-full h-full object-cover"
                  @error="handleImageError"
                >
              </transition>
              
              <!-- Controls -->
              <div v-if="product.images.length > 1" class="absolute inset-0 flex items-center justify-between px-6 opacity-0 group-hover/modal:opacity-100 transition-opacity pointer-events-none">
                <button @click.stop="prevSlide" class="w-14 h-14 bg-white/10 hover:bg-white/20 text-white rounded-lg flex items-center justify-center backdrop-blur-md pointer-events-auto transition-all">
                  <LucideChevronLeft :size="32" />
                </button>
                <button @click.stop="nextSlide" class="w-14 h-14 bg-white/10 hover:bg-white/20 text-white rounded-lg flex items-center justify-center backdrop-blur-md pointer-events-auto transition-all">
                  <LucideChevronRight :size="32" />
                </button>
              </div>

              <!-- Indicators -->
              <div v-if="product.images.length > 1" class="absolute bottom-6 left-1/2 -translate-x-1/2 flex gap-2">
                <div v-for="(_, idx) in product.images" :key="idx" 
                  :class="['h-1.5 rounded-full transition-all', idx === activeSlideIdx ? 'w-8 bg-brand-gold' : 'w-2 bg-white/30']">
                </div>
              </div>
            </div>

            <!-- Side Info -->
            <div class="hidden lg:flex w-80 bg-white p-10 flex-col shrink-0">
              <div class="mb-auto">
                <div class="text-brand-gold font-bold text-[10px] uppercase tracking-[0.2em] mb-4">Быстрый просмотр</div>
                <h3 class="font-serif text-3xl text-brand-brown mb-4 leading-tight">{{ product.name }}</h3>
                <div class="text-2xl font-medium text-brand-brown mb-6">{{ formatPrice(product.price) }}</div>
                <p class="text-brand-brown/60 text-sm leading-relaxed line-clamp-[10]">
                  {{ product.description }}
                </p>
              </div>
              <button @click="goToProduct" class="w-full bg-brand-brown text-white py-4 rounded-lg font-bold hover:bg-brand-gold transition-all shadow-lg active:scale-95"> ПОДРОБНЕЕ </button>
            </div>

            <!-- Close Button -->
            <button @click.stop="isQuickViewOpen = false" class="absolute top-6 right-6 w-12 h-12 bg-white/10 hover:bg-red-500 text-white rounded-lg flex items-center justify-center backdrop-blur-md transition-all z-10">
              <LucideX :size="24" />
            </button>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped>
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.4s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }

.slide-next-enter-active, .slide-next-leave-active,
.slide-prev-enter-active, .slide-prev-leave-active {
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-next-enter-from { transform: translateX(100%); }
.slide-next-leave-to { transform: translateX(-100%); }

.slide-prev-enter-from { transform: translateX(-100%); }
.slide-prev-leave-to { transform: translateX(100%); }
</style>

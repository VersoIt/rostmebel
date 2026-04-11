<script setup lang="ts">
import { LucideHeart, LucideEye, LucideChevronLeft, LucideChevronRight, LucideX } from 'lucide-vue-next';
import type { Product } from '@/types';
import { useFavorites } from '@/composables/useFavorites';
import { useRouter } from 'vue-router';
import { ref } from 'vue';

const props = defineProps<{
  product: Product;
}>();

const router = useRouter();
const { toggleFavorite, isFavorite } = useFavorites();

// Fast Base64 Placeholder (Zero network request)
const placeholder = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 800 800'%3E%3Crect width='800' height='800' fill='%23F5F5F0'/%3E%3Cpath d='M350 400l50-50 50 50M300 500h200' stroke='%23C9A84C' stroke-width='2' fill='none'/%3E%3Ctext x='50%25' y='60%25' text-anchor='middle' font-family='serif' font-size='24' fill='%232C1810' opacity='0.2'%3EPOCT MEБEЛЬ%3C/text%3E%3C/svg%3E";

const isQuickViewOpen = ref(false);
const activeSlideIdx = ref(0);

const nextSlide = () => {
  if (props.product.images.length > 0) {
    activeSlideIdx.value = (activeSlideIdx.value + 1) % props.product.images.length;
  }
};

const prevSlide = () => {
  if (props.product.images.length > 0) {
    activeSlideIdx.value = (activeSlideIdx.value - 1 + props.product.images.length) % props.product.images.length;
  }
};

const handleImageError = (e: Event) => {
  (e.target as HTMLImageElement).src = placeholder;
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
    class="group bg-white rounded-2xl overflow-hidden hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 border border-brand-brown/5 cursor-pointer h-full flex flex-col"
  >
    <div class="relative aspect-square overflow-hidden bg-brand-gray shrink-0">
      <img 
        :src="product.images[0]?.url || placeholder" 
        @error="handleImageError"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        :alt="product.name"
      >
      
      <div class="absolute top-4 right-4 flex flex-col gap-2 translate-x-12 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 transition-all duration-300">
        <button 
          @click.stop="toggleFavorite(product)"
          :class="['w-10 h-10 rounded-full flex items-center justify-center shadow-lg transition-colors', isFavorite(product.id) ? 'bg-brand-gold text-white' : 'bg-white text-brand-brown hover:text-brand-gold']"
        >
          <LucideHeart :size="20" :fill="isFavorite(product.id) ? 'currentColor' : 'none'" />
        </button>
        <button 
          @click.stop="openQuickView" 
          class="w-10 h-10 bg-white rounded-full flex items-center justify-center text-brand-brown hover:text-brand-gold shadow-lg"
        >
          <LucideEye :size="20" />
        </button>
      </div>

      <div v-if="product.price_old" class="absolute top-4 left-4 bg-brand-gold text-white px-3 py-1 rounded-lg text-sm font-medium">
        SALE
      </div>
    </div>

    <div class="p-6 flex-1 flex flex-col">
      <div class="text-xs text-brand-brown/40 uppercase tracking-widest mb-2">
        {{ product.ai_tags?.split(',')[0] || 'Мебель' }}
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

    <!-- Quick View Modal (Teleported to root to avoid parent clipping) -->
    <Teleport to="body">
      <transition 
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div v-if="isQuickViewOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-12">
          <div class="absolute inset-0 bg-brand-brown/95 backdrop-blur-xl" @click.stop="isQuickViewOpen = false"></div>
          
          <div class="relative w-full max-w-6xl aspect-video bg-black rounded-[2rem] overflow-hidden shadow-[0_0_100px_rgba(0,0,0,0.5)] group/modal flex">
            <div class="flex-1 relative bg-neutral-900 flex items-center justify-center">
              <img 
                :src="product.images[activeSlideIdx]?.url || placeholder" 
                class="max-w-full max-h-full object-contain"
                @error="handleImageError"
              >
              
              <!-- Controls -->
              <button 
                v-if="product.images.length > 1"
                @click.stop="prevSlide" 
                class="absolute left-6 top-1/2 -translate-y-1/2 w-16 h-16 bg-white/5 hover:bg-white/10 text-white rounded-full flex items-center justify-center backdrop-blur-md transition-all"
              >
                <LucideChevronLeft :size="40" />
              </button>
              <button 
                v-if="product.images.length > 1"
                @click.stop="nextSlide" 
                class="absolute right-6 top-1/2 -translate-y-1/2 w-16 h-16 bg-white/5 hover:bg-white/10 text-white rounded-full flex items-center justify-center backdrop-blur-md transition-all"
              >
                <LucideChevronRight :size="40" />
              </button>
            </div>

            <!-- Side Info -->
            <div class="hidden lg:flex w-80 bg-white p-10 flex-col">
              <div class="mb-auto">
                <div class="text-brand-gold font-bold text-xs uppercase tracking-widest mb-4">Quick View</div>
                <h3 class="font-serif text-3xl text-brand-brown mb-4">{{ product.name }}</h3>
                <div class="text-2xl font-medium text-brand-brown mb-6">{{ formatPrice(product.price) }}</div>
                <p class="text-brand-brown/60 text-sm leading-relaxed line-clamp-6">
                  {{ product.description }}
                </p>
              </div>
              <button @click="goToProduct" class="w-full bg-brand-brown text-white py-4 rounded-xl font-bold hover:bg-brand-gold transition-all"> ПОДРОБНЕЕ </button>
            </div>

            <!-- Close -->
            <button 
              @click.stop="isQuickViewOpen = false" 
              class="absolute top-6 right-6 w-12 h-12 bg-white/10 hover:bg-red-500 text-white rounded-full flex items-center justify-center backdrop-blur-md transition-all"
            >
              <LucideX :size="24" />
            </button>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

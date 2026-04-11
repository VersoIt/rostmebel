<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useProductStore } from '@/stores/products';
import { LucideCheckCircle, LucideShieldCheck, LucideTruck } from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';
import type { Product } from '@/types';

const route = useRoute();
const productStore = useProductStore();
const product = ref<Product | null>(null);
const activeImage = ref('');
const isOrderModalOpen = ref(false);

const placeholder = 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?q=80&w=800&auto=format&fit=crop';

const handleImageError = (e: Event) => {
  (e.target as HTMLImageElement).src = placeholder;
};

onMounted(async () => {
  const p = await productStore.fetchProduct(route.params.id as string);
  if (p) {
    product.value = p;
    activeImage.value = p.images[0]?.url || placeholder;
  }
});

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(price);
};
</script>

<template>
  <div v-if="product" class="bg-white min-h-screen pt-32 pb-24">
    <div class="max-w-7xl mx-auto px-4">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-16">
        <!-- Gallery -->
        <div class="space-y-6">
          <div class="aspect-square rounded-3xl overflow-hidden bg-brand-gray border border-brand-brown/5">
            <img 
              :src="activeImage" 
              @error="handleImageError"
              class="w-full h-full object-cover" 
              :alt="product.name"
            >
          </div>
          <div class="flex gap-4 overflow-x-auto no-scrollbar">
            <button 
              v-for="img in product.images" 
              :key="img.url"
              @click="activeImage = img.url"
              :class="['w-24 h-24 rounded-2xl overflow-hidden border-2 transition-all shrink-0', activeImage === img.url ? 'border-brand-gold shadow-lg' : 'border-transparent']"
            >
              <img :src="img.url" @error="handleImageError" class="w-full h-full object-cover">
            </button>
          </div>
        </div>

        <!-- Info -->
        <div class="flex flex-col">
          <div class="mb-8">
            <h1 class="font-serif text-5xl text-brand-brown mb-4">{{ product.name }}</h1>
            <div class="flex items-center gap-4 mb-6">
              <span class="text-3xl font-medium text-brand-gold">{{ formatPrice(product.price) }}</span>
              <span v-if="product.price_old" class="text-xl text-brand-brown/30 line-through">
                {{ formatPrice(product.price_old) }}
              </span>
            </div>
            <p class="text-lg text-brand-brown/70 leading-relaxed">
              {{ product.description }}
            </p>
          </div>

          <!-- Specs -->
          <div class="bg-brand-gray/30 p-8 rounded-3xl border border-brand-brown/5 mb-10">
            <h3 class="font-serif text-2xl mb-6">Характеристики</h3>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
              <div v-for="(value, key) in product.specs" :key="key" class="flex flex-col">
                <span class="text-xs text-brand-brown/40 uppercase tracking-widest">{{ key }}</span>
                <span class="text-lg font-medium text-brand-brown">{{ value }}</span>
              </div>
            </div>
          </div>

          <!-- CTA & Benefits -->
          <div class="space-y-8 mt-auto">
            <button 
              @click="isOrderModalOpen = true"
              class="w-full bg-brand-brown text-white py-6 rounded-2xl text-xl font-medium hover:bg-brand-gold transition-all shadow-xl hover:shadow-2xl active:scale-[0.98]"
            >
              Оставить заявку
            </button>
            
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
              <div class="flex items-center gap-3 text-sm text-brand-brown/60">
                <LucideShieldCheck class="text-brand-gold" :size="20" />
                Гарантия 2 года
              </div>
              <div class="flex items-center gap-3 text-sm text-brand-brown/60">
                <LucideTruck class="text-brand-gold" :size="20" />
                Быстрая доставка
              </div>
              <div class="flex items-center gap-3 text-sm text-brand-brown/60">
                <LucideCheckCircle class="text-brand-gold" :size="20" />
                Сборка на месте
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Order Modal -->
    <div v-if="isOrderModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-brand-brown/60 backdrop-blur-sm" @click="isOrderModalOpen = false"></div>
      <div class="relative bg-white w-full max-w-lg rounded-3xl shadow-2xl p-8 max-h-[90vh] overflow-y-auto">
        <button @click="isOrderModalOpen = false" class="absolute top-6 right-6 text-brand-brown/40 hover:text-brand-brown text-2xl">&times;</button>
        <h2 class="font-serif text-3xl mb-2">Оформить заявку</h2>
        <p class="text-brand-brown/60 mb-8">Мы свяжемся с вами в течение 15 минут</p>
        <OrderForm :product-id="product.id" @success="isOrderModalOpen = false" />
      </div>
    </div>
  </div>
</template>

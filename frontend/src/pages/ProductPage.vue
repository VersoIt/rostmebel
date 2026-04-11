<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useProductStore } from '@/stores/products';
import { LucideCheckCircle, LucideShieldCheck, LucideTruck, LucideSearch, LucideX } from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';
import type { Product } from '@/types';

const route = useRoute();
const productStore = useProductStore();
const product = ref<Product | null>(null);
const activeImage = ref('');
const isOrderModalOpen = ref(false);
const isLightboxOpen = ref(false);

const placeholder = 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?q=80&w=800&auto=format&fit=crop';

const handleImageError = (e: Event) => {
  (e.target as HTMLImageElement).src = placeholder;
};

const openLightbox = (url: string) => {
  activeImage.value = url;
  isLightboxOpen.value = true;
};

const updateSchema = (p: Product) => {
  const schema = {
    "@context": "https://schema.org/",
    "@type": "Product",
    "name": p.name,
    "image": p.images.map(i => i.url),
    "description": p.description,
    "offers": {
      "@type": "Offer",
      "url": window.location.href,
      "priceCurrency": "RUB",
      "price": p.price,
      "availability": "https://schema.org/InStock",
      "seller": {
        "@type": "Organization",
        "name": "РОСТ Мебель"
      }
    }
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

onMounted(async () => {
  const p = await productStore.fetchProduct(route.params.id as string);
  if (p) {
    product.value = p;
    activeImage.value = p.images[0]?.url || placeholder;
    updateSchema(p);
    
    const description = `${p.name} — реализованный проект мебели на заказ. ${p.description.substring(0, 150)}...`;
    const metaDesc = document.querySelector('meta[name="description"]');
    if (metaDesc) {
      metaDesc.setAttribute('content', description);
    }
    document.title = `${p.name} — РОСТ Мебель`;
  }
});

onUnmounted(() => {
  const script = document.getElementById('schema-product');
  if (script) script.remove();
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
          <div 
            @click="openLightbox(activeImage)"
            class="aspect-square rounded-3xl overflow-hidden bg-brand-gray border border-brand-brown/5 cursor-zoom-in group relative"
          >
            <img 
              :src="activeImage" 
              @error="handleImageError"
              class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110" 
              :alt="product.name"
            >
            <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors flex items-center justify-center">
              <LucideSearch :size="48" class="text-white opacity-0 group-hover:opacity-100 transition-opacity" />
            </div>
          </div>
          <div class="flex gap-4 overflow-x-auto no-scrollbar pb-2">
            <button 
              v-for="img in product.images" 
              :key="img.url"
              @click="activeImage = img.url"
              :class="['w-24 h-24 rounded-2xl overflow-hidden border-2 transition-all shrink-0 hover:scale-105', activeImage === img.url ? 'border-brand-gold shadow-lg' : 'border-transparent opacity-60 hover:opacity-100']"
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
            <h3 class="font-serif text-2xl mb-6">Детали проекта</h3>
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
              Хочу такой же проект
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
    <Teleport to="body">
      <transition name="fade">
        <div v-if="isOrderModalOpen" class="fixed inset-0 z-[100] bg-black/95 backdrop-blur-xl flex items-center justify-center p-4" @click.self="isOrderModalOpen = false">
          <div class="relative bg-white w-full max-w-lg rounded-3xl shadow-2xl p-10 overflow-hidden transform transition-all">
            <button @click="isOrderModalOpen = false" class="absolute top-6 right-6 text-brand-brown/40 hover:text-brand-brown transition-colors">
              <LucideX :size="28" />
            </button>
            
            <div class="mb-8">
              <h2 class="font-serif text-4xl mb-3 text-brand-brown">Оформить заявку</h2>
              <p class="text-brand-brown/60">Мы свяжемся с вами в течение 15 минут для обсуждения деталей проекта</p>
            </div>

            <OrderForm :product-id="product.id" @success="isOrderModalOpen = false" />
            
            <!-- Decorative element -->
            <div class="absolute -bottom-12 -right-12 w-32 h-32 bg-brand-gold/5 rounded-full blur-3xl"></div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Lightbox -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLightboxOpen" class="fixed inset-0 z-[100] bg-black/95 backdrop-blur-xl flex items-center justify-center p-4 md:p-12" @click="isLightboxOpen = false">
          <button class="absolute top-8 right-8 text-white/50 hover:text-white transition-colors">
            <LucideX :size="40" />
          </button>
          <img :src="activeImage" class="max-w-full max-h-full object-contain shadow-2xl rounded-lg border border-white/10 transition-transform duration-500">
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { 
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); 
}
.fade-enter-from, .fade-leave-to { 
  opacity: 0; 
  backdrop-filter: blur(0px);
}

/* Modal Content Animation */
.fade-enter-active .relative {
  animation: modal-in 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes modal-in {
  from { opacity: 0; transform: scale(0.9) translateY(20px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}
</style>

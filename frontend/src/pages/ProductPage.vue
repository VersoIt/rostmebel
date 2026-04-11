<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProductStore } from '@/stores/products';
import { 
  LucideCheckCircle, 
  LucideShieldCheck, 
  LucideTruck, 
  LucideSearch, 
  LucideX,
  LucideChevronLeft,
  LucideArrowRight
} from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';
import ProductCard from '@/components/catalog/ProductCard.vue';
import type { Product } from '@/types';

const route = useRoute();
const router = useRouter();
const productStore = useProductStore();
const product = ref<Product | null>(null);
const relatedProjects = ref<Product[]>([]);
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

const loadProjectData = async () => {
  const id = route.params.id as string;
  const p = await productStore.fetchProduct(id);
  if (p) {
    product.value = p;
    activeImage.value = p.images[0]?.url || placeholder;
    updateSchema(p);
    
    // SEO
    document.title = `${p.name} — РОСТ Мебель`;
    
    // Fetch related
    await productStore.fetchProducts({ 
      category_id: p.category_id, 
      limit: 4,
      status: 'published'
    });
    relatedProjects.value = productStore.products.filter(item => item.id !== p.id).slice(0, 3);
  }
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

watch(() => route.params.id, loadProjectData);

onMounted(loadProjectData);

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
      "seller": { "@type": "Organization", "name": "РОСТ Мебель" }
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

onUnmounted(() => {
  const script = document.getElementById('schema-product');
  if (script) script.remove();
});

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency', currency: 'RUB', maximumFractionDigits: 0,
  }).format(price);
};
</script>

<template>
  <div v-if="product" class="bg-white min-h-screen">
    <!-- Project Header -->
    <div class="pt-32 pb-12 bg-brand-cream/30">
      <div class="max-w-7xl mx-auto px-6">
        <button @click="router.push('/catalog')" class="flex items-center gap-2 text-brand-brown/40 hover:text-brand-gold transition-colors font-bold text-xs uppercase tracking-widest mb-8 group">
          <LucideChevronLeft :size="16" class="group-hover:-translate-x-1 transition-transform" />
          Назад к проектам
        </button>
        <h1 class="font-serif text-5xl md:text-7xl text-brand-brown mb-4">{{ product.name }}</h1>
        <div class="flex items-center gap-4">
          <span class="bg-brand-gold text-brand-brown px-4 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest">
            {{ productStore.categories.find(c => c.id === product?.category_id)?.name || 'Проект' }}
          </span>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-6 py-20">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-20">
        <!-- Gallery -->
        <div class="space-y-8">
          <div 
            @click="openLightbox(activeImage)"
            class="aspect-square rounded-[2.5rem] overflow-hidden bg-brand-gray border border-brand-brown/5 cursor-zoom-in group relative shadow-2xl"
          >
            <img :src="activeImage" @error="handleImageError" class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110" :alt="product.name">
            <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-colors flex items-center justify-center">
              <LucideSearch :size="48" class="text-white opacity-0 group-hover:opacity-100 transition-opacity" />
            </div>
          </div>
          <div class="flex gap-4 overflow-x-auto no-scrollbar pb-4">
            <button 
              v-for="img in product.images" :key="img.url"
              @click="activeImage = img.url"
              :class="['w-24 h-24 rounded-2xl overflow-hidden border-2 transition-all shrink-0 hover:scale-105 shadow-md', activeImage === img.url ? 'border-brand-gold ring-4 ring-brand-gold/10' : 'border-transparent opacity-60 hover:opacity-100']"
            >
              <img :src="img.url" @error="handleImageError" class="w-full h-full object-cover">
            </button>
          </div>
        </div>

        <!-- Info -->
        <div class="flex flex-col">
          <div class="mb-12">
            <div class="flex items-center gap-6 mb-8">
              <div class="flex flex-col">
                <span class="text-xs text-brand-brown/40 uppercase tracking-[0.2em] font-bold mb-1">Бюджет реализации</span>
                <span class="text-4xl font-serif text-brand-gold">{{ formatPrice(product.price) }}</span>
              </div>
              <div v-if="product.price_old" class="flex flex-col">
                <span class="text-xs text-brand-brown/40 uppercase tracking-[0.2em] font-bold mb-1">Ориентировочно</span>
                <span class="text-xl text-brand-brown/20 line-through">{{ formatPrice(product.price_old) }}</span>
              </div>
            </div>
            <p class="text-xl text-brand-brown/70 leading-relaxed font-light">
              {{ product.description }}
            </p>
          </div>

          <!-- Details -->
          <div class="bg-brand-cream/50 p-10 rounded-[2rem] border border-brand-brown/5 mb-12 shadow-inner">
            <h3 class="font-serif text-2xl mb-8 flex items-center gap-3">
              <span class="w-8 h-px bg-brand-gold"></span>
              Детали проекта
            </h3>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-8">
              <div v-for="(value, key) in product.specs" :key="key" class="flex flex-col border-b border-brand-brown/5 pb-4">
                <span class="text-[10px] text-brand-gold uppercase tracking-[0.2em] font-black mb-1">{{ key }}</span>
                <span class="text-lg font-bold text-brand-brown">{{ value }}</span>
              </div>
            </div>
          </div>

          <!-- CTA -->
          <div class="mt-auto space-y-10">
            <button 
              @click="isOrderModalOpen = true"
              class="w-full bg-brand-brown text-white py-7 rounded-2xl text-xl font-bold hover:bg-brand-gold transition-all shadow-2xl hover:shadow-brand-gold/20 active:scale-[0.98] uppercase tracking-widest"
            >
              Хочу такой же проект
            </button>
            
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-8">
              <div v-for="b in [
                { icon: LucideShieldCheck, t: 'Гарантия 2 года' },
                { icon: LucideTruck, t: 'Монтаж в Крыму' },
                { icon: LucideCheckCircle, t: 'Технадзор' }
              ]" :key="b.t" class="flex items-center gap-3 text-[10px] font-black uppercase tracking-widest text-brand-brown/40">
                <component :is="b.icon" class="text-brand-gold" :size="18" />
                {{ b.t }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Related Projects -->
      <div v-if="relatedProjects.length" class="mt-40">
        <div class="flex items-end justify-between mb-12">
          <h2 class="font-serif text-4xl text-brand-brown">Вам может понравиться</h2>
          <router-link to="/catalog" class="text-brand-gold font-bold text-sm hover:underline flex items-center gap-2 group">
            Смотреть все
            <LucideArrowRight :size="18" class="group-hover:translate-x-1 transition-transform" />
          </router-link>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-10">
          <ProductCard v-for="rp in relatedProjects" :key="rp.id" :product="rp" />
        </div>
      </div>
    </div>

    <!-- Modals (Lightbox & Order) -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="isOrderModalOpen" class="fixed inset-0 z-[200] bg-black/95 backdrop-blur-xl flex items-center justify-center p-4" @click.self="isOrderModalOpen = false">
          <div class="relative bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl p-12 overflow-hidden transform transition-all">
            <button @click="isOrderModalOpen = false" class="absolute top-8 right-8 text-brand-brown/20 hover:text-brand-brown transition-colors">
              <LucideX :size="32" />
            </button>
            <h2 class="font-serif text-4xl mb-4 text-brand-brown">Заявка</h2>
            <p class="text-brand-brown/60 mb-10 font-medium">Обсудим ваш будущий проект?</p>
            <OrderForm :product-id="product.id" @success="isOrderModalOpen = false" />
          </div>
        </div>
      </transition>
    </Teleport>

    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLightboxOpen" class="fixed inset-0 z-[200] bg-black/95 backdrop-blur-xl flex items-center justify-center p-4 md:p-12" @click="isLightboxOpen = false">
          <button class="absolute top-8 right-8 text-white/40 hover:text-white transition-colors">
            <LucideX :size="40" />
          </button>
          <img :src="activeImage" class="max-w-full max-h-full object-contain shadow-2xl rounded-2xl border border-white/10 transition-all duration-500">
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1); }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>

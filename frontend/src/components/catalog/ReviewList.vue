<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { LucideStar, LucideCheckCircle, LucideMessageSquare, LucideX, LucideSearch } from 'lucide-vue-next';
import api from '@/api/client';
import type { ReviewResponse } from '@/types';

const props = defineProps<{
  projectId?: number;
}>();

const reviews = ref<ReviewResponse[]>([]);
const loading = ref(true);

// Lightbox state
const isLightboxOpen = ref(false);
const activeImage = ref('');

const openLightbox = (url: string) => {
  activeImage.value = url;
  isLightboxOpen.value = true;
};

const fetchReviews = async () => {
  try {
    const url = props.projectId 
      ? `/projects/${props.projectId}/reviews` 
      : '/reviews';
    const { data } = await api.get(url);
    reviews.value = data;
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchReviews);

defineExpose({ refresh: fetchReviews });
</script>

<template>
  <div class="space-y-12">
    <div v-if="loading" class="flex justify-center py-20">
      <div class="w-10 h-10 border-4 border-brand-gold border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="reviews.length === 0" class="text-center py-20 bg-brand-cream/20 rounded-[2rem] border border-dashed border-brand-brown/10">
      <LucideMessageSquare :size="48" class="mx-auto text-brand-brown/10 mb-4" />
      <p class="text-brand-brown/40 font-medium">Пока нет отзывов. Станьте первым!</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div 
        v-for="rev in reviews" :key="rev.id"
        class="bg-white p-10 rounded-[2.5rem] shadow-xl border border-brand-brown/5 flex flex-col hover:shadow-2xl transition-all duration-500"
      >
        <div class="flex justify-between items-start mb-8">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 bg-brand-cream rounded-2xl flex items-center justify-center text-brand-gold font-serif text-2xl font-bold">
              {{ rev.client_name[0] }}
            </div>
            <div>
              <h4 class="font-bold text-brand-brown text-lg">{{ rev.client_name }}</h4>
              <div class="flex items-center gap-2 text-green-600">
                <LucideCheckCircle :size="14" />
                <span class="text-[10px] font-black uppercase tracking-widest">Заказ подтвержден</span>
              </div>
            </div>
          </div>
          <div class="flex gap-0.5">
            <LucideStar 
              v-for="i in 5" :key="i"
              :size="16"
              :class="[ i <= rev.rating ? 'text-brand-gold fill-brand-gold' : 'text-brand-gray fill-brand-gray' ]"
            />
          </div>
        </div>

        <p class="text-brand-brown/70 leading-relaxed mb-8 italic">"{{ rev.comment }}"</p>

        <!-- Review Images -->
        <div v-if="rev.images && rev.images.length > 0" class="flex gap-3 mt-auto pt-6 border-t border-brand-brown/5 overflow-x-auto no-scrollbar">
          <div 
            v-for="(img, idx) in rev.images" :key="idx"
            @click="openLightbox(img.url)"
            class="relative w-24 h-24 rounded-2xl overflow-hidden shrink-0 border border-brand-brown/5 cursor-zoom-in group"
          >
            <img :src="img.url" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110">
            <div class="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors flex items-center justify-center">
              <LucideSearch :size="20" class="text-white opacity-0 group-hover:opacity-100 transition-opacity" />
            </div>
          </div>
        </div>
        
        <div class="mt-6 flex justify-between items-center text-[10px] font-black uppercase tracking-widest text-brand-brown/20">
          <span>{{ new Date(rev.created_at).toLocaleDateString() }}</span>
          <span v-if="rev.project_name" class="text-brand-gold">{{ rev.project_name }}</span>
        </div>
      </div>
    </div>

    <!-- Review Lightbox -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLightboxOpen" class="fixed inset-0 z-[250] bg-black/95 backdrop-blur-xl flex items-center justify-center p-4 md:p-12" @click="isLightboxOpen = false">
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

.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { LucideStar, LucideCheckCircle, LucideMessageSquare, LucideImage } from 'lucide-vue-next';
import api from '@/api/client';
import type { ReviewResponse } from '@/types';

const props = defineProps<{
  projectId?: number;
}>();

const reviews = ref<ReviewResponse[]>([]);
const loading = ref(true);

const fetchReviews = async () => {
  try {
    const url = props.projectId 
      ? `/projects/${props.projectId}/reviews` 
      : '/reviews'; // We'll need a global reviews route if we want all reviews on home
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
            class="w-20 h-20 rounded-xl overflow-hidden shrink-0 border border-brand-brown/5 cursor-pointer hover:scale-105 transition-transform"
          >
            <img :src="img.url" class="w-full h-full object-cover">
          </div>
        </div>
        
        <div class="mt-6 flex justify-between items-center text-[10px] font-black uppercase tracking-widest text-brand-brown/20">
          <span>{{ new Date(rev.created_at).toLocaleDateString() }}</span>
          <span v-if="rev.project_name" class="text-brand-gold">{{ rev.project_name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

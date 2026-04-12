<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { LucideCheckCircle, LucideMessageSquare, LucideSearch, LucideStar, LucideX } from 'lucide-vue-next';
import api from '@/api/client';
import type { ReviewResponse } from '@/types';
import { useNotificationStore } from '@/stores/notifications';
import { getApiErrorMessage } from '@/api/errors';

const props = defineProps<{
  projectId?: number;
}>();

const reviews = ref<ReviewResponse[]>([]);
const loading = ref(true);
const notificationStore = useNotificationStore();
const isLightboxOpen = ref(false);
const activeImage = ref('');

const openLightbox = (url: string) => {
  activeImage.value = url;
  isLightboxOpen.value = true;
};

const fetchReviews = async () => {
  loading.value = true;
  try {
    const url = props.projectId ? `/projects/${props.projectId}/reviews` : '/reviews';
    const { data } = await api.get(url);
    reviews.value = data;
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  } finally {
    loading.value = false;
  }
};

onMounted(fetchReviews);

defineExpose({ refresh: fetchReviews });
</script>

<template>
  <div>
    <div v-if="loading" class="ui-empty py-16">
      <div class="mx-auto h-9 w-9 animate-spin rounded-full border-4 border-brand-gold border-t-transparent"></div>
    </div>

    <div v-else-if="reviews.length === 0" class="ui-empty py-14">
      <LucideMessageSquare :size="44" class="mx-auto mb-4 text-brand-brown/12" />
      <p class="font-medium text-brand-brown/45">Отзывов пока нет.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-5 md:grid-cols-2">
      <article v-for="review in reviews" :key="review.id" class="ui-card p-5">
        <div class="mb-5 flex items-start justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-brand-gray font-serif text-xl font-bold text-brand-gold">
              {{ review.client_name?.[0] || '?' }}
            </div>
            <div>
              <h3 class="font-bold text-brand-brown">{{ review.client_name }}</h3>
              <div class="mt-1 flex items-center gap-2 text-green-700">
                <LucideCheckCircle :size="14" />
                <span class="text-[10px] font-black uppercase tracking-widest">Заказ проверен</span>
              </div>
            </div>
          </div>
          <div class="flex gap-0.5">
            <LucideStar
              v-for="i in 5"
              :key="i"
              :size="16"
              :class="[i <= review.rating ? 'fill-brand-gold text-brand-gold' : 'fill-brand-gray text-brand-gray']"
            />
          </div>
        </div>

        <p class="mb-5 leading-7 text-brand-brown/70">"{{ review.comment }}"</p>

        <div v-if="review.images?.length" class="mt-auto flex gap-3 overflow-x-auto border-t border-brand-brown/10 pt-4 no-scrollbar">
          <button
            v-for="image in review.images"
            :key="image.url"
            type="button"
            class="group relative h-20 w-20 shrink-0 overflow-hidden rounded-lg border border-brand-brown/10"
            @click="openLightbox(image.url)"
          >
            <img :src="image.url" class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-[1.035]" alt="">
            <span class="absolute inset-0 flex items-center justify-center bg-black/0 transition-colors group-hover:bg-black/20">
              <LucideSearch :size="18" class="text-white opacity-0 transition-opacity group-hover:opacity-100" />
            </span>
          </button>
        </div>

        <div class="mt-5 flex items-center justify-between gap-3 text-[11px] font-black uppercase tracking-widest text-brand-brown/25">
          <span>{{ new Date(review.created_at).toLocaleDateString('ru-RU') }}</span>
          <span v-if="review.project_name" class="text-brand-gold">{{ review.project_name }}</span>
        </div>
      </article>
    </div>

    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLightboxOpen" class="ui-modal-backdrop" @click="isLightboxOpen = false">
          <button type="button" class="absolute right-5 top-5 rounded-lg bg-white/10 p-3 text-white transition-colors hover:bg-white hover:text-brand-brown">
            <LucideX :size="28" />
          </button>
          <img :src="activeImage" class="z-10 max-h-full max-w-full rounded-lg object-contain shadow-2xl" alt="">
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

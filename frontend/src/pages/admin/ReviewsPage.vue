<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue';
import api from '@/api/client';
import {
  LucideCheckCircle,
  LucideChevronLeft,
  LucideChevronRight,
  LucideImage,
  LucideStar,
  LucideTrash2,
  LucideXCircle,
} from 'lucide-vue-next';
import type { ReviewResponse } from '@/types';
import { useNotificationStore } from '@/stores/notifications';
import { useConfirmStore } from '@/stores/confirm';
import { getApiErrorMessage } from '@/api/errors';

type ReviewStatus = ReviewResponse['status'];
type StatusFilter = ReviewStatus | 'all';

const reviews = ref<ReviewResponse[]>([]);
const total = ref(0);
const absoluteTotal = ref(0);
const statusFilter = ref<StatusFilter>('pending');
const notificationStore = useNotificationStore();
const confirmStore = useConfirmStore();
const currentPage = ref(1);
const limit = 10;
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)));

const statusOptions: Array<{ key: StatusFilter; label: string }> = [
  { key: 'pending', label: 'На проверке' },
  { key: 'approved', label: 'Опубликованы' },
  { key: 'rejected', label: 'Отклонены' },
  { key: 'all', label: 'Все' },
];

const statusMap: Record<ReviewStatus, string> = {
  pending: 'На проверке',
  approved: 'Опубликован',
  rejected: 'Отклонен',
};

const fetchReviews = async () => {
  try {
    const params: Record<string, string | number> = {
      limit,
      offset: (currentPage.value - 1) * limit,
    };

    if (statusFilter.value !== 'all') {
      params.status = statusFilter.value;
    }

    const { data } = await api.get('/admin/reviews', { params });
    reviews.value = data.items;
    total.value = data.total;
    absoluteTotal.value = data.absolute_total;
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

onMounted(fetchReviews);
watch([statusFilter, currentPage], fetchReviews);

const setStatusFilter = (status: StatusFilter) => {
  statusFilter.value = status;
  currentPage.value = 1;
};

const moderate = async (id: number, approved: boolean) => {
  try {
    await api.patch(`/admin/reviews/${id}/status`, { approved });
    notificationStore.show(approved ? 'Отзыв опубликован' : 'Отзыв отклонен', 'success');
    fetchReviews();
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

const deleteReview = async (id: number) => {
  const confirmed = await confirmStore.request({
    title: 'Удалить отзыв?',
    message: 'Отзыв будет удален навсегда и исчезнет из истории модерации.',
    confirmLabel: 'Удалить',
    tone: 'danger',
  });

  if (!confirmed) return;

  try {
    await api.delete(`/admin/reviews/${id}`);
    notificationStore.show('Отзыв удален', 'info');
    fetchReviews();
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

const formatDate = (value: string) => {
  return new Date(value).toLocaleString('ru-RU', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const getStatusClass = (status: ReviewStatus) => {
  if (status === 'approved') return 'bg-green-50 text-green-700 ring-green-100';
  if (status === 'rejected') return 'bg-red-50 text-red-700 ring-red-100';
  return 'bg-yellow-50 text-yellow-700 ring-yellow-100';
};
</script>

<template>
  <div class="space-y-6">
    <section>
      <p class="text-xs font-black uppercase tracking-widest text-brand-gold">Доверие и репутация</p>
      <h1 class="mt-1 font-serif text-3xl font-bold leading-tight text-brand-brown sm:text-4xl">Модерация отзывов</h1>
      <p class="mt-2 text-sm text-brand-brown/55">Всего отзывов в системе: {{ absoluteTotal }}</p>
    </section>

    <section class="-mx-4 overflow-x-auto px-4 sm:mx-0 sm:px-0 no-scrollbar" aria-label="Фильтр отзывов по статусу">
      <div class="flex min-w-max gap-2 rounded-lg border border-brand-brown/10 bg-white p-1.5 shadow-sm">
        <button
          v-for="status in statusOptions"
          :key="status.key"
          type="button"
          @click="setStatusFilter(status.key)"
          :class="[
            'rounded-lg px-4 py-2.5 text-xs font-black uppercase tracking-widest transition-all',
            statusFilter === status.key
              ? 'bg-brand-brown text-white shadow-md'
              : 'text-brand-brown/45 hover:bg-brand-gray hover:text-brand-brown',
          ]"
        >
          {{ status.label }}
        </button>
      </div>
    </section>

    <section v-if="reviews.length" class="space-y-3" aria-label="Отзывы">
      <article
        v-for="review in reviews"
        :key="review.id"
        class="rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm transition-all hover:shadow-md sm:p-5"
      >
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <div class="flex gap-0.5">
                <LucideStar
                  v-for="i in 5"
                  :key="i"
                  :size="16"
                  :class="[i <= review.rating ? 'fill-brand-gold text-brand-gold' : 'fill-brand-gray text-brand-gray']"
                />
              </div>
              <span :class="['rounded-full px-2.5 py-1 text-[11px] font-black uppercase ring-1', getStatusClass(review.status)]">
                {{ statusMap[review.status] }}
              </span>
              <span class="text-xs text-brand-brown/40">{{ formatDate(review.created_at) }}</span>
            </div>

            <div class="mt-3">
              <h2 class="font-bold text-brand-brown">{{ review.client_name }}</h2>
              <p class="text-sm font-semibold text-brand-gold">{{ review.project_name || 'Общий отзыв' }}</p>
            </div>

            <p class="mt-4 leading-7 text-brand-brown/70">
              "{{ review.comment }}"
            </p>

            <div v-if="review.images?.length" class="mt-4 flex flex-wrap gap-2">
              <a
                v-for="image in review.images"
                :key="image.url"
                :href="image.url"
                target="_blank"
                class="block h-16 w-16 overflow-hidden rounded-lg border border-brand-brown/10 bg-brand-gray"
                aria-label="Открыть фото отзыва"
              >
                <img :src="image.url" class="h-full w-full object-cover" alt="">
              </a>
            </div>
            <div v-else class="mt-4 inline-flex items-center gap-2 text-xs font-semibold text-brand-brown/35">
              <LucideImage :size="15" />
              Без фотографий
            </div>
          </div>

          <div class="grid gap-2 sm:grid-cols-3 lg:w-44 lg:grid-cols-1">
            <button
              v-if="review.status !== 'approved'"
              type="button"
              @click="moderate(review.id, true)"
              class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg bg-green-50 px-3 py-2 text-sm font-bold text-green-700 transition-all hover:bg-green-600 hover:text-white"
            >
              <LucideCheckCircle :size="17" />
              Опубликовать
            </button>
            <button
              v-if="review.status !== 'rejected'"
              type="button"
              @click="moderate(review.id, false)"
              class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg bg-yellow-50 px-3 py-2 text-sm font-bold text-yellow-700 transition-all hover:bg-yellow-600 hover:text-white"
            >
              <LucideXCircle :size="17" />
              Отклонить
            </button>
            <button
              type="button"
              @click="deleteReview(review.id)"
              class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg bg-red-50 px-3 py-2 text-sm font-bold text-red-700 transition-all hover:bg-red-600 hover:text-white"
            >
              <LucideTrash2 :size="17" />
              Удалить
            </button>
          </div>
        </div>
      </article>
    </section>

    <section
      v-if="reviews.length === 0"
      class="rounded-lg border border-dashed border-brand-brown/15 bg-white p-10 text-center text-brand-brown/35"
    >
      Отзывы не найдены
    </section>

    <section class="flex flex-col gap-3 rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between">
      <span class="text-sm text-brand-brown/50">
        Показано {{ reviews.length }} из {{ total }} в текущем фильтре
      </span>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="flex h-10 w-10 items-center justify-center rounded-lg border border-brand-brown/10 bg-white text-brand-brown transition-colors hover:bg-brand-gray disabled:opacity-30"
          aria-label="Предыдущая страница"
        >
          <LucideChevronLeft :size="20" />
        </button>
        <div class="flex h-10 min-w-10 items-center justify-center rounded-lg bg-brand-brown px-3 text-sm font-bold text-white">
          {{ currentPage }} / {{ totalPages }}
        </div>
        <button
          type="button"
          @click="currentPage++"
          :disabled="currentPage >= totalPages"
          class="flex h-10 w-10 items-center justify-center rounded-lg border border-brand-brown/10 bg-white text-brand-brown transition-colors hover:bg-brand-gray disabled:opacity-30"
          aria-label="Следующая страница"
        >
          <LucideChevronRight :size="20" />
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import api from '@/api/client';
import { 
  LucideCheckCircle, 
  LucideXCircle, 
  LucideTrash2, 
  LucideStar,
  LucideClock,
  LucideChevronLeft,
  LucideChevronRight,
  LucideImage
} from 'lucide-vue-next';
import { useNotificationStore } from '@/stores/notifications';

const reviews = ref<any[]>([]);
const total = ref(0);
const absoluteTotal = ref(0);
const statusFilter = ref('pending');
const notificationStore = useNotificationStore();
const currentPage = ref(1);
const limit = 10;

const fetchReviews = async () => {
  try {
    const params: any = { 
      limit, 
      offset: (currentPage.value - 1) * limit 
    };
    if (statusFilter.value !== 'all') params.status = statusFilter.value;
    const { data } = await api.get('/admin/reviews', { params });
    reviews.value = data.items;
    total.value = data.total;
    absoluteTotal.value = data.absolute_total;
  } catch (err) {
    console.error(err);
  }
};

onMounted(fetchReviews);
watch([statusFilter, currentPage], fetchReviews);

const moderate = async (id: number, approved: boolean) => {
  try {
    await api.patch(`/admin/reviews/${id}/status`, { approved });
    notificationStore.show(approved ? 'Отзыв опубликован' : 'Отзыв отклонен', 'success');
    fetchReviews();
  } catch (err) {
    notificationStore.show('Ошибка модерации', 'error');
  }
};

const deleteReview = async (id: number) => {
  if (confirm('Удалить отзыв навсегда?')) {
    try {
      await api.delete(`/admin/reviews/${id}`);
      notificationStore.show('Отзыв удален', 'info');
      fetchReviews();
    } catch (err) {
      notificationStore.show('Ошибка удаления', 'error');
    }
  }
};
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-12">
      <div>
        <h1 class="font-serif text-4xl text-brand-brown mb-2">Модерация отзывов</h1>
        <p class="text-brand-brown/40">Всего отзывов в системе: {{ absoluteTotal }}</p>
      </div>
    </div>

    <div class="flex gap-2 mb-8 bg-white p-1.5 rounded-2xl border border-brand-brown/5 w-fit shadow-sm">
      <button 
        v-for="s in [
          { k: 'pending', l: 'На проверке' },
          { k: 'approved', l: 'Опубликованы' },
          { k: 'rejected', l: 'Отклонены' },
          { k: 'all', l: 'Все' }
        ]" :key="s.k"
        @click="statusFilter = s.k; currentPage = 1"
        :class="['px-6 py-2.5 rounded-xl text-sm font-bold uppercase tracking-widest transition-all', statusFilter === s.k ? 'bg-brand-brown text-white shadow-md' : 'text-brand-brown/40 hover:bg-brand-gray']"
      >
        {{ s.l }}
      </button>
    </div>

    <div class="space-y-6">
      <div 
        v-for="rev in reviews" :key="rev.id"
        class="bg-white p-8 rounded-3xl shadow-sm border border-brand-brown/5 flex flex-col md:flex-row gap-8 items-start hover:shadow-md transition-all"
      >
        <!-- Client Info -->
        <div class="md:w-1/4 space-y-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-brand-cream rounded-xl flex items-center justify-center text-brand-gold font-bold">
              {{ rev.client_name ? rev.client_name[0] : '?' }}
            </div>
            <div>
              <div class="font-bold text-brand-brown">{{ rev.client_name }}</div>
              <div class="text-[10px] uppercase font-black text-brand-brown/20 tracking-widest">Клиент</div>
            </div>
          </div>
          <div class="p-4 bg-brand-gray/30 rounded-2xl">
            <div class="text-[10px] uppercase font-black text-brand-brown/40 mb-1">Проект</div>
            <div class="text-sm font-bold text-brand-brown">{{ rev.project_name || 'Общий отзыв' }}</div>
          </div>
        </div>

        <!-- Content -->
        <div class="flex-1 space-y-4">
          <div class="flex items-center gap-4">
            <div class="flex gap-0.5">
              <LucideStar 
                v-for="i in 5" :key="i"
                :size="16"
                :class="[ i <= rev.rating ? 'text-brand-gold fill-brand-gold' : 'text-brand-gray fill-brand-gray' ]"
              />
            </div>
            <span class="text-xs text-brand-brown/30">{{ new Date(rev.created_at).toLocaleString() }}</span>
          </div>
          <p class="text-brand-brown/70 leading-relaxed italic">"{{ rev.comment }}"</p>
          
          <div v-if="rev.images?.length" class="flex gap-2">
            <div v-for="img in rev.images" :key="img.url" class="w-16 h-16 rounded-lg overflow-hidden border border-brand-brown/5">
              <img :src="img.url" class="w-full h-full object-cover">
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex md:flex-col gap-3">
          <button 
            v-if="rev.status !== 'approved'"
            @click="moderate(rev.id, true)"
            class="p-3 bg-green-50 text-green-600 rounded-xl hover:bg-green-100 transition-all shadow-sm" title="Опубликовать"
          >
            <LucideCheckCircle :size="20" />
          </button>
          <button 
            v-if="rev.status !== 'rejected'"
            @click="moderate(rev.id, false)"
            class="p-3 bg-yellow-50 text-yellow-600 rounded-xl hover:bg-yellow-100 transition-all shadow-sm" title="Отклонить"
          >
            <LucideXCircle :size="20" />
          </button>
          <button 
            @click="deleteReview(rev.id)"
            class="p-3 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition-all shadow-sm" title="Удалить навсегда"
          >
            <LucideTrash2 :size="20" />
          </button>
        </div>
      </div>

      <div v-if="reviews.length === 0" class="py-20 text-center bg-white rounded-3xl border border-dashed border-brand-brown/10 text-brand-brown/20 italic">
        Отзывов не найдено
      </div>

      <!-- Pagination -->
      <div class="flex justify-between items-center bg-white p-6 rounded-3xl shadow-sm border border-brand-brown/5">
        <span class="text-sm text-brand-brown/40">Показано {{ reviews.length }} из {{ total }} в текущем фильтре</span>
        <div class="flex gap-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="p-2 rounded-lg bg-brand-gray/50 disabled:opacity-30"><LucideChevronLeft :size="20"/></button>
          <div class="px-4 py-2 bg-brand-brown text-white rounded-lg font-bold">{{ currentPage }}</div>
          <button @click="currentPage++" :disabled="currentPage * limit >= total" class="p-2 rounded-lg bg-brand-gray/50 disabled:opacity-30"><LucideChevronRight :size="20"/></button>
        </div>
      </div>
    </div>
  </div>
</template>

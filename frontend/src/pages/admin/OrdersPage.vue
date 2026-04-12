<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import api from '@/api/client';
import { 
  LucideCheckCircle, 
  LucideXCircle, 
  LucideClock, 
  LucideDownload,
  LucidePlayCircle,
  LucideChevronLeft,
  LucideChevronRight,
  LucideBan,
  LucideRefreshCw
} from 'lucide-vue-next';
import type { Order } from '@/types';
import { useNotificationStore } from '@/stores/notifications';
import { useConfirmStore } from '@/stores/confirm';
import { getApiErrorMessage } from '@/api/errors';
import { downloadFile } from '@/utils/download';

const orders = ref<Order[]>([]);
const total = ref(0);
const absoluteTotal = ref(0);
const statusFilter = ref('new');
const notificationStore = useNotificationStore();
const confirmStore = useConfirmStore();

const currentPage = ref(1);
const limit = 10;

const statusMap: Record<string, string> = {
  'new': 'Новая',
  'processing': 'В работе',
  'done': 'Завершена',
  'rejected': 'Отклонена',
  'spam': 'Спам'
};

const contactMethodMap: Record<string, string> = {
  phone: 'Звонок',
  whatsapp: 'WhatsApp',
  telegram: 'Telegram',
  email: 'Email',
};

const briefItems = (order: Order) => {
  return [
    order.project_type,
    order.budget_range,
    order.city,
    contactMethodMap[order.contact_method] || order.contact_method,
  ].filter(Boolean);
};

const fetchOrders = async () => {
  try {
    const params: any = { 
      limit, 
      offset: (currentPage.value - 1) * limit 
    };
    if (statusFilter.value !== 'all') params.status = statusFilter.value;
    
    const { data } = await api.get('/admin/orders', { params });
    orders.value = data.items;
    total.value = data.total;
    absoluteTotal.value = data.absolute_total;
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

onMounted(fetchOrders);
watch([statusFilter, currentPage], fetchOrders);

const updateStatus = async (id: number, status: string) => {
  try {
    await api.patch(`/admin/orders/${id}/status`, { status });
    notificationStore.show(`Статус заявки #${id} изменен на "${statusMap[status]}"`, 'success');
    fetchOrders();
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

const markAsSpam = async (id: number) => {
  const confirmed = await confirmStore.request({
    title: 'Пометить заявку как спам?',
    message: 'IP-адрес клиента будет временно заблокирован, а заявка уйдет в раздел спама.',
    confirmLabel: 'Заблокировать',
    tone: 'danger',
  });
  if (!confirmed) return;

  try {
    await api.post(`/admin/orders/${id}/spam`);
    notificationStore.show('Заявка помечена как спам', 'info');
    fetchOrders();
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

const exportExcel = () => {
  downloadFile('/admin/orders/export', 'orders.xlsx');
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'new': return LucideClock;
    case 'processing': return LucidePlayCircle;
    case 'done': return LucideCheckCircle;
    case 'rejected': return LucideXCircle;
    case 'spam': return LucideBan;
    default: return LucideClock;
  }
};

const getStatusClass = (status: string) => {
  switch (status) {
    case 'new': return 'bg-blue-100 text-blue-700';
    case 'processing': return 'bg-yellow-100 text-yellow-700';
    case 'done': return 'bg-green-100 text-green-700';
    case 'rejected': return 'bg-red-100 text-red-700';
    case 'spam': return 'bg-gray-100 text-gray-700';
    default: return 'bg-gray-100 text-gray-700';
  }
};
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-12">
      <div>
        <h1 class="font-serif text-4xl text-brand-brown mb-2">Заявки клиентов</h1>
        <p class="text-brand-brown/40">Всего заявок в системе: {{ absoluteTotal }}</p>
      </div>
      <button @click="exportExcel" class="bg-brand-brown text-white px-6 py-3 rounded-lg font-medium hover:bg-brand-gold transition-all flex items-center gap-2 shadow-lg">
        <LucideDownload :size="20" />
        Экспорт EXCEL
      </button>
    </div>

    <!-- Status Filters -->
    <div class="flex gap-2 mb-8 bg-white p-1.5 rounded-lg border border-brand-brown/5 w-fit shadow-sm overflow-x-auto max-w-full">
      <button 
        v-for="s in ['new', 'processing', 'done', 'rejected', 'spam', 'all']" 
        :key="s"
        @click="statusFilter = s; currentPage = 1"
        :class="['px-6 py-2.5 rounded-lg text-sm font-bold uppercase tracking-widest transition-all whitespace-nowrap', statusFilter === s ? 'bg-brand-brown text-white shadow-md' : 'text-brand-brown/40 hover:bg-brand-gray']"
      >
        {{ s === 'all' ? 'Все' : (statusMap[s] || s) }}
      </button>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-lg shadow-sm border border-brand-brown/5 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-[980px] w-full text-left border-collapse">
        <thead>
          <tr class="bg-brand-gray/20 text-brand-brown/40 text-xs uppercase tracking-widest">
            <th class="px-8 py-4 font-semibold">Клиент</th>
            <th class="px-8 py-4 font-semibold">Проект</th>
            <th class="px-8 py-4 font-semibold">Контакты</th>
            <th class="px-8 py-4 font-semibold">Статус</th>
            <th class="px-8 py-4 font-semibold">Дата</th>
            <th class="px-8 py-4 font-semibold text-right">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-brand-brown/5">
          <tr v-for="o in orders" :key="o.id" class="hover:bg-brand-gray/5 transition-colors group">
            <!-- 1. Клиент -->
            <td class="px-8 py-6">
              <div class="font-bold text-brand-brown flex items-center gap-2">
                <span class="text-brand-brown/20 font-mono text-xs">#{{ o.id }}</span>
                {{ o.client_name }}
              </div>
              <div class="text-sm text-brand-brown/60 mt-1 max-w-xs italic line-clamp-1">"{{ o.comment || 'Без комментария' }}"</div>
              <div v-if="briefItems(o).length" class="mt-3 flex flex-wrap gap-2">
                <span
                  v-for="item in briefItems(o)"
                  :key="item"
                  class="rounded-lg bg-brand-gray px-2.5 py-1 text-[11px] font-semibold text-brand-brown/60"
                >
                  {{ item }}
                </span>
              </div>
            </td>

            <!-- 2. Проект -->
            <td class="px-8 py-6">
              <div v-if="o.project_name" class="flex flex-col">
                <span class="text-sm font-bold text-brand-brown">{{ o.project_name }}</span>
                <span class="text-[10px] uppercase font-black text-brand-gold tracking-widest">Индивидуальный заказ</span>
              </div>
              <div v-else class="text-brand-brown/30 text-xs italic">Общая консультация</div>
            </td>

            <!-- 3. Контакты -->
            <td class="px-8 py-6">
              <div class="font-medium text-brand-brown whitespace-nowrap">{{ o.client_phone }}</div>
              <div class="text-xs text-brand-brown/40">{{ o.client_email || 'email не указан' }}</div>
              <div v-if="o.contact_method" class="mt-1 text-[10px] uppercase tracking-widest text-brand-gold">
                {{ contactMethodMap[o.contact_method] || o.contact_method }}
              </div>
            </td>

            <!-- 4. Статус -->
            <td class="px-8 py-6">
              <div :class="['inline-flex items-center gap-2 px-3 py-1 rounded-full text-[10px] font-black uppercase', getStatusClass(o.status)]">
                <component :is="getStatusIcon(o.status)" :size="12" />
                {{ statusMap[o.status] || o.status }}
              </div>
            </td>

            <!-- 5. Дата -->
            <td class="px-8 py-6 text-sm text-brand-brown/40 whitespace-nowrap">
              {{ new Date(o.created_at).toLocaleDateString() }}<br>
              <span class="text-[10px] uppercase font-bold">{{ new Date(o.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</span>
            </td>

            <!-- 6. Действия (Интеллектуальные кнопки) -->
            <td class="px-8 py-6 text-right min-w-[180px]">
              <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                
                <!-- Для НОВЫХ: В работу, Отклонить, Спам -->
                <template v-if="o.status === 'new'">
                  <button @click="updateStatus(o.id, 'processing')" class="p-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-600 hover:text-white transition-all" title="В работу">
                    <LucidePlayCircle :size="18" />
                  </button>
                  <button @click="updateStatus(o.id, 'rejected')" class="p-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-600 hover:text-white transition-all" title="Отклонить">
                    <LucideXCircle :size="18" />
                  </button>
                  <button @click="markAsSpam(o.id)" class="p-2 bg-gray-100 text-gray-600 rounded-lg hover:bg-gray-600 hover:text-white transition-all" title="Это спам">
                    <LucideBan :size="18" />
                  </button>
                </template>

                <!-- Для В РАБОТЕ: Только Завершить -->
                <template v-else-if="o.status === 'processing'">
                  <button @click="updateStatus(o.id, 'done')" class="p-2 bg-green-50 text-green-600 rounded-lg hover:bg-green-600 hover:text-white transition-all flex items-center gap-2 px-4" title="Завершить проект">
                    <LucideCheckCircle :size="18" />
                    <span class="text-[10px] font-bold uppercase">Завершить</span>
                  </button>
                </template>

                <!-- Для СПАМА: Вернуть в новые -->
                <template v-else-if="o.status === 'spam'">
                  <button @click="updateStatus(o.id, 'new')" class="p-2 bg-brand-gold/10 text-brand-gold rounded-lg hover:bg-brand-gold hover:text-white transition-all flex items-center gap-2 px-4" title="Восстановить">
                    <LucideRefreshCw :size="18" />
                    <span class="text-[10px] font-bold uppercase">Восстановить</span>
                  </button>
                </template>

                <!-- Для ЗАВЕРШЕННЫХ и ОТКЛОНЕННЫХ: Ничего (пусто) -->
                <template v-else>
                  <span class="text-[10px] text-brand-brown/20 font-bold uppercase italic tracking-widest">Архив</span>
                </template>

              </div>
            </td>
          </tr>
        </tbody>
        </table>
      </div>
      
      <div v-if="orders.length === 0" class="p-20 text-center text-brand-brown/20 italic">
        Заявок с таким статусом не найдено
      </div>

      <!-- Pagination -->
      <div class="p-6 bg-brand-gray/10 flex items-center justify-between">
        <span class="text-sm text-brand-brown/40">Показано {{ orders.length }} из {{ total }} в текущем фильтре</span>
        <div class="flex gap-2">
          <button 
            @click="currentPage--" 
            :disabled="currentPage === 1"
            class="p-2 rounded-lg bg-white border border-brand-brown/5 disabled:opacity-30"
          >
            <LucideChevronLeft :size="20" />
          </button>
          <div class="px-4 py-2 bg-brand-brown text-white rounded-lg font-bold text-sm">{{ currentPage }}</div>
          <button 
            @click="currentPage++" 
            :disabled="currentPage * limit >= total"
            class="p-2 rounded-lg bg-white border border-brand-brown/5 disabled:opacity-30"
          >
            <LucideChevronRight :size="20" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

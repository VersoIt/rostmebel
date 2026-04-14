<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue';
import type { Component } from 'vue';
import api from '@/api/client';
import {
  LucideBan,
  LucideCheckCircle,
  LucideChevronLeft,
  LucideChevronRight,
  LucideClock,
  LucideDownload,
  LucideMail,
  LucidePhone,
  LucidePlayCircle,
  LucideRefreshCw,
  LucideXCircle,
} from 'lucide-vue-next';
import type { Order, OrderStatus } from '@/types';
import { useNotificationStore } from '@/stores/notifications';
import { useConfirmStore } from '@/stores/confirm';
import { getApiErrorMessage } from '@/api/errors';
import { downloadFile } from '@/utils/download';
import { contactMethodLabels } from '@/utils/orderOptions';

type StatusFilter = OrderStatus | 'all';

type OrderAction = {
  key: string;
  label: string;
  title: string;
  icon: Component;
  tone: 'blue' | 'green' | 'red' | 'gray' | 'gold';
  kind: 'status' | 'spam';
  status?: OrderStatus;
};

const orders = ref<Order[]>([]);
const total = ref(0);
const absoluteTotal = ref(0);
const statusFilter = ref<StatusFilter>('new');
const notificationStore = useNotificationStore();
const confirmStore = useConfirmStore();

const currentPage = ref(1);
const limit = 10;
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)));

const statusOptions: StatusFilter[] = ['new', 'processing', 'done', 'rejected', 'spam', 'all'];

const statusMap: Record<OrderStatus, string> = {
  new: 'Новая',
  processing: 'В работе',
  done: 'Завершена',
  rejected: 'Отклонена',
  spam: 'Спам',
};

const briefItems = (order: Order) => {
  return [
    order.project_type,
    order.budget_range,
    order.city,
    contactMethodLabels[order.contact_method] || order.contact_method,
  ].filter(Boolean);
};

const formatDate = (value: string) => {
  return new Date(value).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });
};

const formatTime = (value: string) => {
  return new Date(value).toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  });
};

const fetchOrders = async () => {
  try {
    const params: Record<string, string | number> = {
      limit,
      offset: (currentPage.value - 1) * limit,
    };

    if (statusFilter.value !== 'all') {
      params.status = statusFilter.value;
    }

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

const setStatusFilter = (status: StatusFilter) => {
  statusFilter.value = status;
  currentPage.value = 1;
};

const updateStatus = async (id: number, status: OrderStatus) => {
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

const getStatusIcon = (status: OrderStatus) => {
  switch (status) {
    case 'new':
      return LucideClock;
    case 'processing':
      return LucidePlayCircle;
    case 'done':
      return LucideCheckCircle;
    case 'rejected':
      return LucideXCircle;
    case 'spam':
      return LucideBan;
    default:
      return LucideClock;
  }
};

const getStatusClass = (status: OrderStatus) => {
  switch (status) {
    case 'new':
      return 'bg-blue-50 text-blue-700 ring-blue-100';
    case 'processing':
      return 'bg-yellow-50 text-yellow-700 ring-yellow-100';
    case 'done':
      return 'bg-green-50 text-green-700 ring-green-100';
    case 'rejected':
      return 'bg-red-50 text-red-700 ring-red-100';
    case 'spam':
      return 'bg-gray-100 text-gray-700 ring-gray-200';
    default:
      return 'bg-gray-100 text-gray-700 ring-gray-200';
  }
};

const getActionClass = (tone: OrderAction['tone']) => {
  const base = 'inline-flex min-h-11 items-center justify-center gap-2 rounded-lg px-3 py-2 text-sm font-bold transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50';

  switch (tone) {
    case 'blue':
      return `${base} bg-blue-50 text-blue-700 hover:bg-blue-600 hover:text-white focus:ring-blue-300`;
    case 'green':
      return `${base} bg-green-50 text-green-700 hover:bg-green-600 hover:text-white focus:ring-green-300`;
    case 'red':
      return `${base} bg-red-50 text-red-700 hover:bg-red-600 hover:text-white focus:ring-red-300`;
    case 'gray':
      return `${base} bg-gray-100 text-gray-700 hover:bg-gray-700 hover:text-white focus:ring-gray-300`;
    case 'gold':
      return `${base} bg-brand-gold/10 text-brand-gold hover:bg-brand-gold hover:text-white focus:ring-brand-gold/30`;
    default:
      return base;
  }
};

const getOrderActions = (order: Order): OrderAction[] => {
  if (order.status === 'new') {
    return [
      {
        key: 'processing',
        label: 'В работу',
        title: 'Перевести заявку в работу',
        icon: LucidePlayCircle,
        tone: 'blue',
        kind: 'status',
        status: 'processing',
      },
      {
        key: 'rejected',
        label: 'Отклонить',
        title: 'Отклонить заявку',
        icon: LucideXCircle,
        tone: 'red',
        kind: 'status',
        status: 'rejected',
      },
      {
        key: 'spam',
        label: 'Спам',
        title: 'Пометить заявку как спам',
        icon: LucideBan,
        tone: 'gray',
        kind: 'spam',
      },
    ];
  }

  if (order.status === 'processing') {
    return [
      {
        key: 'done',
        label: 'Завершить',
        title: 'Завершить заявку',
        icon: LucideCheckCircle,
        tone: 'green',
        kind: 'status',
        status: 'done',
      },
    ];
  }

  if (order.status === 'spam') {
    return [
      {
        key: 'restore',
        label: 'Восстановить',
        title: 'Вернуть заявку в новые',
        icon: LucideRefreshCw,
        tone: 'gold',
        kind: 'status',
        status: 'new',
      },
    ];
  }

  return [];
};

const runOrderAction = (order: Order, action: OrderAction) => {
  if (action.kind === 'spam') {
    markAsSpam(order.id);
    return;
  }

  if (action.status) {
    updateStatus(order.id, action.status);
  }
};
</script>

<template>
  <div class="space-y-6">
    <section class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <p class="text-xs font-black uppercase tracking-widest text-brand-gold">Операционный центр</p>
        <h1 class="mt-1 font-serif text-3xl font-bold leading-tight text-brand-brown sm:text-4xl">Заявки клиентов</h1>
        <p class="mt-2 text-sm text-brand-brown/55">Всего заявок в системе: {{ absoluteTotal }}</p>
      </div>

      <button
        type="button"
        @click="exportExcel"
        class="inline-flex h-12 w-full items-center justify-center gap-2 rounded-lg bg-brand-brown px-5 text-sm font-bold text-white shadow-lg shadow-brand-brown/10 transition-all hover:bg-brand-gold sm:w-auto"
      >
        <LucideDownload :size="19" />
        Экспорт Excel
      </button>
    </section>

    <section class="-mx-4 overflow-x-auto px-4 sm:mx-0 sm:px-0 no-scrollbar" aria-label="Фильтр заявок по статусу">
      <div class="flex min-w-max gap-2 rounded-lg border border-brand-brown/10 bg-white p-1.5 shadow-sm">
        <button
          v-for="status in statusOptions"
          :key="status"
          type="button"
          @click="setStatusFilter(status)"
          :class="[
            'rounded-lg px-4 py-2.5 text-xs font-black uppercase tracking-widest transition-all',
            statusFilter === status
              ? 'bg-brand-brown text-white shadow-md'
              : 'text-brand-brown/45 hover:bg-brand-gray hover:text-brand-brown',
          ]"
        >
          {{ status === 'all' ? 'Все' : statusMap[status] }}
        </button>
      </div>
    </section>

    <section v-if="orders.length" class="space-y-3 xl:hidden" aria-label="Заявки">
      <article
        v-for="order in orders"
        :key="order.id"
        class="rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="mb-1 flex flex-wrap items-center gap-2">
              <span class="font-mono text-xs font-bold text-brand-brown/35">#{{ order.id }}</span>
              <span :class="['inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-black uppercase ring-1', getStatusClass(order.status)]">
                <component :is="getStatusIcon(order.status)" :size="12" />
                {{ statusMap[order.status] }}
              </span>
            </div>
            <h2 class="truncate text-lg font-bold text-brand-brown">{{ order.client_name }}</h2>
          </div>

          <div class="shrink-0 text-right text-xs text-brand-brown/45">
            <div class="font-semibold">{{ formatDate(order.created_at) }}</div>
            <div>{{ formatTime(order.created_at) }}</div>
          </div>
        </div>

        <p class="mt-3 line-clamp-3 text-sm leading-relaxed text-brand-brown/65">
          {{ order.comment || 'Комментарий не указан' }}
        </p>

        <div class="mt-4 grid gap-3 text-sm sm:grid-cols-2">
          <div>
            <div class="text-[11px] font-black uppercase tracking-widest text-brand-brown/35">Проект</div>
            <div class="mt-1 font-semibold text-brand-brown">{{ order.project_name || 'Общая консультация' }}</div>
          </div>

          <div>
            <div class="text-[11px] font-black uppercase tracking-widest text-brand-brown/35">Контакты</div>
            <div class="mt-1 flex flex-wrap gap-2">
              <a
                :href="`tel:${order.client_phone}`"
                class="inline-flex items-center gap-1.5 rounded-lg bg-brand-gray px-2.5 py-1.5 font-semibold text-brand-brown"
              >
                <LucidePhone :size="14" />
                {{ order.client_phone }}
              </a>
              <a
                v-if="order.client_email"
                :href="`mailto:${order.client_email}`"
                class="inline-flex items-center gap-1.5 rounded-lg bg-brand-gray px-2.5 py-1.5 font-semibold text-brand-brown"
              >
                <LucideMail :size="14" />
                Email
              </a>
            </div>
          </div>
        </div>

        <div v-if="briefItems(order).length" class="mt-4 flex flex-wrap gap-2">
          <span
            v-for="item in briefItems(order)"
            :key="item"
            class="rounded-lg bg-brand-gray px-2.5 py-1 text-[11px] font-bold text-brand-brown/60"
          >
            {{ item }}
          </span>
        </div>

        <div class="mt-4 border-t border-brand-brown/10 pt-4">
          <div class="mb-2 text-[11px] font-black uppercase tracking-widest text-brand-brown/35">Действия</div>
          <div v-if="getOrderActions(order).length" class="grid gap-2 sm:grid-cols-3">
            <button
              v-for="action in getOrderActions(order)"
              :key="action.key"
              type="button"
              :title="action.title"
              :class="getActionClass(action.tone)"
              @click="runOrderAction(order, action)"
            >
              <component :is="action.icon" :size="17" />
              {{ action.label }}
            </button>
          </div>
          <div v-else class="rounded-lg bg-brand-gray px-3 py-2 text-sm font-semibold text-brand-brown/45">
            Заявка в архивном состоянии
          </div>
        </div>
      </article>
    </section>

    <section v-if="orders.length" class="hidden overflow-hidden rounded-lg border border-brand-brown/10 bg-white shadow-sm xl:block" aria-label="Таблица заявок">
      <div class="overflow-x-auto">
        <table class="w-full min-w-[1040px] border-collapse text-left">
          <thead>
            <tr class="bg-brand-gray/50 text-xs uppercase tracking-widest text-brand-brown/45">
              <th class="px-6 py-4 font-semibold">Клиент</th>
              <th class="px-6 py-4 font-semibold">Проект</th>
              <th class="px-6 py-4 font-semibold">Контакты</th>
              <th class="px-6 py-4 font-semibold">Статус</th>
              <th class="px-6 py-4 font-semibold">Дата</th>
              <th class="px-6 py-4 text-right font-semibold">Действия</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-brand-brown/10">
            <tr v-for="order in orders" :key="order.id" class="transition-colors hover:bg-brand-gray/30">
              <td class="px-6 py-5 align-top">
                <div class="flex items-center gap-2 font-bold text-brand-brown">
                  <span class="font-mono text-xs text-brand-brown/30">#{{ order.id }}</span>
                  {{ order.client_name }}
                </div>
                <div class="mt-1 max-w-sm truncate text-sm italic text-brand-brown/55">
                  "{{ order.comment || 'Комментарий не указан' }}"
                </div>
                <div v-if="briefItems(order).length" class="mt-3 flex flex-wrap gap-2">
                  <span
                    v-for="item in briefItems(order)"
                    :key="item"
                    class="rounded-lg bg-brand-gray px-2.5 py-1 text-[11px] font-semibold text-brand-brown/60"
                  >
                    {{ item }}
                  </span>
                </div>
              </td>

              <td class="px-6 py-5 align-top">
                <div v-if="order.project_name" class="flex flex-col">
                  <span class="text-sm font-bold text-brand-brown">{{ order.project_name }}</span>
                  <span class="text-[10px] font-black uppercase tracking-widest text-brand-gold">Индивидуальный заказ</span>
                </div>
                <div v-else class="text-xs italic text-brand-brown/35">Общая консультация</div>
              </td>

              <td class="px-6 py-5 align-top">
                <a :href="`tel:${order.client_phone}`" class="font-semibold text-brand-brown hover:text-brand-gold">
                  {{ order.client_phone }}
                </a>
                <div class="text-xs text-brand-brown/45">{{ order.client_email || 'email не указан' }}</div>
                <div v-if="order.contact_method" class="mt-1 text-[10px] font-black uppercase tracking-widest text-brand-gold">
                  {{ contactMethodLabels[order.contact_method] || order.contact_method }}
                </div>
              </td>

              <td class="px-6 py-5 align-top">
                <span :class="['inline-flex items-center gap-2 rounded-full px-3 py-1 text-[10px] font-black uppercase ring-1', getStatusClass(order.status)]">
                  <component :is="getStatusIcon(order.status)" :size="12" />
                  {{ statusMap[order.status] }}
                </span>
              </td>

              <td class="px-6 py-5 align-top text-sm text-brand-brown/50">
                <div class="whitespace-nowrap">{{ formatDate(order.created_at) }}</div>
                <div class="text-[10px] font-bold uppercase">{{ formatTime(order.created_at) }}</div>
              </td>

              <td class="px-6 py-5 align-top text-right">
                <div v-if="getOrderActions(order).length" class="flex flex-wrap items-center justify-end gap-2">
                  <button
                    v-for="action in getOrderActions(order)"
                    :key="action.key"
                    type="button"
                    :title="action.title"
                    :class="getActionClass(action.tone)"
                    @click="runOrderAction(order, action)"
                  >
                    <component :is="action.icon" :size="17" />
                    <span>{{ action.label }}</span>
                  </button>
                </div>
                <span v-else class="text-[11px] font-bold uppercase tracking-widest text-brand-brown/25">Архив</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section
      v-if="orders.length === 0"
      class="rounded-lg border border-dashed border-brand-brown/15 bg-white p-10 text-center text-brand-brown/35"
    >
      Заявок с таким статусом не найдено
    </section>

    <section class="flex flex-col gap-3 rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between">
      <span class="text-sm text-brand-brown/50">
        Показано {{ orders.length }} из {{ total }} в текущем фильтре
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

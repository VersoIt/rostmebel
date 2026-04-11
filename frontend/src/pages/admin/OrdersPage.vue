<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '@/api/client';
import { LucideCheckCircle, LucideXCircle, LucideClock, LucideDownload } from 'lucide-vue-next';
import type { Order } from '@/types';

const orders = ref<Order[]>([]);
const total = ref(0);

onMounted(async () => {
  try {
    const { data } = await api.get('/admin/orders');
    orders.value = data.items;
    total.value = data.total;
  } catch (err) {
    console.error(err);
  }
});

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'new': return LucideClock;
    case 'done': return LucideCheckCircle;
    case 'rejected': return LucideXCircle;
    default: return LucideClock;
  }
};
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-12">
      <h1 class="font-serif text-4xl text-brand-brown">Заявки клиентов</h1>
      <button class="bg-brand-gray text-brand-brown px-6 py-3 rounded-xl font-medium hover:bg-brand-brown hover:text-white transition-all flex items-center gap-2">
        <LucideDownload :size="20" />
        Экспорт CSV
      </button>
    </div>

    <div class="bg-white rounded-3xl shadow-sm border border-brand-brown/5 overflow-hidden">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-brand-gray/20 text-brand-brown/40 text-xs uppercase tracking-widest">
            <th class="px-8 py-4 font-semibold">ID</th>
            <th class="px-8 py-4 font-semibold">Клиент</th>
            <th class="px-8 py-4 font-semibold">Контакты</th>
            <th class="px-8 py-4 font-semibold">Статус</th>
            <th class="px-8 py-4 font-semibold">Дата</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-brand-brown/5">
          <tr v-for="o in orders" :key="o.id" class="hover:bg-brand-gray/10 transition-colors">
            <td class="px-8 py-4 font-mono text-xs text-brand-brown/40">#{{ o.id }}</td>
            <td class="px-8 py-4">
              <div class="font-medium">{{ o.client_name }}</div>
              <div class="text-xs text-brand-brown/40">{{ o.comment }}</div>
            </td>
            <td class="px-8 py-4">
              <div>{{ o.client_phone }}</div>
              <div class="text-sm text-brand-brown/40">{{ o.client_email }}</div>
            </td>
            <td class="px-8 py-4">
              <div class="flex items-center gap-2">
                <component :is="getStatusIcon(o.status)" :size="16" class="text-brand-gold" />
                <span class="capitalize">{{ o.status }}</span>
              </div>
            </td>
            <td class="px-8 py-4 text-sm text-brand-brown/60">
              {{ new Date(o.created_at).toLocaleDateString() }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

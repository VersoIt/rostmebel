<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useProductStore } from '@/stores/products';
import { 
  LucidePlus, 
  LucideEdit, 
  LucideTrash2, 
  LucideSearch, 
  LucideDownload,
  LucideChevronLeft,
  LucideChevronRight
} from 'lucide-vue-next';
import ProductModal from '@/components/admin/ProductModal.vue';
import type { Product } from '@/types';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { useNotificationStore } from '@/stores/notifications';
import { useConfirmStore } from '@/stores/confirm';
import { downloadFile } from '@/utils/download';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';

const productStore = useProductStore();
const notificationStore = useNotificationStore();
const confirmStore = useConfirmStore();
const searchQuery = ref('');
const isModalOpen = ref(false);
const editingProduct = ref<Product | null>(null);

const currentPage = ref(1);
const limit = 10;

const fetch = () => {
  productStore.fetchProducts({ 
    search: searchQuery.value, 
    limit, 
    offset: (currentPage.value - 1) * limit 
  });
};

onMounted(async () => {
  await productStore.fetchCategories();
  fetch();
});

watch([searchQuery, currentPage], fetch);

const openCreate = () => {
  editingProduct.value = null;
  isModalOpen.value = true;
};

const openEdit = (p: Product) => {
  editingProduct.value = p;
  isModalOpen.value = true;
};

const deleteProduct = async (id: number) => {
  const confirmed = await confirmStore.request({
    title: 'Удалить проект?',
    message: 'Проект будет скрыт из каталога. Это действие нельзя отменить из интерфейса.',
    confirmLabel: 'Удалить',
    tone: 'danger',
  });
  if (!confirmed) return;

  try {
    await api.delete(`/admin/projects/${id}`);
    notificationStore.show('Проект удален', 'info');
    fetch();
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

const handleSaved = () => {
  isModalOpen.value = false;
  fetch();
};

const exportProducts = () => {
  downloadFile('/admin/projects/export', 'projects.xlsx');
};

const handleImgError = (e: Event) => {
  (e.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
};
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-12">
      <h1 class="font-serif text-4xl text-brand-brown">Управление проектами</h1>
      <div class="flex gap-4">
        <button @click="exportProducts" class="bg-brand-gray text-brand-brown px-6 py-3 rounded-xl font-medium hover:bg-brand-brown hover:text-white transition-all flex items-center gap-2">
          <LucideDownload :size="20" />
          Экспорт
        </button>
        <button @click="openCreate" class="bg-brand-brown text-white px-6 py-3 rounded-xl font-medium hover:bg-brand-gold transition-all flex items-center gap-2 shadow-lg">
          <LucidePlus :size="20" />
          Добавить проект
        </button>
      </div>
    </div>

    <div class="bg-white rounded-3xl shadow-sm border border-brand-brown/5 overflow-hidden">
      <div class="p-6 border-b border-brand-brown/5 flex items-center gap-4">
        <div class="relative flex-1 max-w-md">
          <input 
            v-model="searchQuery"
            type="text"
            placeholder="Поиск по названию..."
            class="w-full pl-12 pr-4 py-3 rounded-xl bg-brand-gray/50 border-none outline-none focus:ring-2 ring-brand-gold/20"
          >
          <LucideSearch class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/30" :size="20" />
        </div>
      </div>

      <table class="w-full text-left">
        <thead>
          <tr class="bg-brand-gray/20 text-brand-brown/40 text-xs uppercase tracking-widest">
            <th class="px-8 py-4 font-semibold">Проект</th>
            <th class="px-8 py-4 font-semibold">Категория</th>
            <th class="px-8 py-4 font-semibold">Бюджет</th>
            <th class="px-8 py-4 font-semibold">Статус</th>
            <th class="px-8 py-4 font-semibold text-right">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-brand-brown/5">
          <tr v-for="p in productStore.products" :key="p.id" class="hover:bg-brand-gray/5 transition-colors">
            <td class="px-8 py-4">
              <div class="flex items-center gap-4">
                <img 
                  :src="p.images[0]?.url || PLACEHOLDER_IMAGE" 
                  @error="handleImgError"
                  class="w-12 h-12 rounded-lg object-cover bg-brand-gray"
                >
                <span class="font-medium">{{ p.name }}</span>
              </div>
            </td>
            <td class="px-8 py-4 text-brand-brown/60 text-sm">
              {{ productStore.categories.find(c => c.id === p.project_category_id)?.name || '—' }}
            </td>
            <td class="px-8 py-4 font-semibold">
              {{ p.price.toLocaleString() }} ₽
            </td>
            <td class="px-8 py-4">
              <span :class="[
                'px-3 py-1 rounded-full text-[10px] font-bold uppercase',
                p.status === 'published' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
              ]">
                {{ p.status }}
              </span>
            </td>
            <td class="px-8 py-4 text-right">
              <div class="flex items-center justify-end gap-2">
                <button @click="openEdit(p)" class="p-2 text-brand-brown/40 hover:text-brand-gold transition-colors">
                  <LucideEdit :size="18" />
                </button>
                <button @click="deleteProduct(p.id)" class="p-2 text-brand-brown/40 hover:text-red-500 transition-colors">
                  <LucideTrash2 :size="18" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div class="p-6 bg-brand-gray/10 flex items-center justify-between">
        <span class="text-sm text-brand-brown/40">Показано {{ productStore.products.length }} из {{ productStore.total }} проектов</span>
        <div class="flex gap-2">
          <button 
            @click="currentPage--" 
            :disabled="currentPage === 1"
            class="p-2 rounded-lg bg-white border border-brand-brown/5 disabled:opacity-30 disabled:cursor-default"
          >
            <LucideChevronLeft :size="20" />
          </button>
          <div class="px-4 py-2 bg-brand-brown text-white rounded-lg font-bold text-sm">{{ currentPage }}</div>
          <button 
            @click="currentPage++" 
            :disabled="currentPage * limit >= productStore.total"
            class="p-2 rounded-lg bg-white border border-brand-brown/5 disabled:opacity-30 disabled:cursor-default"
          >
            <LucideChevronRight :size="20" />
          </button>
        </div>
      </div>
    </div>

    <ProductModal 
      v-if="isModalOpen" 
      :product="editingProduct" 
      :categories="productStore.categories"
      @close="isModalOpen = false" 
      @saved="handleSaved" 
    />
  </div>
</template>

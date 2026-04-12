<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import {
  LucideChevronLeft,
  LucideChevronRight,
  LucideDownload,
  LucideEdit,
  LucidePlus,
  LucideSearch,
  LucideTrash2,
} from 'lucide-vue-next';
import ProductModal from '@/components/admin/ProductModal.vue';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { useConfirmStore } from '@/stores/confirm';
import { useNotificationStore } from '@/stores/notifications';
import { useProductStore } from '@/stores/products';
import type { Product, ProjectStatus } from '@/types';
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

const statusMap: Record<ProjectStatus, string> = {
  published: 'Опубликован',
  draft: 'Черновик',
  archived: 'Архив',
};

const totalPages = computed(() => Math.max(1, Math.ceil(productStore.total / limit)));

const fetchProjects = () => {
  void productStore.fetchProducts({
    search: searchQuery.value.trim(),
    limit,
    offset: (currentPage.value - 1) * limit,
  });
};

onMounted(async () => {
  await productStore.fetchCategories();
  fetchProjects();
});

watch(searchQuery, () => {
  if (currentPage.value !== 1) {
    currentPage.value = 1;
    return;
  }
  fetchProjects();
});

watch(currentPage, fetchProjects);

const openCreate = () => {
  editingProduct.value = null;
  isModalOpen.value = true;
};

const openEdit = (product: Product) => {
  editingProduct.value = product;
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
    fetchProjects();
  } catch (err) {
    notificationStore.show(getApiErrorMessage(err), 'error');
  }
};

const handleSaved = () => {
  isModalOpen.value = false;
  fetchProjects();
};

const exportProducts = () => {
  downloadFile('/admin/projects/export', 'projects.xlsx');
};

const handleImgError = (event: Event) => {
  (event.target as HTMLImageElement).src = PLACEHOLDER_IMAGE;
};

const categoryName = (product: Product) => {
  return productStore.categories.find((category) => category.id === product.project_category_id)?.name || 'Без категории';
};

const formatPrice = (price: number) => {
  return `${price.toLocaleString('ru-RU')} ₽`;
};

const getStatusClass = (status: ProjectStatus) => {
  if (status === 'published') return 'bg-green-50 text-green-700 ring-green-100';
  if (status === 'archived') return 'bg-gray-100 text-gray-700 ring-gray-200';
  return 'bg-yellow-50 text-yellow-700 ring-yellow-100';
};
</script>

<template>
  <div class="space-y-6">
    <section class="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
      <div>
        <p class="text-xs font-black uppercase text-brand-gold">Каталог</p>
        <h1 class="mt-1 font-serif text-3xl font-bold leading-tight text-brand-brown sm:text-4xl">Управление проектами</h1>
        <p class="mt-2 max-w-2xl text-sm text-brand-brown/55">
          Проекты, которые клиент видит в портфолио, каталоге и AI-поиске.
        </p>
      </div>

      <div class="grid gap-2 sm:grid-cols-2 xl:flex xl:gap-3">
        <button
          type="button"
          @click="exportProducts"
          class="inline-flex h-12 items-center justify-center gap-2 rounded-lg bg-white px-5 text-sm font-bold text-brand-brown ring-1 ring-brand-brown/10 transition-all hover:bg-brand-brown hover:text-white"
        >
          <LucideDownload :size="19" />
          Экспорт
        </button>
        <button
          type="button"
          @click="openCreate"
          class="inline-flex h-12 items-center justify-center gap-2 rounded-lg bg-brand-brown px-5 text-sm font-bold text-white shadow-lg shadow-brand-brown/10 transition-all hover:bg-brand-gold"
        >
          <LucidePlus :size="19" />
          Добавить проект
        </button>
      </div>
    </section>

    <section class="rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm">
      <div class="relative max-w-xl">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Поиск по названию"
          class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/40 py-3 pl-11 pr-4 outline-none transition-all focus:border-brand-gold focus:ring-2 focus:ring-brand-gold/20"
        >
        <LucideSearch class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/35" :size="19" />
      </div>
    </section>

    <section v-if="productStore.loading" class="rounded-lg border border-brand-brown/10 bg-white p-10 text-center text-sm font-semibold text-brand-brown/50">
      Загружаем проекты...
    </section>

    <section v-if="!productStore.loading && productStore.products.length" class="space-y-3 xl:hidden" aria-label="Проекты">
      <article
        v-for="product in productStore.products"
        :key="product.id"
        class="rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm"
      >
        <div class="flex gap-3">
          <img
            :src="product.images[0]?.url || PLACEHOLDER_IMAGE"
            class="h-20 w-20 shrink-0 rounded-lg bg-brand-gray object-cover"
            alt=""
            @error="handleImgError"
          >

          <div class="min-w-0 flex-1">
            <div class="mb-1 flex flex-wrap items-center gap-2">
              <span class="font-mono text-xs font-bold text-brand-brown/35">#{{ product.id }}</span>
              <span :class="['rounded-full px-2.5 py-1 text-[11px] font-black uppercase ring-1', getStatusClass(product.status)]">
                {{ statusMap[product.status] || product.status }}
              </span>
            </div>
            <h2 class="line-clamp-2 font-bold text-brand-brown">{{ product.name }}</h2>
            <p class="mt-1 text-sm text-brand-brown/55">{{ categoryName(product) }}</p>
          </div>
        </div>

        <div class="mt-4 grid grid-cols-2 gap-3 text-sm">
          <div>
            <div class="text-[11px] font-black uppercase text-brand-brown/35">Бюджет</div>
            <div class="mt-1 font-bold text-brand-brown">{{ formatPrice(product.price) }}</div>
          </div>
          <div>
            <div class="text-[11px] font-black uppercase text-brand-brown/35">Заявки</div>
            <div class="mt-1 font-bold text-brand-brown">{{ product.orders_count }}</div>
          </div>
        </div>

        <div class="mt-4 grid grid-cols-2 gap-2 border-t border-brand-brown/10 pt-4">
          <button
            type="button"
            @click="openEdit(product)"
            class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg bg-brand-gray px-3 py-2 text-sm font-bold text-brand-brown transition-all hover:bg-brand-brown hover:text-white"
          >
            <LucideEdit :size="17" />
            Изменить
          </button>
          <button
            type="button"
            @click="deleteProduct(product.id)"
            class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg bg-red-50 px-3 py-2 text-sm font-bold text-red-700 transition-all hover:bg-red-600 hover:text-white"
          >
            <LucideTrash2 :size="17" />
            Удалить
          </button>
        </div>
      </article>
    </section>

    <section v-if="!productStore.loading && productStore.products.length" class="hidden overflow-hidden rounded-lg border border-brand-brown/10 bg-white shadow-sm xl:block" aria-label="Таблица проектов">
      <table class="w-full text-left">
        <thead>
          <tr class="bg-brand-gray/50 text-xs uppercase text-brand-brown/45">
            <th class="px-6 py-4 font-semibold">Проект</th>
            <th class="px-6 py-4 font-semibold">Категория</th>
            <th class="px-6 py-4 font-semibold">Бюджет</th>
            <th class="px-6 py-4 font-semibold">Статус</th>
            <th class="px-6 py-4 text-right font-semibold">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-brand-brown/10">
          <tr v-for="product in productStore.products" :key="product.id" class="transition-colors hover:bg-brand-gray/30">
            <td class="px-6 py-4">
              <div class="flex items-center gap-4">
                <img
                  :src="product.images[0]?.url || PLACEHOLDER_IMAGE"
                  class="h-12 w-12 rounded-lg bg-brand-gray object-cover"
                  alt=""
                  @error="handleImgError"
                >
                <div>
                  <div class="font-bold text-brand-brown">{{ product.name }}</div>
                  <div class="font-mono text-xs text-brand-brown/35">#{{ product.id }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-brand-brown/60">{{ categoryName(product) }}</td>
            <td class="px-6 py-4 font-semibold text-brand-brown">{{ formatPrice(product.price) }}</td>
            <td class="px-6 py-4">
              <span :class="['rounded-full px-3 py-1 text-[10px] font-black uppercase ring-1', getStatusClass(product.status)]">
                {{ statusMap[product.status] || product.status }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <div class="flex items-center justify-end gap-2">
                <button
                  type="button"
                  @click="openEdit(product)"
                  class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-brand-gray text-brand-brown transition-all hover:bg-brand-brown hover:text-white"
                  title="Изменить проект"
                >
                  <LucideEdit :size="18" />
                </button>
                <button
                  type="button"
                  @click="deleteProduct(product.id)"
                  class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-red-50 text-red-700 transition-all hover:bg-red-600 hover:text-white"
                  title="Удалить проект"
                >
                  <LucideTrash2 :size="18" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section
      v-if="!productStore.loading && productStore.products.length === 0"
      class="rounded-lg border border-dashed border-brand-brown/15 bg-white p-10 text-center text-brand-brown/45"
    >
      Проекты не найдены
    </section>

    <section
      v-if="productStore.total > 0"
      class="flex flex-col gap-3 rounded-lg border border-brand-brown/10 bg-white p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between"
    >
      <span class="text-sm text-brand-brown/50">
        Показано {{ productStore.products.length }} из {{ productStore.total }} проектов
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

    <ProductModal
      v-if="isModalOpen"
      :product="editingProduct"
      :categories="productStore.categories"
      @close="isModalOpen = false"
      @saved="handleSaved"
    />
  </div>
</template>

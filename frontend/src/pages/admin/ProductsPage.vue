<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useProductStore } from '@/stores/products';
import { LucidePlus, LucideEdit, LucideTrash2, LucideSearch } from 'lucide-vue-next';
import ProductModal from '@/components/admin/ProductModal.vue';
import type { Product } from '@/types';
import api from '@/api/client';

const productStore = useProductStore();
const searchQuery = ref('');
const isModalOpen = ref(false);
const editingProduct = ref<Product | null>(null);

onMounted(async () => {
  await productStore.fetchCategories();
  await productStore.fetchProducts();
});

watch(searchQuery, (val) => {
  productStore.fetchProducts({ search: val });
});

const openCreate = () => {
  editingProduct.value = null;
  isModalOpen.value = true;
};

const openEdit = (p: Product) => {
  editingProduct.value = p;
  isModalOpen.value = true;
};

const deleteProduct = async (id: number) => {
  if (confirm('Вы уверены, что хотите удалить этот товар?')) {
    try {
      await api.delete(`/admin/products/${id}`);
      await productStore.fetchProducts();
    } catch (err) {
      alert('Ошибка при удалении');
    }
  }
};

const handleSaved = () => {
  isModalOpen.value = false;
  productStore.fetchProducts();
};
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-12">
      <h1 class="font-serif text-4xl text-brand-brown">Управление товарами</h1>
      <button @click="openCreate" class="bg-brand-brown text-white px-6 py-3 rounded-xl font-medium hover:bg-brand-gold transition-all flex items-center gap-2 shadow-lg">
        <LucidePlus :size="20" />
        Добавить товар
      </button>
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
            <th class="px-8 py-4 font-semibold">Товар</th>
            <th class="px-8 py-4 font-semibold">Категория</th>
            <th class="px-8 py-4 font-semibold">Цена</th>
            <th class="px-8 py-4 font-semibold">Статус</th>
            <th class="px-8 py-4 font-semibold text-right">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-brand-brown/5">
          <tr v-for="p in productStore.products" :key="p.id" class="hover:bg-brand-gray/10 transition-colors">
            <td class="px-8 py-4">
              <div class="flex items-center gap-4">
                <img :src="p.images[0]?.url" class="w-12 h-12 rounded-lg object-cover bg-brand-gray">
                <span class="font-medium">{{ p.name }}</span>
              </div>
            </td>
            <td class="px-8 py-4 text-brand-brown/60 text-sm">
              {{ productStore.categories.find(c => c.id === p.category_id)?.name || '—' }}
            </td>
            <td class="px-8 py-4 font-semibold">
              {{ p.price.toLocaleString() }} ₽
            </td>
            <td class="px-8 py-4">
              <span :class="[
                'px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-tighter',
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

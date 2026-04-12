<script setup lang="ts">
import { useFavorites } from '@/composables/useFavorites';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideHeartOff } from 'lucide-vue-next';

const { favorites } = useFavorites();
</script>

<template>
  <div class="min-h-screen bg-brand-cream pt-28">
    <div class="ui-container ui-section">
      <header class="mb-10">
        <p class="ui-eyebrow mb-3">Избранное</p>
        <h1 class="ui-title-xl">Сохраненные проекты</h1>
        <p v-if="favorites.length > 0" class="ui-copy-lg mt-4">
          Проекты, к которым можно вернуться перед расчетом.
        </p>
      </header>

      <div v-if="favorites.length > 0" class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <ProductCard v-for="product in favorites" :key="product.id" :product="product" />
      </div>

      <div v-else class="ui-empty py-16">
        <LucideHeartOff :size="56" class="mx-auto mb-5 text-brand-brown/12" />
        <h2 class="ui-title-md mb-2">Список пока пуст</h2>
        <p class="mx-auto mb-6 max-w-md text-brand-brown/55">
          Сохраняйте проекты из каталога, чтобы быстро вернуться к ним перед консультацией.
        </p>
        <router-link to="/catalog" class="ui-button ui-button-primary">
          Смотреть проекты
        </router-link>
      </div>
    </div>
  </div>
</template>

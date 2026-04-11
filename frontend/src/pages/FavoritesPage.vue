<script setup lang="ts">
import { useFavorites } from '@/composables/useFavorites';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideHeartOff } from 'lucide-vue-next';

const { favorites } = useFavorites();
</script>

<template>
  <div class="bg-brand-cream min-h-screen pt-32 pb-24 px-4">
    <div class="max-w-7xl mx-auto">
      <header class="mb-12">
        <h1 class="font-serif text-5xl text-brand-brown mb-4">Избранное</h1>
        <p v-if="favorites.length > 0" class="text-brand-brown/60">Ваши понравившиеся товары</p>
      </header>

      <div v-if="favorites.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        <ProductCard v-for="p in favorites" :key="p.id" :product="p" />
      </div>

      <div v-else class="text-center py-32 bg-white rounded-3xl border border-dashed border-brand-brown/10">
        <LucideHeartOff :size="64" class="mx-auto text-brand-brown/10 mb-6" />
        <h3 class="text-2xl font-serif text-brand-brown mb-2">Здесь пока пусто</h3>
        <p class="text-brand-brown/40 mb-8">Добавляйте товары в избранное, чтобы не потерять их</p>
        <router-link to="/catalog" class="bg-brand-brown text-white px-8 py-4 rounded-xl font-medium hover:bg-brand-gold transition-colors">
          Перейти в каталог
        </router-link>
      </div>
    </div>
  </div>
</template>

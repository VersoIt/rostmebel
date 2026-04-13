<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import BrandMark from '@/components/common/BrandMark.vue';
import {
  LucideClipboardList,
  LucideExternalLink,
  LucideLayoutDashboard,
  LucideLogOut,
  LucideMessageSquare,
  LucidePackage,
} from 'lucide-vue-next';

const authStore = useAuthStore();
const route = useRoute();

const menuItems = [
  { name: 'Дашборд', shortName: 'Обзор', icon: LucideLayoutDashboard, to: '/admin' },
  { name: 'Проекты', shortName: 'Проекты', icon: LucidePackage, to: '/admin/projects' },
  { name: 'Заявки', shortName: 'Заявки', icon: LucideClipboardList, to: '/admin/orders' },
  { name: 'Отзывы', shortName: 'Отзывы', icon: LucideMessageSquare, to: '/admin/reviews' },
];

const isActive = (path: string) => {
  if (path === '/admin') {
    return route.path === '/admin';
  }
  return route.path === path || route.path.startsWith(`${path}/`);
};

const currentPage = computed(() => {
  return menuItems.find((item) => isActive(item.to))?.name || 'Админка';
});

const logout = () => {
  authStore.logout();
};
</script>

<template>
  <div class="min-h-screen bg-brand-gray/40 text-brand-brown">
    <aside class="fixed inset-y-0 left-0 z-50 hidden w-72 flex-col bg-brand-brown text-brand-white shadow-2xl lg:flex">
      <div class="flex items-center gap-2 border-b border-white/10 p-6">
        <BrandMark class="shrink-0 text-brand-gold" size="2.65rem" />
        <div>
          <div class="font-serif text-xl font-bold uppercase leading-none">РОСТ</div>
          <div class="mt-1 text-xs font-bold uppercase tracking-widest text-brand-gold">Админ-панель</div>
        </div>
      </div>

      <nav class="flex-1 space-y-2 p-4">
        <router-link
          v-for="item in menuItems"
          :key="item.name"
          :to="item.to"
          class="admin-sidebar-link flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold transition-all"
          :class="[
            isActive(item.to)
              ? 'bg-brand-gold text-brand-brown shadow-lg shadow-brand-gold/20'
              : 'text-white/65 hover:bg-white/10 hover:text-white',
          ]"
        >
          <component :is="item.icon" :size="20" />
          {{ item.name }}
        </router-link>
      </nav>

      <div class="space-y-2 border-t border-white/10 p-4">
        <a
          href="/"
          target="_blank"
          class="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-semibold text-white/55 transition-all hover:bg-white/10 hover:text-brand-gold"
        >
          <LucideExternalLink :size="20" />
          На сайт
        </a>
        <button
          type="button"
          @click="logout"
          class="flex w-full items-center gap-3 rounded-lg px-4 py-3 text-left text-sm font-semibold text-white/55 transition-all hover:bg-red-500/10 hover:text-red-300"
        >
          <LucideLogOut :size="20" />
          Выйти
        </button>
      </div>
    </aside>

    <header class="sticky top-0 z-40 border-b border-brand-brown/10 bg-white/95 backdrop-blur lg:hidden">
      <div class="flex items-center justify-between gap-3 px-4 py-3">
        <div class="flex min-w-0 items-center gap-2">
          <BrandMark class="shrink-0 text-brand-gold" size="2.65rem" />
          <div class="min-w-0">
            <div class="text-[11px] font-black uppercase tracking-widest text-brand-gold">РОСТ Админ</div>
            <div class="truncate font-serif text-xl font-bold text-brand-brown">{{ currentPage }}</div>
          </div>
        </div>

        <div class="flex shrink-0 items-center gap-2">
          <a
            href="/"
            target="_blank"
            class="flex h-10 w-10 items-center justify-center rounded-lg border border-brand-brown/10 bg-brand-gray/60 text-brand-brown transition-colors hover:bg-brand-brown hover:text-white"
            aria-label="Открыть сайт"
          >
            <LucideExternalLink :size="18" />
          </a>
          <button
            type="button"
            @click="logout"
            class="flex h-10 w-10 items-center justify-center rounded-lg border border-red-100 bg-red-50 text-red-600 transition-colors hover:bg-red-600 hover:text-white"
            aria-label="Выйти"
          >
            <LucideLogOut :size="18" />
          </button>
        </div>
      </div>
    </header>

    <main class="min-h-screen px-4 pb-[calc(8rem+env(safe-area-inset-bottom))] pt-6 sm:px-6 lg:ml-72 lg:p-10 xl:p-12">
      <router-view />
    </main>

    <nav class="fixed inset-x-0 bottom-0 z-50 border-t border-brand-brown/10 bg-white/95 px-2 pt-2 shadow-[0_-12px_30px_rgba(23,33,29,0.08)] backdrop-blur lg:hidden">
      <div class="grid grid-cols-4 gap-1 pb-[calc(env(safe-area-inset-bottom)+0.5rem)]">
        <router-link
          v-for="item in menuItems"
          :key="item.name"
          :to="item.to"
          class="flex h-14 flex-col items-center justify-center gap-1 rounded-lg text-[11px] font-bold transition-all"
          :class="[
            isActive(item.to)
              ? 'bg-brand-brown text-white shadow-lg shadow-brand-brown/15'
              : 'text-brand-brown/55 hover:bg-brand-gray hover:text-brand-brown',
          ]"
        >
          <component :is="item.icon" :size="19" />
          <span class="leading-none">{{ item.shortName }}</span>
        </router-link>
      </div>
    </nav>
  </div>
</template>

<style scoped>
.admin-sidebar-link {
  text-decoration: none;
  user-select: none;
}
</style>

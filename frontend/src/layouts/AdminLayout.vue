<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { 
  LucideLayoutDashboard, 
  LucidePackage, 
  LucideClipboardList, 
  LucideLogOut, 
  LucideExternalLink
} from 'lucide-vue-next';

const authStore = useAuthStore();

const menuItems = [
  { name: 'Дашборд', icon: LucideLayoutDashboard, to: '/admin' },
  { name: 'Проекты', icon: LucidePackage, to: '/admin/products' },
  { name: 'Заявки', icon: LucideClipboardList, to: '/admin/orders' },
];
</script>

<template>
  <div class="flex min-h-screen bg-brand-gray/30">
    <!-- Sidebar -->
    <aside class="w-72 bg-brand-brown text-brand-white flex flex-col fixed inset-y-0 shadow-2xl z-50">
      <div class="p-8 flex items-center gap-3 border-b border-white/5">
        <div class="w-10 h-10 bg-brand-gold rounded-xl flex items-center justify-center text-white font-serif text-2xl">Р</div>
        <span class="font-serif text-xl tracking-tighter uppercase font-bold">РОСТ <span class="text-brand-gold">Админ</span></span>
      </div>

      <nav class="flex-1 p-6 space-y-2 mt-4">
        <router-link 
          v-for="item in menuItems" 
          :key="item.name"
          :to="item.to"
          v-slot="{ isExactActive, navigate }"
          custom
        >
          <div 
            @click="navigate"
            class="admin-sidebar-link flex items-center gap-4 px-6 py-4 rounded-xl transition-all group cursor-pointer"
            :class="[
              isExactActive 
                ? 'bg-brand-gold !text-brand-brown font-semibold shadow-lg shadow-brand-gold/20' 
                : 'hover:bg-white/5 text-white/60 hover:text-white'
            ]"
          >
            <component :is="item.icon" :size="20" />
            {{ item.name }}
          </div>
        </router-link>
      </nav>

      <div class="p-6 space-y-2 border-t border-white/5">
        <a href="/" target="_blank" class="flex items-center gap-4 px-6 py-4 rounded-xl text-white/40 hover:text-brand-gold hover:bg-white/5 transition-all group">
          <LucideExternalLink :size="20" />
          На сайт
        </a>
        <button 
          @click="authStore.logout()"
          class="w-full flex items-center gap-4 px-6 py-4 rounded-xl text-white/40 hover:text-red-400 hover:bg-red-500/5 transition-all group"
        >
          <LucideLogOut :size="20" class="group-hover:-translate-x-1 transition-transform" />
          Выйти
        </button>
      </div>
    </aside>

    <!-- Content -->
    <main class="flex-1 ml-72 p-12">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.admin-sidebar-link {
  text-decoration: none;
  user-select: none;
}
</style>

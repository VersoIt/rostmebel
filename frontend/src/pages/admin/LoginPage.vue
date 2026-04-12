<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { LucideArrowRight, LucideLock, LucideUser } from 'lucide-vue-next';
import BrandMark from '@/components/common/BrandMark.vue';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');

const handleLogin = async () => {
  const success = await authStore.login({
    username: username.value,
    password: password.value,
  });

  if (success) {
    router.push({ name: 'admin' });
  }
};
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-brand-cream px-4 py-10">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <BrandMark class="mx-auto mb-5 text-brand-gold" size="4rem" />
        <h1 class="font-serif text-3xl font-bold text-brand-brown">Панель управления</h1>
        <p class="mt-2 text-sm text-brand-brown/55">Введите данные администратора, чтобы продолжить</p>
      </div>

      <section class="rounded-lg border border-brand-brown/10 bg-white p-5 shadow-2xl sm:p-8">
        <form class="space-y-5" @submit.prevent="handleLogin">
          <div>
            <label class="mb-2 block text-xs font-black uppercase tracking-widest text-brand-brown/55">Логин</label>
            <div class="relative">
              <input
                v-model="username"
                type="text"
                required
                autocomplete="username"
                class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/40 py-4 pl-12 pr-4 outline-none transition-all focus:border-brand-gold focus:ring-2 focus:ring-brand-gold/20"
                placeholder="admin"
              >
              <LucideUser class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/35" :size="20" />
            </div>
          </div>

          <div>
            <label class="mb-2 block text-xs font-black uppercase tracking-widest text-brand-brown/55">Пароль</label>
            <div class="relative">
              <input
                v-model="password"
                type="password"
                required
                autocomplete="current-password"
                class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/40 py-4 pl-12 pr-4 outline-none transition-all focus:border-brand-gold focus:ring-2 focus:ring-brand-gold/20"
                placeholder="••••••••"
              >
              <LucideLock class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/35" :size="20" />
            </div>
          </div>

          <div v-if="authStore.error" class="rounded-lg bg-red-50 p-4 text-sm font-medium text-red-700">
            {{ authStore.error }}
          </div>

          <button
            type="submit"
            :disabled="authStore.loading"
            class="flex h-12 w-full items-center justify-center gap-2 rounded-lg bg-brand-brown px-5 font-bold text-white shadow-lg shadow-brand-brown/10 transition-all hover:bg-brand-gold disabled:cursor-not-allowed disabled:opacity-50"
          >
            {{ authStore.loading ? 'Входим...' : 'Войти' }}
            <LucideArrowRight v-if="!authStore.loading" :size="20" />
          </button>
        </form>
      </section>

      <div class="mt-7 text-center">
        <router-link to="/" class="text-sm font-semibold text-brand-brown/50 transition-colors hover:text-brand-brown">
          Вернуться на сайт
        </router-link>
      </div>
    </div>
  </div>
</template>

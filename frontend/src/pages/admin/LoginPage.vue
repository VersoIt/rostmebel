<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { LucideLock, LucideUser, LucideArrowRight } from 'lucide-vue-next';

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
  <div class="min-h-screen bg-brand-cream flex items-center justify-center px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <div class="w-20 h-20 bg-brand-gold rounded-3xl flex items-center justify-center text-white font-serif text-4xl mx-auto mb-6 shadow-xl">Р</div>
        <h1 class="font-serif text-3xl text-brand-brown">Панель управления</h1>
        <p class="text-brand-brown/40 mt-2">Введите ваши данные для входа</p>
      </div>

      <div class="bg-white p-8 rounded-3xl shadow-2xl border border-brand-brown/5">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-brand-brown/60 mb-2 uppercase tracking-widest">Логин</label>
            <div class="relative">
              <input 
                v-model="username"
                type="text"
                required
                class="w-full pl-12 pr-4 py-4 rounded-xl border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30 transition-all"
                placeholder="admin"
              >
              <LucideUser class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/30" :size="20" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-brand-brown/60 mb-2 uppercase tracking-widest">Пароль</label>
            <div class="relative">
              <input 
                v-model="password"
                type="password"
                required
                class="w-full pl-12 pr-4 py-4 rounded-xl border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30 transition-all"
                placeholder="••••••••"
              >
              <LucideLock class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/30" :size="20" />
            </div>
          </div>

          <div v-if="authStore.error" class="text-red-500 text-sm bg-red-50 p-4 rounded-xl">
            {{ authStore.error }}
          </div>

          <button 
            type="submit"
            :disabled="authStore.loading"
            class="w-full bg-brand-brown text-white py-4 rounded-xl font-medium hover:bg-brand-gold disabled:opacity-50 transition-all flex items-center justify-center gap-2 shadow-lg"
          >
            {{ authStore.loading ? 'Вход...' : 'Войти' }}
            <LucideArrowRight v-if="!authStore.loading" :size="20" />
          </button>
        </form>
      </div>

      <div class="text-center mt-8">
        <router-link to="/" class="text-brand-brown/40 hover:text-brand-brown transition-colors">
          &larr; Вернуться на сайт
        </router-link>
      </div>
    </div>
  </div>
</template>

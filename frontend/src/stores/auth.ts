import { defineStore } from 'pinia';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import type { TokenPair } from '@/types';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isAuthenticated: !!localStorage.getItem('access_token'),
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async login(credentials: any) {
      this.loading = true;
      this.error = null;
      try {
        const { data } = await api.post<TokenPair>('/admin/auth/login', credentials);
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        this.isAuthenticated = true;
        return true;
      } catch (err: any) {
        this.error = getApiErrorMessage(err);
        return false;
      } finally {
        this.loading = false;
      }
    },
    async logout() {
      try {
        await api.post('/admin/auth/logout');
      } catch (err) {
        console.error(err);
      } finally {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        this.isAuthenticated = false;
        window.location.href = '/admin/login';
      }
    }
  },
});

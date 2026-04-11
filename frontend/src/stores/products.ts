import { defineStore } from 'pinia';
import api from '@/api/client';
import type { Product, Category } from '@/types';

export const useProductStore = defineStore('products', {
  state: () => ({
    products: [] as Product[],
    total: 0,
    categories: [] as Category[],
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async fetchProducts(params: any = {}) {
      this.loading = true;
      try {
        const { data } = await api.get('/products', { params });
        this.products = data.items;
        this.total = data.total;
      } catch (err: any) {
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    },
    async fetchCategories() {
      try {
        const { data } = await api.get('/categories');
        this.categories = data;
      } catch (err: any) {
        console.error(err);
      }
    },
    async aiSearch(query: string) {
      this.loading = true;
      try {
        const { data } = await api.post('/ai/search', { query });
        return data as Product[];
      } catch (err: any) {
        this.error = err.message;
        return [];
      } finally {
        this.loading = false;
      }
    },
    async fetchProduct(idOrSlug: string | number) {
      try {
        const { data } = await api.get(`/products/${idOrSlug}`);
        return data as Product;
      } catch (err: any) {
        console.error(err);
        return null;
      }
    }
  },
});

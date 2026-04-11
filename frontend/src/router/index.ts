import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/PublicLayout.vue'),
      children: [
        { path: '', name: 'home', component: () => import('@/pages/HomePage.vue') },
        { path: 'catalog', name: 'catalog', component: () => import('@/pages/CatalogPage.vue') },
        { path: 'product/:id', name: 'product', component: () => import('@/pages/ProductPage.vue') },
        { path: 'contact', name: 'contact', component: () => import('@/pages/ContactPage.vue') },
        { path: 'favorites', name: 'favorites', component: () => import('@/pages/FavoritesPage.vue') },
      ]
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('@/pages/admin/LoginPage.vue'),
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'admin', component: () => import('@/pages/admin/DashboardPage.vue') },
        { path: 'products', name: 'admin-products', component: () => import('@/pages/admin/ProductsPage.vue') },
        { path: 'orders', name: 'admin-orders', component: () => import('@/pages/admin/OrdersPage.vue') },
      ]
    },
  ],
});

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();
  
  if (to.matched.some(record => record.meta.requiresAuth) && !authStore.isAuthenticated) {
    next({ name: 'admin-login' });
  } else if (to.name === 'admin-login' && authStore.isAuthenticated) {
    next({ name: 'admin' });
  } else {
    next();
  }
});

export default router;

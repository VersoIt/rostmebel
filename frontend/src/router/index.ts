import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior(to, _from, savedPosition) {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth',
        top: 100,
      };
    }
    if (savedPosition) {
      return savedPosition;
    }
    return { top: 0 };
  },
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/PublicLayout.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/pages/HomePage.vue'),
          meta: {
            title: 'РОСТ Мебель — кухни и корпусная мебель по размеру',
            description: 'Проектируем, производим и устанавливаем кухни, шкафы, гардеробные и корпусную мебель по Крыму.',
          },
        },
        {
          path: 'catalog',
          name: 'catalog',
          component: () => import('@/pages/CatalogPage.vue'),
          meta: {
            title: 'Проекты — РОСТ Мебель',
            description: 'Портфолио реализованных кухонь, шкафов и систем хранения с ценами, материалами и деталями.',
          },
        },
        {
          path: 'product/:id',
          name: 'product',
          component: () => import('@/pages/ProductPage.vue'),
          meta: {
            title: 'Детали проекта — РОСТ Мебель',
            description: 'Описание реализованного проекта мебели на заказ: бюджет, материалы, детали и заявка на расчет.',
          },
        },
        {
          path: 'contact',
          name: 'contact',
          component: () => import('@/pages/ContactPage.vue'),
          meta: {
            title: 'Контакты — РОСТ Мебель',
            description: 'Свяжитесь с РОСТ Мебель для консультации, предварительного расчета или замера. Работаем по Крыму.',
          },
        },
        {
          path: 'favorites',
          name: 'favorites',
          component: () => import('@/pages/FavoritesPage.vue'),
          meta: {
            title: 'Избранное — РОСТ Мебель',
            description: 'Сохраненные проекты мебели, к которым можно вернуться перед консультацией.',
          },
        },
      ],
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
        { path: 'projects', name: 'admin-projects', component: () => import('@/pages/admin/ProjectsPage.vue') },
        { path: 'orders', name: 'admin-orders', component: () => import('@/pages/admin/OrdersPage.vue') },
        { path: 'reviews', name: 'admin-reviews', component: () => import('@/pages/admin/ReviewsPage.vue') },
      ],
    },
  ],
});

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();

  if (to.matched.some((record) => record.meta.requiresAuth) && !authStore.isAuthenticated) {
    next({ name: 'admin-login' });
    return;
  }

  if (to.name === 'admin-login' && authStore.isAuthenticated) {
    next({ name: 'admin' });
    return;
  }

  next();
});

router.afterEach((to) => {
  const defaultTitle = 'РОСТ Мебель — кухни и корпусная мебель по размеру';
  const title = (to.meta.title as string) || defaultTitle;
  const description = (to.meta.description as string) || '';

  document.title = title;

  const metaDesc = document.querySelector('meta[name="description"]');
  if (metaDesc) {
    metaDesc.setAttribute('content', description);
  } else {
    const meta = document.createElement('meta');
    meta.name = 'description';
    meta.content = description;
    document.head.appendChild(meta);
  }
});

export default router;

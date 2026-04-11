import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth',
        top: 100, // Offset for sticky header
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
            title: 'РОСТ Мебель — Премиальная мебель на заказ по индивидуальным размерам', 
            description: 'Изготовление кухонь, шкафов и авторской мебели на заказ. 15 лет опыта, гарантия качества и уникальный дизайн вашего интерьера.' 
          }
        },
        { 
          path: 'catalog', 
          name: 'catalog', 
          component: () => import('@/pages/CatalogPage.vue'),
          meta: { 
            title: 'Портфолио проектов — РОСТ Мебель', 
            description: 'Смотрите примеры реализованных нами проектов: кухни, гостиные, спальни и другая мебель на заказ с ценами и деталями.' 
          }
        },
        { 
          path: 'product/:id', 
          name: 'product', 
          component: () => import('@/pages/ProductPage.vue'),
          meta: { 
            title: 'Детали проекта — РОСТ Мебель', 
            description: 'Подробное описание реализованного проекта мебели на заказ.' 
          }
        },
        { 
          path: 'contact', 
          name: 'contact', 
          component: () => import('@/pages/ContactPage.vue'),
          meta: { 
            title: 'Контакты — РОСТ Мебель', 
            description: 'Свяжитесь с нами для консультации или замера. Мы находимся в Севастополе, работаем по всему Крыму.' 
          }
        },
        { 
          path: 'favorites', 
          name: 'favorites', 
          component: () => import('@/pages/FavoritesPage.vue'),
          meta: { 
            title: 'Избранные идеи — РОСТ Мебель', 
            description: 'Проекты и идеи мебели, которые вы сохранили для своего будущего интерьера.' 
          }
        },
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
        { path: 'projects', name: 'admin-projects', component: () => import('@/pages/admin/ProjectsPage.vue') },
        { path: 'orders', name: 'admin-orders', component: () => import('@/pages/admin/OrdersPage.vue') },
        { path: 'reviews', name: 'admin-reviews', component: () => import('@/pages/admin/ReviewsPage.vue') },
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

router.afterEach((to) => {
  const defaultTitle = 'РОСТ Мебель — Премиальная мебель на заказ';
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

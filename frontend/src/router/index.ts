import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { setPageSeo } from '@/utils/seo';

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
            title: 'РОСТ Мебель — кухни и корпусная мебель по размеру в Крыму',
            description: 'Кухни, шкафы, гардеробные и системы хранения по размеру: проект, производство, доставка и монтаж по Крыму.',
            canonicalPath: '/',
          },
        },
        {
          path: 'catalog',
          name: 'catalog',
          component: () => import('@/pages/CatalogPage.vue'),
          meta: {
            title: 'Проекты кухонь и корпусной мебели — РОСТ Мебель',
            description: 'Портфолио реализованных кухонь, шкафов и систем хранения: фотографии, бюджеты, материалы и детали проектов.',
            canonicalPath: '/catalog',
          },
        },
        {
          path: 'product/:id',
          name: 'product',
          component: () => import('@/pages/ProductPage.vue'),
          meta: {
            title: 'Детали проекта — РОСТ Мебель',
            description: 'Карточка реализованного проекта мебели на заказ: фотографии, бюджет, материалы, детали и заявка на расчет.',
          },
        },
        {
          path: 'contact',
          name: 'contact',
          component: () => import('@/pages/ContactPage.vue'),
          meta: {
            title: 'Контакты — РОСТ Мебель',
            description: 'Контакты РОСТ Мебель: консультация, предварительный расчет, замер, доставка и монтаж мебели по Крыму.',
            canonicalPath: '/contact',
          },
        },
        {
          path: 'favorites',
          name: 'favorites',
          component: () => import('@/pages/FavoritesPage.vue'),
          meta: {
            title: 'Избранное — РОСТ Мебель',
            description: 'Сохраненные проекты мебели, к которым можно вернуться перед консультацией.',
            canonicalPath: '/favorites',
            robots: 'noindex,nofollow',
          },
        },
      ],
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('@/pages/admin/LoginPage.vue'),
      meta: { robots: 'noindex,nofollow' },
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true, robots: 'noindex,nofollow' },
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
  setPageSeo({
    title: to.meta.title as string | undefined,
    description: to.meta.description as string | undefined,
    path: (to.meta.canonicalPath as string | undefined) || to.path,
    robots: (to.meta.robots as string | undefined) || 'index,follow',
  });
});

export default router;
